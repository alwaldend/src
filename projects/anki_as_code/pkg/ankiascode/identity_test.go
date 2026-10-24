package ankiascode

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
	"github.com/ankitects/anki/proto/anki/cards"
	ankicol "github.com/ankitects/anki/proto/anki/collection"
	ankiconfig "github.com/ankitects/anki/proto/anki/config"
	"github.com/ankitects/anki/proto/anki/deck_config"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/ankitects/anki/proto/anki/notes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// fixtureSources records editable source content for independent read-only verification.
func fixtureSources(t *testing.T, root string) map[string]string {
	t.Helper()
	files := map[string]string{}
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk sources: %w", err)
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read source: %w", err)
		}
		files[path] = string(data)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return files
}

// TestRequiredDeckRemovalWorkflow refuses deletion by native identity before applying any writes.
func TestRequiredDeckRemovalWorkflow(t *testing.T) {
	for _, format := range []string{"database", "package"} {
		t.Run(format, func(t *testing.T) {
			ctx := t.Context()
			root := t.TempDir()
			database := filepath.Join(root, "collection.anki2")
			client := connectFixture(t, database)
			deck, err := client.decks.GetDeck(ctx, &decks.DeckId{Did: 1})
			if err != nil {
				t.Fatal(err)
			}
			deck.Name = "Required"
			if _, err := client.decks.UpdateDeck(ctx, deck); err != nil {
				t.Fatal(err)
			}
			if err := client.close(); err != nil {
				t.Fatal(err)
			}
			input := database
			if format == "package" {
				input = filepath.Join(root, "collection.colpkg")
				packageFixture(t, database, input)
			}
			text := filepath.Join(root, "text")
			if err := prepareFixtureConfigs(t, ctx, input, text); err != nil {
				t.Fatal(err)
			}
			before, err := fixtureContent(t, input)
			if err != nil {
				t.Fatal(err)
			}
			marker := filepath.Join(text, "decks", "Required", "deck.anki.toml")
			original, err := os.ReadFile(marker)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.RemoveAll(filepath.Dir(marker)); err != nil {
				t.Fatal(err)
			}
			additional := filepath.Join(text, "decks", "Additional", "deck.anki.toml")
			writeFixtureTOML(t, additional, &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}})
			sources := fixtureSources(t, text)
			if _, err := planForFixture(ctx, input, text, ""); err == nil {
				t.Error("plan accepted removing the renamed required default deck")
			}
			if err := applyForFixture(ctx, input, text, ""); err == nil {
				t.Error("apply accepted removing the renamed required default deck")
			}
			assertFixtureSources(t, text, sources)
			after, err := fixtureContent(t, input)
			if err != nil || string(before) != string(after) {
				t.Fatal("required-deck refusal changed native resources or stored bindings")
			}
			if id := stateFixtureDump(t, input).GetResources().GetDecks()["decks/Required/deck.anki.toml"].GetId(); id != 1 {
				t.Fatalf("required deck lost its stored identity: %d", id)
			}
			if err := os.RemoveAll(filepath.Dir(additional)); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Dir(marker), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(marker, original, 0o600); err != nil {
				t.Fatal(err)
			}
			assertConverged(t, input, text)
			saveFixtureArtifact(t, input, "required-deck-preserved"+filepath.Ext(input))
		})
	}
}

// TestUntrackedObjectsWorkflow preserves objects added in Anki after export and rejects destructive shared edits.
func TestUntrackedObjectsWorkflow(t *testing.T) {
	fixture := newFixture(t)
	text := filepath.Join(t.TempDir(), "text")
	ctx := context.Background()
	if err := prepareFixtureConfigs(t, ctx, fixture.base, text); err != nil {
		t.Fatal(err)
	}
	client := inspectFixture(t, fixture.base)
	added, err := client.notes.AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{
		Guid: "untracked-note", NotetypeId: fixture.typeID, Fields: []string{"Native", "Anki", "Untracked"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.decks.AddDeck(ctx, &decks.Deck{Name: "grammar::Untracked", Kind: &decks.Deck_Normal_{Normal: &decks.Deck_Normal{ConfigId: 1}}}); err != nil {
		t.Fatal(err)
	}
	client.exportFixture(t, fixture.base)
	if err := client.close(); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(fixtureNotePath(text)); err != nil {
		t.Fatal(err)
	}
	plan := fixturePlan(t, fixture.base, text)
	if plan.NotesDeleted != 1 {
		t.Fatalf("untracked note included in deletion: %v", plan)
	}
	if err := applyForFixture(ctx, fixture.base, text, ""); err != nil {
		t.Fatal(err)
	}
	assertConverged(t, fixture.base, text)
	client = inspectFixture(t, fixture.base)
	if value, err := client.notes.GetNote(ctx, &notes.NoteId{Nid: added.NoteId}); err != nil || value.GetGuid() != "untracked-note" {
		t.Fatalf("untracked note lost: %v / %v", value, err)
	}
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	before, err := fixtureContent(t, fixture.base)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(text, "decks", "grammar"), filepath.Join(text, "decks", "Language")); err != nil {
		t.Fatal(err)
	}
	writeFixtureTOML(t, filepath.Join(text, "collection.anki.toml"), &ankiascodepb.CollectionConfig{Version: 1, Moved: []*ankiascodepb.ResourceMove{{From: "decks/grammar", To: "decks/Language"}}})
	if err := applyForFixture(ctx, fixture.base, text, ""); err == nil {
		t.Fatal("renaming parent of untracked deck was accepted")
	}
	after, err := fixtureContent(t, fixture.base)
	if err != nil || string(after) != string(before) {
		t.Fatal("rejected operation changed input")
	}
	saveFixtureArtifact(t, fixture.base, "untracked-objects.colpkg")
}

// TestAnkiPackageStateRoundTrip verifies state through Anki's actual full-package import and export APIs.
func TestAnkiPackageStateRoundTrip(t *testing.T) {
	fixture := newFixture(t)
	text := filepath.Join(t.TempDir(), "text")
	if err := prepareFixtureConfigs(t, context.Background(), fixture.database, text); err != nil {
		t.Fatal(err)
	}
	assertConverged(t, fixture.database, text)
	client := connectFixture(t, fixture.database)
	connection, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	packages := import_export.NewBackendImportExportServiceClient(connection)
	first := filepath.Join(t.TempDir(), "first.colpkg")
	if _, err := packages.ExportCollectionPackage(context.Background(), &import_export.ExportCollectionPackageRequest{OutPath: first}); err != nil {
		t.Fatal(err)
	}
	client.opened = false
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	database := filepath.Join(t.TempDir(), "imported.anki2")
	client = connectFixture(t, database)
	connection, err = client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.collections.CloseCollection(context.Background(), &ankicol.CloseCollectionRequest{}); err != nil {
		t.Fatal(err)
	}
	client.opened = false
	packages = import_export.NewBackendImportExportServiceClient(connection)
	if _, err := packages.ImportCollectionPackage(context.Background(), &import_export.ImportCollectionPackageRequest{
		ColPath: database, BackupPath: first, MediaFolder: database + ".media", MediaDb: database + ".media.db",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.collections.OpenCollection(context.Background(), &ankicol.OpenCollectionRequest{
		CollectionPath: database, MediaFolderPath: database + ".media", MediaDbPath: database + ".media.db",
	}); err != nil {
		t.Fatal(err)
	}
	second := filepath.Join(t.TempDir(), "roundtrip.colpkg")
	client.opened = true
	if _, err := packages.ExportCollectionPackage(context.Background(), &import_export.ExportCollectionPackageRequest{OutPath: second}); err != nil {
		t.Fatal(err)
	}
	client.opened = false
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	assertConverged(t, second, text)
	saveFixtureArtifact(t, second, "anki-roundtrip-state.colpkg")
}

// assertFixtureSources checks exact file sets and bytes after collection operations.
func assertFixtureSources(t *testing.T, root string, before map[string]string) {
	t.Helper()
	after := fixtureSources(t, root)
	if len(before) != len(after) {
		t.Fatal("source file set changed")
	}
	for path, data := range before {
		if after[path] != data {
			t.Fatalf("source changed: %s", path)
		}
	}
}

// TestPathStateWorkflow verifies identity-free export, persisted moves, copied notes, and source preservation.
func TestPathStateWorkflow(t *testing.T) {
	fixture := newFixture(t)
	text := filepath.Join(t.TempDir(), "text")
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
		t.Fatal(err)
	}
	files, err := filepath.Glob(filepath.Join(text, "decks", "grammar", "notes", "*.note.anki.toml"))
	if err != nil || len(files) != 1 {
		t.Fatalf("exported notes: %v / %v", files, err)
	}
	notePath := files[0]
	for _, path := range []string{notePath, filepath.Join(text, "decks", "grammar", "deck.anki.toml")} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "id =") || strings.HasPrefix(line, "guid =") || strings.HasPrefix(line, "config_id =") {
				t.Fatalf("native identity exported at %s: %s", path, line)
			}
		}
	}
	assertConverged(t, fixture.base, text)
	note := &ankiascodepb.NoteConfig{}
	readFixtureTOML(t, notePath, note)
	note.Fields["hanzi"] = "Updated content"
	writeFixtureTOML(t, notePath, note)
	if err := os.Rename(filepath.Join(text, "decks", "grammar"), filepath.Join(text, "decks", "Language")); err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(text, "collection.anki.toml")
	if err := os.WriteFile(root, []byte("version = 1\n[[moved]]\nfrom = 'decks/grammar'\nto = 'decks/Language'\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	sources := fixtureSources(t, text)
	plan := fixturePlan(t, fixture.base, text)
	if plan.DecksAdded != 0 || plan.DecksDeleted != 0 || plan.DecksChanged != 1 || plan.NotesAdded != 0 || plan.NotesDeleted != 0 || plan.NotesChanged != 1 {
		t.Fatalf("move replaced identities: %v", plan)
	}
	if err := applyForFixture(context.Background(), fixture.base, text, ""); err != nil {
		t.Fatal(err)
	}
	assertFixtureSources(t, text, sources)
	assertConverged(t, fixture.base, text)
	client := inspectFixture(t, fixture.base)
	changed, err := client.notes.GetNote(context.Background(), &notes.NoteId{Nid: fixture.noteID})
	if err != nil || changed.GetFields()[0] != "Updated content" || changed.GetGuid() != "fixture-note" {
		t.Fatalf("note identity changed: %v / %v", changed, err)
	}
	renamed, err := client.decks.GetDeckIdByName(context.Background(), &generic.String{Val: "Language"})
	if err != nil || renamed.GetDid() != 1 {
		t.Fatalf("deck identity changed: %v / %v", renamed, err)
	}
	for _, id := range fixture.cardIDs {
		if _, err := client.cards.GetCard(context.Background(), &cards.CardId{Cid: id}); err != nil {
			t.Fatalf("retained card lost: %v", err)
		}
	}
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	movedNote := filepath.Join(text, "decks", "Language", "notes", filepath.Base(notePath))
	copyPath := filepath.Join(filepath.Dir(movedNote), "copy.note.anki.toml")
	data, err := os.ReadFile(movedNote)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(copyPath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	added := fixturePlan(t, fixture.base, text)
	if added.NotesAdded != 1 || added.CardsAdded != 0 || !added.CardGenerationDeferred {
		t.Fatalf("copied file did not create a new note: %v", added)
	}
	sources = fixtureSources(t, text)
	result := filepath.Join(t.TempDir(), "state.colpkg")
	if err := buildForFixture(context.Background(), fixture.base, text, result, ""); err != nil {
		t.Fatal(err)
	}
	assertFixtureSources(t, text, sources)
	assertConverged(t, result, text)
	before, err := fixtureContent(t, result)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(root, []byte("version=1\n[[moved]]\nfrom='decks/Language/notes/"+filepath.Base(notePath)+"'\nto='decks/Language/notes/copy.note.anki.toml'\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := applyForFixture(context.Background(), result, text, ""); err == nil {
		t.Fatal("conflicting move accepted")
	}
	after, err := fixtureContent(t, result)
	if err != nil || string(before) != string(after) {
		t.Fatal("rejected move changed collection")
	}
	writeFixtureTOML(t, filepath.Join(text, "decks", "grammar", "deck.anki.toml"), &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}})
	writeFixtureTOML(t, root, &ankiascodepb.CollectionConfig{Version: 1, Moved: []*ankiascodepb.ResourceMove{{From: "decks/grammar", To: "decks/Language"}}})
	if err := applyForFixture(context.Background(), result, text, ""); err == nil {
		t.Fatal("a move source still declared in configs was accepted")
	}
	after, err = fixtureContent(t, result)
	if err != nil || string(before) != string(after) {
		t.Fatal("declared move source changed collection")
	}
	saveFixtureArtifact(t, result, "path-state.colpkg")
	saveFixtureArtifact(t, movedNote, "path-state.note.anki.toml")
}

// fixtureBindings reads persisted identity bindings independently through Anki.
func fixtureBindings(t *testing.T, client *fixtureClient) *ankiascodepb.State {
	t.Helper()
	connection, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	values, err := ankiconfig.NewConfigServiceClient(connection).GetAllConfig(context.Background(), &generic.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	var entries map[string]json.RawMessage
	if err := json.Unmarshal(values.Json, &entries); err != nil {
		t.Fatal(err)
	}
	state := &ankiascodepb.State{}
	if err := protojson.Unmarshal(entries["anki_as_code/state/v1"], state); err != nil {
		t.Fatal(err)
	}
	return state
}

// TestPresetAliasExportWorkflow verifies imported aliases disambiguate native presets without modifying them.
func TestPresetAliasExportWorkflow(t *testing.T) {
	fixture := newFixture(t)
	ctx := context.Background()
	client := connectFixture(t, fixture.database)
	original, err := client.presets.GetDeckConfig(ctx, &deck_config.DeckConfigId{Dcid: 1})
	if err != nil {
		t.Fatal(err)
	}
	duplicate := proto.Clone(original).(*deck_config.DeckConfig)
	duplicate.Id = 0
	if _, err := client.presets.UpdateDeckConfigs(ctx, &deck_config.UpdateDeckConfigsRequest{
		TargetDeckId: 1, Configs: []*deck_config.DeckConfig{duplicate},
	}); err != nil {
		t.Fatal(err)
	}
	first, err := client.decks.GetDeck(ctx, &decks.DeckId{Did: 1})
	if err != nil {
		t.Fatal(err)
	}
	presetID := first.GetNormal().ConfigId
	if presetID == 1 {
		t.Fatal("fixture did not create a distinct preset")
	}
	duplicate, err = client.presets.GetDeckConfig(ctx, &deck_config.DeckConfigId{Dcid: presetID})
	if err != nil {
		t.Fatal(err)
	}
	if duplicate.Name != original.Name {
		t.Fatal("fixture did not create duplicate preset names")
	}
	second, err := client.decks.NewDeck(ctx, &generic.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	second.Name = "other"
	second.GetNormal().ConfigId = 1
	if _, err := client.decks.AddDeck(ctx, second); err != nil {
		t.Fatal(err)
	}
	client.exportFixture(t, fixture.base)
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	text := filepath.Join(t.TempDir(), "text")
	if err := exportForFixture(ctx, fixture.base, text, ""); err == nil {
		t.Fatal("ambiguous native preset names accepted without aliases")
	}
	writeFixtureTOML(t, filepath.Join(text, "collection.anki.toml"), &ankiascodepb.CollectionConfig{Version: 1})
	runFixtureCLI(t, "col", "state", "import", "--collection", fixture.base, "--text", text, "presets/primary="+fmt.Sprint(presetID))
	runFixtureCLI(t, "col", "state", "import", "--collection", fixture.base, "--text", text, "presets/secondary=1")
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
		t.Fatal(err)
	}
	for title, alias := range map[string]string{"grammar": "primary", "other": "secondary"} {
		value := &ankiascodepb.DeckConfig{}
		readFixtureTOML(t, filepath.Join(text, "decks", title, "deck.anki.toml"), value)
		if value.Preset != alias {
			t.Fatalf("deck %q preset reference: %q, want %q", title, value.Preset, alias)
		}
	}
	client = inspectFixture(t, fixture.base)
	state := fixtureBindings(t, client)
	for alias, id := range map[string]int64{"primary": presetID, "secondary": 1} {
		if state.GetResources().GetPresets()[alias].GetId() != id {
			t.Fatalf("export lost preset alias %q", alias)
		}
	}
	for _, expected := range []*deck_config.DeckConfig{original, duplicate} {
		actual, err := client.presets.GetDeckConfig(ctx, &deck_config.DeckConfigId{Dcid: expected.Id})
		if err != nil || !proto.Equal(expected, actual) {
			t.Fatalf("export modified native scheduling preset %d: %v", expected.Id, err)
		}
	}
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	before, err := fixtureContent(t, fixture.base)
	if err != nil {
		t.Fatal(err)
	}
	sources := fixtureSources(t, text)
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
		t.Fatal(err)
	}
	assertFixtureSources(t, text, sources)
	after, err := fixtureContent(t, fixture.base)
	if err != nil || string(before) != string(after) {
		t.Fatal("repeated aliased export changed collection content")
	}
	assertConverged(t, fixture.base, text)
	saveFixtureArtifact(t, fixture.base, "preset-aliases.colpkg")
	saveFixtureArtifact(t, filepath.Join(text, "decks", "grammar", "deck.anki.toml"), "preset-alias.deck.anki.toml")
}

// TestLiteralNoteTypePathWorkflow verifies that # in another filename does not imply member ownership.
func TestLiteralNoteTypePathWorkflow(t *testing.T) {
	fixture := newFixture(t)
	text := filepath.Join(t.TempDir(), "configs")
	if err := prepareFixtureConfigs(t, context.Background(), fixture.base, text); err != nil {
		t.Fatal(err)
	}
	configs, err := config.New(text, filesystem.FS{}).Load()
	if err != nil {
		t.Fatal(err)
	}
	original := "note_types/Chinese.notetype.anki.toml"
	other := ""
	for _, entry := range configs.NoteTypes {
		if entry.Path != original {
			other = entry.Path
			break
		}
	}
	if other == "" {
		t.Fatal("fixture requires an independent note type")
	}
	literal := original + "#literal.notetype.anki.toml"
	if err := os.Rename(filepath.Join(text, other), filepath.Join(text, literal)); err != nil {
		t.Fatal(err)
	}
	writeFixtureTOML(t, filepath.Join(text, "collection.anki.toml"), &ankiascodepb.CollectionConfig{Version: 1, Moved: []*ankiascodepb.ResourceMove{{From: other, To: literal}}})
	runFixtureCLI(t, "col", "apply", "--collection", fixture.base, "--text", text)
	bound := stateFixtureDump(t, fixture.base).GetResources().GetNoteTypes()[literal].GetId()
	renamed := "note_types/Renamed.notetype.anki.toml"
	if err := os.Rename(filepath.Join(text, original), filepath.Join(text, renamed)); err != nil {
		t.Fatal(err)
	}
	for _, deck := range configs.Decks {
		for _, note := range deck.Notes {
			if note.Config.NoteType == original {
				note.Config.NoteType = renamed
				writeFixtureTOML(t, filepath.Join(text, note.Path), note.Config)
			}
		}
	}
	writeFixtureTOML(t, filepath.Join(text, "collection.anki.toml"), &ankiascodepb.CollectionConfig{Version: 1, Moved: []*ankiascodepb.ResourceMove{{From: original, To: renamed}}})
	runFixtureCLI(t, "col", "apply", "--collection", fixture.base, "--text", text)
	actual := stateFixtureDump(t, fixture.base)
	if actual.GetResources().GetNoteTypes()[literal].GetId() != bound || actual.GetResources().GetNoteTypes()[renamed].GetId() != fixture.typeID {
		t.Fatal("literal filename was moved as a member address")
	}
	assertConverged(t, fixture.base, text)
	saveFixtureArtifact(t, fixture.base, "literal-note-type-path.colpkg")
}
