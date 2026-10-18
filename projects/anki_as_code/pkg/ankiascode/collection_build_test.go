package ankiascode

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"google.golang.org/grpc"
)

// TestCollectionScopeBuild verifies copy-and-apply publishes independently of an open source scope.
func TestCollectionScopeBuild(t *testing.T) {
	fixture := newFixture(t)
	application := fixtureApplication()
	ctx := t.Context()
	configs, output := t.TempDir(), filepath.Join(t.TempDir(), "built.colpkg")
	writeFixtureTOML(t, filepath.Join(configs, "collection.anki.toml"), &collectionpb.CollectionConfig{Version: 1})
	original, err := os.ReadFile(fixture.base)
	if err != nil {
		t.Fatal(err)
	}
	if err := application.collections.WithCollection(ctx, fixture.base, func(connection grpc.ClientConnInterface) (bool, error) {
		if err := application.ColBuild(ctx, fixture.base, configs, output, ""); err != nil {
			return false, fmt.Errorf("build shared collection: %w", err)
		}
		if _, err := os.Stat(output); err != nil {
			return false, fmt.Errorf("build did not publish its independent output: %w", err)
		}
		value := &collectionpb.State{Version: 1, Resources: &collectionpb.Resources{Presets: map[string]*collectionpb.PresetState{"after-build": {Id: 1}}}}
		states := state.New(connection, "")
		changed, err := states.Replace(ctx, value)
		if err != nil {
			return false, fmt.Errorf("update state after build: %w", err)
		}
		if _, err := states.Dump(ctx); err != nil {
			return false, fmt.Errorf("read collection after build: %w", err)
		}
		return changed, nil
	}); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(fixture.base)
	if err != nil || string(actual) == string(original) {
		t.Fatalf("later source changes were not published to their own archive: %v", err)
	}
	stored, err := application.ColStateDump(ctx, output, "")
	if err != nil || stored.GetResources().GetPresets()["after-build"] != nil {
		t.Fatalf("independent build included later source changes: value=%v error=%v", stored, err)
	}
	stored, err = application.ColStateDump(ctx, fixture.base, "")
	if err != nil || stored.GetResources().GetPresets()["after-build"].GetId() != 1 {
		t.Fatalf("source scope lost its later changes: value=%v error=%v", stored, err)
	}
	saveFixtureArtifact(t, output, "scoped-build.colpkg")
}

// TestBuildFailureDestination verifies rejected builds preserve destination files and remove staging.
func TestBuildFailureDestination(t *testing.T) {
	for _, existing := range []bool{false, true} {
		fixture := newFixture(t)
		application := fixtureApplication()
		directory := t.TempDir()
		output := filepath.Join(directory, "result.colpkg")
		original, err := os.ReadFile(fixture.base)
		if err != nil {
			t.Fatal(err)
		}
		if existing {
			if err := os.WriteFile(output, original, 0o600); err != nil {
				t.Fatal(err)
			}
		}
		if err := application.ColBuild(t.Context(), fixture.base, filepath.Join(t.TempDir(), "missing"), output, ""); err == nil {
			t.Fatal("build accepted missing configurations")
		}
		if existing {
			actual, err := os.ReadFile(output)
			if err != nil || string(actual) != string(original) {
				t.Fatalf("failed build replaced its destination: %v", err)
			}
			saveFixtureArtifact(t, output, "failed-build-existing.colpkg")
		} else if _, err := os.Stat(output); !os.IsNotExist(err) {
			t.Fatalf("failed build created its destination: %v", err)
		}
		entries, err := os.ReadDir(directory)
		if err != nil {
			t.Fatal(err)
		}
		want := 0
		if existing {
			want = 1
		}
		if len(entries) != want {
			t.Fatalf("failed build retained staging: %v", entries)
		}
	}
}
