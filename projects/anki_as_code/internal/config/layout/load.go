package layout

import (
	"fmt"
	"path/filepath"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"google.golang.org/protobuf/proto"
)

// readNote reads a protobuf-defined note declaration from its source file.
func readNote(path string) (*collectionpb.NoteConfig, error) {
	note := &collectionpb.NoteConfig{}
	if err := readTOML(path, note); err != nil {
		return nil, fmt.Errorf("read note %q: %w", path, err)
	}
	return note, nil
}

// readNotes loads direct note files and applies only the declaring deck's nonempty default.
func readNotes(root, home, collectionRoot string) ([]*collectionpb.NoteConfigFile, error) {
	files, err := discoverFiles(root, noteSuffix)
	if err != nil {
		return nil, fmt.Errorf("discover notes: %w", err)
	}
	notes := make([]*collectionpb.NoteConfigFile, 0, len(files))
	for _, file := range files {
		note, err := readNote(file)
		if err != nil {
			return nil, fmt.Errorf("load note text: %w", err)
		}
		for _, card := range note.Cards {
			if card.Deck == nil && home != "" {
				card.Deck = proto.String(home)
			}
		}
		relative, err := filepath.Rel(collectionRoot, file)
		if err != nil {
			return nil, fmt.Errorf("identify note path: %w", err)
		}
		notes = append(notes, &collectionpb.NoteConfigFile{Path: filepath.ToSlash(relative), Config: note})
	}
	return notes, nil
}

// readDecks reads each deck's configuration and notes, deriving titles from the directory hierarchy.
func readDecks(root, parent, collectionRoot string) ([]*collectionpb.DeckConfigSet, error) {
	entries, err := directoryEntries(root)
	if err != nil {
		return nil, fmt.Errorf("discover decks: %w", err)
	}
	var decks []*collectionpb.DeckConfigSet
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		directory := filepath.Join(root, entry.Name())
		marker := filepath.Join(directory, deckMarker)
		if err := regularPath(marker); err != nil {
			return nil, fmt.Errorf("inspect deck marker: %w", err)
		}
		deck := &collectionpb.DeckConfig{}
		if err := readTOML(marker, deck); err != nil {
			return nil, fmt.Errorf("load deck marker: %w", err)
		}
		title := entry.Name()
		if parent != "" {
			title = parent + "::" + title
		}
		relative, err := filepath.Rel(collectionRoot, filepath.Join(directory, deckMarker))
		if err != nil {
			return nil, fmt.Errorf("identify deck path: %w", err)
		}
		markerPath := filepath.ToSlash(relative)
		notes, err := readNotes(filepath.Join(directory, "notes"), markerPath, collectionRoot)
		if err != nil {
			return nil, fmt.Errorf("load deck %q notes: %w", title, err)
		}
		decks = append(decks, &collectionpb.DeckConfigSet{Title: title, Path: markerPath, Deck: deck, Notes: notes})
		children, err := readDecks(filepath.Join(directory, "decks"), title, collectionRoot)
		if err != nil {
			return nil, fmt.Errorf("load child decks of %q: %w", title, err)
		}
		decks = append(decks, children...)
	}
	return decks, nil
}

// Load parses a root config file or its directory into the grouped protobuf configuration set.
// Syntax, schema types, and filesystem layout are checked; callers validate semantics.
// Source files remain unchanged, and no Anki collection is opened.
func Load(input string) (*collectionpb.CollectionConfigSet, error) {
	configs := &collectionpb.CollectionConfigSet{Config: &collectionpb.CollectionConfig{}}
	config, err := collectionConfig(input)
	if err != nil {
		return nil, fmt.Errorf("resolve collection config: %w", err)
	}
	if err := readTOML(config, configs.Config); err != nil {
		return nil, fmt.Errorf("load root config: %w", err)
	}
	if configs.Config.Version != Version {
		return nil, fmt.Errorf("unsupported text format %d; use format %d", configs.Config.Version, Version)
	}
	root := filepath.Dir(config)
	configs.NoteTypes, err = readNoteTypes(filepath.Join(root, "note_types"), root)
	if err != nil {
		return nil, fmt.Errorf("load note types: %w", err)
	}
	configs.Decks, err = readDecks(filepath.Join(root, "decks"), "", root)
	if err != nil {
		return nil, fmt.Errorf("load decks: %w", err)
	}
	configs.Notes, err = readNotes(filepath.Join(root, "notes"), "", root)
	if err != nil {
		return nil, fmt.Errorf("load collection notes: %w", err)
	}
	return configs, nil
}
