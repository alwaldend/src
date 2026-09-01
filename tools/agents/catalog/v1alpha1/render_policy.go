package catalogv1alpha1

import (
	"fmt"
	"strings"
)

// RenderPolicyMarkdown renders a human projection of the policy catalog from
// the same validated data. It states the JSON digest and is not a second
// authority.
func RenderPolicyMarkdown(catalog PolicyCatalog) string {
	var builder strings.Builder
	builder.WriteString("# Policy catalog\n\n")
	builder.WriteString("> Generated deterministic projection. The JSON document at ")
	builder.WriteString("`tools/agents/catalogs/policy.json` is authoritative.\n\n")
	fmt.Fprintf(&builder, "- ID: `%s`\n", catalog.ID)
	fmt.Fprintf(&builder, "- Schema: `%s`\n", catalog.Schema)
	fmt.Fprintf(&builder, "- Derivation: `%s`\n", catalog.DerivationVersion)
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
	builder.WriteString("\n## Conflicts\n\n")
	if len(catalog.Conflicts) == 0 {
		builder.WriteString("None.\n")
	} else {
		for _, conflict := range catalog.Conflicts {
			fmt.Fprintf(&builder, "- `%s` (%s): %s\n", conflict.ID, conflict.Code,
				strings.Join(conflict.SourcePaths, ", "))
		}
	}
	builder.WriteString("\n## Policies\n\n")
	for _, policy := range catalog.Policies {
		fmt.Fprintf(&builder, "- `%s` covers `%s` at precedence %d\n",
			policy.ID, policy.PathPrefix, policy.Precedence)
		for _, axis := range policy.Axes {
			fmt.Fprintf(&builder, "  - `%s`: %s (from `%s`)\n",
				axis.Name, axis.Value, axis.Source)
		}
	}
	return builder.String()
}
