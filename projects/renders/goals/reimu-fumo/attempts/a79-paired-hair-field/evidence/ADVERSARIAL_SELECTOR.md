# A79 adversarial pre-render selector

## Decision

No published variant is eligible for Blender construction or rendering in its
current form.

| Variant | Verdict | Decisive reason |
| --- | --- | --- |
| A | **REVISE** | The three-part paired-skin ownership is credible, but its current v5 nets fail 15 executable preflight checks: front/crown convergence reaches `2.586 px`, rear convergence reaches `1.674 px`, and front/crown and leaf boundary spans exceed `3 px`. The preflight also silently accepts control-net distortion and omits roots and registered ownership. |
| B | **NO-GO** | The ear-clipped, private-centroid padding creates a triangulated quilt/facet mechanism, while the deliberately nearly planar rear pocket remains a thin card. Eight aperture-root points and a rear coverage-only mask can still pass visibly wrong geometry. |

Variant A is the only current path worth revising. Do not merge B's triangle
pocket or planar rear-face construction into A. This is a representation
eligibility decision, not a visual pass: neither variant has candidate pixels.

## Review isolation and controlling evidence

The verdicts use only the two published representation/interface specs, the
controlling references, the immutable rung-003 pixels, and the recorded
A74/A77/A78 rejection evidence. The current builders, object names as a proxy
for quality, topology cleanliness, and author self-verdicts were not used to
credit either path.

- Variant A spec SHA-256:
  `e51f1761580a9d7597cb66fd5e12afd60af913f1a7e17f1461040bd19279a52b`.
- Variant B spec SHA-256:
  `5140b74417d71b4e378a72288c74a9be7454a6f48c3fe4793f0ac23764c3300d`.
- Variant A rejected v1 literal-net SHA-256:
  `c961eddd61833b8b47b2f1169e4adc25ffc3b5fed12b10f9a1a3f9510cf1e55f`.
- Variant A superseded v2 literal-net SHA-256:
  `28d46aaee95408c33b145ff09fb4def708e45cc8bdcf702947a6996b6d31f354`.
- Variant A current schema-v5 literal-net SHA-256 (rear and leaf form revision
  v4):
  `a660ad4dfd022eb36fe6f83e8a3c7e0e36a5ce656d6a449c250e0cb8e33d39f2`.
- Variant A current pure-preflight program SHA-256:
  `4dedb474e350e93d636403e0bee42c7c73768a075c3d31c7e8e3fedee914967e`.
- Corrected interface contract SHA-256:
  `044d72a747c8779e91e4f04bcf7139a35019997eafdad705931483a08d0de886`.
- Independent density report SHA-256:
  `0e0e68b90df167cfbeb7e5d2471ea9c5f8e419720558c745c35a54d4268c9359`.
- Visual-envelope report SHA-256:
  `0c4fd2f61745ef1d77ccb004cd3533346307e4a6472224cd6f9d3a04a9d59504`.
- `canonical_front_25cm.png` controls the front crown, opening, fringe
  relationship, and face exposure.
- All 30 canonical-turn frames control continuous profile, rear, layer order,
  and the unequal lower lobes. The physical front/side photos control padded
  panel construction and edge roll only.
- A74 proves that a mechanically valid visible envelope can still be a canopy,
  card, blades, and bald rear (`2/10`, `1/10`, `1/10`).
- A77 proves that a near-concentric mantle and uniformly thick ruled leaf can
  pass coarse topology/depth checks while remaining canopy/card mechanisms.
  Its unsigned minima and control-root tests were false-pass routes.
- A78 proves that hidden receiver work cannot fix the dominant pixels. The
  complete brown crown/profile/rear owner must change atomically.

The pre-render question is therefore not whether a mesh can be closed. It is
whether the proposed representation can own every required pixel without
recreating a helmet, canopy, flat card, or bald rear, and whether its gates
would actually reject those outcomes.

## Variant A: REVISE

### Strongest case

A is materially different from the closed failures in several important
ways. It preserves the seated receiver and face, replaces exactly the six
helmet/rear owners, supplies front/crown, compact rear-base, and diagonal-leaf
roles at the first gate, and authors outer and inner skins separately without
receiver-derived Solidify or a live deformation modifier. The rear base is not
merged with the crown field, and the leaf is independently seated. Pointwise
signed roots, explicit bridge
faces, candidate/receiver crossing checks, complete rear coverage, protected
hashes, and clean reopen are all appropriate foundations.

Those facts make revision worthwhile. They do not yet authorize a build.

### 0. The first literal rear-base net failed the representation contract

The later `control_nets.json` publication reopens the decision with direct
disconfirming evidence. Its rear base is a receiver-concentric half shell:

- at `z=160 mm`, the receiver rear center is `y=51.399 mm`, while the A inner
  point is `51.927 mm` and its outer point is `55.524 mm`;
- at `z=140 mm`, the receiver is `51.227 mm`, while A is `51.799/55.206 mm`;
- the inner lateral perimeter stays near `y=0.62 mm` while its X/Z outline
  follows the receiver extrema; and
- all nine courses are nearly constant-height elliptical X/Y arcs. Across
  each full course, Z varies by only about `0.6--1.0 mm`.

Thus the inner field tracks the receiver at roughly sub-millimetre clearance
and the visible field repeats that shell only a few millimetres farther out.
Freezing generated values into JSON does not make the visible surface
independently authored. This is the A79-forbidden near-concentric,
per-height-loft mechanism and would restore a smooth rear helmet beneath the
leaf.

That v1 net was correctly discarded rather than densified. V2 limited receiver
registration to its narrow crown course and moved the non-root hanging field
away from receiver slices. Around `z=160 mm`, its inner/outer center is about
`y=58.2/61.8 mm` rather than `51.4 mm`, and its broad sides no longer wrap to
receiver tangency. It also introduces unequal lower-lobe heights and both
transverse and longitudinal variation. V2 therefore clears the specific
near-concentric blocker at the coarse-control level.

V3 increased the non-root transverse camber to roughly `14--20 mm` over the
broad middle rows while retaining longitudinal variation and the unequal
lower lobe. V4 narrowed the leaf and added convex upper-side support; these
changes are preserved in the current schema-v5 net. They materially reduce
the v2 card and oversized-leaf risks. The v5 rear has a median transverse
Y range of `15.01 mm`, best-plane residual RMS of `11.36 mm`, and all eight
non-root course centres at least `10.92 mm` beyond the receiver rear. This is
not the rejected v1 receiver copy. It still does not earn GO from coordinates
alone: evaluated refinement must preserve that broad camber, create controlled
edge roll, and avoid quilted or taut-shell pinching.

The interface contract was correctly changed from a closed `105 x 105 mm`
rear receiver perimeter to one open crown course with bounded upper-side tabs.
That correction is mandatory. The current Variant A representation spec has
not caught up: it still describes the rejected closed perimeter and an
unmodified final grid. A frozen re-entry submission must make the spec, literal
nets, root contract, applied geometry, and report describe exactly the same
construction. Seat the remaining sides through declared front/crown-to-rear-
base overlap bands, and use the semantic rear mask—not a hidden closed root—to
guarantee complete pixel coverage.

### 0a. Current v5 executable preflight is an honest REVISE, but incomplete

An implementation-independent rerun of the published pure-Python preflight on
the exact v5 net returned `REVISE` with 15 hard failures:

- front/crown level-to-level silhouette motion is `2.5857 px` in front and
  rear, `0.8164 px` in one three-quarter, and `1.1221 px` in the mirror;
- front/crown maximum straight silhouette spans are `3.1580--3.5206 px` in
  the failing views;
- rear-base motion is `1.6741 px` in the mirrored three-quarter; and
- the leaf maximum span is `3.3406--3.6682 px` in all five review views.

These are direct re-entry failures; they are not tolerances to relax. The
report is also not sufficient even if those 15 values become green. Its base
quality code permits angles of `7..173 degrees` and aspect ratio `<=28`, while
the selector requires `25..155 degrees` and `<=4:1`. The current meshes reach
aspect ratios `7.79`, `5.30`, and `6.84`, so all three would falsely pass this
weakened check. The program does not yet establish reference-registered chord
error or component masks, arc-length root seating, open-tab/overlap seating,
self/correspondence crossings, connected thickness-failure patches, or the
hash of a literal applied/frozen output mesh. It evaluates face-centre pairing
rather than the required area samples. Consequently `PASS` from this program
alone must never be interpreted as selector `GO`.

### 1. The unmodified grids are visibly too coarse

Smooth shading does not move a silhouette. The front/crown horseshoe has only
eight boundary spans along its left-temple-to-crown-to-right-temple direction.
The canonical front field is about `128.8 mm` wide. Even the straight-width
lower bound gives an average span of `16.1 mm`; under the frozen `292 mm`
orthographic scale at 512 px, that is about **28 px per straight edge**. The
actual horseshoe arc is longer, so its crown/temple facets would be larger.

The independent conservative density calculation gives even larger average
spans once the complete controlling directions are included: about `50.9 x
32.8 px` for front/crown, `35.0 x 28.4 px` for rear base, and `44.3 x 21.2 px`
for the leaf. Normal interpolation can conceal face shading changes but cannot
round these outer contours, nape lobes, the diagonal leaf tip, or overlap
boundaries.

Required revision: no identity-owning curved boundary may have a projected
straight span longer than 3 px at 512 px, and at most 2 px is required around
crown extrema, temple turns, lobe tips, and the diagonal leaf edge. Uniform
tensor tessellation would need about `137 x 67`, `95 x 77`, and `90 x 44`
evaluated vertices per skin, respectively, but those are sampling equivalents,
not mandatory authoring-cage counts. A sparse Catmull/B-spline cage is valid
when adaptive evaluated tessellation actually passes every projected span,
chord, convergence, and anti-pinching gate. Add local knots/support rows only
where a measured feature fails. Subdivision is sampling only: both paired
cages remain independently authored, the saved pocket has no live modifier,
and refinement may not create the form, a receiver offset, or a new extremum.

### 2. A rectangular grid is not yet a horseshoe proof

A topological rectangle can parameterize a horseshoe strip, but the prose
indexing alone does not show that this particular one can. Inner opening and
outer crown arcs have materially different lengths, and the temple turns
compress one side of the quads while stretching the other. With only eight
arc intervals, the result can fold, pinch, fan into armor panels, or develop a
square/canopy transition while retaining valid Euler counts.

Required revision: publish every literal outer and inner control coordinate
and a pure feasibility report covering:

- consistent signed quad Jacobians and zero foldovers on both skins;
- minimum/maximum quad angles and aspect ratios, including both temple turns;
- monotone courses from the opening arc to rear overlap;
- outer/inner correspondence segments that stay locally inward, never cross
  either skin, and do not make the bridge self-intersect;
- arc-length and projected-pixel spacing on every visible boundary; and
- the intended `1.5--3.5 mm` depth step in both profiles and three-quarters.

Manifoldness, Euler characteristic, and corresponding-pair distance do not
answer these questions.

### 3. Complete ownership is asserted, not yet bounded

The literal nets now exist, but neither the spec nor the current preflight has
reference-registered component masks. Per-object projected vertex bounds are
not a silhouette comparison and contain no occlusion or ownership evidence. A
large rear shield can pass
"zero beige rear pixels" while becoming a curtain; an over-wide crown panel
can pass coverage while remaining a helmet; a leaf can own the farthest pixel
while still being a flat blade.

Required revision: preflight the literal nets through the frozen camera
matrices and record, before Blender:

- front scanline widths and allowed face exposure;
- both asymmetric profile partitions, with compact field and free leaf
  reported separately rather than only as a union;
- rear width `.94 +/- .04 Wh`, height `1.16 +/- .05 Wh`, maximum-width height
  `.70 +/- .04` of rear height, unequal lobe extrema, and diagonal overlap;
- object-ID ownership for crown, rear base, and leaf in front, both
  three-quarters, both profiles, and rear; and
- zero rear face/applique/receiver visibility, without counting the bow,
  retained brown parts, background, or an oversized candidate as coverage.

Use unfiltered per-object ID masks for the mechanical test. A shared brown
material-ID mask cannot prove which part owns a pixel or distinguish a leak
hidden by another brown object.

The canonical-turn masks are per-view evidence, not a calibrated 3D hull.
Their 15 antipodal pairs disagree in width by `12.54%` median and `38.35%`
maximum, so intersecting them into a volume would encode bow occlusion and
perspective error as false geometry. After valid per-view alignment, use them
only to reject an inward miss, an outward overshoot above roughly `5%` of the
front hair width (`13 px` in the source turn), loss of frame-18 continuous
rear occupancy, or loss of the separate trailing profile island. The RGB
source does not resolve same-brown internal panel order; construction
references control that relation.

### 4. The root and thickness gates retain false-pass routes

`80%` signed coverage permits a contiguous fifth of a root to float. On a
roughly `129 mm` course that can be about `26 mm`, or roughly 45 px at the
512 px review scale. The separate one-pixel sentence does not define sample
density, weighting, projection, or treatment of a contiguous failed run.

Likewise, `95%` paired thickness at corresponding points can pass a folded
pair, a locally inverted skin, a thin bridge, or the wrong nearest layer. A
Euclidean pair distance alone does not show that the inner skin is behind the
outer skin or that thickness remains valid over triangle interiors.

Required revision:

- sample roots uniformly by evaluated arc length, not vertex count;
- report signed `min/P05/P50/P95/max`, total passing arc length, endpoints,
  and the longest contiguous failing run in millimetres and every review view;
- require every root sample positive, no candidate/target crossings, and no
  failing run projecting wider than one unfiltered review pixel;
- check the retained fringe/temple composite as well as the receiver, because
  receiver clearance alone cannot prove a seated visible root;
- measure thickness over area-weighted triangle samples using the declared
  outer/inner correspondence and local orientation, not unrestricted nearest
  distance;
- require positive outer-to-inner normal direction, no correspondence segment
  crossing, separate bridge-width/roll checks, and no connected out-of-band
  patch wider than one review pixel; and
- run every check on the exact saved, clean-reopened evaluated geometry.

A may receive **GO** after all four revisions pass without changing its useful
three-part, independently authored paired-skin representation.

### Exact A re-entry gates

The revised A spec and pure report must expose these values directly. A prose
claim or a later beauty render is not a substitute.

| Gate | Required result before Blender |
| --- | --- |
| Artifact coherence | Freeze and hash one revised spec, literal control-net file, executable report, exact applied/frozen geometry, and corrected interface contract. They must describe the same topology, root bands, refinement level, coordinates, and no-live-modifier state. A stale or moving input is `REVISE`, and a pure-report `PASS` is not itself selector `GO`. |
| Boundary density | In every frozen 512 px view, each straight evaluated boundary span is `<=3 px`; spans containing a crown extremum, temple turn, nape-lobe tip, leaf tip, or overlap corner are `<=2 px`. Uniform-tensor equivalents are front/crown `137 x 67`, rear base `95 x 77`, and leaf `90 x 44` evaluated vertices per skin, but a sparse cage may pass through adaptive applied tessellation. Counts never override the measured projection. |
| Chord error | Against the camera-aligned controlling silhouette, maximum normal chord error is `<=1 px` at the identity-critical extrema above and `<=2 px` on every other visible curved boundary. Report maximum and P95 per component and view. |
| Horseshoe map | Both skins have strictly positive signed Jacobians, zero folded or degenerate quads, angles in `25..155 degrees`, aspect ratio `<=4:1`, and adjacent-face area ratio `<=4:1`. Courses remain ordered and monotone from the face root to the rear overlap. |
| Anti-pinching | No extraordinary/refinement vertex lies on an identity-critical silhouette, root, temple turn, lobe tip, or overlap corner. At all other extraordinary vertices, evaluated one-ring normals stay within `20 degrees` of their area-weighted mean and no incident face has less than `25%` of the one-ring median area. |
| Non-concentric rear field | Receiver fitting is limited to the declared open crown root and bounded upper-side tabs. Fit the best single receiver-normal offset to all non-root outer samples and report area-weighted residuals; reject when `>=80%` of visible area lies within `+/-1.0 mm` of that fitted offset. Publish non-root coordinates before fitting and prove that no inner/outer course is obtained from receiver height slices. |
| Refined closed pocket | If any geometric subdivision or adaptive refinement is introduced, validate the applied/evaluated result, not its cage: one closed two-manifold component per part, zero boundary/wire/overused edges, consistent positive volume, zero self-intersections, and explicit bridge continuity. Refinement may move a frozen-view silhouette by at most `1 px` from its frozen revised-form target and may not create a new extremum. One additional diagnostic evaluation level must change every silhouette by `<=0.5 px`, create no new extremum/fold/pinch, and leave the thickness pass/failure classification unchanged. |
| Evaluated thickness | After all refinement is applied, area-weighted declared outer-to-inner thickness is `2.0..5.0 mm` on at least `95%` of visible area; all samples remain positive; `min/P05/P50/P95/max` are reported; no connected failure patch projects wider than `1 px`; correspondence segments and bridge faces do not cross either skin. |
| Root seating | Uniform evaluated arc-length samples are all positive; at least `80%` lies inside the component band; the longest contiguous out-of-band run is `<=0.57 mm` in the view where it projects widest and never exceeds one unfiltered pixel; endpoints pass; candidate/receiver and candidate/retained-root crossings are zero. |
| Rear attachment | Rear-to-receiver contact is an open crown course plus bounded upper-side tabs, never a closed perimeter or a `105 x 105 mm` wrap. Both open-course endpoints pass. Declared front/crown-to-rear-base side bands supply the remaining support, with zero floating component and the same contiguous-run rule. Rear pixel coverage is tested separately by masks. |
| Front | Registered scanlines and face exposure remain inside the frozen A79 bands; the retained fringe, not a new rectangular edge, owns the opening silhouette. |
| Profile partition | Compact field depth is provisionally `.77..85 Wh`; complete brown depth is `1.14 +/- .05 Wh` on side A and inside `1.13..1.29 Wh` on side B. The independent leaf owns the farthest depth over a contiguous visible height of at least `.25 Wh`; the rear base remains separately visible. |
| Rear bounds | Brown union width is `.90..98 Wh`; crown-to-lowest lobe is `1.11..1.21 Wh`; maximum width occurs at `.66..74` of rear height. The dominant leaf is `.45..50 Wh` wide with an approximately `.10 Wh` off-centre bias, treated as a low-confidence bracket rather than a construction coordinate. |
| Layer order | The front/crown-to-rear-base depth step evaluates to `1.5..3.5 mm` and projects as a distinct `2..6 px` boundary in both three-quarters and both profiles. The leaf, rear base, and front/crown component-ID masks each own a nonzero contiguous region; no union-only mask may pass this gate. Require the source-supported bow-over-hair and trailing-leaf-beyond-field relations; do not claim same-brown hidden order from the RGB turn. |
| Rear ownership | Exact unfiltered object-ID masks show zero receiver, face-applique, or background island inside the registered target hair region; no gap or disconnected island exceeds one pixel. After valid per-view alignment, reject an inward miss, an outward miss above `5%` of front hair width, or loss of the continuous frame-18 middle/lower coverage and scalloped taper. Overshoot beyond the registered rear/profile bands is a failure, so an oversized shield cannot pass coverage. |

The profile and leaf bands inherited from the turntable audit are provisional
perspective brackets. Camera alignment must be within the landmark tolerance
before they can pass; they may not be converted directly into hidden receiver
coordinates.

## Variant B: NO-GO

### Strongest case

B correctly chooses the complete six-object replacement boundary and includes
all three required visual roles before a render. Closed paired pockets, exact
hide preservation, signed roots, explicit bridges, zero Solidify/pressure,
candidate crossing checks, and a clean-reopen gate are improvements over A74
and A77. Literal cut patterns are also plausible in the abstract for plush
hair.

The actual surface mechanism defeats that case.

### 1. Private triangle centroids create the wrong surface class

Ear clipping chooses implementation-dependent diagonals across a concave
outline. B then gives every original triangle its own private displaced
centroid and fans three planar triangles around it. Adjacent original
triangles share an edge but not a camber control. This produces tangent breaks,
triangle-scale pillows/dimples, and triangulation-direction bias. With no
geometric refinement, smooth shading changes normals only; it cannot remove
the piecewise-planar silhouette or the centroid quilt from diagnostic light.

That is not a controlled low-frequency stuffed panel. It risks reading as
faceted armor or a quilt of unrelated pads, especially where ear clipping
threads the concave U and unequal rear lobes. Passing genus, edge-incidence,
and `2--5 mm` paired-point checks would not veto that visual failure.

### 2. The rear base is explicitly still a card

The spec calls the complete rear visible face "broad and nearly planar."
Closing a planar sheet with a second skin and bridge gives it finite thickness,
but does not supply the transverse and longitudinal stuffing camber, edge roll,
tension, or overlap deformation needed to stop a card read. A74's face card
and A77's ruled Solidify leaf were rejected for the same visual reason, not
because they happened to have an open boundary.

The independent diagonal leaf does not rescue a planar shield that owns most
rear pixels. This is a representation-level recurrence, so local density or
gate edits are insufficient.

### 3. Eight front-root points can falsely pass

The protected aperture loop has 196 vertices and a `459.762 mm` 3D perimeter.
Eight selected points replace curved contact with long chords. Interpolating
64 samples along those chords only measures the coarse approximation more
often; it does not recover the skipped aperture or retained-fringe geometry.
The band can satisfy receiver distance while cutting across a curved root,
leaving visible gaps, or seating against the wrong layer.

The same `80%` loophole permits long contiguous failures, and the rear mask
only rejects beige gaps. An oversized planar pocket can pass it while failing
reference width, maximum-width height, lobes, diagonal overlap, and layer
separation.

### Re-entry condition

B must not proceed by adding more centroid triangles or padding the planar
shield. Reconsider it only after a representation reset that:

1. replaces private triangle-centroid camber with one shared, smoothly varying
   outer field and a separately authored inner field;
2. gives the rear base measured transverse and longitudinal camber, edge roll,
   unequal lobe tension, and non-card silhouette ownership;
3. samples the aperture/root by arc length densely enough for <=3 px projected
   spans and validates it against retained fringe/temple layers;
4. adds the reference-registered multi-view ownership bounds and contiguous
   root/thickness failure limits required above; and
5. proves the revised literal surfaces before any Blender build.

Those changes replace B's defining topology algorithm. Until then the verdict
is **NO-GO**, not a bounded correction.

## Selection and next action

Revise A only. Re-run this selector on its frozen revised spec and pure
feasibility result before authorizing Blender. If A cannot produce smooth,
non-folded literal paired nets while meeting the reference-registered
ownership and contiguous-failure gates, change representation again. Do not
spend a render to discover a defect already implied by the control surface.
