package plan

import (
	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"github.com/ankitects/anki/proto/anki/notetypes"
)

// desiredState holds resolved declarations indexed by their configuration paths.
type desiredState struct {
	// bindings contains moved stored bindings used to resolve existing identities.
	bindings *collectionpb.State
	// presetBindings holds aliases for existing unmanaged scheduling presets.
	presetBindings map[string]int64
	// notes indexes desired notes by path; new notes have zero native IDs.
	notes map[string]runtimeNote
	// noteTypes indexes declared note types by path.
	noteTypes map[string]*runtimeNoteType
	// decks contains declared decks ordered with parents before children.
	decks []runtimeDeck
}

// runtimeDeck combines an editable deck message with its derived Anki title.
type runtimeDeck struct {
	// id is the existing native identity, or zero for a new declaration.
	id int64
	// path identifies the bound configuration marker, or is empty for untracked decks.
	path string
	// config contains the protobuf-defined editable deck declaration.
	config *collectionpb.DeckConfig
	// title is the full Anki deck name derived from its directory hierarchy.
	title string
}

// runtimeNote combines editable configuration with ordered field values and deck ownership.
type runtimeNote struct {
	// path identifies the bound configuration file, or is empty for untracked notes.
	path string
	// id is the existing native note identity, or zero for a new declaration.
	id int64
	// typeID is the existing native type identity, or zero for a new type.
	typeID int64
	// config contains the protobuf-defined editable note and generated-card assignments.
	config *collectionpb.NoteConfig
	// bodies holds field content in note-type order for comparison and Anki API requests.
	bodies []string
	// homeDeck records the declaring deck or export destination.
	homeDeck string
}

// cardKey identifies an existing generated card by its owning note and ordinal.
type cardKey struct {
	// noteID is the native identity of the card's owning note.
	noteID int64
	// ordinal is the template index, or cloze number minus one.
	ordinal int64
}

// runtimeNoteType combines a native identity with its configuration and state address.
type runtimeNoteType struct {
	// id is the existing native type identity, or zero for a new declaration.
	id int64
	// path is the bound configuration file, or empty for an untracked type.
	path string
	// definition contains editable values used by comparison and native operation preparation.
	definition *collectionpb.NoteTypeConfig
	// native contains captured or prepared native values with resolved member ordinals.
	native *notetypes.Notetype
}
