---
title: Mermaid
description: Bazel-managed Mermaid diagram renderer
---

This package exposes the Mermaid CLI as a repository-wide Bazel tool. Bazel
provisions the pinned Node toolchain, JavaScript dependency graph, Mermaid CLI,
and Chrome-for-Testing browser used by Puppeteer. Rendering runs as a Bazel
target or action; it does not use a host browser, download one during an npm
lifecycle hook, or fetch scripts or fonts from a CDN.

[`theme.json`](theme.json) owns the diagram appearance. It applies a
documentation palette, typography, layout, and the `handDrawn` look to every
diagram, so rendered work shares one visual contract instead of Mermaid's
built-in default. Shapes stay semantic: `()` and other inherently rounded
shapes render rounded, while `[]` boxes render as sketched rectangles. Nothing
forces every node into a rounded outline.

The `mermaid_svg` rule and the `mmdc` target both apply the theme by default.
A diagram can override any part of it with its own Mermaid directive, which
Mermaid merges over the file config:

```text
%%{init: {"look": "classic"}}%%
flowchart LR
    a[Square] --> b(Rounded)
```

Node label paragraphs have `2px` of horizontal padding. The pinned handwriting
font can paint beyond a string's measured advance: `openai.com` measures 78px
but its final glyph reaches 79px. Mermaid clips HTML labels to their measured
`foreignObject`, so the padding reserves room for that overhang before layout.
Label `foreignObject` elements also use `overflow: visible`: if an SVG viewer
falls back to a wider font, its glyphs can extend into the node's existing
padding instead of being cut off at the narrower measured frame. A fallback
Arial render reproduces right-edge clipping on the hostname without this rule.
The rule selects `.label > *` because Mermaid lowercases element selectors;
`foreignobject` would not match the case-sensitive SVG `foreignObject` name.

`themeCSS` boxes every edge label and every cluster title, so an annotation
reads as a caption instead of sitting on the edge line. The rule styles the
`<p>` inside a real label and leaves `.labelBkg` transparent; styling
`.labelBkg` directly would draw a stray box wherever an edge carries no text.
Mermaid emits an empty label plate for unlabelled edges under the dagre
renderer but omits it under ELK, so the rule is what the test asserts rather
than any particular empty element.

Each boxed label carries a `3px` margin. Under the dagre renderer Mermaid sizes
a label's `<p>` to the text only, so the frame clips the outer half of the 1px
border and a label appears to have a missing right edge; the margin keeps the
border inside the clip rect. Under ELK the frame already includes the margin,
so the value there is breathing room rather than a correctness fix.

The cluster title plate pins `font-size: 13px`. ELK sizes a container from the
width of its own title plate, so a child container whose title is wider than
its parent's escapes the parent on both sides. Nesting a hostname title inside
a short parent title is the case that exposes this: at the inherited `16px`
the child measured 333px against a 302px parent and escaped 20px per side,
while `13px` fits it with 7px of clearance per side. The plate is a container
label, not body text, so the smaller size also keeps it visually subordinate.

`flowchart.defaultRenderer` selects `dagre-wrapper`, Mermaid's Dagre renderer.
Diagrams use left-to-right flow by default. Dagre's `nodeSpacing` and
`rankSpacing` control spacing between nodes and ranks, while the title styling
below keeps nested containers readable.

`flowchart.subGraphTitleMargin.top` is negative on purpose. Under dagre a
nested cluster keeps a fixed 20px inset from its parent that does not grow with
`nodeSpacing`, `rankSpacing`, `padding`, or this margin, so the parent title
would be overlapped by the container nested inside it, and the negative margin
is the only thing that lifts it clear. ELK does not need the lift, but both
renderers share this one config, so the value stays and keeps the title on the
cluster border. Remove it only after re-checking a diagram with nested
subgraphs.

The cluster title additionally forces `line-height: 1`, which keeps the plate
short enough to sit on the border consistently instead of spilling past it.

`handDrawnSeed` is fixed, so the sketch jitter is stable rather than different
on every render. Pass `theme` on the rule or a later `-c` argument to the CLI
to replace the whole config for one diagram; because a second `-c` file
replaces rather than merges, override only the specific keys you need with a
directive.

## Paint order

Mermaid emits a container title in `<g class="cluster-label">` inside
`<g class="clusters">`, and emits that group before the edges, the edge labels,
and the nodes. SVG paints in document order and offers no other layer control:
`z-index` needs a CSS-positioned box and does not reorder SVG elements, and
`paint-order` only orders fill, stroke, and markers within a single element.
An edge that crosses a container border therefore paints over the title plate
and appears to run through the title text. Mermaid exposes no z-order, layer,
or paint-order configuration key, so this is fixed by rewriting the SVG rather
than by a theme setting.

The current theme already does everything CSS can do here, so the remaining
crossing is the one case CSS cannot reach. The title plate is the topmost paint
at its own plate coordinates, and on the host_bot diagram one edge crosses the
`External services` plate and two cross the `LLM BM` plate. A rasterised check
draws a white plate and then a red line over it and reads the pixel at the
crossing: reordering the document turns that pixel white, while `z-index`,
`z-index` with `position:relative` or `position:absolute`, `isolation:isolate`,
`paint-order`, `transform:translateZ(0)`, `will-change:transform`, a
`filter:drop-shadow` stacking context, and `opacity:0.999` all leave it red.
Only document order moves an SVG group above a later sibling, which is why the
lift is a rewrite and not a rule.

`cmd/render/paint_order.mjs` moves every title group to the end of the outer
root group after Mermaid writes the file. Dagre places nested containers in
translated groups; the lift composes their ancestor transforms onto each title
so its position stays unchanged. Empty, self-closing paint layers do not affect
the group scan. A browser test checks both title coordinates and visible
stacking. Both entry points apply the lift, so a diagram rendered through
`mermaid_svg` and the same diagram rendered through `mmdc` paint identically.
`mmdc` rewrites the file named by `-o` or `--output`; an SVG printed to stdout
is not rewritten.

## Diagram direction

Write `flowchart LR` for left-to-right flow. Direction belongs to the diagram
source; the shared theme selects Dagre but does not override an explicit
diagram direction. Use another direction when the diagram or the user requires
it.

Subgraphs can declare `direction TB`, but Mermaid ignores that direction when
their internal nodes connect outside the group. Both Dagre and ELK exhibit
this limitation on the T3 Code graph, so verify the rendered arrangement rather
than assuming a direction declaration took effect.

## Actors

An actor is an ordinary rounded node marked with the `actor` class. The theme
owns the figure: a `.actor div` rule paints a hand-drawn stick figure as a
`background-image` and reserves room with `padding-left`. A diagram therefore
carries only its label text:

```
remote("Actor")

classDef actor font-style:normal;
class remote,intranet actor
```

The `classDef` is required even though the theme supplies the appearance.
Mermaid emits a node's class attribute only when a `classDef` declares it, so a
bare `class` statement leaves the node unstyled and the figure does not paint.

Keeping the figure in the theme rather than in the diagram means the ~1.2 KB of
base64 SVG is written once, every actor node shares it, and the diagram source
states topology only. Sizing is fixed by the theme, so all actors match: the
figure plus the renderer's padding lands an actor node at the same height as a
single-line service node.

The rule rides in every rendered document, including diagrams with no actor
node, where it adds about 1.3 KB to the inline stylesheet. That is the cost of
owning appearance centrally; it is acceptable because the rule is inert without
an `actor` node and every diagram already carries a much larger stylesheet. The wobble is baked into the path coordinates with a
fixed seed rather than applied through an SVG filter, so the shared figure
carries no `<defs>` id that would collide when the background repeats.

## Branching edges

Use a visible junction node when several outgoing edges share one label:

```
codex ---|"Inference"| inference@{ shape: junction }
inference --> openai
inference --> openrouter
```

The incoming segment carries the shared label without an arrowhead, and the
outgoing segments point to their respective targets. The junction is a real
flowchart node and affects layout, so inspect the rendered spacing. This is
the form used by the maintained T3 Code diagram. The separate
`architecture-beta` diagram type has its own junction declaration syntax.

## Layout width

Dagre uses `nodeSpacing` between nodes in a rank and `rankSpacing` between
ranks. For an LR diagram, rank spacing contributes to width and node spacing
contributes to height. These settings do not impose a maximum canvas width;
labels, topology, and nested groups also affect the result. Check the rendered
diagram at its intended display size before adjusting spacing or direction.

## Plain renders

`mermaid_svg(..., plain = True)` and `MERMAID_PLAIN=1` render with Mermaid's own
defaults instead of the repository appearance: no theme, no paint-order
post-processing, and no pinned fonts. The output is deliberately neither themed
nor hermetic, and exists to compare a diagram against the upstream appearance.
A plain action declares neither the theme nor the fonts as inputs, so it cannot
depend on them by accident; `//users/simeonwarren/host_bot:rendered_plain` is a
worked example. Plain output is not published and not committed.

[`upstream-theme.json`](upstream-theme.json) is the themed counterpart to that
opt-out. It selects Mermaid's built-in `default` theme and pins the label family
to a font the renderer already carries, so the result is hermetic like any other
render and shows the upstream appearance rather than a host-dependent one. Use
it when a consumer has to show the contrast between Mermaid's own look and the
repository theme, as the blog post on this work does:

```starlark
mermaid_svg(
    name = "example_default",
    src = "example.mmd",
    out = "rendered/example-default.svg",
    theme = "//tools/mermaid:upstream-theme.json",
)
```

This is a comparison input, not a second appearance the repository maintains. It
exists so a before-and-after pair renders from one source through the same
pipeline; a diagram's own look still comes from `theme.json`.

## Native site diagrams

The main site's `mermaid_site_svg` macro lives in
[`projects/alwaldend.com/pkg/mermaid/defs.bzl`](../../projects/alwaldend.com/pkg/mermaid/defs.bzl).
It selects `mermaid_svg(native = True)` with the site's own configuration,
compiled palette, and pinned Liberation Sans font. It does not inherit the
historical theme or actor decoration. Its own preset restores rounded caption
plates using Mermaid's `themeCSS` and flowchart settings. The site's compact
preset uses basis curves and `flowchart.edgeSpacing`; the
[pinned dependency patch](../../third_party/com_github_mermaid_js_mermaid/README.md)
passes the latter through to Dagre before layout, including nested subgraphs.

Native mode requires a `palette` CSS input and `dark_out` SVG output. The
stylesheet exposes `--mermaid-*` custom properties named after Mermaid theme
variables. Chrome resolves them in light and dark media, and the renderer
passes the resulting values to Mermaid's native configuration and exposes the
same `--mermaid-*` properties on its diagram root for input `themeCSS`. Font embedding
uses the CLI's CSS input. The returned SVG bytes are written directly: there
is no color rewriting, paint-order change, or stylesheet insertion afterwards.

Both outputs are packaged beside each other. The site selects the `.dark.svg`
sibling using its existing theme attribute and always prints the light image.
A diagram used by documentation and a historical blog post can expose both
render targets from one authoritative `.mmd` source. Blog consumers keep their
original `mermaid_svg` target and its existing appearance.

## Tests

`test/renderer/renderer_test.go` renders fixture diagrams through the
`mermaid_svg` rule and asserts the contract of the theme: the documentation
palette and hand-drawn node class survive, the built-in lavender colors do not,
labels are measured through the pinned fonts, and the label plates keep their
border slack. It also asserts that every container title is written after the
nodes, which is the paint order `cmd/render/paint_order.mjs` produces and the
property that keeps an edge from drawing over a title. `cmd/mmdc/mmdc_test.go`
drives the interactive CLI and re-renders under a Fontconfig pointed at a
monospace-only pool; a byte-identical result shows the CLI resolved text
through the pinned fonts.

## Hermetic rendering

Mermaid measures label text through Fontconfig, so a renderer that follows the
host would produce different canvas geometry on every machine. Both entry
points instead expose only the repository's pinned fonts to Chrome:
`cmd/render/fontconfig.mjs` writes a Fontconfig file whose sole `<dir>` values
are the directory holding the pinned handwriting face and the directory holding
the pinned Liberation set, and whose aliases map the families the theme requests
onto them. The same file is applied by the build action and by an interactive
`mmdc` run, so both resolve identical text metrics.

The generated Fontconfig must use absolute paths. Chrome is a separate process
from the launcher, and a relative `<dir>` resolves against whatever working
directory it inherits rather than the action's execroot. Passing a caller's own
`-p`/`--puppeteerConfigFile` hands full control of the browser launch, and
including the font aliases in it, to that caller.

Pinned fonts fix the metrics a render measures with, but not the metrics a
consumer displays those glyphs with. A documentation page embeds an SVG through
`<img>`, and an SVG loaded that way is an isolated document: it cannot see the
page's webfonts, loads no network font, and falls back to whatever the viewer's
system provides. A host font is wider than the pinned handwriting face, so
`github.com` overflowed the label box Mermaid sized with the pinned font and
appeared clipped, even though the file was correct. `cmd/render/embed_font.mjs`
therefore inlines the pinned face as a `@font-face` data URI in every themed
document, so the glyphs travel with the SVG and the label metrics match
everywhere, including in a viewer that can load no webfont at all. The rule
adds about 50 KB to a rendered diagram; that is the cost of an image that is
correct wherever it is displayed. Its family name is read from the font's own
`name` table rather than hard-coded, so a font swap cannot silently desync it
from the theme.

The JavaScript launcher runs from Bazel's output tree. Pass absolute paths when
rendering source-tree files directly:

```sh
repo_root="$PWD"
bazel_agent bazel run //tools/mermaid:mmdc -- \
  -i "${repo_root}/path/to/diagram.mmd" \
  -o "${repo_root}/path/to/diagram.svg"
```

Declare a maintained diagram as a `mermaid_svg` action and feed the output to
its consumer, which uses declared Bazel paths and needs no absolute-path
handling. The [diagram skill](skills/mermaid-diagrams/SKILL.md) defines direct
documentation consumption and checked-in publication assets for blog posts.
The `.mmd` remains the authoritative diagram source in both cases.

The pinned upstream binary is owned by
[`third_party/com_google_chrome_headless_shell`](../../third_party/com_google_chrome_headless_shell/README.md).
