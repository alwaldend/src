package ankiclient

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"path/filepath"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// ClientGRPCOptions selects a running server without acquiring its process or collection lifetime.
type ClientGRPCOptions struct {
	// Address is a TCP host:port, tcp:// endpoint, unix:// endpoint or absolute Unix socket path.
	Address string
}

// ClientGRPC owns only a gRPC client connection to an explicitly selected Anki server.
// Each instance allows one Start attempt; Stop leaves remote storage and processes running.
type ClientGRPC struct {
	// options are copied constructor settings.
	options ClientGRPCOptions
	// mutex protects connection readiness and terminal shutdown.
	mutex sync.Mutex
	// connection is the owned network client.
	connection *grpc.ClientConn
	// started prevents repeated startup attempts.
	started bool
	// ready permits calls after the initial connection handshake.
	ready bool
	// stopped makes shutdown terminal, including before Start.
	stopped bool
}

// NewGRPC configures a network backend without contacting the server.
func NewGRPC(options ClientGRPCOptions) *ClientGRPC {
	return &ClientGRPC{options: options}
}

// Start creates the client and waits for readiness within ctx; ctx owns startup only.
func (backend *ClientGRPC) Start(ctx context.Context) error {
	backend.mutex.Lock()
	if backend.started || backend.stopped {
		backend.mutex.Unlock()
		return fmt.Errorf("start Anki gRPC backend: already started or stopped")
	}
	backend.started = true
	address := backend.options.Address
	transport := "tcp"
	if filepath.IsAbs(address) {
		address = "unix://" + address
	}
	if strings.HasPrefix(address, "unix://") {
		transport = "unix"
		address = strings.TrimPrefix(address, "unix://")
		if address == "" {
			backend.mutex.Unlock()
			return fmt.Errorf("parse Anki backend address: Unix socket path is empty")
		}
	} else {
		address = strings.TrimPrefix(address, "tcp://")
		if _, _, err := net.SplitHostPort(address); err != nil {
			backend.mutex.Unlock()
			return fmt.Errorf("parse Anki backend address: %w", err)
		}
	}
	connection, err := grpc.NewClient("passthrough:///anki-network",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			connection, err := (&net.Dialer{}).DialContext(ctx, transport, address)
			if err != nil {
				return nil, fmt.Errorf("dial Anki %s backend %q: %w", transport, address, err)
			}
			return connection, nil
		}),
		grpc.WithDisableRetry(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)),
		grpc.WithUnaryInterceptor(backend.intercept),
	)
	backend.connection = connection
	backend.mutex.Unlock()
	if err != nil {
		return fmt.Errorf("create Anki network client: %w", err)
	}
	connection.Connect()
	for state := connection.GetState(); state != connectivity.Ready; state = connection.GetState() {
		if state == connectivity.Shutdown || !connection.WaitForStateChange(ctx, state) {
			cause := context.Cause(ctx)
			if cause == nil {
				cause = fmt.Errorf("client closed before readiness")
			}
			return fmt.Errorf("wait for Anki network readiness: %w", errors.Join(cause, backend.Stop(context.WithoutCancel(ctx))))
		}
	}
	backend.mutex.Lock()
	defer backend.mutex.Unlock()
	if backend.stopped {
		return fmt.Errorf("start Anki network backend: stopped during startup")
	}
	backend.ready = true
	return nil
}

// Connection lends the client after successful startup and before shutdown.
func (backend *ClientGRPC) Connection() (grpc.ClientConnInterface, error) {
	backend.mutex.Lock()
	defer backend.mutex.Unlock()
	if !backend.ready || backend.stopped {
		return nil, fmt.Errorf("get Anki network connection: backend is not ready")
	}
	return backend.connection, nil
}

// intercept preserves native error details and rejects retained clients after Stop.
func (backend *ClientGRPC) intercept(ctx context.Context, method string, input, output any, connection *grpc.ClientConn, invoke grpc.UnaryInvoker, options ...grpc.CallOption) error {
	if _, err := backend.Connection(); err != nil {
		return fmt.Errorf("call Anki %s: %w", method, err)
	}
	var trailers metadata.MD
	err := invoke(ctx, method, input, output, connection, append(options, grpc.Trailer(&trailers))...)
	return decodeCallError(ctx, method, err, trailers)
}

// Stop closes the client connection without invoking any server RPC.
func (backend *ClientGRPC) Stop(context.Context) error {
	backend.mutex.Lock()
	defer backend.mutex.Unlock()
	if backend.stopped {
		return nil
	}
	backend.stopped, backend.ready = true, false
	if backend.connection != nil {
		if err := backend.connection.Close(); err != nil {
			return fmt.Errorf("close Anki network client: %w", err)
		}
	}
	return nil
}
