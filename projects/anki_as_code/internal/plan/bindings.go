package plan

import (
	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"github.com/ankitects/anki/proto/anki/notetypes"
)

// bindingInputs constructs protobuf input data from resolved native identities.
func bindingInputs(desired desiredState) *collectionpb.Resources {
	inputs := &collectionpb.Resources{
		Decks:     map[string]*collectionpb.DeckState{},
		Notes:     map[string]*collectionpb.NoteState{},
		NoteTypes: map[string]*collectionpb.NoteTypeState{},
		Presets:   map[string]*collectionpb.PresetState{},
	}
	for _, deck := range desired.decks {
		inputs.Decks[deck.path] = &collectionpb.DeckState{Id: deck.id}
	}
	for _, note := range desired.notes {
		inputs.Notes[note.path] = &collectionpb.NoteState{Id: note.id}
	}
	for _, definition := range desired.noteTypes {
		inputs.NoteTypes[definition.path] = noteTypeBindingInput(definition.id, definition.native)
	}
	for alias, id := range desired.presetBindings {
		inputs.Presets[alias] = &collectionpb.PresetState{Id: id}
	}
	return inputs
}

// noteTypeBindingInput projects native identities from a config into protobuf input data.
func noteTypeBindingInput(id int64, definition *notetypes.Notetype) *collectionpb.NoteTypeState {
	value := &collectionpb.NoteTypeState{
		Id:        id,
		Fields:    map[string]*collectionpb.FieldState{},
		Templates: map[string]*collectionpb.TemplateState{},
	}
	for ordinal, field := range definition.Fields {
		value.Fields[field.Name] = &collectionpb.FieldState{Id: field.Config.GetId(), Ordinal: uint32(ordinal)}
	}
	for ordinal, template := range definition.Templates {
		value.Templates[template.Name] = &collectionpb.TemplateState{Id: template.Config.GetId(), Ordinal: uint32(ordinal)}
	}
	return value
}
