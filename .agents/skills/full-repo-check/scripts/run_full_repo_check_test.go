package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func writeWorkspaceModules(t *testing.T, repoRoot string) {
	t.Helper()
	for _, candidate := range repositoryWorkspaces {
		modulePath := filepath.Join(repoRoot, candidate.path, "MODULE.bazel")
		if err := os.MkdirAll(filepath.Dir(modulePath), 0o755); err != nil {
			t.Fatalf("os.MkdirAll() error = %v", err)
		}
		if err := os.WriteFile(modulePath, []byte("module()\n"), 0o644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}
	}
}

func TestExecuteContinuesAndWritesPrivateArtifacts(t *testing.T) {
	repoRoot := t.TempDir()
	writeWorkspaceModules(t, repoRoot)
	getenv := func(name string) string {
		if name == "BUILD_WORKSPACE_DIRECTORY" {
			return repoRoot
		}
		return ""
	}

	commands := 0
	newCommand := func(name string, args ...string) *exec.Cmd {
		commands++
		exitCode := 0
		if commands == 1 {
			exitCode = 9
		}
		command := exec.Command(
			os.Args[0],
			"-test.run=TestHelperProcess",
			"--",
		)
		command.Env = append(
			os.Environ(),
			"GO_WANT_HELPER_PROCESS=1",
			fmt.Sprintf("HELPER_EXIT_CODE=%d", exitCode),
		)
		return command
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if got := execute(getenv, newCommand, &stdout, &stderr); got != 1 {
		t.Fatalf("execute() = %d, want 1", got)
	}
	if commands != 10 {
		t.Errorf("commands executed = %d, want 10", commands)
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}

	runDirectories, err := filepath.Glob(
		filepath.Join(repoRoot, "out", "full-repo-check", "run.*"),
	)
	if err != nil {
		t.Fatalf("filepath.Glob() error = %v", err)
	}
	if len(runDirectories) != 1 {
		t.Fatalf("run directories = %v, want one", runDirectories)
	}
	runDirectory := runDirectories[0]
	assertPermissions(t, runDirectory, 0o700)
	assertPermissions(t, filepath.Join(runDirectory, "logs"), 0o700)

	reportPath := filepath.Join(runDirectory, "report.md")
	assertPermissions(t, reportPath, 0o600)
	report, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("os.ReadFile(report) error = %v", err)
	}
	for _, want := range []string{
		"| root | build | `bazel build --config=agent //...` | FAIL (exit 9)",
		"| projects/rules_template | test |",
		"## Failed commands",
		"| root | build | 9 |",
	} {
		if !strings.Contains(string(report), want) {
			t.Errorf("report missing %q\n%s", want, report)
		}
	}

	logs, err := os.ReadDir(filepath.Join(runDirectory, "logs"))
	if err != nil {
		t.Fatalf("os.ReadDir(logs) error = %v", err)
	}
	if len(logs) != 10 {
		t.Fatalf("log count = %d, want 10", len(logs))
	}
	firstLog := filepath.Join(runDirectory, "logs", "root.build.log")
	assertPermissions(t, firstLog, 0o600)
	contents, err := os.ReadFile(firstLog)
	if err != nil {
		t.Fatalf("os.ReadFile(log) error = %v", err)
	}
	for _, want := range []string{"helper stdout", "helper stderr"} {
		if !strings.Contains(string(contents), want) {
			t.Errorf("combined log missing %q: %s", want, contents)
		}
	}
}

func TestExecuteRejectsMissingWorkspace(t *testing.T) {
	repoRoot := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(repoRoot, "MODULE.bazel"),
		[]byte("module()\n"),
		0o644,
	); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	getenv := func(name string) string {
		if name == "BUILD_WORKSPACE_DIRECTORY" {
			return repoRoot
		}
		return ""
	}
	newCommand := func(name string, args ...string) *exec.Cmd {
		t.Fatal("command started despite missing workspace")
		return nil
	}

	var stderr bytes.Buffer
	if got := execute(getenv, newCommand, &bytes.Buffer{}, &stderr); got != 1 {
		t.Fatalf("execute() = %d, want 1", got)
	}
	if !strings.Contains(
		stderr.String(),
		"projects/rules_binary_toolchain/MODULE.bazel",
	) {
		t.Errorf("stderr = %q, want missing workspace", stderr.String())
	}
}

func assertPermissions(t *testing.T, path string, want os.FileMode) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("os.Stat(%q) error = %v", path, err)
	}
	if got := info.Mode().Perm(); got != want {
		t.Errorf("permissions for %q = %o, want %o", path, got, want)
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintln(os.Stdout, "helper stdout")
	fmt.Fprintln(os.Stderr, "helper stderr")
	exitCode, err := strconv.Atoi(os.Getenv("HELPER_EXIT_CODE"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	os.Exit(exitCode)
}
