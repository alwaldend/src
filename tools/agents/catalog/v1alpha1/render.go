package catalogv1alpha1

import (
	"fmt"
	"strings"
)

// RenderTopologyMarkdown renders a human projection of the topology catalog
// from the same validated data. It states the JSON digest and is not a second
// authority.
func RenderTopologyMarkdown(catalog TopologyCatalog) string {
	var builder strings.Builder
	builder.WriteString("# Topology catalog\n\n")
	builder.WriteString("> Generated deterministic projection. The JSON document at ")
	builder.WriteString("`tools/agents/catalogs/topology.json` is authoritative.\n\n")
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
	builder.WriteString("\n## Trees\n\n")
	for _, tree := range catalog.Trees {
		fmt.Fprintf(&builder, "- `%s` (%s): `%s` — %s\n",
			tree.ID, tree.Path, tree.ReadmePath, tree.Boundary)
	}
	builder.WriteString("\n## Components\n\n")
	for _, component := range catalog.Components {
		fmt.Fprintf(&builder, "- `%s` — %s (%s, %s); lifecycle `%s`\n",
			component.ID, component.Title, component.Path, component.DocsState,
			component.Lifecycle)
	}
	builder.WriteString("\n## Workspaces\n\n")
	for _, workspace := range catalog.Workspaces {
		fmt.Fprintf(&builder, "- `%s` — module `%s` at `%s`\n",
			workspace.ID, workspace.ModuleName, workspace.ModulePath)
	}
	return builder.String()
}
