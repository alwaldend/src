package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestBazelArguments(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "no command",
			want: []string{},
		},
		{
			name: "help passes through",
			args: []string{"--help"},
			want: []string{"--help", "--config=agent"},
		},
		{
			name: "build",
			args: []string{"build", "//path/to/package:all"},
			want: []string{
				"build",
				"--config=agent",
				"//path/to/package:all",
			},
		},
		{
			name: "run separator",
			args: []string{"run", "//path/to:tool", "--", "--tool-flag"},
			want: []string{
				"run",
				"--config=agent",
				"//path/to:tool",
				"--",
				"--tool-flag",
			},
		},
		{
			name: "later options stay later",
			args: []string{"test", "--config=local", "//..."},
			want: []string{
				"test",
				"--config=agent",
				"--config=local",
				"//...",
			},
		},
		{
			name: "validation belongs to Bazel",
			args: []string{"--not-a-command", "build", "//..."},
			want: []string{
				"--not-a-command",
				"--config=agent",
				"build",
				"//...",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := bazelArguments(test.args)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("bazelArguments() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestRunReplacesProcess(t *testing.T) {
	environment := []string{
		"PATH=/bin",
		"TEMP=/old/temp",
		"TMP=/old/tmp",
		"TMPDIR=/old/tmpdir",
	}
	var gotPath string
	var gotArgs, gotEnvironment []string
	err := run(
		[]string{"test", "//pkg:all"},
		environment,
		func(name string) (string, error) {
			if name != "bazel" {
				t.Fatalf("lookPath(%q), want bazel", name)
			}
			return "/bin/bazel", nil
		},
		func(path string, args, env []string) error {
			gotPath = path
			gotArgs = args
			gotEnvironment = env
			return nil
		},
	)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if gotPath != "/bin/bazel" {
		t.Fatalf("path = %q, want /bin/bazel", gotPath)
	}
	wantArgs := []string{
		"/bin/bazel",
		"test",
		"--config=agent",
		"//pkg:all",
	}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %q, want %q", gotArgs, wantArgs)
	}
	if !reflect.DeepEqual(gotEnvironment, environment) {
		t.Fatalf("environment = %q, want unchanged %q", gotEnvironment, environment)
	}
}

func TestRunReportsLookupFailure(t *testing.T) {
	wantErr := errors.New("not found")
	err := run(
		[]string{"build", "//..."},
		nil,
		func(string) (string, error) { return "", wantErr },
		func(string, []string, []string) error {
			t.Fatal("replaceProcess called after lookup failure")
			return nil
		},
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf("run() error = %v, want wrapped %v", err, wantErr)
	}
}

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
	source := filepath.Join(root, "bazel-bin", "projects", "bazel_agent", "cmd", "bazel_agent", "bazel_agent_", "bazel_agent")
	if err := os.MkdirAll(filepath.Dir(source), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(source, []byte("built"), 0o700); err != nil {
		t.Fatal(err)
	}
	report, err := buildDoctorReport(root, "out/task-one/run-one", runner, "/tools/bazel")
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
