## Context

See proposal.md. `tools/bzl/rules/al_bzl_target_doc.bzl` already captures
`streamed_jsonproto` using `native.genquery`. Existing site Mermaid macros
render generated source as SVG or WebP.

## Goals / Non-Goals

Keep capture, conversion, and rendering independently composable. Use the
unconfigured target graph: `select()` includes possible branches, not a
configured action graph. No new dependency, nested module, or deployment.

## Decisions

Proceed with a macro over `genquery` and a sandboxed conversion action.
An aspect would recreate the graph schema; recursive Bazel would conflict
with build locking and hermeticity. Reuse the existing native capture pattern.
Include implicit build inputs, exclude execution tools and non-dependency
attributes such as visibility, and retain the unconfigured input graph.
Use standard-library JSON decoding and existing Cobra for the command.

Keep parsing, graph selection, and Mermaid emission in an internal package;
the command owns filesystem and argument handling. Sort nodes and deduplicate
edges. Traverse before display filtering, retain an explicitly requested root,
and document omitted nodes. Report the selected depth and filters in comments.

Provide opt-in `external_leaves` in the rule and command. It includes the first
external target on each dependency path, stops traversal there, and suppresses
all of its outgoing edges, including to independently selected targets. Existing
external visibility and expansion defaults remain unchanged. Tool classification
belongs to Bazel rule attributes rather than converter naming heuristics.

Provide opt-in `group_by_package` for presentation. Full package labels title
subgraphs, short target names label nodes, and external repository groups
contain their fully qualified package groups. Keep generated node identifiers
global and derive sorted group identifiers independently from target captions.
Grouping preserves the selected graph and all edge identities; flat output
remains the default.

Package grouping selects the already pinned ELK layout. Native Dagre renders
of the nested fixture crossed group captions in both LR and TB directions;
the larger site preview also produced excessive empty space. The ELK fixture
keeps captions clear with the same nodes, edges, site palette, and fonts.
Decision: proceed with ELK only for grouped output. No new dependency or shared
theme change is needed. Revisit this choice if native layout regresses.

Place all visible external targets inside one `External` container, with
repository and package containers nested inside it when grouping is enabled.
Use ELK network-simplex placement and merged edge routes for grouped graphs.
On the site's 174-target external-boundary projection this uses 16% less canvas
area than default ELK with the same hierarchy. Selecting TB explicitly reduces
area by 39%, but routes edges through package title areas; retain LR as the
public default and recommended preview. These comparisons preserve all
174 targets, 173 dependency edges, 104 containers, and the existing fonts.

Offer `diagram_type = "mindmap"` / `--diagram-type=mindmap` as an explicit
ownership overview, alongside the default flowchart. Reuse graph selection,
escaping, sorted identifiers, and package ownership in both serializers.
Mindmap lists each selected target once under a build-input inventory root.
Ungrouped mode keeps full target captions; grouped mode nests short captions
under local packages or `External` / repository / full package branches.
The inventory title and documentation distinguish ownership from dependency
edges. Flow direction applies only to flowcharts. Select the already pinned
`tidy-tree` layout for mindmap, overriding the shared Dagre setting. The radial
force layout overlaps labels when rule kinds are included; tidy-tree keeps the
ownership branches separate on a larger canvas.

The user accepts large canvases and considers both flowchart and mindmap
usable; prioritize readable labels and clear structure over canvas reduction.
Treemap and block remain local experiments: their clipping or extreme height
adds no useful view. No new dependency or shared theme change is needed.

The site's graph targets analyze `:site` at a documented depth and remain companion
artifacts: consuming them from `:site` would create a dependency cycle.
Project docs embed an example graph of the converter, which has no dependency
on those docs. Site packaging uses existing shared layouts and palette.

## Risks / Trade-offs

- Large closures produce unreadable diagrams → expose depth and display filters.
- Missing records at a query boundary → omit dangling edges without inventing nodes.
- JSON metadata changes → decode only graph fields and tolerate extra metadata.
- Label syntax can affect Mermaid → generated identifiers and entity escaping.
- Unconfigured graph can overstate platform dependencies → state this limitation.
