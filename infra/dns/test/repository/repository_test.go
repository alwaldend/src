package repository_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"git.alwaldend.com/alwaldend/src/infra/dns/internal/lint"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

var workspaceMarker = flag.String("workspace-marker", "", "root module runfile")

func TestRepositoryDNSOwnership(t *testing.T) {
	marker, err := runfiles.Rlocation(*workspaceMarker)
	if err != nil {
		t.Fatal(err)
	}
	marker, err = filepath.EvalSymlinks(marker)
	if err != nil {
		t.Fatal(err)
	}
	report, err := lint.Scan(os.DirFS(filepath.Dir(marker)), "alwaldend.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Sources) == 0 {
		t.Fatal("no dnsconfig.json files discovered in the repository checkout")
	}
	t.Logf("Checked ownership across %d DNS source files", len(report.Sources))
}
