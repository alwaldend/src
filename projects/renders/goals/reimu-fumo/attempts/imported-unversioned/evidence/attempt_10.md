# Reimu Fumo attempt 10

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 10 — sculpted stuffing and occupied fabric volumes

**Candidates:** planned `out/reimu_fumo_attempt_010_head.blend` and
`out/reimu_fumo_attempt_010_right_loop.blend`, SHA-256 values pending, neutral
construction stage, review packet `attempt-010-sculpted-volumes`.

**Failure targeted:** Attempt 9 improved the front image but still placed the
front seam too near the visible front extremum. The beige surface therefore
became a side-view sliver while one deep brown block owned the head depth.
Copied thin hair surfaces produced a separated helmet plate and slab locks.
The second continuous parametric U-fold bow again averaged its intended folds
into a smooth upright petal, so bent-sheet bow construction is exhausted.

**Hypothesis:** Sculpting the occupied stuffing volume first will establish a
shallow face lobe, variable gusset, and separately readable rear lobe before
any render topology constrains them. A shared structural shell fitted to that
accepted volume can give the beige front `0.12–0.18 Wh` of visible side roll.
Closed, locally padded hair panels with deeply buried roots can then overlap
like constructed fabric rather than form a hood. Separately, sculpting the
space occupied by a compressed bow lobe before wrapping welded fabric panels
around it can preserve broad folds that a deformed sheet erased.

**Plan written before implementation:**

1. Start both proofs from factory-empty Blender data and copy no Attempt 3–9
   render mesh. Retain only the frozen traces, landmarks, cameras, numeric
   constraints, and the evidence that a full beige structural front reads
   better than a mask. The temporary sculpt scaffolds must never render or be
   duplicated into a visible surface.
2. Sculpt one hidden watertight head scaffold from a low-resolution spherified
   cube, using about `0.03 Wh` voxel resolution and only broad grab, flatten,
   and compression operations. Target width `0.96–0.98 Wh`, height `1.03 Wh`,
   and depth `0.70–0.74 Wh`; keep front and rear near `y = -0.36/+0.36 Wh`.
3. Give the scaffold three independently readable depth contributions: beige
   front lobe `0.12–0.18 Wh`, variable gusset `0.30–0.36 Wh`, and shallow rear
   lobe `0.17–0.22 Wh`. At the middle section, target the front plane near
   `-0.36 Wh`, front seam near `-0.20 Wh`, rear seam near `+0.14 Wh`, and rear
   plane near `+0.35 Wh`. Keep central face and rear-plane variation below
   `0.020 Wh`, crown compression `0.015–0.025 Wh`, lower settling about
   `0.020 Wh`, and asymmetry below `0.012 Wh`; do not globally smooth after
   establishing the planes.
4. Render scaffold silhouettes and measured cross-sections near
   `z = 1.55, 1.25, 0.85 Wh`. Reject before retopology if any depth
   contribution is absent or if the silhouette reads as a box, flat roof,
   egg, balloon, or uninterrupted rear wall.
5. Retopologize the accepted scaffold as one manifold structural shell. Use
   matching `13 × 13` Coons-style beige-front and brown-rear quad patches,
   48-vertex perimeters, five intermediate 48-vertex gusset rings at depth
   parameters about `0, .08, .27, .50, .73, .92, 1`, and at most one
   subdivision level. The front perimeter must sit behind the visible cheek
   extrema so its outer loops create the required side roll; the rear patch
   must own a shallow plane distinct from the gusset.
6. Form seams with shared topology, never tubes: valley depth
   `0.005–0.008 Wh` and recovery width `0.018–0.025 Wh`. Turn the lower front
   seam under the chin, expose the rear seam in grazing light, and reserve its
   lower portion for buried nape roots. Gate the shell in material-ID, neutral
   clay, side, rear, three-quarter, and grazing views. Require beige side
   ownership of at least `0.12 Wh`, no straight silhouette run over
   `0.25 Wh`, separately readable gusset and rear plane, and at least 7/10 for
   silhouette and sewn construction with no mask, shell, chin rim, or rear
   block.
7. After that pass only, build a closed crown yoke over the front
   `0.12–0.16 Wh` of the crown, width `0.84–0.88 Wh`, padded thickness
   `0.018–0.028 Wh`, and structural-seam overlap `0.06–0.08 Wh`. It may not
   wrap around the rear or surround the face.
8. Build three overlapping closed fringe panels: left field, asymmetric broad
   swept center lock, and right field. Give each distinct front and return
   grids joined through three or four turned-edge rows. Use approximate
   visible grids `9 × 7` center and `7 × 6` sides, lateral overlap
   `0.025–0.045 Wh`, yoke-root overlap `0.06–0.08 Wh`, compressed root
   thickness `0.008–0.012 Wh`, padded field thickness `0.025–0.038 Wh`, and
   free turned-edge thickness `0.014–0.020 Wh`. Preserve the broad measured
   center tip without a symmetric three-tooth outline or uniform inflated rim.
9. Render eye-off construction evidence plus a separately scored eye-on
   proportional packet. Temporary flat diagnostic eyes may use centers
   `x = ±0.18 Wh`, width `0.21 Wh`, height `0.32 Wh`, and projection
   `0.002–0.005 Wh`; they cannot affect construction scores. Require a face
   opening `0.78 ± 0.04 Wh` wide and `0.60 ± 0.04 Wh` high, with no beige
   crown crescent, crack, exposed root, or helmet opening.
10. Build closed front/return cheek-lock panels about `11 × 7` per side.
    Expand from `0.08–0.10 Wh` buried roots to `0.15–0.17 Wh` curved mids and
    `0.10–0.12 Wh` blunt tips over `0.55–0.65 Wh` length. Bury each root
    `0.06–0.10 Wh`, keep mid thickness `0.028–0.040 Wh`, curve outward and
    backward around the cheek, and put the long seam on the outer rear edge.
    Add five separate closed nape tabs whose roots overlap behind the rear
    lower seam by at least `0.05 Wh`; only their asymmetric lower portions may
    break silhouette.
11. Gate the completed head construction at outer width `0.97–1.03 Wh` and
    outer depth `0.74–0.82 Wh`. Require all evaluated meshes to have zero
    boundary, non-manifold, and degenerate geometry. Intended penetration is
    allowed only in recorded root-overlap bands. Require every applicable
    independent visual score to reach 8/10 with no helmet, mask, shell, card,
    slab, crack, floating root, or tangent root.
12. Build the bow as a temporary watertight occupied-volume scaffold, not a
    bent sheet. Use about 18 longitudinal rings by 12 vertices to form an
    asymmetric rounded wedge: width `0.50–0.56 Wh`, height `0.49–0.56 Wh`,
    depth `0.11–0.16 Wh`, and compressed root length `0.10–0.13 Wh`. Establish
    a broad shoulder, at least `0.28 Wh` of readable outer vertical edge, three
    large fabric fields, two broad `0.035–0.055 Wh` valleys
    `0.018–0.026 Wh` deep, and one weaker diagonal compression fold.
13. Reject the isolated untrimmed scaffold unless it scores at least 7/10 and
    reads as a broad compressed ribbon lobe rather than a petal, card, fin,
    wing, paddle, pill, or bladder. Only after that gate, retopologize semantic
    front, turned-edge, and back surfaces around the scaffold. Target an
    approximately `22 × 14` front, four to six turned-edge rows through
    `120–170` degrees, and an approximately `18 × 11` inset back, then weld
    matching boundaries into one manifold before one subdivision level.
14. Author five radial pleat wedges directly inside the root zone, with
    valleys `0.012–0.020 Wh`, ridges `0.008–0.014 Wh`, fan angle
    `75–95` degrees, root compression `35–45%`, and root thickness below
    `0.07 Wh`. Use no collar, cuff, bridge, wrap, rectangle, pocket, or
    separate return surface. Place the passing lobe with outward yaw
    `12–18` degrees, backward pitch `8–14` degrees, and forward tip roll
    `5–10` degrees, with its root behind the future hair.
15. Require the complete untrimmed two-loop bow and simple head proxy to score
    at least 8/10 before trim. Then add a separate genuinely scalloped strip
    with 20–24 irregular gathers, rounded lobes, width `0.030–0.042 Wh`, and
    thickness `0.006–0.010 Wh`, seating its inner boundary on the welded seam.
    Reject any piping, teeth, sine-wave zigzag, spikes, slit, or visible gap.
16. Produce fixed front, side, rear, three-quarter, grazing, material-ID,
    silhouette, overlay, and contact-sheet evidence at every earliest gate.
    Stop each component immediately on a failed gate, record the exact visible
    defect and strategy audit, and copy no render mesh into the tracked asset
    until every applicable score is at least 8/10.

**Pre-implementation reference audit correction:** A pixels-first audit of the
physical front and oblique-side images, turn GIF frames 2–3, and sofa GIF
frames 20, 30, 60, and 70 found two material conflicts in the plan above
before either builder produced an accepted render. The following corrections
supersede affected steps 2–15; the original text remains as an audit trail.

- The references do not show a beige structural front joined to a deep brown
  structural gusset and rear. They show one shallow stuffed cushion whose side
  and rear are covered by thin constructed brown hair. Build the structural
  shell as one unified cushion with total depth `0.66–0.82 Wh`; do not assign
  its side or rear to visible brown hair material.
- The forward beige cheek may own `0.12–0.19 Wh` of side depth only in the
  lower-cheek region beneath hair. It must taper below `0.05 Wh` at the crown,
  not remain a visible beige side band or separate front slab. Gate the
  unified cushion for a soft shallow side/rear silhouette, not separately
  visible structural gusset and rear planes.
- The visible side and rear hair must be closed overlays about
  `0.02–0.05 Wh` thick, with shallow local bulges up to `0.03–0.07 Wh`. Hide
  lock and cap roots beneath at least `0.13–0.21 Wh` of overlap. The physical
  front cap/fringe reads substantially continuous, with only a tiny crown
  dart; distinct local panels may define topology but their visible overlaps
  may not fragment the broad cap. Shape the rear hair as one softly flattened
  pear/trapezoid panel, widest near its middle and tapering into four or five
  rounded lower points.
- The default physical bow lobes are larger than the initial scaffold target.
  Use asymmetric core widths `0.64–0.72 Wh`, heights `0.51–0.64 Wh`, and at
  least `0.38 Wh` of readable outer near-vertical or convex run. A temporary
  scaffold may occupy more depth, but the retopologized broad field must return
  to `0.03–0.05 Wh` core thickness and `0.05–0.10 Wh` front/back separation;
  only local fold and root stacks may reach `0.10–0.18 Wh`.
- Confine the pinched throat to `0.06–0.12 Wh`, then continue a visible gather
  fan across `0.20–0.30 Wh`. Replace five equal radial wedges with two or
  three irregular broad root channels plus at most one dominant diagonal
  fold. Keep each lower hanging tail a separate thin trapezoidal or rhombic
  panel behind the hair; it may not be fused into an upper lobe.
- When trim becomes eligible, use irregular scallop protrusion
  `0.03–0.07 Wh` and pitch `0.04–0.08 Wh`, seated on the edge. These bounds
  supersede any fixed gather count that would make regular piping or teeth.
- Re-run the pixels-only audit at the unified-cushion scaffold and resized
  bow-scaffold gates. If either image still depends on implementation notes to
  explain its construction, reject it before retopology or hair.

**Work performed:** A pixels-first audit of all supplied references ran before
either proof rendered. It caught that the initial plan still treated the brown
hair as structural head depth and undersized the default bow lobe. The plan was
corrected and delivered before geometry execution. Both proofs then started
from factory-empty Blender data and reused no rejected render mesh.

The head proof built one watertight unified cushion scaffold, rendered fixed
front, side, rear, and three-quarter views plus three measured cross-sections,
and stopped at the first pixel gate. It did not build retopology or hair. The
bow proof built one reference-sized watertight swept occupied volume, rendered
fixed and close front, side, three-quarter, grazing, material-ID, and
silhouette evidence, and also stopped at its first pixel gate. It did not
build fabric retopology, the second lobe, tails, or trim. Neither proof entered
the tracked asset.

**Evidence:** The head measures `0.9700 × 0.7100 × 1.0299 Wh`. Its lower,
middle, and crown front rolls are `0.1121`, `0.1158`, and `0.0411 Wh`, with a
`0.165 Wh` peak at `z = 1.05475 Wh`. The evaluated mesh has 1,123 vertices and
1,120 faces with zero boundary, non-manifold, or degenerate geometry. Its
internal pixel gate scored silhouette 4/10, sewn construction 3/10, intended
medium 4/10, and presentation 7/10. The implementation-blind review scored
silhouette 5/10, sewn construction 2/10, intended medium 3/10, reference
plausibility 3/10, and presentation 6/10. Both saw a mattress side, flat roof,
near-vertical walls, uninterrupted rear, and a rigid two-lobe three-quarter
transition.

The bow measures `0.6900 × 0.222463 × 0.599746 Wh`, with a `0.090 Wh`
throat, `0.250 Wh` gather-fan allocation, `0.390 Wh` outer run, broad-field
local separations `0.092–0.102 Wh`, and root/outer local stacks
`0.150/0.130 Wh`. Its 288-vertex mesh is watertight with zero boundary,
non-manifold, or degenerate geometry. The internal visual gate scored 3/10.
The implementation-blind review scored silhouette 2/10, sewn construction
1/10, intended medium 2/10, reference plausibility 1/10, and presentation
5/10. It reported balloon, bladder, pill, paddle, rounded-box, hollow-tube,
and thick-card reads: the front was one smooth shield, the valleys and gather
fan disappeared, and side/three-quarter views became a deep hollow rim.

**Criterion results:** Reference auditing and corrected numeric bounds pass.
Both earliest meshes are closed and technically clean. All applicable visual
criteria fail by large margins, so complete likeness, sewn construction,
presentation, reuse, animation readiness, and final integrity remain
unverified. Repository delivery still applies only to the rejected migrated
baseline.

**Decision:** Reject both Attempt 10 scaffolds. Copy no geometry, modifier,
surface, or topology from either proof. Retain only the corrected reference
bounds, fixed review boards, camera contract, and evidence that the first
topology-independent gates still accidentally coupled their macro forms to
longitudinal generation.

### Progress and approach audit after attempt 10

- **Improved procedurally:** the separate reference audit caught two wrong
  assumptions before rendering. Both builders stopped at their first failed
  pixels, avoiding hair, bow trim, body work, and tracked binary churn. The
  default bow's corrected width, height, outer run, throat, and localized depth
  ranges now match the measured physical reference.
- **Improved technically:** both neutral proof meshes were watertight and clean,
  and their metrics explicitly separate local stack depth from overall turned
  centerline depth.
- **Regressed or unchanged visually:** head silhouette remained 4–5/10 and bow
  silhouette fell to 2/10. No proof is remotely close to the 7/10 earliest
  component gate or 8/10 acceptance gate.
- **Absolute result:** neither proof is a viable modeling base. The corrected
  semantic description did not prevent the builders from encoding the same
  unwanted coupling in their base topology.
- **Head representation failure:** the object was called a sculpt scaffold but
  remained a front patch, longitudinal side rings, and rear patch. That
  construction mathematically preserved long front/rear walls, a flat roof,
  and a hard front-to-rear shoulder. It did not permit crown, cheek, underside,
  and rear masses to move independently.
- **Bow representation failure:** one swept closed cross-section coupled front
  fields, back, outer turn, root, and both caps. Sparse angular samples erased
  the valleys, while the smoothed end cap created a false cavity. Refining the
  sweep cannot produce independent folded cloth fields.
- **Highest-leverage problem:** remove longitudinal correspondence from both
  macro blockouts. The next gate must judge genuinely freeform occupied
  volumes before any semantic panel topology is derived.
- **Next approach:** sculpt one head from a low-resolution spherified cube or
  voxel-remeshed rounded mass, using independent broad crown, cheek, lower,
  and rear edits until a pear/trapezoid side passes. Build the bow from three
  separately readable front fabric masses plus a saddle and compressed throat,
  then voxel-union/remesh them into one temporary scaffold. Retopology remains
  forbidden until those neutral pixels pass.
