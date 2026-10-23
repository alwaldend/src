package ankiclient

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

// client owns one transport and its current optional collection phase.
type client struct {
	// transport owns process, network or in-process communication.
	transport Transport
	// lifecycle protects selection of the current collection phase.
	lifecycle sync.Mutex
	// collection holds native storage state, never another Client.
	collection *collectionState
	// started permits one transport startup attempt.
	started atomic.Bool
	// ready permits connection and collection operations after startup.
	ready atomic.Bool
	// stopping rejects new operations once terminal cleanup begins.
	stopping atomic.Bool
	// stopOnce initiates terminal cleanup once.
	stopOnce sync.Once
	// stopped closes after bounded cleanup attempts store stopError.
	stopped chan struct{}
	// done closes after final native transport completion.
	done chan struct{}
	// stopError retains terminal cleanup diagnostics after stopped closes.
	stopError error
}

// NewClient takes ownership of an unstarted transport without performing I/O.
// An existing Client is already the lifecycle owner and is returned unchanged.
func NewClient(transport Transport) Client {
	if existing, ok := transport.(Client); ok {
		return existing
	}
	return &client{transport: transport, stopped: make(chan struct{}), done: make(chan struct{})}
}

// Start prepares the transport once without tying its lifetime to the startup context.
func (client *client) Start(ctx context.Context) error {
	if client.stopping.Load() || !client.started.CompareAndSwap(false, true) {
		return fmt.Errorf("start Anki client: already started or stopped")
	}
	var err error
	if client.transport == nil {
		err = fmt.Errorf("Anki client requires a transport")
	} else if startErr := client.transport.Start(ctx); startErr != nil {
		err = fmt.Errorf("start Anki transport: %w", startErr)
	}
	if err != nil {
		return fmt.Errorf("prepare Anki client: %w", errors.Join(err, client.Stop(context.WithoutCancel(ctx))))
	}
	client.ready.Store(true)
	return nil
}

// Connection explicitly lends current collection access or the ready raw transport.
func (client *client) Connection() (grpc.ClientConnInterface, error) {
	if !client.ready.Load() || client.stopping.Load() {
		return nil, fmt.Errorf("Anki client is not ready")
	}
	client.lifecycle.Lock()
	collection := client.collection
	client.lifecycle.Unlock()
	var connection grpc.ClientConnInterface
	var err error
	if collection != nil {
		connection, err = collection.connection()
	} else {
		connection, err = client.transport.Connection()
	}
	if err != nil {
		return nil, fmt.Errorf("get Anki client connection: %w", err)
	}
	return connection, nil
}

// OpenCollection attaches storage to the existing client without another lifecycle owner.
// Failed preparation remains available to CloseCollection or Stop for cleanup.
func (client *client) OpenCollection(ctx context.Context, options CollectionOptions) error {
	client.lifecycle.Lock()
	if !client.ready.Load() || client.stopping.Load() {
		client.lifecycle.Unlock()
		return fmt.Errorf("open collection: Anki client is not ready")
	}
	if client.collection != nil {
		client.lifecycle.Unlock()
		return fmt.Errorf("open collection: close the current collection first")
	}
	state := newCollectionState(client.transport, options)
	client.collection = state
	client.lifecycle.Unlock()
	if err := state.open(ctx); err != nil {
		return fmt.Errorf("open Anki collection: %w", err)
	}
	return nil
}

// CloseCollection closes the current storage phase while retaining the transport.
func (client *client) CloseCollection(ctx context.Context) error {
	client.lifecycle.Lock()
	state := client.collection
	client.lifecycle.Unlock()
	if state == nil {
		return nil
	}
	if err := state.close(ctx); err != nil {
		return fmt.Errorf("close Anki collection: %w", err)
	}
	client.lifecycle.Lock()
	if client.collection == state {
		client.collection = nil
	}
	client.lifecycle.Unlock()
	return nil
}

// Stop publishes and closes storage before stopping the transport under independent budgets.
// The caller's context bounds waiting; Done observes final native completion.
func (client *client) Stop(ctx context.Context) error {
	client.stopOnce.Do(func() {
		client.stopping.Store(true)
		client.ready.Store(false)
		go func() {
			client.lifecycle.Lock()
			state := client.collection
			client.lifecycle.Unlock()
			budget := 10 * time.Minute
			if state != nil && state.options.CloseTimeout > 0 {
				budget = state.options.CloseTimeout
			}
			cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), budget)
			collectionErr := client.CloseCollection(cleanup)
			cancel()
			var transportErr error
			if client.transport != nil {
				cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
				transportErr = client.transport.Stop(cleanup)
				cancel()
			}
			client.stopError = errors.Join(collectionErr, transportErr)
			close(client.stopped)
			if asynchronous, ok := client.transport.(AsyncClient); ok {
				<-asynchronous.Done()
			}
			close(client.done)
		}()
	})
	select {
	case <-client.stopped:
	case <-ctx.Done():
		return fmt.Errorf("wait for Anki client cleanup: %w", ctx.Err())
	}
	if client.stopError != nil {
		return fmt.Errorf("stop Anki client: %w", client.stopError)
	}
	select {
	case <-client.done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait for Anki transport completion: %w", ctx.Err())
	}
}

// Done observes final native completion, including after a cleanup error or cancelled wait.
func (client *client) Done() <-chan struct{} { return client.done }
