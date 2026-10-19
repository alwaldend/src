package wasm

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// sqliteGuestRunfile selects the WASI SQLite command built with the pinned sysroot.
var sqliteGuestRunfile = flag.String("sqlite-wasm", "", "WASI SQLite probe runfile")

// sqliteNativeRunfile selects the native command that verifies guest-written storage.
var sqliteNativeRunfile = flag.String("sqlite-native", "", "native SQLite probe runfile")

// sqliteGuest runs the real command in a fresh pure-Go WASI instance.
func sqliteGuest(t *testing.T, action, path string) (string, error) {
	t.Helper()
	guestPath, err := runfiles.Rlocation(*sqliteGuestRunfile)
	if err != nil {
		return "", fmt.Errorf("resolve SQLite WASI command: %w", err)
	}
	code, err := os.ReadFile(guestPath)
	if err != nil {
		return "", fmt.Errorf("read SQLite WASI command: %w", err)
	}
	runtime := wazero.NewRuntimeWithConfig(t.Context(), wazero.NewRuntimeConfig().WithCloseOnContextDone(true))
	t.Cleanup(func() {
		if err := runtime.Close(context.Background()); err != nil {
			t.Errorf("close SQLite WASI runtime: %v", err)
		}
	})
	if _, err := wasi_snapshot_preview1.Instantiate(t.Context(), runtime); err != nil {
		return "", fmt.Errorf("initialize SQLite WASI host: %w", err)
	}
	var output bytes.Buffer
	config := wazero.NewModuleConfig().WithArgs("sqlite-probe", action, path).
		WithFSConfig(wazero.NewFSConfig().WithDirMount("/", "/")).
		WithSysWalltime().WithSysNanotime().WithStdout(&output).WithStderr(&output)
	if _, err := runtime.InstantiateWithConfig(t.Context(), code, config); err != nil {
		return output.String(), fmt.Errorf("run SQLite WASI %q: %w", action, err)
	}
	return output.String(), nil
}

// TestSQLiteRoundTrip verifies durable guest writes with an independent native opener.
func TestSQLiteRoundTrip(t *testing.T) {
	directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if directory == "" {
		directory = t.TempDir()
	}
	path := filepath.Join(directory, "wasi-sqlite.anki2")
	if output, err := sqliteGuest(t, "write", path); err != nil || output != "sqlite-roundtrip-ok\n" {
		t.Fatalf("WASI database write: %v; %s", err, output)
	}
	info, err := os.Stat(path)
	if err != nil || info.Size() <= 10*1024*1024 {
		t.Fatalf("guest database did not persist more than 10 MiB: %v; %v", info, err)
	}
	nativePath, err := runfiles.Rlocation(*sqliteNativeRunfile)
	if err != nil {
		t.Fatalf("resolve native SQLite verifier: %v", err)
	}
	output, err := exec.CommandContext(t.Context(), nativePath, "verify", path).CombinedOutput()
	if err != nil || string(output) != "sqlite-roundtrip-ok\n" {
		t.Fatalf("native database reopen: %v; %s", err, output)
	}
	if err := os.WriteFile(filepath.Join(directory, "sqlite-verification.txt"), output, 0o600); err != nil {
		t.Fatalf("retain native verification artifact: %v", err)
	}
	if output, err := sqliteGuest(t, "verify", path); err != nil || output != "sqlite-roundtrip-ok\n" {
		t.Fatalf("second WASI database reopen: %v; %s", err, output)
	}
	if _, err := sqliteGuest(t, "write", filepath.Join(t.TempDir(), "missing", "db.anki2")); err == nil {
		t.Fatal("WASI SQLite accepted a nonexistent parent directory")
	}
}
