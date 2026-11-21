package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// phaseResult records the outcome of one Bazel invocation.
type phaseResult struct {
	// duration is the elapsed time of the invocation.
	duration time.Duration
	// err is the command failure with workspace and phase context.
	err error
}

// summaryCell formats the phase result for the Forgejo job summary.
func (result phaseResult) summaryCell() string {
	status := "Passed"
	if result.err != nil {
		status = "**Failed**"
	}
	return fmt.Sprintf("%s (%s)", status, result.duration.Round(time.Millisecond))
}

// ciResult records the root workspace's CI outcomes.
type ciResult struct {
	// build holds the outcome of bazel build.
	build phaseResult
	// test holds the outcome of bazel test.
	test phaseResult
	// lint holds the outcome of the repository's semantic lint profile.
	lint phaseResult
	// shutdownErr records a failure to retire the Forgejo job's Bazel server.
	shutdownErr error
}

// check runs the root CI phases in order and combines failures.
func check(root string, output io.Writer, grouped bool) (*ciResult, error) {
	info, err := os.Stat(filepath.Join(root, "MODULE.bazel"))
	if err != nil {
		return nil, fmt.Errorf("inspect repository MODULE.bazel: %w", err)
	}
	if !info.Mode().IsRegular() {
		return nil, errors.New("repository MODULE.bazel is not a regular file")
	}
	run := func(arguments ...string) error {
		cmd := exec.CommandContext(context.Background(), "bazel", arguments...)
		cmd.Dir = root
		cmd.Stdin = os.Stdin
		if grouped && arguments[0] != "shutdown" {
			if err := runWithForgejoLogs(cmd, output, arguments[0] == "test"); err != nil {
				return fmt.Errorf("execute Bazel: %w", err)
			}
			return nil
		}
		cmd.Stdout, cmd.Stderr = output, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("execute Bazel: %w", err)
		}
		return nil
	}
	runPhase := func(title, command string, options ...string) phaseResult {
		if grouped {
			// Bazel can leave a terminal escape sequence without a newline.
			fmt.Fprintf(output, "\n##[group]Root / %s\n", title)
			defer fmt.Fprint(output, "\n##[endgroup]\n")
		}
		arguments := append([]string{command, "--config=ci"}, options...)
		arguments = append(arguments, "//...")
		fmt.Fprintf(output, "CI %s\n", strings.Join(arguments, " "))
		start := time.Now()
		err := run(arguments...)
		result := phaseResult{duration: time.Since(start)}
		if err != nil {
			result.err = fmt.Errorf("repository %s: %w", strings.ToLower(title), err)
		}
		return result
	}
	result := &ciResult{
		build: runPhase("Build", "build"),
		test:  runPhase("Test", "test"),
		lint:  runPhase("Lint", "build", "--config=lint"),
	}
	if grouped {
		if err := run("shutdown"); err != nil {
			result.shutdownErr = fmt.Errorf("shut down repository Bazel server: %w", err)
		}
	}
	return result, errors.Join(result.build.err, result.test.err, result.lint.err, result.shutdownErr)
}

// appendSummary writes the Markdown produced by CI to the current step's summary.
func appendSummary(path string, result *ciResult) error {
	summary := "\n## Repository CI\n\n"
	if result == nil {
		summary += "Repository workspace unavailable. See the job log for details.\n"
	} else {
		summary += "| Workspace | Build | Test | Lint |\n| --- | --- | --- | --- |\n"
		summary += fmt.Sprintf("| Root | %s | %s | %s |\n", result.build.summaryCell(), result.test.summaryCell(), result.lint.summaryCell())
		if result.shutdownErr != nil {
			summary += "\nBazel server shutdown failed in Root. See the job log for details.\n"
		}
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if err != nil {
		return fmt.Errorf("open step summary: %w", err)
	}
	_, writeErr := io.WriteString(file, summary)
	if err := errors.Join(writeErr, file.Close()); err != nil {
		return fmt.Errorf("append step summary: %w", err)
	}
	return nil
}

// main runs repository CI with Bazel directly and propagates failed phases.
func main() {
	root := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if root == "" || len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: bazel run --config=ci //tools/ci")
		os.Exit(1)
	}
	grouped := os.Getenv("FORGEJO_ACTIONS") == "true"
	result, err := check(root, os.Stdout, grouped)
	if grouped {
		summaryPath := os.Getenv("FORGEJO_STEP_SUMMARY")
		if summaryPath == "" {
			summaryPath = os.Getenv("GITHUB_STEP_SUMMARY")
		}
		if summaryPath != "" {
			if summaryErr := appendSummary(summaryPath, result); summaryErr != nil {
				fmt.Fprintln(os.Stderr, "CI could not write the job summary:", summaryErr)
			}
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "CI failed:", err)
		os.Exit(1)
	}
}
