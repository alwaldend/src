package ankiclient

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	ankibackend "github.com/ankitects/anki/proto/anki/backend"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// ankiErrorTrailer carries Anki's serialized BackendError on failed native calls.
const ankiErrorTrailer = "anki-error-bin"

// Client manages Anki resource readiness and its connection for generated service clients.
type Client interface {
	// Start prepares owned resources; ctx bounds startup only.
	Start(ctx context.Context) error
	// Connection returns the managed connection, or an error when unavailable.
	Connection() (grpc.ClientConnInterface, error)
	// Stop uses ctx to bound graceful cleanup; ClientSTD always completes process reaping.
	Stop(ctx context.Context) error
}

// ClientSTDOptions configures the process and shutdown policy used by NewSTD.
type ClientSTDOptions struct {
	// Binary is the server executable path or a name resolved through PATH.
	Binary string
	// Arguments are the server's CLI arguments; NewSTD copies them for the backend.
	Arguments []string
	// ShutdownTimeout bounds graceful shutdown before forced termination.
	// Zero uses five seconds; Start rejects negative values.
	// A shorter Stop context takes precedence.
	ShutdownTimeout time.Duration
}

// ClientSTD owns one Anki server process communicating over stdin/stdout.
// Use NewSTD to configure it; ClientSTD must not be copied after use.
type ClientSTD struct {
	// options are owned process settings with a resolved default shutdown timeout.
	options ClientSTDOptions
	// lifecycle serializes process configuration with shutdown and connection access.
	lifecycle sync.Mutex
	// started records the single startup attempt permitted for this instance.
	started bool
	// ready records a completed handshake before exposing the connection.
	ready bool
	// command is the child process owned and reaped by this backend.
	command *exec.Cmd
	// cancelProcess ends the command's independent lifetime context to force termination.
	cancelProcess context.CancelFunc
	// connection carries gRPC HTTP/2 over the child's pipes.
	connection *grpc.ClientConn
	// pipes are the Go-owned transport endpoints.
	pipes *pipeConn
	// done closes after Wait stores the child's exit result.
	done chan struct{}
	// exitError is read only after done closes.
	exitError error
	// stopped rejects calls once shutdown begins, including before Start.
	stopped atomic.Bool
	// stopOnce serializes shutdown and preserves its result.
	stopOnce sync.Once
	// stopError is the result shared by repeated Stop calls.
	stopError error
}

// NewSTD configures a stdin/stdout backend without starting its server.
// It copies options.Arguments; the child inherits the parent environment and stderr.
func NewSTD(options ClientSTDOptions) *ClientSTD {
	options.Arguments = append([]string(nil), options.Arguments...)
	if options.ShutdownTimeout == 0 {
		options.ShutdownTimeout = 5 * time.Second
	}
	return &ClientSTD{options: options}
}

// Start launches the configured server and waits for gRPC readiness.
// ctx controls startup only. Each instance accepts one attempt;
// failed startup releases its resources and reaps any child it started.
func (backend *ClientSTD) Start(ctx context.Context) error {
	backend.lifecycle.Lock()
	if backend.started || backend.stopped.Load() {
		backend.lifecycle.Unlock()
		return fmt.Errorf("start Anki worker: backend has already started or stopped")
	}
	backend.started = true
	err := backend.startProcess(ctx)
	if err == nil {
		err = backend.connect()
	}
	backend.lifecycle.Unlock()
	if err == nil {
		err = backend.waitForReady(ctx)
	}
	if err != nil {
		return fmt.Errorf("start Anki worker: %w", errors.Join(err, backend.Stop(ctx)))
	}
	backend.lifecycle.Lock()
	defer backend.lifecycle.Unlock()
	if backend.stopped.Load() {
		return fmt.Errorf("start Anki worker: backend stopped during startup")
	}
	backend.ready = true
	return nil
}

// startProcess creates the pipes and launches the child while lifecycle is locked.
func (backend *ClientSTD) startProcess(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("check Anki startup context: %w", err)
	}
	if backend.options.ShutdownTimeout < 0 {
		return fmt.Errorf("configure Anki worker: shutdown timeout must not be negative")
	}
	pipes, input, output, err := newPipeConn()
	if err != nil {
		return fmt.Errorf("prepare Anki transport: %w", err)
	}
	processContext, cancel := context.WithCancel(context.Background())
	command := exec.CommandContext(processContext, backend.options.Binary, backend.options.Arguments...)
	command.Stdin, command.Stdout = input, output
	command.Stderr = os.Stderr
	startError := command.Start()
	childPipeError := errors.Join(input.Close(), output.Close())
	if startError != nil {
		cancel()
		return fmt.Errorf("execute Anki worker: %w", errors.Join(startError, childPipeError, pipes.Close()))
	}
	backend.command, backend.pipes, backend.done = command, pipes, make(chan struct{})
	backend.cancelProcess = cancel
	go func() {
		backend.exitError = command.Wait()
		cancel()
		close(backend.done)
	}()
	if childPipeError != nil {
		return fmt.Errorf("close Anki child pipe endpoints: %w", childPipeError)
	}
	return nil
}

// connect creates a single-use gRPC connection while lifecycle is locked.
func (backend *ClientSTD) connect() error {
	var dialed atomic.Bool
	connection, err := grpc.NewClient("passthrough:///anki-worker",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDisableRetry(),
		grpc.WithUnaryInterceptor(backend.intercept),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			if !dialed.CompareAndSwap(false, true) {
				return nil, fmt.Errorf("Anki worker connection cannot be reopened")
			}
			return backend.pipes, nil
		}),
	)
	backend.connection = connection
	if err != nil {
		return fmt.Errorf("create Anki gRPC connection: %w", err)
	}
	return nil
}

// waitForReady bounds the HTTP/2 handshake by cancellation and worker exit.
func (backend *ClientSTD) waitForReady(ctx context.Context) error {
	startup, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	go func() {
		select {
		case <-backend.done:
			cancel(errors.Join(fmt.Errorf("Anki worker exited before gRPC readiness"), backend.exitError))
		case <-startup.Done():
		}
	}()
	backend.connection.Connect()
	for state := backend.connection.GetState(); state != connectivity.Ready; state = backend.connection.GetState() {
		if !backend.connection.WaitForStateChange(startup, state) {
			return fmt.Errorf("wait for Anki gRPC readiness: %w", context.Cause(startup))
		}
	}
	if err := context.Cause(startup); err != nil {
		return fmt.Errorf("wait for Anki gRPC readiness: %w", err)
	}
	return nil
}

// Connection supplies the managed connection to generated Anki service clients.
// It returns an error until startup succeeds, while stopping, or after worker exit.
func (backend *ClientSTD) Connection() (grpc.ClientConnInterface, error) {
	backend.lifecycle.Lock()
	defer backend.lifecycle.Unlock()
	if backend.stopped.Load() {
		return nil, fmt.Errorf("get Anki connection: backend is stopped")
	}
	if !backend.ready {
		return nil, fmt.Errorf("get Anki connection: backend is not ready")
	}
	select {
	case <-backend.done:
		if backend.exitError != nil {
			return nil, fmt.Errorf("get Anki connection: worker exited: %w", backend.exitError)
		}
		return nil, fmt.Errorf("get Anki connection: worker exited")
	default:
		return backend.connection, nil
	}
}

// intercept adds lifecycle checks and native error details to generated calls.
// Cancellation stops waiting; it cannot roll back a mutation already in Anki.
// Failed calls are never replayed or sent to a replacement worker.
func (backend *ClientSTD) intercept(ctx context.Context, method string, input, output any, connection *grpc.ClientConn, invoke grpc.UnaryInvoker, options ...grpc.CallOption) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("call Anki %s: %w", method, err)
	}
	if backend.stopped.Load() {
		return fmt.Errorf("call Anki %s: backend is stopped", method)
	}
	var trailers metadata.MD
	err := invoke(ctx, method, input, output, connection, append(options, grpc.Trailer(&trailers))...)
	return decodeCallError(ctx, method, err, trailers)
}

// decodeCallError preserves gRPC status and native BackendError details for every transport.
func decodeCallError(ctx context.Context, method string, err error, trailers metadata.MD) error {
	if err == nil {
		return nil
	}
	if ctx.Err() != nil {
		return fmt.Errorf("call Anki %s: %w", method, ctx.Err())
	}
	data := trailers.Get(ankiErrorTrailer)
	if len(data) != 1 {
		return fmt.Errorf("call Anki %s: %w", method, err)
	}
	var failure ankibackend.BackendError
	if decodeError := proto.Unmarshal([]byte(data[0]), &failure); decodeError != nil {
		return fmt.Errorf("decode Anki %s failure: %w", method, errors.Join(err, decodeError))
	}
	return fmt.Errorf("call Anki %s: %s: %s: %w", method, failure.Kind, failure.Message, err)
}

// Stop rejects calls, closes transport, and reaps the child. When ctx expires or
// after options.ShutdownTimeout, it cancels the command's context for forced termination.
// It reaps the child even after cancellation.
// Close an open collection through Anki first.
// Stop is safe before Start; repeated calls preserve unexpected exit errors.
func (backend *ClientSTD) Stop(ctx context.Context) error {
	backend.stopOnce.Do(func() {
		backend.lifecycle.Lock()
		backend.stopped.Store(true)
		backend.lifecycle.Unlock()
		if backend.command == nil {
			return
		}
		shutdown, finish := context.WithTimeout(ctx, backend.options.ShutdownTimeout)
		defer finish()
		stopCancellation := context.AfterFunc(shutdown, backend.cancelProcess)
		defer stopCancellation()
		inputError := backend.pipes.closeInput()
		<-backend.done
		var connectionError error
		if backend.connection != nil {
			connectionError = backend.connection.Close()
		}
		pipeError := backend.pipes.Close()
		if err := errors.Join(inputError, connectionError, pipeError, ctx.Err(), backend.exitError); err != nil {
			backend.stopError = fmt.Errorf("stop Anki worker: %w", err)
		}
	})
	return backend.stopError
}
