package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func runTestGit(t *testing.T, directory string, args ...string) string {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = directory
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %q error = %v\n%s", args, err, output)
	}
	return strings.TrimSpace(string(output))
}

func runTestGitInput(
	t *testing.T,
	directory string,
	stdin string,
	args ...string,
) string {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = directory
	command.Stdin = strings.NewReader(stdin)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %q error = %v\n%s", args, err, output)
	}
	return strings.TrimSpace(string(output))
}

func writeTestFile(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestExactLeaseRejectsConcurrentRemoteChange(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	remote := filepath.Join(root, "remote.git")
	seed := filepath.Join(root, "seed")
	work := filepath.Join(root, "work")
	competitor := filepath.Join(root, "competitor")

	runTestGit(t, root, "init", "--bare", remote)
	runTestGit(t, root, "init", "--initial-branch=master", seed)
	configureTestRepository(t, seed)
	writeTestFile(t, filepath.Join(seed, ".gitignore"), "/out/\n")
	writeTestFile(t, filepath.Join(seed, "base.txt"), "base\n")
	runTestGit(t, seed, "add", ".gitignore", "base.txt")
	runTestGit(t, seed, "commit", "-m", "base")
	runTestGit(t, seed, "remote", "add", "origin", remote)
	runTestGit(t, seed, "push", "origin", "master")

	runTestGit(t, root, "clone", remote, work)
	configureTestRepository(t, work)
	runTestGit(t, work, "switch", "-c", "feature")
	writeTestFile(t, filepath.Join(work, "feature.txt"), "first\n")
	runTestGit(t, work, "add", "feature.txt")
	runTestGit(t, work, "commit", "-m", "feature")
	runTestGit(t, work, "push", "origin", "feature")

	repository, err := openGitRepository(
		ctx,
		work,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	snapshot, err := repository.snapshot(ctx, remote, "master", "feature")
	if err != nil {
		t.Fatalf("snapshot() error = %v", err)
	}
	oldRemoteHead := snapshot.RemoteHeadOID
	snapshot.close(ctx)

	writeTestFile(t, filepath.Join(work, "local.txt"), "local\n")
	runTestGit(t, work, "add", "local.txt")
	runTestGit(t, work, "commit", "-m", "local")

	runTestGit(t, root, "clone", remote, competitor)
	configureTestRepository(t, competitor)
	runTestGit(t, competitor, "switch", "feature")
	writeTestFile(t, filepath.Join(competitor, "remote.txt"), "remote\n")
	runTestGit(t, competitor, "add", "remote.txt")
	runTestGit(t, competitor, "commit", "-m", "remote")
	runTestGit(t, competitor, "push", "origin", "feature")
	competitorHead := runTestGit(t, competitor, "rev-parse", "HEAD")

	err = repository.push(
		ctx,
		remote,
		runTestGit(t, work, "rev-parse", "HEAD"),
		"refs/heads/feature",
		oldRemoteHead,
		"refs/heads/master",
		snapshot.BaseOID,
	)
	if err == nil {
		t.Fatal("push() unexpectedly overwrote concurrent remote change")
	}
	remoteHead := strings.Fields(runTestGit(
		t,
		work,
		"ls-remote",
		"--heads",
		"origin",
		"refs/heads/feature",
	))[0]
	if remoteHead != competitorHead {
		t.Fatalf("remote head = %s, want competitor %s", remoteHead, competitorHead)
	}
}

func TestStatusReportsBothSidesOfStagedRename(t *testing.T) {
	ctx := context.Background()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	if err := os.MkdirAll(filepath.Join(repositoryPath, "outside"), 0o700); err != nil {
		t.Fatalf("create outside directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(repositoryPath, "outside", "source.txt"),
		"source\n",
	)
	runTestGit(t, repositoryPath, "add", "outside/source.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	if err := os.MkdirAll(filepath.Join(repositoryPath, "allowed"), 0o700); err != nil {
		t.Fatalf("create allowed directory: %v", err)
	}
	runTestGit(
		t,
		repositoryPath,
		"mv",
		"outside/source.txt",
		"allowed/destination.txt",
	)
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	status, err := repository.status(ctx)
	if err != nil {
		t.Fatalf("status() error = %v", err)
	}
	want := map[string]bool{
		"outside/source.txt":      true,
		"allowed/destination.txt": true,
	}
	for _, value := range status.Staged {
		delete(want, value)
	}
	if len(want) != 0 {
		t.Fatalf("status().Staged = %q, missing rename paths %v", status.Staged, want)
	}
	if err := refusePathsOutside(status.Staged, []string{"allowed"}, "staged"); err == nil {
		t.Fatal("rename from outside allowed path was not refused")
	}
}

func TestRequireDefaultIndexFlagsRejectsStatusHidingEntries(t *testing.T) {
	for _, test := range []struct {
		name string
		set  func(*testing.T, string)
	}{
		{
			name: "assume_unchanged",
			set: func(t *testing.T, directory string) {
				runTestGit(t, directory, "update-index", "--assume-unchanged", "tracked.txt")
			},
		},
		{
			name: "skip_worktree",
			set: func(t *testing.T, directory string) {
				runTestGit(t, directory, "update-index", "--skip-worktree", "tracked.txt")
			},
		},
		{
			name: "intent_to_add",
			set: func(t *testing.T, directory string) {
				writeTestFile(t, filepath.Join(directory, "intent.txt"), "intent\n")
				runTestGit(t, directory, "add", "-N", "intent.txt")
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			directory := filepath.Join(t.TempDir(), "repository")
			runTestGit(t, t.TempDir(), "init", "--initial-branch=master", directory)
			configureTestRepository(t, directory)
			writeTestFile(t, filepath.Join(directory, "tracked.txt"), "tracked\n")
			writeTestFile(t, filepath.Join(directory, "outside.txt"), "outside\n")
			runTestGit(t, directory, "add", "tracked.txt", "outside.txt")
			runTestGit(t, directory, "commit", "-m", "base")
			test.set(t, directory)
			repository, err := openGitRepository(
				ctx,
				directory,
				"git",
				func(string) string { return "" },
				&execRunner{},
			)
			if err != nil {
				t.Fatalf("openGitRepository() error = %v", err)
			}
			paths := []string{"tracked.txt"}
			if test.name == "intent_to_add" {
				paths = []string{"intent.txt"}
			}
			if err := repository.requireDefaultIndexFlags(ctx, paths); err == nil ||
				!strings.Contains(err.Error(), "non-default index flags") {
				t.Fatalf("requireDefaultIndexFlags() error = %v", err)
			}
			if err := repository.requireDefaultIndexFlags(
				ctx,
				[]string{"outside.txt"},
			); err != nil {
				t.Fatalf("outside scoped flag gate error = %v", err)
			}
		})
	}
}

func TestRequireDefaultIndexFlagsSupportsLargeIndexOutput(t *testing.T) {
	ctx := context.Background()
	directory := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", directory)
	configureTestRepository(t, directory)
	object := runTestGitInput(t, directory, "shared\n", "hash-object", "-w", "--stdin")
	var entries strings.Builder
	for index := 0; index < 1_024; index++ {
		_, _ = fmt.Fprintf(
			&entries,
			"100644 %s\tlarge/%04d-%s.txt\n",
			object,
			index,
			strings.Repeat("x", 64),
		)
	}
	runTestGitInput(
		t,
		directory,
		entries.String(),
		"update-index",
		"--add",
		"--index-info",
	)
	output := runTestGit(
		t,
		directory,
		"ls-files",
		"--cached",
		"--stage",
		"--debug",
		"-z",
		"--full-name",
		"--no-recurse-submodules",
		"--",
		".",
	)
	if len(output) <= commandOutputLimit {
		t.Fatalf(
			"large index debug output = %d bytes, want more than generic limit %d",
			len(output),
			commandOutputLimit,
		)
	}
	if len(output) >= indexFlagInspectionOutputLimit {
		t.Fatalf(
			"large index fixture = %d bytes, want below operation limit %d",
			len(output),
			indexFlagInspectionOutputLimit,
		)
	}
	repository := &gitRepository{
		directory:  directory,
		executable: "git",
		runner:     &execRunner{},
	}
	if err := repository.requireDefaultIndexFlags(ctx, []string{"."}); err != nil {
		t.Fatalf("requireDefaultIndexFlags() error = %v", err)
	}
}

func TestRequireDefaultIndexFlagsRefusesOperationCeilingOverflow(t *testing.T) {
	runner := &indexFlagOverflowRunner{}
	repository := &gitRepository{
		directory:  t.TempDir(),
		executable: "git",
		runner:     runner,
	}
	err := repository.requireDefaultIndexFlags(
		context.Background(),
		[]string{"."},
	)
	if err == nil || !strings.Contains(err.Error(), "refusing truncated data") {
		t.Fatalf(
			"requireDefaultIndexFlags() error = %v, want bounded-output refusal",
			err,
		)
	}
	if runner.calls != 1 {
		t.Fatalf("index scan calls = %d, want 1 after overflow", runner.calls)
	}
}

type indexFlagOverflowRunner struct {
	calls int
}

func (r *indexFlagOverflowRunner) Run(
	_ context.Context,
	request command,
) (commandResult, error) {
	r.calls++
	if request.OutputLimit != indexFlagInspectionOutputLimit {
		return commandResult{}, fmt.Errorf(
			"index scan output limit = %d, want %d",
			request.OutputLimit,
			indexFlagInspectionOutputLimit,
		)
	}
	result := commandResult{
		// Emulate execRunner after it retained only the bounded prefix.
		Stdout:    "bounded prefix",
		Truncated: true,
		ExitCode:  0,
	}
	return result, &commandError{
		Command: command{Name: "git", OutputLimit: request.OutputLimit},
		Result:  result,
		Err:     fmt.Errorf("command output exceeded safety limit"),
	}
}

func TestOperationStateUsesLinkedWorktreeGitPath(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	mainWorktree := filepath.Join(root, "main")
	linkedWorktree := filepath.Join(root, "linked")
	runTestGit(t, root, "init", "--initial-branch=master", mainWorktree)
	configureTestRepository(t, mainWorktree)
	writeTestFile(t, filepath.Join(mainWorktree, "base.txt"), "base\n")
	runTestGit(t, mainWorktree, "add", "base.txt")
	runTestGit(t, mainWorktree, "commit", "-m", "base")
	runTestGit(
		t,
		mainWorktree,
		"worktree",
		"add",
		"-b",
		"feature",
		linkedWorktree,
	)
	marker := runTestGit(
		t,
		linkedWorktree,
		"rev-parse",
		"--git-path",
		"CHERRY_PICK_HEAD",
	)
	if !filepath.IsAbs(marker) {
		marker = filepath.Join(linkedWorktree, marker)
	}
	writeTestFile(t, marker, strings.Repeat("a", 40)+"\n")
	repository, err := openGitRepository(
		ctx,
		linkedWorktree,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	active, err := repository.operationState(ctx)
	if err != nil {
		t.Fatalf("operationState() error = %v", err)
	}
	if !containsString(active, "CHERRY_PICK_HEAD") {
		t.Fatalf("operationState() = %q, want CHERRY_PICK_HEAD", active)
	}
	if err := repository.requireNoOperation(ctx); err == nil {
		t.Fatal("requireNoOperation() unexpectedly accepted marker")
	}
}

func TestRepositoryIgnoresReplacementObjectsAndGitEnvironmentInjection(
	t *testing.T,
) {
	ctx := context.Background()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, "base.txt"), "base\n")
	runTestGit(t, repositoryPath, "add", "base.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	base := runTestGit(t, repositoryPath, "rev-parse", "HEAD")

	runTestGit(t, repositoryPath, "switch", "-c", "feature")
	writeTestFile(t, filepath.Join(repositoryPath, "feature.txt"), "original\n")
	runTestGit(t, repositoryPath, "add", "feature.txt")
	runTestGit(
		t,
		repositoryPath,
		"commit",
		"-m",
		"original feature",
		"-m",
		"Original body.\n\n"+commitDisclaimer,
	)
	feature := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	originalTree := runTestGit(t, repositoryPath, "rev-parse", "HEAD^{tree}")

	writeTestFile(t, filepath.Join(repositoryPath, "feature.txt"), "replacement\n")
	runTestGit(t, repositoryPath, "add", "feature.txt")
	replacementTree := runTestGit(t, repositoryPath, "write-tree")
	replacementCommit := "tree " + replacementTree + "\n" +
		"author Replacement <replacement@example.invalid> 1700000000 +0000\n" +
		"committer Replacement <replacement@example.invalid> 1700000000 +0000\n" +
		"gpgsig untrusted-replacement-signature\n\n" +
		"replacement feature\n\nReplacement body.\n"
	replacement := runTestGitInput(
		t,
		repositoryPath,
		replacementCommit,
		"hash-object",
		"-t",
		"commit",
		"-w",
		"--stdin",
	)
	runTestGit(t, repositoryPath, "reset", "--hard", feature)
	runTestGit(t, repositoryPath, "replace", feature, replacement)
	if got := runTestGit(
		t,
		repositoryPath,
		"show",
		"--no-patch",
		"--format=%s",
		feature,
	); got != "replacement feature" {
		t.Fatalf("replacement fixture is inactive: title = %q", got)
	}

	injectedConfig := filepath.Join(repositoryPath, "injected.gitconfig")
	writeTestFile(t, injectedConfig, "[core]\n\thooksPath = injected-hooks\n")
	injectedShallow := filepath.Join(repositoryPath, "injected-shallow")
	writeTestFile(t, injectedShallow, feature+"\n")
	injectedGrafts := filepath.Join(repositoryPath, "injected-grafts")
	writeTestFile(t, injectedGrafts, feature+" "+replacement+"\n")
	t.Setenv("GIT_CONFIG", injectedConfig)
	t.Setenv("GIT_GRAFT_FILE", injectedGrafts)
	t.Setenv("GIT_NO_REPLACE_OBJECTS", "0")
	t.Setenv("GIT_REPLACE_REF_BASE", "refs/replace/")
	t.Setenv("GIT_SHALLOW_FILE", injectedShallow)

	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	tree, err := repository.tree(ctx, feature)
	if err != nil {
		t.Fatalf("headTree() error = %v", err)
	}
	if tree != originalTree {
		t.Fatalf("headTree() = %s, want exact tree %s", tree, originalTree)
	}
	parents, err := repository.commitParents(ctx, feature)
	if err != nil {
		t.Fatalf("commitParents() error = %v", err)
	}
	if len(parents) != 1 || parents[0] != base {
		t.Fatalf("commitParents() = %q, want exact parent %s", parents, base)
	}
	projection, err := repository.projection(ctx, feature)
	if err != nil {
		t.Fatalf("projection() error = %v", err)
	}
	if projection.Title != "original feature" ||
		!strings.Contains(projection.Body, "Original body.") ||
		strings.Contains(projection.Body, "Replacement body.") {
		t.Fatalf("projection() observed replacement: %#v", projection)
	}
	signed, err := repository.commitHasSignature(ctx, feature)
	if err != nil {
		t.Fatalf("commitHasSignature() error = %v", err)
	}
	if signed {
		t.Fatal("commitHasSignature() observed replacement signature")
	}
	count, err := repository.commitCount(ctx, base, feature)
	if err != nil {
		t.Fatalf("commitCount() error = %v", err)
	}
	ancestor, err := repository.isAncestor(ctx, base, feature)
	if err != nil {
		t.Fatalf("isAncestor() error = %v", err)
	}
	commits, err := repository.featureCommits(ctx, base, feature)
	if err != nil {
		t.Fatalf("featureCommits() error = %v", err)
	}
	if count != 1 || !ancestor || len(commits) != 1 ||
		commits[0].OID != feature || commits[0].Title != "original feature" ||
		!commits[0].HasDisclaimer {
		t.Fatalf(
			"range observations used replacement: count=%d ancestor=%v commits=%#v",
			count,
			ancestor,
			commits,
		)
	}
}

func TestOpenGitRepositoryRejectsShallowClone(t *testing.T) {
	root := t.TempDir()
	remote := filepath.Join(root, "remote.git")
	seed := filepath.Join(root, "seed")
	shallow := filepath.Join(root, "shallow")
	runTestGit(t, root, "init", "--bare", remote)
	runTestGit(t, root, "init", "--initial-branch=master", seed)
	configureTestRepository(t, seed)
	writeTestFile(t, filepath.Join(seed, "base.txt"), "base\n")
	runTestGit(t, seed, "add", "base.txt")
	runTestGit(t, seed, "commit", "-m", "base")
	writeTestFile(t, filepath.Join(seed, "second.txt"), "second\n")
	runTestGit(t, seed, "add", "second.txt")
	runTestGit(t, seed, "commit", "-m", "second")
	runTestGit(t, seed, "remote", "add", "origin", remote)
	runTestGit(t, seed, "push", "origin", "master")
	runTestGit(t, root, "clone", "--depth=1", "file://"+remote, shallow)
	if got := runTestGit(
		t,
		shallow,
		"rev-parse",
		"--is-shallow-repository",
	); got != "true" {
		t.Fatalf("shallow fixture state = %q, want true", got)
	}

	_, err := openGitRepository(
		context.Background(),
		shallow,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err == nil || !strings.Contains(err.Error(), "shallow repositories") {
		t.Fatalf("openGitRepository() error = %v, want shallow refusal", err)
	}
}

func TestOpenGitRepositoryRejectsPartialAndPromisorConfiguration(t *testing.T) {
	for _, test := range []struct {
		name      string
		configure func(*testing.T, string)
	}{
		{
			name: "partial clone extension",
			configure: func(t *testing.T, repositoryPath string) {
				runTestGit(
					t,
					repositoryPath,
					"config",
					"core.repositoryFormatVersion",
					"1",
				)
				runTestGit(
					t,
					repositoryPath,
					"config",
					"extensions.partialClone",
					"origin",
				)
			},
		},
		{
			name: "promisor remote",
			configure: func(t *testing.T, repositoryPath string) {
				runTestGit(
					t,
					repositoryPath,
					"config",
					"remote.origin.promisor",
					"true",
				)
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			repositoryPath := filepath.Join(t.TempDir(), "repository")
			runTestGit(
				t,
				t.TempDir(),
				"init",
				"--initial-branch=master",
				repositoryPath,
			)
			configureTestRepository(t, repositoryPath)
			writeTestFile(t, filepath.Join(repositoryPath, "base.txt"), "base\n")
			runTestGit(t, repositoryPath, "add", "base.txt")
			runTestGit(t, repositoryPath, "commit", "-m", "base")
			test.configure(t, repositoryPath)

			repository, err := openGitRepository(
				context.Background(),
				repositoryPath,
				"git",
				func(string) string { return "" },
				&execRunner{},
			)
			if err != nil {
				t.Fatalf("openGitRepository() error = %v", err)
			}
			_, err = repository.currentBranch(context.Background())
			if err == nil || !strings.Contains(err.Error(), "partial or promisor") {
				t.Fatalf(
					"currentBranch() error = %v, want partial-clone refusal",
					err,
				)
			}
		})
	}
}

func TestOpenGitRepositoryRejectsLegacyGraftSpoof(t *testing.T) {
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, "base.txt"), "base\n")
	runTestGit(t, repositoryPath, "add", "base.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	base := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	writeTestFile(t, filepath.Join(repositoryPath, "feature.txt"), "feature\n")
	runTestGit(t, repositoryPath, "add", "feature.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "feature")
	feature := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	tree := runTestGit(t, repositoryPath, "rev-parse", "HEAD^{tree}")
	unrelated := runTestGit(
		t,
		repositoryPath,
		"commit-tree",
		tree,
		"-m",
		"unrelated root",
	)
	writeTestFile(
		t,
		filepath.Join(repositoryPath, ".git", "info", "grafts"),
		feature+" "+unrelated+"\n",
	)
	if got := runTestGit(
		t,
		repositoryPath,
		"-c",
		"advice.graftFileDeprecated=false",
		"rev-list",
		"--count",
		base+".."+feature,
	); got != "2" {
		t.Fatalf("grafted range count = %q, want spoofed count 2", got)
	}
	if _, err := runTestGitAllowFailure(
		repositoryPath,
		"-c",
		"advice.graftFileDeprecated=false",
		"merge-base",
		"--is-ancestor",
		base,
		feature,
	); err == nil {
		t.Fatal("graft fixture did not spoof ancestry")
	}

	_, err := openGitRepository(
		context.Background(),
		repositoryPath,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err == nil || !strings.Contains(err.Error(), "legacy Git grafts") {
		t.Fatalf("openGitRepository() error = %v, want graft refusal", err)
	}
}

func TestCommitCheckedPreservesSignedAmend(t *testing.T) {
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		t.Skip("ssh-keygen is unavailable")
	}
	ctx := context.Background()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, ".gitignore"), "/out/\n")
	key := testIgnoredPath(t, repositoryPath, "signing-key")
	command := exec.Command(
		"ssh-keygen",
		"-q",
		"-t",
		"ed25519",
		"-N",
		"",
		"-f",
		key,
	)
	if output, err := command.CombinedOutput(); err != nil {
		t.Skipf("cannot generate SSH signing key: %v: %s", err, output)
	}
	runTestGit(t, repositoryPath, "config", "gpg.format", "ssh")
	runTestGit(t, repositoryPath, "config", "user.signingkey", key)
	runTestGit(t, repositoryPath, "config", "commit.gpgsign", "true")
	writeTestFile(t, filepath.Join(repositoryPath, "feature.txt"), "feature\n")
	runTestGit(t, repositoryPath, "add", ".gitignore", "feature.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "signed feature")
	oldHead := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	oldAuthor := runTestGit(t, repositoryPath, "show", "-s", "--format=%an <%ae> %aI")
	runTestGit(t, repositoryPath, "config", "commit.gpgsign", "false")
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	oldContents := runTestGit(t, repositoryPath, "cat-file", "commit", oldHead)
	if !strings.Contains(oldContents, "gpgsig ") {
		t.Fatal("test setup did not create a signed commit")
	}
	signed, err := repository.commitHasSignature(ctx, oldHead)
	if err != nil {
		t.Fatalf("commitHasSignature() error = %v", err)
	}
	if !signed {
		status := runTestGit(
			t,
			repositoryPath,
			"show",
			"--no-patch",
			"--format=%G?",
			oldHead,
		)
		t.Fatalf("commitHasSignature() = false for raw signature; status %q", status)
	}
	tree, err := repository.tree(ctx, oldHead)
	if err != nil {
		t.Fatalf("headTree() error = %v", err)
	}
	if _, err := repository.commitChecked(
		ctx,
		"Updated signed feature\n",
		true,
		oldHead,
		tree,
		"master",
		nil,
	); err != nil {
		t.Fatalf("commitChecked() error = %v", err)
	}
	newHead := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	contents := runTestGit(t, repositoryPath, "cat-file", "commit", newHead)
	if !strings.Contains(contents, "gpgsig ") {
		t.Fatal("amended commit does not contain a signature")
	}
	if err := repository.verifyRequiredSignature(ctx, newHead); err != nil {
		t.Fatalf("verifyRequiredSignature() error = %v", err)
	}
	newAuthor := runTestGit(t, repositoryPath, "show", "-s", "--format=%an <%ae> %aI")
	if newAuthor != oldAuthor {
		t.Fatalf("amended author = %q, want %q", newAuthor, oldAuthor)
	}
}

func TestCommitCheckedBlocksConcurrentHeadChange(t *testing.T) {
	ctx := context.Background()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, ".gitignore"), "/out/\n")
	writeTestFile(t, filepath.Join(repositoryPath, "base.txt"), "base\n")
	runTestGit(t, repositoryPath, "add", ".gitignore", "base.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	oldHead := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	runTestGit(t, repositoryPath, "branch", "intruder", oldHead)
	writeTestFile(t, filepath.Join(repositoryPath, "feature.txt"), "feature\n")
	runTestGit(t, repositoryPath, "add", "feature.txt")
	expectedTree := runTestGit(t, repositoryPath, "write-tree")
	runner := &attachedIndexHookRunner{
		delegate: &execRunner{},
		hook: func() error {
			command := exec.Command("git", "switch", "intruder")
			command.Dir = repositoryPath
			if output, err := command.CombinedOutput(); err == nil {
				return fmt.Errorf(
					"concurrent checkout unexpectedly succeeded: %s",
					output,
				)
			}
			return nil
		},
	}
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		runner,
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	newHead, err := repository.commitChecked(
		ctx,
		"feature\n",
		false,
		oldHead,
		expectedTree,
		"master",
		nil,
	)
	if err != nil {
		t.Fatalf("commitChecked() error = %v", err)
	}
	if !runner.fired {
		t.Fatal("concurrent checkout attempt did not run under transaction locks")
	}
	branch, head, err := repository.branchHead(ctx)
	if err != nil {
		t.Fatalf("branchHead() error = %v", err)
	}
	if branch != "master" || head != newHead {
		t.Fatalf("checkout = %s@%s, want master@%s", branch, head, newHead)
	}
	if _, err := os.Lstat(filepath.Join(repositoryPath, ".git", "HEAD.lock")); !os.IsNotExist(err) {
		t.Fatalf("HEAD lock was not cleaned up: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(repositoryPath, ".git", "index.lock")); !os.IsNotExist(err) {
		t.Fatalf("index lock was not cleaned up: %v", err)
	}
}

func TestDetachExactBranchRejectsAnotherSameOIDBranch(t *testing.T) {
	ctx := context.Background()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=feature", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, "base.txt"), "base\n")
	runTestGit(t, repositoryPath, "add", "base.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	head := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	runTestGit(t, repositoryPath, "branch", "intruder", head)
	runTestGit(t, repositoryPath, "symbolic-ref", "HEAD", "refs/heads/intruder")
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	if err := repository.detachExactBranch(ctx, "feature", head); err == nil ||
		!strings.Contains(err.Error(), "symbolic HEAD changed") {
		t.Fatalf("detachExactBranch() error = %v, want same-OID branch refusal", err)
	}
	if got := runTestGit(t, repositoryPath, "symbolic-ref", "--short", "HEAD"); got != "intruder" {
		t.Fatalf("HEAD = %q, want preserved intruder branch", got)
	}
}

func TestHEADLockReleasePreservesAReplacementOwner(t *testing.T) {
	ctx := context.Background()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, "base.txt"), "base\n")
	runTestGit(t, repositoryPath, "add", "base.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	head := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	lock, err := repository.lockHEAD(ctx)
	if err != nil {
		t.Fatalf("lockHEAD() error = %v", err)
	}
	if err := lock.commitOID(head); err != nil {
		t.Fatalf("commitOID() error = %v", err)
	}
	replacement, err := os.OpenFile(
		lock.path,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0o600,
	)
	if err != nil {
		t.Fatalf("create replacement HEAD lock: %v", err)
	}
	if _, err := replacement.WriteString("replacement owner\n"); err != nil {
		t.Fatalf("write replacement HEAD lock: %v", err)
	}
	if err := replacement.Close(); err != nil {
		t.Fatalf("close replacement HEAD lock: %v", err)
	}
	if err := lock.release(); err != nil {
		t.Fatalf("release() error = %v", err)
	}
	contents, err := os.ReadFile(lock.path)
	if err != nil {
		t.Fatalf("replacement HEAD lock was removed: %v", err)
	}
	if string(contents) != "replacement owner\n" {
		t.Fatalf("replacement HEAD lock = %q", contents)
	}
	if err := os.Remove(lock.path); err != nil {
		t.Fatalf("remove replacement HEAD lock: %v", err)
	}
}

func TestCommitCheckedRestoresOriginalBranchWhenAttachmentFails(t *testing.T) {
	ctx := context.Background()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, ".gitignore"), "/out/\n")
	writeTestFile(t, filepath.Join(repositoryPath, "base.txt"), "base\n")
	runTestGit(t, repositoryPath, "add", ".gitignore", "base.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	originalHead := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	writeTestFile(t, filepath.Join(repositoryPath, "feature.txt"), "feature\n")
	runTestGit(t, repositoryPath, "add", "feature.txt")
	expectedTree := runTestGit(t, repositoryPath, "write-tree")
	runner := &failAttachedIndexInspectionRunner{delegate: &execRunner{}}
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		runner,
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	_, err = repository.commitChecked(
		ctx,
		"feature\n",
		false,
		originalHead,
		expectedTree,
		"master",
		nil,
	)
	if err == nil || !strings.Contains(err.Error(), "attached branch index gate") {
		t.Fatalf("commitChecked() error = %v, want pre-commit index refusal", err)
	}
	branch, head, err := repository.branchHead(ctx)
	if err != nil {
		t.Fatalf("branchHead() error = %v", err)
	}
	if branch != "master" || head != originalHead {
		t.Fatalf("checkout = %s@%s, want master@%s", branch, head, originalHead)
	}
	if got := runTestGit(t, repositoryPath, "rev-parse", "master"); got != originalHead {
		t.Fatalf("master = %s, want restored %s", got, originalHead)
	}
	if runner.failures != 1 {
		t.Fatalf("injected index-inspection failures = %d, want 1", runner.failures)
	}
}

func TestRebasePreservesCapturedSignatureRequirement(t *testing.T) {
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		t.Skip("ssh-keygen is unavailable")
	}
	ctx := context.Background()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, ".gitignore"), "/out/\n")
	writeTestFile(t, filepath.Join(repositoryPath, "base.txt"), "base\n")
	runTestGit(t, repositoryPath, "add", ".gitignore", "base.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	key := testIgnoredPath(t, repositoryPath, "signing-key")
	command := exec.Command(
		"ssh-keygen",
		"-q",
		"-t",
		"ed25519",
		"-N",
		"",
		"-f",
		key,
	)
	if output, err := command.CombinedOutput(); err != nil {
		t.Skipf("cannot generate SSH signing key: %v: %s", err, output)
	}
	runTestGit(t, repositoryPath, "config", "gpg.format", "ssh")
	runTestGit(t, repositoryPath, "config", "user.signingkey", key)
	runTestGit(t, repositoryPath, "config", "commit.gpgsign", "true")
	runTestGit(t, repositoryPath, "switch", "-c", "feature")
	writeTestFile(t, filepath.Join(repositoryPath, "feature.txt"), "feature\n")
	runTestGit(t, repositoryPath, "add", "feature.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "signed feature")
	oldFeature := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	runTestGit(t, repositoryPath, "config", "commit.gpgsign", "false")
	runTestGit(t, repositoryPath, "switch", "master")
	writeTestFile(t, filepath.Join(repositoryPath, "advance.txt"), "advance\n")
	runTestGit(t, repositoryPath, "add", "advance.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "advance")
	base := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	runTestGit(t, repositoryPath, "switch", "feature")

	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	if _, err := repository.rebase(ctx, base, "feature", oldFeature, nil); err != nil {
		t.Fatalf("rebase() error = %v", err)
	}
	newFeature, err := repository.head(ctx)
	if err != nil {
		t.Fatalf("head() error = %v", err)
	}
	if newFeature == oldFeature {
		t.Fatal("rebase did not rewrite signed feature")
	}
	signed, err := repository.commitHasSignature(ctx, newFeature)
	if err != nil {
		t.Fatalf("commitHasSignature() error = %v", err)
	}
	if !signed {
		t.Fatal("rebased commit lost captured signature requirement")
	}
}

func TestRequireSingleCommitRejectsEmptyAggregateDiff(t *testing.T) {
	ctx := context.Background()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, "base.txt"), "base\n")
	runTestGit(t, repositoryPath, "add", "base.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	base := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	runTestGit(t, repositoryPath, "commit", "--allow-empty", "-m", "empty")
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	delivery := &delivery{repository: repository}
	head := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	if err := delivery.requireSingleCommit(ctx, base, head); err == nil ||
		!strings.Contains(err.Error(), "empty aggregate diff") {
		t.Fatalf("requireSingleCommit() error = %v, want empty diff refusal", err)
	}
}

func TestRebaseDoesNotUpdateSiblingRefs(t *testing.T) {
	ctx := context.Background()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, ".gitignore"), "/out/\n")
	writeTestFile(t, filepath.Join(repositoryPath, "base.txt"), "base\n")
	runTestGit(t, repositoryPath, "add", ".gitignore", "base.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	runTestGit(t, repositoryPath, "switch", "-c", "feature")
	writeTestFile(t, filepath.Join(repositoryPath, "feature.txt"), "feature\n")
	runTestGit(t, repositoryPath, "add", "feature.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "feature")
	oldFeature := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	runTestGit(t, repositoryPath, "branch", "sibling", oldFeature)
	runTestGit(t, repositoryPath, "switch", "master")
	writeTestFile(t, filepath.Join(repositoryPath, "advance.txt"), "advance\n")
	runTestGit(t, repositoryPath, "add", "advance.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "advance")
	newBase := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	runTestGit(t, repositoryPath, "switch", "feature")
	runTestGit(t, repositoryPath, "config", "rebase.updateRefs", "true")
	hookDirectory := filepath.Join(repositoryPath, ".git", "hooks")
	preRebaseHook := filepath.Join(hookDirectory, "pre-rebase")
	postRewriteHook := filepath.Join(hookDirectory, "post-rewrite")
	referenceTransactionHook := filepath.Join(
		hookDirectory,
		"reference-transaction",
	)
	writeTestFile(t, preRebaseHook, "#!/bin/sh\nexit 91\n")
	writeTestFile(
		t,
		postRewriteHook,
		"#!/bin/sh\ntouch post-rewrite-ran\n",
	)
	writeTestFile(t, referenceTransactionHook, "#!/bin/sh\nexit 92\n")
	for _, hook := range []string{
		preRebaseHook,
		postRewriteHook,
		referenceTransactionHook,
	} {
		if err := os.Chmod(hook, 0o700); err != nil {
			t.Fatalf("make hook executable: %v", err)
		}
	}
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	rebasedHead, err := repository.rebase(
		ctx,
		newBase,
		"feature",
		oldFeature,
		nil,
	)
	if err != nil {
		t.Fatalf("rebase() error = %v", err)
	}
	if err := repository.requireNoOperation(ctx); err != nil {
		t.Fatalf("successful rebase left active operation state: %v", err)
	}
	if got := runTestGit(t, repositoryPath, "rev-parse", "sibling"); got != oldFeature {
		t.Fatalf("sibling ref = %s, want unchanged %s", got, oldFeature)
	}
	if got := runTestGit(t, repositoryPath, "rev-parse", "HEAD"); got != rebasedHead {
		t.Fatalf("HEAD = %s, want returned rebased OID %s", got, rebasedHead)
	}
	if got := runTestGit(
		t,
		repositoryPath,
		"symbolic-ref",
		"--short",
		"HEAD",
	); got != "feature" {
		t.Fatalf("current branch = %q, want feature", got)
	}
	if _, err := os.Stat(filepath.Join(repositoryPath, "post-rewrite-ran")); !os.IsNotExist(err) {
		t.Fatalf("post-rewrite hook ran or could not be inspected: %v", err)
	}
}

func TestAttachedRebaseRejectsConcurrentMainCheckout(t *testing.T) {
	ctx := context.Background()
	repositoryPath, base, oldFeature := newRebaseFixture(t)
	runTestGit(t, repositoryPath, "branch", "intruder", base)
	runner := &attachedTransitionHookRunner{
		delegate: &execRunner{},
		hook: func() error {
			command := exec.Command("git", "switch", "intruder")
			command.Dir = repositoryPath
			if output, err := command.CombinedOutput(); err == nil {
				return fmt.Errorf(
					"concurrent main checkout unexpectedly succeeded: %s",
					output,
				)
			}
			return nil
		},
	}
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		runner,
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	rebasedHead, err := repository.rebase(ctx, base, "feature", oldFeature, nil)
	if err != nil {
		t.Fatalf("rebase() error = %v", err)
	}
	if !runner.fired {
		t.Fatal("concurrent main checkout attempt did not run during attached transition")
	}
	branch, head, err := repository.branchHead(ctx)
	if err != nil {
		t.Fatalf("branchHead() error = %v", err)
	}
	if branch != "feature" || head != rebasedHead {
		t.Fatalf("checkout = %s@%s, want feature@%s", branch, head, rebasedHead)
	}
}

func TestAttachedRebaseRejectsSiblingWorktreeClaim(t *testing.T) {
	ctx := context.Background()
	repositoryPath, base, oldFeature := newRebaseFixture(t)
	sibling := filepath.Join(t.TempDir(), "sibling")
	runTestGit(
		t,
		repositoryPath,
		"worktree",
		"add",
		"-b",
		"sibling_checkout",
		sibling,
		base,
	)
	runner := &attachedTransitionHookRunner{
		delegate: &execRunner{},
		hook: func() error {
			command := exec.Command("git", "switch", "feature")
			command.Dir = sibling
			if output, err := command.CombinedOutput(); err == nil {
				return fmt.Errorf(
					"sibling worktree unexpectedly claimed feature: %s",
					output,
				)
			}
			return nil
		},
	}
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		runner,
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	rebasedHead, err := repository.rebase(ctx, base, "feature", oldFeature, nil)
	if err != nil {
		t.Fatalf("rebase() error = %v", err)
	}
	if !runner.fired {
		t.Fatal("sibling claim attempt did not run during attached transition")
	}
	if got := runTestGit(t, sibling, "symbolic-ref", "--short", "HEAD"); got != "sibling_checkout" {
		t.Fatalf("sibling branch = %q, want sibling_checkout", got)
	}
	branch, head, err := repository.branchHead(ctx)
	if err != nil {
		t.Fatalf("branchHead() error = %v", err)
	}
	if branch != "feature" || head != rebasedHead {
		t.Fatalf("checkout = %s@%s, want feature@%s", branch, head, rebasedHead)
	}
}

func TestAttachedRebaseVerificationFailureRollsBackExactly(t *testing.T) {
	ctx := context.Background()
	repositoryPath, base, oldFeature := newRebaseFixture(t)
	originalTree := runTestGit(t, repositoryPath, "write-tree")
	runner := &failAttachedTransitionVerificationRunner{
		delegate: &execRunner{},
	}
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		runner,
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	_, err = repository.rebase(ctx, base, "feature", oldFeature, nil)
	if err == nil || !strings.Contains(err.Error(), "rolled back") {
		t.Fatalf("rebase() error = %v, want exact transition rollback", err)
	}
	if !runner.failed {
		t.Fatal("attached transition verification failure did not fire")
	}
	branch, head, err := repository.branchHead(ctx)
	if err != nil {
		t.Fatalf("branchHead() error = %v", err)
	}
	if branch != "feature" || head != oldFeature {
		t.Fatalf("checkout = %s@%s, want feature@%s", branch, head, oldFeature)
	}
	if got := runTestGit(t, repositoryPath, "write-tree"); got != originalTree {
		t.Fatalf("index tree = %s, want restored %s", got, originalTree)
	}
	if got := runTestGit(t, repositoryPath, "status", "--porcelain"); got != "" {
		t.Fatalf("restored status = %q, want clean", got)
	}
	if _, err := os.Lstat(filepath.Join(repositoryPath, "advance.txt")); !os.IsNotExist(err) {
		t.Fatalf("base-only file survived exact transition rollback: %v", err)
	}
}

func TestRebaseConflictLeavesMainWorktreeUntouchedAndCleansIsolation(
	t *testing.T,
) {
	ctx := context.Background()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, ".gitignore"), "/out/\n")
	writeTestFile(t, filepath.Join(repositoryPath, "shared.txt"), "base\n")
	runTestGit(t, repositoryPath, "add", ".gitignore", "shared.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	runTestGit(t, repositoryPath, "switch", "-c", "feature")
	writeTestFile(t, filepath.Join(repositoryPath, "shared.txt"), "feature\n")
	runTestGit(t, repositoryPath, "add", "shared.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "feature")
	oldFeature := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	runTestGit(t, repositoryPath, "switch", "master")
	writeTestFile(t, filepath.Join(repositoryPath, "shared.txt"), "master\n")
	runTestGit(t, repositoryPath, "add", "shared.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "advance")
	base := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	runTestGit(t, repositoryPath, "switch", "feature")
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	_, err = repository.rebase(ctx, base, "feature", oldFeature, nil)
	if err == nil || !strings.Contains(err.Error(), "isolated exact-OID rebase") {
		t.Fatalf("rebase() error = %v, want isolated conflict", err)
	}
	branch, head, err := repository.branchHead(ctx)
	if err != nil {
		t.Fatalf("branchHead() error = %v", err)
	}
	if branch != "feature" || head != oldFeature {
		t.Fatalf("main checkout = %s@%s, want feature@%s", branch, head, oldFeature)
	}
	status, err := repository.status(ctx)
	if err != nil {
		t.Fatalf("status() error = %v", err)
	}
	if !status.clean() {
		t.Fatalf("main worktree changed after isolated conflict: %#v", status)
	}
	if err := repository.requireNoOperation(ctx); err != nil {
		t.Fatalf("main worktree has rebase state: %v", err)
	}
	worktrees := runTestGit(t, repositoryPath, "worktree", "list", "--porcelain")
	if strings.Contains(worktrees, "out/repo_delivery/rebase-") {
		t.Fatalf("isolated worktree was not removed:\n%s", worktrees)
	}
}

func TestRebasePinsExactHeadAcrossAmbientBranchChange(t *testing.T) {
	ctx := context.Background()
	repositoryPath, base, oldFeature := newRebaseFixture(t)
	runTestGit(t, repositoryPath, "branch", "intruder", oldFeature)
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("find git: %v", err)
	}
	wrapper := testIgnoredPath(t, repositoryPath, "git-wrapper")
	writeTestFile(t, wrapper, fmt.Sprintf(`#!/bin/sh
for argument in "$@"; do
    if [ "$argument" = rebase ]; then
        %q symbolic-ref HEAD refs/heads/intruder || exit $?
        break
    fi
done
exec %q "$@"
`, realGit, realGit))
	if err := os.Chmod(wrapper, 0o700); err != nil {
		t.Fatalf("make Git wrapper executable: %v", err)
	}
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		wrapper,
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	if _, err := repository.rebase(ctx, base, "feature", oldFeature, nil); err != nil {
		t.Fatalf("rebase() error = %v", err)
	}
	if got := runTestGit(
		t,
		repositoryPath,
		"symbolic-ref",
		"--short",
		"HEAD",
	); got != "feature" {
		t.Fatalf("current branch = %q, want feature", got)
	}
	if got := runTestGit(
		t,
		repositoryPath,
		"rev-parse",
		"intruder",
	); got != oldFeature {
		t.Fatalf("intruder ref = %s, want unchanged %s", got, oldFeature)
	}
	if got := runTestGit(
		t,
		repositoryPath,
		"rev-parse",
		"feature^",
	); got != base {
		t.Fatalf("rebased parent = %s, want exact base %s", got, base)
	}
}

func TestRebaseBranchCASPreservesConcurrentRefUpdate(t *testing.T) {
	ctx := context.Background()
	repositoryPath, base, oldFeature := newRebaseFixture(t)
	runner := &attachedTransitionHookRunner{
		delegate: &execRunner{},
		hook: func() error {
			command := exec.Command(
				"git",
				"update-ref",
				"refs/heads/feature",
				base,
				oldFeature,
			)
			command.Dir = repositoryPath
			if output, err := command.CombinedOutput(); err == nil {
				return fmt.Errorf(
					"concurrent branch update unexpectedly succeeded: %s",
					output,
				)
			}
			return nil
		},
	}
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		runner,
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	rebasedHead, err := repository.rebase(ctx, base, "feature", oldFeature, nil)
	if err != nil {
		t.Fatalf("rebase() error = %v", err)
	}
	if !runner.fired {
		t.Fatal("concurrent branch update did not run under the branch lock")
	}
	if got := runTestGit(t, repositoryPath, "rev-parse", "feature"); got != rebasedHead {
		t.Fatalf("feature ref = %s, want rebased %s", got, rebasedHead)
	}
}

func TestRebaseHeadCASPreservesConcurrentHeadChange(t *testing.T) {
	ctx := context.Background()
	repositoryPath, base, oldFeature := newRebaseFixture(t)
	runTestGit(t, repositoryPath, "branch", "intruder", oldFeature)
	runner := &attachedTransitionHookRunner{
		delegate: &execRunner{},
		hook: func() error {
			command := exec.Command("git", "switch", "intruder")
			command.Dir = repositoryPath
			if output, err := command.CombinedOutput(); err == nil {
				return fmt.Errorf(
					"concurrent checkout unexpectedly succeeded: %s",
					output,
				)
			}
			return nil
		},
	}
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		runner,
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	rebasedHead, err := repository.rebase(ctx, base, "feature", oldFeature, nil)
	if err != nil {
		t.Fatalf("rebase() error = %v", err)
	}
	if !runner.fired {
		t.Fatal("concurrent checkout attempt did not run under transaction locks")
	}
	if got := runTestGit(t, repositoryPath, "rev-parse", "feature"); got != rebasedHead {
		t.Fatalf("feature ref = %s, want rebased %s", got, rebasedHead)
	}
	if got := runTestGit(
		t,
		repositoryPath,
		"symbolic-ref",
		"--short",
		"HEAD",
	); got != "feature" {
		t.Fatalf("HEAD branch = %q, want continuously attached feature", got)
	}
	if got := runTestGit(t, repositoryPath, "rev-parse", "intruder"); got != oldFeature {
		t.Fatalf("intruder ref = %s, want unchanged %s", got, oldFeature)
	}
}

func TestRebaseRestoresOriginalBranchWhenAttachmentFails(t *testing.T) {
	ctx := context.Background()
	repositoryPath, base, originalHead := newRebaseFixture(t)
	runner := &failAttachedTransitionVerificationRunner{delegate: &execRunner{}}
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		runner,
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	_, err = repository.rebase(ctx, base, "feature", originalHead, nil)
	if err == nil || !strings.Contains(err.Error(), "was rolled back") {
		t.Fatalf("rebase() error = %v, want exact transition rollback", err)
	}
	branch, head, err := repository.branchHead(ctx)
	if err != nil {
		t.Fatalf("branchHead() error = %v", err)
	}
	if branch != "feature" || head != originalHead {
		t.Fatalf("checkout = %s@%s, want feature@%s", branch, head, originalHead)
	}
	if got := runTestGit(t, repositoryPath, "rev-parse", "feature"); got != originalHead {
		t.Fatalf("feature = %s, want restored %s", got, originalHead)
	}
	originalTree := runTestGit(t, repositoryPath, "rev-parse", originalHead+"^{tree}")
	if got := runTestGit(t, repositoryPath, "write-tree"); got != originalTree {
		t.Fatalf("restored index tree = %s, want %s", got, originalTree)
	}
	if got := runTestGit(t, repositoryPath, "status", "--porcelain"); got != "" {
		t.Fatalf("restored checkout status = %q, want clean", got)
	}
	if _, err := os.Lstat(filepath.Join(repositoryPath, "advance.txt")); !os.IsNotExist(err) {
		t.Fatalf("base-only file survived failed rebase rollback: %v", err)
	}
	contents, err := os.ReadFile(filepath.Join(repositoryPath, "feature.txt"))
	if err != nil || string(contents) != "feature\n" {
		t.Fatalf("feature content after rollback = %q, %v", contents, err)
	}
	if !runner.failed {
		t.Fatal("attached transition verification failure was not injected")
	}
}

func TestRebaseCleanupFailureLeavesOriginalBranchUntouched(t *testing.T) {
	ctx := context.Background()
	repositoryPath, base, originalHead := newRebaseFixture(t)
	runner := &failFirstWorktreeRemoveRunner{delegate: &execRunner{}}
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		runner,
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	_, err = repository.rebase(ctx, base, "feature", originalHead, nil)
	if err == nil || !strings.Contains(err.Error(), "clean isolated rebase worktree") {
		t.Fatalf("rebase() error = %v, want cleanup refusal", err)
	}
	branch, head, err := repository.branchHead(ctx)
	if err != nil {
		t.Fatalf("branchHead() error = %v", err)
	}
	if branch != "feature" || head != originalHead {
		t.Fatalf("checkout = %s@%s, want feature@%s", branch, head, originalHead)
	}
	if runner.attempts != 2 {
		t.Fatalf("worktree cleanup attempts = %d, want explicit attempt and deferred retry", runner.attempts)
	}
	entries, err := os.ReadDir(filepath.Join(repositoryPath, "out", "repo_delivery"))
	if err != nil {
		t.Fatalf("read isolated rebase root: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("isolated rebase residue = %v, want none", entries)
	}
}

func TestRebaseDisablesConfiguredAutostash(t *testing.T) {
	ctx := context.Background()
	repositoryPath, base, oldFeature := newRebaseFixture(t)
	runTestGit(t, repositoryPath, "config", "rebase.autoStash", "true")
	writeTestFile(t, filepath.Join(repositoryPath, "feature.txt"), "dirty\n")
	repository, err := openGitRepository(
		ctx,
		repositoryPath,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	_, err = repository.rebase(ctx, base, "feature", oldFeature, nil)
	if err == nil || !strings.Contains(err.Error(), "clean worktree") {
		t.Fatalf("rebase() error = %v, want clean-worktree refusal", err)
	}
	if got := runTestGit(t, repositoryPath, "stash", "list"); got != "" {
		t.Fatalf("stash list = %q, want empty", got)
	}
	contents, err := os.ReadFile(filepath.Join(repositoryPath, "feature.txt"))
	if err != nil {
		t.Fatalf("read dirty file: %v", err)
	}
	if string(contents) != "dirty\n" {
		t.Fatalf("dirty file = %q, want preserved", contents)
	}
}

func TestBoundTransportIgnoresRepositoryAndGlobalRoutingConfiguration(t *testing.T) {
	fixture := newPushSafetyFixture(t)
	attacker := filepath.Join(t.TempDir(), "attacker.git")
	globalConfig := filepath.Join(t.TempDir(), "global.gitconfig")
	writeTestFile(
		t,
		globalConfig,
		"[url \""+attacker+"\"]\n\tinsteadOf = "+fixture.remote+"\n",
	)
	t.Setenv("GIT_CONFIG_GLOBAL", globalConfig)
	for _, name := range gitTransportUnsetEnvironment {
		t.Setenv(name, "redirected-by-test")
	}
	t.Setenv("GIT_NO_LAZY_FETCH", "0")
	t.Setenv("SSH_AUTH_SOCK", "preserved-agent-socket")
	runTestGit(
		t,
		fixture.work,
		"config",
		"url."+attacker+".insteadOf",
		fixture.remote,
	)
	runTestGit(
		t,
		fixture.work,
		"config",
		"url."+attacker+".pushInsteadOf",
		fixture.remote,
	)
	runTestGit(t, fixture.work, "config", "core.sshCommand", "attacker-ssh")
	runTestGit(t, fixture.work, "config", "extensions.worktreeConfig", "true")
	runTestGit(
		t,
		fixture.work,
		"config",
		"--worktree",
		"url."+attacker+".insteadOf",
		fixture.remote,
	)
	audit := &transportEnvironmentAuditRunner{
		delegate:   &execRunner{},
		repository: fixture.work,
	}
	fixture.repository.runner = audit
	snapshot, err := fixture.repository.snapshot(
		context.Background(),
		fixture.remote,
		"master",
		"feature",
	)
	if err != nil {
		t.Fatalf("snapshot() through isolated transport error = %v", err)
	}
	if snapshot.BaseOID != fixture.baseOID {
		t.Fatalf("snapshot base = %s, want intended %s", snapshot.BaseOID, fixture.baseOID)
	}
	if err := snapshot.close(context.Background()); err != nil {
		t.Fatalf("snapshot.close() error = %v", err)
	}
	if !audit.sawTransport {
		t.Fatal("transport environment audit did not observe network Git")
	}
}

func TestNewDeliveryCapturesRawSSHAndRejectsHTTPSBeforeNetwork(t *testing.T) {
	t.Run("raw SSH ignores rewrite configuration", func(t *testing.T) {
		repositoryPath := newDeliveryConstructorRepository(t)
		fetchEndpoint := "git@github.com:owner/repository.git"
		pushEndpoint := "ssh://git@github.com/owner/repository.git"
		runTestGit(t, repositoryPath, "remote", "add", "origin", fetchEndpoint)
		runTestGit(
			t,
			repositoryPath,
			"remote",
			"set-url",
			"--push",
			"origin",
			pushEndpoint,
		)
		runTestGit(
			t,
			repositoryPath,
			"config",
			"url.ssh://git@attacker.invalid/owner/repository.git.insteadOf",
			fetchEndpoint,
		)
		runTestGit(t, repositoryPath, "config", "extensions.worktreeConfig", "true")
		runTestGit(
			t,
			repositoryPath,
			"config",
			"--worktree",
			"url.ssh://git@attacker.invalid/owner/repository.git.pushInsteadOf",
			pushEndpoint,
		)
		runner := &constructorIsolationRunner{delegate: &execRunner{}}
		delivery, err := newDelivery(
			context.Background(),
			deliveryConfig{
				Repository: repositoryPath,
				Remote:     "origin",
				ForgeCLI:   "gh-test",
			},
			func(string) string { return "" },
			runner,
		)
		if err != nil {
			t.Fatalf("newDelivery() SSH error = %v", err)
		}
		if delivery.fetchEndpoint != fetchEndpoint ||
			delivery.pushEndpoint != pushEndpoint {
			t.Fatalf(
				"captured endpoints = %q, %q; want raw SSH endpoints",
				delivery.fetchEndpoint,
				delivery.pushEndpoint,
			)
		}
		if runner.forgeCalls != 1 || runner.gitNetworkCalls != 0 {
			t.Fatalf(
				"constructor calls: forge=%d network Git=%d",
				runner.forgeCalls,
				runner.gitNetworkCalls,
			)
		}
		runTestGit(
			t,
			repositoryPath,
			"remote",
			"set-url",
			"--push",
			"origin",
			"ssh://git@attacker.invalid/owner/repository.git",
		)
		if delivery.pushEndpoint != pushEndpoint {
			t.Fatal("captured push endpoint changed with mutable repository config")
		}
	})

	t.Run("HTTPS refuses before forge or Git network", func(t *testing.T) {
		repositoryPath := newDeliveryConstructorRepository(t)
		secret := "constructor-secret"
		runTestGit(
			t,
			repositoryPath,
			"remote",
			"add",
			"origin",
			"https://user:"+secret+"@github.com/owner/repository.git",
		)
		runner := &constructorIsolationRunner{delegate: &execRunner{}}
		_, err := newDelivery(
			context.Background(),
			deliveryConfig{
				Repository: repositoryPath,
				Remote:     "origin",
				ForgeCLI:   "gh-test",
			},
			func(string) string { return "" },
			runner,
		)
		if err == nil || !strings.Contains(err.Error(), "requires SSH") {
			t.Fatalf("newDelivery() HTTPS error = %v, want SSH refusal", err)
		}
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("newDelivery() error leaked credential: %v", err)
		}
		if runner.forgeCalls != 0 || runner.gitNetworkCalls != 0 {
			t.Fatalf(
				"refused constructor calls: forge=%d network Git=%d",
				runner.forgeCalls,
				runner.gitNetworkCalls,
			)
		}
	})
}

func newDeliveryConstructorRepository(t *testing.T) string {
	t.Helper()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, ".gitignore"), "/out/\n")
	writeTestFile(t, filepath.Join(repositoryPath, "base.txt"), "base\n")
	runTestGit(t, repositoryPath, "add", ".gitignore", "base.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	return repositoryPath
}

func TestPushIgnoresConfiguredPushInsteadOf(t *testing.T) {
	fixture := newPushSafetyFixture(t)
	attacker := filepath.Join(t.TempDir(), "attacker.git")
	runTestGit(
		t,
		fixture.work,
		"config",
		"url."+attacker+".pushInsteadOf",
		fixture.remote,
	)
	if err := fixture.repository.push(
		context.Background(),
		fixture.remote,
		fixture.localHead,
		"refs/heads/feature",
		fixture.remoteHead,
		"refs/heads/master",
		fixture.baseOID,
	); err != nil {
		t.Fatalf("push() through isolated transport error = %v", err)
	}
	if got := remoteRefOID(
		t,
		fixture.work,
		fixture.remote,
		"refs/heads/feature",
	); got != fixture.localHead {
		t.Fatalf("intended remote head = %s, want %s", got, fixture.localHead)
	}
}

func TestSnapshotCleansPrivateRefsAfterAtomicFetchFailure(t *testing.T) {
	fixture := newPushSafetyFixture(t)
	runner := &partialFetchFailureRunner{delegate: &execRunner{}}
	fixture.repository.runner = runner

	_, err := fixture.repository.snapshot(
		context.Background(),
		fixture.remote,
		"master",
		"feature",
	)
	if err == nil || !strings.Contains(err.Error(), "fetch exact base") {
		t.Fatalf("snapshot() error = %v, want fetch failure", err)
	}
	if !runner.sawAtomic || !runner.sawNoPrune || !runner.sawNoPruneTags {
		t.Fatalf(
			"fetch flags: atomic=%v no-prune=%v no-prune-tags=%v",
			runner.sawAtomic,
			runner.sawNoPrune,
			runner.sawNoPruneTags,
		)
	}
	if got := runTestGit(
		t,
		fixture.work,
		"for-each-ref",
		"--format=%(refname)",
		"refs/repo-delivery/",
	); got != "" {
		t.Fatalf("failed snapshot left private refs: %q", got)
	}
	assertNoTransportResidue(t, fixture.work)
}

func TestSnapshotLostAcknowledgementCleansOwnedMainRef(t *testing.T) {
	fixture := newPushSafetyFixture(t)
	runner := &lostAckSnapshotInstallRunner{delegate: &execRunner{}}
	fixture.repository.runner = runner
	_, err := fixture.repository.snapshot(
		context.Background(),
		fixture.remote,
		"master",
		"feature",
	)
	if err == nil || !strings.Contains(err.Error(), "reported failure after installing") {
		t.Fatalf("snapshot() lost-ack error = %v", err)
	}
	if !runner.fired {
		t.Fatal("snapshot ref lost-ack injection did not run")
	}
	if got := runTestGit(
		t,
		fixture.work,
		"for-each-ref",
		"--format=%(refname)",
		"refs/repo-delivery/",
	); got != "" {
		t.Fatalf("lost-ack snapshot left private refs: %q", got)
	}
	assertNoTransportResidue(t, fixture.work)
}

func TestSnapshotMakesRemoteOnlyObjectsVisibleAndCleansTransport(t *testing.T) {
	fixture := newPushSafetyFixture(t)
	writeTestFile(t, filepath.Join(fixture.seed, "remote-only.txt"), "remote only\n")
	runTestGit(t, fixture.seed, "add", "remote-only.txt")
	runTestGit(t, fixture.seed, "commit", "-m", "remote-only base")
	runTestGit(t, fixture.seed, "push", "origin", "master")
	remoteOnlyOID := runTestGit(t, fixture.seed, "rev-parse", "HEAD")
	if _, err := runTestGitAllowFailure(
		fixture.work,
		"cat-file",
		"-e",
		remoteOnlyOID+"^{commit}",
	); err == nil {
		t.Fatal("remote-only object unexpectedly existed before isolated fetch")
	}
	snapshot, err := fixture.repository.snapshot(
		context.Background(),
		fixture.remote,
		"master",
		"feature",
	)
	if err != nil {
		t.Fatalf("snapshot() remote-only object error = %v", err)
	}
	if snapshot.BaseOID != remoteOnlyOID {
		t.Fatalf("snapshot base = %s, want %s", snapshot.BaseOID, remoteOnlyOID)
	}
	if _, err := runTestGitAllowFailure(
		fixture.work,
		"cat-file",
		"-e",
		remoteOnlyOID+"^{commit}",
	); err != nil {
		t.Fatalf("fetched object is not visible in the main object database: %v", err)
	}
	assertNoTransportResidue(t, fixture.work)
	if err := snapshot.close(context.Background()); err != nil {
		t.Fatalf("snapshot.close() error = %v", err)
	}
}

func TestSnapshotCleanupDoesNotDeleteAConcurrentReplacement(t *testing.T) {
	fixture := newPushSafetyFixture(t)
	snapshot, err := fixture.repository.snapshot(
		context.Background(),
		fixture.remote,
		"master",
		"feature",
	)
	if err != nil {
		t.Fatalf("snapshot() error = %v", err)
	}
	target := snapshot.privateRefs[0]
	runTestGit(t, fixture.work, "update-ref", target, fixture.localHead)
	cleanupErr := snapshot.close(context.Background())
	if cleanupErr == nil || !strings.Contains(cleanupErr.Error(), "snapshot-owned") {
		t.Fatalf("snapshot.close() error = %v, want exact-CAS cleanup refusal", cleanupErr)
	}
	if got := runTestGit(t, fixture.work, "rev-parse", "--verify", target); got != fixture.localHead {
		t.Fatalf("replacement private ref = %s, want preserved %s", got, fixture.localHead)
	}
	runTestGit(t, fixture.work, "update-ref", "-d", target, fixture.localHead)
}

func TestSnapshotCleanupRefusesSymbolicReplacement(t *testing.T) {
	fixture := newPushSafetyFixture(t)
	snapshot, err := fixture.repository.snapshot(
		context.Background(),
		fixture.remote,
		"master",
		"feature",
	)
	if err != nil {
		t.Fatalf("snapshot() error = %v", err)
	}
	target := snapshot.privateRefs[0]
	master := runTestGit(t, fixture.work, "rev-parse", "refs/heads/master")
	runTestGit(t, fixture.work, "symbolic-ref", target, "refs/heads/master")
	cleanupErr := snapshot.close(context.Background())
	if cleanupErr == nil || !strings.Contains(cleanupErr.Error(), "became symbolic") {
		t.Fatalf("snapshot.close() error = %v, want symbolic-ref refusal", cleanupErr)
	}
	if got := runTestGit(t, fixture.work, "rev-parse", "refs/heads/master"); got != master {
		t.Fatalf("symbolic cleanup changed master to %s, want %s", got, master)
	}
	if got := runTestGit(t, fixture.work, "symbolic-ref", target); got != "refs/heads/master" {
		t.Fatalf("temporary ref target = %q, want preserved symbolic ref", got)
	}
	runTestGit(t, fixture.work, "symbolic-ref", "--delete", target)
}

func TestPushReconcilesAmbiguousSuccessAndRunsWhenUpToDate(t *testing.T) {
	fixture := newPushSafetyFixture(t)
	runner := &pushInterceptRunner{
		delegate:             &execRunner{},
		failAfterFirstPush:   true,
		capturePushArguments: true,
	}
	fixture.repository.runner = runner

	if err := fixture.repository.push(
		context.Background(),
		fixture.remote,
		fixture.localHead,
		"refs/heads/feature",
		fixture.remoteHead,
		"refs/heads/master",
		fixture.baseOID,
	); err != nil {
		t.Fatalf("push() did not reconcile accepted update: %v", err)
	}
	if runner.pushCount != 1 ||
		!containsString(runner.firstPushArgs, "--atomic") ||
		!containsString(runner.firstPushArgs, "--no-push-option") {
		t.Fatalf("push audit = count %d args %q", runner.pushCount, runner.firstPushArgs)
	}

	runner = &pushInterceptRunner{
		delegate:             &execRunner{},
		capturePushArguments: true,
	}
	fixture.repository.runner = runner
	if err := fixture.repository.push(
		context.Background(),
		fixture.remote,
		fixture.localHead,
		"refs/heads/feature",
		fixture.localHead,
		"refs/heads/master",
		fixture.baseOID,
	); err != nil {
		t.Fatalf("up-to-date exact-lease push error = %v", err)
	}
	if runner.pushCount != 1 {
		t.Fatalf("up-to-date publication ran %d pushes, want one", runner.pushCount)
	}
}

func TestPushRollsBackHeadWhenBaseAdvancesDuringPush(t *testing.T) {
	fixture := newPushSafetyFixture(t)
	runner := &pushInterceptRunner{
		delegate: &execRunner{},
		beforeFirstPush: func() {
			writeTestFile(
				t,
				filepath.Join(fixture.seed, "base-race.txt"),
				"advanced during push\n",
			)
			runTestGit(t, fixture.seed, "add", "base-race.txt")
			runTestGit(t, fixture.seed, "commit", "-m", "advance during push")
			runTestGit(t, fixture.seed, "push", "origin", "master")
		},
	}
	fixture.repository.runner = runner

	err := fixture.repository.push(
		context.Background(),
		fixture.remote,
		fixture.localHead,
		"refs/heads/feature",
		fixture.remoteHead,
		"refs/heads/master",
		fixture.baseOID,
	)
	if err == nil || !strings.Contains(err.Error(), "was restored") {
		t.Fatalf("push() error = %v, want confirmed rollback", err)
	}
	if runner.pushCount != 2 {
		t.Fatalf("push count = %d, want feature update and rollback", runner.pushCount)
	}
	if got := remoteRefOID(t, fixture.work, "origin", "refs/heads/feature"); got != fixture.remoteHead {
		t.Fatalf("remote feature = %s, want restored %s", got, fixture.remoteHead)
	}
	if got := remoteRefOID(t, fixture.work, "origin", "refs/heads/master"); got == fixture.baseOID {
		t.Fatalf("remote base did not advance from %s", fixture.baseOID)
	}
}

type transportEnvironmentAuditRunner struct {
	delegate     commandRunner
	repository   string
	sawTransport bool
}

type constructorIsolationRunner struct {
	delegate        commandRunner
	forgeCalls      int
	gitNetworkCalls int
}

func (r *constructorIsolationRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if request.Name == "gh-test" {
		r.forgeCalls++
		return commandResult{Stdout: `{
            "id":"repository-node",
            "nameWithOwner":"owner/repository",
            "url":"https://github.com/owner/repository"
        }`}, nil
	}
	for _, operation := range []string{"ls-remote", "fetch", "push"} {
		if containsString(request.Args, operation) {
			r.gitNetworkCalls++
			return commandResult{}, fmt.Errorf("unexpected Git network operation")
		}
	}
	return r.delegate.Run(ctx, request)
}

func (r *transportEnvironmentAuditRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if containsString(request.Env, "GIT_TERMINAL_PROMPT=0") {
		r.sawTransport = true
		environment := mergeEnvironment(
			os.Environ(),
			request.Env,
			request.UnsetEnv,
			request.UnsetEnvPrefixes,
		)
		values := environmentMap(environment)
		for _, name := range gitTransportUnsetEnvironment {
			if _, exists := values[name]; exists {
				return commandResult{}, fmt.Errorf(
					"transport environment retained %s",
					name,
				)
			}
		}
		if values["SSH_AUTH_SOCK"] != "preserved-agent-socket" {
			return commandResult{}, fmt.Errorf("SSH_AUTH_SOCK was not preserved")
		}
		if values["GIT_NO_LAZY_FETCH"] != "1" {
			return commandResult{}, fmt.Errorf("lazy fetching was not disabled")
		}
		gitDir := values["GIT_DIR"]
		transportRoot := filepath.Join(r.repository, "out", "repo_delivery")
		relative, relativeErr := filepath.Rel(transportRoot, gitDir)
		info, statErr := os.Lstat(gitDir)
		if relativeErr != nil || !strings.HasPrefix(relative, "transport-") ||
			strings.ContainsRune(relative, filepath.Separator) ||
			request.Dir != gitDir || statErr != nil || !info.IsDir() ||
			info.Mode().Perm() != 0o700 || info.Mode()&os.ModeSymlink != 0 {
			return commandResult{}, fmt.Errorf("transport did not use a private isolated Git directory")
		}
		if _, err := os.Lstat(filepath.Join(gitDir, "config")); !os.IsNotExist(err) {
			return commandResult{}, fmt.Errorf("isolated Git directory unexpectedly has config")
		}
		wantObjectDir, err := filepath.EvalSymlinks(
			filepath.Join(r.repository, ".git", "objects"),
		)
		if err != nil || values["GIT_OBJECT_DIRECTORY"] != wantObjectDir {
			return commandResult{}, fmt.Errorf("transport did not use the canonical main object directory")
		}
		configuration, err := commandConfiguration(values)
		if err != nil {
			return commandResult{}, err
		}
		for key, want := range map[string]string{
			"core.bare":            "true",
			"core.gitProxy":        "none",
			"core.sshCommand":      isolatedSSHCommand,
			"http.followRedirects": "false",
			"http.proxy":           "",
			"http.sslVerify":       "true",
			"protocol.allow":       "never",
			"protocol.ext.allow":   "never",
			"protocol.file.allow":  "always",
			"protocol.ssh.allow":   "always",
			"ssh.variant":          "ssh",
		} {
			if configuration[key] != want {
				return commandResult{}, fmt.Errorf(
					"safe transport config %s = %q, want %q",
					key,
					configuration[key],
					want,
				)
			}
		}
		for key := range configuration {
			if strings.HasPrefix(key, "remote.") || strings.HasPrefix(key, "url.") {
				return commandResult{}, fmt.Errorf("transport retained mutable remote routing config")
			}
		}
		if values["GIT_TERMINAL_PROMPT"] != "0" ||
			values["GIT_CONFIG_GLOBAL"] != os.DevNull ||
			values["GIT_CONFIG_SYSTEM"] != os.DevNull {
			return commandResult{}, fmt.Errorf(
				"transport isolation environment is incomplete",
			)
		}
	}
	return r.delegate.Run(ctx, request)
}

func commandConfiguration(environment map[string]string) (map[string]string, error) {
	count, err := strconv.Atoi(environment["GIT_CONFIG_COUNT"])
	if err != nil || count <= 0 {
		return nil, fmt.Errorf("invalid command-scoped Git configuration count")
	}
	result := make(map[string]string, count)
	for index := 0; index < count; index++ {
		suffix := strconv.Itoa(index)
		key, keyFound := environment["GIT_CONFIG_KEY_"+suffix]
		value, valueFound := environment["GIT_CONFIG_VALUE_"+suffix]
		_, duplicate := result[key]
		if !keyFound || !valueFound || key == "" || duplicate {
			return nil, fmt.Errorf("malformed command-scoped Git configuration")
		}
		result[key] = value
	}
	return result, nil
}

type attachedIndexHookRunner struct {
	delegate commandRunner
	hook     func() error
	fired    bool
}

func (r *attachedIndexHookRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if !r.fired && hasEnvironmentPrefix(request.Env, "GIT_INDEX_FILE=") &&
		containsString(request.Args, "write-tree") {
		r.fired = true
		if err := r.hook(); err != nil {
			return commandResult{}, err
		}
	}
	return r.delegate.Run(ctx, request)
}

type failAttachedIndexInspectionRunner struct {
	delegate commandRunner
	failures int
}

func (r *failAttachedIndexInspectionRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if r.failures == 0 &&
		hasEnvironmentPrefix(request.Env, "GIT_INDEX_FILE=") &&
		containsString(request.Args, "write-tree") {
		r.failures++
		return syntheticCommandFailure(
			request,
			74,
			"synthetic attached index inspection failure",
		)
	}
	return r.delegate.Run(ctx, request)
}

type attachedTransitionHookRunner struct {
	delegate commandRunner
	hook     func() error
	fired    bool
}

func (r *attachedTransitionHookRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if !r.fired && containsSequence(request.Args, "read-tree", "-m", "-u") {
		r.fired = true
		if err := r.hook(); err != nil {
			return commandResult{}, err
		}
	}
	return r.delegate.Run(ctx, request)
}

type failAttachedTransitionVerificationRunner struct {
	delegate     commandRunner
	transitioned bool
	failed       bool
}

func (r *failAttachedTransitionVerificationRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if r.transitioned && !r.failed && containsString(request.Args, "write-tree") {
		r.failed = true
		return syntheticCommandFailure(
			request,
			75,
			"synthetic attached-transition verification failure",
		)
	}
	result, err := r.delegate.Run(ctx, request)
	if err == nil && !r.transitioned &&
		containsSequence(request.Args, "read-tree", "-m", "-u") {
		r.transitioned = true
	}
	return result, err
}

type failFirstWorktreeRemoveRunner struct {
	delegate commandRunner
	attempts int
}

func (r *failFirstWorktreeRemoveRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if containsSequence(request.Args, "worktree", "remove") {
		r.attempts++
		if r.attempts == 1 {
			return syntheticCommandFailure(
				request,
				72,
				"synthetic worktree cleanup failure",
			)
		}
	}
	return r.delegate.Run(ctx, request)
}

type replaceBeforeRefDeleteRunner struct {
	delegate    commandRunner
	targetRef   string
	replacement string
	replaced    bool
}

func (r *replaceBeforeRefDeleteRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if !r.replaced && containsString(request.Args, "update-ref") &&
		containsString(request.Args, "-d") &&
		containsString(request.Args, r.targetRef) {
		r.replaced = true
		command := exec.Command("git", "update-ref", r.targetRef, r.replacement)
		command.Dir = request.Dir
		if output, err := command.CombinedOutput(); err != nil {
			return commandResult{}, fmt.Errorf(
				"install concurrent replacement ref: %v: %s",
				err,
				output,
			)
		}
	}
	return r.delegate.Run(ctx, request)
}

func syntheticCommandFailure(
	request command,
	exitCode int,
	message string,
) (commandResult, error) {
	result := commandResult{Stderr: message, ExitCode: exitCode}
	return result, &commandError{
		Command: request,
		Result:  result,
		Err:     fmt.Errorf("synthetic command failure: %s", message),
	}
}

func containsSequence(values []string, sequence ...string) bool {
	if len(sequence) == 0 || len(sequence) > len(values) {
		return false
	}
	for index := 0; index <= len(values)-len(sequence); index++ {
		matched := true
		for offset, want := range sequence {
			if values[index+offset] != want {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
}

type partialFetchFailureRunner struct {
	delegate       commandRunner
	sawAtomic      bool
	sawNoPrune     bool
	sawNoPruneTags bool
}

type lostAckSnapshotInstallRunner struct {
	delegate commandRunner
	fired    bool
}

func (r *lostAckSnapshotInstallRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if r.fired || !containsString(request.Args, "update-ref") ||
		!containsString(
			request.Args,
			"repo_delivery: install exact fetched snapshot",
		) {
		return r.delegate.Run(ctx, request)
	}
	r.fired = true
	result, err := r.delegate.Run(ctx, request)
	if err != nil {
		return result, err
	}
	result.ExitCode = 73
	result.Stderr = "synthetic snapshot ref lost acknowledgement"
	return result, &commandError{
		Command: request,
		Result:  result,
		Err:     fmt.Errorf("synthetic snapshot ref lost acknowledgement"),
	}
}

func (r *partialFetchFailureRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if !containsString(request.Args, "fetch") {
		return r.delegate.Run(ctx, request)
	}
	r.sawAtomic = containsString(request.Args, "--atomic")
	r.sawNoPrune = containsString(request.Args, "--no-prune")
	r.sawNoPruneTags = containsString(request.Args, "--no-prune-tags")
	result, err := r.delegate.Run(ctx, request)
	if err != nil {
		return result, err
	}
	result.Stderr = "synthetic fetch failure"
	result.ExitCode = 41
	return result, &commandError{
		Command: request,
		Result:  result,
		Err:     fmt.Errorf("synthetic fetch failure"),
	}
}

type pushInterceptRunner struct {
	delegate             commandRunner
	beforeFirstPush      func()
	failAfterFirstPush   bool
	capturePushArguments bool
	pushCount            int
	firstPushArgs        []string
}

func (r *pushInterceptRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if !containsString(request.Args, "push") {
		return r.delegate.Run(ctx, request)
	}
	r.pushCount++
	if r.pushCount == 1 {
		if r.capturePushArguments {
			r.firstPushArgs = append([]string(nil), request.Args...)
		}
		if r.beforeFirstPush != nil {
			r.beforeFirstPush()
		}
	}
	result, err := r.delegate.Run(ctx, request)
	if r.pushCount == 1 && r.failAfterFirstPush && err == nil {
		result.ExitCode = 42
		return result, &commandError{
			Command: request,
			Result:  result,
			Err:     fmt.Errorf("synthetic lost acknowledgement"),
		}
	}
	return result, err
}

type pushSafetyFixture struct {
	remote     string
	seed       string
	work       string
	repository *gitRepository
	baseOID    string
	remoteHead string
	localHead  string
}

func newPushSafetyFixture(t *testing.T) pushSafetyFixture {
	t.Helper()
	root := t.TempDir()
	remote := filepath.Join(root, "remote.git")
	seed := filepath.Join(root, "seed")
	work := filepath.Join(root, "work")
	runTestGit(t, root, "init", "--bare", remote)
	runTestGit(t, root, "init", "--initial-branch=master", seed)
	configureTestRepository(t, seed)
	writeTestFile(t, filepath.Join(seed, ".gitignore"), "/out/\n")
	writeTestFile(t, filepath.Join(seed, "base.txt"), "base\n")
	runTestGit(t, seed, "add", ".gitignore", "base.txt")
	runTestGit(t, seed, "commit", "-m", "base")
	runTestGit(t, seed, "remote", "add", "origin", remote)
	runTestGit(t, seed, "push", "origin", "master")
	runTestGit(t, root, "clone", remote, work)
	configureTestRepository(t, work)
	runTestGit(t, work, "switch", "-c", "feature")
	writeTestFile(t, filepath.Join(work, "feature.txt"), "feature\n")
	runTestGit(t, work, "add", "feature.txt")
	runTestGit(t, work, "commit", "-m", "feature")
	runTestGit(t, work, "push", "origin", "feature")
	repository, err := openGitRepository(
		context.Background(),
		work,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	snapshot, err := repository.snapshot(
		context.Background(),
		remote,
		"master",
		"feature",
	)
	if err != nil {
		t.Fatalf("snapshot() error = %v", err)
	}
	baseOID := snapshot.BaseOID
	remoteHead := snapshot.RemoteHeadOID
	snapshot.close(context.Background())
	writeTestFile(t, filepath.Join(work, "local.txt"), "local\n")
	runTestGit(t, work, "add", "local.txt")
	runTestGit(t, work, "commit", "-m", "local")
	return pushSafetyFixture{
		remote:     remote,
		seed:       seed,
		work:       work,
		repository: repository,
		baseOID:    baseOID,
		remoteHead: remoteHead,
		localHead:  runTestGit(t, work, "rev-parse", "HEAD"),
	}
}

func remoteRefOID(
	t *testing.T,
	directory string,
	remote string,
	ref string,
) string {
	t.Helper()
	value := runTestGit(t, directory, "ls-remote", "--heads", remote, ref)
	if value == "" {
		return ""
	}
	return strings.Fields(value)[0]
}

func assertNoTransportResidue(t *testing.T, repositoryPath string) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(repositoryPath, "out", "repo_delivery"))
	if err != nil {
		t.Fatalf("read delivery scratch directory: %v", err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "transport-") {
			t.Fatalf("isolated transport residue = %q", entry.Name())
		}
	}
}

func environmentMap(environment []string) map[string]string {
	result := make(map[string]string, len(environment))
	for _, value := range environment {
		name, content, found := strings.Cut(value, "=")
		if found {
			result[name] = content
		}
	}
	return result
}

func newRebaseFixture(t *testing.T) (string, string, string) {
	t.Helper()
	repositoryPath := filepath.Join(t.TempDir(), "repository")
	runTestGit(t, t.TempDir(), "init", "--initial-branch=master", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, ".gitignore"), "/out/\n")
	writeTestFile(t, filepath.Join(repositoryPath, "base.txt"), "base\n")
	runTestGit(t, repositoryPath, "add", ".gitignore", "base.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "base")
	runTestGit(t, repositoryPath, "switch", "-c", "feature")
	writeTestFile(t, filepath.Join(repositoryPath, "feature.txt"), "feature\n")
	runTestGit(t, repositoryPath, "add", "feature.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "feature")
	oldFeature := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	runTestGit(t, repositoryPath, "switch", "master")
	writeTestFile(t, filepath.Join(repositoryPath, "advance.txt"), "advance\n")
	runTestGit(t, repositoryPath, "add", "advance.txt")
	runTestGit(t, repositoryPath, "commit", "-m", "advance")
	base := runTestGit(t, repositoryPath, "rev-parse", "HEAD")
	runTestGit(t, repositoryPath, "switch", "feature")
	return repositoryPath, base, oldFeature
}

func testIgnoredPath(t *testing.T, repositoryPath string, name string) string {
	t.Helper()
	directory := filepath.Join(repositoryPath, "out", "test_artifacts")
	if err := os.MkdirAll(directory, 0o700); err != nil {
		t.Fatalf("create ignored test artifact directory: %v", err)
	}
	return filepath.Join(directory, name)
}

func runTestGitAllowFailure(
	directory string,
	args ...string,
) (string, error) {
	command := exec.Command("git", args...)
	command.Dir = directory
	output, err := command.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func configureTestRepository(t *testing.T, directory string) {
	t.Helper()
	runTestGit(t, directory, "config", "user.name", "Delivery Test")
	runTestGit(t, directory, "config", "user.email", "delivery@example.invalid")
	runTestGit(t, directory, "config", "commit.gpgsign", "false")
}
