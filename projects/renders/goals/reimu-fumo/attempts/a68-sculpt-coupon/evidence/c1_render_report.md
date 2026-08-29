# A68 C1 native-Grab fixed-view review

## Verdict

**Reject C1 and do not build the next checkpoint on it.** The native Grab
stroke is technically real but visually ineffective. It does not reduce
HD-01 (the oversized spherical/helmet-like cap and near-vertical rear wall).
The edited vertices were dragged toward world-positive Y, which is rearward
in this file, so the edit's geometric direction is contrary to the intended
depth reduction. At visible scale the result is indistinguishable from C0;
the only above-1% pixel difference is a tiny one-sided patch in the ordinary
three-quarter view.

The next stroke should use a materially broader rear/crown footprint and move
the rear silhouette forward (world-negative Y), with side plus both
three-quarter views as hard gates. It must preserve full rear coverage.

## Review protocol

I first inspected all five C1 renders against the canonical front image, the
canonical turntable front/side/rear/three-quarter frames, and the physical
side construction photo. I did not inspect C0 or the stroke report until this
absolute image review was complete. I then inspected the matching C0 packet
and only afterward consulted the implementation diagnostics.

The five C1 images are complete, nonblank, correctly framed 640 by 640,
8-bit RGBA renders. Every image was inspected at original resolution.

## Absolute reference review (C1 only)

- Unlabeled same-subject recognition: partial. The costume reads as Reimu,
  but the head/hair construction does not read as the same physical plush
  variant without explanation.
- Intended medium: fails. The cap remains a smooth, rigid helmet mass rather
  than thin constructed hair panels conforming to a shallow stuffed head.
- Overall reference likeness: 4/10.
- Macro silhouette and proportions: 3/10.
- Constructed-plush construction: 3/10.
- Identity-defining features: 6/10.
- Contact, attachment, and occlusion: 5/10.
- Intended-medium read: 3/10.
- Presentation readability: 7/10.
- Major visible failure: yes, HD-01.
- Absolute decision: reject.

Largest discrepancies relevant to this checkpoint:

1. In side view the brown cap projects far behind the face and ends in a
   broad, almost vertical rear wall; the references show a shallower head with
   layered rear hair panels.
2. Both three-quarter views retain the oversized spherical helmet silhouette
   and excessive crown-to-rear volume.
3. The cap is a uniform smooth shell; the reference has fabric-panel flow,
   lower rear locks, and visibly softer, flatter construction.
4. The cap-to-fringe and cap-to-rear-lock transitions remain disconnected in
   construction language.
5. The bow is still much stiffer and broader than the physical reference,
   though it was outside the C1 edit scope.

Front and rear coverage remain intact: C1 creates no bald patch and does not
visibly disturb the fringe, bow, face, locks, body, cameras, or lighting.

## C1 versus frozen C0

The comparison used identical scene, cameras, lighting, render settings,
resolution, and render-spec bytes. ImageMagick RMSE values on the 0--1
normalized scale were:

| View | Normalized RMSE | Pixels differing by more than 1% | HD-01 result |
| --- | ---: | ---: | --- |
| Front | 0.0000181 | 0 | No visible change |
| Side | 0.0001235 | 0 | No visible depth reduction |
| Rear | 0.0000309 | 0 | Coverage preserved; no useful reshaping |
| Three-quarter | 0.0000991 | 257 of 409,600 (0.063%) | Minute one-sided change only |
| Mirrored three-quarter | 0.0000110 | 0 | No visible change |

This is well below a meaningful silhouette correction. The side view, which
most strongly controls HD-01, has **zero pixels differing by more than 1%**
from C0. Direct inspection likewise shows the same rear wall and spherical
depth. C1 therefore does not improve HD-01. Because its actual displacement
is rearward and unilateral, it is directionally a slight regression even
though the regression is too small to be perceptually important at 640 px.

No other material visual regression is detectable. The important failure is
that the dominant defect survives unchanged.

## Implementation diagnostics revealed after visual review

- Edited object: `A68_BackCap_Sculpt`.
- Native operator: `bpy.ops.sculpt.brush_stroke`, result `FINISHED`.
- Moved vertices: 35 of 9,449 (0.37%).
- Maximum displacement: 0.001194 m (1.194 mm), approximately 0.89% of the
  134.04 mm cap width.
- Mean displacement: 0.000371 m (0.371 mm), approximately 0.28% of cap width.
- Moved region: one positive-X, positive-Y crown/rear patch, not the broad
  rear silhouette responsible for HD-01.
- Frozen non-target objects changed/added/removed: none.

These diagnostics explain the blind visual result: the footprint is far too
small and the drag direction increases rear depth instead of reducing it.

## Frozen inputs and tool identity

- C1 blend:
  `out/reimu_fumo_attempt_068_sculpt_coupon/c1_native_grab/reimu_fumo_a68_c1.blend`
- C1 SHA-256 before and after render:
  `67c9dbf7787749038ca168215647991f6b4df422f081c35a1b410852b9931557`
- C0 comparison blend SHA-256:
  `26c8613fe3eb17a1ddfcf7c8b596ed2aa264162b86d2b1e81acf7033d1fa75ba`
- Render spec SHA-256 (identical to C0):
  `d377222d84dd64aacaf7edb071f50929ecd880f781392c0fdb4060256659d1d8`
- C1 manifest SHA-256:
  `8131f55d9155a1cf58a6f25dedcbb9c06604ebc4843b4bdebb18aa97f0446800`
- Bazel target:
  `//projects/renders/cmd/fumo_review:render_packet`
- Blender: `5.2.1 LTS`, build hash `9e2066aef7ef`.
- Render exit status: 0.

Equal before/after C1 hashes confirm that the packet render was read-only.

## Pixel artifacts

| View | Camera | Bytes | SHA-256 |
| --- | --- | ---: | --- |
| Front | `Review_front_Camera` | 522273 | `1b74b9073096a13c10c49e3a04945ad7df58d5f95505570eaeeca291bb11abb0` |
| Side | `Review_side_Camera` | 474448 | `80dc68759520c28ed6afb0d8d1a511ca6b2dfa91719e1783506131d9a42af34c` |
| Rear | `Review_rear_Camera` | 506980 | `d91e0c410fe79e8463eda7744b7408453c6d5aba37b2d081d3650642fc0bc94e` |
| Three-quarter | `Review_three_quarter_Camera` | 503577 | `45a224c84a6f543de8861353df981691863cdd5a9672acc09e41a9998ff5dc02` |
| Mirrored three-quarter | `Review_three_quarter_mirror_Camera` | 508778 | `5bb9379403d41a3f660e528a06dfe8fef1e63c7901a6d23a8a38a31c9261ae76` |

The machine-readable packet manifest is `packet/manifest.json`.
