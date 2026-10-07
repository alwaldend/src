package textstore

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/ankiformat"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/validation"
)

const (
	// CollectionMarker is the default root config filename used for directory inputs.
	CollectionMarker = "collection.anki.toml"
	// DeckMarker is the exact deck marker filename discovered recursively.
	DeckMarker = "deck.anki.toml"
	// noteSuffix identifies editable note declaration filenames.
	noteSuffix = ".note.anki.toml"
	// NoteTypeSuffix identifies editable note-type appearance declaration filenames.
	NoteTypeSuffix = ".notetype.anki.toml"
)

// Paths are relative to their declaring file, never to a prescribed layout.
func referencePath(config, value string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("missing path in %q", config)
	}
	value, err := expandHome(value)
	if err != nil {
		return "", fmt.Errorf("expand resource path: %w", err)
	}
	if !filepath.IsAbs(value) {
		value = filepath.Join(filepath.Dir(config), value)
	}
	path, err := filepath.Abs(value)
	if err != nil {
		return "", fmt.Errorf("resolve reference %q: %w", value, err)
	}
	return path, nil
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
	path, err := ExpandPath(input)
	if err != nil {
		return "", fmt.Errorf("resolve collection path: %w", err)
	}
	info, err := os.Lstat(path)
	if err != nil {
		return "", fmt.Errorf("inspect collection path %q: %w", path, err)
	}
	if info.IsDir() {
		path = filepath.Join(path, CollectionMarker)
	}
	if err := regularPath(path); err != nil {
		return "", fmt.Errorf("inspect root config: %w", err)
	}
	return path, nil
}

// Unrelated files are ignored. Missing note roots represent empty note sets.
func discoverFiles(root, suffix string, allowMissing bool) ([]string, error) {
	resolved, err := filepath.EvalSymlinks(root)
	if os.IsNotExist(err) && allowMissing {
		// Validate the nearest existing parent too, so missing children cannot hide links.
		parent := filepath.Dir(root)
		for {
			resolved, err = filepath.EvalSymlinks(parent)
			if !os.IsNotExist(err) {
				break
			}
			next := filepath.Dir(parent)
			if next == parent {
				break
			}
			parent = next
		}
		if err != nil {
			return nil, fmt.Errorf("resolve notes parent %q: %w", parent, err)
		}
		if resolved != parent {
			return nil, fmt.Errorf("symlink in notes parent %q", parent)
		}
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("resolve discovery root %q: %w", root, err)
	}
	if resolved != root {
		return nil, fmt.Errorf("symlink in discovery root %q", root)
	}
	info, err := os.Lstat(root)
	if err != nil {
		return nil, fmt.Errorf("inspect discovery root: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("discovery root %q must be a directory", root)
	}
	files := []string{}
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk %q: %w", path, walkErr)
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink in discovery tree %q", path)
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), suffix) {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("declaration %q must be a regular file", path)
		}
		files = append(files, path)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("discover declarations: %w", err)
	}
	return files, nil
}

// discoverDecks reads deck markers and resolves their configured note locations.
func discoverDecks(config string, manifest *collectionpb.CollectionConfig) ([]model.Deck, error) {
	root, err := referencePath(config, manifest.DecksPath)
	if err != nil {
		return nil, fmt.Errorf("resolve decks root: %w", err)
	}
	markers, err := discoverFiles(root, DeckMarker, false)
	if err != nil {
		return nil, fmt.Errorf("find deck markers: %w", err)
	}
	decks := make([]model.Deck, 0, len(markers))
	for _, path := range markers {
		// The deck marker has an exact filename, unlike named note suffixes.
		if filepath.Base(path) != DeckMarker {
			continue
		}
		deck := model.Deck{Deck: &collectionpb.Deck{}}
		if err := ReadTOML(path, deck.Deck); err != nil {
			return nil, fmt.Errorf("read deck marker: %w", err)
		}
		if err := validation.ValidDeck(deck.Title); err != nil {
			return nil, fmt.Errorf("read deck title: %w", err)
		}
		if _, err := referencePath(path, deck.NotesPath); err != nil {
			return nil, fmt.Errorf("read deck notes path: %w", err)
		}
		if err := ankiformat.EncodeDeck(&deck); err != nil {
			return nil, fmt.Errorf("read deck settings: %w", err)
		}
		deck.Path = path
		decks = append(decks, deck)
	}
	return decks, nil
}

// Identity reservation uses referenced note roots, including locations outside
// the collection config's directory, without requiring copied IDs to be unique.
func referencedNotes(root string) ([]string, error) {
	config, err := collectionConfig(root)
	if err != nil {
		return nil, fmt.Errorf("find identity root: %w", err)
	}
	manifest := &collectionpb.CollectionConfig{}
	if err := ReadTOML(config, manifest); err != nil {
		return nil, fmt.Errorf("read identity root: %w", err)
	}
	decks, err := discoverDecks(config, manifest)
	if err != nil {
		return nil, fmt.Errorf("read identity decks: %w", err)
	}
	paths := []string{}
	if manifest.GetNotesPath() != "" {
		path, err := referencePath(config, manifest.GetNotesPath())
		if err != nil {
			return nil, fmt.Errorf("resolve identity notes: %w", err)
		}
		paths = append(paths, path)
	}
	for _, deck := range decks {
		path, err := referencePath(deck.Path, deck.NotesPath)
		if err != nil {
			return nil, fmt.Errorf("resolve deck identity notes: %w", err)
		}
		paths = append(paths, path)
	}
	files := []string{}
	seen := map[string]bool{}
	for _, path := range paths {
		discovered, err := discoverFiles(path, noteSuffix, true)
		if err != nil {
			return nil, fmt.Errorf("discover reserved identities: %w", err)
		}
		for _, file := range discovered {
			if !seen[file] {
				seen[file] = true
				files = append(files, file)
			}
		}
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

// ExpandPath expands current-user home notation in a collection resource path.
func ExpandPath(path string) (string, error) {
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
