package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mcpTestTransport retains the physical connection to distinguish disconnect from graceful close.
type mcpTestTransport struct {
	// command starts the real installed MCP executable.
	command *mcp.CommandTransport
	// connection permits closing stdio while a request is still active.
	connection mcp.Connection
}

// Connect starts and retains the physical stdio connection for the SDK client.
func (transport *mcpTestTransport) Connect(ctx context.Context) (mcp.Connection, error) {
	connection, err := transport.command.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("start MCP fixture transport: %w", err)
	}
	transport.connection = connection
	return connection, nil
}

// connectMCP starts the installed executable with real SDK transport and bounded shutdown.
func connectMCP(t *testing.T, binary string) (*mcp.ClientSession, func(), func()) {
	t.Helper()
	command := exec.Command(binary, "mcp", "serve")
	command.Env = append(os.Environ(), "PATH=")
	command.Stderr = os.Stderr
	client := mcp.NewClient(&mcp.Implementation{Name: "anki-mcp-acceptance", Version: "1"}, nil)
	transport := &mcpTestTransport{command: &mcp.CommandTransport{Command: command}}
	session, err := client.Connect(t.Context(), transport, nil)
	if err != nil {
		t.Fatalf("connect installed MCP server: %v", err)
	}
	var once sync.Once
	stop := func() {
		once.Do(func() {
			if err := session.Close(); err != nil {
				t.Errorf("close MCP session: %v", err)
			}
			if command.ProcessState == nil || !command.ProcessState.Success() {
				t.Errorf("MCP process did not exit successfully: %v", command.ProcessState)
			}
		})
	}
	t.Cleanup(stop)
	return session, stop, func() {
		if err := transport.connection.Close(); err != nil {
			t.Errorf("disconnect MCP transport: %v", err)
		}
	}
}

// callMCP requires a protocol-valid tool result and its structured/text JSON agreement.
func callMCP(t *testing.T, session *mcp.ClientSession, name string, arguments map[string]any, success bool) map[string]any {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 90*time.Second)
	defer cancel()
	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: arguments})
	if err != nil {
		t.Fatalf("call MCP tool %s: %v", name, err)
	}
	if result.IsError == success {
		var diagnostics []string
		for _, content := range result.Content {
			if text, ok := content.(*mcp.TextContent); ok {
				diagnostics = append(diagnostics, text.Text)
			}
		}
		t.Fatalf("MCP tool %s success=%v: %s", name, success, strings.Join(diagnostics, "; "))
	}
	if !success {
		if len(result.Content) == 0 {
			t.Fatalf("MCP tool %s omitted error diagnostics", name)
		}
		return nil
	}
	encoded, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var value map[string]any
	if err := json.Unmarshal(encoded, &value); err != nil || value == nil {
		t.Fatalf("MCP tool %s omitted structured object: %s %v", name, encoded, err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("MCP tool %s omitted JSON text fallback", name)
	}
	text, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("MCP tool %s returned a non-text fallback", name)
	}
	var fallback any
	if err := json.Unmarshal([]byte(text.Text), &fallback); err != nil {
		t.Fatal(err)
	}
	normalized, err := json.Marshal(fallback)
	if err != nil || !bytes.Equal(encoded, normalized) {
		t.Fatalf("MCP tool %s text and structured outputs differ: %v", name, err)
	}
	return value
}

// TestMCPDiscovery verifies canonical typed discovery without a backend, PATH or writable temp storage.
// Advertised oneof schemas must not require a selection, count null as selected,
// reject a single non-null member beside null, or accept two non-null members.
func TestMCPDiscovery(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_NO_ENGINE_CLI"))
	data, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(root, "anki_as_code")
	if err := os.WriteFile(binary, data, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TMPDIR", filepath.Join(root, "absent"))
	session, stop, _ := connectMCP(t, binary)
	validator := mcp.NewServer(&mcp.Implementation{Name: "advertised-schema-consumer", Version: "1"}, nil)
	names := map[string]bool{}
	var catalog []*mcp.Tool
	for tool, err := range session.Tools(t.Context(), nil) {
		if err != nil {
			t.Fatal(err)
		}
		if names[tool.Name] || tool.InputSchema == nil || tool.OutputSchema == nil {
			t.Fatalf("invalid advertised MCP tool: %+v", tool)
		}
		names[tool.Name] = true
		// An independent SDK consumer must resolve every advertised input/output schema.
		mcp.AddTool(validator, tool, func(context.Context, *mcp.CallToolRequest, map[string]any) (*mcp.CallToolResult, map[string]any, error) {
			return nil, map[string]any{}, nil
		})
		catalog = append(catalog, tool)
	}
	for _, name := range []string{"config_dump", "config_export", "col_export", "col_plan", "col_apply", "col_build", "col_state_dump", "col_state_import", "col_state_update", "col_state_replace", "col_open", "col_close", "api_notes_get_note", "api_tags_all_tags", "api_search_build_search_string", "api_backend_import_export_export_collection_package"} {
		if !names[name] {
			t.Fatalf("missing MCP tool %q", name)
		}
	}
	if names["api_backend_tags_all_tags"] {
		t.Fatal("forwarded backend alias duplicates the declared tags tool")
	}
	if len(catalog) < 281 {
		t.Fatalf("native catalog incomplete: %d tools", len(catalog))
	}
	consumerContext, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := validator.Connect(consumerContext, serverTransport, nil)
	if err != nil {
		t.Fatalf("connect advertised-schema validator: %v", err)
	}
	defer func() {
		if err := serverSession.Close(); err != nil {
			t.Errorf("close advertised-schema validator: %v", err)
		}
	}()
	client := mcp.NewClient(&mcp.Implementation{Name: "advertised-schema-probe", Version: "1"}, nil)
	consumer, err := client.Connect(consumerContext, clientTransport, nil)
	if err != nil {
		t.Fatalf("connect advertised-schema consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			t.Errorf("close advertised-schema consumer: %v", err)
		}
	}()
	for name, request := range map[string]map[string]any{
		"neither selected":   {},
		"normal null":        {"normal": nil},
		"filtered null":      {"filtered": nil},
		"both null":          {"normal": nil, "filtered": nil},
		"normal selected":    {"normal": map[string]any{}},
		"filtered selected":  {"filtered": map[string]any{}},
		"normal with null":   {"normal": map[string]any{}, "filtered": nil},
		"filtered with null": {"normal": nil, "filtered": map[string]any{}},
	} {
		t.Run(name, func(t *testing.T) {
			callMCP(t, consumer, "api_decks_add_deck", map[string]any{"collection": "absent.anki2", "request": request}, true)
		})
	}
	t.Run("both selected", func(t *testing.T) {
		callMCP(t, consumer, "api_decks_add_deck", map[string]any{
			"collection": "absent.anki2", "request": map[string]any{"normal": map[string]any{}, "filtered": map[string]any{}},
		}, false)
	})
	callMCP(t, session, "api_tags_all_tags", map[string]any{"collection": "absent.anki2", "request": map[string]any{"unknown": true}}, false)
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := json.MarshalIndent(catalog, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "mcp-tools.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	stop()
}

// TestMCPWorkflow exercises every workflow and native mutation through installed stdio.
// Failures must not open storage for invalid input, corrupt protocol framing,
// lose integer identity, publish disabled archive writes, or prevent later calls.
func TestMCPWorkflow(t *testing.T) {
	binary := installedAPI(t)
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
	session, stop, _ := connectMCP(t, binary)
	native := func(name string, request any, success bool) map[string]any {
		return callMCP(t, session, "api_"+name, map[string]any{"collection": database, "request": request}, success)
	}
	native("tags_all_tags", map[string]any{"unknown": true}, false)
	native("i18n_format_timespan", map[string]any{"seconds": 1, "context": "UNKNOWN_ENUM"}, false)
	native("notes_get_note", map[string]any{"nid": "9223372036854775808"}, false)
	native("decks_add_deck", map[string]any{"normal": map[string]any{}, "filtered": map[string]any{}}, false)
	callMCP(t, session, "col_state_dump", map[string]any{"collection": database, "unknown": true}, false)
	if info, err := os.Stat(database); err != nil || info.Size() != 0 {
		t.Fatalf("invalid MCP input opened collection: %v %v", info, err)
	}
	callMCP(t, session, "col_state_dump", map[string]any{"collection": database}, true)
	deck := native("decks_new_deck", map[string]any{}, true)
	deck["name"] = "MCP native deck"
	native("decks_add_deck", deck, true)
	search := native("search_build_search_string", map[string]any{"group": map[string]any{
		"joiner": "AND", "nodes": []any{map[string]any{"deck": "MCP native deck"}, map[string]any{"negated": map[string]any{"tag": "excluded"}}},
	}}, true)
	if !strings.Contains(search["val"].(string), "excluded") {
		t.Fatalf("recursive search payload was not preserved: %v", search)
	}
	identity := native("decks_get_deck_id_by_name", map[string]any{"val": "MCP native deck"}, true)["did"]
	if _, ok := identity.(string); !ok {
		t.Fatalf("native identity is not a ProtoJSON string: %v", identity)
	}
	native("media_add_media_file", map[string]any{"desired_name": "mcp.txt", "data": base64.StdEncoding.EncodeToString([]byte("MCP fixture media"))}, true)
	text := filepath.Join(root, "configs")
	snapshot := filepath.Join(root, "state.anki.toml")
	callMCP(t, session, "col_export", map[string]any{"collection": database, "output": text, "state_output": snapshot}, true)
	configs := callMCP(t, session, "config_dump", map[string]any{"text": text}, true)
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
	callMCP(t, session, "col_state_replace", map[string]any{"collection": database, "state": bindings}, true)
	configs["decks"] = append(configs["decks"].([]any), map[string]any{"title": "MCP workflow deck", "path": "decks/MCP workflow deck/deck.anki.toml", "deck": map[string]any{"preset": "Default", "normal": map[string]any{}}})
	callMCP(t, session, "config_export", map[string]any{"output": text, "configs": configs}, true)
	plan := callMCP(t, session, "col_plan", map[string]any{"collection": database, "text": text}, true)
	if changes, ok := plan["changes"].([]any); !ok || len(changes) == 0 {
		t.Fatalf("MCP plan omitted the declared new deck: %v", plan)
	}
	callMCP(t, session, "col_apply", map[string]any{"collection": database, "text": text}, true)
	plan = callMCP(t, session, "col_plan", map[string]any{"collection": database, "text": text}, true)
	if changes, ok := plan["changes"].([]any); ok && len(changes) != 0 {
		t.Fatalf("MCP apply did not converge: %v", plan)
	}
	bindings = callMCP(t, session, "col_state_dump", map[string]any{"collection": database}, true)
	callMCP(t, session, "col_state_update", map[string]any{"collection": database, "state": bindings}, true)
	callMCP(t, session, "col_state_import", map[string]any{"collection": database, "text": text, "bindings": map[string]any{"decks/MCP native deck/deck.anki.toml": identity}}, true)
	if got := native("decks_get_deck_id_by_name", map[string]any{"val": "MCP native deck"}, true)["did"]; got != identity {
		t.Fatalf("workflow changed native identity: %v != %v", got, identity)
	}
	archive := filepath.Join(root, "base.colpkg")
	native("backend_import_export_export_collection_package", map[string]any{"out_path": archive, "include_media": true}, true)
	base, err := os.ReadFile(archive)
	if err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(root, "built.colpkg")
	callMCP(t, session, "col_build", map[string]any{"base": archive, "text": text, "output": output}, true)
	if current, err := os.ReadFile(archive); err != nil || !bytes.Equal(base, current) {
		t.Fatalf("MCP build changed base archive: %v", err)
	}
	callMCP(t, session, "api_decks_get_deck_id_by_name", map[string]any{"collection": output, "request": map[string]any{"val": "MCP workflow deck"}}, true)
	deck["name"] = "Discarded MCP archive deck"
	callMCP(t, session, "api_decks_add_deck", map[string]any{"collection": archive, "write_back": false, "request": deck}, true)
	if current, err := os.ReadFile(archive); err != nil || !bytes.Equal(base, current) {
		t.Fatalf("disabled write_back published archive mutation: %v", err)
	}
	native("notes_get_note", map[string]any{"nid": "9223372036854775806"}, false)
	native("tags_all_tags", map[string]any{}, true)
	stop()
	if entries, err := os.ReadDir(scratch); err != nil || len(entries) != 0 {
		t.Fatalf("MCP workflow retained collection scratch: %v %v", entries, err)
	}
	if outputDir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); outputDir != "" {
		for _, artifact := range []struct {
			// name identifies the stable public output artifact.
			name string
			// value is the protocol-observed verification evidence.
			value any
		}{{"mcp-state.json", bindings}, {"mcp-plan.json", plan}} {
			data, err := json.MarshalIndent(artifact.value, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(outputDir, artifact.name), data, 0o600); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// TestMCPCancellation verifies cancellation and disconnect reach an active native RPC.
func TestMCPCancellation(t *testing.T) {
	for _, disconnect := range []bool{false, true} {
		t.Run(map[bool]string{false: "cancel", true: "disconnect"}[disconnect], func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			started, cancelled := make(chan struct{}), make(chan struct{})
			server := grpc.NewServer(grpc.UnknownServiceHandler(func(_ any, stream grpc.ServerStream) error {
				close(started)
				<-stream.Context().Done()
				close(cancelled)
				return status.Error(codes.Canceled, "fixture observed MCP cancellation")
			}))
			t.Cleanup(server.Stop)
			go func() { _ = server.Serve(listener) }()
			session, stop, closeTransport := connectMCP(t, installedAPI(t))
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			done := make(chan error, 1)
			go func() {
				result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "api_i18n_format_timespan", Arguments: map[string]any{"address": listener.Addr().String(), "request": map[string]any{"seconds": 65}}})
				if err == nil && result.IsError {
					err = fmt.Errorf("native tool reported cancellation")
				}
				done <- err
			}()
			select {
			case <-started:
			case <-time.After(15 * time.Second):
				t.Fatal("MCP native request did not reach cancellation fixture")
			}
			if disconnect {
				closeTransport()
			} else {
				cancel()
			}
			select {
			case <-cancelled:
			case <-time.After(15 * time.Second):
				t.Fatal("native RPC did not observe cancellation")
			}
			if err := <-done; err == nil {
				t.Fatal("cancelled MCP call succeeded")
			}
			if !disconnect {
				result := callMCP(t, session, "api_i18n_format_timespan", map[string]any{"address": "inprocess://after-cancellation", "request": map[string]any{"seconds": 65}}, true)
				if !strings.Contains(result["val"].(string), "minute") {
					t.Fatalf("session unusable after cancellation: %v", result)
				}
			}
			stop()
		})
	}
}
