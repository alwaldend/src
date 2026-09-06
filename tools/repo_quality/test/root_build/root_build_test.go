package root_build_test

import (
	"flag"
	"os"
	"testing"

	"github.com/bazelbuild/buildtools/build"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

var rootBuild = flag.String("root-build", "", "runfile path of the root BUILD.bazel")

func loadStatements(t *testing.T, source []byte) []*build.LoadStmt {
	t.Helper()
	file, err := build.ParseBuild("BUILD.bazel", source)
	if err != nil {
		t.Fatal(err)
	}
	var loads []*build.LoadStmt
	for _, statement := range file.Stmt {
		if load, ok := statement.(*build.LoadStmt); ok {
			loads = append(loads, load)
		}
	}
	return loads
}

func TestRootBuildHasNoLoads(t *testing.T) {
	path, err := runfiles.Rlocation(*rootBuild)
	if err != nil {
		t.Fatal(err)
	}
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, load := range loadStatements(t, source) {
		start, _ := load.Span()
		t.Errorf("BUILD.bazel:%d: load(%q) makes root package loading evaluate unrelated dependencies; move the implementation to its owning package and keep a native alias or test_suite here", start.Line, load.Module.Value)
	}
}

func TestLoadDetection(t *testing.T) {
	for _, test := range []struct {
		name   string
		source string
		want   int
	}{
		{"native alias", `alias(name = "tool", actual = "//tools/tool")`, 0},
		{"comment", "# load(\"//tools:defs.bzl\", \"tool\")\n", 0},
		{"string", `exports_files(["load(example)"])`, 0},
		{"load", `load("//tools:defs.bzl", "tool")`, 1},
		{"multiline load", "load(\n    \"//tools:defs.bzl\",\n    \"tool\",\n)\n", 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := len(loadStatements(t, []byte(test.source))); got != test.want {
				t.Errorf("load count = %d, want %d", got, test.want)
			}
		})
	}
}
