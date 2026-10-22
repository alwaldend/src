Observed 2026-10-11 in the dedicated diagram-conversion feature worktree.

The existing renderer suite passed all three tests after adapting expectations
before implementation. It exercises exact default ratio, direct encoder parity,
native natural-canvas opt-out, explicit scales and ratio overrides, retained
content, and theme backgrounds. Artifact comparisons use the pinned browser.

An intermediate 3x-ui WebP target built successfully at 2850 × 1140 pixels, exactly
5:2. Its complete render was visually inspected: the overview is above the two
detail sections, the labels and links remain readable, and no content is cropped
or stretched. Required final check plans and delivery receipts are under ignored
out/convert-infra-diagrams; receipts own exact-candidate validation/publication.

The intermediate preview target was subsequently removed. The repository's
3x-ui documentation consumes build-time SVGs; WebPs remain user previews in
ignored task scratch. The shared raster default remains 5:2 as requested.
