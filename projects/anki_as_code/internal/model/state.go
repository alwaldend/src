package model

import collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"

// CollectionState is the current snapshot shared by export, validation, and planning.
type CollectionState struct {
	// Bindings contains the path-to-native mappings persisted in the collection.
	Bindings *collectionpb.ReconciliationState
	// Presets indexes existing scheduling preset names by native identity; they remain unmanaged.
	Presets map[int64]string
	// Notes indexes notes by GUID, excluding filtered cards from their declarations.
	Notes map[string]Note
	// NoteTypes indexes existing note types and their configurations by identity.
	NoteTypes map[int64]*NoteType
	// Decks contains the current normal and filtered deck declarations.
	Decks []Deck
	// CardDeckIDs indexes the current database deck identity of each managed note/ordinal pair.
	CardDeckIDs map[CardKey]int64
	// OccupiedDecks records resident card occupancy for validating deck-kind changes.
	OccupiedDecks map[int64]bool
}
