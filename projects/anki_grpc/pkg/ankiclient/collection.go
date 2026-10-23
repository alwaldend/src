package ankiclient

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/import_export"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// CollectionOptions configures one collection on a client-owned transport.
type CollectionOptions struct {
	// Path names an .anki2 database or .colpkg package. Anki may create a missing database.
	Path string
	// RestoreDirectory overrides a package's restore destination. It must be empty and is retained.
	// An empty value creates an owned temporary restore; database inputs reject this option.
	RestoreDirectory string
	// WriteBack enables automatic package publication; its zero value leaves source archives unchanged.
	// Database writes always commit in place.
	WriteBack bool
	// CloseTimeout bounds publication during Client.Stop; nonpositive values select ten minutes.
	CloseTimeout time.Duration
}

// collectionState tracks storage preparation, publication and native call completion.
// It is private state of one Client and does not implement a client lifecycle.
type collectionState struct {
	// options are copied storage and publication settings.
	options CollectionOptions
	// transport supplies native RPC access owned by the surrounding client.
	transport Transport
	// lifecycle protects client registration, readiness and restore-directory access.
	lifecycle sync.Mutex
	// operation serializes startup, shutdown and native consuming calls.
	operation chan struct{}
	// active counts unary calls under lifecycle so shutdown can wait for their results.
	active int
	// drained closes once shutdown has begun and all registered unary callers have returned.
	drained chan struct{}
	// ready permits generated-client calls only after successful native open.
	ready bool
	// stopped invalidates access as soon as shutdown begins.
	stopped bool
	// directory is the resolved package destination, including after preparation failure.
	directory string
	// temporary distinguishes an owned restore from caller-supplied storage.
	temporary bool
	// opened records successful native open until close or package export consumes it.
	opened bool
	// prepared records successful startup even after a client close or export consumes storage.
	prepared bool
	// source is the absolute archive publication path selected during startup.
	source string
	// openRequest retains native paths so consumed package storage can be reopened for publication.
	openRequest *collection.OpenCollectionRequest
	// uncertain records any unconfirmed native call; successful later calls cannot clear it.
	uncertain bool
	// rpc gates generated calls while retaining the underlying transport.
	rpc *collectionConnection
	// stopOnce retains one terminal shutdown result for repeated callers.
	stopOnce sync.Once
	// closed publishes the close result without blocking other callers past their deadlines.
	closed chan struct{}
	// stopError preserves native close and file cleanup errors.
	stopError error
}

// newCollectionState allocates one storage phase without touching its transport or files.
func newCollectionState(transport Transport, options CollectionOptions) *collectionState {
	return &collectionState{transport: transport, options: options, operation: make(chan struct{}, 1), drained: make(chan struct{}), closed: make(chan struct{})}
}

// open prepares collection storage on the surrounding client's ready transport.
func (backend *collectionState) open(ctx context.Context) error {
	if err := backend.acquire(ctx); err != nil {
		return fmt.Errorf("prepare Anki collection: %w", err)
	}
	defer backend.release()
	backend.lifecycle.Lock()
	stopped := backend.stopped
	backend.lifecycle.Unlock()
	if stopped {
		return fmt.Errorf("start Anki collection: stopped before preparation")
	}
	if backend.transport == nil {
		return fmt.Errorf("open Anki collection: a prepared transport is required")
	}
	connection, err := backend.transport.Connection()
	if err != nil {
		return fmt.Errorf("borrow prepared Anki backend: %w", err)
	}
	if backend.options.Path == "" {
		return fmt.Errorf("start Anki collection: path is required")
	}
	source, err := filepath.Abs(backend.options.Path)
	if err != nil {
		return fmt.Errorf("resolve Anki collection path: %w", err)
	}
	request, err := backend.prepare(source)
	if err != nil {
		return fmt.Errorf("prepare Anki collection files: %w", err)
	}
	backend.source, backend.openRequest = source, request
	backend.lifecycle.Lock()
	backend.rpc = &collectionConnection{owner: backend, ClientConnInterface: connection}
	backend.lifecycle.Unlock()
	if err := os.MkdirAll(request.MediaFolderPath, 0o700); err != nil {
		return fmt.Errorf("prepare Anki media directory: %w", err)
	}
	if strings.EqualFold(filepath.Ext(source), ".colpkg") {
		var trailers metadata.MD
		backend.uncertain = true
		_, err := import_export.NewBackendImportExportServiceClient(connection).ImportCollectionPackage(ctx, &import_export.ImportCollectionPackageRequest{
			BackupPath: source, ColPath: request.CollectionPath,
			MediaFolder: request.MediaFolderPath, MediaDb: request.MediaDbPath,
		}, grpc.Trailer(&trailers))
		backend.uncertain = err != nil && len(trailers.Get(ErrorTrailer)) != 1
		if err != nil {
			return fmt.Errorf("restore Anki collection package: %w", err)
		}
	}
	var trailers metadata.MD
	backend.uncertain = true
	_, err = collection.NewBackendCollectionServiceClient(connection).OpenCollection(ctx, request, grpc.Trailer(&trailers))
	backend.uncertain = err != nil && len(trailers.Get(ErrorTrailer)) != 1
	if err != nil {
		return fmt.Errorf("open Anki collection: %w", err)
	}
	backend.opened = true
	backend.prepared = true
	backend.lifecycle.Lock()
	defer backend.lifecycle.Unlock()
	if backend.stopped {
		return fmt.Errorf("start Anki collection: stopped during preparation")
	}
	backend.ready = true
	return nil
}

// prepare selects direct database paths or an empty native package restore destination.
func (backend *collectionState) prepare(source string) (*collection.OpenCollectionRequest, error) {
	database := source
	switch strings.ToLower(filepath.Ext(source)) {
	case ".anki2":
		if backend.options.RestoreDirectory != "" {
			return nil, fmt.Errorf("restore directory applies only to collection packages")
		}
	case ".colpkg":
		directory := backend.options.RestoreDirectory
		var err error
		if directory == "" {
			directory, err = os.MkdirTemp("", "anki-collection-")
			if err != nil {
				return nil, fmt.Errorf("create temporary restore directory: %w", err)
			}
			backend.temporary = true
		} else {
			directory, err = filepath.Abs(directory)
			if err != nil {
				return nil, fmt.Errorf("resolve restore directory: %w", err)
			}
			if err := os.MkdirAll(directory, 0o700); err != nil {
				return nil, fmt.Errorf("create restore directory: %w", err)
			}
			entries, err := os.ReadDir(directory)
			if err != nil {
				return nil, fmt.Errorf("inspect restore directory: %w", err)
			}
			if len(entries) != 0 {
				return nil, fmt.Errorf("restore directory must be empty")
			}
		}
		backend.lifecycle.Lock()
		backend.directory = directory
		backend.lifecycle.Unlock()
		database = filepath.Join(directory, "collection.anki2")
	default:
		return nil, fmt.Errorf("collection path must end in .anki2 or .colpkg")
	}
	media := strings.TrimSuffix(database, filepath.Ext(database)) + ".media"
	return &collection.OpenCollectionRequest{
		CollectionPath:  database,
		MediaFolderPath: media,
		MediaDbPath:     media + ".db2",
	}, nil
}

// connection lends native access only while this storage phase is prepared.
// Retained generated clients reject after close, consuming calls or an unconfirmed response.
func (backend *collectionState) connection() (grpc.ClientConnInterface, error) {
	backend.lifecycle.Lock()
	ready, connection := backend.ready && !backend.stopped, backend.rpc
	backend.lifecycle.Unlock()
	if !ready {
		return nil, fmt.Errorf("get Anki collection connection: collection is not prepared")
	}
	if _, err := backend.transport.Connection(); err != nil {
		return nil, fmt.Errorf("get Anki collection backend: %w", err)
	}
	return connection, nil
}

// close invalidates access, writes enabled packages back and removes completed temporary restores.
// It leaves the transport running. Failed or cancelled native work retains restore
// files until Client.Done confirms that transport cleanup has completed.
// Each storage phase retains one close result, including after failed preparation.
func (backend *collectionState) close(ctx context.Context) error {
	performClose := false
	backend.stopOnce.Do(func() {
		performClose = true
		backend.lifecycle.Lock()
		backend.stopped, backend.ready = true, false
		if backend.active == 0 {
			close(backend.drained)
		}
		backend.lifecycle.Unlock()
	})
	var closeErr error
	if performClose {
		defer close(backend.closed)
		backend.stopError = backend.closeStorage(ctx)
	} else {
		select {
		case <-backend.closed:
		case <-ctx.Done():
			closeErr = fmt.Errorf("wait for collection close: %w", ctx.Err())
		}
	}
	if closeErr == nil {
		closeErr = backend.stopError
	}
	if closeErr != nil {
		backend.lifecycle.Lock()
		directory := backend.directory
		backend.lifecycle.Unlock()
		if directory != "" {
			return fmt.Errorf("stop collection backend (restore %q retained): %w", directory, closeErr)
		}
		return fmt.Errorf("stop collection backend: %w", closeErr)
	}
	return nil
}

// closeStorage drains registered calls, publishes archives and releases confirmed storage.
func (backend *collectionState) closeStorage(ctx context.Context) error {
	if err := backend.acquire(ctx); err != nil {
		return fmt.Errorf("stop Anki collection: %w", err)
	}
	defer backend.release()
	select {
	case <-backend.drained:
	case <-ctx.Done():
		return fmt.Errorf("wait for active Anki collection calls: %w", ctx.Err())
	}
	if backend.uncertain {
		return fmt.Errorf("native collection completion is unconfirmed; stop the client before cleanup")
	}
	if backend.prepared && backend.options.WriteBack && strings.EqualFold(filepath.Ext(backend.source), ".colpkg") {
		connection := backend.rpc.ClientConnInterface
		if !backend.opened {
			if _, err := collection.NewBackendCollectionServiceClient(connection).OpenCollection(ctx, backend.openRequest); err != nil {
				return fmt.Errorf("reopen Anki collection for write-back: %w", err)
			}
			backend.opened = true
		}
		_, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: backend.source, IncludeMedia: true})
		if err != nil {
			return fmt.Errorf("write back Anki collection package %q: %w", backend.source, err)
		}
		backend.opened = false
	}
	if backend.opened {
		_, err := collection.NewBackendCollectionServiceClient(backend.rpc.ClientConnInterface).CloseCollection(ctx, &collection.CloseCollectionRequest{})
		if err != nil {
			return fmt.Errorf("close Anki collection: %w", err)
		}
		backend.opened = false
	}
	if backend.temporary {
		if err := os.RemoveAll(backend.directory); err != nil {
			return fmt.Errorf("remove temporary Anki restore: %w", err)
		}
	}
	return nil
}

// acquire serializes native lifecycle operations while respecting cancellation.
func (backend *collectionState) acquire(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("check collection operation context: %w", err)
	}
	select {
	case backend.operation <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait for collection operation: %w", ctx.Err())
	}
}

// release finishes a serialized native lifecycle operation.
func (backend *collectionState) release() {
	<-backend.operation
}
