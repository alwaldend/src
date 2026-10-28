# Monospace verification

The new equal-glyph-advance browser assertion was added before switching the
font. It rejected the previous proportional labels. The final renders load
the embedded Liberation Mono face and pass that assertion for all six gallery
types and the edge-semantics fixture in both modes at 728px and 360px.

At the unchanged 16px base size, the class canvas grew to 1462px and the
flowchart to 1331px. Their previous 1.75px paths became 0.87px and 0.96px at
article width. Normal graph paths now use 2.25px, restoring approximately
1.12px and 1.23px respectively. Thick paths remain 3.5px; invisible and dashed
semantics remain distinct. Mindmap branches, sequence guides, message weights,
and the active Sass palette retain their prior settings.

The nested site fixture measures 1714 × 739px with 44px nesting inset and 9px
title clearance. Its old 1700px width guard was tied to proportional labels;
the new 1800px bound retains a compact canvas relative to the former
2120 × 790px layout. Existing containment, title-clearance, caption-fit,
mobile, print, and historical-blog checks remain in place and pass.

Passing checks before delivery preparation:

- `//projects/alwaldend.com/test/mermaid:readability_test`
- `//projects/alwaldend.com/test/mermaid:mermaid_test`
- `//tools/mermaid/test/renderer:native_font_test`
- `//tools/mermaid/test/renderer:native_raster_test`
- `//tools/mermaid/test/renderer:aspect_ratio_test`

Their Bazel undeclared outputs contain screenshots, rendered references, and
measurements. All twelve gallery SVGs were visually inspected. The four
catalog entries preserve all 48 colors from commit
`fd25e40c4bad58c65511cf4ebaa60c651669d823`; the active palette is byte-identical
to that commit. Palette values now have one maintained owner in the reusable
package, and proposal metadata references its IDs.

Final candidate gates and publication state belong to delivery receipts under
`out/mermaid-monospace/`. Dense diagrams still need native-size viewing on
small screens.

The main integration risks were measured rather than inferred: wider glyphs
weakened scaled paths, and a Sans-specific canvas bound rejected otherwise
intact geometry. Existing browser artifacts exposed both. Preview evidence
uses workspace copies of public SVG outputs because the inline preview cannot
resolve the external Bazel-output symlinks directly.
