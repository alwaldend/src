## Context

The user rejected both Hugo and JavaScript SVG rewriting and approved native
Mermaid layout, ordinary actor labels, native rounded nodes, and a small CSS
rule for rounded containers. Historical blog diagrams must remain unchanged.

## Goals / Non-Goals

Use Mermaid's native configuration and CSS input for documentation images;
follow explicit, saved, and system themes; retain readable fonts, mobile
layout, and light print output. Do not preserve the custom border-straddling
captions, introduce a Mermaid fork, change blog sources, or deploy the site.

## Decisions

Decision-review verdict: proceed with two native renders. The pinned Mermaid
theme engine computes colors from concrete values and has no public theme
plugin registration. Two static variants avoid both SVG rewriting and a
client-side rendering dependency, at the cost of a second generated asset.

The site owns a thin render macro, its native configuration, and a Sass palette
compiled with the existing Dart Sass toolchain. Chrome resolves the generated
CSS custom properties into Mermaid theme variables before rendering, so there
is no Sass or SVG parser and no independent copy of the site palette.

The shared Mermaid renderer retains its historical default path. Its native
path passes the font stylesheet through Mermaid CLI's CSS option and writes
the returned SVG bytes directly. It does not lift labels, rewrite colors, or
insert styles into a finished SVG. Use Dagre's native layout, normal labels, and
the rounded node syntax already present in the source diagram.

The user explicitly retained Dagre after the first native candidate used ELK.
The layout change was reverted before publication; Dagre is part of the
accepted configuration.

Markdown and the SVG shortcode share an image partial that recognizes an
existing `.dark.svg` sibling. The site's theme attribute chooses the visible
image; print always uses light. Blog pages bypass the paired-image behavior.

## Failure cases to verify

- Palette extraction loses a color or resolves the wrong mode.
- Native rendering clips labels, overlaps nested container titles, or changes
  geometry between modes; sans-serif fonts must remain embedded and measured
  against pinned files.
- Theme selection follows the OS instead of the explicit menu choice, fails
  after reload, or leaves two accessible images visible.
- The dark sibling is not packaged or linked correctly; shortcode and Markdown
  images differ; print selects dark or mobile content overflows.
- Historical SVG bytes change because the default render path was modified.

## Verification

On 2026-09-24 the browser E2E passed explicit, saved, and system selection,
standalone palettes, native label bounds, rounded shapes, mobile layout,
light print output, shortcode fallback, and unchanged blog publication files.
Light and dark page screenshots were inspected. The scoped repository suite
passed all 37 tests, including historical SVG freshness; affected-target lint
passed. After restoring Dagre, the browser E2E passed again. Final delivery
revalidates the prepared candidate and records its OID
under `out/mermaid-theme/` before updating PR #103.

The browser check produces PNGs and a JSON report in its Bazel undeclared
outputs. Task-local full-page screenshots and raw logs live under
`out/mermaid-theme/`.

## Session review

The first five native browser runs included four failures before the passing run.
`native-wrapper-collapse` was a real layout defect, fixed by sizing the image
wrapper explicitly. `lazy-image-observation` was a test timing defect, fixed
by waiting for a newly displayed lazy image to load before decoding it.
`native-shape-selection` was an assertion defect: Mermaid's empty label
rectangles are not node outlines, so corner checks now select direct shape
children. These findings are covered by the existing E2E; no new unit tests
or shared workflow changes were needed. The subsequent full test selection
took 13 seconds; cached lint completed in under a second.
