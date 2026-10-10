package collection

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/textstore"
)

// TestRefreshAndHomePathsWorkflow verifies refreshed exports and home-relative paths.
func TestRefreshAndHomePathsWorkflow(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	t.Setenv("HOME", root)
	base := fixture(t, root)
	input := "~/" + filepath.Base(base)
	if err := Export(ctx, input, "~/text"); err != nil {
		t.Fatal(err)
	}
	text := filepath.Join(root, "text")
	unrelated := filepath.Join(text, "README.txt")
	if err := os.WriteFile(unrelated, []byte("keep me"), 0o600); err != nil {
		t.Fatal(err)
	}
	stale := textstore.NotePath(text, "grammar", 999)
	if err := os.WriteFile(stale, []byte("stale declaration"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Export(ctx, input, "~/text"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatalf("stale note survived refresh: %v", err)
	}
	value, err := os.ReadFile(unrelated)
	if err != nil || string(value) != "keep me" {
		t.Fatalf("unrelated file changed: %q %v", value, err)
	}
	config := filepath.Join(text, textstore.CollectionMarker)
	before, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	if err := Export(ctx, "~/missing.colpkg", "~/text"); err == nil {
		t.Fatal("invalid input accepted")
	}
	after, err := os.ReadFile(config)
	if err != nil || string(before) != string(after) {
		t.Fatal("failed export modified existing declarations")
	}
	alternate := filepath.Join(text, "custom-root.toml")
	if err := os.Rename(config, alternate); err != nil {
		t.Fatal(err)
	}
	plan, err := Plan(ctx, input, "~/text/custom-root.toml")
	if err != nil || len(plan.Changes) != 0 {
		t.Fatalf("home-relative arbitrary config did not converge: %+v %v", plan, err)
	}
	if err := Build(ctx, input, "~/text/custom-root.toml", "~/rebuilt.colpkg"); err != nil {
		t.Fatal(err)
	}
	if err := Apply(ctx, "~/rebuilt.colpkg", "~/text/custom-root.toml"); err != nil {
		t.Fatal(err)
	}
	note := textstore.NotePath(text, "grammar", 100)
	if _, err := textstore.GenerateIDs(ctx, []string{strings.Replace(note, root, "~", 1)}); err != nil {
		t.Fatal(err)
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := os.WriteFile(filepath.Join(artifacts, "refresh-home-paths.txt"), []byte("Repeated export removes stale declarations, retains unrelated files, and preserves the previous export on invalid input; home-relative export/plan/build/apply/identity paths and arbitrary root filenames converge.\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
