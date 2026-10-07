package validation

import (
	"fmt"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/model"
)

// ValidateBodies rejects field separators reserved by Anki's database representation.
func ValidateBodies(n model.Note) error {
	for i, body := range n.Bodies {
		if strings.Contains(body, "\x1f") || strings.Contains(body, "\x00") {
			return fmt.Errorf("note %d field %d contains reserved separators", n.NoteId, i)
		}
	}
	return nil
}

// ValidDeck checks that a deck title has valid nonempty hierarchy components.
func ValidDeck(name string) error {
	for _, part := range strings.Split(name, "::") {
		if strings.TrimSpace(part) == "" || strings.TrimSpace(part) != part || strings.ContainsAny(part, "\x1f\x00\r\n") {
			return fmt.Errorf("invalid deck name %q", name)
		}
	}
	return nil
}
