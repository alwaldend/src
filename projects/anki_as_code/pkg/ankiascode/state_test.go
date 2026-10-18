package ankiascode

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/deck_config"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// stateFixtureConfigs writes independent declarations without requiring collection export.
func stateFixtureConfigs(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFixtureTOML(t, filepath.Join(root, "collection.anki.toml"), &ankiascodepb.CollectionConfig{Version: 1})
	writeFixtureTOML(t, filepath.Join(root, "decks", "grammar", "deck.anki.toml"), &ankiascodepb.DeckConfig{Preset: "Default", Normal: &decks.Deck_Normal{}})
	writeFixtureTOML(t, filepath.Join(root, "note_types", "Chinese.notetype.anki.toml"), &ankiascodepb.NoteTypeConfig{Name: "Chinese"})
	for _, name := range []string{"example", "other"} {
		writeFixtureTOML(t, filepath.Join(root, "decks", "grammar", "notes", name+".note.anki.toml"), &ankiascodepb.NoteConfig{NoteType: "note_types/Chinese.notetype.anki.toml"})
	}
	return root
}

// stateFixtureDump decodes default TOML output independently and keeps status off stdout.
func stateFixtureDump(t *testing.T, input string) *ankiascodepb.State {
	t.Helper()
	return stateFixtureDumpKey(t, input, "")
}

// stateFixtureDumpKey reads one selected namespace through the built CLI and an independent TOML parser.
func stateFixtureDumpKey(t *testing.T, input, key string) *ankiascodepb.State {
	t.Helper()
	file := filepath.Join(t.TempDir(), "state.anki.toml")
	if err := os.WriteFile(file, runFixtureCLI(t, "col", "state", "dump", "--state-key", key, "--collection", input), 0o600); err != nil {
		t.Fatal(err)
	}
	value := &ankiascodepb.State{}
	readFixtureTOML(t, file, value)
	return value
}

// rejectStateCLI verifies rejected native operations emit only an error stream and returns it.
func rejectStateCLI(t *testing.T, arguments ...string) string {
	t.Helper()
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, binary, arguments...)
	var stdout, stderr strings.Builder
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err == nil || stdout.Len() != 0 || stderr.Len() == 0 {
		t.Fatalf("invalid state operation accepted: error=%v stdout=%q stderr=%q", err, stdout.String(), stderr.String())
	}
	return stderr.String()
}

// TestStateCLIWorkflow verifies stored-state dump, batched import, replacement, and unchanged inputs on rejection.
func TestStateCLIWorkflow(t *testing.T) {
	for _, format := range []string{"package", "database"} {
		t.Run(format, func(t *testing.T) {
			fixture := newFixture(t)
			input := fixture.base
			if format == "database" {
				input = fixture.database
			}
			text := stateFixtureConfigs(t)
			sources := fixtureSources(t, text)
			before, err := fixtureContent(t, input)
			if err != nil {
				t.Fatal(err)
			}
			empty := stateFixtureDump(t, input)
			if empty.Version != 1 || (len(empty.GetResources().GetDecks())+len(empty.GetResources().GetNotes())+len(empty.GetResources().GetNoteTypes())) != 0 {
				t.Fatal("dump initialized state")
			}
			after, err := fixtureContent(t, input)
			if err != nil || string(before) != string(after) {
				t.Fatal("dump modified collection content")
			}
			arguments := []string{
				"col", "state", "import", "--collection", input, "--text", text,
				"decks/grammar/deck.anki.toml=1",
				fmt.Sprintf("decks/grammar/notes/example.note.anki.toml=%d", fixture.noteID),
				fmt.Sprintf("note_types/Chinese.notetype.anki.toml=%d", fixture.typeID),
			}
			runFixtureCLI(t, arguments...)
			stored := stateFixtureDump(t, input)
			if stored.GetResources().GetNotes()["decks/grammar/notes/example.note.anki.toml"].GetId() != fixture.noteID || stored.GetResources().GetNoteTypes()["note_types/Chinese.notetype.anki.toml"].GetFields()["hanzi"] == nil {
				t.Fatal("batch did not bind roots and native members")
			}
			patch := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"updated": {Id: 1}}}}
			patchFile := filepath.Join(t.TempDir(), "update.anki.toml")
			writeFixtureTOML(t, patchFile, patch)
			runFixtureCLI(t, "col", "state", "update", "--collection", input, patchFile)
			stored.GetResources().Presets["updated"] = &ankiascodepb.PresetState{Id: 1}
			if !proto.Equal(stored, stateFixtureDump(t, input)) {
				t.Fatal("CLI update dropped omitted bindings")
			}
			patchJSON, err := protojson.Marshal(patch)
			if err != nil {
				t.Fatal(err)
			}
			patchJSONFile := filepath.Join(t.TempDir(), "update.anki.json")
			if err := os.WriteFile(patchJSONFile, patchJSON, 0o600); err != nil {
				t.Fatal(err)
			}
			beforeUpdate := stateFixtureDump(t, input)
			runFixtureCLI(t, "col", "state", "update", "--collection", input, patchFile)
			runFixtureCLI(t, "col", "state", "update", "--format", "json", "--collection", input, patchJSONFile)
			invalidPatch := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Notes: map[string]*ankiascodepb.NoteState{"decks/grammar/deck.anki.toml": {Id: 999}}}}
			writeFixtureTOML(t, patchFile, invalidPatch)
			rejectStateCLI(t, "col", "state", "update", "--collection", input, patchFile)
			afterUpdate := stateFixtureDump(t, input)
			if !proto.Equal(beforeUpdate, afterUpdate) {
				t.Fatal("no-op or rejected CLI update changed stored state")
			}
			customKey := "anki_as_code/state/alternate"
			custom := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"alternate": {Id: 1}}}}
			customFile := filepath.Join(t.TempDir(), "alternate.anki.toml")
			writeFixtureTOML(t, customFile, custom)
			runFixtureCLI(t, "col", "state", "replace", "--state-key", customKey, "--collection", input, customFile)
			runFixtureCLI(t, "col", "state", "import", "--state-key", customKey, "--collection", input, "--text", text, "presets/second=1")
			custom.GetResources().GetPresets()["second"] = &ankiascodepb.PresetState{Id: 1}
			if !proto.Equal(custom, stateFixtureDumpKey(t, input, customKey)) || !proto.Equal(stored, stateFixtureDump(t, input)) {
				t.Fatal("CLI state namespaces interfered")
			}
			before, err = fixtureContent(t, input)
			if err != nil {
				t.Fatal(err)
			}
			runFixtureCLI(t, arguments...)
			failure := rejectStateCLI(t, "col", "state", "import", "--collection", input, "--text", text, "notes/missing.note.anki.toml=opaque-id")
			if !strings.Contains(failure, "is not declared") {
				t.Fatalf("identifier checked before resource: %s", failure)
			}
			for _, identifier := range []string{"opaque-id", "", "0", "-1", "9223372036854775808"} {
				rejectStateCLI(t, "col", "state", "import", "--collection", input, "--text", text, "presets/new="+identifier)
			}
			rejectStateCLI(t, "col", "state", "import", "--collection", input, "--text", text, "presets/new=1", "presets/missing=9223372036854775807")
			rejectStateCLI(t, "col", "state", "import", "--collection", input, "--text", text, "presets/new=1", fmt.Sprintf("decks/grammar/notes/other.note.anki.toml=%d", fixture.noteID))
			after, err = fixtureContent(t, input)
			if err != nil || string(before) != string(after) {
				t.Fatal("no-op or rejected batch changed collection")
			}
			if !proto.Equal(stored, stateFixtureDump(t, input)) {
				t.Fatal("partial batch was published")
			}
			file := filepath.Join(t.TempDir(), "state.anki.toml")
			replacement := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"only": {Id: 1}}}}
			writeFixtureTOML(t, file, replacement)
			runFixtureCLI(t, "col", "state", "replace", "--collection", input, file)
			if !proto.Equal(replacement, stateFixtureDump(t, input)) {
				t.Fatal("replace merged rather than replaced state")
			}
			before, err = fixtureContent(t, input)
			if err != nil {
				t.Fatal(err)
			}
			runFixtureCLI(t, "col", "state", "replace", "--collection", input, file)
			after, err = fixtureContent(t, input)
			if err != nil || string(before) != string(after) {
				t.Fatal("identical replacement changed collection content")
			}
			jsonData := runFixtureCLI(t, "col", "state", "dump", "--format", "json", "--collection", input)
			jsonState := &ankiascodepb.State{}
			if err := protojson.Unmarshal(jsonData, jsonState); err != nil || !proto.Equal(replacement, jsonState) {
				t.Fatalf("JSON dump differs from TOML: %v", err)
			}
			jsonFile := filepath.Join(t.TempDir(), "state.anki.json")
			if err := os.WriteFile(jsonFile, jsonData, 0o600); err != nil {
				t.Fatal(err)
			}
			runFixtureCLI(t, "col", "state", "replace", "--format", "json", "--collection", input, jsonFile)
			tomlData := runFixtureCLI(t, "col", "state", "dump", "--collection", input)
			binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			command := exec.CommandContext(ctx, binary, "col", "state", "replace", "--collection", input, "-")
			command.Stdin = strings.NewReader(string(tomlData))
			if output, err := command.CombinedOutput(); err != nil {
				t.Fatalf("replace TOML from stdin: %v: %s", err, output)
			}
			invalid := filepath.Join(t.TempDir(), "invalid.toml")
			if err := os.WriteFile(invalid, []byte("version = ["), 0o600); err != nil {
				t.Fatal(err)
			}
			rejectStateCLI(t, "col", "state", "replace", "--collection", input, invalid)
			rejectStateCLI(t, "col", "state", "replace", "--format", "yaml", "--collection", input, file)
			rejectStateCLI(t, "col", "state", "dump", "--format", "yaml", "--collection", input)
			after, err = fixtureContent(t, input)
			if err != nil || string(before) != string(after) {
				t.Fatal("encoding round trips or invalid state input changed collection content")
			}
			if !proto.Equal(custom, stateFixtureDumpKey(t, input, customKey)) {
				t.Fatal("default state replacement changed the other namespace")
			}
			assertFixtureSources(t, text, sources)
			saveFixtureArtifact(t, file, "state-"+format+".toml")
			saveFixtureArtifact(t, jsonFile, "state-"+format+".json")
			saveFixtureArtifact(t, input, "state-"+format+filepath.Ext(input))
		})
	}
}

// TestApplicationStateWorkflows verifies typed application operations own native publication and preserve state isolation.
func TestApplicationStateWorkflows(t *testing.T) {
	for _, format := range []string{"package", "database"} {
		t.Run(format, func(t *testing.T) {
			fixture := newFixture(t)
			input := fixture.base
			if format == "database" {
				input = fixture.database
			}
			application := fixtureApplication()
			ctx := t.Context()
			text := stateFixtureConfigs(t)
			empty, err := application.ColStateDump(ctx, ankiascodeapi.CollectionOptions{Path: input})
			if err != nil || empty.GetVersion() != 1 || len(empty.GetResources().GetNotes()) != 0 {
				t.Fatalf("dump empty application state: value=%v error=%v", empty, err)
			}
			bindings := map[string]string{
				"decks/grammar/deck.anki.toml":               "1",
				"decks/grammar/notes/example.note.anki.toml": fmt.Sprint(fixture.noteID),
				"note_types/Chinese.notetype.anki.toml":      fmt.Sprint(fixture.typeID),
			}
			if err := application.ColStateImport(ctx, ankiascodeapi.CollectionOptions{Path: input}, text, bindings); err != nil {
				t.Fatal(err)
			}
			stored, err := application.ColStateDump(ctx, ankiascodeapi.CollectionOptions{StateKey: state.DefaultKey, Path: input})
			if err != nil || stored.GetResources().GetNotes()["decks/grammar/notes/example.note.anki.toml"].GetId() != fixture.noteID || stored.GetResources().GetNoteTypes()["note_types/Chinese.notetype.anki.toml"].GetFields()["hanzi"] == nil {
				t.Fatalf("application import did not persist native bindings: value=%v error=%v", stored, err)
			}
			patch := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"added": {Id: 1}}}}
			originalPatch := proto.Clone(patch)
			if err := application.ColStateUpdate(ctx, ankiascodeapi.CollectionOptions{Path: input}, patch); err != nil {
				t.Fatal(err)
			}
			stored.Resources.Presets["added"] = &ankiascodepb.PresetState{Id: 1}
			actual, err := application.ColStateDump(ctx, ankiascodeapi.CollectionOptions{Path: input})
			if err != nil || !proto.Equal(stored, actual) || !proto.Equal(patch, originalPatch) {
				t.Fatalf("application update lost bindings or changed its input: value=%v error=%v", actual, err)
			}
			key := "anki_as_code/state/application"
			if err := application.ColStateReplace(ctx, ankiascodeapi.CollectionOptions{StateKey: key, Path: input}, patch); err != nil {
				t.Fatal(err)
			}
			before, err := fixtureContent(t, input)
			if err != nil {
				t.Fatal(err)
			}
			if err := application.ColStateUpdate(ctx, ankiascodeapi.CollectionOptions{Path: input}, patch); err != nil {
				t.Fatal(err)
			}
			if err := application.ColStateImport(ctx, ankiascodeapi.CollectionOptions{Path: input}, text, map[string]string{"presets/new": "1", "presets/missing": "9223372036854775807"}); err == nil {
				t.Fatal("application accepted a partially invalid import")
			}
			if err := application.ColStateReplace(ctx, ankiascodeapi.CollectionOptions{Path: input}, &ankiascodepb.State{Version: 2}); err == nil {
				t.Fatal("application accepted an unsupported state version")
			}
			cancelled, cancel := context.WithCancel(ctx)
			cancel()
			if err := application.ColStateReplace(cancelled, ankiascodeapi.CollectionOptions{Path: input}, patch); !errors.Is(err, context.Canceled) {
				t.Fatalf("application lost cancellation: %v", err)
			}
			after, err := fixtureContent(t, input)
			if err != nil || string(before) != string(after) {
				t.Fatalf("no-op or rejected application operation modified the source: %v", err)
			}
			actual, err = application.ColStateDump(ctx, ankiascodeapi.CollectionOptions{Path: input})
			if err != nil || !proto.Equal(stored, actual) {
				t.Fatalf("rejected application operation wrote partial state or retained a reservation: value=%v error=%v", actual, err)
			}
			if err := application.ColStateReplace(ctx, ankiascodeapi.CollectionOptions{Path: input}, patch); err != nil {
				t.Fatal(err)
			}
			for _, namespace := range []string{"", key} {
				actual, err := application.ColStateDump(ctx, ankiascodeapi.CollectionOptions{StateKey: namespace, Path: input})
				if err != nil || !proto.Equal(patch, actual) {
					t.Fatalf("application replacement merged old state or crossed namespaces: value=%v error=%v", actual, err)
				}
			}
			saveFixtureArtifact(t, input, "application-state-"+format+filepath.Ext(input))
		})
	}
}

// TestStateAPIWorkflow verifies all state methods on one existing backend connection and a native package round trip.
func TestStateAPIWorkflow(t *testing.T) {
	fixture := newFixture(t)
	ctx := context.Background()
	client := connectFixture(t, fixture.database)
	connection, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	store := state.New(connection, "")
	before, err := store.Dump(ctx)
	if err != nil || before.Version != 1 || (len(before.GetResources().GetDecks())+len(before.GetResources().GetNotes())+len(before.GetResources().GetNoteTypes())) != 0 {
		t.Fatalf("empty state: %v %v", before, err)
	}
	configs, err := config.New(stateFixtureConfigs(t)).Load()
	if err != nil {
		t.Fatal(err)
	}
	note, err := client.notes.GetNote(ctx, &notes.NoteId{Nid: fixture.noteID})
	if err != nil {
		t.Fatal(err)
	}
	card, err := client.cards.GetCard(ctx, &cards.CardId{Cid: fixture.cardIDs[0]})
	if err != nil {
		t.Fatal(err)
	}
	history, err := client.stats.GetReviewLogs(ctx, &cards.CardId{Cid: fixture.cardIDs[0]})
	if err != nil {
		t.Fatal(err)
	}
	preset, err := client.presets.GetDeckConfig(ctx, &deck_config.DeckConfigId{Dcid: 1})
	if err != nil {
		t.Fatal(err)
	}
	bindings := map[string]string{"decks/grammar/deck.anki.toml": "1", "decks/grammar/notes/example.note.anki.toml": fmt.Sprint(fixture.noteID), "note_types/Chinese.notetype.anki.toml": fmt.Sprint(fixture.typeID)}
	if changed, err := store.Import(ctx, configs, nil); err != nil || changed {
		t.Fatalf("empty import: %v %v", changed, err)
	}
	if changed, err := store.Update(ctx, &ankiascodepb.State{Version: 1}); err != nil || changed {
		t.Fatalf("empty update: %v %v", changed, err)
	}
	if changed, err := store.Import(ctx, configs, bindings); err != nil || !changed {
		t.Fatalf("batch import: %v %v", changed, err)
	}
	stored, err := store.Dump(ctx)
	if err != nil {
		t.Fatal(err)
	}
	namedDefault, err := state.New(connection, "anki_as_code/state/v1").Dump(ctx)
	if err != nil || !proto.Equal(stored, namedDefault) {
		t.Fatal("empty API key did not select the versioned default")
	}
	for _, invalid := range []*ankiascodepb.State{
		{Version: 1, Resources: &ankiascodepb.Resources{Decks: map[string]*ankiascodepb.DeckState{"same": {Id: 1}}, Notes: map[string]*ankiascodepb.NoteState{"same": {Id: fixture.noteID}}}},
		{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"invalid": {}}}},
	} {
		if changed, err := store.Replace(ctx, invalid); err == nil || changed {
			t.Fatal("invalid typed state accepted")
		}
	}
	customKey := "anki_as_code/state/independent"
	custom := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"another": {Id: 1}}, NoteTypes: map[string]*ankiascodepb.NoteTypeState{
		"note_types/Chinese.notetype.anki.toml": {
			Id:        fixture.typeID,
			Fields:    map[string]*ankiascodepb.FieldState{"name/#/日本語": {Id: -7, Ordinal: 9}},
			Templates: map[string]*ankiascodepb.TemplateState{"Card/#/1": {Ordinal: 3}},
		},
	}}}
	if changed, err := state.New(connection, customKey).Replace(ctx, custom); err != nil || !changed {
		t.Fatalf("replace independent API namespace: %v %v", changed, err)
	}
	patch := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{
		Presets:   map[string]*ankiascodepb.PresetState{"added": {Id: 1}},
		NoteTypes: map[string]*ankiascodepb.NoteTypeState{"note_types/Chinese.notetype.anki.toml": {Id: fixture.typeID, Fields: map[string]*ankiascodepb.FieldState{"replacement/#/日本語": {Id: -9, Ordinal: 2}}}},
	}}
	originalPatch := proto.Clone(patch)
	independentStore := state.New(connection, customKey)
	if changed, err := independentStore.Update(ctx, patch); err != nil || !changed {
		t.Fatalf("update independent namespace: %v %v", changed, err)
	}
	custom.GetResources().Presets["added"] = &ankiascodepb.PresetState{Id: 1}
	custom.GetResources().NoteTypes["note_types/Chinese.notetype.anki.toml"] = proto.Clone(patch.Resources.NoteTypes["note_types/Chinese.notetype.anki.toml"]).(*ankiascodepb.NoteTypeState)
	updated, err := independentStore.Dump(ctx)
	if err != nil || !proto.Equal(updated, custom) || !proto.Equal(patch, originalPatch) {
		t.Fatal("update merge semantics or input ownership changed")
	}
	if changed, err := independentStore.Update(ctx, patch); err != nil || changed {
		t.Fatalf("no-op update: %v %v", changed, err)
	}
	conflict := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Notes: map[string]*ankiascodepb.NoteState{"notes/duplicate.note.anki.toml": {Id: fixture.noteID}}}}
	if changed, err := store.Update(ctx, conflict); err == nil || changed {
		t.Fatal("update accepted duplicate native ownership")
	}
	if changed, err := store.Update(ctx, &ankiascodepb.State{Version: 2}); err == nil || changed {
		t.Fatal("update accepted unsupported version")
	}
	if changed, err := store.Import(ctx, configs, bindings); err != nil || changed {
		t.Fatalf("no-op import: %v %v", changed, err)
	}
	if changed, err := store.Import(ctx, configs, map[string]string{"presets/A-new": "1", "presets/Z-missing": "9223372036854775807"}); err == nil || changed {
		t.Fatal("partial API import accepted")
	}
	if changed, err := store.Import(ctx, configs, map[string]string{"notes/missing.note.anki.toml": "opaque-id"}); err == nil || changed || !strings.Contains(err.Error(), "is not declared") {
		t.Fatalf("API identifier checked before resource: %v %v", changed, err)
	}
	actual, err := store.Dump(ctx)
	if err != nil || !proto.Equal(stored, actual) {
		t.Fatal("failed API import wrote state")
	}
	replacement := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"only": {Id: 1}}}}
	if changed, err := store.Replace(ctx, replacement); err != nil || !changed {
		t.Fatalf("replace: %v %v", changed, err)
	}
	if changed, err := store.Replace(ctx, replacement); err != nil || changed {
		t.Fatalf("no-op replace: %v %v", changed, err)
	}
	if changed, err := store.Replace(ctx, &ankiascodepb.State{Version: 2}); err == nil || changed {
		t.Fatal("unsupported replacement accepted")
	}
	if changed, err := store.Replace(ctx, &ankiascodepb.State{Version: 1}); err != nil || !changed {
		t.Fatalf("clear state: %v %v", changed, err)
	}
	actual, err = store.Dump(ctx)
	if err != nil || (len(actual.GetResources().GetDecks())+len(actual.GetResources().GetNotes())+len(actual.GetResources().GetNoteTypes())) != 0 || len(actual.GetResources().GetPresets()) != 0 {
		t.Fatal("empty replacement retained old bindings")
	}
	if _, err := store.Replace(ctx, stored); err != nil {
		t.Fatal(err)
	}
	client.exportFixture(t, fixture.base)
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	client = inspectFixture(t, fixture.base)
	connection, err = client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	actual, err = state.New(connection, "").Dump(ctx)
	if err != nil || !proto.Equal(stored, actual) {
		t.Fatal("native package round trip lost state")
	}
	independent, err := state.New(connection, customKey).Dump(ctx)
	if err != nil || !proto.Equal(custom, independent) {
		t.Fatal("native package round trip or another namespace lost nested bindings")
	}
	for address := range actual.GetResources().GetNoteTypes() {
		if strings.Contains(address, "#fields/") || strings.Contains(address, "#templates/") || strings.HasPrefix(address, "presets/") {
			t.Fatalf("synthetic address persisted: %q", address)
		}
	}
	actualNote, err := client.notes.GetNote(ctx, &notes.NoteId{Nid: fixture.noteID})
	if err != nil || !proto.Equal(note, actualNote) {
		t.Fatal("state operations changed note")
	}
	actualCard, err := client.cards.GetCard(ctx, &cards.CardId{Cid: fixture.cardIDs[0]})
	if err != nil || !proto.Equal(card, actualCard) {
		t.Fatal("state operations changed card")
	}
	actualHistory, err := client.stats.GetReviewLogs(ctx, &cards.CardId{Cid: fixture.cardIDs[0]})
	if err != nil || !proto.Equal(history, actualHistory) {
		t.Fatal("state operations changed review history")
	}
	actualPreset, err := client.presets.GetDeckConfig(ctx, &deck_config.DeckConfigId{Dcid: 1})
	if err != nil || !proto.Equal(preset, actualPreset) {
		t.Fatal("state operations changed preset")
	}
	assertFixtureMedia(t, client)
	if err := client.close(); err != nil {
		t.Fatal(err)
	}
	saveFixtureArtifact(t, fixture.base, "state-api-roundtrip.colpkg")
}

// prepareFixtureConfigs exports configs and separately binds their native fixture objects.
func prepareFixtureConfigs(t *testing.T, ctx context.Context, input, output string) error {
	t.Helper()
	if err := exportForFixture(ctx, input, output, ""); err != nil {
		return fmt.Errorf("export fixture configs: %w", err)
	}
	root, err := expandPath(output)
	if err != nil {
		return fmt.Errorf("resolve fixture export: %w", err)
	}
	if err := bindFixtureConfigs(ctx, input, root, ""); err != nil {
		return fmt.Errorf("bind fixture configs: %w", err)
	}
	return nil
}

// bindFixtureConfigs prepares state independently by matching exact native fixture values.
func bindFixtureConfigs(ctx context.Context, input, root, key string) error {
	configs, err := config.New(root).Load()
	if err != nil {
		return fmt.Errorf("load fixture configs: %w", err)
	}
	if err := runFixtureCollection(ctx, input, input, func(connection grpc.ClientConnInterface) (bool, error) {
		snapshot, err := resource.New(connection).Snapshot(ctx)
		if err != nil {
			return false, fmt.Errorf("read fixture resources: %w", err)
		}
		resources := &ankiascodepb.Resources{
			Decks: map[string]*ankiascodepb.DeckState{}, Notes: map[string]*ankiascodepb.NoteState{},
			NoteTypes: map[string]*ankiascodepb.NoteTypeState{}, Presets: map[string]*ankiascodepb.PresetState{},
		}
		decksByName := map[string]*decks.Deck{}
		for _, native := range snapshot.Decks {
			decksByName[native.Name] = native
		}
		files := slices.Clone(configs.Notes)
		for _, group := range configs.Decks {
			native := decksByName[group.Title]
			resources.Decks[group.Path] = &ankiascodepb.DeckState{Id: native.Id}
			files = append(files, group.Notes...)
			if group.Deck.Preset != "" {
				resources.Presets[group.Deck.Preset] = &ankiascodepb.PresetState{Id: native.GetNormal().GetConfigId()}
			}
		}
		typesByName := map[string]*notetypes.Notetype{}
		for _, native := range snapshot.NoteTypes {
			typesByName[native.Name] = native
		}
		definitions := map[string]*ankiascodepb.NoteTypeConfig{}
		for _, file := range configs.NoteTypes {
			definitions[file.Path] = file.Config
			native := typesByName[file.Config.Name]
			binding := &ankiascodepb.NoteTypeState{Id: native.Id, Fields: map[string]*ankiascodepb.FieldState{}, Templates: map[string]*ankiascodepb.TemplateState{}}
			for index, field := range native.Fields {
				binding.Fields[field.Name] = &ankiascodepb.FieldState{Id: field.Config.GetId(), Ordinal: uint32(index)}
			}
			for index, template := range native.Templates {
				binding.Templates[template.Name] = &ankiascodepb.TemplateState{Id: template.Config.GetId(), Ordinal: uint32(index)}
			}
			resources.NoteTypes[file.Path] = binding
		}
		candidates := map[string][]int64{}
		for _, native := range snapshot.Notes {
			identity := fmt.Sprintf("%d:%s", native.NotetypeId, strings.Join(native.Fields, "\x1f"))
			candidates[identity] = append(candidates[identity], native.Id)
		}
		for _, file := range files {
			var fields []string
			for _, field := range definitions[file.Config.NoteType].Fields {
				fields = append(fields, file.Config.Fields[field.Name])
			}
			identity := fmt.Sprintf("%d:%s", resources.NoteTypes[file.Config.NoteType].Id, strings.Join(fields, "\x1f"))
			matches := candidates[identity]
			if len(matches) == 0 {
				return false, fmt.Errorf("no native fixture note matches %q", file.Path)
			}
			resources.Notes[file.Path] = &ankiascodepb.NoteState{Id: matches[0]}
			candidates[identity] = matches[1:]
		}
		changed, err := state.New(connection, key).Replace(ctx, &ankiascodepb.State{Version: 1, Resources: resources})
		if err != nil {
			return false, fmt.Errorf("install fixture bindings: %w", err)
		}
		return changed, nil
	}); err != nil {
		return fmt.Errorf("bind native fixture: %w", err)
	}
	return nil
}

// TestReconciledExportWorkflow verifies config-only output, read-only export and independent state dumps.
func TestReconciledExportWorkflow(t *testing.T) {
	for _, inputFormat := range []string{"package", "database"} {
		t.Run(inputFormat, func(t *testing.T) {
			fixture := newFixture(t)
			input := fixture.base
			if inputFormat == "database" {
				input = fixture.database
			}
			before, err := fixtureContent(t, input)
			if err != nil {
				t.Fatal(err)
			}
			text := filepath.Join(t.TempDir(), "text")
			runFixtureCLI(t, "col", "export", "--collection", input, "--output", text)
			after, err := fixtureContent(t, input)
			if err != nil || string(before) != string(after) {
				t.Fatal("export modified source collection content")
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
			if err := bindFixtureConfigs(context.Background(), input, text, ""); err != nil {
				t.Fatal(err)
			}
			stored := stateFixtureDump(t, input)
			if stored.GetResources().GetNotes()["decks/grammar/notes/你好.note.anki.toml"].GetId() != fixture.noteID {
				t.Fatal("separate state setup did not bind the native note")
			}
			for _, format := range []string{"toml", "json"} {
				data := runFixtureCLI(t, "col", "state", "dump", "--format", format, "--collection", input)
				value := &ankiascodepb.State{}
				if err := codec.Decode(data, codec.Format(format), value); err != nil || !proto.Equal(value, stored) {
					t.Fatalf("separate %s state dump: %v", format, err)
				}
				file := filepath.Join(t.TempDir(), "state.anki."+format)
				if err := os.WriteFile(file, data, 0o600); err != nil {
					t.Fatal(err)
				}
				saveFixtureArtifact(t, file, "dumped-state-"+inputFormat+"."+format)
			}
			assertConverged(t, input, text)
		})
	}
}

// TestNamespacedCollectionWorkflow verifies selected-key moves, export naming, planning, and native builds.
func TestNamespacedCollectionWorkflow(t *testing.T) {
	fixture := newFixture(t)
	key := "anki_as_code/state/workflow"
	text := filepath.Join(t.TempDir(), "configs")
	runFixtureCLI(t, "col", "export", "--state-key", key, "--collection", fixture.base, "--output", text)
	if err := bindFixtureConfigs(context.Background(), fixture.base, text, key); err != nil {
		t.Fatal(err)
	}
	oldPath := "decks/grammar/notes/你好.note.anki.toml"
	newPath := "decks/grammar/notes/renamed.note.anki.toml"
	if err := os.Rename(filepath.Join(text, oldPath), filepath.Join(text, newPath)); err != nil {
		t.Fatal(err)
	}
	writeFixtureTOML(t, filepath.Join(text, "collection.anki.toml"), &ankiascodepb.CollectionConfig{Version: 1, Moved: []*ankiascodepb.ResourceMove{{From: oldPath, To: newPath}}})
	plan := &ankiascodepb.CollectionChangePlan{}
	if err := protojson.Unmarshal(runFixtureCLI(t, "col", "plan", "--state-key", key, "--collection", fixture.base, "--text", text), plan); err != nil {
		t.Fatal(err)
	}
	if len(plan.Changes) == 0 || plan.NotesAdded != 0 || plan.NotesDeleted != 0 {
		t.Fatal("namespaced move did not preserve note identity")
	}
	runFixtureCLI(t, "col", "apply", "--state-key", key, "--collection", fixture.base, "--text", text)
	stored := stateFixtureDumpKey(t, fixture.base, key)
	if stored.GetResources().GetNotes()[newPath].GetId() != fixture.noteID || stored.GetResources().GetNotes()[oldPath] != nil {
		t.Fatal("namespaced apply did not move the binding")
	}
	if err := protojson.Unmarshal(runFixtureCLI(t, "col", "plan", "--state-key", key, "--collection", fixture.base, "--text", text), plan); err != nil || len(plan.Changes) != 0 {
		t.Fatal("namespaced plan did not converge")
	}
	reexported := filepath.Join(t.TempDir(), "reexported")
	runFixtureCLI(t, "col", "export", "--state-key", key, "--collection", fixture.base, "--output", reexported)
	if _, err := os.Stat(filepath.Join(reexported, newPath)); err != nil {
		t.Fatal("namespaced export lost the bound note filename")
	}
	output := filepath.Join(t.TempDir(), "built.colpkg")
	runFixtureCLI(t, "col", "build", "--state-key", key, "--base", fixture.base, "--text", text, "--output", output)
	if !proto.Equal(stored, stateFixtureDumpKey(t, output, key)) {
		t.Fatal("namespaced build lost state")
	}
	defaultState := stateFixtureDump(t, output)
	if (len(defaultState.GetResources().GetDecks())+len(defaultState.GetResources().GetNotes())+len(defaultState.GetResources().GetNoteTypes())) != 0 || len(defaultState.GetResources().GetPresets()) != 0 {
		t.Fatal("namespaced collection commands initialized the default key")
	}
	saveFixtureArtifact(t, output, "namespaced-reconciliation.colpkg")
	stateFile := filepath.Join(t.TempDir(), "namespaced-state.anki.toml")
	writeFixtureTOML(t, stateFile, stored)
	saveFixtureArtifact(t, stateFile, "namespaced-state.anki.toml")
}
