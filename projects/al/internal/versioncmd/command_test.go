package versioncmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

func TestNoArgumentsPrintsUsage(t *testing.T) {
	var output bytes.Buffer
	if err := execute(nil, time.Now, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "nightly-tag") {
		t.Fatalf("usage = %q", output.String())
	}
}

func TestTopLevelHelpPrintsUsageAndSucceeds(t *testing.T) {
	for _, help := range []string{"--help", "-h"} {
		var output bytes.Buffer
		if err := execute([]string{help}, time.Now, &output); err != nil {
			t.Fatalf("execute(%q) error = %v", help, err)
		}
		if !strings.Contains(output.String(), "nightly-tag") {
			t.Fatalf("execute(%q) usage = %q", help, output.String())
		}
		if !strings.Contains(output.String(), "--channel") ||
			!strings.Contains(output.String(), "--release") {
			t.Fatalf("execute(%q) omitted context options: %q", help, output.String())
		}
	}
}

func TestEverySubcommandHelpPrintsUsageAndSucceeds(t *testing.T) {
	commands := []string{
		"show",
		"nightly-tag",
		"release-tag",
		"release-plan",
		"release-publish",
		"bazel-status",
	}
	for _, command := range commands {
		for _, help := range []string{"--help", "-h"} {
			var output bytes.Buffer
			if err := execute([]string{command, help}, time.Now, &output); err != nil {
				t.Fatalf("execute(%q, %q) error = %v", command, help, err)
			}
			if !strings.Contains(output.String(), command) {
				t.Fatalf("execute(%q, %q) usage = %q", command, help, output.String())
			}
			if !strings.Contains(output.String(), "--channel") ||
				!strings.Contains(output.String(), "--release") {
				t.Fatalf("execute(%q, %q) omitted context options: %q", command, help, output.String())
			}
		}
	}
}

func TestSubcommandHelpAfterOptionsPrintsUsageAndSucceeds(t *testing.T) {
	tests := map[string][]string{
		"show":            {"--format=json"},
		"nightly-tag":     {"--date=2026-08-30", "--dry-run"},
		"release-tag":     {"--dry-run"},
		"release-plan":    {"--plan=out/plan.json"},
		"release-publish": {"--plan=out/plan.json", "--receipt=out/receipt.json"},
	}
	for command, options := range tests {
		args := append([]string{command}, options...)
		args = append(args, "--help")
		var output bytes.Buffer
		if err := execute(args, time.Now, &output); err != nil {
			t.Fatalf("execute(%q) error = %v", args, err)
		}
		if got := output.String(); !strings.Contains(got, command) || !strings.Contains(got, "--repository") {
			t.Fatalf("execute(%q) omitted command or repository flag: %q", args, got)
		}
	}
}

func TestUnknownCommandFails(t *testing.T) {
	if err := execute([]string{"unknown"}, time.Now, &bytes.Buffer{}); err == nil {
		t.Fatal("unknown command succeeded")
	}
}

func TestMutationRejectsInspectionContextOptions(t *testing.T) {
	for _, args := range [][]string{
		{"--channel", "release", "nightly-tag", "--dry-run"},
		{"--release", "2026.35", "release-tag", "--dry-run"},
	} {
		if err := execute(args, time.Now, &bytes.Buffer{}); err == nil ||
			!strings.Contains(err.Error(), "cannot be used with ref mutations") {
			t.Fatalf("execute(%q) error = %v", args, err)
		}
	}
}

func TestInvalidGlobalContextFailsBeforeDispatch(t *testing.T) {
	if err := execute(
		[]string{"--channel", "invalid", "nightly-tag", "--dry-run"},
		time.Now,
		&bytes.Buffer{},
	); err == nil || !strings.Contains(err.Error(), "invalid channel") {
		t.Fatalf("execute(invalid channel) error = %v", err)
	}
}

// TestReleaseCreationAndWorkspaceStatus exercises guarded local ref creation,
// explicit checkout, and installed-AL stamping, retaining verified ref evidence.
func TestReleaseCreationAndWorkspaceStatus(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	binary, err := runfiles.Rlocation(os.Getenv("AL_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	statusScript, err := runfiles.Rlocation(os.Getenv("WORKSPACE_STATUS_SCRIPT"))
	if err != nil {
		t.Fatal(err)
	}
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	installed := t.TempDir()
	for name, source := range map[string]string{"al": binary, "git": git} {
		if err := os.Symlink(source, filepath.Join(installed, name)); err != nil {
			t.Fatal(err)
		}
	}
	repository := t.TempDir()
	runTestGit(t, repository, "init", "--initial-branch=master")
	runTestGit(t, repository, "config", "user.name", "Versioning Test")
	runTestGit(t, repository, "config", "user.email", "versioning@example.invalid")
	runTestGit(t, repository, "config", "commit.gpgsign", "false")
	payload := filepath.Join(repository, "payload.txt")
	if err := os.WriteFile(payload, []byte("release payload\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runTestGit(t, repository, "add", "payload.txt")
	runTestGit(t, repository, "commit", "-m", "Initialize release fixture")
	commit := runTestGit(t, repository, "rev-parse", "HEAD")
	environment := []string{
		"PATH=" + installed,
		"GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null",
		"BUILD_WORKSPACE_DIRECTORY=" + t.TempDir(),
	}
	create := func(dryRun bool) ([]byte, error) {
		args := []string{"release", "create", "--repository", repository, "--date", "2026-08-30"}
		if dryRun {
			args = append(args, "--dry-run")
		}
		command := exec.CommandContext(ctx, binary, args...)
		command.Env = environment
		output, err := command.CombinedOutput()
		if err != nil {
			return output, fmt.Errorf("create release: %w", err)
		}
		return output, nil
	}
	refsBefore := runTestGit(t, repository, "show-ref")
	refusals := map[string]string{}
	for _, scenario := range []string{"dirty", "branch-exists", "tag-exists"} {
		switch scenario {
		case "dirty":
			if err := os.WriteFile(payload, []byte("uncommitted payload\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		case "branch-exists":
			runTestGit(t, repository, "branch", "releases/2026.35")
		case "tag-exists":
			runTestGit(t, repository, "tag", "v2026.35.0")
		}
		before := runTestGit(t, repository, "show-ref")
		for _, dryRun := range []bool{true, false} {
			output, err := create(dryRun)
			if err == nil {
				t.Fatalf("release create accepted %s with dry-run=%v: %s", scenario, dryRun, output)
			}
			refusals[fmt.Sprintf("%s/dry-run=%v", scenario, dryRun)] = string(output)
			if refs := runTestGit(t, repository, "show-ref"); refs != before {
				t.Fatalf("refusal changed refs: %s\n%s", before, refs)
			}
			if head := runTestGit(t, repository, "symbolic-ref", "HEAD"); head != "refs/heads/master" {
				t.Fatalf("refusal changed HEAD: %s", head)
			}
		}
		switch scenario {
		case "dirty":
			if err := os.WriteFile(payload, []byte("release payload\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		case "branch-exists":
			runTestGit(t, repository, "branch", "-D", "releases/2026.35")
		case "tag-exists":
			runTestGit(t, repository, "tag", "-d", "v2026.35.0")
		}
	}
	dryRun, err := create(true)
	if err != nil || string(dryRun) != "releases/2026.35\n" {
		t.Fatalf("release dry-run: %v\n%s", err, dryRun)
	}
	if refs := runTestGit(t, repository, "show-ref"); refs != refsBefore {
		t.Fatalf("dry-run changed refs: %s", refs)
	}
	created, err := create(false)
	if err != nil || string(created) != "releases/2026.35\n" {
		t.Fatalf("create release: %v\n%s", err, created)
	}
	for _, ref := range []string{"refs/heads/releases/2026.35", "refs/tags/v2026.35.0"} {
		if value := runTestGit(t, repository, "rev-parse", ref); value != commit {
			t.Fatalf("created ref %s points to %s, want %s", ref, value, commit)
		}
	}
	createdHEAD := runTestGit(t, repository, "symbolic-ref", "HEAD")
	if createdHEAD != "refs/heads/master" {
		t.Fatalf("release creation changed HEAD: %s", createdHEAD)
	}
	refsCreated := runTestGit(t, repository, "show-ref")
	runTestGit(t, repository, "checkout", "releases/2026.35")
	runStatus := func() ([]byte, error) {
		command := exec.CommandContext(ctx, "/bin/sh", statusScript)
		command.Dir = repository
		command.Env = environment
		output, err := command.CombinedOutput()
		if err != nil {
			return output, fmt.Errorf("run workspace status adapter: %w", err)
		}
		return output, nil
	}
	status, err := runStatus()
	if err != nil {
		t.Fatalf("stamp clean release: %v\n%s", err, status)
	}
	want := "STABLE_VERSION 2026.35.0\nSTABLE_VERSION_CHANNEL release\n" +
		"STABLE_GIT_COMMIT " + commit + "\nSTABLE_GIT_TREE_STATE clean\n"
	if string(status) != want {
		t.Fatalf("release status = %q, want %q", status, want)
	}
	if err := os.WriteFile(payload, []byte("modified payload\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dirty, err := runStatus()
	if err == nil || !strings.Contains(string(dirty), "dirty") {
		t.Fatalf("dirty release status = %q, error = %v", dirty, err)
	}
	runTestGit(t, repository, "add", "payload.txt")
	runTestGit(t, repository, "commit", "-m", "Update release fixture")
	tag := exec.CommandContext(ctx, binary, "version", "release-tag", "--repository", repository)
	tag.Env = environment
	patchTag, err := tag.CombinedOutput()
	if err != nil || string(patchTag) != "v2026.35.1\n" {
		t.Fatalf("tag patch after explicit checkout: %v\n%s", err, patchTag)
	}
	if patch := runTestGit(t, repository, "rev-parse", "refs/tags/v2026.35.1"); patch != runTestGit(t, repository, "rev-parse", "HEAD") {
		t.Fatalf("patch tag points to %s instead of checked-out release HEAD", patch)
	}
	if err := os.Remove(filepath.Join(installed, "al")); err != nil {
		t.Fatal(err)
	}
	missing, err := runStatus()
	if err == nil || !strings.Contains(string(missing), "bazel run //tools/al:install") {
		t.Fatalf("missing AL status = %q, error = %v", missing, err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		creation, err := json.MarshalIndent(map[string]any{
			"commit": commit, "head_after_create": createdHEAD,
			"refs_before": refsBefore, "refs_created": refsCreated,
			"dry_run": string(dryRun), "refusals": refusals, "patch_tag": string(patchTag),
		}, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		for name, body := range map[string][]byte{
			"release-status.txt": status, "dirty-status-error.txt": dirty, "missing-al-error.txt": missing,
			"release-creation.json": append(creation, '\n'),
		} {
			if err := os.WriteFile(filepath.Join(output, name), body, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestReleasePlanWritesReviewedJSON(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	runTestGit(t, directory, "init", "--initial-branch=master")
	runTestGit(t, directory, "config", "user.name", "Versioning Test")
	runTestGit(t, directory, "config", "user.email", "versioning@example.invalid")
	runTestGit(t, directory, "commit", "--allow-empty", "-m", "base")
	runTestGit(t, directory, "tag", "v2026.35.0-nightly.20260830")
	var output bytes.Buffer
	if err := execute([]string{
		"--repository", directory,
		"release-plan",
		"--plan", filepath.Join(directory, "out", "plan.json"),
	}, func() time.Time {
		return time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC)
	}, &output); err != nil {
		t.Fatalf("release-plan error = %v", err)
	}
	if !strings.Contains(output.String(), `"schema": "versioning/release_ref_plan/v1"`) {
		t.Fatalf("release-plan output missing schema: %q", output.String())
	}
}

func TestReleasePlanRequiresPlanPath(t *testing.T) {
	var output bytes.Buffer
	if err := execute([]string{"--repository", t.TempDir(), "release-plan"}, time.Now, &output); err == nil ||
		!strings.Contains(err.Error(), "--plan") {
		t.Fatalf("release-plan() error = %v, want --plan refusal", err)
	}
}

func TestReleasePublishRequiresPlanPath(t *testing.T) {
	var output bytes.Buffer
	if err := execute([]string{"--repository", t.TempDir(), "release-publish"}, time.Now, &output); err == nil ||
		!strings.Contains(err.Error(), "--plan") {
		t.Fatalf("release-publish() error = %v, want --plan refusal", err)
	}
}

func TestReleasePublishRejectsMismatchedAtomicFlag(t *testing.T) {
	directory := t.TempDir()
	plan := `{
		"schema": "versioning/release_ref_plan/v1",
		"kind": "ReleaseRefPlan",
		"version": "2026.35.0",
		"channel": "release",
		"commit": "` + strings.Repeat("a", 40) + `",
		"treeState": "clean",
		"refs": [
			{"ref": "refs/heads/releases/2026.35", "expected": "` + strings.Repeat("a", 40) + `", "operation": "create"},
			{"ref": "refs/tags/v2026.35.0", "expected": "` + strings.Repeat("a", 40) + `", "operation": "create"}
		],
		"atomic": true,
		"lease": "lease"
	}`
	planPath := filepath.Join(directory, "plan.json")
	if err := os.WriteFile(planPath, []byte(plan), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := execute([]string{
		"--repository", directory,
		"release-publish",
		"--plan", planPath,
		"--receipt", filepath.Join(directory, "receipt.json"),
		"--atomic=false",
	}, time.Now, &output); err == nil || !strings.Contains(err.Error(), "--atomic must match") {
		t.Fatalf("release-publish() error = %v, want atomic mismatch refusal", err)
	}
}

func runTestGit(t *testing.T, directory string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", directory}, args...)...)
	var stdout bytes.Buffer
	command.Stdout = &stdout
	if err := command.Run(); err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(stdout.String())
}

// execute drives the native command tree with a fixture clock and output stream.
func execute(args []string, now func() time.Time, stdout io.Writer) error {
	cmd := newCommand(now)
	cmd.SetArgs(args)
	cmd.SetOut(stdout)
	cmd.SetErr(io.Discard)
	cmd.SilenceUsage, cmd.SilenceErrors = true, true
	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("execute version fixture: %w", err)
	}
	return nil
}
