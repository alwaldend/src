package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

// App wires configuration and state operations around one managed collection.
// Configure Options before Start and serialize lifecycle calls. State operations
// borrow the prepared backend until Stop; App must not be copied after startup.
type App struct {
	// Options configures collection initialization and is read only by Start.
	Options Options
	// worker owns the native server process configured by Start.
	worker *backend.BackendSTD
	// collection owns the native collection prepared over worker's connection.
	collection *backend.BackendCollection
}

// Options configures collection initialization performed by Start.
type Options struct {
	// Input is the existing .anki2 database opened during application startup.
	Input string
}

// StateOptions selects the namespace used by WithState.
// The command layer binds flags directly to these application options.
type StateOptions struct {
	// Key selects the metadata namespace; an empty key uses state.DefaultKey.
	Key string
}

// Export writes an editable text representation of a collection.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Export(ctx context.Context, input, output string) error {
	return errors.New("export is not implemented")
}

// Plan describes changes needed to reconcile a collection with its text.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Plan(ctx context.Context, input, text string) (any, error) {
	return nil, errors.New("plan is not implemented")
}

// Apply reconciles an offline collection with its text representation.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Apply(ctx context.Context, input, text string) error {
	return errors.New("apply is not implemented")
}

// Build creates an archive from a base collection and desired text.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Build(ctx context.Context, base, text, output string) error {
	return errors.New("build is not implemented")
}

// Dump prints parsed configurations through the pure text adapter, without Anki access.
func (application *App) Dump(ctx context.Context, input, format string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("dump configuration: %w", err)
	}
	data, err := config.New(input).Dump(format)
	if err != nil {
		return nil, fmt.Errorf("dump configuration: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("finish configuration dump: %w", err)
	}
	return data, nil
}

// Start prepares one managed server and collection using Options.
// Call Stop even after failure. A configured backend cannot be started again.
func (application *App) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("prepare application startup: %w", err)
	}
	if application.worker != nil {
		return fmt.Errorf("application backend has already been initialized")
	}
	input := application.Options.Input
	if !strings.EqualFold(filepath.Ext(input), ".anki2") {
		return fmt.Errorf("collection input %q must be an existing .anki2 database", input)
	}
	path, err := filepath.Abs(input)
	if err != nil {
		return fmt.Errorf("resolve collection input %q: %w", input, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("inspect collection input %q: %w", input, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("collection input %q must be a regular database file", input)
	}
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		return fmt.Errorf("locate Anki server: %w", err)
	}
	application.worker = backend.New(backend.BackendSTDOptions{Binary: binary})
	application.collection = backend.NewCollection(backend.BackendCollectionOptions{Backend: application.worker, Path: path})
	if err := application.worker.Start(ctx); err != nil {
		return fmt.Errorf("start Anki server: %w", err)
	}
	if err := application.collection.Start(ctx); err != nil {
		return fmt.Errorf("open application collection: %w", err)
	}
	return nil
}

// Stop closes the collection before its worker and preserves both cleanup errors.
// It is safe before Start and after failed startup; ctx bounds shutdown.
func (application *App) Stop(ctx context.Context) error {
	var failure error
	if application.collection != nil {
		if err := application.collection.Stop(ctx); err != nil {
			failure = errors.Join(failure, fmt.Errorf("close collection: %w", err))
		}
	}
	if application.worker != nil {
		if err := application.worker.Stop(ctx); err != nil {
			failure = errors.Join(failure, fmt.Errorf("stop Anki server: %w", err))
		}
	}
	if failure != nil {
		return fmt.Errorf("stop application: %w", failure)
	}
	return nil
}

// WithState lends the prepared backend connection to a state operation.
// It owns no lifecycle or input path; use Start before calling and Stop afterward.
func (application *App) WithState(ctx context.Context, options StateOptions, operation func(*state.Store) error) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("prepare state operation: %w", err)
	}
	if application.collection == nil {
		return fmt.Errorf("application collection is not prepared")
	}
	connection, err := application.collection.Connection()
	if err != nil {
		return fmt.Errorf("get managed state connection: %w", err)
	}
	if err := operation(state.New(connection, options.Key)); err != nil {
		return fmt.Errorf("run state operation: %w", err)
	}
	return nil
}
