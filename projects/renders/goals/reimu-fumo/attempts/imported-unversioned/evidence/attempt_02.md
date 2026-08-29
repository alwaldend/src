# Reimu Fumo attempt 2

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 2 — frozen reference calibration

**Candidate:** `LANDMARKS.md`, SHA-256
`c36ea42d27f72f616c123d3c21e12a5a5372c893e7700a2fc38d8c84c727ee66`,
reference-contract stage, review packet `reference-metrics-v1`.

**Failure targeted:** Earlier variants guessed dimensions, optimized the front
view first, treated oblique images as orthographic, and averaged visibly
different plush poses and specimens.

**Hypothesis:** Assigning each source a controlling role, freezing normalized
landmarks and camera alignment before geometry, and recording uncertainty and
variant conflicts will prevent the next model from drifting toward convenient
but false proportions.

**Plan:** Extract every GIF frame; select exact controlling frames; annotate
the stills and rotation views; measure the outer head, face, eyes, fringe, bow,
garment, sleeves, skirt, feet, depth, and contact; normalize to front head
width `Wh`; freeze source authority, exclusions, tolerances, and cameras; and
review conflicts instead of averaging them.

**Work performed:** Measured both 474-pixel front stills, the oblique physical
side, all four coalesced rotation frames, and all 91 sofa frames. The resulting
contract records 39 controlled targets, exact GIF frames, source uncertainty,
camera coordinates, alignment rules, and variant handling.

**Evidence:**

- Local measurement dossier SHA-256:
  `f6d303c0c6c11befc0da286931e852e3c5c1906f326014ddf0eba955d222f0a1`.
- Turn contact-sheet SHA-256:
  `736884a1866aafe890f6c4bf748b1d1a3c2bef8cba25ef47e3bbcf91ba718500`.
- Sofa contact-sheet SHA-256:
  `e46c85843e7cd2314f66dcdcd374413233806f068d021834d06032b8ca106119`.
- Physical-front annotation SHA-256:
  `d17068ee444187e046bf623c4d1438360f041f06e02c05ce0ddfc79a405a45ce`.
- Clean-front annotation SHA-256:
  `621874d8c63e5bf254af6c5c4c5af44f3fa326a4e6a6b1a6b93bcdcb1b9be08e`.
- Physical-side annotation SHA-256:
  `3e44761c12b31f639a2dfe3ba1665de2c16ee0dda6da19c5a84eca15b038b25d`.
- Default bow-span conflict was measured as `1.37 Wh` in the rotation GIF,
  `1.69 Wh` in the physical front, and `1.94 Wh` in the clean front. The
  physical-front state now controls the base; other states are rig targets.

**Criterion results:** The evidence contract for criterion 2 passes as a
prerequisite, but criterion 2 remains unverified for model geometry. Reference
source coverage is now sufficient for a first macro-form candidate. Visual,
construction, reuse, animation, and integrity criteria remain failed or
unverified.

**Decision:** Accept and freeze the calibration. Start one clean macro-form
candidate without retaining baseline render geometry.

### Progress and approach audit after attempt 2

- **Improved:** source authority, measurement repeatability, side/rear control,
  uncertainty disclosure, and protection against post-hoc landmark changes.
- **Unchanged:** the actual model, which remains an absolute 3/10 baseline.
- **Absolute result:** no visual improvement is claimed from measurement work.
- **Approach evidence:** the depth target `0.74 Wh`, flat central face, large
  head/body overlap, compact garment, and hem-hidden foot pods contradict the
  egg head, tall torso, cone skirt, and tangent feet of the old representation.
- **Repeated-defect diagnosis:** all major form defects have survived at least
  two primitive-derived attempts, so the goal's strategy-reset rule forbids
  another parameterized refinement of that geometry.
- **Highest-leverage problem:** produce a black-silhouette candidate whose head
  cushion, compact body, bow pockets, sleeve envelopes, skirt, and foot pods
  pass front, side, rear, and three-quarter ratios before identity detail.
- **Next approach:** create low-resolution sewn-pattern cages in a clean local
  Blender file, sculpt broad stuffing only, and reject the candidate before
  detail if any helmet, cone, tube, wing, human, or disconnected read remains.
