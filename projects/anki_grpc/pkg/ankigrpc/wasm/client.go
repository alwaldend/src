package wasm

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/internal/wasmhost"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"google.golang.org/grpc"
)

// Client owns one embedded Rust server and an in-memory gRPC connection.
// Use New; do not copy a Client. Each instance permits one Start attempt.
type Client struct {
	// mutex protects lifecycle state and the terminal execution result.
	mutex sync.Mutex
	// transport owns gRPC dialing, readiness and error decoding over memory stdio.
	transport *ankiclient.ClientGRPC
	// cancel interrupts the independent guest lifetime when Stop begins.
	cancel context.CancelFunc
	// started records the single permitted startup attempt.
	started bool
	// stopped rejects new calls once terminal cleanup begins.
	stopped bool
	// done closes after guest execution and all owned resources are released.
	done chan struct{}
	// exitError preserves unexpected guest and cleanup failures.
	exitError error
}

// New constructs an inert client with its own backend state and no process or listener.
func New() *Client { return &Client{done: make(chan struct{})} }

// Start executes the embedded server and waits for readiness; ctx owns startup only.
// Call Stop after any startup attempt; Done observes cleanup beyond a Stop deadline.
func (client *Client) Start(ctx context.Context) error {
	client.mutex.Lock()
	if client.started || client.stopped {
		client.mutex.Unlock()
		return fmt.Errorf("start in-process WASM client: already started or stopped")
	}
	client.started = true
	if err := ctx.Err(); err != nil {
		close(client.done)
		client.mutex.Unlock()
		return fmt.Errorf("start in-process WASM client: %w", err)
	}
	caller, server := net.Pipe()
	var dialed atomic.Bool
	client.transport = ankiclient.NewGRPCWithDialer(func(context.Context) (net.Conn, error) {
		if !dialed.CompareAndSwap(false, true) {
			return nil, fmt.Errorf("in-process WASM connection cannot be reopened")
		}
		return caller, nil
	})
	lifetime, cancel := context.WithCancel(context.Background())
	client.cancel = cancel
	go client.execute(lifetime, caller, server)
	client.mutex.Unlock()

	if err := client.transport.Start(ctx); err != nil {
		return fmt.Errorf("start in-process WASM client: %w", errors.Join(err, client.Stop(ctx)))
	}
	client.mutex.Lock()
	failure := context.Cause(ctx)
	select {
	case <-client.done:
		failure = errors.Join(failure, fmt.Errorf("WASM server exited during startup"), client.exitError)
	default:
	}
	if client.stopped {
		failure = errors.Join(failure, fmt.Errorf("WASM client stopped during startup"))
	}
	client.mutex.Unlock()
	if failure != nil {
		return fmt.Errorf("start in-process WASM client: %w", errors.Join(failure, client.Stop(ctx)))
	}
	return nil
}

// execute owns guest execution and closes the memory stream and gRPC connection once.
func (client *Client) execute(ctx context.Context, caller, server net.Conn) {
	err := wasmhost.RunCommand(ctx, []string{"anki_grpc", "serve"}, guest, server, server, os.Stderr)
	cleanup := errors.Join(client.transport.Stop(context.Background()), caller.Close(), server.Close())
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if client.stopped && errors.Is(err, context.Canceled) {
		err = nil
	}
	if err != nil {
		err = fmt.Errorf("execute in-process WASM server: %w", err)
	}
	if cleanup != nil {
		cleanup = fmt.Errorf("close in-process WASM transport: %w", cleanup)
	}
	client.exitError = errors.Join(err, cleanup)
	client.cancel()
	close(client.done)
}

// Connection lends the gRPC client only after successful startup and before shutdown.
func (client *Client) Connection() (grpc.ClientConnInterface, error) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if client.transport == nil || client.stopped {
		return nil, fmt.Errorf("get in-process WASM connection: client is not ready")
	}
	connection, err := client.transport.Connection()
	if err != nil {
		return nil, fmt.Errorf("get in-process WASM connection: %w", err)
	}
	return connection, nil
}

// Stop rejects calls and cancels execution, using ctx only to bound waiting for cleanup.
// Close an open collection through Anki first. Repeated Stop preserves terminal errors.
func (client *Client) Stop(ctx context.Context) error {
	client.mutex.Lock()
	if !client.stopped {
		client.stopped = true
		if client.cancel != nil {
			client.cancel()
		} else if !client.started {
			close(client.done)
		}
	}
	transport := client.transport
	client.mutex.Unlock()
	var transportError error
	if transport != nil {
		transportError = transport.Stop(ctx)
	}
	select {
	case <-client.done:
	default:
		select {
		case <-client.done:
		case <-ctx.Done():
			return fmt.Errorf("wait for in-process WASM shutdown: %w", errors.Join(ctx.Err(), transportError))
		}
	}
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if client.exitError != nil {
		return fmt.Errorf("stop in-process WASM client: %w", client.exitError)
	}
	return nil
}

// Done closes after terminal cleanup, including when Stop's caller stopped waiting.
func (client *Client) Done() <-chan struct{} { return client.done }
