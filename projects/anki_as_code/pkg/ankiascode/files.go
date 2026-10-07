package ankiascode

import (
	"errors"
	"fmt"
	"io"
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

// copyArchive streams a source package into a newly created staging file and reports close errors.
func copyArchive(source, destination string) (err error) {
	input, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open source archive: %w", err)
	}
	defer func() {
		if closeErr := input.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close source archive: %w", closeErr))
		}
	}()
	output, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("create staging archive: %w", err)
	}
	defer func() {
		if closeErr := output.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close staging archive: %w", closeErr))
		}
	}()
	if _, err := io.Copy(output, input); err != nil {
		return fmt.Errorf("copy archive bytes: %w", err)
	}
	return nil
}
