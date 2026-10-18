package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// apiCLI runs the installed application and retains its separated response and diagnostics.
func apiCLI(t *testing.T, binary, input string, success bool, arguments ...string) []byte {
	t.Helper()
	command := exec.CommandContext(t.Context(), binary, arguments...)
	command.Stdin = strings.NewReader(input)
	command.Env = append(os.Environ(), "PATH=")
	var output, diagnostics bytes.Buffer
	command.Stdout, command.Stderr = &output, &diagnostics
	err := command.Run()
	if success && err != nil || !success && err == nil {
		t.Fatalf("CLI %v: %v; stderr=%s stdout=%s", arguments, err, diagnostics.String(), output.String())
	}
	if !success && (output.Len() != 0 || diagnostics.Len() == 0) {
		t.Fatalf("failed CLI polluted response: stdout=%s stderr=%s", output.String(), diagnostics.String())
	}
	return output.Bytes()
}

// installedAPI copies only the application executable outside Bazel's output and runfiles trees.
func installedAPI(t *testing.T) string {
	t.Helper()
	source := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	data, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(t.TempDir(), "anki_as_code")
	if err := os.WriteFile(binary, data, 0o700); err != nil {
		t.Fatal(err)
	}
	return binary
}

// TestAPIWorkflow exercises schema discovery, installed startup, formats and archive publication.
func TestAPIWorkflow(t *testing.T) {
	binary := installedAPI(t)
	for _, service := range []string{"tags", "backend-tags", "i18n", "ankidroid", "anki-hub", "ankiweb", "frontend", "github", "image-occlusion"} {
		apiCLI(t, binary, "", true, "api", service, "--help")
	}
	root := t.TempDir()
	scratch := filepath.Join(root, "temporary")
	if err := os.Mkdir(scratch, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TMPDIR", scratch)
	database := filepath.Join(root, "collection.anki2")
	if err := os.WriteFile(database, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, arguments := range [][]string{
		{"api", "tags", "all-tags"},
		{"api", "tags", "all-tags", "--collection", database, "--backend", "unix:///unused"},
		{"api", "tags", "all-tags", "--collection", database, "--input-format", "yaml"},
		{"api", "backend-collection", "open-collection", "--collection", database},
	} {
		apiCLI(t, binary, "{}", false, arguments...)
	}
	apiCLI(t, binary, `{"unknown":true}`, false, "api", "tags", "all-tags", "--collection", database)
	if info, err := os.Stat(database); err != nil || info.Size() != 0 {
		t.Fatalf("invalid requests opened collection: %v %v", info, err)
	}
	apiCLI(t, binary, "{}", true, "api", "backend-tags", "all-tags", "--collection", database)
	deck := apiCLI(t, binary, "{}", true, "api", "decks", "new-deck", "--collection", database)
	var request map[string]any
	if err := json.Unmarshal(deck, &request); err != nil {
		t.Fatal(err)
	}
	request["name"] = "CLI native deck"
	data, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	apiCLI(t, binary, string(data), true, "api", "decks", "add-deck", "--collection", database)
	result := apiCLI(t, binary, "val = 'CLI native deck'\n", true, "api", "decks", "get-deck-id-by-name", "--collection", database, "--input-format", "toml", "--output-format", "toml")
	if !strings.Contains(string(result), "did = \"") {
		t.Fatalf("TOML response lost protobuf integer encoding: %s", result)
	}
	requestFile := filepath.Join(root, "request.json")
	if err := os.WriteFile(requestFile, []byte(`{"val":"CLI native deck"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	apiCLI(t, binary, "", true, "api", "decks", "get-deck-id-by-name", "--collection", database, "--input", requestFile)
	archive := filepath.Join(root, "collection.colpkg")
	apiCLI(t, binary, fmt.Sprintf(`{"out_path":%q,"include_media":true}`, archive), true, "api", "backend-import-export", "export-collection-package", "--collection", database)
	original, err := os.ReadFile(archive)
	if err != nil {
		t.Fatal(err)
	}
	for index, flags := range [][]string{{"--write-back", "--no-write-back"}, {"--no-write-back", "--write-back=false"}} {
		request["name"] = fmt.Sprintf("Discarded %d", index)
		data, err := json.Marshal(request)
		if err != nil {
			t.Fatal(err)
		}
		arguments := append([]string{"api", "decks", "add-deck", "--collection", archive}, flags...)
		apiCLI(t, binary, string(data), true, arguments...)
		current, err := os.ReadFile(archive)
		if err != nil || !bytes.Equal(original, current) {
			t.Fatalf("disabled publication changed source: %v", err)
		}
	}
	request["name"] = "Persisted CLI deck"
	data, err = json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	apiCLI(t, binary, string(data), true, "api", "decks", "add-deck", "--collection", archive, "--no-write-back", "--write-back")
	observed := apiCLI(t, binary, `{"val":"Persisted CLI deck"}`, true, "api", "decks", "get-deck-id-by-name", "--collection", archive)
	apiCLI(t, binary, `{"nid":"9223372036854775806"}`, false, "api", "notes", "get-note", "--collection", archive)
	apiCLI(t, binary, `{"val":"Persisted CLI deck"}`, true, "api", "decks", "get-deck-id-by-name", "--collection", archive)
	if entries, err := os.ReadDir(scratch); err != nil || len(entries) != 0 {
		t.Fatalf("completed CLI retained temporary worker or collection files: %v %v", entries, err)
	}
	t.Run("help-without-temporary-storage", func(t *testing.T) {
		t.Setenv("TMPDIR", filepath.Join(root, "unavailable"))
		apiCLI(t, binary, "", true, "api", "tags", "all-tags", "--help")
	})
	if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
		data, err := os.ReadFile(archive)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, "api-collection.colpkg"), data, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, "api-response.toml"), result, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, "api-response.json"), observed, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestAPINetworkWorkflow verifies client disconnection preserves remote process and collection lifetime.
func TestAPINetworkWorkflow(t *testing.T) {
	binary := installedAPI(t)
	for _, transport := range []string{"tcp", "unix"} {
		t.Run(transport, func(t *testing.T) {
			root := t.TempDir()
			// Relative socket addresses avoid Unix sockaddr limits in long sandbox paths.
			t.Chdir(root)
			address := "127.0.0.1:0"
			if transport == "unix" {
				address = "anki.sock"
			}
			worker := exec.CommandContext(t.Context(), binary, "server", "serve", "--transport", transport, "--listen", address)
			worker.Dir = root
			stderr, err := worker.StderrPipe()
			if err != nil {
				t.Fatal(err)
			}
			if err := worker.Start(); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				_ = worker.Process.Kill()
				_ = worker.Wait()
			})
			ready := make(chan string, 1)
			go func() {
				scanner := bufio.NewScanner(stderr)
				var diagnostics strings.Builder
				for scanner.Scan() {
					line := scanner.Text()
					if strings.HasPrefix(line, "Listening on ") {
						ready <- line
						for scanner.Scan() {
						}
						return
					}
					diagnostics.WriteString(line + "\n")
				}
				ready <- fmt.Sprintf("server exited before readiness: %s (%v)", diagnostics.String(), scanner.Err())
			}()
			var endpoint string
			select {
			case line := <-ready:
				var ok bool
				endpoint, ok = strings.CutPrefix(line, "Listening on ")
				if !ok {
					t.Fatal(line)
				}
			case <-time.After(15 * time.Second):
				t.Fatal("server readiness timed out")
			}
			request := fmt.Sprintf(`{"collection_path":%q,"media_folder_path":%q,"media_db_path":%q}`, filepath.Join(root, "remote.anki2"), filepath.Join(root, "remote.media"), filepath.Join(root, "remote.media.db2"))
			t.Setenv("TMPDIR", filepath.Join(root, "unavailable"))
			apiCLI(t, binary, request, true, "api", "backend-collection", "open-collection", "--backend", endpoint)
			for range 2 {
				apiCLI(t, binary, "{}", true, "api", "tags", "all-tags", "--backend", endpoint)
			}
			formatted := apiCLI(t, binary, `{"seconds":65,"context":"PRECISE"}`, true, "api", "i18n", "format-timespan", "--backend", endpoint)
			if !strings.Contains(string(formatted), "val") {
				t.Fatalf("i18n response is empty: %s", formatted)
			}
			apiCLI(t, binary, "{}", true, "api", "backend-collection", "close-collection", "--backend", endpoint)
		})
	}
}

// TestAPIInputCancellation bounds incomplete stdin before any native collection startup.
func TestAPIInputCancellation(t *testing.T) {
	binary := installedAPI(t)
	for _, mode := range []string{"timeout", "interrupt"} {
		t.Run(mode, func(t *testing.T) {
			database := filepath.Join(t.TempDir(), "unopened.anki2")
			if err := os.WriteFile(database, nil, 0o600); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			arguments := []string{"api", "tags", "all-tags", "--collection", database, "--input", "-"}
			if mode == "timeout" {
				arguments = append(arguments, "--timeout", "200ms")
			}
			command := exec.CommandContext(ctx, binary, arguments...)
			input, err := command.StdinPipe()
			if err != nil {
				t.Fatal(err)
			}
			defer input.Close()
			var output, diagnostics bytes.Buffer
			command.Stdout, command.Stderr = &output, &diagnostics
			if err := command.Start(); err != nil {
				t.Fatal(err)
			}
			if _, err := input.Write([]byte("{}")); err != nil {
				t.Fatal(err)
			}
			if mode == "interrupt" {
				time.Sleep(300 * time.Millisecond)
				if err := command.Process.Signal(os.Interrupt); err != nil {
					t.Fatal(err)
				}
			}
			err = command.Wait()
			if ctx.Err() != nil {
				t.Fatalf("stalled input ignored %s: %s", mode, diagnostics.String())
			}
			expected := "context deadline exceeded"
			if mode == "interrupt" {
				expected = "context canceled"
			}
			if err == nil || output.Len() != 0 || !strings.Contains(diagnostics.String(), expected) {
				t.Fatalf("cancellation result: %v; stdout=%s stderr=%s", err, output.String(), diagnostics.String())
			}
			if info, err := os.Stat(database); err != nil || info.Size() != 0 {
				t.Fatalf("incomplete request opened collection: %v %v", info, err)
			}
			if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
				if err := os.WriteFile(filepath.Join(destination, "api-input-"+mode+".txt"), diagnostics.Bytes(), 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
