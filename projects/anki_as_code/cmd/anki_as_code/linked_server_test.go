package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// installedLinkedCLI copies only the application, without runfiles or a companion worker.
func installedLinkedCLI(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	source := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	data, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(root, "anki_as_code")
	if err := os.WriteFile(binary, data, 0o700); err != nil {
		t.Fatal(err)
	}
	return binary, root
}

// linkedCommand runs an installed CLI with a bounded lifetime and separated streams.
func linkedCommand(t *testing.T, binary string, success bool, arguments ...string) []byte {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, binary, arguments...)
	command.Stdin = strings.NewReader("")
	var output, diagnostics bytes.Buffer
	command.Stdout, command.Stderr = &output, &diagnostics
	err := command.Run()
	if ctx.Err() != nil {
		t.Fatalf("%v exceeded its deadline: %v", arguments, ctx.Err())
	}
	if success && err != nil || !success && err == nil {
		t.Fatalf("%v: %v; stderr=%s", arguments, err, diagnostics.String())
	}
	if !success && diagnostics.Len() == 0 {
		t.Fatal("failed server did not report diagnostics")
	}
	return output.Bytes()
}

// TestLinkedWorker verifies one installed executable works without PATH or temporary storage.
func TestLinkedWorker(t *testing.T) {
	binary, root := installedLinkedCLI(t)
	t.Setenv("PATH", "")
	t.Setenv("TMPDIR", filepath.Join(root, "unavailable"))
	for _, arguments := range [][]string{{"server"}, {"server", "--help"}} {
		help := linkedCommand(t, binary, true, arguments...)
		if !strings.Contains(string(help), "Available Commands:") || !strings.Contains(string(help), "run") || !strings.Contains(string(help), "serve") {
			t.Fatalf("server group did not list run and serve: %s", help)
		}
	}
	for _, arguments := range [][]string{{"server", "run"}, {"server", "run", "--help"}} {
		help := linkedCommand(t, binary, true, arguments...)
		if !strings.Contains(string(help), "Usage: anki_as_code server run <COMMAND>") || !strings.Contains(string(help), "serve") {
			t.Fatalf("native CLI did not list serve: %s", help)
		}
	}
	for _, arguments := range [][]string{{"server", "run", "serve", "--help"}, {"server", "serve", "--help"}} {
		help := linkedCommand(t, binary, true, arguments...)
		if !strings.Contains(string(help), "Usage: anki_as_code server run serve [--transport") {
			t.Fatalf("serve did not delegate native help: %s", help)
		}
	}
	for _, arguments := range [][]string{{"server", "unknown"}, {"server", "--transport", "stdio"}, {"server", "run", "unknown"}, {"server", "run", "--transport", "stdio"}, {"server", "run", "--engine=cgo", "--help"}, {"server", "serve", "--engine=cgo", "--help"}} {
		if output := linkedCommand(t, binary, false, arguments...); len(output) != 0 {
			t.Fatalf("invalid command arguments wrote %d protocol bytes", len(output))
		}
	}
	for _, transport := range []string{"unknown", "tcp"} {
		for _, arguments := range [][]string{{"server", "run", "serve"}, {"server", "serve"}} {
			arguments = append(arguments, "--transport", transport)
			if output := linkedCommand(t, binary, false, arguments...); len(output) != 0 {
				t.Fatalf("invalid serve arguments wrote %d protocol bytes", len(output))
			}
		}
	}
	// An incomplete handshake fails after HTTP/2 may have emitted its settings.
	linkedCommand(t, binary, false, "server", "run", "serve", "--transport", "stdio")
	database := filepath.Join(root, "linked.anki2")
	if err := os.WriteFile(database, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	response := linkedCommand(t, binary, true, "col", "state", "dump", "--collection", database, "--format", "json")
	var state map[string]any
	if err := json.Unmarshal(response, &state); err != nil {
		t.Fatal(err)
	}
	if state["version"] != "1" {
		t.Fatalf("unexpected native state: %s", response)
	}
	if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
		data, err := os.ReadFile(database)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, "linked.anki2"), data, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, "linked-state.json"), response, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestLinkedServerInterrupt verifies both linked listeners terminate when interrupted.
func TestLinkedServerInterrupt(t *testing.T) {
	for _, mode := range []string{"run", "serve"} {
		for _, transport := range []string{"tcp", "unix"} {
			t.Run(mode+"/"+transport, func(t *testing.T) {
				binary, root := installedLinkedCLI(t)
				t.Setenv("PATH", "")
				t.Setenv("TMPDIR", filepath.Join(root, "unavailable"))
				address := "127.0.0.1:0"
				if transport == "unix" {
					address = "anki.sock"
				}
				ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
				defer cancel()
				arguments := []string{"server", mode}
				if mode == "run" {
					arguments = append(arguments, "serve")
				}
				arguments = append(arguments, "--transport", transport, "--listen", address)
				command := exec.CommandContext(ctx, binary, arguments...)
				command.Dir = root
				stderr, err := command.StderrPipe()
				if err != nil {
					t.Fatal(err)
				}
				if err := command.Start(); err != nil {
					t.Fatal(err)
				}
				done := make(chan error, 1)
				go func() { done <- command.Wait() }()
				t.Cleanup(func() { _ = command.Process.Kill() })
				ready := make(chan string, 1)
				go func() {
					scanner := bufio.NewScanner(stderr)
					for scanner.Scan() {
						if strings.HasPrefix(scanner.Text(), "Listening on ") {
							ready <- scanner.Text()
							return
						}
					}
					ready <- "server exited before readiness"
				}()
				select {
				case line := <-ready:
					if !strings.HasPrefix(line, "Listening on "+transport+"://") {
						t.Fatal(line)
					}
				case <-ctx.Done():
					t.Fatal("server readiness timed out")
				}
				if err := command.Process.Signal(os.Interrupt); err != nil {
					t.Fatal(err)
				}
				select {
				case err := <-done:
					if err == nil {
						t.Fatal("interrupt did not produce an interrupted exit")
					}
				case <-time.After(3 * time.Second):
					t.Fatal("linked server ignored interrupt")
				}
			})
		}
	}
}
