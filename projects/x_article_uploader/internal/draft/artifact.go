package draft

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/markdown"
)

// Artifact is the draft artifact the converter emitted.
type Artifact struct {
	Title       string                  `json:"title"`
	Document    map[string]any          `json:"content_state"`
	Locators    []markdown.ImageLocator `json:"image_locators"`
	Banner      *markdown.ImageSource   `json:"banner_locator,omitempty"`
	Diagnostics []markdown.Diagnostic   `json:"diagnostics"`
}

// ReadArtifact loads a converted draft artifact.
func ReadArtifact(path string) (*Artifact, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read artifact: %w", err)
	}
	var artifact Artifact
	if err := json.Unmarshal(raw, &artifact); err != nil {
		return nil, fmt.Errorf("decode artifact %s: %w", path, err)
	}
	return &artifact, nil
}

// ImageBytes resolves a locator against the post package it records and returns
// the image's bytes and resolved path, refusing a path that escapes the post
// directory or bytes whose digest no longer matches.
//
// Containment is enforced by an os.Root anchored at the workspace, not by a
// pathname check followed by a read: the recorded package could itself contain
// `..`, and an apparently safe relative path could be a symlink whose target
// lies outside the post package. os.Root resolves every component of a name
// through descriptors and refuses a name whose symlink escapes the root, so the
// object whose containment is checked is the object read; a writer that swaps
// the canonical target for a symlink after a check cannot redirect the read,
// because no pathname is re-opened between the check and the read. The returned
// path is still the canonical target the read resolved to.
func ImageBytes(locator markdown.ImageSource, workspaceRoot string) ([]byte, string, error) {
	packagePath := filepath.FromSlash(locator.PostPackage)
	if filepath.IsAbs(packagePath) {
		return nil, "", fmt.Errorf("locator records absolute post package %q", locator.PostPackage)
	}
	cleaned := filepath.Clean(locator.Path)
	if filepath.IsAbs(cleaned) || escapes(cleaned) {
		return nil, "", fmt.Errorf("locator %q escapes the post package %s", locator.Path, locator.PostPackage)
	}

	// Anchor the whole resolution at the workspace root. OpenRoot follows
	// symlinks in the workspace path itself but confines every subsequent name
	// to that directory, so both the recorded package and the image reference
	// are resolved through descriptors and cannot escape the workspace.
	root, err := os.OpenRoot(workspaceRoot)
	if err != nil {
		return nil, "", fmt.Errorf("open workspace root %s: %w", workspaceRoot, err)
	}
	defer root.Close()

	packageRoot, err := root.OpenRoot(filepath.ToSlash(packagePath))
	if err != nil {
		return nil, "", fmt.Errorf("locator post package %s is not inside the workspace: %w", locator.PostPackage, err)
	}
	defer packageRoot.Close()

	content, err := packageRoot.ReadFile(cleaned)
	if err != nil {
		return nil, "", fmt.Errorf("locator %q is not readable inside the post package %s: %w", locator.Path, locator.PostPackage, err)
	}
	sum := sha256.Sum256(content)
	digest := hex.EncodeToString(sum[:])
	if digest != locator.Digest {
		return nil, "", fmt.Errorf("image %q digest %s does not match the artifact's %s", locator.Path, digest, locator.Digest)
	}
	if err := markdown.ValidateImageType(locator.Path, content, locator.MediaType); err != nil {
		return nil, "", fmt.Errorf("validate image %q: %w", locator.Path, err)
	}
	// The returned path is informational only: the bytes above were read through
	// the root, so the read never depended on this pathname. It is still
	// canonicalized so callers that report it name the file the read used.
	resolved := filepath.Join(workspaceRoot, packagePath, cleaned)
	if canonical, err := filepath.EvalSymlinks(resolved); err == nil {
		resolved = canonical
	}
	return content, resolved, nil
}

// escapes reports whether a cleaned relative path leaves its directory.
func escapes(cleaned string) bool {
	return cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator))
}
