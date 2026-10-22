package collection

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	ankicol "github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// nativeCompletion observes final release of resources used by owned native work.
type nativeCompletion interface {
	// Done closes after owned native work and transport resources are released.
	Done() <-chan struct{}
}

// scope owns one application collection from preparation through native cleanup.
type scope struct {
	// rawClient owns the selected transport independently of collection storage.
	rawClient ankiclient.StartableClient
	// filesystem shares collection-directory operations with the application's storage.
	filesystem filesystem.FS
	// source is the reserved database or archive path.
	source string
	// writeBack enables publication to source for archive inputs.
	writeBack bool
	// closeTimeout bounds native publication and collection close.
	closeTimeout time.Duration
	// directory holds an owned archive restore, or is empty for direct databases.
	directory string
	// request records successful preparation and permits reopening after explicit export.
	request *ankicol.OpenCollectionRequest
	// rpc tracks calls made through the prepared collection connection.
	rpc *scopeConnection
	// mutex protects admission and completion of borrowed calls.
	mutex sync.Mutex
	// closing rejects new calls once shutdown begins.
	closing bool
	// active counts borrowed calls whose responses have not returned.
	active int
	// drained closes after shutdown begins and all borrowed calls return.
	drained chan struct{}
	// opened permits calls until native close or export consumes the collection.
	opened bool
	// uncertain retains any unconfirmed response until transport shutdown.
	uncertain bool
	// stopOnce starts exactly one cleanup attempt.
	stopOnce sync.Once
	// stopped publishes the bounded cleanup result to independent Stop callers.
	stopped chan struct{}
	// done releases Store reservations only after confirmed native completion.
	done chan struct{}
	// stopError retains the cleanup result after stopped closes.
	stopError error
}

// newScope allocates a single-use collection owner without starting native work.
func newScope(client ankiclient.StartableClient, storage filesystem.FS, source string, writeBack bool, closeTimeout time.Duration) *scope {
	if closeTimeout <= 0 {
		closeTimeout = 10 * time.Minute
	}
	return &scope{
		rawClient: client, filesystem: storage, source: source, writeBack: writeBack, closeTimeout: closeTimeout,
		drained: make(chan struct{}), stopped: make(chan struct{}), done: make(chan struct{}),
	}
}

// open prepares one collection before Store exposes the scope to callers.
func (owner *scope) open(ctx context.Context) error {
	database := owner.source
	switch strings.ToLower(filepath.Ext(owner.source)) {
	case ".anki2":
	case ".colpkg":
		parent := os.TempDir()
		if _, host := owner.filesystem.Afero().(*afero.OsFs); !host {
			parent = "/tmp"
			if err := owner.filesystem.MkdirAll(parent, 0o700); err != nil {
				return fmt.Errorf("create virtual temporary directory %q: %w", parent, err)
			}
		}
		directory, err := owner.filesystem.TempDir(parent, "anki-collection-")
		if err != nil {
			return fmt.Errorf("create temporary restore directory: %w", err)
		}
		owner.directory = directory
		database = filepath.Join(directory, "collection.anki2")
	default:
		return fmt.Errorf("collection path must end in .anki2 or .colpkg")
	}
	media := strings.TrimSuffix(database, filepath.Ext(database)) + ".media"
	request := &ankicol.OpenCollectionRequest{CollectionPath: database, MediaFolderPath: media, MediaDbPath: media + ".db2"}
	if err := owner.filesystem.MkdirAll(media, 0o700); err != nil {
		return fmt.Errorf("prepare Anki media directory: %w", err)
	}
	if err := owner.rawClient.Start(ctx); err != nil {
		return fmt.Errorf("start collection client: %w", err)
	}
	connection, err := owner.rawClient.Connection()
	if err != nil {
		return fmt.Errorf("get collection transport: %w", err)
	}
	owner.rpc = &scopeConnection{owner: owner, ClientConnInterface: connection}
	if owner.directory != "" {
		var trailers metadata.MD
		owner.uncertain = true
		_, err := import_export.NewBackendImportExportServiceClient(connection).ImportCollectionPackage(ctx, &import_export.ImportCollectionPackageRequest{
			BackupPath: owner.source, ColPath: database, MediaFolder: media, MediaDb: request.MediaDbPath,
		}, grpc.Trailer(&trailers))
		owner.uncertain = err != nil && len(trailers.Get(ankiclient.ErrorTrailer)) != 1
		if err != nil {
			return fmt.Errorf("restore Anki collection package: %w", err)
		}
	}
	var trailers metadata.MD
	owner.uncertain = true
	_, err = ankicol.NewBackendCollectionServiceClient(connection).OpenCollection(ctx, request, grpc.Trailer(&trailers))
	owner.uncertain = err != nil && len(trailers.Get(ankiclient.ErrorTrailer)) != 1
	if err != nil {
		return fmt.Errorf("open Anki collection: %w", err)
	}
	owner.request, owner.opened = request, true
	return nil
}

// Connection lends access until close, export, cancellation or shutdown consumes it.
func (owner *scope) Connection() (grpc.ClientConnInterface, error) {
	owner.mutex.Lock()
	ready := owner.opened && !owner.closing && !owner.uncertain
	owner.mutex.Unlock()
	if !ready {
		return nil, fmt.Errorf("get Anki collection connection: collection is not open")
	}
	if _, err := owner.rawClient.Connection(); err != nil {
		return nil, fmt.Errorf("get Anki collection transport: %w", err)
	}
	return owner.rpc, nil
}

// Stop starts owned cleanup once; each caller independently bounds its wait.
func (owner *scope) Stop(ctx context.Context) error {
	owner.stopOnce.Do(func() {
		owner.mutex.Lock()
		owner.closing = true
		if owner.active == 0 {
			close(owner.drained)
		}
		owner.mutex.Unlock()
		go owner.finish(context.WithoutCancel(ctx))
	})
	select {
	case <-owner.stopped:
	case <-ctx.Done():
		return fmt.Errorf("wait for collection cleanup: %w", errors.Join(ctx.Err(), context.Cause(ctx)))
	}
	if owner.stopError != nil {
		return fmt.Errorf("stop collection scope: %w", owner.stopError)
	}
	select {
	case <-owner.done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait for native collection completion: %w", errors.Join(ctx.Err(), context.Cause(ctx)))
	}
}

// finish publishes bounded cleanup errors and separately observes native completion.
func (owner *scope) finish(ctx context.Context) {
	closeContext, cancelClose := context.WithTimeout(ctx, owner.closeTimeout)
	released, collectionErr := owner.closeStorage(closeContext)
	cancelClose()
	stopContext, cancelStop := context.WithTimeout(ctx, 5*time.Second)
	transportErr := owner.rawClient.Stop(stopContext)
	cancelStop()
	if collectionErr != nil && owner.directory != "" {
		collectionErr = fmt.Errorf("restore %q retained: %w", owner.directory, collectionErr)
	}
	if transportErr != nil {
		transportErr = fmt.Errorf("stop collection transport: %w", transportErr)
	}
	completion, observesNative := owner.rawClient.(nativeCompletion)
	if !observesNative && !released {
		collectionErr = errors.Join(collectionErr, fmt.Errorf("external server completion is unconfirmed; this Store retains the source reservation"))
	}
	owner.stopError = errors.Join(collectionErr, transportErr)
	if observesNative {
		select {
		case <-completion.Done():
			close(owner.done)
			close(owner.stopped)
			return
		default:
		}
		close(owner.stopped)
		<-completion.Done()
		close(owner.done)
		return
	}
	if released {
		close(owner.done)
	}
	close(owner.stopped)
}

// closeStorage drains borrowed calls, publishes archives and removes released restores.
// Its release result remains true if native cleanup succeeds but local removal fails.
func (owner *scope) closeStorage(ctx context.Context) (bool, error) {
	select {
	case <-owner.drained:
	case <-ctx.Done():
		return false, fmt.Errorf("wait for active Anki collection calls: %w", ctx.Err())
	}
	if owner.uncertain {
		return false, fmt.Errorf("native collection completion is unconfirmed")
	}
	opened := owner.opened
	if owner.request != nil && owner.writeBack && owner.directory != "" {
		connection := owner.rpc.ClientConnInterface
		if !opened {
			if _, err := ankicol.NewBackendCollectionServiceClient(connection).OpenCollection(ctx, owner.request); err != nil {
				return false, fmt.Errorf("reopen Anki collection for write-back: %w", err)
			}
		}
		if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: owner.source, IncludeMedia: true}); err != nil {
			return false, fmt.Errorf("write back Anki collection package %q: %w", owner.source, err)
		}
		opened = false
	}
	if opened {
		if _, err := ankicol.NewBackendCollectionServiceClient(owner.rpc.ClientConnInterface).CloseCollection(ctx, &ankicol.CloseCollectionRequest{}); err != nil {
			return false, fmt.Errorf("close Anki collection: %w", err)
		}
	}
	if owner.directory != "" {
		if err := owner.filesystem.RemoveAll(owner.directory); err != nil {
			return true, fmt.Errorf("remove temporary Anki restore: %w", err)
		}
	}
	return true, nil
}
