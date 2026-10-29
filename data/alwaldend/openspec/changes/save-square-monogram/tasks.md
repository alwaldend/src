## 1. Selected artwork

- [x] 1.1 Save S8 as a canonical SVG and compare its path coordinates with
      the approved study.
- [x] 1.2 Render and inspect a 512 by 512 dark lossless WebP; verify its
      dimensions, palette, and full underline.
- [x] 1.3 Document and expose the assets with repository-internal Bazel
      visibility; verify declared files and visibility against data policy.
- [x] 1.4 Center the complete mark and verify equal opposing margins in
      the SVG geometry and rendered WebP bounds.
- [x] 1.5 Update the documented recipe to render at 448 by 448 and pad to
      512 by 512; execute it and compare with the delivered padded WebP.

## 2. Small-size refinement

- [x] 2.1 Compare restrained refinements with S8 at 16, 24, and 32 pixels
      in both palettes; preserve its large-size character and centered bounds.
- [x] 2.2 Apply the selected geometry to the canonical SVG, update its
      documentation, and regenerate the ignored padded dark WebP.
- [x] 2.3 Verify the generated source archive and site favicon references
      match the refined canonical SVG.

## Evidence

- The initial SVG matched the selected S8 study. The small-size refinement
  retains its topology and external bounds; body bottoms align at y=45.5.
- A 2.5-unit horizontal translation centers the complete mark within the
  SVG, leaving equal 8.5-unit side margins and 8-unit vertical margins.
- The documented padded preview has foreground bounds of 330 by 336 at
  offset 91/88 in a 512-square image, leaving equal opposing margins.
  Executing the README command reproduced the delivered padded WebP
  byte-for-byte, as verified with `cmp`.
- The preview was visually inspected, including the full underline.
- ImageMagick identifies a 512 by 512 WebP with background RGB 33/33/33
  and ink RGB 237/237/237. Its `VP8L` chunk confirms lossless encoding.
- The WebP remains in ignored output and is excluded from Git. The source
  filegroup uses the existing data-tree visibility pattern.
- Publication gates and candidate identity are recorded by repo_delivery
  under ignored `out/s8-favicon/`.
- Browser comparisons covered the original, a thicker-center candidate,
  and the selected wider-channel candidate at 16, 24, and 32 pixels in
  both palettes. The selected version preserves the full underline and
  upright, increases the minimum stroke from 2 to 3 units, and opens the
  outer-to-middle channels to 4 units. Its minimum clearance is 3.5 units
  above the terminal. The outer vertical strokes narrow from 6 to 5 units
  to accommodate the clearer winding within the same bounds.
- At 16 pixels, the terminal grows from 0.5 to 0.75 pixels and the main
  channels reach 1 pixel. Antialiasing remains visible at standard density;
  the refinement improves separation without promising pixel-perfect edges.
- The refined padded WebP was regenerated from the documented command and
  inspected. The site source archive contains one canonical favicon; the
  homepage, blog, docs, projects, project detail, and 404 pages each reference
  its existing content-fingerprinted resource with byte-identical SVG content.
