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

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
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
	configs, err := config.New(root, filesystem.FS{}).Load()
	if err != nil {
		t.Fatal(err)
	}
	store := state.New(connection, "")
	empty, err := store.Dump(ctx)
	if err != nil || !proto.Equal(empty, &ankiascodepb.State{Version: 1}) {
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
	patch := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"added": {Id: 1}}}}
	patch.Resources.NoteTypes = map[string]*ankiascodepb.NoteTypeState{
		"note_types/Basic.notetype.anki.toml": {
			Id:     typeID.Ntid,
			Fields: map[string]*ankiascodepb.FieldState{"Front": proto.Clone(stored.Resources.NoteTypes["note_types/Basic.notetype.anki.toml"].Fields["Front"]).(*ankiascodepb.FieldState)},
		},
	}
	originalPatch := proto.Clone(patch)
	if changed, err := store.Update(ctx, patch); err != nil || !changed || !proto.Equal(patch, originalPatch) {
		t.Fatalf("update or input ownership: %v %v", changed, err)
	}
	stored.Resources.Presets["added"] = &ankiascodepb.PresetState{Id: 1}
	stored.Resources.NoteTypes["note_types/Basic.notetype.anki.toml"] = proto.Clone(patch.Resources.NoteTypes["note_types/Basic.notetype.anki.toml"]).(*ankiascodepb.NoteTypeState)
	if changed, err := store.Update(ctx, patch); err != nil || changed {
		t.Fatalf("identical update: %v %v", changed, err)
	}
	conflict := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Notes: map[string]*ankiascodepb.NoteState{"notes/other.note.anki.toml": {Id: added.NoteId}}}}
	if changed, err := store.Update(ctx, conflict); err == nil || changed {
		t.Fatal("merged identity conflict accepted")
	}
	independent := state.New(connection, "anki_as_code/state/independent")
	if changed, err := independent.Replace(ctx, patch); err != nil || !changed {
		t.Fatalf("independent namespace: %v %v", changed, err)
	}
	for _, format := range []codec.Format{codec.TOML, codec.JSON} {
		dumped := runStateCLI(t, connection, open, "", true, nil, "dump", "--format", string(format))
		value := &ankiascodepb.State{}
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
	stored.Resources.Presets["cli"] = &ankiascodepb.PresetState{Id: 1}
	message := runStateCLI(t, connection, open, "", false, nil, "import", "--text", root, "notes/missing.note.anki.toml=opaque")
	if !strings.Contains(string(message), "is not declared") {
		t.Fatal("identifier parsed before resolving resource")
	}
	runStateCLI(t, connection, open, "", false, nil, "import", "--text", root, "presets/partial=1", "presets/missing=9223372036854775807")
	runStateCLI(t, connection, open, "", false, []byte("version = ["), "replace", "-")
	runStateCLI(t, connection, open, "", false, nil, "dump", "--format", "yaml")
	runStateCLI(t, connection, open, "anki_as_code/state/independent", true, []byte("version = '1'\n"), "replace", "-")
	if value, err := independent.Dump(ctx); err != nil || !proto.Equal(value, &ankiascodepb.State{Version: 1}) {
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
	if changed, err := store.Replace(ctx, &ankiascodepb.State{Version: 1}); err == nil || changed {
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

// checkStateRecovery repairs invalid native metadata through both replacement APIs.
func checkStateRecovery(t *testing.T, ctx context.Context, connection grpc.ClientConnInterface, open *collection.OpenCollectionRequest) {
	t.Helper()
	const key = "anki_as_code/state/recovery"
	configs := ankiconfig.NewConfigServiceClient(connection)
	store := state.New(connection, key)
	replacement := &ankiascodepb.State{Version: 1}
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
		if changed, err := store.Replace(ctx, &ankiascodepb.State{Version: 2}); err == nil || changed {
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
	worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: binary, Arguments: []string{"serve"}})
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
