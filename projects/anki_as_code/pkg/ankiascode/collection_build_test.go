package ankiascode

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"google.golang.org/grpc"
)

// TestCollectionBuildParentTraversal publishes absent destinations after resolving parent links.
func TestCollectionBuildParentTraversal(t *testing.T) {
	for _, destination := range []string{"missing", "dangling", "missing-parent"} {
		t.Run(destination, func(t *testing.T) {
			fixture := newFixture(t)
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, "other", "deep"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink("other/deep", filepath.Join(root, "redirect")); err != nil {
				t.Fatal(err)
			}
			original, err := os.ReadFile(fixture.base)
			if err != nil {
				t.Fatal(err)
			}
			configs := t.TempDir()
			configPath := filepath.Join(configs, "collection.anki.toml")
			writeFixtureTOML(t, configPath, &ankiascodepb.CollectionConfig{Version: 1})
			configuration, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatal(err)
			}
			sentinel := filepath.Join(root, "built.colpkg")
			if err := os.WriteFile(sentinel, original, 0o600); err != nil {
				t.Fatal(err)
			}
			intended := filepath.Join(root, "other", "built.colpkg")
			output := root + "/redirect/../built.colpkg"
			if destination == "dangling" {
				if err := os.Symlink("missing.colpkg", intended); err != nil {
					t.Fatal(err)
				}
			} else if destination == "missing-parent" {
				intended = filepath.Join(root, "other", "absent", "built.colpkg")
				output = root + "/redirect/../absent/built.colpkg"
			}
			err = fixtureApplication().ColBuild(t.Context(), ankiascodeapi.CollectionOptions{Path: fixture.base}, configs, output)
			if destination == "missing-parent" {
				if !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("build accepted a missing output parent: %v", err)
				}
				if _, err := os.Stat(intended); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("failed build created output: %v", err)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if info, err := os.Lstat(intended); err != nil || !info.Mode().IsRegular() {
					t.Fatalf("build did not publish the intended regular package: %v", err)
				}
				if _, err := fixtureContent(t, intended); err != nil {
					t.Fatalf("intended package cannot be opened: %v", err)
				}
				saveFixtureArtifact(t, intended, "parent-traversal-built-"+destination+".colpkg")
			}
			for path, wanted := range map[string][]byte{fixture.base: original, sentinel: original, configPath: configuration} {
				if data, err := os.ReadFile(path); err != nil || string(data) != string(wanted) {
					t.Fatalf("build changed preserved input or lexical destination %q: %v", path, err)
				}
			}
			if _, err := os.Stat(filepath.Join(root, "other", "missing.colpkg")); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("build followed a dangling output link: %v", err)
			}
			saveFixtureArtifact(t, sentinel, "parent-traversal-preserved-"+destination+".colpkg")
		})
	}
}

// TestCollectionScopeBuild verifies copy-and-apply publishes independently of an open source scope.
func TestCollectionScopeBuild(t *testing.T) {
	fixture := newFixture(t)
	application := fixtureApplication()
	ctx := t.Context()
	configs, output := t.TempDir(), filepath.Join(t.TempDir(), "built.colpkg")
	writeFixtureTOML(t, filepath.Join(configs, "collection.anki.toml"), &ankiascodepb.CollectionConfig{Version: 1})
	original, err := os.ReadFile(fixture.base)
	if err != nil {
		t.Fatal(err)
	}
	if err := application.collections.WithCollection(ctx, ankiascodeapi.ColOpenOptions{Path: fixture.base, WriteBack: true}, func(connection grpc.ClientConnInterface) error {
		if err := application.ColBuild(ctx, ankiascodeapi.CollectionOptions{Path: fixture.base}, configs, output); err != nil {
			return fmt.Errorf("build shared collection: %w", err)
		}
		if _, err := os.Stat(output); err != nil {
			return fmt.Errorf("build did not publish its independent output: %w", err)
		}
		value := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"after-build": {Id: 1}}}}
		states := state.New(connection, "")
		if _, err := states.Replace(ctx, value); err != nil {
			return fmt.Errorf("update state after build: %w", err)
		}
		if _, err := states.Dump(ctx); err != nil {
			return fmt.Errorf("read collection after build: %w", err)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(fixture.base)
	if err != nil || string(actual) == string(original) {
		t.Fatalf("later source changes were not published to their own archive: %v", err)
	}
	stored, err := application.ColStateDump(ctx, ankiascodeapi.CollectionOptions{Path: output})
	if err != nil || stored.GetResources().GetPresets()["after-build"] != nil {
		t.Fatalf("independent build included later source changes: value=%v error=%v", stored, err)
	}
	stored, err = application.ColStateDump(ctx, ankiascodeapi.CollectionOptions{Path: fixture.base})
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
		if err := application.ColBuild(t.Context(), ankiascodeapi.CollectionOptions{Path: fixture.base}, filepath.Join(t.TempDir(), "missing"), output); err == nil {
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
