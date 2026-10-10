package model

import collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"

// Deck combines an editable deck message with persistence metadata.
type Deck struct {
	// Deck contains the protobuf-defined editable deck declaration.
	*collectionpb.Deck
	// Path records the marker filename and is omitted from serialization.
	Path string
	// Kind carries the upstream deck-kind protobuf between persistence and encoding.
	Kind string
}
