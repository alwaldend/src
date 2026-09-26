package icon_generator_test

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestWriteFailureCleanup(t *testing.T) {
	dir := t.TempDir()
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, self, "-test.run=^TestFileSizeLimitHelper$")
	cmd.Dir = dir
	cmd.Env = []string{"PATH=", "ICON_GENERATOR_LIMIT_HELPER=" + executable(t)}
	output, err := cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(output), "encode") {
		t.Fatalf("expected encoding error: %v %s", err, output)
	}
	if _, err := os.Stat(filepath.Join(dir, "icon.png")); !os.IsNotExist(err) {
		t.Fatalf("incomplete output remains: %v", err)
	}
}

func TestFileSizeLimitHelper(t *testing.T) {
	binary := os.Getenv("ICON_GENERATOR_LIMIT_HELPER")
	if binary == "" {
		return
	}
	signal.Ignore(syscall.SIGXFSZ)
	if err := syscall.Setrlimit(syscall.RLIMIT_FSIZE, &syscall.Rlimit{Cur: 0, Max: 0}); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Exec(binary, []string{binary, "--seed", "42"}, []string{"PATH="}); err != nil {
		t.Fatal(err)
	}
}
