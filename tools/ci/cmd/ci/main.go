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

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/workspaces"
)

// commandRunner executes a literal argument vector in one workspace.
type commandRunner func(directory string, arguments ...string) error

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

// workspaceResult records the CI outcomes in one Bazel workspace.
type workspaceResult struct {
	// name is the workspace's display name.
	name string
	// build holds the outcome of bazel build.
	build phaseResult
	// test holds the outcome of bazel test.
	test phaseResult
}

// check attempts every workspace's build and test phases and combines failures.
func check(root string, run commandRunner, output io.Writer, grouped bool) ([]workspaceResult, error) {
	paths, err := workspaces.Discover(root)
	if err != nil {
		return nil, fmt.Errorf("discover CI workspaces: %w", err)
	}
	results := make([]workspaceResult, 0, len(paths))
	var failures []error
	for _, path := range paths {
		name := strings.ReplaceAll(path, "\r", " ")
		if path == "." {
			name = "Root"
		}
		runPhase := func(phase, title string) phaseResult {
			if grouped {
				fmt.Fprintf(output, "##[group]%s / %s\n", name, title)
				// Start the marker on a new line even if Bazel left a partial line.
				defer fmt.Fprint(output, "\n##[endgroup]\n")
			}
			fmt.Fprintf(output, "CI %s: %s --config=ci //...\n", path, phase)
			start := time.Now()
			err := run(filepath.Join(root, path), "bazel", phase, "--config=ci", "//...")
			result := phaseResult{duration: time.Since(start)}
			if err != nil {
				result.err = fmt.Errorf("%s %s: %w", path, phase, err)
				failures = append(failures, result.err)
			}
			return result
		}
		build := runPhase("build", "Build")
		test := runPhase("test", "Test")
		results = append(results, workspaceResult{name: name, build: build, test: test})
	}
	return results, errors.Join(failures...)
}

// appendSummary writes the Markdown produced by CI to the current step's summary.
func appendSummary(path string, results []workspaceResult) error {
	var summary strings.Builder
	summary.WriteString("\n## Repository CI\n\n")
	if len(results) == 0 {
		summary.WriteString("Workspace discovery failed. See the job log for details.\n")
	} else {
		escapeCell := strings.NewReplacer(
			"&", "&amp;", "<", "&lt;", ">", "&gt;", "|", "&#124;", "\\", "&#92;",
			"`", "&#96;", "*", "&#42;", "_", "&#95;", "~", "&#126;", "[", "&#91;", "]", "&#93;",
		)
		summary.WriteString("| Workspace | Build | Test |\n| --- | --- | --- |\n")
		for _, result := range results {
			fmt.Fprintf(&summary, "| %s | %s | %s |\n", escapeCell.Replace(result.name), result.build.summaryCell(), result.test.summaryCell())
		}
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if err != nil {
		return fmt.Errorf("open step summary: %w", err)
	}
	_, writeErr := io.WriteString(file, summary.String())
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
	run := func(directory string, arguments ...string) error {
		cmd := exec.CommandContext(context.Background(), arguments[0], arguments[1:]...)
		cmd.Dir = directory
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if grouped {
			cmd.Stderr = cmd.Stdout
		}
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("execute Bazel: %w", err)
		}
		return nil
	}
	results, err := check(root, run, os.Stdout, grouped)
	if grouped {
		summaryPath := os.Getenv("FORGEJO_STEP_SUMMARY")
		if summaryPath == "" {
			summaryPath = os.Getenv("GITHUB_STEP_SUMMARY")
		}
		if summaryPath != "" {
			if summaryErr := appendSummary(summaryPath, results); summaryErr != nil {
				fmt.Fprintln(os.Stderr, "CI could not write the job summary:", summaryErr)
			}
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "CI failed:", err)
		os.Exit(1)
	}
}
