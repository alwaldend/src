---
name: mermaid-diagrams
description: >-
  Write and render Mermaid diagrams in this repository with the checked-in
  theme, hermetic Bazel rendering, and the maintained `mermaid_svg` and
  `mermaid_webp` rules. Use when adding or editing a diagram source, or wiring
  a rendered diagram into a build consumer; do not use to redesign the shared
  theme.
---

# Write repository diagrams

## Preserve the requested scope

For an application overview, start with actors, application processes, storage
and external services. Add internal packages, job registries, generators or other
implementation detail only when requested. Preserve supplied node names, grouping,
arrow directions and labels. Use a flowchart for structural relationships; a state
diagram describes states and transitions. Keep planned relationships identifiable
in the accompanying prose rather than claiming that drawing them implements them.

## Build the owner's normal target

1. Inspect the owning `.mmd` and `BUILD.bazel`. Reuse the existing render target;
   if the requested format or color variant is missing, add a small normal target
   beside it. Use this workflow for local previews too. Honor a restriction on
   publication separately from permission to edit/build local files.
2. Match the owner's existing macro and theme. "Regular repository styling"
   means that preset, not Mermaid's built-in `default`/`dark` or the upstream
   comparison theme. For native site diagrams, pair `mermaid_site_svg` with
   `mermaid_site_webp`; historical consumers use `mermaid_svg`/`mermaid_webp`.
3. Run `bazel build //<owner>:<render_target>` from the workspace root. The rule
   supplies the pinned browser, fonts, palette, font subsetter and runfiles. Do
   not reconstruct private renderer commands from `aquery`, cached executables,
   discovered font paths or manually assembled launchers for ordinary diagram
   work. Diagnose a failed target through its declared inputs and build log.
4. Inspect the actual built asset at the intended display size. Show that image
   through the available display mechanism and provide the file in the requested
   format. An HTML preview or screenshot does not replace the requested WebP.

For example, add a dark raster beside an existing native site SVG, sharing its
authoritative source:

```starlark
load("//projects/alwaldend.com/pkg/mermaid:defs.bzl", "mermaid_site_webp")

mermaid_site_webp(
    name = "architecture_webp",
    src = "docs/architecture.mmd",
    out = "rendered/architecture.dark.webp",
    aspect_ratio = [],
    color_scheme = "dark",
    scale = 2,
)
```

Build `//<owner>:architecture_webp`; its file is under
`bazel-bin/<owner>/rendered/architecture.dark.webp`. Keep disposable copies under
ignored task scratch when a stable attachment path is needed. Simple diagrams
should be readable directly; do not add navigation or zoom controls to compensate
for unnecessary detail. See the relevant sections of
[`tools/mermaid/README.md`](../../README.md) for renderer and theme mechanics.

## Connect the render to its consumer

For diagrams published in the main site's documentation, use
`mermaid_site_svg` from `//projects/alwaldend.com/pkg/mermaid:defs.bzl`.
It renders light/dark SVG siblings with native Mermaid layout and the site's
embedded typography. Package both outputs; reference the light `.svg` in Markdown or the SVG
shortcode, and the site selects its `.dark.svg` sibling for dark mode.
It uses no SVG post-processing or custom actor decoration.
Keep historical blog assets on their existing `mermaid_svg` render targets;
when one source serves both, declare a separate site render from that source.

For README documentation, consume the rendered target directly instead of
committing an SVG or adding a `write_source_files` copy. Point the README image
at the path `docs_filegroup` publishes and let the build produce the file.

Blog posts keep their images beside `index.md` as checked-in publication
assets. Use `write_source_files` to copy hermetic `mermaid_svg` outputs into the
post package, publish them through its `docs_filegroup`, and run the generated
freshness tests. The `.mmd` stays authoritative; never hand-edit the SVG. When
another package owns the diagram, copy its render target instead of creating
a second diagram source. The blog and that package's documentation may publish
the same render under their respective URLs.

Other consumers that cannot consume a Bazel target may also use a checked-in
projection through `write_source_files`; document the reason in the package.

## Render a raster for a service upload

A consumer that uploads the image to a service may need a raster, because some
upload endpoints accept only a fixed set of image media types that excludes
`image/svg+xml`. The X Articles media upload endpoints, for example, accept a
fixed raster set that `projects/x_article_uploader/README.md` owns and enumerates;
`mermaid_webp` produces `image/webp`, which is in that set. Declare a
`mermaid_webp` action beside the SVG render:

```starlark
load("//tools/mermaid:defs.bzl", "mermaid_webp")

mermaid_webp(
    name = "rendered_webp",
    src = "architecture.mmd",
    out = "rendered/architecture.webp",
)
```

By default the rule encodes the historical shared theme with its pinned browser,
fonts, and paint-order pass. For a site-themed blog or social image, use
`mermaid_site_webp` from `//projects/alwaldend.com/pkg/mermaid:defs.bzl`; it selects
native rendering with the site's palette and embedded typography. Select
`color_scheme = "dark"` for the site's dark appearance; `"light"` is the default.
The raster has that fixed appearance in page embeds and social previews. Other
native consumers set `native = True` with the SVG target's `palette`, `theme`,
and `label_font`. Native WebP follows the selected SVG without post-processing;
both modes use the pinned browser encoder with no image-conversion toolchain.
The raster uses the diagram's natural size and the rule's `scale` (default `1`).
Set `scale = 2` when the consumer needs higher pixel density.
WebP canvases default to 5:2 proportions with centered, theme-colored padding.
Set `aspect_ratio = []` for natural proportions or supply another positive integer
width/height pair. The CLI equivalent for natural output is `--aspect-ratio natural`.
There is no `plain` raster: use the SVG rule's `plain` opt-out for comparison.

For a consumer that keeps the image as a checked-in asset, copy the rule output
into the package with `write_source_files` instead of committing a hand-made
raster, so the asset can be regenerated and compared against a fresh render.
The `.mmd` source stays authoritative. A blog post that is intended for
syndication references the raster render beside `index.md` rather than the SVG.

## Default to Dagre and `flowchart LR`

The shared theme selects Dagre. Write `flowchart LR` for left-to-right flow
unless the diagram or the user calls for another direction. Direction is
declared in the source, not forced by the theme. Dagre's `nodeSpacing` and
`rankSpacing` adjust spacing, but do not guarantee a maximum canvas width.

Subgraph directions are ignored when internal nodes connect outside the group;
switching between Dagre and ELK does not fix that case. Verify the actual
render before claiming that groups follow their declared direction.

## Keep appearance in the owning theme

`tools/mermaid/theme.json` owns the historical shared appearance; native site
macros select their bundled site theme. Keep that existing selection when changing
diagram content or choosing a dark variant. A diagram carries structure and labels only. A
diagram that needs a different look should not carry a local `themeCSS` block:
change the theme when the look should apply everywhere.

A diagram may still override a specific key with a Mermaid directive, because
Mermaid merges a directive over the file config:

```text
%%{init: {"look": "classic"}}%%
flowchart LR
    a[Square] --> b(Rounded)
```

Shapes stay semantic. `()` and the other inherently rounded shapes render
rounded, `[]` renders as a sketched rectangle, and nothing forces every node
into one outline. Pick the shape for what the node means rather than for the
look.

Keep nested container titles short and retain the shared title styling. Dagre
uses a fixed inset between nested containers, so the theme lifts title plates
to avoid overlap. Verify title and container containment from the rendered SVG
rather than assuming that a successful render has readable geometry.

## Render hermetically

The theme is applied by default, both entry points pin Chrome, and both expose
only the repository's pinned fonts through Fontconfig. Never add a CDN script,
a `<script type="module">` import, a downloaded font, or a network fetch to a
diagram or theme: rendering must resolve everything from the Bazel graph.

Hermeticity includes the delivered font. Historical output embeds its pinned
handwriting face; native site output embeds the font selected by its bundled
theme. An SVG shown through `<img>` cannot see the page's webfonts, so retain
that embedding rather than relying on an installed font or a network font fetch.

Set `plain = True` on `mermaid_svg`, or `MERMAID_PLAIN=1` for `mmdc`, to render
Mermaid's own defaults for comparison. That output is neither themed nor
hermetic by design; never publish or commit it.

When a consumer has to publish an upstream-appearance comparison, use
`theme = "//tools/mermaid:upstream-theme.json"` instead: it shows Mermaid's own
appearance while keeping the render hermetic. The post at
`//projects/alwaldend.com/content/blog/diagrams-in-ac-era` is the worked
example.

The public `//tools/mermaid:mmdc` target remains available for an explicitly
requested CLI render or focused CLI troubleshooting. It requires absolute input
and output paths because its launcher runs from Bazel's output tree; the README
contains the invocation. It applies the historical shared theme, pinned fonts and
paint-order fix. Its SVG output does not substitute for building a native site
variant or the maintained target for the requested format before handoff.

## Verify

Render the diagram and inspect the output rather than the source. A change is
complete when the container titles are readable, no title escapes its parent
container, and no edge crosses a title plate. `//tools/mermaid/test/renderer`
holds the fixtures that pin these properties; extend them when a change
introduces a new condition.
