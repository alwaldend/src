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

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/export"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/plan"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/reconcile"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/generic"
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
	// reject injects a native-operation failure after earlier phases have executed.
	reject string
}

// Invoke rejects snapshot discovery so execution cannot silently plan again.
func (connection executionConnection) Invoke(ctx context.Context, method string, input, output any, options ...grpc.CallOption) error {
	if connection.reject != "" && strings.HasSuffix(method, connection.reject) {
		return fmt.Errorf("injected execution failure at %s", method)
	}
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
	directory := t.TempDir()
	worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: binary, Arguments: []string{"serve"}})
	var collections collection.BackendCollectionServiceClient
	opened := false
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if opened {
			if _, err := collections.CloseCollection(cleanup, &collection.CloseCollectionRequest{}); err != nil {
				t.Error(err)
			}
		}
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
	collections = collection.NewBackendCollectionServiceClient(connection)
	if _, err := collections.OpenCollection(ctx, &collection.OpenCollectionRequest{
		CollectionPath:  filepath.Join(directory, "collection.anki2"),
		MediaFolderPath: filepath.Join(directory, "collection.media"),
		MediaDbPath:     filepath.Join(directory, "collection.media.db2"),
	}); err != nil {
		t.Fatal(err)
	}
	opened = true
	resources, states := resource.New(connection), state.New(connection, "saved-plan")
	initial, err := resources.Snapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	emptyState, err := states.Dump(ctx)
	if err != nil {
		t.Fatal(err)
	}
	emptyPlan, err := plan.Build(&ankiascodepb.CollectionConfigSet{Config: &ankiascodepb.CollectionConfig{Version: 1}}, initial, emptyState)
	if err != nil {
		t.Fatal(err)
	}
	if changed, err := reconcile.New(resources, states).Apply(ctx, emptyPlan); err != nil || changed {
		t.Fatalf("empty plan changed native state: changed=%v error=%v", changed, err)
	}
	unchanged, err := states.Dump(ctx)
	if err != nil || !proto.Equal(unchanged, emptyState) {
		t.Fatalf("empty namespace was rewritten: %v %v", unchanged, err)
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
	spare, err := resources.Decks.AddDeck(ctx, &decks.Deck{Name: "Spare", Kind: &decks.Deck_Normal_{Normal: &decks.Deck_Normal{ConfigId: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	if err := export.New(resources, states, filesystem.FS{}).Export(ctx, root, ""); err != nil {
		t.Fatal(err)
	}
	configs, err := config.New(root, filesystem.FS{}).Load()
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
	var definition *ankiascodepb.NoteTypeConfig
	for _, file := range configs.NoteTypes {
		for _, native := range initial.NoteTypes {
			if native.Name == file.Config.Name {
				bindings[file.Path] = fmt.Sprint(native.Id)
			}
		}
		if file.Config.Name == "Basic" {
			definition = proto.Clone(file.Config).(*ankiascodepb.NoteTypeConfig)
		}
	}
	if _, err := states.Import(ctx, configs, bindings); err != nil {
		t.Fatal(err)
	}
	stored, err := states.Dump(ctx)
	if err != nil {
		t.Fatal(err)
	}
	replacementPath := "note_types/Replacement.notetype.anki.toml"
	var replacementID int64
	var replacementName string
	for _, file := range configs.NoteTypes {
		if file.Config.Name != "Basic" {
			replacementID = stored.Resources.NoteTypes[file.Path].Id
			replacementName, file.Path = file.Config.Name, replacementPath
			break
		}
	}
	if replacementID == 0 {
		t.Fatal("native replacement fixture is missing")
	}
	replacementDeckPath := "decks/Replacement/deck.anki.toml"
	for _, group := range configs.Decks {
		if group.Title == "Spare" {
			group.Path = replacementDeckPath
		}
	}
	configs.Decks[0].Notes[0].Config.Fields["Front"] = "After"
	deckPath, typePath, notePath := "decks/Created/deck.anki.toml", "note_types/Created.notetype.anki.toml", "decks/Created/notes/Added.note.anki.toml"
	definition.Name = "Created"
	definition.Templates[0].TargetDeck = deckPath
	configs.NoteTypes = append(configs.NoteTypes, &ankiascodepb.NoteTypeConfigFile{Path: typePath, Config: definition})
	configs.Decks = append(configs.Decks, &ankiascodepb.DeckConfigSet{
		Path: deckPath, Title: "Created", Deck: &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}},
		Notes: []*ankiascodepb.NoteConfigFile{{Path: notePath, Config: &ankiascodepb.NoteConfig{NoteType: typePath, Fields: map[string]string{"Front": "Added", "Back": "Answer"}}}},
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
	detached := &ankiascodepb.CollectionPlan{}
	if err := codec.Decode(saved, codec.JSON, detached); err != nil {
		t.Fatal(err)
	}
	frozen := proto.Clone(detached)
	execution := executionConnection{ClientConnInterface: connection}
	reconciler := reconcile.New(resource.New(execution), state.New(execution, "saved-plan"))
	changed, err := reconciler.Apply(ctx, detached)
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
	replaced, err := resources.NoteTypes.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: recorded.Resources.NoteTypes[replacementPath].Id})
	if err != nil || replaced.GetId() == replacementID || replaced.GetName() != replacementName {
		t.Fatalf("note type was not replaced with its original name: %v %v", replaced, err)
	}
	if _, err := resources.NoteTypes.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: replacementID}); err == nil {
		t.Fatal("replaced native note type remains")
	}
	replacedDeck, err := resources.Decks.GetDeck(ctx, &decks.DeckId{Did: recorded.Resources.Decks[replacementDeckPath].GetId()})
	if err != nil || replacedDeck.GetId() == spare.Id || replacedDeck.GetName() != "Spare" {
		t.Fatalf("deck was not replaced with its original name: %v %v", replacedDeck, err)
	}
	if _, err := resources.Decks.GetDeck(ctx, &decks.DeckId{Did: spare.Id}); err == nil {
		t.Fatal("replaced native deck remains")
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
	if _, err := reconciler.Apply(ctx, detached); err == nil {
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
	checkRootNoteCreation(t, resources, states, connection, configs)
	checkMixedCaseHierarchy(t, resources, connection)
	checkFailedPlans(t, resources, connection)
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "execution-plan.json"), saved, 0o600); err != nil {
			t.Fatal(err)
		}
		emptyData, err := codec.Encode(emptyPlan, codec.JSON)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "empty-execution.json"), emptyData, 0o600); err != nil {
			t.Fatal(err)
		}
		typeData, err := codec.Encode(replaced, codec.JSON)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "replacement-type.json"), typeData, 0o600); err != nil {
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
		opened = false
	}
}

// checkFailedPlans preserves caller-owned plans and stored bindings after partial native execution.
func checkFailedPlans(t *testing.T, resources *resource.Store, connection grpc.ClientConnInterface) {
	t.Helper()
	for _, stage := range []struct {
		// name distinguishes native resources and the state namespace for each failure.
		name string
		// method identifies the native mutation rejected after earlier operations succeed.
		method string
	}{
		{name: "FailureType", method: "/AddNotetype"},
		{name: "FailureNotes", method: "/AddNotes"},
		{name: "FailureCards", method: "/SetDeck"},
	} {
		t.Run(stage.name, func(t *testing.T) {
			ctx := t.Context()
			states := state.New(connection, stage.name)
			stored, err := states.Dump(ctx)
			if err != nil {
				t.Fatal(err)
			}
			snapshot, err := resources.Snapshot(ctx)
			if err != nil {
				t.Fatal(err)
			}
			source := "decks/" + stage.name + "/deck.anki.toml"
			target := "decks/" + stage.name + "Target/deck.anki.toml"
			definition := "note_types/" + stage.name + ".notetype.anki.toml"
			configs := &ankiascodepb.CollectionConfigSet{
				Config: &ankiascodepb.CollectionConfig{Version: 1},
				Decks: []*ankiascodepb.DeckConfigSet{
					{
						Path: source, Title: stage.name,
						Deck: &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}},
						Notes: []*ankiascodepb.NoteConfigFile{{
							Path: "decks/" + stage.name + "/notes/example.note.anki.toml",
							Config: &ankiascodepb.NoteConfig{
								NoteType: definition,
								Fields:   map[string]string{"Front": stage.name, "Back": "Answer"},
								Cards:    []*ankiascodepb.CardConfig{{Ordinal: 0, Deck: proto.String(target)}},
							},
						}},
					},
					{Path: target, Title: stage.name + "Target", Deck: &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}}},
				},
				NoteTypes: []*ankiascodepb.NoteTypeConfigFile{{
					Path: definition,
					Config: &ankiascodepb.NoteTypeConfig{
						Name:   stage.name,
						Fields: []*notetypes.Notetype_Field{{Name: "Front"}, {Name: "Back"}},
						Templates: []*ankiascodepb.TemplateConfig{{
							Name: "Card", TargetDeck: source,
							Config: &notetypes.Notetype_Template_Config{QFormat: "{{Front}}", AFormat: "{{Back}}"},
						}},
					},
				}},
			}
			planned, err := plan.Build(configs, snapshot, stored)
			if err != nil {
				t.Fatal(err)
			}
			frozen := proto.CloneOf(planned)
			failed := executionConnection{ClientConnInterface: connection, reject: stage.method}
			if _, err := reconcile.New(resource.New(failed), states).Apply(ctx, planned); err == nil || !strings.Contains(err.Error(), "injected execution failure") {
				t.Fatalf("native failure was not reached: %v", err)
			}
			if !proto.Equal(planned, frozen) {
				t.Fatal("failed execution changed its input plan")
			}
			actual, err := states.Dump(ctx)
			if err != nil || !proto.Equal(actual, stored) {
				t.Fatalf("failed execution persisted partial bindings: %v %v", actual, err)
			}
			created, err := resources.Decks.GetDeckIdByName(ctx, &generic.String{Val: stage.name})
			if err != nil || created.GetDid() == 0 {
				t.Fatalf("failure did not follow a committed native operation: %v", err)
			}
			if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
				data, err := codec.Encode(planned, codec.JSON)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(output, stage.name+"-plan.json"), data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// checkMixedCaseHierarchy executes accepted parent/child names with different casing.
func checkMixedCaseHierarchy(t *testing.T, resources *resource.Store, connection grpc.ClientConnInterface) {
	t.Helper()
	t.Run("mixed-case hierarchy", func(t *testing.T) {
		ctx := t.Context()
		states := state.New(connection, "mixed-case-hierarchy")
		stored, err := states.Dump(ctx)
		if err != nil {
			t.Fatal(err)
		}
		snapshot, err := resources.Snapshot(ctx)
		if err != nil {
			t.Fatal(err)
		}
		configs := &ankiascodepb.CollectionConfigSet{
			Config: &ankiascodepb.CollectionConfig{Version: 1},
			Decks: []*ankiascodepb.DeckConfigSet{
				{Path: "decks/a/decks/child/deck.anki.toml", Title: "A::child", Deck: &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}}},
				{Path: "decks/a/deck.anki.toml", Title: "a", Deck: &ankiascodepb.DeckConfig{Normal: &decks.Deck_Normal{}}},
			},
		}
		prepared, err := plan.Build(configs, snapshot, stored)
		if err != nil {
			t.Fatal(err)
		}
		reconciler := reconcile.New(resource.New(executionConnection{ClientConnInterface: connection}), states)
		if _, err := reconciler.Apply(ctx, prepared); err != nil {
			t.Fatal(err)
		}
		recorded, err := states.Dump(ctx)
		if err != nil {
			t.Fatal(err)
		}
		for _, group := range configs.Decks {
			id := recorded.GetResources().GetDecks()[group.Path].GetId()
			value, err := resources.Decks.GetDeck(ctx, &decks.DeckId{Did: id})
			if err != nil || id == 0 || !strings.EqualFold(value.GetName(), group.Title) {
				t.Fatalf("native hierarchy binding %q: %v %v", group.Path, value, err)
			}
		}
		if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
			for name, value := range map[string]proto.Message{"mixed-case-deck-plan.json": prepared, "mixed-case-deck-state.json": recorded} {
				data, err := codec.Encode(value, codec.JSON)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(output, name), data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		}
	})
}

// checkRootNoteCreation verifies native deck selection and explicit moves after re-export.
func checkRootNoteCreation(t *testing.T, resources *resource.Store, states *state.Store, connection grpc.ClientConnInterface, configs *ankiascodepb.CollectionConfigSet) {
	t.Helper()
	t.Run("root note creation", func(t *testing.T) {
		ctx := t.Context()
		desired := proto.Clone(configs).(*ankiascodepb.CollectionConfigSet)
		rootPath := "notes/root.note.anki.toml"
		var typePath string
		for _, file := range desired.NoteTypes {
			if file.Config.Name == "Basic" {
				typePath = file.Path
			}
		}
		if typePath == "" {
			t.Fatal("native Basic type is missing")
		}
		desired.Notes = append(desired.Notes, &ankiascodepb.NoteConfigFile{Path: rootPath, Config: &ankiascodepb.NoteConfig{NoteType: typePath, Fields: map[string]string{"Front": "Root", "Back": "Answer"}}})
		stored, err := states.Dump(ctx)
		if err != nil {
			t.Fatal(err)
		}
		snapshot, err := resources.Snapshot(ctx)
		if err != nil {
			t.Fatal(err)
		}
		prepared, err := plan.Build(desired, snapshot, stored)
		if err != nil {
			t.Fatal(err)
		}
		reconciler := reconcile.New(resource.New(executionConnection{ClientConnInterface: connection}), states)
		if _, err := reconciler.Apply(ctx, prepared); err != nil {
			t.Fatal(err)
		}
		recorded, err := states.Dump(ctx)
		if err != nil {
			t.Fatal(err)
		}
		noteID := recorded.Resources.Notes[rootPath].GetId()
		snapshot, err = resources.Snapshot(ctx)
		if err != nil {
			t.Fatal(err)
		}
		converged, err := plan.Build(desired, snapshot, recorded)
		if err != nil || len(converged.GetReport().GetChanges()) != 0 {
			t.Fatalf("root-note creation did not converge: %v %v", converged, err)
		}
		if changed, err := reconciler.Apply(ctx, converged); err != nil || changed {
			t.Fatalf("converged root note changed: %v %v", changed, err)
		}
		output := filepath.Join(t.TempDir(), "export")
		if err := export.New(resources, states, filesystem.FS{}).Export(ctx, output, ""); err != nil {
			t.Fatal(err)
		}
		exported, err := config.New(output, filesystem.FS{}).Load()
		if err != nil {
			t.Fatal(err)
		}
		exportedPath := "decks/Default/notes/root.note.anki.toml"
		if _, err := os.Stat(filepath.Join(output, filepath.FromSlash(exportedPath))); err != nil {
			t.Fatalf("Anki did not choose its default deck: %v", err)
		}
		exported.Config.Moved = []*ankiascodepb.ResourceMove{
			{From: rootPath, To: exportedPath},
			{From: "decks/Replacement", To: "decks/Spare"},
		}
		moved, err := plan.Build(exported, snapshot, recorded)
		if err != nil || moved.GetReport().GetNotesAdded() != 0 || moved.GetReport().GetNotesDeleted() != 0 || moved.GetReport().GetDecksAdded() != 0 || moved.GetReport().GetDecksDeleted() != 0 || moved.GetProposedBindings().GetNotes()[exportedPath].GetId() != noteID {
			t.Fatalf("exported path move lost root-note identity: %v %v", moved, err)
		}
		if _, err := reconciler.Apply(ctx, moved); err != nil {
			t.Fatal(err)
		}
		rebound, err := states.Dump(ctx)
		if err != nil || rebound.GetResources().GetNotes()[exportedPath].GetId() != noteID || rebound.GetResources().GetNotes()[rootPath] != nil {
			t.Fatalf("root-note move changed native identity: %v %v", rebound, err)
		}
		if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
			for name, value := range map[string]proto.Message{"root-note-created.json": recorded, "root-note-converged.json": converged, "root-note-export-move.json": moved, "root-note-rebound.json": rebound} {
				data, err := codec.Encode(value, codec.JSON)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(artifacts, name), data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		}
	})
}
