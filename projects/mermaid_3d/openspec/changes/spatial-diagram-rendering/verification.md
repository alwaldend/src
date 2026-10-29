## Acceptance evidence

On 2026-10-29, the six-example browser test passed with remote requests blocked.
It checked spatial node volumes, group containment, independent route/node
intersection checks, light/dark theme and embedded font loading, selection,
orbit/pan/zoom/reset, mobile controls, and sequence event order. Its artifacts
are twelve screenshots and a JSON report with HTML SHA-256 identities.

The site build packages all six HTML files beside the existing SVG assets, and
the gallery contains all six links. Existing native raster, theme readability,
and site integration tests passed after the Mermaid upgrade. The other existing
renderer checks also passed during upgrade validation. Final candidate checks
and publication state belong to the repository delivery receipts. The delivery
validation plan passed all 34 selected tests, including repository quality, and
semantic lint over the 64 targets selected from affected packages.

Before publication, master advanced to `976eb5ffa`. The single task commit was
rebased; the only conflict was the project index table, resolved by retaining
both independently added project entries. The maintained Gazelle workflow
generated package policies for the new packages. Two exact tool-consumer
exceptions cover the renderer execution tool and the browser test; no product
runtime package receives a general tool dependency exception. The delivery
plan is rerun on the synchronized candidate.

The visual follow-up replaced front-face connections with side ports and direct
clear routes, with detours only where needed. The default projection also clears
unrelated node faces. Labels now use bounded projected collision handling and
leader lines, with endpoint cardinalities near their connections. Class/entity
cards distinguish headings and compartments; group volumes use the maintained
actor-line contrast.

All six samples passed expanded browser acceptance, including at least 12px
text at reset scale, projected label separation, full geometry framing in the
scrollable canvas and fitted overview, front view, phone/desktop resize, fresh
touch-device load, and native swipe scrolling. Each run saves full-scene and
page screenshots, including phone, scrolled, and fitted views. Both appearances
and all six full-scene outputs were visually inspected. The existing delivery
plan retains the SVG/raster/site consumer tests, repository quality, and semantic
lint; receipts identify the exact candidate and final gate results.

The readable view deliberately scrolls when the scene cannot fit at a useful
text size. Fit diagram supplies an overview, which can have smaller text. Dense
graphs and arbitrary camera angles can exceed the bounded label-placement area;
camera movement and selection remain available for inspection.

## Architecture documentation

The architecture source in `docs/architecture.mmd` builds light/dark SVGs,
a preview raster, and interactive HTML. Both the README and landing consume the
same packaged assets. The site build passed, and generated page references were
checked against their image, source, and HTML files. Render inspection prompted
removal of a group-title/edge crossing. The 3D view opened offline with seven
nodes, eight edges, one group, the embedded font, and no browser errors or remote
requests. Both appearances were captured. Delivery checks also cover the new
docs package, site validation, and the repository's DNS configuration test.

## Curved connections

The browser acceptance was extended before implementation and failed against the
straight-edge candidate (`class: visibly curved connections`). It then passed
all six samples with smooth sampled curves, no sharp direction changes, and the
existing independent route/node checks. Curve samples also drive tube geometry,
framing, and arrowhead direction; dashes follow cumulative path length.

All six full-scene outputs were inspected for curve shape, markers, dashes, and
sequence self-calls. The architecture view also opened in both appearances with
the embedded font, no browser errors, and no remote requests. Shared routes and
the monochrome palette remain intentional. Final delivery checks retain the
39-test suite, semantic lint, and the built site; receipts own exact results.

## Review corrections

The review reproduced transitive selection highlighting, arrowheads on plain
sequence lines, notes outside block bounds, rejected note-only blocks, and
reordered numeric participant IDs. Browser coverage was added before the fixes;
the new fixture initially failed with `Empty sequence blocks are not supported`.

After correction, the six maintained samples and the additional sequence fixture
passed browser acceptance. Actual material opacity checks cover direct edge and
node selection, nested block selection, and note selection. The fixture also
checks participant order, message markers, note membership, and complete group
containment. Its light and dark full-scene outputs were inspected. The final
delivery plan retains the 39-test suite, semantic lint, architecture renders,
and site build; delivery receipts own exact candidate results and publication.

## Rendering refinements

Browser acceptance was extended before implementation and initially failed its
offline size budget. All eight views then passed with measured text containment,
participant separation, note/message spacing, and compact phone overview checks.
The flowchart HTML decreased from 6,633,626 to approximately 3,446,000 bytes;
the sequence sample decreased from 545 to 80 draw calls. The JSON report records
actual payload sizes and draw calls for each input. Budgets are 4 MiB and 150
draw calls for these maintained fixtures, not claims about arbitrary diagrams.

Light and dark outputs were inspected for multiline cards, diamond faces,
continuous curves and dashes, and quieter rear group frames. The phone overview
keeps its existing scale while reducing excess panel height. The rebuilt
architecture view opened offline with its font, seven nodes, eight edges, and
no browser errors. Final checks retain the 39-test consumer suite, semantic lint,
architecture outputs, and site build. Delivery receipts own exact candidate
identities, base synchronization, and final results.

## Adapter correctness

Three additional browser fixtures cover decoded entities and generic class
headings, multiline decision containment, invisible placement links, local
subgraph direction, and explicit/shorthand nested activations. Acceptance was
extended before implementation; its first run reproduced the diamond defect.
All eleven views then passed. Corrected graph, class, and sequence screenshots
were inspected, including the maintained flowchart. The suite records a maximum
of 80 draw calls and less than 3.5 MB per offline view for this collection.

The delivery plan retains the 39-test consumer suite, semantic lint of all
affected packages, architecture renders, and the site build. Exact-candidate
results and publication are owned by delivery receipts; repeatable visual
artifacts and hashes are produced by the browser target.

## Element and layout boundaries

The review reproduced three defects through production rendering: shared node
and connection IDs lost a selectable object; graph self-loops crossed their own
front faces; long group titles escaped their frame and phone canvas. Three new
browser fixtures cover these conditions. Selection checks now distinguish
element kinds and inspect connection details as well as material opacity.
The route check independently tests the self-loop's front projection.

All fourteen views pass browser acceptance. Inspected light/dark scenes and
phone overviews confirm the corrected loop, multiline title, and sequence
identifier case. The final loop label sits above its curve. The maintained
flowchart, class, state, and sequence views retain their theme and framing.
The collection remains below 3.5 MB per HTML and at most 80 draw calls.
The delivery plan retains the 39-test consumer suite, semantic lint, architecture
outputs, and site build. Receipts own the rebased candidate and publication.

## Source structure

Build-time orchestration remains in `cmd/render`; model adaptation, sequence
handling, routing, and shared text layout live under `internal/`. The browser
controller, scene objects, labels, document, and styles live together under
`internal/viewer`. Rendering follows controller updates explicitly, removing
the redundant render after rebuilding the scene. One fixture list declares
browser inputs, and one input list supplies both test arguments and data.

All fourteen views passed the existing browser suite after the split. Their
model payloads equal the saved prior candidate. Of 113 screenshots, 106 match
byte for byte, including all 28 light/dark scene images. The seven differences
are touch-scroll offsets: their overlapping content matches pixel for pixel
after horizontal translation. The native swipe check permits momentum rather
than fixing a scroll position. No layout, theme, or interaction changes were
needed. Final validation retains the consumer suite, semantic lint including
both new internal packages, architecture outputs, and site build. Delivery
receipts own exact-candidate results and publication.

## Shared theme and static snapshots

Shared JSON theme data replaces browser recovery of CSS colors. All fourteen
existing model payloads and stable browser screenshots retain their prior
appearance. The website's generated CSS is byte-identical before and after
the theme migration.

The snapshot suite checks six production image targets across PNG/WebP,
light/dark appearances, front/spatial camera views, and multiple aspect ratios.
It compares decoded pixels against a fresh viewer canvas and independently
checks visible content, labels, opaque background, bounds, fonts, and camera
direction. Negative cases cover invalid options, missing fonts, external and
local-file dependencies, readiness timeout, and controls painted over the page.
A camera-distance case exercises clipping-plane fitting beyond the old fixed
far plane. Repeatable images and `verification.json` are emitted by the test.

All six sample diagrams expose light/dark PNG pairs. The project README and
landing embed a portrait architecture snapshot and retain the interactive link.
Website browser checks exercise both PNG variants through the real theme menu
and compare the image canvas with the page background. Both appearances were
inspected. The first integration check exposed a missing Chrome runfiles input;
the snapshot binary now declares it like the render binary. A test-side decode
race was corrected by waiting for lazy-loaded images before reading pixels.

Exact-candidate delivery retains the existing consumer suite, adds snapshot
acceptance and new package lint, and builds sample snapshots and the website.

## Dependency compatibility

- Mermaid 12.2.0 exposes shared graph data for five supported types. Sequence
  uses a separate event database and adapter, included at the user's request.
- CLI 12.0.0 replaces `myCSS` with licensed `customFontCSS`. Its local loader's
  `.tff` typo requires the maintained one-line `.ttf` correction.
- Mermaid's new 120px flowchart wrapping/minimum-width defaults changed existing
  layout. Explicit 200px wrapping and zero minimum width preserve prior themes.

## Session ergonomics

- `CLI12_FONT_INPUT`: inspecting the complete custom-font loading and embedding
  path before migration would have avoided repeated failing checks. Narrowing
  to the existing native font fixture made the final diagnosis and verification
  inexpensive. The patch README records the removal condition.
- `THREE_PACKAGE_LAYOUT`: the pinned release omits the assumed minified bundles
  and does not export `package.json`. Resolve its main entry and inspect the
  actual package files on upgrades; browser acceptance covers offline imports.
- `PREVIEW_LOCALHOST`: the shared preview returned `ERR_CONNECTION_REFUSED`
  twice. Repository browser screenshots provided repeatable visual evidence.
- No shared host, skill, or runtime changes were made. Owner-local OpenSpec,
  browser artifacts, and delivery receipts retain the continuation evidence.
- `PROJECTED_VISIBILITY`: DOM overflow and spatial collision checks alone missed
  canvas clipping and projected overlaps. Browser acceptance now checks actual
  projected scene/label bounds and exercises native touch scrolling. This closes
  the observed coverage gap without duplicating the layout algorithm in tests.
- `BOUNDED_WORKTREE_DISCOVERY`: a full worktree inventory produced about 19,000
  tokens of irrelevant output during the curve follow-up. The current checkout's
  `rev-parse` fields were sufficient; use those for future continuation. Existing
  workspace guidance already owns that rule, so no shared policy change is needed.
- `SELECTION_OBSERVABILITY`: checking only the selected identifier missed
  transitive highlighting. The browser fixture now verifies each rendered
  object's material opacity against independently stated selection expectations.
  This adds observable behavior to the existing inspection API and closes the
  demonstrated acceptance gap without introducing a separate test framework.
- `RENDERING_COST_OBSERVABILITY`: the visual review measured duplicated module
  bytes and hundreds of independent dash draw calls despite passing geometry
  checks. File-size and draw-call budgets now exercise those user-visible costs
  in the existing browser suite, with values retained in its report.
- `ADAPTER_SEMANTIC_COVERAGE`: valid parser inputs exposed seven display/layout
  defects beyond the original gallery assertions. Three end-to-end fixtures now
  cover the demonstrated conditions without introducing unit-level copies of
  layout calculations. The first new narrow fixture also exposed a horizontal
  scrolling assumption; the touch check now chooses the overflowing axis.
- `ELEMENT_BOUNDARY_COVERAGE`: follow-up review probes exposed identifier
  collisions, self-loop projection, and group-title bounds outside the gallery's
  original coverage. Three additional browser fixtures now exercise these
  conditions through the production pipeline. No additional test framework,
  dependency, or shared instruction changes were needed.
