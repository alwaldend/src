# A75 face-cushion topology design

## Verdict

The one-object A75 boundary is defensible only as a **registered visible
face-field retopology**. It is not a credible way to rebuild the macro head,
reduce the profile depth, change the rear volume, or repair hair seating.

The frozen A44 cap is an almost exact offset copy of the current receiver. Its
used base vertices are `0.7794-0.7800 mm` from the old receiver, and its
evaluated surfaces are only `0.5396-1.0200 mm` outside it. The fixed eyes are
only `0.076/0.080 mm` clear of the receiver. A materially different whole
cushion would therefore either cross the cap and appliques or leave them
floating.

This is also an observability limit, not merely a geometry inconvenience.
Earlier isolated-versus-assembled masks exposed only `27.73%` of the receiver
in front, `16.25/16.12%` in the two three-quarter views, `3.56%` in profile,
and `0%` at the rear. Front and three-quarter renders can judge an improved
face carrier. Profile and rear renders can only veto leakage, clipping, or
integration regressions. If A75 is expected to improve the macro profile or
rear silhouette, stop before Blender and reject the module boundary.

The conditional topology below is the smallest cage that can independently
support both fixed eye footprints, the mouth, a broad cheek roll, and a hidden
transition collar without a center pole. It must not be promoted merely for
passing topology checks.

## Controlling evidence

Use `Wh = 132 mm`, the current receiver width. Front is world negative `Y` and
`Z` is up.

| Source | What it controls for A75 | What it does not prove |
| --- | --- | --- |
| `canonical_front_25cm.png` | Exact-variant face exposure, frontal identity, graphics, and the 25 cm whole-plush scale | Hidden cushion perimeter, depth, or seams |
| `canonical_turn_180.gif` | Composite layer order and the absence of beige leakage in both profiles, both three-quarters, and rear | An isolated cushion outline or factory pattern |
| `physical_front.png` | Soft cheek/chin roll and low-frequency stuffing | Exact-variant dimensions |
| `physical_side.png` | Near-planar face, thin cap-over-cushion seating, and soft overlap | Orthographic depth |
| `clean_front.png` | Eye, lid, and mouth placement cross-check | Contact, compression, or hidden construction |
| `sofa.gif` | Pile-softened edges, applique seating, and broad compression | Metric silhouette |
| `turn.gif` | Qualitative continuity through side and rear | Fine dimensions or exact variant identity |

The source-measured visible beige field is `.603 +/- .030 Wh` in both axes,
or `79.6 +/- 4.0 mm`. That is a visible-pixel region, not a literal panel seam.
The larger cap-free receiver region is only a construction bound:
`X=-54.300..54.300`, `Y=-54.201..-13.977`, and
`Z=92.200..177.462 mm` (`.823 Wh` wide and `.646 Wh` high).

## Parent and frozen interfaces

| Item | Evaluated world-space bounds or contact |
| --- | --- |
| `Head_Cushion_Manual_Target` | `X=-66..66`, `Y=-54.201..51.399`, `Z=91.500..220.200 mm`; `132 x 105.6 x 128.7 mm` |
| A44 cap | `X=-67.020..67.020`, `Y=-55.062..52.419`, `Z=90.480..221.220 mm`; `134.04 x 107.48 x 130.74 mm` |
| Left eye applique | `X=-39.940..-11.102`, `Y=-54.184..-48.214`, `Z=124.769..147.998 mm`; nearest receiver clearance `0.076 mm` |
| Right eye applique | `X=11.102..39.940`, `Y=-54.187..-48.211`, `Z=124.769..147.998 mm`; nearest receiver clearance `0.080 mm` |
| Lid stitches | `Z=149.048..150.447 mm` |
| Upper expression stitches | `Z=157.940..159.161 mm` |
| Mouth | `X=-1.802..1.802`, `Y=-46.894..-46.700`, `Z=110.685..111.047 mm`; nearest receiver clearance `0.320 mm` |

The old receiver already intersects temple panels, cheek locks, bow loops,
some rear leaves, the crown seam, and the bodice. A75 must not claim that
preserving those crossings is clean construction. It may only require that no
new connected crossing is introduced and that every baseline count stays the
same or decreases.

## Exact base representation

Create one closed mesh object named `A75_Face_Cushion`:

- one `9 x 9` all-quad front chart: `81` vertices and `64` quads;
- one `9 x 9` all-quad rear chart: `81` vertices and `64` quads;
- two independent `32`-vertex gusset loops between the front and rear
  boundaries: `64` vertices;
- three `32`-quad bridge strips: front to gusset A, gusset A to gusset B, and
  gusset B to rear, for `96` quads; and
- total base topology: **226 vertices, 448 edges, and 224 quads**.

The Euler characteristic is `226 - 448 + 224 = 2`. The mesh is one closed
genus-zero component. It has no triangles, ngons, center vertex, radial fan,
or UV-sphere pole. Front-field vertices have valence four. The only eight
valence-three grid corners are on the hidden front and rear seam collars.

Map each chart from a regular parameter grid with the concentric square-to-disk
mapping

```text
p = s * sqrt(1 - t^2 / 2)
q = t * sqrt(1 - s^2 / 2)
```

before fitting it to the measured receiver. This supplies a regular
non-radial quad field. Do not start from a cube, UV sphere, beveled extrusion,
superellipsoid, homothetic outline rings, or a center-fan disk. Those
representations repeat the square-card, mattress, egg, helmet, or pole-pinch
failures already observed in A44, A72, A73, and A74.

Use nonuniform front control rows/columns so the cage explicitly represents
the two eye footprints, eye-to-mouth field, lower cheek, and hidden collar.
Nine stations are the minimum that permits a frozen outer band, a transition
band, two independent eye domains, a central mouth domain, and a free cheek
field without sharing one pole or one rigid plane.

## Deterministic fitting algorithm

1. Open the exact protected rung-003 parent and verify its SHA-256 before any
   candidate work. Hide only `Head_Cushion_Manual_Target` by exact name; keep
   it in the file as a hidden registration target and rollback witness.
2. Create the 226-vertex cage in world coordinates with identity object
   transform. Give it a new single-user mesh and a new single-user copy of the
   beige material.
3. Sample the old evaluated receiver for the rear chart, both gusset loops,
   the complete front perimeter, and the outer two front-grid bands. These
   vertices own only hidden cap/hair registration, not a new silhouette.
4. Fit the central front field with a smooth elliptical weight whose value and
   first derivative are zero at the transition collar. The face is the front
   skin of the same volume; it is not an inset disk, face plate, or second
   object.
5. Keep the old evaluated surface as exact interpolation anchors beneath both
   eye footprints, lids, expression stitches, and mouth. Do not infer their
   required support depth from object origins.
6. Author a broad, mildly convex face plateau directly in the base cage. Its
   low-frequency relief across the exposed face should be `.005-.019 Wh`
   (`0.7-2.5 mm`) and must not exceed `.025 Wh` (`3.3 mm`). Roll into the
   cheek over at least `.08 Wh` (`10.6 mm`); do not form a vertical wall, lip,
   or rim.
7. Use the old receiver only as a build-time registration target. Iterate a
   fixed limit-surface fit on hidden controls until the evaluated subdivision
   surface reaches the registration gates below. Do not leave a live
   Shrinkwrap dependency in the saved P0 candidate.
8. Permit at most `.012 Wh` (`1.6 mm`) of low-frequency left/right cheek
   asymmetry, and only if every frozen-applique clearance still passes. Do not
   add random noise or decorative wrinkles at P0.

The following are directional flattening gates, not claims about an isolated
source seam. Measured on the evaluated face from its frontmost point, require
at least:

| Depth behind frontmost surface | Available face width |
| ---: | ---: |
| `2 mm` | `.58 Wh` (`76.6 mm`) |
| `4 mm` | `.68 Wh` (`89.8 mm`) |
| `6 mm` | `.74 Wh` (`97.7 mm`) |

The old receiver supplies only `.448`, `.525`, and `.603 Wh` at the same
depths. Failure to improve those widths without violating the frozen
interfaces proves the one-object module invalid; it is not a reason to push
the cushion through the hair.

## Seams and openings

- The mesh is closed and has no geometric face opening.
- Treat the front-chart/gusset boundary as a plausible hidden sewn collar
  beneath the hair, but do not claim that the references reveal this factory
  seam.
- Store named vertex groups for `A75_face_plateau`, `A75_transition`,
  `A75_hair_contact`, `A75_applique_support`, `A75_gusset`, and `A75_rear`.
- Add no visible seam curve, piping, groove, torus, inset rim, or material
  border at P0. If the form later passes and a seam becomes useful, represent
  it as a broad `0.5-0.8 mm` compression within the hidden collar.

## Modifier plan

The P0 object should have exactly one geometry modifier:

1. Catmull-Clark Subdivision Surface, viewport level `1`, render level `2`.

Use smooth shading. Do not add Solidify, Bevel, Remesh, Displace/noise,
Corrective Smooth, Weighted Normal, or a live Shrinkwrap. The object is already
a closed stuffed volume; Solidify would create another nested shell. The base
cage must own the plateau, cheek roll, crown/nape taper, and gusset transitions.
Subdivision is relaxation, not a mechanism for turning a box into a cushion.

## Mechanical gates

Reject before beauty rendering unless all of these pass on the evaluated
candidate:

1. Exactly one new visible mesh object and no other visible additions.
2. Base topology is exactly `226 V / 448 E / 224 F`; every face is a quad;
   one connected component; zero boundary, wire, non-manifold, duplicate, or
   degenerate elements.
3. No center pole, vertex fan, or base valence greater than four.
4. Every frozen object and transform is byte-equivalent to the parent. The old
   receiver alone changes visibility.
5. Hidden collar, gusset, and rear evaluated samples target the old receiver
   with `<=0.15 mm` RMS and `<=0.25 mm` maximum displacement.
6. Cap/receiver crossings remain zero. Sampled cap clearance stays
   `0.30-1.30 mm`, with no coherent gap or bulge relative to the baseline
   `0.5396-1.0200 mm` band.
7. Both eyes remain noncrossing and `0.05-0.25 mm` from the evaluated face.
   The mouth remains noncrossing and `0.15-0.50 mm` away. Lid and expression
   stitches acquire no new crossing.
8. No new connected intersection component occurs against any frozen object;
   each existing per-object triangle-crossing count stays equal or decreases.
9. A material-ID mask reports zero beige pixels outside the frozen brown-hair
   silhouette in front, rear, both three-quarters, and both profiles.
10. Fixed-camera difference masks show no new visible gap wider than one pixel
    at review resolution.

## Render and early-stop gates

Render an immutable P0 in the complete parent context from front, both
three-quarters, both profiles, and rear. Attach the canonical front, all
relevant canonical-turn neighborhoods, physical front/side, clean front,
sofa, and older-turn evidence to the blind review packet.

The positive preference test is the front plus the worse three-quarter view.
Profile and rear are regression tests because the receiver is almost entirely
occluded there. Stop immediately for any of the following:

- square corner, straight plaque edge, inset rim, card-thin profile, egg
  bulge, mattress side wall, helmet read, pole dimple, or visible faceting;
- beige crown/rear leakage, a visible face-plate edge, or a new hard tangency;
- floating or crossing eyes, mouth, hair, bow, collar, or bodice;
- a broad face can be produced only by moving frozen hair/appliques or by
  violating their clearances;
- the face region is not clearly preferred to rung 003 in front and the worse
  three-quarter view; or
- whole-model failure is still dominated by the unchanged hair shell.

That last outcome rejects the **module boundary**, not merely the parameter
values. Do not reshape the cushion to compensate for frozen hair. An internal
P0 may advance only when macro identity, plush construction, and contact each
score at least `6/10`, with no categorical veto. Final approval still requires
the repository visual-quality gate: unlabeled recognition, every applicable
category at least `8/10`, and no major visible failure.

## Evidence used

- A66 head/receiver inventory and exact contact audit;
- A69 interface, observability, and reference reports;
- A44 head-envelope and head-only FFD rejection records;
- A72 and A73 receiver build/review packets;
- A74 measurement packet and blind rejection;
- the tracked canonical and supporting reference packet; and
- two independent A75 read-only dimension/topology reviews.
