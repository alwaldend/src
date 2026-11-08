package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/protobuf/encoding/protojson"
)

// TestDictionaryChild is the real AL child process for the dictionary fixture.
func TestDictionaryChild(t *testing.T) {
	if os.Getenv("AL_DICTIONARY_CHILD") != "1" {
		return
	}
	values := map[string]string{"key": os.Getenv("FIXTURE_API_KEY"), "last": os.Getenv("FIXTURE_LAST"), "vault_token": os.Getenv("VAULT_TOKEN")}
	if err := json.NewEncoder(os.Stdout).Encode(values); err != nil {
		t.Fatal(err)
	}
	os.Exit(0)
}

// TestDictionaryConfiguration runs generated configuration through the actual
// AL executable, injector RPC process, resource graph and child environment.
func TestDictionaryConfiguration(t *testing.T) {
	runtimeFiles, err := runfiles.New()
	if err != nil {
		t.Fatal(err)
	}
	configPath, err := runtimeFiles.Rlocation(os.Getenv("DICTIONARY_CONFIG"))
	if err != nil {
		t.Fatal(err)
	}
	binary, err := runtimeFiles.Rlocation(os.Getenv("AL_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.Header.Get("X-Vault-Token") != "" {
			http.Error(w, "unexpected fixture authentication or method", http.StatusBadRequest)
			return
		}
		values := map[string]any{}
		switch r.URL.Path {
		case "/v1/secrets/data/fixture/credentials":
			values["api_key"] = "synthetic-key"
		case "/v1/secrets/data/fixture/credentials_override":
			values["api_key"] = "synthetic-override"
		case "/v1/secrets/data/fixture/alpha":
			values["value"] = "alpha"
		case "/v1/secrets/data/fixture/omega":
			values["value"] = "omega"
		default:
			http.Error(w, "unselected or unknown fixture resource", http.StatusBadRequest)
			return
		}
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"data": values, "metadata": map[string]any{"version": 1}}}); err != nil {
			return
		}
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	config, err := al.LoadConfigs(ctx, configPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(config.VaultConn) != 1 {
		t.Fatal("fixture lost its Vault connection")
	}
	config.VaultConn[0].Config.Address = server.URL
	content, err := protojson.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	selectedConfig := filepath.Join(directory, "config.json")
	if err := os.WriteFile(selectedConfig, content, 0o600); err != nil {
		t.Fatal(err)
	}
	child, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	command := exec.CommandContext(ctx, binary, "run", "--config", selectedConfig, "--plugin_label", "dictionary=selected", "--", child, "-test.run=^TestDictionaryChild$")
	command.Dir = directory
	command.Env = append([]string{"PATH=" + os.Getenv("PATH"), "HOME=" + directory, "TMPDIR=" + directory, "AL_DICTIONARY_CHILD=1"}, runtimeFiles.Env()...)
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil {
		t.Fatalf("run generated dictionary configuration: %v\n%s", err, stderr.String())
	}
	actual := map[string]string{}
	if err := json.Unmarshal(stdout.Bytes(), &actual); err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{"key": "synthetic-override", "last": "alpha", "vault_token": "synthetic-child-token"}
	if !reflect.DeepEqual(actual, expected) || requests.Load() != 3 {
		t.Fatalf("unexpected selected child environment or request count: %v, %d", actual, requests.Load())
	}
	for label, diagnostic := range map[string]string{"missing": "missing dependency", "cycle": "dependency cycle"} {
		failed := exec.CommandContext(ctx, binary, "run", "--config", selectedConfig, "--plugin_label", "dictionary="+label, "--", child, "-test.run=^TestDictionaryChild$")
		failed.Dir, failed.Env = directory, command.Env
		var output, errors bytes.Buffer
		failed.Stdout, failed.Stderr = &output, &errors
		if err := failed.Run(); err == nil || output.Len() != 0 || !strings.Contains(errors.String(), diagnostic) {
			t.Fatalf("%s dependency fixture did not fail before child startup with the expected diagnostic", label)
		}
	}
	output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if output == "" {
		t.Fatal("missing repeatable fixture output directory")
	}
	result, err := json.MarshalIndent(map[string]any{"verified": true, "requests": requests.Load(), "environment": actual}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "dictionary-results.json"), append(result, '\n'), 0o600); err != nil {
		t.Fatal(fmt.Errorf("write dictionary fixture result: %w", err))
	}
}
