# Reimu Fumo attempt 11

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 11 — direct clay sculpt and open cloth relief

**Candidates:** planned `out/reimu_fumo_attempt_011_head_sculpt.blend` and
`out/reimu_fumo_attempt_011_bow_field_sculpt.blend`, SHA-256 values pending,
neutral construction stage, review packet `attempt-011-direct-sculpt`.

**Failure targeted:** Attempt 10 used corrected dimensions but still encoded
the head with front-to-rear correspondence and the bow with one coupled closed
cross-section. Those generators had too few independent artistic degrees of
freedom: moving one region inevitably changed unrelated surfaces, subdivision
averaged away identifying asymmetry and fold hierarchy, and clean topology
could not prevent mattress, shield, or bladder reads.

**Hypothesis:** Actual local Sculpt Mode strokes on one isotropic head clay can
shape crown, face, cheek, underside, and rear independently while the artist
alternates fixed views. Actual local sculpt or art-directed cloth brushes on
one open bow front patch can establish the controlling physical lobe's outline,
normal flow, broad fold valleys, and root gathering before a backing surface
couples them to thickness. Delaying topology and thickness removes the two
remaining mathematical causes of the repeated failures.

**Reference-interpretation correction before implementation:** The first
dormant head draft treated `Wh = 1` as the bare beige cushion width. The frozen
contract defines `Wh` as the outer head-and-hair envelope. A reference-only
front, side, rear, and three-quarter audit therefore rejected the draft before
execution: its `0.96–0.98 Wh` bare cushion would leave almost no construction
space for the padded cap and locks, and its front-first, fixed-pixel stroke
order would recreate a rounded square before the side profile existed. The
corrected proof targets a `0.90–0.95 Wh` bare cushion, treats hair-covered crown
and rear sections as loose containment checks, establishes crown, underside,
and unequal front/rear depth before cheek width, and calibrates every brush
radius as a fraction of the current projected cushion width. The audit itself
is not geometry evidence and cannot pass the gate.

**Plan written before implementation:**

1. Start from factory-empty Blender data and copy no Attempt 3–10 render mesh.
   Load the bundled Blender 5.1 essential mesh-sculpt brush assets locally;
   use no external service or generated asset. Before touching either
   candidate, run a disposable technical probe through the live interactive
   Blender MCP `VIEW_3D` context: voxel-remesh one test mass, enter Sculpt
   Mode, apply recorded `Grab`, `Elastic Grab`, `Scrape/Fill`, and local
   `Smooth` strokes with object-space locations, and prove that one stroke
   group changes only the intended local vertex neighborhood. Calibrate cloth
   brushes independently on one centered open patch. Save the probe only under
   `out/`; background execution is rejected because it crashed on actual
   `bpy.ops.sculpt.brush_stroke` calls.
2. Reject the workflow itself if actual `bpy.ops.sculpt.brush_stroke` calls
   cannot execute reliably or if recorded locality cannot be verified. Do not
   replace them with scripted vertex formulas, profile interpolation, a
   lattice, proportional-edit emulation, or another generator and call the
   result a sculpt. Diagnose the MCP/Blender context first.
3. For the head, use one disposable cube only as initial clay. Subdivide and
   voxel-remesh it isotropically at `0.025–0.035 Wh`; after remesh no accepted
   surface may retain cube-face, front/rear-patch, ring, trace, or bevel
   correspondence. Keep one connected beige object and no visible seam,
   structural hair color, modifier-driven silhouette, or secondary shell.
4. Shape the head only with broad local Sculpt Mode strokes, alternating
   front, side, rear, and three-quarter cameras after each small stroke family.
   Measure large, medium, and small brush radii as `0.27–0.33`, `0.18–0.24`,
   and `0.10–0.15` of the current projected bare-cushion width rather than as
   fixed pixels. Establish the side crown descent, tucked underside, and
   unequal front/rear depth extrema before expanding the lower cheek in front
   view. Author a broad gently padded face and a narrower soft rear support,
   not evidence-free flat planes. Use symmetry only for the first macro
   blockout, disable it for center-plane side strokes, then add
   `0.004–0.012 Wh` stuffing asymmetry. Local smoothing may remove brush
   artifacts but cannot globally average the accepted planes or silhouette.
5. Target overall bare-cushion width `0.90–0.95 Wh`, height
   `1.015–1.045 Wh`, and depth `0.68–0.74 Wh`. The previously listed crown,
   middle, and lower slice ranges are loose containment vetoes because hair
   hides most crown and rear construction; they cannot design or pass the
   silhouette. The three sections must differ in curvature and centroid as
   well as scale; affine copies fail. A non-rendering future hair guide may
   mark the remaining outer envelope and lower-cheek cover depth
   `0.12–0.19 Wh`, tapering below `0.05 Wh` at the crown, but cannot influence
   or conceal this gate.
6. Render neutral clay, MatCap, silhouette, grazing, and measured sections in
   fixed front, side, rear, and three-quarter views immediately after each
   major head stroke group. Hard-reject any front/rear or top/bottom silhouette
   pair parallel for more than `0.20 Wh`; any side roof longer than `0.18 Wh`
   within `0.008 Wh` height; any two-lobe shoulder, notch, or indentation over
   `0.012 Wh`; any egg, balloon, rounded-box, helmet, foam-block, or molded-toy
   read; or a three-quarter highlight that breaks instead of flowing through
   face, cheek, and side.
7. The bare head gate scores only silhouette, sewn-cushion construction, plush
   medium, and presentation; it cannot fairly score Reimu identity. Require
   every applicable category to reach at least 7/10 under an
   implementation-blind review. If it fails, discard the clay and restart from
   the initial remesh with one changed stroke family; do not patch, retopologize,
   add hair, or add surface detail. A passing clay becomes a later retopology
   target; seam topology may not redefine its volume.
8. For the bow, choose the viewer-right lobe in the physical-front reference
   as the controlling asymmetric subject. Trace its outer contour, throat, and
   two or three visible fold centerlines into non-rendering review guides. Do
   not average the two lobes or use the alternate clean-front bow state as the
   primary shape.
9. Build one open, sufficiently subdivided quad patch in the front plane as a
   disposable cloth-relief canvas. Its topology supplies local edit freedom
   only: it is not a bent U, front/back pair, filled outline extrusion, pocket,
   volume union, or final retopology. Pin or mask only the root neighborhood.
   Shape the outer contour and surface normals with actual Sculpt Mode `Grab`,
   `Grab Cloth`, `Pinch Folds Cloth`, and local smooth strokes, alternating
   close front, grazing, side, and three-quarter review. Keep the proven root
   mask active for cloth strokes because the solver propagates across most of
   an unmasked patch.
10. Sculpt two or three unequal broad channels that converge into the root,
    one dominant diagonal fold, and broad quiet fields between them. Keep the
    pinched throat `0.06–0.12 Wh` long and make its irregular gather fan remain
    visible for `0.20–0.30 Wh`. Fold channels must differ in width, depth, and
    direction, change the neutral normal field, and never terminate as isolated
    engraved grooves or equal accordion pleats.
11. Gate the open red relief before any backing, perimeter bridge, Solidify,
    trim, second lobe, tail, knot, or material detail. Target width
    `0.64–0.72 Wh`, height `0.51–0.64 Wh`, and at least `0.38 Wh` of readable
    outer near-vertical or convex run. Require traced contour landmarks and
    inflections within `0.05 Wh`, at least two readable broad fold valleys
    under neutral MatCap and grazing light, and a root that disappears behind
    a simple non-scoring head proxy rather than ending in a blunt cap.
12. Hard-reject any petal, fin, paddle, bladder, card, tire, hollow-rim,
    open-shell, shield, caterpillar, inflated-pocket, glued-cushion,
    molded-toy, primitive-intersection, collar, cuff, or explanation-dependent
    read. Open boundary edges are intentional at this relief gate, but there
    may be no non-manifold or degenerate geometry. Require at least 7/10 for
    silhouette, cloth construction, plush medium, fold hierarchy, and root
    attachment under an implementation-blind review.
13. Only after the red front relief passes may a later gate duplicate and
    locally offset a thin backing, construct a continuous narrow turned edge,
    and bridge the perimeter. Broad-field finished thickness must remain
    `0.03–0.05 Wh`; only root and folded-edge stacks may reach
    `0.10–0.18 Wh`. The passing front surface and fold normals must remain
    unchanged. Retopology and trim remain forbidden until that thickness gate
    also passes.
14. Forbid throughout Attempt 11: scripted spherified-cube displacement,
    profile lofts, longitudinal rings, Coons patches, superellipsoids, section
    sweeps, voxel unions of bow blobs, filled-outline extrusions, U-folds,
    global subdivision or smoothing after folds, free cloth simulation as the
    silhouette designer, and trim, fuzz, color, or lighting used to imply
    construction absent from neutral geometry.
15. Pair every acceptance contact sheet with the Attempt 10 reference boards,
    retain exact hashes and recorded sculpt strokes, and run one independent
    pixels-only review before advancing. Dimensions and topology are veto
    checks, never evidence sufficient to pass.

**Preflight result before candidate geometry:** Background Blender 5.1.1
crashed on actual `bpy.ops.sculpt.brush_stroke`, and the first two interactive
probes used invalid shared viewport assumptions. Probe V3 locked one centered
target, front orthographic view, `view_distance = 3.0`, projected target width,
region-local `mouse`, window-space `mouse_event`, exact Essentials activation,
matched VIEW/sample sizes, and disabled unified-size override. Standard
`Grab` was exactly repeatable: both `300 px` tests moved the same `27/642`
vertices with identical coordinate and index hashes; `150 px` moved `7/642`,
so locality scaled monotonically. Local `Smooth` changed `53/642` vertices on
Grab-deformed clay, while a multi-sample `Scrape/Fill` path changed `8/642`.
Those three brushes are allowed for the head.

`Elastic Grab` moved `642/642` vertices at both tested sizes, `Grab 2D` reached
the full depth of the test sphere, and `Grab Silhouette` moved none. They are
therefore forbidden for candidate head shaping even though the probe's
technical aggregate passed. `Grab Cloth` and `Pinch Folds Cloth` propagated to
`1341/1353` and `970/1353` open-patch vertices. A derived-region mask kept all
selected masked vertices unchanged while `915` unmasked vertices still moved;
cloth brushes are therefore allowed only with that mask discipline and an
immediate bow checkpoint. The probe contact renders were poorly framed and do
not count as visual deformation evidence; the exact coordinate hashes and
metrics prove execution only. No Fumo candidate geometry was created by the
probe.

**Head family-0 work actually performed:** A dormant, family-at-a-time builder
first failed safely on Blender 5.1 API differences: the Eevee enum was
`BLENDER_EEVEE`, and scene-level ray casting ignored the active Sculpt object.
The temporary driver was corrected to use object-local ray casts with a
distance covering the orthographic origin. Those failures landed no stroke and
their setup files were preserved only under `out/`. The accepted driver then
created one `0.92 x 0.70 x 1.00 Wh`, `0.025 Wh` voxel-remeshed cube with
`7,560` vertices and no modifier or prior render mesh. Initial planar Smooth
paths changed no coordinates and were rejected before checkpointing. A closer
three-pass calibration made all 24 exact Essentials `Smooth` calls return
`FINISHED` and move local neighborhoods; the final coordinate hash was
`54a2e7dbec047506c48137a32d7b2c5b0ea1ad9ed56b08263e4c47b4105749cc`.
The mesh remained one closed manifold and its width/depth vetoes passed.

**Head family-0 raw review evidence:** The fixed neutral, silhouette, and
grazing renders were saved as a 4-by-3 contact sheet with SHA-256
`bb44017ef63c422d6d6210ea323ffbc0eb65ec066cf4e13cb00d353849bb4687`.
An independent context-light reviewer saw only the controlling head references,
the candidate pixels, and the neutral-cushion stage. Scores were front `3/10`,
side `2/10`, rear `2/10`, and three-quarter `1/10`. The five ordered defects
were a rigid box/mattress silhouette, a hard vertical three-quarter division,
front lower-center/left/top puckers, absent crown and underside roll in side
and rear, and underexposed grazing evidence. The numeric gate passed, proving
again that dimensions and manifold topology cannot substitute for likeness.

**Head family-0 criterion result and decision:** Reference likeness,
constructed-plush read, and presentation remain **fail**; measured outer bounds
and temporary manifold hygiene pass only as veto checks; reusable structure,
animation, final materials, clean reuse, and delivery remain **unverified**.
The checkpoint is rejected. No family is approved, no rejected coordinate may
be reused, and no side-depth, cheek, face, hair, or bow work may advance from
it.

### Progress and approach audit after Attempt 11 head family 0

- **Improved:** actual interactive Sculpt Mode execution, exact brush
  provenance, checkpoint hashing, and family interlocks are now verified.
  The candidate also remained one connected manifold inside the bare-cushion
  width/depth containment bounds.
- **Regressed or unchanged:** absolute head pixels did not improve. The side
  and rear stayed rectangular, the three-quarter split became explicit, and
  Smooth introduced new face notches and dimples. Presentation regressed in
  grazing views because most of the side/rear signal fell to black.
- **Absolute versus relative:** this was not an absolute visual improvement
  over Attempt 10. It was only a workflow proof attached to another rejected
  block.
- **Evidence for continuing the workflow:** V3 and all 24 recorded calls prove
  that local native strokes are deterministic and auditable. The failure came
  from choosing a smoothing stroke on planar seed regions, not from another
  procedural representation or an unreliable operator.
- **Repeated-defect diagnosis:** a planar cube cannot become a cushion by
  averaging its interior face vertices. That preserves silhouette extrema and
  creates concave face artifacts. More Smooth strokes would strengthen the
  same mattress/notch defects and are forbidden.
- **Highest-leverage unresolved problem:** move the macro silhouette extrema
  themselves into a crown and tucked underside before any planar-face or depth
  refinement.
- **Approach decision:** revise the first stroke family but retain the proven
  direct-sculpt execution, evidence, and interlocks. Restart from the untouched
  remesh; do not patch family 0.

**Reset B plan written before implementation:**

1. Create a new temporary candidate directory and a byte-new untouched
   `0.92 x 0.70 x 1.00 Wh` voxel-remeshed cube using the same verified probe.
   Copy no rejected coordinates or checkpoint datablock.
2. Replace all initial Smooth calls with eight single standard-`Grab` strokes:
   mirrored front and rear upper/lower corner tucks plus separate side front
   and rear upper/lower tucks. Start on the visible silhouette within
   `0.01–0.02 Wh` of each seed corner; move each endpoint inward and vertically
   toward the target crown or underside by `0.05–0.09 Wh`.
3. Use one medium `0.24–0.28` projected-width radius per tuck, no repetition,
   no Elastic/2D/Silhouette brush, and no Smooth. Require each exact native
   stroke to move a nonzero local set without changing topology.
4. Rerender fixed neutral, silhouette, and grazing front, side, rear, and
   three-quarter views. Keep a dim diagnostic fill and world contribution in
   grazing mode so every surface remains readable.
5. Accept this internal family only if it removes the full-height/full-width
   corner runs without a notch, crater, two-lobe split, new plate, or loss of
   the `0.90–0.95 Wh` middle-width gate. An independent pixels-only review must
   score all four views at least `7/10`; otherwise discard the reset and change
   the seed strategy before another sculpt family.

**Reset B work actually performed:** The interactive driver removed
`view_selected`, set Blender's smooth-view duration to zero while changing
views, fixed the region at the audited orthographic location and distance, and
forced one redraw before projection or ray casting. It created a byte-new
`0.92 × 0.70 × 1.00 Wh` sharp voxel-remeshed cube with 7,560 vertices. Eight
single standard-`Grab` strokes then pulled mirrored front, rear, and side
upper/lower corners inward. No stroke repeated and no Smooth or disallowed
brush ran. The driver rendered and stopped at family 0.

**Reset B raw evidence:** Projected width remained exactly `668.278 px` in
front/rear and `508.472 px` at the side for every relevant stroke. Each stroke
moved only 235–243 of 7,560 vertices, with maximum displacement
`0.0338–0.0431 Wh` and strong-influence radius about `0.178 Wh`. The evaluated
mesh remained one closed manifold at `0.920 × 0.700 × 1.000 Wh`. The exact
checkpoint and contact-sheet hashes are
`4e1c12e1a9bbc1516e58eb9cc29d47432a4c642af6477b7b23a9cc6dccfe9eec`
and
`6912501cd4c176ba1e4dd2a4bf05ce0e99ea9cd9ac01c04d1cf8b76c459837db`.

The context-light reviewer saw only the controlling references, exact
candidate pixels, and neutral-cushion stage. It scored front `4/10`, side
`2/10`, rear `3/10`, three-quarter `2/10`, silhouette `3/10`, sewn-cushion
construction `2/10`, plush-medium read `3/10`, and presentation `7/10`. It
reported long planar faces, nearly square corners, a major full-height
three-quarter cleft, a rectangular side, tight corner pinches, and absent
crown/underside stuffing flow.

**Reset B criterion result and decision:** Native-stroke reliability,
projected-view stability, numeric containment, and temporary manifold hygiene
pass. Reference likeness, constructed-plush read, and presentation readiness
fail; reuse, animation, materials, and final delivery remain unverified. Reject
the exact checkpoint and every coordinate it contains. No later family may
advance from it.

### Progress and approach audit after Attempt 11 reset B

- **Improved:** the corrected viewport eliminated delayed projection drift;
  every Grab stroke was local, reproducible in scale, and auditable. The
  grazing packet was readable enough for an independent absolute rejection.
- **Regressed or unchanged:** the head remained a padded rectangular block.
  The side stayed at `2/10`, and a new full-height two-lobe cleft became the
  dominant three-quarter defect. No user-visible acceptance criterion passed.
- **Absolute result:** this candidate is not closer to approval in absolute
  terms. It proves only that reliable corner tucks cannot repair the seed.
- **Evidence against continuing:** both local Smooth and broad Grab families
  preserved the sharp seed's planar mass. Their different artifacts share one
  cause: the cube starts at final depth with six persistent planes, so local
  subtraction leaves long walls and overlapping view edits form axial
  trenches.
- **Highest-leverage unresolved problem:** begin with a continuous rolled
  perimeter and deliberately insufficient depth, then sculpt front and rear
  fullness outward at different heights before touching the front silhouette.
- **Approach decision:** discard the sharp-cube seed and its complete family
  order. Keep only the proven interactive context and ordinary Grab operator.
  Start Attempt 12 from an under-depth pillow puck; do not use another cube or
  generic ellipsoid.
