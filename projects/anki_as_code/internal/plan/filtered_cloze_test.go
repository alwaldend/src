package plan_test

import (
	"fmt"
	"path/filepath"
	"testing"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/export"
	planning "git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/plan"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// captureFilteredCloze records real resident cloze cards and canonical exported configs.
func captureFilteredCloze(t *testing.T, resources *resource.Store, states *state.Store, filteredID int64) (*ankiascodepb.CollectionConfigSet, *resource.Snapshot, *ankiascodepb.State) {
	t.Helper()
	ctx := t.Context()
	created, err := resources.NoteTypes.AddNotetype(ctx, &notetypes.Notetype{
		Name: "Filtered cloze", Config: &notetypes.Notetype_Config{Kind: notetypes.Notetype_Config_KIND_CLOZE},
		Fields:    []*notetypes.Notetype_Field{{Name: "Text", Config: &notetypes.Notetype_Field_Config{}}, {Name: "Extra", Config: &notetypes.Notetype_Field_Config{}}},
		Templates: []*notetypes.Notetype_Template{{Name: "Cloze", Config: &notetypes.Notetype_Template_Config{QFormat: "{{cloze:Text}}", AFormat: "{{cloze:Text}}<br>{{Extra}}"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	note, err := resources.Notes.NewNote(ctx, &notetypes.NotetypeId{Ntid: created.Id})
	if err != nil {
		t.Fatal(err)
	}
	note.Fields = []string{"{{c1::First}} {{c2::Second}}", "Answer"}
	added, err := resources.Notes.AddNote(ctx, &notes.AddNoteRequest{Note: note, DeckId: 1})
	if err != nil {
		t.Fatal(err)
	}
	identities, err := resources.Notes.CardsOfNote(ctx, &notes.NoteId{Nid: added.NoteId})
	if err != nil || len(identities.GetCids()) != 2 {
		t.Fatalf("native cloze fixture: %v", err)
	}
	for _, id := range identities.Cids {
		card, err := resources.Cards.GetCard(ctx, &cards.CardId{Cid: id})
		if err != nil {
			t.Fatal(err)
		}
		card.OriginalDeckId, card.DeckId, card.OriginalDue = 1, filteredID, card.Due
		if _, err := resources.Cards.UpdateCards(ctx, &cards.UpdateCardsRequest{Cards: []*cards.Card{card}, SkipUndoEntry: true}); err != nil {
			t.Fatal(err)
		}
	}
	root := filepath.Join(t.TempDir(), "cloze")
	if err := export.New(resources, states).Export(ctx, root, ""); err != nil {
		t.Fatal(err)
	}
	all, err := config.New(root).Load()
	if err != nil {
		t.Fatal(err)
	}
	configs := &ankiascodepb.CollectionConfigSet{Config: all.Config}
	for _, file := range all.NoteTypes {
		if file.Config.Name == "Filtered cloze" {
			configs.NoteTypes = append(configs.NoteTypes, file)
		}
	}
	if len(configs.NoteTypes) != 1 {
		t.Fatal("export omitted the native cloze type")
	}
	files := append([]*ankiascodepb.NoteConfigFile{}, all.Notes...)
	for _, group := range all.Decks {
		files = append(files, group.Notes...)
	}
	for _, file := range files {
		if file.Config.NoteType == configs.NoteTypes[0].Path {
			file.Path = "notes/filtered-cloze.note.anki.toml"
			configs.Notes = append(configs.Notes, file)
		}
	}
	if len(configs.Notes) != 1 || len(configs.Notes[0].Config.Cards) != 0 {
		t.Fatal("export did not leave resident cloze cards unmanaged")
	}
	if _, err := states.Import(ctx, configs, map[string]string{configs.NoteTypes[0].Path: fmt.Sprint(created.Id), configs.Notes[0].Path: fmt.Sprint(added.NoteId)}); err != nil {
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
	return configs, snapshot, stored
}

// checkFilteredCloze protects the shared template identity across all cloze card ordinals.
func checkFilteredCloze(t *testing.T, configs *ankiascodepb.CollectionConfigSet, snapshot *resource.Snapshot, stored *ankiascodepb.State, output string) {
	t.Helper()
	for _, action := range []string{"unchanged", "replace", "remove", "move"} {
		t.Run("filtered cloze template "+action, func(t *testing.T) {
			desired := proto.Clone(configs).(*ankiascodepb.CollectionConfigSet)
			file := desired.NoteTypes[0]
			if action == "replace" || action == "move" {
				name := file.Config.Templates[0].Name
				file.Config.Templates[0].Name = "Renamed"
				if action == "move" {
					desired.Config.Moved = []*ankiascodepb.ResourceMove{{From: file.Path + "#templates/" + name, To: file.Path + "#templates/Renamed"}}
				}
			} else if action == "remove" {
				file.Config.Templates = nil
			}
			prepared, err := planning.Build(desired, snapshot, stored)
			if action == "replace" || action == "remove" {
				if err == nil || prepared != nil {
					t.Fatal("resident cloze cards lost protection against template replacement")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(prepared.CardAssignments) != 0 {
				t.Fatal("resident cloze cards acquired managed assignments")
			}
			if action == "move" && (len(prepared.NoteTypes) != 1 || prepared.NoteTypes[0].Value.Templates[0].Ord == nil || prepared.NoteTypes[0].Value.Templates[0].Ord.Val != 0) {
				t.Fatal("explicit cloze template move lost its native identity")
			}
			if output != "" {
				writePlanningArtifact(t, filepath.Join(output, "filtered-cloze-"+action+".json"), prepared)
			}
		})
	}
}
