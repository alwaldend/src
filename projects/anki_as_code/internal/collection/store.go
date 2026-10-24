package collection

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	"google.golang.org/grpc"
)

// Options configures managed native scopes and their collection publication budget.
type Options struct {
	// FS supplies source inspection and collection-directory operations.
	// The server must independently access the same storage and paths.
	FS filesystem.FS
	// NewClient constructs clients and selects its default when an operation address is empty.
	NewClient ankigrpc.ClientFactory
	// CloseTimeout bounds collection publication and close, independently of request cancellation.
	// Nonpositive values select ten minutes; transport shutdown has a separate five-second budget.
	CloseTimeout time.Duration
}

// sourceReservation retains physical identity until its scope finishes native cleanup.
type sourceReservation struct {
	// identity identifies filesystem aliases of the reserved source.
	identity os.FileInfo
	// scope observes native completion after construction; nil reserves startup selection.
	scope *scope
}

// Store reserves collection sources within explicitly configured native scopes.
// Use New for native access. Store must not be copied after use.
type Store struct {
	// filesystem supplies source inspection, alias resolution and collection-directory storage.
	filesystem filesystem.FS
	// newClient delegates address interpretation to the client library.
	newClient ankigrpc.ClientFactory
	// closeTimeout allows native archive publication to outlive transport shutdown's short budget.
	closeTimeout time.Duration
	// mutex protects source reservations, never native work.
	mutex sync.Mutex
	// entries retains sources from reservation through native cleanup.
	entries map[string]*sourceReservation
}

// New configures managed collection scopes; it starts no client or process.
func New(options Options) Store {
	return Store{filesystem: options.FS, newClient: options.NewClient, closeTimeout: options.CloseTimeout}
}

// createClient delegates connection-string interpretation and construction to anki_grpc.
func (store *Store) createClient(address string) (ankiclient.StartableClient, error) {
	constructor := store.newClient
	if constructor == nil {
		if address == "" {
			return nil, fmt.Errorf("collection access requires a configured client factory or per-operation client address")
		}
		constructor = ankigrpc.NewClient
	}
	client, err := constructor(address)
	if err != nil {
		return nil, fmt.Errorf("construct Anki client: %w", err)
	}
	return client, nil
}

// Open reserves a source and returns one already-open application collection.
// Cancellation requests cleanup; source reservations remain until native completion.
func (store *Store) Open(ctx context.Context, options ankiascodeapi.ColOpenOptions) (_ ankiclient.Client, err error) {
	source, reservation, err := store.reserve(ctx, options.Path)
	if err != nil {
		return nil, fmt.Errorf("reserve collection: %w", err)
	}
	client, err := store.createClient(options.Address)
	if err != nil {
		store.releaseSource(source, reservation)
		return nil, fmt.Errorf("select collection client: %w", err)
	}
	owner := newScope(client, store.filesystem, source, options.WriteBack, store.closeTimeout)
	store.mutex.Lock()
	reservation.scope = owner
	store.mutex.Unlock()
	opened := false
	defer func() {
		if !opened {
			err = errors.Join(err, owner.Stop(context.WithoutCancel(ctx)))
		}
		store.observe(ctx, source, reservation)
	}()
	if err := owner.open(ctx); err != nil {
		return nil, fmt.Errorf("open managed collection: %w", err)
	}
	opened = true
	return owner, nil
}

// WithCollection lends an open connection and waits for cleanup after callback completion or panic.
// Operation errors and cancellation retain enabled archive publication and cleanup diagnostics.
func (store *Store) WithCollection(ctx context.Context, options ankiascodeapi.ColOpenOptions, operation func(grpc.ClientConnInterface) error) (err error) {
	if operation == nil {
		return fmt.Errorf("collection access requires an operation")
	}
	client, err := store.Open(ctx, options)
	if err != nil {
		return fmt.Errorf("open collection scope: %w", err)
	}
	defer func() {
		if closeErr := client.Stop(context.WithoutCancel(ctx)); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close collection scope: %w", closeErr))
		}
	}()
	connection, err := client.Connection()
	if err != nil {
		return fmt.Errorf("get open collection connection: %w", errors.Join(err, ctx.Err(), context.Cause(ctx)))
	}
	if err := operation(connection); err != nil {
		return fmt.Errorf("run collection callback: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("finish collection scope: %w", errors.Join(err, context.Cause(ctx)))
	}
	return nil
}

// observe handles session cancellation and removes only this reservation after native completion.
func (store *Store) observe(ctx context.Context, source string, reservation *sourceReservation) {
	go func() {
		select {
		case <-reservation.scope.done:
		case <-ctx.Done():
			_ = reservation.scope.Stop(context.Background())
			<-reservation.scope.done
		}
		store.releaseSource(source, reservation)
	}()
}

// releaseSource cannot remove a newer reservation for the same path.
func (store *Store) releaseSource(source string, reservation *sourceReservation) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if store.entries[source] == reservation {
		delete(store.entries, source)
	}
}

// reserve records physical source identity before startup so concurrent opens cannot bypass the guard.
func (store *Store) reserve(ctx context.Context, path string) (string, *sourceReservation, error) {
	if err := ctx.Err(); err != nil {
		return "", nil, fmt.Errorf("reserve cancelled collection: %w", errors.Join(err, context.Cause(ctx)))
	}
	source, info, err := resolveCollectionSource(store.filesystem, path)
	if err != nil {
		return "", nil, fmt.Errorf("resolve collection source: %w", err)
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()
	for reserved, reservation := range store.entries {
		if reservation.scope != nil {
			select {
			case <-reservation.scope.done:
				delete(store.entries, reserved)
				continue
			default:
			}
		}
		identity := reservation.identity
		if source == reserved || os.SameFile(info, identity) {
			return "", nil, fmt.Errorf("collection %q is already open", source)
		}
	}
	if store.entries == nil {
		store.entries = make(map[string]*sourceReservation)
	}
	reservation := &sourceReservation{identity: info}
	store.entries[source] = reservation
	return source, reservation, nil
}

// resolveCollectionSource resolves aliases and requires an existing regular source file.
func resolveCollectionSource(storage filesystem.FS, path string) (string, os.FileInfo, error) {
	source, err := storage.ResolvePath(path, false)
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
