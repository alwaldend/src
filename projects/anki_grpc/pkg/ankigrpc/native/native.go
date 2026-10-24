//go:build cgo && linux && amd64

package native

import (
	"context"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/internal/directclient"
	"google.golang.org/grpc"
)

// Client owns one in-process connection and its native resource lifetime.
// It constructs the connection for generated Anki clients without starting a process or listener.
// Use New or the zero value; do not copy a client after use.
type Client struct {
	// connection is constructed with the client and retained until native cleanup finishes.
	connection directclient.Client
}

// New constructs an unstarted client using upstream Anki initialization defaults.
func New() *Client { return &Client{} }

// Start initializes the owned connection synchronously, checking ctx before and after.
// Cancellation waits for native initialization and failure cleanup to finish.
func (client *Client) Start(ctx context.Context) error {
	if err := client.connection.Start(ctx, func() (directclient.Backend, error) {
		server, err := newServer()
		if err != nil {
			return nil, fmt.Errorf("initialize native Anki server: %w", err)
		}
		return &nativeBackend{server: server}, nil
	}); err != nil {
		return fmt.Errorf("start in-process Anki client: %w", err)
	}
	return nil
}

// Connection lends the client's connection after startup and before shutdown.
func (client *Client) Connection() (grpc.ClientConnInterface, error) {
	connection, err := client.connection.Connection()
	if err != nil {
		return nil, fmt.Errorf("get in-process Anki connection: %w", err)
	}
	return connection, nil
}

// Stop rejects new calls and releases native state after active calls finish.
// A deadline bounds waiting while cleanup continues; a later Stop can await it.
// Close an open collection through its generated client before stopping the client.
func (client *Client) Stop(ctx context.Context) error {
	if err := client.connection.Stop(ctx); err != nil {
		return fmt.Errorf("stop in-process Anki client: %w", err)
	}
	return nil
}

// Done closes after Stop finishes startup and native cleanup, including cancelled calls.
// It returns the same channel before startup, while stopping and after cleanup.
func (client *Client) Done() <-chan struct{} {
	return client.connection.Done()
}

// nativeBackend adapts the private native handle to the shared command owner.
type nativeBackend struct {
	// server owns the C ABI handle and supplies native command execution.
	server *server
}

// Command forwards borrowed request bytes to the private native handle.
func (backend *nativeBackend) Command(method string, request []byte) ([]byte, error) {
	response, err := backend.server.command(method, request)
	if err != nil {
		return nil, fmt.Errorf("execute native Anki command: %w", err)
	}
	return response, nil
}

// Close releases native state; the native destructor has no fallible result.
func (backend *nativeBackend) Close() error {
	backend.server.close()
	return nil
}
