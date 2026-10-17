package collection

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
)

// backendCleanupTimeout bounds each native owner cleanup independently of callback cancellation.
const backendCleanupTimeout = 5 * time.Second

// Store reserves collection sources and publishes changed archives within managed native scopes.
// Its zero value is ready to use. It must not be copied after use.
type Store struct {
	// backendOptions configures each owned worker; an empty Binary uses bundled-server/PATH discovery.
	backendOptions backend.BackendSTDOptions
	// mutex protects source reservations, never native work.
	mutex sync.Mutex
	// entries retains sources from reservation through native cleanup.
	entries map[string]os.FileInfo
}

// WithCollection reserves a source and lends its managed native connection.
// A true callback result publishes a package back to its source path before cleanup;
// errors and cancellation skip publication. Callers exporting elsewhere through
// Anki return false to skip source publication. Database writes persist in place.
// Callers construct services on the connection and finish all work before returning.
// Store starts and stops both native owners; retained connections reject later calls.
func (store *Store) WithCollection(ctx context.Context, path string, operation func(grpc.ClientConnInterface) (bool, error)) error {
	if operation == nil {
		return fmt.Errorf("collection access requires an operation")
	}
	source, err := store.reserve(ctx, path)
	if err != nil {
		return fmt.Errorf("reserve collection: %w", err)
	}
	defer func() {
		store.mutex.Lock()
		delete(store.entries, source)
		store.mutex.Unlock()
	}()
	worker, err := store.newBackend()
	if err != nil {
		return fmt.Errorf("configure collection worker: %w", err)
	}
	if err := withBackend(ctx, worker, func() error {
		prepared := backend.NewCollection(backend.BackendCollectionOptions{Backend: worker, Path: source})
		if err := withBackend(ctx, prepared, func() error {
			connection, err := prepared.Connection()
			if err != nil {
				return fmt.Errorf("get prepared collection connection: %w", err)
			}
			changed, err := operation(connection)
			if err != nil {
				return fmt.Errorf("run collection callback: %w", err)
			}
			if err := ctx.Err(); err != nil {
				return fmt.Errorf("finish collection scope: %w", errors.Join(err, context.Cause(ctx)))
			}
			if changed && !strings.EqualFold(filepath.Ext(source), ".anki2") {
				if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{
					OutPath: source, IncludeMedia: true,
				}); err != nil {
					return fmt.Errorf("publish collection scope: %w", err)
				}
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

// newBackend configures a fresh worker using explicit options or bundled-server/PATH discovery.
func (store *Store) newBackend() (*backend.BackendSTD, error) {
	options := store.backendOptions
	if options.Binary == "" {
		binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
		if err != nil {
			binary, err = exec.LookPath("anki_grpc")
		}
		if err != nil {
			return nil, fmt.Errorf("locate Anki server: %w", err)
		}
		options.Binary = binary
	}
	return backend.New(options), nil
}

// withBackend starts one owned backend and stops it after every startup attempt or callback exit.
// Cleanup is bounded independently of cancellation; callback and cleanup errors remain discoverable.
func withBackend(ctx context.Context, owned backend.Backend, operation func() error) (err error) {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start backend scope: %w", errors.Join(err, context.Cause(ctx)))
	}
	defer func() {
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), backendCleanupTimeout)
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
	source, info, err := resolveCollectionSource(path)
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
func resolveCollectionSource(path string) (string, os.FileInfo, error) {
	expanded, err := expandPath(path)
	if err != nil {
		return "", nil, fmt.Errorf("expand source path: %w", err)
	}
	source, err := filepath.EvalSymlinks(expanded)
	if err != nil {
		return "", nil, fmt.Errorf("resolve source aliases: %w", err)
	}
	info, err := os.Stat(source)
	if err != nil {
		return "", nil, fmt.Errorf("inspect source: %w", err)
	}
	if !info.Mode().IsRegular() {
		return "", nil, fmt.Errorf("collection source must be a regular file")
	}
	return source, info, nil
}
