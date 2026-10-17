package plan_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/export"
	planning "git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/plan"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/protobuf/proto"
)

// TestPlanningWorkflow plans exported native data after stopping the backend.
// It verifies resource ownership, moves, errors, determinism and unchanged inputs.
func TestPlanningWorkflow(t *testing.T) {
	ctx := t.Context()
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	database := filepath.Join(t.TempDir(), "collection.anki2")
	worker := backend.New(backend.BackendSTDOptions{Binary: binary})
	collection := backend.NewCollection(backend.BackendCollectionOptions{Backend: worker, Path: database})
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
	resources := resource.New(connection)
	states := state.New(connection, "planning-test")
	created, err := resources.NoteTypes.AddNotetype(ctx, &notetypes.Notetype{
		Name: "Planning", Config: &notetypes.Notetype_Config{Css: ".card { color: blue; }"},
		Fields: []*notetypes.Notetype_Field{
			{Name: "Front", Config: &notetypes.Notetype_Field_Config{}},
			{Name: "Back", Config: &notetypes.Notetype_Field_Config{}},
		},
		Templates: []*notetypes.Notetype_Template{{Name: "Card", Config: &notetypes.Notetype_Template_Config{QFormat: "{{Front}}", AFormat: "{{Back}}"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	nativeNote, err := resources.Notes.NewNote(ctx, &notetypes.NotetypeId{Ntid: created.Id})
	if err != nil {
		t.Fatal(err)
	}
	nativeNote.Fields = []string{"Question", "Answer"}
	// Native IDs identify notes even when the API accepts empty GUID metadata.
	nativeNote.Guid = ""
	note, err := resources.Notes.AddNote(ctx, &notes.AddNoteRequest{DeckId: 1, Note: nativeNote})
	if err != nil {
		t.Fatal(err)
	}
	home, err := resources.Decks.AddDeck(ctx, &decks.Deck{Name: "Home", Kind: &decks.Deck_Normal_{Normal: &decks.Deck_Normal{ConfigId: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	filtered, err := resources.Decks.AddDeck(ctx, &decks.Deck{Name: "Filtered", Kind: &decks.Deck_Filtered_{Filtered: &decks.Deck_Filtered{}}})
	if err != nil {
		t.Fatal(err)
	}
	target, err := resources.Decks.AddDeck(ctx, &decks.Deck{Name: "Target", Kind: &decks.Deck_Normal_{Normal: &decks.Deck_Normal{ConfigId: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	memberTypes := map[string]int64{}
	for _, name := range []string{"Reordered", "Identified"} {
		definition, err := resources.NoteTypes.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: created.Id})
		if err != nil {
			t.Fatal(err)
		}
		definition.Id, definition.Name = 0, name
		for _, field := range definition.Fields {
			field.Ord = nil
		}
		definition.Templates[0].Ord, definition.Templates[0].Name = nil, "Forward"
		definition.Templates = append(definition.Templates, &notetypes.Notetype_Template{Name: "Reverse", Config: &notetypes.Notetype_Template_Config{QFormat: "{{Back}}", AFormat: "{{Front}}"}})
		if name == "Identified" {
			for index, field := range definition.Fields {
				field.Config.Id = proto.Int64(int64(100 + index))
			}
			for index, template := range definition.Templates {
				template.Config.Id = proto.Int64(int64(200 + index))
			}
		}
		added, err := resources.NoteTypes.AddNotetype(ctx, definition)
		if err != nil {
			t.Fatal(err)
		}
		memberTypes[name] = added.Id
	}
	root := filepath.Join(t.TempDir(), "configs")
	if _, err := resources.Decks.AddDeck(ctx, &decks.Deck{Name: "Spare", Kind: &decks.Deck_Normal_{Normal: &decks.Deck_Normal{ConfigId: 1}}}); err != nil {
		t.Fatal(err)
	}
	if err := export.New(resources, states).Export(ctx, root); err != nil {
		t.Fatal(err)
	}
	configs, err := config.New(root).Load()
	if err != nil {
		t.Fatal(err)
	}
	initial, err := resources.Snapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	bindings := map[string]string{}
	notePath := ""
	for _, deck := range configs.Decks {
		for _, native := range initial.Decks {
			if native.Name == deck.Title {
				bindings[deck.Path] = fmt.Sprint(native.Id)
			}
		}
		for _, file := range deck.Notes {
			notePath = file.Path
			bindings[file.Path] = fmt.Sprint(note.NoteId)
		}
	}
	for _, file := range configs.NoteTypes {
		for _, native := range initial.NoteTypes {
			if native.Name == file.Config.Name {
				bindings[file.Path] = fmt.Sprint(native.Id)
			}
		}
	}
	if notePath == "" {
		t.Fatal("export omitted the native note")
	}
	if _, err := states.Import(ctx, configs, bindings); err != nil {
		t.Fatal(err)
	}
	// An omitted note sharing the managed type and deck must remain untracked.
	untracked, err := resources.Notes.NewNote(ctx, &notetypes.NotetypeId{Ntid: created.Id})
	if err != nil {
		t.Fatal(err)
	}
	untracked.Fields = []string{"Untracked", "Retained"}
	untracked.Guid = ""
	residentNote, err := resources.Notes.AddNote(ctx, &notes.AddNoteRequest{DeckId: home.Id, Note: untracked})
	if err != nil {
		t.Fatal(err)
	}
	generated, err := resources.Notes.CardsOfNote(ctx, &notes.NoteId{Nid: residentNote.NoteId})
	if err != nil || len(generated.GetCids()) != 1 {
		t.Fatalf("resident card fixture: %v", err)
	}
	resident, err := resources.Cards.GetCard(ctx, &cards.CardId{Cid: generated.Cids[0]})
	if err != nil {
		t.Fatal(err)
	}
	resident.OriginalDeckId, resident.DeckId, resident.OriginalDue = home.Id, filtered.Id, resident.Due
	if _, err := resources.Cards.UpdateCards(ctx, &cards.UpdateCardsRequest{Cards: []*cards.Card{resident}, SkipUndoEntry: true}); err != nil {
		t.Fatal(err)
	}
	// An untracked template can reference an otherwise empty managed deck.
	if _, err := resources.NoteTypes.AddNotetype(ctx, &notetypes.Notetype{
		Name: "Untracked target", Config: &notetypes.Notetype_Config{},
		Fields: []*notetypes.Notetype_Field{{Name: "Front", Config: &notetypes.Notetype_Field_Config{}}},
		Templates: []*notetypes.Notetype_Template{{Name: "Card", Config: &notetypes.Notetype_Template_Config{
			QFormat: "{{Front}}", AFormat: "{{Front}}", TargetDeckId: target.Id,
		}}},
	}); err != nil {
		t.Fatal(err)
	}
	stored, err := states.Dump(ctx)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := resources.Snapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// Capture actual native drift while retaining the previously imported state.
	for _, name := range []string{"Reordered", "Identified"} {
		definition, err := resources.NoteTypes.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: memberTypes[name]})
		if err != nil {
			t.Fatal(err)
		}
		if name == "Reordered" {
			slices.Reverse(definition.Fields)
			slices.Reverse(definition.Templates)
		} else {
			if definition.Fields[0].Config.GetId() != 100 || definition.Templates[0].Config.GetId() != 200 {
				t.Fatal("native stable-member fixture lost its identities")
			}
			definition.Fields = definition.Fields[1:]
			definition.Templates = definition.Templates[1:]
			definition.Templates[0].Config.AFormat = "{{Back}}"
		}
		if _, err := resources.NoteTypes.UpdateNotetype(ctx, definition); err != nil {
			t.Fatal(err)
		}
	}
	drifted, err := resources.Snapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	clozeConfigs, clozeSnapshot, clozeState := captureFilteredCloze(t, resources, state.New(connection, "filtered-cloze-plan"), filtered.Id)
	emptyStates := map[string]*collectionpb.State{}
	for _, name := range []string{"absent", "empty"} {
		store := state.New(connection, "empty-planning-"+name)
		if name == "empty" {
			if _, err := store.Replace(ctx, &collectionpb.State{Version: 1, Resources: &collectionpb.Resources{}}); err != nil {
				t.Fatal(err)
			}
		}
		value, err := store.Dump(ctx)
		if err != nil {
			t.Fatal(err)
		}
		emptyStates[name] = value
	}
	if emptyStates["absent"].Resources != nil || emptyStates["empty"].Resources == nil {
		t.Fatal("native empty-state fixtures lost their resource presence distinction")
	}
	output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if output != "" {
		if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: filepath.Join(output, "planning.colpkg")}); err != nil {
			t.Fatal(err)
		}
	}
	if err := collection.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	if err := worker.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(database)
	if err != nil {
		t.Fatal(err)
	}
	originalConfigs, originalState := proto.Clone(configs), proto.Clone(stored)
	originalSnapshot, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		// name describes the intended collection change.
		name string
		// edit changes only detached desired configs.
		edit func(*collectionpb.CollectionConfigSet)
		// check verifies externally visible plan behavior.
		check func(*testing.T, *collectionpb.CollectionChangePlan)
	}{
		{name: "converged", check: func(t *testing.T, plan *collectionpb.CollectionChangePlan) {
			if len(plan.Changes) != 0 || plan.StateChanged || plan.CardGenerationDeferred {
				t.Fatalf("nonempty converged plan: %v", plan)
			}
		}},
		{name: "content", edit: func(value *collectionpb.CollectionConfigSet) {
			value.Decks[0].Notes[0].Config.Fields["Front"] = "Changed"
			value.Decks[0].Notes[0].Config.Tags = []string{"lesson"}
		}, check: func(t *testing.T, plan *collectionpb.CollectionChangePlan) {
			if plan.NotesChanged != 1 || !plan.CardGenerationDeferred || plan.StateChanged {
				t.Fatalf("content plan: %v", plan)
			}
		}},
		{name: "creation", edit: func(value *collectionpb.CollectionConfigSet) {
			file := proto.Clone(value.Decks[0].Notes[0]).(*collectionpb.NoteConfigFile)
			file.Path = "decks/Default/notes/new.note.anki.toml"
			value.Decks[0].Notes = append(value.Decks[0].Notes, file)
		}, check: func(t *testing.T, plan *collectionpb.CollectionChangePlan) {
			if plan.NotesAdded != 1 || plan.CardsAdded != 0 || !plan.CardGenerationDeferred || !plan.StateChanged {
				t.Fatalf("creation plan: %v", plan)
			}
		}},
		{name: "filtered kind projection", edit: func(value *collectionpb.CollectionConfigSet) {
			for _, group := range value.Decks {
				if group.Title == "Spare" {
					group.Deck = &collectionpb.DeckConfig{Filtered: &decks.Deck_Filtered{}}
				}
			}
			file := proto.Clone(value.Decks[0].Notes[0]).(*collectionpb.NoteConfigFile)
			file.Path = "decks/Default/notes/generated.note.anki.toml"
			value.Decks[0].Notes = append(value.Decks[0].Notes, file)
		}, check: func(t *testing.T, plan *collectionpb.CollectionChangePlan) {
			if plan.DecksChanged != 1 || plan.NotesAdded != 1 || !plan.CardGenerationDeferred {
				t.Fatalf("filtered kind projection: %v", plan)
			}
		}},
		{name: "deletion", edit: func(value *collectionpb.CollectionConfigSet) { value.Decks[0].Notes = nil }, check: func(t *testing.T, plan *collectionpb.CollectionChangePlan) {
			if plan.NotesDeleted != 1 || plan.CardsDeleted != 1 || !plan.StateChanged {
				t.Fatalf("managed deletion plan: %v", plan)
			}
		}},
		{name: "move", edit: func(value *collectionpb.CollectionConfigSet) {
			destination := "decks/Default/notes/renamed.note.anki.toml"
			value.Config.Moved = []*collectionpb.ResourceMove{{From: notePath, To: destination}}
			value.Decks[0].Notes[0].Path = destination
		}, check: func(t *testing.T, plan *collectionpb.CollectionChangePlan) {
			if plan.NotesAdded != 0 || plan.NotesDeleted != 0 || plan.NotesChanged != 0 || !plan.StateChanged || len(plan.Changes) != 1 || plan.Changes[0].Action != "move" || plan.Changes[0].From != notePath {
				t.Fatalf("identity-preserving move: %v", plan)
			}
		}},
		{name: "independent directory moves", edit: func(value *collectionpb.CollectionConfigSet) {
			value.Config.Moved = []*collectionpb.ResourceMove{
				{From: "decks/Default", To: "decks/Renamed"},
				{From: "decks/Spare", To: "decks/Relocated"},
			}
			for _, group := range value.Decks {
				for _, move := range value.Config.Moved {
					if strings.HasPrefix(group.Path, move.From+"/") {
						group.Path = move.To + strings.TrimPrefix(group.Path, move.From)
						for _, note := range group.Notes {
							note.Path = move.To + strings.TrimPrefix(note.Path, move.From)
							for _, card := range note.Config.Cards {
								if strings.HasPrefix(card.GetDeck(), move.From+"/") {
									card.Deck = proto.String(move.To + strings.TrimPrefix(card.GetDeck(), move.From))
								}
							}
						}
					}
				}
			}
		}, check: func(t *testing.T, plan *collectionpb.CollectionChangePlan) {
			if len(plan.Changes) != 3 || !plan.StateChanged || plan.CardGenerationDeferred {
				t.Fatalf("directory moves did not retain native resources: %v", plan)
			}
			for _, change := range plan.Changes {
				if change.Action != "move" {
					t.Fatalf("directory move replaced a resource: %v", change)
				}
			}
		}},
		{name: "native deck normalization", edit: func(value *collectionpb.CollectionConfigSet) {
			value.Decks[0].Title = " Default "
		}, check: func(t *testing.T, plan *collectionpb.CollectionChangePlan) {
			if plan.DecksChanged != 1 {
				t.Fatalf("native deck normalization was validated locally: %v", plan)
			}
		}},
		{name: "mixed-case hierarchy", edit: func(value *collectionpb.CollectionConfigSet) {
			value.Decks = append(value.Decks,
				&collectionpb.DeckConfigSet{Path: "decks/a/decks/child/deck.anki.toml", Title: "A::child", Deck: &collectionpb.DeckConfig{Normal: &decks.Deck_Normal{}}},
				&collectionpb.DeckConfigSet{Path: "decks/a/deck.anki.toml", Title: "a", Deck: &collectionpb.DeckConfig{Normal: &decks.Deck_Normal{}}},
			)
		}, check: func(t *testing.T, plan *collectionpb.CollectionChangePlan) {
			if plan.DecksAdded != 2 {
				t.Fatalf("mixed-case hierarchy plan: %v", plan)
			}
		}},
		{name: "presentation", edit: func(value *collectionpb.CollectionConfigSet) {
			for _, file := range value.NoteTypes {
				if file.Config.Name == "Planning" {
					file.Config.Config.Css = ".card { color: red; }"
				}
			}
		}, check: func(t *testing.T, plan *collectionpb.CollectionChangePlan) {
			if plan.NoteTypesChanged != 1 || !plan.CardGenerationDeferred {
				t.Fatalf("presentation plan: %v", plan)
			}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
			if test.edit != nil {
				test.edit(desired)
			}
			saved := proto.Clone(desired)
			plan, err := planning.Build(desired, snapshot, stored)
			if err != nil {
				t.Fatal(err)
			}
			test.check(t, plan.Report)
			if test.name == "mixed-case hierarchy" {
				if len(plan.Decks) != 2 || plan.Decks[0].Value.Name != "a" || plan.Decks[1].Value.Name != "A::child" {
					t.Fatalf("child deck precedes its declared parent: %v", plan.Decks)
				}
			}
			if test.name == "converged" || test.name == "move" {
				if len(plan.CardAssignments)+len(plan.Notes)+len(plan.NoteTypes)+len(plan.Decks) != 0 {
					t.Fatal("unchanged native resources still have executable operations")
				}
			}
			if (test.name == "creation" || test.name == "content") && len(plan.CardAssignments) != 1 {
				t.Fatalf("assignments include unrelated notes: %v", plan.CardAssignments)
			}
			if len(plan.CardAssignments) > 0 && !slices.Contains(plan.FilteredDecks, filtered.Id) {
				t.Fatal("card execution lost the existing filtered-deck exclusion")
			}
			if test.name == "filtered kind projection" && !slices.Contains(plan.FilteredDecks, stored.Resources.Decks["decks/Spare/deck.anki.toml"].Id) {
				t.Fatal("card execution did not use the resulting filtered kind")
			}
			permuted := proto.Clone(desired).(*collectionpb.CollectionConfigSet)
			slices.Reverse(permuted.Decks)
			slices.Reverse(permuted.NoteTypes)
			slices.Reverse(permuted.Notes)
			slices.Reverse(permuted.Config.Moved)
			for _, group := range permuted.Decks {
				slices.Reverse(group.Notes)
			}
			reordered, err := planning.Build(permuted, snapshot, stored)
			if err != nil || !proto.Equal(plan, reordered) {
				t.Fatalf("declaration order changed the executable plan: %v", err)
			}
			repeated, err := planning.Build(desired, snapshot, stored)
			if err != nil || !proto.Equal(plan, repeated) {
				t.Fatalf("nondeterministic plan: %v", err)
			}
			if !proto.Equal(desired, saved) {
				t.Fatal("planning changed supplied configs")
			}
			if output != "" {
				writePlanningArtifact(t, filepath.Join(output, test.name+".json"), plan.Report)
				writePlanningArtifact(t, filepath.Join(output, test.name+"-execution.json"), plan)
				if test.name == "content" {
					if err := config.New(filepath.Join(output, "desired")).Export(desired); err != nil {
						t.Fatal(err)
					}
				}
			}
		})
	}
	for name, value := range emptyStates {
		t.Run("empty namespace "+name, func(t *testing.T) {
			original := proto.Clone(value)
			empty := &collectionpb.CollectionConfigSet{Config: &collectionpb.CollectionConfig{Version: 1}}
			prepared, err := planning.Build(empty, snapshot, value)
			if err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(prepared.Report, &collectionpb.CollectionChangePlan{}) {
				t.Fatalf("empty namespace produced changes: %v", prepared.Report)
			}
			if !proto.Equal(value, original) {
				t.Fatal("planning mutated the empty stored state")
			}
			if output != "" {
				writePlanningArtifact(t, filepath.Join(output, "empty-namespace-"+name+".json"), prepared)
			}
		})
	}
	t.Run("replace tracked deck", func(t *testing.T) {
		desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
		previous, replacement := "decks/Spare/deck.anki.toml", "decks/Replacement/deck.anki.toml"
		var declaration *collectionpb.DeckConfigSet
		for _, group := range desired.Decks {
			if group.Path == previous {
				group.Path, declaration = replacement, group
			}
		}
		if declaration == nil {
			t.Fatal("native deck replacement fixture is missing")
		}
		prepared, err := planning.Build(desired, snapshot, stored)
		if err != nil {
			t.Fatal(err)
		}
		if prepared.Report.DecksAdded != 1 || prepared.Report.DecksDeleted != 1 || !slices.Contains(prepared.RemovedDecks, stored.Resources.Decks[previous].Id) {
			t.Fatalf("same-name deck replacement was not planned: %v", prepared.Report)
		}
		if prepared.ProposedBindings.Decks[previous] != nil || (prepared.ProposedBindings.Decks[replacement] == nil || prepared.ProposedBindings.Decks[replacement].GetId() != 0) {
			t.Fatal("replacement retained the deleted deck binding")
		}
		unbound := &collectionpb.CollectionConfigSet{Config: &collectionpb.CollectionConfig{Version: 1}, Decks: []*collectionpb.DeckConfigSet{declaration}}
		if _, err := planning.Build(unbound, snapshot, emptyStates["absent"]); err == nil {
			t.Fatal("deck replacement claimed an untracked native name")
		}
		if output != "" {
			writePlanningArtifact(t, filepath.Join(output, "deck-replacement.json"), prepared)
		}
	})
	t.Run("replace tracked note type", func(t *testing.T) {
		desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
		path := "note_types/Replacement.notetype.anki.toml"
		var previous string
		for _, file := range desired.NoteTypes {
			if file.Config.Name == "Reordered" {
				previous, file.Path = file.Path, path
				break
			}
		}
		if previous == "" {
			t.Fatal("native replacement fixture is missing")
		}
		prepared, err := planning.Build(desired, snapshot, stored)
		if err != nil {
			t.Fatal(err)
		}
		if prepared.Report.NoteTypesAdded != 1 || prepared.Report.NoteTypesDeleted != 1 || !slices.Contains(prepared.RemovedNoteTypes, stored.Resources.NoteTypes[previous].Id) {
			t.Fatalf("same-name replacement was not planned: %v", prepared.Report)
		}
		if prepared.ProposedBindings.NoteTypes[previous] != nil || prepared.ProposedBindings.NoteTypes[path] == nil || prepared.ProposedBindings.NoteTypes[path].Id != 0 {
			t.Fatal("replacement reused the deleted note-type binding")
		}
		if _, err := planning.Build(desired, snapshot, emptyStates["absent"]); err == nil {
			t.Fatal("replacement claimed an untracked native name")
		}
		if output != "" {
			writePlanningArtifact(t, filepath.Join(output, "note-type-replacement.json"), prepared)
		}
	})
	checkMemberDrift(t, configs, drifted, stored, output)
	checkNewReferences(t, configs, snapshot, stored, output)
	checkFilteredCloze(t, clozeConfigs, clozeSnapshot, clozeState, output)
	t.Run("destination only", func(t *testing.T) {
		desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
		desired.Decks[0].Notes[0].Config.Cards[0].Deck = proto.String("decks/Home/deck.anki.toml")
		prepared, err := planning.Build(desired, snapshot, stored)
		if err != nil {
			t.Fatal(err)
		}
		if prepared.Report.CardsChanged != 1 || len(prepared.CardAssignments) != 1 || len(prepared.Notes) != 0 {
			t.Fatal("destination-only edit was lost or rewrites note content")
		}
	})
	t.Run("root note without destination", func(t *testing.T) {
		desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
		desired.Notes, desired.Decks[0].Notes = desired.Decks[0].Notes, nil
		desired.Notes[0].Config.Cards = nil
		desired.Notes[0].Config.Fields["Front"] = "Root content"
		prepared, err := planning.Build(desired, snapshot, stored)
		if err != nil {
			t.Fatal(err)
		}
		if len(prepared.Notes) != 1 || len(prepared.CardAssignments) != 0 {
			t.Fatal("root note without destinations contains assignment operations")
		}
	})
	// No cards occupy Target: only the untracked template retains this deck.
	for _, action := range []string{"remove", "convert"} {
		t.Run("untracked template deck "+action, func(t *testing.T) {
			desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
			for index, group := range desired.Decks {
				if group.Title != "Target" {
					continue
				}
				if action == "remove" {
					desired.Decks = append(desired.Decks[:index], desired.Decks[index+1:]...)
				} else {
					group.Deck = &collectionpb.DeckConfig{Filtered: &decks.Deck_Filtered{}}
				}
				break
			}
			if plan, err := planning.Build(desired, snapshot, stored); err == nil || plan != nil {
				t.Fatal("untracked template lost its normal-deck target")
			}
		})
	}
	for _, test := range []struct {
		// name identifies an invalid request.
		name string
		// edit constructs the invalid detached request.
		edit func(*collectionpb.CollectionConfigSet)
	}{
		{name: "filtered home deck conversion", edit: func(value *collectionpb.CollectionConfigSet) {
			for _, group := range value.Decks {
				if group.Title == "Home" {
					group.Deck = &collectionpb.DeckConfig{Filtered: &decks.Deck_Filtered{}}
				}
			}
		}},
		{name: "unknown type", edit: func(value *collectionpb.CollectionConfigSet) {
			value.Decks[0].Notes[0].Config.NoteType = "note_types/missing.notetype.anki.toml"
		}},
		{name: "missing field", edit: func(value *collectionpb.CollectionConfigSet) { delete(value.Decks[0].Notes[0].Config.Fields, "Back") }},
		{name: "untracked reference", edit: func(value *collectionpb.CollectionConfigSet) {
			value.Notes = value.Decks[0].Notes
			value.Decks = nil
			value.Notes[0].Config.Cards = nil
		}},
		{name: "cross-resource move collision", edit: func(value *collectionpb.CollectionConfigSet) {
			value.Config.Moved = []*collectionpb.ResourceMove{{From: notePath, To: "decks/Spare/deck.anki.toml"}}
			value.Decks[0].Notes = nil
		}},
		{name: "occupied move destination", edit: func(value *collectionpb.CollectionConfigSet) {
			value.Config.Moved = []*collectionpb.ResourceMove{{From: "decks/Default", To: "decks/Spare"}}
			value.Decks = slices.DeleteFunc(value.Decks, func(group *collectionpb.DeckConfigSet) bool { return group.Title == "Default" })
		}},
		{name: "overlapping moves", edit: func(value *collectionpb.CollectionConfigSet) {
			value.Config.Moved = []*collectionpb.ResourceMove{{From: "notes/old", To: "notes/new"}, {From: "notes/new", To: "notes/third"}}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
			test.edit(desired)
			if plan, err := planning.Build(desired, snapshot, stored); err == nil || plan != nil {
				t.Fatal("invalid request returned a plan")
			}
		})
	}
	// State member addresses cannot distinguish duplicate declared names.
	for _, kind := range []string{"field", "template"} {
		t.Run("duplicate "+kind+" reference", func(t *testing.T) {
			desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
			for _, file := range desired.NoteTypes {
				if file.Config.Name != "Planning" {
					continue
				}
				if kind == "field" {
					file.Config.Fields[1].Name = file.Config.Fields[0].Name
				} else {
					file.Config.Templates = append(file.Config.Templates, proto.Clone(file.Config.Templates[0]).(*collectionpb.TemplateConfig))
				}
			}
			if plan, err := planning.Build(desired, initial, stored); err == nil || plan != nil {
				t.Fatal("ambiguous member state address accepted")
			}
		})
	}
	if plan, err := planning.Build(nil, snapshot, stored); err == nil || plan != nil {
		t.Fatal("missing configs accepted")
	}
	if plan, err := planning.Build(configs, nil, stored); err == nil || plan != nil {
		t.Fatal("missing snapshot accepted")
	}
	if plan, err := planning.Build(configs, snapshot, nil); err == nil || plan != nil {
		t.Fatal("missing state accepted")
	}
	checkMemberMoveOwnership(t, configs, initial, stored)
	checkPayloadOwnership(t, configs, snapshot, stored, output)
	checkPresetBindings(t, configs, snapshot, stored, output)
	afterSnapshot, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(database)
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(configs, originalConfigs) || !proto.Equal(stored, originalState) || !bytes.Equal(originalSnapshot, afterSnapshot) || !bytes.Equal(before, after) {
		t.Fatal("planning modified supplied values or the closed database")
	}
	if output != "" {
		writePlanningArtifact(t, filepath.Join(output, "state.json"), stored)
		if err := os.WriteFile(filepath.Join(output, "snapshot.json"), originalSnapshot, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// checkPresetBindings preserves existing aliases and discovers referenced native preset names.
func checkPresetBindings(t *testing.T, configs *collectionpb.CollectionConfigSet, snapshot *resource.Snapshot, stored *collectionpb.State, output string) {
	t.Helper()
	t.Run("preset binding ownership", func(t *testing.T) {
		desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
		bindings := proto.Clone(stored).(*collectionpb.State)
		bindings.Resources.Presets = map[string]*collectionpb.PresetState{"Saved alias": {Id: 1}, "Unused alias": {Id: 1}}
		for _, group := range desired.Decks {
			if group.Deck.Normal != nil {
				group.Deck.Preset = snapshot.Presets[1]
				if group.Title == "Default" {
					group.Deck.Preset = "Saved alias"
				}
			}
		}
		original := proto.Clone(bindings)
		prepared, err := planning.Build(desired, snapshot, bindings)
		if err != nil {
			t.Fatal(err)
		}
		for _, alias := range []string{"Saved alias", "Unused alias", snapshot.Presets[1]} {
			if prepared.ProposedBindings.Presets[alias].GetId() != 1 {
				t.Fatalf("lost native preset alias %q", alias)
			}
		}
		if output != "" {
			writePlanningArtifact(t, filepath.Join(output, "preset-bindings.json"), prepared)
		}
		prepared.ProposedBindings.Presets["Saved alias"].Id++
		if !proto.Equal(bindings, original) {
			t.Fatal("resolved preset bindings alias supplied state")
		}
	})
}

// checkPayloadOwnership edits returned payloads without changing caller-owned inputs.
func checkPayloadOwnership(t *testing.T, configs *collectionpb.CollectionConfigSet, snapshot *resource.Snapshot, stored *collectionpb.State, output string) {
	t.Helper()
	t.Run("returned payload ownership", func(t *testing.T) {
		desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
		desired.Decks[0].Title += " Renamed"
		desired.Decks[0].Notes[0].Config.Fields["Front"] = "Owned content"
		for _, file := range desired.NoteTypes {
			file.Config.Name += " Renamed"
		}
		saved := proto.Clone(desired)
		prepared, err := planning.Build(desired, snapshot, stored)
		if err != nil {
			t.Fatal(err)
		}
		if len(prepared.Decks) == 0 || len(prepared.NoteTypes) == 0 || len(prepared.Notes) == 0 || len(prepared.CardAssignments) == 0 {
			t.Fatal("ownership fixture lacks native operations")
		}
		if output != "" {
			writePlanningArtifact(t, filepath.Join(output, "owned-payloads.json"), prepared)
		}
		for _, deck := range prepared.Decks {
			deck.Value.Name = "Edited result"
			if deck.Value.GetNormal() != nil {
				deck.Value.GetNormal().ConfigId++
			}
		}
		definitions := slices.Clone(prepared.NoteTypeRenames)
		for _, operation := range prepared.NoteTypes {
			definitions = append(definitions, operation.Value)
		}
		for _, definition := range definitions {
			definition.Config.Css = "Edited result"
			for _, field := range definition.Fields {
				field.Name = "Edited result"
				if field.Config.Id != nil {
					(*field.Config.Id)++
				}
			}
			for _, template := range definition.Templates {
				template.Config.QFormat = "Edited result"
				if template.Config.Id != nil {
					(*template.Config.Id)++
				}
			}
		}
		prepared.Notes[0].Value.Fields[0] = "Edited result"
		for _, assignments := range prepared.CardAssignments {
			for _, card := range assignments.Cards {
				card.Ordinal++
			}
		}
		proto.Reset(prepared.ExpectedState.Resources)
		proto.Reset(prepared.ProposedBindings)
		if !proto.Equal(desired, saved) {
			t.Fatal("returned operations alias supplied configs")
		}
	})
}

// checkNewReferences verifies typed path references without provisional Anki identities.
func checkNewReferences(t *testing.T, configs *collectionpb.CollectionConfigSet, snapshot *resource.Snapshot, stored *collectionpb.State, output string) {
	t.Helper()
	t.Run("new object path references", func(t *testing.T) {
		desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
		deckPath, typePath, notePath := "decks/New/deck.anki.toml", "note_types/New.notetype.anki.toml", "decks/New/notes/New.note.anki.toml"
		var definition *collectionpb.NoteTypeConfig
		for _, file := range desired.NoteTypes {
			if file.Config.Name == "Planning" {
				definition = proto.Clone(file.Config).(*collectionpb.NoteTypeConfig)
			}
		}
		if definition == nil {
			t.Fatal("native note type fixture missing")
		}
		definition.Name, definition.Templates[0].TargetDeck = "New", deckPath
		desired.NoteTypes = append(desired.NoteTypes, &collectionpb.NoteTypeConfigFile{Path: typePath, Config: definition})
		desired.Decks = append(desired.Decks, &collectionpb.DeckConfigSet{Path: deckPath, Title: "New", Deck: &collectionpb.DeckConfig{Normal: &decks.Deck_Normal{}}, Notes: []*collectionpb.NoteConfigFile{{Path: notePath, Config: &collectionpb.NoteConfig{NoteType: typePath, Fields: map[string]string{"Front": "New", "Back": "Answer"}}}}})
		prepared, err := planning.Build(desired, snapshot, stored)
		if err != nil {
			t.Fatal(err)
		}
		if !proto.Equal(prepared.ExpectedState, stored) || prepared.ProposedBindings.Decks[deckPath].Id != 0 || prepared.ProposedBindings.NoteTypes[typePath].Id != 0 || prepared.ProposedBindings.Notes[notePath].Id != 0 {
			t.Fatal("proposed bindings invented native identities")
		}
		if len(prepared.Notes) != 1 || prepared.Notes[0].NoteType != typePath || prepared.Notes[0].Value.NotetypeId != 0 {
			t.Fatal("new note did not retain its note-type path reference")
		}
		if len(prepared.NoteTypes) != 1 || prepared.NoteTypes[0].TemplateDecks["Card"] != deckPath || prepared.NoteTypes[0].Value.Templates[0].Config.TargetDeckId != 0 {
			t.Fatal("new template did not retain its target-deck path reference")
		}
		data, err := codec.Encode(prepared, codec.JSON)
		if err != nil {
			t.Fatal(err)
		}
		decoded := &collectionpb.CollectionPlan{}
		if err := codec.Decode(data, codec.JSON, decoded); err != nil || !proto.Equal(prepared, decoded) {
			t.Fatalf("path-based plan did not round-trip: %v", err)
		}
		if output != "" {
			writePlanningArtifact(t, filepath.Join(output, "new-references.json"), prepared)
		}
	})
}

// checkMemberDrift verifies native member mappings after independent Anki edits.
func checkMemberDrift(t *testing.T, configs *collectionpb.CollectionConfigSet, snapshot *resource.Snapshot, stored *collectionpb.State, output string) {
	t.Helper()
	for _, rename := range []bool{false, true} {
		t.Run(fmt.Sprintf("native member drift rename=%t", rename), func(t *testing.T) {
			desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
			if rename {
				for _, file := range desired.NoteTypes {
					if file.Config.Name == "Reordered" {
						desired.Config.Moved = []*collectionpb.ResourceMove{
							{From: file.Path + "#fields/Front", To: file.Path + "#fields/Question"},
							{From: file.Path + "#templates/Forward", To: file.Path + "#templates/Prompt"},
						}
						file.Config.Fields[0].Name = "Question"
						file.Config.Templates[0].Name = "Prompt"
						file.Config.Templates[0].Config.QFormat = "{{Question}}"
						file.Config.Templates[1].Config.AFormat = "{{Question}}"
					}
				}
			}
			prepared, err := planning.Build(desired, snapshot, stored)
			if err != nil {
				t.Fatal(err)
			}
			seen := 0
			for _, operation := range prepared.NoteTypes {
				value := operation.Value
				switch value.Name {
				case "Reordered":
					seen++
					if value.Fields[0].Ord == nil || value.Fields[0].Ord.Val != 1 || value.Templates[0].Ord == nil || value.Templates[0].Ord.Val != 1 {
						t.Error("saved ordinals exchanged reordered native members")
					}
				case "Identified":
					seen++
					if value.Fields[0].Ord != nil || value.Templates[0].Ord != nil || value.Fields[1].Ord == nil || value.Fields[1].Ord.Val != 0 || value.Templates[1].Ord == nil || value.Templates[1].Ord.Val != 0 {
						t.Error("removed stable member claimed a surviving native ordinal")
					}
				}
			}
			if seen != 2 {
				t.Fatal("member drift did not prepare both affected types")
			}
			if output != "" {
				writePlanningArtifact(t, filepath.Join(output, fmt.Sprintf("member-drift-%t.json", rename)), prepared)
			}
		})
	}
}

// checkMemberMoveOwnership rejects identity transfers between native note types.
func checkMemberMoveOwnership(t *testing.T, configs *collectionpb.CollectionConfigSet, snapshot *resource.Snapshot, stored *collectionpb.State) {
	t.Helper()
	for _, kind := range []string{"fields", "templates"} {
		t.Run("cross-note-type "+kind+" move", func(t *testing.T) {
			desired := proto.Clone(configs).(*collectionpb.CollectionConfigSet)
			var source, destination *collectionpb.NoteTypeConfigFile
			for _, file := range desired.NoteTypes {
				if file.Config.Name == "Basic" {
					destination = file
				}
				if kind == "fields" && file.Config.Name == "Planning" || kind == "templates" && len(file.Config.Templates) > 1 {
					source = file
				}
			}
			if source == nil || destination == nil {
				t.Fatal("member move fixture lacks the native source or destination type")
			}
			name := ""
			if kind == "fields" {
				field := source.Config.Fields[1]
				name = field.Name
				source.Config.Fields = source.Config.Fields[:1]
				for _, group := range desired.Decks {
					for _, note := range group.Notes {
						if note.Config.NoteType == source.Path {
							delete(note.Config.Fields, name)
						}
					}
				}
				source.Config.Templates[0].Config.AFormat = "{{Front}}"
				field.Name = "Transferred"
				destination.Config.Fields = append(destination.Config.Fields, field)
			} else {
				template := source.Config.Templates[1]
				name = template.Name
				source.Config.Templates = source.Config.Templates[:1]
				template.Name = "Transferred"
				destination.Config.Templates = append(destination.Config.Templates, template)
			}
			desired.Config.Moved = []*collectionpb.ResourceMove{{From: source.Path + "#" + kind + "/" + name, To: destination.Path + "#" + kind + "/Transferred"}}
			if plan, err := planning.Build(desired, snapshot, stored); err == nil || plan != nil {
				t.Fatal("member identity transferred between unrelated native note types")
			}
		})
	}
}

// writePlanningArtifact serializes a repeatable protobuf result for inspection.
func writePlanningArtifact(t *testing.T, filename string, message proto.Message) {
	t.Helper()
	data, err := codec.Encode(message, codec.JSON)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, data, 0o600); err != nil {
		t.Fatal(err)
	}
}
