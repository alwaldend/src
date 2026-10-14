package model

import collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"

// Deck combines an editable deck message with its derived Anki title.
type Deck struct {
	// ID is the native Anki identity stored in collection state.
	ID int64
	// Path identifies the bound configuration marker, or is empty for untracked decks.
	Path string
	// Deck contains the protobuf-defined editable deck declaration.
	Deck *collectionpb.DeckConfig
	// Title is the full Anki deck name derived from its directory hierarchy.
	Title string
}
