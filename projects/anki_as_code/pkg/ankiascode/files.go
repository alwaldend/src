package ankiascode

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// expandPath resolves collection/archive paths without depending on configuration internals.
func expandPath(input string) (string, error) {
	if input == "~" || strings.HasPrefix(input, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory: %w", err)
		}
		input = filepath.Join(home, strings.TrimPrefix(strings.TrimPrefix(input, "~"), "/"))
	}
	path, err := filepath.Abs(input)
	if err != nil {
		return "", fmt.Errorf("resolve collection path: %w", err)
	}
	return path, nil
}
