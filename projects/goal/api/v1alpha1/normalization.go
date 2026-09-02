package v1alpha1

import "sort"

// Normalized returns a copy in the canonical persistence and digest form.
// Relationship order is semantic-free, so references are sorted by name and
// nil lists become non-nil empty lists. The receiver and its backing slices
// are never mutated.
func (r Relationships) Normalized() Relationships {
	normalized := Relationships{
		DependsOnGoalRefs:  normalizeGoalReferences(r.DependsOnGoalRefs),
		SupersedesGoalRefs: normalizeGoalReferences(r.SupersedesGoalRefs),
	}
	if r.ParentGoalRef != nil {
		parent := *r.ParentGoalRef
		normalized.ParentGoalRef = &parent
	}
	return normalized
}

// Normalized returns a copy whose relationship fields are in canonical
// persistence and digest form.
func (s GoalSpec) Normalized() GoalSpec {
	s.Relationships = s.Relationships.Normalized()
	return s
}

// Normalized returns a copy whose relationship fields are in canonical
// persistence and digest form.
func (g Goal) Normalized() Goal {
	g.Spec = g.Spec.Normalized()
	g.Status.Plans = normalizePlanSummaries(g.Status.Plans)
	return g
}

func hasPlanSummary(values []PlanSummary, planID string) bool {
	if planID == "" {
		return false
	}
	for _, plan := range values {
		if plan.PlanID == planID {
			return true
		}
	}
	return false
}

func normalizeGoalReferences(values []GoalReference) []GoalReference {
	normalized := append([]GoalReference{}, values...)
	sort.Slice(normalized, func(left int, right int) bool {
		return normalized[left].Name < normalized[right].Name
	})
	return normalized
}

func normalizePlanSummaries(values []PlanSummary) []PlanSummary {
	return append([]PlanSummary{}, values...)
}
