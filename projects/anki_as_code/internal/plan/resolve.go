package plan

import (
	"fmt"
	"math"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"github.com/ankitects/anki/proto/anki/decks"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// resolveDesired prepares owned protobuf payloads and indexes declarations in dependency order.
func resolveDesired(configs *collectionpb.CollectionConfigSet, current snapshot, bindings *collectionpb.Resources) (desiredState, error) {
	configs = proto.CloneOf(configs)
	desired := desiredState{
		bindings:    bindings,
		decksByPath: map[string]*collectionpb.PlannedDeck{},
		decksByID:   map[int64]*collectionpb.PlannedDeck{},
		noteTypes:   map[string]*collectionpb.PlannedNoteType{},
		typesByID:   map[int64]*collectionpb.PlannedNoteType{},
		notes:       map[string]*collectionpb.PlannedNote{},
		notesByID:   map[int64]*collectionpb.PlannedNote{},
		cards:       map[string]*collectionpb.PlannedCardAssignments{},
	}
	if bindings.Presets == nil {
		bindings.Presets = map[string]*collectionpb.PresetState{}
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

// resolveDecks binds paths, resolves presets and prepares native decks with shared indexes.
func resolveDecks(desired *desiredState, groups []*collectionpb.DeckConfigSet, current snapshot, paths map[string]bool) error {
	names := map[string]bool{}
	for _, group := range groups {
		if desired.bindings.Notes[group.Path] != nil || desired.bindings.NoteTypes[group.Path] != nil {
			return fmt.Errorf("deck path %q is bound to another entity", group.Path)
		}
		id := desired.bindings.Decks[group.Path].GetId()
		owners := current.deckNames[folded(group.Title)]
		if current.decks[id] == nil {
			id = 0
		}
		for _, owner := range owners {
			if owner != id && current.deckPaths[owner] == "" {
				return fmt.Errorf("deck title %q belongs to an untracked object; import it first", group.Title)
			}
		}
		if names[folded(group.Title)] || paths[group.Path] {
			return fmt.Errorf("duplicate deck path or title %q", group.Path)
		}
		names[folded(group.Title)], paths[group.Path] = true, true
		settings := group.Deck
		if settings == nil || (settings.Normal == nil) == (settings.Filtered == nil) {
			return fmt.Errorf("deck %q requires exactly one normal or filtered configuration", group.Path)
		}
		if settings.Filtered != nil && settings.Preset != "" {
			return fmt.Errorf("filtered decks do not reference scheduling presets")
		}
		value := &decks.Deck{Id: id, Name: group.Title}
		if settings.Normal != nil {
			if settings.Normal.ConfigId != 0 {
				return fmt.Errorf("normal.config_id is native-only; use preset")
			}
			presetID, err := resolvePreset(settings.Preset, current.presets, desired.bindings.Presets)
			if err != nil {
				return fmt.Errorf("deck %q: %w", group.Path, err)
			}
			settings.Normal.ConfigId = presetID
			value.Kind = &decks.Deck_Normal_{Normal: settings.Normal}
		} else {
			value.Kind = &decks.Deck_Filtered_{Filtered: settings.Filtered}
		}
		if old := current.decks[id]; old != nil && (old.GetFiltered() != nil) != (value.GetFiltered() != nil) && current.occupiedDecks[id] {
			return fmt.Errorf("empty occupied deck in Anki before changing its kind")
		}
		operation := &collectionpb.PlannedDeck{Path: group.Path, Value: value}
		desired.decksByPath[group.Path] = operation
		if id != 0 {
			desired.decksByID[id] = operation
		}
	}
	for _, deck := range desired.decksByPath {
		parts := strings.Split(deck.Value.Name, "::")
		for i := 1; i < len(parts); i++ {
			if !names[folded(strings.Join(parts[:i], "::"))] {
				return fmt.Errorf("deck %q lacks a declared parent", deck.Value.Name)
			}
		}
	}
	return nil
}

// resolvePreset maps a readable alias to an existing unmanaged scheduling preset.
func resolvePreset(alias string, presets map[int64]string, bindings map[string]*collectionpb.PresetState) (int64, error) {
	id := int64(1)
	if alias != "" {
		id = bindings[alias].GetId()
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
			bindings[alias] = &collectionpb.PresetState{Id: id}
		}
	}
	if _, exists := presets[id]; !exists {
		return 0, fmt.Errorf("preset %q is missing in Anki", alias)
	}
	return id, nil
}

// resolveNoteTypes binds paths and prepares member mappings and unresolved deck references once.
func resolveNoteTypes(desired *desiredState, files []*collectionpb.NoteTypeConfigFile, current snapshot, paths map[string]bool) error {
	typeNames := map[string]bool{}
	for _, file := range files {
		if desired.bindings.Decks[file.Path] != nil || desired.bindings.Notes[file.Path] != nil {
			return fmt.Errorf("note-type path %q is bound to another entity", file.Path)
		}
		if err := prepareNoteType(file.Config); err != nil {
			return fmt.Errorf("note type %q: %w", file.Path, err)
		}
		id := desired.bindings.NoteTypes[file.Path].GetId()
		previous := current.noteTypes[id]
		for _, owner := range current.typeNames[folded(file.Config.Name)] {
			if owner != previous.GetId() && current.typePaths[owner] == "" {
				return fmt.Errorf("note-type name %q belongs to an untracked object; import it first", file.Config.Name)
			}
		}
		if paths[file.Path] || typeNames[folded(file.Config.Name)] {
			return fmt.Errorf("duplicate note-type path or name")
		}
		paths[file.Path], typeNames[folded(file.Config.Name)] = true, true
		if err := resolveTypeReferences(file.Config, desired.decksByPath); err != nil {
			return fmt.Errorf("resolve type references: %w", err)
		}
		original := current.bindings.GetResources().GetNoteTypes()[current.typePaths[id]]
		value, err := prepareNativeType(file.Config, previous, desired.bindings.NoteTypes[file.Path], original)
		if err != nil {
			return fmt.Errorf("prepare note type %q: %w", file.Path, err)
		}
		operation := &collectionpb.PlannedNoteType{Path: file.Path, Value: value, SortFieldIndex: file.Config.Config.SortFieldIdx}
		for _, template := range file.Config.Templates {
			if template.TargetDeck != "" {
				if operation.TemplateDecks == nil {
					operation.TemplateDecks = map[string]string{}
				}
				operation.TemplateDecks[template.Name] = template.TargetDeck
			}
		}
		desired.noteTypes[file.Path] = operation
		if value.Id != 0 {
			desired.typesByID[value.Id] = operation
		}
	}
	return nil
}

// resolveNotes orders named content and prepares native notes and destination assignments.
func resolveNotes(desired *desiredState, files []*collectionpb.NoteConfigFile, home string, current snapshot, paths map[string]bool) error {
	for _, file := range files {
		if paths[file.Path] {
			return fmt.Errorf("duplicate resource path %q", file.Path)
		}
		paths[file.Path] = true
		if desired.bindings.Decks[file.Path] != nil || desired.bindings.NoteTypes[file.Path] != nil {
			return fmt.Errorf("note path %q is bound to another entity", file.Path)
		}
		id := desired.bindings.Notes[file.Path].GetId()
		note, err := prepareNote(file, home, *desired, current.notes[id])
		if err != nil {
			return fmt.Errorf("note %q: %w", file.Path, err)
		}
		desired.notes[file.Path] = note
		desired.cards[file.Path] = &collectionpb.PlannedCardAssignments{Path: file.Path, Deck: home, Cards: file.Config.Cards}
		if note.Value.Id != 0 {
			desired.notesByID[note.Value.Id] = note
		}
	}
	return nil
}

// prepareNote resolves type and deck paths and translates named fields into native order.
func prepareNote(file *collectionpb.NoteConfigFile, home string, desired desiredState, previous *notes.Note) (*collectionpb.PlannedNote, error) {
	config := file.Config
	if config == nil {
		return nil, fmt.Errorf("note configuration is required")
	}
	definition := desired.noteTypes[config.NoteType]
	if definition == nil {
		return nil, fmt.Errorf("unknown declared note-type path %q", config.NoteType)
	}
	settings := definition.Value
	if previous != nil && previous.NotetypeId != settings.Id {
		return nil, fmt.Errorf("retained note cannot change note type")
	}
	if len(config.Fields) != len(settings.Fields) {
		return nil, fmt.Errorf("field count differs from declared note type")
	}
	value := &notes.Note{Id: previous.GetId(), Tags: config.Tags}
	for _, field := range settings.Fields {
		body, exists := config.Fields[field.Name]
		if !exists {
			return nil, fmt.Errorf("missing field %q", field.Name)
		}
		value.Fields = append(value.Fields, body)
	}
	ordinals := map[int64]bool{}
	for _, card := range config.Cards {
		if card.Ordinal < 0 || card.Ordinal > math.MaxInt32 || ordinals[card.Ordinal] {
			return nil, fmt.Errorf("invalid or duplicate card ordinal")
		}
		ordinals[card.Ordinal] = true
		if settings.Config.Kind != notetypes.Notetype_Config_KIND_CLOZE && card.Ordinal >= int64(len(settings.Templates)) {
			return nil, fmt.Errorf("no template at card ordinal %d", card.Ordinal)
		}
		target := desired.decksByPath[card.GetDeck()]
		if target == nil || target.Value.GetNormal() == nil {
			return nil, fmt.Errorf("card requires declared normal-deck path %q", card.GetDeck())
		}
	}
	if home != "" {
		target := desired.decksByPath[home]
		if target == nil || target.Value.GetNormal() == nil {
			return nil, fmt.Errorf("default deck must be normal")
		}
	}
	return &collectionpb.PlannedNote{Path: file.Path, Value: value, NoteType: config.NoteType, Deck: home}, nil
}

// resolveTypeReferences converts readable references into native settings for comparison.
func resolveTypeReferences(definition *collectionpb.NoteTypeConfig, declared map[string]*collectionpb.PlannedDeck) error {
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
			deck := declared[template.TargetDeck]
			if deck == nil || deck.Value.GetNormal() == nil {
				return fmt.Errorf("unknown normal-deck path %q", template.TargetDeck)
			}
			template.Config.TargetDeckId = deck.Value.Id
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
	moved, err := moveResources(stored.GetResources(), configs.GetConfig().GetMoved())
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
