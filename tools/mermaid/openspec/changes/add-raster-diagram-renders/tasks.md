## 1. Render rule

- [x] 1.1 Add a `mermaid_webp` rule beside `mermaid_svg` that renders one diagram source to a `.webp` output, reusing the theme, pinned browser and fonts, and paint-order pass, and encoding through the pinned browser; verify a fixture diagram renders to a file whose bytes begin with a WebP RIFF header
- [x] 1.2 Prove one appearance across formats: render one source through both rules and assert the raster's canvas reproduces the SVG render's geometry at the rule's scale rather than introducing a second layout; verify the check fails when the raster stops following the maintained document
- [x] 1.3 Confirm the rule declares no image-conversion tool or library, so the pinned browser is the only encoder

## 2. Documentation

- [x] 2.1 Extend `tools/mermaid/README.md` with the raster rule, the accepted media types it exists for, and the reason it is a projection rather than a second appearance
- [x] 2.2 Extend the `mermaid-diagrams` skill with the raster rule, its `write_source_files` use for a checked-in asset, and the accepted media types

## 3. Validation

- [x] 3.1 Run `bazel_agent bazel test //tools/mermaid/...` and confirm the SVG contract checks and the new raster check pass
- [x] 3.2 Run `bazel_agent bazel run //tools/openspec -- validate --all --strict --no-interactive` with `OPENSPEC_PROJECT=tools/mermaid` and confirm this change validates strictly
