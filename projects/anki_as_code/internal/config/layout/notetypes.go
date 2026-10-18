package layout

import (
	"fmt"
	"path/filepath"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
)

// readNoteTypes discovers and reads note-type configurations.
func readNoteTypes(root, collectionRoot string) ([]*ankiascodepb.NoteTypeConfigFile, error) {
	files, err := discoverFiles(root, noteTypeSuffix)
	if err != nil {
		return nil, fmt.Errorf("discover note-type configuration: %w", err)
	}
	noteTypes := make([]*ankiascodepb.NoteTypeConfigFile, 0, len(files))
	for _, file := range files {
		config := &ankiascodepb.NoteTypeConfig{}
		if err := readTOML(file, config); err != nil {
			return nil, fmt.Errorf("read note-type configuration: %w", err)
		}
		relative, err := filepath.Rel(collectionRoot, file)
		if err != nil {
			return nil, fmt.Errorf("identify note-type path: %w", err)
		}
		noteTypes = append(noteTypes, &ankiascodepb.NoteTypeConfigFile{Path: filepath.ToSlash(relative), Config: config})
	}
	return noteTypes, nil
}
