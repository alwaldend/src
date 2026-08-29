package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type workspace struct {
	name string
	path string
}

var repositoryWorkspaces = []workspace{
	{name: "root", path: "."},
	{
		name: "projects/rules_binary_toolchain",
		path: "projects/rules_binary_toolchain",
	},
	{name: "projects/rules_docs", path: "projects/rules_docs"},
	{name: "projects/rules_docs_gazelle", path: "projects/rules_docs_gazelle"},
	{name: "projects/rules_promptfoo", path: "projects/rules_promptfoo"},
	{
		name: "projects/rules_promptfoo_gazelle",
		path: "projects/rules_promptfoo_gazelle",
	},
	{name: "projects/rules_skill", path: "projects/rules_skill"},
	{
		name: "projects/rules_skill_gazelle",
		path: "projects/rules_skill_gazelle",
	},
	{name: "projects/rules_template", path: "projects/rules_template"},
}

var checkPhases = []string{"build", "test"}

type checkResult struct {
	workspace       string
	phase           string
	command         []string
	exitCode        int
	durationSeconds float64
	logPath         string
}

type commandFactory func(name string, args ...string) *exec.Cmd

func repositoryRoot(getenv func(string) string) (string, error) {
	root := getenv("BUILD_WORKSPACE_DIRECTORY")
	if root == "" {
		return "", errors.New(
			"BUILD_WORKSPACE_DIRECTORY is unset; use the generated Bazel launcher",
		)
	}
	absolute, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("could not resolve repository root: %w", err)
	}
	return filepath.Clean(absolute), nil
}

func validateWorkspaces(repoRoot string) error {
	var missing []string
	for _, candidate := range repositoryWorkspaces {
		modulePath := filepath.Join(repoRoot, candidate.path, "MODULE.bazel")
		info, err := os.Stat(modulePath)
		if err != nil || !info.Mode().IsRegular() {
			missing = append(
				missing,
				filepath.Join(candidate.path, "MODULE.bazel"),
			)
		}
	}
	if len(missing) != 0 {
		return fmt.Errorf(
			"missing Bazel workspace module(s): %s",
			strings.Join(missing, ", "),
		)
	}
	return nil
}

func privateDirectory(path string) error {
	if err := os.MkdirAll(path, 0o700); err != nil {
		return err
	}
	return os.Chmod(path, 0o700)
}

func createRunDirectory(repoRoot string) (string, error) {
	outputRoot := filepath.Join(repoRoot, "out", "full-repo-check")
	if err := privateDirectory(outputRoot); err != nil {
		return "", fmt.Errorf("could not create output directory: %w", err)
	}
	runDirectory, err := os.MkdirTemp(outputRoot, "run.")
	if err != nil {
		return "", fmt.Errorf("could not create run directory: %w", err)
	}
	logsDirectory := filepath.Join(runDirectory, "logs")
	if err := privateDirectory(logsDirectory); err != nil {
		return "", fmt.Errorf("could not create log directory: %w", err)
	}
	return runDirectory, nil
}

func displayCommand(result checkResult) string {
	return strings.Join(result.command, " ")
}

func runCheck(
	repoRoot string,
	runDirectory string,
	candidate workspace,
	phase string,
	newCommand commandFactory,
	progress io.Writer,
) checkResult {
	command := []string{"bazel_agent", phase, "//..."}
	safeWorkspace := strings.ReplaceAll(candidate.name, "/", "__")
	logPath := filepath.Join(
		runDirectory,
		"logs",
		fmt.Sprintf("%s.%s.log", safeWorkspace, phase),
	)
	result := checkResult{
		workspace: candidate.name,
		phase:     phase,
		command:   command,
		exitCode:  127,
		logPath:   logPath,
	}

	fmt.Fprintf(
		progress,
		"[START] %s %s: %s\n",
		candidate.name,
		phase,
		displayCommand(result),
	)
	started := time.Now()
	logFile, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_EXCL|os.O_WRONLY,
		0o600,
	)
	if err == nil {
		process := newCommand(command[0], command[1:]...)
		process.Dir = filepath.Join(repoRoot, candidate.path)
		process.Stdout = logFile
		process.Stderr = logFile
		err = process.Run()
		switch {
		case err == nil:
			result.exitCode = 0
		default:
			var exitError *exec.ExitError
			if errors.As(err, &exitError) {
				result.exitCode = exitError.ExitCode()
			} else {
				fmt.Fprintf(logFile, "could not execute bazel_agent: %v\n", err)
			}
		}
		if closeError := logFile.Close(); closeError != nil && err == nil {
			result.exitCode = 1
		}
	} else {
		fmt.Fprintf(progress, "could not create log %s: %v\n", logPath, err)
	}
	result.durationSeconds = time.Since(started).Seconds()

	status := "PASS"
	if result.exitCode != 0 {
		status = "FAIL"
	}
	fmt.Fprintf(
		progress,
		"[%s] %s %s: exit %d, %.1fs\n",
		status,
		candidate.name,
		phase,
		result.exitCode,
		result.durationSeconds,
	)
	return result
}

func runChecks(
	repoRoot string,
	runDirectory string,
	newCommand commandFactory,
	progress io.Writer,
) []checkResult {
	results := make([]checkResult, 0, len(repositoryWorkspaces)*len(checkPhases))
	for _, candidate := range repositoryWorkspaces {
		for _, phase := range checkPhases {
			results = append(
				results,
				runCheck(
					repoRoot,
					runDirectory,
					candidate,
					phase,
					newCommand,
					progress,
				),
			)
		}
	}
	return results
}

func relativeLogPath(runDirectory string, result checkResult) (string, error) {
	relative, err := filepath.Rel(runDirectory, result.logPath)
	if err != nil {
		return "", fmt.Errorf("could not make log path relative: %w", err)
	}
	return filepath.ToSlash(relative), nil
}

func writeReport(
	runDirectory string,
	results []checkResult,
	generatedAt time.Time,
) (string, error) {
	var report strings.Builder
	fmt.Fprintln(&report, "# Full repository check")
	fmt.Fprintln(&report)
	fmt.Fprintf(
		&report,
		"Generated at: `%s`\n\n",
		generatedAt.UTC().Format(time.RFC3339),
	)
	fmt.Fprintln(
		&report,
		"Scope: normal `//...` expansion through `bazel_agent`. The runner",
	)
	fmt.Fprintln(
		&report,
		"enables batch mode and the agent configuration. This excludes `manual`",
	)
	fmt.Fprintln(
		&report,
		"targets, incompatible targets, and optional configuration matrices.",
	)
	fmt.Fprintln(&report)
	fmt.Fprintln(
		&report,
		"| Workspace | Phase | Command | Result | Duration | Log |",
	)
	fmt.Fprintln(&report, "|---|---|---|---|---:|---|")
	for _, result := range results {
		status := "PASS"
		if result.exitCode != 0 {
			status = "FAIL"
		}
		relativeLog, err := relativeLogPath(runDirectory, result)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(
			&report,
			"| %s | %s | `%s` | %s (exit %d) | %.1fs | [%s](%s) |\n",
			result.workspace,
			result.phase,
			displayCommand(result),
			status,
			result.exitCode,
			result.durationSeconds,
			relativeLog,
			relativeLog,
		)
	}

	fmt.Fprintln(&report)
	fmt.Fprintln(&report, "## Failed commands")
	fmt.Fprintln(&report)
	fmt.Fprintln(
		&report,
		"| Workspace | Phase | Exit code | Diagnostic log |",
	)
	fmt.Fprintln(&report, "|---|---|---:|---|")
	failures := 0
	for _, result := range results {
		if result.exitCode == 0 {
			continue
		}
		failures++
		relativeLog, err := relativeLogPath(runDirectory, result)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(
			&report,
			"| %s | %s | %d | [%s](%s) |\n",
			result.workspace,
			result.phase,
			result.exitCode,
			relativeLog,
			relativeLog,
		)
	}
	if failures == 0 {
		fmt.Fprintln(&report, "| None | None | 0 | None |")
	}

	reportPath := filepath.Join(runDirectory, "report.md")
	reportFile, err := os.OpenFile(
		reportPath,
		os.O_CREATE|os.O_EXCL|os.O_WRONLY,
		0o600,
	)
	if err != nil {
		return "", fmt.Errorf("could not create report: %w", err)
	}
	if _, err := io.WriteString(reportFile, report.String()); err != nil {
		reportFile.Close()
		return "", fmt.Errorf("could not write report: %w", err)
	}
	if err := reportFile.Close(); err != nil {
		return "", fmt.Errorf("could not close report: %w", err)
	}
	return reportPath, nil
}

func execute(
	getenv func(string) string,
	newCommand commandFactory,
	stdout io.Writer,
	stderr io.Writer,
) int {
	repoRoot, err := repositoryRoot(getenv)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := validateWorkspaces(repoRoot); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	runDirectory, err := createRunDirectory(repoRoot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	results := runChecks(
		repoRoot,
		runDirectory,
		newCommand,
		stdout,
	)
	reportPath, err := writeReport(runDirectory, results, time.Now())
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Report: %s\n", reportPath)
	for _, result := range results {
		if result.exitCode != 0 {
			return 1
		}
	}
	return 0
}

func main() {
	os.Exit(execute(os.Getenv, exec.Command, os.Stdout, os.Stderr))
}
