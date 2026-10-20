package ankigrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"google.golang.org/grpc"
)

// ClientFactory constructs unstarted managed clients from connection strings.
type ClientFactory interface {
	// NewClient creates a client whose startup and cleanup belong to the caller.
	NewClient(address string) (ankiclient.Client, error)
}

// ClientFactoryFunc adapts a constructor with concrete result T to ClientFactory.
type ClientFactoryFunc[T ankiclient.Client] func(string) (T, error)

// NewClient invokes the adapted constructor without starting its client.
func (factory ClientFactoryFunc[T]) NewClient(address string) (ankiclient.Client, error) {
	if factory == nil {
		return nil, fmt.Errorf("Anki client factory requires a constructor function")
	}
	client, err := factory(address)
	if err != nil {
		return nil, fmt.Errorf("construct Anki client with factory function: %w", err)
	}
	return client, nil
}

// isNil reports an absent constructor across all generic adapter instantiations.
func (factory ClientFactoryFunc[T]) isNil() bool {
	return factory == nil
}

// WithClient constructs a client through factory, starts it and lends its ready connection.
// A nil factory selects NewClient. Cleanup follows every startup attempt or callback
// exit, including panic unwinding, with a five-second budget independent of ctx.
// It closes only resources owned by that client; it neither opens nor reserves a collection.
func WithClient(ctx context.Context, factory ClientFactory, address string, operation func(grpc.ClientConnInterface) error) (err error) {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start Anki client scope: %w", errors.Join(err, context.Cause(ctx)))
	}
	if operation == nil {
		return fmt.Errorf("Anki client scope requires an operation")
	}
	if factory == nil {
		factory = ClientFactoryFunc[ankiclient.Client](NewClient)
	}
	client, err := factory.NewClient(address)
	if err != nil {
		return fmt.Errorf("construct scoped Anki client: %w", err)
	}
	defer func() {
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		if stopErr := client.Stop(cleanup); stopErr != nil {
			err = errors.Join(err, fmt.Errorf("stop scoped Anki client: %w", stopErr))
		}
	}()
	if err := client.Start(ctx); err != nil {
		return fmt.Errorf("start scoped Anki client: %w", err)
	}
	connection, err := client.Connection()
	if err != nil {
		return fmt.Errorf("get scoped Anki connection: %w", err)
	}
	if err := operation(connection); err != nil {
		return fmt.Errorf("run scoped Anki operation: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("finish Anki client scope: %w", errors.Join(err, context.Cause(ctx)))
	}
	return nil
}

// clientFactoriesMutex protects registration and lookup without enclosing factory calls.
var clientFactoriesMutex sync.RWMutex

// clientFactories maps URL schemes to constructors; it never retains client instances.
var clientFactories = map[string]ClientFactory{
	"tcp":  ClientFactoryFunc[*ankiclient.ClientGRPC](newNetworkClient),
	"unix": ClientFactoryFunc[*ankiclient.ClientGRPC](newNetworkClient),
}

// RegisterClientFactory installs a constructor for a lowercase URL scheme.
// Executables and optional direct adapters register their owned schemes. It panics for nil
// factories, invalid schemes or duplicate registrations, preserving existing owners.
func RegisterClientFactory(scheme string, factory ClientFactory) {
	parsed, err := url.Parse(scheme + ":")
	function, isFunction := factory.(interface{ isNil() bool })
	if factory == nil || isFunction && function.isNil() || err != nil || scheme == "" || parsed.Scheme != scheme {
		panic(fmt.Sprintf("register Anki client factory: invalid scheme or factory %q", scheme))
	}
	clientFactoriesMutex.Lock()
	defer clientFactoriesMutex.Unlock()
	if _, exists := clientFactories[scheme]; exists {
		panic(fmt.Sprintf("register Anki client factory: scheme %q is already registered", scheme))
	}
	clientFactories[scheme] = factory
}

// NewClient constructs an unstarted lifecycle owner using its registered URL scheme.
// TCP and Unix are built in. Executables register their worker schemes explicitly;
// optional in-process adapters register through opt-in imports. An empty address
// is invalid because this library owns no default worker entrypoint.
// The caller owns the returned client's startup and cleanup.
func NewClient(address string) (ankiclient.Client, error) {
	parsed, err := parseClientAddress(address)
	if err != nil {
		return nil, fmt.Errorf("select Anki client: %w", err)
	}
	clientFactoriesMutex.RLock()
	factory := clientFactories[parsed.Scheme]
	clientFactoriesMutex.RUnlock()
	if factory == nil {
		return nil, fmt.Errorf("Anki client scheme %q is unavailable; its owner must register a client factory", parsed.Scheme)
	}
	client, err := factory.NewClient(parsed.String())
	if err != nil {
		return nil, fmt.Errorf("construct Anki client for scheme %q: %w", parsed.Scheme, err)
	}
	return client, nil
}

// parseClientAddress normalizes shorthand independently of registered constructors.
func parseClientAddress(address string) (*url.URL, error) {
	if address == "" {
		return nil, fmt.Errorf("Anki client address requires an explicit endpoint or registered scheme")
	}
	if filepath.IsAbs(address) {
		return &url.URL{Scheme: "unix", Path: address}, nil
	}
	// Brackets keep zone characters inside the endpoint; only the port can
	// contain URL suffixes. Unbracketed engine-qualified URLs still use parsing.
	if host, port, err := net.SplitHostPort(address); err == nil && !strings.ContainsAny(port, "/?#") && (strings.HasPrefix(address, "[") || !strings.ContainsAny(host, "/?#+")) {
		return &url.URL{Scheme: "tcp", Host: address}, nil
	}
	parsed, err := parseClientURL(address)
	if err != nil {
		return nil, fmt.Errorf("parse Anki client address: %w", err)
	}
	if parsed.Scheme == "" {
		return &url.URL{Scheme: "tcp", Host: address}, nil
	}
	return parsed, nil
}

// parseClientURL decodes URLs and scoped TCP zones without inferring address shorthand.
func parseClientURL(address string) (*url.URL, error) {
	if strings.HasPrefix(address, "tcp://[") {
		endpoint := strings.TrimPrefix(address, "tcp://")
		if host, port, err := net.SplitHostPort(endpoint); err == nil && strings.Contains(host, "%") && !strings.ContainsAny(port, "/?#") {
			zone := strings.IndexByte(host, '%')
			if strings.HasPrefix(host[zone:], "%25") {
				// Mask escapes to validate URI syntax without net/url rejecting
				// encoded zone characters. Raw delimiters retain their legacy spelling.
				if _, err := url.Parse(strings.ReplaceAll(address, "%", "%25")); err == nil {
					if decoded, err := url.PathUnescape(host); err == nil {
						host = decoded
					}
				}
			}
			return &url.URL{Scheme: "tcp", Host: net.JoinHostPort(host, port)}, nil
		}
	}
	parsed, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("parse Anki client URL: %w", err)
	}
	return parsed, nil
}

// newNetworkClient constructs a borrowed-server connection without contacting it.
func newNetworkClient(address string) (*ankiclient.ClientGRPC, error) {
	parsed, err := parseClientURL(address)
	if err != nil {
		return nil, fmt.Errorf("parse network address: %w", err)
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.Opaque != "" || parsed.Scheme == "tcp" && parsed.Path != "" {
		return nil, fmt.Errorf("network address requires a TCP endpoint or Unix socket path without credentials, parameters or a fragment")
	}
	target := parsed.Host
	if parsed.Scheme == "unix" {
		target = "unix://" + parsed.Host + parsed.Path
	}
	return ankiclient.NewGRPC(ankiclient.ClientGRPCOptions{Address: target}), nil
}

// NewSubprocessClient constructs an unstarted worker from a normalized subprocess URL.
// The executable owner registers this constructor for its supported worker schemes.
// The URL must identify an executable; arguments are passed exactly, including zero values.
func NewSubprocessClient(address string) (*ankiclient.ClientSTD, error) {
	parsed, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("parse subprocess address: %w", err)
	}
	if parsed.User != nil || parsed.Fragment != "" || parsed.Opaque != "" || parsed.Host != "" && parsed.Path != "" {
		return nil, fmt.Errorf("subprocess address requires one executable path without credentials or a fragment")
	}
	query, err := url.ParseQuery(parsed.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("parse subprocess parameters: %w", err)
	}
	for key := range query {
		if key != "binary" && key != "arg" && key != "shutdown_timeout" {
			return nil, fmt.Errorf("unknown subprocess parameter %q", key)
		}
	}
	binary := parsed.Path
	if parsed.Host != "" {
		binary = parsed.Host
	}
	if values, exists := query["binary"]; exists {
		if len(values) != 1 || binary != "" || values[0] == "" {
			return nil, fmt.Errorf("subprocess binary must specify one nonempty executable")
		}
		binary = values[0]
	}
	if binary == "" {
		return nil, fmt.Errorf("subprocess address requires an explicit executable")
	}
	arguments := query["arg"]
	var shutdown time.Duration
	if values, exists := query["shutdown_timeout"]; exists {
		if len(values) != 1 {
			return nil, fmt.Errorf("subprocess shutdown_timeout must appear once")
		}
		shutdown, err = time.ParseDuration(values[0])
		if err != nil {
			return nil, fmt.Errorf("parse subprocess shutdown_timeout: %w", err)
		}
		if shutdown < 0 {
			return nil, fmt.Errorf("subprocess shutdown_timeout must be nonnegative")
		}
	}
	return ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: binary, Arguments: arguments, ShutdownTimeout: shutdown}), nil
}

// SubprocessAddress encodes an executable and argument vector without a shell.
// The executable and arguments are copied into query parameters, preserving
// relative paths and reserved characters. Omitting arg values means zero arguments.
func SubprocessAddress(binary string, arguments ...string) string {
	query := make(url.Values)
	if binary != "" {
		query.Set("binary", binary)
	}
	for _, argument := range arguments {
		query.Add("arg", argument)
	}
	if len(query) == 0 {
		return "subprocess://"
	}
	return "subprocess://?" + query.Encode()
}
