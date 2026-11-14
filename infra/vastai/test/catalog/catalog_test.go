package catalog

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// terraformPath locates the packaged offline Terraform launcher.
var terraformPath = flag.String("terraform", "", "Terraform launcher runfile")

// configurationPath locates the actual catalog root and its adjacent inputs.
var configurationPath = flag.String("configuration", "", "Catalog root runfile")

// TestCatalog exercises the account, team and rental roots without live operations.
func TestCatalog(t *testing.T) {
	terraform, err := runfiles.Rlocation(*terraformPath)
	if err != nil {
		t.Fatal(err)
	}
	configuration, err := runfiles.Rlocation(*configurationPath)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	directory := filepath.Join(root, "infra", "vastai", "tf")
	teamDirectory := filepath.Join(root, "infra", "vastai", "alwaldend")
	rentalDirectory := filepath.Join(root, "infra", "vastai", "rentals_test")
	owner := filepath.Dir(filepath.Dir(configuration))
	inputs := map[string]string{
		"infra/vastai/tf/main.tf":                      configuration,
		"infra/vastai/config.json":                     filepath.Join(owner, "config.json"),
		"infra/vastai/tf/tests/catalog.tftest.hcl":     filepath.Join(filepath.Dir(configuration), "tests", "catalog.tftest.hcl"),
		"infra/vastai/alwaldend/main.tf":               filepath.Join(owner, "alwaldend", "main.tf"),
		"infra/vastai/alwaldend/tests/keys.tftest.hcl": filepath.Join(owner, "alwaldend", "tests", "keys.tftest.hcl"),
		"infra/vastai/rentals_test/main.tf":            filepath.Join(owner, "rentals_test", "main.tf"),
	}
	for relative, source := range inputs {
		data, err := os.ReadFile(source)
		if err != nil {
			t.Fatalf("read declared fixture input %s: %v", relative, err)
		}
		destination := filepath.Join(root, relative)
		if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(destination, data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	run := func(success bool, args ...string) string {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, terraform, append([]string{"--chdir", directory}, args...)...)
		cmd.WaitDelay = 5 * time.Second
		cmd.Dir = directory
		cmd.Env = append(os.Environ(), "TMPDIR=.", "VAST_API_KEY=fixture_issuer", "VAULT_TOKEN=fixture_component_token", "VAULT_ADDR=http://127.0.0.1:1")
		output, err := cmd.CombinedOutput()
		if (err == nil) != success {
			t.Fatalf("Terraform %v: %v\n%s", args, err, output)
		}
		return string(output)
	}
	run(true, "init", "-backend=false", "-no-color")
	run(true, "test", "-no-color")
	// Exercise metadata rotation and missing/invalid metadata through the real CLI.
	fixturePath := filepath.Join(directory, "tests", "catalog.tftest.hcl")
	originalFixture, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}
	field := regexp.MustCompile(`data_json\s*=\s*"(?:[^"\\]|\\.)*"`)
	for _, metadata := range []string{
		`{"metadata":{"sshpubkey":"ssh-rsa Zml4dHVyZS1yb3RhdGlvbg== rotated-fixture"}}`,
		`{"metadata":{}}`,
		`{"metadata":{"sshpubkey":"invalid-key-material"}}`,
	} {
		encoded, err := json.Marshal(metadata)
		if err != nil {
			t.Fatal(err)
		}
		changed := field.ReplaceAll(originalFixture, append([]byte("data_json = "), encoded...))
		if string(changed) == string(originalFixture) {
			t.Fatal("metadata fixture was not changed")
		}
		if err := os.WriteFile(fixturePath, changed, 0o600); err != nil {
			t.Fatal(err)
		}
		valid := strings.Contains(metadata, "rotated-fixture")
		output := run(valid, "test", "-no-color")
		if !valid && !strings.Contains(output, "must contain one OpenSSH public key") {
			t.Fatal("missing or invalid user metadata bypassed the key validation")
		}
	}
	if err := os.WriteFile(fixturePath, originalFixture, 0o600); err != nil {
		t.Fatal(err)
	}

	directory = teamDirectory
	run(true, "init", "-backend=false", "-no-color")
	run(true, "test", "-no-color")
	directory = rentalDirectory
	run(true, "init", "-backend=false", "-no-color")
	run(true, "validate", "-no-color")
	configurationFile := filepath.Join(root, "infra", "vastai", "config.json")
	data, err := os.ReadFile(configurationFile)
	if err != nil {
		t.Fatal(err)
	}
	var catalog map[string]any
	if err := json.Unmarshal(data, &catalog); err != nil {
		t.Fatal(err)
	}
	catalog["environment_variables"].(map[string]any)["WEB_ENABLE_AUTH"] = "false"
	changed, err := json.Marshal(catalog)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configurationFile, changed, 0o600); err != nil {
		t.Fatal(err)
	}
	directory = filepath.Join(root, "infra", "vastai", "tf")
	run(true, "test", "-no-color")
	directory = teamDirectory
	run(true, "test", "-no-color")
	catalog["teams"].(map[string]any)["alwaldend"].(map[string]any)["api_keys"] = []string{"../outside"}
	invalid, err := json.Marshal(catalog)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configurationFile, invalid, 0o600); err != nil {
		t.Fatal(err)
	}
	if output := run(false, "test", "-no-color"); !strings.Contains(output, "Invalid team key name") {
		t.Fatalf("invalid key configuration bypassed its guard: %s", output)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "catalog-root.json"), []byte(`{"result":"passed","account_teams":1,"team_roots":1,"rental_roots_validated":1,"invalid_configurations_rejected":1,"ssh_metadata_rotation_checked":true,"invalid_ssh_metadata_rejected":2,"live_operations":0}`), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
