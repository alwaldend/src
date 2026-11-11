package evaluation

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// TestConfiguredInputs compiles the actual public run matrix without API credentials.
func TestConfiguredInputs(t *testing.T) {
	for _, group := range []string{"", "plush"} {
		t.Run("matrix_"+group, func(t *testing.T) { checkConfiguredInputs(t, group) })
	}
}

// checkConfiguredInputs exercises either the complete matrix or its GIF-only group.
func checkConfiguredInputs(t *testing.T, group string) {
	t.Helper()
	binary, err := runfiles.Rlocation(os.Getenv("EVALUATOR_BINARY"))
	if err != nil {
		t.Fatalf("resolve configuration checker: %v", err)
	}
	directory := t.TempDir()
	var arguments []string
	if err := json.Unmarshal([]byte(os.Getenv("EVALUATION_ARGUMENTS")), &arguments); err != nil {
		t.Fatalf("decode configured argument projection: %v", err)
	}
	arguments = append([]string{"check"}, arguments...)
	arguments = append(arguments, "--output", directory)
	if group != "" {
		var selection []string
		if err := json.Unmarshal([]byte(os.Getenv("EVALUATION_GIF_ARGUMENTS")), &selection); err != nil {
			t.Fatalf("decode GIF argument projection: %v", err)
		}
		arguments = append(arguments, selection...)
	}
	command := exec.Command(binary, arguments...)
	command.Env = append(os.Environ(), "VAST_API_KEY=", "VAULT_TOKEN=")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("check configured evals: %v\n%s", err, output)
	}
	data, err := os.ReadFile(filepath.Join(directory, "configuration-check.json"))
	if err != nil {
		t.Fatalf("read configured-input artifact: %v", err)
	}
	var receipt struct {
		// Rentals is the real matrix pooled by installed model and hardware.
		Rentals []json.RawMessage `json:"rentals"`
		// Workloads retains all expanded workflow/case coordinates.
		Workloads []json.RawMessage `json:"workloads"`
		// LiveOperations must stay zero for source acceptance.
		LiveOperations int `json:"live_operations"`
	}
	if err := json.Unmarshal(data, &receipt); err != nil {
		t.Fatalf("decode configuration artifact: %v", err)
	}
	if len(receipt.Rentals) == 0 || len(receipt.Workloads) < len(receipt.Rentals) || receipt.LiveOperations != 0 {
		t.Fatal("configured-input acceptance did not cover the matrix offline")
	}
	if group != "" {
		for _, entry := range receipt.Workloads {
			var workload map[string]any
			if err := json.Unmarshal(entry, &workload); err != nil {
				t.Fatalf("decode selected workload: %v", err)
			}
			if workload["case_id"] != "plush_petting" {
				t.Fatal("GIF selection included a still-image case")
			}
		}
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		filename := "configuration-check.json"
		if group != "" {
			filename = "gif-configuration-check.json"
		}
		if err := os.WriteFile(filepath.Join(output, filename), data, 0o600); err != nil {
			t.Fatalf("preserve configuration artifact: %v", err)
		}
	}
}
