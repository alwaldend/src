package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/markdown"
)

// TestRunDerivesRelativePostPackage asserts that omitting --post-package
// with an absolute --source and an explicit --workspace records a
// workspace-relative package in each locator, rather than the absolute source
// directory that publication refuses.
func TestRunDerivesRelativePostPackage(t *testing.T) {
	// An explicit workspace, so the derived package does not depend on the
	// test's own checkout or on any repository marker.
	workspace := t.TempDir()
	relPackage := filepath.Join("content", "blog", "example")
	dir := filepath.Join(workspace, relPackage)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create post directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "diagram.png"), tinyPNGForCommandTest, 0o644); err != nil {
		t.Fatalf("write image: %v", err)
	}
	source := filepath.Join(dir, "index.md")
	body := "---\ntitle: Example\n---\n\n![diagram](./diagram.png)\n"
	if err := os.WriteFile(source, []byte(body), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	out := filepath.Join(dir, "artifact.json")

	var stdout, stderr bytes.Buffer
	if err := run(source, out, "", workspace, &stdout, &stderr); err != nil {
		t.Fatalf("convert: %v (stderr: %s)", err, stderr.String())
	}
	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	var artifact struct {
		Locators []struct {
			PostPackage string `json:"post_package"`
			Path        string `json:"path"`
		} `json:"image_locators"`
	}
	if err := json.Unmarshal(raw, &artifact); err != nil {
		t.Fatalf("decode artifact: %v", err)
	}
	if len(artifact.Locators) != 1 {
		t.Fatalf("expected one locator, got %d", len(artifact.Locators))
	}
	if filepath.IsAbs(artifact.Locators[0].PostPackage) {
		t.Errorf("locator records absolute post package %q, which publication refuses",
			artifact.Locators[0].PostPackage)
	}
	if want := filepath.ToSlash(relPackage); artifact.Locators[0].PostPackage != want {
		t.Errorf("locator records post package %q, expected the workspace-relative %q",
			artifact.Locators[0].PostPackage, want)
	}
}

// tinyPNGForCommandTest is a complete, decodable one-pixel PNG.
var tinyPNGForCommandTest = []byte{
	0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A,
	0x00, 0x00, 0x00, 0x0d, 'I', 'H', 'D', 'R',
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	0x89, 0x00, 0x00, 0x00, 0x0a, 'I', 'D', 'A', 'T',
	0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00, 0x05,
	0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00,
	0x00, 0x00, 'I', 'E', 'N', 'D', 0xae, 0x42, 0x60, 0x82,
}

// TestPostPackageWithinRefusesEscapingSymlink asserts a post directory reached
// through an in-workspace symlink whose target lies outside the workspace is
// refused rather than recorded as a package publication's descriptor-anchored
// read would reject.
func TestPostPackageWithinRefusesEscapingSymlink(t *testing.T) {
	workspace := t.TempDir()
	// The real post lives outside the workspace; the link inside it is an
	// apparently safe relative path whose canonical target escapes the root.
	outside := filepath.Join(t.TempDir(), "example")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatalf("create outside post: %v", err)
	}
	link := filepath.Join(workspace, "content", "blog", "example")
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatalf("create link parent: %v", err)
	}
	if err := os.Symlink(outside, link); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	if _, handle, err := postPackageWithin(workspace, link); err == nil {
		handle.Close()
		t.Fatal("expected a post directory symlinked outside the workspace to be refused")
	}
}

// TestPostPackageWithinCanonicalizesInWorkspaceSymlink asserts an in-workspace
// symlink still derives the relative package of the directory it resolves to,
// so the recorded package matches the object publication reads.
func TestPostPackageWithinCanonicalizesInWorkspaceSymlink(t *testing.T) {
	workspace := t.TempDir()
	real := filepath.Join(workspace, "content", "blog", "real")
	if err := os.MkdirAll(real, 0o755); err != nil {
		t.Fatalf("create post: %v", err)
	}
	link := filepath.Join(workspace, "content", "blog", "alias")
	if err := os.Symlink(real, link); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	got, handle, err := postPackageWithin(workspace, link)
	if err != nil {
		t.Fatalf("postPackageWithin: %v", err)
	}
	defer handle.Close()
	if want := "content/blog/real"; got != want {
		t.Errorf("derived package %q, expected the canonical %q", got, want)
	}
}

// TestPostPackageWithinRefusesLinkIntoAnotherCheckout asserts a post directory
// reached through a link into a separate checkout that carries its own
// repository marker is refused. The workspace is the caller's, so no marker in
// the linked tree can be mistaken for it.
func TestPostPackageWithinRefusesLinkIntoAnotherCheckout(t *testing.T) {
	workspace := t.TempDir()
	// A separate checkout with its own repository marker, linked into this one.
	external := t.TempDir()
	if err := os.Mkdir(filepath.Join(external, ".git"), 0o755); err != nil {
		t.Fatalf("write external marker: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(external, "example"), 0o755); err != nil {
		t.Fatalf("create external post: %v", err)
	}
	link := filepath.Join(workspace, "content", "blog", "example")
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatalf("create link parent: %v", err)
	}
	if err := os.Symlink(filepath.Join(external, "example"), link); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	if _, handle, err := postPackageWithin(workspace, link); err == nil {
		handle.Close()
		t.Fatal("expected a link into another checkout to be refused rather than recorded relative to the external tree")
	}
}

// TestPostPackageWithinRefusesIntermediateSymlinkEscape asserts a post
// directory reached through an intermediate symlink into another checkout is
// refused. The post directory itself is a real directory reached through a
// symlinked parent, so a check of its own components alone would miss the
// escape; the descriptor-anchored traversal refuses it at the parent.
func TestPostPackageWithinRefusesIntermediateSymlinkEscape(t *testing.T) {
	workspace := t.TempDir()
	// The intermediate alias points at another checkout, and the post directory
	// beneath it is a real directory that itself carries a repository marker.
	external := t.TempDir()
	if err := os.Mkdir(filepath.Join(external, ".git"), 0o755); err != nil {
		t.Fatalf("write external marker: %v", err)
	}
	postDir := filepath.Join(external, "post")
	if err := os.MkdirAll(filepath.Join(postDir, ".git"), 0o755); err != nil {
		t.Fatalf("create external post marker: %v", err)
	}
	alias := filepath.Join(workspace, "alias")
	if err := os.Symlink(external, alias); err != nil {
		t.Fatalf("create alias symlink: %v", err)
	}

	got, handle, err := postPackageWithin(workspace, filepath.Join(alias, "post"))
	if err == nil {
		handle.Close()
		t.Fatalf("expected a post reached through an intermediate symlink to be refused, derived package %q", got)
	}
}

// TestPostPackageWithinRefusesLinkAtAnotherRoot asserts a link that points
// directly at another checkout's root is refused rather than collapsing the
// derived package to ".".
func TestPostPackageWithinRefusesLinkAtAnotherRoot(t *testing.T) {
	workspace := t.TempDir()
	external := t.TempDir()
	if err := os.Mkdir(filepath.Join(external, ".git"), 0o755); err != nil {
		t.Fatalf("write external marker: %v", err)
	}
	link := filepath.Join(workspace, "content", "blog", "example")
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatalf("create link parent: %v", err)
	}
	if err := os.Symlink(external, link); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	got, handle, err := postPackageWithin(workspace, link)
	if err == nil {
		handle.Close()
		t.Fatalf("expected a link directly at another repository root to be refused, derived package %q", got)
	}
}

// TestPostPackageWithinRefusesDirectoryOutsideWorkspace asserts a post
// directory named outside the workspace is refused before any traversal.
func TestPostPackageWithinRefusesDirectoryOutsideWorkspace(t *testing.T) {
	workspace := t.TempDir()
	outside := filepath.Join(t.TempDir(), "example")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatalf("create outside post: %v", err)
	}

	if _, handle, err := postPackageWithin(workspace, outside); err == nil {
		handle.Close()
		t.Fatal("expected a post directory outside the workspace to be refused")
	}
}

// TestPostPackageWithinHandleSurvivesSymlinkRetarget asserts the handle the
// derivation returns stays anchored to the directory it validated, so a symlink
// retargeted after validation cannot redirect an image read into another tree
// while the artifact still records the validated package.
func TestPostPackageWithinHandleSurvivesSymlinkRetarget(t *testing.T) {
	workspace := t.TempDir()
	validated := filepath.Join(workspace, "content", "blog", "real")
	if err := os.MkdirAll(validated, 0o755); err != nil {
		t.Fatalf("create validated post: %v", err)
	}
	inside := append([]byte{}, tinyPNGForCommandTest...)
	inside = append(inside, 0x01, 0x02, 0x03)
	if err := os.WriteFile(filepath.Join(validated, "diagram.png"), inside, 0o644); err != nil {
		t.Fatalf("write validated image: %v", err)
	}
	link := filepath.Join(workspace, "content", "blog", "alias")
	if err := os.Symlink(validated, link); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	got, handle, err := postPackageWithin(workspace, link)
	if err != nil {
		t.Fatalf("postPackageWithin: %v", err)
	}
	defer handle.Close()
	if want := "content/blog/real"; got != want {
		t.Errorf("derived package %q, expected the canonical %q", got, want)
	}

	// Retarget the symlink after validation, so a read through the pathname
	// would reach the replacement tree while the recorded package still names
	// the validated directory.
	replacement := filepath.Join(t.TempDir(), "real")
	if err := os.MkdirAll(replacement, 0o755); err != nil {
		t.Fatalf("create replacement post: %v", err)
	}
	outside := append([]byte{}, tinyPNGForCommandTest...)
	outside = append(outside, 0xFF, 0xFE)
	if err := os.WriteFile(filepath.Join(replacement, "diagram.png"), outside, 0o644); err != nil {
		t.Fatalf("write replacement image: %v", err)
	}
	if err := os.Remove(link); err != nil {
		t.Fatalf("remove symlink: %v", err)
	}
	if err := os.Symlink(replacement, link); err != nil {
		t.Fatalf("retarget symlink: %v", err)
	}

	read, err := handle.ReadFile("diagram.png")
	if err != nil {
		t.Fatalf("read through validated handle: %v", err)
	}
	if !bytes.Equal(read, inside) {
		t.Errorf("handle read %x, expected the validated package's bytes %x", read, inside)
	}
}

// TestPackageHandleRefusesPackageNamingAnotherDirectory asserts an explicitly
// supplied package is only accepted when it names the post directory the
// locator paths are relative to, because publication resolves those paths
// against the package: a package naming a different directory would make
// publication read a different object than conversion did.
func TestPackageHandleRefusesPackageNamingAnotherDirectory(t *testing.T) {
	workspace := t.TempDir()
	postDir := filepath.Join(workspace, "content", "blog", "example")
	if err := os.MkdirAll(postDir, 0o755); err != nil {
		t.Fatalf("create post: %v", err)
	}
	other := filepath.Join(workspace, "content", "blog", "other")
	if err := os.MkdirAll(other, 0o755); err != nil {
		t.Fatalf("create other: %v", err)
	}

	if _, err := packageHandle(workspace, "content/blog/other", postDir); err == nil {
		t.Fatal("expected a package naming a different directory than the post to be refused")
	}
	handle, err := packageHandle(workspace, "content/blog/example", postDir)
	if err != nil {
		t.Fatalf("packageHandle for the matching package: %v", err)
	}
	defer handle.Close()
}

// TestPackageHandleRefusesEscapingPackage asserts an absolute or escaping
// package is refused rather than opened.
func TestPackageHandleRefusesEscapingPackage(t *testing.T) {
	workspace := t.TempDir()
	postDir := t.TempDir()
	if _, err := packageHandle(workspace, "/etc", postDir); err == nil {
		t.Fatal("expected an absolute package to be refused")
	}
	if _, err := packageHandle(workspace, "content/../../etc", postDir); err == nil {
		t.Fatal("expected an escaping package to be refused")
	}
}

// TestOpenWithinComparesOpenedIdentity asserts the package accepted is the
// directory the handle reads, compared by file identity, and that a package
// naming another directory is refused even when both are inside the workspace.
func TestOpenWithinComparesOpenedIdentity(t *testing.T) {
	workspace := t.TempDir()
	postDir := filepath.Join(workspace, "content", "blog", "example")
	if err := os.MkdirAll(postDir, 0o755); err != nil {
		t.Fatalf("create post: %v", err)
	}
	other := filepath.Join(workspace, "content", "blog", "other")
	if err := os.MkdirAll(other, 0o755); err != nil {
		t.Fatalf("create other: %v", err)
	}
	canonical, err := filepath.EvalSymlinks(workspace)
	if err != nil {
		t.Fatalf("canonicalize workspace: %v", err)
	}

	if _, err := openWithin(canonical, "content/blog/other", postDir); err == nil {
		t.Fatal("expected a package naming a different directory than the post to be refused")
	}
	// A package naming the post directory itself is accepted.
	handle, err := openWithin(canonical, "content/blog/example", postDir)
	if err != nil {
		t.Fatalf("openWithin for the post's own package: %v", err)
	}
	defer handle.Close()
	// A symlinked package is refused by the anchored open, which is the same
	// behavior publication's own anchored resolution has.
	alias := filepath.Join(workspace, "alias")
	if err := os.Symlink(postDir, alias); err != nil {
		t.Fatalf("create symlink: %v", err)
	}
	if _, err := openWithin(canonical, "alias", postDir); err == nil {
		t.Fatal("expected a symlinked package to be refused by the anchored open")
	}
}

// TestConversionReadsSourceThroughValidatedHandle asserts the source is read
// through the validated package handle rather than by its pathname. A source
// that is itself a symlink out of the package is refused instead of read, so
// the article text cannot come from outside the boundary the locators name.
func TestConversionReadsSourceThroughValidatedHandle(t *testing.T) {
	workspace := t.TempDir()
	dir := filepath.Join(workspace, "content", "blog", "example")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create post directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "diagram.png"), tinyPNGForCommandTest, 0o644); err != nil {
		t.Fatalf("write image: %v", err)
	}
	// The source file is a symlink to a post outside the package.
	outside := filepath.Join(t.TempDir(), "index.md")
	if err := os.WriteFile(outside,
		[]byte("---\ntitle: Outside Post\n---\n\n![diagram](./diagram.png)\n"), 0o644); err != nil {
		t.Fatalf("write outside source: %v", err)
	}
	source := filepath.Join(dir, "index.md")
	if err := os.Symlink(outside, source); err != nil {
		t.Fatalf("create source symlink: %v", err)
	}

	out := filepath.Join(t.TempDir(), "artifact.json")
	var stdout, stderr bytes.Buffer
	err := run(source, out, "", workspace, &stdout, &stderr)
	if err == nil {
		t.Fatalf("expected a source symlinked out of the package to be refused, got success (artifact %s)", out)
	}
	if !strings.Contains(err.Error(), "read source") {
		t.Errorf("error did not report the refused source read: %v", err)
	}
}

// TestConversionReadsImagesThroughValidatedHandle reproduces the retarget the
// descriptor-anchored handle exists to defeat: after the package is validated,
// the post's symlink is pointed at a replacement tree whose image differs. A
// read through the source pathname would take the replacement's bytes while the
// locator still records the validated package, producing an artifact
// publication rejects; a read through the handle takes the validated bytes.
func TestConversionReadsImagesThroughValidatedHandle(t *testing.T) {
	workspace := t.TempDir()
	validated := filepath.Join(workspace, "content", "blog", "real")
	if err := os.MkdirAll(validated, 0o755); err != nil {
		t.Fatalf("create validated post: %v", err)
	}
	validatedBytes := tinyPNGForCommandTest
	if err := os.WriteFile(filepath.Join(validated, "diagram.png"), validatedBytes, 0o644); err != nil {
		t.Fatalf("write validated image: %v", err)
	}
	link := filepath.Join(workspace, "content", "blog", "alias")
	if err := os.Symlink(validated, link); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	got, handle, err := postPackageWithin(workspace, link)
	if err != nil {
		t.Fatalf("postPackageWithin: %v", err)
	}
	defer handle.Close()

	replacement := filepath.Join(t.TempDir(), "real")
	if err := os.MkdirAll(replacement, 0o755); err != nil {
		t.Fatalf("create replacement post: %v", err)
	}
	var replacementImage bytes.Buffer
	if err := png.Encode(&replacementImage, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatalf("encode replacement image: %v", err)
	}
	replacementBytes := replacementImage.Bytes()
	if err := os.WriteFile(filepath.Join(replacement, "diagram.png"), replacementBytes, 0o644); err != nil {
		t.Fatalf("write replacement image: %v", err)
	}
	if err := os.Remove(link); err != nil {
		t.Fatalf("remove symlink: %v", err)
	}
	if err := os.Symlink(replacement, link); err != nil {
		t.Fatalf("retarget symlink: %v", err)
	}

	converter := markdown.New(link, got)
	converter.PackageRoot = handle
	converter.PostSource = filepath.Join(link, "index.md")
	body := []byte("---\ntitle: Example\n---\n\n![diagram](./diagram.png)\n")
	article, err := converter.Convert(body)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if len(article.Locators) != 1 {
		t.Fatalf("expected one locator, got %d", len(article.Locators))
	}
	var decoded struct {
		Locators []struct {
			PostPackage string `json:"post_package"`
			Digest      string `json:"digest"`
		} `json:"image_locators"`
	}
	encoded, err := article.JSON()
	if err != nil {
		t.Fatalf("encode artifact: %v", err)
	}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("decode artifact: %v", err)
	}
	if want := "content/blog/real"; decoded.Locators[0].PostPackage != want {
		t.Errorf("locator records package %q, expected the validated %q", decoded.Locators[0].PostPackage, want)
	}
	validatedSum := sha256.Sum256(validatedBytes)
	if want := hex.EncodeToString(validatedSum[:]); decoded.Locators[0].Digest != want {
		t.Errorf("locator digest %s is not the validated image's %s; the read followed the retargeted symlink",
			decoded.Locators[0].Digest, want)
	}
}

// TestRunReportsContinuingDiagnosticsOnFailure asserts a run that both fails and
// loses formatting reports the continuing diagnostics too, because the returned
// error carries only the failing ones.
func TestRunReportsContinuingDiagnosticsOnFailure(t *testing.T) {
	workspace := t.TempDir()
	// The supplied package must name the post directory the locator paths are
	// relative to, so the post lives at its package path under the workspace.
	relPackage := filepath.Join("content", "blog", "example")
	dir := filepath.Join(workspace, relPackage)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create post directory: %v", err)
	}
	// Inline code keeps its text but loses its style (continuing); an SVG image
	// cannot be uploaded (failing).
	if err := os.WriteFile(filepath.Join(dir, "diagram.svg"), []byte("<svg></svg>"), 0o644); err != nil {
		t.Fatalf("write svg: %v", err)
	}
	source := filepath.Join(dir, "index.md")
	body := "---\ntitle: Example\n---\n\nKeep `inline code` and ![diagram](./diagram.svg).\n"
	if err := os.WriteFile(source, []byte(body), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err := run(source, filepath.Join(dir, "artifact.json"), filepath.ToSlash(relPackage), workspace, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected conversion of an unacceptable image to fail")
	}
	if !strings.Contains(stderr.String(), "inline-code-style-lost") {
		t.Errorf("stderr did not report the continuing diagnostic: %q", stderr.String())
	}
	if !strings.Contains(err.Error(), "image-media-type-rejected") {
		t.Errorf("error did not report the failing diagnostic: %v", err)
	}
}
