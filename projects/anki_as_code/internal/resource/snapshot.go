package resource

import (
	"github.com/ankitects/anki/proto/anki/cards"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
)

// Snapshot holds complete native protobuf values for one resource read.
// Sequential reads require exclusive offline access; this is not a database transaction.
type Snapshot struct {
	// Decks contains normal and filtered decks in title order.
	Decks []*decks.Deck
	// NoteTypes contains native definitions in identity order.
	NoteTypes []*notetypes.Notetype
	// Notes contains shared note content in identity order.
	Notes []*notes.Note
	// Cards contains generated cards in identity order, including unmanaged filtered cards.
	Cards []*cards.Card
	// Presets maps existing scheduling identities to names without managing their values.
	Presets map[int64]string
}
