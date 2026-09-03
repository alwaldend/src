package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/hashicorp/vault/api"
)

func TestVaultDevServerWritableFixture(t *testing.T) {
	vaultPath, ok := bazel.FindBinary("tools/vault/injector", "vault_fixture")
	if !ok {
		t.Fatal("could not find Vault binary")
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	address := fmt.Sprintf("http://%s", listener.Addr())
	if err := listener.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, vaultPath, "server", "-dev", "-dev-root-token-id=fixture", "-dev-listen-address="+listener.Addr().String())
	cmd.Env = append(
		os.Environ(),
		"HOME="+t.TempDir(),
		"XDG_CACHE_HOME="+t.TempDir(),
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	var stderr bytes.Buffer
	var stdoutBuffer bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start Vault: %v", err)
	}
	defer func() {
		cancel()
		_ = cmd.Wait()
	}()

	ready := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			stdoutBuffer.Write(scanner.Bytes())
			stdoutBuffer.WriteByte('\n')
			if bytes.Contains(scanner.Bytes(), []byte("Unseal Key")) {
				ready <- nil
				return
			}
		}
		ready <- fmt.Errorf("Vault exited before readiness: %w", scanner.Err())
	}()
	select {
	case err := <-ready:
		if err != nil {
			t.Fatalf("stdout:\n%s\nstderr:\n%s", stdoutBuffer.String(), stderr.String())
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("timed out waiting for Vault readiness; output:\n%s", stderr.String())
	}

	client, err := api.NewClient(&api.Config{Address: address})
	if err != nil {
		t.Fatalf("create Vault client: %v", err)
	}
	client.SetToken("fixture")
	if _, err := client.Logical().Write("sys/mounts/secrets", map[string]any{
		"type":    "kv",
		"options": map[string]any{"version": "2"},
	}); err != nil {
		t.Fatalf("enable secrets mount: %v", err)
	}
	mount := client.KVv2("secrets")
	if _, err := mount.Put(ctx, "injector/fixture", map[string]any{
		"value": "writable",
	}); err != nil {
		t.Fatalf("put secret: %v", err)
	}
	secret, err := mount.Get(ctx, "injector/fixture")
	if err != nil {
		t.Fatalf("get secret: %v", err)
	}
	if secret.Data["value"] != "writable" {
		t.Fatalf("secret value = %#v, want writable", secret.Data["value"])
	}
	if _, err := client.Logical().Delete("sys/mounts/secrets"); err != nil {
		t.Fatalf("delete mount: %v", err)
	}
}
