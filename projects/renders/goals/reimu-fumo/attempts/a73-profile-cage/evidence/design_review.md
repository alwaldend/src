# A73 receiver profile-loft design review

## Verdict

Use one **asymmetric profile loft with independent front, rear, and width
curves**, pole-free quad caps, and an explicitly sampled front-to-rear section.
This is a bounded reset from A72's separable superellipsoid: no global vertical
power, no global section exponent, no single top or bottom pole, and no
front/rear symmetry.

This remains a construction hypothesis. The receiver is largely hidden by the
face and hair in every source, so the references do not reveal its exact seam
or naked outline. They do bound the visible front plane, total depth, crown,
rear-leaf reserve, and lower-rear turn-in. A73 must be judged as a viable rest
surface, not presented as source-proven manufacture.

The ring-loft family still has a material risk of becoming a generic cushion.
The following contract makes the difference from A72 testable before any hair
is added. If the first packet still reads as an egg, mattress, or rounded box,
close the profile-loft family instead of tuning it.

## Evidence map

- Canonical 25 cm front: controls `Wh`, receiver scale, broad frontal mass,
  crown/underside asymmetry, and the shallow visible face plane.
- Canonical slow-180 frames 03, 07, and 29: control front and both
  front-three-quarter transitions.
- Canonical frames 10/11 and 25/26: control the two profiles, high rear volume,
  face-plane shallowness, and the depth reserved for free rear hair leaves.
- Canonical frame 18 and frames 14/15/22: control rear and rear-three-quarter
  taper. Confidence is lower because the receiver is fully occluded.
- Physical front and clean front: support a slightly flattened stuffed face,
  not a sphere, skull, or raised mask.
- Physical side: supports a shallow receiving cushion followed by separately
  draping hair; it vetoes filling the posterior hair silhouette with the
  receiver.
- Older turn and sofa frames: qualitative cross-checks for soft panel contact,
  edge roll, and non-rigid transitions. They do not override the canonical
  pair dimensionally.
- A72 five-view packet: negative evidence. It shows the exact rounded-mattress
  front, symmetric egg profile, low-centered rear fullness, and lower pole that
  this contract must remove.
- The historically accepted A25 support-cushion V3 is not a geometry parent.
  Under the current absolute gate its front/rear are rounded rectangles and
  its side is a symmetric capsule, repeating the same family as A72.

## Coordinate and scale contract

- `Wh = 132 mm` is the evaluated maximum receiver width.
- `H = 0.880 Wh = 116.16 mm` is the evaluated receiver height, within the
  frozen `0.84--0.90 Wh` support band.
- Let `t = 0` at the evaluated underside and `t = 1` at the evaluated crown.
- `W(t)` is half-width; `F(t)` is positive reach from the datum toward the
  face; `R(t)` is positive reach from the datum toward the rear.
- Geometry coordinates are `x = +/- W`, `y_front = -F`, `y_rear = +R`.
- The evaluated maximum depth target is `0.800 Wh = 105.6 mm`, inside the
  frozen `0.77--0.85 Wh` band. The remaining roughly `0.35 Wh` of the layered
  side envelope belongs to independent rear hair, not this receiver.

The terminal profile rings sit at `t=.015` and `.985`; shallow all-quad cap
patches extend their centers to the actual `t=0` and `1` extrema.

## Supported source-pixel witnesses and uncertainty

These are the pixel measurements already frozen in the reference dossier and
Attempt-53 calibration. They constrain the surrounding visible head/hair
envelope; none exposes the naked receiver, so none is converted into a precise
receiver ring.

- Canonical 1000 px front: outer head/hair `Wh = 368 +/- 4 px`, center
  `x = 485 px`, crown top `y = 231 px`. The central fringe tip is recorded as
  `(.588, .677) Wh` from the hair-left/crown datums. These control final front
  registration, not hidden depth.
- The turntable uses one unrecentered, unscaled `498 x 498 px` canvas. Its
  common projected-axis estimate is `x = 247.5 px`, with `8 px` standard
  deviation and a `36.5 px` support range; it is only a coarse initialization.
  Nominal frame yaw uncertainty is `+/-12 degrees`.
- Turn frame 005, near front: visible head/hair conservative bbox
  `[100,64]--[346,349]` (`247 x 286 px`); permissive bbox
  `[94,63]--[351,350]` (`258 x 288 px`).
- Turn frame 009, three-quarter: conservative bbox
  `[142,63]--[409,323]` (`268 x 261 px`); permissive bbox
  `[136,62]--[415,329]` (`280 x 268 px`).
- Turn frame 012, nominal side: conservative bbox
  `[169,69]--[442,337]` (`274 x 269 px`); permissive bbox
  `[164,63]--[448,344]` (`285 x 282 px`).
- Turn frame 019, nominal rear: conservative bbox
  `[145,66]--[386,349]` (`242 x 284 px`); permissive bbox
  `[119,64]--[399,350]` (`281 x 287 px`).
- The calibration records `1--3 px` dark-hair/red-cloth segmentation
  uncertainty. Bow-dominated crown/root pixels are explicitly invalid for
  base head/hair inference, and the visible mask contains outer hair as well
  as any hidden support.

Therefore these bboxes can veto a receiver that protrudes through the outer
envelope, but they cannot justify a naked-receiver silhouette overlay or the
exact per-height values below. The dimensional profile remains a reversible
P0 hypothesis inside the frozen `0.75--0.85 Wh` depth and `0.84--0.90 Wh`
height bands.

## Explicit height profiles

These values are the initial P0 cage contract. Interpolate each column with a
shape-preserving monotone cubic or connect these exact rings and let one level
of Catmull-Clark provide continuity. Do not replace the table with fitted
powers or a superellipse.

| `t` | `z` mm | `W/Wh` | `F/Wh` | `R/Wh` | full depth mm |
| ---: | ---: | ---: | ---: | ---: | ---: |
| .015 | -56.3 | .180 | .120 | .140 | 34.3 |
| .045 | -52.9 | .280 | .200 | .230 | 56.8 |
| .100 | -46.5 | .380 | .280 | .300 | 76.6 |
| .180 | -37.2 | .450 | .335 | .350 | 90.4 |
| .280 | -25.6 | .485 | .350 | .380 | 96.4 |
| .420 |  -9.3 | .500 | .355 | .410 | 101.0 |
| .560 |   7.0 | .500 | .355 | .435 | 104.3 |
| .690 |  22.1 | .490 | .355 | .445 | 105.6 |
| .780 |  32.5 | .465 | .335 | .425 | 100.3 |
| .860 |  41.8 | .420 | .300 | .370 | 88.4 |
| .920 |  48.8 | .340 | .235 | .285 | 68.6 |
| .965 |  54.0 | .230 | .150 | .175 | 42.9 |
| .985 |  56.3 | .140 | .085 | .095 | 23.8 |

Consequences are deliberate:

- The front-center reach varies by only `.005 Wh` from `t=.28` through `.69`,
  so the face side remains shallow and cushion-like rather than egg-shaped.
- Rear reach peaks at `t=.69`, not at mid-height or low on the form.
- The lower rear at `t=.28` is `.065 Wh` forward of that maximum, meeting the
  requested `.05--.08 Wh` nape undercut before the stronger underside turn-in.
- Maximum width occurs around `t=.42--.56`; crown and underside have different
  taper rates, preventing vertical mirror symmetry.

## Per-ring front-to-rear section

A72's section exponent of `4` created squared shoulders. A73 should sample the
positive-X half-section from these explicit anchors, then mirror it across X:

| station | `x/W` | normalized Y |
| --- | ---: | ---: |
| front center | .00 | `-1.000 F` |
| front inner | .38 | `-0.995 F` |
| front outer | .68 | `-0.900 F` |
| front shoulder | .88 | `-0.620 F` |
| temple maximum | 1.00 | `+0.060 R` |
| rear shoulder | .90 | `+0.580 R` |
| rear outer | .68 | `+0.840 R` |
| rear inner | .36 | `+0.970 R` |
| rear center | .00 | `+1.000 R` |

Use a C1 cubic through these anchors, with horizontal tangents at front and
rear center, and sample 32 perimeter vertices. Do not use a shared signed-power
formula. The broad front span is intentional; the temple maximum sits slightly
behind the datum so the three-quarter highlight rolls around instead of
forming a squared front corner.

Near the two caps, blend this mid-section toward an ellipse so the broad face
does not continue as a hard shoulder into the crown or underside. One explicit
weight is sufficient:

```text
q(t) = smoothstep(.10, .28, t)
       * (1 - smoothstep(.78, .94, t))
section(t) = q * mid_section + (1 - q) * ellipse_section
```

This blend changes only the azimuthal section. It must not alter the independent
height profiles above.

## Pole-free topology and caps

- Keep 32 corresponding vertices around each of the 13 profile rings and join
  adjacent rings with quads.
- Close each 32-vertex terminal ring with a directly generated `9 x 9` quad
  patch. A `9 x 9` grid has exactly 32 perimeter edges, so its boundary can be
  stitched one-to-one without a fan, triangle pole, or high-valence center.
- Map the cap grid to the terminal rounded section and give its interior a
  broad dome of only `.010--.013 Wh` beyond the boundary plane. The cap center
  should have ordinary valence four; there must be no axis vertex shared by
  all perimeter edges.
- Generate the cap grid deterministically rather than relying on an ambiguous
  edit-mode fill. A square-to-disc map or a Coons patch is acceptable as long
  as the 32 boundary coordinates are byte-identical to the terminal ring.
- Expected control topology is about 514 vertices and 512 all-quad faces:
  416 ring vertices, 98 new cap-interior vertices, 384 ring-strip quads, and
  128 cap quads.
- Preserve longitudinal loops at the front center, front shoulder, temple,
  rear shoulder, and rear center as named vertex groups. They are later
  attachment coordinates, not visible seams.

Use smooth shading and at most one viewport/two render subdivision levels.
Because subdivision shrinks extrema, evaluate the final cage and apply only
global X/Y/Z precompensation needed to recover `132 x 105.6 x 116.16 mm`.
Do not tune individual profile rows after looking at one favored view.

## Frozen P0 landmark gates

### Numeric

- Evaluated width: `1.000 +/- .010 Wh`.
- Evaluated height: `.880 +/- .020 Wh`.
- Evaluated maximum depth: `.800 +/- .015 Wh`, and always within
  `.77--.85 Wh`.
- Rear maximum occurs at `t=.64--.73`.
- `R(.69) - R(.28) = .05--.08 Wh`; target `.065 Wh`.
- Front-center reach over `t=.28--.69` varies by at most `.010 Wh`.
- Rear reach is strictly increasing from `t=.18` to the high maximum and
  strictly decreasing from the maximum to `t=.985`.
- No top or bottom axis pole; cap-center valence is four and cap patches have
  no triangles or n-gons.
- No evaluated self-intersections, degenerate faces, or non-manifold edges.

### Fixed-view pixels

- Front: broad stuffed support, but no long rounded-rectangle wall, flat
  mattress shoulder, pointed underside, or vertical mirror symmetry.
- Side: the front centerline is shallow through the face zone; the decisive
  fullness is above center and rearward; the lower rear visibly turns forward
  before the underside. A symmetric egg is an immediate reset.
- Rear: broadest supported mass is upper-mid, with a visible narrowing toward
  the lower rear. A low-centered bulge or vertical rear wall is a reset.
- Both three-quarters: the face-to-temple highlight bends continuously; there
  is no superellipse corner, box shoulder, slab, or shield read.
- Top and bottom: closures are broad domed patches with no dimple, nipple,
  star, fan, or subdivision pinch.

The isolated receiver must score at least `6/10` for macro-silhouette viability
in an implementation-blind review. This bounded gate does not claim hair,
identity, materials, or final plush approval.

## Implementation order and early exits

1. Plot the front and side polylines from the table before Blender. Reject the
   build script if either becomes vertically symmetric or moves the rear
   maximum away from `t=.69`.
2. Build only the receiver in a clean candidate copy. Keep the canonical front,
   all controlling turn frames, physical front, and physical side in the
   review packet; do not judge against a single image.
3. Render front, side, rear, both three-quarters, top, and bottom under the A72
   fixed neutral setup. The profile and both three-quarter views are controlling;
   the front alone cannot pass the candidate.
4. Run pixel review before exposing topology or the table to the reviewer.
5. If the form is directionally viable but misses one frozen landmark, allow
   one bounded row correction and rerun every view. If it still reads as an
   egg/mattress/rounded box, close the ring/profile-loft representation.
6. Only after P0 passes may another attempt derive a thin hair skin and free
   rear leaves from surface-relative anchors. Never thicken or texture this
   hidden receiver to disguise a failed silhouette.

## Main regression risk

The broad front plane plus central width plateau can still become a mattress
if an implementation smooths away the height-dependent rear shift or uses a
single exponent for all cross-sections. The non-negotiable evidence that A73
is a real reset is the combination of a high rear maximum, `.065 Wh` lower
undercut, asymmetric crown/underside taper, explicit section anchors, and
pole-free caps. Missing any one of those invalidates the test before rendering.
