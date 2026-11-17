package main_test

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// TestBootstrap exercises the packaged command against an isolated TLS Vault.
func TestBootstrap(t *testing.T) {
	var record map[string]any
	writes := 0
	reads := 0
	failure := ""
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "fixture-vault-token" {
			t.Error("missing fixture authentication")
			w.WriteHeader(403)
			return
		}
		if failure == "denied" {
			w.WriteHeader(403)
			return
		}
		if failure == "redirect" {
			w.Header().Set("Location", "/redirect-target")
			w.WriteHeader(307)
			return
		}
		if r.Method == "GET" {
			reads++
			if strings.Contains(r.URL.Path, "/metadata/") {
				if failure == "history" {
					fmt.Fprint(w, `{"data":{"current_version":1}}`)
					return
				}
				w.WriteHeader(404)
				return
			}
			if record == nil {
				w.WriteHeader(404)
				return
			}
			if err := json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"data": record, "metadata": map[string]any{"version": 1}}}); err != nil {
				t.Error(err)
			}
			return
		}
		if r.Method != "PUT" && r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
		var payload struct {
			// Data carries the candidate pair only inside this private fixture.
			Data map[string]any `json:"data"`
			// Options enforces the immutable Vault record contract.
			Options map[string]int `json:"options"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Error(err)
			return
		}
		cas, present := payload.Options["cas"]
		if !present || cas != 0 || record != nil || failure == "conflict" {
			w.WriteHeader(400)
			return
		}
		record = payload.Data
		writes++
		fmt.Fprint(w, `{"data":{"version":1}}`)
	}))
	defer server.Close()
	directory := t.TempDir()
	ca := filepath.Join(directory, "ca.pem")
	if err := os.WriteFile(ca, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw}), 0o600); err != nil {
		t.Fatal(err)
	}
	tool, err := runfiles.Rlocation(os.Getenv("SSH_KEY_TOOL"))
	if err != nil {
		t.Fatal(err)
	}
	environment := []string{"VAULT_ADDR=" + server.URL, "VAULT_TOKEN=fixture-vault-token", "VAULT_CACERT=" + ca, "VAULT_SSH_KEY_MOUNT=secrets", "VAULT_SSH_KEY_PATH=fixture/signing"}
	for _, variable := range os.Environ() {
		if strings.HasPrefix(variable, "RUNFILES_") || strings.HasPrefix(variable, "TEST_") || strings.HasPrefix(variable, "PATH=") {
			environment = append(environment, variable)
		}
	}
	invoke := func(success bool) string {
		t.Helper()
		command := exec.Command(tool, "ensure")
		command.Env = environment
		output, err := command.CombinedOutput()
		for _, secret := range []string{"fixture-vault-token", "BEGIN OPENSSH PRIVATE KEY"} {
			if strings.Contains(string(output), secret) {
				t.Fatal("credential leaked in output")
			}
		}
		if (err == nil) != success {
			t.Fatalf("unexpected command result: %v: %s", err, output)
		}
		return string(output)
	}
	for _, mode := range []string{"denied", "redirect", "history", "conflict"} {
		failure = mode
		invoke(false)
		if writes != 0 {
			t.Fatal("failed prerequisite wrote a key")
		}
	}
	failure = ""
	invoke(true)
	if writes != 1 || record["private_key"] == "" || record["public_key"] == "" {
		t.Fatal("pair not persisted")
	}
	private := record["private_key"]
	invoke(true)
	if writes != 1 || record["private_key"] != private {
		t.Fatal("cache rotated")
	}
	for _, field := range []string{"private_key", "public_key", "fingerprint", "algorithm"} {
		original := record[field]
		record[field] = "invalid"
		invoke(false)
		record[field] = original
	}
	// The real Terraform declaration must have no plan effects, run once on
	// creation, reuse Vault on replacement, and omit private values from state.
	source, err := runfiles.Rlocation(os.Getenv("BOOTSTRAP_SOURCE"))
	if err != nil {
		t.Fatal(err)
	}
	declaration, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "bootstrap.tf"), declaration, 0o600); err != nil {
		t.Fatal(err)
	}
	configuration := map[string]any{"locals": map[string]any{"signing_key_tool": tool, "signing_key_mount": "secrets", "signing_key_path": "fixture/signing"}}
	raw, err := json.Marshal(configuration)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "fixture.tf.json"), raw, 0o600); err != nil {
		t.Fatal(err)
	}
	terraform, err := runfiles.Rlocation(os.Getenv("TERRAFORM"))
	if err != nil {
		t.Fatal(err)
	}
	record = nil
	for index, args := range [][]string{{"init", "-backend=false"}, {"plan"}, {"apply"}, {"apply"}, {"apply", "-replace=terraform_data.commit_signing"}} {
		previousReads, previousWrites := reads, writes
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		command := exec.CommandContext(ctx, terraform, append([]string{"--chdir", directory}, args...)...)
		command.Env = append(environment, "HOME="+directory, "TF_IN_AUTOMATION=1")
		command.Stdin = strings.NewReader("yes\n")
		output, err := command.CombinedOutput()
		cancel()
		if strings.Contains(string(output), "OPENSSH PRIVATE KEY") {
			t.Fatal("key leaked in Terraform output")
		}
		if strings.Contains(string(output), "fixture-vault-token") {
			t.Fatal("token leaked in Terraform output")
		}
		if err != nil {
			t.Fatalf("Terraform %s failed: %v: %s", args[0], err, output)
		}
		if (index < 2 || index == 3) && (reads != previousReads || writes != previousWrites) {
			t.Fatal("plan or repeat apply invoked bootstrap")
		}
		if index == 2 && writes != previousWrites+1 {
			t.Fatal("initial apply did not create a key")
		}
		if index == 4 && (reads != previousReads+1 || writes != previousWrites) {
			t.Fatal("replacement did not reuse cache")
		}
	}
	state, err := os.ReadFile(filepath.Join(directory, "terraform.tfstate"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(state), "OPENSSH PRIVATE KEY") || strings.Contains(string(state), "fixture-vault-token") {
		t.Fatal("key leaked in Terraform state")
	}
	if writes != 2 {
		t.Fatal("Terraform rotated the cached key")
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "signing-bootstrap.json"), []byte(`{"cli":true,"cache":true,"failures":true,"terraform":true,"secret_exclusion":true}`), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
