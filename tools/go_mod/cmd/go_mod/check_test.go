package main

import (
	"os"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

func TestParseVersion(t *testing.T) {
	for _, test := range []struct {
		name string
		data string
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

func TestWorkspaceMatchesConfiguredVersion(t *testing.T) {
	root, err := rootDir(*workspaceMarker)
	if err != nil {
		t.Fatal(err)
	}
	gitPath, err := runfiles.Rlocation(*gitRunfile)
	if err != nil {
		t.Fatal(err)
	}
	if err := check(root, gitPath, *goVersion); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateFileRewritesVersion(t *testing.T) {
	data := []byte("module example.com/root\n\ngo 1.21.0\n")
	if err := updateFile("go.mod", data, "1.26.5"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("go.mod")
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
