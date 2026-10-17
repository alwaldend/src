package plan

import (
	"fmt"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
)

// validateDesired protects references that Anki would otherwise migrate or remove.
func validateDesired(desired desiredState, current snapshot) error {
	notesByID := map[int64]runtimeNote{}
	typesByID := map[int64]*runtimeNoteType{}
	for _, note := range desired.notes {
		if note.id != 0 {
			notesByID[note.id] = note
		}
	}
	for _, definition := range desired.noteTypes {
		if definition.id != 0 {
			typesByID[definition.id] = definition
		}
	}
	removedDecks, err := validateDeckOwnership(desired, current)
	if err != nil {
		return fmt.Errorf("protect decks: %w", err)
	}
	if err := validateCardDestinations(notesByID, current, removedDecks); err != nil {
		return fmt.Errorf("protect card destinations: %w", err)
	}
	if err := validateNoteOwnership(typesByID, current); err != nil {
		return fmt.Errorf("protect untracked notes: %w", err)
	}
	if err := validateFilteredCards(desired, notesByID, current, removedDecks); err != nil {
		return fmt.Errorf("protect filtered cards: %w", err)
	}
	return nil
}

// validateDeckOwnership protects untracked descendants and template targets, returning removed decks.
func validateDeckOwnership(desired desiredState, current snapshot) (map[int64]string, error) {
	retainedDecks := map[int64]runtimeDeck{}
	for _, deck := range desired.decks {
		if deck.id != 0 {
			retainedDecks[deck.id] = deck
		}
	}
	removedDecks := map[int64]string{}
	for _, old := range current.decks {
		if old.path == "" {
			continue
		}
		if _, retained := retainedDecks[old.id]; !retained {
			removedDecks[old.id] = old.path
		}
		next := retainedDecks[old.id].title
		if next != old.title {
			for _, child := range current.decks {
				if child.path == "" && strings.HasPrefix(folded(child.title), folded(old.title)+"::") {
					return nil, fmt.Errorf("deck %q has an untracked descendant %q; import it before renaming or deleting the parent", old.path, child.title)
				}
			}
		}
	}
	for _, definition := range current.noteTypes {
		if definition.path != "" {
			continue
		}
		for _, template := range definition.definition.Templates {
			id := template.Config.TargetDeckId
			if _, removing := removedDecks[id]; removing {
				return nil, fmt.Errorf("untracked note type %q targets removed deck %q", definition.definition.Name, removedDecks[id])
			}
			if deck, retained := retainedDecks[id]; retained && deck.config.Normal == nil {
				return nil, fmt.Errorf("untracked note type %q requires normal deck %q", definition.definition.Name, deck.path)
			}
		}
	}
	return removedDecks, nil
}

// validateCardDestinations requires a surviving destination before a retained card's deck is removed.
func validateCardDestinations(notesByID map[int64]runtimeNote, current snapshot, removedDecks map[int64]string) error {
	for key, deckID := range current.cardDeckIDs {
		deckPath, removing := removedDecks[deckID]
		if !removing {
			continue
		}
		note, retained := notesByID[key.noteID]
		if !retained {
			if current.notes[key.noteID].path == "" {
				return fmt.Errorf("untracked note %d has cards in deck %q", key.noteID, deckPath)
			}
			continue
		}
		destination := cardDestination(note, key.ordinal)
		if destination == "" {
			return fmt.Errorf("retained note %q has a card in removed deck %q; declare a surviving normal-deck destination", note.path, deckPath)
		}
	}
	return nil
}

// validateNoteOwnership rejects type changes that would migrate untracked notes or their cards.
func validateNoteOwnership(typesByID map[int64]*runtimeNoteType, current snapshot) error {
	for id, note := range current.notes {
		if note.path != "" {
			continue
		}
		previous := current.noteTypes[note.typeID]
		wanted := typesByID[note.typeID]
		if wanted == nil && previous.path != "" {
			return fmt.Errorf("untracked note %d still uses note type %q", id, previous.path)
		}
		if wanted != nil && !sameCardStructure(previous.definition, wanted) {
			return fmt.Errorf("note type %q has untracked notes; import them before changing fields or card generation", wanted.path)
		}
	}
	return nil
}

// validateFilteredCards leaves resident filtered cards and both of their deck references unmanaged.
func validateFilteredCards(desired desiredState, notesByID map[int64]runtimeNote, current snapshot, removedDecks map[int64]string) error {
	for _, card := range current.filteredCards {
		note, exists := notesByID[card.noteID]
		if !exists && current.notes[card.noteID].path != "" {
			return fmt.Errorf("cannot remove note with unmanaged filtered cards; empty the filtered deck in Anki first")
		}
		if exists {
			for _, declared := range note.config.Cards {
				if declared.Ordinal == int64(card.ordinal) {
					return fmt.Errorf("note %q already has an unmanaged filtered card at ordinal %d", note.path, card.ordinal)
				}
			}
		}
		if err := validateFilteredTemplate(current.notes[card.noteID].typeID, card.ordinal, desired.noteTypes, current.noteTypes); err != nil {
			return fmt.Errorf("protect filtered card %d: %w", card.id, err)
		}
		for _, id := range []int64{card.deckID, card.homeDeckID} {
			if _, removing := removedDecks[id]; id != 0 && removing {
				return fmt.Errorf("cannot remove a deck referenced by unmanaged filtered card %d; empty it in Anki first", card.id)
			}
		}
	}
	return nil
}

// sameCardStructure checks schema and generation edits that would mutate untracked notes/cards.
func sameCardStructure(old *collectionpb.NoteTypeConfig, desired *runtimeNoteType) bool {
	wanted := desired.definition
	if old.Config.Kind != wanted.Config.Kind || len(old.Fields) != len(wanted.Fields) || len(old.Templates) != len(wanted.Templates) {
		return false
	}
	for index, field := range old.Fields {
		if field.Name != wanted.Fields[index].Name || desired.native.Fields[index].Ord == nil || desired.native.Fields[index].Ord.Val != uint32(index) {
			return false
		}
	}
	for index, template := range old.Templates {
		next := wanted.Templates[index]
		if template.Name != next.Name || template.Config.QFormat != next.Config.QFormat || desired.native.Templates[index].Ord == nil || desired.native.Templates[index].Ord.Val != uint32(index) {
			return false
		}
	}
	return true
}
