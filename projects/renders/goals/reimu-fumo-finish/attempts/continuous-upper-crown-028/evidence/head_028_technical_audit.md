# Frozen 028 differential technical audit

Candidate: `head_028_candidate.blend`, SHA-256
`c4ab72a53eb12e64f7f5d2bb216ea1a1734f0bb43cf8e19393f532624aa671b6`.
Differential reference: `head_027c_candidate.blend`, SHA-256
`19ce0cb14c7d679750422702d6df97753480cf8a4db7cd73f1203b0f28bf7416`.
Construction source: frozen 026, SHA-256
`56efb16739c746153c5a562195b221645865e0ae4a6c78a5f491783b2c700882`.
All three files remained unchanged. Pinned Blender 5.2.1 LTS independently
reopened each. No asset save or render was performed.

## Protected state and integrity

All 71 non-head controls, inspected rig state and recorded appearance
match independently reopened 027c and 026. All 14 retired original targets
remain recoverable, hidden and intrinsically unchanged. No expected visible
geometry is missing and no unexpected visible geometry was introduced.

All 15 new meshes have valid persistent database material slots and unit
Head weights. Their base and evaluated geometry have finite coordinates,
no edges below 1e-10 m and no faces below 1e-14 square meters. Evaluated
meshes have two incident faces per edge, consistent winding and positive
signed volume. The open hood/fringe bases close under Solidify as expected.
All corresponding evaluated mesh vertex counts and face-index order match
027c.

## What moved

Source front/rear face membership was independently recovered by disposable
tagged evaluation of 026, rather than inferred from the new shape.

| Core region or control | Differential result |
| --- | --- |
| All 3,837 core vertices at world Z <= 125 mm | Exactly unchanged |
| 4,896 rear-panel faces | Coordinate hash and winding identical |
| 1,248 appended gusset faces | Coordinate hash and winding identical |
| Every core vertex X/Z | Exactly unchanged |
| 628 underside faces wholly below 89 mm | Coordinates, winding and material names exactly unchanged |
| Mouth and all three separate rear cloth meshes | Exactly unchanged |

The core's Y displacement ranges from -0.671878 to +18.085353 mm. Its
world bounds are X [-58.169626, 58.169626] mm, Y [-53.164937, 31.984441] mm,
Z [78.579240, 193.220809] mm. Overall core depth therefore changes from
84.484439 to 85.149378 mm; upper-front recession is not an equal reduction
in total depth.

Literal preservation of every derived cloth/detail point below 125 mm does
not hold. Recomputed surface normals, Solidify thickness and receiver
triangles propagate smaller changes across that cutoff:

| Derived region at reference world Z <= 125 mm | Maximum movement, mm |
| --- | ---: |
| Fringe midsurface without Solidify | 0 |
| Fringe including inward Solidify layer | 0.148487 |
| Hood outer surface and full shell | 0.018900 |
| Left eye applique | 0.001043 |
| Right eye applique | 0.005752 |
| Left cheek lock | 0 |
| Right cheek lock | 0.000075 |

The fringe midsurface preserves X/Z everywhere. Derived normal-offset
geometry does not: maximum whole-mesh X/Z displacement is 0.321269 mm for
the hood and 1.043882 mm for the fringe's inner thickness layer. These
facts narrow the preservation claim; they are not new acceptance thresholds.

## Lower-face grid and contact

All 63 fixed core-only frontal rays hit. The first 54 hit coordinates are
identical to 027c. Only the nine rays at Z = 125 mm change, by at most
0.022218 mm in Y, where evaluated triangles span the cutoff. Maximum
corresponding triangle-normal change is 1.326332 degrees.

Plane Y-residual RMS is 0.945594 mm versus 0.945571 mm in 027c; maximum
absolute residual is 4.136358 versus 4.134376 mm. Orthogonal RMS is
0.885255 versus 0.885176 mm. These remain fixed-region witnesses, not
reference normals, differential curvature or a shape-quality passing score.

All 639 collar/tie root samples and nearest-head gaps are exactly unchanged:
0.297082 to 0.810112 mm overall. No collar/tie triangle-edge crossings into
the core were found. The only new-cloth/collar pair within a 2 mm AABB
lower bound remains the right cheek/right collar: neither edge direction
intersects, and minimum sampled separation is 2.136873 mm.

Hood outer samples remain at signed nearest-core distances of zero to
0.550002 mm, with 101 samples within 0.01 mm and no negative outer samples
beyond 0.01 mm. The inward layer still overlaps the core by 0.042135 to
0.701589 mm. This is reported thickness overlap, not an intersection-free
assembly claim.

## Upper depth witnesses

Actual core-only centerline rays, not analytic helper predictions:

| World Z, mm | 027c front Y, mm | 028 front Y, mm | 028 front-to-rear depth, mm |
| --- | ---: | ---: | ---: |
| 125 | -51.145211 | -51.141024 | 81.970602 |
| 140 | -52.500010 | -53.005278 | 84.825307 |
| 150 | -52.500010 | -50.806701 | 80.865264 |
| 160 | -52.500010 | -46.422675 | 72.681054 |
| 170 | -52.500010 | -40.293545 | 60.336843 |
| 180 | -49.923956 | -32.211810 | 42.760730 |
| 190 | -25.624037 | -22.955447 | 17.103970 |

Rear hits are unchanged in all seven sections. The 180 mm front recession
is 17.712146 mm. The implementation regenerated explicit front/rear meshes;
it did not apply the alternative through-depth volumetric deformation, so
no volumetric-field Jacobian is claimed. Positive sampled depths and
unchanged rear/gusset correspondence are reported directly.

[Durable JSON bundle](head_028_machine_evidence.md) contains full correspondence,
topology, binding, source-preservation, contact and sample records. No visual,
reference-likeness, full-stage or animated-rig acceptance is implied.
