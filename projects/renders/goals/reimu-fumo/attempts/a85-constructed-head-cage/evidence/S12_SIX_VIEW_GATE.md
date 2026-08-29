# A85 S12 six-view regression gate

## Scope and binding

This is a pixel-only regression review of the pinned S12 six-view packet. I
did not inspect Blender geometry, cameras, or authoring intent. Silhouettes
were isolated from the beauty renders at the established 55% luminance
threshold, with a 32--80% sweep to separate projection boundaries from
lighting and anti-aliasing.

- Blend SHA-256:
  `982da6404ea6edcbb4432903e67dad4ee5c130a203a5a5727a374b773fc9ad8a`
- Packet spec SHA-256:
  `404129616f9f9c90d64d4fab52f96a3aa4057a232cd0550e62ff74fdf9694a85`
- Front:
  `2d4a9f4c81780a71a4575a857ae144b250abb0e5043697f07c37af36ecf30d93`
- Rear:
  `d240a82a3f66e46bcccdeecf8fee1b19ca9408c61657b95a916930aa4301d1ad`
- Profile left/right:
  `84943201026e9ce9e14752a06d30d9140a490104fca6f662b7c0190b6f8d5d7a`,
  `a80a0b2a14482321c9e4039d5815dc5064d487c9c9c9361e042ce756a09d4ee9`
- Three-quarter pair:
  `a1213326eacc79719d6976c737444a4be70ea41f12b903d7af7af59b3ff501e3`,
  `758e19d19d671274e0ae62eaa1065aa98c2b7afed3250c310525336284df9a53`

## Fixed-threshold bounds

At 55%, the silhouette bounds are:

| View | Bound | Horizontal center |
| --- | --- | ---: |
| Front | `382 x 373 +65 +61` | `255.5 px` |
| Rear | `382 x 371 +65 +61` | `255.5 px` |
| Profile left | `240 x 374 +136 +61` | `255.5 px` |
| Profile right | `240 x 372 +136 +61` | `255.5 px` |
| Three-quarter | `410 x 373 +51 +61` | `255.5 px` |
| Three-quarter mirror | `410 x 374 +51 +61` | `255.5 px` |

Every paired view has exactly the same width, horizontal offset, and center.
Across the 32--80% sweep, these widths and offsets remain invariant. Pairwise
height differences vary by one to three pixels as the threshold moves,
indicating that the differences belong to the anti-aliased, shadowed bottom
edge rather than a stable projected side boundary.

## Mirrored silhouette comparison

Each left view was flipped horizontally before comparing its binary mask with
the corresponding right view. The lower contact/self-shadow begins near
`y=400`; therefore both full-frame and upper-silhouette mismatch counts are
reported.

| Pair | Full 512 x 512 mismatch | Mismatch above `y=400` | Upper-crop rate |
| --- | ---: | ---: | ---: |
| Profile left vs flipped profile right | `2800 px` | `3 px` | `.0017%` |
| Three-quarter vs flipped three-quarter mirror | `2234 px` | `13 px` | `.0075%` |
| Front vs flipped rear | `423 px` | `7 px` | `.0040%` |

The profile pair is therefore equal to within three thresholded pixels over
the unshadowed silhouette. The mirrored three-quarter pair is equal to within
13 pixels, also negligible at 512 x 512. The full-frame mismatches are
localized to the bottom 34 pixels: the light and self-shadow reverse between
the paired cameras, causing the luminance threshold to cut into opposite sides
of the underside. Visual inspection shows the same reversal, rather than a
stable outline displacement.

## Rear versus front

Front and rear share the exact `382 px` width, `x=65` offset, `255.5 px`
center, and `y=61` top at 55%. Rear height is two pixels shorter at that
threshold. The sweep changes this difference from one to three pixels while
leaving the width and center exact; above the bottom-shadow band the mirrored
masks differ by only seven pixels. There is no pixel evidence of a rear bulge,
horizontal shift, or different crown/side projection. The bottom two-pixel
height difference is lighting-confounded and is not a justified rear-form
claim.

## Verdict

**Pass the S12 six-view projection-regression gate.** No left/right profile
asymmetry, mirrored three-quarter inconsistency, front/rear bound regression,
or projection asymmetry beyond lighting is detected. The unshadowed masks are
effectively mirrored, and every pair preserves exact horizontal registration.
The only unresolved area is the lowest shadowed/contact band, where beauty
pixels cannot distinguish surface projection from illumination without side
component-ID renders; it contains no stable cross-threshold evidence of a
geometry mismatch.
