package main

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"google.golang.org/protobuf/encoding/protojson"
)

const fixtureSecret = "synthetic command cleanup fixture"

// Run as the command's child process, so the assertions exercise the real
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

func TestExecuteInjectorCleanup(t *testing.T) {
	for _, mode := range []string{"success", "startup_failure", "cancel"} {
		t.Run(mode, func(t *testing.T) {
			injector, ok := bazel.FindBinary("tools/vault/injector", "injector")
			if !ok {
				t.Fatal("injector binary missing from runfiles")
			}
			injector, err := filepath.Abs(injector)
			if err != nil {
				t.Fatal(err)
			}
			self, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			scratch, secrets := t.TempDir(), t.TempDir()
			receipt := filepath.Join(scratch, "consumer-receipt")
			value := "{{ index .Last.Files 0 }}"
			if mode == "startup_failure" {
				value = "{{"
			}
			data, err := al.ToPbJson(map[string]any{"res": []any{
				map[string]any{"name": "fixture", "file": map[string]any{"value": fixtureSecret}},
				map[string]any{"name": "AL_TEST_SECRET_PATH", "deps": []any{"fixture"}, "env": map[string]any{"value": value}},
			}}).Get()
			if err != nil {
				t.Fatal(err)
			}
			config, err := protojson.Marshal(&al_proto.Config{Plugins: []*al_proto.PluginConfig{{
				Name: "fixture", Bin: injector, Data: data,
				Labels: map[string]string{"test": "cleanup"}, Env: map[string]string{"TMPDIR": secrets},
			}}})
			if err != nil {
				t.Fatal(err)
			}
			configPath := filepath.Join(scratch, "config.json")
			if err := os.WriteFile(configPath, config, 0o600); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			var logs bytes.Buffer
			cmdCtx := al.NewCmdCtx(ctx, "")
			defer cmdCtx.RequestShutdown()
			cmdCtx.Logger = log.New(&logs, "", 0)
			cmdCtx.Stdout, cmdCtx.Stderr = &logs, &logs
			cmdCtx.Args = []string{"al", "run", "--config", configPath, "--plugin_label", "test=cleanup", "--", self, "-test.run=^TestSecretConsumer$", "--", "al-secret-consumer", receipt, mode}
			result := make(chan int, 1)
			go func() { result <- Execute(cmdCtx) }()
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
		})
	}
}
