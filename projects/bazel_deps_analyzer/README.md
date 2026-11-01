# Bazel dependency analyzer

Bazel rules and a Go command that analyze Bazel's streamed JSON dependency
graph as individual targets or aggregated packages. Target analysis produces
deterministic Mermaid flowcharts and ownership mindmaps. Package analysis
produces Mermaid or JSON from one shared aggregate. Generated diagrams compose
with the repository's existing SVG and WebP render rules.

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
| `build_deps_only`  | `False`       | Select explicit `deps`/`srcs` values using the capture contract below                      |
| `out`              | `<name>.mmd`  | Generated Mermaid output path                                                              |

## Command

Capture a graph and convert it independently:

```sh
al bazel query 'deps(//projects/alwaldend.com:site, 3)' \
  --output=streamed_jsonproto --implicit_deps --notool_deps --nonodep_deps \
  > "$PWD/out/site-dependencies.ndjson"

al bazel run //projects/bazel_deps_analyzer -- analyze target \
  --input "$PWD/out/site-dependencies.ndjson" \
  --output "$PWD/out/site-dependencies.mmd" \
  --root //projects/alwaldend.com:site --depth 3
```

The bare command and bare `analyze` print help and exit successfully. Choose
`analyze target` for the existing target conversion or `analyze package` for
aggregation. Conversion flags belong to the selected command; previous
root-level invocations must insert `analyze target`. The Bazel target and
JSON-input rules select that subcommand automatically.

For `analyze target`, `--input` and `--output` default to `-` for standard input
and standard output. `--root` is required; `--depth`, `--direction`, `--include-files`, and
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
rather than `aquery`'s JSON action graph. Extra record metadata is ignored by default.
Empty, malformed, conflicting, or rootless input fails before an output file is
created. Records and edges are sorted, and repeated edges are deduplicated.

## Package dependency graphs

`analyze package` constructs an overview from every record in one workspace's
capture, including disconnected components. It needs no root target or
synthetic selection node. Each owner is a graph node, rather than a container
of target nodes.

```sh
al bazel run //projects/bazel_deps_analyzer -- analyze package \
  --input "$PWD/out/all-targets.ndjson" \
  --output "$PWD/out/packages.mmd" --package-depth 1
```

`--package-depth` is a positive number of local package path segments and
defaults to one. For example, `//tools/compiler/frontend:library` belongs to
`//tools` at depth one, `//tools/compiler` at depth two, and
`//tools/compiler/frontend` at depth three or higher. Shorter package paths
remain intact. The main `//` package is captioned `root`; local `//root` and
`//External` remain distinct. All external repositories share one `External`
node at every depth. Package hierarchy depth is separate from dependency-hop
`--depth`, which belongs only to `analyze target`.

Files and external records are included by default. Set `--include-files=false`
or `--include-external=false` to exclude them and their incident edges. Nodes
report selected target counts and internal target-edge counts; internal edges
are not drawn as self-loops. Every distinct directed relationship between
selected groups is retained without edge captions. Reciprocal relationships
share a double-headed arrow by default; `--merge-reciprocal=false` keeps
separate arrows. Automatic ELK layout and edge routing are emitted by the
analyzer, with no manual node placement. `--direction` defaults to `LR` and
also accepts `RL`, `TB`, and `BT`.

Use `--format=json` to consume exactly the same constructed aggregate:

```sh
al bazel run //projects/bazel_deps_analyzer -- analyze package \
  --input "$PWD/out/all-targets.ndjson" \
  --output "$PWD/out/packages.json" --package-depth 2 --format json
```

The JSON contract is:

- `groups`: sorted by `caption`, each containing `caption`, `targetCount`,
  `internalEdges`, and sorted normalized target labels in `members`.
- `edges`: sorted by `source`, then `dependency`. These are zero-based group
  array indices; `targetEdges` counts the original distinct target edges.
  Reciprocal directions are both present regardless of `--merge-reciprocal`.
- `meta`: `packageDepth`, `selectedTargets`, `selectedEdges`,
  `unresolvedReferences`, `includeFiles`, and `includeExternal`. Selected edge
  counts include both internal and inter-group target edges. Opt-in filtering
  adds `buildDepsOnly: true`; that field is omitted by default.

Mermaid is the default format. JSON memberships let viewers use the analyzer's
construction directly rather than implement a second ownership algorithm.
Both formats report distinct unresolved dependency labels from selected source
records. Missing records are omitted, not invented. Aggregate paths do not
establish target-level reachability, and the capture's original query limits
still apply. Keep workspace captures separate: identical labels in separate
workspaces do not establish a shared target identity.

For a declared render pipeline, use the dedicated JSON-input rule:

```starlark
load("//projects/bazel_deps_analyzer:defs.bzl", "bazel_deps_packages_from_json")

bazel_deps_packages_from_json(
    name = "packages",
    src = "all-targets.ndjson",
    package_depth = 1,
)
```

This rule accepts `direction`, `include_files`, `include_external`,
`merge_reciprocal`, `build_deps_only`, and `format` with the command's defaults. Its output is
`<name>.mmd`, or `<name>.json` for `format = "json"`; `out` overrides that path.
Feed the Mermaid output directly to the existing SVG or WebP render rules.

## Explicit deps/srcs filtering

Both `analyze target` and `analyze package` accept `--build-deps-only`.
The capture macro and both JSON-input rules expose `build_deps_only = False`.
Enabling it keeps only explicitly specified dependency-label values from
attributes named `deps` or `srcs`, intersected with the captured `ruleInput`
labels. It excludes non-explicit, nodep, and aspect-owned attributes and
selector condition labels. Scalar/list labels, label dictionaries (values or
keys according to type), and every unconfigured selector alternative are
supported. String and output attributes contribute no dependencies.

This is a structural filter, not a complete semantic build classifier. It
omits ordinary `tools`, `data`, visibility, alias `actual`, `exports`, `headers`,
and custom dependency attributes even when they describe necessary build
inputs. A label also present in valid explicit `deps`/`srcs` remains selected.
Tool exclusion is performed by the required upstream `--notool_deps` query
flag, not inferred from attribute names, paths, or hidden target status.
Bazel's collapsed `ruleInput` cannot prove per-attribute tool provenance: for
example, a tool-transition `deps` label also present in non-tool `data` can
survive the intersection. The format has no per-attribute transition metadata
that would resolve this ambiguity.

Generated-file → generating-rule edges remain as structural build-producer
links, including when the generating rule is reached through a hidden file.
Package-group and environment-group metadata edges are removed. Filtering
happens before target traversal/depth and package aggregation. Package mode
still selects every captured record, including disconnected metadata, tool,
runtime, and source-file nodes; it filters edges rather than classifying nodes.
Existing file/external display options still apply. Filtered Mermaid contains
`%% build_deps_only=true`; filtered package JSON adds `meta.buildDepsOnly`.
Disabled filtering preserves existing output and accepts older captures.

Filtered input requires a new capture with all these options:

```sh
al bazel query 'deps(//projects/alwaldend.com:site, 3)' \
  --output=streamed_jsonproto --consistent_labels \
  --noimplicit_deps --notool_deps --nonodep_deps \
  --proto:output_rule_attrs=name,deps,srcs --proto:default_values=true \
  --proto:flatten_selects=false --proto:rule_inputs_and_outputs=true \
  > "$PWD/out/site-build-dependencies.ndjson"

al bazel run //projects/bazel_deps_analyzer -- analyze package \
  --input "$PWD/out/site-build-dependencies.ndjson" \
  --build-deps-only --format=json \
  --output "$PWD/out/site-build-packages.json"
```

Every RULE must contain an attribute named `name`, even rules without deps or
srcs. Native Bazel emits that marker with the options above. Old flat captures,
including the previously saved full root capture, lack the marker and fail with
exact recapture guidance before output replacement. The marker alone cannot
prove query options: callers must provide the complete requested metadata.
In particular, do not use `--proto:default_values=false`: it can omit an
explicitly supplied value equal to its default. Keep unflattened selectors to
retain scalar and dictionary alternatives. Omitted repeated value fields mean
empty collections; an ordinary LABEL without `stringValue` represents None.
Malformed or unsupported relevant encodings and conflicting records fail
before output replacement. The parser accepts native emitted JSON, not explicit
JSON null substitutes for omitted protobuf fields.

The capture macro sets these options automatically when `build_deps_only` is
true. With the default false it preserves its original query options, including
implicit inputs and omitted attributes. JSON-input rules cannot repair an old
capture; their caller supplies the required metadata and upstream filtering.

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
al bazel build //projects/alwaldend.com:deps_mermaid \
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
al bazel test //projects/bazel_deps_analyzer/test/e2e:e2e_test
```

The E2E command checks graph semantics and failure handling, then verifies a
real Bazel diamond graph rendered through Mermaid. It preserves the diagram,
WebP, and JSON report in its undeclared test outputs.
