package collection

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/textstore"
)

// TestLargeCollectionWorkflow verifies export and convergence for a large note inventory.
func TestLargeCollectionWorkflow(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root)
	source := filepath.Join(root, "fixture.sqlite")
	const databaseSize = 513 << 20
	if err := os.Truncate(source, databaseSize); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "text")
	if err := Export(ctx, source, dir); err != nil {
		t.Fatal(err)
	}
	note, err := textstore.ReadNote(textstore.NotePath(dir, "grammar", 100))
	if err != nil || note.Fields["hanzi"] != "<b>你好</b>" {
		t.Fatalf("large SQLite export failed: %v", err)
	}
	largeArchive := filepath.Join(root, "large.colpkg")
	archiveWithDatabase(t, base, source, largeArchive)
	archiveText := filepath.Join(root, "archive-text")
	if err := Export(ctx, largeArchive, archiveText); err != nil {
		t.Fatal(err)
	}
	rebuilt := filepath.Join(root, "large-rebuilt.colpkg")
	if err := Build(ctx, largeArchive, archiveText, rebuilt); err != nil {
		t.Fatal(err)
	}
	archivePlan, err := Plan(ctx, rebuilt, archiveText)
	if err != nil || len(archivePlan.Changes) != 0 {
		t.Fatalf("large archive did not converge: %+v, %v", archivePlan, err)
	}
	t.Run("concurrent-archive-update", func(t *testing.T) {
		// The large database keeps recompression in progress long enough to
		// observe its output file and simulate a writer after the initial check.
		stop := make(chan struct{})
		updated := make(chan error, 1)
		go func() {
			ticker := time.NewTicker(time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-stop:
					updated <- fmt.Errorf("archive output was not observed")
					return
				case <-ticker.C:
					paths, err := filepath.Glob(filepath.Join(root, ".anki-output-*"))
					if err != nil {
						updated <- err
						return
					}
					if len(paths) == 0 {
						continue
					}
					f, err := os.OpenFile(largeArchive, os.O_WRONLY|os.O_APPEND, 0o600)
					if err == nil {
						_, err = f.WriteString("concurrent-update")
						err = errors.Join(err, f.Close())
					}
					updated <- err
					return
				}
			}
		}()
		err := Apply(ctx, largeArchive, archiveText)
		close(stop)
		if updateErr := <-updated; updateErr != nil {
			t.Fatal(updateErr)
		}
		if err == nil || !strings.Contains(err.Error(), "collection changed during apply") {
			t.Fatalf("late archive update was not refused: %v", err)
		}
		content, err := os.ReadFile(largeArchive)
		if err != nil || !strings.HasSuffix(string(content), "concurrent-update") {
			t.Fatal("concurrent archive update was overwritten")
		}
	})
	if err := Apply(ctx, source, dir); err != nil {
		t.Fatal(err)
	}
	plan, err := Plan(ctx, source, dir)
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("large collection did not converge: %+v, %v", plan, err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "large-database.txt"), []byte("513 MiB SQLite and Zstandard archive exported and reconciled; empty plans\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestLargeNoteFieldWorkflow verifies large field content survives export and reconciliation.
func TestLargeNoteFieldWorkflow(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root)
	dir := filepath.Join(root, "text")
	if err := Export(ctx, base, dir); err != nil {
		t.Fatal(err)
	}
	path := textstore.NotePath(dir, "grammar", 100)
	note, err := textstore.ReadNote(path)
	if err != nil {
		t.Fatal(err)
	}
	const fieldSize = 17 << 20
	note.Fields["hanzi"] = strings.Repeat("x", fieldSize)
	if err := textstore.WriteNote(path, note); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(root, "large-field.colpkg")
	if err := Build(ctx, base, dir, output); err != nil {
		t.Fatal(err)
	}
	plan, err := Plan(ctx, output, dir)
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("large note field did not converge: %+v, %v", plan, err)
	}
	archive, err := openPackage(ctx, output, root)
	if err != nil {
		t.Fatal(err)
	}
	defer archive.Close()
	var fields string
	if err := archive.db.QueryRow("select flds from notes where id=100").Scan(&fields); err != nil {
		t.Fatal(err)
	}
	if len(strings.Split(fields, "\x1f")[0]) != fieldSize {
		t.Fatal("large field content was not preserved")
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := os.WriteFile(filepath.Join(artifacts, "large-field.txt"), []byte("17 MiB note field rebuilt without truncation; empty plan\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
