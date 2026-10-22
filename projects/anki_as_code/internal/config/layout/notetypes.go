package layout

import (
	"fmt"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
)

// readNoteTypes discovers and reads note-type configurations.
func readNoteTypes(storage filesystem.FS, root, collectionRoot string) ([]*ankiascodepb.NoteTypeConfigFile, error) {
	files, err := discoverFiles(storage, root, noteTypeSuffix)
	if err != nil {
		return nil, fmt.Errorf("discover note-type configuration: %w", err)
	}
	noteTypes := make([]*ankiascodepb.NoteTypeConfigFile, 0, len(files))
	for _, file := range files {
		config := &ankiascodepb.NoteTypeConfig{}
		if err := readTOML(storage, file, config); err != nil {
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
