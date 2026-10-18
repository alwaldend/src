package ankiascode

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/collection"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"google.golang.org/protobuf/proto"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"github.com/ankitects/anki/proto/anki/import_export"
	"google.golang.org/grpc"
)

// TestApplicationServerBoundary verifies native delegation, startup cancellation and failure causes at the application boundary.
func TestApplicationServerBoundary(t *testing.T) {
	for _, serve := range []bool{false, true} {
		for _, fail := range []bool{false, true} {
			var captured []string
			failure := errors.New("native runner failed")
			application := New(Options{ServerRunner: func(arguments []string) error {
				captured = arguments
				if fail {
					return failure
				}
				return nil
			}})
			operation := application.ServerRun
			expected := []string{"application server run"}
			if serve {
				operation = application.ServerServe
				expected = append(expected, "serve")
			}
			arguments := []string{"--listen", "socket with spaces", "--future-native-option"}
			err := operation(t.Context(), expected[0], arguments)
			expected = append(expected, arguments...)
			if fail && !errors.Is(err, failure) || !fail && err != nil || strings.Join(captured, "\x00") != strings.Join(expected, "\x00") {
				t.Fatalf("native delegation serve=%t fail=%t: args=%v error=%v", serve, fail, captured, err)
			}
			if strings.Join(arguments, "\x00") != "--listen\x00socket with spaces\x00--future-native-option" {
				t.Fatal("application changed caller arguments")
			}
		}
	}
	application := New(Options{ServerRunner: func([]string) error {
		t.Fatal("cancelled startup invoked the native runner")
		return nil
	}})
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	for _, operation := range []func(context.Context, string, []string) error{application.ServerRun, application.ServerServe} {
		if err := operation(ctx, "application", nil); !errors.Is(err, context.Canceled) {
			t.Fatalf("server startup lost cancellation: %v", err)
		}
	}
	application = &App{}
	for _, operation := range []func(context.Context, string, []string) error{application.ServerRun, application.ServerServe} {
		if err := operation(t.Context(), "application", nil); err == nil {
			t.Fatal("unconfigured application accepted native hosting")
		}
	}
}

// TestDirectDatabaseWorkflow verifies native writes persist without a staging copy, including before a caller error.
func TestDirectDatabaseWorkflow(t *testing.T) {
	fixture := newFixture(t)
	scratch := t.TempDir()
	t.Setenv("TMPDIR", scratch)
	ctx := context.Background()
	expected := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"direct": {Id: 1}}}}
	for _, fail := range []bool{false, true} {
		if fail {
			expected.Resources.Presets["committed"] = &ankiascodepb.PresetState{Id: 1}
		}
		failure := errors.New("operation failed after its native write")
		err := runFixtureCollection(ctx, fixture.database, fixture.database, func(connection grpc.ClientConnInterface) (bool, error) {
			store := state.New(connection, "")
			changed, err := store.Replace(ctx, expected)
			if err != nil {
				return false, fmt.Errorf("replace direct state: %w", err)
			}
			if fail {
				return false, failure
			}
			return changed, nil
		})
		if fail && !errors.Is(err, failure) || !fail && err != nil {
			t.Fatalf("direct operation result: %v", err)
		}
		client := connectFixture(t, fixture.database)
		connection, err := client.worker.Connection()
		if err != nil {
			t.Fatal(err)
		}
		stored, err := state.New(connection, "").Dump(ctx)
		if err != nil || !proto.Equal(stored, expected) {
			t.Fatal("direct native write did not survive close")
		}
		if err := client.close(); err != nil {
			t.Fatal(err)
		}
	}
	files, err := os.ReadDir(scratch)
	if err != nil || len(files) != 0 {
		t.Fatalf("direct database workflow created private files: %v %v", files, err)
	}
	saveFixtureArtifact(t, fixture.database, "direct.anki2")
}

// TestNativePublicationWorkflow verifies completed writes survive operation failures and startup failures preserve sources.
func TestNativePublicationWorkflow(t *testing.T) {
	for _, scenario := range []string{"operation error", "cancelled operation", "cancelled open", "invalid package"} {
		t.Run(scenario, func(t *testing.T) {
			fixture := newFixture(t)
			input := fixture.base
			if scenario == "invalid package" {
				input = filepath.Join(t.TempDir(), "invalid.colpkg")
				if err := os.WriteFile(input, []byte("invalid collection package"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			original, err := os.ReadFile(input)
			if err != nil {
				t.Fatal(err)
			}
			scratch := t.TempDir()
			t.Setenv("TMPDIR", scratch)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if scenario == "cancelled open" {
				cancel()
			}
			called := false
			failure := errors.New("operation failed after a native write")
			err = runFixtureCollection(ctx, input, input, func(connection grpc.ClientConnInterface) (bool, error) {
				called = true
				expected := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"discarded": {Id: 1}}}}
				if _, err := state.New(connection, "").Replace(ctx, expected); err != nil {
					return false, fmt.Errorf("write restored collection: %w", err)
				}
				if scenario == "operation error" {
					return false, failure
				}
				cancel()
				return true, nil
			})
			if err == nil || scenario == "operation error" && !errors.Is(err, failure) {
				t.Fatalf("failed operation result: %v", err)
			}
			if called != (scenario == "operation error" || scenario == "cancelled operation") {
				t.Fatal("collection preparation did not gate the operation")
			}
			if called {
				actual := stateFixtureDump(t, input)
				if actual.GetResources().GetPresets()["discarded"].GetId() != 1 {
					t.Fatal("request failure suppressed committed archive writes")
				}
			} else if actual, err := os.ReadFile(input); err != nil || string(actual) != string(original) {
				t.Fatalf("failed startup changed the source archive: %v", err)
			}
			files, err := os.ReadDir(scratch)
			if err != nil || len(files) != 0 {
				t.Fatalf("normal native cleanup retained restore files: %v %v", files, err)
			}
			saveFixtureArtifact(t, input, "failed-"+scenario+".colpkg")
		})
	}
}

// TestCollectionLifecycleWorkflow verifies native persistence and connection shutdown through the app.
func TestCollectionLifecycleWorkflow(t *testing.T) {
	for _, format := range []string{"package", "database"} {
		t.Run(format, func(t *testing.T) {
			fixture := newFixture(t)
			input := fixture.base
			output := filepath.Join(t.TempDir(), "lifecycle.colpkg")
			if format == "database" {
				input, output = fixture.database, fixture.database
			}
			scratch := t.TempDir()
			t.Setenv("TMPDIR", scratch)
			ctx := context.Background()
			expected := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"lifecycle": {Id: 1}}}}
			var borrowed grpc.ClientConnInterface
			if err := runFixtureCollection(ctx, input, output, func(connection grpc.ClientConnInterface) (bool, error) {
				borrowed = connection
				changed, err := state.New(connection, "").Replace(ctx, expected)
				if err != nil {
					return false, fmt.Errorf("write collection state: %w", err)
				}
				return changed, nil
			}); err != nil {
				t.Fatal(err)
			}
			if _, err := state.New(borrowed, "").Dump(ctx); err == nil {
				t.Fatal("borrowed connection remained usable after the workflow")
			}
			files, err := os.ReadDir(scratch)
			if err != nil || len(files) != 0 {
				t.Fatalf("completed workflow retained restore files: %v %v", files, err)
			}
			if !proto.Equal(expected, stateFixtureDump(t, output)) {
				t.Fatal("native state did not survive collection shutdown")
			}
			if format == "package" {
				if !proto.Equal(expected, stateFixtureDump(t, input)) {
					t.Fatal("source write-back lost native changes")
				}
			} else {
				output = filepath.Join(t.TempDir(), "lifecycle.colpkg")
				packageFixture(t, input, output)
			}
			restored := inspectFixture(t, output)
			assertFixtureMedia(t, restored)
			if err := restored.close(); err != nil {
				t.Fatal(err)
			}
			saveFixtureArtifact(t, output, "lifecycle-"+format+".colpkg")
		})
	}
}

// TestStateInputFailures rejects unavailable or malformed sources without creating collections.
func TestStateInputFailures(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "collection.colpkg")
	if err := os.WriteFile(archive, []byte("invalid collection archive"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, input := range []string{filepath.Join(root, "missing.anki2"), root, archive} {
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		value, err := fixtureApplication().ColStateDump(ctx, ankiascodeapi.CollectionPath(input), "")
		cancel()
		if err == nil || value != nil {
			t.Fatalf("invalid collection input %q: error=%v state=%v", input, err, value)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "missing.anki2")); !os.IsNotExist(err) {
		t.Fatalf("state command created a missing database: %v", err)
	}
}

// runFixtureCollection exercises native callbacks through collection.Store.
func runFixtureCollection(ctx context.Context, input, output string, operation func(grpc.ClientConnInterface) (bool, error)) error {
	application := fixtureApplication()
	if err := application.collections.WithCollection(ctx, collection.CollectionOptions{Path: input, WriteBack: true}, func(connection grpc.ClientConnInterface) (bool, error) {
		changed, err := operation(connection)
		if err != nil {
			return false, fmt.Errorf("run native fixture operation: %w", err)
		}
		if changed && output != "" && !strings.EqualFold(filepath.Ext(input), ".anki2") {
			if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: output, IncludeMedia: true}); err != nil {
				return false, fmt.Errorf("export fixture collection: %w", err)
			}
		}
		return false, nil
	}); err != nil {
		return fmt.Errorf("run fixture collection: %w", err)
	}
	return nil
}
