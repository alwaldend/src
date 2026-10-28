# Verification

The before-render used gallery revision
`ac8753c65c97b92df135468d5cbcba9541e588f1`. The new browser check was written
before changing the preset. It reproduced subpixel graph strokes at 728px,
competing sequence guides, and low-contrast mindmap branches/root text.

The implemented candidate passes `:readability_test` and `:mermaid_test` in
`//projects/alwaldend.com/test/mermaid`. Their Bazel undeclared outputs retain
screenshots and JSON measurements. The former covers fourteen rendered SVGs
at 728px and 360px, with network requests blocked except embedded data. The
latter exercises the built site's theme selection, mobile and print behavior,
and historical blog diagrams.

At 728px, minimum normal path widths in the light gallery are 1.06 CSS px
(class), 1.07px (flowchart), 1.33px (ER), and 1.75px (state). Mindmap branches
use the site ink in both modes. Sequence guides are 1px in SVG coordinates,
with approximately 5.74:1 light and 5.65:1 dark canvas contrast, below messages.

The semantic fixture verifies solid/dashed distinction, arrow direction,
resolving markers, stronger thick edges, and invisible edges. It caught that
Mermaid emits both normal and thick classes on thick paths. The production
selectors explicitly exclude thick and invisible classes; thick paths retain
their native 3.5px width. Mode comparisons verify rendered bounds, labels,
fonts, dashes, and markers. Classic shapes can have randomized collinear Bézier
control points, so path-string equality would incorrectly reject equivalent
rendered geometry.

The measurements do not guarantee readable labels on every phone or arbitrarily
large graph. Mobile screenshots are retained; dense diagrams still need native
size viewing. Appearance alternatives and their limitations are recorded in
[appearance-options.md](appearance-options.md). Final candidate gates and
publication state belong to delivery receipts, not this document.

## Session ergonomics

- `mermaid-semantic-class-overlap` — fixture-tested: selecting the normal
  class alone also changes thick edges. Keep the semantic fixture and explicit
  exclusions; source class names alone are insufficient evidence.
- `mermaid-classic-path-serialization` — fixture-tested: six baseline mode
  comparisons falsely failed on randomized collinear control points despite
  identical bounds. Compare rendered geometry and semantic marker structure.
- `mermaid-action-path-resolution` — live: one scratch measurement launch
  wrongly treated a generated browser path like an external repository path.
  Resolve `external/` through the output base and generated paths through the
  execroot. No shared runtime or dependency changes were needed.
