//go:build cgo && linux && amd64

package native

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// Client owns one in-process connection and its native resource lifetime.
// It constructs the connection for generated Anki clients without starting a process or listener.
// Use New or the zero value; do not copy a client after use.
type Client struct {
	// connection is constructed with the client and retained until native cleanup finishes.
	connection connection
}

// New constructs an unstarted client using upstream Anki initialization defaults.
func New() *Client { return &Client{} }

// Start initializes the owned connection synchronously, checking ctx before and after.
// Cancellation does not interrupt native initialization or own the client lifetime.
func (client *Client) Start(ctx context.Context) error {
	if err := client.connection.start(ctx); err != nil {
		return fmt.Errorf("start in-process Anki client: %w", err)
	}
	return nil
}

// Connection lends the client's connection after startup and before shutdown.
func (client *Client) Connection() (grpc.ClientConnInterface, error) {
	client.connection.mutex.Lock()
	defer client.connection.mutex.Unlock()
	if client.connection.native == nil || client.connection.stopped {
		return nil, fmt.Errorf("get in-process Anki connection: client is not ready")
	}
	return &client.connection, nil
}

// Stop rejects new calls and releases native state after active calls finish.
// A deadline bounds waiting while cleanup continues; a later Stop can await it.
// Close an open collection through its generated client before stopping the client.
func (client *Client) Stop(ctx context.Context) error {
	if err := client.connection.stop(ctx); err != nil {
		return fmt.Errorf("stop in-process Anki client: %w", err)
	}
	return nil
}

// Done closes after Stop finishes startup and native cleanup, including cancelled calls.
// It returns the same channel before startup, while stopping and after cleanup.
func (client *Client) Done() <-chan struct{} {
	client.connection.mutex.Lock()
	defer client.connection.mutex.Unlock()
	if client.connection.done == nil {
		client.connection.done = make(chan struct{})
	}
	return client.connection.done
}
