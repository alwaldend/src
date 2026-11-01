package versioning

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestCalendarUsesISOWeekYear(t *testing.T) {
	date := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)
	if got := CalendarForDate(date); got != (Calendar{Year: 2026, Week: 53}) {
		t.Fatalf("CalendarForDate() = %+v, want 2026.53", got)
	}
}

func TestUnpaddedWeekIsSemVerCompatible(t *testing.T) {
	calendar, err := ParseCalendar("2026.8")
	if err != nil || calendar != (Calendar{Year: 2026, Week: 8}) {
		t.Fatalf("ParseCalendar(valid) = %+v, %v", calendar, err)
	}
	if _, err := ParseCalendar("2026.08"); err == nil {
		t.Fatal("ParseCalendar() accepted a zero-padded numeric week")
	}
}

func TestCalendarRejectsNonexistentISOWeek(t *testing.T) {
	if _, err := ParseCalendar("2027.53"); err == nil || !strings.Contains(err.Error(), "has no week 53") {
		t.Fatalf("ParseCalendar() error = %v", err)
	}
	if _, err := ParseVersionTag("v2027.53.0"); err == nil || !strings.Contains(err.Error(), "has no week 53") {
		t.Fatalf("ParseVersionTag() error = %v", err)
	}
	if _, err := ParseCalendar("2026.53"); err != nil {
		t.Fatalf("ParseCalendar(2026.53) error = %v", err)
	}
}

func TestNightlyTagMustMatchISOWeek(t *testing.T) {
	version := NightlyVersion(time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC))
	if got := version.Tag(); got != "v2026.35.0-nightly.20260830" {
		t.Fatalf("NightlyVersion().Tag() = %q", got)
	}
	if _, err := ParseVersionTag("v2026.34.0-nightly.20260830"); err == nil {
		t.Fatal("ParseVersionTag() accepted a mismatched nightly week")
	}
}

func TestParseDateRejectsUnrepresentableCalendarVersions(t *testing.T) {
	for _, value := range []string{"0001-01-01", "0999-12-31"} {
		if _, err := ParseDate(value); err == nil ||
			!strings.Contains(err.Error(), "four-digit calendar version") {
			t.Fatalf("ParseDate(%q) error = %v", value, err)
		}
	}
}

func TestOversizedPatchReturnsIntegrityError(t *testing.T) {
	tag := "v2026.35." + strings.Repeat("9", 1000)
	if _, err := ParseVersionTag(tag); err == nil || !strings.Contains(err.Error(), "release patch") {
		t.Fatalf("ParseVersionTag() error = %v", err)
	}

	runner := &fakeRunner{responses: map[string]string{
		"rev-parse HEAD":                              "abc123",
		"symbolic-ref --quiet HEAD":                   "refs/heads/master",
		"status --porcelain --untracked-files=normal": "",
		"tag --points-at abc123 --list v*":            tag,
	}}
	_, err := (Repository{Git: runner, TrunkBranch: "master"}).Inspect()
	if err == nil || !strings.Contains(err.Error(), "release patch") {
		t.Fatalf("Inspect() error = %v", err)
	}
}

type fakeRunner struct {
	responses      map[string]string
	responseLists  map[string][]string
	errors         map[string]error
	operationPaths map[string]bool
	calls          []string
	inputs         []string
}

func (f *fakeRunner) GitPathExists(name string) (bool, error) {
	return f.operationPaths[name], nil
}

func (f *fakeRunner) Run(input string, args ...string) (string, error) {
	key := strings.Join(args, " ")
	f.calls = append(f.calls, key)
	f.inputs = append(f.inputs, input)
	if err, ok := f.errors[key]; ok {
		return "", err
	}
	if responses := f.responseLists[key]; len(responses) != 0 {
		f.responseLists[key] = responses[1:]
		return responses[0], nil
	}
	if response, ok := f.responses[key]; ok {
		return response, nil
	}
	return "", fmt.Errorf("unexpected Git call %q", key)
}

func releaseRunner(t *testing.T) *fakeRunner {
	t.Helper()
	return &fakeRunner{responses: map[string]string{
		"rev-parse HEAD":                                   "abc123",
		"symbolic-ref --quiet HEAD":                        "refs/heads/releases/2026.35",
		"status --porcelain --untracked-files=normal":      "",
		"tag --points-at abc123 --list v*":                 "",
		"tag --list v2026.35.*":                            "",
		"rev-parse --verify refs/tags/v2026.35.0^{commit}": "base123",
		"rev-list --first-parent abc123":                   "abc123\nparent1\nbase123\nroot123",
	}, errors: map[string]error{
		"show-ref --verify --quiet refs/tags/v2026.35.2": missingRefError(t),
	}}
}

func TestReleaseBaseUsesExplicitTagNamespace(t *testing.T) {
	runner := releaseRunner(t)
	if _, err := (Repository{Git: runner, TrunkBranch: "master"}).Inspect(); err != nil {
		t.Fatal(err)
	}
	for _, call := range runner.calls {
		if call == "rev-parse --verify v2026.35.0^{commit}" {
			t.Fatalf("Inspect() used an ambiguous release base: %v", runner.calls)
		}
	}
}

func TestReleaseBaseMustBeOnFirstParentChain(t *testing.T) {
	runner := releaseRunner(t)
	runner.responses["rev-list --first-parent abc123"] = "abc123\nparent1\nroot123"
	_, err := (Repository{Git: runner, TrunkBranch: "master"}).Inspect()
	if err == nil || !strings.Contains(err.Error(), "first-parent chain") {
		t.Fatalf("Inspect() error = %v", err)
	}
}

func TestInspectRejectsMalformedOwnedVersionTags(t *testing.T) {
	for _, tag := range []string{
		"v2026.08.0",
		"v2026.35.01",
		"v2026.54.0",
		"v2026.35.0-nightly.invalid",
	} {
		runner := &fakeRunner{responses: map[string]string{
			"rev-parse HEAD":                              "abc123",
			"symbolic-ref --quiet HEAD":                   "refs/heads/master",
			"status --porcelain --untracked-files=normal": "",
			"tag --points-at abc123 --list v*":            tag,
		}}
		if _, err := (Repository{Git: runner, TrunkBranch: "master"}).Inspect(); err == nil || !strings.Contains(err.Error(), "validate version tag") {
			t.Fatalf("Inspect(%q) error = %v", tag, err)
		}
	}
}

func TestInspectPropagatesBranchResolutionFailure(t *testing.T) {
	runner := &fakeRunner{
		responses: map[string]string{"rev-parse HEAD": "abc123"},
		errors: map[string]error{
			"symbolic-ref --quiet HEAD": errors.New("branch database failed"),
		},
	}
	if _, err := (Repository{Git: runner, TrunkBranch: "master"}).Inspect(); err == nil || !strings.Contains(err.Error(), "branch database failed") {
		t.Fatalf("Inspect() error = %v", err)
	}
}

func TestInspectRejectsMovingHead(t *testing.T) {
	runner := &fakeRunner{
		responses: map[string]string{
			"symbolic-ref --quiet HEAD":                   "refs/heads/feature",
			"status --porcelain --untracked-files=normal": "",
			"tag --points-at abc123 --list v*":            "",
		},
		responseLists: map[string][]string{
			"rev-parse HEAD": {"abc123", "def456"},
		},
	}
	if _, err := (Repository{Git: runner, TrunkBranch: "master"}).Inspect(); err == nil || !strings.Contains(err.Error(), "HEAD moved") {
		t.Fatalf("Inspect() error = %v", err)
	}
}

func TestReleasePatchCountsFirstParentCommits(t *testing.T) {
	runner := releaseRunner(t)
	state, err := (Repository{Git: runner, TrunkBranch: "master"}).Inspect()
	if err != nil {
		t.Fatal(err)
	}
	if state.Version != "2026.35.2" || state.Channel != "release" {
		t.Fatalf("Inspect() = %+v", state)
	}
}

func TestReleaseTagMustMatchCalculatedPatch(t *testing.T) {
	runner := releaseRunner(t)
	runner.responses["tag --points-at abc123 --list v*"] = "v2026.35.1"
	if _, err := (Repository{Git: runner, TrunkBranch: "master"}).Inspect(); err == nil {
		t.Fatal("Inspect() accepted a stale release tag")
	}
}

func TestBazelStatusRejectsCalculatedReleaseTagPointingElsewhere(t *testing.T) {
	runner := releaseRunner(t)
	delete(runner.errors, "show-ref --verify --quiet refs/tags/v2026.35.2")
	runner.responses["show-ref --verify --quiet refs/tags/v2026.35.2"] = "old123"
	runner.responses["rev-parse --verify refs/tags/v2026.35.2^{commit}"] = "old123"
	_, err := (Repository{Git: runner, TrunkBranch: "master"}).BazelStatus()
	if err == nil || !strings.Contains(err.Error(), "not HEAD") {
		t.Fatalf("BazelStatus() error = %v", err)
	}
}

func TestReleaseRejectsEarlierTagOutsideRewrittenFirstParentHistory(t *testing.T) {
	runner := releaseRunner(t)
	runner.responses["tag --list v2026.35.*"] = "v2026.35.1"
	runner.responses["show-ref --verify --quiet refs/tags/v2026.35.1"] = "old123"
	runner.responses["rev-parse --verify refs/tags/v2026.35.1^{commit}"] = "old123"
	_, err := (Repository{Git: runner, TrunkBranch: "master"}).Inspect()
	if err == nil || !strings.Contains(err.Error(), "not first-parent commit") {
		t.Fatalf("Inspect() error = %v", err)
	}
}

func TestNightlyRefRequiresCleanTrunk(t *testing.T) {
	runner := &fakeRunner{responses: map[string]string{
		"symbolic-ref --quiet HEAD":                   "refs/heads/master",
		"status --porcelain --untracked-files=normal": " M file",
	}}
	_, err := (Repository{Git: runner, TrunkBranch: "master"}).CreateNightly(
		time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC),
		false,
	)
	if err == nil || !strings.Contains(err.Error(), "clean worktree") {
		t.Fatalf("CreateNightly() error = %v", err)
	}
	for _, call := range runner.calls {
		if strings.HasPrefix(call, "update-ref") {
			t.Fatalf("dirty nightly attempted mutation: %v", runner.calls)
		}
	}
}

func TestNightlyRefRequiresConfiguredTrunk(t *testing.T) {
	runner := &fakeRunner{responses: map[string]string{
		"symbolic-ref --quiet HEAD": "refs/heads/feature",
	}}
	_, err := (Repository{Git: runner, TrunkBranch: "master"}).CreateNightly(
		time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC),
		false,
	)
	if err == nil || !strings.Contains(err.Error(), "requires trunk") {
		t.Fatalf("CreateNightly() error = %v", err)
	}
}

func TestExactNightlyTagSelectsNightlyChannelOnTrunk(t *testing.T) {
	runner := &fakeRunner{responses: map[string]string{
		"rev-parse HEAD":                              "abc123",
		"symbolic-ref --quiet HEAD":                   "refs/heads/master",
		"status --porcelain --untracked-files=normal": "",
		"tag --points-at abc123 --list v*":            "v2026.35.0-nightly.20260830",
	}}
	state, err := (Repository{Git: runner, TrunkBranch: "master"}).Inspect()
	if err != nil {
		t.Fatal(err)
	}
	if state.Channel != "nightly" || state.Version != "2026.35.0-nightly.20260830" {
		t.Fatalf("Inspect() = %+v", state)
	}
}

func TestMixedChannelTagsRequireBranchContext(t *testing.T) {
	tags := "v2026.35.0\nv2026.35.0-nightly.20260830"
	detached := &fakeRunner{responses: map[string]string{
		"rev-parse HEAD": "abc123",
		"status --porcelain --untracked-files=normal": "",
		"tag --points-at abc123 --list v*":            tags,
	}, errors: map[string]error{
		"symbolic-ref --quiet HEAD": missingRefError(t),
	}}
	_, err := (Repository{Git: detached, TrunkBranch: "master"}).Inspect()
	if err == nil || !strings.Contains(err.Error(), "both release tag") {
		t.Fatalf("Inspect(detached) error = %v", err)
	}

	trunk := &fakeRunner{responses: map[string]string{
		"rev-parse HEAD":                              "abc123",
		"symbolic-ref --quiet HEAD":                   "refs/heads/master",
		"status --porcelain --untracked-files=normal": "",
		"tag --points-at abc123 --list v*":            tags,
	}}
	state, err := (Repository{Git: trunk, TrunkBranch: "master"}).Inspect()
	if err != nil {
		t.Fatal(err)
	}
	if state.Channel != "nightly" || state.Version != "2026.35.0-nightly.20260830" {
		t.Fatalf("Inspect(trunk) = %+v", state)
	}

	detached = &fakeRunner{responses: map[string]string{
		"rev-parse HEAD": "abc123",
		"status --porcelain --untracked-files=normal":      "",
		"tag --points-at abc123 --list v*":                 tags,
		"rev-parse --verify refs/tags/v2026.35.0^{commit}": "abc123",
		"rev-list --first-parent abc123":                   "abc123\nroot123",
		"tag --list v2026.35.*":                            tags,
		"show-ref --verify --quiet refs/tags/v2026.35.0":   "abc123",
	}, errors: map[string]error{
		"symbolic-ref --quiet HEAD": missingRefError(t),
	}}
	state, err = (Repository{
		Git:              detached,
		TrunkBranch:      "master",
		RequestedChannel: "release",
	}).Inspect()
	if err != nil {
		t.Fatal(err)
	}
	if state.Channel != "release" || state.Version != "2026.35.0" {
		t.Fatalf("Inspect(detached release) = %+v", state)
	}
}

func TestNightlyRefRejectsExistingNightlyAtHead(t *testing.T) {
	for _, dryRun := range []bool{false, true} {
		runner := &fakeRunner{responses: map[string]string{
			"symbolic-ref --quiet HEAD":                   "refs/heads/master",
			"status --porcelain --untracked-files=normal": "",
			"rev-parse HEAD":                              "abc123",
			"tag --points-at abc123 --list v*":            "v2026.35.0-nightly.20260830",
		}}
		_, err := (Repository{Git: runner, TrunkBranch: "master"}).CreateNightly(
			time.Date(2026, time.August, 31, 0, 0, 0, 0, time.UTC),
			dryRun,
		)
		if err == nil || !strings.Contains(err.Error(), "v2026.35.0-nightly.20260830") {
			t.Fatalf("CreateNightly(dryRun=%v) error = %v", dryRun, err)
		}
		for _, call := range runner.calls {
			if call == "update-ref --stdin" {
				t.Fatalf("CreateNightly(dryRun=%v) mutated Git: %v", dryRun, runner.calls)
			}
		}
	}
}

func TestReleaseStartRejectsExistingReleaseAtHead(t *testing.T) {
	for _, dryRun := range []bool{false, true} {
		runner := &fakeRunner{responses: map[string]string{
			"symbolic-ref --quiet HEAD":                   "refs/heads/master",
			"status --porcelain --untracked-files=normal": "",
			"rev-parse HEAD":                              "abc123",
			"tag --points-at abc123 --list v*":            "v2026.35.0",
		}}
		_, err := (Repository{Git: runner, TrunkBranch: "master"}).StartRelease(
			time.Date(2026, time.August, 31, 0, 0, 0, 0, time.UTC),
			false,
			dryRun,
		)
		if err == nil || !strings.Contains(err.Error(), "v2026.35.0") {
			t.Fatalf("StartRelease(dryRun=%v) error = %v", dryRun, err)
		}
		for _, call := range runner.calls {
			if call == "update-ref --stdin" {
				t.Fatalf("StartRelease(dryRun=%v) mutated Git: %v", dryRun, runner.calls)
			}
		}
	}
}

func TestStartReleaseUsesAtomicRefTransaction(t *testing.T) {
	runner := &fakeRunner{responses: map[string]string{
		"symbolic-ref --quiet HEAD":                   "refs/heads/master",
		"status --porcelain --untracked-files=normal": "",
		"tag --points-at abc123 --list v*":            "",
		"rev-parse HEAD":                              "abc123",
		"update-ref --stdin":                          "",
		"switch releases/2026.35":                     "",
	}, errors: map[string]error{
		"show-ref --verify --quiet refs/tags/v2026.35.0":        missingRefError(t),
		"show-ref --verify --quiet refs/heads/releases/2026.35": missingRefError(t),
	}}
	branch, err := (Repository{Git: runner, TrunkBranch: "master"}).StartRelease(
		time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC),
		true,
		false,
	)
	if err != nil {
		t.Fatal(err)
	}
	if branch != "releases/2026.35" {
		t.Fatalf("StartRelease() = %q", branch)
	}
	transaction := ""
	for index, call := range runner.calls {
		if call == "update-ref --stdin" {
			transaction = runner.inputs[index]
			break
		}
	}
	if !strings.Contains(transaction, "create refs/tags/v2026.35.0 abc123") ||
		!strings.Contains(transaction, "create refs/heads/releases/2026.35 abc123") ||
		!strings.Contains(transaction, "verify refs/heads/master abc123") ||
		!strings.Contains(transaction, "prepare\ncommit\n") {
		t.Fatalf("unexpected update-ref transaction:\n%s", transaction)
	}
}

func TestNightlyTransactionBindsValidatedCommit(t *testing.T) {
	const tagRef = "refs/tags/v2026.35.0-nightly.20260830"
	runner := &fakeRunner{responses: map[string]string{
		"symbolic-ref --quiet HEAD":                   "refs/heads/master",
		"status --porcelain --untracked-files=normal": "",
		"rev-parse HEAD":                              "abc123",
		"tag --points-at abc123 --list v*":            "",
		"update-ref --stdin":                          "",
		"rev-parse --verify " + tagRef + "^{commit}":  "abc123",
	}, errors: map[string]error{
		"show-ref --verify --quiet " + tagRef: missingRefError(t),
	}}
	if _, err := (Repository{Git: runner, TrunkBranch: "master"}).CreateNightly(
		time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC),
		false,
	); err != nil {
		t.Fatal(err)
	}
	transaction := runner.inputs[len(runner.inputs)-2]
	if !strings.Contains(transaction, "verify refs/heads/master abc123") ||
		!strings.Contains(transaction, "create "+tagRef+" abc123") {
		t.Fatalf("unexpected nightly transaction:\n%s", transaction)
	}
}

func TestNightlyDryRunDetectsExistingTag(t *testing.T) {
	runner := &fakeRunner{responses: map[string]string{
		"symbolic-ref --quiet HEAD":                                       "refs/heads/master",
		"status --porcelain --untracked-files=normal":                     "",
		"tag --points-at abc123 --list v*":                                "",
		"rev-parse HEAD":                                                  "abc123",
		"show-ref --verify --quiet refs/tags/v2026.35.0-nightly.20260830": "abc123",
	}}
	_, err := (Repository{Git: runner, TrunkBranch: "master"}).CreateNightly(
		time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC),
		true,
	)
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("CreateNightly(dry-run) error = %v", err)
	}
}

func TestReleaseStartDryRunDetectsExistingBranch(t *testing.T) {
	runner := &fakeRunner{responses: map[string]string{
		"symbolic-ref --quiet HEAD":                             "refs/heads/master",
		"status --porcelain --untracked-files=normal":           "",
		"rev-parse HEAD":                                        "abc123",
		"tag --points-at abc123 --list v*":                      "",
		"show-ref --verify --quiet refs/heads/releases/2026.35": "abc123",
	}, errors: map[string]error{
		"show-ref --verify --quiet refs/tags/v2026.35.0": missingRefError(t),
	}}
	_, err := (Repository{Git: runner, TrunkBranch: "master"}).StartRelease(
		time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC),
		true,
		true,
	)
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("StartRelease(dry-run) error = %v", err)
	}
	for _, call := range runner.calls {
		if call == "update-ref --stdin" || strings.HasPrefix(call, "switch ") {
			t.Fatalf("dry-run collision mutated Git: %v", runner.calls)
		}
	}
}

func missingRefError(t *testing.T) error {
	t.Helper()
	err := exec.Command("sh", "-c", "exit 1").Run()
	if err == nil {
		t.Fatal("missing-ref fixture unexpectedly succeeded")
	}
	return err
}

func TestBazelStatusIsDeterministicForDevelopment(t *testing.T) {
	runner := &fakeRunner{responses: map[string]string{
		"rev-parse HEAD":                              "abc123",
		"symbolic-ref --quiet HEAD":                   "refs/heads/feature",
		"status --porcelain --untracked-files=normal": " M file",
		"tag --points-at abc123 --list v*":            "",
	}}
	status, err := (Repository{Git: runner, TrunkBranch: "master"}).BazelStatus()
	if err != nil {
		t.Fatal(err)
	}
	want := "STABLE_VERSION 0.0.0-dev\n" +
		"STABLE_VERSION_CHANNEL development\n" +
		"GIT_COMMIT abc123\n" +
		"GIT_TREE_STATE modified\n"
	if status != want {
		t.Fatalf("BazelStatus() = %q, want %q", status, want)
	}
}

func TestDevelopmentStableStatusIsInvariant(t *testing.T) {
	statusFor := func(commit string, treeState string) string {
		runner := &fakeRunner{responses: map[string]string{
			"rev-parse HEAD":                              commit,
			"symbolic-ref --quiet HEAD":                   "refs/heads/feature",
			"status --porcelain --untracked-files=normal": treeState,
			"tag --points-at " + commit + " --list v*":    "",
		}}
		status, err := (Repository{Git: runner, TrunkBranch: "master"}).BazelStatus()
		if err != nil {
			t.Fatal(err)
		}
		var stable []string
		for _, line := range strings.Split(status, "\n") {
			if strings.HasPrefix(line, "STABLE_") {
				stable = append(stable, line)
			}
		}
		return strings.Join(stable, "\n")
	}
	first := statusFor("abc123", "")
	second := statusFor("def456", " M file")
	if first != second {
		t.Fatalf("development stable status changed:\n%s\n---\n%s", first, second)
	}
}

func TestBazelStatusRefusesDirtyRelease(t *testing.T) {
	runner := releaseRunner(t)
	runner.responses["status --porcelain --untracked-files=normal"] = "?? dirty"
	_, err := (Repository{Git: runner, TrunkBranch: "master"}).BazelStatus()
	if err == nil || !strings.Contains(err.Error(), "refusing to stamp dirty release") {
		t.Fatalf("BazelStatus() error = %v", err)
	}
}
