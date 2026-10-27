# Verification

Task: add the Bazel dependency analyzer and its site-root WebP consumer.
Observations below are from 2026-10-27; exact final candidate checks and Git
publication state belong to the repository delivery receipts.
The external-leaf, package-grouping, mindmap, and Bazel 9 additions belong to
the current candidate. Publication was initially withheld and later explicitly
authorized; delivery receipts own its exact commit and remote state.

## Observed behavior

- The converter E2E passes shared dependencies, cycles, reordering, depth,
  generated files, groups, filters, malformed input, missing roots, conflicting
  records, and label escaping. Its real diamond fixture renders through native
  genquery and the maintained Mermaid WebP rule. Test outputs retain the Mermaid,
  image, and JSON report. No generated image is checked into Git.
- `//projects/alwaldend.com:deps_webp`, `:deps_dark_webp`, and `:deps_svg` build
  through native genquery against `:site`. The current depth-eight overview
  contains 36 targets and 37 dependency edges, reaching content and documentation
  owners while excluding execution tools and visibility metadata.
  Independent Bazel DOT agrees with the JSON input graph after omitting its
  toolchain-type declaration edges, which `ruleInput` does not represent.
  The dark WebP was visually inspected; labels and arrows are clear.
- `//projects/alwaldend.com:site_test` passes. Representative HTML contains the
  analyzer landing, reference documentation, project-index entry, and Bazel
  taxonomy entry. The reference page packages both SVG appearances and retains
  its command/rule examples. Local browser DOM inspection confirms page
  identity, documentation links, and successful loading of both SVG resources.
- The analyzer's OpenSpec check, semantic lint, repository quality suite, and
  four DNS checks passed during implementation. Final aggregate validation
  records the selected checks against the delivered candidate.

## Local external and package options

- `external_leaves` includes the first external targets and suppresses their
  traversal and outgoing edges. `group_by_package` uses full package captions,
  short target captions, and external repository containers. Both default to
  false; existing flat and fully expanded external modes remain available.
- Capture retains implicit non-tool inputs and excludes visibility metadata.
  The Hugo tools attribute now uses execution configuration, so the standard
  tool filter excludes Sass. The actual site and both Hugo action tests pass.
- The refreshed input graph contains 6,418 targets and 9,461 edges, including 60
  external boundary targets when files are shown. Grouped output preserves
  every identity and edge. The earlier capture matched independent DOT after
  excluding six toolchain-type declarations; its report identifies them. The
  new serializer source files add two local input nodes and edges.
- The dark preview keeps one shortest input path to every external boundary:
  174 targets and 173 edges, organized into 75 packages and 28 external
  repositories. This projection does not claim to show every local input path;
  complete flat and grouped Mermaid files remain under ignored task output.
- The E2E command verifies boundaries, independent shared children, root and
  depth cases, grouping identity, nested repositories, escaping, and ordering.
  A native dark fixture retains the grouped Mermaid and WebP as test artifacts.
- Local validation passed 31 root-workspace tests, all three standalone Hugo
  module tests, semantic lint of affected packages, and the new fixture's
  explicit buildifier check. Raw logs and the structured graph comparison are
  under `out/bazel_deps_analyzer/`.

## Local layout comparison

- All visible external targets now belong to one `External` container. The
  package-grouped preview has 104 containers: that parent, 28 repositories, and
  75 packages. E2E covers both flat and grouped membership and empty-container
  suppression. The native dark fixture has readable nested captions.
- At unchanged font sizes, the 174-target/173-edge preview measures
  8,394.5 × 9,686.5 with default ELK, and 7,010.5 × 9,686.5 with network-simplex
  placement and merged edge routes: 16.5% less canvas area. The grouped output
  now selects the latter. Flat layout and public option defaults are unchanged.
- The TB variant measures 15,480 × 3,212.5, reducing area by 38.8%, but its
  incoming paths intersect package title areas. It is shown as an alternative,
  not the recommended preview. Sampling the recommended LR paths every four
  SVG units found no intersection with the interiors of its 104 group captions.
  This geometric diagnostic supplements visual inspection of the real fixture.
- Radial mindmap (3,210 × 1,911) and treemap (2,396 × 772) prototypes retain the
  174 targets as package/repository inventories, including the External branch;
  they omit dependency relationships. The mindmap overlaps captions and the
  treemap clips long labels. An initial mindmap using the shared Dagre layout
  became a 42,663 × 464 strip. A block graph retains 173 dependency edges but
  becomes 2,165 × 161,474. Treemap and block remain local experiments. The
  accepted mindmap option now uses the larger tree layout described below.
- The local comparison gallery shows all nine attempts with a zoom viewer.
  Dark WebPs are retained for the compact LR/TB flowcharts, radial mindmap,
  and treemap. Mermaid, SVG, measurements, and renderer logs are preserved in
  `out/bazel_deps_analyzer/layout/`. Complete graph sources were regenerated
  with the current converter. These local experiments preceded the later publication authorization.

## Selectable mindmap

- Both rules expose `diagram_type`, and the command exposes `--diagram-type`.
  Flowchart remains the default; explicit flowchart output equals the default.
  Mindmap is a build-input inventory with ownership branches, unique selected
  targets, rule kinds, optional package grouping, and a common External branch.
  Selection and package ownership are shared between the two serializers.
- E2E cases were written before implementation and failed on the absent flag.
  The completed command and both native pipelines pass depth, file and external
  filters, external roots, shared/cyclic targets, stable record ordering,
  escaped labels, and invalid types without creating or overwriting output.
  Real query, grouped JSON, and hostile-label fixtures render as dark WebPs;
  diagrams, images, and the pipeline report are retained in test outputs.
- Native site companions reuse the existing depth-eight JSON capture. Their
  selected target sets independently match at 36 targets. SVG siblings and the
  dark mindmap WebP build through the existing site renderer; the main site
  still builds with no capture cycle.
- The current external-boundary preview has the same 174 targets and all 60
  external boundaries in both views. Flowchart retains its 173 dependency
  edges, and its Mermaid bytes are unchanged from the preceding preview.
  Complete sources retain all 6,418 targets in both formats. Captures, source
  hashes, conversion/render arguments, and scope are recorded in
  `out/bazel_deps_analyzer/mindmap/report.json`.
- The mindmap uses the already pinned tidy-tree layout. The radial layout
  crowded rule-kind captions; the selected tree measures 2,872 × 9,900 and has
  no overlaps among its 280 target/ownership node boxes at a one-pixel tolerance.
  The flowchart WebP measures 7,011 × 9,687. Large canvases are accepted; caption
  readability and clear ownership take priority over canvas reduction.
- The two dark views were shown with an interactive zoom/scroll viewer. Local
  validation passed 30 tests, including E2E, both owner OpenSpec checks, site
  build, and repository quality. Semantic lint and explicit formatting of the
  new Go files also pass. These observations precede the Bazel 9 upgrade;
  final delivery checks validate its exact candidate.

## Repaired integration defects

- `SITE-GRAPH-SCOPE-001`: native genquery loads the full unconfigured scope even
  for a bounded expression. The pinned uv importer assigned Linux source builds
  to macOS-only packages. The initial importer patch was replaced with uv's
  supported complementary macOS/non-macOS resolution environments. Generated
  package markers now express the platform boundary for the unmodified importer.
  The user approved the generated tomli 2.5.0 addition for Python 3.10. Existing
  package identities and distribution hashes are preserved.
- `SITE-CRATE-CACHE-001`: a clean isolated test reproduces Bazel 8.7's absolute
  Git strip-prefix links becoming invalid after cache reuse from a second
  output base. Bazel 9.3.1 emits relative links and passes the same test with
  identical manifest, Rust source, and both license files. Upgrade the pin and
  use the original crate_universe Git recipe. The temporary Anki archive,
  projections, templates, update targets, and formatting exception are removed.
  Cargo.lock and rules_rust remain unchanged. Probe records are under
  `out/bazel_deps_analyzer/anki-acquisition/`.
- Bazel 9 requires explicit Java provider imports in the existing transitive
  bazel_jar_jar dependency; select upstream 0.1.15 instead of 0.1.7. Refresh the
  generated Bazel preset and normal module metadata. The independent Python
  fixture established that Bazel 9 alone does not fix native Python import.
  The subsequent uv environment partition replaces the local importer patch.
- The site also lacked `//third_party/org_golang_google_grpc:docs`. Gazelle's
  docs language generated the standard packaging target without changing pins.

The shared dependency repair is maintained in
`infra/src/openspec/changes/repair-site-dependency-loading`. Raw logs, captures,
and the review WebP are under ignored `out/bazel_deps_analyzer/`.

## Session ergonomics

- `GROUP-LAYOUT-001`: Dagre crossed nested captions in both LR and TB fixture
  renders and wasted space in the large preview. The already pinned ELK layout
  fixes the fixture; select it only for grouped diagrams and retain flat layout.
- `PREVIEW-RUNFILES-001`: a configured action binary lacked interactive runfiles,
  and source font paths were absent from the execroot after cached actions.
  Use the Bazel renderer entry point, an owned font-subsetter launcher, and
  physical source paths for local previews. Failed attempts remain in task logs.
- `LAYOUT-CAPTION-001`: smaller canvas area alone did not establish a better
  layout. The TB comparison routes edges across title plates; retain LR despite
  its larger area. Mindmap/treemap compactness also omits dependency semantics.
- `MINDMAP-READABILITY-001`: the force layout overlaps captions when rule kinds
  are included. The pinned tidy-tree layout preserves ownership with no node
  box overlaps on the actual preview. Keep scope and relationship semantics
  explicit, and accept the larger canvas rather than dropping useful labels.
- `VALIDATION-LABEL-001`: an initial combined check named a nonexistent
  owner-local validation subpackage. Its two valid test targets passed, but the
  command failed target parsing. The declared analyzer validation target is
  `//projects/bazel_deps_analyzer/openspec:validate_test`; its separate rerun and
  26 repository quality checks passed.
- `DIAGRAM-SCOPE-001`: the initial depth-three view had seven targets and
  stopped at the source archive. Calling it the full graph overstated its
  scope. The site now exposes a depth-eight overview and a dark WebP target;
  documentation names the depth and omitted dependency categories explicitly.
- `DISCOVERY-OUTPUT-001`: early combined skill and README reads exceeded output
  budgets. Focused file ranges and JSON field extraction reduced later output.
- `GENQUERY-CLOSURE-001`: a successful bounded command-line capture did not
  establish native genquery acceptance. Retain the real target build as the gate.
- `CACHE-REFETCH-001`: force-fetching one repository restarted broader module
  inputs. Isolate the exact generated Git rule and test cache reuse directly;
  prefer the verified Bazel fix over a permanent crate-source override.
- `CACHE-PROBE-001`: the home bazelrc overrides workspace cache settings. Use
  command-line cache paths and assert the resolved cache location before
  interpreting a clean-cache probe. The invalid initial probe was preserved
  and its source paths restored; later probes used task-owned caches.
- `RENDER-PATH-001`: Node normalized `..` before a Bazel symlink was followed.
  Use declared render targets, or physical absolute paths for diagnostic runs.
- Native preview snapshots timed out, while DOM queries and local asset reads
  worked. Both SVG resources decoded after the diagnostic tab requested eager
  loading; product sources retain their normal lazy loading.
- No subagents, new software dependencies, host configuration changes, or live
  deployments were used. Maintained task records and delivery receipts provide
  the continuation point.

## Bazel 9 upgrade acceptance

- Bazel 9.3.1 builds the original locked Git crate and both native dark diagram
  targets. The analyzer E2E and patched Python reachability tests pass. The full
  site test passes after the initial cold protobuf build exceeded its deadline;
  the cached continuation completed in 91.6 seconds.
- All 13 standalone modules inherit the root Bazel pin and successfully refresh
  their module metadata through `bazel mod deps --lockfile_mode=update`. Their
  generated lockfiles migrate to format 28 without changing owner constraints.
  All 62 selected non-manual tests across these modules pass under the new pin.
- `VALIDATION-DEADLINE-001` (live): the initial 240-second site check expired
  during cold C++ compilation, without a failed assertion. Use a longer bounded
  deadline for the Bazel major-version migration; retain both logs.
- `DELIVERY-ENTRY-001` (live): the optional `al` command is unavailable. The
  documented task-local repo-delivery launcher works; refresh it under the new
  Bazel pin and retain that entry point through preparation and publication.
- The user authorized publishing the completed revision to PR 206. Exact final
  validation and remote state are recorded by repo-delivery; no live site
  deployment is requested.

## Native uv replacement acceptance

- Complementary `tool.uv.environments` entries for Darwin and its complement
  make package resolution markers sufficient for the unmodified rules_python
  2.4.2 importer. The patch, override, and five patch-specific files are removed.
  Marker-qualified overrides did not fix package records. Requirements-file
  compilation passed for Linux/macOS/Windows, but the native partition avoids
  extra exports and platform-specific wheel tag configuration.
- The normal `//projects/autoscroll:uv_lock.update` workflow preserves all 25
  existing package identities and distribution hashes. uv reforks Python marker
  ranges, selecting the existing PySide6 6.12.0 branch from Python 3.10. The user
  approved its tomli 2.5.0 dependency for Python 3.10. No requirements files are
  added. Repeated sandboxed lock generation is byte-identical to source.
- A native three-platform pip import probe yields all five PyObjC package
  aliases only for `cp314_osx_aarch64`; PySide6 retains Linux, macOS, and Windows
  aliases. The temporary import-platform list was restored exactly. Artifacts
  and commands are in `out/bazel_deps_analyzer/uv-platform-import/`, including
  `native-platform-report.json` and `repeat-lock-report.json`.
- Autoscroll, analyzer E2E, and the full-site test pass with the unmodified
  importer. Both native dark diagram targets build. This establishes Linux
  execution and cross-platform import metadata; macOS/Windows binaries were
  not executed on the Linux host. Exact final gates belong to delivery receipts.
- `UV-IMPORT-ALTERNATIVE-001`: the initial repair implemented generic graph
  reachability inside an external importer before testing supported uv
  resolution partitions. The public setting removes that maintenance burden.
  Test producer options against actual consumer metadata before carrying a
  private dependency patch.
- `DELIVERY-ENVIRONMENT-001`: per-command `CODEX_TOOL_CALL_ID` changes the
  delivery environment digest. Run exact-candidate validation and publication
  sequentially with one captured environment; keep the safety guard unchanged.
