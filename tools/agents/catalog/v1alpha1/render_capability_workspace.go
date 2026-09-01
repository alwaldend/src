package catalogv1alpha1

import (
	"fmt"
	"strings"
)

// RenderCapabilityMarkdown renders a human projection of the capability
// catalog from the same validated data. It states the JSON digest and is not
// a second authority.
func RenderCapabilityMarkdown(catalog CapabilityCatalog) string {
	var builder strings.Builder
	builder.WriteString("# Capability catalog\n\n")
	builder.WriteString("> Generated deterministic projection. The JSON document at ")
	builder.WriteString("`tools/agents/catalogs/capability.json` is authoritative.\n\n")
	fmt.Fprintf(&builder, "- ID: `%s`\n", catalog.ID)
	fmt.Fprintf(&builder, "- Schema: `%s`\n", catalog.Schema)
	fmt.Fprintf(&builder, "- Derivation: `%s`\n", catalog.DerivationVersion)
	fmt.Fprintf(&builder, "- Producer: `%s`\n", catalog.ProducerRef)
	fmt.Fprintf(&builder, "- Source revision: `%s`\n", catalog.SourceRevision)
	fmt.Fprintf(&builder, "- Completeness: `%s`\n", catalog.Completeness)
	fmt.Fprintf(&builder, "- JSON digest: `%s`\n", catalog.Digest)
	builder.WriteString("\n## Limitations\n\n")
	if len(catalog.Limitations) == 0 {
		builder.WriteString("None.\n")
	} else {
		for _, limitation := range catalog.Limitations {
			fmt.Fprintf(&builder, "- %s\n", limitation)
		}
	}
	builder.WriteString("\n## Providers\n\n")
	for _, provider := range catalog.Providers {
		fmt.Fprintf(&builder, "- `%s` (%s, owned by `%s`) — %s\n",
			provider.ID, provider.Kind, provider.Owner, provider.SourcePath)
	}
	builder.WriteString("\n## Skills\n\n")
	for _, skill := range catalog.Skills {
		fmt.Fprintf(&builder, "- `%s` (owned by `%s`): layer `%s`, activation `%s`, cost `%s`\n",
			skill.ID, skill.Owner, skill.Layer, skill.Activation, skill.ContextCost)
		if len(skill.Exclusions) > 0 {
			fmt.Fprintf(&builder, "  - exclusions: %s\n", strings.Join(skill.Exclusions, ", "))
		}
		if len(skill.CapabilityRefs) > 0 {
			fmt.Fprintf(&builder, "  - capabilities: %s\n",
				strings.Join(skill.CapabilityRefs, ", "))
		}
		if len(skill.Dependencies) > 0 {
			fmt.Fprintf(&builder, "  - dependencies: %s\n",
				strings.Join(skill.Dependencies, ", "))
		}
	}
	return builder.String()
}

// RenderWorkspaceCheckMarkdown renders a human projection of the
// workspace-check catalog from the same validated data. It states the JSON
// digest and is not a second authority.
func RenderWorkspaceCheckMarkdown(catalog WorkspaceCheckCatalog) string {
	var builder strings.Builder
	builder.WriteString("# Workspace-check catalog\n\n")
	builder.WriteString("> Generated deterministic projection. The JSON document at ")
	builder.WriteString("`tools/agents/catalogs/workspace-check.json` is authoritative.\n\n")
	fmt.Fprintf(&builder, "- ID: `%s`\n", catalog.ID)
	fmt.Fprintf(&builder, "- Schema: `%s`\n", catalog.Schema)
	fmt.Fprintf(&builder, "- Derivation: `%s`\n", catalog.DerivationVersion)
	fmt.Fprintf(&builder, "- Producer: `%s`\n", catalog.ProducerRef)
	fmt.Fprintf(&builder, "- Source revision: `%s`\n", catalog.SourceRevision)
	fmt.Fprintf(&builder, "- Completeness: `%s`\n", catalog.Completeness)
	fmt.Fprintf(&builder, "- JSON digest: `%s`\n", catalog.Digest)
	builder.WriteString("\n## Limitations\n\n")
	if len(catalog.Limitations) == 0 {
		builder.WriteString("None.\n")
	} else {
		for _, limitation := range catalog.Limitations {
			fmt.Fprintf(&builder, "- %s\n", limitation)
		}
	}
	builder.WriteString("\n## Workspaces\n\n")
	for _, workspace := range catalog.Workspaces {
		fmt.Fprintf(&builder, "- `%s` — module `%s` at `%s`\n",
			workspace.ID, workspace.ModuleName, workspace.ModulePath)
		fmt.Fprintf(&builder, "  - bazelIgnore: %t, rootOverride: %t, docsAggregation: %t, fullCheck: %t\n",
			workspace.Projections.BazelIgnore,
			workspace.Projections.RootOverride,
			workspace.Projections.DocsAggregation,
			workspace.Projections.FullCheck)
		for _, phase := range workspace.Phases {
			fmt.Fprintf(&builder, "  - phase `%s` via `%s`: %s\n",
				phase.ID, phase.ProviderRef, phase.CommandTemplate)
		}
	}
	return builder.String()
}
