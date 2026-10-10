package model

import (
	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"github.com/ankitects/anki/proto/anki/notetypes"
)

// NoteType describes an upstream Anki note-type definition, its card ordinals,
// and editable appearance shared by validation and persistence.
type NoteType struct {
	// Notetype provides the upstream identity, configuration, and ordered field definitions.
	*notetypes.Notetype
	// Ordinals records the existing template ordinals accepted by validation.
	Ordinals map[int]bool
	// Appearance contains the current editable CSS and card-template settings.
	Appearance *collectionpb.NoteTypeAppearance
}

// CollectionState is the current collection snapshot shared by export, validation,
// and planning. OccupiedDecks records deck occupancy for deck-kind validation.
type CollectionState struct {
	// Notes indexes current notes by Anki identity, excluding filtered cards from their declarations.
	Notes map[int64]Note
	// Types indexes existing upstream note types by identity.
	Types map[int64]*NoteType
	// Decks contains the current normal and filtered deck declarations.
	Decks []Deck
	// Appearance contains the current note-type appearance settings in identity order.
	Appearance []*collectionpb.NoteTypeAppearance
	// CardDeckIDs indexes the current database deck identity of each managed card.
	CardDeckIDs map[int64]int64
	// OccupiedDecks records resident card occupancy for validating deck-kind changes.
	OccupiedDecks map[int64]bool
}
