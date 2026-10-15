package app_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/app"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// TestStateLifecycleWorkflow shares one prepared application across state operations.
// It covers unavailable/stopped access, repeated startup, nested calls, cancellation,
// callback errors and native persistence after explicit shutdown.
func TestStateLifecycleWorkflow(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	database := filepath.Join(t.TempDir(), "collection.anki2")
	createCollection(t, ctx, database)
	application := &app.App{Options: app.Options{Input: database}}
	t.Cleanup(func() { stopApplication(t, application) })
	options := app.StateOptions{Key: "anki_as_code/state/application"}
	called := false
	if err := application.WithState(ctx, options, func(*state.Store) error {
		called = true
		return nil
	}); err == nil || called {
		t.Fatal("unprepared application invoked a state operation")
	}
	if err := application.Start(ctx); err != nil {
		t.Fatal(err)
	}
	if err := application.Start(ctx); err == nil {
		t.Fatal("second startup replaced an active application")
	}
	expected := &collectionpb.State{Version: 1, Resources: &collectionpb.Resources{Presets: map[string]*collectionpb.PresetState{"shared": {Id: 1}}}}
	if err := application.WithState(ctx, options, func(store *state.Store) error {
		if _, err := store.Replace(ctx, expected); err != nil {
			return fmt.Errorf("write shared state: %w", err)
		}
		if err := application.WithState(ctx, options, func(nested *state.Store) error {
			value, err := nested.Dump(ctx)
			if err != nil {
				return fmt.Errorf("read nested state: %w", err)
			}
			if !proto.Equal(value, expected) {
				return fmt.Errorf("nested operation did not share the collection")
			}
			return nil
		}); err != nil {
			return fmt.Errorf("reuse prepared application: %w", err)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	failure := errors.New("state callback failure")
	if err := application.WithState(ctx, options, func(*state.Store) error {
		return fmt.Errorf("injected operation failure: %w", failure)
	}); !errors.Is(err, failure) {
		t.Fatalf("operation error cause was lost: %v", err)
	}
	cancelled, stop := context.WithCancel(ctx)
	stop()
	if err := application.WithState(cancelled, options, func(*state.Store) error {
		called = true
		return nil
	}); !errors.Is(err, context.Canceled) || called {
		t.Fatalf("cancelled state access: called=%v error=%v", called, err)
	}
	checkState(t, ctx, application, options, expected)
	checkState(t, ctx, application, app.StateOptions{}, &collectionpb.State{Version: 1})
	stopApplication(t, application)
	if err := application.WithState(ctx, options, func(*state.Store) error {
		called = true
		return nil
	}); err == nil || called {
		t.Fatal("stopped application invoked a state operation")
	}
	reopened := &app.App{Options: app.Options{Input: database}}
	t.Cleanup(func() { stopApplication(t, reopened) })
	if err := reopened.Start(ctx); err != nil {
		t.Fatal(err)
	}
	checkState(t, ctx, reopened, options, expected)
	stopApplication(t, reopened)
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := protojson.Marshal(expected)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "state.anki.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
		data, err = os.ReadFile(database)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "state.anki2"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// checkState reads a namespace without changing the application's lifetime.
func checkState(t *testing.T, ctx context.Context, application *app.App, options app.StateOptions, expected *collectionpb.State) {
	t.Helper()
	if err := application.WithState(ctx, options, func(store *state.Store) error {
		value, err := store.Dump(ctx)
		if err != nil {
			return fmt.Errorf("read application state: %w", err)
		}
		if !proto.Equal(value, expected) {
			return fmt.Errorf("application state differs: %v", value)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

// stopApplication bounds cleanup independently of the operation context.
func stopApplication(t *testing.T, application *app.App) {
	t.Helper()
	cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := application.Stop(cleanup); err != nil {
		t.Error(err)
	}
}

// createCollection creates a native fixture through the existing backend library.
func createCollection(t *testing.T, ctx context.Context, database string) {
	t.Helper()
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	worker := backend.New(backend.BackendSTDOptions{Binary: binary})
	collection := backend.NewCollection(backend.BackendCollectionOptions{Backend: worker, Path: database})
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := collection.Stop(cleanup); err != nil {
			t.Error(err)
		}
		if err := worker.Stop(cleanup); err != nil {
			t.Error(err)
		}
	}()
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	if err := collection.Start(ctx); err != nil {
		t.Fatal(err)
	}
}
