package main

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

var (
	workspaceMarker = flag.String("workspace-marker", "", "root module runfile")
	gitRunfile      = flag.String("git", "", "Git executable runfile")
)

func TestRootModuleIsFresh(t *testing.T) {
	marker, err := runfiles.Rlocation(*workspaceMarker)
	if err != nil {
		t.Fatal(err)
	}
	marker, err = filepath.EvalSymlinks(marker)
	if err != nil {
		t.Fatal(err)
	}
	if err := update(filepath.Dir(marker), *gitRunfile, true); err != nil {
		t.Fatal(err)
	}
}
