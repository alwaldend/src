package plan

import (
	"fmt"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
)

// snapshot indexes native messages without converting them to editable configs.
type snapshot struct {
	// bindings contains the expected stored state before declared moves.
	bindings *ankiascodepb.State
	// presets indexes existing unmanaged scheduling presets by native ID.
	presets map[int64]string
	// notes indexes captured native notes by ID.
	notes map[int64]*notes.Note
	// noteTypes indexes captured native definitions by ID.
	noteTypes map[int64]*notetypes.Notetype
	// decks indexes captured normal and filtered decks by ID.
	decks map[int64]*decks.Deck
	// notePaths locates tracked native notes; absent IDs are untracked.
	notePaths map[int64]string
	// typePaths locates tracked native definitions.
	typePaths map[int64]string
	// deckPaths locates tracked native deck markers.
	deckPaths map[int64]string
	// deckNames indexes folded native names for collision checks.
	deckNames map[string][]int64
	// typeNames indexes folded native names for untracked-type protection.
	typeNames map[string][]int64
	// cardDeckIDs indexes normal-card destinations by note ID and ordinal.
	cardDeckIDs map[cardKey]int64
	// occupiedDecks records resident and filtered-home deck references.
	occupiedDecks map[int64]bool
	// filteredCards retains native references needed to protect unmanaged cards.
	filteredCards []*cards.Card
}

// readSnapshot indexes supplied native values and bindings without mutating either input.
func readSnapshot(native *resource.Snapshot, bindings *ankiascodepb.State) (snapshot, error) {
	current := snapshot{
		bindings:      bindings,
		presets:       native.Presets,
		notes:         map[int64]*notes.Note{},
		noteTypes:     map[int64]*notetypes.Notetype{},
		decks:         map[int64]*decks.Deck{},
		notePaths:     map[int64]string{},
		typePaths:     map[int64]string{},
		deckPaths:     map[int64]string{},
		deckNames:     map[string][]int64{},
		typeNames:     map[string][]int64{},
		cardDeckIDs:   map[cardKey]int64{},
		occupiedDecks: map[int64]bool{},
	}
	for address, deck := range bindings.GetResources().GetDecks() {
		current.deckPaths[deck.Id] = address
	}
	for address, note := range bindings.GetResources().GetNotes() {
		current.notePaths[note.Id] = address
	}
	for address, definition := range bindings.GetResources().GetNoteTypes() {
		current.typePaths[definition.Id] = address
	}
	for _, definition := range native.NoteTypes {
		current.noteTypes[definition.Id] = definition
		name := folded(definition.Name)
		current.typeNames[name] = append(current.typeNames[name], definition.Id)
	}
	for _, deck := range native.Decks {
		current.decks[deck.Id] = deck
		name := folded(deck.Name)
		current.deckNames[name] = append(current.deckNames[name], deck.Id)
	}
	for _, note := range native.Notes {
		definition := current.noteTypes[note.NotetypeId]
		if definition == nil || len(note.Fields) != len(definition.Fields) {
			return current, fmt.Errorf("note %d has an unavailable field layout", note.Id)
		}
		current.notes[note.Id] = note
	}
	for _, card := range native.Cards {
		if current.notes[card.NoteId] == nil {
			return current, fmt.Errorf("card %d references unavailable note %d", card.Id, card.NoteId)
		}
		deck := current.decks[card.DeckId]
		if deck == nil {
			return current, fmt.Errorf("card %d references unavailable deck %d", card.Id, card.DeckId)
		}
		current.occupiedDecks[card.DeckId] = true
		if card.OriginalDeckId != 0 {
			current.occupiedDecks[card.OriginalDeckId] = true
		}
		if card.OriginalDeckId != 0 || deck.GetNormal() == nil {
			current.filteredCards = append(current.filteredCards, card)
			continue
		}
		current.cardDeckIDs[cardKey{noteID: card.NoteId, ordinal: int64(card.TemplateIdx)}] = card.DeckId
	}
	return current, nil
}
