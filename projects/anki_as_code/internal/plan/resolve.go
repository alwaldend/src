package plan

import (
	"fmt"
	"math"
	"slices"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"

	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// resolveDesired clones declarations and resolves each resource group in dependency order.
func resolveDesired(configs *collectionpb.CollectionConfigSet, current snapshot, state *collectionpb.State) (desiredState, error) {
	configs = proto.Clone(configs).(*collectionpb.CollectionConfigSet)
	desired := desiredState{
		notes:          map[string]runtimeNote{},
		noteTypes:      map[string]*runtimeNoteType{},
		bindings:       state,
		presetBindings: map[string]int64{},
	}
	for alias, preset := range state.Resources.Presets {
		desired.presetBindings[alias] = preset.Id
	}
	paths := map[string]bool{}
	if err := resolveDecks(&desired, configs.Decks, current, paths); err != nil {
		return desired, fmt.Errorf("resolve decks: %w", err)
	}
	if err := resolveNoteTypes(&desired, configs.NoteTypes, current, paths); err != nil {
		return desired, fmt.Errorf("resolve note types: %w", err)
	}
	if err := resolveNotes(&desired, configs.Notes, "", current, paths); err != nil {
		return desired, fmt.Errorf("resolve root notes: %w", err)
	}
	for _, group := range configs.Decks {
		if err := resolveNotes(&desired, group.Notes, group.Path, current, paths); err != nil {
			return desired, fmt.Errorf("resolve deck %q notes: %w", group.Path, err)
		}
	}
	return desired, nil
}

// resolveDecks binds deck paths, resolves presets and orders declared parents before children.
func resolveDecks(desired *desiredState, groups []*collectionpb.DeckConfigSet, current snapshot, paths map[string]bool) error {
	names := map[string]bool{}
	for _, group := range groups {
		if desired.bindings.Resources.Notes[group.Path] != nil || desired.bindings.Resources.NoteTypes[group.Path] != nil {
			return fmt.Errorf("deck path %q is bound to another entity", group.Path)
		}
		id := desired.bindings.Resources.Decks[group.Path].GetId()
		if !slices.ContainsFunc(current.decks, func(old runtimeDeck) bool { return old.id == id }) {
			id = 0
			if slices.ContainsFunc(current.decks, func(old runtimeDeck) bool { return folded(old.title) == folded(group.Title) }) {
				return fmt.Errorf("deck %q exists without this path binding; import it first", group.Title)
			}
		}
		for _, old := range current.decks {
			if old.id != id && old.path == "" && folded(old.title) == folded(group.Title) {
				return fmt.Errorf("deck title %q belongs to an untracked object; import it first", group.Title)
			}
		}
		if names[folded(group.Title)] || paths[group.Path] {
			return fmt.Errorf("duplicate deck path or title %q", group.Path)
		}
		names[folded(group.Title)], paths[group.Path] = true, true
		if (group.Deck.Normal == nil) == (group.Deck.Filtered == nil) {
			return fmt.Errorf("deck %q requires exactly one normal or filtered configuration", group.Path)
		}
		if group.Deck.Filtered != nil && group.Deck.Preset != "" {
			return fmt.Errorf("filtered decks do not reference scheduling presets")
		}
		if group.Deck.Normal != nil {
			if group.Deck.Normal.ConfigId != 0 {
				return fmt.Errorf("normal.config_id is native-only; use preset")
			}
			presetID, err := resolvePreset(group.Deck.Preset, current.presets, desired.presetBindings)
			if err != nil {
				return fmt.Errorf("deck %q: %w", group.Path, err)
			}
			group.Deck.Normal.ConfigId = presetID
		}
		deck := runtimeDeck{id: id, path: group.Path, title: group.Title, config: group.Deck}
		for _, old := range current.decks {
			if old.id == id && (old.config.Filtered != nil) != (deck.config.Filtered != nil) && current.occupiedDecks[id] {
				return fmt.Errorf("empty occupied deck in Anki before changing its kind")
			}
		}
		desired.decks = append(desired.decks, deck)
	}
	for _, deck := range desired.decks {
		parts := strings.Split(deck.title, "::")
		for i := 1; i < len(parts); i++ {
			if !names[folded(strings.Join(parts[:i], "::"))] {
				return fmt.Errorf("deck %q lacks a declared parent", deck.title)
			}
		}
	}
	slices.SortFunc(desired.decks, func(a, b runtimeDeck) int { return strings.Compare(a.title, b.title) })
	return nil
}

// resolvePreset maps a readable alias to an existing unmanaged scheduling preset.
func resolvePreset(alias string, presets map[int64]string, bindings map[string]int64) (int64, error) {
	id := int64(1)
	if alias != "" {
		id = bindings[alias]
		if id == 0 {
			for candidate, name := range presets {
				if name == alias {
					if id != 0 {
						return 0, fmt.Errorf("preset %q is ambiguous; import its binding", alias)
					}
					id = candidate
				}
			}
			if id == 0 {
				return 0, fmt.Errorf("unknown existing preset %q", alias)
			}
			bindings[alias] = id
		}
	}
	if _, exists := presets[id]; !exists {
		return 0, fmt.Errorf("preset %q is missing in Anki", alias)
	}
	return id, nil
}

// resolveNoteTypes binds type paths and translates field, template and deck references.
func resolveNoteTypes(desired *desiredState, files []*collectionpb.NoteTypeConfigFile, current snapshot, paths map[string]bool) error {
	typeNames := map[string]bool{}
	for _, file := range files {
		if desired.bindings.Resources.Decks[file.Path] != nil || desired.bindings.Resources.Notes[file.Path] != nil {
			return fmt.Errorf("note-type path %q is bound to another entity", file.Path)
		}
		id := desired.bindings.Resources.NoteTypes[file.Path].GetId()
		if current.noteTypes[id] == nil {
			id = 0
			for _, old := range current.noteTypes {
				if folded(old.definition.Name) == folded(file.Config.Name) {
					return fmt.Errorf("note type %q exists without this path binding; import it first", file.Config.Name)
				}
			}
		}
		for _, old := range current.noteTypes {
			if old.id != id && old.path == "" && folded(old.definition.Name) == folded(file.Config.Name) {
				return fmt.Errorf("note-type name %q belongs to an untracked object; import it first", file.Config.Name)
			}
		}
		if paths[file.Path] || typeNames[folded(file.Config.Name)] {
			return fmt.Errorf("duplicate note-type path or name")
		}
		paths[file.Path], typeNames[folded(file.Config.Name)] = true, true
		if err := prepareNoteType(file.Config); err != nil {
			return fmt.Errorf("note type %q: %w", file.Path, err)
		}
		if err := resolveTypeReferences(file.Config, desired.decks); err != nil {
			return fmt.Errorf("resolve type references: %w", err)
		}
		entry := &runtimeNoteType{id: id, path: file.Path, definition: file.Config}
		var original *collectionpb.NoteTypeState
		if old := current.noteTypes[id]; old != nil {
			original = current.bindings.Resources.NoteTypes[old.path]
		}
		prepared, err := prepareNativeType(entry, current.noteTypes[id], desired.bindings.Resources.NoteTypes[file.Path], original)
		if err != nil {
			return fmt.Errorf("prepare note type %q: %w", file.Path, err)
		}
		entry.native = prepared
		desired.noteTypes[file.Path] = entry
	}
	return nil
}

// resolveNotes binds desired paths to existing notes and orders their named fields.
func resolveNotes(desired *desiredState, files []*collectionpb.NoteConfigFile, home string, current snapshot, paths map[string]bool) error {
	for _, file := range files {
		if paths[file.Path] {
			return fmt.Errorf("duplicate resource path %q", file.Path)
		}
		paths[file.Path] = true
		if desired.bindings.Resources.Decks[file.Path] != nil || desired.bindings.Resources.NoteTypes[file.Path] != nil {
			return fmt.Errorf("note path %q is bound to another entity", file.Path)
		}
		id := desired.bindings.Resources.Notes[file.Path].GetId()
		note := runtimeNote{path: file.Path, config: file.Config, homeDeck: home}
		previous := current.notes[id]
		note.id = previous.id
		if err := validateNote(&note, *desired, previous); err != nil {
			return fmt.Errorf("note %q: %w", file.Path, err)
		}
		desired.notes[file.Path] = note
	}
	return nil
}

// validateNote resolves its path reference and checks fields and generated-card assignments.
func validateNote(note *runtimeNote, desired desiredState, previous runtimeNote) error {
	definition := desired.noteTypes[note.config.NoteType]
	if definition == nil {
		return fmt.Errorf("unknown declared note-type path %q", note.config.NoteType)
	}
	note.typeID = definition.id
	settings := definition.definition
	if note.id != 0 && previous.typeID != note.typeID {
		return fmt.Errorf("retained note cannot change note type")
	}
	if len(note.config.Fields) != len(settings.Fields) {
		return fmt.Errorf("field count differs from declared note type")
	}
	for _, field := range settings.Fields {
		body, exists := note.config.Fields[field.Name]
		if !exists {
			return fmt.Errorf("missing field %q", field.Name)
		}
		note.bodies = append(note.bodies, body)
	}
	ordinals := map[int64]bool{}
	for _, card := range note.config.Cards {
		if card.Ordinal < 0 || card.Ordinal > math.MaxInt32 || ordinals[card.Ordinal] {
			return fmt.Errorf("invalid or duplicate card ordinal")
		}
		ordinals[card.Ordinal] = true
		if settings.Config.Kind != notetypes.Notetype_Config_KIND_CLOZE && card.Ordinal >= int64(len(settings.Templates)) {
			return fmt.Errorf("no template at card ordinal %d", card.Ordinal)
		}
		if !slices.ContainsFunc(desired.decks, func(deck runtimeDeck) bool { return deck.path == card.GetDeck() && deck.config.Normal != nil }) {
			return fmt.Errorf("card requires declared normal-deck path %q", card.GetDeck())
		}
	}
	if note.homeDeck != "" && !slices.ContainsFunc(desired.decks, func(deck runtimeDeck) bool { return deck.path == note.homeDeck && deck.config.Normal != nil }) {
		return fmt.Errorf("default deck must be normal")
	}
	return nil
}

// resolveTypeReferences converts readable references into private native API settings.
func resolveTypeReferences(definition *collectionpb.NoteTypeConfig, decks []runtimeDeck) error {
	if definition.Config.OriginalId != nil || definition.Config.TargetDeckIdUnused != 0 || definition.Config.SortFieldIdx != 0 || len(definition.Config.Reqs) != 0 || len(definition.Config.Other) != 0 {
		return fmt.Errorf("native identity, provenance, and derived type fields are not editable")
	}
	if definition.SortField != "" {
		found := false
		for ordinal, field := range definition.Fields {
			if field.Name == definition.SortField {
				definition.Config.SortFieldIdx = uint32(ordinal)
				found = true
			}
		}
		if !found {
			return fmt.Errorf("unknown sort field %q", definition.SortField)
		}
	}
	for _, field := range definition.Fields {
		if field.Config.Id != nil || len(field.Config.Other) != 0 {
			return fmt.Errorf("field config.id is native-only")
		}
	}
	for _, template := range definition.Templates {
		if template.Config.Id != nil || template.Config.TargetDeckId != 0 || len(template.Config.Other) != 0 {
			return fmt.Errorf("template native IDs are not editable; use target_deck")
		}
		if template.TargetDeck != "" {
			found := false
			for _, deck := range decks {
				if deck.path == template.TargetDeck && deck.config.Normal != nil {
					template.Config.TargetDeckId = deck.id
					found = true
				}
			}
			if !found {
				return fmt.Errorf("unknown normal-deck path %q", template.TargetDeck)
			}
		}

	}
	return nil
}

// resolve binds desired declarations to a supplied native snapshot and state.
func resolve(configs *collectionpb.CollectionConfigSet, native *resource.Snapshot, stored *collectionpb.State) (desiredState, snapshot, error) {
	if configs == nil || configs.GetConfig().GetVersion() != 1 {
		return desiredState{}, snapshot{}, fmt.Errorf("configs must have version 1")
	}
	if native == nil {
		return desiredState{}, snapshot{}, fmt.Errorf("native snapshot is required")
	}
	if stored == nil || stored.Version != 1 {
		return desiredState{}, snapshot{}, fmt.Errorf("state must have version 1")
	}
	if err := validateMoveSources(configs); err != nil {
		return desiredState{}, snapshot{}, fmt.Errorf("resolve move sources: %w", err)
	}
	moved, err := moveState(stored, configs.GetConfig().GetMoved())
	if err != nil {
		return desiredState{}, snapshot{}, fmt.Errorf("resolve state moves: %w", err)
	}
	before, err := readSnapshot(native, stored)
	if err != nil {
		return desiredState{}, snapshot{}, fmt.Errorf("read supplied snapshot: %w", err)
	}
	desired, err := resolveDesired(configs, before, moved)
	if err != nil {
		return desiredState{}, snapshot{}, fmt.Errorf("resolve desired resources: %w", err)
	}
	if err := validateDesired(desired, before); err != nil {
		return desiredState{}, snapshot{}, fmt.Errorf("validate reconciliation ownership: %w", err)
	}
	return desired, before, nil
}
