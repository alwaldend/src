package ankiclient

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// ClientGRPC owns only a gRPC client connection to an explicitly selected Anki server.
// Each instance allows one Start attempt; Stop leaves remote storage and processes running.
type ClientGRPC struct {
	// transport is the validated TCP or Unix network used for dialing.
	transport string
	// address is a decoded endpoint; Start never parses it again.
	address string
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

// NewGRPC validates an endpoint and configures a client without contacting its server.
func NewGRPC(address string) (*ClientGRPC, error) {
	transport, endpoint, err := parseNetworkAddress(address)
	if err != nil {
		return nil, fmt.Errorf("configure Anki network client: %w", err)
	}
	return &ClientGRPC{transport: transport, address: endpoint}, nil
}

// parseNetworkAddress decodes endpoint syntax once before any client startup.
// Unix URL suffixes are paths, including relative paths containing spaces.
func parseNetworkAddress(address string) (string, string, error) {
	if filepath.IsAbs(address) {
		return "unix", address, nil
	}
	if len(address) >= len("unix://") && strings.EqualFold(address[:len("unix://")], "unix://") {
		path := address[len("unix://"):]
		if path == "" || strings.ContainsAny(path, "?#") || strings.Contains(strings.SplitN(path, "/", 2)[0], "@") {
			return "", "", fmt.Errorf("Unix address requires a socket path without credentials, parameters or a fragment")
		}
		decoded, err := url.PathUnescape(path)
		if err != nil {
			return "", "", fmt.Errorf("decode Unix socket path: %w", err)
		}
		return "unix", decoded, nil
	}
	explicit := len(address) >= len("tcp://") && strings.EqualFold(address[:len("tcp://")], "tcp://")
	endpoint := address
	if explicit {
		endpoint = address[len("tcp://"):]
	}
	host, port, err := net.SplitHostPort(endpoint)
	if err != nil {
		return "", "", fmt.Errorf("parse TCP endpoint: %w", err)
	}
	if strings.ContainsAny(port, "/?#") || !strings.HasPrefix(endpoint, "[") && strings.ContainsAny(host, "/?#@") {
		return "", "", fmt.Errorf("TCP address requires an endpoint without credentials, a path, parameters or a fragment")
	}
	if explicit {
		if !strings.HasPrefix(endpoint, "[") || !strings.Contains(host, "%") {
			if _, err := url.Parse(address); err != nil {
				return "", "", fmt.Errorf("validate TCP URL: %w", err)
			}
		}
		if zone := strings.IndexByte(host, '%'); zone >= 0 && strings.HasPrefix(host[zone:], "%25") {
			// Accept legacy raw zones and decode escaped URL zones exactly once.
			// Mask escapes for URI syntax validation without net/url rejecting
			// reserved characters in a bracketed Linux interface zone.
			if _, err := url.Parse(strings.ReplaceAll(address, "%", "%25")); err == nil {
				if decoded, err := url.PathUnescape(host); err == nil {
					host = decoded
				}
			}
		}
	}
	return "tcp", net.JoinHostPort(host, port), nil
}

// Start creates the client and waits for readiness within ctx; ctx owns startup only.
func (backend *ClientGRPC) Start(ctx context.Context) error {
	backend.mutex.Lock()
	if backend.started || backend.stopped {
		backend.mutex.Unlock()
		return fmt.Errorf("start Anki gRPC backend: already started or stopped")
	}
	backend.started = true
	transport, address := backend.transport, backend.address
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
	return DecodeCallError(ctx, method, err, trailers)
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
