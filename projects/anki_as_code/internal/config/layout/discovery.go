package layout

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// IsConfigFile identifies declaration filenames for callers managing export directories.
// It performs no filesystem inspection or configuration validation.
func IsConfigFile(name string) bool {
	return name == collectionMarker || name == deckMarker || strings.HasSuffix(name, noteSuffix) || strings.HasSuffix(name, noteTypeSuffix)
}

// regularPath rejects paths containing symlinks or unsuitable filesystem entries.
func regularPath(path string) error {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return fmt.Errorf("resolve %q: %w", path, err)
	}
	if resolved != path {
		return fmt.Errorf("symlink in path %q", path)
	}
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("inspect %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%q must be a regular file", path)
	}
	return nil
}

// collectionConfig resolves a root config file from a file or directory input.
func collectionConfig(input string) (string, error) {
	path, err := expandPath(input)
	if err != nil {
		return "", fmt.Errorf("resolve collection path: %w", err)
	}
	info, err := os.Lstat(path)
	if err != nil {
		return "", fmt.Errorf("inspect collection path %q: %w", path, err)
	}
	if info.IsDir() {
		path = filepath.Join(path, collectionMarker)
	}
	if err := regularPath(path); err != nil {
		return "", fmt.Errorf("inspect root config: %w", err)
	}
	return path, nil
}

// directoryEntries reads a fixed resource directory; missing directories represent empty sets.
// Managed entries cannot be symlinks, including dangling links.
func directoryEntries(root string) ([]os.DirEntry, error) {
	info, err := os.Lstat(root)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("inspect resource directory %q: %w", root, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("resource path %q must be a directory", root)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read resource directory %q: %w", root, err)
	}
	for _, entry := range entries {
		if entry.Type()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("symlink in resource directory %q", filepath.Join(root, entry.Name()))
		}
	}
	return entries, nil
}

// discoverFiles lists direct declarations and ignores unrelated files and subdirectories.
func discoverFiles(root, suffix string) ([]string, error) {
	entries, err := directoryEntries(root)
	if err != nil {
		return nil, fmt.Errorf("discover declarations: %w", err)
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), suffix) {
			continue
		}
		path := filepath.Join(root, entry.Name())
		if !entry.Type().IsRegular() {
			return nil, fmt.Errorf("declaration %q must be a regular file", path)
		}
		files = append(files, path)
	}
	return files, nil
}

// expandHome expands current-user home notation without supporting named-user expansion.
func expandHome(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, strings.TrimPrefix(strings.TrimPrefix(path, "~"), "/")), nil
}

// expandPath expands current-user home notation in a collection resource path.
func expandPath(path string) (string, error) {
	expanded, err := expandHome(path)
	if err != nil {
		return "", fmt.Errorf("expand path %q: %w", path, err)
	}
	absolute, err := filepath.Abs(expanded)
	if err != nil {
		return "", fmt.Errorf("resolve path %q: %w", expanded, err)
	}
	return absolute, nil
}
