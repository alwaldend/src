package ankiascode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// TestBasePathCollectionExport publishes configuration and state through rooted directory and dangling aliases.
func TestBasePathCollectionExport(t *testing.T) {
	fixture := newFixture(t)
	client := inspectFixture(t, fixture.base)
	connection, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	for _, namespace := range []string{"rooted-absolute", "rooted-relative", "nested-absolute", "nested-relative"} {
		t.Run(namespace, func(t *testing.T) {
			storage := afero.NewBasePathFs(afero.NewOsFs(), t.TempDir())
			if strings.HasPrefix(namespace, "nested-") {
				storage = afero.NewBasePathFs(storage, "scope")
			}
			if err := storage.MkdirAll("target", 0o755); err != nil {
				t.Fatal(err)
			}
			linker := storage.(afero.Linker)
			if err := linker.SymlinkIfPossible("target", "alias"); err != nil {
				t.Fatal(err)
			}
			if err := linker.SymlinkIfPossible("target/state.toml", "state-link.toml"); err != nil {
				t.Fatal(err)
			}
			prefix := ""
			if strings.HasSuffix(namespace, "-absolute") {
				prefix = string(filepath.Separator)
			}
			application := New(Options{FS: storage})
			if err := application.ColExport(t.Context(), ankiascodeapi.CollectionOptions{Connection: connection}, prefix+"alias/config", prefix+"state-link.toml"); err != nil {
				t.Fatal(err)
			}
			if _, err := application.ConfigDump(t.Context(), prefix+"target/config", "json"); err != nil {
				t.Fatal(err)
			}
			if _, err := application.ConfigDump(t.Context(), prefix+"alias/config", "json"); err == nil {
				t.Fatal("config loading accepted a symlink ancestor")
			}
			data, err := afero.ReadFile(storage, "target/state.toml")
			if err != nil {
				t.Fatal(err)
			}
			bindings := &ankiascodepb.State{}
			if err := codec.Decode(data, codec.TOML, bindings); err != nil || len(bindings.GetResources().GetNotes()) == 0 {
				t.Fatalf("rooted state output lost bindings: %v", err)
			}
			if info, _, err := storage.(afero.Lstater).LstatIfPossible("state-link.toml"); err != nil || info.Mode()&os.ModeSymlink == 0 {
				t.Fatalf("dangling output alias was replaced: %v", err)
			}
			entries, err := afero.ReadDir(storage, "target")
			if err != nil || len(entries) != 2 || entries[0].Name() != "config" || entries[1].Name() != "state.toml" {
				t.Fatalf("state publication left temporary output: %v %v", entries, err)
			}
			if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
				if err := os.WriteFile(filepath.Join(artifacts, "basepath-export-state-"+namespace+".toml"), data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestInjectedCollectionConfiguration keeps every config and state-output operation in memory.
func TestInjectedCollectionConfiguration(t *testing.T) {
	fixture := newFixture(t)
	client := inspectFixture(t, fixture.base)
	connection, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	storage := afero.NewMemMapFs()
	application := New(Options{FS: storage})
	options := ankiascodeapi.CollectionOptions{Connection: connection, StateKey: "injected-storage"}
	root := filepath.Join(t.TempDir(), "memory-only")
	text, statePath := filepath.Join(root, "text"), filepath.Join(root, "state.toml")
	if err := application.ColExport(t.Context(), options, text, statePath); err != nil {
		t.Fatal(err)
	}
	data, err := afero.ReadFile(storage, statePath)
	if err != nil {
		t.Fatal(err)
	}
	bindings := &ankiascodepb.State{}
	if err := codec.Decode(data, codec.TOML, bindings); err != nil {
		t.Fatal(err)
	}
	if err := application.ColStateUpdate(t.Context(), options, bindings); err != nil {
		t.Fatal(err)
	}
	path := fixtureNotePath(text)
	data, err = afero.ReadFile(storage, path)
	if err != nil {
		t.Fatal(err)
	}
	note := &ankiascodepb.NoteConfig{}
	if err := codec.Decode(data, codec.TOML, note); err != nil {
		t.Fatal(err)
	}
	note.Tags = append(note.Tags, "injected-storage")
	data, err = codec.Encode(note, codec.TOML)
	if err != nil {
		t.Fatal(err)
	}
	if err := afero.WriteFile(storage, path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	address, err := filepath.Rel(text, path)
	if err != nil {
		t.Fatal(err)
	}
	if err := application.ColStateImport(t.Context(), options, text, map[string]string{filepath.ToSlash(address): strconv.FormatInt(fixture.noteID, 10)}); err != nil {
		t.Fatal(err)
	}
	plan, err := application.ColPlan(t.Context(), options, text)
	if err != nil || len(plan.GetChanges()) == 0 {
		t.Fatalf("memory plan: %v %v", plan, err)
	}
	if err := application.ColApply(t.Context(), options, text); err != nil {
		t.Fatal(err)
	}
	plan, err = application.ColPlan(t.Context(), options, text)
	if err != nil || len(plan.GetChanges()) != 0 {
		t.Fatalf("memory apply did not converge: %v %v", plan, err)
	}
	if _, err := os.Stat(root); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("collection configuration escaped to disk: %v", err)
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := os.WriteFile(filepath.Join(artifacts, "memory-applied.note.anki.toml"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestBorrowedCollectionWorkflow verifies reuse, independent sources and caller-owned publication.
func TestBorrowedCollectionWorkflow(t *testing.T) {
	const key = "borrowed/primary"
	fixture := newFixture(t)
	first, second := inspectFixture(t, fixture.base), inspectFixture(t, fixture.base)
	firstConnection, err := first.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	secondConnection, err := second.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	firstOptions := ankiascodeapi.CollectionOptions{Connection: firstConnection, StateKey: key}
	secondOptions := ankiascodeapi.CollectionOptions{Connection: secondConnection, StateKey: key}
	application := &App{}
	ctx := t.Context()
	root := t.TempDir()
	text, stateOutput := filepath.Join(root, "text"), filepath.Join(root, "state.anki.toml")
	original, err := os.ReadFile(fixture.base)
	if err != nil {
		t.Fatal(err)
	}
	if err := application.ColExport(ctx, firstOptions, text, stateOutput); err != nil {
		t.Fatal(err)
	}
	exported := &ankiascodepb.State{}
	readFixtureTOML(t, stateOutput, exported)
	if err := application.ColStateUpdate(ctx, firstOptions, exported); err != nil {
		t.Fatal(err)
	}
	notePath := fixtureNotePath(text)
	address, err := filepath.Rel(text, notePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := application.ColStateImport(ctx, firstOptions, text, map[string]string{filepath.ToSlash(address): strconv.FormatInt(fixture.noteID, 10)}); err != nil {
		t.Fatal(err)
	}
	note := &ankiascodepb.NoteConfig{}
	readFixtureTOML(t, notePath, note)
	note.Tags = append(note.Tags, "borrowed-workflow")
	writeFixtureTOML(t, notePath, note)
	report, err := application.ColPlan(ctx, firstOptions, text)
	if err != nil || len(report.GetChanges()) == 0 {
		t.Fatalf("borrowed plan did not find the edited note: %v %v", report, err)
	}
	if err := application.ColApply(ctx, firstOptions, text); err != nil {
		t.Fatal(err)
	}
	report, err = application.ColPlan(ctx, firstOptions, text)
	if err != nil || len(report.GetChanges()) != 0 {
		t.Fatalf("borrowed apply did not converge: %v %v", report, err)
	}
	if err := application.ColStateReplace(ctx, secondOptions, &ankiascodepb.State{Version: 1}); err != nil {
		t.Fatal(err)
	}
	other, err := application.ColStateDump(ctx, secondOptions)
	if err != nil || len(other.GetResources().GetNotes()) != 0 {
		t.Fatalf("second source inherited the first collection state: %v %v", other, err)
	}
	if err := application.ColStateReplace(ctx, firstOptions, &ankiascodepb.State{Version: 2}); err == nil {
		t.Fatal("invalid borrowed state replacement succeeded")
	}
	cause := errors.New("borrowed operation cancelled")
	cancelled, cancel := context.WithCancelCause(ctx)
	cancel(cause)
	if _, err := application.ColStateDump(cancelled, firstOptions); !errors.Is(err, context.Canceled) || !errors.Is(err, cause) {
		t.Fatalf("borrowed cancellation lost its cause: %v", err)
	}
	actual, err := application.ColStateDump(ctx, firstOptions)
	if err != nil || !proto.Equal(actual, exported) {
		t.Fatalf("borrowed source changed or closed after another source, failure or cancellation: %v %v", actual, err)
	}
	otherNote, err := second.notes.GetNote(ctx, &notes.NoteId{Nid: fixture.noteID})
	if err != nil || slices.Contains(otherNote.GetTags(), "borrowed-workflow") {
		t.Fatalf("apply changed or closed the second collection: %v %v", otherNote, err)
	}
	after, err := os.ReadFile(fixture.base)
	if err != nil || string(after) != string(original) {
		t.Fatal("borrowed operations published the source archive")
	}
	output := filepath.Join(root, "borrowed.colpkg")
	first.exportFixture(t, output)
	restored := inspectFixture(t, output)
	persisted, err := restored.notes.GetNote(ctx, &notes.NoteId{Nid: fixture.noteID})
	if err != nil || !slices.Contains(persisted.GetTags(), "borrowed-workflow") {
		t.Fatalf("owner publication lost the applied note: %v %v", persisted, err)
	}
	assertFixtureMedia(t, restored)
	saveFixtureArtifact(t, output, "borrowed-collection.colpkg")
	saveFixtureArtifact(t, stateOutput, "borrowed-state.anki.toml")
}

// TestCollectionOptionsFailures rejects conflicting access fields before native access.
func TestCollectionOptionsFailures(t *testing.T) {
	root := t.TempDir()
	client := connectFixture(t, filepath.Join(root, "collection.anki2"))
	connection, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	application := fixtureApplication()
	unopened := filepath.Join(root, "unopened.anki2")
	if err := os.WriteFile(unopened, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	var nilConnection *grpc.ClientConn
	failures := map[string]string{}
	for name, options := range map[string]ankiascodeapi.CollectionOptions{
		"empty":                  {},
		"typed-nil":              {Connection: nilConnection},
		"path-and-connection":    {Path: unopened, Connection: connection},
		"address-and-connection": {Address: "unix:///unused", Connection: connection},
		"all-access-fields":      {Path: unopened, Address: "unix:///unused", Connection: connection},
		"address-without-path":   {Address: "unix:///unused"},
	} {
		value, err := application.ColStateDump(t.Context(), options)
		if err == nil || !strings.Contains(err.Error(), "exactly one") || value != nil {
			t.Fatalf("%s reached native access: state=%v error=%v", name, value, err)
		}
		failures[name] = err.Error()
	}
	output := filepath.Join(root, "unpublished.colpkg")
	if err := application.ColBuild(t.Context(), ankiascodeapi.CollectionOptions{Connection: connection}, "unused", output); err == nil || !strings.Contains(err.Error(), "path") {
		t.Fatalf("build accepted a borrowed connection: %v", err)
	} else {
		failures["build-borrowed-connection"] = err.Error()
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("rejected build published output: %v", err)
	}
	if info, err := os.Stat(unopened); err != nil || info.Size() != 0 {
		t.Fatalf("invalid options opened the collection: %v %v", info, err)
	}
	if _, err := application.ColStateDump(t.Context(), ankiascodeapi.CollectionOptions{Connection: connection}); err != nil {
		t.Fatalf("option validation changed caller-owned access: %v", err)
	}
	if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
		data, err := json.MarshalIndent(failures, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, "collection-options-validation.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestBorrowedAPIWorkflow invokes the generic application operation against caller-owned native resources.
func TestBorrowedAPIWorkflow(t *testing.T) {
	fixture := newFixture(t)
	client := inspectFixture(t, fixture.base)
	connection, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile(fixture.base)
	if err != nil {
		t.Fatal(err)
	}
	application := &App{}
	request := ankiascodeapi.APIRequest{
		Collection: ankiascodeapi.CollectionOptions{Connection: connection}, WriteBack: true,
		Method: "/anki.notes.NotesService/GetNote", Input: []byte(fmt.Sprintf(`{"nid":"%d"}`, fixture.noteID)), InputFormat: "json", OutputFormat: "json",
	}
	output, err := application.APICall(t.Context(), request)
	if err != nil {
		t.Fatal(err)
	}
	note := &notes.Note{}
	if err := protojson.Unmarshal(output, note); err != nil || note.GetId() != fixture.noteID {
		t.Fatalf("raw native response: %s %v", output, err)
	}
	request.Input = []byte(`{"unknown":true}`)
	if _, err := application.APICall(t.Context(), request); err == nil {
		t.Fatal("invalid raw input accepted")
	}
	request.Input = []byte(fmt.Sprintf(`{"nid":"%d"}`, fixture.noteID))
	request.Collection.Address = "unix:///unused"
	if _, err := application.APICall(t.Context(), request); err == nil || !strings.Contains(err.Error(), "exactly one") {
		t.Fatalf("ambiguous raw API options accepted: %v", err)
	}
	if _, err := client.notes.GetNote(t.Context(), &notes.NoteId{Nid: fixture.noteID}); err != nil {
		t.Fatalf("raw API closed caller connection: %v", err)
	}
	current, err := os.ReadFile(fixture.base)
	if err != nil || string(current) != string(original) {
		t.Fatalf("raw API published borrowed archive: %v", err)
	}
	database := filepath.Join(t.TempDir(), "unopened.anki2")
	if err := os.WriteFile(database, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	request = ankiascodeapi.APIRequest{Collection: ankiascodeapi.CollectionOptions{Path: database}, Method: "/anki.collection.BackendCollectionService/OpenCollection", Input: []byte("{}"), InputFormat: "json", OutputFormat: "json"}
	if _, err := application.APICall(t.Context(), request); err == nil || !strings.Contains(err.Error(), "collection lifetime") {
		t.Fatalf("managed lifetime bypass: %v", err)
	}
	if info, err := os.Stat(database); err != nil || info.Size() != 0 {
		t.Fatalf("rejected lifecycle call opened collection: %v %v", info, err)
	}
	if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
		if err := os.WriteFile(filepath.Join(destination, "borrowed-api-note.json"), output, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
