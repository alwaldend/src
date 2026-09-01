package catalogv1alpha1

import (
	"fmt"
	"strings"
)

// RenderActionMarkdown renders a human projection of the action catalog from
// the same validated data. It states the JSON digest and is not a second
// authority.
func RenderActionMarkdown(catalog ActionCatalog) string {
	var builder strings.Builder
	builder.WriteString("# Action catalog\n\n")
	builder.WriteString("> Generated deterministic projection. The JSON document at ")
	builder.WriteString("`tools/agents/catalogs/action.json` is authoritative.\n\n")
	fmt.Fprintf(&builder, "- ID: `%s`\n", catalog.ID)
	fmt.Fprintf(&builder, "- Schema: `%s`\n", catalog.Schema)
	fmt.Fprintf(&builder, "- Derivation: `%s`\n", catalog.DerivationVersion)
	fmt.Fprintf(&builder, "- Source revision: `%s`\n", catalog.SourceRevision)
	fmt.Fprintf(&builder, "- Completeness: `%s`\n", catalog.Completeness)
	fmt.Fprintf(&builder, "- JSON digest: `%s`\n", catalog.Digest)
	builder.WriteString("\n## Providers\n\n")
	for _, provider := range catalog.Providers {
		fmt.Fprintf(&builder, "- `%s` (owned by `%s`) — %s\n",
			provider.ID, provider.Owner, provider.DefinitionPath)
	}
	builder.WriteString("\n## Actions\n\n")
	for _, action := range catalog.Actions {
		fmt.Fprintf(&builder, "- `%s` (%s.%s): %s\n",
			action.ID, action.ProviderRef, action.Selector, action.Classification)
	}
	builder.WriteString("\n## Aliases\n\n")
	if len(catalog.Aliases) == 0 {
		builder.WriteString("None.\n")
	} else {
		for _, alias := range catalog.Aliases {
			extra := ""
			if alias.ReplacementRef != "" {
				extra = " -> " + alias.ReplacementRef
			}
			fmt.Fprintf(&builder, "- `%s.%s` %s%s\n",
				alias.ProviderRef, alias.Selector, alias.State, extra)
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
