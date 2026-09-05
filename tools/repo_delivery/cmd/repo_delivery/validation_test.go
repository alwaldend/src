package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const (
	testValidationReceipt = "out/delivery/prepare.json"
	testValidationPlan    = "out/delivery/checks.json"
)

type validationTestRunner struct {
	delegate commandRunner
	checks   []command
	onCheck  func(command)
	fail     bool
	truncate bool
}

func (r *validationTestRunner) Run(ctx context.Context, cmd command) (commandResult, error) {
	if cmd.Name != "bazel_agent" {
		return r.delegate.Run(ctx, cmd)
	}
	r.checks = append(r.checks, cmd)
	if _, ok := ctx.Deadline(); !ok {
		return commandResult{}, errors.New("missing check deadline")
	}
	if r.onCheck != nil {
		r.onCheck(cmd)
	}
	if r.fail {
		return commandResult{ExitCode: 1, Stderr: "test failed"}, errors.New("test failed")
	}
	return commandResult{Stdout: "test passed\n", Truncated: r.truncate}, nil
}

func newValidationFixture(t *testing.T) (integrationDeliveryFixture, *prepareReport, *validationTestRunner, validationPlan) {
	t.Helper()
	fixture := newIntegrationDeliveryFixture(t)
	for _, directory := range []string{"owned", "nested"} {
		if err := os.MkdirAll(filepath.Join(fixture.seed, directory), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	writeTestFile(t, filepath.Join(fixture.seed, "MODULE.bazel"), "module(name = \"fixture\")\n")
	writeTestFile(t, filepath.Join(fixture.seed, "BUILD.bazel"), "# Fixture root package\n")
	writeTestFile(t, filepath.Join(fixture.seed, "owned", "BUILD.bazel"), "# Fixture package\n")
	writeTestFile(t, filepath.Join(fixture.seed, "nested", "MODULE.bazel"), "module(name = \"nested\")\n")
	runTestGit(t, fixture.seed, "add", "MODULE.bazel", "BUILD.bazel", "owned", "nested")
	runTestGit(t, fixture.seed, "commit", "-m", "Add fixture workspaces")
	runTestGit(t, fixture.seed, "push", "origin", "master")
	prepared := fixture.prepare(t)
	runner := &validationTestRunner{delegate: fixture.delivery.repository.runner}
	fixture.delivery.repository.runner = runner
	plan := validationPlan{Schema: validationPlanSchema, Checks: []validationCheck{
		{Workspace: ".", Kind: "test", Targets: []string{"//:repo_quality_test"}, TimeoutSeconds: 30},
		{Workspace: ".", Kind: "lint", Targets: []string{"//owned:all"}, TimeoutSeconds: 30},
		{Workspace: "nested", Kind: "test", Targets: []string{"//:test"}, TimeoutSeconds: 30},
	}, GapDecisions: []validationGapDecision{{Path: "feature.txt", Reason: "Plain text fixture; root quality covers formatting."}}}
	writeValidationPlan(t, fixture, plan)
	return fixture, prepared, runner, plan
}

func writeValidationPlan(t *testing.T, fixture integrationDeliveryFixture, plan validationPlan) {
	t.Helper()
	if err := fixture.delivery.writeAtomicIgnoredJSON(context.Background(), testValidationPlan, "test plan", plan); err != nil {
		t.Fatal(err)
	}
}

func requireValidated(t *testing.T, fixture integrationDeliveryFixture) *validationReport {
	t.Helper()
	report, err := fixture.delivery.validateCandidate(context.Background(), testValidationReceipt, testValidationPlan)
	if err != nil {
		t.Fatalf("validateCandidate() = %v, report %#v", err, report)
	}
	if report.Status != "validated" {
		t.Fatalf("validation = %#v", report)
	}
	return report
}

func TestValidationAndContinuation(t *testing.T) {
	fixture, prepared, runner, _ := newValidationFixture(t)
	report := requireValidated(t, fixture)
	if report.HeadOID != prepared.HeadOID || len(runner.checks) != 3 {
		t.Fatalf("validation = %#v, checks %d", report, len(runner.checks))
	}
	for index, want := range [][]string{{"bazel", "test", "//:repo_quality_test"}, {"bazel", "build", "--config=lint", "//owned:all"}, {"bazel", "test", "//:test"}} {
		if !reflect.DeepEqual(runner.checks[index].Args, want) {
			t.Fatalf("check %d args = %v", index, runner.checks[index].Args)
		}
	}
	if runner.checks[0].Dir != fixture.work || runner.checks[2].Dir != filepath.Join(fixture.work, "nested") {
		t.Fatalf("workspace routing = %#v", runner.checks)
	}
	state, err := fixture.delivery.readValidationState(context.Background(), testValidationReceipt)
	if err != nil {
		t.Fatal(err)
	}
	if state.ReceiptRevision != prepared.Receipt.RevisionNonce || state.TreeOID != prepared.TreeOID {
		t.Fatalf("state = %#v", state)
	}
	if _, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, false, false); err != nil {
		t.Fatal(err)
	}
	if fixture.forge.pull != nil {
		t.Fatal("readiness inspection created a pull request")
	}
	published, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, true, false)
	if err != nil {
		t.Fatalf("continue = %v, report %#v", err, published)
	}
	if published.Status != "verified" || published.Publication == nil || !published.Publication.Verified {
		t.Fatalf("publication = %#v", published)
	}
	if _, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, true, false); err == nil {
		t.Fatal("repeated publication was permitted")
	}
	verified, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, false, false)
	if err != nil || verified.Status != "verified" {
		t.Fatalf("verification = %#v, %v", verified, err)
	}
}

func TestValidationPlanRefusesUnsafeOrIncompleteChecks(t *testing.T) {
	fixture, prepared, _, plan := newValidationFixture(t)
	for _, test := range []struct {
		name   string
		mutate func(*validationPlan)
	}{
		{"run", func(p *validationPlan) { p.Checks[0].Kind = "run" }},
		{"flag", func(p *validationPlan) { p.Checks[0].Targets = []string{"--config=unsafe"} }},
		{"recursive", func(p *validationPlan) { p.Checks[0].Targets = []string{"//..."} }},
		{"root_all", func(p *validationPlan) { p.Checks[0].Targets = []string{"//:all"} }},
		{"timeout", func(p *validationPlan) { p.Checks[0].TimeoutSeconds = 3601 }},
		{"no_quality", func(p *validationPlan) { p.Checks = p.Checks[1:] }},
		{"missing_gap", func(p *validationPlan) { p.GapDecisions = nil }},
		{"extra_gap", func(p *validationPlan) {
			p.GapDecisions = append(p.GapDecisions, validationGapDecision{Path: "other", Reason: "not changed"})
		}},
		{"workspace_escape", func(p *validationPlan) { p.Checks[0].Workspace = "../other" }},
		{"not_workspace", func(p *validationPlan) { p.Checks[0].Workspace = "owned" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			encoded, _ := json.Marshal(plan)
			var candidate validationPlan
			if err := json.Unmarshal(encoded, &candidate); err != nil {
				t.Fatal(err)
			}
			test.mutate(&candidate)
			if err := fixture.delivery.checkValidationPlan(candidate, prepared.Receipt); err == nil {
				t.Fatal("unsafe/incomplete plan accepted")
			}
		})
	}
	changed := prepared.Receipt
	changed.Scope.AggregatePaths = []string{"owned/file.txt"}
	plan.Checks = plan.Checks[:1]
	plan.GapDecisions = []validationGapDecision{}
	if err := fixture.delivery.checkValidationPlan(plan, changed); err == nil || !strings.Contains(err.Error(), "lint affected package") {
		t.Fatalf("missing lint error = %v", err)
	}
}

func TestValidationPlanJSONIsStrict(t *testing.T) {
	for _, value := range []string{
		`{"schema":"repo_delivery/validation_plan/v1","schema":"other","checks":[],"gap_decisions":[]}`,
		`{"schema":"repo_delivery/validation_plan/v1","checks":[],"gap_decisions":[],"command":"sh"}`,
		`{"schema":"repo_delivery/validation_plan/v1","checks":null,"gap_decisions":[]}`,
		`{"schema":"repo_delivery/validation_plan/v1","checks":[],"gap_decisions":[]} {}`,
	} {
		var plan validationPlan
		if err := decodeValidationJSON([]byte(value), validationPlanShape(), &plan); err == nil {
			t.Fatalf("accepted %s", value)
		}
	}
}

func TestValidationRejectsInputOutputAliasesBeforeWriting(t *testing.T) {
	for _, output := range []string{"receipt", "state", "lock", "log"} {
		t.Run(output, func(t *testing.T) {
			fixture, _, runner, plan := newValidationFixture(t)
			path := testValidationReceipt
			switch output {
			case "state":
				path = validationStatePath(testValidationReceipt)
			case "lock":
				path = testValidationReceipt + ".lock"
			case "log":
				path = validationStatePath(testValidationReceipt) + ".check-01.log"
			}
			if output != "receipt" {
				if err := fixture.delivery.writeAtomicIgnoredJSON(context.Background(), path, "aliased plan", plan); err != nil {
					t.Fatal(err)
				}
			}
			absolute := filepath.Join(fixture.work, path)
			before, err := os.ReadFile(absolute)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := fixture.delivery.validateCandidate(context.Background(), testValidationReceipt, absolute); err == nil {
				t.Fatal("aliased plan accepted")
			}
			after, err := os.ReadFile(absolute)
			if err != nil {
				t.Fatal(err)
			}
			if string(before) != string(after) || len(runner.checks) != 0 {
				t.Fatal("refusal wrote input or ran checks")
			}
		})
	}
}

func TestFailedAndTruncatedValidationCannotPublish(t *testing.T) {
	for _, truncate := range []bool{false, true} {
		t.Run(map[bool]string{false: "failure", true: "truncated"}[truncate], func(t *testing.T) {
			fixture, _, runner, _ := newValidationFixture(t)
			runner.fail = !truncate
			runner.truncate = truncate
			report, err := fixture.delivery.validateCandidate(context.Background(), testValidationReceipt, testValidationPlan)
			if err == nil || report.Status != "failed" || len(report.Results) != 1 || len(runner.checks) != 1 {
				t.Fatalf("failed validation = %#v, %v", report, err)
			}
			if _, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, true, false); err == nil {
				t.Fatal("failed validation published")
			}
			if fixture.forge.pull != nil {
				t.Fatal("failed validation created pull request")
			}
		})
	}
}

func TestContinuationRefusesChangedValidationInputs(t *testing.T) {
	for _, change := range []string{"head", "dirty", "plan", "receipt", "environment", "incomplete", "log"} {
		t.Run(change, func(t *testing.T) {
			fixture, _, _, _ := newValidationFixture(t)
			requireValidated(t, fixture)
			switch change {
			case "head":
				runTestGit(t, fixture.work, "commit", "--amend", "-m", "Changed commit", "-m", commitDisclaimer)
			case "dirty":
				writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "changed\n")
			case "plan":
				path := filepath.Join(fixture.work, testValidationPlan)
				contents, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, append(contents, '\n'), 0o600); err != nil {
					t.Fatal(err)
				}
			case "receipt":
				receipt, err := fixture.delivery.readReceipt(context.Background(), testValidationReceipt)
				if err != nil {
					t.Fatal(err)
				}
				receipt.RevisionNonce = strings.Repeat("f", 64)
				if err := fixture.delivery.writeAtomicIgnoredJSON(context.Background(), testValidationReceipt, "receipt", receipt); err != nil {
					t.Fatal(err)
				}
			case "environment":
				t.Setenv("REPO_DELIVERY_TEST_CHANGED_INPUT", "changed")
			case "log":
				if err := os.WriteFile(filepath.Join(fixture.work, validationStatePath(testValidationReceipt)+".check-01.log"), []byte("changed evidence"), 0o600); err != nil {
					t.Fatal(err)
				}
			case "incomplete":
				state, err := fixture.delivery.readValidationState(context.Background(), testValidationReceipt)
				if err != nil {
					t.Fatal(err)
				}
				state.Results = state.Results[:1]
				if err := fixture.delivery.saveValidationState(context.Background(), testValidationReceipt, state); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, true, false); err == nil {
				t.Fatal("changed inputs published")
			}
			if fixture.forge.pull != nil {
				t.Fatal("changed inputs created pull request")
			}
		})
	}
}

func TestValidationRejectsInputsChangedDuringChecks(t *testing.T) {
	for _, change := range []string{"dirty", "plan"} {
		t.Run(change, func(t *testing.T) {
			fixture, _, runner, _ := newValidationFixture(t)
			runner.onCheck = func(command) {
				if change == "dirty" {
					writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "dirty\n")
				} else {
					writeTestFile(t, filepath.Join(fixture.work, testValidationPlan), "{}\n")
				}
			}
			report, err := fixture.delivery.validateCandidate(context.Background(), testValidationReceipt, testValidationPlan)
			if err == nil || report.Status != "failed" {
				t.Fatalf("changed input = %#v, %v", report, err)
			}
		})
	}
}

func TestContinuationRebaseRequiresFreshValidation(t *testing.T) {
	fixture, prepared, runner, plan := newValidationFixture(t)
	requireValidated(t, fixture)
	fixture.advanceBase(t)
	report, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, true, false)
	if err == nil || report.Status != "revalidation_required" || report.HeadOID == prepared.HeadOID {
		t.Fatalf("rebase = %#v, %v", report, err)
	}
	if len(report.Results) != 0 {
		t.Fatal("rebased candidate reports the prior candidate's passing checks")
	}
	if fixture.forge.pull != nil {
		t.Fatal("rebase published before validation")
	}
	if _, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, true, false); err == nil {
		t.Fatal("rebased candidate published without checks")
	}
	writeValidationPlan(t, fixture, plan)
	requireValidated(t, fixture)
	if len(runner.checks) != 6 {
		t.Fatalf("new candidate reused previous checks: %d", len(runner.checks))
	}
	published, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, true, false)
	if err != nil || published.Status != "verified" {
		t.Fatalf("validated rebase = %#v, %v", published, err)
	}
}

func TestContinuationDoesNotRepeatPartialPublication(t *testing.T) {
	fixture, prepared, _, _ := newValidationFixture(t)
	requireValidated(t, fixture)
	if _, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, true, false); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "revised\n")
	writeTestFile(t, filepath.Join(fixture.work, "out/delivery/commit.md"), "Revise fixture\n\nExercise partial publication.\n")
	second, err := fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md", ReceiptFile: testValidationReceipt,
		Paths: []string{"feature.txt"}, RewriteOID: prepared.HeadOID,
	})
	if err != nil {
		t.Fatal(err)
	}
	requireValidated(t, fixture)
	fixture.forge.failUpdates = 1
	report, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, true, false)
	if err == nil || report.Status != "publication_attempted" {
		t.Fatalf("partial publication = %#v, %v", report, err)
	}
	if got := runTestGit(t, fixture.remote, "rev-parse", "feature"); got != second.HeadOID {
		t.Fatalf("remote %s, want %s", got, second.HeadOID)
	}
	if _, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, true, false); err == nil {
		t.Fatal("partial publication repeated automatically")
	}
	if _, err := fixture.delivery.validateCandidate(context.Background(), testValidationReceipt, testValidationPlan); err == nil {
		t.Fatal("validation reset uncertain publication")
	}
	if _, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, false, false); err == nil {
		t.Fatal("partial metadata update verified")
	}
	if _, err := fixture.delivery.publish(context.Background(), publishOptions{ReceiptFile: testValidationReceipt, ValidatedHead: second.HeadOID}); err != nil {
		t.Fatalf("explicit diagnosed recovery: %v", err)
	}
	if _, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, false, false); err != nil {
		t.Fatal(err)
	}
}
