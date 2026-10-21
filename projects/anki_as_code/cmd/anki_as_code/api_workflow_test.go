package main

import (
	"bufio"
	"bytes"
	"context"
	"debug/elf"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/pkg/ankiascode"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
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

// registerLocalAPIClient registers the native fixture once across repeated test runs.
var registerLocalAPIClient sync.Once

// TestLocalAPIAddresses verifies local strategy selection for managed and backend-only calls.
func TestLocalAPIAddresses(t *testing.T) {
	binary := installedAPI(t)
	root := t.TempDir()
	t.Setenv("TMPDIR", root)
	database := filepath.Join(root, "inprocess.anki2")
	if err := os.WriteFile(database, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	subprocess := ankigrpc.SubprocessAddress(binary, "server", "serve")
	addresses := []string{subprocess}
	unavailable := []string{"inprocess+://invalid"}
	if os.Getenv("ANKI_TEST_WASM") == "true" {
		addresses = append(addresses, "inprocess+wasm://local", "inprocess+wasm://another-instance")
		addresses = append(addresses, strings.Replace(subprocess, "subprocess://", "subprocess+wasm://", 1))
	} else {
		unavailable = append(unavailable, "inprocess+wasm://unavailable")
		unavailable = append(unavailable, "subprocess+wasm://unavailable")
	}
	if os.Getenv("ANKI_TEST_CGO") == "true" {
		addresses = append(addresses, "inprocess://local", "inprocess://any-other-name", "inprocess+cgo://local")
		addresses = append(addresses, strings.Replace(subprocess, "subprocess://", "subprocess+cgo://", 1))
	} else {
		unavailable = append(unavailable, "inprocess+cgo://unavailable", "subprocess+cgo://unavailable")
		if os.Getenv("ANKI_TEST_WASM") == "true" {
			addresses = append(addresses, "inprocess://local")
		} else {
			unavailable = append(unavailable, "inprocess://unavailable")
		}
	}
	for _, address := range addresses {
		apiCLI(t, binary, "{}", true, "api", "tags", "all-tags", "--collection", database, "--address", address)
		formatted := apiCLI(t, binary, "{\"seconds\":65,\"context\":\"PRECISE\"}", true, "api", "i18n", "format-timespan", "--address", address)
		if !strings.Contains(string(formatted), "minute") {
			t.Fatalf("direct backend formatting failed: %s", formatted)
		}
		apiCLI(t, binary, "{}", false, "api", "backend-collection", "close-collection", "--collection", database, "--address", address)
		apiCLI(t, binary, "{}", false, "api", "tags", "all-tags", "--address", address)
		usesWASM := strings.HasPrefix(address, "subprocess+wasm://") || strings.HasPrefix(address, "inprocess+wasm://") || os.Getenv("ANKI_TEST_CGO") != "true" && strings.HasPrefix(address, "subprocess://")
		if os.Getenv("ANKI_TEST_WASM") == "true" && usesWASM {
			command := exec.CommandContext(t.Context(), binary, "api", "backend-sync", "sync-collection", "--address", address)
			command.Stdin = strings.NewReader("{}")
			output, err := command.CombinedOutput()
			if err == nil || !strings.Contains(string(output), "unsupported") {
				t.Fatalf("configured worker did not use offline WASM: %v; %s", err, output)
			}
		}
	}
	for _, address := range unavailable {
		command := exec.CommandContext(t.Context(), binary, "api", "i18n", "format-timespan", "--address", address)
		command.Stdin = strings.NewReader("{}")
		output, err := command.CombinedOutput()
		if err == nil {
			t.Fatal("unavailable engine was accepted")
		}
		if !strings.Contains(string(output), "Anki client scheme") || !strings.Contains(string(output), "unavailable") {
			t.Fatalf("unavailable client scheme error is unclear: %s", output)
		}
	}
	registerLocalAPIClient.Do(func() {
		ankigrpc.RegisterClientFactory("native-api-fixture", ankigrpc.ClientFactoryFunc[ankiclient.Client](func(string) (ankiclient.Client, error) {
			client, err := ankigrpc.NewClient("inprocess://registered-api")
			if err != nil {
				return nil, fmt.Errorf("construct registered API fixture: %w", err)
			}
			return client, nil
		}))
	})
	registered, err := ankiascode.New(ankiascode.Options{}).APICall(t.Context(), ankiascodeapi.APIRequest{
		Collection: ankiascodeapi.CollectionOptions{Path: database, Address: "native-api-fixture://managed"},
		Method:     "/anki.tags.TagsService/AllTags", Input: []byte("{}"),
	})
	if err != nil || !json.Valid(registered) {
		t.Fatalf("registered native scheme could not prepare an API collection: %s %v", registered, err)
	}
	for _, address := range []string{"private-api-fixture://managed", "tcp://injected-owned-backend"} {
		calls := 0
		application := ankiascode.New(ankiascode.Options{NewClient: ankigrpc.ClientFactoryFunc[ankiclient.Client](func(input string) (ankiclient.Client, error) {
			calls++
			if input != address {
				return nil, fmt.Errorf("injected API address changed: %q", input)
			}
			client, err := ankigrpc.NewClient("inprocess://injected-api")
			if err != nil {
				return nil, fmt.Errorf("construct injected API fixture: %w", err)
			}
			return client, nil
		})})
		output, err := application.APICall(t.Context(), ankiascodeapi.APIRequest{
			Collection: ankiascodeapi.CollectionOptions{Path: database, Address: address},
			Method:     "/anki.tags.TagsService/AllTags", Input: []byte("{}"),
		})
		if err != nil || !json.Valid(output) || calls != 1 {
			t.Errorf("injected managed API bypassed its factory: address=%s calls=%d output=%s error=%v", address, calls, output, err)
		}
	}
	if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
		data, err := os.ReadFile(database)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, "registered-local-api.json"), registered, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, "inprocess-api.anki2"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestExplicitWorkerEngines verifies qualified schemes select the named compiled engine.
// An invalid WASM cache distinguishes a real cgo worker from an accidental WASM fallback.
func TestExplicitWorkerEngines(t *testing.T) {
	binary := installedAPI(t)
	blocked := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(blocked, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ANKI_WASM_CACHE_DIR", blocked)
	var report strings.Builder
	alias := filepath.Join(t.TempDir(), "application-worker")
	if err := os.Symlink(binary, alias); err != nil {
		t.Fatal(err)
	}
	for _, engine := range []string{"cgo", "wasm"} {
		enabled := os.Getenv("ANKI_TEST_"+strings.ToUpper(engine)) == "true"
		for _, executable := range []string{binary, os.Getenv("ANKI_NO_ENGINE_CLI")} {
			if executable != binary {
				executable = filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), executable)
			}
			available := enabled && executable == binary
			help := linkedCommand(t, executable, available, "--server-engine="+engine, "server", "run", "--help")
			if available && !strings.Contains(string(help), "serve") {
				t.Fatalf("selected %s engine did not expose Rust CLI help: %s", engine, help)
			}
		}
		address := strings.Replace(ankigrpc.SubprocessAddress(binary, "server", "serve"), "subprocess://", "subprocess+"+engine+"://", 1)
		command := exec.CommandContext(t.Context(), binary, "api", "i18n", "format-timespan", "--address", address)
		command.Stdin = strings.NewReader(`{"seconds":65,"context":"PRECISE"}`)
		output, err := command.CombinedOutput()
		switch {
		case !enabled:
			if err == nil || !strings.Contains(string(output), "scheme") || !strings.Contains(string(output), "unavailable") {
				t.Fatalf("unavailable %s engine was not rejected: %v; %s", engine, err, output)
			}
		case engine == "cgo":
			if err != nil || !strings.Contains(string(output), "minute") {
				t.Fatalf("explicit cgo worker used WASM or failed: %v; %s", err, output)
			}
		default:
			if err == nil || !strings.Contains(string(output), "compilation cache") {
				t.Fatalf("explicit WASM worker bypassed WASM cache failure: %v; %s", err, output)
			}
		}
		fmt.Fprintf(&report, "%s: enabled=%t; selected engine verified with unusable WASM cache\n", engine, enabled)
		if enabled {
			t.Run("standalone/"+engine, func(t *testing.T) {
				t.Setenv("ANKI_WASM_CACHE_DIR", "")
				standalone := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_STANDALONE_CLI"))
				address := strings.Replace(ankigrpc.SubprocessAddress(standalone, "serve"), "subprocess://", "subprocess+"+engine+"://", 1)
				output := apiCLI(t, binary, `{"seconds":65,"context":"PRECISE"}`, true, "api", "i18n", "format-timespan", "--address", address)
				if !strings.Contains(string(output), "minute") {
					t.Fatalf("standalone worker arguments were not preserved: %s", output)
				}
			})
			fmt.Fprintf(&report, "%s: standalone worker served with its original arguments\n", engine)
			aliasAddress := strings.Replace(ankigrpc.SubprocessAddress(alias, "server", "serve"), "subprocess://", "subprocess+"+engine+"://", 1)
			command := exec.CommandContext(t.Context(), binary, "api", "i18n", "format-timespan", "--address", aliasAddress)
			command.Stdin = strings.NewReader(`{"seconds":65,"context":"PRECISE"}`)
			output, err := command.CombinedOutput()
			if engine == "cgo" {
				if err != nil || !strings.Contains(string(output), "minute") {
					t.Fatalf("aliased cgo worker failed: %v; %s", err, output)
				}
			} else if err == nil || !strings.Contains(string(output), "compilation cache") {
				t.Fatalf("aliased WASM worker bypassed engine selection: %v; %s", err, output)
			}
			conflict := strings.Replace(ankigrpc.SubprocessAddress(binary, "--server-engine="+engine, "server", "serve"), "subprocess://", "subprocess+"+engine+"://", 1)
			command = exec.CommandContext(t.Context(), binary, "api", "i18n", "format-timespan", "--address", conflict)
			command.Stdin = strings.NewReader("{}")
			output, err = command.CombinedOutput()
			if err == nil || !strings.Contains(string(output), "scheme owns engine selection") {
				t.Fatalf("qualified %s scheme accepted a second engine selector: %v; %s", engine, err, output)
			}
		}
	}
	database := filepath.Join(t.TempDir(), "preferred.anki2")
	if err := os.WriteFile(database, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, address := range []string{"", ankigrpc.SubprocessAddress(binary, "server", "serve"), "inprocess://preferred"} {
		arguments := []string{"api", "i18n", "format-timespan", "--address", address}
		if address == "" {
			arguments = append(arguments, "--collection", database)
		}
		command := exec.CommandContext(t.Context(), binary, arguments...)
		command.Stdin = strings.NewReader(`{"seconds":65,"context":"PRECISE"}`)
		output, err := command.CombinedOutput()
		if os.Getenv("ANKI_TEST_CGO") == "true" {
			if err != nil || !strings.Contains(string(output), "minute") {
				t.Fatalf("generic strategy %q did not prefer cgo: %v; %s", address, err, output)
			}
		} else if err == nil || !strings.Contains(string(output), "compilation cache") {
			t.Fatalf("generic strategy %q did not fall back to WASM: %v; %s", address, err, output)
		}
	}
	fmt.Fprintln(&report, "generic subprocess and in-process strategies prefer cgo, otherwise WASM")
	linkedCommand(t, binary, false, "--server-engine=unknown", "server", "run", "--help")
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		if err := os.WriteFile(filepath.Join(directory, "worker-engine-selection.txt"), []byte(report.String()), 0o600); err != nil {
			t.Fatalf("retain worker selection verification: %v", err)
		}
	}
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
		{"api", "tags", "all-tags", "--collection", database, "--input-format", "yaml"},
		{"api", "backend-collection", "open-collection", "--collection", database},
	} {
		apiCLI(t, binary, "{}", false, arguments...)
	}
	apiCLI(t, binary, `{"unknown":true}`, false, "api", "tags", "all-tags", "--collection", database)
	if info, err := os.Stat(database); err != nil || info.Size() != 0 {
		t.Fatalf("invalid requests opened collection: %v %v", info, err)
	}
	apiCLI(t, binary, "{}", true, "api", "backend-tags", "all-tags", "--collection", database, "--address=")
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

// TestAPINetworkWorkflow verifies raw access borrows collection lifetime while
// managed paths close storage and preserve the remote server and committed data.
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
			apiCLI(t, binary, request, true, "api", "backend-collection", "open-collection", "--address", endpoint)
			for range 2 {
				apiCLI(t, binary, "{}", true, "api", "tags", "all-tags", "--address", endpoint)
			}
			formatted := apiCLI(t, binary, `{"seconds":65,"context":"PRECISE"}`, true, "api", "i18n", "format-timespan", "--address", endpoint)
			if !strings.Contains(string(formatted), "val") {
				t.Fatalf("i18n response is empty: %s", formatted)
			}
			response := apiCLI(t, binary, "{}", true, "api", "backend-collection", "close-collection", "--address", endpoint)
			managed := filepath.Join(root, "managed.anki2")
			if err := os.WriteFile(managed, nil, 0o600); err != nil {
				t.Fatal(err)
			}
			deck := apiCLI(t, binary, "{}", true, "api", "decks", "new-deck", "--collection", managed, "--address", endpoint)
			var managedRequest map[string]any
			if err := json.Unmarshal(deck, &managedRequest); err != nil {
				t.Fatal(err)
			}
			managedRequest["name"] = "Managed " + transport
			encoded, err := json.Marshal(managedRequest)
			if err != nil {
				t.Fatal(err)
			}
			apiCLI(t, binary, string(encoded), true, "api", "decks", "add-deck", "--collection", managed, "--address", endpoint)
			lookup := fmt.Sprintf(`{"val":%q}`, managedRequest["name"])
			managedResponse := apiCLI(t, binary, lookup, true, "api", "decks", "get-deck-id-by-name", "--collection", managed, "--address", endpoint)
			apiCLI(t, binary, "{}", false, "api", "tags", "all-tags", "--address", endpoint)
			apiCLI(t, binary, `{"seconds":65,"context":"PRECISE"}`, true, "api", "i18n", "format-timespan", "--address", endpoint)
			request = fmt.Sprintf(`{"collection_path":%q,"media_folder_path":%q,"media_db_path":%q}`, managed, filepath.Join(root, "managed.media"), filepath.Join(root, "managed.media.db2"))
			apiCLI(t, binary, request, true, "api", "backend-collection", "open-collection", "--address", endpoint)
			persisted := apiCLI(t, binary, lookup, true, "api", "decks", "get-deck-id-by-name", "--address", endpoint)
			if !bytes.Equal(persisted, managedResponse) {
				t.Fatalf("managed network write did not persist: %s != %s", persisted, managedResponse)
			}
			apiCLI(t, binary, "{}", true, "api", "backend-collection", "close-collection", "--address", endpoint)
			if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
				if err := os.WriteFile(filepath.Join(destination, "api-"+transport+"-managed.json"), managedResponse, 0o600); err != nil {
					t.Fatal(err)
				}
				data, err := os.ReadFile(managed)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(destination, "api-"+transport+"-managed.anki2"), data, 0o600); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(destination, "api-"+transport+"-response.json"), formatted, 0o600); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(destination, "api-"+transport+"-close.json"), response, 0o600); err != nil {
					t.Fatal(err)
				}
			}
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

// TestExternalOnlyExecutable verifies pure-Go external access with both engines absent.
func TestExternalOnlyExecutable(t *testing.T) {
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_NO_ENGINE_CLI"))
	image, err := elf.Open(binary)
	if err != nil {
		t.Fatal(err)
	}
	libraries, libraryErr := image.ImportedLibraries()
	closeErr := image.Close()
	if libraryErr != nil || closeErr != nil || len(libraries) != 0 {
		t.Fatalf("external-only command links native code: %v; %v; %v", libraries, libraryErr, closeErr)
	}
	root := t.TempDir()
	address := externalServer(t, installedAPI(t), root, "tcp", "")
	response := apiCLI(t, binary, `{ "seconds":65, "context":"PRECISE" }`, true, "api", "i18n", "format-timespan", "--address", address)
	if !strings.Contains(string(response), "minute") {
		t.Fatalf("external-only command lost external clients: %s", response)
	}
	database := filepath.Join(root, "external-only.anki2")
	if err := os.WriteFile(database, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	apiCLI(t, binary, "{}", true, "api", "tags", "all-tags", "--collection", database, "--address", address)
	for _, arguments := range [][]string{
		{"server", "serve"},
		{"api", "i18n", "format-timespan", "--address", "inprocess://unavailable"},
		{"api", "i18n", "format-timespan", "--address", "inprocess+cgo://unavailable"},
	} {
		command := exec.CommandContext(t.Context(), binary, arguments...)
		command.Stdin = strings.NewReader("{}")
		output, err := command.CombinedOutput()
		if err == nil || !strings.Contains(string(output), "unavailable") && !strings.Contains(string(output), "no local engines") {
			t.Fatalf("external-only local hosting did not fail clearly: %v; %s", err, output)
		}
	}
	if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
		if err := os.WriteFile(filepath.Join(destination, "external-only-response.json"), response, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
