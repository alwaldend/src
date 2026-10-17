package ankiascode

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"github.com/ankitects/anki/proto/anki/decks"
	"google.golang.org/protobuf/proto"
)

// TestExportConfigOnlyWorkflow verifies config-only output, read-only export and independent state dumps.
func TestExportConfigOnlyWorkflow(t *testing.T) {
	for _, inputFormat := range []string{"package", "database"} {
		t.Run(inputFormat, func(t *testing.T) {
			fixture := newFixture(t)
			input := fixture.base
			if inputFormat == "database" {
				input = fixture.database
			}
			before, err := os.ReadFile(input)
			if err != nil {
				t.Fatal(err)
			}
			text := filepath.Join(t.TempDir(), "text")
			runFixtureCLI(t, "col", "export", "--collection", input, "--output", text)
			after, err := os.ReadFile(input)
			if err != nil || string(before) != string(after) {
				t.Fatal("export modified source collection bytes")
			}
			initial := stateFixtureDump(t, input)
			if len(initial.GetResources().GetDecks())+len(initial.GetResources().GetNotes())+len(initial.GetResources().GetNoteTypes()) != 0 {
				t.Fatal("export initialized source state")
			}
			for _, format := range []string{"toml", "json"} {
				file := filepath.Join(text, "state.anki."+format)
				if _, err := os.Stat(file); !os.IsNotExist(err) {
					t.Fatalf("export generated %s state: %v", format, err)
				}
				if err := os.WriteFile(file, []byte("caller state"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			runFixtureCLI(t, "col", "export", "--collection", input, "--output", text)
			for _, format := range []string{"toml", "json"} {
				if data, err := os.ReadFile(filepath.Join(text, "state.anki."+format)); err != nil || string(data) != "caller state" {
					t.Fatalf("export changed caller-owned %s state: %v", format, err)
				}
			}
			runFixtureCLI(t, "col", "state", "import", "--collection", input, "--text", text,
				"decks/grammar/deck.anki.toml=1",
				fmt.Sprintf("note_types/Chinese.notetype.anki.toml=%d", fixture.typeID),
				fmt.Sprintf("decks/grammar/notes/你好.note.anki.toml=%d", fixture.noteID))
			stored := stateFixtureDump(t, input)
			if stored.GetResources().GetNotes()["decks/grammar/notes/你好.note.anki.toml"].GetId() != fixture.noteID {
				t.Fatal("separate state setup did not bind the native note")
			}
			for _, format := range []string{"toml", "json"} {
				data := runFixtureCLI(t, "col", "state", "dump", "--format", format, "--collection", input)
				value := &collectionpb.State{}
				if err := codec.Decode(data, codec.Format(format), value); err != nil || !proto.Equal(value, stored) {
					t.Fatalf("separate %s state dump: %v", format, err)
				}
				file := filepath.Join(t.TempDir(), "state.anki."+format)
				if err := os.WriteFile(file, data, 0o600); err != nil {
					t.Fatal(err)
				}
				saveFixtureArtifact(t, file, "dumped-state-"+inputFormat+"."+format)
			}
			key := "export-fixture/custom"
			custom := proto.Clone(stored).(*collectionpb.State)
			path := "decks/grammar/notes/custom.note.anki.toml"
			custom.Resources.Notes[path] = custom.Resources.Notes["decks/grammar/notes/你好.note.anki.toml"]
			delete(custom.Resources.Notes, "decks/grammar/notes/你好.note.anki.toml")
			stateFile := filepath.Join(t.TempDir(), "custom-state.anki.toml")
			writeFixtureTOML(t, stateFile, custom)
			runFixtureCLI(t, "col", "state", "replace", "--collection", input, "--state-key", key, stateFile)
			beforeExport := fixtureSources(t, text)
			runFixtureCLI(t, "col", "export", "--collection", input, "--output", text)
			assertFixtureSources(t, text, beforeExport)
			selected := filepath.Join(t.TempDir(), "selected")
			runFixtureCLI(t, "col", "export", "--collection", input, "--output", selected, "--state-key", key)
			if _, err := os.Stat(filepath.Join(selected, filepath.FromSlash(path))); err != nil {
				t.Fatal("export ignored the selected state namespace", err)
			}
			if !proto.Equal(custom, stateFixtureDumpKey(t, input, key)) || !proto.Equal(stored, stateFixtureDump(t, input)) {
				t.Fatal("export changed stored bindings")
			}
			configs, err := config.New(selected).Load()
			if err != nil || len(configs.GetDecks()) == 0 || len(configs.GetNoteTypes()) == 0 {
				t.Fatalf("exported hierarchy is not loadable: %v", err)
			}
			saveFixtureArtifact(t, filepath.Join(selected, filepath.FromSlash(path)), "exported-note-"+inputFormat+".toml")
		})
	}
}

// TestExportWritesWorkflow verifies in-place output and caller-owned files/directories.
func TestExportWritesWorkflow(t *testing.T) {
	fixture := newFixture(t)
	root := t.TempDir()
	t.Setenv("HOME", root)
	text := filepath.Join(root, "text")
	reference := filepath.Join(root, "permissions")
	if err := os.Mkdir(reference, 0o755); err != nil {
		t.Fatal(err)
	}
	defaultMode, err := os.Stat(reference)
	if err != nil {
		t.Fatal(err)
	}
	if err := exportForFixture(context.Background(), fixture.base, "~/text", ""); err != nil {
		t.Fatal(err)
	}
	created, err := os.Stat(text)
	if err != nil || created.Mode().Perm() != defaultMode.Mode().Perm() {
		t.Fatalf("new export directory permissions differ from normal creation: %v: %v", created, err)
	}
	if err := os.Chmod(text, 0o775); err != nil {
		t.Fatal(err)
	}
	unrelated := filepath.Join(text, "README.txt")
	if err := os.WriteFile(unrelated, []byte("keep me"), 0o600); err != nil {
		t.Fatal(err)
	}
	private := filepath.Join(text, "private")
	if err := os.Mkdir(private, 0o700); err != nil {
		t.Fatal(err)
	}
	privateFile := filepath.Join(private, "personal.txt")
	if err := os.WriteFile(privateFile, []byte("unrelated private content"), 0o600); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(text, "stale.note.anki.toml")
	if err := os.WriteFile(stale, []byte("stale"), 0o600); err != nil {
		t.Fatal(err)
	}
	staleDeck := filepath.Join(text, "decks", "Removed")
	marker := filepath.Join(staleDeck, "decks", "Child", "deck.anki.toml")
	writeFixtureTOML(t, marker, &collectionpb.DeckConfig{Normal: &decks.Deck_Normal{}})
	writeFixtureTOML(t, filepath.Join(filepath.Dir(marker), "notes", "old.note.anki.toml"), &collectionpb.NoteConfig{})
	staleType := filepath.Join(text, "note_types", "Removed.notetype.anki.toml")
	writeFixtureTOML(t, staleType, &collectionpb.NoteTypeConfig{Name: "Removed"})
	if err := exportForFixture(context.Background(), fixture.base, "~/text", ""); err != nil {
		t.Fatal(err)
	}
	refreshed, err := os.Stat(text)
	if err != nil || refreshed.Mode().Perm() != 0o775 {
		t.Fatalf("export refresh changed directory permissions: %v: %v", refreshed, err)
	}
	if data, err := os.ReadFile(stale); err != nil || string(data) != "stale" {
		t.Fatal("export changed an unselected declaration")
	}
	for _, retained := range []string{marker, staleType} {
		if _, err := os.Stat(retained); err != nil {
			t.Fatalf("export removed an unselected resource at %q: %v", retained, err)
		}
	}
	data, err := os.ReadFile(unrelated)
	if err != nil || string(data) != "keep me" {
		t.Fatal("unrelated file changed")
	}
	for path, mode := range map[string]os.FileMode{private: 0o700, privateFile: 0o600, unrelated: 0o600} {
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != mode {
			t.Errorf("export changed private permissions at %q: %v: %v", path, info, err)
		}
	}
	privateContent, err := os.ReadFile(privateFile)
	if err != nil || string(privateContent) != "unrelated private content" {
		t.Fatal("unrelated private content changed")
	}
	configPath := filepath.Join(text, "collection.anki.toml")
	before, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := exportForFixture(context.Background(), filepath.Join(root, "missing.colpkg"), text, ""); err == nil {
		t.Fatal("invalid input accepted")
	}
	after, err := os.ReadFile(configPath)
	if err != nil || string(before) != string(after) {
		t.Fatal("failed export changed existing declarations")
	}
	fresh := filepath.Join(root, "fresh")
	if err := exportForFixture(context.Background(), fixture.base, fresh, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := config.New(fresh).Load(); err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile(fixture.database)
	if err != nil {
		t.Fatal(err)
	}
	if err := exportForFixture(context.Background(), fixture.database, filepath.Dir(fixture.database), ""); err == nil {
		t.Fatal("export accepted a destination containing its source")
	}
	if after, err := os.ReadFile(fixture.database); err != nil || string(after) != string(original) {
		t.Fatal("rejected export changed its source database")
	}
	for _, scenario := range []string{"source symlink", "destination symlink"} {
		t.Run(scenario, func(t *testing.T) {
			sourceDirectory := filepath.Dir(fixture.database)
			link := filepath.Join(t.TempDir(), "alias")
			if err := os.Symlink(sourceDirectory, link); err != nil {
				t.Fatal(err)
			}
			input, output := fixture.database, sourceDirectory
			if scenario == "source symlink" {
				input = filepath.Join(link, filepath.Base(fixture.database))
			} else {
				output = link
			}
			before := fixtureSources(t, sourceDirectory)
			if err := exportForFixture(context.Background(), input, output, ""); err == nil {
				t.Fatal("export accepted a destination physically containing its source")
			}
			assertFixtureSources(t, sourceDirectory, before)
		})
	}

	saveFixtureArtifact(t, configPath, "exported.anki.toml")
}
