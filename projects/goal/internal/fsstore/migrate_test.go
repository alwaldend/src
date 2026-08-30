package fsstore

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func pauseGoalForPromotion(
	t *testing.T,
	store *Store,
	goalDir string,
	expectedResourceVersion string,
) string {
	t.Helper()
	reference, err := store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: expectedResourceVersion,
		Execution:               "paused",
	})
	if err != nil {
		t.Fatal(err)
	}
	return reference.ResourceVersion
}

func TestPromotionRejectsSymlinkAliasToSourceWithoutDeadlock(t *testing.T) {
	store, root := newTestStore(t)
	source := initTestGoal(t, store, root, "alias-promotion")
	resourceVersion := pauseGoalForPromotion(t, store, source, "1")
	destinationRoot := filepath.Join(root, "project", "goals")
	if err := os.MkdirAll(destinationRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(source, filepath.Join(destinationRoot, filepath.Base(source))); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	result := make(chan error, 1)
	go func() {
		_, err := store.Promote(PromoteOptions{
			GoalDir:                 source,
			DestinationGoalsRoot:    destinationRoot,
			ExpectedResourceVersion: resourceVersion,
			OwnerRoot:               filepath.Join(root, "project"),
		})
		result <- err
	}()
	select {
	case err := <-result:
		if err == nil || !strings.Contains(err.Error(), "same") {
			t.Fatalf("Promote() error = %v, want canonical identity rejection", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Promote() deadlocked while acquiring the same canonical lock twice")
	}
}

func TestPromotionRejectsOverlappingGoalPathsBeforeCreatingDestination(t *testing.T) {
	t.Run("destination nested under source", func(t *testing.T) {
		store, root := newTestStore(t)
		source := initTestGoal(t, store, root, "nested-destination")
		resourceVersion := pauseGoalForPromotion(t, store, source, "1")
		destinationRoot := filepath.Join(source, "goals")
		_, err := store.Promote(PromoteOptions{
			GoalDir:                 source,
			DestinationGoalsRoot:    destinationRoot,
			ExpectedResourceVersion: resourceVersion,
			OwnerRoot:               source,
		})
		if err == nil || !strings.Contains(err.Error(), "overlap") {
			t.Fatalf("Promote() error = %v, want overlap rejection", err)
		}
		if pathExists(destinationRoot) {
			t.Fatal("promotion created its destination before rejecting overlap")
		}
	})

	t.Run("destination is ancestor of source", func(t *testing.T) {
		store, root := newTestStore(t)
		ownerRoot := filepath.Join(root, "project")
		destinationRoot := filepath.Join(ownerRoot, "goals")
		nestedGoalsRoot := filepath.Join(
			destinationRoot,
			"ancestor-destination",
			"child",
		)
		if _, err := store.Init(InitOptions{
			GoalsRoot: nestedGoalsRoot,
			Title:     "Nested source",
			GoalID:    "ancestor-destination",
			Scope:     "workspace",
			Criteria:  []string{"The result is verified."},
		}); err != nil {
			t.Fatal(err)
		}
		source := filepath.Join(nestedGoalsRoot, "ancestor-destination")
		resourceVersion := pauseGoalForPromotion(t, store, source, "1")
		_, err := store.Promote(PromoteOptions{
			GoalDir:                 source,
			DestinationGoalsRoot:    destinationRoot,
			ExpectedResourceVersion: resourceVersion,
			OwnerRoot:               ownerRoot,
		})
		if err == nil || !strings.Contains(err.Error(), "overlap") {
			t.Fatalf("Promote() error = %v, want overlap rejection", err)
		}
	})
}

func TestPromotionAllowsNonstandardDestinationAndRecordsOwner(t *testing.T) {
	store, root := newTestStore(t)
	source := initTestGoal(t, store, root, "owner-validation")
	resourceVersion := pauseGoalForPromotion(t, store, source, "1")
	ownerRoot := filepath.Join(root, "project")
	destinationRoot := filepath.Join(ownerRoot, "maintained-records")
	options := PromoteOptions{
		GoalDir:                 source,
		DestinationGoalsRoot:    destinationRoot,
		ExpectedResourceVersion: resourceVersion,
		OwnerRoot:               ownerRoot,
	}
	if _, err := store.Promote(options); err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(destinationRoot, filepath.Base(source))
	existing, err := store.readGoalManifest(target)
	if err != nil {
		t.Fatal(err)
	}
	if got := existing.Metadata.Annotations[localOwnerRootAnnotation]; got != "project" {
		t.Fatalf("owner annotation = %q, want project", got)
	}
	if _, err := store.Promote(options); err != nil {
		t.Fatalf("idempotent Promote() failed: %v", err)
	}

	existing.Metadata.Annotations[localOwnerRootAnnotation] = "other-project"
	if err := store.writeYAML(filepath.Join(target, "goal.yaml"), existing); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Promote(options); err == nil ||
		!strings.Contains(err.Error(), "different owner root") {
		t.Fatalf("Promote() error = %v, want owner mismatch rejection", err)
	}
}

func TestPromotionRejectsUnknownRecordAndAttemptFiles(t *testing.T) {
	t.Run("source root", func(t *testing.T) {
		store, root := newTestStore(t)
		source := initTestGoal(t, store, root, "unknown-root-file")
		resourceVersion := pauseGoalForPromotion(t, store, source, "1")
		writeTestFile(t, filepath.Join(source, "notes.md"), "unexpected\n")
		_, err := store.Promote(PromoteOptions{
			GoalDir:                 source,
			DestinationGoalsRoot:    filepath.Join(root, "project", "goals"),
			ExpectedResourceVersion: resourceVersion,
		})
		if err == nil || !strings.Contains(err.Error(), "unexpected goal record entry") {
			t.Fatalf("Promote() error = %v, want root allowlist rejection", err)
		}
	})

	t.Run("source attempt", func(t *testing.T) {
		store, root := newTestStore(t)
		source := initTestGoal(t, store, root, "unknown-attempt-file")
		if _, err := store.Checkpoint(CheckpointOptions{
			GoalDir:                 source,
			ExpectedResourceVersion: "1",
			AttemptID:               "attempt-1",
		}); err != nil {
			t.Fatal(err)
		}
		review := writeRefineReview(t, root, "promotion-review.yaml")
		if _, err := store.Checkpoint(CheckpointOptions{
			GoalDir:                 source,
			ExpectedResourceVersion: "2",
			AttemptID:               "attempt-1",
			CloseAttempt:            true,
			ReviewFile:              review,
		}); err != nil {
			t.Fatal(err)
		}
		resourceVersion := pauseGoalForPromotion(t, store, source, "3")
		writeTestFile(
			t,
			filepath.Join(source, "attempts", "attempt-1", "notes.md"),
			"unexpected\n",
		)
		_, err := store.Promote(PromoteOptions{
			GoalDir:                 source,
			DestinationGoalsRoot:    filepath.Join(root, "project", "goals"),
			ExpectedResourceVersion: resourceVersion,
		})
		if err == nil ||
			(!strings.Contains(err.Error(), "unexpected entry") &&
				!strings.Contains(err.Error(), "must contain only")) {
			t.Fatalf("Promote() error = %v, want attempt allowlist rejection", err)
		}
	})

	t.Run("existing destination", func(t *testing.T) {
		store, root := newTestStore(t)
		source := initTestGoal(t, store, root, "unknown-target-file")
		resourceVersion := pauseGoalForPromotion(t, store, source, "1")
		destinationRoot := filepath.Join(root, "project", "goals")
		options := PromoteOptions{
			GoalDir:                 source,
			DestinationGoalsRoot:    destinationRoot,
			ExpectedResourceVersion: resourceVersion,
		}
		if _, err := store.Promote(options); err != nil {
			t.Fatal(err)
		}
		writeTestFile(
			t,
			filepath.Join(destinationRoot, filepath.Base(source), "notes.md"),
			"unexpected\n",
		)
		if _, err := store.Promote(options); err == nil ||
			!strings.Contains(err.Error(), "layout is invalid") {
			t.Fatalf("Promote() error = %v, want target allowlist rejection", err)
		}
	})
}

func TestLengthFramedRecordDigestsAreDeterministicAndUnambiguous(t *testing.T) {
	store, root := newTestStore(t)
	source := initTestGoal(t, store, root, "deterministic-digest")
	files, err := inspectPromotionRecord(source)
	if err != nil {
		t.Fatal(err)
	}
	recordDigest, err := digestRecord(source, files)
	if err != nil {
		t.Fatal(err)
	}
	repeatedRecordDigest, err := digestRecord(source, files)
	if err != nil {
		t.Fatal(err)
	}
	if recordDigest != repeatedRecordDigest {
		t.Fatal("promoted record digest is not deterministic")
	}

	first := map[string][]byte{
		"a": []byte("x\x00b\x00y"),
	}
	second := map[string][]byte{
		"a": []byte("x"),
		"b": []byte("y"),
	}
	firstDigest := digestLegacyFiles(first)
	if firstDigest != digestLegacyFiles(first) {
		t.Fatal("legacy record digest is not deterministic")
	}
	if firstDigest == digestLegacyFiles(second) {
		t.Fatal("length-framed digest aliases different file boundaries")
	}
}

func writeLegacyGoal(t *testing.T, source string, title string) {
	t.Helper()
	writeTestFile(
		t,
		filepath.Join(source, "README.md"),
		"# "+title+"\n\n## Acceptance criteria\n\n- The import is inspectable.\n",
	)
}

func TestMigrationBoundsLegacyEvidenceBeforeReadingContent(t *testing.T) {
	store, root := newTestStore(t)
	source := filepath.Join(root, "legacy", "bounded-import")
	writeLegacyGoal(t, source, "Bounded import")
	for index := 0; index <= maxEvidenceFiles; index++ {
		path := filepath.Join(source, fmt.Sprintf("evidence-%03d.md", index))
		content := []byte("# Evidence\n")
		if index == 0 {
			// A content-first implementation reports invalid UTF-8 here instead
			// of rejecting the directory's cardinality before allocating bodies.
			content = []byte{0xff}
		}
		if err := os.WriteFile(path, content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	destinationRoot := filepath.Join(root, "project", "goals")
	_, err := store.Migrate(MigrateOptions{
		SourceGoalDir:        source,
		DestinationGoalsRoot: destinationRoot,
	})
	if err == nil || !strings.Contains(
		err.Error(),
		fmt.Sprintf("legacy evidence cardinality exceeds %d", maxEvidenceFiles),
	) {
		t.Fatalf("Migrate() error = %v, want pre-read cardinality rejection", err)
	}
	if pathExists(filepath.Join(destinationRoot, "bounded-import")) {
		t.Fatal("over-cardinality migration published a target")
	}
}

func TestMigrationUsesCanonicalMarkdownByteLimits(t *testing.T) {
	t.Run("evidence may exceed the plan limit", func(t *testing.T) {
		store, root := newTestStore(t)
		source := filepath.Join(root, "legacy", "large-evidence-import")
		writeLegacyGoal(t, source, "Large evidence import")
		evidence := bytes.Repeat([]byte("x"), maxPlanResultBytes+1)
		if err := os.WriteFile(filepath.Join(source, "evidence.md"), evidence, 0o644); err != nil {
			t.Fatal(err)
		}
		destinationRoot := filepath.Join(root, "project", "goals")
		if _, err := store.Migrate(MigrateOptions{
			SourceGoalDir:        source,
			DestinationGoalsRoot: destinationRoot,
		}); err != nil {
			t.Fatalf("Migrate() rejected canonical-size evidence: %v", err)
		}
		imported, err := os.Stat(filepath.Join(
			destinationRoot,
			"large-evidence-import",
			"attempts",
			"imported-unversioned",
			"evidence",
			"evidence.md",
		))
		if err != nil {
			t.Fatal(err)
		}
		if imported.Size() != int64(len(evidence)) {
			t.Fatalf("imported evidence size = %d, want %d", imported.Size(), len(evidence))
		}
	})

	t.Run("evidence cannot exceed its canonical limit", func(t *testing.T) {
		store, root := newTestStore(t)
		source := filepath.Join(root, "legacy", "oversize-evidence-import")
		writeLegacyGoal(t, source, "Oversize evidence import")
		path := filepath.Join(source, "evidence.md")
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			t.Fatal(err)
		}
		if err := file.Truncate(maxEvidenceFileBytes + 1); err != nil {
			_ = file.Close()
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
		_, err = store.Migrate(MigrateOptions{
			SourceGoalDir:        source,
			DestinationGoalsRoot: filepath.Join(root, "project", "goals"),
		})
		if err == nil || !strings.Contains(
			err.Error(),
			fmt.Sprintf("at most %d bytes", maxEvidenceFileBytes),
		) {
			t.Fatalf("Migrate() error = %v, want evidence byte-limit rejection", err)
		}
	})
}

func legacySourceDigest(t *testing.T, source string) string {
	t.Helper()
	files, err := inspectUnversionedRecord(source)
	if err != nil {
		t.Fatal(err)
	}
	return digestLegacyFiles(files)
}

func TestMigrationRejectsCanonicalPathOverlapBeforePublication(t *testing.T) {
	t.Run("target symlink aliases source", func(t *testing.T) {
		store, root := newTestStore(t)
		source := filepath.Join(root, "legacy", "alias-import")
		writeLegacyGoal(t, source, "Alias import")
		destinationRoot := filepath.Join(root, "project", "goals")
		if err := os.MkdirAll(destinationRoot, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(
			source,
			filepath.Join(destinationRoot, "alias-import"),
		); err != nil {
			t.Skipf("symlinks unavailable: %v", err)
		}
		result := make(chan error, 1)
		go func() {
			_, err := store.Migrate(MigrateOptions{
				SourceGoalDir:        source,
				DestinationGoalsRoot: destinationRoot,
			})
			result <- err
		}()
		select {
		case err := <-result:
			if err == nil || !strings.Contains(err.Error(), "same") {
				t.Fatalf("Migrate() error = %v, want canonical equality rejection", err)
			}
		case <-time.After(time.Second):
			t.Fatal("Migrate() deadlocked while acquiring an aliased source/target lock")
		}
	})

	t.Run("target nested under source", func(t *testing.T) {
		store, root := newTestStore(t)
		source := filepath.Join(root, "legacy", "nested-import")
		writeLegacyGoal(t, source, "Nested import")
		destinationRoot := filepath.Join(source, "goals")
		_, err := store.Migrate(MigrateOptions{
			SourceGoalDir:        source,
			DestinationGoalsRoot: destinationRoot,
			OwnerRoot:            source,
		})
		if err == nil || !strings.Contains(err.Error(), "overlap") {
			t.Fatalf("Migrate() error = %v, want nested target rejection", err)
		}
		if pathExists(destinationRoot) {
			t.Fatal("migration created an overlapping destination root")
		}
	})

	t.Run("target is source ancestor", func(t *testing.T) {
		store, root := newTestStore(t)
		ownerRoot := filepath.Join(root, "project")
		destinationRoot := filepath.Join(ownerRoot, "goals")
		source := filepath.Join(
			destinationRoot,
			"ancestor-import",
			"legacy",
			"ancestor-import",
		)
		writeLegacyGoal(t, source, "Ancestor import")
		_, err := store.Migrate(MigrateOptions{
			SourceGoalDir:        source,
			DestinationGoalsRoot: destinationRoot,
			OwnerRoot:            ownerRoot,
		})
		if err == nil || !strings.Contains(err.Error(), "overlap") {
			t.Fatalf("Migrate() error = %v, want ancestor target rejection", err)
		}
	})
}

func TestMigrationFailureLeavesSourceUnchangedAndTargetAbsent(t *testing.T) {
	store, root := newTestStore(t)
	source := filepath.Join(root, "legacy", "failed-import")
	writeLegacyGoal(t, source, "Failed import")
	before := legacySourceDigest(t, source)
	destinationRoot := filepath.Join(root, "project", "goals")
	target := filepath.Join(destinationRoot, "failed-import")
	store.beforeRename = func(path string) error {
		if path == target {
			return errors.New("stop before publication")
		}
		return nil
	}
	_, err := store.Migrate(MigrateOptions{
		SourceGoalDir:        source,
		DestinationGoalsRoot: destinationRoot,
	})
	if err == nil || !strings.Contains(err.Error(), "stop before publication") {
		t.Fatalf("Migrate() error = %v, want injected publication failure", err)
	}
	if pathExists(target) {
		t.Fatal("failed migration published a target")
	}
	if after := legacySourceDigest(t, source); after != before {
		t.Fatal("failed migration changed the source")
	}
	entries, err := os.ReadDir(destinationRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("ordinary failure left staging residue: %v", entries)
	}
}

func TestMigrationRedigestsSourceImmediatelyBeforePublication(t *testing.T) {
	store, root := newTestStore(t)
	source := filepath.Join(root, "legacy", "changing-import")
	writeLegacyGoal(t, source, "Changing import")
	destinationRoot := filepath.Join(root, "project", "goals")
	target := filepath.Join(destinationRoot, "changing-import")
	store.beforeRename = func(path string) error {
		if path == target {
			return os.WriteFile(
				filepath.Join(source, "README.md"),
				[]byte("# Changed during import\n"),
				0o644,
			)
		}
		return nil
	}
	_, err := store.Migrate(MigrateOptions{
		SourceGoalDir:        source,
		DestinationGoalsRoot: destinationRoot,
	})
	if err == nil || !strings.Contains(err.Error(), "source changed") {
		t.Fatalf("Migrate() error = %v, want changed-source rejection", err)
	}
	if pathExists(target) {
		t.Fatal("migration published a target from a changing source")
	}
}

func TestMigrationFinalRenameNeverOverwritesTarget(t *testing.T) {
	store, root := newTestStore(t)
	source := filepath.Join(root, "legacy", "no-overwrite-import")
	writeLegacyGoal(t, source, "No overwrite import")
	destinationRoot := filepath.Join(root, "project", "goals")
	target := filepath.Join(destinationRoot, "no-overwrite-import")
	store.beforeRename = func(path string) error {
		if path == target {
			return os.Mkdir(target, 0o755)
		}
		return nil
	}
	_, err := store.Migrate(MigrateOptions{
		SourceGoalDir:        source,
		DestinationGoalsRoot: destinationRoot,
	})
	if err == nil || !strings.Contains(err.Error(), "target appeared") {
		t.Fatalf("Migrate() error = %v, want no-overwrite publication failure", err)
	}
	if !pathExists(target) || pathExists(filepath.Join(target, "goal.yaml")) {
		t.Fatal("migration overwrote the concurrently created target directory")
	}
}

func TestMigrationIdempotenceRequiresMatchingProvenanceAndOptions(t *testing.T) {
	store, root := newTestStore(t)
	source := filepath.Join(root, "legacy", "idempotent-import")
	writeLegacyGoal(t, source, "Idempotent import")
	ownerRoot := filepath.Join(root, "project")
	destinationRoot := filepath.Join(ownerRoot, "maintained-imports")
	options := MigrateOptions{
		SourceGoalDir:        source,
		DestinationGoalsRoot: destinationRoot,
		OwnerRoot:            ownerRoot,
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
		t.Fatalf("idempotent migration results differ: %+v %+v", first, second)
	}
	migrated, err := store.readGoalManifest(
		filepath.Join(destinationRoot, filepath.Base(source)),
	)
	if err != nil {
		t.Fatal(err)
	}
	if got := migrated.Metadata.Annotations[localOwnerRootAnnotation]; got != "project" {
		t.Fatalf("owner annotation = %q, want project", got)
	}
	for _, test := range []struct {
		name    string
		change  func(*MigrateOptions)
		message string
	}{
		{
			name: "title",
			change: func(changed *MigrateOptions) {
				changed.Title = "Different title"
			},
			message: "different provenance or options",
		},
		{
			name: "scope",
			change: func(changed *MigrateOptions) {
				changed.Scope = "project"
			},
			message: "different provenance or options",
		},
		{
			name: "owner",
			change: func(changed *MigrateOptions) {
				changed.OwnerRoot = filepath.Join(root, "other-project")
			},
			message: "different provenance or options",
		},
		{
			name: "criteria",
			change: func(changed *MigrateOptions) {
				changed.Criteria = []string{"Different criterion."}
			},
			message: "different criteria options",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			changed := options
			test.change(&changed)
			if _, err := store.Migrate(changed); err == nil ||
				!strings.Contains(err.Error(), test.message) {
				t.Fatalf("Migrate() changed-option error = %v", err)
			}
		})
	}
	writeTestFile(t, filepath.Join(source, "NOTES.md"), "# New source input\n")
	if _, err := store.Migrate(options); err == nil ||
		!strings.Contains(err.Error(), "different provenance or options") {
		t.Fatalf("Migrate() changed-provenance error = %v", err)
	}
}

func TestMigrationNeverOverwritesExistingUnrelatedTarget(t *testing.T) {
	t.Run("valid goal", func(t *testing.T) {
		store, root := newTestStore(t)
		source := filepath.Join(root, "legacy", "unrelated-target")
		writeLegacyGoal(t, source, "Unrelated target")
		destinationRoot := filepath.Join(root, "project", "goals")
		if _, err := store.Init(InitOptions{
			GoalsRoot: destinationRoot,
			Title:     "Existing unrelated goal",
			GoalID:    "unrelated-target",
			Scope:     "workspace",
		}); err != nil {
			t.Fatal(err)
		}
		targetManifest := filepath.Join(destinationRoot, "unrelated-target", "goal.yaml")
		before, err := os.ReadFile(targetManifest)
		if err != nil {
			t.Fatal(err)
		}
		_, err = store.Migrate(MigrateOptions{
			SourceGoalDir:        source,
			DestinationGoalsRoot: destinationRoot,
		})
		if err == nil || !strings.Contains(err.Error(), "different provenance") {
			t.Fatalf("Migrate() error = %v, want unrelated-target rejection", err)
		}
		after, err := os.ReadFile(targetManifest)
		if err != nil {
			t.Fatal(err)
		}
		if string(after) != string(before) {
			t.Fatal("migration changed an unrelated valid target")
		}
	})

	t.Run("invalid directory", func(t *testing.T) {
		store, root := newTestStore(t)
		source := filepath.Join(root, "legacy", "invalid-target")
		writeLegacyGoal(t, source, "Invalid target")
		destinationRoot := filepath.Join(root, "project", "goals")
		target := filepath.Join(destinationRoot, "invalid-target")
		marker := filepath.Join(target, "marker.txt")
		writeTestFile(t, marker, "do not replace\n")
		_, err := store.Migrate(MigrateOptions{
			SourceGoalDir:        source,
			DestinationGoalsRoot: destinationRoot,
		})
		if err == nil || !strings.Contains(err.Error(), "target is invalid") {
			t.Fatalf("Migrate() error = %v, want invalid-target rejection", err)
		}
		content, err := os.ReadFile(marker)
		if err != nil || string(content) != "do not replace\n" {
			t.Fatalf("migration changed invalid target marker: %q, %v", content, err)
		}
	})
}
