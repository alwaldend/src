package collection

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/textstore"
)

// TestSchemaPreservationWorkflow verifies unmanaged upstream configuration survives reconciliation.
func TestSchemaPreservationWorkflow(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	original := fixture(t, root)
	database := filepath.Join(root, "fixture.sqlite")
	db, err := openDB(database)
	if err != nil {
		t.Fatal(err)
	}
	// Captured wire values include explicit defaults, reordered fields, optional
	// zero limits, nested unmanaged settings, and an unknown nonminimal varint.
	noteConfig := []byte{0x2a, 6, 'p', 'r', 'e', 'f', 'i', 'x', 0xa0, 6, 0x87, 0}
	templateConfig := []byte{0x28, 9, 0xa0, 6, 0x87, 0}
	deckKind := []byte{0x0a, 17, 0x22, 0, 8, 1, 0xa0, 6, 7, 0x30, 0, 0x42, 4, 8, 1, 16, 2, 0x28, 0, 0xa8, 6, 9}
	for _, update := range []struct {
		// query updates one stored configuration in the fixture database.
		query string
		// value supplies the captured upstream protobuf payload.
		value []byte
	}{
		{"UPDATE notetypes SET config=? WHERE id=10", noteConfig},
		{"UPDATE templates SET config=? WHERE ntid=10 AND ord=0", templateConfig},
		{"UPDATE decks SET kind=? WHERE id=1", deckKind},
	} {
		if _, err := db.Exec(update.query, update.value); err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(root, "schema.colpkg")
	archiveWithDatabase(t, original, database, base)
	textRoot := filepath.Join(root, "text")
	if err := Export(ctx, base, textRoot); err != nil {
		t.Fatal(err)
	}
	plan, err := Plan(ctx, base, textRoot)
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("unmodified schema export must converge: %+v, %v", plan, err)
	}
	unchanged := filepath.Join(root, "unchanged.colpkg")
	if err := Build(ctx, base, textRoot, unchanged); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"decks", "notetypes", "templates"} {
		assertRows(t, base, unchanged, root, "SELECT * FROM "+table+" ORDER BY 1", true)
	}
	marker := filepath.Join(textstore.DeckPath(textRoot, "grammar"), textstore.DeckMarker)
	deck := model.Deck{Deck: &collectionpb.Deck{}}
	if err := textstore.ReadTOML(marker, deck.Deck); err != nil {
		t.Fatal(err)
	}
	deck.Normal.ExtendNew = 25
	if err := textstore.WriteTOML(marker, deck.Deck); err != nil {
		t.Fatal(err)
	}
	appearancePath := filepath.Join(textRoot, "note_types", "Chinese"+textstore.NoteTypeSuffix)
	appearance := &collectionpb.NoteTypeAppearance{}
	if err := textstore.ReadTOML(appearancePath, appearance); err != nil {
		t.Fatal(err)
	}
	appearance.Css = ".card { color: blue; }"
	appearance.Templates[0].Front = "{{hanzi}}"
	if err := textstore.WriteTOML(appearancePath, appearance); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(root, "edited.colpkg")
	if err := Build(ctx, base, textRoot, output); err != nil {
		t.Fatal(err)
	}
	result, err := openPackage(ctx, output, root)
	if err != nil {
		t.Fatal(err)
	}
	defer result.Close()
	for _, check := range []struct {
		// query selects a configuration payload after reconciliation.
		query string
		// parts lists byte sequences that must survive the update.
		parts [][]byte
	}{
		{"SELECT config FROM notetypes WHERE id=10", [][]byte{noteConfig[:8], noteConfig[8:]}},
		{"SELECT config FROM templates WHERE ntid=10 AND ord=0", [][]byte{templateConfig[:2], templateConfig[2:]}},
		{"SELECT kind FROM decks WHERE id=1", [][]byte{{0x30, 0}, {0x42, 4, 8, 1, 16, 2}, {0xa0, 6, 7}, {0xa8, 6, 9}}},
	} {
		var actual []byte
		if err := result.db.QueryRow(check.query).Scan(&actual); err != nil {
			t.Fatal(err)
		}
		for _, part := range check.parts {
			if !bytes.Contains(actual, part) {
				t.Fatalf("unmanaged setting %x lost from %x", part, actual)
			}
		}
	}
	for _, table := range []string{"cards", "config", "revlog"} {
		assertRows(t, base, output, root, "SELECT * FROM "+table+" ORDER BY 1", true)
	}
	plan, err = Plan(ctx, output, textRoot)
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("schema edits must converge: %+v, %v", plan, err)
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := os.WriteFile(filepath.Join(artifacts, "schema-preservation.txt"), []byte("Upstream-schema reconciliation preserves untouched blobs, optional zero limits, nested unmanaged deck settings, unknown wire bytes, scheduling, history, and convergence.\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
