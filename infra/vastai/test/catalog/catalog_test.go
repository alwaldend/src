package catalog

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// terraformPath locates the packaged offline Terraform launcher.
var terraformPath = flag.String("terraform", "", "Terraform launcher runfile")

// configurationPath locates the actual catalog root and its adjacent inputs.
var configurationPath = flag.String("configuration", "", "Catalog root runfile")

// teamModulePath locates the shared team module used by the actual root.
var teamModulePath = flag.String("team-module", "", "Team module runfile")

// TestCatalog exercises actual configuration with distinct mock Vault identities.
func TestCatalog(t *testing.T) {
	terraform, err := runfiles.Rlocation(*terraformPath)
	if err != nil {
		t.Fatal(err)
	}
	configuration, err := runfiles.Rlocation(*configurationPath)
	if err != nil {
		t.Fatal(err)
	}
	teamModule, err := runfiles.Rlocation(*teamModulePath)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	directory := filepath.Join(root, "infra", "vastai", "tf")
	inputs := map[string]string{
		"infra/vastai/tf/main.tf":                  configuration,
		"infra/vastai/tf/teams.json":               filepath.Join(filepath.Dir(configuration), "teams.json"),
		"infra/vastai/tf/tests/catalog.tftest.hcl": filepath.Join(filepath.Dir(configuration), "tests", "catalog.tftest.hcl"),
		"projects/tf_modules/vastai_team/main.tf":  teamModule,
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
		cmd.Env = append(os.Environ(), "TMPDIR=.", "VAST_API_KEY=fixture_issuer", "VAULT_TOKEN=fixture_project_token", "TF_VAR_component_vault_token=fixture_component_token", "VAULT_ADDR=http://127.0.0.1:1")
		output, err := cmd.CombinedOutput()
		if (err == nil) != success {
			t.Fatalf("Terraform %v: %v\n%s", args, err, output)
		}
		return string(output)
	}
	run(true, "init", "-backend=false", "-no-color")
	run(true, "test", "-no-color")
	data, err := os.ReadFile(filepath.Join(directory, "teams.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"duplicate-destination", "invalid-identity"} {
		var teams map[string]any
		if err := json.Unmarshal(data, &teams); err != nil {
			t.Fatal(err)
		}
		key := teams["alwaldend"].(map[string]any)["api_keys"].(map[string]any)["rentals"].(map[string]any)
		if mode == "duplicate-destination" {
			key["secret_name"] = "alwaldend.com/vault1/approles/src_projects_vastai_terraform_provider/vastai"
		} else {
			key["vault_auth"] = "unsupported_identity"
		}
		invalid, err := json.Marshal(teams)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, "teams.json"), invalid, 0o600); err != nil {
			t.Fatal(err)
		}
		if output := run(false, "test", "-no-color"); !strings.Contains(output, "Resource precondition failed") {
			t.Fatalf("%s failed without the catalog guard: %s", mode, output)
		}
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "catalog-root.json"), []byte(`{"result":"passed","provider_aliases_verified":2,"catalog_key":"alwaldend/rentals","invalid_configurations_rejected":2,"live_operations":0}`), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
