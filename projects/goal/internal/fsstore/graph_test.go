package fsstore

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	v1alpha1 "git.alwaldend.com/alwaldend/src/projects/goal/api/v1alpha1"
)

func TestGraphIsDeterministicAndDerivesDependencyState(t *testing.T) {
	store, root := newTestStore(t)
	goalDirs := map[string]string{}
	for _, name := range []string{"waiting", "unknown", "ready", "achieved"} {
		goalDirs[name] = initTestGoal(t, store, root, name)
	}
	achieveGraphGoal(t, store, root, goalDirs["achieved"])
	rewriteGraphGoal(t, store, goalDirs["ready"], func(goal *v1alpha1.Goal) {
		goal.Spec.Relationships.DependsOnGoalRefs = []v1alpha1.GoalReference{
			{Name: "achieved"},
		}
	})
	rewriteGraphGoal(t, store, goalDirs["waiting"], func(goal *v1alpha1.Goal) {
		goal.Spec.Relationships.DependsOnGoalRefs = []v1alpha1.GoalReference{
			{Name: "ready"},
		}
	})
	rewriteGraphGoal(t, store, goalDirs["unknown"], func(goal *v1alpha1.Goal) {
		goal.Spec.Relationships.DependsOnGoalRefs = []v1alpha1.GoalReference{
			{Name: "not-created-yet"},
		}
	})

	goalsRoot := filepath.Join(root, "out", "task", "goals")
	first, err := store.Graph(goalsRoot)
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Graph(goalsRoot)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("graph output is not deterministic:\nfirst=%#v\nsecond=%#v", first, second)
	}
	if first.State != v1alpha1.GoalGraphStateUnknown {
		t.Fatalf("graph state = %s, want Unknown", first.State)
	}
	wantOrder := []string{"achieved", "ready", "unknown", "waiting"}
	wantStates := map[string]v1alpha1.GoalDependencyState{
		"achieved": v1alpha1.GoalDependencyStateClosed,
		"ready":    v1alpha1.GoalDependencyStateReady,
		"unknown":  v1alpha1.GoalDependencyStateUnknown,
		"waiting":  v1alpha1.GoalDependencyStateWaiting,
	}
	gotOrder := make([]string, 0, len(first.Nodes))
	for _, node := range first.Nodes {
		gotOrder = append(gotOrder, node.GoalRef.Name)
		if node.DependencyState != wantStates[node.GoalRef.Name] {
			t.Errorf(
				"%s dependency state = %s, want %s",
				node.GoalRef.Name,
				node.DependencyState,
				wantStates[node.GoalRef.Name],
			)
		}
	}
	if !reflect.DeepEqual(gotOrder, wantOrder) {
		t.Fatalf("node order = %#v, want %#v", gotOrder, wantOrder)
	}
}

func TestGraphRejectsInvalidCompleteRecordBeforeDerivingReadiness(
	t *testing.T,
) {
	store, root := newTestStore(t)
	achievedDir := initTestGoal(t, store, root, "achieved")
	dependentDir := initTestGoal(t, store, root, "dependent")
	achieveGraphGoal(t, store, root, achievedDir)
	rewriteGraphGoal(t, store, dependentDir, func(goal *v1alpha1.Goal) {
		goal.Spec.Relationships.DependsOnGoalRefs = []v1alpha1.GoalReference{
			{Name: "achieved"},
		}
	})
	if err := os.WriteFile(
		filepath.Join(achievedDir, "attempts", "attempt-001", "result.md"),
		[]byte("tampered\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	goalsRoot := filepath.Join(root, "out", "task", "goals")
	if _, err := store.Graph(goalsRoot); err == nil ||
		!strings.Contains(err.Error(), "artifact digest") {
		t.Fatalf("invalid complete-record graph error = %v", err)
	}
}

func TestSetRelationshipsRefreshesProjectionAndRejectsStaleAndCycles(
	t *testing.T,
) {
	store, root := newTestStore(t)
	alphaDir := initTestGoal(t, store, root, "alpha")
	betaDir := initTestGoal(t, store, root, "beta")
	_ = initTestGoal(t, store, root, "parent")
	criteriaPath := filepath.Join(alphaDir, "criteria.yaml")
	readmePath := filepath.Join(alphaDir, "README.md")
	criteriaBefore, err := os.ReadFile(criteriaPath)
	if err != nil {
		t.Fatal(err)
	}
	readmeBefore, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatal(err)
	}

	reference, err := store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 alphaDir,
		ExpectedResourceVersion: "1",
		ParentGoal:              "parent",
		DependsOn:               []string{"not-created-yet", "beta"},
		Supersedes:              []string{"parent"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if reference.ResourceVersion != "2" {
		t.Fatalf("resourceVersion = %q, want 2", reference.ResourceVersion)
	}
	alpha, err := store.readGoalManifest(alphaDir)
	if err != nil {
		t.Fatal(err)
	}
	if alpha.Metadata.Generation != 2 || alpha.Metadata.ResourceVersion != "2" {
		t.Fatalf("unexpected updated metadata: %#v", alpha.Metadata)
	}
	if alpha.Spec.Relationships.ParentGoalRef == nil ||
		alpha.Spec.Relationships.ParentGoalRef.Name != "parent" {
		t.Fatalf("unexpected parent relationship: %#v", alpha.Spec.Relationships)
	}
	wantDependencies := []v1alpha1.GoalReference{
		{Name: "beta"},
		{Name: "not-created-yet"},
	}
	if !reflect.DeepEqual(
		alpha.Spec.Relationships.DependsOnGoalRefs,
		wantDependencies,
	) {
		t.Fatalf(
			"dependencies = %#v, want %#v",
			alpha.Spec.Relationships.DependsOnGoalRefs,
			wantDependencies,
		)
	}
	criteriaAfter, err := os.ReadFile(criteriaPath)
	if err != nil {
		t.Fatal(err)
	}
	readmeAfter, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(criteriaBefore, criteriaAfter) {
		t.Fatal("relationship update changed criteria")
	}
	if bytes.Equal(readmeBefore, readmeAfter) {
		t.Fatal("relationship update did not refresh README")
	}
	for _, want := range []string{
		"- Resource version: `2`",
		"- Parent: `parent`",
		"- Depends on: `beta`, `not-created-yet`",
		"- Supersedes: `parent`",
	} {
		if !bytes.Contains(readmeAfter, []byte(want)) {
			t.Errorf("refreshed README does not contain %q", want)
		}
	}

	_, err = store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 alphaDir,
		ExpectedResourceVersion: "2",
		ClearParent:             true,
		DependsOn:               []string{"beta"},
	})
	if err != nil {
		t.Fatal(err)
	}
	alpha, err = store.readGoalManifest(alphaDir)
	if err != nil {
		t.Fatal(err)
	}
	if alpha.Spec.Relationships.ParentGoalRef != nil ||
		len(alpha.Spec.Relationships.SupersedesGoalRefs) != 0 ||
		alpha.Metadata.ResourceVersion != "3" ||
		alpha.Metadata.Generation != 3 {
		t.Fatalf("parent was not cleared cleanly: %#v", alpha)
	}
	alphaBytes, err := os.ReadFile(filepath.Join(alphaDir, "goal.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 alphaDir,
		ExpectedResourceVersion: "2",
	})
	if err == nil || !strings.Contains(err.Error(), "stale resourceVersion") {
		t.Fatalf("stale update error = %v", err)
	}
	unchangedAlphaBytes, err := os.ReadFile(filepath.Join(alphaDir, "goal.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(alphaBytes, unchangedAlphaBytes) {
		t.Fatal("stale relationship update changed goal.yaml")
	}

	betaBytes, err := os.ReadFile(filepath.Join(betaDir, "goal.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 betaDir,
		ExpectedResourceVersion: "1",
		DependsOn:               []string{"alpha"},
	})
	if err == nil || !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("cycle update error = %v", err)
	}
	unchangedBetaBytes, err := os.ReadFile(filepath.Join(betaDir, "goal.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(betaBytes, unchangedBetaBytes) {
		t.Fatal("rejected cycle changed goal.yaml")
	}
}

func TestSetRelationshipsNoOpStillAdvancesVersions(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "no-op")

	reference, err := store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if reference.ResourceVersion != "2" {
		t.Fatalf("resourceVersion = %q, want 2", reference.ResourceVersion)
	}
	goal, err := store.readGoalManifest(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if goal.Metadata.ResourceVersion != "2" || goal.Metadata.Generation != 2 {
		t.Fatalf("no-op request did not advance metadata: %#v", goal.Metadata)
	}
	readme, err := os.ReadFile(filepath.Join(goalDir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(readme, []byte("- Resource version: `2`")) {
		t.Fatal("no-op request did not refresh README")
	}
}

func TestSetRelationshipsRejectsActiveAttempt(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "active-attempt")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
		AttemptID:               "attempt-001",
		WorkType:                "change",
	}); err != nil {
		t.Fatal(err)
	}
	goalBefore, err := os.ReadFile(filepath.Join(goalDir, "goal.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "2",
		DependsOn:               []string{"another-goal"},
	})
	if err == nil || !strings.Contains(err.Error(), "no active attempt") {
		t.Fatalf("active-attempt update error = %v", err)
	}
	goalAfter, readErr := os.ReadFile(filepath.Join(goalDir, "goal.yaml"))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !bytes.Equal(goalBefore, goalAfter) {
		t.Fatal("rejected active-attempt update changed goal.yaml")
	}
}

func TestSetRelationshipsValidatesCompleteTargetRecord(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "invalid-record")
	criteriaPath := filepath.Join(goalDir, "criteria.yaml")
	var criteria CriteriaManifest
	if err := store.readYAML(criteriaPath, &criteria); err != nil {
		t.Fatal(err)
	}
	criteria.Spec.GoalRef.Name = "different-goal"
	if err := store.writeYAML(criteriaPath, criteria); err != nil {
		t.Fatal(err)
	}
	goalBefore, err := os.ReadFile(filepath.Join(goalDir, "goal.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
		DependsOn:               []string{"another-goal"},
	})
	if err == nil || !strings.Contains(err.Error(), "does not match goal") {
		t.Fatalf("invalid-record update error = %v", err)
	}
	goalAfter, readErr := os.ReadFile(filepath.Join(goalDir, "goal.yaml"))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !bytes.Equal(goalBefore, goalAfter) {
		t.Fatal("invalid-record update changed goal.yaml")
	}
}

func TestSetRelationshipsReportsCommittedProjectionFailure(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "projection-failure")
	readmePath := filepath.Join(goalDir, "README.md")
	projectionError := errors.New("injected projection failure")
	store.beforeRename = func(target string) error {
		if target == readmePath {
			return projectionError
		}
		return nil
	}

	reference, err := store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
		DependsOn:               []string{"another-goal"},
	})
	if !errors.Is(err, projectionError) ||
		!strings.Contains(err.Error(), "relationships committed") {
		t.Fatalf("projection failure error = %v", err)
	}
	if reference.ResourceVersion != "2" {
		t.Fatalf("committed reference = %#v, want resourceVersion 2", reference)
	}
	goal, readErr := store.readGoalManifest(goalDir)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if goal.Metadata.ResourceVersion != "2" {
		t.Fatalf("committed goal resourceVersion = %q, want 2", goal.Metadata.ResourceVersion)
	}
}

func TestSetRelationshipsRejectsBoundAndDuplicatesDeterministically(
	t *testing.T,
) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "bounded")
	goalBefore, err := os.ReadFile(filepath.Join(goalDir, "goal.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	overBound := make(
		[]string,
		v1alpha1.MaxGoalRelationshipReferences+1,
	)
	for index := range overBound {
		overBound[index] = fmt.Sprintf("dependency-%03d", index)
	}
	_, err = store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
		DependsOn:               overBound,
	})
	if err == nil || !strings.Contains(err.Error(), "cardinality") {
		t.Fatalf("over-bound update error = %v", err)
	}

	duplicateErrors := make([]string, 0, 2)
	for _, references := range [][]string{
		{"duplicate", "other", "duplicate"},
		{"duplicate", "duplicate", "other"},
	} {
		_, err := store.SetRelationships(SetRelationshipsOptions{
			GoalDir:                 goalDir,
			ExpectedResourceVersion: "1",
			DependsOn:               references,
		})
		if err == nil || !strings.Contains(err.Error(), "duplicate") {
			t.Fatalf("duplicate update error = %v", err)
		}
		duplicateErrors = append(duplicateErrors, err.Error())
	}
	if duplicateErrors[0] != duplicateErrors[1] {
		t.Fatalf(
			"duplicate errors depend on input order: %#v",
			duplicateErrors,
		)
	}
	goalAfter, readErr := os.ReadFile(filepath.Join(goalDir, "goal.yaml"))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !bytes.Equal(goalBefore, goalAfter) {
		t.Fatal("rejected bounded or duplicate update changed goal.yaml")
	}
}

func TestRelationshipProjectionIsBoundedAndDeterministic(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "projection-order")
	goal, err := store.readGoalManifest(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	var criteria CriteriaManifest
	if err := store.readYAML(
		filepath.Join(goalDir, "criteria.yaml"),
		&criteria,
	); err != nil {
		t.Fatal(err)
	}
	first := goal
	first.Spec.Relationships.DependsOnGoalRefs = []LocalGoalReference{
		{Name: "charlie"},
		{Name: "alpha"},
		{Name: "bravo"},
	}
	second := goal
	second.Spec.Relationships.DependsOnGoalRefs = []LocalGoalReference{
		{Name: "bravo"},
		{Name: "charlie"},
		{Name: "alpha"},
	}
	firstProjection, err := renderREADME(first, criteria, nil, 2)
	if err != nil {
		t.Fatal(err)
	}
	secondProjection, err := renderREADME(second, criteria, nil, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(firstProjection, secondProjection) {
		t.Fatal("relationship projection depends on semantic-free input order")
	}
	want := "- Depends on: `alpha`, `bravo`, … 1 omitted"
	if !bytes.Contains(firstProjection, []byte(want)) {
		t.Fatalf("bounded relationship projection does not contain %q", want)
	}
}

func TestSetRelationshipsRejectsVisibleCycleForEveryEdgeKind(t *testing.T) {
	for _, relationship := range []string{"parent", "dependency", "supersession"} {
		t.Run(relationship, func(t *testing.T) {
			store, root := newTestStore(t)
			alphaDir := initTestGoal(t, store, root, "alpha")
			betaDir := initTestGoal(t, store, root, "beta")
			first := relationshipOptions(
				relationship,
				alphaDir,
				"1",
				"beta",
			)
			if _, err := store.SetRelationships(first); err != nil {
				t.Fatal(err)
			}
			betaBefore, err := os.ReadFile(filepath.Join(betaDir, "goal.yaml"))
			if err != nil {
				t.Fatal(err)
			}
			second := relationshipOptions(
				relationship,
				betaDir,
				"1",
				"alpha",
			)
			if _, err := store.SetRelationships(second); err == nil ||
				!strings.Contains(err.Error(), "cycle") {
				t.Fatalf("%s cycle error = %v", relationship, err)
			}
			betaAfter, err := os.ReadFile(filepath.Join(betaDir, "goal.yaml"))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(betaBefore, betaAfter) {
				t.Fatalf("rejected %s cycle changed goal.yaml", relationship)
			}
		})
	}
}

func TestSetRelationshipsRepairsOneOfTwoDisjointCycles(t *testing.T) {
	store, root := newTestStore(t)
	directories := map[string]string{}
	for _, name := range []string{"cycle-a", "cycle-b", "cycle-c", "cycle-d"} {
		directories[name] = initTestGoal(t, store, root, name)
	}
	for _, edge := range []struct {
		source string
		target string
	}{
		{source: "cycle-a", target: "cycle-b"},
		{source: "cycle-b", target: "cycle-a"},
		{source: "cycle-c", target: "cycle-d"},
		{source: "cycle-d", target: "cycle-c"},
	} {
		rewriteGraphGoal(t, store, directories[edge.source], func(goal *v1alpha1.Goal) {
			goal.Spec.Relationships.DependsOnGoalRefs =
				[]v1alpha1.GoalReference{{Name: edge.target}}
		})
	}

	reference, err := store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 directories["cycle-a"],
		ExpectedResourceVersion: "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if reference.ResourceVersion != "2" {
		t.Fatalf("resourceVersion = %q, want 2", reference.ResourceVersion)
	}
	analysis, err := store.Graph(filepath.Join(root, "out", "task", "goals"))
	if err != nil {
		t.Fatal(err)
	}
	if analysis.State != v1alpha1.GoalGraphStateInvalid {
		t.Fatalf("catalog state = %s, want Invalid", analysis.State)
	}
	wantCycles := [][]string{{"cycle-c", "cycle-d"}}
	if got := graphCycleNames(
		t,
		analysis,
		v1alpha1.GoalGraphRelationshipDependency,
	); !reflect.DeepEqual(got, wantCycles) {
		t.Fatalf("remaining cycles = %#v, want %#v", got, wantCycles)
	}

	goalBefore, err := os.ReadFile(
		filepath.Join(directories["cycle-a"], "goal.yaml"),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 directories["cycle-a"],
		ExpectedResourceVersion: "2",
		DependsOn:               []string{"cycle-b"},
	})
	if err == nil || !strings.Contains(err.Error(), "introduces") {
		t.Fatalf("reintroduced cycle error = %v", err)
	}
	goalAfter, readErr := os.ReadFile(
		filepath.Join(directories["cycle-a"], "goal.yaml"),
	)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !bytes.Equal(goalBefore, goalAfter) {
		t.Fatal("rejected reintroduced cycle changed goal.yaml")
	}
}

func TestSetRelationshipsAllowsReducedCyclicComponent(t *testing.T) {
	store, root := newTestStore(t)
	aDir := initTestGoal(t, store, root, "cycle-a")
	bDir := initTestGoal(t, store, root, "cycle-b")
	cDir := initTestGoal(t, store, root, "cycle-c")
	rewriteGraphGoal(t, store, aDir, func(goal *v1alpha1.Goal) {
		goal.Spec.Relationships.DependsOnGoalRefs =
			[]v1alpha1.GoalReference{{Name: "cycle-b"}}
	})
	rewriteGraphGoal(t, store, bDir, func(goal *v1alpha1.Goal) {
		goal.Spec.Relationships.DependsOnGoalRefs =
			[]v1alpha1.GoalReference{{Name: "cycle-c"}}
	})
	rewriteGraphGoal(t, store, cDir, func(goal *v1alpha1.Goal) {
		goal.Spec.Relationships.DependsOnGoalRefs =
			[]v1alpha1.GoalReference{{Name: "cycle-a"}}
	})

	if _, err := store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 aDir,
		ExpectedResourceVersion: "1",
		DependsOn:               []string{"cycle-c"},
	}); err != nil {
		t.Fatal(err)
	}
	analysis, err := store.Graph(filepath.Join(root, "out", "task", "goals"))
	if err != nil {
		t.Fatal(err)
	}
	wantCycles := [][]string{{"cycle-a", "cycle-c"}}
	if got := graphCycleNames(
		t,
		analysis,
		v1alpha1.GoalGraphRelationshipDependency,
	); !reflect.DeepEqual(got, wantCycles) {
		t.Fatalf("reduced cycles = %#v, want %#v", got, wantCycles)
	}
}

func TestSetRelationshipsRejectsWorsenedCycle(t *testing.T) {
	store, root := newTestStore(t)
	aDir := initTestGoal(t, store, root, "cycle-a")
	bDir := initTestGoal(t, store, root, "cycle-b")
	cDir := initTestGoal(t, store, root, "cycle-c")
	rewriteGraphGoal(t, store, aDir, func(goal *v1alpha1.Goal) {
		goal.Spec.Relationships.DependsOnGoalRefs =
			[]v1alpha1.GoalReference{{Name: "cycle-b"}}
	})
	rewriteGraphGoal(t, store, bDir, func(goal *v1alpha1.Goal) {
		goal.Spec.Relationships.DependsOnGoalRefs =
			[]v1alpha1.GoalReference{{Name: "cycle-a"}}
	})
	rewriteGraphGoal(t, store, cDir, func(goal *v1alpha1.Goal) {
		goal.Spec.Relationships.DependsOnGoalRefs =
			[]v1alpha1.GoalReference{{Name: "cycle-a"}}
	})

	_, err := store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 aDir,
		ExpectedResourceVersion: "1",
		DependsOn:               []string{"cycle-b", "cycle-c"},
	})
	if err == nil || !strings.Contains(err.Error(), "worsens") {
		t.Fatalf("worsened cycle error = %v", err)
	}
}

func TestSetRelationshipsComparesCyclesByRelationshipKind(t *testing.T) {
	store, root := newTestStore(t)
	aDir := initTestGoal(t, store, root, "cycle-a")
	bDir := initTestGoal(t, store, root, "cycle-b")
	rewriteGraphGoal(t, store, aDir, func(goal *v1alpha1.Goal) {
		goal.Spec.Relationships.DependsOnGoalRefs =
			[]v1alpha1.GoalReference{{Name: "cycle-b"}}
	})
	rewriteGraphGoal(t, store, bDir, func(goal *v1alpha1.Goal) {
		goal.Spec.Relationships.DependsOnGoalRefs =
			[]v1alpha1.GoalReference{{Name: "cycle-a"}}
		goal.Spec.Relationships.ParentGoalRef =
			&v1alpha1.GoalReference{Name: "cycle-a"}
	})

	_, err := store.SetRelationships(SetRelationshipsOptions{
		GoalDir:                 aDir,
		ExpectedResourceVersion: "1",
		ParentGoal:              "cycle-b",
	})
	if err == nil || !strings.Contains(err.Error(), "Parent cycle") {
		t.Fatalf("cross-kind cycle error = %v", err)
	}
}

func TestGraphRejectsCatalogAboveBoundBeforeLoadingMembers(t *testing.T) {
	store, root := newTestStore(t)
	goalsRoot := filepath.Join(root, "out", "task", "goals")
	if err := os.MkdirAll(goalsRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	for index := 0; index <= maxGoals; index++ {
		if err := os.Mkdir(
			filepath.Join(goalsRoot, fmt.Sprintf("goal-%04d", index)),
			0o755,
		); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := store.Graph(goalsRoot); err == nil ||
		!strings.Contains(err.Error(), "goal cardinality") {
		t.Fatalf("over-bound catalog error = %v", err)
	}
}

func relationshipOptions(
	relationship string,
	goalDir string,
	resourceVersion string,
	target string,
) SetRelationshipsOptions {
	options := SetRelationshipsOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: resourceVersion,
	}
	switch relationship {
	case "parent":
		options.ParentGoal = target
	case "dependency":
		options.DependsOn = []string{target}
	case "supersession":
		options.Supersedes = []string{target}
	}
	return options
}

func graphCycleNames(
	t *testing.T,
	analysis v1alpha1.GoalGraphAnalysis,
	relationship v1alpha1.GoalGraphRelationship,
) [][]string {
	t.Helper()
	for _, relation := range analysis.Relationships {
		if relation.Relationship != relationship {
			continue
		}
		names := make([][]string, 0, len(relation.CyclicComponents))
		for _, cycle := range relation.CyclicComponents {
			members := make([]string, 0, len(cycle.GoalRefs))
			for _, reference := range cycle.GoalRefs {
				members = append(members, reference.Name)
			}
			names = append(names, members)
		}
		return names
	}
	t.Fatalf("missing %s relationship analysis", relationship)
	return nil
}

func achieveGraphGoal(t *testing.T, store *Store, root string, dir string) {
	t.Helper()
	result := filepath.Join(
		root,
		"out",
		"task",
		filepath.Base(dir)+"-result.md",
	)
	review := filepath.Join(
		root,
		"out",
		"task",
		filepath.Base(dir)+"-review.yaml",
	)
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
		GoalDir:                 dir,
		ExpectedResourceVersion: "1",
		AttemptID:               "attempt-001",
		ResultFile:              result,
		ReviewFile:              review,
		CloseAttempt:            true,
		Outcome:                 "achieved",
	}); err != nil {
		t.Fatal(err)
	}
}

func rewriteGraphGoal(
	t *testing.T,
	store *Store,
	dir string,
	mutate func(*v1alpha1.Goal),
) {
	t.Helper()
	manifest, err := store.readGoalManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	goal := v1alpha1.Goal(manifest)
	mutate(&goal)
	goal = goal.Normalized()
	if err := goal.Validate(); err != nil {
		t.Fatalf("invalid rewritten test goal: %v", err)
	}
	if err := store.writeYAML(
		filepath.Join(dir, "goal.yaml"),
		GoalManifest(goal),
	); err != nil {
		t.Fatal(err)
	}
}
