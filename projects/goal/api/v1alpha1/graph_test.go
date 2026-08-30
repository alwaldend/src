package v1alpha1

import (
	"reflect"
	"strings"
	"testing"
)

func graphGoal(name string, outcome string) Goal {
	goal := validGoal()
	goal.Metadata.Name = name
	goal.Metadata.Generation = 1
	goal.Spec.Title = name
	goal.Status.Outcome = outcome
	goal.Status.ActiveAttemptID = ""
	goal.Status.AcceptedAttemptID = ""
	goal.Status.AcceptedResultDigest = ""
	switch outcome {
	case "open":
		goal.Status.Execution = "active"
	case "achieved":
		goal.Status.Execution = "paused"
		goal.Status.AcceptedAttemptID = "attempt-001"
		goal.Status.AcceptedResultDigest = "sha256:" +
			strings.Repeat("a", 64)
	case "abandoned", "superseded":
		goal.Status.Execution = "paused"
	}
	return goal
}

func goalReferences(names ...string) []GoalReference {
	references := make([]GoalReference, 0, len(names))
	for _, name := range names {
		references = append(references, GoalReference{Name: name})
	}
	return references
}

func relationshipAnalysis(
	t *testing.T,
	analysis GoalGraphAnalysis,
	relationship GoalGraphRelationship,
) GoalRelationshipAnalysis {
	t.Helper()
	for _, candidate := range analysis.Relationships {
		if candidate.Relationship == relationship {
			return candidate
		}
	}
	t.Fatalf("missing %s relationship analysis", relationship)
	return GoalRelationshipAnalysis{}
}

func dependencyStates(analysis GoalGraphAnalysis) map[string]GoalDependencyState {
	states := make(map[string]GoalDependencyState, len(analysis.Nodes))
	for _, node := range analysis.Nodes {
		states[node.GoalRef.Name] = node.DependencyState
	}
	return states
}

func cycleNames(cycles []GoalGraphCycle) [][]string {
	names := make([][]string, 0, len(cycles))
	for _, cycle := range cycles {
		members := make([]string, 0, len(cycle.GoalRefs))
		for _, reference := range cycle.GoalRefs {
			members = append(members, reference.Name)
		}
		names = append(names, members)
	}
	return names
}

func TestAnalyzeGoalGraphDerivesDependencyStates(t *testing.T) {
	achieved := graphGoal("achieved", "achieved")
	open := graphGoal("open", "open")
	abandoned := graphGoal("abandoned", "abandoned")
	superseded := graphGoal("superseded", "superseded")
	ready := graphGoal("ready", "open")
	ready.Spec.Relationships.DependsOnGoalRefs = goalReferences("achieved")
	waiting := graphGoal("waiting", "open")
	waiting.Spec.Relationships.DependsOnGoalRefs = goalReferences(
		"open",
		"achieved",
	)
	blocked := graphGoal("blocked", "open")
	blocked.Spec.Relationships.DependsOnGoalRefs = goalReferences(
		"missing",
		"abandoned",
	)
	blockedBySupersession := graphGoal("blocked-by-supersession", "open")
	blockedBySupersession.Spec.Relationships.DependsOnGoalRefs =
		goalReferences("superseded")
	unknown := graphGoal("unknown", "open")
	unknown.Spec.Relationships.DependsOnGoalRefs =
		goalReferences("missing")

	analysis, err := AnalyzeGoalGraph([]Goal{
		unknown,
		waiting,
		superseded,
		ready,
		open,
		blockedBySupersession,
		blocked,
		achieved,
		abandoned,
	})
	if err != nil {
		t.Fatal(err)
	}
	if analysis.State != GoalGraphStateUnknown {
		t.Fatalf("catalog state = %s, want Unknown", analysis.State)
	}
	want := map[string]GoalDependencyState{
		"abandoned":               GoalDependencyStateClosed,
		"achieved":                GoalDependencyStateClosed,
		"blocked":                 GoalDependencyStateBlocked,
		"blocked-by-supersession": GoalDependencyStateBlocked,
		"open":                    GoalDependencyStateReady,
		"ready":                   GoalDependencyStateReady,
		"superseded":              GoalDependencyStateClosed,
		"unknown":                 GoalDependencyStateUnknown,
		"waiting":                 GoalDependencyStateWaiting,
	}
	if got := dependencyStates(analysis); !reflect.DeepEqual(got, want) {
		t.Fatalf("dependency states = %#v, want %#v", got, want)
	}
	dependencies := relationshipAnalysis(
		t,
		analysis,
		GoalGraphRelationshipDependency,
	)
	if dependencies.State != GoalGraphStateUnknown {
		t.Fatalf("dependency graph state = %s, want Unknown", dependencies.State)
	}
	unresolved := 0
	for _, edge := range dependencies.Edges {
		if !edge.Resolved {
			unresolved++
			if edge.ToGoalRef.Name != "missing" {
				t.Fatalf("unexpected unresolved edge: %#v", edge)
			}
		}
	}
	if unresolved != 2 {
		t.Fatalf("unresolved edge count = %d, want 2", unresolved)
	}
}

func TestAnalyzeGoalGraphDetectsEachRelationshipCycleIndependently(
	t *testing.T,
) {
	parentA := graphGoal("parent-a", "open")
	parentB := graphGoal("parent-b", "open")
	parentA.Spec.Relationships.ParentGoalRef =
		&GoalReference{Name: "parent-b"}
	parentB.Spec.Relationships.ParentGoalRef =
		&GoalReference{Name: "parent-a"}

	dependencyA := graphGoal("dependency-a", "open")
	dependencyB := graphGoal("dependency-b", "open")
	dependencyA.Spec.Relationships.DependsOnGoalRefs =
		goalReferences("dependency-b")
	dependencyB.Spec.Relationships.DependsOnGoalRefs =
		goalReferences("dependency-a")

	supersessionA := graphGoal("supersession-a", "open")
	supersessionB := graphGoal("supersession-b", "open")
	supersessionA.Spec.Relationships.SupersedesGoalRefs =
		goalReferences("supersession-b")
	supersessionB.Spec.Relationships.SupersedesGoalRefs =
		goalReferences("supersession-a")

	analysis, err := AnalyzeGoalGraph([]Goal{
		supersessionB,
		dependencyB,
		parentB,
		supersessionA,
		dependencyA,
		parentA,
	})
	if err != nil {
		t.Fatal(err)
	}
	if analysis.State != GoalGraphStateInvalid {
		t.Fatalf("catalog state = %s, want Invalid", analysis.State)
	}
	for _, expected := range []struct {
		relationship GoalGraphRelationship
		members      [][]string
	}{
		{GoalGraphRelationshipParent, [][]string{{"parent-a", "parent-b"}}},
		{GoalGraphRelationshipDependency, [][]string{{"dependency-a", "dependency-b"}}},
		{GoalGraphRelationshipSupersession, [][]string{{"supersession-a", "supersession-b"}}},
	} {
		relation := relationshipAnalysis(t, analysis, expected.relationship)
		if relation.State != GoalGraphStateInvalid {
			t.Errorf("%s state = %s, want Invalid", expected.relationship, relation.State)
		}
		if got := cycleNames(relation.CyclicComponents); !reflect.DeepEqual(got, expected.members) {
			t.Errorf("%s cycles = %#v, want %#v", expected.relationship, got, expected.members)
		}
	}
	states := dependencyStates(analysis)
	if states["dependency-a"] != GoalDependencyStateUnknown ||
		states["dependency-b"] != GoalDependencyStateUnknown {
		t.Fatalf("dependency cycle members were not Unknown: %#v", states)
	}
	for _, name := range []string{
		"parent-a",
		"parent-b",
		"supersession-a",
		"supersession-b",
	} {
		if states[name] != GoalDependencyStateReady {
			t.Errorf("%s dependency state = %s, want Ready", name, states[name])
		}
	}
}

func TestAnalyzeGoalGraphDoesNotCombineRelationshipKindsForCycles(
	t *testing.T,
) {
	first := graphGoal("mixed-a", "open")
	second := graphGoal("mixed-b", "open")
	first.Spec.Relationships.ParentGoalRef =
		&GoalReference{Name: "mixed-b"}
	second.Spec.Relationships.DependsOnGoalRefs =
		goalReferences("mixed-a")

	analysis, err := AnalyzeGoalGraph([]Goal{second, first})
	if err != nil {
		t.Fatal(err)
	}
	if analysis.State != GoalGraphStateValid {
		t.Fatalf("mixed-kind cycle produced state %s, want Valid", analysis.State)
	}
	for _, relation := range analysis.Relationships {
		if relation.State != GoalGraphStateValid ||
			len(relation.CyclicComponents) != 0 {
			t.Fatalf("mixed edges formed a %s cycle: %#v", relation.Relationship, relation)
		}
	}
}

func TestAnalyzeGoalGraphScopesUnknownToUnresolvedRelationshipKind(
	t *testing.T,
) {
	goal := graphGoal("source", "open")
	goal.Spec.Relationships.ParentGoalRef =
		&GoalReference{Name: "missing-parent"}
	goal.Spec.Relationships.SupersedesGoalRefs =
		goalReferences("missing-supersession")

	analysis, err := AnalyzeGoalGraph([]Goal{goal})
	if err != nil {
		t.Fatal(err)
	}
	if analysis.State != GoalGraphStateUnknown {
		t.Fatalf("catalog state = %s, want Unknown", analysis.State)
	}
	if got := relationshipAnalysis(
		t,
		analysis,
		GoalGraphRelationshipParent,
	).State; got != GoalGraphStateUnknown {
		t.Errorf("parent state = %s, want Unknown", got)
	}
	if got := relationshipAnalysis(
		t,
		analysis,
		GoalGraphRelationshipDependency,
	).State; got != GoalGraphStateValid {
		t.Errorf("dependency state = %s, want Valid", got)
	}
	if got := relationshipAnalysis(
		t,
		analysis,
		GoalGraphRelationshipSupersession,
	).State; got != GoalGraphStateUnknown {
		t.Errorf("supersession state = %s, want Unknown", got)
	}
	if got := dependencyStates(analysis)["source"]; got != GoalDependencyStateReady {
		t.Errorf("source dependency state = %s, want Ready", got)
	}
}

func TestAnalyzeGoalGraphOutputIsStableAndReportsAllComponents(
	t *testing.T,
) {
	a := graphGoal("cycle-a", "open")
	b := graphGoal("cycle-b", "open")
	c := graphGoal("cycle-c", "open")
	d := graphGoal("cycle-d", "open")
	e := graphGoal("cycle-e", "open")
	fanout := graphGoal("fanout", "open")
	a.Spec.Relationships.DependsOnGoalRefs = goalReferences("cycle-b")
	b.Spec.Relationships.DependsOnGoalRefs = goalReferences("cycle-a")
	c.Spec.Relationships.DependsOnGoalRefs = goalReferences("cycle-d")
	d.Spec.Relationships.DependsOnGoalRefs = goalReferences("cycle-e")
	e.Spec.Relationships.DependsOnGoalRefs = goalReferences("cycle-c")
	fanout.Spec.Relationships.DependsOnGoalRefs = goalReferences(
		"cycle-c",
		"cycle-a",
	)

	first, err := AnalyzeGoalGraph([]Goal{fanout, e, c, a, d, b})
	if err != nil {
		t.Fatal(err)
	}
	fanout.Spec.Relationships.DependsOnGoalRefs = goalReferences(
		"cycle-a",
		"cycle-c",
	)
	second, err := AnalyzeGoalGraph([]Goal{b, d, a, c, e, fanout})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("analysis depends on input ordering:\nfirst: %#v\nsecond: %#v", first, second)
	}
	dependencies := relationshipAnalysis(
		t,
		first,
		GoalGraphRelationshipDependency,
	)
	wantCycles := [][]string{
		{"cycle-a", "cycle-b"},
		{"cycle-c", "cycle-d", "cycle-e"},
	}
	if got := cycleNames(dependencies.CyclicComponents); !reflect.DeepEqual(got, wantCycles) {
		t.Fatalf("cycles = %#v, want %#v", got, wantCycles)
	}
	if got := dependencyStates(first)["fanout"]; got != GoalDependencyStateWaiting {
		t.Fatalf("fanout dependency state = %s, want Waiting", got)
	}
	wantNodeOrder := []string{
		"cycle-a",
		"cycle-b",
		"cycle-c",
		"cycle-d",
		"cycle-e",
		"fanout",
	}
	gotNodeOrder := make([]string, 0, len(first.Nodes))
	for _, node := range first.Nodes {
		gotNodeOrder = append(gotNodeOrder, node.GoalRef.Name)
	}
	if !reflect.DeepEqual(gotNodeOrder, wantNodeOrder) {
		t.Fatalf("node order = %#v, want %#v", gotNodeOrder, wantNodeOrder)
	}
}

func TestAnalyzeGoalGraphRejectsDuplicateCatalogNames(t *testing.T) {
	first := graphGoal("duplicate", "open")
	second := graphGoal("duplicate", "open")
	first.Metadata.Namespace = "first"
	second.Metadata.Namespace = "second"
	for _, goals := range [][]Goal{{first, second}, {second, first}} {
		_, err := AnalyzeGoalGraph(goals)
		if err == nil || err.Error() !=
			"catalog contains duplicate goal name \"duplicate\"" {
			t.Fatalf("duplicate catalog error = %v", err)
		}
	}
}

func TestAnalyzeGoalGraphEmptyCatalogIsStableAndValid(t *testing.T) {
	analysis, err := AnalyzeGoalGraph(nil)
	if err != nil {
		t.Fatal(err)
	}
	if analysis.State != GoalGraphStateValid || analysis.Nodes == nil ||
		len(analysis.Nodes) != 0 || len(analysis.Relationships) != 3 {
		t.Fatalf("unexpected empty analysis: %#v", analysis)
	}
	for _, relation := range analysis.Relationships {
		if relation.State != GoalGraphStateValid || relation.Edges == nil ||
			relation.CyclicComponents == nil {
			t.Fatalf("unstable empty relation: %#v", relation)
		}
	}
}
