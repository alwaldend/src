package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestCLIReferenceWorkflow exercises the generator executable without a native
// server, checks the help hierarchy, and retains the generated pages as evidence.
func TestCLIReferenceWorkflow(t *testing.T) {
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("CLI_DOCS_GENERATOR"))
	directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if directory == "" {
		directory = t.TempDir()
	}
	output := filepath.Join(directory, "cli")
	blocked := filepath.Join(directory, "blocked")
	if err := os.WriteFile(blocked, []byte("output directory blocker\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, arguments := range [][]string{
		nil,
		{"--output", output, "unexpected"},
		{"--output", blocked},
	} {
		ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
		command := exec.CommandContext(ctx, binary, arguments...)
		if diagnostics, err := command.CombinedOutput(); err == nil {
			t.Fatalf("invalid invocation %q succeeded: %s", arguments, diagnostics)
		}
		cancel()
	}
	var previous map[string][]byte
	for range 2 {
		ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
		command := exec.CommandContext(ctx, binary, "--output", output)
		diagnostics, err := command.CombinedOutput()
		cancel()
		if err != nil {
			t.Fatalf("generate reference: %v\n%s", err, diagnostics)
		}
		if len(diagnostics) != 0 {
			t.Fatalf("unexpected generator output: %s", diagnostics)
		}
		files, err := os.ReadDir(output)
		if err != nil {
			t.Fatal(err)
		}
		pages := map[string][]byte{}
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
				t.Fatalf("unexpected reference entry %q", file.Name())
			}
			page, err := os.ReadFile(filepath.Join(output, file.Name()))
			if err != nil {
				t.Fatal(err)
			}
			pages[file.Name()] = page
			if file.Name() != "_index.md" {
				prefix := "## anki_as_code " + strings.TrimSuffix(file.Name(), ".md")
				for line := range strings.SplitSeq(string(page), "\n") {
					if strings.HasPrefix(line, "## ") && line != prefix && !strings.HasPrefix(line, prefix+" ") {
						t.Errorf("%s contains another command group's help: %s", file.Name(), line)
					}
				}
			}
			if previous != nil && !bytes.Equal(previous[file.Name()], page) {
				t.Errorf("successive generations of %s differ", file.Name())
			}
		}
		index := string(pages["_index.md"])
		if len(index) > 4096 || strings.Contains(index, "```") {
			t.Error("reference index should be a short description and child links")
		}
		for _, name := range []string{"api", "col", "config", "help", "server"} {
			if !strings.Contains(index, "]("+name+"/)") {
				t.Errorf("index omits the %s child link", name)
			}
		}
		for name, expected := range map[string][]string{
			"_index.md": {"CLI reference", "Export and reconcile Anki collections as TOML"},
			"col.md": {
				"## anki_as_code col export\n", "## anki_as_code col state replace\n",
				"col, collection", "--state-key", "--state-output", "-h, --help",
			},
			"config.md": {"## anki_as_code config export\n", "config, conf"},
			"api.md": {
				"## anki_as_code api notes get-note\n", "## anki_as_code api tags all-tags\n",
				"--input-format", "--output-format", "--timeout", "--write-back", "Request fields:", "Response:",
			},
			"server.md": {"## anki_as_code server run\n", "## anki_as_code server serve\n", "anki_as_code server serve --help"},
			"help.md":   {"## anki_as_code help\n"},
		} {
			for _, content := range expected {
				if !strings.Contains(string(pages[name]), content) {
					t.Errorf("%s missing %q", name, content)
				}
			}
		}
		export := strings.SplitN(string(pages["col.md"]), "## anki_as_code col export\n", 2)
		if len(export) != 2 || !strings.Contains(strings.SplitN(export[1], "\n## ", 2)[0], "--state-key") {
			t.Error("collection export help omits the inherited namespace flag")
		}
		if previous != nil && len(previous) != len(pages) {
			t.Error("successive generations have different page inventories")
		}
		previous = pages
	}
}
