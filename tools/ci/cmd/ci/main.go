package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// The root ignore file owns the boundaries of independently built workspaces.
func workspaces(root string) ([]string, error) {
	if _, err := os.Stat(filepath.Join(root, "MODULE.bazel")); err != nil {
		return nil, fmt.Errorf("repository MODULE.bazel: %w", err)
	}
	file, err := os.Open(filepath.Join(root, ".bazelignore"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	paths := []string{"."}
	seen := map[string]bool{".": true}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		path := strings.TrimSpace(scanner.Text())
		if path == "" || strings.HasPrefix(path, "#") {
			continue
		}
		path = filepath.Clean(path)
		if !filepath.IsLocal(path) {
			return nil, fmt.Errorf("nonlocal workspace boundary %q", path)
		}
		if seen[path] {
			continue
		}
		info, err := os.Stat(filepath.Join(root, path, "MODULE.bazel"))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("invalid MODULE.bazel in %s", path)
		}
		paths = append(paths, path)
		seen[path] = true
	}
	return paths, scanner.Err()
}

type commandRunner func(directory string, arguments ...string) error

func check(root string, run commandRunner, output io.Writer) error {
	paths, err := workspaces(root)
	if err != nil {
		return err
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
