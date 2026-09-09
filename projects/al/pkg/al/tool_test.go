package al

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
)

func TestLoadToolConfiguration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "al.lua")
	content := `
tool_method({
    name = "build",
    kind = "bazel",
    bazel = {
        target = "//test:tool",
        output = "bazel-bin/test/tool_/tool",
    },
})

tool({
    name = "fixture",
    method = "build",
})
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := LoadConfigs(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if len(config.GetTools()) != 1 || config.GetTools()[0].GetName() != "fixture" || config.GetTools()[0].GetMethod() != "build" {
		t.Fatalf("tools were not parsed: %s", config.GetTools())
	}
	if len(config.GetToolMethods()) != 1 || config.GetToolMethods()[0].GetName() != "build" || config.GetToolMethods()[0].GetBazel().GetTarget() != "//test:tool" {
		t.Fatalf("tool methods were not parsed: %s", config.GetToolMethods())
	}
}

func TestToolLookup(t *testing.T) {
	config := &al_proto.Config{
		Tools:       []*al_proto.Tool{{Name: "fixture", Method: "build"}},
		ToolMethods: []*al_proto.ToolMethod{{Name: "build", Kind: "bazel", Options: &al_proto.ToolMethod_Bazel{Bazel: &al_proto.BazelToolMethod{Target: "//test:tool", Output: "bin"}}}},
	}
	tool, err := ToolByName(config, "fixture")
	if err != nil {
		t.Fatal(err)
	}
	if tool.GetMethod() != "build" {
		t.Fatalf("method = %q, want build", tool.GetMethod())
	}
	if _, err := ToolByName(config, "missing"); err == nil {
		t.Fatal("missing tool did not fail")
	}
	if _, err := ToolByName(&al_proto.Config{Tools: []*al_proto.Tool{{Name: "fixture", Method: "build"}, {Name: "fixture", Method: "build"}}}, "fixture"); err == nil {
		t.Fatal("duplicate tool did not fail")
	}
	if _, err := ToolMethodByName(config, "build"); err != nil {
		t.Fatal(err)
	}
}

type fakeBuild struct {
	calls   int
	targets []string
}

func (build *fakeBuild) build(workspace string, method *al_proto.ToolMethod, environment []string) error {
	build.calls++
	build.targets = append(build.targets, method.GetBazel().GetTarget())
	output := filepath.Join(workspace, method.GetBazel().GetOutput())
	if err := os.MkdirAll(filepath.Dir(output), 0o700); err != nil {
		return err
	}
	return os.WriteFile(output, []byte("#!/bin/sh\nexit 0\n"), 0o755)
}

func TestEnsureToolCacheEntry(t *testing.T) {
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, "MODULE.bazel"), []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workspace, ".bazelrc"), []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "al.lua"), []byte("config({})\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	method := &al_proto.ToolMethod{Name: "fixture", Kind: "bazel", Options: &al_proto.ToolMethod_Bazel{Bazel: &al_proto.BazelToolMethod{Target: "//test:tool", Output: "bazel-bin/tool"}}}
	var logs strings.Builder
	build := fakeBuild{}
	cacheRoot := t.TempDir()
	executable, digest, installed, err := ensureToolCacheEntry(&logs, workspace, cacheRoot, method, nil, build.build)
	if err != nil {
		t.Fatal(err)
	}
	if !installed || build.calls != 1 || !strings.Contains(logs.String(), "cache miss") {
		t.Fatalf("unexpected miss result: installed=%t calls=%d logs=%q", installed, build.calls, logs.String())
	}
	output := filepath.Join(workspace, method.GetBazel().GetOutput())
	if _, err := os.Stat(output); err != nil {
		t.Fatalf("cache installation removed build output: %v", err)
	}
	info, err := os.Stat(executable)
	if err != nil || info.Mode().IsRegular() == false || info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("invalid executable: mode=%v err=%v", info.Mode(), err)
	}
	logs.Reset()
	calls := build.calls
	executable, sameDigest, installed, err := ensureToolCacheEntry(&logs, workspace, cacheRoot, method, nil, build.build)
	if err != nil || installed || build.calls != calls || logs.Len() != 0 {
		t.Fatalf("unexpected hit result: installed=%t calls=%d logs=%q err=%v", installed, build.calls, logs.String(), err)
	}
	if sameDigest != digest {
		t.Fatalf("digest changed without input change: %q to %q", digest, sameDigest)
	}
}

func TestRunToolInvalidCacheOutput(t *testing.T) {
	workspace := t.TempDir()
	for _, name := range []string{"MODULE.bazel", ".bazelrc", "al.lua"} {
		if err := os.WriteFile(filepath.Join(workspace, name), []byte(""), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	method := &al_proto.ToolMethod{Name: "fixture", Kind: "bazel", Options: &al_proto.ToolMethod_Bazel{Bazel: &al_proto.BazelToolMethod{Target: "//test:tool", Output: "missing"}}}
	err := buildToolWithBazel(workspace, method, nil)
	if err == nil || !strings.Contains(err.Error(), "find bazel in PATH") {
		t.Fatalf("invalid output behavior not covered without host bazel: %v", err)
	}
}
