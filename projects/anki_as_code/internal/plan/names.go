package plan

import "golang.org/x/text/cases"

// folded matches Anki's case-insensitive resource names for ownership resolution.
func folded(value string) string {
	return cases.Fold().String(value)
}
