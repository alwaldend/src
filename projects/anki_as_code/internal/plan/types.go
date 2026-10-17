package plan

import collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"

// desiredState holds prepared operations and shared indexes for resolved declarations.
type desiredState struct {
	// bindings contains moved stored bindings used to resolve existing identities.
	bindings *collectionpb.Resources
	// decksByPath resolves declared deck references, including new decks with zero IDs.
	decksByPath map[string]*collectionpb.PlannedDeck
	// decksByID indexes retained decks for ownership checks and comparison.
	decksByID map[int64]*collectionpb.PlannedDeck
	// noteTypes indexes prepared definitions by their configuration paths.
	noteTypes map[string]*collectionpb.PlannedNoteType
	// typesByID indexes retained definitions without rebuilding ownership maps.
	typesByID map[int64]*collectionpb.PlannedNoteType
	// notes indexes prepared content by its configuration path.
	notes map[string]*collectionpb.PlannedNote
	// notesByID indexes retained notes for ownership checks and card comparison.
	notesByID map[int64]*collectionpb.PlannedNote
	// cards holds declaring-deck defaults and explicit destinations for each note.
	cards map[string]*collectionpb.PlannedCardAssignments
}

// cardKey identifies an existing generated card by its owning note and ordinal.
type cardKey struct {
	// noteID is the native identity of the card's owning note.
	noteID int64
	// ordinal is the template index, or cloze number minus one.
	ordinal int64
}
