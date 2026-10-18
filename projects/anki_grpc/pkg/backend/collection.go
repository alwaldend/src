package backend

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/import_export"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// BackendCollectionOptions configures a collection on an already-started Backend.
type BackendCollectionOptions struct {
	// Backend lends its prepared connection; the caller retains process lifetime.
	Backend Backend
	// Path names an .anki2 database or .colpkg package. Anki may create a missing database.
	Path string
	// RestoreDirectory overrides a package's restore destination. It must be empty and is retained.
	// An empty value creates an owned temporary restore; database inputs reject this option.
	RestoreDirectory string
	// WriteBack enables automatic package publication; its zero value leaves source archives unchanged.
	// Database writes always commit in place.
	WriteBack bool
}

// BackendCollection implements Backend for collection lifetime on a borrowed server.
// Start opens storage; Stop writes packages back when explicitly enabled and closes storage.
// Neither starts nor stops options.Backend.
// Database writes commit in place. Callers stop this owner before its borrowed backend.
// Use NewCollection; a BackendCollection must not be copied after use.
type BackendCollection struct {
	// options are copied constructor settings; Backend remains borrowed.
	options BackendCollectionOptions
	// lifecycle protects client registration, readiness and restore-directory access.
	lifecycle sync.Mutex
	// operation serializes startup, shutdown and native consuming calls.
	operation chan struct{}
	// active counts unary calls under lifecycle so shutdown can wait for their results.
	active int
	// drained closes once shutdown has begun and all registered unary callers have returned.
	drained chan struct{}
	// started limits the owner to one preparation attempt.
	started bool
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
	// connection gates generated calls while retaining the underlying transport.
	connection *collectionConnection
	// stopOnce retains one terminal shutdown result for repeated callers.
	stopOnce sync.Once
	// stopError preserves native close and file cleanup errors.
	stopError error
}

// NewCollection configures a collection owner without touching its backend or files.
func NewCollection(options BackendCollectionOptions) *BackendCollection {
	return &BackendCollection{options: options, operation: make(chan struct{}, 1), drained: make(chan struct{})}
}

// Start borrows a prepared backend and opens the configured native collection once.
// Failed startup exposes no connection. Call Stop even after startup failure,
// then stop the borrowed backend; Start does not run a separate cleanup lifecycle.
func (backend *BackendCollection) Start(ctx context.Context) error {
	backend.lifecycle.Lock()
	if backend.started || backend.stopped {
		backend.lifecycle.Unlock()
		return fmt.Errorf("start Anki collection: already started or stopped")
	}
	backend.started = true
	backend.lifecycle.Unlock()
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
	if backend.options.Backend == nil {
		return fmt.Errorf("start Anki collection: a prepared backend is required")
	}
	connection, err := backend.options.Backend.Connection()
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
	backend.connection = &collectionConnection{owner: backend, ClientConnInterface: connection}
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
		backend.uncertain = err != nil && len(trailers.Get(ankiErrorTrailer)) != 1
		if err != nil {
			return fmt.Errorf("restore Anki collection package: %w", err)
		}
	}
	var trailers metadata.MD
	backend.uncertain = true
	_, err = collection.NewBackendCollectionServiceClient(connection).OpenCollection(ctx, request, grpc.Trailer(&trailers))
	backend.uncertain = err != nil && len(trailers.Get(ankiErrorTrailer)) != 1
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
func (backend *BackendCollection) prepare(source string) (*collection.OpenCollectionRequest, error) {
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

// Connection lends the gated connection only while the native collection is prepared.
// Retained generated clients reject after Stop, consuming calls or an unconfirmed response.
func (backend *BackendCollection) Connection() (grpc.ClientConnInterface, error) {
	backend.lifecycle.Lock()
	ready, connection := backend.ready && !backend.stopped, backend.connection
	backend.lifecycle.Unlock()
	if !ready {
		return nil, fmt.Errorf("get Anki collection connection: collection is not prepared")
	}
	if _, err := backend.options.Backend.Connection(); err != nil {
		return nil, fmt.Errorf("get Anki collection backend: %w", err)
	}
	return connection, nil
}

// Stop invalidates access, writes enabled packages back and removes completed temporary restores.
// It never stops the borrowed backend. Failed/cancelled native work retains restore
// files until the caller has stopped that backend and can safely remove them.
// Stop is terminal and idempotent, including before Start and after failed startup.
func (backend *BackendCollection) Stop(ctx context.Context) error {
	backend.stopOnce.Do(func() {
		backend.lifecycle.Lock()
		backend.stopped, backend.ready = true, false
		if backend.active == 0 {
			close(backend.drained)
		}
		backend.lifecycle.Unlock()
		if err := backend.acquire(ctx); err != nil {
			backend.stopError = fmt.Errorf("stop Anki collection: %w", err)
			return
		}
		defer backend.release()
		select {
		case <-backend.drained:
		case <-ctx.Done():
			backend.stopError = fmt.Errorf("wait for active Anki collection calls: %w", ctx.Err())
			return
		}
		if backend.uncertain {
			backend.stopError = fmt.Errorf("native collection completion is unconfirmed; stop the borrowed backend before cleanup")
			return
		}
		if backend.prepared && backend.options.WriteBack && strings.EqualFold(filepath.Ext(backend.source), ".colpkg") {
			connection := backend.connection.ClientConnInterface
			if !backend.opened {
				if _, err := collection.NewBackendCollectionServiceClient(connection).OpenCollection(ctx, backend.openRequest); err != nil {
					backend.stopError = fmt.Errorf("reopen Anki collection for write-back: %w", err)
					return
				}
				backend.opened = true
			}
			_, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: backend.source, IncludeMedia: true})
			if err != nil {
				backend.stopError = fmt.Errorf("write back Anki collection package %q: %w", backend.source, err)
				return
			}
			backend.opened = false
		}
		if backend.opened {
			_, err := collection.NewBackendCollectionServiceClient(backend.connection.ClientConnInterface).CloseCollection(ctx, &collection.CloseCollectionRequest{})
			if err != nil {
				backend.stopError = fmt.Errorf("close Anki collection: %w", err)
				return
			}
			backend.opened = false
		}
		if backend.temporary {
			if err := os.RemoveAll(backend.directory); err != nil {
				backend.stopError = fmt.Errorf("remove temporary Anki restore: %w", err)
			}
		}
	})
	if backend.stopError != nil {
		backend.lifecycle.Lock()
		directory := backend.directory
		backend.lifecycle.Unlock()
		if directory != "" {
			return fmt.Errorf("stop collection backend (restore %q retained): %w", directory, backend.stopError)
		}
		return fmt.Errorf("stop collection backend: %w", backend.stopError)
	}
	return nil
}

// acquire serializes native lifecycle operations while respecting cancellation.
func (backend *BackendCollection) acquire(ctx context.Context) error {
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
func (backend *BackendCollection) release() {
	<-backend.operation
}
