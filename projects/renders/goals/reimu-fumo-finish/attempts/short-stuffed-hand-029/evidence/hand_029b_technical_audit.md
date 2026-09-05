# Frozen 029b winding repair and contact addendum

Candidate `hand_029b_candidate.blend`, SHA-256
`9ad353c57147831cd9440ec8ef7836f95dfb8c719da7f14fe1d122802f16f37d`.
Compared with frozen 029, SHA-256
`93b2c163bff3e0a6e0d478b7b2312dd8578610d618b3570140cdf76ca400f386`.
Pinned Blender 5.2.1 LTS independently reopened both. Both files remain
unchanged; no asset save or render was performed.

## Verified correction and unchanged state

All 84 non-hand controls, inspected rig state and recorded appearance match
029. Both hands retain exactly the same evaluated vertex positions
(maximum displacement zero), per-face material regions, per-vertex weights,
object matrices and armature targets. Persistent material bindings remain
valid. Both closed meshes have 338 vertices, 672 edges and 336 faces, with
no nonfinite coordinates, zero-length edges, degenerate faces or inconsistent
local winding.

All 336 left polygons retain their winding. All 336 right polygons reverse
winding, correcting the inherited inward orientation. Signed volumes are
now positive: left +5.400756e-6 and right +5.440189e-6 cubic meters.

The 59 left and 58 right protected proximal source-indexed faces retain
bit-identical point coordinates and material regions. The right-hand claim
is coordinate preservation, not winding preservation.

## Why right contact evidence was rescanned

Left evaluated triangles are exactly identical when orientation is ignored,
so the previous left contact/penetration evidence carries forward exactly.
Although every right polygon retains the same four corners, reversing its
order changed Blender's quad diagonals: all 672 unoriented evaluated
triangles differ. The small difference in right-hand volume magnitude is
consistent with this changed piecewise-triangular surface. Therefore prior
right contact measurements were not reused as current measurements.

A bounded right-only rescan checked the same sleeve panels, cuff edges/folds,
sewn joins and skirt, plus sampled trim burial. No global scan was added.

## Updated right contact evidence

Every stable cloth-edge intersection identity remains the same as in 029;
the greatest corresponding witness movement is 0.051432 mm at the rear
folded cuff edge. The front cuff has 19 triangle-edge crossing witnesses
instead of 18, reflecting the changed hand triangulation; the rear cuff
still has 18. All cuff-edge intersections remain within the measured cuff
axial band. No checked intersection appears beyond its distal edge.

Some seam-band intersections are first surfaces from the front or distal
cuff-axis view, so this is not an intersection-free assembly or a claim
that every overlap is hidden. The rescan found no new exposed-distal
crossing region. Material assignment is identical; the corrected right
normal orientation intentionally affects shading, and this audit makes no
pixel-equality or visual-acceptance claim.

Sampled trim-vertex burial within the closed right hand is now:

| Right trim | 029 maximum, mm | 029b maximum, mm |
| --- | ---: | ---: |
| Front folded cuff edge | 1.192275 | 1.181324 |
| Rear folded cuff edge | 0.828219 | 0.815146 |
| Lower pinched cuff fold | 1.510815 | 1.490561 |
| Upper pinched cuff fold | No interior sample | No interior sample |

The three winding-independent ray-parity directions agree for every tested
classification. These are sampled burial depths, not exact maximum
interpenetration. Existing trim surfaces are open, so reverse “skin inside
trim volume” remains undefined and is not asserted.

The right hand still has no tested skirt intersection in either edge
direction, with unchanged minimum sampled separation of 4.331653 mm.
The exact-carried left skirt separation remains 4.239505 mm.

Evidence:

- [Selected machine evidence](hand_029_machine_evidence.md) includes the
  final differential audit, contact-evidence resolution and right-rescan
  digest. The rescan selection includes changes and trim probes; its large
  raw contact-witness array remains task-local and is bound by file hash.
- [029 audit and carried left evidence](hand_029_technical_audit.md).

No shape, reference-likeness, global collision, animation or full-stage
acceptance is inferred from these bounded checks.
