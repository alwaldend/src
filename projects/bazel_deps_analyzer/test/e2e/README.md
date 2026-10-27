# Dependency diagram acceptance

Failure cases are defined before the implementation:

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
