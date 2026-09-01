package catalogv1alpha1

import (
	"fmt"
	"strings"
)

// RenderGoalMarkdown renders a human projection of the goal catalog from the
// same validated data. It states the JSON digest and is not a second
// authority.
func RenderGoalMarkdown(catalog GoalCatalog) string {
	var builder strings.Builder
	builder.WriteString("# Goal catalog\n\n")
	builder.WriteString("> Generated deterministic projection. The JSON document at ")
	builder.WriteString("`tools/agents/catalogs/goal.json` is authoritative.\n\n")
	fmt.Fprintf(&builder, "- ID: `%s`\n", catalog.ID)
	fmt.Fprintf(&builder, "- Schema: `%s`\n", catalog.Schema)
	fmt.Fprintf(&builder, "- Derivation: `%s`\n", catalog.DerivationVersion)
	fmt.Fprintf(&builder, "- Source revision: `%s`\n", catalog.SourceRevision)
	fmt.Fprintf(&builder, "- Completeness: `%s`\n", catalog.Completeness)
	fmt.Fprintf(&builder, "- JSON digest: `%s`\n", catalog.Digest)
	builder.WriteString("\n## Goals\n\n")
	for _, goal := range catalog.Goals {
		fmt.Fprintf(&builder, "- `%s`: %s\n", goal.CandidatePath, goal.Availability)
		if goal.Identity != nil {
			fmt.Fprintf(&builder, "  - Name: `%s`\n", goal.Identity.Name)
			fmt.Fprintf(&builder, "  - Owner root: `%s`\n", goal.Identity.OwnerRoot)
			fmt.Fprintf(&builder, "  - Scope: `%s`\n", goal.Identity.Scope)
		}
		if goal.CoarseStatus != nil {
			fmt.Fprintf(&builder, "  - Outcome: `%s`\n", goal.CoarseStatus.Outcome)
			fmt.Fprintf(&builder, "  - Execution: `%s`\n", goal.CoarseStatus.Execution)
		}
		if goal.Reason != "" {
			fmt.Fprintf(&builder, "  - Reason: %s\n", goal.Reason)
		}
	}
	builder.WriteString("\n## Limitations\n\n")
	if len(catalog.Limitations) == 0 {
		builder.WriteString("None.\n")
	} else {
		for _, limitation := range catalog.Limitations {
			fmt.Fprintf(&builder, "- %s\n", limitation)
		}
	}
	return builder.String()
}

// RenderIndexMarkdown renders a human projection of the agent system index
// from the same validated data. It states the JSON digest and is not a
// second authority.
func RenderIndexMarkdown(index AgentSystemIndex) string {
	var builder strings.Builder
	builder.WriteString("# Agent system index\n\n")
	builder.WriteString("> Generated deterministic projection. The JSON document at ")
	builder.WriteString("`tools/agents/catalogs/index.json` is authoritative.\n\n")
	fmt.Fprintf(&builder, "- ID: `%s`\n", index.ID)
	fmt.Fprintf(&builder, "- Schema: `%s`\n", index.Schema)
	fmt.Fprintf(&builder, "- Derivation: `%s`\n", index.DerivationVersion)
	fmt.Fprintf(&builder, "- Source revision: `%s`\n", index.SourceRevision)
	fmt.Fprintf(&builder, "- Completeness: `%s`\n", index.Completeness)
	fmt.Fprintf(&builder, "- JSON digest: `%s`\n", index.Digest)
	builder.WriteString("\n## Catalogs\n\n")
	for _, catalog := range index.Catalogs {
		fmt.Fprintf(&builder, "- `%s` (%s): %s, digest `%s`\n",
			catalog.ID, catalog.Kind, catalog.Completeness, catalog.Digest)
		for _, route := range catalog.QueryRoutes {
			fmt.Fprintf(&builder, "  - route: `%s`\n", route)
		}
	}
	builder.WriteString("\n## Conflicts\n\n")
	if len(index.CatalogEnvelope.Conflicts) == 0 {
		builder.WriteString("None.\n")
	} else {
		for _, conflict := range index.CatalogEnvelope.Conflicts {
			fmt.Fprintf(&builder, "- `%s` (%s): %s\n", conflict.ID,
				conflict.Code,
				strings.Join(conflict.SourcePaths, ", "))
		}
	}
	builder.WriteString("\n## Limitations\n\n")
	if len(index.Limitations) == 0 {
		builder.WriteString("None.\n")
	} else {
		for _, limitation := range index.Limitations {
			fmt.Fprintf(&builder, "- %s\n", limitation)
		}
	}
	return builder.String()
}
