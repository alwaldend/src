package ankiascode

import (
	"context"
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
	"github.com/ankitects/anki/proto/anki/notes"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// TestBorrowedCollectionWorkflow verifies reuse, independent sources and caller-owned publication.
func TestBorrowedCollectionWorkflow(t *testing.T) {
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
	firstSource := ankiascodeapi.CollectionConnection(firstConnection)
	secondSource := ankiascodeapi.CollectionConnection(secondConnection)
	application := &App{}
	ctx := t.Context()
	const key = "borrowed/primary"
	root := t.TempDir()
	text, stateOutput := filepath.Join(root, "text"), filepath.Join(root, "state.anki.toml")
	original, err := os.ReadFile(fixture.base)
	if err != nil {
		t.Fatal(err)
	}
	if err := application.ColExport(ctx, firstSource, text, stateOutput, key); err != nil {
		t.Fatal(err)
	}
	exported := &ankiascodepb.State{}
	readFixtureTOML(t, stateOutput, exported)
	if err := application.ColStateUpdate(ctx, firstSource, key, exported); err != nil {
		t.Fatal(err)
	}
	notePath := fixtureNotePath(text)
	address, err := filepath.Rel(text, notePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := application.ColStateImport(ctx, firstSource, text, key, map[string]string{filepath.ToSlash(address): strconv.FormatInt(fixture.noteID, 10)}); err != nil {
		t.Fatal(err)
	}
	note := &ankiascodepb.NoteConfig{}
	readFixtureTOML(t, notePath, note)
	note.Tags = append(note.Tags, "borrowed-workflow")
	writeFixtureTOML(t, notePath, note)
	report, err := application.ColPlan(ctx, firstSource, text, key)
	if err != nil || len(report.GetChanges()) == 0 {
		t.Fatalf("borrowed plan did not find the edited note: %v %v", report, err)
	}
	if err := application.ColApply(ctx, firstSource, text, key); err != nil {
		t.Fatal(err)
	}
	report, err = application.ColPlan(ctx, firstSource, text, key)
	if err != nil || len(report.GetChanges()) != 0 {
		t.Fatalf("borrowed apply did not converge: %v %v", report, err)
	}
	if err := application.ColStateReplace(ctx, secondSource, key, &ankiascodepb.State{Version: 1}); err != nil {
		t.Fatal(err)
	}
	other, err := application.ColStateDump(ctx, secondSource, key)
	if err != nil || len(other.GetResources().GetNotes()) != 0 {
		t.Fatalf("second source inherited the first collection state: %v %v", other, err)
	}
	if err := application.ColStateReplace(ctx, firstSource, key, &ankiascodepb.State{Version: 2}); err == nil {
		t.Fatal("invalid borrowed state replacement succeeded")
	}
	cause := errors.New("borrowed operation cancelled")
	cancelled, cancel := context.WithCancelCause(ctx)
	cancel(cause)
	if _, err := application.ColStateDump(cancelled, firstSource, key); !errors.Is(err, context.Canceled) || !errors.Is(err, cause) {
		t.Fatalf("borrowed cancellation lost its cause: %v", err)
	}
	actual, err := application.ColStateDump(ctx, firstSource, key)
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

// TestCollectionSourceFailures rejects empty and ambiguous sources before opening or native reads.
func TestCollectionSourceFailures(t *testing.T) {
	client := connectFixture(t, filepath.Join(t.TempDir(), "collection.anki2"))
	connection, err := client.worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	application := &App{}
	var nilConnection *grpc.ClientConn
	var nilSource *conflictingCollectionSource
	for _, source := range []ankiascodeapi.CollectionSource{
		nil,
		ankiascodeapi.CollectionPath(""),
		ankiascodeapi.CollectionConnection(nil),
		ankiascodeapi.CollectionConnection(nilConnection),
		nilSource,
		conflictingCollectionSource{path: "unused.anki2", connection: connection},
	} {
		if value, err := application.ColStateDump(t.Context(), source, ""); err == nil || !strings.Contains(err.Error(), "exactly one") || value != nil {
			t.Fatalf("invalid source reached native access: state=%v error=%v", value, err)
		}
	}
	if _, err := application.ColStateDump(t.Context(), ankiascodeapi.CollectionConnection(connection), ""); err != nil {
		t.Fatalf("source validation affected caller-owned access: %v", err)
	}
}

// conflictingCollectionSource supplies both alternatives to exercise boundary validation.
type conflictingCollectionSource struct {
	// path names the invalid additional path alternative.
	path string
	// connection identifies an already-open caller-owned collection.
	connection grpc.ClientConnInterface
}

// Path returns the conflicting path alternative.
func (source conflictingCollectionSource) Path() string { return source.path }

// Connection returns the conflicting borrowed connection alternative.
func (source conflictingCollectionSource) Connection() grpc.ClientConnInterface {
	return source.connection
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
		Collection: ankiascodeapi.CollectionConnection(connection), WriteBack: true,
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
	request.Backend = "unix:///unused"
	if _, err := application.APICall(t.Context(), request); err == nil {
		t.Fatal("ambiguous raw API source accepted")
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
	request = ankiascodeapi.APIRequest{Collection: ankiascodeapi.CollectionPath(database), Method: "/anki.collection.BackendCollectionService/OpenCollection", Input: []byte("{}"), InputFormat: "json", OutputFormat: "json"}
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
