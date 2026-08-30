package versioning

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type gitFixture struct {
	directory  string
	repository Repository
}

func newGitFixture(t *testing.T, objectFormat string) gitFixture {
	t.Helper()
	directory := filepath.Join(t.TempDir(), "repo")
	args := []string{"init", "--initial-branch=master"}
	if objectFormat != "" {
		args = append(args, "--object-format="+objectFormat)
	}
	args = append(args, directory)
	command := exec.Command("git", args...)
	command.Env = hardenedGitEnvironment(os.Environ())
	if output, err := command.CombinedOutput(); err != nil {
		if objectFormat == "sha256" &&
			(strings.Contains(string(output), "unknown option") ||
				strings.Contains(string(output), "unsupported")) {
			t.Skipf("installed Git lacks SHA-256 repositories: %s", output)
		}
		t.Fatalf("git %s: %s: %v", strings.Join(args, " "), output, err)
	}
	runFixtureGit(t, directory, "config", "user.name", "Versioning Test")
	runFixtureGit(t, directory, "config", "user.email", "versioning@example.invalid")
	runFixtureGit(t, directory, "commit", "--allow-empty", "-m", "base")
	return gitFixture{
		directory: directory,
		repository: Repository{
			Git:         ExecRunner{Directory: directory},
			TrunkBranch: "master",
		},
	}
}

func runFixtureGit(t *testing.T, directory string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", directory}, args...)...)
	command.Env = hardenedGitEnvironment(os.Environ())
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		t.Fatalf("git %s: %s: %v", strings.Join(args, " "), stderr.String(), err)
	}
	return strings.TrimSpace(stdout.String())
}

func fixtureRefExists(t *testing.T, directory string, ref string) bool {
	t.Helper()
	command := exec.Command("git", "-C", directory, "show-ref", "--verify", "--quiet", ref)
	command.Env = hardenedGitEnvironment(os.Environ())
	err := command.Run()
	if err == nil {
		return true
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
		return false
	}
	t.Fatalf("inspect ref %s: %v", ref, err)
	return false
}

func TestRealGitRejectsSecondNightlyOnUnchangedCommit(t *testing.T) {
	fixture := newGitFixture(t, "")
	first := time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC)
	second := time.Date(2026, time.August, 31, 0, 0, 0, 0, time.UTC)
	if _, err := fixture.repository.CreateNightly(first, false); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.repository.CreateNightly(second, false); err == nil ||
		!strings.Contains(err.Error(), "already has nightly tag") {
		t.Fatalf("second CreateNightly() error = %v", err)
	}
	if fixtureRefExists(t, fixture.directory, "refs/tags/v2026.36.0-nightly.20260831") {
		t.Fatal("second nightly tag was created")
	}
	state, err := fixture.repository.Inspect()
	if err != nil {
		t.Fatal(err)
	}
	if state.Version != "2026.35.0-nightly.20260830" {
		t.Fatalf("Inspect().Version = %q", state.Version)
	}
}

func TestRealGitRejectsSecondReleaseOnUnchangedCommit(t *testing.T) {
	fixture := newGitFixture(t, "")
	first := time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC)
	second := time.Date(2026, time.August, 31, 0, 0, 0, 0, time.UTC)
	if _, err := fixture.repository.StartRelease(first, false, false); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.repository.StartRelease(second, false, false); err == nil ||
		!strings.Contains(err.Error(), "already has release tag") {
		t.Fatalf("second StartRelease() error = %v", err)
	}
	if fixtureRefExists(t, fixture.directory, "refs/tags/v2026.36.0") ||
		fixtureRefExists(t, fixture.directory, "refs/heads/releases/2026.36") {
		t.Fatal("second release refs were created")
	}
	if _, err := fixture.repository.Inspect(); err != nil {
		t.Fatalf("Inspect() after refused release = %v", err)
	}
}

func TestRealGitRejectsNightlyAfterReleaseOnSameCommit(t *testing.T) {
	fixture := newGitFixture(t, "")
	date := time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC)
	if _, err := fixture.repository.StartRelease(date, false, false); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.repository.CreateNightly(date, false); err == nil ||
		!strings.Contains(err.Error(), "release tag") {
		t.Fatalf("CreateNightly() after release error = %v", err)
	}
	if fixtureRefExists(
		t,
		fixture.directory,
		"refs/tags/v2026.35.0-nightly.20260830",
	) {
		t.Fatal("nightly tag was created after the release tag")
	}
}

func TestRealGitRejectsEarlierReleaseTagOutsideRewrittenHistory(t *testing.T) {
	fixture := newGitFixture(t, "")
	date := time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC)
	if _, err := fixture.repository.StartRelease(date, true, false); err != nil {
		t.Fatal(err)
	}
	base := runFixtureGit(t, fixture.directory, "rev-parse", "HEAD")
	runFixtureGit(t, fixture.directory, "commit", "--allow-empty", "-m", "patch one")
	if _, err := fixture.repository.TagRelease(false); err != nil {
		t.Fatal(err)
	}
	runFixtureGit(t, fixture.directory, "reset", "--hard", base)
	runFixtureGit(t, fixture.directory, "commit", "--allow-empty", "-m", "replacement one")
	runFixtureGit(t, fixture.directory, "commit", "--allow-empty", "-m", "replacement two")
	_, err := fixture.repository.Inspect()
	if err == nil || !strings.Contains(err.Error(), "not first-parent commit") {
		t.Fatalf("Inspect() after rewritten release error = %v", err)
	}
}

func TestRealGitCreatesNightlyInSHA256Repository(t *testing.T) {
	fixture := newGitFixture(t, "sha256")
	tag, err := fixture.repository.CreateNightly(
		time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC),
		false,
	)
	if err != nil {
		t.Fatal(err)
	}
	want := runFixtureGit(t, fixture.directory, "rev-parse", "HEAD")
	got := runFixtureGit(t, fixture.directory, "rev-parse", "refs/tags/"+tag+"^{commit}")
	if got != want {
		t.Fatalf("nightly tag points to %s, want %s", got, want)
	}
}

func TestExecRunnerIgnoresRepositoryRoutingEnvironment(t *testing.T) {
	selected := newGitFixture(t, "")
	other := newGitFixture(t, "")
	runFixtureGit(t, other.directory, "commit", "--allow-empty", "-m", "other")
	t.Setenv("GIT_DIR", filepath.Join(other.directory, ".git"))
	t.Setenv("GIT_WORK_TREE", other.directory)
	t.Setenv("GIT_INDEX_FILE", filepath.Join(other.directory, ".git", "index"))

	selectedHead := runFixtureGit(t, selected.directory, "rev-parse", "HEAD")
	runnerHead, err := selected.repository.Git.Run("", "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if runnerHead != selectedHead {
		t.Fatalf("ExecRunner HEAD = %s, want selected repository %s", runnerHead, selectedHead)
	}
	if _, err := selected.repository.CreateNightly(
		time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC),
		false,
	); err != nil {
		t.Fatal(err)
	}
	const ref = "refs/tags/v2026.35.0-nightly.20260830"
	if !fixtureRefExists(t, selected.directory, ref) || fixtureRefExists(t, other.directory, ref) {
		t.Fatal("nightly mutation escaped the selected repository")
	}
}

func TestCanonicalBranchNameIgnoresSameNamedTag(t *testing.T) {
	fixture := newGitFixture(t, "")
	runFixtureGit(t, fixture.directory, "tag", "master")
	state, err := fixture.repository.Inspect()
	if err != nil {
		t.Fatal(err)
	}
	if state.Branch != "master" {
		t.Fatalf("Inspect().Branch = %q, want master", state.Branch)
	}
	if _, err := fixture.repository.CreateNightly(
		time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC),
		true,
	); err != nil {
		t.Fatalf("nightly dry-run with same-named tag: %v", err)
	}
}

func TestReleasePatchIgnoresReplaceRefs(t *testing.T) {
	fixture := newGitFixture(t, "")
	if _, err := fixture.repository.StartRelease(
		time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC),
		true,
		false,
	); err != nil {
		t.Fatal(err)
	}
	base := runFixtureGit(t, fixture.directory, "rev-parse", "HEAD")
	runFixtureGit(t, fixture.directory, "commit", "--allow-empty", "-m", "patch one")
	runFixtureGit(t, fixture.directory, "commit", "--allow-empty", "-m", "patch two")
	head := runFixtureGit(t, fixture.directory, "rev-parse", "HEAD")
	tree := runFixtureGit(t, fixture.directory, "rev-parse", "HEAD^{tree}")
	replacement := runFixtureGit(
		t,
		fixture.directory,
		"commit-tree",
		tree,
		"-p",
		base,
		"-m",
		"replacement",
	)
	runFixtureGit(t, fixture.directory, "replace", head, replacement)
	state, err := fixture.repository.Inspect()
	if err != nil {
		t.Fatal(err)
	}
	if state.Version != "2026.35.2" {
		t.Fatalf("Inspect().Version with replace ref = %q, want 2026.35.2", state.Version)
	}
}

func TestMutationsRejectPathCleanMerge(t *testing.T) {
	fixture := newGitFixture(t, "")
	runFixtureGit(t, fixture.directory, "switch", "-c", "topic")
	runFixtureGit(t, fixture.directory, "commit", "--allow-empty", "-m", "topic")
	runFixtureGit(t, fixture.directory, "switch", "master")
	runFixtureGit(t, fixture.directory, "merge", "--no-commit", "--no-ff", "topic")
	if status := runFixtureGit(
		t,
		fixture.directory,
		"status",
		"--porcelain",
		"--untracked-files=normal",
	); status != "" {
		t.Fatalf("merge fixture is not path-clean: %q", status)
	}
	_, err := fixture.repository.CreateNightly(
		time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC),
		true,
	)
	if err == nil || !strings.Contains(err.Error(), "MERGE_HEAD") {
		t.Fatalf("CreateNightly() during merge error = %v", err)
	}
}

func TestDetachedReleaseSelectionAndReleaseLine(t *testing.T) {
	fixture := newGitFixture(t, "")
	date := time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC)
	if _, err := fixture.repository.CreateNightly(date, false); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.repository.StartRelease(date, true, false); err != nil {
		t.Fatal(err)
	}
	runFixtureGit(t, fixture.directory, "switch", "--detach", "v2026.35.0")
	releaseAtTag := fixture.repository
	releaseAtTag.RequestedChannel = "release"
	state, err := releaseAtTag.Inspect()
	if err != nil {
		t.Fatal(err)
	}
	if state.Version != "2026.35.0" || state.Channel != "release" {
		t.Fatalf("detached patch-zero state = %+v", state)
	}

	runFixtureGit(t, fixture.directory, "switch", "releases/2026.35")
	runFixtureGit(t, fixture.directory, "commit", "--allow-empty", "-m", "patch one")
	runFixtureGit(t, fixture.directory, "switch", "--detach", "HEAD")
	releaseCommit := fixture.repository
	releaseCommit.RequestedRelease = "2026.35"
	state, err = releaseCommit.Inspect()
	if err != nil {
		t.Fatal(err)
	}
	if state.Version != "2026.35.1" || state.Channel != "release" {
		t.Fatalf("detached release-line state = %+v", state)
	}
}

func TestDetachedReleaseChannelRejectsTagAtWrongPatchPosition(t *testing.T) {
	fixture := newGitFixture(t, "")
	date := time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC)
	if _, err := fixture.repository.StartRelease(date, true, false); err != nil {
		t.Fatal(err)
	}
	runFixtureGit(t, fixture.directory, "commit", "--allow-empty", "-m", "patch one")
	runFixtureGit(t, fixture.directory, "tag", "v2026.35.9")
	runFixtureGit(t, fixture.directory, "switch", "--detach", "HEAD")
	repository := fixture.repository
	repository.RequestedChannel = "release"
	_, err := repository.Inspect()
	if err == nil || !strings.Contains(err.Error(), "ahead of calculated patch") {
		t.Fatalf("Inspect() with misplaced detached tag error = %v", err)
	}
}

func TestDryRunChecksRefNamespaceCollisions(t *testing.T) {
	fixture := newGitFixture(t, "")
	runFixtureGit(t, fixture.directory, "branch", "releases/2026.35/child")
	_, err := fixture.repository.StartRelease(
		time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC),
		false,
		true,
	)
	if err == nil {
		t.Fatal("release dry-run accepted a branch namespace collision")
	}
	if fixtureRefExists(t, fixture.directory, "refs/tags/v2026.35.0") ||
		fixtureRefExists(t, fixture.directory, "refs/heads/releases/2026.35") {
		t.Fatal("release dry-run created refs")
	}
}
