package plan

import (
	"fmt"
	"strings"

	"github.com/ankitects/anki/proto/anki/notetypes"
)

// validateDesired protects references that Anki would otherwise migrate or remove.
func validateDesired(desired desiredState, current snapshot) error {
	removedDecks, err := validateDeckOwnership(desired, current)
	if err != nil {
		return fmt.Errorf("protect decks: %w", err)
	}
	if err := validateCardDestinations(desired, current, removedDecks); err != nil {
		return fmt.Errorf("protect card destinations: %w", err)
	}
	if err := validateNoteOwnership(desired, current); err != nil {
		return fmt.Errorf("protect untracked notes: %w", err)
	}
	if err := validateFilteredCards(desired, current, removedDecks); err != nil {
		return fmt.Errorf("protect filtered cards: %w", err)
	}
	return nil
}

// validateDeckOwnership protects required decks, untracked descendants and template targets, returning removed decks.
func validateDeckOwnership(desired desiredState, current snapshot) (map[int64]string, error) {
	removedDecks := map[int64]string{}
	for _, old := range current.decks {
		path := current.deckPaths[old.Id]
		if path == "" {
			continue
		}
		next := desired.decksByID[old.Id]
		if next == nil {
			if old.Id == 1 {
				return nil, fmt.Errorf("required default deck %q (native ID 1) cannot be removed; retain its declaration or move it explicitly", path)
			}
			removedDecks[old.Id] = path
		}
		if next == nil || next.Value.Name != old.Name {
			for _, child := range current.decks {
				if current.deckPaths[child.Id] == "" && strings.HasPrefix(folded(child.Name), folded(old.Name)+"::") {
					return nil, fmt.Errorf("deck %q has an untracked descendant %q; import it before renaming or deleting the parent", path, child.Name)
				}
			}
		}
	}
	for _, definition := range current.noteTypes {
		if current.typePaths[definition.Id] != "" {
			continue
		}
		for _, template := range definition.Templates {
			id := template.Config.TargetDeckId
			if _, removing := removedDecks[id]; removing {
				return nil, fmt.Errorf("untracked note type %q targets removed deck %q", definition.Name, removedDecks[id])
			}
			if deck := desired.decksByID[id]; deck != nil && deck.Value.GetNormal() == nil {
				return nil, fmt.Errorf("untracked note type %q requires normal deck %q", definition.Name, deck.Path)
			}
		}
	}
	return removedDecks, nil
}

// validateCardDestinations requires a surviving destination before a retained card's deck is removed.
func validateCardDestinations(desired desiredState, current snapshot, removedDecks map[int64]string) error {
	for key, deckID := range current.cardDeckIDs {
		deckPath, removing := removedDecks[deckID]
		if !removing {
			continue
		}
		note := desired.notesByID[key.noteID]
		if note == nil {
			if current.notePaths[key.noteID] == "" {
				return fmt.Errorf("untracked note %d has cards in deck %q", key.noteID, deckPath)
			}
			continue
		}
		if cardDestination(desired.cards[note.Path], key.ordinal) == "" {
			return fmt.Errorf("retained note %q has a card in removed deck %q; declare a surviving normal-deck destination", note.Path, deckPath)
		}
	}
	return nil
}

// validateNoteOwnership rejects type changes that would migrate untracked notes or their cards.
func validateNoteOwnership(desired desiredState, current snapshot) error {
	for id, note := range current.notes {
		if current.notePaths[id] != "" {
			continue
		}
		previous := current.noteTypes[note.NotetypeId]
		wanted := desired.typesByID[note.NotetypeId]
		if wanted == nil && current.typePaths[previous.Id] != "" {
			return fmt.Errorf("untracked note %d still uses note type %q", id, current.typePaths[previous.Id])
		}
		if wanted != nil && !sameCardStructure(previous, wanted.Value) {
			return fmt.Errorf("note type %q has untracked notes; import them before changing fields or card generation", wanted.Path)
		}
	}
	return nil
}

// validateFilteredCards leaves resident filtered cards and both of their deck references unmanaged.
func validateFilteredCards(desired desiredState, current snapshot, removedDecks map[int64]string) error {
	for _, card := range current.filteredCards {
		note := desired.notesByID[card.NoteId]
		if note == nil && current.notePaths[card.NoteId] != "" {
			return fmt.Errorf("cannot remove note with unmanaged filtered cards; empty the filtered deck in Anki first")
		}
		if note != nil {
			for _, declared := range desired.cards[note.Path].Cards {
				if declared.Ordinal == int64(card.TemplateIdx) {
					return fmt.Errorf("note %q already has an unmanaged filtered card at ordinal %d", note.Path, card.TemplateIdx)
				}
			}
		}
		if err := validateFilteredTemplate(current.notes[card.NoteId].NotetypeId, int(card.TemplateIdx), desired, current); err != nil {
			return fmt.Errorf("protect filtered card %d: %w", card.Id, err)
		}
		for _, id := range []int64{card.DeckId, card.OriginalDeckId} {
			if _, removing := removedDecks[id]; id != 0 && removing {
				return fmt.Errorf("cannot remove a deck referenced by unmanaged filtered card %d; empty it in Anki first", card.Id)
			}
		}
	}
	return nil
}

// sameCardStructure checks schema and generation edits that would mutate untracked notes/cards.
func sameCardStructure(old, wanted *notetypes.Notetype) bool {
	if old.Config.Kind != wanted.Config.Kind || len(old.Fields) != len(wanted.Fields) || len(old.Templates) != len(wanted.Templates) {
		return false
	}
	for index, field := range old.Fields {
		next := wanted.Fields[index]
		if field.Name != next.Name || next.Ord == nil || next.Ord.Val != uint32(index) {
			return false
		}
	}
	for index, template := range old.Templates {
		next := wanted.Templates[index]
		if template.Name != next.Name || template.Config.QFormat != next.Config.QFormat || next.Ord == nil || next.Ord.Val != uint32(index) {
			return false
		}
	}
	return true
}
