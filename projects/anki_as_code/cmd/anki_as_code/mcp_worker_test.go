package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
)

// runMCPWorkerFixture forwards a real native worker and records its process lifetime.
func runMCPWorkerFixture(binary, events string) error {
	command := exec.Command(binary, "server", "serve")
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := command.Start(); err != nil {
		return fmt.Errorf("start recorded native worker: %w", err)
	}
	if err := recordMCPWorkerEvent(events, "started", command.Process.Pid); err != nil {
		_ = command.Process.Kill()
		_ = command.Wait()
		return fmt.Errorf("record worker startup: %w", err)
	}
	if err := command.Wait(); err != nil {
		return fmt.Errorf("wait for recorded native worker: %w", err)
	}
	if err := recordMCPWorkerEvent(events, "stopped", command.Process.Pid); err != nil {
		return fmt.Errorf("record worker shutdown: %w", err)
	}
	return nil
}

// recordMCPWorkerEvent appends public fixture process evidence without writing protocol bytes.
func recordMCPWorkerEvent(path, event string, pid int) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open worker evidence: %w", err)
	}
	_, writeErr := fmt.Fprintf(file, "%s %d\n", event, pid)
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return fmt.Errorf("append and close worker evidence: %w", err)
	}
	return nil
}

// TestMCPCollectionSessions verifies explicit handles, deferred publication and disconnect cleanup.
func TestMCPCollectionSessions(t *testing.T) {
	binary, root := installedLinkedCLI(t)
	events := filepath.Join(root, "workers.log")
	address := ankigrpc.SubprocessAddress(os.Args[0], "mcp-worker-fixture", binary, events)
	database := filepath.Join(root, "collection.anki2")
	archive := filepath.Join(root, "collection.colpkg")
	externalArchive(t, ankigrpc.SubprocessAddress(binary, "server", "serve"), database, archive)
	requestFile := filepath.Join(root, "get-deck.json")
	if err := os.WriteFile(requestFile, []byte(`{"did":"1"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	session, stop, disconnect := connectMCP(t, binary)
	if _, err := os.Stat(events); !os.IsNotExist(err) {
		t.Fatalf("discovery started a worker: %v", err)
	}
	opened := callMCP(t, session, "col_open", map[string]any{"collection": archive, "address": address}, true)
	handle, ok := opened["handle"].(string)
	if !ok || handle == "" {
		t.Fatalf("open omitted its handle: %v", opened)
	}
	args := map[string]any{"handle": handle, "request": map[string]any{"did": "1"}}
	deck := callMCP(t, session, "api_decks_get_deck", args, true)
	originalName := deck["name"]
	deck["name"] = "Published on close"
	args["request"] = deck
	callMCP(t, session, "api_decks_update_deck", args, true)
	callMCP(t, session, "api_decks_get_deck", map[string]any{"handle": handle, "request": map[string]any{"did": "-1"}}, false)
	callMCP(t, session, "api_decks_get_deck", map[string]any{"handle": handle, "request": map[string]any{"did": "1"}}, true)
	text := filepath.Join(root, "text")
	snapshot := filepath.Join(root, "state.toml")
	callMCP(t, session, "col_export", map[string]any{"handle": handle, "output": text, "state_output": snapshot}, true)
	snapshotBytes, err := os.ReadFile(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	state := &ankiascodepb.State{}
	if err := codec.Decode(snapshotBytes, codec.TOML, state); err != nil {
		t.Fatal(err)
	}
	encoded, err := codec.Encode(state, codec.JSON)
	if err != nil {
		t.Fatal(err)
	}
	var bindings map[string]any
	if err := json.Unmarshal(encoded, &bindings); err != nil {
		t.Fatal(err)
	}
	callMCP(t, session, "col_state_replace", map[string]any{"handle": handle, "state": bindings}, true)
	plan := callMCP(t, session, "col_plan", map[string]any{"handle": handle, "text": text}, true)
	if len(plan) != 0 {
		t.Fatalf("handle-backed export did not converge: %v", plan)
	}
	callMCP(t, session, "col_apply", map[string]any{"handle": handle, "text": text}, true)
	callMCP(t, session, "col_state_dump", map[string]any{"handle": handle}, true)
	callMCP(t, session, "col_open", map[string]any{"collection": archive, "address": address}, false)
	callMCP(t, session, "api_decks_get_deck", map[string]any{"handle": handle, "collection": archive, "request": map[string]any{"did": "1"}}, false)
	verifyPublished := func(want any) {
		data := linkedCommand(t, binary, true, "api", "decks", "get-deck", "--collection", archive, "--input", requestFile)
		var actual map[string]any
		if err := json.Unmarshal(data, &actual); err != nil || actual["name"] != want {
			t.Fatalf("archive publication mismatch: %v %s", err, data)
		}
	}
	verifyPublished(originalName)
	assertMCPWorkerEvents(t, events, 1, 0)
	callMCP(t, session, "col_close", map[string]any{"handle": handle}, true)
	assertMCPWorkerEvents(t, events, 1, 1)
	verifyPublished("Published on close")
	callMCP(t, session, "col_state_dump", map[string]any{"handle": handle}, false)
	readonly := callMCP(t, session, "col_open", map[string]any{"collection": archive, "address": address, "write_back": false}, true)
	args = map[string]any{"handle": readonly["handle"], "request": map[string]any{"did": "1"}}
	deck = callMCP(t, session, "api_decks_get_deck", args, true)
	deck["name"] = "Unpublished session"
	args["request"] = deck
	callMCP(t, session, "api_decks_update_deck", args, true)
	callMCP(t, session, "col_close", map[string]any{"handle": readonly["handle"]}, true)
	verifyPublished("Published on close")
	second := callMCP(t, session, "col_open", map[string]any{"collection": archive, "address": address}, true)
	args = map[string]any{"handle": second["handle"], "request": map[string]any{"did": "1"}}
	deck = callMCP(t, session, "api_decks_get_deck", args, true)
	deck["name"] = "Published on disconnect"
	args["request"] = deck
	callMCP(t, session, "api_decks_update_deck", args, true)
	// Independent handles own independent worker/collection lifetimes.
	third := callMCP(t, session, "col_open", map[string]any{"collection": database, "address": address}, true)
	callMCP(t, session, "api_decks_get_deck", map[string]any{"handle": third["handle"], "request": map[string]any{"did": "1"}}, true)
	assertMCPWorkerEvents(t, events, 4, 2)
	disconnect()
	stop()
	assertMCPWorkerEvents(t, events, 4, 4)
	verifyPublished("Published on disconnect")
	data, err := os.ReadFile(events)
	if err != nil {
		t.Fatal(err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "mcp-worker-lifetimes.log"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// assertMCPWorkerEvents checks complete process starts and graceful shutdowns.
func assertMCPWorkerEvents(t *testing.T, path string, starts, stops int) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(data), "started ") != starts || strings.Count(string(data), "stopped ") != stops {
		t.Fatalf("worker reuse or shutdown mismatch: want %d starts/%d stops, got %s", starts, stops, data)
	}
}
