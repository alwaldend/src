package model

import collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"

// CollectionState is the current snapshot shared by export, validation, and planning.
type CollectionState struct {
	// Notes indexes current notes, excluding filtered cards from their declarations.
	Notes map[int64]Note
	// NoteTypes indexes existing note types and their settings by identity.
	NoteTypes map[int64]*collectionpb.NoteTypeSettings
	// Decks contains the current normal and filtered deck declarations.
	Decks []Deck
	// CardDeckIDs indexes the current database deck identity of each managed card.
	CardDeckIDs map[int64]int64
	// OccupiedDecks records resident card occupancy for validating deck-kind changes.
	OccupiedDecks map[int64]bool
}
