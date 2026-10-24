package rebuild_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/rebuild"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/pkg/ankiascode"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// TestRootedBuildSourceConflict rejects two virtual spellings of the same memory archive.
func TestRootedBuildSourceConflict(t *testing.T) {
	for _, source := range []string{"base.colpkg", "/base.colpkg"} {
		t.Run(source, func(t *testing.T) {
			storage := afero.NewBasePathFs(afero.NewMemMapFs(), "/scope")
			if err := storage.MkdirAll(".", 0o755); err != nil {
				t.Fatal(err)
			}
			if err := afero.WriteFile(storage, source, []byte("original"), 0o600); err != nil {
				t.Fatal(err)
			}
			called := false
			builder := rebuild.New(func(_ context.Context, path string) error {
				called = true
				if err := afero.WriteFile(storage, path, []byte("updated"), 0o600); err != nil {
					return fmt.Errorf("edit stage: %w", err)
				}
				return nil
			}, filesystem.New(storage))
			output := "/base.colpkg"
			if source == output {
				output = "base.colpkg"
			}
			err := builder.Build(t.Context(), source, output)
			if err == nil || !strings.Contains(err.Error(), "output must differ") || called {
				t.Errorf("source alias was not rejected before staging: called=%t error=%v", called, err)
			}
			if data, err := afero.ReadFile(storage, source); err != nil || string(data) != "original" {
				t.Errorf("source alias changed original archive: %q %v", data, err)
			}
			if entries, err := afero.ReadDir(storage, "."); err != nil || len(entries) != 1 {
				t.Errorf("source alias left staging files: %v %v", entries, err)
			}
		})
	}
}

// TestBuildSymlinkParentTraversal follows intermediate links before their parent components.
func TestBuildSymlinkParentTraversal(t *testing.T) {
	for _, sourceKind := range []string{"link", "input"} {
		for _, rooted := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s-%t", sourceKind, rooted), func(t *testing.T) {
				root := t.TempDir()
				if err := os.MkdirAll(filepath.Join(root, "other", "deep"), 0o755); err != nil {
					t.Fatal(err)
				}
				for path, data := range map[string]string{"base.colpkg": "wrong", "other/base.colpkg": "original"} {
					if err := os.WriteFile(filepath.Join(root, path), []byte(data), 0o600); err != nil {
						t.Fatal(err)
					}
				}
				if err := os.Symlink("other/deep", filepath.Join(root, "redirect")); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("redirect/../base.colpkg", filepath.Join(root, "source.colpkg")); err != nil {
					t.Fatal(err)
				}
				var storage afero.Fs = afero.NewOsFs()
				source, output := filepath.Join(root, "source.colpkg"), filepath.Join(root, "built.colpkg")
				if rooted {
					storage = afero.NewBasePathFs(storage, root)
					source, output = "source.colpkg", "built.colpkg"
				}
				if sourceKind == "input" {
					source = root + string(filepath.Separator) + "redirect/../base.colpkg"
					if rooted {
						source = "redirect/../base.colpkg"
					}
				}
				builder := rebuild.New(func(_ context.Context, path string) error {
					data, err := afero.ReadFile(storage, path)
					if err != nil {
						return fmt.Errorf("read stage: %w", err)
					}
					if string(data) != "original" {
						return fmt.Errorf("resolved the wrong source archive: %q", data)
					}
					return nil
				}, filesystem.New(storage))
				if err := builder.Build(t.Context(), source, output); err != nil {
					t.Fatal(err)
				}
				if data, err := afero.ReadFile(storage, output); err != nil || string(data) != "original" {
					t.Fatalf("published wrong archive: %q %v", data, err)
				}
			})
		}
	}
}

// TestBasePathArchiveBuild verifies source and destination aliases through rooted filesystem adapters.
func TestBasePathArchiveBuild(t *testing.T) {
	for _, namespace := range []string{"absolute", "relative", "nested", "relative-link"} {
		t.Run(namespace, func(t *testing.T) {
			root := t.TempDir()
			storage := afero.NewBasePathFs(afero.NewOsFs(), root)
			if namespace == "nested" {
				storage = afero.NewBasePathFs(storage, "scope")
			}
			for _, directory := range []string{"archives", "results"} {
				if err := storage.MkdirAll(directory, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			if err := afero.WriteFile(storage, "archives/base.colpkg", []byte("original"), 0o600); err != nil {
				t.Fatal(err)
			}
			if namespace == "relative-link" {
				if err := os.Symlink("archives/base.colpkg", filepath.Join(root, "source.colpkg")); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("results", filepath.Join(root, "destination")); err != nil {
					t.Fatal(err)
				}
			} else {
				linker := storage.(afero.Linker)
				if err := linker.SymlinkIfPossible("archives/base.colpkg", "source.colpkg"); err != nil {
					t.Fatal(err)
				}
				if err := linker.SymlinkIfPossible("results", "destination"); err != nil {
					t.Fatal(err)
				}
			}
			source, output := "source.colpkg", "destination/built.colpkg"
			if namespace != "relative" {
				source, output = "/"+source, "/"+output
			}
			var staging string
			builder := rebuild.New(func(_ context.Context, path string) error {
				staging = filepath.Dir(path)
				data, err := afero.ReadFile(storage, path)
				if err != nil {
					return fmt.Errorf("read rooted stage: %w", err)
				}
				if string(data) != "original" {
					return fmt.Errorf("rooted stage lost source bytes")
				}
				if err := afero.WriteFile(storage, path, []byte("updated"), 0o600); err != nil {
					return fmt.Errorf("edit rooted stage: %w", err)
				}
				return nil
			}, filesystem.New(storage))
			if err := builder.Build(t.Context(), source, output); err != nil {
				t.Fatal(err)
			}
			for path, wanted := range map[string]string{"archives/base.colpkg": "original", "results/built.colpkg": "updated"} {
				if data, err := afero.ReadFile(storage, path); err != nil || string(data) != wanted {
					t.Fatalf("rooted archive %s: %q %v", path, data, err)
				}
			}
			if _, err := storage.Stat(staging); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("rooted staging retained: %v", err)
			}
			if info, _, err := storage.(afero.Lstater).LstatIfPossible("destination"); err != nil || info.Mode()&os.ModeSymlink == 0 {
				t.Fatalf("destination alias was replaced: %v", err)
			}
			if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
				data, err := json.MarshalIndent(map[string]string{"namespace": namespace, "source": "original", "output": "updated"}, "", "  ")
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(artifacts, "basepath-build-"+namespace+".json"), data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestMemoryArchiveBuild verifies staging, failure preservation and cleanup entirely in injected storage.
func TestMemoryArchiveBuild(t *testing.T) {
	for _, outcome := range []string{"success", "apply-error", "cancel"} {
		t.Run(outcome, func(t *testing.T) {
			storage := afero.NewMemMapFs()
			if err := storage.MkdirAll("/archives", 0o755); err != nil {
				t.Fatal(err)
			}
			base, output := "/archives/base.colpkg", "/archives/output.colpkg"
			for _, path := range []string{base, output} {
				if err := afero.WriteFile(storage, path, []byte("original"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			ctx, cancel := context.WithCancelCause(t.Context())
			defer cancel(nil)
			failure := errors.New("memory apply failed")
			var staged string
			builder := rebuild.New(func(_ context.Context, path string) error {
				staged = filepath.Dir(path)
				if err := afero.WriteFile(storage, path, []byte("updated"), 0o600); err != nil {
					return fmt.Errorf("edit memory stage: %w", err)
				}
				switch outcome {
				case "apply-error":
					return failure
				case "cancel":
					cancel(failure)
				}
				return nil
			}, filesystem.New(storage))
			err := builder.Build(ctx, base, output)
			if outcome == "success" && err != nil || outcome != "success" && !errors.Is(err, failure) {
				t.Fatalf("%s: %v", outcome, err)
			}
			want := "original"
			if outcome == "success" {
				want = "updated"
			}
			if data, err := afero.ReadFile(storage, output); err != nil || string(data) != want {
				t.Fatalf("published memory output: %q %v", data, err)
			}
			if data, err := afero.ReadFile(storage, base); err != nil || string(data) != "original" {
				t.Fatalf("source changed: %q %v", data, err)
			}
			if _, err := storage.Stat(staged); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("memory staging retained: %v", err)
			}
		})
	}
}

// TestArchiveBuild verifies staged publication using real Anki archives and native state writes.
func TestArchiveBuild(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "base.colpkg")
	createArchive(t, base)
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	factory := ankigrpc.NewFactory(ankigrpc.SubprocessAddress(binary, "serve"))
	factory.RegisterClientFactory("subprocess", ankigrpc.ToClientFactory(ankigrpc.NewSubprocessClient))
	application := ankiascode.New(ankiascode.Options{NewClient: factory.NewClient})
	wanted := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{
		Decks: map[string]*ankiascodepb.DeckState{"decks/Default/deck.anki.toml": {Id: 1}},
	}}
	const key = "archive-build-test"
	var staged string
	apply := func(ctx context.Context, path string) error {
		staged = path
		if filepath.Dir(path) == root || path == base {
			return fmt.Errorf("apply did not receive a private archive")
		}
		if err := application.ColStateReplace(ctx, ankiascodeapi.CollectionOptions{StateKey: key, Path: path}, wanted); err != nil {
			return fmt.Errorf("update staged native state: %w", err)
		}
		return nil
	}
	output := filepath.Join(root, "built.colpkg")
	if err := rebuild.New(apply, filesystem.FS{}).Build(t.Context(), base, output); err != nil {
		t.Fatal(err)
	}
	stored, err := application.ColStateDump(t.Context(), ankiascodeapi.CollectionOptions{StateKey: key, Path: output})
	if err != nil || !proto.Equal(stored, wanted) {
		t.Fatalf("published archive lost native state: %v %v", stored, err)
	}
	initial, err := application.ColStateDump(t.Context(), ankiascodeapi.CollectionOptions{StateKey: key, Path: base})
	if err != nil || len(initial.GetResources().GetDecks()) != 0 {
		t.Fatalf("build changed source state: %v %v", initial, err)
	}
	original := readArchive(t, base)
	assertRemoved(t, filepath.Dir(staged))
	unchanged := filepath.Join(root, "unchanged.colpkg")
	if err := rebuild.New(func(ctx context.Context, path string) error {
		if _, err := application.ColStateDump(ctx, ankiascodeapi.CollectionOptions{StateKey: key, Path: path}); err != nil {
			return fmt.Errorf("inspect unchanged stage: %w", err)
		}
		return nil
	}, filesystem.FS{}).Build(t.Context(), base, unchanged); err != nil {
		t.Fatal(err)
	}
	unchangedState, err := application.ColStateDump(t.Context(), ankiascodeapi.CollectionOptions{StateKey: key, Path: unchanged})
	if err != nil || !proto.Equal(initial, unchangedState) {
		t.Fatalf("no-op build changed native state: %v", err)
	}

	for _, scenario := range []string{"apply-error", "cancel", "publish-error", "panic"} {
		t.Run(scenario, func(t *testing.T) {
			ctx, cancel := context.WithCancelCause(t.Context())
			defer cancel(nil)
			failure := errors.New("apply failed after native publication")
			before := readArchive(t, output)
			destination := output
			if scenario == "publish-error" {
				destination = t.TempDir()
			}
			builder := rebuild.New(func(ctx context.Context, path string) error {
				if err := apply(ctx, path); err != nil {
					return fmt.Errorf("prepare failure stage: %w", err)
				}
				switch scenario {
				case "apply-error":
					return failure
				case "cancel":
					cancel(failure)
				case "panic":
					panic(failure)
				}
				return nil
			}, filesystem.FS{})
			var recovered any
			func() {
				defer func() { recovered = recover() }()
				err = builder.Build(ctx, base, destination)
			}()
			if scenario == "panic" {
				if recovered != failure {
					t.Fatalf("lost callback panic: %v", recovered)
				}
			} else {
				if err == nil {
					t.Fatal("failed build succeeded")
				}
				if (scenario == "apply-error" || scenario == "cancel") && !errors.Is(err, failure) {
					t.Fatalf("lost original failure: %v", err)
				}
				if scenario == "cancel" && !errors.Is(err, context.Canceled) {
					t.Fatalf("lost context cancellation: %v", err)
				}
			}
			if !bytes.Equal(before, readArchive(t, output)) {
				t.Fatal("failed build replaced destination")
			}
			assertRemoved(t, filepath.Dir(staged))
		})
	}

	missing := filepath.Join(root, "absent.colpkg")
	if err := rebuild.New(func(ctx context.Context, path string) error {
		if err := apply(ctx, path); err != nil {
			return fmt.Errorf("update unpublished stage: %w", err)
		}
		return context.Canceled
	}, filesystem.FS{}).Build(t.Context(), base, missing); !errors.Is(err, context.Canceled) {
		t.Fatalf("absent-output failure lost its cause: %v", err)
	}
	assertRemoved(t, missing)
	assertRemoved(t, filepath.Dir(staged))

	alias := filepath.Join(root, "alias.colpkg")
	link := filepath.Join(root, "link.colpkg")
	fifo := filepath.Join(root, "fifo.colpkg")
	if err := os.Link(base, alias); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(base, link); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	called := false
	reject := rebuild.New(func(context.Context, string) error {
		called = true
		return nil
	}, filesystem.FS{})
	canceled, cancel := context.WithCancelCause(t.Context())
	cause := errors.New("canceled before starting archive build")
	cancel(cause)
	if err := reject.Build(canceled, base, missing); !errors.Is(err, context.Canceled) || !errors.Is(err, cause) || called {
		t.Fatalf("initial cancellation lost its cause or called apply: %v", err)
	}
	for _, paths := range [][2]string{
		{base, base},
		{base, alias},
		{base, link},
		{fifo, missing},
		{root, missing},
		{missing, output},
		{base, filepath.Join(root, "missing", "output.colpkg")},
	} {
		if err := reject.Build(t.Context(), paths[0], paths[1]); err == nil || called {
			t.Fatalf("invalid paths reached apply: %v: %v", paths, err)
		}
	}
	if !bytes.Equal(original, readArchive(t, base)) {
		t.Fatal("build workflow changed base archive bytes")
	}
	for name, path := range map[string]string{"base.colpkg": base, "built.colpkg": output, "unchanged.colpkg": unchanged} {
		writeArtifact(t, name, readArchive(t, path))
	}
	stateJSON, err := protojson.Marshal(stored)
	if err != nil {
		t.Fatal(err)
	}
	writeArtifact(t, "built-state.json", stateJSON)
}

// createArchive creates a native collection package through explicit collection RPCs.
func createArchive(t *testing.T, path string) {
	t.Helper()
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: binary, Arguments: []string{"serve"}})
	var collections collection.BackendCollectionServiceClient
	opened := false
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if opened {
			if _, err := collections.CloseCollection(ctx, &collection.CloseCollectionRequest{}); err != nil {
				t.Error(err)
			}
		}
		if err := worker.Stop(ctx); err != nil {
			t.Error(err)
		}
	})
	if err := worker.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	mediaDirectory := filepath.Join(directory, "collection.media")
	if err := os.Mkdir(mediaDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	connection, err := worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	collections = collection.NewBackendCollectionServiceClient(connection)
	if _, err := collections.OpenCollection(t.Context(), &collection.OpenCollectionRequest{
		CollectionPath: filepath.Join(directory, "collection.anki2"), MediaFolderPath: mediaDirectory, MediaDbPath: mediaDirectory + ".db2",
	}); err != nil {
		t.Fatal(err)
	}
	opened = true
	if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(t.Context(), &import_export.ExportCollectionPackageRequest{
		OutPath: path, IncludeMedia: true,
	}); err != nil {
		t.Fatal(err)
	}
	opened = false
}

// readArchive reads fixture bytes for source and destination preservation checks.
func readArchive(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// assertRemoved verifies that failed outputs and private staging no longer exist.
func assertRemoved(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("temporary or failed output remains: %s: %v", path, err)
	}
}

// writeArtifact preserves synthetic native packages and decoded state for inspection.
func writeArtifact(t *testing.T, name string, data []byte) {
	t.Helper()
	if root := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); root != "" {
		if err := os.WriteFile(filepath.Join(root, name), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
