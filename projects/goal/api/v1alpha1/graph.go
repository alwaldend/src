package v1alpha1

import (
	"fmt"
	"sort"
)

// GoalGraphState describes whether a graph can be resolved and is acyclic.
type GoalGraphState string

const (
	GoalGraphStateValid   GoalGraphState = "Valid"
	GoalGraphStateUnknown GoalGraphState = "Unknown"
	GoalGraphStateInvalid GoalGraphState = "Invalid"
)

// GoalGraphRelationship identifies one independently analyzed edge graph.
type GoalGraphRelationship string

const (
	GoalGraphRelationshipParent       GoalGraphRelationship = "Parent"
	GoalGraphRelationshipDependency   GoalGraphRelationship = "Dependency"
	GoalGraphRelationshipSupersession GoalGraphRelationship = "Supersession"
)

// GoalDependencyState is a derived view of whether an open Goal's direct
// dependencies currently permit work. It is not dispatch or scheduler state.
type GoalDependencyState string

const (
	GoalDependencyStateReady   GoalDependencyState = "Ready"
	GoalDependencyStateWaiting GoalDependencyState = "Waiting"
	GoalDependencyStateBlocked GoalDependencyState = "Blocked"
	GoalDependencyStateUnknown GoalDependencyState = "Unknown"
	GoalDependencyStateClosed  GoalDependencyState = "Closed"
)

// GoalGraphAnalysis contains the complete, stable analysis of a catalog.
// Callers may bound rendered projections, but analysis never truncates input.
type GoalGraphAnalysis struct {
	State         GoalGraphState             `json:"state" yaml:"state"`
	Nodes         []GoalGraphNode            `json:"nodes" yaml:"nodes"`
	Relationships []GoalRelationshipAnalysis `json:"relationships" yaml:"relationships"`
}

// GoalGraphNode contains the derived dependency state for one Goal.
type GoalGraphNode struct {
	GoalRef         GoalReference       `json:"goalRef" yaml:"goalRef"`
	Outcome         string              `json:"outcome" yaml:"outcome"`
	DependencyState GoalDependencyState `json:"dependencyState" yaml:"dependencyState"`
}

// GoalRelationshipAnalysis reports one relationship graph independently.
// Invalid means at least one cycle; Unknown means at least one unresolved
// reference and no cycle; Valid means every reference resolves and no cycle
// exists.
type GoalRelationshipAnalysis struct {
	Relationship     GoalGraphRelationship `json:"relationship" yaml:"relationship"`
	State            GoalGraphState        `json:"state" yaml:"state"`
	Edges            []GoalGraphEdge       `json:"edges" yaml:"edges"`
	CyclicComponents []GoalGraphCycle      `json:"cyclicComponents" yaml:"cyclicComponents"`
}

// GoalGraphEdge is directed from the declaring Goal to its referenced Goal.
type GoalGraphEdge struct {
	FromGoalRef GoalReference `json:"fromGoalRef" yaml:"fromGoalRef"`
	ToGoalRef   GoalReference `json:"toGoalRef" yaml:"toGoalRef"`
	Resolved    bool          `json:"resolved" yaml:"resolved"`
}

// GoalGraphCycle is one cyclic strongly connected component. It avoids the
// potentially exponential cost of enumerating every simple cycle.
type GoalGraphCycle struct {
	GoalRefs []GoalReference `json:"goalRefs" yaml:"goalRefs"`
}

// AnalyzeGoalGraph validates and analyzes a complete catalog in memory. Goal
// names must be unique across the catalog because GoalReference is name-only.
// The function performs no I/O, never schedules work, and returns stable,
// fully sorted slices independent of input and relationship order.
func AnalyzeGoalGraph(goals []Goal) (GoalGraphAnalysis, error) {
	normalized := append([]Goal{}, goals...)
	for index := range normalized {
		normalized[index] = normalized[index].Normalized()
	}
	sort.Slice(normalized, func(left int, right int) bool {
		return normalized[left].Metadata.Name <
			normalized[right].Metadata.Name
	})

	catalog := make(map[string]Goal, len(normalized))
	for index, goal := range normalized {
		if index > 0 && normalized[index-1].Metadata.Name ==
			goal.Metadata.Name {
			return GoalGraphAnalysis{}, fmt.Errorf(
				"catalog contains duplicate goal name %q",
				goal.Metadata.Name,
			)
		}
		if err := goal.Validate(); err != nil {
			return GoalGraphAnalysis{}, fmt.Errorf(
				"goal %q is invalid: %w",
				goal.Metadata.Name,
				err,
			)
		}
		catalog[goal.Metadata.Name] = goal
	}

	relationships := []GoalGraphRelationship{
		GoalGraphRelationshipParent,
		GoalGraphRelationshipDependency,
		GoalGraphRelationshipSupersession,
	}
	analysis := GoalGraphAnalysis{
		State:         GoalGraphStateValid,
		Nodes:         make([]GoalGraphNode, 0, len(normalized)),
		Relationships: make([]GoalRelationshipAnalysis, 0, len(relationships)),
	}
	dependencyCycles := map[string]bool{}
	for _, relationship := range relationships {
		relation := analyzeRelationshipGraph(
			normalized,
			catalog,
			relationship,
		)
		analysis.Relationships = append(analysis.Relationships, relation)
		if relationship == GoalGraphRelationshipDependency {
			for _, cycle := range relation.CyclicComponents {
				for _, reference := range cycle.GoalRefs {
					dependencyCycles[reference.Name] = true
				}
			}
		}
		analysis.State = combineGraphState(analysis.State, relation.State)
	}
	for _, goal := range normalized {
		analysis.Nodes = append(analysis.Nodes, GoalGraphNode{
			GoalRef: GoalReference{Name: goal.Metadata.Name},
			Outcome: goal.Status.Outcome,
			DependencyState: dependencyState(
				goal,
				catalog,
				dependencyCycles,
			),
		})
	}
	return analysis, nil
}

func analyzeRelationshipGraph(
	goals []Goal,
	catalog map[string]Goal,
	relationship GoalGraphRelationship,
) GoalRelationshipAnalysis {
	relation := GoalRelationshipAnalysis{
		Relationship:     relationship,
		State:            GoalGraphStateValid,
		Edges:            []GoalGraphEdge{},
		CyclicComponents: []GoalGraphCycle{},
	}
	adjacency := make(map[string][]string, len(goals))
	for _, goal := range goals {
		adjacency[goal.Metadata.Name] = []string{}
		for _, reference := range relationshipReferences(
			goal.Spec.Relationships,
			relationship,
		) {
			_, resolved := catalog[reference.Name]
			relation.Edges = append(relation.Edges, GoalGraphEdge{
				FromGoalRef: GoalReference{Name: goal.Metadata.Name},
				ToGoalRef:   reference,
				Resolved:    resolved,
			})
			if resolved {
				adjacency[goal.Metadata.Name] = append(
					adjacency[goal.Metadata.Name],
					reference.Name,
				)
			} else {
				relation.State = GoalGraphStateUnknown
			}
		}
	}
	relation.CyclicComponents = cyclicComponents(adjacency)
	if len(relation.CyclicComponents) != 0 {
		relation.State = GoalGraphStateInvalid
	}
	return relation
}

func relationshipReferences(
	relationships Relationships,
	relationship GoalGraphRelationship,
) []GoalReference {
	switch relationship {
	case GoalGraphRelationshipParent:
		if relationships.ParentGoalRef == nil {
			return []GoalReference{}
		}
		return []GoalReference{*relationships.ParentGoalRef}
	case GoalGraphRelationshipDependency:
		return relationships.DependsOnGoalRefs
	case GoalGraphRelationshipSupersession:
		return relationships.SupersedesGoalRefs
	default:
		return []GoalReference{}
	}
}

func dependencyState(
	goal Goal,
	catalog map[string]Goal,
	cycleMembers map[string]bool,
) GoalDependencyState {
	if goal.Status.Outcome != "open" {
		return GoalDependencyStateClosed
	}
	if cycleMembers[goal.Metadata.Name] {
		return GoalDependencyStateUnknown
	}
	hasUnknown := false
	hasWaiting := false
	hasBlocked := false
	for _, reference := range goal.Spec.Relationships.DependsOnGoalRefs {
		dependency, resolved := catalog[reference.Name]
		if !resolved {
			hasUnknown = true
			continue
		}
		switch dependency.Status.Outcome {
		case "achieved":
		case "abandoned", "superseded":
			hasBlocked = true
		case "open":
			hasWaiting = true
		}
	}
	if hasBlocked {
		return GoalDependencyStateBlocked
	}
	if hasUnknown {
		return GoalDependencyStateUnknown
	}
	if hasWaiting {
		return GoalDependencyStateWaiting
	}
	return GoalDependencyStateReady
}

func combineGraphState(
	current GoalGraphState,
	next GoalGraphState,
) GoalGraphState {
	if current == GoalGraphStateInvalid || next == GoalGraphStateInvalid {
		return GoalGraphStateInvalid
	}
	if current == GoalGraphStateUnknown || next == GoalGraphStateUnknown {
		return GoalGraphStateUnknown
	}
	return GoalGraphStateValid
}

func cyclicComponents(adjacency map[string][]string) []GoalGraphCycle {
	names := make([]string, 0, len(adjacency))
	for name := range adjacency {
		names = append(names, name)
		sort.Strings(adjacency[name])
	}
	sort.Strings(names)

	index := 0
	indices := map[string]int{}
	lowLinks := map[string]int{}
	onStack := map[string]bool{}
	stack := []string{}
	cycles := []GoalGraphCycle{}
	var visit func(string)
	visit = func(name string) {
		indices[name] = index
		lowLinks[name] = index
		index++
		stack = append(stack, name)
		onStack[name] = true

		for _, neighbor := range adjacency[name] {
			neighborIndex, visited := indices[neighbor]
			if !visited {
				visit(neighbor)
				if lowLinks[neighbor] < lowLinks[name] {
					lowLinks[name] = lowLinks[neighbor]
				}
			} else if onStack[neighbor] && neighborIndex < lowLinks[name] {
				lowLinks[name] = neighborIndex
			}
		}

		if lowLinks[name] != indices[name] {
			return
		}
		component := []string{}
		for {
			last := len(stack) - 1
			member := stack[last]
			stack = stack[:last]
			onStack[member] = false
			component = append(component, member)
			if member == name {
				break
			}
		}
		sort.Strings(component)
		if len(component) == 1 &&
			!containsName(adjacency[component[0]], component[0]) {
			return
		}
		cycle := GoalGraphCycle{GoalRefs: make([]GoalReference, 0, len(component))}
		for _, member := range component {
			cycle.GoalRefs = append(cycle.GoalRefs, GoalReference{Name: member})
		}
		cycles = append(cycles, cycle)
	}
	for _, name := range names {
		if _, visited := indices[name]; !visited {
			visit(name)
		}
	}
	sort.Slice(cycles, func(left int, right int) bool {
		return cycles[left].GoalRefs[0].Name < cycles[right].GoalRefs[0].Name
	})
	return cycles
}

func containsName(values []string, name string) bool {
	index := sort.SearchStrings(values, name)
	return index < len(values) && values[index] == name
}
