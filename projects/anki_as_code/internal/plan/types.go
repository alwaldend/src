package plan

import ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"

// desiredState holds prepared operations and shared indexes for resolved declarations.
type desiredState struct {
	// bindings contains moved stored bindings used to resolve existing identities.
	bindings *ankiascodepb.Resources
	// decksByPath resolves declared deck references, including new decks with zero IDs.
	decksByPath map[string]*ankiascodepb.PlannedDeck
	// decksByID indexes retained decks for ownership checks and comparison.
	decksByID map[int64]*ankiascodepb.PlannedDeck
	// noteTypes indexes prepared definitions by their configuration paths.
	noteTypes map[string]*ankiascodepb.PlannedNoteType
	// typesByID indexes retained definitions without rebuilding ownership maps.
	typesByID map[int64]*ankiascodepb.PlannedNoteType
	// notes indexes prepared content by its configuration path.
	notes map[string]*ankiascodepb.PlannedNote
	// notesByID indexes retained notes for ownership checks and card comparison.
	notesByID map[int64]*ankiascodepb.PlannedNote
	// cards holds declaring-deck defaults and explicit destinations for each note.
	cards map[string]*ankiascodepb.PlannedCardAssignments
}

// cardKey identifies an existing generated card by its owning note and ordinal.
type cardKey struct {
	// noteID is the native identity of the card's owning note.
	noteID int64
	// ordinal is the template index, or cloze number minus one.
	ordinal int64
}
