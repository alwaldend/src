## Context

See proposal.md for motivation. The CLI already uses Cobra and the graph
package owns JSON parsing, label normalization, filtering, and Mermaid escaping.
The current full capture has roughly 287,000 targets and 390,000 edges.

## Goals / Non-Goals

The new command owns aggregation and readability choices. Rendering remains
with the existing Mermaid tool. The command consumes captures without launching
Bazel or introducing dependency thresholds. Existing target-level selection
remains unchanged unless explicit dependency filtering is enabled.

## Decisions

- Add `analyze target` and `analyze package`. Move all existing conversion
  flags to `analyze target`; the bare root and bare `analyze` print help.
  Reuse file handling through a converter callback. Existing Bazel rules pass
  `analyze target`, and a
  dedicated packages-from-JSON rule exposes the new mode to render consumers.
- Read one workspace capture with the existing parser. Include all records
  by default instead of requiring the task's former synthetic selection root.
- Use captions `root`, `//<local-path>`, and `External` to prevent collisions
  between the root, local directories, and external repositories. Truncate
  local package paths at a positive `--package-depth` (default one); retain
  shorter paths and keep external repositories together at every depth.
  Dependency-hop `--depth` remains an `analyze target`-only option.
- Aggregate selected nodes and dependency pairs in linear passes, then sort
  only the output identities and relationships. Count internal target edges
  in node captions; preserve cross-package direction without drawing self-loops.
- Let the tool emit `<-->` for reciprocal relationships and reuse the existing
  ELK/network-simplex/merged-route directive. Separate directed arrows remain
  available with `--merge-reciprocal=false`. Manual placement and threshold-based
  edge removal would weaken reproducibility or hide real relationships.
- Count distinct unresolved labels in the selected source records and describe
  capture and quotient-graph limits in the generated Mermaid comments.
- Construct one aggregate containing sorted groups, normalized member labels,
  and sorted directed group edges with original target-edge counts. Serialize
  that graph as Mermaid or JSON; reciprocal merging is Mermaid presentation
  only. The viewer consumes JSON memberships and does not reconstruct groups.

## Risks / Trade-offs

- Dense aggregate graphs still have crossings → verify a dark render from the
  real capture and preserve a separate-arrow mode for relationship inspection.
- Missing records prevent complete reachability → report unresolved references
  and retain the existing partial-capture disclosure with the published image.
- Large captures require memory for the parsed target graph → reuse the current
  representation and avoid additional per-target diagram strings.

## Validation failures to cover before implementation

Disconnected records, duplicate records and inputs, reciprocal edges, internal
edges, root and external namespace collisions, file/external filters, unknown
references, malformed labels and input, invalid flags, output preservation,
escaped captions, deterministic ordering, and unchanged existing command modes.
Also cover root help without stdin decoding, migration of target invocations,
depth one/two/beyond a path's length, invalid depths, and flag isolation.
JSON coverage checks membership, original-edge counts, equivalent expanded
Mermaid relationships, format validation, and independence from arrow merging.

## Explicit dependency filtering

The graph parser owns an opt-in structural filter before target traversal or
package aggregation. It selects only explicitly specified `deps` and `srcs`
label values, excludes nodep and aspect attributes and selector condition
labels, and intersects those values with captured `ruleInput`. Supported label
encodings include scalar/list labels, label dictionaries, and unflattened
selector alternatives. String and output attributes are not dependencies.
Generated-file producer links remain; package/environment metadata edges do
not. All-record package selection and source-file behavior remain unchanged.

Require a `name` attribute marker on every RULE to distinguish a rule without
deps/srcs from an old attribute-less capture. Required capture options are
`--noimplicit_deps --notool_deps --nonodep_deps
--proto:output_rule_attrs=name,deps,srcs --proto:default_values=true
--proto:flatten_selects=false --proto:rule_inputs_and_outputs=true`.
The marker cannot prove those options were used; callers must provide complete
metadata, including defaults and empty attributes. The capture macro sets the
options only for the opt-in mode. Tool exclusion belongs to Bazel's upstream
query flags. Flattened rule inputs cannot prove per-attribute tool provenance
when a label also occurs in a retained non-tool attribute such as data. No
attribute-name or hidden-target heuristic can fix that ambiguity. This mode
is not a complete semantic build classifier; actual, exports, headers, and
custom dependency attributes are deliberately omitted.

New preimplementation failure cases: explicit/default and nodep/aspect values;
tools/data/visibility and same-label overlaps; selector branches versus
condition labels, including scalar and dictionary alternatives; malformed and
unsupported selected encodings; duplicate and conflicting metadata; missing
name markers even on disconnected rules; old-capture rejection and output
preservation; alias omission; generated producers; metadata groups; filtering
before depth; normalized, deduplicated, deterministic JSON/Mermaid counts;
unchanged default output; both commands and native capture/JSON-input rules.
