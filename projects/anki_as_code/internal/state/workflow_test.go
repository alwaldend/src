package state_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/collection"
	ankiconfig "github.com/ankitects/anki/proto/anki/config"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// TestStateWorkflow exercises CLI and Go state operations on a caller-owned collection.
// Commands own managed backends and must persist bindings without changing native objects.
func TestStateWorkflow(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	connection := startBackend(t, ctx)
	collections := collection.NewBackendCollectionServiceClient(connection)
	database := filepath.Join(t.TempDir(), "collection.anki2")
	open := &collection.OpenCollectionRequest{CollectionPath: database, MediaFolderPath: strings.TrimSuffix(database, ".anki2") + ".media", MediaDbPath: strings.TrimSuffix(database, ".anki2") + ".media.db2"}
	if _, err := collections.OpenCollection(ctx, open, grpc.WaitForReady(true)); err != nil {
		t.Fatal(err)
	}
	opened := true
	t.Cleanup(func() {
		if opened {
			cleanup, stop := context.WithTimeout(context.Background(), 5*time.Second)
			defer stop()
			if _, err := collections.CloseCollection(cleanup, &collection.CloseCollectionRequest{}); err != nil {
				t.Error(err)
			}
		}
	})
	types := notetypes.NewNotetypesServiceClient(connection)
	typeID, err := types.GetNotetypeIdByName(ctx, &generic.String{Val: "Basic"})
	if err != nil {
		t.Fatal(err)
	}
	items := notes.NewNotesServiceClient(connection)
	added, err := items.AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{NotetypeId: typeID.Ntid, Fields: []string{"Front", "Back"}}})
	if err != nil {
		t.Fatal(err)
	}
	beforeNote, err := items.GetNote(ctx, &notes.NoteId{Nid: added.NoteId})
	if err != nil {
		t.Fatal(err)
	}
	generated, err := items.CardsOfNote(ctx, &notes.NoteId{Nid: added.NoteId})
	if err != nil || len(generated.GetCids()) != 1 {
		t.Fatalf("fixture cards: %v", err)
	}
	cardService := cards.NewCardsServiceClient(connection)
	beforeCard, err := cardService.GetCard(ctx, &cards.CardId{Cid: generated.Cids[0]})
	if err != nil {
		t.Fatal(err)
	}
	beforeType, err := types.GetNotetype(ctx, typeID)
	if err != nil {
		t.Fatal(err)
	}
	root := writeConfigs(t)
	configs, err := config.New(root).Load()
	if err != nil {
		t.Fatal(err)
	}
	store := state.New(connection, "")
	empty, err := store.Dump(ctx)
	if err != nil || !proto.Equal(empty, &collectionpb.State{Version: 1}) {
		t.Fatalf("empty state: %v %v", empty, err)
	}
	bindings := map[string]string{"decks/Default/deck.anki.toml": "1", "notes/example.note.anki.toml": fmt.Sprint(added.NoteId), "note_types/Basic.notetype.anki.toml": fmt.Sprint(typeID.Ntid)}
	if changed, err := store.Import(ctx, configs, bindings); err != nil || !changed {
		t.Fatalf("batched import: %v %v", changed, err)
	}
	stored, err := store.Dump(ctx)
	if err != nil || stored.GetResources().GetNoteTypes()["note_types/Basic.notetype.anki.toml"].GetFields()["Front"] == nil {
		t.Fatalf("native members were not bound: %v", err)
	}
	checkStatePreview(t, ctx, store, stored, beforeType)
	if changed, err := store.Import(ctx, configs, bindings); err != nil || changed {
		t.Fatalf("identical import: %v %v", changed, err)
	}
	for _, invalid := range []map[string]string{
		{"notes/missing.note.anki.toml": "opaque"},
		{"presets/new": "1", "presets/missing": "9223372036854775807"},
		{"notes/other.note.anki.toml": fmt.Sprint(added.NoteId)},
	} {
		if changed, err := store.Import(ctx, configs, invalid); err == nil || changed {
			t.Fatal("invalid batch accepted")
		}
	}
	patch := &collectionpb.State{Version: 1, Resources: &collectionpb.Resources{Presets: map[string]*collectionpb.PresetState{"added": {Id: 1}}}}
	patch.Resources.NoteTypes = map[string]*collectionpb.NoteTypeState{
		"note_types/Basic.notetype.anki.toml": {
			Id:     typeID.Ntid,
			Fields: map[string]*collectionpb.FieldState{"Front": proto.Clone(stored.Resources.NoteTypes["note_types/Basic.notetype.anki.toml"].Fields["Front"]).(*collectionpb.FieldState)},
		},
	}
	originalPatch := proto.Clone(patch)
	if changed, err := store.Update(ctx, patch); err != nil || !changed || !proto.Equal(patch, originalPatch) {
		t.Fatalf("update or input ownership: %v %v", changed, err)
	}
	stored.Resources.Presets["added"] = &collectionpb.PresetState{Id: 1}
	stored.Resources.NoteTypes["note_types/Basic.notetype.anki.toml"] = proto.Clone(patch.Resources.NoteTypes["note_types/Basic.notetype.anki.toml"]).(*collectionpb.NoteTypeState)
	if changed, err := store.Update(ctx, patch); err != nil || changed {
		t.Fatalf("identical update: %v %v", changed, err)
	}
	conflict := &collectionpb.State{Version: 1, Resources: &collectionpb.Resources{Notes: map[string]*collectionpb.NoteState{"notes/other.note.anki.toml": {Id: added.NoteId}}}}
	if changed, err := store.Update(ctx, conflict); err == nil || changed {
		t.Fatal("merged identity conflict accepted")
	}
	independent := state.New(connection, "anki_as_code/state/independent")
	if changed, err := independent.Replace(ctx, patch); err != nil || !changed {
		t.Fatalf("independent namespace: %v %v", changed, err)
	}
	for _, format := range []codec.Format{codec.TOML, codec.JSON} {
		dumped := runStateCLI(t, connection, open, "", true, nil, "dump", "--format", string(format))
		value := &collectionpb.State{}
		if err := codec.Decode(dumped, format, value); err != nil || !proto.Equal(value, stored) {
			t.Fatalf("CLI dump %s: %v", format, err)
		}
		runStateCLI(t, connection, open, "", true, dumped, "replace", "--format", string(format), "-")
		file := filepath.Join(t.TempDir(), "state.anki."+string(format))
		if err := os.WriteFile(file, dumped, 0o600); err != nil {
			t.Fatal(err)
		}
		runStateCLI(t, connection, open, "", true, nil, "update", "--format", string(format), file)
		if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
			if err := os.WriteFile(filepath.Join(output, "state.anki."+string(format)), dumped, 0o600); err != nil {
				t.Fatal(err)
			}
		}
	}
	runStateCLI(t, connection, open, "", true, nil, "import", "--text", root, "presets/cli=1", "decks/Default/deck.anki.toml=1")
	stored.Resources.Presets["cli"] = &collectionpb.PresetState{Id: 1}
	message := runStateCLI(t, connection, open, "", false, nil, "import", "--text", root, "notes/missing.note.anki.toml=opaque")
	if !strings.Contains(string(message), "is not declared") {
		t.Fatal("identifier parsed before resolving resource")
	}
	runStateCLI(t, connection, open, "", false, nil, "import", "--text", root, "presets/partial=1", "presets/missing=9223372036854775807")
	runStateCLI(t, connection, open, "", false, []byte("version = ["), "replace", "-")
	runStateCLI(t, connection, open, "", false, nil, "dump", "--format", "yaml")
	runStateCLI(t, connection, open, "anki_as_code/state/independent", true, []byte("version = '1'\n"), "replace", "-")
	if value, err := independent.Dump(ctx); err != nil || !proto.Equal(value, &collectionpb.State{Version: 1}) {
		t.Fatal("replace retained omitted bindings")
	}
	checkStateRecovery(t, ctx, connection, open)
	if value, err := store.Dump(ctx); err != nil || !proto.Equal(value, stored) {
		t.Fatal("rejected operations or another namespace changed bindings")
	}
	afterNote, noteErr := items.GetNote(ctx, &notes.NoteId{Nid: added.NoteId})
	afterCard, cardErr := cardService.GetCard(ctx, &cards.CardId{Cid: generated.Cids[0]})
	afterType, typeErr := types.GetNotetype(ctx, typeID)
	if noteErr != nil || cardErr != nil || typeErr != nil || !proto.Equal(beforeNote, afterNote) || !proto.Equal(beforeCard, afterCard) || !proto.Equal(beforeType, afterType) {
		t.Fatal("state operations changed native objects or closed the caller's collection")
	}
	if _, err := collections.CloseCollection(ctx, &collection.CloseCollectionRequest{}); err != nil {
		t.Fatal(err)
	}
	opened = false
	if _, err := store.Dump(ctx); err == nil {
		t.Fatal("dump accepted an unprepared collection")
	}
	if changed, err := store.Replace(ctx, &collectionpb.State{Version: 1}); err == nil || changed {
		t.Fatal("replacement accepted an unprepared collection")
	}
	if _, err := collections.OpenCollection(ctx, open); err != nil {
		t.Fatal(err)
	}
	opened = true
	if value, err := store.Dump(ctx); err != nil || !proto.Equal(value, stored) {
		t.Fatal("state did not persist after the caller reopened the collection")
	}
	if _, err := collections.CloseCollection(ctx, &collection.CloseCollectionRequest{}); err != nil {
		t.Fatal(err)
	}
	opened = false
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := os.ReadFile(database)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "state.anki2"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// checkStatePreview verifies detached protobuf inputs and moves against persisted native state.
func checkStatePreview(t *testing.T, ctx context.Context, store *state.Store, stored *collectionpb.State, definition *notetypes.Notetype) {
	t.Helper()
	const typePath = "note_types/Basic.notetype.anki.toml"
	configs := &collectionpb.CollectionConfigSet{Config: &collectionpb.CollectionConfig{Version: 1, Moved: []*collectionpb.ResourceMove{
		{From: "decks/Default", To: "decks/Renamed"},
		{From: typePath + "#fields/Front", To: typePath + "#fields/Question"},
	}}}
	originalConfigs := proto.Clone(configs)
	preview, err := store.Preview(ctx, configs, nil)
	if err != nil || preview.GetResources().GetDecks()["decks/Renamed/deck.anki.toml"].GetId() != 1 || preview.GetResources().GetNoteTypes()[typePath].GetFields()["Question"] == nil {
		t.Fatalf("preview moved bindings: %v %v", preview, err)
	}
	if preview.Resources.Decks["decks/Default/deck.anki.toml"] != nil || preview.Resources.NoteTypes[typePath].Fields["Front"] != nil || !proto.Equal(configs, originalConfigs) {
		t.Fatal("preview retained move sources or changed input configs")
	}
	preview.Resources.Decks["decks/Renamed/deck.anki.toml"].Id = 999
	if value, err := store.Dump(ctx); err != nil || !proto.Equal(value, stored) {
		t.Fatal("preview or its returned snapshot changed persisted state")
	}
	inputs := proto.Clone(stored.Resources).(*collectionpb.Resources)
	inputs.Notes = map[string]*collectionpb.NoteState{"notes/new.note.anki.toml": {}, "notes/another.note.anki.toml": {}}
	inputs.Presets = map[string]*collectionpb.PresetState{"existing": {Id: 1}}
	originalInputs := proto.Clone(inputs)
	preview, err = store.Preview(ctx, nil, inputs)
	if err != nil || !proto.Equal(preview.GetResources(), inputs) || preview.Resources.NoteTypes[typePath].Fields["Back"].Ordinal != 1 || preview.Resources.NoteTypes[typePath].Fields["Front"].Id != definition.Fields[0].Config.GetId() {
		t.Fatalf("preview protobuf bindings: %v", err)
	}
	if changed, err := store.Replace(ctx, preview); err == nil || changed {
		t.Fatal("persisted unresolved preview identities")
	}
	preview.Resources.NoteTypes[typePath].Fields["Front"].Id = 0
	preview.Resources.Notes["notes/new.note.anki.toml"].Id = 999
	if !proto.Equal(inputs, originalInputs) {
		t.Fatal("preview aliased or changed its input resources")
	}
	retained, err := store.Preview(ctx, nil, nil)
	if err != nil || !proto.Equal(retained, stored) {
		t.Fatalf("nil resources did not retain stored bindings: %v", err)
	}
	cleared, err := store.Preview(ctx, nil, &collectionpb.Resources{})
	if err != nil || !proto.Equal(cleared.GetResources(), &collectionpb.Resources{}) {
		t.Fatalf("empty resources did not clear all groups: %v", err)
	}
	for _, invalid := range []*collectionpb.Resources{
		{Notes: map[string]*collectionpb.NoteState{"../outside": {Id: 1}}},
		{Notes: map[string]*collectionpb.NoteState{"notes/one": {Id: 1}, "notes/two": {Id: 1}}},
		{Decks: map[string]*collectionpb.DeckState{"decks/invalid/deck.anki.toml": nil}},
		{Notes: map[string]*collectionpb.NoteState{"notes/invalid.note.anki.toml": nil}},
		{NoteTypes: map[string]*collectionpb.NoteTypeState{typePath: nil}},
		{NoteTypes: map[string]*collectionpb.NoteTypeState{typePath: {Id: definition.Id, Fields: map[string]*collectionpb.FieldState{"Front": nil}}}},
		{NoteTypes: map[string]*collectionpb.NoteTypeState{typePath: {Id: definition.Id, Templates: map[string]*collectionpb.TemplateState{"Card 1": nil}}}},
		{Presets: map[string]*collectionpb.PresetState{"invalid": {Id: 0}}},
	} {
		if _, err := store.Preview(ctx, nil, invalid); err == nil {
			t.Fatal("invalid preview bindings accepted")
		}
	}
	for _, moves := range [][]*collectionpb.ResourceMove{
		{{From: "decks/Default/deck.anki.toml", To: "notes/example.note.anki.toml"}},
		{{From: "decks/Default", To: "decks/Other"}, {From: "decks/Other", To: "decks/Third"}},
	} {
		invalid := &collectionpb.CollectionConfigSet{Config: &collectionpb.CollectionConfig{Version: 1, Moved: moves}}
		if _, err := store.Preview(ctx, invalid, nil); err == nil {
			t.Fatal("conflicting preview moves accepted")
		}
	}
	configs.Decks = []*collectionpb.DeckConfigSet{{Path: "decks/Default/deck.anki.toml"}}
	if _, err := store.Preview(ctx, configs, nil); err == nil {
		t.Fatal("move source still declared was accepted")
	}
	if changed, err := store.Replace(ctx, cleared); err != nil || !changed {
		t.Fatalf("persist resolved preview: %v %v", changed, err)
	}
	if changed, err := store.Replace(ctx, stored); err != nil || !changed {
		t.Fatalf("restore original bindings: %v %v", changed, err)
	}
	if value, err := store.Dump(ctx); err != nil || !proto.Equal(value, stored) {
		t.Fatal("preview failure or rejected write changed persisted bindings")
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := codec.Encode(cleared, codec.TOML)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "state-preview.anki.toml"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// checkStateRecovery repairs invalid native metadata through both replacement APIs.
func checkStateRecovery(t *testing.T, ctx context.Context, connection grpc.ClientConnInterface, open *collection.OpenCollectionRequest) {
	t.Helper()
	const key = "anki_as_code/state/recovery"
	configs := ankiconfig.NewConfigServiceClient(connection)
	store := state.New(connection, key)
	replacement := &collectionpb.State{Version: 1}
	for _, saved := range []string{
		`{"version":"1","resources":{"notes":{"notes/broken.note.anki.toml":{"id":"0"}}}}`,
		`{"version":"99"}`,
		`{"version":"1","unknown":true}`,
	} {
		request := &ankiconfig.SetConfigJsonRequest{Key: key, ValueJson: []byte(saved)}
		if _, err := configs.SetConfigJson(ctx, request); err != nil {
			t.Fatal(err)
		}
		if _, err := store.Dump(ctx); err == nil {
			t.Fatal("invalid saved state was accepted")
		}
		before, err := configs.GetAllConfig(ctx, &generic.Empty{})
		if err != nil {
			t.Fatal(err)
		}
		if changed, err := store.Replace(ctx, &collectionpb.State{Version: 2}); err == nil || changed {
			t.Fatal("invalid replacement accepted")
		}
		after, err := configs.GetAllConfig(ctx, &generic.Empty{})
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(decodeConfigJSON(t, before.Json), decodeConfigJSON(t, after.Json)) {
			t.Fatal("invalid replacement changed collection config")
		}
		if changed, err := store.Replace(ctx, replacement); err != nil || !changed {
			t.Fatalf("replace invalid saved state %s: changed=%v error=%v", saved, changed, err)
		}
		if changed, err := store.Replace(ctx, replacement); err != nil || changed {
			t.Fatalf("identical recovered replacement: changed=%v error=%v", changed, err)
		}
		if _, err := configs.SetConfigJson(ctx, request); err != nil {
			t.Fatal(err)
		}
		runStateCLI(t, connection, open, key, true, []byte("version = '1'\n"), "replace", "-")
		if value, err := store.Dump(ctx); err != nil || !proto.Equal(value, replacement) {
			t.Fatalf("CLI did not repair saved state: %v", err)
		}
	}
}

// decodeConfigJSON preserves native numbers while ignoring JSON object-key ordering.
func decodeConfigJSON(t *testing.T, data []byte) any {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		t.Fatal(err)
	}
	return value
}

// startBackend owns one managed native worker and lends its connection to the fixture.
func startBackend(t *testing.T, ctx context.Context) grpc.ClientConnInterface {
	t.Helper()
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	worker := backend.New(backend.BackendSTDOptions{Binary: binary})
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := worker.Stop(cleanup); err != nil {
			t.Error(err)
		}
	})
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	connection, err := worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	return connection
}

// runStateCLI lends a closed database to the managed CLI, then reopens it for API verification.
func runStateCLI(t *testing.T, connection grpc.ClientConnInterface, open *collection.OpenCollectionRequest, key string, success bool, input []byte, args ...string) []byte {
	t.Helper()
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	collections := collection.NewBackendCollectionServiceClient(connection)
	if _, err := collections.CloseCollection(ctx, &collection.CloseCollectionRequest{}); err != nil {
		t.Fatal(err)
	}
	defer func() {
		cleanup, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		if _, err := collections.OpenCollection(cleanup, open); err != nil {
			t.Error(err)
		}
	}()
	arguments := append([]string{"col", "state", "--collection", open.CollectionPath, "--state-key", key}, args...)
	command := exec.CommandContext(ctx, binary, arguments...)
	command.Env = append(os.Environ(), "PATH=")
	command.Stdin = bytes.NewReader(input)
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	err := command.Run()
	if (err == nil) != success {
		t.Fatalf("state %v: error=%v stderr=%s", args, err, stderr.String())
	}
	if !success {
		if stdout.Len() != 0 || stderr.Len() == 0 {
			t.Fatal("failed command polluted stdout or omitted its error")
		}
		return stderr.Bytes()
	}
	if args[0] != "dump" && stdout.Len() != 0 {
		t.Fatal("state status polluted stdout")
	}
	return stdout.Bytes()
}

// writeConfigs declares importable paths without relying on collection export.
func writeConfigs(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for name, content := range map[string]string{
		"collection.anki.toml":                "version = 1\n",
		"decks/Default/deck.anki.toml":        "preset = 'Default'\n[normal]\n",
		"note_types/Basic.notetype.anki.toml": "name = 'Basic'\n",
		"notes/example.note.anki.toml":        "note_type = 'note_types/Basic.notetype.anki.toml'\n",
		"notes/other.note.anki.toml":          "note_type = 'note_types/Basic.notetype.anki.toml'\n",
	} {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
