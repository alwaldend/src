package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	validationPlanSchema  = "repo_delivery/validation_plan/v1"
	validationStateSchema = "repo_delivery/validation/v1"
)

type validationCheck struct {
	Workspace      string   `json:"workspace"`
	Kind           string   `json:"kind"`
	Targets        []string `json:"targets"`
	TimeoutSeconds int      `json:"timeout_seconds"`
}

type validationGapDecision struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

type validationPlan struct {
	Schema       string                  `json:"schema"`
	Checks       []validationCheck       `json:"checks"`
	GapDecisions []validationGapDecision `json:"gap_decisions"`
}

type validationCheckResult struct {
	Index     int    `json:"index"`
	Status    string `json:"status"`
	ExitCode  int    `json:"exit_code"`
	Truncated bool   `json:"truncated"`
	LogFile   string `json:"log_file"`
	LogDigest string `json:"log_sha256"`
}

type validationState struct {
	Schema                string                  `json:"schema"`
	Status                string                  `json:"status"`
	ReceiptRevision       string                  `json:"receipt_revision"`
	RepositoryFingerprint string                  `json:"repository_fingerprint"`
	HeadOID               string                  `json:"head_oid"`
	TreeOID               string                  `json:"tree_oid"`
	PlanFile              string                  `json:"plan_file"`
	PlanDigest            string                  `json:"plan_sha256"`
	EnvironmentDigest     string                  `json:"environment_sha256"`
	StartedAt             string                  `json:"started_at"`
	CompletedAt           string                  `json:"completed_at"`
	Results               []validationCheckResult `json:"results"`
}

type validationReport struct {
	Status      string                  `json:"status"`
	StateFile   string                  `json:"state_file"`
	HeadOID     string                  `json:"head_oid"`
	TreeOID     string                  `json:"tree_oid"`
	NextAction  string                  `json:"next_action"`
	Results     []validationCheckResult `json:"results"`
	Publication *publishReport          `json:"publication,omitempty"`
}

func newValidateCommand(ctx context.Context, config *deliveryConfig,
	getenv func(string) string, stdout io.Writer, runner commandRunner,
) *cobra.Command {
	var receiptFile, planFile string
	command := &cobra.Command{
		Use: "validate", Short: "Run an explicit bounded validation plan for a prepared candidate",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			delivery, err := deliveryFromConfig(ctx, config, getenv, runner)
			if err != nil {
				return err
			}
			report, err := delivery.validateCandidate(ctx, receiptFile, planFile)
			if report != nil {
				if outputErr := writeJSON(stdout, report); outputErr != nil {
					return outputErr
				}
			}
			return err
		},
	}
	command.Flags().StringVar(&receiptFile, "receipt-file", "", "ignored preparation receipt")
	command.Flags().StringVar(&planFile, "plan-file", "", "explicit validation plan in the same ignored out/<task>/ directory")
	_ = command.MarkFlagRequired("receipt-file")
	_ = command.MarkFlagRequired("plan-file")
	return command
}

func newContinueCommand(ctx context.Context, config *deliveryConfig,
	getenv func(string) string, stdout io.Writer, runner commandRunner,
) *cobra.Command {
	var receiptFile string
	var publish, noPullRequest bool
	command := &cobra.Command{
		Use: "continue", Short: "Inspect validation, explicitly publish, or verify an attempted publication",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if noPullRequest && !publish {
				return fmt.Errorf("--no-pull-request requires --publish")
			}
			delivery, err := deliveryFromConfig(ctx, config, getenv, runner)
			if err != nil {
				return err
			}
			report, err := delivery.continueCandidate(ctx, receiptFile, publish, noPullRequest)
			if report != nil {
				if outputErr := writeJSON(stdout, report); outputErr != nil {
					return outputErr
				}
			}
			return err
		},
	}
	command.Flags().StringVar(&receiptFile, "receipt-file", "", "ignored preparation receipt with adjacent validation results")
	command.Flags().BoolVar(&publish, "publish", false, "explicitly publish the recorded validated candidate once")
	command.Flags().BoolVar(&noPullRequest, "no-pull-request", false, "publish without creating a pull request when none exists")
	_ = command.MarkFlagRequired("receipt-file")
	return command
}

func validationStatePath(receiptFile string) string { return receiptFile + ".validation.json" }

func (s validationState) report(receiptFile string) *validationReport {
	next := "diagnose the recorded result before retrying"
	switch s.Status {
	case "validated":
		next = "continue --receipt-file " + receiptFile + " --publish"
	case "revalidation_required":
		next = "validate --receipt-file " + receiptFile + " --plan-file " + s.PlanFile
	case "publication_attempted":
		next = "continue --receipt-file " + receiptFile + " (verification only; no automatic re-push)"
	case "verified":
		next = "complete review and handoff"
	}
	return &validationReport{
		Status: s.Status, StateFile: validationStatePath(receiptFile),
		HeadOID: s.HeadOID, TreeOID: s.TreeOID, NextAction: next, Results: s.Results,
	}
}

func validationPlanShape() *canonicalJSONShape {
	s := &canonicalJSONShape{scalarKind: "string"}
	return &canonicalJSONShape{fields: map[string]*canonicalJSONShape{
		"schema": s,
		"checks": {element: &canonicalJSONShape{fields: map[string]*canonicalJSONShape{
			"workspace": s, "kind": s, "targets": {element: s},
			"timeout_seconds": {scalarKind: "number"},
		}}},
		"gap_decisions": {element: &canonicalJSONShape{fields: map[string]*canonicalJSONShape{
			"path": s, "reason": s,
		}}},
	}}
}

func validationStateShape() *canonicalJSONShape {
	s := &canonicalJSONShape{scalarKind: "string"}
	n := &canonicalJSONShape{scalarKind: "number"}
	return &canonicalJSONShape{fields: map[string]*canonicalJSONShape{
		"schema": s, "status": s, "receipt_revision": s, "repository_fingerprint": s,
		"head_oid": s, "tree_oid": s, "plan_file": s, "plan_sha256": s, "environment_sha256": s,
		"started_at": s, "completed_at": s,
		"results": {element: &canonicalJSONShape{fields: map[string]*canonicalJSONShape{
			"index": n, "status": s, "exit_code": n, "truncated": {scalarKind: "boolean"},
			"log_file": s, "log_sha256": s,
		}}},
	}}
}

func decodeValidationJSON(contents []byte, shape *canonicalJSONShape, value any) error {
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.UseNumber()
	if err := validateCanonicalJSONValue(decoder, shape, "validation"); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return fmt.Errorf("validation contains trailing JSON")
	}
	return json.Unmarshal(contents, value)
}

func (d *delivery) readValidationPlan(ctx context.Context, path string) (validationPlan, string, error) {
	var plan validationPlan
	absolute, err := d.receiptPath(ctx, path, true)
	if err != nil {
		return plan, "", err
	}
	contents, err := readStableReceiptFile(absolute)
	if err != nil {
		return plan, "", err
	}
	if err := decodeValidationJSON(contents, validationPlanShape(), &plan); err != nil {
		return plan, "", err
	}
	return plan, digestStrings(validationPlanSchema, string(contents)), nil
}

var validationLabelPattern = regexp.MustCompile(`^//([A-Za-z0-9_.+\-/]+)?:[A-Za-z0-9_.+\-/]+$`)

func (d *delivery) checkValidationPlan(plan validationPlan, receipt preparationReceipt) error {
	if plan.Schema != validationPlanSchema || len(plan.Checks) == 0 || len(plan.Checks) > 32 {
		return fmt.Errorf("validation plan requires schema %s and 1–32 checks", validationPlanSchema)
	}
	quality := false
	linted := make(map[string]bool)
	for _, check := range plan.Checks {
		if check.Kind != "test" && check.Kind != "build" && check.Kind != "lint" {
			return fmt.Errorf("unsupported validation kind %q", check.Kind)
		}
		if check.TimeoutSeconds < 1 || check.TimeoutSeconds > 3600 || len(check.Targets) == 0 || len(check.Targets) > 128 {
			return fmt.Errorf("validation checks require 1–128 targets and timeout_seconds between 1 and 3600")
		}
		if _, err := d.validationWorkspace(check.Workspace); err != nil {
			return err
		}
		seen := make(map[string]bool)
		for _, target := range check.Targets {
			if !validationLabelPattern.MatchString(target) || strings.Contains(target, "..") || target == "//:all" || seen[target] {
				return fmt.Errorf("validation requires unique explicit //package:target labels, not %q", target)
			}
			seen[target] = true
			if check.Workspace == "." && check.Kind == "test" && target == "//:repo_quality_test" {
				quality = true
			}
			if check.Workspace == "." && check.Kind == "lint" {
				linted[target] = true
			}
		}
	}
	if !quality {
		return fmt.Errorf("validation plan must include the root //:repo_quality_test test")
	}
	selection, err := bazelValidationForPaths(d.repository.directory, receipt.Scope.AggregatePaths)
	if err != nil {
		return err
	}
	for _, label := range selection.Labels {
		if !linted[label] {
			return fmt.Errorf("validation plan must lint affected package %s", label)
		}
	}
	decisions := make(map[string]bool)
	for _, decision := range plan.GapDecisions {
		if strings.TrimSpace(decision.Reason) == "" || len(decision.Reason) > 4096 || decisions[decision.Path] {
			return fmt.Errorf("gap decisions require a unique path and a nonempty bounded reason")
		}
		decisions[decision.Path] = true
	}
	for _, gap := range selection.Gaps {
		if !decisions[gap.Path] {
			return fmt.Errorf("validation plan requires a gap decision for %s (%s)", gap.Path, gap.Reason)
		}
		delete(decisions, gap.Path)
	}
	if len(decisions) != 0 {
		return fmt.Errorf("validation plan contains decisions for paths outside its current validation gaps")
	}
	return nil
}

func (d *delivery) validationWorkspace(path string) (string, error) {
	if path != "." {
		validated, err := validateTaskPaths([]string{path})
		if err != nil || len(validated) != 1 || validated[0] != path {
			return "", fmt.Errorf("invalid validation workspace %q", path)
		}
	}
	absolute := filepath.Join(d.repository.directory, filepath.FromSlash(path))
	canonical, err := canonicalExistingPath(absolute)
	if err != nil || canonical != absolute {
		return "", fmt.Errorf("validation workspace must be an existing non-symlink path: %s", path)
	}
	for _, name := range []string{"MODULE.bazel", "WORKSPACE.bazel", "WORKSPACE"} {
		info, err := os.Lstat(filepath.Join(absolute, name))
		if err == nil && info.Mode().IsRegular() {
			return absolute, nil
		}
		if err == nil {
			return "", fmt.Errorf("validation workspace marker must be a regular non-symlink file")
		}
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("read validation workspace marker: %w", err)
		}
	}
	return "", fmt.Errorf("validation workspace %s has no regular Bazel workspace marker", path)
}

func (d *delivery) requireValidationCandidate(ctx context.Context, receipt preparationReceipt) error {
	branch, head, err := d.repository.branchHead(ctx)
	if err != nil {
		return err
	}
	tree, err := d.repository.tree(ctx, head)
	if err != nil {
		return err
	}
	base := strings.TrimPrefix(receipt.BaseRef, "refs/heads/")
	if d.base != "" && d.base != base {
		return fmt.Errorf("validation base differs from preparation receipt")
	}
	if err := d.validateReceiptContext(ctx, receipt, &inspection{
		Base: base, HeadRef: "refs/heads/" + branch, LocalHeadOID: head, LocalTreeOID: tree,
	}); err != nil {
		return err
	}
	status, err := d.repository.status(ctx)
	if err != nil {
		return err
	}
	if !status.clean() {
		return fmt.Errorf("recorded validation requires a fully clean worktree and index")
	}
	if err := d.repository.requireDefaultIndexFlags(ctx, []string{"."}); err != nil {
		return err
	}
	return d.repository.requireIndexTree(ctx, receipt.PreparedTreeOID)
}

func (d *delivery) readValidationState(ctx context.Context, receiptFile string) (validationState, error) {
	var state validationState
	absolute, err := d.receiptPath(ctx, validationStatePath(receiptFile), true)
	if err != nil {
		return state, err
	}
	contents, err := readStableReceiptFile(absolute)
	if err != nil {
		return state, err
	}
	if err := decodeValidationJSON(contents, validationStateShape(), &state); err != nil {
		return state, err
	}
	if state.Schema != validationStateSchema || !isObjectID(state.HeadOID) || !isObjectID(state.TreeOID) ||
		!validReceiptRevisionNonce(state.ReceiptRevision) || !validSHA256Digest(state.PlanDigest) || !validSHA256Digest(state.EnvironmentDigest) ||
		!validSHA256Digest(state.RepositoryFingerprint) || len(state.Results) > 32 {
		return state, fmt.Errorf("invalid recorded validation state")
	}
	if _, err := time.Parse(time.RFC3339Nano, state.StartedAt); err != nil {
		return state, fmt.Errorf("invalid validation start time")
	}
	if state.CompletedAt != "" {
		if _, err := time.Parse(time.RFC3339Nano, state.CompletedAt); err != nil {
			return state, fmt.Errorf("invalid validation completion time")
		}
	}
	switch state.Status {
	case "running", "failed", "validated", "publication_attempted", "revalidation_required", "verified":
	default:
		return state, fmt.Errorf("unknown validation state %q", state.Status)
	}
	return state, nil
}

func (d *delivery) saveValidationState(ctx context.Context, receiptFile string, state validationState) error {
	return d.writeAtomicIgnoredJSON(ctx, validationStatePath(receiptFile), "validation state", state)
}

func (d *delivery) validateCandidate(ctx context.Context, receiptFile, planFile string) (report *validationReport, returnErr error) {
	if err := d.requireSameTaskOutputs(receiptFile, planFile); err != nil {
		return nil, err
	}
	var err error
	receiptFile, err = d.receiptPath(ctx, receiptFile, true)
	if err != nil {
		return nil, err
	}
	planFile, err = d.receiptPath(ctx, planFile, true)
	if err != nil {
		return nil, err
	}
	for _, output := range []string{receiptFile, receiptFile + ".lock", validationStatePath(receiptFile)} {
		if output == planFile {
			return nil, fmt.Errorf("validation plan must be distinct from receipt, lock, and state outputs")
		}
	}
	tx, err := d.beginReceiptTransaction(ctx, receiptFile, true)
	if err != nil {
		return nil, err
	}
	defer func() { returnErr = errors.Join(returnErr, tx.close()) }()
	receipt, err := tx.read()
	if err != nil {
		return nil, err
	}
	statePath := validationStatePath(receiptFile)
	absolute, err := d.receiptPath(ctx, statePath, false)
	if err != nil {
		return nil, err
	}
	if _, err := os.Lstat(absolute); err == nil {
		previous, err := d.readValidationState(ctx, receiptFile)
		if err != nil {
			return nil, err
		}
		if previous.HeadOID == receipt.PreparedHeadOID && (previous.Status == "publication_attempted" || previous.Status == "verified") {
			return previous.report(receiptFile), fmt.Errorf("publication was already attempted for this candidate; use continue for verification and diagnose before manual recovery")
		}
	}
	plan, planDigest, err := d.readValidationPlan(ctx, planFile)
	if err != nil {
		return nil, err
	}
	if err := d.requireValidationCandidate(ctx, receipt); err != nil {
		return nil, err
	}
	if err := d.checkValidationPlan(plan, receipt); err != nil {
		return nil, err
	}
	for index := range plan.Checks {
		logPath := fmt.Sprintf("%s.check-%02d.log", statePath, index+1)
		if planFile == logPath {
			return nil, fmt.Errorf("validation plan must be distinct from check log outputs")
		}
		if _, err := d.receiptPath(ctx, logPath, false); err != nil {
			return nil, err
		}
	}
	state := validationState{
		Schema: validationStateSchema, Status: "running",
		ReceiptRevision: receipt.RevisionNonce, RepositoryFingerprint: receipt.RepositoryFingerprint,
		HeadOID: receipt.PreparedHeadOID, TreeOID: receipt.PreparedTreeOID,
		PlanFile: planFile, PlanDigest: planDigest, StartedAt: time.Now().UTC().Format(time.RFC3339Nano),
		EnvironmentDigest: validationEnvironmentDigest(),
		Results:           []validationCheckResult{},
	}
	if err := d.saveValidationState(ctx, receiptFile, state); err != nil {
		return nil, err
	}
	defer func() {
		if returnErr != nil && state.Status == "running" {
			state.Status = "failed"
		}
		state.CompletedAt = time.Now().UTC().Format(time.RFC3339Nano)
		returnErr = errors.Join(returnErr, d.saveValidationState(context.WithoutCancel(ctx), receiptFile, state))
		report = state.report(receiptFile)
	}()
	if _, err := d.repository.run(ctx, "diff", "--check", receipt.BaseOID, receipt.PreparedHeadOID, "--"); err != nil {
		return nil, err
	}
	for index, check := range plan.Checks {
		if err := tx.requireCurrent(ctx); err != nil {
			return nil, err
		}
		if err := d.requireValidationCandidate(ctx, receipt); err != nil {
			return nil, err
		}
		workspace, err := d.validationWorkspace(check.Workspace)
		if err != nil {
			return nil, err
		}
		args := []string{"bazel", check.Kind}
		if check.Kind == "lint" {
			args = []string{"bazel", "build", "--config=lint"}
		}
		args = append(args, check.Targets...)
		stepCtx, cancel := context.WithTimeout(ctx, time.Duration(check.TimeoutSeconds)*time.Second)
		result, runErr := d.repository.runner.Run(stepCtx, command{
			Name: "bazel_agent", Args: args, Dir: workspace,
			Env:              []string{"TMPDIR=" + filepath.Dir(absolute), "TMP=" + filepath.Dir(absolute), "TEMP=" + filepath.Dir(absolute), "PWD=" + workspace},
			UnsetEnv:         append(append([]string{}, gitUnsetEnvironment...), "BUILD_WORKSPACE_DIRECTORY", "BUILD_WORKING_DIRECTORY", "OLDPWD", "SHLVL", "_"),
			UnsetEnvPrefixes: gitUnsetEnvironmentPrefixes, OutputLimit: 96 * 1024,
		})
		cancel()
		logPath := fmt.Sprintf("%s.check-%02d.log", statePath, index+1)
		log := "stdout:\n" + redactCredentials(result.Stdout) + "\nstderr:\n" + redactCredentials(result.Stderr)
		if err := d.writeValidationLog(ctx, logPath, []byte(log)); err != nil {
			return nil, err
		}
		outcome := validationCheckResult{
			Index: index, Status: "passed", ExitCode: result.ExitCode,
			Truncated: result.Truncated, LogFile: logPath, LogDigest: digestStrings("repo_delivery validation log v1", log),
		}
		if runErr != nil || result.ExitCode != 0 || result.Truncated {
			outcome.Status = "failed"
		}
		state.Results = append(state.Results, outcome)
		if err := d.saveValidationState(ctx, receiptFile, state); err != nil {
			return nil, err
		}
		if outcome.Status != "passed" {
			return nil, fmt.Errorf("validation check %d failed; inspect %s", index+1, logPath)
		}
		if err := d.requireValidationCandidate(ctx, receipt); err != nil {
			return nil, err
		}
	}
	_, currentDigest, err := d.readValidationPlan(ctx, planFile)
	if err != nil {
		return nil, err
	}
	if currentDigest != planDigest {
		return nil, fmt.Errorf("validation plan changed while checks ran")
	}
	if err := tx.requireCurrent(ctx); err != nil {
		return nil, err
	}
	if validationEnvironmentDigest() != state.EnvironmentDigest {
		return nil, fmt.Errorf("validation environment changed while checks ran")
	}
	state.Status = "validated"
	return nil, nil
}

func (d *delivery) writeValidationLog(ctx context.Context, path string, contents []byte) error {
	absolute, err := d.receiptPath(ctx, path, false)
	if err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(absolute), ".validation-log-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if _, err := file.Write(contents); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), absolute)
}

func (d *delivery) continueCandidate(ctx context.Context, receiptFile string, publish, noPullRequest bool) (report *validationReport, returnErr error) {
	var err error
	receiptFile, err = d.receiptPath(ctx, receiptFile, true)
	if err != nil {
		return nil, err
	}
	tx, err := d.beginReceiptTransaction(ctx, receiptFile, true)
	if err != nil {
		return nil, err
	}
	defer func() { returnErr = errors.Join(returnErr, tx.close()) }()
	receipt, err := tx.read()
	if err != nil {
		return nil, err
	}
	state, err := d.readValidationState(ctx, receiptFile)
	if err != nil {
		return nil, err
	}
	report = state.report(receiptFile)
	if state.HeadOID != receipt.PreparedHeadOID || state.TreeOID != receipt.PreparedTreeOID || state.RepositoryFingerprint != receipt.RepositoryFingerprint {
		return report, fmt.Errorf("recorded validation is for a different prepared candidate; run validate")
	}
	if err := d.requireValidationCandidate(ctx, receipt); err != nil {
		return report, err
	}
	if state.Status == "publication_attempted" || state.Status == "verified" {
		if publish {
			return report, fmt.Errorf("publication was already attempted; continue without --publish only verifies; diagnose failures before manual recovery")
		}
		if d.base == "" {
			d.base = strings.TrimPrefix(receipt.BaseRef, "refs/heads/")
		}
		verified, err := d.verify(ctx, &receipt)
		if err != nil {
			return report, err
		}
		if !verified.Verified {
			return report, fmt.Errorf("publication is not verified")
		}
		state.Status = "verified"
		state.ReceiptRevision = receipt.RevisionNonce
		if err := d.saveValidationState(ctx, receiptFile, state); err != nil {
			return report, err
		}
		return state.report(receiptFile), nil
	}
	if state.Status != "validated" {
		return report, fmt.Errorf("validation state %s does not permit publication", state.Status)
	}
	if state.ReceiptRevision != receipt.RevisionNonce {
		return report, fmt.Errorf("preparation receipt changed after validation; run validate")
	}
	if err := d.requireSameTaskOutputs(receiptFile, state.PlanFile); err != nil {
		return report, err
	}
	plan, digest, err := d.readValidationPlan(ctx, state.PlanFile)
	if err != nil {
		return report, err
	}
	if digest != state.PlanDigest {
		return report, fmt.Errorf("validation plan changed after validation; run validate")
	}
	if err := d.checkValidationPlan(plan, receipt); err != nil {
		return report, err
	}
	if len(state.Results) != len(plan.Checks) || state.CompletedAt == "" {
		return report, fmt.Errorf("validation results are incomplete")
	}
	for index, result := range state.Results {
		if result.Index != index || result.Status != "passed" || result.ExitCode != 0 || result.Truncated || !validSHA256Digest(result.LogDigest) {
			return report, fmt.Errorf("validation results do not establish passing checks")
		}
		wantLog := fmt.Sprintf("%s.check-%02d.log", validationStatePath(receiptFile), index+1)
		if result.LogFile != wantLog {
			return report, fmt.Errorf("validation result log path differs from its check")
		}
		absolute, err := d.receiptPath(ctx, wantLog, true)
		if err != nil {
			return report, err
		}
		contents, err := readStableReceiptFile(absolute)
		if err != nil {
			return report, err
		}
		if digestStrings("repo_delivery validation log v1", string(contents)) != result.LogDigest {
			return report, fmt.Errorf("validation result log changed")
		}
	}
	if validationEnvironmentDigest() != state.EnvironmentDigest {
		return report, fmt.Errorf("validation environment changed; run validate")
	}
	if !publish {
		return report, nil
	}
	if err := tx.requireCurrent(ctx); err != nil {
		return report, err
	}
	state.Status = "publication_attempted"
	if err := d.saveValidationState(ctx, receiptFile, state); err != nil {
		return report, err
	}
	report = state.report(receiptFile)
	publication, publishErr := d.publishWithReceiptTransaction(ctx, publishOptions{
		ReceiptFile: receiptFile, ValidatedHead: state.HeadOID, NoPullRequest: noPullRequest,
	}, tx)
	var revalidation *revalidationRequiredError
	if errors.As(publishErr, &revalidation) {
		state.Status = "revalidation_required"
		state.HeadOID = revalidation.Report.HeadOID
		state.TreeOID = revalidation.Report.TreeOID
		state.ReceiptRevision = revalidation.Report.Receipt.RevisionNonce
		state.Results = []validationCheckResult{}
	} else if publishErr == nil && publication != nil && publication.Verified {
		state.Status = "verified"
	}
	if err := d.saveValidationState(context.WithoutCancel(ctx), receiptFile, state); err != nil {
		return report, errors.Join(publishErr, err)
	}
	report = state.report(receiptFile)
	report.Publication = publication
	if revalidation != nil {
		return report, fmt.Errorf("publication paused before push: revalidate the updated receipt with the selected plan")
	}
	return report, publishErr
}

// Only a digest is retained: inherited environment values can contain secrets.
// Workspace/scratch overrides are fixed by the executor and excluded here so
// the generated launcher's incidental working directory does not invalidate it.
func validationEnvironmentDigest() string {
	unset := append(append([]string{}, gitUnsetEnvironment...),
		"BUILD_WORKSPACE_DIRECTORY", "BUILD_WORKING_DIRECTORY", "TMPDIR", "TMP", "TEMP", "PWD", "OLDPWD", "SHLVL", "_")
	environment := mergeEnvironment(os.Environ(), nil, unset, gitUnsetEnvironmentPrefixes)
	sort.Strings(environment)
	return digestStrings("repo_delivery validation environment v1", environment...)
}
