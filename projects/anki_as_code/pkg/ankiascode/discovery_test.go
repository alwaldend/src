package ankiascode

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// TestConfigDrivenCollectionWorkflow verifies alternate root filenames and nested deck-derived titles.
func TestConfigDrivenCollectionWorkflow(t *testing.T) {
	fixture := newFixture(t)
	text := filepath.Join(t.TempDir(), "text")
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(text, "collection.anki.toml")
	alternate := filepath.Join(text, "custom-root.toml")
	if err := os.Rename(source, alternate); err != nil {
		t.Fatal(err)
	}
	assertConverged(t, fixture.base, alternate)
	marker := filepath.Join(text, "decks", "grammar", "decks", "Lesson", "deck.anki.toml")
	writeFixtureTOML(t, marker, &collectionpb.DeckConfig{Normal: &decks.Deck_Normal{}})
	old := fixtureNotePath(text)
	note := &collectionpb.NoteConfig{}
	readFixtureTOML(t, old, note)
	note.Cards = nil
	target := filepath.Join(filepath.Dir(marker), "notes", "renamed.note.anki.toml")
	writeFixtureTOML(t, target, note)
	if err := os.Remove(old); err != nil {
		t.Fatal(err)
	}
	writeFixtureTOML(t, alternate, &collectionpb.CollectionConfig{Version: 1, Moved: []*collectionpb.ResourceMove{{From: "decks/grammar/notes/你好.note.anki.toml", To: "decks/grammar/decks/Lesson/notes/renamed.note.anki.toml"}}})
	before := fixtureSources(t, text)
	plan := fixturePlan(t, fixture.base, alternate)
	if plan.NotesAdded != 0 || plan.NotesDeleted != 0 || plan.CardsChanged != 2 || plan.DecksAdded != 1 {
		t.Fatalf("state/path resolution: %v", plan)
	}
	result := filepath.Join(t.TempDir(), "hierarchy.colpkg")
	if err := buildForFixture(context.Background(), fixture.base, alternate, result, ""); err != nil {
		t.Fatal(err)
	}
	assertConverged(t, result, alternate)
	assertFixtureSources(t, text, before)
	saveFixtureArtifact(t, result, "hierarchy.colpkg")
}

// TestNoteTypeExportNamesWorkflow verifies tracked filenames, native discovery and stable referenced exports.
func TestNoteTypeExportNamesWorkflow(t *testing.T) {
	cases := []struct {
		// name identifies the native collision scenario and its emitted artifacts.
		name string
		// nativeName names the untracked type competing for a tracked filename.
		nativeName string
		// filename is the expected newly allocated filename, except for the dynamic fallback case.
		filename string
	}{
		{name: "renamed", nativeName: "Chinese", filename: "Chinese-2.notetype.anki.toml"},
		{name: "occupied-suffix", nativeName: "Chinese", filename: "Chinese-3.notetype.anki.toml"},
		{name: "earlier-native-id", nativeName: "Basic", filename: "Basic-2.notetype.anki.toml"},
		{name: "case-insensitive", nativeName: "chinese", filename: "chinese-2.notetype.anki.toml"},
		{name: "fallback", nativeName: "Native/Type"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			fixture := newFixture(t)
			ctx := context.Background()
			root := t.TempDir()
			initial := filepath.Join(root, "initial")
			if err := prepareFixtureConfigs(t, ctx, fixture.base, initial); err != nil {
				t.Fatal(err)
			}
			stored := stateFixtureDump(t, fixture.base)
			stateFile := filepath.Join(root, "state.anki.toml")
			types := stored.Resources.NoteTypes
			trackedPath := "note_types/Chinese.notetype.anki.toml"
			tracked := types[trackedPath]
			basicPath := "note_types/Basic.notetype.anki.toml"
			basic := types[basicPath]
			if tracked == nil || basic == nil || basic.Id >= tracked.Id {
				t.Fatal("fixture requires a tracked Chinese type and an earlier Basic type")
			}
			nativeCount := len(types)
			switch test.name {
			case "occupied-suffix":
				delete(types, basicPath)
				types["note_types/Chinese-2.notetype.anki.toml"] = basic
			case "earlier-native-id":
				delete(types, basicPath)
				delete(types, trackedPath)
				trackedPath = basicPath
				types[trackedPath] = tracked
			case "fallback":
				delete(types, trackedPath)
				trackedPath = fmt.Sprintf("note_types/type-%d.notetype.anki.toml", nativeCount+1)
				types[trackedPath] = tracked
			}
			client := inspectFixture(t, fixture.base)
			definition, err := client.types.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: fixture.typeID})
			if err != nil {
				t.Fatal(err)
			}
			definition.Name = "Renamed"
			if _, err := client.types.UpdateNotetype(ctx, definition); err != nil {
				t.Fatal(err)
			}
			newID := basic.Id
			filename := test.filename
			if test.name != "earlier-native-id" {
				added, err := client.types.AddNotetype(ctx, &notetypes.Notetype{
					Name: test.nativeName, Config: &notetypes.Notetype_Config{},
					Fields: []*notetypes.Notetype_Field{
						{Name: "Front", Config: &notetypes.Notetype_Field_Config{}},
						{Name: "Back", Config: &notetypes.Notetype_Field_Config{}},
					},
					Templates: []*notetypes.Notetype_Template{
						{Name: "Card", Config: &notetypes.Notetype_Template_Config{QFormat: "{{Front}}", AFormat: "{{Back}}"}},
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				newID = added.Id
				for _, binding := range types {
					if binding.Id >= newID {
						t.Fatal("new native type must sort after the tracked types")
					}
				}
				if test.name == "fallback" {
					filename = fmt.Sprintf("type-%d-2.notetype.anki.toml", nativeCount+1)
				}
				nativeCount++
			}
			_, err = client.notes.AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{
				Guid: "discovered-type-note", NotetypeId: newID, Fields: []string{"Discovered", "Back"},
			}})
			if err != nil {
				t.Fatal(err)
			}
			client.exportFixture(t, fixture.base)
			if err := client.close(); err != nil {
				t.Fatal(err)
			}
			writeFixtureTOML(t, stateFile, stored)
			runFixtureCLI(t, "col", "state", "replace", "--collection", fixture.base, stateFile)
			before, err := os.ReadFile(fixture.base)
			if err != nil {
				t.Fatal(err)
			}
			text := filepath.Join(root, "exported")
			runFixtureCLI(t, "col", "export", "--collection", fixture.base, "--output", text)
			after, err := os.ReadFile(fixture.base)
			if err != nil || string(before) != string(after) || !proto.Equal(stored, stateFixtureDump(t, fixture.base)) {
				t.Fatal("export changed the source collection or its stored bindings")
			}
			exported, err := config.New(text).Load()
			if err != nil || len(exported.GetNoteTypes()) != nativeCount {
				t.Fatalf("exported type count: %v", err)
			}
			newPath := "note_types/" + filename
			for address, name := range map[string]string{trackedPath: "Renamed", newPath: test.nativeName} {
				value := &collectionpb.NoteTypeConfig{}
				readFixtureTOML(t, filepath.Join(text, address), value)
				if value.Name != name {
					t.Fatalf("type filename %q contains %q, want %q", address, value.Name, name)
				}
			}
			notePath := "decks/grammar/notes/Discovered.note.anki.toml"
			note := &collectionpb.NoteConfig{}
			readFixtureTOML(t, filepath.Join(text, notePath), note)
			if note.NoteType != newPath || note.Fields["Front"] != "Discovered" {
				t.Fatal("discovered note lost its content or allocated type reference")
			}
			beforeSources := fixtureSources(t, text)
			if err := bindFixtureConfigs(ctx, fixture.base, text, ""); err != nil {
				t.Fatal(err)
			}
			runFixtureCLI(t, "col", "export", "--collection", fixture.base, "--output", text)
			assertFixtureSources(t, text, beforeSources)
			assertConverged(t, fixture.base, text)
			saveFixtureArtifact(t, fixture.base, "type-names-"+test.name+".colpkg")
			writeFixtureTOML(t, stateFile, stateFixtureDump(t, fixture.base))
			saveFixtureArtifact(t, stateFile, "type-names-"+test.name+".state.anki.toml")
			saveFixtureArtifact(t, filepath.Join(text, newPath), "type-names-"+test.name+".notetype.anki.toml")
		})
	}
}

// TestReconciledExportWritesWorkflow verifies in-place output and caller-owned files/directories.
func TestReconciledExportWritesWorkflow(t *testing.T) {
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
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, "~/text"); err != nil {
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
	config := filepath.Join(text, "collection.anki.toml")
	before, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	if err := exportForFixture(context.Background(), filepath.Join(root, "missing.colpkg"), text, ""); err == nil {
		t.Fatal("invalid input accepted")
	}
	after, err := os.ReadFile(config)
	if err != nil || string(before) != string(after) {
		t.Fatal("failed export changed existing declarations")
	}
	fresh := filepath.Join(root, "fresh")
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, fresh); err != nil {
		t.Fatal(err)
	}
	assertConverged(t, fixture.base, fresh)
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

	saveFixtureArtifact(t, config, "exported.anki.toml")
}

// TestReadableNoteNamesWorkflow verifies safe sort-field labels and state-backed re-export through the CLI.
func TestReadableNoteNamesWorkflow(t *testing.T) {
	fixture := newFixture(t)
	ctx := context.Background()
	client := connectFixture(t, fixture.database)
	definition, err := client.types.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: fixture.typeID})
	if err != nil {
		t.Fatal(err)
	}
	definition.Config.SortFieldIdx = 1
	if _, err := client.types.UpdateNotetype(ctx, definition); err != nil {
		t.Fatal(err)
	}
	original, err := client.notes.GetNote(ctx, &notes.NoteId{Nid: fixture.noteID})
	if err != nil {
		t.Fatal(err)
	}
	original.Fields[1] = "<b>Hello &amp; 世界</b>"
	if _, err := client.notes.UpdateNotes(ctx, &notes.UpdateNotesRequest{Notes: []*notes.Note{original}}); err != nil {
		t.Fatal(err)
	}
	labels := []string{
		"<div>Hello &amp; 世界</div>", "hello 世界", "../../A/B\\C:*?\"<>|\n D",
		"<img src=\"image.png\">", "", strings.Repeat("长", 80),
		"<script>hidden()</script><style>.hidden {}</style><b>Visible</b>", "CON", "[sound:audio.mp3]", "Bound title",
	}
	var importedID int64
	for index, label := range labels {
		added, err := client.notes.AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{
			Guid: fmt.Sprintf("readable-name-%d", index), NotetypeId: fixture.typeID,
			Fields: []string{"你好", label, "hello"},
		}})
		if err != nil {
			t.Fatal(err)
		}
		importedID = added.NoteId
	}
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	packageFixture(t, fixture.database, fixture.base)
	text := filepath.Join(t.TempDir(), "text")
	writeFixtureTOML(t, filepath.Join(text, "collection.anki.toml"), &collectionpb.CollectionConfig{Version: 1})
	writeFixtureTOML(t, filepath.Join(text, "decks", "grammar", "deck.anki.toml"), &collectionpb.DeckConfig{Normal: &decks.Deck_Normal{}})
	writeFixtureTOML(t, filepath.Join(text, "decks", "grammar", "notes", "Hello-世界.note.anki.toml"), &collectionpb.NoteConfig{})
	runFixtureCLI(t, "col", "state", "import", "--collection", fixture.base, "--text", text,
		"decks/grammar/notes/Hello-世界.note.anki.toml="+fmt.Sprint(importedID))
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"Hello-世界.note.anki.toml":   "Bound title",
		"Hello-世界-2.note.anki.toml": original.Fields[1],
		"Hello-世界-3.note.anki.toml": labels[0], "hello-世界-4.note.anki.toml": labels[1],
		"A-B-C-D.note.anki.toml": labels[2], "note.note.anki.toml": labels[3],
		"note-2.note.anki.toml": labels[4], strings.Repeat("长", 40) + ".note.anki.toml": labels[5],
		"Visible.note.anki.toml": labels[6], "note-CON.note.anki.toml": labels[7],
		"note-3.note.anki.toml": labels[8],
	}
	noteDirectory := filepath.Join(text, "decks", "grammar", "notes")
	files, err := os.ReadDir(noteDirectory)
	if err != nil || len(files) != len(want) {
		t.Fatalf("exported note count: %d, want %d: %v", len(files), len(want), err)
	}
	for name, label := range want {
		value := &collectionpb.NoteConfig{}
		readFixtureTOML(t, filepath.Join(noteDirectory, name), value)
		if value.Fields["article-title"] != label {
			t.Fatalf("filename %q changed sort-field content: %q", name, value.Fields["article-title"])
		}
	}
	assertConverged(t, fixture.base, text)
	before := fixtureSources(t, text)
	archiveBefore, err := os.ReadFile(fixture.base)
	if err != nil {
		t.Fatal(err)
	}
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
		t.Fatal(err)
	}
	assertFixtureSources(t, text, before)
	archiveAfter, err := os.ReadFile(fixture.base)
	if err != nil || string(archiveBefore) != string(archiveAfter) {
		t.Fatal("repeated export changed source collection")
	}
	editedPath := filepath.Join(noteDirectory, "Hello-世界-2.note.anki.toml")
	edited := &collectionpb.NoteConfig{}
	readFixtureTOML(t, editedPath, edited)
	edited.Fields["article-title"] = "Different title"
	writeFixtureTOML(t, editedPath, edited)
	runFixtureCLI(t, "col", "apply", "--collection", fixture.base, "--text", text)
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
		t.Fatal(err)
	}
	readFixtureTOML(t, editedPath, edited)
	if edited.Fields["article-title"] != "Different title" {
		t.Fatal("edited note did not retain its bound filename")
	}
	client = inspectFixture(t, fixture.base)
	if _, err := client.notes.AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{
		Guid: "native-after-export", NotetypeId: fixture.typeID, Fields: []string{"你好", "Hello 世界", "hello"},
	}}); err != nil {
		t.Fatal(err)
	}
	client.exportFixture(t, fixture.base)
	if err := client.close(); err != nil {
		t.Fatal(err)
	}

	before = fixtureSources(t, text)
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
		t.Fatal(err)
	}
	after := fixtureSources(t, text)
	if len(after) != len(before)+1 {
		t.Fatal("new native note changed the existing file set")
	}
	for path, data := range before {
		if after[path] != data {
			t.Fatalf("new native note changed an existing config: %s", path)
		}
	}
	added := &collectionpb.NoteConfig{}
	readFixtureTOML(t, filepath.Join(noteDirectory, "Hello-世界-5.note.anki.toml"), added)
	if added.Fields["article-title"] != "Hello 世界" {
		t.Fatal("new native note did not receive an unoccupied content-derived name")
	}
	assertConverged(t, fixture.base, text)
	saveFixtureArtifact(t, fixture.base, "readable-note-names.colpkg")
	saveFixtureArtifact(t, editedPath, "readable-note-name.note.anki.toml")
}
