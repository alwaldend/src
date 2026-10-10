package model

import collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"

// CollectionDesiredState contains the complete managed collection state loaded from text files.
type CollectionDesiredState struct {
	// Config declares the text format and managed resource locations.
	Config *collectionpb.CollectionConfig
	// Notes contains the complete desired set of notes.
	Notes []Note
	// Decks contains the complete desired set of deck declarations.
	Decks []Deck
	// Appearance contains the desired CSS and existing card-template settings.
	Appearance []*collectionpb.NoteTypeAppearance
}
