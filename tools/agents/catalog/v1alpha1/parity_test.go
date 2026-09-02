package catalogv1alpha1

import (
	"strings"
	"testing"
)

func TestTopologyJSONMarkdownParity(t *testing.T) {
	catalog := TopologyCatalog{
		CatalogEnvelope: sampleEnvelope(),
		Trees: []TopologyTree{
			{
				ID: "projects", Path: "projects", ReadmePath: "projects/README.md",
				Boundary: TreeBoundaryProduct,
			},
		},
		Components: []TopologyComponent{
			{
				ID: "agents", Path: "projects/agents", OwnerReadme: "projects/agents/README.md",
				BuildPath: "projects/agents/BUILD.bazel", Title: "Agents",
				Description: "Repository-wide agent system",
				Lifecycle:   "active", DocsState: "owned",
			},
		},
		Workspaces: []TopologyWorkspace{
			{ID: "root", Path: ".", ModulePath: "MODULE.bazel", ModuleName: "com_alwaldend_src"},
		},
	}
	content, err := CanonicalJSONTopology(catalog)
	if err != nil {
		t.Fatalf("canonical encode: %v", err)
	}
	decoded, err := DecodeTopologyStrict(content)
	if err != nil {
		t.Fatalf("decode round-trip: %v", err)
	}
	markdown := RenderTopologyMarkdown(decoded)
	if !strings.Contains(markdown, "JSON digest: `"+decoded.Digest+"`") {
		t.Fatalf("markdown lacks the JSON digest: %s", markdown)
	}
	for _, expected := range []string{
		"`projects` (projects): `projects/README.md`",
		"`agents` — Agents (projects/agents, owned); lifecycle `active`",
		"`root` — module `com_alwaldend_src` at `MODULE.bazel`",
	} {
		if !strings.Contains(markdown, expected) {
			t.Fatalf("markdown missing %q:\n%s", expected, markdown)
		}
	}
}
