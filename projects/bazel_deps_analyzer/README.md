# Bazel dependency analyzer

Bazel rules and a Go command that convert Bazel's streamed JSON dependency
graph into deterministic Mermaid flowcharts and ownership mindmaps. Generated diagrams compose with
the repository's existing SVG and WebP render rules.

## Bazel rules

Load the macro in a BUILD file and choose a root target:

```starlark
load("//projects/bazel_deps_analyzer:defs.bzl", "bazel_deps_mermaid")
load("//tools/mermaid:defs.bzl", "mermaid_webp")

bazel_deps_mermaid(
    name = "dependencies",
    target = ":application",
    depth = 3,
)

mermaid_webp(
    name = "dependencies_webp",
    src = ":dependencies",
    out = "dependencies.webp",
    aspect_ratio = [],
)
```

The macro captures `deps(target)` with native `genquery`, following the pattern
in `tools/bzl`, then runs the converter in a sandbox. It never launches Bazel
inside an action. The capture is available as `<name>.query.ndjson`; the main
target produces `<name>.mmd`, or the path selected with `out`.

`bazel_deps_mermaid_from_json` accepts an existing capture through `src`, with
the complete target label supplied as `root`. It takes the same conversion
options as the target macro:

| Attribute          | Default       | Behavior                                                                                   |
| ------------------ | ------------- | ------------------------------------------------------------------------------------------ |
| `depth`            | `-1`          | Maximum shortest-path dependency distance; `-1` is unlimited and `0` shows the root only   |
| `diagram_type`     | `"flowchart"` | Dependency `flowchart` or ownership `mindmap` inventory                                    |
| `direction`        | `"LR"`        | `LR`, `RL`, `TB`, or `BT` flowchart direction; ignored for mindmap                         |
| `include_files`    | `False`       | Include source files and generated files                                                   |
| `include_external` | `False`       | Include external repository targets                                                        |
| `external_leaves`  | `False`       | Show the first external targets as leaves and omit their dependencies                      |
| `group_by_package` | `False`       | Group short target names under full package names, nesting external packages by repository |
| `out`              | `<name>.mmd`  | Generated Mermaid output path                                                              |

## Command

Capture a graph and convert it independently:

```sh
bazel_agent bazel query 'deps(//projects/alwaldend.com:site, 3)' \
  --output=streamed_jsonproto --implicit_deps --notool_deps --nonodep_deps \
  > "$PWD/out/site-dependencies.ndjson"

bazel_agent bazel run //projects/bazel_deps_analyzer -- \
  --input "$PWD/out/site-dependencies.ndjson" \
  --output "$PWD/out/site-dependencies.mmd" \
  --root //projects/alwaldend.com:site --depth 3
```

`--input` and `--output` default to `-` for standard input and standard output.
`--root` is required; `--depth`, `--direction`, `--include-files`, and
`--include-external` select the same behavior as their rule attributes.
Use `--external-leaves` (or `external_leaves = True` in a rule) to show the
first external target on each path without traversing its dependencies.
This option includes external targets without requiring `--include-external`;
when both are set, the leaf boundary applies. Both options default to false,
and `--include-external` alone retains its full external expansion behavior.

Use `--group-by-package` (or `group_by_package = True` in a rule) to label
package subgraphs with their full names and nodes with short target names.
All visible external targets belong to an `External` subgraph, including in
flat mode. With package grouping, it contains repository groups, which contain
their fully qualified package groups. This changes the
diagram's presentation without changing its selected targets or edges.
Grouping defaults to false and composes with `external_leaves` and file filters.
Grouped diagrams select Mermaid's ELK layout with network-simplex placement
and merged edge routes; the repository renderer already includes it. Flat
diagrams keep their existing layout. Direction remains explicit; inspect
package captions when changing it because incoming edges can cross titles in
top-to-bottom layouts. Neither layout setting removes nodes or dependency edges.

The command consumes the stream of JSON `Target` records emitted by
[`--output=streamed_jsonproto`](https://bazel.build/versions/8.7.0/query/language#output-formats),
rather than `aquery`'s JSON action graph. Extra record metadata is ignored.
Empty, malformed, conflicting, or rootless input fails before an output file is
created. Records and edges are sorted, and repeated edges are deduplicated.

## Mindmap overview

Select `diagram_type = "mindmap"` in either rule, or `--diagram-type=mindmap`
on the command. Flowchart remains the default. Both views select the same
targets using the existing depth, file, and external options.

The mindmap is a **build-input inventory**: branches describe ownership and
list each selected target once. They do not represent dependency edges.
Use the flowchart to trace build relationships, shared inputs, or cycles.
With package grouping enabled, local packages and external repositories form
separate branches; external packages are nested under their repository within
`External`. Without grouping, target captions retain their full labels and
external targets still share the `External` branch. Rule kinds remain visible.
Mindmap uses the renderer's pinned tidy-tree layout. Flow direction applies only
to flowcharts. Large canvases are acceptable; prioritize readable captions and
clear ownership when selecting a view.

```starlark
bazel_deps_mermaid(
    name = "input_inventory",
    target = ":application",
    diagram_type = "mindmap",
    external_leaves = True,
    group_by_package = True,
    include_files = True,
)
```

## Graph semantics

In flowcharts, arrows point from a target to its dependencies. Ungrouped nodes
show full target labels; grouped nodes show short names beneath full package
captions. Both views retain rule kinds. Generated files depend on their generating rule. Labels use
generated identifiers and escaped text so target names cannot become Mermaid
syntax or HTML.

The target macro includes explicit and implicit build inputs while excluding
execution tools and non-dependency attributes such as visibility. This is Bazel's
unconfigured target graph: `select()` can include multiple possible branches.
It does not represent one platform's configured targets or action graph.

Depth measures actual graph edges before display filters. Filtering removes
nodes and their incident edges without inventing replacement edges, so a node
reached through a hidden file or external target can appear disconnected.
The explicitly requested root always remains visible. References to records
outside the input capture are omitted, including at a bounded query frontier.
With `external_leaves`, traversal and all outgoing edges stop at external
targets. A descendant can still appear when reached through a separate local
path. The depth and file filters continue to apply; use `include_files` to
display dependencies that are source or generated files.

## Generated example

This view is built from the converter command itself at depth two:

![Converter target dependencies](docs/assets/converter.svg)

## Alwaldend.com build graph

The site package exposes companion targets that analyze `:site` at depth eight:

```sh
bazel_agent bazel build //projects/alwaldend.com:deps_mermaid \
  //projects/alwaldend.com:deps_svg //projects/alwaldend.com:deps_webp \
  //projects/alwaldend.com:deps_dark_webp
```

Outputs are `dependencies.mmd`, light/dark `dependencies.svg` siblings,
light `dependencies.webp`, and dark `dependencies.dark.webp` under
`bazel-bin/projects/alwaldend.com/`. This overview reaches site content packages
and repository documentation owners while omitting files, external targets,
tool dependencies, and visibility metadata. It is not a complete transitive
graph. The images use the site's shared palette and renderer.
These outputs stay outside `:site`'s own
inputs, since embedding a graph into the target it queries creates a dependency
cycle. The separately generated converter example above is packaged into the
site's documentation.

The companion `:deps_mindmap`, `:deps_mindmap_svg`, and
`:deps_mindmap_dark_webp` targets reuse that depth-eight JSON capture and target
selection for a package ownership inventory. They produce
`dependencies.mindmap.mmd`, light/dark SVG siblings, and
`dependencies.mindmap.dark.webp`. As with the flowchart, these are build
artifacts outside the site's own input closure.

`genquery` loads the complete scope closure before evaluating a depth limit.
Every package in that closure must therefore load successfully, including
platform branches absent from the rendered diagram. If a foreign-platform
dependency cannot load on the host, use a bounded command-line query and the
JSON-input rule or command to inspect the available portion of the graph.

## Validation

```sh
bazel_agent bazel test //projects/bazel_deps_analyzer/test/e2e:e2e_test
```

The E2E command checks graph semantics and failure handling, then verifies a
real Bazel diamond graph rendered through Mermaid. It preserves the diagram,
WebP, and JSON report in its undeclared test outputs.
