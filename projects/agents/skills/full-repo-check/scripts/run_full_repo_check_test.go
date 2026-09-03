package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
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

func TestExecuteContinuesAndWritesRestrictedArtifacts(t *testing.T) {
	repoRoot := t.TempDir()
	writeWorkspaceModules(t, repoRoot)
	getenv := func(name string) string {
		if name == "BUILD_WORKSPACE_DIRECTORY" {
			return repoRoot
		}
		return ""
	}

	commands := 0
	queries := 0
	var processes []*exec.Cmd
	var checkProcesses []*exec.Cmd
	newCommand := func(name string, args ...string) *exec.Cmd {
		if name == "bazel_agent" && len(args) == 4 &&
			args[0] == "bazel" && args[1] == "query" {
			queries++
			command := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--")
			command.Env = append(
				os.Environ(),
				"GO_WANT_HELPER_PROCESS=1",
				"HELPER_EXIT_CODE=0",
				"HELPER_OUTPUT=kind(\"//target\")\n",
			)
			processes = append(processes, command)
			return command
		}
		commands++
		wantArgs := []string{"bazel", "build", "//..."}
		if commands%2 == 0 {
			wantArgs[1] = "test"
		}
		if name != "bazel_agent" ||
			strings.Join(args, " ") != strings.Join(wantArgs, " ") {
			t.Fatalf("command = %q %q, want bazel_agent %q", name, args, wantArgs)
		}
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
		checkProcesses = append(checkProcesses, command)
		return command
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if got := execute(
		getenv,
		newCommand,
		[]string{"--json-report", filepath.Join(repoRoot, "out", "full.json")},
		&stdout,
		&stderr,
	); got != 1 {
		t.Fatalf("execute() = %d, want 1", got)
	}
	if commands != 20 {
		t.Errorf("commands executed = %d, want 20", commands)
	}
	if queries != len(repositoryWorkspaces) {
		t.Errorf("query count = %d, want %d", queries, len(repositoryWorkspaces))
	}
	if got := len(processes); got != len(repositoryWorkspaces) {
		t.Fatalf("query process count = %d, want %d", got, len(repositoryWorkspaces))
	}
	if got := len(checkProcesses); got != 20 {
		t.Fatalf("check process count = %d, want 20", got)
	}
	for index, process := range checkProcesses {
		candidate := repositoryWorkspaces[index/len(checkPhases)]
		wantDir := filepath.Join(repoRoot, candidate.path)
		if process.Dir != wantDir {
			t.Errorf(
				"process %d directory = %q, want %q",
				index,
				process.Dir,
				wantDir,
			)
		}
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
		"| root | build | `bazel_agent bazel build //...` | FAIL (exit 9)",
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
	if len(logs) != 20 {
		t.Fatalf("log count = %d, want 20", len(logs))
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
	if got := execute(getenv, newCommand, nil, &bytes.Buffer{}, &stderr); got != 1 {
		t.Fatalf("execute() = %d, want 1", got)
	}
	if !strings.Contains(
		stderr.String(),
		"projects/rules_binary_toolchain/MODULE.bazel",
	) {
		t.Errorf("stderr = %q, want missing workspace", stderr.String())
	}
}

func TestJSONReportEmission(t *testing.T) {
	repoRoot := t.TempDir()
	writeWorkspaceModules(t, repoRoot)
	getenv := func(name string) string {
		if name == "BUILD_WORKSPACE_DIRECTORY" {
			return repoRoot
		}
		return ""
	}
	commandCount := 0
	newCommand := func(name string, args ...string) *exec.Cmd {
		commandCount++
		command := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--")
		command.Env = append(
			os.Environ(),
			"GO_WANT_HELPER_PROCESS=1",
			"HELPER_EXIT_CODE=0",
		)
		if name == "bazel_agent" &&
			len(args) == 3 &&
			args[0] == "query" {
			command.Env = append(command.Env, "HELPER_OUTPUT=kind(\"//target\")\n")
		}
		return command
	}
	reportPath := filepath.Join(repoRoot, "out", "report.json")
	var stderr bytes.Buffer
	if got := execute(
		getenv,
		newCommand,
		[]string{"--json-report", reportPath},
		&bytes.Buffer{},
		&stderr,
	); got != 0 {
		t.Fatalf("execute() = %d, stderr = %s", got, stderr.String())
	}
	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("os.ReadFile(report) error = %v", err)
	}
	var report checkReport
	if err := json.Unmarshal(content, &report); err != nil {
		t.Fatalf("json.Unmarshal() error = %v: %s", err, content)
	}
	if report.APIVersion != reportAPIVersion {
		t.Errorf("APIVersion = %q, want %q", report.APIVersion, reportAPIVersion)
	}
	if report.Kind != reportKind {
		t.Errorf("Kind = %q, want %q", report.Kind, reportKind)
	}
	if report.State != reportComplete {
		t.Errorf("State = %q, want %q", report.State, reportComplete)
	}
	if report.Inputs.Profile != "full" {
		t.Errorf("Profile = %q, want full", report.Inputs.Profile)
	}
	if got := len(report.Workspaces); got != len(repositoryWorkspaces) {
		t.Errorf("Workspaces = %d, want %d", got, len(repositoryWorkspaces))
	}
	if report.TargetUniverse <= 0 {
		t.Errorf("TargetUniverse = %d, want positive", report.TargetUniverse)
	}
	if got := len(report.Phases); got != len(repositoryWorkspaces)*len(checkPhases) {
		t.Errorf("Phases = %d, want %d", got, len(repositoryWorkspaces)*len(checkPhases))
	}
	completed := 0
	withResults := 0
	for _, phase := range report.Phases {
		if phase.Status == "pass" {
			completed++
		}
		if phase.Result != nil {
			withResults++
			if phase.Result.DurationMS < 0 {
				t.Errorf("phase %s/%s negative duration", phase.Workspace, phase.Phase)
			}
		}
	}
	if completed != len(report.Phases) {
		t.Errorf("pass phases = %d, want %d", completed, len(report.Phases))
	}
	if withResults != len(report.Phases) {
		t.Errorf("phases with results = %d, want %d", withResults, len(report.Phases))
	}
}

func TestResumeWithIdenticalInputs(t *testing.T) {
	repoRoot := t.TempDir()
	writeWorkspaceModules(t, repoRoot)
	getenv := func(name string) string {
		if name == "BUILD_WORKSPACE_DIRECTORY" {
			return repoRoot
		}
		return ""
	}
	newCommand := func(name string, args ...string) *exec.Cmd {
		command := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--")
		command.Env = append(
			os.Environ(),
			"GO_WANT_HELPER_PROCESS=1",
			"HELPER_EXIT_CODE=0",
		)
		return command
	}
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	reportDir := filepath.Join(repoRoot, "out")
	if err := os.MkdirAll(reportDir, 0o755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}
	reportPath := filepath.Join(reportDir, "full.json")
	if got := execute(
		getenv,
		newCommand,
		[]string{"--json-report", reportPath},
		&stdout,
		&stderr,
	); got != 0 {
		t.Fatalf("initial execute() = %d, stderr = %s", got, stderr.String())
	}
	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	var report checkReport
	if err := json.Unmarshal(content, &report); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	// The resume input is the report's own inputs; identical, so it resumes.
	// Simulate an interrupted run so the report is actually resumable.
	report.State = reportUnfinished
	if got := validateResume(&report, report.Inputs); got != nil {
		t.Errorf("validateResume(identical) error = %v", got)
	}
}

func TestResumeWithDifferentInputsRefuses(t *testing.T) {
	report := checkReport{
		APIVersion: reportAPIVersion,
		Kind:       reportKind,
		ID:         "full-repo-check.test",
		Inputs: runInputs{
			Workspaces:      []string{"root"},
			Profile:         "full",
			TargetUniverses: map[string]int64{"root": 42},
		},
		TargetUniverse: 42,
		State:          reportUnfinished,
	}
	if err := validateResume(&report, runInputs{}); err == nil {
		t.Error("validateResume() accepted different inputs")
	} else if !strings.Contains(err.Error(), "inputs do not match") {
		t.Errorf("error = %v, want inputs do not match", err)
	}
}

func TestResumeRejectsTargetUniverseReduction(t *testing.T) {
	report := checkReport{
		APIVersion: reportAPIVersion,
		Kind:       reportKind,
		ID:         "full-repo-check.test",
		Inputs: runInputs{
			Workspaces:      []string{"root"},
			Profile:         "full",
			TargetUniverses: map[string]int64{"root": 100},
		},
		TargetUniverse: 100,
		State:          "unfinished",
	}
	err := validateUniverseGrowth(
		&report,
		runInputs{TargetUniverses: map[string]int64{"root": 80}},
		false,
	)
	if err == nil {
		t.Error("validateUniverseGrowth() accepted reduction")
	} else if !strings.Contains(err.Error(), "was reduced") {
		t.Errorf("error = %v, want reduction message", err)
	}
	// Growth is refused unless explicitly allowed.
	err = validateUniverseGrowth(
		&report,
		runInputs{TargetUniverses: map[string]int64{"root": 120}},
		false,
	)
	if err == nil {
		t.Error("validateUniverseGrowth() accepted growth without allow")
	}
	if err := validateUniverseGrowth(
		&report,
		runInputs{TargetUniverses: map[string]int64{"root": 120}},
		true,
	); err != nil {
		t.Errorf("validateUniverseGrowth() with allow error = %v", err)
	}
}

func TestResumeMarksUnexecutedPhases(t *testing.T) {
	report := checkReport{
		APIVersion: reportAPIVersion,
		Kind:       reportKind,
		ID:         "full-repo-check.test",
		Phases: []jsonPhase{
			{
				Workspace: "root",
				Phase:     "build",
				Status:    "pass",
				Result: &jsonCheckResult{
					Status:     "pass",
					ExitCode:   0,
					DurationMS: 12,
				},
			},
			{Workspace: "root", Phase: "test", Status: "unexecuted"},
		},
		State: reportUnfinished,
	}
	results, composed, err := runChecks(
		"",
		"",
		func(name string, args ...string) *exec.Cmd {
			t.Fatal("command started for unexecuted phase")
			return nil
		},
		&bytes.Buffer{},
		report.Inputs,
		&report,
		time.Now(),
	)
	if err != nil {
		t.Fatalf("runChecks() error = %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("results = %d, want 2", len(results))
	}
	if results[0].status != "pass" {
		t.Errorf("resumed status = %q, want pass", results[0].status)
	}
	if results[1].status != "unexecuted" {
		t.Errorf("unexecuted status = %q, want unexecuted", results[1].status)
	}
	if composed.State != reportUnfinished {
		t.Errorf("composed state = %q, want %q", composed.State, reportUnfinished)
	}
	if composed.Phases[0].Result == nil {
		t.Error("resumed pass phase has no result")
	}
	if composed.Phases[1].Result != nil {
		t.Error("unexecuted phase must not carry a synthetic result")
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
	if output := os.Getenv("HELPER_OUTPUT"); output != "" {
		fmt.Fprint(os.Stdout, output)
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
