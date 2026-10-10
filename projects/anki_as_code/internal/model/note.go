package model

import collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"

// Note combines editable note content with its resolved layout and source location.
type Note struct {
	// Note contains the protobuf-defined editable note and card declarations.
	*collectionpb.Note
	// OrderedFields records field names in note-type order for reconciliation.
	OrderedFields []string
	// Bodies holds field content in note-type order for comparison and persistence.
	Bodies []string
	// Path records the note filename for identity persistence.
	Path string
	// HomeDeck records the declaring deck or export destination.
	HomeDeck string
}
