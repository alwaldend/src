package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func writeTestFile(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
}

func testToolWorkspace(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, relative := range commonToolSourcePaths {
		path := filepath.Join(root, filepath.FromSlash(relative))
		if filepath.Ext(relative) == ".bazelrc" ||
			filepath.Ext(relative) == ".lock" ||
			filepath.Base(relative) == "MODULE.bazel" ||
			filepath.Base(relative) == ".bazeliskrc" ||
			filepath.Base(relative) == ".bazelrc" {
			writeTestFile(t, path, relative+"\n", 0o600)
			continue
		}
		if err := os.MkdirAll(path, 0o700); err != nil {
			t.Fatal(err)
		}
		writeTestFile(t, filepath.Join(path, "settings.bazelrc"), "settings\n", 0o600)
	}
	writeTestFile(
		t,
		filepath.Join(root, "tools/repo_delivery/BUILD.bazel"),
		"alias(name = \"repo_delivery\")\n",
		0o600,
	)
	writeTestFile(
		t,
		filepath.Join(root, "tools/repo_delivery/cmd/repo_delivery/main.go"),
		"package main\n",
		0o600,
	)
	return root
}

func testToolEnvironment(root string) []string {
	return []string{
		"HOME=" + filepath.Join(root, "home"),
		"PATH=" + os.Getenv("PATH"),
	}
}

func fakeRepoDeliveryBuild(t *testing.T, calls *atomic.Int32) toolBuildFunc {
	t.Helper()
	return func(workspace, _ string, spec cachedToolSpec, _ []string) error {
		calls.Add(1)
		writeTestFile(
			t,
			filepath.Join(workspace, filepath.FromSlash(spec.Output)),
			"#!/bin/sh\nexit 0\n",
			0o700,
		)
		return nil
	}
}

func TestParseToolRun(t *testing.T) {
	options, name, arguments, err := parseToolRun([]string{
		"--workspace-root",
		"/workspace",
		"--cache-root",
		"/cache",
		"repo_delivery",
		"--",
		"provider",
	})
	if err != nil {
		t.Fatal(err)
	}
	if options.WorkspaceRoot != "/workspace" || options.CacheRoot != "/cache" {
		t.Fatalf("options = %#v", options)
	}
	if name != "repo_delivery" || !reflect.DeepEqual(arguments, []string{"provider"}) {
		t.Fatalf("name = %q, arguments = %q", name, arguments)
	}
}

func TestOverrideEnvironmentReplacesRegisteredValues(t *testing.T) {
	got := overrideEnvironment(
		[]string{"PATH=/bin", "BAZEL_BINDIR=old", "HOME=/home/example"},
		[]string{"BAZEL_BINDIR=."},
	)
	want := []string{"PATH=/bin", "HOME=/home/example", "BAZEL_BINDIR=."}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("overrideEnvironment() = %q, want %q", got, want)
	}
}

func TestMcpCordisRegistersRulesJSExecutionEnvironment(t *testing.T) {
	want := []string{"BAZEL_BINDIR=."}
	if got := cachedToolSpecs["mcp_cordis"].Environment; !reflect.DeepEqual(got, want) {
		t.Fatalf("Cordis environment = %q, want %q", got, want)
	}
}

func TestBazelToolBuildUsesGeneratedLauncher(t *testing.T) {
	arguments := bazelToolBuildArguments(
		"/workspace/out/tool-build/launch",
		cachedToolSpecs["mcp_cordis"],
	)
	want := []string{
		"run",
		"--config=agent",
		"--script_path=/workspace/out/tool-build/launch",
		"//projects/mcp_cordis:mcp_cordis",
	}
	if !reflect.DeepEqual(arguments, want) {
		t.Fatalf("bazelToolBuildArguments() = %q, want %q", arguments, want)
	}
}

func TestToolSourceDigestTracksDeclaredInputsAcrossWorktrees(t *testing.T) {
	first := testToolWorkspace(t)
	second := testToolWorkspace(t)
	spec := cachedToolSpecs["repo_delivery"]
	firstDigest, err := toolSourceDigest(first, spec, testToolEnvironment(first))
	if err != nil {
		t.Fatal(err)
	}
	secondDigest, err := toolSourceDigest(second, spec, testToolEnvironment(second))
	if err != nil {
		t.Fatal(err)
	}
	if firstDigest != secondDigest {
		t.Fatalf("identical worktrees have different digests: %s != %s", firstDigest, secondDigest)
	}
	writeTestFile(
		t,
		filepath.Join(second, "tools/repo_delivery/cmd/repo_delivery/main.go"),
		"package main\n// changed\n",
		0o600,
	)
	changedDigest, err := toolSourceDigest(second, spec, testToolEnvironment(second))
	if err != nil {
		t.Fatal(err)
	}
	if changedDigest == firstDigest {
		t.Fatal("declared source change did not change the cache digest")
	}
}

func TestToolSourceDigestTracksTransitiveBazelRCImports(t *testing.T) {
	root := testToolWorkspace(t)
	writeTestFile(
		t,
		filepath.Join(root, ".bazelrc"),
		"import %workspace%/config/outer.bazelrc\n",
		0o600,
	)
	writeTestFile(
		t,
		filepath.Join(root, "config/outer.bazelrc"),
		"try-import '%workspace%/config/inner bazelrc'\n",
		0o600,
	)
	inner := filepath.Join(root, "config/inner bazelrc")
	writeTestFile(t, inner, "build --define=value=first\n", 0o600)
	spec := cachedToolSpecs["repo_delivery"]
	first, err := toolSourceDigest(root, spec, testToolEnvironment(root))
	if err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, inner, "build --define=value=second\n", 0o600)
	second, err := toolSourceDigest(root, spec, testToolEnvironment(root))
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("transitively imported bazelrc did not change the cache digest")
	}
}

func TestToolSourceDigestTracksHomeBazelRCImports(t *testing.T) {
	root := testToolWorkspace(t)
	home := filepath.Join(root, "home")
	settings := filepath.Join(home, "settings.bazelrc")
	writeTestFile(
		t,
		filepath.Join(home, ".bazelrc"),
		"import "+settings+"\n",
		0o600,
	)
	writeTestFile(t, settings, "build --define=value=first\n", 0o600)
	spec := cachedToolSpecs["repo_delivery"]
	first, err := toolSourceDigest(root, spec, testToolEnvironment(root))
	if err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, settings, "build --define=value=second\n", 0o600)
	second, err := toolSourceDigest(root, spec, testToolEnvironment(root))
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("home bazelrc import did not change the cache digest")
	}
}

func TestToolSourceDigestTracksEnvironmentBazelRC(t *testing.T) {
	root := testToolWorkspace(t)
	rc := filepath.Join(root, "config/environment.bazelrc")
	writeTestFile(t, rc, "build --define=value=first\n", 0o600)
	environment := append(
		testToolEnvironment(root),
		"BAZELRC="+rc,
	)
	spec := cachedToolSpecs["repo_delivery"]
	first, err := toolSourceDigest(root, spec, environment)
	if err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, rc, "build --define=value=second\n", 0o600)
	second, err := toolSourceDigest(root, spec, environment)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("BAZELRC environment file did not change the cache digest")
	}
}

func TestMcpCordisDigestExcludesWorkspacePackages(t *testing.T) {
	root := testToolWorkspace(t)
	for _, relative := range cachedToolSpecs["mcp_cordis"].SourcePaths {
		path := filepath.Join(root, filepath.FromSlash(relative))
		if relative == "projects/mcp_cordis/internal" ||
			relative == "projects/mcp_cordis/patches" {
			writeTestFile(t, filepath.Join(path, "input"), relative+"\n", 0o600)
			continue
		}
		writeTestFile(t, path, relative+"\n", 0o600)
	}
	plugin := filepath.Join(root, "projects/mcp_cordis/plugins/example.mjs")
	writeTestFile(t, plugin, "export default {};\n", 0o600)
	spec := cachedToolSpecs["mcp_cordis"]
	first, err := toolSourceDigest(root, spec, testToolEnvironment(root))
	if err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, plugin, "export default { changed: true };\n", 0o600)
	second, err := toolSourceDigest(root, spec, testToolEnvironment(root))
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatal("workspace plugin change invalidated the stable Cordis runtime")
	}
	writeTestFile(
		t,
		filepath.Join(root, "projects/mcp_cordis/internal/input"),
		"runtime changed\n",
		0o600,
	)
	third, err := toolSourceDigest(root, spec, testToolEnvironment(root))
	if err != nil {
		t.Fatal(err)
	}
	if third == second {
		t.Fatal("Cordis runtime change did not invalidate the cache")
	}
}

func TestEnsureCachedToolBuildsOnceAndReusesExactEntry(t *testing.T) {
	workspace := testToolWorkspace(t)
	cache := filepath.Join(t.TempDir(), "cache")
	var calls atomic.Int32
	options := toolOptions{WorkspaceRoot: workspace, CacheRoot: cache}
	environment := testToolEnvironment(workspace)
	executable, digest, installed, err := ensureCachedTool(
		"repo_delivery",
		options,
		environment,
		fakeRepoDeliveryBuild(t, &calls),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !installed || calls.Load() != 1 || !strings.Contains(executable, digest) {
		t.Fatalf(
			"executable = %q, digest = %q, installed = %t, calls = %d",
			executable,
			digest,
			installed,
			calls.Load(),
		)
	}
	secondExecutable, secondDigest, installed, err := ensureCachedTool(
		"repo_delivery",
		options,
		environment,
		fakeRepoDeliveryBuild(t, &calls),
	)
	if err != nil {
		t.Fatal(err)
	}
	if installed || calls.Load() != 1 || secondExecutable != executable ||
		secondDigest != digest {
		t.Fatalf(
			"second result = %q %q %t, calls = %d",
			secondExecutable,
			secondDigest,
			installed,
			calls.Load(),
		)
	}
}

func TestEnsureCachedToolSerializesConcurrentMisses(t *testing.T) {
	workspace := testToolWorkspace(t)
	cache := filepath.Join(t.TempDir(), "cache")
	options := toolOptions{WorkspaceRoot: workspace, CacheRoot: cache}
	environment := testToolEnvironment(workspace)
	var calls atomic.Int32
	build := func(workspace, _ string, spec cachedToolSpec, _ []string) error {
		calls.Add(1)
		time.Sleep(50 * time.Millisecond)
		output := filepath.Join(workspace, filepath.FromSlash(spec.Output))
		if err := os.MkdirAll(filepath.Dir(output), 0o700); err != nil {
			return err
		}
		return os.WriteFile(output, []byte("#!/bin/sh\nexit 0\n"), 0o700)
	}
	const workers = 4
	var group sync.WaitGroup
	errors := make(chan error, workers)
	for index := 0; index < workers; index++ {
		group.Add(1)
		go func() {
			defer group.Done()
			_, _, _, err := ensureCachedTool(
				"repo_delivery",
				options,
				environment,
				build,
			)
			errors <- err
		}()
	}
	group.Wait()
	close(errors)
	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("build calls = %d, want 1", calls.Load())
	}
}

func TestEnsureCachedToolRetriesWhenInputsChangeDuringBuild(t *testing.T) {
	workspace := testToolWorkspace(t)
	cache := filepath.Join(t.TempDir(), "cache")
	options := toolOptions{WorkspaceRoot: workspace, CacheRoot: cache}
	environment := testToolEnvironment(workspace)
	var calls atomic.Int32
	build := func(workspace, _ string, spec cachedToolSpec, _ []string) error {
		call := calls.Add(1)
		writeTestFile(
			t,
			filepath.Join(workspace, filepath.FromSlash(spec.Output)),
			"#!/bin/sh\nexit 0\n",
			0o700,
		)
		if call == 1 {
			writeTestFile(
				t,
				filepath.Join(
					workspace,
					"tools/repo_delivery/cmd/repo_delivery/main.go",
				),
				"package main\n// changed during build\n",
				0o600,
			)
		}
		return nil
	}
	_, _, installed, err := ensureCachedTool(
		"repo_delivery",
		options,
		environment,
		build,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !installed || calls.Load() != 2 {
		t.Fatalf("installed = %t, build calls = %d, want true and 2", installed, calls.Load())
	}
}

func TestCopyTreeMaterializesExternalAndPreservesInternalSymlinks(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	destination := filepath.Join(root, "destination")
	external := filepath.Join(root, "external")
	writeTestFile(t, filepath.Join(source, "data/value"), "internal\n", 0o600)
	if err := os.Chmod(filepath.Join(source, "data"), 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(filepath.Join(source, "data"), 0o700)
	})
	if err := os.Symlink("data/value", filepath.Join(source, "internal")); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, external, "external\n", 0o600)
	if err := os.Symlink(external, filepath.Join(source, "external")); err != nil {
		t.Fatal(err)
	}
	if err := copyTree(source, destination); err != nil {
		t.Fatal(err)
	}
	if target, err := os.Readlink(filepath.Join(destination, "internal")); err != nil ||
		target != "data/value" {
		t.Fatalf("internal symlink = %q, %v", target, err)
	}
	info, err := os.Lstat(filepath.Join(destination, "external"))
	if err != nil {
		t.Fatal(err)
	}
	if !info.Mode().IsRegular() {
		t.Fatalf("external link mode = %s, want materialized file", info.Mode())
	}
	directoryInfo, err := os.Stat(filepath.Join(destination, "data"))
	if err != nil {
		t.Fatal(err)
	}
	if directoryInfo.Mode().Perm() != 0o755 {
		t.Fatalf("copied directory mode = %o, want 0755", directoryInfo.Mode().Perm())
	}
}

func TestRunToolExecutesCachedBinaryWithOriginalEnvironment(t *testing.T) {
	workspace := testToolWorkspace(t)
	cache := filepath.Join(t.TempDir(), "cache")
	environment := testToolEnvironment(workspace)
	var calls atomic.Int32
	var gotPath string
	var gotArguments, gotEnvironment []string
	err := runTool(
		[]string{
			"run",
			"--workspace-root",
			workspace,
			"--cache-root",
			cache,
			"repo_delivery",
			"--",
			"provider",
		},
		environment,
		func(path string, arguments, environment []string) error {
			gotPath = path
			gotArguments = arguments
			gotEnvironment = environment
			return nil
		},
		fakeRepoDeliveryBuild(t, &calls),
	)
	if err != nil {
		t.Fatal(err)
	}
	if gotPath == "" || !reflect.DeepEqual(gotArguments, []string{gotPath, "provider"}) {
		t.Fatalf("path = %q, arguments = %q", gotPath, gotArguments)
	}
	if !reflect.DeepEqual(gotEnvironment, environment) {
		t.Fatalf("environment = %q, want %q", gotEnvironment, environment)
	}
}

func TestToolCacheRejectsWritableRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "cache")
	if err := os.Mkdir(root, 0o777); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(root, 0o777); err != nil {
		t.Fatal(err)
	}
	err := ensurePrivateDirectory(root)
	if err == nil || !strings.Contains(err.Error(), "must not be writable") {
		t.Fatalf("ensurePrivateDirectory() error = %v", err)
	}
}
