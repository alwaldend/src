package bazelcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildDoctorReportIsBoundedAndDetectsStaleInstall(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".bazeliskrc"), []byte(
		"USE_BAZEL_VERSION=8.3.1\nBAZELISK_VERIFY_SHA256=abc123\n",
	), 0o600); err != nil {
		t.Fatal(err)
	}
	runner := filepath.Join(root, "installed")
	if err := os.WriteFile(runner, []byte("installed"), 0o700); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(root, "bazel-bin", "external", "rules_al+", "cmd", "al", "al_", "al")
	if err := os.MkdirAll(filepath.Dir(source), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(source, []byte("built"), 0o700); err != nil {
		t.Fatal(err)
	}
	report, err := buildDoctorReport(doctorOptions{workspaceRoot: root, scratchPath: "out/task-one/run-one"}, runner, "/tools/bazel")
	if err != nil {
		t.Fatal(err)
	}
	if report.Schema != "agents.alwaldend.com/bazel-agent-doctor/v1alpha1" ||
		report.Source.StaleInstall == nil || !*report.Source.StaleInstall {
		t.Fatalf("unexpected doctor report: %#v", report)
	}
	if !report.Scratch.Namespaced || report.Scratch.Reason != "" {
		t.Fatalf("scratch = %#v, want namespaced", report.Scratch)
	}
	if strings.Contains(strings.Join(report.RCFiles, " "), root) {
		t.Fatalf("rc composition leaked fixture root: %#v", report.RCFiles)
	}
}

func TestTaskScratchRejectsSharedOrEscapingPaths(t *testing.T) {
	root := t.TempDir()
	for _, candidate := range []string{"out", "out/shared", "../outside"} {
		got := taskScratch(root, candidate)
		if got.Namespaced || got.Reason == "" {
			t.Fatalf("taskScratch(%q) = %#v, want classified refusal", candidate, got)
		}
	}
}
