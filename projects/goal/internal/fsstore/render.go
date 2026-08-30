package fsstore

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

func (s *Store) Render(goalDir string, expectedResourceVersion string, limit int) error {
	dir, err := s.resolveInsideWorkspace(goalDir)
	if err != nil {
		return err
	}
	if _, err := parseResourceVersion(expectedResourceVersion); err != nil {
		return fmt.Errorf("expected resource version: %w", err)
	}
	limit, err = validateLimit(limit)
	if err != nil {
		return err
	}
	lock, err := s.acquireGoalLock(dir)
	if err != nil {
		return err
	}
	defer lock.release()
	goal, criteria, attempts, err := s.loadAndValidate(dir)
	if err != nil {
		return err
	}
	if goal.Metadata.ResourceVersion != expectedResourceVersion {
		return fmt.Errorf(
			"stale resourceVersion: expected %s, current is %s",
			expectedResourceVersion,
			goal.Metadata.ResourceVersion,
		)
	}
	content, err := renderREADME(goal, criteria, attempts, limit)
	if err != nil {
		return err
	}
	return s.atomicWrite(filepath.Join(dir, "README.md"), content, 0o644)
}

func renderREADME(
	goal GoalManifest,
	criteria CriteriaManifest,
	attempts []AttemptManifest,
	limit int,
) ([]byte, error) {
	limit, err := validateLimit(limit)
	if err != nil {
		return nil, err
	}
	var output bytes.Buffer
	fmt.Fprintf(&output, "# %s\n\n", markdownText(goal.Spec.Title, 200))
	output.WriteString(
		"> Generated bounded projection. Edit `goal.yaml`, `criteria.yaml`, and " +
			"attempt records only through the goal tool.\n\n",
	)
	fmt.Fprintf(&output, "- Goal ID: `%s`\n", goal.Metadata.Name)
	fmt.Fprintf(&output, "- API version: `%s`\n", goal.APIVersion)
	fmt.Fprintf(&output, "- Resource version: `%s`\n", goal.Metadata.ResourceVersion)
	fmt.Fprintf(&output, "- Generation: `%d`\n", goal.Metadata.Generation)
	fmt.Fprintf(&output, "- Lifecycle generation: `%d`\n", goal.Status.LifecycleGeneration)
	fmt.Fprintf(&output, "- Scope: `%s`\n", goal.Spec.Scope)
	fmt.Fprintf(&output, "- Outcome: `%s`\n", goal.Status.Outcome)
	fmt.Fprintf(&output, "- Execution: `%s`\n", goal.Status.Execution)
	fmt.Fprintf(&output, "- Active attempt: `%s`\n", valueOrDash(goal.Status.ActiveAttemptID))

	output.WriteString("\n## Relationships\n\n")
	relationships := goal.Spec.Relationships.Normalized()
	parent := "—"
	if relationships.ParentGoalRef != nil {
		parent = "`" + relationships.ParentGoalRef.Name + "`"
	}
	fmt.Fprintf(&output, "- Parent: %s\n", parent)
	fmt.Fprintf(
		&output,
		"- Depends on: %s\n",
		markdownGoalReferences(
			relationships.DependsOnGoalRefs,
			limit,
		),
	)
	fmt.Fprintf(
		&output,
		"- Supersedes: %s\n",
		markdownGoalReferences(
			relationships.SupersedesGoalRefs,
			limit,
		),
	)

	output.WriteString("\n## Acceptance criteria\n\n")
	criteriaLimit := min(limit, len(criteria.Spec.Items))
	for _, item := range criteria.Spec.Items[:criteriaLimit] {
		required := "optional"
		if item.Required {
			required = "required"
		}
		fmt.Fprintf(
			&output,
			"- `%s` (r%d, %s): %s — Evidence: %s\n",
			item.CriterionID,
			item.Revision,
			required,
			markdownText(item.Statement, 240),
			markdownText(item.EvidenceMethod, 240),
		)
	}
	if omitted := len(criteria.Spec.Items) - criteriaLimit; omitted > 0 {
		fmt.Fprintf(&output, "- … %d more criteria omitted by the projection limit.\n", omitted)
	}
	if len(criteria.Spec.Items) == 0 {
		output.WriteString("- No acceptance criteria recorded.\n")
	}

	output.WriteString("\n## Recent attempts\n\n")
	sorted := append([]AttemptManifest(nil), attempts...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Status.ObservedAt == sorted[j].Status.ObservedAt {
			return sorted[i].Metadata.Name > sorted[j].Metadata.Name
		}
		return sorted[i].Status.ObservedAt > sorted[j].Status.ObservedAt
	})
	attemptLimit := min(limit, len(sorted))
	for _, attempt := range sorted[:attemptLimit] {
		fmt.Fprintf(
			&output,
			"- [`%s`](attempts/%s/) — `%s`, `%s`, resource version `%s`, criteria r%d\n",
			attempt.Metadata.Name,
			attempt.Metadata.Name,
			attempt.Spec.WorkType,
			attempt.Status.State,
			attempt.Metadata.ResourceVersion,
			attempt.Spec.CriteriaRevision,
		)
	}
	if omitted := len(sorted) - attemptLimit; omitted > 0 {
		fmt.Fprintf(&output, "- … %d older attempts omitted by the projection limit.\n", omitted)
	}
	if len(sorted) == 0 {
		output.WriteString("- No attempts recorded.\n")
	}

	output.WriteString("\n## Record map\n\n")
	output.WriteString("- [`goal.yaml`](goal.yaml): machine-authoritative goal state\n")
	output.WriteString("- [`criteria.yaml`](criteria.yaml): versioned acceptance criteria\n")
	output.WriteString("- [`attempts/`](attempts/): isolated attempt records and evidence\n")
	return output.Bytes(), nil
}

func markdownText(value string, maximumRunes int) string {
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "`", "'")
	value = strings.TrimSpace(value)
	if utf8.RuneCountInString(value) <= maximumRunes {
		return value
	}
	runes := []rune(value)
	return strings.TrimSpace(string(runes[:maximumRunes-1])) + "…"
}

func valueOrDash(value string) string {
	if value == "" {
		return "—"
	}
	return value
}

func markdownGoalReferences(
	references []LocalGoalReference,
	limit int,
) string {
	if len(references) == 0 {
		return "—"
	}
	returned := min(limit, len(references))
	values := make([]string, 0, returned+1)
	for _, reference := range references[:returned] {
		values = append(values, "`"+reference.Name+"`")
	}
	if omitted := len(references) - returned; omitted > 0 {
		values = append(values, fmt.Sprintf("… %d omitted", omitted))
	}
	return strings.Join(values, ", ")
}
