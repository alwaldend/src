package layout

import (
	"fmt"
	"path/filepath"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
)

// readNoteTypes discovers and reads note-type configurations.
func readNoteTypes(root, collectionRoot string) ([]*collectionpb.NoteTypeConfigFile, error) {
	files, err := discoverFiles(root, noteTypeSuffix)
	if err != nil {
		return nil, fmt.Errorf("discover note-type configuration: %w", err)
	}
	noteTypes := make([]*collectionpb.NoteTypeConfigFile, 0, len(files))
	for _, file := range files {
		config := &collectionpb.NoteTypeConfig{}
		if err := readTOML(file, config); err != nil {
			return nil, fmt.Errorf("read note-type configuration: %w", err)
		}
		relative, err := filepath.Rel(collectionRoot, file)
		if err != nil {
			return nil, fmt.Errorf("identify note-type path: %w", err)
		}
		noteTypes = append(noteTypes, &collectionpb.NoteTypeConfigFile{Path: filepath.ToSlash(relative), Config: config})
	}
	return noteTypes, nil
}
