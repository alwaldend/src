package layout

import (
	"fmt"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"google.golang.org/protobuf/proto"
)

// readNote reads a protobuf-defined note declaration from its source file.
func readNote(storage filesystem.FS, path string) (*ankiascodepb.NoteConfig, error) {
	note := &ankiascodepb.NoteConfig{}
	if err := readTOML(storage, path, note); err != nil {
		return nil, fmt.Errorf("read note %q: %w", path, err)
	}
	return note, nil
}

// readNotes loads direct note files and applies only the declaring deck's nonempty default.
func readNotes(storage filesystem.FS, root, home, collectionRoot string) ([]*ankiascodepb.NoteConfigFile, error) {
	files, err := discoverFiles(storage, root, noteSuffix)
	if err != nil {
		return nil, fmt.Errorf("discover notes: %w", err)
	}
	notes := make([]*ankiascodepb.NoteConfigFile, 0, len(files))
	for _, file := range files {
		note, err := readNote(storage, file)
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
		notes = append(notes, &ankiascodepb.NoteConfigFile{Path: filepath.ToSlash(relative), Config: note})
	}
	return notes, nil
}

// readDecks reads each deck's configuration and notes, deriving titles from the directory hierarchy.
func readDecks(storage filesystem.FS, root, parent, collectionRoot string) ([]*ankiascodepb.DeckConfigSet, error) {
	entries, err := directoryEntries(storage, root)
	if err != nil {
		return nil, fmt.Errorf("discover decks: %w", err)
	}
	var decks []*ankiascodepb.DeckConfigSet
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		directory := filepath.Join(root, entry.Name())
		marker := filepath.Join(directory, deckMarker)
		if err := regularPath(storage, marker); err != nil {
			return nil, fmt.Errorf("inspect deck marker: %w", err)
		}
		deck := &ankiascodepb.DeckConfig{}
		if err := readTOML(storage, marker, deck); err != nil {
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
		notes, err := readNotes(storage, filepath.Join(directory, "notes"), markerPath, collectionRoot)
		if err != nil {
			return nil, fmt.Errorf("load deck %q notes: %w", title, err)
		}
		decks = append(decks, &ankiascodepb.DeckConfigSet{Title: title, Path: markerPath, Deck: deck, Notes: notes})
		children, err := readDecks(storage, filepath.Join(directory, "decks"), title, collectionRoot)
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
func Load(storage filesystem.FS, input string) (*ankiascodepb.CollectionConfigSet, error) {
	configs := &ankiascodepb.CollectionConfigSet{Config: &ankiascodepb.CollectionConfig{}}
	config, err := collectionConfig(storage, input)
	if err != nil {
		return nil, fmt.Errorf("resolve collection config: %w", err)
	}
	if err := readTOML(storage, config, configs.Config); err != nil {
		return nil, fmt.Errorf("load root config: %w", err)
	}
	if configs.Config.Version != Version {
		return nil, fmt.Errorf("unsupported text format %d; use format %d", configs.Config.Version, Version)
	}
	root := filepath.Dir(config)
	configs.NoteTypes, err = readNoteTypes(storage, filepath.Join(root, "note_types"), root)
	if err != nil {
		return nil, fmt.Errorf("load note types: %w", err)
	}
	configs.Decks, err = readDecks(storage, filepath.Join(root, "decks"), "", root)
	if err != nil {
		return nil, fmt.Errorf("load decks: %w", err)
	}
	configs.Notes, err = readNotes(storage, filepath.Join(root, "notes"), "", root)
	if err != nil {
		return nil, fmt.Errorf("load collection notes: %w", err)
	}
	return configs, nil
}
