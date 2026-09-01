package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
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
	{name: "projects/rules_hugo", path: "projects/rules_hugo"},
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
	status          string
	durationSeconds float64
	logPath         string
}

func resumedResult(candidate checkStep, phase jsonPhase) checkResult {
	return checkResult{
		workspace: candidate.workspace,
		phase:     candidate.phase,
		command:   []string{"bazel_agent", candidate.phase, "//..."},
		exitCode:  phase.Result.ExitCode,
		status:    phase.Result.Status,
		logPath:   phase.Result.Log,
	}
}

func isCompleted(phase jsonPhase) bool {
	return phase.Result != nil &&
		(phase.Result.Status == "pass" || phase.Result.Status == "fail")
}

func isUnexecuted(phase jsonPhase) bool {
	return phase.Status == "unexecuted"
}

func resumeResult(
	step checkStep,
	prior jsonPhase,
) checkResult {
	result := resumedResult(step, prior)
	result.durationSeconds = float64(prior.Result.DurationMS) / 1000
	result.logPath = prior.Result.Log
	return result
}

type jsonCheckResult struct {
	Workspace       string   `json:"workspace"`
	Phase           string   `json:"phase"`
	Command         []string `json:"command"`
	Status          string   `json:"status"`
	ExitCode        int      `json:"exitCode,omitempty"`
	DurationMS      int64    `json:"durationMs"`
	Log             string   `json:"log,omitempty"`
	TargetCount     int64    `json:"targetCount,omitempty"`
	UniverseReduced bool     `json:"universeReduced,omitempty"`
}

type jsonPhase struct {
	Workspace   string           `json:"workspace"`
	Phase       string           `json:"phase"`
	Status      string           `json:"status"`
	StartedAt   string           `json:"startedAt,omitempty"`
	CompletedAt string           `json:"completedAt,omitempty"`
	ElapsedMS   int64            `json:"elapsedMs,omitempty"`
	TargetCount int64            `json:"targetCount,omitempty"`
	Error       string           `json:"error,omitempty"`
	Result      *jsonCheckResult `json:"result,omitempty"`
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
		status:    "pending",
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
	if err != nil {
		fmt.Fprintf(progress, "could not create log %s: %v\n", logPath, err)
		result.status = "unexecuted"
		result.durationSeconds = time.Since(started).Seconds()
		return result
	}
	process := newCommand(command[0], command[1:]...)
	process.Dir = filepath.Join(repoRoot, candidate.path)
	process.Stdout = logFile
	process.Stderr = logFile
	runErr := process.Run()
	if runErr == nil {
		result.exitCode = 0
	} else {
		var exitError *exec.ExitError
		if errors.As(runErr, &exitError) {
			result.exitCode = exitError.ExitCode()
		} else {
			fmt.Fprintf(logFile, "could not execute bazel_agent: %v\n", runErr)
		}
	}
	if closeError := logFile.Close(); closeError != nil && runErr == nil {
		result.exitCode = 1
	}
	result.durationSeconds = time.Since(started).Seconds()
	if result.exitCode == 0 {
		result.status = "pass"
	} else {
		result.status = "fail"
	}

	fmt.Fprintf(
		progress,
		"[%s] %s %s: exit %d, %.1fs\n",
		strings.ToUpper(result.status),
		candidate.name,
		phase,
		result.exitCode,
		result.durationSeconds,
	)
	return result
}

func relativeLogPath(runDirectory string, result checkResult) (string, error) {
	relative, err := filepath.Rel(runDirectory, result.logPath)
	if err != nil {
		return "", fmt.Errorf("could not make log path relative: %w", err)
	}
	return filepath.ToSlash(relative), nil
}

func checkPlan() []checkStep {
	plan := make([]checkStep, 0, len(repositoryWorkspaces)*len(checkPhases))
	for _, candidate := range repositoryWorkspaces {
		for _, phase := range checkPhases {
			plan = append(plan, checkStep{
				workspace: candidate.name,
				phase:     phase,
			})
		}
	}
	return plan
}

func planInputs(
	workspaceNames []string,
	profile string,
	targetUniverses map[string]int64,
) runInputs {
	return runInputs{
		Workspaces:      workspaceNames,
		Profile:         profile,
		TargetUniverses: cloneCounts(targetUniverses),
	}
}

func digestInputs(inputs runInputs) (string, error) {
	canonical, err := json.Marshal(inputs)
	if err != nil {
		return "", fmt.Errorf("could not encode run inputs: %w", err)
	}
	sum := sha256.Sum256(canonical)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func cloneCounts(source map[string]int64) map[string]int64 {
	clone := make(map[string]int64, len(source))
	for key, value := range source {
		clone[key] = value
	}
	return clone
}

func runChecks(
	repoRoot string,
	runDirectory string,
	newCommand commandFactory,
	progress io.Writer,
	inputs runInputs,
	prior *checkReport,
	startedAt time.Time,
) ([]checkResult, *checkReport, error) {
	plan := checkPlan()
	if prior != nil {
		plan = planFromReport(prior)
	}
	trackingStatus := ""
	if prior != nil {
		trackingStatus = prior.State
	}
	results := make([]checkResult, 0, len(plan))
	priorByKey := map[string]jsonPhase{}
	if prior != nil {
		for _, phase := range prior.Phases {
			priorByKey[phase.Workspace+"\x00"+phase.Phase] = phase
		}
	}
	for _, step := range plan {
		key := step.workspace + "\x00" + step.phase
		priorPhase, hasPrior := priorByKey[key]
		var result checkResult
		switch {
		case hasPrior && isCompleted(priorPhase):
			result = resumeResult(step, priorPhase)
			fmt.Fprintf(
				progress,
				"[RESUME] %s %s: keeping %s result\n",
				step.workspace,
				step.phase,
				priorPhase.Result.Status,
			)
		case hasPrior && isUnexecuted(priorPhase):
			result = runCheck(
				repoRoot,
				runDirectory,
				workspaceFor(step.workspace),
				step.phase,
				newCommand,
				progress,
			)
		default:
			result = runCheck(
				repoRoot,
				runDirectory,
				workspaceFor(step.workspace),
				step.phase,
				newCommand,
				progress,
			)
		}
		results = append(results, result)
	}
	report := composeReport(results, inputs, startedAt, trackingStatus)
	if prior != nil && prior.State != reportComplete {
		report.ResumeMarker = prior.ID
		stillUnexecuted := false
		for _, phase := range report.Phases {
			if phase.Status == "unexecuted" {
				stillUnexecuted = true
				break
			}
		}
		if stillUnexecuted {
			report.State = reportUnfinished
		} else {
			report.State = reportComplete
		}
	}
	return results, &report, nil
}

func planFromReport(report *checkReport) []checkStep {
	plan := make([]checkStep, 0, len(report.Phases))
	for _, phase := range report.Phases {
		plan = append(plan, checkStep{
			workspace: phase.Workspace,
			phase:     phase.Phase,
		})
	}
	return plan
}

func workspaceFor(name string) workspace {
	for _, candidate := range repositoryWorkspaces {
		if candidate.name == name {
			return candidate
		}
	}
	for _, candidate := range repositoryWorkspaces {
		if candidate.path == name {
			return candidate
		}
	}
	return workspace{name: name, path: name}
}

func statusFor(result *checkResult) string {
	if result == nil {
		return "pending"
	}
	return result.status
}

func errorMessageFor(result *checkResult) string {
	if result == nil {
		return ""
	}
	return result.logPath + " resumed"
}

func composeReport(
	results []checkResult,
	inputs runInputs,
	startedAt time.Time,
	trackingStatus string,
) checkReport {
	phases := make([]jsonPhase, 0, len(results))
	for _, result := range results {
		completedAt := ""
		if result.status != "unexecuted" {
			completedAt = startedAt.Add(
				time.Duration(result.durationSeconds * float64(time.Second)),
			).UTC().Format(time.RFC3339Nano)
		}
		phase := jsonPhase{
			Workspace:   result.workspace,
			Phase:       result.phase,
			Status:      result.status,
			StartedAt:   startedAt.UTC().Format(time.RFC3339Nano),
			CompletedAt: completedAt,
			ElapsedMS:   int64(result.durationSeconds * 1000),
			TargetCount: inputs.TargetUniverses[result.workspace],
		}
		if result.status == "unexecuted" {
			phase.Error = "previously unexecuted; interrupted work is not inferred"
		} else {
			phase.Result = &jsonCheckResult{
				Workspace:  result.workspace,
				Phase:      result.phase,
				Command:    result.command,
				Status:     result.status,
				ExitCode:   result.exitCode,
				DurationMS: phase.ElapsedMS,
				Log:        result.logPath,
			}
		}
		phases = append(phases, phase)
	}
	counts := cloneCounts(inputs.TargetUniverses)
	total := int64(0)
	for _, count := range counts {
		total += count
	}
	state := reportComplete
	if trackingStatus == reportUnfinished {
		state = reportUnfinished
	}
	report := checkReport{
		APIVersion:     reportAPIVersion,
		Kind:           reportKind,
		ID:             "full-repo-check." + startedAt.UTC().Format("20060102T150405Z"),
		StartedAt:      startedAt.UTC().Format(time.RFC3339Nano),
		Inputs:         inputs,
		Workspaces:     inputs.Workspaces,
		Phases:         phases,
		TargetUniverse: total,
		State:          state,
	}
	return report
}

const (
	reportAPIVersion = "agents.alwaldend.com/full-repo-check/v1alpha1"
	reportKind       = "FullRepoCheckReport"
	reportUnfinished = "unfinished"
	reportComplete   = "complete"
)

func writeJSONReport(
	reportPath string,
	report checkReport,
) (string, error) {
	content, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("could not encode JSON report: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(reportPath), 0o700); err != nil {
		return "", fmt.Errorf("could not create JSON report directory: %w", err)
	}
	reportFile, err := os.OpenFile(
		reportPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0o600,
	)
	if err != nil {
		return "", fmt.Errorf("could not create JSON report: %w", err)
	}
	if _, err := reportFile.Write(append(content, '\n')); err != nil {
		reportFile.Close()
		return "", fmt.Errorf("could not write JSON report: %w", err)
	}
	if err := reportFile.Close(); err != nil {
		return "", fmt.Errorf("could not close JSON report: %w", err)
	}
	return reportPath, nil
}

func writeManifest(
	runDirectory string,
	reportPath string,
	inputs runInputs,
) (string, error) {
	digest, err := digestInputs(inputs)
	if err != nil {
		return "", err
	}
	manifest := struct {
		ReportPath string `json:"reportPath"`
		Inputs     runInputs
		Digest     string
	}{
		ReportPath: reportPath,
		Inputs:     inputs,
		Digest:     digest,
	}
	content, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("could not encode manifest: %w", err)
	}
	manifestPath := filepath.Join(runDirectory, "manifest.json")
	manifestFile, err := os.OpenFile(
		manifestPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0o600,
	)
	if err != nil {
		return "", fmt.Errorf("could not create manifest: %w", err)
	}
	if _, err := manifestFile.Write(append(content, '\n')); err != nil {
		manifestFile.Close()
		return "", fmt.Errorf("could not write manifest: %w", err)
	}
	if err := manifestFile.Close(); err != nil {
		return "", fmt.Errorf("could not close manifest: %w", err)
	}
	return manifestPath, nil
}

func readManifest(path string) (manifest, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return manifest{}, err
	}
	var result manifest
	if err := json.Unmarshal(content, &result); err != nil {
		return manifest{}, fmt.Errorf("could not decode manifest: %w", err)
	}
	return result, nil
}

type manifest struct {
	ReportPath string
	Inputs     runInputs
	Digest     string
}

func readCheckReport(path string) (*checkReport, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var report checkReport
	if err := json.Unmarshal(content, &report); err != nil {
		return nil, fmt.Errorf("could not decode check report: %w", err)
	}
	if report.APIVersion != reportAPIVersion || report.Kind != reportKind {
		return nil, fmt.Errorf(
			"unsupported report type %q %q",
			report.APIVersion,
			report.Kind,
		)
	}
	return &report, nil
}

func validateResume(report *checkReport, current runInputs) error {
	if report == nil {
		return nil
	}
	if report.State != reportUnfinished {
		return fmt.Errorf(
			"resume file %q is already complete; it cannot be resumed",
			report.ID,
		)
	}
	priorDigest, err := digestInputs(report.Inputs)
	if err != nil {
		return err
	}
	currentDigest, err := digestInputs(current)
	if err != nil {
		return err
	}
	if priorDigest != currentDigest {
		return fmt.Errorf(
			"resume inputs do not match: recorded %s, current %s",
			priorDigest,
			currentDigest,
		)
	}
	previousCount := report.TargetUniverse
	currentCount := int64(0)
	for _, count := range current.TargetUniverses {
		currentCount += count
	}
	if currentCount > previousCount {
		return fmt.Errorf(
			"target universe grew from %d to %d; permit with --allow-universe-growth",
			previousCount,
			currentCount,
		)
	}
	if currentCount != previousCount {
		return fmt.Errorf(
			"target universe changed from %d to %d; refusing to resume against "+
				"different inputs",
			previousCount,
			currentCount,
		)
	}
	return nil
}

type runOptions struct {
	jsonReportPath      string
	resumePath          string
	allowUniverseGrowth bool
}

func parseRunFlags(args []string) (runOptions, error) {
	flags := flag.NewFlagSet("run_full_repo_check", flag.ContinueOnError)
	var opts runOptions
	flags.StringVar(
		&opts.jsonReportPath,
		"json-report",
		"",
		"write the versioned JSON report to this path",
	)
	flags.StringVar(
		&opts.resumePath,
		"resume",
		"",
		"resume an unfinished report from this path",
	)
	flags.BoolVar(
		&opts.allowUniverseGrowth,
		"allow-universe-growth",
		false,
		"resume even when the target universe grew",
	)
	if err := flags.Parse(args); err != nil {
		return runOptions{}, err
	}
	if opts.resumePath == "" && opts.allowUniverseGrowth {
		return runOptions{}, fmt.Errorf(
			"--allow-universe-growth requires --resume",
		)
	}
	return opts, nil
}

// runInputs record the exact binding a check ran or will run against.
type runInputs struct {
	Workspaces      []string         `json:"workspaces"`
	Profile         string           `json:"profile,omitempty"`
	TargetUniverses map[string]int64 `json:"targetUniverses"`
}

type checkStep struct {
	workspace string
	phase     string
}

// checkReport is the versioned, resumable JSON report schema.
type checkReport struct {
	APIVersion     string      `json:"apiVersion"`
	Kind           string      `json:"kind"`
	ID             string      `json:"id"`
	StartedAt      string      `json:"startedAt,omitempty"`
	Inputs         runInputs   `json:"inputs"`
	Workspaces     []string    `json:"workspaces,omitempty"`
	Phases         []jsonPhase `json:"phases"`
	TargetUniverse int64       `json:"targetUniverse"`
	State          string      `json:"state"`
	ResumeMarker   string      `json:"resumeMarker,omitempty"`
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
		"applies the agent configuration. This excludes `manual`",
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
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
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

func workspacePaths(repoRoot string) []string {
	paths := make([]string, 0, len(repositoryWorkspaces))
	for _, candidate := range repositoryWorkspaces {
		if candidate.path == "." {
			paths = append(paths, filepath.Clean(repoRoot))
		} else {
			paths = append(
				paths,
				filepath.Join(repoRoot, filepath.FromSlash(candidate.path)),
			)
		}
	}
	return paths
}

func countTargets(
	newCommand commandFactory,
	candidate workspace,
	repoRoot string,
) (int64, error) {
	process := newCommand(
		"bazel_agent",
		"query",
		"//...",
		"--output=label_kind",
	)
	process.Dir = filepath.Join(repoRoot, candidate.path)
	var output strings.Builder
	process.Stdout = &output
	process.Stderr = io.Discard
	if err := process.Run(); err != nil {
		return 0, fmt.Errorf(
			"could not expand target universe for %s: %w",
			candidate.name,
			err,
		)
	}
	return int64(strings.Count(output.String(), "\n")), nil
}

func collectTargetUniverses(
	newCommand commandFactory,
	repoRoot string,
	progress io.Writer,
) (map[string]int64, error) {
	fmt.Fprintf(progress, "Expanding target universes with bazel_agent query ...\n")
	counts := make(map[string]int64, len(repositoryWorkspaces))
	for _, candidate := range repositoryWorkspaces {
		count, err := countTargets(newCommand, candidate, repoRoot)
		if err != nil {
			return nil, err
		}
		counts[candidate.name] = count
		fmt.Fprintf(progress, "%s: %d targets\n", candidate.name, count)
	}
	return counts, nil
}

func validateUniverseGrowth(
	report *checkReport,
	current runInputs,
	allow bool,
) error {
	if report == nil {
		return nil
	}
	previous := report.TargetUniverse
	currentTotal := int64(0)
	for _, count := range current.TargetUniverses {
		currentTotal += count
	}
	if currentTotal == previous {
		return nil
	}
	if currentTotal > previous {
		if !allow {
			return fmt.Errorf(
				"target universe grew from %d to %d; pass "+
					"--allow-universe-growth to allow it",
				previous,
				currentTotal,
			)
		}
		return nil
	}
	return fmt.Errorf(
		"target universe was reduced from %d to %d; refusing to resume",
		previous,
		currentTotal,
	)
}

func execute(
	getenv func(string) string,
	newCommand commandFactory,
	args []string,
	stdout io.Writer,
	stderr io.Writer,
) int {
	opts, err := parseRunFlags(args)
	if err != nil {
		fmt.Fprintln(stderr, "could not parse flags:", err)
		return 2
	}
	repoRoot, err := repositoryRoot(getenv)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := validateWorkspaces(repoRoot); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	startedAt := time.Now()

	var prior *checkReport
	if opts.resumePath != "" {
		prior, err = readCheckReport(opts.resumePath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}

	inputs := runInputs{
		Workspaces:      workspaceNames(),
		Profile:         "full",
		TargetUniverses: nil,
	}
	if prior != nil {
		inputs = prior.Inputs
	}
	universes, err := collectTargetUniverses(
		newCommand,
		repoRoot,
		stdout,
	)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	inputs.TargetUniverses = universes
	if err := validateUniverseGrowth(prior, inputs, opts.allowUniverseGrowth); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := validateResume(prior, inputs); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	runDirectory, err := createRunDirectory(repoRoot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	results, report, err := runChecks(
		repoRoot,
		runDirectory,
		newCommand,
		stdout,
		inputs,
		prior,
		startedAt,
	)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	reportPath, err := writeReport(runDirectory, results, startedAt)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Report: %s\n", reportPath)
	if opts.jsonReportPath != "" {
		if _, err := writeJSONReport(
			opts.jsonReportPath,
			*report,
		); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	if _, err := writeManifest(runDirectory, reportPath, inputs); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	allPassed := report.State == reportComplete
	for _, result := range results {
		if result.status == "fail" || result.status == "unexecuted" {
			allPassed = false
		}
	}
	if !allPassed {
		return 1
	}
	return 0
}

func workspaceNames() []string {
	names := make([]string, 0, len(repositoryWorkspaces))
	for _, candidate := range repositoryWorkspaces {
		names = append(names, candidate.name)
	}
	return names
}

func main() {
	os.Exit(execute(os.Getenv, exec.Command, os.Args[1:], os.Stdout, os.Stderr))
}
