package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
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
			want: []string{"--batch"},
		},
		{
			name: "help passes through",
			args: []string{"--help"},
			want: []string{"--batch", "--help", "--config=agent"},
		},
		{
			name: "build",
			args: []string{"build", "//path/to/package:all"},
			want: []string{
				"--batch",
				"build",
				"--config=agent",
				"//path/to/package:all",
			},
		},
		{
			name: "run separator",
			args: []string{"run", "//path/to:tool", "--", "--tool-flag"},
			want: []string{
				"--batch",
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
				"--batch",
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
				"--batch",
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
	workspace := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(workspace, "MODULE.bazel"),
		[]byte("module(name = \"test\")\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}
	oldWorkingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(workspace); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWorkingDirectory); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	environment := []string{
		"PATH=/bin",
		"TEMP=/old/temp",
		"TMP=/old/tmp",
		"TMPDIR=/old/tmpdir",
	}
	var gotPath string
	var gotArgs, gotEnvironment []string
	err = run(
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
		"--batch",
		"test",
		"--config=agent",
		"//pkg:all",
	}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %q, want %q", gotArgs, wantArgs)
	}
	tmpDirectory := filepath.Join(workspace, temporaryDirectory)
	wantEnvironment := []string{
		"PATH=/bin",
		"TMPDIR=" + tmpDirectory,
		"TMP=" + tmpDirectory,
		"TEMP=" + tmpDirectory,
	}
	if !reflect.DeepEqual(gotEnvironment, wantEnvironment) {
		t.Fatalf("environment = %q, want %q", gotEnvironment, wantEnvironment)
	}
	if info, err := os.Stat(tmpDirectory); err != nil || !info.IsDir() {
		t.Fatalf("temporary directory stat = %v, %v; want directory", info, err)
	}
}

func TestRunReportsLookupFailure(t *testing.T) {
	workspace := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(workspace, "MODULE.bazel"),
		[]byte("module(name = \"test\")\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}
	oldWorkingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(workspace); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWorkingDirectory); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	wantErr := errors.New("not found")
	err = run(
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

func TestFindWorkspaceUsesNearestModule(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	nestedWorkspace := filepath.Join(root, "projects", "nested")
	start := filepath.Join(nestedWorkspace, "package")
	for _, directory := range []string{root, nestedWorkspace} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(
			filepath.Join(directory, "MODULE.bazel"),
			[]byte("module(name = \"test\")\n"),
			0o644,
		); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(start, 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := findWorkspace(start)
	if err != nil {
		t.Fatalf("findWorkspace() error = %v", err)
	}
	if got != nestedWorkspace {
		t.Fatalf("findWorkspace() = %q, want %q", got, nestedWorkspace)
	}
}

func TestFindWorkspaceReportsMissingModule(t *testing.T) {
	t.Parallel()
	_, err := findWorkspace(string(filepath.Separator))
	if err == nil {
		t.Fatal("findWorkspace() error = nil, want an error")
	}
}
