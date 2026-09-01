package fsstore

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

func newTestStore(t *testing.T) (*Store, string) {
	t.Helper()
	root := t.TempDir()
	runtimeRoot := t.TempDir()
	if err := os.Chmod(runtimeRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_RUNTIME_DIR", runtimeRoot)
	store, err := NewStore(root)
	if err != nil {
		t.Fatal(err)
	}
	instant := time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC)
	var clockMu sync.Mutex
	store.now = func() time.Time {
		clockMu.Lock()
		defer clockMu.Unlock()
		value := instant
		instant = instant.Add(time.Second)
		return value
	}
	store.random = bytes.NewReader(bytes.Repeat([]byte{0xab}, 4096))
	return store, root
}

func initTestGoal(t *testing.T, store *Store, root string, id string) string {
	t.Helper()
	_, err := store.Init(InitOptions{
		GoalsRoot: filepath.Join(root, "out", "task", "goals"),
		Title:     "Goal " + id,
		GoalID:    id,
		Scope:     "workspace",
		Criteria:  []string{"The result is verified."},
	})
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(root, "out", "task", "goals", id)
}

func TestNewStoreWithRuntimeDirRejectsInvalidRuntimeDirectory(t *testing.T) {
	workspaceRoot := t.TempDir()
	insecureRuntime := t.TempDir()
	if err := os.Chmod(insecureRuntime, 0o755); err != nil {
		t.Fatal(err)
	}
	workspaceRuntime := filepath.Join(workspaceRoot, "runtime")
	if err := os.Mkdir(workspaceRuntime, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name       string
		runtimeDir string
		want       string
	}{
		{name: "empty", want: "required"},
		{name: "relative", runtimeDir: "runtime", want: "absolute"},
		{name: "insecure mode", runtimeDir: insecureRuntime, want: "mode 0700"},
		{name: "inside workspace", runtimeDir: workspaceRuntime, want: "outside the workspace"},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewStoreWithRuntimeDir(workspaceRoot, test.runtimeDir)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("NewStoreWithRuntimeDir() error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestNewStoreRejectsFinalLockRootInsideWorkspace(t *testing.T) {
	runtimeRoot := t.TempDir()
	if err := os.Chmod(runtimeRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	workspaceRoot := filepath.Join(runtimeRoot, "alwaldend", "goal")
	if err := os.MkdirAll(workspaceRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := NewStoreWithRuntimeDir(workspaceRoot, runtimeRoot)
	if err == nil || !strings.Contains(err.Error(), "locks must be outside") {
		t.Fatalf("NewStoreWithRuntimeDir() error = %v, want final lock-root rejection", err)
	}
	if pathExists(filepath.Join(workspaceRoot, "locks")) {
		t.Fatal("constructor created the rejected lock directory inside the workspace")
	}
}

func TestTwoGoalIsolationAndBoundedCatalog(t *testing.T) {
	store, root := newTestStore(t)
	first := initTestGoal(t, store, root, "first-goal")
	second := initTestGoal(t, store, root, "second-goal")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: first, ExpectedResourceVersion: "1", AttemptID: "attempt-a",
		WorkType: "investigation",
	}); err != nil {
		t.Fatal(err)
	}
	secondGoal, err := store.readGoalManifest(second)
	if err != nil {
		t.Fatal(err)
	}
	if secondGoal.Metadata.ResourceVersion != "1" || secondGoal.Status.ActiveAttemptID != "" {
		t.Fatalf("second goal changed: %+v", secondGoal)
	}
	list, err := store.List(filepath.Join(root, "out", "task", "goals"), 1)
	if err != nil {
		t.Fatal(err)
	}
	if list.Returned != 1 || list.Total != 2 || !list.Truncated {
		t.Fatalf("unexpected bounded list: %+v", list)
	}
}

func TestGoalWithoutTrackedAttemptsDirectoryValidatesAndStartsAttempt(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "checked-out-goal")
	plan := filepath.Join(root, "out", "task", "plan.md")
	writeTestFile(t, plan, "# Checked-out plan\n")
	attemptsDir := filepath.Join(goalDir, "attempts")
	if err := os.Remove(attemptsDir); err != nil {
		t.Fatal(err)
	}
	if err := store.ValidateGoal(goalDir); err != nil {
		t.Fatalf("ValidateGoal() without empty attempts directory: %v", err)
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "attempt-1",
		PlanFile: filepath.Join("out", "task", "plan.md"),
	}); err != nil {
		t.Fatalf("Checkpoint() did not recreate attempts directory: %v", err)
	}
	if !pathExists(filepath.Join(attemptsDir, "attempt-1", "attempt.yaml")) {
		t.Fatal("Checkpoint() did not publish the first attempt")
	}
	content, err := os.ReadFile(filepath.Join(attemptsDir, "attempt-1", "plan.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "# Checked-out plan\n" {
		t.Fatalf("unexpected workspace-relative plan: %q", content)
	}
}

func TestFreshStoreResumesExplicitSessionWithPortableReference(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "resumable-goal")
	sessionRoot := filepath.Join(root, "out", "task", "goal-sessions")
	binding, err := store.Attach(AttachOptions{
		SessionRoot: sessionRoot,
		SessionID:   "session-1",
		GoalDir:     goalDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	goalRef := binding.Metadata.Annotations[localGoalReferenceAnnotation]
	if filepath.IsAbs(goalRef) || strings.Contains(goalRef, root) {
		t.Fatalf("session leaked an absolute path: %q", goalRef)
	}
	fresh, err := NewStore(root)
	if err != nil {
		t.Fatal(err)
	}
	view, err := fresh.ShowSession(sessionRoot, "session-1", 5)
	if err != nil {
		t.Fatal(err)
	}
	if view.Goal.Metadata.Name != "resumable-goal" || view.Session == nil || view.SessionStale {
		t.Fatalf("unexpected resumed view: %+v", view)
	}
	if pathExists(filepath.Join(sessionRoot, ".locks")) {
		t.Fatal("session lock files must remain outside the workspace")
	}
	if _, err := store.Attach(AttachOptions{
		SessionRoot: filepath.Join(root, "out", "goal-sessions"),
		SessionID:   "bad-session",
		GoalDir:     goalDir,
	}); err == nil {
		t.Fatal("expected non-task-specific session root to be rejected")
	}
}

func TestConcurrentCheckpointRejectsOneStaleWriter(t *testing.T) {
	first, root := newTestStore(t)
	goalDir := initTestGoal(t, first, root, "concurrent-goal")
	second, err := NewStore(root)
	if err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	errorsByWriter := make(chan error, 2)
	for index, store := range []*Store{first, second} {
		store := store
		id := fmt.Sprintf("attempt-%c", 'a'+index)
		go func() {
			<-start
			_, err := store.Checkpoint(CheckpointOptions{
				GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: id,
			})
			errorsByWriter <- err
		}()
	}
	close(start)
	successes := 0
	failures := 0
	for range 2 {
		if err := <-errorsByWriter; err == nil {
			successes++
		} else if strings.Contains(err.Error(), "stale resourceVersion") {
			failures++
		} else {
			t.Fatalf("unexpected concurrent error: %v", err)
		}
	}
	if successes != 1 || failures != 1 {
		t.Fatalf("successes=%d failures=%d", successes, failures)
	}
	if _, err := first.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", Execution: "paused",
	}); err == nil || !strings.Contains(err.Error(), "stale resourceVersion") {
		t.Fatalf("expected stale writer rejection, got %v", err)
	}
}

func TestPerGoalLockSerializesTwoStoresThroughWorkspaceAlias(t *testing.T) {
	first, root := newTestStore(t)
	goalDir := initTestGoal(t, first, root, "locked-goal")
	alias := filepath.Join(root, "workspace-alias")
	if err := os.Symlink(root, alias); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	second, err := NewStore(alias)
	if err != nil {
		t.Fatal(err)
	}
	aliasGoalDir := filepath.Join(alias, "out", "task", "goals", "locked-goal")
	goalPath := filepath.Join(goalDir, "goal.yaml")
	firstEntered := make(chan struct{})
	releaseFirst := make(chan struct{})
	first.beforeRename = func(target string) error {
		if target == goalPath {
			close(firstEntered)
			<-releaseFirst
		}
		return nil
	}
	secondEntered := make(chan struct{}, 1)
	second.beforeRename = func(target string) error {
		if target == goalPath {
			secondEntered <- struct{}{}
		}
		return nil
	}
	firstResult := make(chan error, 1)
	go func() {
		_, err := first.Checkpoint(CheckpointOptions{
			GoalDir: goalDir, ExpectedResourceVersion: "1", Execution: "paused",
		})
		firstResult <- err
	}()
	<-firstEntered
	secondResult := make(chan error, 1)
	go func() {
		_, err := second.Checkpoint(CheckpointOptions{
			GoalDir: aliasGoalDir, ExpectedResourceVersion: "1", Execution: "waiting",
		})
		secondResult <- err
	}()

	select {
	case <-secondEntered:
		close(releaseFirst)
		<-firstResult
		<-secondResult
		t.Fatal("second store reached the rename while the goal lock was held")
	case err := <-secondResult:
		close(releaseFirst)
		<-firstResult
		t.Fatalf("second store returned before the goal lock was released: %v", err)
	case <-time.After(100 * time.Millisecond):
	}
	close(releaseFirst)
	if err := <-firstResult; err != nil {
		t.Fatal(err)
	}
	if err := <-secondResult; err == nil ||
		!strings.Contains(err.Error(), "stale resourceVersion") {
		t.Fatalf("expected serialized stale writer rejection, got %v", err)
	}
}

func TestPerGoalLockKeysDifferentWorkspacesIndependently(t *testing.T) {
	first, root := newTestStore(t)
	firstGoal := initTestGoal(t, first, root, "same-relative-goal")
	secondRoot := filepath.Join(root, "second-workspace")
	if err := os.MkdirAll(secondRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	second, err := NewStore(secondRoot)
	if err != nil {
		t.Fatal(err)
	}
	secondGoal := initTestGoal(t, second, secondRoot, "same-relative-goal")
	firstGoalPath := filepath.Join(firstGoal, "goal.yaml")
	firstEntered := make(chan struct{})
	releaseFirst := make(chan struct{})
	first.beforeRename = func(target string) error {
		if target == firstGoalPath {
			close(firstEntered)
			<-releaseFirst
		}
		return nil
	}
	firstResult := make(chan error, 1)
	go func() {
		_, err := first.Checkpoint(CheckpointOptions{
			GoalDir: firstGoal, ExpectedResourceVersion: "1", Execution: "paused",
		})
		firstResult <- err
	}()
	<-firstEntered
	secondResult := make(chan error, 1)
	go func() {
		_, err := second.Checkpoint(CheckpointOptions{
			GoalDir: secondGoal, ExpectedResourceVersion: "1", Execution: "paused",
		})
		secondResult <- err
	}()

	select {
	case err := <-secondResult:
		close(releaseFirst)
		if err != nil {
			<-firstResult
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		close(releaseFirst)
		<-firstResult
		t.Fatal("equal relative goal paths in different workspaces shared a lock")
	}
	if err := <-firstResult; err != nil {
		t.Fatal(err)
	}
}

func TestPerGoalLockDoesNotSerializeSiblingGoals(t *testing.T) {
	first, root := newTestStore(t)
	firstGoal := initTestGoal(t, first, root, "first-lock-goal")
	secondGoal := initTestGoal(t, first, root, "second-lock-goal")
	second, err := NewStore(root)
	if err != nil {
		t.Fatal(err)
	}
	firstGoalPath := filepath.Join(firstGoal, "goal.yaml")
	firstEntered := make(chan struct{})
	releaseFirst := make(chan struct{})
	first.beforeRename = func(target string) error {
		if target == firstGoalPath {
			close(firstEntered)
			<-releaseFirst
		}
		return nil
	}
	firstResult := make(chan error, 1)
	go func() {
		_, err := first.Checkpoint(CheckpointOptions{
			GoalDir: firstGoal, ExpectedResourceVersion: "1", Execution: "paused",
		})
		firstResult <- err
	}()
	<-firstEntered
	secondResult := make(chan error, 1)
	go func() {
		_, err := second.Checkpoint(CheckpointOptions{
			GoalDir: secondGoal, ExpectedResourceVersion: "1", Execution: "paused",
		})
		secondResult <- err
	}()

	select {
	case err := <-secondResult:
		close(releaseFirst)
		if err != nil {
			<-firstResult
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		close(releaseFirst)
		<-firstResult
		t.Fatal("a lock for one goal blocked a sibling goal")
	}
	if err := <-firstResult; err != nil {
		t.Fatal(err)
	}
}

func TestGoalLockPersistsOverwritesPIDAndReleases(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "pid-lock-goal")
	lockDir := filepath.Join(store.runtimeRoot, "alwaldend", "goal", "locks")
	entries, err := os.ReadDir(lockDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].IsDir() {
		t.Fatalf("unexpected per-goal lock entries: %+v", entries)
	}
	lockPath := filepath.Join(lockDir, entries[0].Name())
	if err := os.WriteFile(lockPath, []byte("stale-holder\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(lockPath)
	if err != nil {
		t.Fatal(err)
	}
	lock, err := store.acquireGoalLock(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(lockPath)
	if err != nil {
		_ = lock.release()
		t.Fatal(err)
	}
	if want := fmt.Sprintf("%d\n", os.Getpid()); string(content) != want {
		_ = lock.release()
		t.Fatalf("lock holder content = %q, want %q", content, want)
	}
	if err := lock.release(); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(lockPath)
	if err != nil {
		t.Fatalf("released lock file did not persist: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("lock permissions = %o, want 600", got)
	}
	beforeStat := before.Sys().(*syscall.Stat_t)
	afterStat := info.Sys().(*syscall.Stat_t)
	if beforeStat.Ino != afterStat.Ino {
		t.Fatal("clean release replaced the persistent lock inode")
	}
	content, err = os.ReadFile(lockPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) != 0 {
		t.Fatalf("released lock metadata = %q, want empty", content)
	}

	acquiredAgain := make(chan error, 1)
	go func() {
		next, err := store.acquireGoalLock(goalDir)
		if err == nil {
			err = next.release()
		}
		acquiredAgain <- err
	}()
	select {
	case err := <-acquiredAgain:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("released goal lock could not be acquired again")
	}
}

func TestGoalLockRejectsUnsafeInodes(t *testing.T) {
	for _, test := range []struct {
		name string
		want string
		make func(string, string) error
	}{
		{
			name: "symlink",
			want: "open goal lock",
			make: func(target string, lockPath string) error {
				return os.Symlink(target, lockPath)
			},
		},
		{
			name: "hard link",
			want: "exactly one link",
			make: func(target string, lockPath string) error {
				return os.Link(target, lockPath)
			},
		},
		{
			name: "fifo",
			want: "regular file",
			make: func(_ string, lockPath string) error {
				return syscall.Mkfifo(lockPath, 0o600)
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			store, root := newTestStore(t)
			goalDir := initTestGoal(t, store, root, "unsafe-lock-goal")
			lockPath, err := store.pathLockPath(goalDir)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(lockPath); err != nil {
				t.Fatal(err)
			}
			target := filepath.Join(root, "must-not-change")
			original := []byte("unrelated same-user data\n")
			if err := os.WriteFile(target, original, 0o600); err != nil {
				t.Fatal(err)
			}
			if err := test.make(target, lockPath); err != nil {
				t.Fatal(err)
			}
			lock, err := store.acquireGoalLock(goalDir)
			if err == nil {
				_ = lock.release()
				t.Fatal("acquired an unsafe lock inode")
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Fatalf("acquireGoalLock() error = %v, want %q", err, test.want)
			}
			content, err := os.ReadFile(target)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(content, original) {
				t.Fatalf("unsafe lock target changed to %q", content)
			}
		})
	}
}

func TestGoalLockCrashLeavesDiagnosticPID(t *testing.T) {
	if os.Getenv("GOAL_LOCK_CRASH_HELPER") == "1" {
		store, err := NewStoreWithRuntimeDir(
			os.Getenv("GOAL_LOCK_CRASH_WORKSPACE"),
			os.Getenv("GOAL_LOCK_CRASH_RUNTIME"),
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		lock, err := store.acquireGoalLock(os.Getenv("GOAL_LOCK_CRASH_PATH"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		_ = lock
		os.Exit(0)
	}

	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "crashed-lock-goal")
	lockPath, err := store.pathLockPath(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	command := exec.Command(os.Args[0], "-test.run=^TestGoalLockCrashLeavesDiagnosticPID$")
	command.Env = append(
		os.Environ(),
		"GOAL_LOCK_CRASH_HELPER=1",
		"GOAL_LOCK_CRASH_WORKSPACE="+root,
		"GOAL_LOCK_CRASH_RUNTIME="+store.runtimeRoot,
		"GOAL_LOCK_CRASH_PATH="+goalDir,
	)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("crash helper failed: %v\n%s", err, output)
	}
	content, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatal(err)
	}
	if want := fmt.Sprintf("%d\n", command.ProcessState.Pid()); string(content) != want {
		t.Fatalf("crashed holder metadata = %q, want %q", content, want)
	}

	lock, err := store.acquireGoalLock(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := lock.release(); err != nil {
		t.Fatal(err)
	}
	content, err = os.ReadFile(lockPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) != 0 {
		t.Fatalf("clean successor release left metadata %q", content)
	}
}

func TestAtomicWriteFailureLeavesGoalManifestUnchanged(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "atomic-goal")
	path := filepath.Join(goalDir, "goal.yaml")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	store.beforeRename = func(target string) error {
		if target == path {
			return errors.New("power loss")
		}
		return nil
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", Execution: "paused",
	}); err == nil {
		t.Fatal("expected injected failure")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("goal manifest changed before atomic rename")
	}
}

func TestInvalidProspectiveTransitionLeavesNoAttempt(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "preflight-goal")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "attempt-1",
		Execution: "paused",
	}); err == nil {
		t.Fatal("expected open attempt plus paused execution to fail preflight")
	}
	if pathExists(filepath.Join(goalDir, "attempts", "attempt-1")) {
		t.Fatal("invalid prospective transition left an attempt")
	}
	if err := store.ValidateGoal(goalDir); err != nil {
		t.Fatal(err)
	}
}

func TestEvidenceFilesMustBeMarkdown(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "evidence-format-goal")
	invalid := filepath.Join(root, "out", "task", "evidence.txt")
	writeTestFile(t, invalid, "plain text\n")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "bad-attempt",
		EvidenceFiles: []string{invalid},
	}); err == nil {
		t.Fatal("expected non-Markdown evidence to be rejected")
	}
	if pathExists(filepath.Join(goalDir, "attempts", "bad-attempt")) {
		t.Fatal("rejected evidence left an attempt")
	}

	valid := filepath.Join(root, "out", "task", "evidence.md")
	writeTestFile(t, valid, "# Evidence\n\nVerified.\n")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "good-attempt",
		EvidenceFiles: []string{valid},
	}); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(
		goalDir,
		"attempts",
		"good-attempt",
		"evidence",
		"evidence.md",
	))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "# Evidence\n\nVerified.\n" {
		t.Fatalf("unexpected copied evidence: %q", content)
	}
}

func TestClosedAttemptIsImmutableAndLateCloseAfterPauseFails(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "closed-attempt-goal")
	plan := filepath.Join(root, "out", "task", "plan.md")
	writeTestFile(t, plan, "# Plan\n")
	review := writeRefineReview(t, root, "close-review.yaml")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "attempt-1",
		PlanFile: plan, WorkType: "change",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "2", AttemptID: "attempt-1",
		CloseAttempt: true, ReviewFile: review,
	}); err != nil {
		t.Fatal(err)
	}
	result := filepath.Join(root, "out", "task", "late-result.md")
	writeTestFile(t, result, "# Late result\n")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "3", AttemptID: "attempt-1",
		ResultFile: result,
	}); err == nil || !strings.Contains(err.Error(), "closed and immutable") {
		t.Fatalf("expected closed-attempt publication to fail, got %v", err)
	}
	attemptBytes, err := os.ReadFile(filepath.Join(goalDir, "attempts", "attempt-1", "result.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(attemptBytes), "Late result") {
		t.Fatal("closed attempt result was modified")
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "3", AttemptID: "attempt-2",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "4", AttemptID: "attempt-2",
		CloseAttempt: true, ReviewFile: review, Execution: "paused",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "5", AttemptID: "attempt-2",
		CloseAttempt: true, ReviewFile: review,
	}); err == nil || !strings.Contains(err.Error(), "requires active execution") {
		t.Fatalf("expected late close after pause to fail, got %v", err)
	}
}

func TestAchievedBindsAcceptedAttemptAndDetectsArtifactMutation(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "accepted-goal")
	result := filepath.Join(root, "out", "task", "accepted-result.md")
	review := filepath.Join(root, "out", "task", "accepted-review.yaml")
	writeTestFile(t, result, "# Result\n\nVerified.\n")
	writeTestFile(t, review, `decision: accept
criteria:
  - criterionID: criterion-001
    criterionRevision: 1
    verdict: pass
    evidenceRefs:
      - result.md
`)
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "attempt-1",
		ResultFile: result, ReviewFile: review, CloseAttempt: true, Outcome: "achieved",
	}); err != nil {
		t.Fatal(err)
	}
	goal, _, _, err := store.loadAndValidate(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if goal.Status.AcceptedAttemptID != "attempt-1" ||
		!validDigest(goal.Status.AcceptedResultDigest) {
		t.Fatalf("accepted result pointer was not recorded: %+v", goal.Status)
	}
	if err := os.WriteFile(
		filepath.Join(goalDir, "attempts", "attempt-1", "result.md"),
		[]byte("tampered\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}
	if err := store.ValidateGoal(goalDir); err == nil ||
		!strings.Contains(err.Error(), "artifact digest") {
		t.Fatalf("closed artifact mutation was not detected: %v", err)
	}
}

func TestCriteriaUpdateVersionsItems(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "criteria-goal")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", Execution: "paused",
	}); err != nil {
		t.Fatal(err)
	}
	desired := filepath.Join(root, "out", "task", "criteria-input.yaml")
	writeTestFile(t, desired, `items:
  - criterionID: criterion-001
    statement: The revised result is verified.
    evidenceMethod: Run the deterministic test.
  - criterionID: criterion-002
    required: false
    statement: The projection is readable.
    evidenceMethod: Inspect README.md.
`)
	result, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "2", CriteriaFile: desired,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ResourceVersion != "3" {
		t.Fatalf("unexpected resource version: %+v", result)
	}
	goal, criteria, _, err := store.loadAndValidate(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if goal.Status.CriteriaRevision != 2 || goal.Status.LifecycleGeneration != 3 ||
		criteria.Spec.Items[0].Revision != 2 || criteria.Spec.Items[1].Revision != 1 ||
		criteria.Spec.Items[1].Required {
		t.Fatalf("criteria versions are wrong: goal=%+v criteria=%+v", goal, criteria)
	}
}

func TestUnversionedMigrationIsNonDestructiveAndIdempotent(t *testing.T) {
	store, root := newTestStore(t)
	source := filepath.Join(root, "out", "task", "legacy", "legacy-goal")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := "# Legacy goal\n\n## Acceptance criteria\n\n- The legacy evidence remains inspectable.\n"
	writeTestFile(t, filepath.Join(source, "README.md"), legacy)
	writeTestFile(t, filepath.Join(source, "FAILURES.md"), "# Failures\n\nNone.\n")
	sourceFiles, err := inspectUnversionedRecord(source)
	if err != nil {
		t.Fatal(err)
	}
	sourceDigest := digestLegacyFiles(sourceFiles)
	destinationRoot := filepath.Join(root, "out", "task", "imported", "goals")
	options := MigrateOptions{
		SourceGoalDir:        source,
		DestinationGoalsRoot: destinationRoot,
		Scope:                "workspace",
	}
	first, err := store.Migrate(options)
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Migrate(options)
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatalf("migration is not idempotent: %+v %+v", first, second)
	}
	target := filepath.Join(destinationRoot, "legacy-goal")
	goal, criteria, _, err := store.loadAndValidate(target)
	if err != nil {
		t.Fatal(err)
	}
	if goal.Status.Migration.SourceFormat != "unversioned" || len(criteria.Spec.Items) != 1 {
		t.Fatalf("migration metadata missing: %+v %+v", goal, criteria)
	}
	if goal.Status.Migration.SourcePath != "out/task/legacy/legacy-goal" {
		t.Fatalf(
			"migration source path = %q, want workspace-relative reference",
			goal.Status.Migration.SourcePath,
		)
	}
	if goal.Status.Migration.MappingVersion != "v1" ||
		goal.Status.Migration.ExtractionMode != "extracted" {
		t.Fatalf(
			"migration mapping provenance missing: %+v",
			goal.Status.Migration,
		)
	}
	if goal.Status.Migration.SourceDigest != sourceDigest {
		t.Fatalf(
			"migration source digest = %q, want %q",
			goal.Status.Migration.SourceDigest,
			sourceDigest,
		)
	}
	plan, err := os.ReadFile(filepath.Join(target, "attempts", "imported-unversioned", "plan.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(plan) != legacy {
		t.Fatal("legacy README was not preserved verbatim")
	}
	latestSourceFiles, err := inspectUnversionedRecord(source)
	if err != nil {
		t.Fatal(err)
	}
	if digestLegacyFiles(latestSourceFiles) != sourceDigest ||
		pathExists(filepath.Join(source, "goal.yaml")) {
		t.Fatal("migration changed the legacy source directory")
	}
}

func TestMutationsLeaveNoLockWALOrJournalArtifactsInCatalog(t *testing.T) {
	store, root := newTestStore(t)
	goalsRoot := filepath.Join(root, "out", "task", "goals")
	goalDir := initTestGoal(t, store, root, "plain-goal")
	desired := filepath.Join(root, "out", "task", "criteria-input.yaml")
	writeTestFile(t, desired, `items:
  - criterionID: criterion-001
    statement: The revised result is verified.
    evidenceMethod: Run the deterministic test.
`)
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", Execution: "paused",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "2", CriteriaFile: desired,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "3", Execution: "active",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "4", AttemptID: "attempt-1",
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.Render(goalDir, "5", 3); err != nil {
		t.Fatal(err)
	}

	legacyDir := filepath.Join(root, "out", "task", "legacy", "plain-legacy")
	writeTestFile(t, filepath.Join(legacyDir, "README.md"), "# Legacy\n")
	if _, err := store.Migrate(MigrateOptions{
		SourceGoalDir:        legacyDir,
		DestinationGoalsRoot: goalsRoot,
		Scope:                "workspace",
	}); err != nil {
		t.Fatal(err)
	}

	if err := filepath.WalkDir(goalsRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := strings.ToLower(entry.Name())
		if name == ".locks" || strings.HasSuffix(name, ".lock") ||
			strings.Contains(name, "wal") || strings.Contains(name, "journal") ||
			strings.Contains(name, "transaction") {
			relative, relErr := filepath.Rel(goalsRoot, path)
			if relErr != nil {
				return relErr
			}
			return fmt.Errorf("unexpected store coordination artifact %q", relative)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestPromotionPreservesIdentityProvenanceAndPortablePaths(t *testing.T) {
	store, root := newTestStore(t)
	source := initTestGoal(t, store, root, "promoted-goal")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: source, ExpectedResourceVersion: "1", Execution: "paused",
	}); err != nil {
		t.Fatal(err)
	}
	destinationRoot := filepath.Join(root, "project", "goals")
	first, err := store.Promote(PromoteOptions{
		GoalDir: source, DestinationGoalsRoot: destinationRoot,
		ExpectedResourceVersion: "2",
	})
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Promote(PromoteOptions{
		GoalDir: source, DestinationGoalsRoot: destinationRoot,
		ExpectedResourceVersion: "2",
	})
	if err != nil {
		t.Fatal(err)
	}
	if first != second || first.GoalID != "promoted-goal" {
		t.Fatalf("promotion not stable: %+v %+v", first, second)
	}
	target := filepath.Join(destinationRoot, "promoted-goal")
	goal, _, _, err := store.loadAndValidate(target)
	if err != nil {
		t.Fatal(err)
	}
	ownerReference := goal.Metadata.Annotations[localOwnerRootAnnotation]
	if goal.Spec.Scope != "project" || goal.Metadata.Generation != 2 ||
		goal.Status.Promotion.SourceGeneration != 1 ||
		!validDigest(goal.Status.Promotion.SourceStateDigest) ||
		ownerReference == "" || filepath.IsAbs(ownerReference) {
		t.Fatalf("bad promoted goal: %+v", goal)
	}
	for _, file := range []string{"goal.yaml", "criteria.yaml"} {
		content, err := os.ReadFile(filepath.Join(target, file))
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(content, []byte(root)) {
			t.Fatalf("%s leaks workspace path", file)
		}
		if file == "goal.yaml" && bytes.Contains(content, []byte("ownerRoot:")) {
			t.Fatal("goal spec still serializes backend-local ownerRoot")
		}
	}
	sourceGoal, err := store.readGoalManifest(source)
	if err != nil {
		t.Fatal(err)
	}
	if sourceGoal.Spec.Scope != "workspace" || sourceGoal.Metadata.ResourceVersion != "2" {
		t.Fatal("promotion mutated its source")
	}
}

func TestRenderAndShowRemainBounded(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "bounded-goal")
	resourceVersion := 1
	review := writeRefineReview(t, root, "bounded-review.yaml")
	for index := range 8 {
		id := fmt.Sprintf("attempt-%02d", index)
		if _, err := store.Checkpoint(CheckpointOptions{
			GoalDir: goalDir, ExpectedResourceVersion: fmt.Sprint(resourceVersion),
			AttemptID: id, WorkType: "validation",
		}); err != nil {
			t.Fatal(err)
		}
		resourceVersion++
		if _, err := store.Checkpoint(CheckpointOptions{
			GoalDir: goalDir, ExpectedResourceVersion: fmt.Sprint(resourceVersion),
			AttemptID: id, CloseAttempt: true, ReviewFile: review,
		}); err != nil {
			t.Fatal(err)
		}
		resourceVersion++
	}
	if err := store.Render(goalDir, fmt.Sprint(resourceVersion), 3); err != nil {
		t.Fatal(err)
	}
	readme, err := os.ReadFile(filepath.Join(goalDir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if links := strings.Count(string(readme), "— `validation`,"); links != 3 {
		t.Fatalf("rendered %d attempt links, want 3", links)
	}
	if !strings.Contains(string(readme), "older attempts omitted") {
		t.Fatal("bounded render does not report omission")
	}
	view, err := store.ShowGoal(goalDir, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(view.Attempts) != 3 || view.Total != 8 || !view.Truncated {
		t.Fatalf("show is not bounded: %+v", view)
	}
}

func TestStrictDecodeAndKubernetesCompatibleNames(t *testing.T) {
	store, root := newTestStore(t)
	if _, err := store.Init(InitOptions{
		GoalsRoot: filepath.Join(root, "out", "task", "goals"),
		Title:     "Bad ID", GoalID: "bad_id",
	}); err == nil {
		t.Fatal("expected underscore in metadata.name to be rejected")
	}
	if _, err := store.Init(InitOptions{
		GoalsRoot: filepath.Join(root, "out", "task", "goals"),
		Title:     "Bad DNS labels", GoalID: "bad-.id",
	}); err == nil {
		t.Fatal("expected an invalid DNS subdomain to be rejected")
	}
	goalDir := initTestGoal(t, store, root, "strict-goal")
	goalPath := filepath.Join(goalDir, "goal.yaml")
	content, err := os.ReadFile(goalPath)
	if err != nil {
		t.Fatal(err)
	}
	content = append(content, []byte("unknownField: true\n")...)
	writeTestFile(t, goalPath, string(content))
	if err := store.ValidateGoal(goalDir); err == nil || !strings.Contains(err.Error(), "strict") {
		t.Fatalf("expected strict unknown-field failure, got %v", err)
	}
}

func TestPersistedMetadataAndLocalOwnerRootAreStrict(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "strict-owner-goal")
	goal, err := store.readGoalManifest(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, ownerRoot := range []string{".", "projects/goal"} {
		candidate := goal
		candidate.Metadata.Annotations = map[string]string{
			localOwnerRootAnnotation: ownerRoot,
		}
		if err := candidate.validate(); err != nil {
			t.Errorf("valid local owner root %q: %v", ownerRoot, err)
		}
	}
	for _, ownerRoot := range []string{
		"",
		"/absolute",
		"../parent",
		"a/../b",
		"trailing/",
		`windows\path`,
	} {
		candidate := goal
		candidate.Metadata.Annotations = map[string]string{
			localOwnerRootAnnotation: ownerRoot,
		}
		if err := candidate.validate(); err == nil {
			t.Errorf("invalid local owner root %q was accepted", ownerRoot)
		}
	}
	missing := goal
	missing.Metadata.Annotations = nil
	if err := missing.validate(); err == nil {
		t.Fatal("missing local owner root annotation was accepted")
	}
	missingVersion := goal
	missingVersion.Metadata.ResourceVersion = ""
	if err := missingVersion.validate(); err == nil {
		t.Fatal("persisted goal without resourceVersion was accepted")
	}
}

func TestSymlinkedPathCannotEscapeWorkspace(t *testing.T) {
	store, root := newTestStore(t)
	external := t.TempDir()
	link := filepath.Join(root, "out", "task", "escape")
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := store.Init(InitOptions{
		GoalsRoot: filepath.Join(link, "goals"),
		GoalID:    "escaped-goal",
		Title:     "Escaped goal",
	}); err == nil || !strings.Contains(err.Error(), "escapes workspace") {
		t.Fatalf("expected symlink escape rejection, got %v", err)
	}
}

func writeTestFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeRefineReview(t *testing.T, root string, name string) string {
	t.Helper()
	path := filepath.Join(root, "out", "task", name)
	writeTestFile(t, path, `decision: refine
criteria:
  - criterionID: criterion-001
    criterionRevision: 1
    verdict: unverified
    evidenceRefs: []
`)
	return path
}
