package backend_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/ankitects/anki/proto/anki/notes"
)

// TestCollectionWriteBack verifies opt-in persistence after RPC failures and default source preservation.
func TestCollectionWriteBack(t *testing.T) {
	worker := backend.New(backend.BackendSTDOptions{Binary: workerPath(t)})
	if err := worker.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := worker.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	})
	root := t.TempDir()
	archive := filepath.Join(root, "source.colpkg")
	managed := backend.NewCollection(backend.BackendCollectionOptions{Backend: worker, Path: filepath.Join(root, "source.anki2")})
	if err := managed.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	connection, err := managed.Connection()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(t.Context(), &import_export.ExportCollectionPackageRequest{OutPath: archive, IncludeMedia: true}); err != nil {
		t.Fatal(err)
	}
	if err := managed.Stop(t.Context()); err != nil {
		t.Fatal(err)
	}
	for _, policy := range []string{"default", "disabled", "enabled"} {
		disabled := policy != "enabled"
		before, err := os.ReadFile(archive)
		if err != nil {
			t.Fatal(err)
		}
		options := backend.BackendCollectionOptions{Backend: worker, Path: archive}
		if policy != "default" {
			options.WriteBack = !disabled
		}
		managed = backend.NewCollection(options)
		if err := managed.Start(t.Context()); err != nil {
			t.Fatal(err)
		}
		connection, err := managed.Connection()
		if err != nil {
			t.Fatal(err)
		}
		client := decks.NewDecksServiceClient(connection)
		deck, err := client.NewDeck(t.Context(), &generic.Empty{})
		if err != nil {
			t.Fatal(err)
		}
		deck.Name = "Write-back verification"
		if _, err := client.AddDeck(t.Context(), deck); err != nil {
			t.Fatal(err)
		}
		if _, err := notes.NewNotesServiceClient(connection).GetNote(t.Context(), &notes.NoteId{Nid: -1}); err == nil {
			t.Fatal("invalid note unexpectedly exists")
		}
		if err := managed.Stop(t.Context()); err != nil {
			t.Fatal(err)
		}
		if disabled {
			after, err := os.ReadFile(archive)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("disabled write-back changed source: %v", err)
			}
		}
		inspection := backend.NewCollection(backend.BackendCollectionOptions{Backend: worker, Path: archive})
		if err := inspection.Start(t.Context()); err != nil {
			t.Fatal(err)
		}
		connection, err = inspection.Connection()
		if err != nil {
			t.Fatal(err)
		}
		_, lookupErr := decks.NewDecksServiceClient(connection).GetDeckIdByName(t.Context(), &generic.String{Val: deck.Name})
		if disabled && lookupErr == nil || !disabled && lookupErr != nil {
			t.Fatalf("write-back disabled=%v, lookup=%v", disabled, lookupErr)
		}
		if err := inspection.Stop(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	managed = backend.NewCollection(backend.BackendCollectionOptions{Backend: worker, Path: archive, WriteBack: true})
	if err := managed.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	connection, err = managed.Connection()
	if err != nil {
		t.Fatal(err)
	}
	restore := collectionRestoreDirectory(t, connection)
	if err := os.Rename(archive, archive+".original"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(archive, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := managed.Stop(t.Context()); err == nil {
		t.Fatal("write-back to a directory unexpectedly succeeded")
	}
	if _, err := os.Stat(filepath.Join(restore, "collection.anki2")); err != nil {
		t.Fatalf("failed publication discarded recoverable storage: %v", err)
	}
	writeReport(t, "collection-write-back.json", map[string]any{"default_enabled": false, "explicit_enable_persists": true, "request_failure_preserves_writes": true, "disabled_preserves_archive": true, "failed_publication_retains_restore": true})
}
