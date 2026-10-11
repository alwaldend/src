package validation

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode/utf8"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/ankiformat"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
	"github.com/ankitects/anki/proto/anki/notetypes"
)

// validateNoteTypes resolves identities and validates complete note-type declarations.
func validateNoteTypes(desired *model.CollectionDesiredState, current model.CollectionState) error {
	nextID := int64(0)
	existingNames := map[string]int64{}
	for id, settings := range current.NoteTypes {
		nextID = max(nextID, id)
		existingNames[ankiformat.Folded(settings.Name)] = id
	}
	for _, settings := range desired.NoteTypes {
		if settings == nil {
			return fmt.Errorf("nil note type declaration")
		}
		nextID = max(nextID, settings.Id)
	}
	ids, names := map[int64]bool{}, map[string]bool{}
	for _, settings := range desired.NoteTypes {
		if settings.Id == 0 {
			settings.Id = existingNames[ankiformat.Folded(settings.Name)]
			if settings.Id == 0 {
				if nextID == math.MaxInt64 {
					return fmt.Errorf("note type identity range exhausted")
				}
				nextID++
				settings.Id = nextID
			}
		}
		if settings.Id <= 0 || ids[settings.Id] || names[ankiformat.Folded(settings.Name)] {
			return fmt.Errorf("invalid or duplicate note type %d/%q", settings.Id, settings.Name)
		}
		ids[settings.Id], names[ankiformat.Folded(settings.Name)] = true, true
		if err := validateNoteType(settings, current.NoteTypes[settings.Id], desired.Decks); err != nil {
			return fmt.Errorf("validate note type %d: %w", settings.Id, err)
		}
	}
	return nil
}

// validateNoteType checks configuration, ordered fields, and template structure.
func validateNoteType(settings, current *collectionpb.NoteTypeSettings, decks []model.Deck) error {
	if err := validateSettingName(settings.Name); err != nil {
		return fmt.Errorf("note type name: %w", err)
	}
	if settings.Config == nil {
		return fmt.Errorf("configuration is required")
	}
	if settings.Config.Kind != notetypes.Notetype_Config_KIND_NORMAL && settings.Config.Kind != notetypes.Notetype_Config_KIND_CLOZE {
		return fmt.Errorf("unsupported note type kind")
	}
	if len(settings.Fields) == 0 || int64(settings.Config.SortFieldIdx) >= int64(len(settings.Fields)) {
		return fmt.Errorf("fields and a valid sort field are required")
	}
	if err := validateFields(settings.Fields); err != nil {
		return fmt.Errorf("validate fields: %w", err)
	}
	if err := validateTemplates(settings.Templates, current, decks); err != nil {
		return fmt.Errorf("validate templates: %w", err)
	}
	if settings.Config.Kind == notetypes.Notetype_Config_KIND_CLOZE && len(settings.Templates) != 1 {
		return fmt.Errorf("cloze note types require exactly one template")
	}
	if err := validatePresentation(settings.Config.Css, settings.Config.LatexPre, settings.Config.LatexPost); err != nil {
		return fmt.Errorf("validate configuration text: %w", err)
	}
	return nil
}

// validateFields checks unique names and configuration in the declared storage order.
func validateFields(fields []*notetypes.Notetype_Field) error {
	names := map[string]bool{}
	for i, field := range fields {
		if field == nil {
			return fmt.Errorf("nil field declaration")
		}
		if err := validateSettingName(field.Name); err != nil {
			return fmt.Errorf("field name: %w", err)
		}
		if names[ankiformat.Folded(field.Name)] {
			return fmt.Errorf("duplicate field name %q", field.Name)
		}
		names[ankiformat.Folded(field.Name)] = true
		if field.Ord != nil && uint64(field.Ord.Val) != uint64(i) {
			return fmt.Errorf("field %q ordinal must match its list position", field.Name)
		}
		field.Ord = nil
		if field.Config == nil {
			field.Config = &notetypes.Notetype_Field_Config{}
		}
		if err := validatePresentation(field.Config.FontName, field.Config.Description); err != nil {
			return fmt.Errorf("field %q configuration: %w", field.Name, err)
		}
	}
	return nil
}

// validateTemplates checks contiguous ordinals, unique names, and configured target decks.
func validateTemplates(templates []*collectionpb.TemplateSettings, current *collectionpb.NoteTypeSettings, decks []model.Deck) error {
	if len(templates) == 0 {
		return fmt.Errorf("at least one template is required")
	}
	for _, template := range templates {
		if template == nil {
			return fmt.Errorf("nil template declaration")
		}
	}
	sort.Slice(templates, func(i, j int) bool { return templates[i].Ordinal < templates[j].Ordinal })
	names := map[string]bool{}
	for i, template := range templates {
		if template.Ordinal != int64(i) {
			return fmt.Errorf("template ordinals must be unique and contiguous from zero")
		}
		if err := validateSettingName(template.Name); err != nil {
			return fmt.Errorf("template name: %w", err)
		}
		if names[ankiformat.Folded(template.Name)] {
			return fmt.Errorf("duplicate template name %q", template.Name)
		}
		names[ankiformat.Folded(template.Name)] = true
		if template.Config == nil {
			return fmt.Errorf("template %q needs configuration", template.Name)
		}
		if err := validateTemplateConfig(template, current, decks); err != nil {
			return fmt.Errorf("template %q: %w", template.Name, err)
		}
	}
	return nil
}

// validateTemplateConfig checks text and optional normal-deck targeting.
func validateTemplateConfig(template *collectionpb.TemplateSettings, current *collectionpb.NoteTypeSettings, decks []model.Deck) error {
	config := template.Config
	if err := validatePresentation(config.QFormat, config.AFormat, config.QFormatBrowser, config.AFormatBrowser, config.BrowserFontName); err != nil {
		return fmt.Errorf("presentation: %w", err)
	}
	if config.TargetDeckId == 0 {
		return nil
	}
	if current != nil {
		for _, old := range current.Templates {
			if old.Name == template.Name && old.Config.TargetDeckId == config.TargetDeckId {
				return nil
			}
		}
	}
	for _, deck := range decks {
		if deck.Id == config.TargetDeckId && deck.Normal != nil {
			return nil
		}
	}
	return fmt.Errorf("target deck %d is not a declared normal deck", config.TargetDeckId)
}

// validateSettingName checks nonempty configuration names without reserved separators.
func validateSettingName(name string) error {
	if strings.TrimSpace(name) == "" || strings.ContainsAny(name, "\x1f\r\n") {
		return fmt.Errorf("invalid name %q", name)
	}
	if err := validatePresentation(name); err != nil {
		return fmt.Errorf("name text: %w", err)
	}
	return nil
}

// validatePresentation rejects invalid UTF-8 and NUL in configuration text.
func validatePresentation(values ...string) error {
	for _, value := range values {
		if !utf8.ValidString(value) || strings.ContainsRune(value, 0) {
			return fmt.Errorf("text must contain valid UTF-8 without NUL")
		}
	}
	return nil
}
