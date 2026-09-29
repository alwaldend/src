// Package e2e holds the project's end-to-end checks. They exercise conversion
// against the real blog posts discovered from the content tree and publication
// against a recorded-response client, so no check contacts the X API or needs
// live credentials.
package e2e

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

var (
	rasterRender = flag.String("raster-render", "", "runfile path of the diagram WebP render")
	rasterSource = flag.String("raster-source", "", "runfile path of the diagram SVG render")
)

// workspaceRoot returns the repository root from the runfiles tree, so a
// post's recorded package resolves to the same bytes the build used.
func workspaceRoot(t *testing.T) string {
	t.Helper()
	root := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if root != "" {
		return root
	}
	// In a sandboxed test the runfiles tree holds the repository below its
	// `_main` directory, so a known declared file locates that root.
	resolved, err := runfiles.Rlocation("_main/projects/alwaldend.com/content/blog/README.md")
	if err != nil {
		t.Fatalf("locate runfiles root: %v", err)
	}
	// The runfiles symlink tree points each declared file at its source in the
	// workspace, so canonicalize before deriving the root: the converter and
	// the draft command both resolve containment on canonical paths, and a root
	// whose files are symlinks elsewhere would make an in-tree post look like it
	// escapes the tree.
	resolved, err = filepath.EvalSymlinks(resolved)
	if err != nil {
		t.Fatalf("canonicalize runfiles root: %v", err)
	}
	// <runfiles>/_main/projects/alwaldend.com/content/blog/README.md -> <runfiles>/_main.
	blogDir := filepath.Dir(resolved)
	return filepath.Clean(filepath.Join(blogDir, "..", "..", "..", ".."))
}

// rlocation resolves a declared runfile.
func rlocation(t *testing.T, path string) string {
	t.Helper()
	resolved, err := runfiles.Rlocation(path)
	if err != nil {
		t.Fatalf("resolve runfile %s: %v", path, err)
	}
	return resolved
}

// listPosts returns every post directory under a content tree that holds an
// index.md, so the covered set is derived from the tree rather than listed.
func listPosts(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	posts := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		index := filepath.Join(root, entry.Name(), "index.md")
		if _, err := os.Stat(index); err != nil {
			continue
		}
		posts = append(posts, filepath.Join(root, entry.Name()))
	}
	return posts, nil
}

// readPost reads a post's source.
func readPost(t *testing.T, dir string) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(dir, "index.md"))
	if err != nil {
		t.Fatalf("read post %s: %v", dir, err)
	}
	return raw
}

// hasDiagnostic reports whether a diagnostic with the code was reported.
func hasDiagnostic(codes []string, code string) bool {
	for _, c := range codes {
		if c == code {
			return true
		}
	}
	return false
}

func joinCodes(codes []string) string { return strings.Join(codes, ", ") }
