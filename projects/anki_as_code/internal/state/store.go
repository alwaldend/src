package state

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	ankiconfig "github.com/ankitects/anki/proto/anki/config"
	"github.com/ankitects/anki/proto/anki/deck_config"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Store manages the application's state key through a caller-owned Anki connection.
type Store struct {
	// key selects the only collection configuration entry this store manages.
	key string
	// configs reads collection settings and writes only the owned state key.
	configs ankiconfig.ConfigServiceClient
	// decks checks imported deck identities without editing them.
	decks decks.DecksServiceClient
	// notes checks imported note identities without reading unrelated notes/cards.
	notes notes.NotesServiceClient
	// noteTypes reads native member identities for imported note types.
	noteTypes notetypes.NotetypesServiceClient
	// presets checks existing scheduling preset references without changing values.
	presets deck_config.DeckConfigServiceClient
}

// New borrows an open connection and selects a namespace; an empty key uses DefaultKey.
// It starts no worker and owns no files.
func New(connection grpc.ClientConnInterface, key string) *Store {
	if key == "" {
		key = DefaultKey
	}
	return &Store{
		key:       key,
		configs:   ankiconfig.NewConfigServiceClient(connection),
		decks:     decks.NewDecksServiceClient(connection),
		notes:     notes.NewNotesServiceClient(connection),
		noteTypes: notetypes.NewNotetypesServiceClient(connection),
		presets:   deck_config.NewDeckConfigServiceClient(connection),
	}
}

// read returns the selected raw metadata value, or nil when its key is absent.
func (store *Store) read(ctx context.Context) (json.RawMessage, error) {
	values, err := store.configs.GetAllConfig(ctx, &generic.Empty{})
	if err != nil {
		return nil, fmt.Errorf("read collection configuration: %w", err)
	}
	var entries map[string]json.RawMessage
	if err := json.Unmarshal(values.Json, &entries); err != nil {
		return nil, fmt.Errorf("decode configuration entries: %w", err)
	}
	return entries[store.key], nil
}

// Dump returns the stored protobuf object without writes; absent state is empty version 1.
func (store *Store) Dump(ctx context.Context) (*ankiascodepb.State, error) {
	data, err := store.read(ctx)
	if err != nil {
		return nil, fmt.Errorf("read collection state: %w", err)
	}
	value := &ankiascodepb.State{Version: 1}
	if len(data) != 0 {
		if err := protojson.Unmarshal(data, value); err != nil {
			return nil, fmt.Errorf("decode collection state: %w", err)
		}
	}
	if err := validateState(value, false); err != nil {
		return nil, fmt.Errorf("read collection state: %w", err)
	}
	return value, nil
}

// Replace replaces the complete stored object, including removal of omitted bindings.
// It validates state metadata, leaves native objects untouched, and reports actual changes.
// Invalid saved state can be replaced; native read errors still fail the operation.
func (store *Store) Replace(ctx context.Context, value *ankiascodepb.State) (bool, error) {
	if err := validateState(value, false); err != nil {
		return false, fmt.Errorf("replacement state: %w", err)
	}
	data, err := store.read(ctx)
	if err != nil {
		return false, fmt.Errorf("read state before replacement: %w", err)
	}
	before := &ankiascodepb.State{Version: 1}
	if len(data) != 0 {
		if err := protojson.Unmarshal(data, before); err != nil {
			// An undecodable discarded value must force a write, not block recovery.
			before = nil
		}
	}
	changed, err := store.write(ctx, before, value)
	if err != nil {
		return false, fmt.Errorf("replace collection state: %w", err)
	}
	return changed, nil
}

// Update merges external state without removing omitted map entries or changing native objects.
// Supplied map keys replace complete typed entries, including note-type members.
// It validates the merged ownership before one write and leaves its input unchanged.
func (store *Store) Update(ctx context.Context, value *ankiascodepb.State) (bool, error) {
	if err := validateState(value, false); err != nil {
		return false, fmt.Errorf("update state: %w", err)
	}
	if resources := value.Resources; resources == nil || len(resources.Decks)+len(resources.Notes)+len(resources.NoteTypes)+len(resources.Presets) == 0 {
		return false, nil
	}
	before, err := store.Dump(ctx)
	if err != nil {
		return false, fmt.Errorf("read state before update: %w", err)
	}
	merged := proto.Clone(before).(*ankiascodepb.State)
	proto.Merge(merged, value)
	if err := validateState(merged, false); err != nil {
		return false, fmt.Errorf("merged state: %w", err)
	}
	changed, err := store.write(ctx, before, merged)
	if err != nil {
		return false, fmt.Errorf("update collection state: %w", err)
	}
	return changed, nil
}

// Import binds declared paths to existing native IDs and saves once after the whole batch succeeds.
// Identifiers are interpreted after resolving declarations. Native note-type members
// are bound automatically; unrelated bindings and native objects remain unchanged.
func (store *Store) Import(ctx context.Context, configs *ankiascodepb.CollectionConfigSet, bindings map[string]string) (bool, error) {
	if configs == nil {
		return false, fmt.Errorf("import requires configurations")
	}
	if len(bindings) == 0 {
		return false, nil
	}
	before, err := store.Dump(ctx)
	if err != nil {
		return false, fmt.Errorf("read imported state: %w", err)
	}
	declaredDecks := map[string]*ankiascodepb.DeckConfig{}
	declaredNotes, declaredTypes := map[string]bool{}, map[string]bool{}
	for _, value := range configs.NoteTypes {
		declaredTypes[value.Path] = true
	}
	for _, value := range configs.Notes {
		declaredNotes[value.Path] = true
	}
	for _, value := range configs.Decks {
		declaredDecks[value.Path] = value.Deck
		for _, note := range value.Notes {
			declaredNotes[note.Path] = true
		}
	}
	value := cloneState(before)
	for _, address := range slices.Sorted(maps.Keys(bindings)) {
		if err := validAddress(address); err != nil {
			return false, fmt.Errorf("import binding %q: %w", address, err)
		}
		isPreset := strings.HasPrefix(address, "presets/")
		if declaredDecks[address] == nil && !declaredNotes[address] && !declaredTypes[address] && !isPreset {
			return false, fmt.Errorf("configuration path %q is not declared", address)
		}
		id, err := strconv.ParseInt(bindings[address], 10, 64)
		if err != nil {
			return false, fmt.Errorf("import at %q requires an int64 native ID: %w", address, err)
		}
		if id <= 0 {
			return false, fmt.Errorf("import at %q requires a positive native ID", address)
		}
		switch {
		case declaredDecks[address] != nil:
			err = store.importDeck(ctx, value, address, id, declaredDecks[address].GetPreset())
		case declaredNotes[address]:
			err = store.importNote(ctx, value, address, id)
		case declaredTypes[address]:
			err = store.importNoteType(ctx, value, address, id)
		case isPreset:
			err = store.importPreset(ctx, value, strings.TrimPrefix(address, "presets/"), id)
		}
		if err != nil {
			return false, fmt.Errorf("import binding %q: %w", address, err)
		}
	}
	if err := validateState(value, false); err != nil {
		return false, fmt.Errorf("validate imported state: %w", err)
	}
	changed, err := store.write(ctx, before, value)
	if err != nil {
		return false, fmt.Errorf("save imported bindings: %w", err)
	}
	return changed, nil
}

// importDeck verifies an existing deck and binds its declared preset alias when names match.
func (store *Store) importDeck(ctx context.Context, value *ankiascodepb.State, address string, id int64, alias string) error {
	if previous := value.Resources.Decks[address].GetId(); previous != 0 && previous != id {
		return fmt.Errorf("deck path is already bound")
	}
	deck, err := store.decks.GetDeck(ctx, &decks.DeckId{Did: id})
	if err != nil {
		return fmt.Errorf("read imported deck: %w", err)
	}
	if deck.Id != id {
		return fmt.Errorf("deck identity %d does not exist", id)
	}
	if normal := deck.GetNormal(); normal != nil && alias != "" && value.Resources.Presets[alias] == nil {
		preset, err := store.presets.GetDeckConfig(ctx, &deck_config.DeckConfigId{Dcid: normal.ConfigId})
		if err != nil {
			return fmt.Errorf("read imported deck preset: %w", err)
		}
		if preset.Name == alias {
			value.Resources.Presets[alias] = &ankiascodepb.PresetState{Id: normal.ConfigId}
		}
	}
	value.Resources.Decks[address] = &ankiascodepb.DeckState{Id: id}
	return nil
}

// importNote verifies an existing note without reading unrelated notes or cards.
func (store *Store) importNote(ctx context.Context, value *ankiascodepb.State, address string, id int64) error {
	if previous := value.Resources.Notes[address].GetId(); previous != 0 && previous != id {
		return fmt.Errorf("note path is already bound")
	}
	note, err := store.notes.GetNote(ctx, &notes.NoteId{Nid: id})
	if err != nil {
		return fmt.Errorf("read imported note: %w", err)
	}
	if note.Id != id {
		return fmt.Errorf("note identity %d does not exist", id)
	}
	value.Resources.Notes[address] = &ankiascodepb.NoteState{Id: id}
	return nil
}

// importNoteType verifies a native definition and records its nested field and template identities.
func (store *Store) importNoteType(ctx context.Context, value *ankiascodepb.State, address string, id int64) error {
	if previous := value.Resources.NoteTypes[address].GetId(); previous != 0 && previous != id {
		return fmt.Errorf("note-type path is already bound")
	}
	definition, err := store.noteTypes.GetNotetype(ctx, &notetypes.NotetypeId{Ntid: id})
	if err != nil {
		return fmt.Errorf("read imported note type: %w", err)
	}
	if definition.Id != id {
		return fmt.Errorf("note type identity %d does not exist", id)
	}
	binding, err := bindType(definition)
	if err != nil {
		return fmt.Errorf("bind imported note type: %w", err)
	}
	value.Resources.NoteTypes[address] = binding
	return nil
}

// importPreset verifies an external scheduling preset and records its alias without editing its settings.
func (store *Store) importPreset(ctx context.Context, value *ankiascodepb.State, alias string, id int64) error {
	if previous := value.Resources.Presets[alias].GetId(); previous != 0 && previous != id {
		return fmt.Errorf("preset alias %q is already bound", alias)
	}
	preset, err := store.presets.GetDeckConfig(ctx, &deck_config.DeckConfigId{Dcid: id})
	if err != nil {
		return fmt.Errorf("read imported preset: %w", err)
	}
	if preset.Id != id {
		return fmt.Errorf("preset identity %d does not exist", id)
	}
	value.Resources.Presets[alias] = &ankiascodepb.PresetState{Id: id}
	return nil
}

// write persists a changed object once, retaining a byte-preserving no-op signal for callers.
func (store *Store) write(ctx context.Context, before, value *ankiascodepb.State) (bool, error) {
	if proto.Equal(before, value) {
		return false, nil
	}
	data, err := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(value)
	if err != nil {
		return false, fmt.Errorf("encode collection state: %w", err)
	}
	if _, err := store.configs.SetConfigJson(ctx, &ankiconfig.SetConfigJsonRequest{Key: store.key, ValueJson: data}); err != nil {
		return false, fmt.Errorf("write collection state: %w", err)
	}
	return true, nil
}
