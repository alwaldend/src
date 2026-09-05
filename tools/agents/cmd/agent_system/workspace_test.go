package agent_system

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

func TestContextResolvesRelativeWorkspaceAgainstBazelWorkspace(t *testing.T) {
	root := t.TempDir()
	t.Setenv("BUILD_WORKSPACE_DIRECTORY", root)
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("workspace policy\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, flags := range [][]string{nil, {"--workspace-root", "."}} {
		var output bytes.Buffer
		if err := Run(flags, &output); err != nil {
			t.Fatal(err)
		}
		capsule, err := v1alpha1.DecodeContextCapsule(output.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if capsule.Identity.WorkspaceRoot != root || len(capsule.Documents) == 0 || capsule.Documents[0].Path != "AGENTS.md" {
			t.Fatalf("read runfiles instead of source workspace: %+v", capsule.Identity)
		}
	}
}

func TestAggregateRelativePathsUseBazelWorkspace(t *testing.T) {
	root := t.TempDir()
	t.Setenv("BUILD_WORKSPACE_DIRECTORY", root)
	aggregateFixtureCatalog(t, root)
	aggregateFixtureCases(t, root)
	var output bytes.Buffer
	if err := Run([]string{"aggregate", "--catalog", "capability.json", "--markdown", "coverage.md"}, &output); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"coverage-matrix.json", "coverage.md"} {
		if _, err := os.Stat(filepath.Join(root, name)); err != nil {
			t.Fatalf("output did not reach source workspace: %v", err)
		}
	}
}

func TestRelativeFileCommandsHonorExplicitWorkspaceAndAbsolutePaths(t *testing.T) {
	root := t.TempDir()
	t.Setenv("BUILD_WORKSPACE_DIRECTORY", filepath.Join(root, "different-run-workspace"))
	for _, command := range []string{"normalize", "coverage", "aggregate"} {
		external := filepath.Join(root, "absolute-output.json")
		args := []string{"--workspace-root", root, "--input", "input.json", "--output", external}
		if command != "normalize" {
			args = append(args, "--catalog", "catalog.json")
		}
		var input, catalog, output string
		switch command {
		case "normalize":
			opts, err := parseNormalizeFlags(args)
			if err != nil {
				t.Fatal(err)
			}
			input, output = opts.input, opts.output
		case "coverage":
			opts, err := parseCoverageFlags(args)
			if err != nil {
				t.Fatal(err)
			}
			input, output, catalog = opts.input, opts.output, opts.catalog
		case "aggregate":
			opts, err := parseAggregateFlags(args)
			if err != nil {
				t.Fatal(err)
			}
			input, output, catalog = opts.input, opts.output, opts.catalog
		}
		if input != filepath.Join(root, "input.json") || output != external || catalog != "" && catalog != filepath.Join(root, "catalog.json") {
			t.Fatalf("%s relative/absolute resolution: input %s output %s catalog %s", command, input, output, catalog)
		}
	}
}

func TestCommandWorkspaceOutsideBazelUsesWorkingDirectory(t *testing.T) {
	t.Setenv("BUILD_WORKSPACE_DIRECTORY", "")
	want, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	got, err := commandWorkspaceRoot("")
	if err != nil || got != want {
		t.Fatalf("root = %s, error %v; want %s", got, err, want)
	}
}
