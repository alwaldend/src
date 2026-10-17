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

// resolveExportOutputPath follows symlinks through an absolute path's nearest existing ancestor.
// The remaining suffix may name directories that Export has yet to create.
func resolveExportOutputPath(path string) (string, error) {
	ancestor := path
	for {
		resolved, err := filepath.EvalSymlinks(ancestor)
		if err == nil {
			suffix, err := filepath.Rel(ancestor, path)
			if err != nil {
				return "", fmt.Errorf("resolve export suffix: %w", err)
			}
			return filepath.Join(resolved, suffix), nil
		}
		parent := filepath.Dir(ancestor)
		if !os.IsNotExist(err) || parent == ancestor {
			return "", fmt.Errorf("resolve export path %q: %w", path, err)
		}
		ancestor = parent
	}
}
