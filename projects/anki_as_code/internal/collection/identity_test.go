package collection

import collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/textstore"
)

// TestGenerateIdentitiesWorkflow verifies copied notes receive distinct usable identities.
func TestGenerateIdentitiesWorkflow(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root)
	dir := filepath.Join(root, "text")
	if err := Export(ctx, base, dir); err != nil {
		t.Fatal(err)
	}
	original := textstore.NotePath(dir, "grammar", 100)
	contents, err := os.ReadFile(original)
	if err != nil {
		t.Fatal(err)
	}
	var files []string
	for _, deck := range []string{"copies", "copies::child"} {
		if err := textstore.WriteDeck(dir, model.Deck{Deck: &collectionpb.Deck{Title: deck}, Kind: "CgIIAQ=="}); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(textstore.DeckPath(dir, deck), "notes", "copied.note.anki.toml")
		if err := os.WriteFile(path, contents, 0o600); err != nil {
			t.Fatal(err)
		}
		files = append(files, path)
	}
	count, err := textstore.GenerateIDs(context.Background(), []string{files[0], textstore.DeckPath(dir, "copies"), files[0]})
	if err != nil || count != 2 {
		t.Fatalf("generate identities: count=%d, error=%v", count, err)
	}
	ids := map[int64]bool{100: true, 200: true, 201: true}
	guids := map[string]bool{"guid": true}
	for _, path := range files {
		note, err := textstore.ReadNote(path)
		if err != nil {
			t.Fatal(err)
		}
		if note.NoteId <= 0 || ids[note.NoteId] || note.Guid == "" || guids[note.Guid] {
			t.Fatal("generated note identities are not unique")
		}
		ids[note.NoteId], guids[note.Guid] = true, true
		if note.NoteTypeId != 10 || note.Fields["hanzi"] != "<b>你好</b>" || len(note.Cards) != 1 {
			t.Fatal("generation changed note content or card declarations")
		}
		for _, card := range note.Cards {
			if card.Id <= 0 || ids[card.Id] {
				t.Fatal("generated card identities are not unique")
			}
			ids[card.Id] = true
		}
	}
	unchanged, err := os.ReadFile(original)
	if err != nil || string(unchanged) != string(contents) {
		t.Fatal("unselected note was modified")
	}
	plan, err := Plan(ctx, base, dir)
	if err != nil || plan.NotesAdded != 2 || plan.CardsAdded != 2 || plan.NotesDeleted != 0 {
		t.Fatalf("incorrect identity plan: %+v, %v", plan, err)
	}
	output := filepath.Join(root, "identities.colpkg")
	if err := Build(ctx, base, dir, output); err != nil {
		t.Fatal(err)
	}
	plan, err = Plan(ctx, output, dir)
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("generated identities did not converge: %+v, %v", plan, err)
	}
	assertRows(t, base, output, root, "select * from cards where nid=100 order by id", true)
	if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
		if err := os.WriteFile(filepath.Join(artifact, "generated-identities.txt"), []byte("two copied notes, four new cards; nested and overlapping paths; original scheduling preserved; rebuilt archive converged\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestGenerateIdentitiesRefusesInvalidBatch verifies invalid batches fail before note files are changed.
func TestGenerateIdentitiesRefusesInvalidBatch(t *testing.T) {
	root := t.TempDir()
	valid := filepath.Join(root, "valid.note.anki.toml")
	contents := "note_id = 1\nguid = 'original'\nnote_type_id = 10\nnote_type = 'Chinese'\n[fields]\nhanzi = 'unchanged'\n"
	if err := os.WriteFile(valid, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	invalid := filepath.Join(root, "invalid.note.anki.toml")
	if err := os.WriteFile(invalid, []byte("invalid TOML ["), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := textstore.GenerateIDs(context.Background(), []string{valid, invalid}); err == nil {
		t.Fatal("malformed batch accepted")
	}
	got, err := os.ReadFile(valid)
	if err != nil || string(got) != contents {
		t.Fatal("invalid batch changed an earlier note")
	}
	if _, err := textstore.GenerateIDs(context.Background(), []string{filepath.Join(root, "missing")}); err == nil {
		t.Fatal("missing path accepted")
	}
	link := filepath.Join(root, "link.note.anki.toml")
	if err := os.Symlink(valid, link); err != nil {
		t.Fatal(err)
	}
	if _, err := textstore.GenerateIDs(context.Background(), []string{link}); err == nil {
		t.Fatal("symlink accepted")
	}
}
