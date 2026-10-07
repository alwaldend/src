package reconcile_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/export"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/plan"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/reconcile"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// executionConnection forbids collection discovery while allowing native execution reads.
type executionConnection struct {
	// ClientConnInterface supplies the real native transport and streaming implementation.
	grpc.ClientConnInterface
}

// Invoke rejects snapshot discovery so execution cannot silently plan again.
func (connection executionConnection) Invoke(ctx context.Context, method string, input, output any, options ...grpc.CallOption) error {
	for _, suffix := range []string{"/SearchNotes", "/SearchCards", "/GetDeckNames", "/GetNotetypeNames"} {
		if strings.HasSuffix(method, suffix) {
			return fmt.Errorf("unexpected resource discovery during execution: %s", method)
		}
	}
	if err := connection.ClientConnInterface.Invoke(ctx, method, input, output, options...); err != nil {
		return fmt.Errorf("invoke native method %s: %w", method, err)
	}
	return nil
}

// TestSavedPlanWorkflow executes a detached protobuf plan without discovering resources.
// It covers dynamic deck/type/note IDs, metadata retention, stale state and convergence.
func TestSavedPlanWorkflow(t *testing.T) {
	ctx := t.Context()
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	worker := backend.New(backend.BackendSTDOptions{Binary: binary})
	collection := backend.NewCollection(backend.BackendCollectionOptions{Backend: worker, Path: filepath.Join(t.TempDir(), "collection.anki2")})
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := collection.Stop(cleanup); err != nil {
			t.Error(err)
		}
		if err := worker.Stop(cleanup); err != nil {
			t.Error(err)
		}
	})
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	if err := collection.Start(ctx); err != nil {
		t.Fatal(err)
	}
	connection, err := collection.Connection()
	if err != nil {
		t.Fatal(err)
	}
	resources, states := resource.New(connection), state.New(connection, "saved-plan")
	initial, err := resources.Snapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var basicID int64
	for _, definition := range initial.NoteTypes {
		if definition.Name == "Basic" {
			basicID = definition.Id
		}
	}
	if basicID == 0 {
		t.Fatal("native Basic type missing")
	}
	original, err := resources.Notes.NewNote(ctx, &notetypes.NotetypeId{Ntid: basicID})
	if err != nil {
		t.Fatal(err)
	}
	original.Fields = []string{"Before", "Answer"}
	added, err := resources.Notes.AddNote(ctx, &notes.AddNoteRequest{Note: original, DeckId: 1})
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(t.TempDir(), "configs")
	if err := export.New(resources, states).Export(ctx, root); err != nil {
		t.Fatal(err)
	}
	configs, err := config.New(root).Load()
	if err != nil {
		t.Fatal(err)
	}
	initial, err = resources.Snapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	bindings := map[string]string{}
	for _, group := range configs.Decks {
		for _, native := range initial.Decks {
			if native.Name == group.Title {
				bindings[group.Path] = fmt.Sprint(native.Id)
			}
		}
		for _, file := range group.Notes {
			bindings[file.Path] = fmt.Sprint(added.NoteId)
		}
	}
	var definition *collectionpb.NoteTypeConfig
	for _, file := range configs.NoteTypes {
		for _, native := range initial.NoteTypes {
			if native.Name == file.Config.Name {
				bindings[file.Path] = fmt.Sprint(native.Id)
			}
		}
		if file.Config.Name == "Basic" {
			definition = proto.Clone(file.Config).(*collectionpb.NoteTypeConfig)
		}
	}
	if _, err := states.Import(ctx, configs, bindings); err != nil {
		t.Fatal(err)
	}
	stored, err := states.Dump(ctx)
	if err != nil {
		t.Fatal(err)
	}
	configs.Decks[0].Notes[0].Config.Fields["Front"] = "After"
	deckPath, typePath, notePath := "decks/Created/deck.anki.toml", "note_types/Created.notetype.anki.toml", "decks/Created/notes/Added.note.anki.toml"
	definition.Name = "Created"
	definition.Templates[0].TargetDeck = deckPath
	configs.NoteTypes = append(configs.NoteTypes, &collectionpb.NoteTypeConfigFile{Path: typePath, Config: definition})
	configs.Decks = append(configs.Decks, &collectionpb.DeckConfigSet{
		Path: deckPath, Title: "Created", Deck: &collectionpb.DeckConfig{Normal: &decks.Deck_Normal{}},
		Notes: []*collectionpb.NoteConfigFile{{Path: notePath, Config: &collectionpb.NoteConfig{NoteType: typePath, Fields: map[string]string{"Front": "Added", "Back": "Answer"}}}},
	})
	source := proto.Clone(configs)
	compiled, err := plan.Build(configs, initial, stored)
	if err != nil {
		t.Fatal(err)
	}
	if compiled.ProposedBindings.Decks[deckPath].Id != 0 || compiled.ProposedBindings.NoteTypes[typePath].Id != 0 || compiled.Notes[len(compiled.Notes)-1].NoteType == "" {
		t.Fatal("new objects use invented identities instead of path references")
	}
	saved, err := codec.Encode(compiled, codec.JSON)
	if err != nil {
		t.Fatal(err)
	}
	detached := &collectionpb.CollectionPlan{}
	if err := codec.Decode(saved, codec.JSON, detached); err != nil {
		t.Fatal(err)
	}
	frozen := proto.Clone(detached)
	execution := executionConnection{ClientConnInterface: connection}
	engine := reconcile.New(resource.New(execution), state.New(execution, "saved-plan"))
	changed, err := engine.Apply(ctx, detached)
	if err != nil || !changed {
		t.Fatalf("execute detached plan: changed=%v error=%v", changed, err)
	}
	if !proto.Equal(detached, frozen) || !proto.Equal(configs, source) {
		t.Fatal("execution changed its plan or source configs")
	}
	recorded, err := states.Dump(ctx)
	if err != nil {
		t.Fatal(err)
	}
	current, err := resources.Notes.GetNote(ctx, &notes.NoteId{Nid: added.NoteId})
	if err != nil {
		t.Fatal(err)
	}
	if current.Guid != original.Guid || current.Fields[0] != "After" {
		t.Fatal("retained note lost identity or content")
	}
	noteID, typeID, deckID := recorded.Resources.Notes[notePath].Id, recorded.Resources.NoteTypes[typePath].Id, recorded.Resources.Decks[deckPath].Id
	created, err := resources.Notes.GetNote(ctx, &notes.NoteId{Nid: noteID})
	if err != nil {
		t.Fatal(err)
	}
	if created.NotetypeId != typeID || created.Guid == "" || created.Fields[0] != "Added" {
		t.Fatal("created note references unresolved identities")
	}
	nativeType, err := resources.NoteTypes.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: typeID})
	if err != nil {
		t.Fatal(err)
	}
	if nativeType.Templates[0].Config.TargetDeckId != deckID {
		t.Fatal("template retained a provisional deck identity")
	}
	snapshot, err := resources.Snapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	converged, err := plan.Build(configs, snapshot, recorded)
	if err != nil || len(converged.GetReport().GetChanges()) != 0 || len(converged.GetCardAssignments()) != 0 {
		t.Fatalf("executed plan did not converge: %v %v", converged, err)
	}
	before, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := engine.Apply(ctx, detached); err == nil {
		t.Fatal("stale saved plan executed again")
	}
	snapshot, err = resources.Snapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	after, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("stale plan rejection changed native resources")
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "execution-plan.json"), saved, 0o600); err != nil {
			t.Fatal(err)
		}
		data, err := codec.Encode(recorded, codec.JSON)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "executed-state.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: filepath.Join(output, "executed-plan.colpkg")}); err != nil {
			t.Fatal(err)
		}
	}
}
