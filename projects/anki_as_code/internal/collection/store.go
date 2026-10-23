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
	// FS supplies source inspection; backend storage remains independently owned.
	FS filesystem.FS
	// NewClient constructs clients and selects its default when an operation address is empty.
	NewClient ankigrpc.ClientFactory
	// CloseTimeout bounds collection publication and close, independently of request cancellation.
	// Nonpositive values select ten minutes; transport shutdown has a separate five-second budget.
	CloseTimeout time.Duration
}

// CollectionOptions uses the application contract for source and publication settings.
type CollectionOptions = ankiascodeapi.ColOpenOptions

// sourceReservation retains physical identity until its client finishes native cleanup.
type sourceReservation struct {
	// identity identifies filesystem aliases of the reserved source.
	identity os.FileInfo
	// client observes native completion after construction; nil reserves startup selection.
	client ankiclient.Client
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
	entries map[string]*sourceReservation
}

// New configures managed collection scopes; it starts no client or process.
func New(options Options) Store {
	return Store{filesystem: options.FS, newClient: options.NewClient, closeTimeout: options.CloseTimeout}
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

// Open reserves a source and starts one client with its optional collection state.
// Cancellation requests cleanup; source reservations remain until native completion.
func (store *Store) Open(ctx context.Context, options CollectionOptions) (_ ankiclient.Client, err error) {
	source, reservation, err := store.reserve(ctx, options.Path)
	if err != nil {
		return nil, fmt.Errorf("reserve collection: %w", err)
	}
	client, err := store.createClient(options.Address)
	if err != nil {
		store.releaseSource(source, reservation)
		return nil, fmt.Errorf("select collection client: %w", err)
	}
	store.mutex.Lock()
	reservation.client = client
	store.mutex.Unlock()
	opened := false
	defer func() {
		if !opened {
			err = errors.Join(err, client.Stop(context.WithoutCancel(ctx)))
		}
		store.observe(ctx, source, reservation)
	}()
	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("start collection client: %w", err)
	}
	if err := client.OpenCollection(ctx, ankiclient.CollectionOptions{Path: source, WriteBack: options.WriteBack, CloseTimeout: store.closeTimeout}); err != nil {
		return nil, fmt.Errorf("open managed collection: %w", err)
	}
	opened = true
	return client, nil
}

// WithCollection lends an open connection and waits for cleanup after callback completion or panic.
// Operation errors and cancellation retain enabled archive publication and cleanup diagnostics.
func (store *Store) WithCollection(ctx context.Context, options CollectionOptions, operation func(grpc.ClientConnInterface) error) (err error) {
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
		case <-reservation.client.Done():
		case <-ctx.Done():
			_ = reservation.client.Stop(context.Background())
			<-reservation.client.Done()
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
		if reservation.client != nil {
			select {
			case <-reservation.client.Done():
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
