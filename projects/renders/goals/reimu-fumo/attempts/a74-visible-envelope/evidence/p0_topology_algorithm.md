# A74 P0 deterministic topology algorithm

## Outcome and non-negotiable boundary

Build exactly four candidate objects in a disposable copy of the protected
rung-003 scene:

1. `A74_CrownTemple_Saddle`;
2. `A74_Left_Cheek_Lock`;
3. `A74_Right_Cheek_Lock`; and
4. `A74_Subordinate_Face_Support`.

The three brown objects own the visible P0 hair pixels.  The beige support is
only a face-applique carrier.  It is not a new head receiver, and P0 makes no
claim about its hidden back, seam, or manufacture.

The builder must hide these 15 legacy objects by exact equality, never by
collection, prefix, substring, material, or type:

```text
Head_Cushion_Manual_Target
A44 continuous hair cap with smooth opening
A44 left temple fringe panel
A44 left temple transition panel
A44 off-center main bang panel
A44 right swept fringe panel
A44 right temple transition panel
A45 left tapered flexible cheek lock
A45 right tapered flexible cheek lock
A42 Left asymmetric rear lock
A42 Off-center main rear lock
A42 Short right rear lock
A42 Main lock left seated seam
A42 Main lock right seated seam
Subtle crown center seam
```

Do not delete, rename, copy geometry from, or derive vertices from those
objects.  Preserve the 41 already-hidden legacy exclusions.  Preserve the
seven exact face witnesses and all 15 frozen bow/root objects listed in
`out/reimu_fumo_attempt_073_profile_loft/asset_interface.md`; record their
object transform, data-block identity, vertex/curve coordinate digest,
material-slot digest, `hide_viewport`, and `hide_render` before construction
and require byte-for-byte field equality before saving.

The protected parent must remain
`sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.
The builder re-hashes it before append and after the candidate is saved.

## Machine-readable inputs

Compile one `p0_contract.json` before Blender starts.  The build script reads
that file; no duplicate coordinate table may live in the build script.

The contract contains:

- `Wh_m = 0.132`, world axes `+X = subject right`, `+Y = rear`, `+Z = up`;
- crown datum `Zc = 0.2202 m` and front center `Xc = 0`;
- all `source_authority` paths and SHA-256 values from
  `head_hair_curves.json`;
- these raw, unrounded observed curves from that file:
  `canonical_front_outer_left`, `canonical_front_outer_right`,
  `canonical_front_hairline`, `canonical_front_left_lock_outer`,
  `canonical_front_right_lock_outer`, and
  `canonical_profile_front_visible`;
- the six independent left/right front scanline rows from
  `front_outer_scanlines.csv`;
- the source hash and schema/version of both files; and
- the 15-object hide allowlist and 22-object frozen-witness allowlist.

Static compilation fails if a source file or controlling image hash differs,
an expected curve is absent, a point count differs, curve normalization has
changed, the source coordinates are not finite, or the two lock curves are
equal under reflection.  The contract labels every unobserved depth value
below `construction_hypothesis`, never `measurement`.

Use source pixels as stored.  Convert canonical-front points with:

```text
x = Wh * (pixel_x - 485) / 368
z = Zc - Wh * (pixel_y - 231) / 368
```

The hairline's published `u` origin is equivalent after subtracting `0.5`
from `u`; do not mix its `301 px` left-bound origin with the signed formula.

## Common deterministic curve operations

Implement the following operations locally, without SciPy or Blender
operators whose output depends on selection or mode.

1. `polyline_arclength(points)` uses double-precision Euclidean cumulative
   length and retains every source point.
2. `resample_polyline(points, fractions)` uses piecewise linear interpolation
   in source order.  It does not smooth, extrapolate, or close an open curve.
3. `smoothstep(t) = t*t*(3 - 2*t)` for `t` clamped to `[0,1]`.
4. All generated coordinates are rounded only for the JSON report, never
   before `Mesh.from_pydata`.
5. Vertex indices are append-order stable.  Faces use explicit index
   arithmetic below; never use `bpy.ops.mesh.*`, triangulation, remesh,
   shrinkwrap, cloth, or mirror.

Smooth shading is presentation state, not geometry.  P0 uses one Catmull-Clark
modifier at viewport/render level 1 on each brown object.  Assign crease
weight `1.0` to source-visible free boundaries and `0.35` to the hidden saddle
root; this makes subdivision smooth the panel interiors without moving the
measured outer edges into a different contour family.  Store and check both
base-cage and evaluated landmark errors.

## Crown/temple saddle

### Domain and source boundaries

The saddle is one rectangular paired-skin grid with `Nu=25` columns and
`Nv=8` courses per skin.  It is a continuous front field and crown return,
not a closed cap.

Form the source-visible outer arc `O` in left-to-right order by concatenating:

1. `canonical_front_outer_left` in reverse order, from its lowest temple
   point to the crown;
2. the crown point once; and
3. `canonical_front_outer_right` from its second point to its lowest temple
   point.

This is exactly 25 source vertices.  Let `s_i` be the normalized cumulative
arc-length fractions of `O`.  Resample the seven-point front hairline at the
same 25 fractions to produce `H_i`.  The resampling is only correspondence;
the seven original H0--H6 points must be inserted as exact additional
fractions if they are not already present.  In that case form the sorted union
of both fraction sets, resample `O` and `H` on that same union, increase `Nu`
deterministically, and report the resulting `Nu`.  The normal case is
therefore 25--30 columns, never an arbitrary fixed six-ring loft.

For each column `i`, define the front-field courses `j=0..4`:

```text
q = j / 4
B(i,j).x = lerp(H_i.x, O_i.x, smoothstep(q))
B(i,j).z = lerp(H_i.z, O_i.z, smoothstep(q))
B(i,j).y = Wh * [
    -0.408
    + 0.020 * (1 - (2*s_i - 1)^2)
    + q * (0.225 + 0.055*cos(pi*(2*s_i - 1)))
    + 0.010 * sin(pi*q) * sin(2*pi*s_i + 0.35)
]
```

The first two terms establish a shallow, slightly stuffed face field rather
than a camera-plane card.  The asymmetric sine is deliberately low amplitude
and fixed; it prevents a mechanically mirrored normal field without moving a
source-controlled X/Z boundary.

Courses `j=5..7` continue over the crown to an open root.  Let
`r=(j-4)/3`, `e=smoothstep(r)`, and `a=2*s_i-1`:

```text
B(i,j).x = O_i.x * [1 - e*(0.10 + 0.07*abs(a))]
B(i,j).y = B(i,4).y + Wh * e * [0.56 - 0.10*abs(a)]
B(i,j).z = O_i.z + Wh * [0.035*sin(pi*r)*(1-0.35*a*a)
                         - e*(0.045 + 0.055*abs(a))]
```

This produces a broad crown turn and a deliberately open, slightly contracted
rear root.  The hidden `j=7` root reaches approximately `+0.36..+0.39 Wh`
while the hairline stays near `-0.41 Wh`, so the evaluated compact field is
eligible for the `.77-.85 Wh` depth gate.  It does not add a rear wall, nape,
or underside.  The depth constants are hypotheses and may receive one bounded
P0 correction only after a registered profile/3Q render; X/Z source boundaries
may not be rescaled to compensate.

### Paired skins and indices

The `B` grid is the outer visible skin.  Create its quads in `(j,i)` order:

```text
outer_index(j,i) = j*Nu + i
outer_face(j,i) = (outer_index(j,i), outer_index(j,i+1),
                   outer_index(j+1,i+1), outer_index(j+1,i))
```

for `j=0..6`, `i=0..Nu-2`.

Compute area-weighted vertex normals from those base faces.  Create an inner
skin in a second contiguous block:

```text
t(i,j) = Wh * [0.012 + 0.004*sin(pi*j/7)^2
               + 0.002*sin(pi*s_i)^2]
I(i,j) = B(i,j) - t(i,j) * normalized_normal(i,j)
inner_index(j,i) = Nv*Nu + j*Nu + i
```

Reverse every inner face winding.  Join outer to inner only along:

- the measured free hairline `j=0` (`Nu-1` quads);
- the left temple return `i=0` (`Nv-1` quads); and
- the right temple return `i=Nu-1` (`Nv-1` quads).

Do **not** join either `j=7` root edge.  Those two root polylines, connected
only at their lateral endpoints by the temple strips, form one explicit open
boundary loop of `2*(Nu-1)` edges.  Expected base counts are:

```text
vertices = 2*Nu*Nv
faces = 2*(Nu-1)*(Nv-1) + (Nu-1) + 2*(Nv-1)
boundary_edges = 2*(Nu-1)
boundary_loops = 1
connected_components = 1
```

No face may have all four vertices on `j=7`; this exact invariant prohibits a
rear closure.  No Solidify modifier is allowed.  The paired skins are authored
geometry, so their thickness and open root can be inspected without modifier
evaluation.

## Independent cheek-lock pockets

Build left and right from their own source polylines.  Do not reflect,
average, share a mesh datablock, instance, or use a Mirror modifier.

Treat each raw lock polyline plus the straight root course from its last point
back to its first point as a closed 2-D boundary.  Select seven `z` stations
from the sorted union of its source-point heights: root maximum, four evenly
spaced interior quantiles, penultimate low station, and exact lowest source
point.  At each station, intersect a horizontal line with the source polygon;
take the leftmost and rightmost intersections.  A station with fewer than two
intersections is a preflight failure, not a reason to extrapolate.

For each station `j=0..6`, build five width samples `i=0..4` using
`u=i/4`, `x=lerp(x_left,x_right,smoothstep(u))`, and the exact station `z`.
Use a distinct fixed depth course for each side:

```text
root_y_left  = -0.360 Wh; tip_y_left  = -0.392 Wh
root_y_right = -0.348 Wh; tip_y_right = -0.405 Wh
c_j = j/6
y_center = lerp(root_y_side, tip_y_side, smoothstep(c_j))
y_front(i,j) = y_center - Wh*[0.008 + 0.010*sin(pi*u)*sin(pi*c_j)]
y_back(i,j)  = y_center + Wh*[0.012 + 0.010*sin(pi*u)*sin(pi*c_j)]
```

This is a shallow stuffed pocket with unequal left/right course and bulge,
not a planar extrusion.  The source polygon independently controls each
front silhouette.  Add 35 front vertices followed by 35 back vertices, create
`4*6` quads per skin with reversed back winding, and join all four boundary
courses with narrow edge strips.  Each lock therefore has exactly:

```text
vertices = 70
faces = 68
boundary_edges = 0
connected_components = 1
all_edges_have_two_faces = true
```

The source-visible lower course is a free stuffed edge: it is closed by its
own narrow strip but has no shared vertex, parent, weld, or bridge to the
saddle/support.  Only the upper root is seated.  After subdivision, require
root-to-saddle clearance `.002-.010 Wh`, free-boundary clearance at least
`.015 Wh` over its central 70%, and zero evaluated triangle crossings.  If a
root misses, adjust only its unobserved Y depth, not its source X/Z outline.

## Subordinate face support

P0 deliberately does not reconstruct another hidden cushion.  Build one open
`9 x 9` beige front patch whose sole purpose is to carry the frozen face
appliques and fill the accepted front opening.

Use bounds `x = +/-0.330 Wh`, `z_top = Zc-0.285 Wh`, and
`z_bottom = Zc-0.920 Wh`.  For `u=i/8`, `v=j/8`:

```text
x = 0.330 Wh * (2u-1)
z = lerp(z_top, z_bottom, v)
y = Wh * [-0.386 - 0.012*(1-(2u-1)^2)*(1-(2v-1)^2)]
```

Create 81 vertices and 64 front-facing quads, reversing the ordinary
left-to-right/top-to-bottom grid winding so every base normal points toward
`-Y`.  Do not add Solidify, a back
skin, side walls, cap, rear surface, or subdivision.  The support must have
one boundary loop of 32 edges and normals toward `-Y`.  Its P0 neutral material
uses a Geometry/Backfacing mix to make rear-facing camera rays transparent.
That is an explicit diagnostic-stage contract, not a proposed final head
construction.  It is necessary because P0 intentionally omits the rear
pockets and must not solve rear leakage by closing the brown saddle.

The support's projected visible pixels are evaluated after brown occlusion,
not from its mesh bounds.  They must form the accepted approximately
`.603 x .603 Wh` beige exposure in the front and contribute zero pixels in
rear, both profiles, and both 3Q support-leak masks.  Frozen face witnesses
remain visible and unchanged; any face/support crossing is reported, not
hidden by moving those witnesses.

## Required modifiers, materials, and object isolation

- Saddle and locks: one Catmull-Clark modifier, levels 1/1, optimal display
  off, UV smoothing fixed explicitly; no other geometry modifier.
- Support: no geometry modifier.
- Every object and mesh datablock is single-user and newly created.
- Assign one new brown neutral-clay material to all three brown objects and
  one new diagnostic beige material to the support.  Do not mutate legacy
  materials.
- Apply no transforms; all candidate matrices remain identity in world space.
- Add no parent, constraint, driver, shape key, armature, vertex group, seam
  object, fiber, stitch, wrinkle, bow change, or rigging artifact at P0.

## Static and clean-reopen mechanical gates

Fail before rendering if any gate below fails.

### Provenance and boundary

1. Parent hash matches before append and after save.
2. Contract input/source/reference hashes match exactly.
3. Exactly the 15 allowlisted objects change from render-visible to hidden;
   their geometry, transforms, datablocks, materials, and names are unchanged.
4. All 22 frozen face/bow witnesses and all 41 pre-hidden exclusions match
   their before snapshots.
5. Exactly four `A74_` P0 objects exist and are render-visible.  There is no
   P1/rear object and no other new render-visible mesh/curve.

### Saddle representation vetoes

1. Counts match the formulas for recorded `Nu`; one connected component, one
   boundary loop, every face a quad, no edge has more than two incident faces.
2. The boundary is exactly the paired `j=7` root courses.  No root-closing
   face, rear wall, bottom, radial fan, pole, or Solidify modifier exists.
3. All seven evaluated hairline anchors are within `.03 Wh`; all six source
   outer scanline extrema are within `.03 Wh`; no measured-boundary segment
   exceeds `.12 Wh`.
4. Evaluated compact depth is `.77-.85 Wh`; maximum thickness is `.020 Wh` or
   less, minimum thickness is `.010 Wh` or more, and no inner/outer crossing
   occurs.
5. Card veto: outer-skin Y range is at least `.70 Wh`, the normal rotates at
   least 55 degrees along the center meridian, lateral and meridional normal
   variance are both nonzero, and no run of three adjacent course rows is
   coplanar within `1e-5 Wh`.
6. Helmet veto: the mesh is non-watertight, its sole boundary is rear-facing,
   and the sum of faces whose four vertices lie on the rear-root course is
   zero.

### Locks and support

1. Each lock has 70 vertices, 68 quad faces, one closed manifold component,
   and no mirror/instance/shared-data relationship.
2. Source-outline evaluated projection error is at most `.03 Wh`; widths are
   separately `.14-.18 Wh`; crown-to-lowest-lock is `1.098 +/- .03 Wh`.
3. RMS distance between the left lock and a reflected/reparameterized right
   lock is at least `.005 Wh`; this catches accidental symmetry without
   requiring decorative asymmetry.
4. Root seating, free-edge clearance, and triangle-crossing gates above pass.
5. Support has 81 vertices, 64 quads, one 32-edge boundary loop, no back/side
   faces, no modifier, and every polygon normal has negative Y.
6. A material-ID mask reports zero support pixels outside the accepted front
   opening in all six fixed views.

### Determinism

Run the builder twice into two separate `out/` candidates.  On clean reopen
under pinned Blender 5.2.1, require identical ordered base-coordinate hashes,
face-index hashes, modifier/property reports, visibility manifests, and
render settings.  `.blend` file bytes need not match because Blender may write
non-geometric metadata; the normalized scene manifest must match.

## First-pixel early return

Only after every mechanical gate passes, render 512--640 px front, left/right
front 3Q, and both profiles, plus the rear support-leak mask in frozen
whole-subject context.  Pair each likeness view with its controlling source
and a registered silhouette overlay.

Return early and reject P0 before any rear pocket, second modeling pass,
material detail, or seam work if any of these occurs:

- outer/hairline/aperture error exceeds the plan tolerance;
- support contributes a non-front pixel;
- the saddle reads as a helmet, mask, visor, slab, card, bridge, rounded box,
  or egg;
- either lock reads as a side wall, mirrored paddle, floating card, or clips;
- there is bald daylight, a root gap over `.01 Wh`, stand-off over `.05 Wh`,
  or accidental tangency;
- either 3Q loses the continuous crown-to-temple seating or correct
  fringe/face/lock order; or
- implementation-blind identity, construction, or contact is below `6/10`.

One bounded correction is allowed only if the representation itself survives
all vetoes.  It may alter the unobserved Y-course hypotheses, support bounds,
or skin thickness.  It may not mirror source curves, rescale X/Z boundaries,
close the saddle, add a receiver/rear wall, bridge a gap, or move a frozen
witness.  A surviving card/helmet or the need for closure resets A74 rather
than authorizing another patch.
