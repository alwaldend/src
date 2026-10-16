# Frozen 029 two-hand audit

Candidate `hand_029_candidate.blend`, SHA-256
`93b2c163bff3e0a6e0d478b7b2312dd8578610d618b3570140cdf76ca400f386`.
Source `head_028_candidate.blend`, SHA-256
`c4ab72a53eb12e64f7f5d2bb216ea1a1734f0bb43cf8e19393f532624aa671b6`.
Both files remained unchanged. Pinned Blender 5.2.1 LTS; no asset save or
render was performed.

## Preserved state and hand geometry

All 84 non-target controls, the inspected rig record and recorded
material-node/light/world/color settings match independently reopened 028.
Both original inserts remain hidden and intrinsically unchanged. There is
no missing or unexpected visible geometry.

Each new hand has 338 vertices, 672 edges and 336 faces, with two incident
faces per edge and Euler characteristic 2. Base and evaluated meshes are
finite and have no zero-length edges, degenerate faces or locally inconsistent
winding. Material slots are valid persistent database IDs, with no missing
face materials. Every vertex has one unit weight on the correct Arm_L or
Arm_R bone, through the intended rig modifier.

The left 59 and right 58 source faces wholly at axial a <= -0.015 m retain
bit-identical evaluated world coordinates, vertex order and material names.
Maximum corresponding proximal vertex movement is zero. Maximum transverse
movement is 3.87 nm left and 4.68 nm right (float precision). Distal axial
extent moves from -6.076753/-6.137869 mm to +5.000003/+4.999998 mm, using
the source probe's actual cuff centers and axes.

### Inherited right-hand winding defect

029 right-hand volume is negative: -5.436137e-6 cubic meters; left is
+5.400756e-6. A separate source-only probe confirms the orientation was
already inward in 028: right volume -2.472699e-6 and front-center normal
Y = +0.990551, versus left volume +2.472700e-6 and front normal
Y = -0.991175. 029 preserved source face order and therefore inherited this
defect. It must not be described as outward-wound merely because the closed
mesh has locally consistent winding. A winding-only correction can preserve
point positions while intentionally changing the winding-specific claim.

## Cuff overlap is not global separation

Both triangle-edge directions were checked against each side's two sleeve
panels, four cuff edge/fold pieces and two sewn joins, plus the red skirt.
The original inserts already crossed some proximal sleeve panels/joins.
The new cuff edge/fold intersections below are new, but all lie within the
measured axial cuff band; no checked crossing extends beyond its distal edge.

| Cuff component | Left crossings, both directions summed | Right crossings |
| --- | ---: | ---: |
| Front folded cuff edge | 19 | 18 |
| Rear folded cuff edge | 18 | 18 |
| Lower pinched cuff fold | 7 | 0 |
| Upper pinched cuff fold | 0 | 0 |

These are seam-band overlaps, not automatic evidence that the distal hand
lobe visibly clips. They are also not all buried proximal contact. For
example, the front cuff intersections near world coordinates
(-49.729586, -15.627852, 37.937541) mm and
(49.731586, -15.314807, 37.807833) mm are on the first hand/cloth surface
from a front ray. Rear-cuff witnesses can be hidden from the front by
approximately 4.8–5.1 mm of hand yet exposed from the distal cuff-axis view.
This bounds the evidence; the audit does not replace the fixed-image review.

### Sampled overlap depth

A winding-independent three-direction ray-parity probe classified cuff-trim
vertices inside each closed hand. All tested classifications agreed across
the three directions. Maximum nearest-hand-boundary distances were:

| Trim component | Left burial, mm | Right burial, mm |
| --- | ---: | ---: |
| Front folded cuff edge | 1.056736 | 1.192275 |
| Rear folded cuff edge | 0.442504 | 0.828219 |
| Lower pinched cuff fold | 1.214331 | 1.510815 |
| Upper pinched cuff fold | No interior sample | No interior sample |

These are sampled trim burial depths, not exact maximum mesh penetration.
The existing trim meshes are open, so reverse “skin inside trim volume” is
undefined and was not asserted. Bidirectional surface-edge intersections
remain valid independently of that volume limitation. No new requirement
to close or rebuild existing trims is inferred.

## Skirt separation and evidence boundary

Neither hand intersects the red skirt in either tested triangle-edge
direction. Minimum sampled bidirectional vertex-to-triangle separations
are 4.239505 mm left and 4.331653 mm right. Hand samples beyond the distal
cuff band have at least 14.543652/14.160210 mm sampled skirt separation.

Evidence:

- Main audit in [machine evidence](hand_029_machine_evidence.md): controls, geometry, materials,
  weights, correspondence, both-direction contact and visibility witnesses.
- Source orientation in the same machine evidence: inherited winding.
- Trim penetration in the same machine evidence: parity method and depth
  witnesses, without inventing a volume for open cloth.

This is a bounded rest-frame audit, not global collision, animation,
reference-likeness or final-stage acceptance.

The separately frozen [029b addendum](hand_029b_technical_audit.md) documents
the right-hand winding correction, exact left evidence carry-forward and
fresh right contact measurements after Blender changed quad triangulation.
