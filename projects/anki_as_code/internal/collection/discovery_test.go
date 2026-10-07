package collection

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/textstore"
)

// TestConfigDrivenCollectionWorkflow verifies configured layouts and explicit resource identities.
func TestConfigDrivenCollectionWorkflow(t *testing.T) {
	root := t.TempDir()
	base := fixture(t, root)
	text := filepath.Join(root, "text")
	if err := Export(context.Background(), base, text); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(text, "collection.anki.toml")
	content, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	content = []byte(strings.ReplaceAll(string(content), `decks_path = "decks"`, `decks_path = "content"`))
	if err := os.WriteFile(config, content, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(text, "decks"), filepath.Join(text, "content")); err != nil {
		t.Fatal(err)
	}
	markerDir := filepath.Join(text, "content", "arbitrary-folder")
	if err := os.Rename(filepath.Join(text, "content", "grammar"), markerDir); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(markerDir, "deck.anki.toml")
	declaration := model.Deck{Deck: &collectionpb.Deck{}}
	if err := textstore.ReadTOML(marker, declaration.Deck); err != nil {
		t.Fatal(err)
	}
	if declaration.Title != "grammar" {
		t.Fatal("deck title is not explicit")
	}
	declaration.NotesPath = "../../shared-notes"
	if err := textstore.WriteTOML(marker, declaration.Deck); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(markerDir, "notes"), filepath.Join(text, "shared-notes")); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(text, "shared-notes", "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(text, "shared-notes", "100.note.anki.toml"), filepath.Join(nested, "greeting.note.anki.toml")); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{filepath.Join(markerDir, "README.md"), filepath.Join(nested, "unrelated.toml")} {
		if err := os.WriteFile(path, []byte("not an Anki declaration"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	rebuilt := filepath.Join(root, "rebuilt.colpkg")
	if err := Build(context.Background(), base, config, rebuilt); err != nil {
		t.Fatal(err)
	}
	plan, err := Plan(context.Background(), rebuilt, config)
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("rearranged collection did not converge: %+v, %v", plan, err)
	}
	assertRows(t, base, rebuilt, root, "SELECT id,name,kind FROM decks ORDER BY id", true)
	assertRows(t, base, rebuilt, root, "SELECT * FROM cards ORDER BY id", true)
	assertRows(t, base, rebuilt, root, "SELECT * FROM config ORDER BY key", true)
	// Missing configured deck roots must fail rather than plan wholesale deletion.
	if err := os.Rename(filepath.Join(text, "content"), filepath.Join(text, "missing")); err != nil {
		t.Fatal(err)
	}
	if _, err := Plan(context.Background(), rebuilt, config); err == nil {
		t.Fatal("missing decks root accepted")
	}
	if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
		if err := os.WriteFile(filepath.Join(artifact, "config-driven-layout.txt"), []byte("Explicit titles, recursive marker/note discovery, arbitrary relative paths, unrelated files, exact card preservation, and convergence verified.\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestCardAppearanceWorkflow verifies CSS and template edits preserve unrelated note-type settings.
func TestCardAppearanceWorkflow(t *testing.T) {
	root := t.TempDir()
	base := fixture(t, root)
	text := filepath.Join(root, "text")
	if err := Export(context.Background(), base, text); err != nil {
		t.Fatal(err)
	}
	settingsRoot := filepath.Join(text, "note_types")
	rootConfig, err := os.ReadFile(filepath.Join(text, textstore.CollectionMarker))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(rootConfig), `note_types_path = "note_types"`) || strings.Contains(string(rootConfig), "settings_path") {
		t.Fatal("root must reference note types without exporting collection settings")
	}
	for _, path := range []string{filepath.Join(text, "settings.anki.toml"), filepath.Join(text, "settings")} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatal("collection settings were exported")
		}
	}
	settingsPath := filepath.Join(settingsRoot, "Chinese.notetype.anki.toml")
	appearance := &collectionpb.NoteTypeAppearance{}
	if err := textstore.ReadTOML(settingsPath, appearance); err != nil {
		t.Fatal(err)
	}
	appearance.Css = ".card {\n color: red;\n}"
	appearance.Templates[0].Front = "<div>{{hanzi}}</div>"
	if err := textstore.WriteAppearance(settingsPath, appearance); err != nil {
		t.Fatal(err)
	}
	plan, err := Plan(context.Background(), base, text)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, change := range plan.Changes {
		if change.Resource == "appearance" {
			found = true
		}
	}
	if !found {
		t.Fatal("appearance edit missing from plan")
	}
	rebuilt := filepath.Join(root, "styled.colpkg")
	if err := Build(context.Background(), base, text, rebuilt); err != nil {
		t.Fatal(err)
	}
	result, err := openPackage(context.Background(), rebuilt, root)
	if err != nil {
		t.Fatal(err)
	}
	defer result.Close()
	var css, template []byte
	if err := result.db.QueryRow("SELECT config FROM notetypes WHERE id=10").Scan(&css); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(css), ".card {\n color: red;\n}") {
		t.Fatal("CSS not applied exactly")
	}
	if err := result.db.QueryRow("SELECT config FROM templates WHERE ntid=10 AND ord=0").Scan(&template); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(template), "<div>{{hanzi}}</div>") {
		t.Fatal("front template not applied")
	}
	for _, encoded := range [][]byte{css, template} {
		if !bytes.Contains(encoded, []byte{0xa0, 0x06, 0x07}) {
			t.Fatal("unknown configuration lost")
		}
	}
	assertRows(t, base, rebuilt, root, "SELECT * FROM cards ORDER BY id", true)
	assertRows(t, base, rebuilt, root, "SELECT * FROM config ORDER BY key", true)
	plan, err = Plan(context.Background(), rebuilt, text)
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("appearance edits did not converge: %+v, %v", plan, err)
	}
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatal(err)
	}
	duplicate := filepath.Join(settingsRoot, "duplicate.notetype.anki.toml")
	if err := os.WriteFile(duplicate, content, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Plan(context.Background(), rebuilt, text); err == nil {
		t.Fatal("duplicate note type declaration accepted")
	}
	if err := os.Remove(duplicate); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(settingsPath); err != nil {
		t.Fatal(err)
	}
	if _, err := Plan(context.Background(), rebuilt, text); err == nil {
		t.Fatal("missing note type declaration accepted")
	}
	nested := filepath.Join(settingsRoot, "appearance", "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "renamed.notetype.anki.toml"), content, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(settingsRoot, filepath.Join(text, "configuration")); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(text, textstore.CollectionMarker)
	configBytes, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	configBytes = []byte(strings.ReplaceAll(string(configBytes), `note_types_path = "note_types"`, `note_types_path = "configuration"`))
	if err := os.WriteFile(config, configBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	plan, err = Plan(context.Background(), rebuilt, text)
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("relocated settings did not converge: %+v, %v", plan, err)
	}
	if err := os.Rename(filepath.Join(text, "configuration"), filepath.Join(text, "missing-note-types")); err != nil {
		t.Fatal(err)
	}
	if _, err := Plan(context.Background(), rebuilt, text); err == nil {
		t.Fatal("missing note types directory accepted")
	}
	if artifact := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifact != "" {
		if err := os.WriteFile(filepath.Join(artifact, "card-appearance.txt"), []byte("Per-note-type files without collection settings; CSS/template edits; unknown protobuf and scheduling preservation; relocated/nested settings; missing/duplicate rejection; convergence verified.\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestExternalNoteIdentityWorkflow verifies identity reservation for externally located notes.
func TestExternalNoteIdentityWorkflow(t *testing.T) {
	root := t.TempDir()
	base := fixture(t, root)
	text := filepath.Join(root, "text")
	if err := Export(context.Background(), base, text); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(text, "decks", "grammar", textstore.DeckMarker)
	deckBytes, err := os.ReadFile(marker)
	if err != nil {
		t.Fatal(err)
	}
	external := filepath.Join(root, "external-notes")
	if err := os.Rename(filepath.Join(text, "decks", "grammar", "notes"), external); err != nil {
		t.Fatal(err)
	}
	deckBytes = []byte(strings.ReplaceAll(string(deckBytes), `notes_path = 'notes'`, `notes_path = "`+filepath.ToSlash(external)+`"`))
	if err := os.WriteFile(marker, deckBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	original := filepath.Join(external, "100.note.anki.toml")
	note, err := textstore.ReadNote(original)
	if err != nil {
		t.Fatal(err)
	}
	// Reserve an explicitly declared future card ID outside the config directory.
	note.Cards[0].Id = 9000000000000000
	if err := textstore.WriteNote(original, note); err != nil {
		t.Fatal(err)
	}
	copied := filepath.Join(external, "copy.note.anki.toml")
	if err := textstore.WriteNote(copied, note); err != nil {
		t.Fatal(err)
	}
	count, err := textstore.GenerateIDsWithConfig(context.Background(), []string{copied}, filepath.Join(text, textstore.CollectionMarker))
	if err != nil || count != 1 {
		t.Fatalf("external identity generation: %d, %v", count, err)
	}
	result, err := textstore.ReadNote(copied)
	if err != nil || result.NoteId <= 9000000000000000 {
		t.Fatalf("external collection identities not reserved: %+v, %v", result, err)
	}
	originalNote, err := textstore.ReadNote(original)
	if err != nil || originalNote.NoteId != 100 || originalNote.Cards[0].Id != 9000000000000000 {
		t.Fatal("unselected external note changed")
	}
}

// TestTemplateOnlySyncWorkflow verifies appearance-only reconciliation and synchronization metadata.
func TestTemplateOnlySyncWorkflow(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := fixture(t, root)
	dir := filepath.Join(root, "text")
	if err := Export(ctx, base, dir); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "note_types", "Chinese.notetype.anki.toml")
	appearance := &collectionpb.NoteTypeAppearance{}
	if err := textstore.ReadTOML(path, appearance); err != nil {
		t.Fatal(err)
	}
	appearance.Templates[0].Front = "<div>{{hanzi}}</div>"
	if err := textstore.WriteAppearance(path, appearance); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(root, "template-only.colpkg")
	if err := Build(ctx, base, dir, output); err != nil {
		t.Fatal(err)
	}
	archive, err := openPackage(ctx, output, root)
	if err != nil {
		t.Fatal(err)
	}
	defer archive.Close()
	var usn, modified int64
	if err := archive.db.QueryRow("SELECT usn,mtime_secs FROM notetypes WHERE id=10").Scan(&usn, &modified); err != nil {
		t.Fatal(err)
	}
	if usn != -1 || modified <= 0 {
		t.Fatalf("template-only edit is not pending sync: usn=%d, modified=%d", usn, modified)
	}
	assertRows(t, base, output, root, "SELECT config FROM notetypes ORDER BY id", true)
	assertRows(t, base, output, root, "SELECT * FROM cards ORDER BY id", true)
	plan, err := Plan(ctx, output, dir)
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("template-only edit did not converge: %+v, %v", plan, err)
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := os.WriteFile(filepath.Join(artifacts, "template-only-sync.txt"), []byte("Template-only edit marks the parent note type for sync; CSS and scheduling preserved; empty plan.\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
