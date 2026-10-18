package plan

import (
	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"github.com/ankitects/anki/proto/anki/notetypes"
)

// bindingInputs constructs protobuf input data from resolved native identities.
func bindingInputs(desired desiredState) *ankiascodepb.Resources {
	inputs := &ankiascodepb.Resources{
		Decks:     map[string]*ankiascodepb.DeckState{},
		Notes:     map[string]*ankiascodepb.NoteState{},
		NoteTypes: map[string]*ankiascodepb.NoteTypeState{},
		Presets:   desired.bindings.Presets,
	}
	for _, deck := range desired.decksByPath {
		inputs.Decks[deck.Path] = &ankiascodepb.DeckState{Id: deck.Value.Id}
	}
	for _, note := range desired.notes {
		inputs.Notes[note.Path] = &ankiascodepb.NoteState{Id: note.Value.Id}
	}
	for _, definition := range desired.noteTypes {
		inputs.NoteTypes[definition.Path] = noteTypeBindingInput(definition.Value)
	}
	return inputs
}

// noteTypeBindingInput projects native identities from a config into protobuf input data.
func noteTypeBindingInput(definition *notetypes.Notetype) *ankiascodepb.NoteTypeState {
	value := &ankiascodepb.NoteTypeState{
		Id:        definition.Id,
		Fields:    map[string]*ankiascodepb.FieldState{},
		Templates: map[string]*ankiascodepb.TemplateState{},
	}
	for ordinal, field := range definition.Fields {
		value.Fields[field.Name] = &ankiascodepb.FieldState{Id: field.Config.GetId(), Ordinal: uint32(ordinal)}
	}
	for ordinal, template := range definition.Templates {
		value.Templates[template.Name] = &ankiascodepb.TemplateState{Id: template.Config.GetId(), Ordinal: uint32(ordinal)}
	}
	return value
}
