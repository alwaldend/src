# Dependency diagram acceptance

Failure cases are defined before the implementation:

- Explicit deps/srcs filtering must reject attribute-less captures and malformed
  relevant metadata before output replacement, including disconnected rules.
  Native capture must emit the name marker and complete default/unflattened
  metadata. Explicit/default, nodep/aspect, dictionary and selector branches,
  tools/data/visibility, overlapping labels, alias omission, generated producer
  links, and package/environment metadata must follow the documented structural
  contract. Duplicate/conflicting records, filtering before depth, unchanged
  defaults, both CLI modes and Bazel rules, and deterministic JSON/Mermaid counts
  must be covered with retained query, diagram, aggregate, and report artifacts.
- Native non-label selector None branches can contain serialized type defaults
  such as `isDefaultValue: true` with `stringValue: ""`. They must contribute no
  edges and must not abort either mode, including when the rule is disconnected.

- Shared dependencies, duplicated edges, cycles, and reordered JSON must not
  duplicate nodes or change diagram bytes.
- Depth zero and one must retain the root and respect shortest-path distance.
- Files, generated-file producers, external labels, and missing boundary
  records must follow the documented filters without dangling edges.
- Quotes, HTML, Mermaid entity syntax, and directive-like text in labels must
  remain label content and render successfully.
- Empty, truncated, unsupported, conflicting, or rootless JSON and invalid
  CLI options must fail without creating a partial output.
- A real Bazel diamond graph must survive native query, Go conversion, and
  WebP rendering. Its diagram, image, and report are retained as test artifacts.
- Opt-in external leaves must retain the first external target on every local
  dependency path and suppress all outgoing edges from external targets, even
  when their children are independently visible. External-only descendants
  must disappear; shared local descendants reached without crossing the
  external boundary must remain. Depth, file filters, and an external root
  must still behave consistently. Defaults and full external expansion must
  remain unchanged.
- Native capture must retain explicit and implicit target inputs, exclude
  execution tools, and omit visibility metadata. Tool classification follows
  Bazel attribute transitions rather than attribute-name heuristics.
- Optional package grouping must retain every selected node and edge while
  using full package captions and short target captions. Identical target names
  in different packages must remain distinct. External packages must be nested
  under their own repositories, including root packages. Input order must not
  change output. Quoted and directive-like package or target names must remain
  escaped label content. Flat output must remain the default, and a real native
  render must verify nested group geometry.
- Every visible external target must be inside one `External` container in
  flat and package-grouped views. Local targets must remain outside it, and
  hidden external targets must not leave an empty container.
- Mindmap must be explicitly selected and identify itself as a build-input
  inventory. It must list every selected target once, including shared or
  cyclic dependencies, using the same depth, file, and external boundaries.
  Optional package branches must preserve full identities, short captions,
  and external repository nesting. Flat mode must remain available; hidden
  external targets must not leave an empty branch. Reordering must be stable.
  Hostile labels must remain escaped content in a real native dark render.
  An unknown diagram type must fail before creating or overwriting output.
  Default output must equal explicit flowchart output. Real target capture and
  JSON-input rules must both produce retained mindmap and WebP artifacts.
- Bare root and `analyze` invocation must show help successfully without
  consuming input; root flags must not invoke conversion. Existing target
  behavior belongs to `analyze target`; superseded `targets` and `packages`
  commands, unknown commands, and flags from the other mode must fail.
- The `analyze package` command must include disconnected captured targets without a
  root, aggregate the requested positive number of package segments, and
  distinguish `root`, `//root`, `//External`, `//external`, and `External`.
  Depth one is the default; larger depths retain shorter package paths and
  keep all external repositories together. Duplicate records and inputs must
  not inflate target or internal-edge counts. Generated files must retain their
  producer dependency.
- Reciprocal package arrows must expand to exactly the same directed pairs as
  separate arrows, without self-loops, edge captions, or added relationships.
  File and external filters must remove only selected targets and incident
  edges. Missing-reference counts must deduplicate unknown labels from selected
  sources, excluding known but filtered targets.
- Package output must preserve ordering under reordered records and inputs,
  escape hostile captions, report aggregate capture and reachability limits,
  and support every documented direction and file/standard-stream interface.
  Malformed or conflicting records, invalid ownership labels, unknown flags,
  unsupported directions, nonpositive package depths, invalid formats, and
  positional arguments must preserve existing output. The command fixture
  retains `packages.mmd`, its dark render, and a JSON verification report in
  undeclared test outputs.
- Package JSON must contain exactly the documented schema with sorted groups,
  normalized sorted members, directed weighted edges, and capture metadata.
  Its counts and relationships must agree with Mermaid at each package depth.
  Reciprocal merging must leave JSON unchanged. Native JSON and Mermaid rules
  must produce the same verified command results with their default extensions.
