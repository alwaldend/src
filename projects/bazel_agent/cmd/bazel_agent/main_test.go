package main

import (
	"errors"
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
	t.Parallel()
	environment := []string{"PATH=/bin"}
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
		"--batch",
		"test",
		"--config=agent",
		"//pkg:all",
	}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %q, want %q", gotArgs, wantArgs)
	}
	if !reflect.DeepEqual(gotEnvironment, environment) {
		t.Fatalf("environment = %q, want %q", gotEnvironment, environment)
	}
}

func TestRunReportsLookupFailure(t *testing.T) {
	t.Parallel()
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
