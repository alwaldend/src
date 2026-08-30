package versioning

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Runner interface {
	Run(input string, args ...string) (string, error)
	GitPathExists(name string) (bool, error)
}

type ExecRunner struct {
	Directory string
}

func (r ExecRunner) Run(input string, args ...string) (string, error) {
	command := exec.Command("git", append([]string{"-C", r.Directory}, args...)...)
	command.Env = hardenedGitEnvironment(os.Environ())
	command.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return "", fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), message, err)
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (r ExecRunner) GitPathExists(name string) (bool, error) {
	path, err := r.Run(
		"",
		"rev-parse",
		"--path-format=absolute",
		"--git-path",
		name,
	)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("inspect Git operation marker %s: %w", name, err)
}

var gitEnvironmentToUnset = map[string]bool{
	"GIT_ALTERNATE_OBJECT_DIRECTORIES": true,
	"GIT_COMMON_DIR":                   true,
	"GIT_CONFIG":                       true,
	"GIT_CONFIG_COUNT":                 true,
	"GIT_CONFIG_GLOBAL":                true,
	"GIT_CONFIG_NOSYSTEM":              true,
	"GIT_CONFIG_PARAMETERS":            true,
	"GIT_CONFIG_SYSTEM":                true,
	"GIT_DIR":                          true,
	"GIT_GRAFT_FILE":                   true,
	"GIT_INDEX_FILE":                   true,
	"GIT_NAMESPACE":                    true,
	"GIT_NO_LAZY_FETCH":                true,
	"GIT_NO_REPLACE_OBJECTS":           true,
	"GIT_OBJECT_DIRECTORY":             true,
	"GIT_OPTIONAL_LOCKS":               true,
	"GIT_QUARANTINE_PATH":              true,
	"GIT_REPLACE_REF_BASE":             true,
	"GIT_SHALLOW_FILE":                 true,
	"GIT_WORK_TREE":                    true,
}

func hardenedGitEnvironment(environment []string) []string {
	result := make([]string, 0, len(environment)+4)
	for _, value := range environment {
		name, _, found := strings.Cut(value, "=")
		if !found || gitEnvironmentToUnset[name] ||
			strings.HasPrefix(name, "GIT_CONFIG_KEY_") ||
			strings.HasPrefix(name, "GIT_CONFIG_VALUE_") {
			continue
		}
		result = append(result, value)
	}
	return append(
		result,
		"GIT_GRAFT_FILE="+os.DevNull,
		"GIT_NO_LAZY_FETCH=1",
		"GIT_NO_REPLACE_OBJECTS=1",
		"GIT_OPTIONAL_LOCKS=0",
	)
}

type Repository struct {
	Git              Runner
	TrunkBranch      string
	RequestedChannel string
	RequestedRelease string
}

type State struct {
	Version   string `json:"version"`
	Channel   string `json:"channel"`
	Branch    string `json:"branch"`
	Commit    string `json:"commit"`
	TreeState string `json:"treeState"`
}

func (r Repository) Inspect() (State, error) {
	if r.Git == nil {
		return State{}, fmt.Errorf("Git runner is required")
	}
	if r.TrunkBranch == "" {
		return State{}, fmt.Errorf("trunk branch is required")
	}
	commit, err := r.Git.Run("", "rev-parse", "HEAD")
	if err != nil {
		return State{}, err
	}
	branch, _, err := r.currentBranch()
	if err != nil {
		return State{}, err
	}
	status, err := r.Git.Run("", "status", "--porcelain", "--untracked-files=normal")
	if err != nil {
		return State{}, err
	}
	treeState := "clean"
	if status != "" {
		treeState = "modified"
	}
	state := State{
		Version:   DevelopmentVersion,
		Channel:   "development",
		Branch:    branch,
		Commit:    commit,
		TreeState: treeState,
	}

	stable, nightly, err := r.versionTagsAt(commit)
	if err != nil {
		return State{}, err
	}
	requestedChannel, err := normalizeRequestedChannel(r.RequestedChannel)
	if err != nil {
		return State{}, err
	}
	var requestedRelease *Calendar
	if r.RequestedRelease != "" {
		calendar, err := ParseCalendar(r.RequestedRelease)
		if err != nil {
			return State{}, fmt.Errorf("parse requested release: %w", err)
		}
		if requestedChannel == "nightly" {
			return State{}, fmt.Errorf("requested release %s conflicts with nightly channel", calendar.String())
		}
		requestedRelease = &calendar
	}

	const releasePrefix = "releases/"
	if strings.HasPrefix(branch, releasePrefix) {
		if requestedChannel == "nightly" {
			return State{}, fmt.Errorf(
				"requested nightly channel conflicts with release branch %q",
				branch,
			)
		}
		calendar, err := ParseCalendar(strings.TrimPrefix(branch, releasePrefix))
		if err != nil {
			return State{}, fmt.Errorf("parse release branch %q: %w", branch, err)
		}
		if requestedRelease != nil && *requestedRelease != calendar {
			return State{}, fmt.Errorf(
				"release branch %q conflicts with requested release %s",
				branch,
				requestedRelease.String(),
			)
		}
		return r.inspectRelease(state, calendar, commit, stable)
	}
	if requestedRelease != nil {
		if branch != "HEAD" {
			return State{}, fmt.Errorf(
				"requested release %s is only valid on detached HEAD or its matching release branch; current branch is %q",
				requestedRelease.String(),
				branch,
			)
		}
		return r.inspectRelease(state, *requestedRelease, commit, stable)
	}

	if requestedChannel == "release" {
		if stable == nil {
			return State{}, fmt.Errorf("requested release channel but HEAD has no release tag")
		}
		return r.inspectRelease(state, stable.Calendar, commit, stable)
	}
	if requestedChannel == "nightly" {
		if nightly == nil {
			return State{}, fmt.Errorf("requested nightly channel but HEAD has no nightly tag")
		}
		state.Version = nightly.String()
		state.Channel = "nightly"
		return r.finishInspect(state, commit)
	}
	if stable != nil && nightly != nil && branch != r.TrunkBranch {
		return State{}, fmt.Errorf(
			"HEAD has both release tag %s and nightly tag %s without trunk or release branch context",
			stable.Tag(),
			nightly.Tag(),
		)
	}
	if nightly != nil && branch == r.TrunkBranch {
		state.Version = nightly.String()
		state.Channel = "nightly"
		return r.finishInspect(state, commit)
	}
	if stable != nil {
		state.Version = stable.String()
		state.Channel = "release"
		return r.finishInspect(state, commit)
	}
	if nightly != nil {
		if branch != r.TrunkBranch && branch != "HEAD" {
			return State{}, fmt.Errorf("nightly tag %s is on %q, not trunk %q", nightly.Tag(), branch, r.TrunkBranch)
		}
		state.Version = nightly.String()
		state.Channel = "nightly"
	}
	return r.finishInspect(state, commit)
}

func (r Repository) inspectRelease(
	state State,
	calendar Calendar,
	commit string,
	stable *Version,
) (State, error) {
	base := ReleaseVersion(calendar, 0)
	baseRef := "refs/tags/" + base.Tag()
	baseCommit, err := r.Git.Run("", "rev-parse", "--verify", baseRef+"^{commit}")
	if err != nil {
		return State{}, fmt.Errorf("resolve release base %s: %w", base.Tag(), err)
	}
	chainOutput, err := r.Git.Run("", "rev-list", "--first-parent", commit)
	if err != nil {
		return State{}, err
	}
	chain := strings.Fields(chainOutput)
	patch := -1
	for index, candidate := range chain {
		if candidate == baseCommit {
			patch = index
			break
		}
	}
	if patch < 0 {
		return State{}, fmt.Errorf(
			"commit %s does not contain %s on its first-parent chain",
			commit,
			base.Tag(),
		)
	}
	lineTags, err := r.Git.Run("", "tag", "--list", "v"+calendar.String()+".*")
	if err != nil {
		return State{}, err
	}
	for _, tagName := range strings.Fields(lineTags) {
		version, err := ParseVersionTag(tagName)
		if err != nil {
			return State{}, fmt.Errorf("validate release-line tag %q: %w", tagName, err)
		}
		if version.Nightly != "" || version.Calendar != calendar {
			continue
		}
		if version.Patch == 0 {
			continue
		}
		if version.Patch > patch {
			return State{}, fmt.Errorf(
				"release tag %s is ahead of calculated patch %d",
				tagName,
				patch,
			)
		}
		tagCommit, _, err := r.refCommitIfPresent("refs/tags/" + tagName)
		if err != nil {
			return State{}, err
		}
		expectedCommit := chain[patch-version.Patch]
		if tagCommit != expectedCommit {
			return State{}, fmt.Errorf(
				"release tag %s points to %s, not first-parent commit %s",
				tagName,
				tagCommit,
				expectedCommit,
			)
		}
	}
	expected := ReleaseVersion(calendar, patch)
	expectedRef := "refs/tags/" + expected.Tag()
	tagCommit, exists, err := r.refCommitIfPresent(expectedRef)
	if err != nil {
		return State{}, err
	}
	if exists && tagCommit != commit {
		return State{}, fmt.Errorf(
			"release tag %s points to %s, not HEAD %s",
			expected.Tag(),
			tagCommit,
			commit,
		)
	}
	if stable != nil && stable.String() != expected.String() {
		return State{}, fmt.Errorf(
			"exact release tag %s disagrees with calculated version %s",
			stable.Tag(),
			expected.String(),
		)
	}
	state.Version = expected.String()
	state.Channel = "release"
	return r.finishInspect(state, commit)
}

func (r Repository) currentBranch() (string, string, error) {
	ref, err := r.Git.Run("", "symbolic-ref", "--quiet", "HEAD")
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
			return "HEAD", "", nil
		}
		return "", "", fmt.Errorf("resolve current branch: %w", err)
	}
	const headPrefix = "refs/heads/"
	if !strings.HasPrefix(ref, headPrefix) || len(ref) == len(headPrefix) {
		return "", "", fmt.Errorf("current symbolic ref %q is not a branch", ref)
	}
	return strings.TrimPrefix(ref, headPrefix), ref, nil
}

func (r Repository) finishInspect(state State, expectedCommit string) (State, error) {
	actualCommit, err := r.Git.Run("", "rev-parse", "HEAD")
	if err != nil {
		return State{}, err
	}
	if actualCommit != expectedCommit {
		return State{}, fmt.Errorf(
			"HEAD moved during inspection from %s to %s; retry",
			expectedCommit,
			actualCommit,
		)
	}
	return state, nil
}

func (r Repository) CreateNightly(date time.Time, dryRun bool) (string, error) {
	if err := validateVersionDate(date); err != nil {
		return "", err
	}
	sourceRef, head, err := r.requireCleanBranch(r.TrunkBranch)
	if err != nil {
		return "", err
	}
	tag := NightlyVersion(date).Tag()
	stable, nightly, err := r.versionTagsAt(head)
	if err != nil {
		return "", err
	}
	if stable != nil {
		return "", fmt.Errorf(
			"HEAD already has release tag %s; refusing to add a nightly afterward",
			stable.Tag(),
		)
	}
	if nightly != nil {
		return "", fmt.Errorf(
			"HEAD already has nightly tag %s; refusing to create %s",
			nightly.Tag(),
			tag,
		)
	}
	tagRef := "refs/tags/" + tag
	if err := r.requireAbsentRef(tagRef); err != nil {
		return "", err
	}
	if err := r.updateRefs(
		[]refUpdate{{Ref: sourceRef, Value: head}},
		[]refUpdate{{Ref: tagRef, Value: head}},
		dryRun,
	); err != nil {
		return "", fmt.Errorf("create nightly tag %s: %w", tag, err)
	}
	if err := r.verifyCreatedRef(tagRef, head, dryRun); err != nil {
		return "", err
	}
	return tag, nil
}

func (r Repository) StartRelease(date time.Time, switchBranch bool, dryRun bool) (string, error) {
	if err := validateVersionDate(date); err != nil {
		return "", err
	}
	sourceRef, head, err := r.requireCleanBranch(r.TrunkBranch)
	if err != nil {
		return "", err
	}
	calendar := CalendarForDate(date)
	branch := "releases/" + calendar.String()
	tag := ReleaseVersion(calendar, 0).Tag()
	stable, _, err := r.versionTagsAt(head)
	if err != nil {
		return "", err
	}
	if stable != nil {
		return "", fmt.Errorf(
			"HEAD already has release tag %s; refusing to create %s",
			stable.Tag(),
			tag,
		)
	}
	if err := r.requireAbsentRef("refs/tags/" + tag); err != nil {
		return "", err
	}
	if err := r.requireAbsentRef("refs/heads/" + branch); err != nil {
		return "", err
	}
	tagRef := "refs/tags/" + tag
	branchRef := "refs/heads/" + branch
	if err := r.updateRefs(
		[]refUpdate{{Ref: sourceRef, Value: head}},
		[]refUpdate{
			{Ref: tagRef, Value: head},
			{Ref: branchRef, Value: head},
		},
		dryRun,
	); err != nil {
		return "", fmt.Errorf("create release refs: %w", err)
	}
	if !dryRun {
		if switchBranch {
			if _, err := r.Git.Run("", "switch", branch); err != nil {
				return "", fmt.Errorf(
					"created release branch %q and tag %s at %s, but switching failed; refs were kept: %w",
					branch,
					tag,
					head,
					err,
				)
			}
		}
	}
	return branch, nil
}

func (r Repository) TagRelease(dryRun bool) (string, error) {
	state, err := r.Inspect()
	if err != nil {
		return "", err
	}
	if state.TreeState != "clean" {
		return "", fmt.Errorf("release tagging requires a clean worktree")
	}
	if state.Channel != "release" || !strings.HasPrefix(state.Branch, "releases/") {
		return "", fmt.Errorf("release tagging requires a releases/YYYY.W branch")
	}
	if err := r.requireNoGitOperation(); err != nil {
		return "", err
	}
	tag := "v" + state.Version
	tagRef := "refs/tags/" + tag
	if err := r.requireAbsentRef(tagRef); err != nil {
		return "", err
	}
	if err := r.updateRefs(
		[]refUpdate{{Ref: "refs/heads/" + state.Branch, Value: state.Commit}},
		[]refUpdate{{Ref: tagRef, Value: state.Commit}},
		dryRun,
	); err != nil {
		return "", fmt.Errorf("create release tag %s: %w", tag, err)
	}
	if err := r.verifyCreatedRef(tagRef, state.Commit, dryRun); err != nil {
		return "", err
	}
	return tag, nil
}

func (r Repository) BazelStatus() (string, error) {
	state, err := r.Inspect()
	if err != nil {
		return "", err
	}
	if state.TreeState != "clean" && state.Channel != "development" {
		return "", fmt.Errorf("refusing to stamp dirty %s version %s", state.Channel, state.Version)
	}
	commitKey := "STABLE_GIT_COMMIT"
	treeStateKey := "STABLE_GIT_TREE_STATE"
	if state.Channel == "development" {
		// Bazel does not invalidate stamped actions when volatile status changes.
		// Keep ordinary development's stable-status inputs static across commits.
		commitKey = "GIT_COMMIT"
		treeStateKey = "GIT_TREE_STATE"
	}
	return fmt.Sprintf(
		"STABLE_VERSION %s\nSTABLE_VERSION_CHANNEL %s\n%s %s\n%s %s\n",
		state.Version,
		state.Channel,
		commitKey,
		state.Commit,
		treeStateKey,
		state.TreeState,
	), nil
}

func (r Repository) requireCleanBranch(expected string) (string, string, error) {
	branch, branchRef, err := r.currentBranch()
	if err != nil {
		return "", "", err
	}
	if branch != expected {
		return "", "", fmt.Errorf(
			"operation requires trunk %q; current branch is %q",
			expected,
			branch,
		)
	}
	status, err := r.Git.Run("", "status", "--porcelain", "--untracked-files=normal")
	if err != nil {
		return "", "", err
	}
	if status != "" {
		return "", "", fmt.Errorf("operation requires a clean worktree")
	}
	if err := r.requireNoGitOperation(); err != nil {
		return "", "", err
	}
	head, err := r.Git.Run("", "rev-parse", "HEAD")
	if err != nil {
		return "", "", err
	}
	return branchRef, head, nil
}

func (r Repository) requireNoGitOperation() error {
	markers := []string{
		"AM_HEAD",
		"BISECT_START",
		"CHERRY_PICK_HEAD",
		"MERGE_HEAD",
		"REVERT_HEAD",
		"rebase-apply",
		"rebase-merge",
		"sequencer",
	}
	for _, marker := range markers {
		exists, err := r.Git.GitPathExists(marker)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf(
				"operation requires no Git operation in progress; found %s",
				marker,
			)
		}
	}
	return nil
}

type refUpdate struct {
	Ref   string
	Value string
}

func (r Repository) updateRefs(
	verifies []refUpdate,
	creates []refUpdate,
	dryRun bool,
) error {
	var transaction strings.Builder
	transaction.WriteString("start\n")
	for _, update := range verifies {
		fmt.Fprintf(&transaction, "verify %s %s\n", update.Ref, update.Value)
	}
	for _, update := range creates {
		fmt.Fprintf(&transaction, "create %s %s\n", update.Ref, update.Value)
	}
	transaction.WriteString("prepare\n")
	if dryRun {
		transaction.WriteString("abort\n")
	} else {
		transaction.WriteString("commit\n")
	}
	if _, err := r.Git.Run(transaction.String(), "update-ref", "--stdin"); err != nil {
		return err
	}
	return nil
}

func (r Repository) verifyCreatedRef(ref string, expected string, dryRun bool) error {
	if dryRun {
		return nil
	}
	actual, err := r.Git.Run("", "rev-parse", "--verify", ref+"^{commit}")
	if err != nil {
		return fmt.Errorf("verify created Git ref %s: %w", ref, err)
	}
	if actual != expected {
		return fmt.Errorf("created Git ref %s points to %s, want %s", ref, actual, expected)
	}
	return nil
}

func (r Repository) versionTagsAt(commit string) (*Version, *Version, error) {
	tagsOutput, err := r.Git.Run("", "tag", "--points-at", commit, "--list", "v*")
	if err != nil {
		return nil, nil, err
	}
	return classifyTags(strings.Fields(tagsOutput))
}

func normalizeRequestedChannel(value string) (string, error) {
	switch value {
	case "", "auto":
		return "auto", nil
	case "release", "nightly":
		return value, nil
	default:
		return "", fmt.Errorf(
			"invalid requested channel %q; want auto, release, or nightly",
			value,
		)
	}
}

func (r Repository) requireAbsentRef(ref string) error {
	_, err := r.Git.Run("", "show-ref", "--verify", "--quiet", ref)
	if err == nil {
		return fmt.Errorf("Git ref %s already exists", ref)
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
		return nil
	}
	return fmt.Errorf("check Git ref %s: %w", ref, err)
}

func (r Repository) refCommitIfPresent(ref string) (string, bool, error) {
	_, err := r.Git.Run("", "show-ref", "--verify", "--quiet", ref)
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
			return "", false, nil
		}
		return "", false, fmt.Errorf("check Git ref %s: %w", ref, err)
	}
	commit, err := r.Git.Run("", "rev-parse", "--verify", ref+"^{commit}")
	if err != nil {
		return "", false, fmt.Errorf("resolve Git ref %s: %w", ref, err)
	}
	return commit, true, nil
}

func classifyTags(tags []string) (*Version, *Version, error) {
	var stable *Version
	var nightly *Version
	for _, tag := range tags {
		version, err := ParseVersionTag(tag)
		if err != nil {
			if ownedVersionTagPattern.MatchString(tag) {
				return nil, nil, fmt.Errorf("validate version tag %s: %w", tag, err)
			}
			continue
		}
		selected := &stable
		if version.Nightly != "" {
			selected = &nightly
		}
		if *selected != nil && (*selected).String() != version.String() {
			return nil, nil, fmt.Errorf("HEAD has conflicting version tags %s and %s", (*selected).Tag(), tag)
		}
		copy := version
		*selected = &copy
	}
	return stable, nightly, nil
}
