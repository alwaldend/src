package agent_system

import (
	"fmt"
	"os"
	"path/filepath"
)

// Bazel run starts in its runfiles tree. Resolve user paths against the
// source workspace without changing the process-wide working directory.
// Outside Bazel, retain ordinary working-directory-relative CLI behavior.
func commandWorkspaceRoot(requested string) (string, error) {
	base := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if base == "" {
		var err error
		base, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("resolve command directory: %w", err)
		}
	}
	if requested == "" {
		requested = base
	} else if !filepath.IsAbs(requested) {
		requested = filepath.Join(base, requested)
	}
	return filepath.Abs(requested)
}

func workspaceFilePaths(workspace string, paths ...*string) error {
	root, err := commandWorkspaceRoot(workspace)
	if err != nil {
		return err
	}
	for _, path := range paths {
		if *path != "" && !filepath.IsAbs(*path) {
			*path = filepath.Join(root, *path)
		}
	}
	return nil
}
