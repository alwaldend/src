package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/workspaces"
)

// commandRunner executes a literal argument vector in one workspace.
type commandRunner func(directory string, arguments ...string) error

// check attempts every workspace's build and test phases and combines failures.
func check(root string, run commandRunner, output io.Writer) error {
	paths, err := workspaces.Discover(root)
	if err != nil {
		return fmt.Errorf("discover CI workspaces: %w", err)
	}
	var failures []error
	for _, path := range paths {
		for _, phase := range []string{"build", "test"} {
			fmt.Fprintf(output, "CI %s: %s --config=ci //...\n", path, phase)
			if err := run(filepath.Join(root, path), "bazel", phase, "--config=ci", "//..."); err != nil {
				failures = append(failures, fmt.Errorf("%s %s: %w", path, phase, err))
			}
		}
	}
	return errors.Join(failures...)
}

// main runs repository CI with Bazel directly and propagates failed phases.
func main() {
	root := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if root == "" || len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: bazel run --config=ci //tools/ci")
		os.Exit(1)
	}
	run := func(directory string, arguments ...string) error {
		cmd := exec.CommandContext(context.Background(), arguments[0], arguments[1:]...)
		cmd.Dir = directory
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("execute Bazel: %w", err)
		}
		return nil
	}
	if err := check(root, run, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "CI failed:", err)
		os.Exit(1)
	}
}
