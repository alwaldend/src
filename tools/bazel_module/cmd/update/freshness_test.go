package main

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// workspaceMarker locates the source worktree through its root module runfile.
var workspaceMarker = flag.String("workspace-marker", "", "root module runfile")

// TestRootModuleIsFresh compares the current checkout against the generated include list.
func TestRootModuleIsFresh(t *testing.T) {
	marker, err := runfiles.Rlocation(*workspaceMarker)
	if err != nil {
		t.Fatal(err)
	}
	marker, err = filepath.EvalSymlinks(marker)
	if err != nil {
		t.Fatal(err)
	}
	if err := update(t.Context(), filepath.Dir(marker), *testGit, true); err != nil {
		t.Fatal(err)
	}
}
