package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type bazelExitError struct {
	err  error
	code int
}

func (e *bazelExitError) Error() string {
	return fmt.Sprintf("run stamped Bazel command: %v", e.err)
}

func (e *bazelExitError) Unwrap() error {
	return e.err
}

func (e *bazelExitError) ExitCode() int {
	return e.code
}

func launchBazel(
	repository string,
	trunk string,
	channel string,
	release string,
	args []string,
) error {
	if len(args) == 0 {
		return fmt.Errorf("bazel requires a Bazel command after --")
	}
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve versioning executable: %w", err)
	}
	executable, err = filepath.Abs(executable)
	if err != nil {
		return fmt.Errorf("make versioning executable absolute: %w", err)
	}
	repository, err = filepath.Abs(repository)
	if err != nil {
		return fmt.Errorf("make repository path absolute: %w", err)
	}
	statusCommand := strings.Join([]string{
		shellQuote(executable),
		"--repo", shellQuote(repository),
		"--trunk", shellQuote(trunk),
		"--channel", shellQuote(channel),
	}, " ")
	if release != "" {
		statusCommand += " --release " + shellQuote(release)
	}
	statusCommand += " bazel-status"
	arguments := insertBeforeSeparator(
		args,
		"--workspace_status_command="+statusCommand,
	)
	command := exec.Command("bazel_agent", arguments...)
	command.Dir = repository
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() > 0 {
			return &bazelExitError{err: err, code: exitError.ExitCode()}
		}
		return fmt.Errorf("run stamped Bazel command: %w", err)
	}
	return nil
}

func insertBeforeSeparator(args []string, value string) []string {
	result := make([]string, 0, len(args)+1)
	for index, argument := range args {
		if argument == "--" {
			result = append(result, value)
			result = append(result, args[index:]...)
			return result
		}
		result = append(result, argument)
	}
	return append(result, value)
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}
