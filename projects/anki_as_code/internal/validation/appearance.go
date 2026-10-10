package validation

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
)

// validateAppearance checks appearance identities and existing template structure.
func validateAppearance(d *model.CollectionDesiredState, s model.CollectionState) error {
	if len(d.Appearance) != len(s.Types) {
		return fmt.Errorf("appearance must declare every existing note type")
	}
	seen := map[int64]bool{}
	for i := range d.Appearance {
		appearance := d.Appearance[i]
		current := s.Types[appearance.Id]
		if current == nil || current.Name != appearance.Name || seen[appearance.Id] {
			return fmt.Errorf("unknown, renamed, or duplicate appearance note type %d", appearance.Id)
		}
		seen[appearance.Id] = true
		if len(appearance.Templates) != len(current.Appearance.Templates) {
			return fmt.Errorf("note type %d template structure cannot change", appearance.Id)
		}
		sort.Slice(appearance.Templates, func(i, j int) bool { return appearance.Templates[i].Ordinal < appearance.Templates[j].Ordinal })
		values := []string{appearance.Css}
		for j, template := range appearance.Templates {
			original := current.Appearance.Templates[j]
			if template.Ordinal != original.Ordinal || template.Name != original.Name {
				return fmt.Errorf("note type %d template identities cannot change", appearance.Id)
			}
			values = append(values, template.Front, template.Back, template.BrowserFront, template.BrowserBack, template.BrowserFont)
		}
		for _, value := range values {
			if !utf8.ValidString(value) || strings.ContainsRune(value, 0) {
				return fmt.Errorf("note type %d appearance must contain valid UTF-8 without NUL", appearance.Id)
			}
		}
	}
	return nil
}
