package layout

import (
	"fmt"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"google.golang.org/protobuf/proto"
)

// Export writes supplied configurations as TOML files in the fixed hierarchy.
// Callers own validation and unique, representable paths. Declared files are
// overwritten; other files are untouched. Failures may leave partial output.
func Export(storage filesystem.FS, configs *ankiascodepb.CollectionConfigSet, output string) error {
	output, err := storage.ExpandPath(output)
	if err != nil {
		return fmt.Errorf("resolve configuration output: %w", err)
	}
	if err := exportFile(storage, filepath.Join(output, collectionMarker), configs.GetConfig()); err != nil {
		return fmt.Errorf("export collection config: %w", err)
	}
	for _, deck := range configs.GetDecks() {
		if err := exportFile(storage, filepath.Join(output, filepath.FromSlash(deck.Path)), deck.Deck); err != nil {
			return fmt.Errorf("export deck %q: %w", deck.Path, err)
		}
		if err := exportNotes(storage, output, deck.Notes); err != nil {
			return fmt.Errorf("export deck notes: %w", err)
		}
	}
	if err := exportNotes(storage, output, configs.Notes); err != nil {
		return fmt.Errorf("export root notes: %w", err)
	}
	for _, definition := range configs.NoteTypes {
		if err := exportFile(storage, filepath.Join(output, filepath.FromSlash(definition.Path)), definition.Config); err != nil {
			return fmt.Errorf("export note type %q: %w", definition.Path, err)
		}
	}
	return nil
}

// exportNotes writes each supplied note at its collection-relative path.
func exportNotes(storage filesystem.FS, root string, notes []*ankiascodepb.NoteConfigFile) error {
	for _, note := range notes {
		if err := exportFile(storage, filepath.Join(root, filepath.FromSlash(note.Path)), note.Config); err != nil {
			return fmt.Errorf("export note %q: %w", note.Path, err)
		}
	}
	return nil
}

// exportFile creates parent directories and writes one encoded configuration.
func exportFile(storage filesystem.FS, path string, message proto.Message) error {
	if err := storage.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config directory for %q: %w", path, err)
	}
	if err := writeTOML(storage, path, message); err != nil {
		return fmt.Errorf("export config file %q: %w", path, err)
	}
	return nil
}
