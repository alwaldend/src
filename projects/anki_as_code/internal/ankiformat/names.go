package ankiformat

import "golang.org/x/text/cases"

// Folded normalizes text for case-insensitive Anki name comparisons.
func Folded(s string) string {
	return cases.Fold().String(s)
}
