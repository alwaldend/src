package command

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/protobuf/encoding/protojson"
)

// fixtureSecret is the expected public value in the generated cleanup fixture.
const fixtureSecret = "synthetic command cleanup fixture"

// TestSecretConsumer runs as the command's child process, exercising the real
// injector binary, RPC transport, environment injection, and process shutdown.
func TestSecretConsumer(t *testing.T) {
	if len(os.Args) < 4 || os.Args[len(os.Args)-3] != "al-secret-consumer" {
		return
	}
	receipt, mode := os.Args[len(os.Args)-2], os.Args[len(os.Args)-1]
	path := os.Getenv("AL_TEST_SECRET_PATH")
	check := func() {
		t.Helper()
		data, err := os.ReadFile(path)
		if err != nil || string(data) != fixtureSecret {
			t.Fatalf("injected file unavailable to consumer: %v", err)
		}
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != 0o600 {
			t.Fatalf("injected file permissions: %v, %v", info, err)
		}
	}
	check()
	signals := make(chan os.Signal, 1)
	if mode == "cancel" {
		signal.Notify(signals, syscall.SIGTERM)
		defer signal.Stop(signals)
	}
	if err := os.WriteFile(receipt, []byte(path), 0o600); err != nil {
		t.Fatal(err)
	}
	if mode == "cancel" {
		select {
		case <-signals:
			check()
			if err := os.WriteFile(receipt+".stopped", []byte("secret remained available during shutdown"), 0o600); err != nil {
				t.Fatal(err)
			}
		case <-time.After(15 * time.Second):
			t.Fatal("consumer did not receive SIGTERM")
		}
	}
}

// TestExecuteInjectorCleanup verifies generated configuration and transitive
// plugin runfiles through successful, failed and cancelled command lifecycles.
func TestExecuteInjectorCleanup(t *testing.T) {
	for _, mode := range []string{"success", "startup_failure", "cancel"} {
		t.Run(mode, func(t *testing.T) {
			configVariable := "CLEANUP_CONFIG"
			if mode == "startup_failure" {
				configVariable = "CLEANUP_FAILURE_CONFIG"
			}
			fixture, err := runfiles.Rlocation(os.Getenv(configVariable))
			if err != nil {
				t.Fatalf("resolve generated cleanup configuration: %v", err)
			}
			self, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			scratch, secrets := t.TempDir(), t.TempDir()
			receipt := filepath.Join(scratch, "consumer-receipt")
			config, err := al.LoadConfigs(context.Background(), fixture)
			if err != nil {
				t.Fatalf("load generated cleanup configuration: %v", err)
			}
			for _, plugin := range config.PluginConfigs {
				plugin.Env = map[string]string{"TMPDIR": secrets}
			}
			content, err := protojson.Marshal(config)
			if err != nil {
				t.Fatal(err)
			}
			configPath := filepath.Join(scratch, "config.json")
			if err := os.WriteFile(configPath, content, 0o600); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			var logs bytes.Buffer
			binary, err := runfiles.Rlocation(os.Getenv("AL_BINARY"))
			if err != nil {
				t.Fatal(err)
			}
			command := exec.CommandContext(ctx, binary, "run", "--config", configPath, "--plugin_label", "test=cleanup", "--", self, "-test.run=^TestSecretConsumer$", "--", "al-secret-consumer", receipt, mode)
			command.Env = append(os.Environ(), "TMPDIR="+secrets)
			command.Stdout, command.Stderr = &logs, &logs
			command.Cancel = func() error { return command.Process.Signal(syscall.SIGTERM) }
			command.WaitDelay = 15 * time.Second
			result := make(chan int, 1)
			go func() {
				code := 0
				if err := command.Run(); err != nil {
					code = 1
				}
				result <- code
			}()
			if mode == "cancel" {
				ticker := time.NewTicker(10 * time.Millisecond)
				defer ticker.Stop()
				ready := false
				for !ready {
					select {
					case code := <-result:
						t.Fatalf("command exited before consumer ready: %d: %s", code, &logs)
					case <-ctx.Done():
						t.Fatal("consumer did not become ready")
					case <-ticker.C:
						_, err := os.Stat(receipt)
						ready = err == nil
					}
				}
				cancel()
			}
			select {
			case code := <-result:
				if (mode == "success" && code != 0) || (mode != "success" && code == 0) {
					t.Fatalf("unexpected exit code %d: %s", code, &logs)
				}
			case <-time.After(30 * time.Second):
				t.Fatal("command did not finish cleanup")
			}
			if mode == "startup_failure" {
				if !strings.Contains(logs.String(), "could not parse template") {
					t.Fatalf("startup failed before the intended template failure: %s", &logs)
				}
				if _, err := os.Stat(receipt); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("consumer ran despite startup failure: %v", err)
				}
			} else {
				path, err := os.ReadFile(receipt)
				if err != nil {
					t.Fatalf("consumer did not verify injected file: %v", err)
				}
				if _, err := os.Stat(string(path)); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("injected secret survived command return: %v", err)
				}
				if mode == "cancel" {
					if _, err := os.Stat(receipt + ".stopped"); err != nil {
						t.Fatalf("consumer shutdown could not access injected file: %v", err)
					}
				}
			}
			remaining, err := os.ReadDir(secrets)
			if err != nil || len(remaining) != 0 {
				t.Fatalf("injector left temporary material: %v, %v", remaining, err)
			}
			output := filepath.Join(os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"), mode+".txt")
			if err := os.WriteFile(output, []byte("command lifecycle and temporary-file cleanup verified\n"), 0o600); err != nil {
				t.Fatalf("write cleanup verification artifact: %v", err)
			}
		})
	}
}
