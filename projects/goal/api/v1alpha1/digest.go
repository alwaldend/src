package v1alpha1

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// GoalStateDigest binds an attempt to portable goal state. Backend resource
// versions, timestamps, local annotations, and provenance are excluded.
func GoalStateDigest(goal Goal) (string, error) {
	goal = goal.Normalized()
	state := struct {
		Domain               string        `json:"domain"`
		APIVersion           string        `json:"apiVersion"`
		Kind                 string        `json:"kind"`
		Name                 string        `json:"name"`
		Generation           uint64        `json:"generation"`
		Title                string        `json:"title"`
		Scope                string        `json:"scope"`
		Retention            Retention     `json:"retention"`
		Relationships        Relationships `json:"relationships"`
		LifecycleGeneration  uint64        `json:"lifecycleGeneration"`
		Outcome              string        `json:"outcome"`
		Execution            string        `json:"execution"`
		ActiveAttemptID      string        `json:"activeAttemptID"`
		AcceptedAttemptID    string        `json:"acceptedAttemptID"`
		AcceptedResultDigest string        `json:"acceptedResultDigest"`
		CriteriaRevision     uint64        `json:"criteriaRevision"`
	}{
		Domain:               "goals.alwaldend.com/portable-goal-state/v1alpha1",
		APIVersion:           goal.APIVersion,
		Kind:                 goal.Kind,
		Name:                 goal.Metadata.Name,
		Generation:           goal.Metadata.Generation,
		Title:                goal.Spec.Title,
		Scope:                goal.Spec.Scope,
		Retention:            goal.Spec.Retention,
		Relationships:        goal.Spec.Relationships,
		LifecycleGeneration:  goal.Status.LifecycleGeneration,
		Outcome:              goal.Status.Outcome,
		Execution:            goal.Status.Execution,
		ActiveAttemptID:      goal.Status.ActiveAttemptID,
		AcceptedAttemptID:    goal.Status.AcceptedAttemptID,
		AcceptedResultDigest: goal.Status.AcceptedResultDigest,
		CriteriaRevision:     goal.Status.CriteriaRevision,
	}
	content, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	return digest(content), nil
}

// CriteriaDigest binds an attempt to the portable criteria specification.
func CriteriaDigest(criteria GoalCriteria) (string, error) {
	projection := struct {
		Domain     string       `json:"domain"`
		APIVersion string       `json:"apiVersion"`
		Kind       string       `json:"kind"`
		Name       string       `json:"name"`
		Spec       CriteriaSpec `json:"spec"`
	}{
		Domain:     "goals.alwaldend.com/portable-criteria/v1alpha1",
		APIVersion: criteria.APIVersion,
		Kind:       criteria.Kind,
		Name:       criteria.Metadata.Name,
		Spec:       criteria.Spec,
	}
	content, err := json.Marshal(projection)
	if err != nil {
		return "", err
	}
	return digest(content), nil
}

// ValidateReviewAgainstCriteria ensures every review result refers to the
// exact criterion revision in the attempt's bound criteria snapshot.
func ValidateReviewAgainstCriteria(
	review CloseReview,
	criteria GoalCriteria,
) error {
	items := make(map[string]Criterion, len(criteria.Spec.Items))
	for _, item := range criteria.Spec.Items {
		items[item.CriterionID] = item
	}
	for _, result := range review.Criteria {
		item, ok := items[result.CriterionID]
		if !ok || item.Revision != result.CriterionRevision {
			return fmt.Errorf(
				"review criterion %q is absent or does not match the bound item revision",
				result.CriterionID,
			)
		}
	}
	return nil
}

// ReviewAcceptsRequired verifies the acceptance gate for an achieved goal.
func ReviewAcceptsRequired(
	review CloseReview,
	criteria GoalCriteria,
) error {
	if review.Decision != "accept" {
		return fmt.Errorf(
			"achieved outcome requires an accept close review",
		)
	}
	requiredCount := 0
	results := make(map[string]CriterionReview, len(review.Criteria))
	for _, result := range review.Criteria {
		results[result.CriterionID] = result
	}
	for _, item := range criteria.Spec.Items {
		if !item.Required {
			continue
		}
		requiredCount++
		result, ok := results[item.CriterionID]
		if !ok || result.CriterionRevision != item.Revision ||
			result.Verdict != "pass" {
			return fmt.Errorf(
				"required criterion %q revision %d is not covered by an exact pass",
				item.CriterionID,
				item.Revision,
			)
		}
	}
	if requiredCount == 0 {
		return fmt.Errorf(
			"achieved outcome requires at least one current required criterion",
		)
	}
	return nil
}

func digest(value []byte) string {
	sum := sha256.Sum256(value)
	return "sha256:" + hex.EncodeToString(sum[:])
}
