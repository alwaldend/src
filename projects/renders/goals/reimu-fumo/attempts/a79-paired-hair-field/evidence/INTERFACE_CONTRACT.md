# A79 paired hair-field interface contract

## Status and scope

This freezes the mechanical boundary for A79 before either representation
variant builds a candidate. It is a contract, not a candidate: no source
object, goal record, `.blend`, camera, material, or animation was changed.

The controlling machine-readable file is
`interface/interface_contract.json`. The measurements come from the protected
rung-003 source at SHA-256
`c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`,
opened read-only by pinned Blender 5.2.1. A78 supplies the exhaustive receiver,
aperture, object, and camera inventory. A79 adds only a bounded old-cap versus
retained-panel probe in `interface/retained_band_baseline.json`; it does not
repeat the broad inventory.

World axes are X left/right, negative Y front, positive Y rear, and Z up.
Blender stores metres; this contract reports millimetres.

## Exact visibility and object boundary

The candidate must set both `hide_viewport` and `hide_render` to `true` on
exactly these six existing objects:

1. `A44 continuous hair cap with smooth opening`
2. `A42 Left asymmetric rear lock`
3. `A42 Off-center main rear lock`
4. `A42 Short right rear lock`
5. `A42 Main lock left seated seam`
6. `A42 Main lock right seated seam`

Those two flags are the only permitted source-field changes. The objects and
their datablocks remain present and otherwise unchanged. Every other one of
the 177 A77 source objects keeps its baseline type, transform, visibility,
materials, base counts/bounds, and modifier identities. The 37 A78 interface
objects additionally keep exact object-state, base-geometry, and evaluated-
geometry hashes. All five review cameras and the `800 x 800` render state keep
their A78 hashes.

The render-visible replacement is exactly three fresh, single-user,
identity-transformed mesh objects in collection `A79 Paired hair field`:

- `A79 Front crown padded field`
- `A79 Compact rear base field`
- `A79 Broad asymmetric rear leaf`

No proxy, helper, seam, receiver, or fourth render-visible A79 object may
enter P0. Helpers belong in builder data, not the saved scene.

## Explicit paired-skin topology

Each component is one connected closed two-manifold with no modifier or shape
key. The base mesh is the evaluated P0 mesh; this prevents a subdivision or
Solidify step from changing the measured envelope after the builder passes.

Every polygon carries FACE/INT attribute `a79_surface_role`: outer `1`, inner
`2`, or bridge `3`. Every vertex carries POINT/INT attributes `a79_pair_id`
and `a79_skin`; each positive pair id occurs once on outer skin (`1`) and once
on inner skin (`-1`). Every outer triangle has one inner triangle with the
same three pair ids and opposite winding. Bridge faces use paired boundary
vertices only.

The validator rejects non-finite coordinates, vertices coincident within
`1e-7 m`, duplicate faces, triangles at or below `1e-12 m^2`, loose or
non-two-face edges, disconnected geometry, inconsistent/negative volume,
unintended self-intersections, shared datablocks, non-identity transforms, or
any modifier. Adjacent triangles are ignored by the self-intersection test
only when they share a mesh vertex; all other BVH overlap candidates receive
an explicit 3D triangle/coplanar overlap test.

Thickness is the paired outer/inner triangle-centroid separation, weighted by
outer world-space triangle area. Every component must place at least 95% of
that area inside `2.0--5.0 mm`. The paired displacement must retain at least
80% of its magnitude along the outer normal, and every sample, including edge
roll, must remain inside the hard `0.5--7.0 mm` guard.

## Root-band semantics

“80% coverage” is never a minimum-vertex shortcut.

- A retained band is the immutable set of all evaluated retained-object
  vertices and triangle centroids whose baseline unsigned distance to the old
  cap was at most `2.0 mm`. The denominator and bounds below are frozen.
- For a retained sample, signed clearance is
  `(sample - nearest_candidate_outer) dot candidate_outer_normal`. Positive
  means the retained fringe, cheek lock, or seam remains on top of the new
  field rather than being swallowed by it.
- Receiver distance uses the receiver's outward winding. Positive means the
  candidate inner skin is outside the receiver.
- Candidate vertex-group roots are resampled uniformly by arc length; raw
  vertex density cannot increase the passing fraction.
- Every fraction rounds upward to an integer requirement.

### Retained front interface

All eight bands target `A79 Front crown padded field`.

| Band | Frozen samples | World bounds min -> max mm | Required positive signed clearance | Required | Old-cap crossing ceiling |
| --- | ---: | --- | ---: | ---: | ---: |
| left temple fringe | 265 | `(-55.992,-54.411,150.955)` -> `(-20.411,-45.382,192.486)` | `0.20--2.50 mm` | 212 | 0 |
| left temple transition | 257 | `(-62.517,-46.072,109.162)` -> `(-54.342,-11.566,166.106)` | `0.20--2.50 mm` | 206 | 377 |
| main bang | 149 | `(-41.988,-55.994,171.760)` -> `(18.273,-48.961,196.084)` | `0.20--2.50 mm` | 120 | 0 |
| right swept fringe | 281 | `(8.143,-56.004,164.110)` -> `(55.155,-45.607,194.930)` | `0.20--2.50 mm` | 225 | 0 |
| right temple transition | 260 | `(54.025,-46.109,109.195)` -> `(62.491,-11.641,166.052)` | `0.20--2.50 mm` | 208 | 388 |
| left cheek root | 83 | `(-61.398,-44.922,143.176)` -> `(-54.439,-31.065,160.039)` | `0.20--2.00 mm` | 67 | 92 |
| right cheek root | 91 | `(54.457,-45.580,143.175)` -> `(62.398,-31.066,160.039)` | `0.20--2.00 mm` | 73 | 114 |
| crown seam | 1,136 | `(-0.100,-54.541,177.950)` -> `(0.100,-25.682,219.094)` | `0.05--0.60 mm` | 909 | 225 |

The crossing ceiling is baseline-relative, not permission to create a new
crossing family. A candidate assembly's counts against each retained object
sum across all three new components and may not exceed the listed old-cap
count. The three zero-baseline retained pairs are hard zero. Positive signed
coverage must pass independently even when a nonzero legacy ceiling exists.

### Receiver and aperture registration

The frozen cap face-opening boundary is one closed 196-point loop, perimeter
`459.761558 mm`, bounds
`(-56.946568,-54.127648,91.123715)` through
`(56.946568,-10.791719,180.005699)`. Every source point was outside the
receiver at `+0.773125..+0.779999 mm`.

`A79_ROOT_RECEIVER_FACE_LOOP` on the front/crown inner skin must register at
least 157 of those 196 anchors within `1.0 mm`. The corresponding candidate
samples must be `+0.20..+0.80 mm` outside the receiver. The band is a single
closed chain; crossings with the receiver are zero.

The rear base is **not** perimeter-glued. Its receiver attachment is three
open courses on the inner skin:

- `A79_ROOT_RECEIVER_REAR_CROWN` is an open crown course, resampled to 48
  equal-arc points. At least 39 points and both endpoints must be
  `+0.20..+0.80 mm` outside the receiver. Its length is `80--160 mm`, X span
  `65--115 mm`, every sample has Z at least `175 mm`, and Z span is at most
  `45 mm`. There is deliberately no absolute Y floor: signed receiver
  clearance controls seating at the high lateral crown, where the receiver
  surface itself is near Y `8 mm`.
- `A79_ROOT_RECEIVER_REAR_TAB_L` and `_R` are independent open upper-side
  tabs, each resampled to 16 equal-arc points. At least 13 points and both
  endpoints must pass `+0.20..+0.80 mm`. Each tab is only `8--35 mm` long and
  spans at most `30 mm` in Z. The left tab stays within X `-66..-35`, Y
  `15..52`, Z `155..210 mm`; the right tab mirrors the X window at `35..66`.

All three courses have zero receiver crossings. The open course and bounded
tabs encode the canonical hanging-panel attachment under the bow/crown; they
cannot descend into a near-concentric receiver wrap. Complete rear coverage
is owned by the semantic silhouette mask below, not by expanding the contact
band around the receiver.

### New-component root bands

The front/crown-to-rear-base and rear-leaf-to-rear-base joints use the exact
vertex-group pairs named in the JSON contract. Each side is independently
resampled to 48 equal-arc points; at least 39 per direction must have positive
`0.20--0.80 mm` signed clearance to the receiving outer skin. Pairwise
triangle crossings are zero. This is a tucked layer gap smaller than roughly
two rear-camera pixels, not a shared/intersecting shell.

## Baseline-relative crossing contract

Every candidate component has zero exact triangle intersections with
`Head_Cushion_Manual_Target`; all three candidate/candidate pairs are also
zero. This deliberately improves on the removed rear locks, whose receiver
crossing counts were 279, 0, and 44. Their bad intersections are context, not
an inherited allowance.

Intersection counts are exact evaluated triangle-pair counts after BVH AABB
broad phase, segment/triangle narrow phase, and coplanar overlap handling. A
nearest vertex, unsigned minimum, AABB-overlap count, or raw BVH overlap count
cannot satisfy this gate.

## Aperture regression

The seven face witnesses remain exact. Under `Review_front_Camera`, the new
assembly may not reduce their A78 evaluated-vertex visibility fractions:

| Witness | Minimum visible fraction |
| --- | ---: |
| left / right eye applique | `0.98742747` / `1.0` |
| left / right half-lid stitch | `1.0` / `1.0` |
| left / right upper expression stitch | `0.27586207` / `0.60344828` |
| mouth dash | `1.0` |

This prevents a mechanically seated crown from stealing face pixels.

## Rear-coverage semantic mask

The mask is geometric, not an RGB threshold. It uses one ray through every
pixel center of the unchanged `Review_rear_Camera` at `800 x 800`; the camera
hash is `6039de82056750b2d754f0c3ea6595d945fd778ae0e33b7ac29e4c1bfd947908`.
The receiver projection bbox is approximately
`(219.178,152.877)--(580.822,505.480)` pixels.

The ROI is the receiver's isolated rear silhouette minus pixels where an
unchanged non-hair source object is closer to the camera. “Hair” means object
identity: the three A79 components plus the retained five fringe panels, two
cheek locks, and crown seam. Hidden legacy objects never count. A
`support_leak` pixel has no hair-owner hit closer than the receiver; a
`beige_leak` is a support leak landing on receiver material
`Face fabric clay`.

Pass requires:

- zero beige-leak pixels in the full ROI;
- zero support-leak pixels after one 8-neighbour erosion of the ROI boundary;
- one 8-connected hair-coverage component after discarding only one-pixel
  speckles; and
- in a two-pixel dilation of every component root, the support-leak mask is
  empty after one erosion, so no root gap is two pixels thick.

The validator must publish `receiver_roi.pgm`, `hair_coverage.pgm`,
`support_leak.pgm`, `beige_leak.pgm`, and `root_gap.pgm` with its JSON report.
This distinguishes brown-looking exposed support from actual hair ownership
and makes “bald rear” machine-detectable.

## Required validation order

1. Verify the contract and all referenced inventory digests.
2. Compare the complete source object set/state to A77; compare interface
   states/geometries and fixed cameras to A78.
3. Enforce the exact visibility delta and exact three-object candidate set.
4. Fail topology, attribute, winding, pairing, thickness, or self-intersection
   defects before contact checks.
5. Check retained, receiver, and component root bands using every frozen or
   equal-arc sample.
6. Check exact crossing ceilings.
7. Check aperture/witness regression and rear semantic masks.
8. Reopen and rerun against the exact saved candidate with pinned Blender;
   command success alone is not a pass.

Any failed row is fatal before rendering. Passing this interface contract is
technical eligibility only; it cannot override the A79 helmet/card/bald-rear
self-veto or the fixed-view implementation-blind visual gate.
