package repository

import "git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"

// snapshot keeps database preservation data private to reconciliation.
type snapshot struct {
	// model.CollectionState contains the shared managed collection state.
	model.CollectionState
	// noteTypes retains original note-type and template configuration bytes by identity.
	noteTypes map[int64]noteTypeStorage
	// newPositions records existing new-card positions by card identity.
	newPositions map[int64]int64
	// nextPosition is the stored new-card allocation cursor.
	nextPosition int64
	// filteredCards retains references needed to protect unmanaged cards in filtered decks.
	filteredCards []filteredCard
}

// noteTypeStorage retains original protobuf bytes for lossless updates.
type noteTypeStorage struct {
	// config retains the original note-type configuration protobuf bytes.
	config []byte
	// templates indexes original template configuration bytes by ordinal.
	templates map[int][]byte
}

// filteredCard contains only the references needed to protect an unmanaged card.
type filteredCard struct {
	// id identifies the unmanaged card.
	id int64
	// noteID identifies the note referenced by the card.
	noteID int64
	// ordinal identifies the card template within its note type.
	ordinal int
	// deckID identifies the filtered deck currently holding the card.
	deckID int64
	// homeDeckID identifies the original deck retained by Anki for restoration.
	homeDeckID int64
}
