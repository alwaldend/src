package state

import (
	"fmt"
	"path"
	"strings"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// DefaultKey is the versioned Anki config namespace selected when callers omit a key.
const DefaultKey = "anki_as_code/state/v1"

// invalidAddressCharacters excludes Windows separators and NUL from collection-relative state paths.
const invalidAddressCharacters = "\\\x00"

// Export generates complete version-1 bindings for projected configs from one native snapshot.
// identities supplies the native ID assigned to each config path by its exporter.
// This operation reads only its inputs; it neither reads nor installs stored state.
func Export(configs *ankiascodepb.CollectionConfigSet, snapshot *resource.Snapshot, identities map[string]int64) (*ankiascodepb.State, error) {
	if configs == nil || snapshot == nil {
		return nil, fmt.Errorf("state export requires configs and a native snapshot")
	}
	value := cloneState(&ankiascodepb.State{Version: 1})
	nativeDecks := map[int64]*decks.Deck{}
	for _, deck := range snapshot.Decks {
		nativeDecks[deck.Id] = deck
	}
	nativeNotes := map[int64]bool{}
	for _, note := range snapshot.Notes {
		nativeNotes[note.Id] = true
	}
	nativeTypes := map[int64]*notetypes.Notetype{}
	for _, definition := range snapshot.NoteTypes {
		nativeTypes[definition.Id] = definition
	}
	for _, deck := range configs.Decks {
		id := identities[deck.Path]
		native := nativeDecks[id]
		if native == nil {
			return nil, fmt.Errorf("exported deck %q has no native identity", deck.Path)
		}
		value.Resources.Decks[deck.Path] = &ankiascodepb.DeckState{Id: id}
		if normal := native.GetNormal(); normal != nil {
			alias := deck.GetDeck().GetPreset()
			if alias == "" || snapshot.Presets[normal.ConfigId] == "" {
				return nil, fmt.Errorf("exported deck %q has no scheduling preset", deck.Path)
			}
			if previous := value.Resources.Presets[alias]; previous != nil && previous.Id != normal.ConfigId {
				return nil, fmt.Errorf("exported preset alias %q identifies multiple presets", alias)
			}
			value.Resources.Presets[alias] = &ankiascodepb.PresetState{Id: normal.ConfigId}
		}
		if err := bindExportedNotes(value, deck.Notes, identities, nativeNotes); err != nil {
			return nil, fmt.Errorf("bind exported deck notes: %w", err)
		}
	}
	if err := bindExportedNotes(value, configs.Notes, identities, nativeNotes); err != nil {
		return nil, fmt.Errorf("bind exported root notes: %w", err)
	}
	for _, definition := range configs.NoteTypes {
		id := identities[definition.Path]
		binding, err := bindType(nativeTypes[id])
		if err != nil {
			return nil, fmt.Errorf("bind exported note type %q: %w", definition.Path, err)
		}
		value.Resources.NoteTypes[definition.Path] = binding
	}
	if err := validateState(value, false); err != nil {
		return nil, fmt.Errorf("validate exported state: %w", err)
	}
	return value, nil
}

// bindExportedNotes resolves declared note paths against identities in the export snapshot.
func bindExportedNotes(value *ankiascodepb.State, notes []*ankiascodepb.NoteConfigFile, identities map[string]int64, nativeNotes map[int64]bool) error {
	for _, note := range notes {
		id := identities[note.Path]
		if !nativeNotes[id] {
			return fmt.Errorf("exported note %q has no native identity", note.Path)
		}
		value.Resources.Notes[note.Path] = &ankiascodepb.NoteState{Id: id}
	}
	return nil
}

// cloneState preserves bindings and initializes typed maps for import.
func cloneState(value *ankiascodepb.State) *ankiascodepb.State {
	result := proto.Clone(value).(*ankiascodepb.State)
	if result.Resources == nil {
		result.Resources = &ankiascodepb.Resources{}
	}
	if result.Resources.Decks == nil {
		result.Resources.Decks = map[string]*ankiascodepb.DeckState{}
	}
	if result.Resources.Notes == nil {
		result.Resources.Notes = map[string]*ankiascodepb.NoteState{}
	}
	if result.Resources.NoteTypes == nil {
		result.Resources.NoteTypes = map[string]*ankiascodepb.NoteTypeState{}
	}
	if result.Resources.Presets == nil {
		result.Resources.Presets = map[string]*ankiascodepb.PresetState{}
	}
	return result
}

// validAddress rejects addresses outside the collection or noncanonical filesystem paths.
func validAddress(address string) error {
	if address == "" || address == "." || strings.HasPrefix(address, "/") || strings.ContainsAny(address, invalidAddressCharacters) || path.Clean(address) != address || address == ".." || strings.HasPrefix(address, "../") {
		return fmt.Errorf("invalid collection-relative address %q", address)
	}
	return nil
}

// claimBinding checks path ownership across maps and native identity ownership within a map.
func claimBinding(address string, id int64, paths map[string]bool, identities map[int64]string, allowUnresolved bool) error {
	if err := validAddress(address); err != nil {
		return fmt.Errorf("state address: %w", err)
	}
	if id < 0 || (id == 0 && !allowUnresolved) {
		return fmt.Errorf("invalid native identity at %q", address)
	}
	if paths[address] {
		return fmt.Errorf("state path %q belongs to multiple entity maps", address)
	}
	if previous := identities[id]; previous != "" {
		return fmt.Errorf("native identity is bound to both %q and %q", previous, address)
	}
	paths[address] = true
	if id != 0 {
		identities[id] = address
	}
	return nil
}

// validateState checks only the application's metadata schema and identity ownership.
func validateState(value *ankiascodepb.State, allowUnresolved bool) error {
	if value == nil || value.Version != 1 {
		return fmt.Errorf("state must have version 1")
	}
	paths := map[string]bool{}
	decks := map[int64]string{}
	for address, deck := range value.GetResources().GetDecks() {
		if deck == nil {
			return fmt.Errorf("nil deck state at %q", address)
		}
		if err := claimBinding(address, deck.GetId(), paths, decks, allowUnresolved); err != nil {
			return fmt.Errorf("deck binding: %w", err)
		}
	}
	notes := map[int64]string{}
	for address, note := range value.GetResources().GetNotes() {
		if note == nil {
			return fmt.Errorf("nil note state at %q", address)
		}
		if err := claimBinding(address, note.GetId(), paths, notes, allowUnresolved); err != nil {
			return fmt.Errorf("note binding: %w", err)
		}
	}
	types := map[int64]string{}
	for address, definition := range value.GetResources().GetNoteTypes() {
		if err := claimBinding(address, definition.GetId(), paths, types, allowUnresolved); err != nil {
			return fmt.Errorf("note-type binding: %w", err)
		}
		if definition == nil {
			return fmt.Errorf("nil note-type state at %q", address)
		}
		for name, field := range definition.Fields {
			if field == nil {
				return fmt.Errorf("nil field state %q at %q", name, address)
			}
		}
		for name, template := range definition.Templates {
			if template == nil {
				return fmt.Errorf("nil template state %q at %q", name, address)
			}
		}
	}
	for alias, preset := range value.GetResources().GetPresets() {
		if alias == "" || preset.GetId() <= 0 {
			return fmt.Errorf("invalid preset binding %q", alias)
		}
	}
	return nil
}

// bindType projects native member identities and ordinals for state import and export.
func bindType(definition *notetypes.Notetype) (*ankiascodepb.NoteTypeState, error) {
	if definition == nil {
		return nil, fmt.Errorf("nil native note type")
	}
	binding := &ankiascodepb.NoteTypeState{
		Id:        definition.Id,
		Fields:    map[string]*ankiascodepb.FieldState{},
		Templates: map[string]*ankiascodepb.TemplateState{},
	}
	for ordinal, field := range definition.Fields {
		if field == nil {
			return nil, fmt.Errorf("nil native field at ordinal %d", ordinal)
		}
		binding.Fields[field.Name] = &ankiascodepb.FieldState{Id: field.Config.GetId(), Ordinal: uint32(ordinal)}
	}
	for ordinal, template := range definition.Templates {
		if template == nil {
			return nil, fmt.Errorf("nil native template at ordinal %d", ordinal)
		}
		binding.Templates[template.Name] = &ankiascodepb.TemplateState{Id: template.Config.GetId(), Ordinal: uint32(ordinal)}
	}
	return binding, nil
}
