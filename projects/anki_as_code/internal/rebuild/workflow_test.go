package rebuild_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/rebuild"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/pkg/ankiascode"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// TestArchiveBuild verifies staged publication using real Anki archives and native state writes.
func TestArchiveBuild(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "base.colpkg")
	createArchive(t, base)
	original := readArchive(t, base)
	application := &ankiascode.App{}
	wanted := &collectionpb.State{Version: 1, Resources: &collectionpb.Resources{
		Decks: map[string]*collectionpb.DeckState{"decks/Default/deck.anki.toml": {Id: 1}},
	}}
	const key = "archive-build-test"
	var staged string
	apply := func(ctx context.Context, path string) error {
		staged = path
		if filepath.Dir(path) == root || path == base {
			return fmt.Errorf("apply did not receive a private archive")
		}
		if err := application.ColStateReplace(ctx, path, key, wanted); err != nil {
			return fmt.Errorf("update staged native state: %w", err)
		}
		return nil
	}
	output := filepath.Join(root, "built.colpkg")
	if err := rebuild.New(apply).Build(t.Context(), base, output); err != nil {
		t.Fatal(err)
	}
	stored, err := application.ColStateDump(t.Context(), output, key)
	if err != nil || !proto.Equal(stored, wanted) {
		t.Fatalf("published archive lost native state: %v %v", stored, err)
	}
	initial, err := application.ColStateDump(t.Context(), base, key)
	if err != nil || len(initial.GetResources().GetDecks()) != 0 {
		t.Fatalf("build changed source state: %v %v", initial, err)
	}
	assertRemoved(t, filepath.Dir(staged))
	unchanged := filepath.Join(root, "unchanged.colpkg")
	if err := rebuild.New(func(ctx context.Context, path string) error {
		if _, err := application.ColStateDump(ctx, path, key); err != nil {
			return fmt.Errorf("inspect unchanged stage: %w", err)
		}
		return nil
	}).Build(t.Context(), base, unchanged); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(original, readArchive(t, unchanged)) {
		t.Fatal("no-op build changed archive bytes")
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
			})
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
	}).Build(t.Context(), base, missing); !errors.Is(err, context.Canceled) {
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
	})
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

// createArchive creates a native collection package through the managed backend APIs.
func createArchive(t *testing.T, path string) {
	t.Helper()
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	worker := backend.New(backend.BackendSTDOptions{Binary: binary})
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := worker.Stop(ctx); err != nil {
			t.Error(err)
		}
	})
	if err := worker.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	if err := os.Mkdir(filepath.Join(directory, "collection.media"), 0o700); err != nil {
		t.Fatal(err)
	}
	collection := backend.NewCollection(backend.BackendCollectionOptions{
		Backend: worker, Path: filepath.Join(directory, "collection.anki2"),
	})
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := collection.Stop(ctx); err != nil {
			t.Error(err)
		}
	})
	if err := collection.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	connection, err := collection.Connection()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(t.Context(), &import_export.ExportCollectionPackageRequest{
		OutPath: path, IncludeMedia: true,
	}); err != nil {
		t.Fatal(err)
	}
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
