---
name: mermaid-diagrams
description: >-
  Write and render Mermaid diagrams in this repository with the checked-in
  theme, hermetic Bazel rendering, and the maintained `mermaid_svg` rule. Use
  when adding or editing a diagram source, or wiring a rendered diagram into a
  build consumer; do not use to redesign the shared theme.
---

# Write repository diagrams

Read [`tools/mermaid/README.md`](../../README.md) for the theme contract and the
measured reasons behind each setting before changing a diagram that depends on
one of them.

## Prefer the Bazel rule

Declare a maintained diagram as a `mermaid_svg` action, and feed that target
straight into whatever consumes the image:

```starlark
load("//tools/mermaid:defs.bzl", "mermaid_svg")

mermaid_svg(
    name = "rendered",
    src = "architecture.mmd",
    out = "rendered/architecture.svg",
)

docs_filegroup(
    name = "docs_assets",
    srcs = [":rendered"],
    prefix = "content/docs/<package>/assets",
)
```

For diagrams published in the main site's documentation, use
`mermaid_site_svg` from `//projects/alwaldend.com/pkg/mermaid:defs.bzl`.
It renders light/dark SVG siblings with native Mermaid layout and sans-serif
labels. Package both outputs; reference the light `.svg` in Markdown or the SVG
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

## Default to Dagre and `flowchart LR`

The shared theme selects Dagre. Write `flowchart LR` for left-to-right flow
unless the diagram or the user calls for another direction. Direction is
declared in the source, not forced by the theme. Dagre's `nodeSpacing` and
`rankSpacing` adjust spacing, but do not guarantee a maximum canvas width.

Subgraph directions are ignored when internal nodes connect outside the group;
switching between Dagre and ELK does not fix that case. Verify the actual
render before claiming that groups follow their declared direction.

## Keep theme changes in `tools/mermaid/theme.json`

The theme owns appearance, so a diagram carries structure and labels only. A
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

The height of hermeticity is the whole pipeline, not just the render. A themed
output also embeds the pinned handwriting face, because a page that shows the
SVG through `<img>` cannot see the page's webfonts and would otherwise lay the
text out in a wider system font than the renderer measured. Do not remove that
embedding or reintroduce a webfont dependency, and do not hand a viewer a
diagram whose labels only fit when a particular font happens to be installed.

Set `plain = True` on `mermaid_svg`, or `MERMAID_PLAIN=1` for `mmdc`, to render
Mermaid's own defaults for comparison. That output is neither themed nor
hermetic by design; never publish or commit it.

When a consumer has to publish an upstream-appearance comparison, use
`theme = "//tools/mermaid:upstream-theme.json"` instead: it shows Mermaid's own
appearance while keeping the render hermetic. The post at
`//projects/alwaldend.com/content/blog/diagrams-in-ac-era` is the worked
example.

For an interactive render, `mmdc` needs absolute paths because the launcher
runs from Bazel's output tree:

```sh
repo_root="$PWD"
bazel_agent bazel run //tools/mermaid:mmdc -- \
  -i "${repo_root}/path/to/diagram.mmd" \
  -o "${repo_root}/path/to/diagram.svg"
```

`mmdc` applies the same theme, pinned fonts, and paint-order fix as a build
action, so an interactive render matches the committed one. It writes only the
file named by `-o` or `--output`, so an SVG printed to stdout is not rewritten.
Always re-render a maintained diagram through its `mermaid_svg` target before
handoff, which is the only path a build action verifies.

## Verify

Render the diagram and inspect the output rather than the source. A change is
complete when the container titles are readable, no title escapes its parent
container, and no edge crosses a title plate. `//tools/mermaid/test/renderer`
holds the fixtures that pin these properties; extend them when a change
introduces a new condition.
