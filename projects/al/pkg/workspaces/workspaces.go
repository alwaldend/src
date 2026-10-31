package workspaces

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Discover returns the root and declared nested modules in ignore-file order.
// Entries without MODULE.bazel are caches or source exclusions, not workspaces.
func Discover(root string) ([]string, error) {
	if info, err := os.Stat(filepath.Join(root, "MODULE.bazel")); err != nil {
		return nil, fmt.Errorf("inspect repository MODULE.bazel: %w", err)
	} else if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("repository MODULE.bazel is not a regular file")
	}
	file, err := os.Open(filepath.Join(root, ".bazelignore"))
	if err != nil {
		return nil, fmt.Errorf("open workspace boundaries .bazelignore: %w", err)
	}
	paths := []string{"."}
	seen := map[string]bool{".": true}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		path := strings.TrimSpace(scanner.Text())
		if path == "" || strings.HasPrefix(path, "#") {
			continue
		}
		path = filepath.Clean(path)
		if !filepath.IsLocal(path) {
			_ = file.Close()
			return nil, fmt.Errorf("nonlocal workspace boundary %q", path)
		}
		if seen[path] {
			continue
		}
		info, err := os.Stat(filepath.Join(root, path, "MODULE.bazel"))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			_ = file.Close()
			return nil, fmt.Errorf("inspect workspace %q: %w", path, err)
		}
		if !info.Mode().IsRegular() {
			_ = file.Close()
			return nil, fmt.Errorf("invalid MODULE.bazel in %s", path)
		}
		paths = append(paths, path)
		seen[path] = true
	}
	if err := scanner.Err(); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("read workspace boundaries: %w", err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close workspace boundaries: %w", err)
	}
	return paths, nil
}
