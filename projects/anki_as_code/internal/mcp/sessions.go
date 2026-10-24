package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"google.golang.org/grpc"
)

// collectionSessions owns protocol handles while App owns their managed native scopes.
type collectionSessions struct {
	// context bounds session lifetimes by the hosting command rather than an opening request.
	context context.Context
	// application opens collections through the existing source reservation and cleanup owner.
	application ankiascodeapi.App
	// mutex protects handle registration, lookup and removal.
	mutex sync.Mutex
	// entries maps opaque protocol handles to prepared collection sessions.
	entries map[string]*collectionHandle
	// next supplies distinct handle names within this MCP connection.
	next uint64
}

// collectionHandle serializes complete borrowed operations on one open collection.
type collectionHandle struct {
	// client owns prepared collection access and cleanup through the application contract.
	client ankiclient.Client
	// cancel ends the independent session context after native cleanup finishes.
	cancel context.CancelFunc
	// operation prevents workflows and close from interleaving on one collection.
	operation chan struct{}
}

// register adds explicit open/close controls without constructing a worker during discovery.
func (sessions *collectionSessions) register(adapter *server) error {
	if err := adapter.add("col_open", "Open a collection once and return a handle for native and workflow calls. Publish enabled archives on col_close or disconnect.",
		toolSchema([]string{"collection", "address", "write_back"}, []string{"collection"}), toolSchema([]string{"handle"}, []string{"handle"}),
		func(ctx context.Context, value arguments) ([]byte, error) {
			handle, err := sessions.open(ctx, value)
			if err != nil {
				return nil, fmt.Errorf("open MCP collection: %w", err)
			}
			data, err := json.Marshal(map[string]string{"handle": handle})
			if err != nil {
				return nil, fmt.Errorf("encode collection handle: %w", err)
			}
			return data, nil
		}); err != nil {
		return fmt.Errorf("register collection open: %w", err)
	}
	if err := adapter.add("col_close", "Close an open handle, publish enabled archives and stop its worker. Errors do not imply rollback.",
		toolSchema([]string{"handle"}, []string{"handle"}), toolSchema(nil, nil),
		func(ctx context.Context, value arguments) ([]byte, error) {
			if err := sessions.close(ctx, value.Handle); err != nil {
				return nil, fmt.Errorf("close MCP collection: %w", err)
			}
			return []byte("{}"), nil
		}); err != nil {
		return fmt.Errorf("register collection close: %w", err)
	}
	return nil
}

// open detaches a successful collection lifetime from its transient tool request.
func (sessions *collectionSessions) open(ctx context.Context, value arguments) (string, error) {
	lifetime, cancel := context.WithCancel(sessions.context)
	stop := context.AfterFunc(ctx, cancel)
	writeBack := value.WriteBack == nil || *value.WriteBack
	client, err := sessions.application.ColOpen(lifetime, ankiascodeapi.ColOpenOptions{Path: value.Collection, Address: value.Address, WriteBack: writeBack})
	stop()
	if err != nil {
		cancel()
		return "", fmt.Errorf("prepare collection handle: %w", err)
	}
	sessions.mutex.Lock()
	if ctx.Err() != nil || lifetime.Err() != nil {
		sessions.mutex.Unlock()
		closeErr := client.Stop(context.Background())
		cancel()
		return "", fmt.Errorf("opening session ended before registration: %w", errors.Join(fmt.Errorf("opening request or MCP session ended"), ctx.Err(), lifetime.Err(), closeErr))
	}
	sessions.next++
	handle := "collection-" + strconv.FormatUint(sessions.next, 10)
	sessions.entries[handle] = &collectionHandle{client: client, cancel: cancel, operation: make(chan struct{}, 1)}
	sessions.mutex.Unlock()
	return handle, nil
}

// borrow lends a handle's connection and a release function for one complete operation.
func (sessions *collectionSessions) borrow(ctx context.Context, handle string) (grpc.ClientConnInterface, func(), error) {
	entry, err := sessions.lookup(handle)
	if err != nil {
		return nil, nil, fmt.Errorf("select open collection: %w", err)
	}
	if err := entry.acquire(ctx); err != nil {
		return nil, nil, fmt.Errorf("borrow open collection: %w", err)
	}
	connection, err := entry.client.Connection()
	if err != nil {
		entry.release()
		return nil, nil, fmt.Errorf("get open collection connection: %w", err)
	}
	return connection, entry.release, nil
}

// close publishes one session and forgets its handle only after successful cleanup.
func (sessions *collectionSessions) close(ctx context.Context, handle string) error {
	entry, err := sessions.lookup(handle)
	if err != nil {
		return fmt.Errorf("select collection to close: %w", err)
	}
	if err := entry.acquire(ctx); err != nil {
		return fmt.Errorf("wait to close collection: %w", err)
	}
	defer entry.release()
	if err := entry.client.Stop(ctx); err != nil {
		return fmt.Errorf("finish open collection: %w", err)
	}
	entry.cancel()
	sessions.mutex.Lock()
	delete(sessions.entries, handle)
	sessions.mutex.Unlock()
	return nil
}

// lookup rejects unknown handles without starting or opening native resources.
func (sessions *collectionSessions) lookup(handle string) (*collectionHandle, error) {
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()
	entry := sessions.entries[handle]
	if entry == nil {
		return nil, fmt.Errorf("unknown collection handle %q", handle)
	}
	return entry, nil
}

// closeAll closes independent collections concurrently under their existing bounded cleanup policy.
// SDK Run has drained request handlers before the hosting command calls this method.
func (sessions *collectionSessions) closeAll() error {
	results := make(chan error, len(sessions.entries))
	for _, entry := range sessions.entries {
		go func(entry *collectionHandle) {
			err := entry.client.Stop(context.Background())
			entry.cancel()
			results <- err
		}(entry)
	}
	var result error
	for range sessions.entries {
		result = errors.Join(result, <-results)
	}
	if result != nil {
		return fmt.Errorf("clean up MCP collection sessions: %w", result)
	}
	return nil
}

// acquire respects cancellation while waiting for a complete collection operation.
func (entry *collectionHandle) acquire(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("acquire collection operation: %w", err)
	}
	select {
	case entry.operation <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait for collection operation: %w", ctx.Err())
	}
}

// release allows the next operation or explicit close to use the collection.
func (entry *collectionHandle) release() { <-entry.operation }
