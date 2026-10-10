package textstore

import (
	"fmt"
	"os"
	"path/filepath"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
)

// readNoteTypes discovers and reads note-type appearance declarations.
func readNoteTypes(root string) ([]*collectionpb.NoteTypeAppearance, error) {
	files, err := discoverFiles(root, NoteTypeSuffix, false)
	if err != nil {
		return nil, fmt.Errorf("discover note type settings: %w", err)
	}
	appearances := make([]*collectionpb.NoteTypeAppearance, 0, len(files))
	for _, file := range files {
		appearance := &collectionpb.NoteTypeAppearance{}
		if err := ReadTOML(file, appearance); err != nil {
			return nil, fmt.Errorf("read note type settings: %w", err)
		}
		appearances = append(appearances, appearance)
	}
	return appearances, nil
}

// WriteNoteTypes writes the collection note-type appearance declarations.
func WriteNoteTypes(root string, appearances []*collectionpb.NoteTypeAppearance) error {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return fmt.Errorf("create note types directory %q: %w", root, err)
	}
	seen := map[string]bool{}
	for _, appearance := range appearances {
		path := filepath.Join(root, escapeDeckPart(appearance.Name)+NoteTypeSuffix)
		if seen[path] {
			return fmt.Errorf("duplicate note type filename %q", path)
		}
		seen[path] = true
		if err := WriteAppearance(path, appearance); err != nil {
			return fmt.Errorf("write note type %d: %w", appearance.Id, err)
		}
	}
	return nil
}

// WriteAppearance writes one note-type appearance declaration.
func WriteAppearance(path string, appearance *collectionpb.NoteTypeAppearance) error {
	if err := WriteTOML(path, appearance); err != nil {
		return fmt.Errorf("write note type appearance %q: %w", path, err)
	}
	return nil
}
