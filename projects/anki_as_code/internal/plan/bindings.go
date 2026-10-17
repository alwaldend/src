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
		Presets:   desired.bindings.Presets,
	}
	for _, deck := range desired.decksByPath {
		inputs.Decks[deck.Path] = &collectionpb.DeckState{Id: deck.Value.Id}
	}
	for _, note := range desired.notes {
		inputs.Notes[note.Path] = &collectionpb.NoteState{Id: note.Value.Id}
	}
	for _, definition := range desired.noteTypes {
		inputs.NoteTypes[definition.Path] = noteTypeBindingInput(definition.Value)
	}
	return inputs
}

// noteTypeBindingInput projects native identities from a config into protobuf input data.
func noteTypeBindingInput(definition *notetypes.Notetype) *collectionpb.NoteTypeState {
	value := &collectionpb.NoteTypeState{
		Id:        definition.Id,
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
