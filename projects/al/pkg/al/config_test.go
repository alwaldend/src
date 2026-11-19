package al

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigDumpPrivateAndTruncated(t *testing.T) {
	dir := t.TempDir()
	config := filepath.Join(dir, "config.json")
	out := filepath.Join(dir, "dump.json")
	if err := os.WriteFile(config, []byte(`{}`), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, exists := range []bool{false, true} {
		if exists {
			if err := os.WriteFile(out, []byte(strings.Repeat("old-synthetic-secret", 100)), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(out, 0o644); err != nil {
				t.Fatal(err)
			}
		}
		if err := DumpConfigs(context.Background(), out, config); err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(out)
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(out)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0o600 || string(data) != "{}\n" {
			t.Fatal("config output retained public permissions or stale contents")
		}
	}
}

func TestConfigErrorsDoNotEchoValues(t *testing.T) {
	dir := t.TempDir()
	for extension, content := range map[string]string{
		"json": `{"vault_auth":[{"no_auth":"synthetic-secret"}]}`,
		"lua":  `error("synthetic-secret")`,
	} {
		path := filepath.Join(dir, "config."+extension)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		_, err := LoadConfigs(context.Background(), path)
		if err == nil || strings.Contains(err.Error(), "synthetic-secret") {
			t.Fatal("config error disclosed an input value")
		}
	}
}
