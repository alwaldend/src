// Command convert compiles a Markdown blog post into the draft artifact the
// draft command consumes. It performs no network access and needs no credentials.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/markdown"
)

func main() {
	source := flag.String("source", "", "path to the post's index.md")
	out := flag.String("out", "", "path to write the draft artifact JSON")
	postPackage := flag.String("post-package", "", "repository package the post's directory belongs to (default: derived from --workspace)")
	workspace := flag.String("workspace", ".", "repository root a derived --post-package is relative to")
	warningsAsErrors := flag.Bool("warnings-as-errors", true, "fail conversion on warnings; false still reports warnings and always fails on errors")
	flag.Parse()

	if err := runWithOptions(*source, *out, *postPackage, *workspace, *warningsAsErrors, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(source, out, postPackage, workspace string, stdout, stderr io.Writer) error {
	if err := runWithOptions(source, out, postPackage, workspace, true, stdout, stderr); err != nil {
		return fmt.Errorf("convert article: %w", err)
	}
	return nil
}

func runWithOptions(source, out, postPackage, workspace string, warningsAsErrors bool, stdout, stderr io.Writer) error {
	if source == "" {
		return fmt.Errorf("--source is required")
	}
	if out == "" {
		return fmt.Errorf("--out is required")
	}
	postDir := filepath.Dir(source)
	var handle *os.Root
	if postPackage == "" {
		// The post package is recorded in every image locator and resolved
		// against the workspace root at publication, so it must be a relative
		// repository path. An absolute --source would otherwise record an
		// absolute package, producing an artifact publication refuses. Derive
		// it from the workspace named by --workspace, refusing when that cannot
		// yield a relative package rather than recording an unusable one.
		derived, validated, err := postPackageWithin(workspace, postDir)
		if err != nil {
			return fmt.Errorf("derive --post-package: %w", err)
		}
		handle = validated
		postPackage = derived
	} else {
		// A supplied package is the directory publication will resolve the
		// locators against, so read the post's images through that same
		// directory handle rather than re-opening the source pathname: a
		// symlink retargeted after validation would otherwise supply bytes from
		// another tree while the locator still names the supplied package.
		validated, err := packageHandle(workspace, postPackage, postDir)
		if err != nil {
			return fmt.Errorf("resolve --post-package: %w", err)
		}
		handle = validated
	}
	defer handle.Close()

	// Read the source through the same descriptor-anchored handle the package
	// was validated with, not by its pathname: a directory symlink retargeted
	// between a pathname read and the package's identity check would otherwise
	// supply one post's Markdown while the artifact records another package's
	// images, and a source symlink out of the package would be read at all.
	raw, err := handle.ReadFile(filepath.Base(source))
	if err != nil {
		return fmt.Errorf("read source: %w", err)
	}

	converter := markdown.New(postDir, postPackage)
	converter.PackageRoot = handle
	converter.PostSource = source
	article, convertErr := converter.Convert(raw)
	if article != nil {
		printDiagnostics(stderr, article.Diagnostics)
	}
	if convertErr != nil {
		return fmt.Errorf("convert source %q: %w", source, convertErr)
	}
	if warningsAsErrors && len(article.Diagnostics) > 0 {
		return fmt.Errorf("convert source %q: conversion failed with %d warnings; use --warnings-as-errors=false to allow warnings", source, len(article.Diagnostics))
	}
	encoded, err := article.JSON()
	if err != nil {
		return fmt.Errorf("encode artifact: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := os.WriteFile(out, append(encoded, '\n'), 0o644); err != nil {
		return fmt.Errorf("write artifact: %w", err)
	}
	fmt.Fprintf(stdout, "converted %s -> %s (%d blocks)\n",
		source, out, len(article.Document.Blocks))
	return nil
}

// openWithin opens a workspace-relative name through a descriptor-anchored
// handle on the canonical workspace and confirms, by file identity, that the
// opened directory is postDir.
//
// Identity of the opened objects is compared rather than pathnames: a name
// could be retargeted between a pathname comparison and the open, and an open
// that follows the new target would read a directory that was never compared.
// Comparing the descriptors instead means the directory accepted for the
// package is the directory the reads use. os.Root refuses a name whose symlink
// escapes the workspace, so a name pointing outside it is refused.
func openWithin(canonicalWorkspace, name, postDir string) (*os.Root, error) {
	root, err := os.OpenRoot(canonicalWorkspace)
	if err != nil {
		return nil, fmt.Errorf("open workspace %s: %w", canonicalWorkspace, err)
	}
	defer root.Close()
	handle, err := root.OpenRoot(filepath.ToSlash(name))
	if err != nil {
		return nil, fmt.Errorf("package %q is not inside the workspace %s: %w", name, canonicalWorkspace, err)
	}
	postHandle, err := os.OpenRoot(postDir)
	if err != nil {
		handle.Close()
		return nil, fmt.Errorf("open post directory %s: %w", postDir, err)
	}
	defer postHandle.Close()
	opened, err := handle.Stat(".")
	if err != nil {
		handle.Close()
		return nil, fmt.Errorf("stat package %q: %w", name, err)
	}
	want, err := postHandle.Stat(".")
	if err != nil {
		handle.Close()
		return nil, fmt.Errorf("stat post directory %s: %w", postDir, err)
	}
	if !os.SameFile(opened, want) {
		handle.Close()
		return nil, fmt.Errorf("package %q names a different directory than the post directory %s; a locator is recorded relative to the post directory and resolved against the package, so the two must be the same directory",
			name, postDir)
	}
	return handle, nil
}

// packageHandle opens an explicitly supplied post package through a
// descriptor-anchored handle on the workspace and confirms it names the
// directory the post's locator paths are relative to.
//
// A locator records a path relative to the post directory and a package
// publication resolves that path against, so the two must denote the same
// directory or publication reads a different object than conversion did. The
// handle is what conversion reads images through, so a symlink retargeted after
// validation cannot supply bytes from another tree while the locator still
// names the supplied package.
func packageHandle(workspace, postPackage, postDir string) (*os.Root, error) {
	absWorkspace, err := filepath.Abs(workspace)
	if err != nil {
		return nil, fmt.Errorf("resolve workspace %s: %w", workspace, err)
	}
	if filepath.IsAbs(postPackage) {
		return nil, fmt.Errorf("post package %q is absolute, which publication refuses", postPackage)
	}
	cleaned := filepath.Clean(postPackage)
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("post package %q escapes the workspace %s", postPackage, absWorkspace)
	}
	canonicalWorkspace, err := filepath.EvalSymlinks(absWorkspace)
	if err != nil {
		return nil, fmt.Errorf("canonicalize workspace %s: %w", workspace, err)
	}
	return openWithin(canonicalWorkspace, cleaned, postDir)
}

// postPackageWithin derives the workspace-relative package for a post directory
// when --post-package is omitted, and returns a descriptor-anchored handle on
// the validated directory so the caller can read the post's images through the
// same traversal it validated. A locator records this value and publication
// resolves it against the workspace root through a descriptor-anchored handle,
// so the recorded package must name the canonical directory relative to the
// workspace publication will use: an absolute source directory would record an
// absolute package publication refuses, and a symlinked directory would record
// the link even though publication's read resolves the canonical target.
//
// The workspace root is taken from the caller rather than rediscovered by
// walking the filesystem. A root found by walking is whatever tree a symbolic
// link happens to reach, so a post directory under an in-workspace link into
// another checkout that carries its own repository marker would be recorded
// relative to that external tree, and publication resolving the package against
// its own workspace would then name an unrelated directory there instead of
// refusing the escape.
//
// Containment is checked on the canonical paths and then established by opening
// the recorded package through an os.Root anchored at the canonical workspace.
// os.Root resolves every component through descriptors and refuses a name whose
// symlink escapes the root, so the directory the recorded package names is the
// directory the traversal reached, and a writer that swaps a directory for a
// link cannot leave a package that was accepted but that publication cannot
// resolve. Relating the two canonical paths, rather than a lexical one, is what
// refuses a link whose target lies outside the workspace instead of recording a
// package relative to the target tree.
//
// The returned handle is the validated directory itself, kept open for the
// caller. The caller reads images through it rather than re-opening the source
// pathname, so a symlink retargeted after validation cannot redirect a read
// into another tree while the artifact still records the validated package.
func postPackageWithin(workspace, postDir string) (string, *os.Root, error) {
	absWorkspace, err := filepath.Abs(workspace)
	if err != nil {
		return "", nil, fmt.Errorf("resolve workspace %s: %w", workspace, err)
	}
	absPost, err := filepath.Abs(postDir)
	if err != nil {
		return "", nil, fmt.Errorf("resolve %s: %w", postDir, err)
	}
	// Canonicalize the workspace and the directory separately, then relate the
	// two canonical paths: a directory re-discovered from the canonical target
	// would move to whatever tree the link points at, so the escape check would
	// pass for a link into another repository.
	canonicalWorkspace, err := filepath.EvalSymlinks(absWorkspace)
	if err != nil {
		return "", nil, fmt.Errorf("canonicalize workspace %s: %w", workspace, err)
	}
	canonical, err := filepath.EvalSymlinks(absPost)
	if err != nil {
		return "", nil, fmt.Errorf("canonicalize %s: %w", postDir, err)
	}
	recorded, err := filepath.Rel(canonicalWorkspace, canonical)
	if err != nil {
		return "", nil, fmt.Errorf("relate %s to workspace %s: %w", postDir, canonicalWorkspace, err)
	}
	if recorded == ".." || strings.HasPrefix(recorded, ".."+string(filepath.Separator)) {
		return "", nil, fmt.Errorf("%s resolves to %s, outside the workspace %s; pass --post-package explicitly",
			postDir, canonical, canonicalWorkspace)
	}

	// Anchor the resolution at the canonical workspace root: opening the
	// recorded package confirms publication can resolve it, and comparing the
	// opened handle with the post directory refuses a package that has since
	// become a link out of the workspace.
	handle, err := openWithin(canonicalWorkspace, recorded, postDir)
	if err != nil {
		return "", nil, fmt.Errorf("%v; pass --post-package explicitly", err)
	}
	return filepath.ToSlash(recorded), handle, nil
}

// printDiagnostics writes every diagnostic to w.
func printDiagnostics(w io.Writer, diagnostics []markdown.Diagnostic) {
	for _, diagnostic := range diagnostics {
		fmt.Fprintln(w, diagnostic.String())
	}
}
