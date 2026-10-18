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

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"google.golang.org/protobuf/proto"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"github.com/ankitects/anki/proto/anki/import_export"
	"google.golang.org/grpc"
)

// TestDirectDatabaseWorkflow verifies native writes persist without a staging copy, including before a caller error.
func TestDirectDatabaseWorkflow(t *testing.T) {
	fixture := newFixture(t)
	scratch := t.TempDir()
	t.Setenv("TMPDIR", scratch)
	ctx := context.Background()
	expected := &collectionpb.State{Version: 1, Resources: &collectionpb.Resources{Presets: map[string]*collectionpb.PresetState{"direct": {Id: 1}}}}
	for _, fail := range []bool{false, true} {
		if fail {
			expected.Resources.Presets["committed"] = &collectionpb.PresetState{Id: 1}
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

// TestNativePublicationWorkflow verifies that failed operations leave source archives unchanged.
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
				expected := &collectionpb.State{Version: 1, Resources: &collectionpb.Resources{Presets: map[string]*collectionpb.PresetState{"discarded": {Id: 1}}}}
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
			actual, err := os.ReadFile(input)
			if err != nil || string(actual) != string(original) {
				t.Fatalf("failed operation changed the source archive: %v", err)
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
			original, err := os.ReadFile(input)
			if err != nil {
				t.Fatal(err)
			}
			scratch := t.TempDir()
			t.Setenv("TMPDIR", scratch)
			ctx := context.Background()
			expected := &collectionpb.State{Version: 1, Resources: &collectionpb.Resources{Presets: map[string]*collectionpb.PresetState{"lifecycle": {Id: 1}}}}
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
				actual, err := os.ReadFile(input)
				if err != nil || string(actual) != string(original) {
					t.Fatal("writing a distinct native package changed its input")
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
		value, err := fixtureApplication().ColStateDump(ctx, input, "")
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
	if err := application.collections.WithCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
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
