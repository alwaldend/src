package repository_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"git.alwaldend.com/alwaldend/src/infra/dns/internal/dump"
	"git.alwaldend.com/alwaldend/src/infra/dns/internal/lint"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

var workspaceMarker = flag.String("workspace-marker", "", "root module runfile")

// checkoutRoot resolves the repository root from the root module runfile.
func checkoutRoot(t *testing.T) string {
	t.Helper()
	marker, err := runfiles.Rlocation(*workspaceMarker)
	if err != nil {
		t.Fatal(err)
	}
	marker, err = filepath.EvalSymlinks(marker)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Dir(marker)
}

func TestRepositoryDNSOwnership(t *testing.T) {
	report, err := lint.Scan(os.DirFS(checkoutRoot(t)), "alwaldend.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Sources) == 0 {
		t.Fatal("no dnsconfig.json files discovered in the repository checkout")
	}
	t.Logf("Checked ownership across %d DNS source files", len(report.Sources))
}

// TestGeneratedDNSPagesAreCurrent keeps the checked-in declaration pages in
// step with the declarations they project.
func TestGeneratedDNSPagesAreCurrent(t *testing.T) {
	root := checkoutRoot(t)
	pages, err := dump.RenderPages(os.DirFS(root))
	if err != nil {
		t.Fatal(err)
	}
	for _, page := range pages {
		path := filepath.Join(root, filepath.FromSlash(page.Path))
		existing, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: %v", page.Path, err)
			continue
		}
		if string(existing) != page.Body {
			t.Errorf("%s is out of date; run bazel run //infra/dns/cmd/dump -- --write", page.Path)
		}
	}
}
