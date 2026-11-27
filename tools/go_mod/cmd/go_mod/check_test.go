package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// goVersion is the expected directive supplied by the Bazel test target.
var (
	goVersion = flag.String("version", "", "Expected Go version")
	// gitRunfile identifies the configured Git executable.
	gitRunfile = flag.String("git", "", "Git executable runfile")
	// workspaceMarker identifies a tracked file at the repository root.
	workspaceMarker = flag.String("workspace-marker", "", "Workspace module runfile")
)

// TestParseVersion checks directives with and without leading comments.
func TestParseVersion(t *testing.T) {
	for _, test := range []struct {
		// name identifies the fixture in test output.
		name string
		// data contains the module text to parse.
		data string
		// want is the expected directive value.
		want string
	}{
		{"module header", "module example.com/root\n\ngo 1.26.5\n", "1.26.5"},
		{"leading comments", "// leading comment\n\ngo 1.24.0\n", "1.24.0"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, hasDirective, err := parseVersion("go.mod", []byte(test.data))
			if err != nil {
				t.Fatal(err)
			}
			if !hasDirective {
				t.Fatal("expected a go directive")
			}
			if got != test.want {
				t.Fatalf("parseVersion() = %q, want %q", got, test.want)
			}
		})
	}
}

// TestWorkspaceMatchesConfiguredVersion checks each tracked module against the build setting.
func TestWorkspaceMatchesConfiguredVersion(t *testing.T) {
	root, err := rootDir(*workspaceMarker)
	if err != nil {
		t.Fatal(err)
	}
	gitPath, err := runfiles.Rlocation(*gitRunfile)
	if err != nil {
		t.Fatal(err)
	}
	if err := check(t.Context(), root, gitPath, *goVersion); err != nil {
		t.Fatal(err)
	}
}

// TestUpdateFileRewritesVersion verifies the persisted Go directive.
func TestUpdateFileRewritesVersion(t *testing.T) {
	data := []byte("module example.com/root\n\ngo 1.21.0\n")
	path := filepath.Join(t.TempDir(), "go.mod")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := updateFile(path, data, "1.26.5"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got, hasDirective, err := parseVersion("go.mod", data)
	if err != nil {
		t.Fatal(err)
	}
	if !hasDirective {
		t.Fatal("expected a go directive")
	}
	if got != "1.26.5" {
		t.Fatalf("updated version = %q, want %q", got, "1.26.5")
	}
}
