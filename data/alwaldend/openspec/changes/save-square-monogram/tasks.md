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

## Evidence

- SVG paths match the selected S8 study; body bottoms align at y=45.5.
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
