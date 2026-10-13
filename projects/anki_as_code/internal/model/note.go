package model

import collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"

// Note combines editable configuration with ordered field values and deck ownership.
type Note struct {
	// Path identifies the bound configuration file, or is empty for untracked notes.
	Path string
	// GUID is Anki-owned synchronization metadata, excluded from editable configurations.
	GUID string
	// ID is the runtime Anki note identity, excluded from editable configuration.
	ID int64
	// TypeID resolves the declared note-type path to its runtime Anki identity.
	TypeID int64
	// Note contains the protobuf-defined editable note and generated-card assignments.
	Note *collectionpb.NoteConfig
	// Bodies holds field content in note-type order for comparison and Anki API requests.
	Bodies []string
	// HomeDeck records the declaring deck or export destination.
	HomeDeck string
}

// CardKey identifies a generated card independently of Anki's numeric storage ID.
type CardKey struct {
	// NoteGUID is the stable identity of the card's owning note.
	NoteGUID string
	// Ordinal is the template index, or cloze number minus one.
	Ordinal int64
}
