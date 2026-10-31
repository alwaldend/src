package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/tools/ci/pkg/workspaces"
)

type commandRunner func(directory string, arguments ...string) error

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

func main() {
	root := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if root == "" || len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: bazel run --config=ci //tools/ci")
		os.Exit(1)
	}
	run := func(directory string, arguments ...string) error {
		cmd := exec.Command("bazel_agent", arguments...)
		cmd.Dir = directory
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return cmd.Run()
	}
	if err := check(root, run, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "CI failed:", err)
		os.Exit(1)
	}
}
