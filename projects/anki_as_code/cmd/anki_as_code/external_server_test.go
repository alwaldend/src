package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/pkg/ankiascode"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"github.com/ankitects/anki/proto/anki/import_export"
)

// externalServer starts a real linked listener and retains process ownership in the test.
func externalServer(t *testing.T, binary, directory, transport string) string {
	t.Helper()
	listen := "127.0.0.1:0"
	if transport == "unix" {
		listen = "external.sock"
	}
	ctx, cancel := context.WithTimeout(context.WithoutCancel(t.Context()), 2*time.Minute)
	command := exec.CommandContext(ctx, binary, "server", "serve", "--transport", transport, "--listen", listen)
	command.Dir = directory
	stderr, err := command.StderrPipe()
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	command.Stdout = &stdout
	if err := command.Start(); err != nil {
		cancel()
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()
	t.Cleanup(func() {
		defer cancel()
		if err := command.Process.Signal(os.Interrupt); err != nil {
			t.Errorf("interrupt external server: %v", err)
		}
		select {
		case err := <-done:
			var exit *exec.ExitError
			if err != nil && (!errors.As(err, &exit) || exit.Sys().(syscall.WaitStatus).Signal() != syscall.SIGINT) {
				t.Errorf("external server shutdown: %v stdout=%q", err, stdout.String())
			}
			if stdout.Len() != 0 {
				t.Errorf("network server wrote protocol bytes to stdout: %q", stdout.String())
			}
		case <-time.After(5 * time.Second):
			cancel()
			<-done
			t.Error("external server did not stop after interruption")
		}
	})
	ready := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(stderr)
		reported := false
		for scanner.Scan() {
			if !reported && strings.HasPrefix(scanner.Text(), "Listening on ") {
				ready <- strings.TrimPrefix(scanner.Text(), "Listening on ")
				reported = true
			}
		}
		if !reported {
			ready <- ""
		}
	}()
	select {
	case address := <-ready:
		if address == "" {
			t.Fatal("external server exited before readiness")
		}
		return address
	case <-ctx.Done():
		t.Fatal("external server did not become ready")
		return ""
	}
}

// externalArchive creates a native archive fixture on the same running server.
func externalArchive(t *testing.T, address, database, archive string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	backend := ankiclient.NewGRPC(ankiclient.BackendGRPCOptions{Address: address})
	defer func() {
		if err := backend.Stop(context.Background()); err != nil {
			t.Errorf("close archive transport: %v", err)
		}
	}()
	if err := backend.Start(ctx); err != nil {
		t.Fatal(err)
	}
	prepared := ankiclient.NewCollection(ankiclient.BackendCollectionOptions{Backend: backend, Path: database})
	defer func() {
		if err := prepared.Stop(context.Background()); err != nil {
			t.Errorf("close archive collection: %v", err)
		}
	}()
	if err := prepared.Start(ctx); err != nil {
		t.Fatal(err)
	}
	connection, err := prepared.Connection()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx,
		&import_export.ExportCollectionPackageRequest{OutPath: archive, IncludeMedia: true}); err != nil {
		t.Fatal(err)
	}
}

// TestApplicationBackendSelection reuses one App across external endpoints and embedded defaults.
// Each external listener stops before the next operation so a retained address cannot succeed.
func TestApplicationBackendSelection(t *testing.T) {
	binary, root := installedLinkedCLI(t)
	t.Chdir(root)
	t.Setenv("TMPDIR", root)
	application := ankiascode.New(ankiascode.Options{BackendBinary: binary, BackendArguments: []string{"server", "serve"}})
	database := filepath.Join(root, "application.anki2")
	if err := os.WriteFile(database, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	wanted := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{
		Decks: map[string]*ankiascodepb.DeckState{"decks/Default/deck.anki.toml": {Id: 1}},
	}}
	for _, transport := range []string{"tcp", "unix", "embedded"} {
		t.Run(transport, func(t *testing.T) {
			options := ankiascodeapi.CollectionOptions{Path: database, StateKey: "application/" + transport}
			if transport != "embedded" {
				options.Address = externalServer(t, binary, root, transport)
			}
			if err := application.ColStateReplace(t.Context(), options, wanted); err != nil {
				t.Fatal(err)
			}
			stored, err := application.ColStateDump(t.Context(), options)
			if err != nil || stored.GetResources().GetDecks()["decks/Default/deck.anki.toml"].GetId() != 1 {
				t.Fatalf("selected backend lost namespace bindings: %v %v", stored, err)
			}
			if options.Address != "" {
				remoteOnly := ankiascode.New(ankiascode.Options{})
				stored, err := remoteOnly.ColStateDump(t.Context(), options)
				if err != nil || stored.GetResources().GetDecks()["decks/Default/deck.anki.toml"].GetId() != 1 {
					t.Fatalf("external options required an embedded executable: %v %v", stored, err)
				}
			}
		})
	}
	stored, err := application.ColStateDump(t.Context(), ankiascodeapi.CollectionOptions{Path: database})
	if err != nil || stored.GetVersion() != 1 || len(stored.GetResources().GetDecks()) != 0 {
		t.Fatalf("zero options retained a previous endpoint or namespace: %v %v", stored, err)
	}
	if _, err := application.ColStateDump(t.Context(), ankiascodeapi.CollectionOptions{Path: database, Address: "invalid"}); err == nil || !strings.Contains(err.Error(), "address") {
		t.Fatalf("invalid operation address fell back to embedded: %v", err)
	}
	if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
		data, err := os.ReadFile(database)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, "application-backends.anki2"), data, 0o600); err != nil {
			t.Fatal(err)
		}
		data, err = json.MarshalIndent(stored, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, "application-default-state.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestExternalServerWorkflow exercises every collection command against reusable TCP and Unix listeners.
// It retains archive, state and plan artifacts for independent inspection.
func TestExternalServerWorkflow(t *testing.T) {
	for _, transport := range []string{"tcp", "unix"} {
		t.Run(transport, func(t *testing.T) {
			binary, root := installedLinkedCLI(t)
			t.Chdir(root)
			t.Setenv("PATH", "")
			t.Setenv("TMPDIR", root)
			address := externalServer(t, binary, root, transport)
			database := filepath.Join(root, "external.anki2")
			if err := os.WriteFile(database, nil, 0o600); err != nil {
				t.Fatal(err)
			}
			linkedCommand(t, binary, true, "col", "state", "dump", "--collection", database)
			linkedCommand(t, binary, true, "col", "state", "dump", "--collection", database, "--address=", "--format", "json")
			text, stateFile := filepath.Join(root, "text"), filepath.Join(root, "state.anki.toml")
			linkedCommand(t, binary, true, "col", "--address", address, "export", "--collection", database, "--output", text, "--state-output", stateFile)
			stateArguments := []string{"col", "state", "--address", address, "--collection", database, "--state-key", "external-test"}
			linkedCommand(t, binary, true, append(append([]string{}, stateArguments...), "update", stateFile)...)
			state := linkedCommand(t, binary, true, append(append([]string{}, stateArguments...), "dump", "--format", "json")...)
			var stored struct {
				// Resources contains the exported native deck identities used by import.
				Resources struct {
					// Decks maps configuration paths to their native identifiers.
					Decks map[string]struct {
						// ID is the existing native deck identifier.
						ID string `json:"id"`
					} `json:"decks"`
				} `json:"resources"`
			}
			if err := json.Unmarshal(state, &stored); err != nil || len(stored.Resources.Decks) == 0 {
				t.Fatalf("state update did not persist native bindings: %s %v", state, err)
			}
			empty := filepath.Join(root, "empty.anki.toml")
			if err := os.WriteFile(empty, []byte("version = \"1\"\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			linkedCommand(t, binary, true, append(append([]string{}, stateArguments...), "replace", empty)...)
			for path, binding := range stored.Resources.Decks {
				linkedCommand(t, binary, true, append(append([]string{}, stateArguments...), "import", "--text", text, path+"="+binding.ID)...)
				break
			}
			linkedCommand(t, binary, true, append(append([]string{}, stateArguments...), "update", stateFile)...)
			archive, rebuilt := filepath.Join(root, "base.colpkg"), filepath.Join(root, "rebuilt.colpkg")
			externalArchive(t, address, database, archive)
			original, err := os.ReadFile(archive)
			if err != nil {
				t.Fatal(err)
			}
			deck := filepath.Join(text, "decks", "External", "deck.anki.toml")
			if err := os.MkdirAll(filepath.Dir(deck), 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(deck, []byte("preset = \"Default\"\n[normal]\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			linkedCommand(t, binary, true, "col", "build", "--address", address, "--base", archive, "--text", text, "--output", rebuilt, "--state-key", "external-test")
			if actual, err := os.ReadFile(archive); err != nil || !bytes.Equal(original, actual) {
				t.Fatalf("build changed the base archive: %v", err)
			}
			linkedCommand(t, binary, true, "col", "apply", "--address", address, "--collection", archive, "--text", text, "--state-key", "external-test")
			plan := linkedCommand(t, binary, true, "col", "plan", "--address", address, "--collection", rebuilt, "--text", text, "--state-key", "external-test")
			var report ankiascodepb.CollectionChangePlan
			if err := protojson.Unmarshal(plan, &report); err != nil || len(report.GetChanges()) != 0 || report.GetStateChanged() {
				t.Fatalf("build did not converge: %s %v", plan, err)
			}
			linkedCommand(t, binary, true, "col", "apply", "--address", address, "--collection", rebuilt, "--text", text, "--state-key", "external-test")
			plan = linkedCommand(t, binary, true, "col", "plan", "--address", address, "--collection", rebuilt, "--text", text, "--state-key", "external-test")
			if err := protojson.Unmarshal(plan, &report); err != nil || len(report.GetChanges()) != 0 || report.GetStateChanged() {
				t.Fatalf("repeated apply did not converge: %s %v", plan, err)
			}
			published, err := os.ReadFile(rebuilt)
			if err != nil {
				t.Fatal(err)
			}
			exported := filepath.Join(root, "rebuilt-text")
			linkedCommand(t, binary, true, "col", "export", "--address", address, "--collection", rebuilt, "--output", exported, "--state-key", "external-test")
			if _, err := os.Stat(filepath.Join(exported, "decks", "External", "deck.anki.toml")); err != nil {
				t.Fatalf("native archive does not contain the new deck: %v", err)
			}
			if transport == "tcp" {
				linkedCommand(t, binary, true, "collection", "state", "dump", "--address", strings.TrimPrefix(address, "tcp://"), "--collection", database)
			} else {
				// The Linux cwd alias keeps absolute addresses within sockaddr limits under long test roots.
				absolute := "/proc/self/cwd/external.sock"
				for _, endpoint := range []string{absolute, "unix://" + absolute} {
					linkedCommand(t, binary, true, "col", "state", "dump", "--address", endpoint, "--collection", database)
				}
			}
			if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
				if err := os.WriteFile(filepath.Join(destination, "external-"+transport+"-plan.json"), plan, 0o600); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(destination, "external-"+transport+"-state.json"), state, 0o600); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(destination, "external-"+transport+".colpkg"), published, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestExternalAddressFailures rejects malformed endpoints and cancels unavailable servers without fallback.
func TestExternalAddressFailures(t *testing.T) {
	binary, root := installedLinkedCLI(t)
	database := filepath.Join(root, "unchanged.anki2")
	if err := os.WriteFile(database, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, address := range []string{"invalid", "unix://"} {
		linkedCommand(t, binary, false, "col", "state", "dump", "--collection", database, "--address", address)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			t.Errorf("close unresponsive endpoint: %v", err)
		}
	}()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, binary, "col", "state", "dump", "--collection", database, "--address", listener.Addr().String())
	command.Cancel = func() error {
		if err := command.Process.Signal(os.Interrupt); err != nil {
			return fmt.Errorf("interrupt waiting client: %w", err)
		}
		return nil
	}
	command.WaitDelay = 2 * time.Second
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	if err := listener.(*net.TCPListener).SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	connection, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := connection.Close(); err != nil {
			t.Errorf("close unresponsive connection: %v", err)
		}
	}()
	if err := command.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
	if err := command.Wait(); err == nil || ctx.Err() != nil || stdout.Len() != 0 || !strings.Contains(stderr.String(), "network readiness") {
		t.Fatalf("unavailable endpoint did not cancel cleanly: %v stdout=%q stderr=%q", err, stdout.String(), stderr.String())
	}
	if data, err := os.ReadFile(database); err != nil || len(data) != 0 {
		t.Fatalf("failed external connection changed the collection: %v", err)
	}
}
