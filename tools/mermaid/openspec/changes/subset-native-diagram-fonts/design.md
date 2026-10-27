## Decisions

Use FontTools 4.67.0 without optional dependencies. Render the first native
variant with the complete pinned font, derive a conservative Unicode set from
the decoded SVG (including style/markup characters), and subset the same font.
Render the delivered variants using that subset through Mermaid's supported
CSS input. Keep the renderer's full pinned Fontconfig face for measurement.
Preserve font naming, hinting, layout closure, and deterministic timestamps.
The extra first render avoids guessing characters Mermaid generates or decodes.

Do not strip embedded fonts, rely on browser fonts, rewrite native SVG geometry,
or add a separate font registry. Existing Python dependency locking and Bazel
actions own the executable; native SVG and WebP rules pass its declared path.

## Failure cases and acceptance before implementation

Check actual native output: font remains embedded and loadable; font family,
supported rendered glyph coverage, advance widths and outlines match the pinned
original; both light and dark documents shrink substantially. Browser raster
and site checks must preserve visible labels and geometry. Subsetter failure
must fail the render without publishing an incomplete output. Historical outputs
must retain their existing behavior. Record website archive size afterward.

## Observed validation

The new native-font E2E failed with complete embedded fonts before implementation
and passed after subsetting. Both variants embed a 22,844-byte font and preserve
the tested 19 glyphs' decomposed outlines and advance widths. Native SVG, native
raster, historical renderer, and generated Python manifest checks pass.
The website's SVG total is approximately 7.16 MB after subsetting. Website
archive size also depends on navigation and minification; the site owner's
change records their separate comparison.

The website browser fixture waits for its lazy-loaded first diagram to finish
loading before decoding. Both themed diagrams, mobile layouts, and native
font outlines/metrics pass the browser and generated-font checks.
