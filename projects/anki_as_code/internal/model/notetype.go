package model

import collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"

// NoteType combines a native identity with its configuration and state address.
type NoteType struct {
	// ID is the native Anki identity kept in collection state.
	ID int64
	// Path is the bound configuration file, or empty for an untracked type.
	Path string
	// Definition contains values used by native reconciliation and export projection.
	Definition *collectionpb.NoteTypeConfig
}
