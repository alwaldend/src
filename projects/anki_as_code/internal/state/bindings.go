package state

import (
	"fmt"
	"path"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"google.golang.org/protobuf/proto"
)

// DefaultKey is the versioned Anki config namespace selected when callers omit a key.
const DefaultKey = "anki_as_code/state/v1"

// cloneState preserves bindings and initializes typed maps for import.
func cloneState(value *collectionpb.State) *collectionpb.State {
	result := proto.Clone(value).(*collectionpb.State)
	if result.Resources == nil {
		result.Resources = &collectionpb.Resources{}
	}
	if result.Resources.Decks == nil {
		result.Resources.Decks = map[string]*collectionpb.DeckState{}
	}
	if result.Resources.Notes == nil {
		result.Resources.Notes = map[string]*collectionpb.NoteState{}
	}
	if result.Resources.NoteTypes == nil {
		result.Resources.NoteTypes = map[string]*collectionpb.NoteTypeState{}
	}
	if result.Resources.Presets == nil {
		result.Resources.Presets = map[string]*collectionpb.PresetState{}
	}
	return result
}

// validAddress rejects addresses outside the collection or noncanonical filesystem paths.
func validAddress(address string) error {
	if address == "" || address == "." || strings.HasPrefix(address, "/") || strings.ContainsAny(address, "\\\x00") || path.Clean(address) != address || address == ".." || strings.HasPrefix(address, "../") {
		return fmt.Errorf("invalid collection-relative address %q", address)
	}
	return nil
}

// claimBinding checks path ownership across maps and native identity ownership within a map.
func claimBinding(address string, id int64, paths map[string]bool, identities map[int64]string) error {
	if err := validAddress(address); err != nil {
		return fmt.Errorf("state address: %w", err)
	}
	if id <= 0 {
		return fmt.Errorf("invalid native identity at %q", address)
	}
	if paths[address] {
		return fmt.Errorf("state path %q belongs to multiple entity maps", address)
	}
	if previous := identities[id]; previous != "" {
		return fmt.Errorf("native identity is bound to both %q and %q", previous, address)
	}
	paths[address], identities[id] = true, address
	return nil
}

// validateState checks only the application's metadata schema and identity ownership.
func validateState(value *collectionpb.State) error {
	if value == nil || value.Version != 1 {
		return fmt.Errorf("state must have version 1")
	}
	paths := map[string]bool{}
	decks := map[int64]string{}
	for address, deck := range value.GetResources().GetDecks() {
		if err := claimBinding(address, deck.GetId(), paths, decks); err != nil {
			return fmt.Errorf("deck binding: %w", err)
		}
	}
	notes := map[int64]string{}
	for address, note := range value.GetResources().GetNotes() {
		if err := claimBinding(address, note.GetId(), paths, notes); err != nil {
			return fmt.Errorf("note binding: %w", err)
		}
	}
	types := map[int64]string{}
	for address, definition := range value.GetResources().GetNoteTypes() {
		if err := claimBinding(address, definition.GetId(), paths, types); err != nil {
			return fmt.Errorf("note-type binding: %w", err)
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

// bindType projects native fields and templates beneath their owning note-type state.
func bindType(value *collectionpb.State, address string, id int64, definition *collectionpb.NoteTypeConfig) {
	binding := &collectionpb.NoteTypeState{
		Id:        id,
		Fields:    map[string]*collectionpb.FieldState{},
		Templates: map[string]*collectionpb.TemplateState{},
	}
	for ordinal, field := range definition.Fields {
		binding.Fields[field.Name] = &collectionpb.FieldState{Id: field.Config.GetId(), Ordinal: uint32(ordinal)}
	}
	for ordinal, template := range definition.Templates {
		binding.Templates[template.Name] = &collectionpb.TemplateState{Id: template.Config.GetId(), Ordinal: uint32(ordinal)}
	}
	value.Resources.NoteTypes[address] = binding
}
