package ankigrpc

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"google.golang.org/grpc"
)

// ClientFactory constructs unstarted clients from connection strings.
// The caller owns startup and cleanup of each returned client.
type ClientFactory func(address string) (ankiclient.StartableClient, error)

// ToClientFactory converts a concrete-returning constructor into a ClientFactory.
// It does not start the client.
// A nil constructor produces a nil factory.
func ToClientFactory[T ankiclient.StartableClient](constructor func(string) (T, error)) ClientFactory {
	if constructor == nil {
		return nil
	}
	return func(address string) (ankiclient.StartableClient, error) {
		client, err := constructor(address)
		if err != nil {
			return nil, fmt.Errorf("construct Anki client with factory function: %w", err)
		}
		return client, nil
	}
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
		factory = defaultFactory.NewClient
	}
	client, err := factory(address)
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

// NewSubprocessClient constructs an unstarted worker from an explicit subprocess URL.
// The executable owner registers this constructor for its supported worker schemes.
// The URL must identify an executable; arguments are passed exactly, including zero values.
func NewSubprocessClient(address string) (*ankiclient.ClientSTD, error) {
	options, err := ParseSubprocessAddress(address)
	if err != nil {
		return nil, fmt.Errorf("parse subprocess client options: %w", err)
	}
	return ankiclient.NewSTD(options), nil
}

// ParseSubprocessAddress validates and decodes an explicit worker URL without starting it.
// Executable owners may adapt the returned options before constructing their client.
func ParseSubprocessAddress(address string) (ankiclient.ClientSTDOptions, error) {
	parsed, err := url.Parse(address)
	if err != nil {
		return ankiclient.ClientSTDOptions{}, fmt.Errorf("parse subprocess address: %w", err)
	}
	if parsed.User != nil || parsed.Fragment != "" || parsed.Opaque != "" || parsed.Host != "" && parsed.Path != "" {
		return ankiclient.ClientSTDOptions{}, fmt.Errorf("subprocess address requires one executable path without credentials or a fragment")
	}
	query, err := url.ParseQuery(parsed.RawQuery)
	if err != nil {
		return ankiclient.ClientSTDOptions{}, fmt.Errorf("parse subprocess parameters: %w", err)
	}
	for key := range query {
		if key != "binary" && key != "arg" && key != "shutdown_timeout" {
			return ankiclient.ClientSTDOptions{}, fmt.Errorf("unknown subprocess parameter %q", key)
		}
	}
	binary := parsed.Path
	if parsed.Host != "" {
		binary = parsed.Host
	}
	if values, exists := query["binary"]; exists {
		if len(values) != 1 || binary != "" || values[0] == "" {
			return ankiclient.ClientSTDOptions{}, fmt.Errorf("subprocess binary must specify one nonempty executable")
		}
		binary = values[0]
	}
	if binary == "" {
		return ankiclient.ClientSTDOptions{}, fmt.Errorf("subprocess address requires an explicit executable")
	}
	arguments := query["arg"]
	var shutdown time.Duration
	if values, exists := query["shutdown_timeout"]; exists {
		if len(values) != 1 {
			return ankiclient.ClientSTDOptions{}, fmt.Errorf("subprocess shutdown_timeout must appear once")
		}
		shutdown, err = time.ParseDuration(values[0])
		if err != nil {
			return ankiclient.ClientSTDOptions{}, fmt.Errorf("parse subprocess shutdown_timeout: %w", err)
		}
		if shutdown < 0 {
			return ankiclient.ClientSTDOptions{}, fmt.Errorf("subprocess shutdown_timeout must be nonnegative")
		}
	}
	return ankiclient.ClientSTDOptions{Binary: binary, Arguments: arguments, ShutdownTimeout: shutdown}, nil
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
