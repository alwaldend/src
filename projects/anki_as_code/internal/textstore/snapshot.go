package textstore

import (
	"fmt"
	"os"
	"path/filepath"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"google.golang.org/protobuf/proto"
)

// Export stages a complete snapshot and preserves unrelated destination files.
func Export(s model.CollectionState, output string) error {
	existing := false
	if info, err := os.Lstat(output); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("export destination %q must be a directory", output)
		}
		existing = true
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("check export destination: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("create export parent: %w", err)
	}
	tmp, err := os.MkdirTemp(filepath.Dir(output), ".anki-export-")
	if err != nil {
		return fmt.Errorf("create export staging: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(tmp)
	}()
	if existing {
		if err := os.CopyFS(tmp, os.DirFS(output)); err != nil {
			return fmt.Errorf("stage existing export: %w", err)
		}
		if err := RemoveExportDeclarations(tmp); err != nil {
			return fmt.Errorf("refresh existing declarations: %w", err)
		}
	}
	for _, deck := range s.Decks {
		if err := WriteDeck(tmp, deck); err != nil {
			return fmt.Errorf("export deck: %w", err)
		}
	}
	manifest := &collectionpb.CollectionConfig{
		FormatVersion: 2,
		DecksPath:     "decks",
		NoteTypesPath: "note_types",
		NotesPath:     proto.String("notes"),
	}
	if err := WriteTOML(filepath.Join(tmp, CollectionMarker), manifest); err != nil {
		return fmt.Errorf("write export root config: %w", err)
	}

	if err := WriteNoteTypes(filepath.Join(tmp, manifest.NoteTypesPath), s.Appearance); err != nil {
		return fmt.Errorf("export note types: %w", err)
	}

	for _, n := range s.Notes {
		if len(n.Cards) > 0 {
			n.HomeDeck = n.Cards[0].GetDeck()
		} else {
			if err := os.MkdirAll(filepath.Join(tmp, "notes"), 0o755); err != nil {
				return fmt.Errorf("create cardless notes directory: %w", err)
			}
		}
		if err := WriteNote(NotePath(tmp, n.HomeDeck, n.NoteId), n); err != nil {
			return fmt.Errorf("export note: %w", err)
		}
	}

	if err := PublishExport(tmp, output, existing); err != nil {
		return fmt.Errorf("publish text export: %w", err)
	}

	return nil
}
