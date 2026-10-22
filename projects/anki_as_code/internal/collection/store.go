package collection

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	"google.golang.org/grpc"
)

// backendCleanupTimeout bounds transport shutdown independently of collection publication.
const backendCleanupTimeout = 5 * time.Second

// Options configures managed native scopes and their collection publication budget.
type Options struct {
	// FS supplies source inspection; backend storage remains independently owned.
	FS filesystem.FS
	// NewClient constructs clients and selects its default when an operation address is empty.
	NewClient ankigrpc.ClientFactory
	// CloseTimeout bounds collection publication and close, independently of request cancellation.
	// Nonpositive values select ten minutes; transport shutdown has a separate five-second budget.
	CloseTimeout time.Duration
}

// CollectionOptions selects the source and publication policy for one offline scope.
type CollectionOptions struct {
	// Address selects a client; an empty address delegates default selection to the factory.
	Address string
	// Path identifies an existing collection database or package archive.
	Path string
	// WriteBack enables archive publication during cleanup; its zero value leaves source archives unchanged.
	// Database writes persist in place regardless of this setting.
	WriteBack bool
}

// Store reserves collection sources within explicitly configured native scopes.
// Use New for native access. Store must not be copied after use.
type Store struct {
	// filesystem supplies source path inspection and alias resolution.
	filesystem filesystem.FS
	// newClient delegates address interpretation to the client library.
	newClient ankigrpc.ClientFactory
	// closeTimeout allows native archive publication to outlive transport shutdown's short budget.
	closeTimeout time.Duration
	// mutex protects source reservations, never native work.
	mutex sync.Mutex
	// entries retains sources from reservation through native cleanup.
	entries map[string]os.FileInfo
}

// New configures managed collection scopes; it starts no client or process.
func New(options Options) Store {
	closeTimeout := options.CloseTimeout
	if closeTimeout <= 0 {
		closeTimeout = 10 * time.Minute
	}
	return Store{filesystem: options.FS, newClient: options.NewClient, closeTimeout: closeTimeout}
}

// createClient delegates connection-string interpretation and construction to anki_grpc.
func (store *Store) createClient(address string) (ankiclient.Client, error) {
	constructor := store.newClient
	if constructor == nil {
		if address == "" {
			return nil, fmt.Errorf("collection access requires a configured client factory or per-operation client address")
		}
		constructor = ankigrpc.ClientFactoryFunc[ankiclient.Client](ankigrpc.NewClient)
	}
	client, err := constructor.NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("construct Anki client: %w", err)
	}
	return client, nil
}

// WithCollection reserves a source and lends its managed native connection.
// ClientCollection publishes archives during cleanup when WriteBack is enabled,
// independently of callback results. Database writes persist in place.
// Callers construct services on the connection and finish all work before returning.
// Store starts and stops both native owners; retained connections reject later calls.
func (store *Store) WithCollection(ctx context.Context, options CollectionOptions, operation func(grpc.ClientConnInterface) error) error {
	if operation == nil {
		return fmt.Errorf("collection access requires an operation")
	}
	source, err := store.reserve(ctx, options.Path)
	if err != nil {
		return fmt.Errorf("reserve collection: %w", err)
	}
	worker, err := store.createClient(options.Address)
	if err != nil {
		store.releaseSource(source)
		return fmt.Errorf("select collection client: %w", err)
	}
	defer func() {
		if asynchronous, ok := worker.(ankiclient.AsyncClient); ok {
			stopped, cancel := context.WithCancel(context.Background())
			cancel()
			// Ensure terminal cleanup even if the scope's context rejected startup.
			// The scope already reports cleanup errors; Done observes completion.
			_ = worker.Stop(stopped)
			done := asynchronous.Done()
			select {
			case <-done:
			default:
				go func() {
					<-done
					store.releaseSource(source)
				}()
				return
			}
		}
		store.releaseSource(source)
	}()
	if err := withBackend(ctx, worker, backendCleanupTimeout, func() error {
		prepared := ankiclient.NewCollection(ankiclient.ClientCollectionOptions{Backend: worker, Path: source, WriteBack: options.WriteBack})
		if err := withBackend(ctx, prepared, store.closeTimeout, func() error {
			connection, err := prepared.Connection()
			if err != nil {
				return fmt.Errorf("get prepared collection connection: %w", err)
			}
			if err := operation(connection); err != nil {
				return fmt.Errorf("run collection callback: %w", err)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("run prepared collection: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("run collection scope: %w", err)
	}
	return nil
}

// releaseSource removes a reservation only after its owned backend has finished cleanup.
func (store *Store) releaseSource(source string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	delete(store.entries, source)
}

// withBackend starts one owned backend and stops it after every startup attempt or callback exit.
// Cleanup is bounded independently of cancellation; callback and cleanup errors remain discoverable.
func withBackend(ctx context.Context, owned ankiclient.Client, closeTimeout time.Duration, operation func() error) (err error) {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start backend scope: %w", errors.Join(err, context.Cause(ctx)))
	}
	defer func() {
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), closeTimeout)
		defer cancel()
		if stopErr := owned.Stop(cleanup); stopErr != nil {
			err = errors.Join(err, fmt.Errorf("stop collection scope backend: %w", stopErr))
		}
	}()
	if err := owned.Start(ctx); err != nil {
		return fmt.Errorf("start collection scope backend: %w", err)
	}
	if err := operation(); err != nil {
		return fmt.Errorf("run prepared backend: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("finish backend scope: %w", errors.Join(err, context.Cause(ctx)))
	}
	return nil
}

// reserve records physical source identity before startup so concurrent opens cannot bypass the guard.
func (store *Store) reserve(ctx context.Context, path string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("reserve cancelled collection: %w", errors.Join(err, context.Cause(ctx)))
	}
	source, info, err := resolveCollectionSource(store.filesystem, path)
	if err != nil {
		return "", fmt.Errorf("resolve collection source: %w", err)
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()
	for reserved, identity := range store.entries {
		if source == reserved || os.SameFile(info, identity) {
			return "", fmt.Errorf("collection %q is already open", source)
		}
	}
	if store.entries == nil {
		store.entries = make(map[string]os.FileInfo)
	}
	store.entries[source] = info
	return source, nil
}

// resolveCollectionSource resolves aliases and requires an existing regular source file.
func resolveCollectionSource(storage filesystem.FS, path string) (string, os.FileInfo, error) {
	expanded, err := storage.ExpandPath(path)
	if err != nil {
		return "", nil, fmt.Errorf("expand source path: %w", err)
	}
	source, err := storage.ResolvePath(expanded, false)
	if err != nil {
		return "", nil, fmt.Errorf("resolve source aliases: %w", err)
	}
	info, err := storage.Stat(source)
	if err != nil {
		return "", nil, fmt.Errorf("inspect source: %w", err)
	}
	if !info.Mode().IsRegular() {
		return "", nil, fmt.Errorf("collection source must be a regular file")
	}
	return source, info, nil
}
