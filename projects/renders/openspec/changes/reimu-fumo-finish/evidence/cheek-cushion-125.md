# Lower cheek cushion study 125

Color-assisted comparison of the canonical front and fixed assembly 122 front
finds a materially broad lower jaw: at reference rows 550/560/570/580, approximate
skin widths are 0.527/0.481/0.380/0.193 head widths, versus
0.577/0.543/0.474/0.392 in the aligned model. Uncertain skin/shadow boundaries
allow roughly 3–5 reference pixels error. Chin height is already close.
Camera elevation could narrow the apparent jaw by raising rearward corners;
these pixels establish a mismatch, not a unique anatomical reconstruction.

## Decision and bounded plan

Test a minimal edit to the existing cushion cage. Move lower front diagonals
inward and upward by 3.5 mm, and lower side corners inward and upward by
1.5 mm. Keep upper rings and the object transform unchanged. Constrain the
evaluated bottom height by a small center-bottom cage correction. Reproject
the existing diagnostic face footprints in depth onto the changed receiver,
preserving their canonical X/Z coordinates and layer offsets. This is not
final facial artwork. Reject missed rays, receiver clipping, a detached mouth,
lower-hair gaps or side/three-quarter shape regressions.

Use protected assembly 116 C, SHA-256
`012713486a2b2c9ce15da049005187b89fa178b36687dcc7179946e73ff2d829`.
Root is the sole native writer, using the same pinned background mesh route,
four threads and passing preflight. Allow one deformation and at most two
causally diagnosed receiver/contact repairs, with no parameter sweep. Save
new bytes, preserve the source, and clean-reopen for fixed front, side, rear
and three-quarter views. The frozen camera contract remains unchanged and no
stage or final acceptance follows from this local test.

## Result

Saved candidate SHA-256 is
`e4471beb3933da7656864ecb9110c6d125c8fbb97ef9ad8bae14531dfb08235a`.
Pinned Blender clean-reopened the unchanged bytes and rendered all four fixed
views. Evaluated chin height changed from 90.861566 to 90.861559 mm. Facial
footprints retained X/Z coordinates; their largest sampled face-center depth
offset error is 8.48 micrometres, below the 15 micrometre diagnostic limit.
This sampling is not a full continuous-surface proof.

Root and independent pixel reviewer `/root/garment_118_evidence` retain this
bounded cheek correction for integration: the lower contour is rounder with
no obvious side or lower-hair regression. The same approximate row-width
measurement now gives 0.547/0.500/0.427/0.340 head widths. The lowest row
remains too broad by approximately 0.147 head widths; this is not head,
camera, stage or final acceptance. No further cage tuning is authorized by
this closed local test. The source 116C is unchanged.
