// Package mmdc_test asserts the contract of the interactive Mermaid CLI: a
// run applies the checked-in repository theme, and it measures label text
// through the repository's pinned fonts rather than through whatever
// Fontconfig the host provides.
//
// The font check is the load-bearing one. Mermaid asks Chrome for text
// metrics, so a CLI run that followed the host would render the same diagram
// with different geometry on every machine.
package mmdc_test

import (
	"bytes"
	"context"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/mermaid/test/svgdoc"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

var (
	mmdc        = flag.String("mmdc", "", "runfile path of the mmdc launcher")
	input       = flag.String("input", "", "runfile path of the diagram source")
	nestedInput = flag.String("nested-input", "", "runfile path of the nested-container diagram source")
	monoFont    = flag.String("mono-anchor", "", "runfile path of the pinned monospace font")
)

// runfilesRoot returns the root the launcher expects to start in. The
// launcher applies the repository theme through a workspace-relative path, so
// it must run from the workspace root of the runfiles tree, as it does under
// `bazel run`. A go_test starts in its own package directory instead, which
// would leave that path unresolved.
func runfilesRoot(t *testing.T) string {
	t.Helper()
	sourceDirectory := os.Getenv("TEST_SRCDIR")
	workspace := os.Getenv("TEST_WORKSPACE")
	if sourceDirectory == "" || workspace == "" {
		t.Fatal("TEST_SRCDIR and TEST_WORKSPACE must be set")
	}
	root := filepath.Join(sourceDirectory, workspace)
	if _, err := os.Stat(filepath.Join(root, "tools", "mermaid", "theme.json")); err != nil {
		t.Fatalf("runfiles workspace root does not hold the theme: %v", err)
	}
	return root
}

// render runs the CLI for the named fixture diagram and returns the SVG it
// wrote.
func renderFixture(t *testing.T, sourcePath, output string, environment []string) []byte {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	launcher, err := runfiles.Rlocation(*mmdc)
	if err != nil {
		t.Fatalf("resolve mmdc: %v", err)
	}
	source, err := runfiles.Rlocation(sourcePath)
	if err != nil {
		t.Fatalf("resolve input: %v", err)
	}

	command := exec.CommandContext(ctx, launcher, "-i", source, "-o", output)
	command.Dir = runfilesRoot(t)
	command.Env = environment
	combined, err := command.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatalf("mmdc exceeded its deadline: %v\n%s", ctx.Err(), combined)
	}
	if err != nil {
		t.Fatalf("mmdc failed: %v\n%s", err, combined)
	}

	svg, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read rendered SVG: %v", err)
	}
	if len(svg) == 0 {
		t.Fatal("mmdc wrote an empty SVG")
	}
	return svg
}

// render runs the CLI for the default fixture diagram and returns the SVG it
// wrote.
func render(t *testing.T, output string, environment []string) []byte {
	t.Helper()
	return renderFixture(t, *input, output, environment)
}

func TestCLIAppliesRepositoryTheme(t *testing.T) {
	svg := string(render(t, filepath.Join(t.TempDir(), "smoke.svg"), os.Environ()))

	if !strings.Contains(svg, "<svg") {
		t.Fatal("mmdc did not write an SVG document")
	}
	for _, expected := range []string{"fill:#ffffff", "stroke:#1e1e1e", "rough-node"} {
		if !strings.Contains(svg, expected) {
			t.Errorf("CLI output does not contain %q", expected)
		}
	}
	for _, forbidden := range []string{"#ECECFF", "#9370DB"} {
		if strings.Contains(svg, forbidden) {
			t.Errorf("CLI output still contains the Mermaid default theme color %s", forbidden)
		}
	}
}

// TestCLIAppliesPaintOrder asserts that an interactive render matches a build
// action on the one property the action fixes after Mermaid runs. Mermaid emits
// a container title ahead of the edges, so an edge crossing a container border
// would draw over the title plate; the CLI has to lift the plates too, or the
// diagram a caller iterates on is not the diagram the build commits.
func TestCLIAppliesPaintOrder(t *testing.T) {
	scratch := t.TempDir()
	path := filepath.Join(scratch, "nested.svg")
	svg := string(renderFixture(t, *nestedInput, path, os.Environ()))

	if !svgdoc.ClusterTitlesPaintLast(svg) {
		t.Error("CLI output does not paint container titles after the node layer, " +
			"so an interactive render differs from a build action")
	}
}

// TestCLIIgnoresHostFontconfig pins the font isolation of the CLI. It points
// the ambient Fontconfig at a monospace-only pool, so labels measured through
// those metrics would render wider and change the canvas. A byte-identical
// result shows the run resolved text through the pinned fonts instead.
func TestCLIIgnoresHostFontconfig(t *testing.T) {
	scratch := t.TempDir()
	pristine := render(t, filepath.Join(scratch, "pristine.svg"), os.Environ())

	anchor, err := runfiles.Rlocation(*monoFont)
	if err != nil {
		t.Fatalf("resolve monospace anchor: %v", err)
	}
	poisonedLibrary := filepath.Join(scratch, "poison", "lib")
	if err := os.MkdirAll(poisonedLibrary, 0o755); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(anchor)
	if err != nil {
		t.Fatalf("read monospace anchor: %v", err)
	}
	if err := os.WriteFile(filepath.Join(poisonedLibrary, filepath.Base(anchor)), content, 0o644); err != nil {
		t.Fatal(err)
	}

	// Guard against a vacuous check: the pool must exclude the proportional
	// font, so a run that consulted it could not produce the same metrics.
	entries, err := os.ReadDir(poisonedLibrary)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "LiberationSans") {
			t.Fatalf("poisoned font pool unexpectedly contains a proportional font: %s", entry.Name())
		}
	}

	poisonedDirectory := filepath.Dir(poisonedLibrary)
	configFile := filepath.Join(poisonedDirectory, "fonts.conf")
	config := `<?xml version="1.0"?>
<fontconfig>
  <dir>` + poisonedLibrary + `</dir>
  <cachedir>` + filepath.Join(poisonedDirectory, "cache") + `</cachedir>
  <alias><family>Helvetica</family><prefer><family>Liberation Mono</family></prefer></alias>
  <alias><family>Arial</family><prefer><family>Liberation Mono</family></prefer></alias>
  <alias><family>sans-serif</family><prefer><family>Liberation Mono</family></prefer></alias>
</fontconfig>
`
	if err := os.WriteFile(configFile, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}

	environment := append(os.Environ(),
		"FONTCONFIG_FILE="+configFile,
		"FONTCONFIG_PATH="+poisonedDirectory,
	)
	poisoned := render(t, filepath.Join(scratch, "poisoned.svg"), environment)

	if !bytes.Equal(pristine, poisoned) {
		t.Errorf("host Fontconfig changed the CLI canvas; rendering is not hermetic:\npristine %d bytes, poisoned %d bytes",
			len(pristine), len(poisoned))
	}
}
