package plan

import (
	"fmt"
	"math"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/protobuf/proto"
)

// stateMember is the shared identity contract of stored field and template bindings.
type stateMember interface {
	comparable
	// GetId returns a stable native member ID when one is available.
	GetId() int64
	// GetOrdinal returns the member's position when its binding was recorded.
	GetOrdinal() uint32
}

// memberReference identifies the original member independently of its current ordinal.
type memberReference struct {
	// name is the original stored name, including before an explicit rename.
	name string
	// id is the authoritative stable identity, or zero when Anki provides none.
	id int64
	// bound distinguishes a missing recorded member from a new declaration.
	bound bool
}

// nativeMember indexes one current member by its native name and optional stable ID.
type nativeMember struct {
	// name is the current native member name.
	name string
	// id is its optional stable native identity.
	id int64
}

// originalMember resolves an explicit rename against the original stored member names.
// Saved ordinals identify original bindings, never arbitrary current native members.
func originalMember[T stateMember](name string, binding T, original map[string]T) (memberReference, error) {
	var missing T
	result := memberReference{name: name}
	if binding == missing {
		return result, nil
	}
	result.id, result.bound = binding.GetId(), true
	if result.id != 0 {
		return result, nil
	}
	found := false
	for oldName, old := range original {
		if old.GetId() == 0 && old.GetOrdinal() == binding.GetOrdinal() {
			if found {
				return result, fmt.Errorf("ambiguous stored member ordinal for %q", name)
			}
			result.name, found = oldName, true
		}
	}
	if !found {
		return result, fmt.Errorf("missing original member binding for %q", name)
	}
	return result, nil
}

// matchMember locates authoritative IDs or original names without trusting stale ordinals.
func matchMember(reference memberReference, members []nativeMember) (int, error) {
	found := -1
	for index, member := range members {
		matches := reference.id != 0 && reference.id == member.id || reference.id == 0 && reference.name == member.name
		if matches {
			if found >= 0 {
				return -1, fmt.Errorf("ambiguous native member %q", reference.name)
			}
			found = index
		}
	}
	if found < 0 && reference.bound && reference.id == 0 {
		return -1, fmt.Errorf("bound member %q is unavailable; import its current native identity before renaming", reference.name)
	}
	return found, nil
}

// prepareNativeType maps members on owned configs and retains detached native metadata.
func prepareNativeType(wanted *collectionpb.NoteTypeConfig, previous *notetypes.Notetype, moved, original *collectionpb.NoteTypeState) (*notetypes.Notetype, error) {
	previous = proto.CloneOf(previous)
	value := &notetypes.Notetype{Id: previous.GetId(), Name: wanted.Name, Config: proto.CloneOf(wanted.Config)}
	fields, templates := []nativeMember{}, []nativeMember{}
	if previous != nil {
		for _, field := range previous.Fields {
			fields = append(fields, nativeMember{name: field.Name, id: field.Config.GetId()})
		}
		for _, template := range previous.Templates {
			templates = append(templates, nativeMember{name: template.Name, id: template.Config.GetId()})
		}
	} else {
		moved, original = nil, nil
	}
	usedFields, usedTemplates := make([]bool, len(fields)), make([]bool, len(templates))
	for index, source := range wanted.Fields {
		reference, err := originalMember(source.Name, moved.GetFields()[source.Name], original.GetFields())
		if err != nil {
			return nil, fmt.Errorf("resolve field: %w", err)
		}
		ordinal, err := matchMember(reference, fields)
		if err != nil {
			return nil, fmt.Errorf("match field: %w", err)
		}
		field := source
		if ordinal >= 0 {
			if usedFields[ordinal] {
				return nil, fmt.Errorf("fields claim the same native member %q", source.Name)
			}
			usedFields[ordinal] = true
			field.Ord = &generic.UInt32{Val: uint32(ordinal)}
			field.Config.Id, field.Config.Other = previous.Fields[ordinal].Config.Id, previous.Fields[ordinal].Config.Other
		}
		if previous != nil && uint32(index) == wanted.Config.SortFieldIdx {
			value.Config.SortFieldIdx = math.MaxUint32
			if field.Ord != nil {
				value.Config.SortFieldIdx = field.Ord.Val
			}
		}
		value.Fields = append(value.Fields, field)
	}
	for _, source := range wanted.Templates {
		reference, err := originalMember(source.Name, moved.GetTemplates()[source.Name], original.GetTemplates())
		if err != nil {
			return nil, fmt.Errorf("resolve template: %w", err)
		}
		ordinal, err := matchMember(reference, templates)
		if err != nil {
			return nil, fmt.Errorf("match template: %w", err)
		}
		template := &notetypes.Notetype_Template{Name: source.Name, Config: source.Config}
		if ordinal >= 0 {
			if usedTemplates[ordinal] {
				return nil, fmt.Errorf("templates claim the same native member %q", source.Name)
			}
			usedTemplates[ordinal] = true
			template.Ord = &generic.UInt32{Val: uint32(ordinal)}
			template.Config.Id, template.Config.Other = previous.Templates[ordinal].Config.Id, previous.Templates[ordinal].Config.Other
		}
		value.Templates = append(value.Templates, template)
	}
	return value, nil
}

// unchangedNoteType compares editable native values and the already prepared member mappings.
func unchangedNoteType(old *notetypes.Notetype, wanted *collectionpb.PlannedNoteType) bool {
	editable := editableType(wanted.Value)
	editable.Config.SortFieldIdx = wanted.SortFieldIndex
	if !proto.Equal(editableType(old), editable) {
		return false
	}
	for index, field := range wanted.Value.Fields {
		if field.Ord == nil || field.Ord.Val != uint32(index) {
			return false
		}
	}
	for index, template := range wanted.Value.Templates {
		if template.Ord == nil || template.Ord.Val != uint32(index) {
			return false
		}
		// A declared new target has no native ID yet, unlike an omitted target.
		if wanted.TemplateDecks[template.Name] != "" && template.Config.TargetDeckId == 0 {
			return false
		}
	}
	return true
}

// validateFilteredTemplate protects unmanaged resident cards from structural template edits.
func validateFilteredTemplate(typeID int64, ordinal int, desired desiredState, current snapshot) error {
	entry := desired.typesByID[typeID]
	if entry == nil {
		if current.typePaths[typeID] == "" {
			return nil
		}
		return fmt.Errorf("filtered card requires existing note type %d", typeID)
	}
	old, settings := current.noteTypes[typeID], entry.Value
	if old.Config.Kind != settings.Config.Kind {
		return fmt.Errorf("cannot change note type kind while filtered cards remain; empty the filtered deck in Anki first")
	}
	// All cloze card ordinals share template zero.
	if old.Config.Kind == notetypes.Notetype_Config_KIND_CLOZE {
		ordinal = 0
	}
	if ordinal >= len(old.Templates) || ordinal >= len(settings.Templates) {
		return fmt.Errorf("cannot remove a filtered card template; empty the filtered deck in Anki first")
	}
	template := settings.Templates[ordinal]
	if template.Ord == nil || template.Ord.Val != uint32(ordinal) {
		return fmt.Errorf("cannot replace or reorder a filtered card template; empty the filtered deck in Anki first")
	}
	return nil
}

// prepareNoteType supplies protobuf defaults used by reference resolution.
// Member names must identify unambiguous state addresses; Anki validates their
// syntax, structural constraints and template bodies during apply.
func prepareNoteType(definition *collectionpb.NoteTypeConfig) error {
	if definition == nil {
		return fmt.Errorf("note-type configuration is required")
	}
	if definition.Config == nil {
		definition.Config = &notetypes.Notetype_Config{}
	}
	fields := map[string]bool{}
	for _, field := range definition.Fields {
		if field == nil {
			return fmt.Errorf("nil field declaration")
		}
		if fields[field.Name] {
			return fmt.Errorf("duplicate field reference %q", field.Name)
		}
		fields[field.Name] = true
		if field.Ord != nil {
			return fmt.Errorf("field ordinal is native-only; use array order")
		}
		if field.Config == nil {
			field.Config = &notetypes.Notetype_Field_Config{}
		}
	}
	templates := map[string]bool{}
	for _, template := range definition.Templates {
		if template == nil {
			return fmt.Errorf("nil template declaration")
		}
		if templates[template.Name] {
			return fmt.Errorf("duplicate template reference %q", template.Name)
		}
		templates[template.Name] = true
		if template.Config == nil {
			template.Config = &notetypes.Notetype_Template_Config{}
		}
	}
	return nil
}
