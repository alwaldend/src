package textstore

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// RemoveExportDeclarations removes managed declarations from an existing export directory.
func RemoveExportDeclarations(root string) error {
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk staged export %q: %w", path, err)
		}
		if entry.IsDir() {
			return nil
		}
		name := entry.Name()
		if name == CollectionMarker || name == DeckMarker || strings.HasSuffix(name, noteSuffix) || strings.HasSuffix(name, NoteTypeSuffix) {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("remove stale declaration %q: %w", path, err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("clear staged declarations: %w", err)
	}
	return nil
}

// PublishExport publishes staged declarations while retaining unrelated destination files.
func PublishExport(staging, output string, existing bool) error {
	if !existing {
		if err := os.Rename(staging, output); err != nil {
			return fmt.Errorf("publish export: %w", err)
		}
		return nil
	}
	backup, err := os.MkdirTemp(filepath.Dir(output), ".anki-export-backup-")
	if err != nil {
		return fmt.Errorf("reserve export backup: %w", err)
	}
	if err := os.Remove(backup); err != nil {
		return fmt.Errorf("prepare export backup: %w", err)
	}
	if err := os.Rename(output, backup); err != nil {
		return fmt.Errorf("retain previous export: %w", err)
	}
	if err := os.Rename(staging, output); err != nil {
		restoreErr := os.Rename(backup, output)
		return fmt.Errorf("publish refreshed export (previous data retained at %q): %w", backup, errors.Join(err, restoreErr))
	}
	if err := os.RemoveAll(backup); err != nil {
		return fmt.Errorf("remove previous export backup %q: %w", backup, err)
	}
	return nil
}
