package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

type (
	lookPathFunc       func(string) (string, error)
	replaceProcessFunc func(string, []string, []string) error
)

const temporaryDirectory = "out/tmp"

func bazelArguments(args []string) []string {
	result := make([]string, 0, len(args)+2)
	result = append(result, "--batch")
	if len(args) == 0 {
		return result
	}
	result = append(result, args[0], "--config=agent")
	return append(result, args[1:]...)
}

func findWorkspace(start string) (string, error) {
	directory, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}
	for {
		module := filepath.Join(directory, "MODULE.bazel")
		if info, statErr := os.Stat(module); statErr == nil && !info.IsDir() {
			return directory, nil
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			return "", fmt.Errorf("find MODULE.bazel above %s", start)
		}
		directory = parent
	}
}

func withTemporaryDirectory(environment []string, directory string) []string {
	result := make([]string, 0, len(environment)+3)
	for _, entry := range environment {
		if strings.HasPrefix(entry, "TMPDIR=") ||
			strings.HasPrefix(entry, "TMP=") ||
			strings.HasPrefix(entry, "TEMP=") {
			continue
		}
		result = append(result, entry)
	}
	for _, name := range []string{"TMPDIR", "TMP", "TEMP"} {
		result = append(result, name+"="+directory)
	}
	return result
}

func run(
	args []string,
	environment []string,
	lookPath lookPathFunc,
	replaceProcess replaceProcessFunc,
) error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}
	workspace, err := findWorkspace(workingDirectory)
	if err != nil {
		return err
	}
	tmpDirectory := filepath.Join(workspace, temporaryDirectory)
	if err := os.MkdirAll(tmpDirectory, 0o755); err != nil {
		return fmt.Errorf("create temporary directory %s: %w", tmpDirectory, err)
	}
	environment = withTemporaryDirectory(environment, tmpDirectory)

	bazelPath, err := lookPath("bazel")
	if err != nil {
		return fmt.Errorf("find bazel in PATH: %w", err)
	}
	processArgs := append([]string{bazelPath}, bazelArguments(args)...)
	if err := replaceProcess(bazelPath, processArgs, environment); err != nil {
		return fmt.Errorf("execute bazel: %w", err)
	}
	return nil
}

func main() {
	err := run(os.Args[1:], os.Environ(), exec.LookPath, syscall.Exec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bazel_agent: %v\n", err)
		os.Exit(1)
	}
}
