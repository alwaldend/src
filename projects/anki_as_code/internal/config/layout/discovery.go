package layout

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
)

const (
	// Version identifies the current unreleased TOML configuration schema.
	Version = 1
	// collectionMarker is the default root config filename used for directory inputs.
	collectionMarker = "collection.anki.toml"
	// deckMarker is the required filename in every deck directory.
	deckMarker = "deck.anki.toml"
	// noteSuffix identifies editable note declaration filenames.
	noteSuffix = ".note.anki.toml"
	// noteTypeSuffix identifies editable note-type definition filenames.
	noteTypeSuffix = ".notetype.anki.toml"
)

// regularPath rejects paths containing symlinks or unsuitable filesystem entries.
func regularPath(storage filesystem.FS, path string) error {
	resolved, err := storage.ResolvePath(path, false)
	if err != nil {
		return fmt.Errorf("resolve %q: %w", path, err)
	}
	if resolved != path {
		return fmt.Errorf("symlink in path %q", path)
	}
	info, err := storage.Lstat(path)
	if err != nil {
		return fmt.Errorf("inspect %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%q must be a regular file", path)
	}
	return nil
}

// collectionConfig resolves a root config file from a file or directory input.
func collectionConfig(storage filesystem.FS, input string) (string, error) {
	path, err := storage.ResolvePathWithoutLinks(input, true)
	if err != nil {
		return "", fmt.Errorf("resolve collection path: %w", err)
	}
	info, err := storage.Lstat(path)
	if err != nil {
		return "", fmt.Errorf("inspect collection path %q: %w", path, err)
	}
	if info.IsDir() {
		path = filepath.Join(path, collectionMarker)
	}
	if err := regularPath(storage, path); err != nil {
		return "", fmt.Errorf("inspect root config: %w", err)
	}
	return path, nil
}

// directoryEntries reads a fixed resource directory; missing directories represent empty sets.
// Managed entries cannot be symlinks, including dangling links.
func directoryEntries(storage filesystem.FS, root string) ([]os.FileInfo, error) {
	info, err := storage.Lstat(root)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("inspect resource directory %q: %w", root, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("resource path %q must be a directory", root)
	}
	entries, err := storage.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read resource directory %q: %w", root, err)
	}
	for _, entry := range entries {
		if entry.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("symlink in resource directory %q", filepath.Join(root, entry.Name()))
		}
	}
	return entries, nil
}

// discoverFiles lists direct declarations and ignores unrelated files and subdirectories.
func discoverFiles(storage filesystem.FS, root, suffix string) ([]string, error) {
	entries, err := directoryEntries(storage, root)
	if err != nil {
		return nil, fmt.Errorf("discover declarations: %w", err)
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), suffix) {
			continue
		}
		path := filepath.Join(root, entry.Name())
		if !entry.Mode().IsRegular() {
			return nil, fmt.Errorf("declaration %q must be a regular file", path)
		}
		files = append(files, path)
	}
	return files, nil
}
