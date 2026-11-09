package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestProtobufPlanFiles checks the public file contract through the real CLI.
// Failure cases are legacy schemas, unknown envelope/payload fields, ambiguous
// operations and an unsupported later operation causing an earlier mutation.
func TestProtobufPlanFiles(t *testing.T) {
	f, server := newFixture(t)
	defer server.Close()
	path := filepath.Join(t.TempDir(), "plan.json")
	output, err := runCLI(t, server, "", "plan", "--config", writePolicy(t, `force_delete_post("7")`), "--out", path)
	if err != nil {
		t.Fatalf("plan: %v: %s", err, output)
	}
	document := readPlanDocument(t, path)
	if _, exists := document["version"]; exists {
		t.Fatal("legacy plan version emitted")
	}
	if _, exists := document["actions"]; exists {
		t.Fatal("legacy action list emitted")
	}
	payload, ok := document["plan"].(map[string]any)
	if !ok {
		t.Fatal("protobuf plan payload missing")
	}
	operations, ok := payload["operations"].([]any)
	if !ok || len(operations) != 1 {
		t.Fatalf("operations: %v", payload)
	}
	operation := operations[0].(map[string]any)["delete_post"].(map[string]any)
	if operation["post_id"] != "7" || operation["forced"] != true {
		t.Fatalf("operation: %v", operation)
	}
	if _, ok := payload["cleanup_report"].(map[string]any); !ok {
		t.Fatal("cleanup report missing")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		if err := os.WriteFile(filepath.Join(directory, "protobuf-plan.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	output, err = runCLI(t, server, "", "apply", "--plan", path)
	if err != nil || !reflect.DeepEqual(f.mutations, []string{"/2/tweets/7"}) {
		t.Fatalf("apply: %v: %s; mutations %v", err, output, f.mutations)
	}
	for name, text := range map[string]string{
		"legacy":              `{"version":1,"account_id":"100","evaluated_at":"2026-11-13T00:00:00Z","actions":[{"kind":"delete_post","target_id":"7"}]}`,
		"unknown-envelope":    `{"account_id":"100","evaluated_at":"2026-11-13T00:00:00Z","extra":true,"plan":{"operations":[]}}`,
		"unknown-payload":     `{"operations":[{"delete_post":{"post_id":"7","extra":true}}]}`,
		"ambiguous-operation": `{"operations":[{"delete_post":{"post_id":"7"},"unrepost":{"post_id":"901"}}]}`,
		"unsupported-later":   `{"operations":[{"delete_post":{"post_id":"7"}},{"unlike":{"post_id":"101"}}]}`,
	} {
		t.Run(name, func(t *testing.T) {
			fixture, endpoint := newFixture(t)
			defer endpoint.Close()
			if name != "legacy" && name != "unknown-envelope" {
				text = `{"account_id":"100","evaluated_at":"2026-11-13T00:00:00Z","plan":` + text + `}`
			}
			input := filepath.Join(t.TempDir(), "rejected.json")
			if !json.Valid([]byte(text)) {
				t.Fatal("invalid acceptance fixture JSON")
			}
			if err := os.WriteFile(input, []byte(text), 0o600); err != nil {
				t.Fatal(err)
			}
			output, err := runCLI(t, endpoint, "", "apply", "--plan", input)
			if err == nil || fixture.requests != 0 {
				t.Fatalf("invalid plan reached X: %v: %s; requests %d", err, output, fixture.requests)
			}
		})
	}
}

// TestPolicyWrapperBoundaries checks full and overlapping long-string delimiters.
// Valid policies may end inside a line comment without a trailing newline.
func TestPolicyWrapperBoundaries(t *testing.T) {
	for name, suffix := range map[string]string{
		"suffix-one":    "]=",
		"suffix-two":    "]=] ]==",
		"suffix-three":  "]=] ]==] ]===",
		"full-closers":  "]=] ]==] ]===]",
		"plain-bracket": "]",
		"newline":       "]=\n",
	} {
		t.Run(name, func(t *testing.T) {
			fixture, server := newFixture(t)
			defer server.Close()
			path := filepath.Join(t.TempDir(), "plan.json")
			source := writePolicy(t, `delete_post("7") -- `+suffix)
			output, err := runCLI(t, server, "", "plan", "--config", source, "--out", path)
			if err != nil {
				t.Fatalf("valid policy wrapper: %v: %s", err, output)
			}
			operations := readPlanDocument(t, path)["plan"].(map[string]any)["operations"].([]any)
			if !reflect.DeepEqual(operations, []any{map[string]any{"delete_post": map[string]any{"post_id": "7", "forced": false}}}) {
				t.Fatalf("wrapper changed selection: %v", operations)
			}
			if fixture.requests != 0 {
				t.Fatalf("planning made %d X requests", fixture.requests)
			}
			if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
				artifact, err := json.Marshal(map[string]any{"case": name, "operations": operations, "x_requests": fixture.requests})
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(directory, "wrapper-"+name+".json"), artifact, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestLargeSavedPlan checks that JSON expansion never makes a saved plan unreadable.
// Failure cases are large unchanged reports and text enlarged by JSON/HTML escaping.
func TestLargeSavedPlan(t *testing.T) {
	for name, text := range map[string]string{
		"ascii":   strings.Repeat("x", 200),
		"escaped": strings.Repeat("<>&\x00\n", 20),
	} {
		t.Run(name, func(t *testing.T) {
			f, server := newFixture(t)
			defer server.Close()
			const count = 40000
			created := time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339)
			rows := make([]map[string]any, count)
			for index := range rows {
				rows[index] = map[string]any{"id": strconv.Itoa(index + 1), "author_id": "100", "created_at": created, "text": text}
			}
			f.pages = [][]map[string]any{rows}
			path := filepath.Join(t.TempDir(), "large-plan.json")
			output, err := runCLI(t, server, "", "plan", "--evaluation-timeout", "60s", "--config", writePolicy(t, `force_delete_post("7")`), "--out", path)
			if err != nil {
				t.Fatalf("save large plan: %v: %s", err, output)
			}
			info, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			if info.Size() <= 16<<20 {
				t.Fatalf("fixture does not exercise JSON expansion: %d bytes", info.Size())
			}
			payload := readPlanDocument(t, path)["plan"].(map[string]any)
			if len(payload["cleanup_report"].(map[string]any)["decisions"].([]any)) != count || len(payload["operations"].([]any)) != 1 {
				t.Fatal("saved plan lost decisions or changed selected operations")
			}
			output, err = runCLI(t, server, "", "apply", "--plan", path)
			if err != nil || !reflect.DeepEqual(f.mutations, []string{"/2/tweets/7"}) || f.reads != 0 {
				t.Fatalf("reload/apply large plan: %v: %s; mutations %v, reads %d", err, output, f.mutations, f.reads)
			}
			if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
				artifact, err := json.Marshal(map[string]any{"decisions": count, "document_bytes": info.Size(), "mutations": f.mutations, "reads": f.reads})
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(directory, fmt.Sprintf("large-plan-%s.json", name)), artifact, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
