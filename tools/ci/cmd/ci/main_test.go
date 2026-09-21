package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for path, content := range map[string]string{
		"MODULE.bazel":               "",
		".bazelignore":               "# Independent modules\nout\nnode_modules\ntools/example\ntools/example\n",
		"tools/example/MODULE.bazel": "",
	} {
		path = filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestBuildFailureStillTestsEveryWorkspace(t *testing.T) {
	root := fixture(t)
	failure := errors.New("build failed")
	var calls []string
	run := func(directory string, arguments ...string) error {
		relative, err := filepath.Rel(root, directory)
		if err != nil {
			t.Fatal(err)
		}
		calls = append(calls, relative+": "+strings.Join(arguments, " "))
		if len(calls) == 1 {
			return failure
		}
		return nil
	}
	var output bytes.Buffer
	if err := check(root, run, &output); !errors.Is(err, failure) {
		t.Fatalf("lost build failure: %v", err)
	}
	want := []string{
		".: bazel build --config=ci //...",
		".: bazel test --config=ci //...",
		"tools/example: bazel build --config=ci //...",
		"tools/example: bazel test --config=ci //...",
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %q, want %q", calls, want)
	}
}

func TestSuccessfulChecks(t *testing.T) {
	if err := check(fixture(t), func(string, ...string) error { return nil }, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
}

func TestInvalidWorkspaceDoesNotStartCommands(t *testing.T) {
	for _, ignore := range []string{"../outside\n", "/outside\n"} {
		root := fixture(t)
		if err := os.WriteFile(filepath.Join(root, ".bazelignore"), []byte(ignore), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := check(root, func(string, ...string) error { t.Fatal("started command"); return nil }, &bytes.Buffer{}); err == nil {
			t.Fatal("accepted invalid boundary")
		}
	}
}
