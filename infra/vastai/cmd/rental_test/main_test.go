package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// commandPath resolves the actual command under test.
var commandPath = flag.String("command", "", "Rental command runfile")

// launcherPath resolves the real packaged launcher used by the cancellation fixture.
var launcherPath = flag.String("launcher", "", "Terraform launcher runfile")

// TestMain makes the test executable a subprocess fixture for the packaged tools.
func TestMain(m *testing.M) {
	if os.Getenv("VASTAI_RENTAL_FIXTURE") == "1" && (filepath.Base(os.Args[0]) == "terraform" || filepath.Base(os.Args[0]) == "vastai") {
		os.Exit(toolFixture())
	}
	os.Exit(m.Run())
}

// toolFixture simulates only the commands used by the interactive lifecycle test.
func toolFixture() int {
	arguments := os.Args[1:]
	if len(arguments) > 0 && strings.HasPrefix(arguments[0], "-chdir=") {
		arguments = arguments[1:]
	}
	operation := strings.Join(arguments, " ")
	path := os.Getenv("VASTAI_RENTAL_CALLS")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return 2
	}
	_, writeErr := fmt.Fprintln(f, operation)
	if closeErr := f.Close(); writeErr != nil || closeErr != nil {
		return 2
	}
	mode, marker := os.Getenv("VASTAI_RENTAL_MODE"), os.Getenv("VASTAI_RENTAL_MARKER")
	if operation == "init" {
		return 0
	}
	if operation == "version -json" {
		fmt.Fprint(os.Stdout, `{"provider_selections":{}}`)
		return 0
	}
	if filepath.Base(os.Args[0]) == "vastai" {
		if strings.Contains(operation, "show instances") {
			if mode == "inventory-failure" {
				fmt.Fprint(os.Stdout, `null`)
			} else {
				fmt.Fprint(os.Stdout, `[{"id":8}]`)
			}
			return 0
		}
		id := 42
		name := "alwaldend"
		if mode == "wrong-team" {
			id = 99
			name = "other_team"
		}
		fmt.Fprintf(os.Stdout, `{"id":%d,"username":%q,"is_team":%t,"credit":1,"private_fixture_field":"must-not-leak"}`, id, name, mode != "personal-account")
		return 0
	}
	if strings.HasPrefix(operation, "state list") {
		if mode == "state-read-failure" {
			fmt.Fprintln(os.Stderr, "Backend unavailable")
			return 1
		}
		if _, err := os.Stat(marker); err == nil {
			fmt.Fprintln(os.Stdout, "vastai_instance.worker")
		} else if mode == "empty-backend" {
			fmt.Fprintln(os.Stderr, "No state file was found!")
			return 1
		}
		return 0
	}
	if strings.Contains(operation, "-auto-approve") || strings.Contains(operation, "-out") {
		return 2
	}
	if strings.Contains(operation, "-destroy") {
		if mode == "cleanup-failure" {
			return 1
		}
		if err := os.Remove(marker); err != nil && !os.IsNotExist(err) {
			return 2
		}
		fmt.Fprintln(os.Stdout, "Apply complete! Resources: 0 added, 0 changed, 1 destroyed.")
		return 0
	}
	if !strings.HasPrefix(operation, "apply") {
		return 2
	}
	if err := os.WriteFile(marker, []byte("91"), 0o600); err != nil {
		return 2
	}
	if mode == "cancel" {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		if err := os.WriteFile(os.Getenv("VASTAI_RENTAL_READY"), []byte("ready"), 0o600); err != nil {
			return 2
		}
		<-ctx.Done()
		if err := os.WriteFile(os.Getenv("VASTAI_RENTAL_CANCELLED"), []byte("interrupted"), 0o600); err != nil {
			return 2
		}
		return 1
	}
	if mode == "creation-failure" {
		return 1
	}
	fmt.Fprintln(os.Stdout, "rental_contract_id = \"91\"\nrental_hourly_price = 0.001")
	return 0
}

// TestCancellationThroughLauncher verifies SIGTERM reaches Terraform before release starts.
func TestCancellationThroughLauncher(t *testing.T) {
	dir := t.TempDir()
	binary, err := runfiles.Rlocation(*commandPath)
	if err != nil {
		t.Fatal(err)
	}
	launcher, err := runfiles.Rlocation(*launcherPath)
	if err != nil {
		t.Fatal(err)
	}
	helper, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"terraform", "vastai"} {
		if err := os.Symlink(helper, filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
	wrapped := filepath.Join(dir, "wrapped-terraform")
	if err := os.Symlink(launcher, wrapped); err != nil {
		t.Fatal(err)
	}
	fixtureDirectory := filepath.Join(dir, "runfiles", "_main", "fixture")
	if err := os.MkdirAll(fixtureDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	configuration, err := json.Marshal(map[string]any{"terraform": filepath.Join(dir, "terraform"), "working_directory": "_main/fixture"})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wrapped+".json", configuration, 0o600); err != nil {
		t.Fatal(err)
	}
	receiptPath, ready, cancelled := filepath.Join(dir, "receipt.json"), filepath.Join(dir, "ready"), filepath.Join(dir, "cancelled")
	cmd := exec.Command(binary, "--terraform", wrapped, "--vastai", filepath.Join(dir, "vastai"), "--receipt", receiptPath)
	cmd.Env = append(os.Environ(), "RUNFILES_DIR="+filepath.Join(dir, "runfiles"), "TMPDIR="+dir, "VASTAI_RENTAL_FIXTURE=1", "VASTAI_RENTAL_MODE=cancel", "VASTAI_RENTAL_MARKER="+filepath.Join(dir, "contract"), "VASTAI_RENTAL_CALLS="+filepath.Join(dir, "calls"), "VASTAI_RENTAL_READY="+ready, "VASTAI_RENTAL_CANCELLED="+cancelled)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	deadline := time.NewTimer(8 * time.Second)
	defer deadline.Stop()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		select {
		case err := <-done:
			t.Fatalf("command exited before creation: %v\n%s", err, output.String())
		case <-deadline.C:
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
			t.Fatalf("creation did not start\n%s", output.String())
		case <-ticker.C:
		}
	}
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-deadline.C:
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-done
		t.Fatalf("cancellation left a descendant holding the output pipe\n%s", output.String())
	}
	if _, err := os.Stat(cancelled); err != nil {
		t.Fatal("Terraform child never received graceful cancellation")
	}
	data, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if result["cleanup_attempted"] != true || result["state_empty"] != true {
		t.Fatalf("interrupted creation did not release its contract: %s", data)
	}
	if outputDir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); outputDir != "" {
		if err := os.WriteFile(filepath.Join(outputDir, "rental-command-cancellation.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestCommandLifecycle verifies cleanup ordering, partial failure, and state isolation.
func TestCommandLifecycle(t *testing.T) {
	for _, mode := range []string{"success", "empty-backend", "state-read-failure", "inventory-failure", "personal-account", "creation-failure", "cleanup-failure", "existing-rental", "wrong-team"} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			binary, err := runfiles.Rlocation(*commandPath)
			if err != nil {
				t.Fatal(err)
			}
			helper, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{"terraform", "vastai"} {
				if err := os.Symlink(helper, filepath.Join(dir, name)); err != nil {
					t.Fatal(err)
				}
			}
			marker, callsPath, receiptPath := filepath.Join(dir, "contract"), filepath.Join(dir, "calls"), filepath.Join(dir, "receipt.json")
			if mode == "existing-rental" {
				if err := os.WriteFile(marker, []byte("prior-contract"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			cmd := exec.Command(binary, "--terraform", filepath.Join(dir, "terraform"), "--vastai", filepath.Join(dir, "vastai"), "--receipt", receiptPath)
			cmd.Env = append(os.Environ(), "VASTAI_RENTAL_FIXTURE=1", "VASTAI_RENTAL_MODE="+mode, "VASTAI_RENTAL_MARKER="+marker, "VASTAI_RENTAL_CALLS="+callsPath)
			output, err := cmd.CombinedOutput()
			success := mode == "success" || mode == "empty-backend"
			if (err == nil) != success {
				t.Fatalf("unexpected result: %v\n%s", err, output)
			}
			if strings.Contains(string(output), "must-not-leak") {
				t.Fatal("account response leaked")
			}
			calls, err := os.ReadFile(callsPath)
			if err != nil {
				t.Fatal(err)
			}
			operations := strings.Split(strings.TrimSpace(string(calls)), "\n")
			rentals, cleanup := 0, 0
			for i, call := range operations {
				if strings.HasPrefix(call, "apply") && !strings.Contains(call, "-destroy") {
					rentals++
					if i+1 >= len(operations) || !strings.Contains(operations[i+1], "-destroy") {
						t.Fatal("release was not the immediate next operation")
					}
				}
				if strings.Contains(call, "-destroy") {
					cleanup++
				}
			}
			if mode == "existing-rental" || mode == "wrong-team" || mode == "personal-account" || mode == "inventory-failure" || mode == "state-read-failure" {
				if rentals != 0 || cleanup != 0 {
					t.Fatal("preflight failure mutated rental state")
				}
			} else if rentals != 1 || cleanup != 1 {
				t.Fatal("test did not attempt exactly one rental and release")
			}
			data, err := os.ReadFile(receiptPath)
			if err != nil {
				t.Fatal(err)
			}
			var result map[string]any
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(data), "must-not-leak") || result["success"] != success {
				t.Fatal("invalid or unsanitized receipt")
			}
			if success && (result["contract_id"] != "91" || result["state_empty"] != true) {
				t.Fatal("receipt did not prove the contract lifecycle")
			}
			if outputDir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); outputDir != "" {
				if err := os.WriteFile(filepath.Join(outputDir, "rental-command-"+mode+".json"), data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
