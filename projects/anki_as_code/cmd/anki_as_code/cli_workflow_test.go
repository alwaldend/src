package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestConfigExportWorkflow round-trips JSON and TOML through the real file-only CLI.
// Malformed input must fail before writes; successful exports retain unrelated files.
func TestConfigExportWorkflow(t *testing.T) {
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	input := `
{
  "config": {
    "version": "1"
  },
  "decks": [
    {
      "deck": {
        "preset": "Default"
      },
      "title": "Lesson",
      "path": "decks/Lesson/deck.anki.toml",
      "notes": [
        {
          "path": "decks/Lesson/notes/greeting.note.anki.toml",
          "config": {
            "note_type": "note_types/Basic.notetype.anki.toml",
            "fields": {
              "Front": "<b>你好</b>\nSecond \"line\""
            }
          }
        }
      ]
    }
  ],
  "note_types": [
    {
      "path": "note_types/Basic.notetype.anki.toml",
      "config": {
        "name": "Basic",
        "fields": [
          {
            "name": "Front"
          }
        ],
        "templates": [
          {
            "name": "Card",
            "config": {
              "q_format": "{{Front}}",
              "a_format": "{{Front}}"
            }
          }
        ]
      }
    }
  ]
}
`
	root := t.TempDir()
	inputPath := filepath.Join(root, "configs.json")
	if err := os.WriteFile(inputPath, []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}
	var expected any
	if err := json.Unmarshal([]byte(input), &expected); err != nil {
		t.Fatal(err)
	}
	for _, source := range []string{"file", "stdin", "dash", "toml"} {
		t.Run(source, func(t *testing.T) {
			output := filepath.Join(root, source)
			if err := os.Mkdir(output, 0o755); err != nil {
				t.Fatal(err)
			}
			unrelated := filepath.Join(output, "README.txt")
			if err := os.WriteFile(unrelated, []byte("retained\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			arguments := []string{"config", "export", "--output", output}
			data := input
			switch source {
			case "file":
				arguments = append(arguments, inputPath)
			case "dash":
				arguments = append(arguments, "-")
			case "toml":
				dump := exec.CommandContext(t.Context(), binary, "config", "dump", "--text", filepath.Join(root, "file"), "--format", "toml")
				encoded, err := dump.Output()
				if err != nil {
					t.Fatal(err)
				}
				data = string(encoded)
				arguments = append(arguments, "--format", "toml")
			}
			command := exec.CommandContext(t.Context(), binary, arguments...)
			command.Stdin = strings.NewReader(data)
			var stdout, stderr bytes.Buffer
			command.Stdout, command.Stderr = &stdout, &stderr
			if err := command.Run(); err != nil || stdout.Len() != 0 || stderr.Len() == 0 {
				t.Fatalf("config export: error=%v stdout=%q stderr=%q", err, stdout.String(), stderr.String())
			}
			dump := exec.CommandContext(t.Context(), binary, "conf", "dump", "--text", output)
			encoded, err := dump.Output()
			if err != nil {
				t.Fatal(err)
			}
			var actual any
			if err := json.Unmarshal(encoded, &actual); err != nil || !reflect.DeepEqual(expected, actual) {
				t.Fatalf("configuration round-trip: %v\n%s", err, encoded)
			}
			if data, err := os.ReadFile(unrelated); err != nil || string(data) != "retained\n" {
				t.Fatalf("unrelated export file changed: %q %v", data, err)
			}
			if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
				if err := os.WriteFile(filepath.Join(destination, "export-"+source+".json"), encoded, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
	for index, data := range []string{`{`, `{"unknown":true}`, `{"config":{"version":"invalid"}}`} {
		output := filepath.Join(root, "invalid", strconv.Itoa(index))
		command := exec.CommandContext(t.Context(), binary, "config", "export", "--output", output)
		command.Stdin = strings.NewReader(data)
		var stdout, stderr bytes.Buffer
		command.Stdout, command.Stderr = &stdout, &stderr
		if err := command.Run(); err == nil || stdout.Len() != 0 || stderr.Len() == 0 {
			t.Fatalf("malformed config input accepted: error=%v stdout=%q stderr=%q", err, stdout.String(), stderr.String())
		}
		if _, err := os.Stat(output); !os.IsNotExist(err) {
			t.Fatalf("malformed input created output: %v", err)
		}
	}
	if data, err := os.ReadFile(inputPath); err != nil || string(data) != input {
		t.Fatalf("export changed its source: %q %v", data, err)
	}
}

// TestStateSubcommandErrors rejects incomplete and unknown operations before backend access.
func TestStateSubcommandErrors(t *testing.T) {
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	for _, args := range [][]string{
		{"col", "state"},
		{"col", "state", "--collection", "unused.anki2"},
		{"col", "state", "nonsense"},
		{"col", "state", "nonsense", "--collection", "unused.anki2"},
	} {
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		command := exec.CommandContext(ctx, binary, args...)
		var stdout, stderr bytes.Buffer
		command.Stdout, command.Stderr = &stdout, &stderr
		err := command.Run()
		cancel()
		if err == nil || stdout.Len() != 0 || stderr.Len() == 0 {
			t.Fatalf("invalid state command %v: error=%v stdout=%q stderr=%q", args, err, stdout.String(), stderr.String())
		}
	}
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, binary, "col", "state", "--help")
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil || !strings.Contains(stdout.String(), "Available Commands:") || stderr.Len() != 0 {
		t.Fatalf("explicit state help: error=%v stdout=%q stderr=%q", err, stdout.String(), stderr.String())
	}
}

// TestStateRequiredFlags verifies missing state collection inputs fail with clean output streams.
func TestStateRequiredFlags(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	for _, args := range [][]string{
		{"col", "state", "import", "--text", "unused", "notes/example.note.anki.toml=1"},
		{"col", "state", "update", "unused"},
	} {
		var stdout, stderr bytes.Buffer
		command := exec.CommandContext(ctx, binary, args...)
		command.Stdout, command.Stderr = &stdout, &stderr
		if err := command.Run(); err == nil || stdout.Len() != 0 || !strings.Contains(stderr.String(), "required flag") {
			t.Fatalf("required state flags: error=%v stdout=%q stderr=%q", err, stdout.String(), stderr.String())
		}
	}
}

// TestStateCommandSyntax verifies binding errors are reported before collection access.
func TestStateCommandSyntax(t *testing.T) {
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	for _, arguments := range [][]string{
		{"bad"},
		{"=opaque"},
		{"notes/a.note.anki.toml=1", "notes/a.note.anki.toml=2"},
	} {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		args := append([]string{"col", "state", "import", "--collection", "unused.anki2", "--text", "missing"}, arguments...)
		command := exec.CommandContext(ctx, binary, args...)
		var stdout, stderr bytes.Buffer
		command.Stdout, command.Stderr = &stdout, &stderr
		err := command.Run()
		cancel()
		if err == nil || stdout.Len() != 0 || !strings.Contains(stderr.String(), "binding") || strings.Contains(stderr.String(), "unknown command") {
			t.Fatalf("invalid binding %v: error=%v stdout=%q stderr=%q", arguments, err, stdout.String(), stderr.String())
		}
	}
}
