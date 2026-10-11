package validation_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/ankiformat"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/validation"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// TestDesiredStateWorkflow verifies resource identities and deck-kind encoding through public packages.
func TestDesiredStateWorkflow(t *testing.T) {
	settings := &collectionpb.NoteTypeSettings{
		Id: 10, Name: "Basic", Config: &notetypes.Notetype_Config{},
		Fields:    []*notetypes.Notetype_Field{{Name: "Front", Config: &notetypes.Notetype_Field_Config{}}},
		Templates: []*collectionpb.TemplateSettings{{Name: "Forward", Config: &notetypes.Notetype_Template_Config{}}},
	}
	deck := model.Deck{Deck: &collectionpb.Deck{Id: 1, Title: "Default", Normal: &decks.Deck_Normal{ConfigId: 1}}}
	encoded, err := ankiformat.EncodeDeck(deck.Deck)
	if err != nil {
		t.Fatal(err)
	}
	deck.Kind = encoded
	current := model.CollectionState{
		NoteTypes: map[int64]*collectionpb.NoteTypeSettings{10: settings},
		Decks:     []model.Deck{deck},
	}
	results := map[string]bool{}
	for _, scenario := range []struct {
		// name identifies the desired-state identity scenario.
		name string
		// secondID assigns the second note's identity independently of its GUID and card.
		secondID int64
		// reject records whether validation must refuse the desired state.
		reject bool
	}{
		{"distinct-identities", 21, false},
		{"duplicate-identity", 20, true},
		{"zero-identity", 0, true},
		{"negative-identity", -1, true},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			desired := model.CollectionDesiredState{
				Decks: []model.Deck{deck}, NoteTypes: []*collectionpb.NoteTypeSettings{settings},
				Notes: []model.Note{
					{Note: &collectionpb.Note{
						NoteId: 20, Guid: "first", NoteTypeId: 10, NoteType: "Basic", Fields: map[string]string{"Front": "first"},
						Cards: []*collectionpb.Card{{Id: 100, Deck: proto.String("Default")}},
					}},
					{Note: &collectionpb.Note{
						NoteId: scenario.secondID, Guid: "second", NoteTypeId: 10, NoteType: "Basic", Fields: map[string]string{"Front": "second"},
						Cards: []*collectionpb.Card{{Id: 101, Deck: proto.String("Default")}},
					}},
				},
			}
			if err := validation.ValidateDesired(&desired, current); (err != nil) != scenario.reject {
				t.Fatalf("reject=%v, validation error=%v", scenario.reject, err)
			}
			results[scenario.name] = true
		})
	}
	t.Run("conflicting-deck-kinds", func(t *testing.T) {
		conflict := &collectionpb.Deck{Normal: &decks.Deck_Normal{}, Filtered: &decks.Deck_Filtered{}}
		if _, err := ankiformat.EncodeDeck(conflict); err == nil {
			t.Fatal("normal and explicitly empty filtered declarations must conflict")
		}
		results["conflicting-deck-kinds"] = true
	})
	t.Run("empty-filtered-round-trip", func(t *testing.T) {
		filtered := &collectionpb.Deck{Filtered: &decks.Deck_Filtered{}}
		encoded, err := ankiformat.EncodeDeck(filtered)
		if err != nil {
			t.Fatal(err)
		}
		decoded, err := ankiformat.DecodeDeck(encoded)
		if err != nil {
			t.Fatal(err)
		}
		if decoded.Normal != nil || decoded.Filtered == nil || !proto.Equal(decoded.Filtered, filtered.Filtered) {
			t.Fatal("empty filtered configuration lost its kind or presence")
		}
		results["empty-filtered-round-trip"] = true
	})
	if t.Failed() {
		return
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "validation-workflow.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestDeckPreservationWorkflow checks nested preservation and refusal of unknown-data loss.
func TestDeckPreservationWorkflow(t *testing.T) {
	unknown := []byte{0xa0, 0x06, 0x07}
	results := map[string]bool{}
	for _, scenario := range []struct {
		// name identifies the configuration edit being exercised.
		name string
		// current supplies the original upstream deck configuration.
		current *collectionpb.Deck
		// desired supplies the edited readable declaration.
		desired *collectionpb.Deck
		// reject indicates an edit would discard unknown nested data.
		reject bool
	}{
		{"removed-search-term", &collectionpb.Deck{Filtered: &decks.Deck_Filtered{SearchTerms: []*decks.Deck_Filtered_SearchTerm{{}, {}}}}, &collectionpb.Deck{Filtered: &decks.Deck_Filtered{SearchTerms: []*decks.Deck_Filtered_SearchTerm{{}}}}, true},
		{"removed-daily-limit", &collectionpb.Deck{Normal: &decks.Deck_Normal{ReviewLimitToday: &decks.Deck_Normal_DayLimit{}}}, &collectionpb.Deck{Normal: &decks.Deck_Normal{}}, true},
		{"edited-search-term", &collectionpb.Deck{Filtered: &decks.Deck_Filtered{SearchTerms: []*decks.Deck_Filtered_SearchTerm{{}}}}, &collectionpb.Deck{Filtered: &decks.Deck_Filtered{SearchTerms: []*decks.Deck_Filtered_SearchTerm{{Search: "tag:language"}}}}, false},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			if scenario.current.Filtered != nil {
				terms := scenario.current.Filtered.SearchTerms
				terms[len(terms)-1].ProtoReflect().SetUnknown(unknown)
			} else {
				scenario.current.Normal.ReviewLimitToday.ProtoReflect().SetUnknown(unknown)
			}
			currentKind, err := ankiformat.EncodeDeck(scenario.current)
			if err != nil {
				t.Fatal(err)
			}
			desiredKind, err := ankiformat.EncodeDeck(scenario.desired)
			if err != nil {
				t.Fatal(err)
			}
			desired := model.Deck{Deck: scenario.desired, Kind: desiredKind}
			err = ankiformat.PreserveDeckKind(model.Deck{Deck: scenario.current, Kind: currentKind}, &desired)
			if (err != nil) != scenario.reject {
				t.Fatalf("reject=%v, preservation error=%v", scenario.reject, err)
			}
			if !scenario.reject {
				kind, err := ankiformat.DecodeDeckKind(desired.Kind)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(kind.GetFiltered().SearchTerms[0].ProtoReflect().GetUnknown(), unknown) {
					t.Fatal("nested unknown data lost")
				}
			}
			results[scenario.name] = true
		})
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" && !t.Failed() {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "deck-preservation-workflow.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
