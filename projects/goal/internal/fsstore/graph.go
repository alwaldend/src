package fsstore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	v1alpha1 "git.alwaldend.com/alwaldend/src/projects/goal/api/v1alpha1"
)

// SetRelationshipsOptions describes a complete replacement of the target
// Goal's list relationships. ParentGoal preserves the current parent when it
// is empty; ClearParent explicitly removes it.
type SetRelationshipsOptions struct {
	GoalDir                 string
	ExpectedResourceVersion string
	ParentGoal              string
	ClearParent             bool
	DependsOn               []string
	Supersedes              []string
}

// Graph derives the relationship analysis for one local goal catalog. The
// analysis is never persisted.
func (s *Store) Graph(goalsRoot string) (v1alpha1.GoalGraphAnalysis, error) {
	root, err := s.resolveInsideWorkspace(goalsRoot)
	if err != nil {
		return v1alpha1.GoalGraphAnalysis{}, err
	}
	goals, err := s.loadGoalCatalog(root)
	if err != nil {
		return v1alpha1.GoalGraphAnalysis{}, err
	}
	return v1alpha1.AnalyzeGoalGraph(goals)
}

// SetRelationships replaces the target Goal's relationships and refreshes its
// derived README projection. Unresolved references are allowed so goals can be
// created independently. A new cycle or an expansion of a cycle involving the
// updated Goal is rejected independently for each relationship kind. Existing
// unrelated cycles do not prevent an update from repairing another cycle.
// Without a catalog lock, a concurrent cross-goal write can race the catalog
// snapshot and is detected by the next Graph call.
func (s *Store) SetRelationships(
	options SetRelationshipsOptions,
) (GoalReference, error) {
	if _, err := parseResourceVersion(
		options.ExpectedResourceVersion,
	); err != nil {
		return GoalReference{}, fmt.Errorf("expected resource version: %w", err)
	}
	if options.ParentGoal != "" && options.ClearParent {
		return GoalReference{}, fmt.Errorf(
			"parent goal and clear parent are mutually exclusive",
		)
	}
	dir, err := s.resolveInsideWorkspace(options.GoalDir)
	if err != nil {
		return GoalReference{}, err
	}
	root := filepath.Dir(dir)

	// Read each catalog member under its own lock, releasing it before moving
	// to the next. Keeping no catalog-wide lock avoids serializing unrelated
	// goal operations and prevents cross-goal lock-order deadlocks.
	catalog, err := s.loadGoalCatalog(root)
	if err != nil {
		return GoalReference{}, err
	}

	lock, err := s.acquireGoalLock(dir)
	if err != nil {
		return GoalReference{}, err
	}
	defer lock.release()
	if err := s.checkNoPendingPublication(dir); err != nil {
		return GoalReference{}, err
	}
	goal, criteria, attempts, err := s.loadAndValidate(dir)
	if err != nil {
		return GoalReference{}, err
	}
	if goal.Metadata.ResourceVersion != options.ExpectedResourceVersion {
		return GoalReference{}, fmt.Errorf(
			"stale resourceVersion: expected %s, current is %s",
			options.ExpectedResourceVersion,
			goal.Metadata.ResourceVersion,
		)
	}
	if goal.Status.ActiveAttemptID != "" {
		return GoalReference{}, fmt.Errorf(
			"relationship update requires no active attempt; close %q first",
			goal.Status.ActiveAttemptID,
		)
	}
	targetIndex := -1
	for index := range catalog {
		if catalog[index].Metadata.Name == goal.Metadata.Name {
			targetIndex = index
			catalog[index] = v1alpha1.Goal(goal)
			break
		}
	}
	if targetIndex == -1 {
		return GoalReference{}, fmt.Errorf(
			"goal %q is not present in its catalog",
			goal.Metadata.Name,
		)
	}
	currentAnalysis, err := v1alpha1.AnalyzeGoalGraph(catalog)
	if err != nil {
		return GoalReference{}, fmt.Errorf(
			"analyze current goal graph: %w",
			err,
		)
	}

	updated := v1alpha1.Goal(goal)
	relationships := updated.Spec.Relationships
	if options.ClearParent {
		relationships.ParentGoalRef = nil
	} else if options.ParentGoal != "" {
		relationships.ParentGoalRef = &v1alpha1.GoalReference{
			Name: options.ParentGoal,
		}
	}
	relationships.DependsOnGoalRefs = goalReferences(options.DependsOn)
	relationships.SupersedesGoalRefs = goalReferences(options.Supersedes)
	updated.Spec.Relationships = relationships.Normalized()
	if updated.Metadata.Generation == ^uint64(0) {
		return GoalReference{}, fmt.Errorf("cannot increment generation")
	}
	resourceVersion, err := incrementResourceVersion(
		updated.Metadata.ResourceVersion,
	)
	if err != nil {
		return GoalReference{}, err
	}
	updated.Metadata.ResourceVersion = resourceVersion
	updated.Metadata.Generation++
	updatedManifest := GoalManifest(updated)
	if err := validateProspectiveRecord(
		updatedManifest,
		criteria,
		attempts,
	); err != nil {
		return GoalReference{}, fmt.Errorf(
			"invalid relationship update: %w",
			err,
		)
	}

	catalog[targetIndex] = updated
	analysis, err := v1alpha1.AnalyzeGoalGraph(catalog)
	if err != nil {
		return GoalReference{}, fmt.Errorf(
			"analyze prospective goal graph: %w",
			err,
		)
	}
	if err := rejectNewOrWorsenedCycles(
		updated.Metadata.Name,
		currentAnalysis,
		analysis,
	); err != nil {
		return GoalReference{}, err
	}
	updatedBytes, err := marshalYAML(GoalManifest(updatedManifest))
	if err != nil {
		return GoalReference{}, err
	}
	reference := GoalReference{
		GoalID:          updated.Metadata.Name,
		GoalRef:         ".",
		ResourceVersion: updated.Metadata.ResourceVersion,
	}
	intent, err := s.beginPublication(
		dir,
		updated.Metadata.Name,
		goal.Metadata.ResourceVersion,
		updated.Metadata.ResourceVersion,
		[]publicationFileEntry{
			{Path: "goal.yaml", BeforeDigest: "", Content: updatedBytes},
		},
		nil,
	)
	if err != nil {
		return GoalReference{}, err
	}
	if err := s.publishIntentFiles(dir, intent, ""); err != nil {
		if !s.goalCommitted(dir, intent) {
			_ = s.finishPublication(dir)
			return GoalReference{}, err
		}
		return reference, &PublicationIncompleteError{
			OperationID:      intent.Spec.OperationID,
			IntendedRevision: reference.ResourceVersion,
			Phase:            "publish",
			Kind:             "relationships",
			Message:          err.Error(),
			Cause:            err,
		}
	}
	if err := s.refreshREADMEProjection(dir); err != nil {
		return reference, &PublicationIncompleteError{
			OperationID:      intent.Spec.OperationID,
			IntendedRevision: reference.ResourceVersion,
			Phase:            "projection",
			Kind:             "relationships",
			Message:          err.Error(),
			Cause:            err,
		}
	}
	if err := s.finishPublication(dir); err != nil {
		return reference, &PublicationIncompleteError{
			OperationID:      intent.Spec.OperationID,
			IntendedRevision: reference.ResourceVersion,
			Phase:            "finish",
			Kind:             "relationships",
			Message:          err.Error(),
			Cause:            err,
		}
	}
	return reference, nil
}

func (s *Store) loadGoalCatalog(root string) ([]v1alpha1.Goal, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read goals root: %w", err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		names = append(names, entry.Name())
		if len(names) > maxGoals {
			return nil, fmt.Errorf("goal cardinality exceeds %d", maxGoals)
		}
	}
	sort.Strings(names)
	goals := make([]v1alpha1.Goal, 0, len(names))
	for _, name := range names {
		dir := filepath.Join(root, name)
		lock, err := s.acquireGoalLock(dir)
		if err != nil {
			return nil, fmt.Errorf("goal %q: %w", name, err)
		}
		goal, _, _, readErr := s.loadAndValidate(dir)
		releaseErr := lock.release()
		if err := errors.Join(readErr, releaseErr); err != nil {
			return nil, fmt.Errorf("goal %q: %w", name, err)
		}
		if goal.Metadata.Name != name {
			return nil, fmt.Errorf(
				"goal directory %q does not match metadata.name",
				name,
			)
		}
		goals = append(goals, v1alpha1.Goal(goal))
	}
	return goals, nil
}

func rejectNewOrWorsenedCycles(
	goalName string,
	current v1alpha1.GoalGraphAnalysis,
	prospective v1alpha1.GoalGraphAnalysis,
) error {
	// SetRelationships changes only edges leaving goalName. Therefore every
	// cyclic component created or expanded by the update must contain that
	// Goal. Components are compared within their relationship kind because
	// parent, dependency, and supersession edges are separate graphs.
	for _, relationship := range []v1alpha1.GoalGraphRelationship{
		v1alpha1.GoalGraphRelationshipParent,
		v1alpha1.GoalGraphRelationshipDependency,
		v1alpha1.GoalGraphRelationshipSupersession,
	} {
		currentCycle, wasCyclic := goalCycleComponent(
			current,
			relationship,
			goalName,
		)
		prospectiveCycle, isCyclic := goalCycleComponent(
			prospective,
			relationship,
			goalName,
		)
		if !isCyclic {
			continue
		}
		if !wasCyclic {
			return fmt.Errorf(
				"prospective goal catalog snapshot introduces a %s cycle involving goal %q",
				relationship,
				goalName,
			)
		}
		if !goalCycleIsSubset(prospectiveCycle, currentCycle) {
			return fmt.Errorf(
				"prospective goal catalog snapshot worsens a %s cycle involving goal %q",
				relationship,
				goalName,
			)
		}
	}
	return nil
}

func goalCycleComponent(
	analysis v1alpha1.GoalGraphAnalysis,
	relationship v1alpha1.GoalGraphRelationship,
	goalName string,
) (v1alpha1.GoalGraphCycle, bool) {
	for _, relation := range analysis.Relationships {
		if relation.Relationship != relationship {
			continue
		}
		for _, cycle := range relation.CyclicComponents {
			for _, reference := range cycle.GoalRefs {
				if reference.Name == goalName {
					return cycle, true
				}
			}
		}
		return v1alpha1.GoalGraphCycle{}, false
	}
	return v1alpha1.GoalGraphCycle{}, false
}

func goalCycleIsSubset(
	subset v1alpha1.GoalGraphCycle,
	superset v1alpha1.GoalGraphCycle,
) bool {
	members := make(map[string]struct{}, len(superset.GoalRefs))
	for _, reference := range superset.GoalRefs {
		members[reference.Name] = struct{}{}
	}
	for _, reference := range subset.GoalRefs {
		if _, ok := members[reference.Name]; !ok {
			return false
		}
	}
	return true
}

func goalReferences(names []string) []v1alpha1.GoalReference {
	references := make([]v1alpha1.GoalReference, 0, len(names))
	for _, name := range names {
		references = append(references, v1alpha1.GoalReference{Name: name})
	}
	return references
}
