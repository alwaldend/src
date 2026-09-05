# Frozen 027c bounded technical audit

Candidate `head_027c_candidate.blend`, SHA-256
`19ce0cb14c7d679750422702d6df97753480cf8a4db7cd73f1203b0f28bf7416`.
Source 026 SHA-256
`56efb16739c746153c5a562195b221645865e0ae4a6c78a5f491783b2c700882`.
Both files remained unchanged. Pinned Blender 5.2.1 LTS cleanly reopened
each; no asset save or render occurred.

## Preserved controls and bindings

All 71 non-target controls match independently reopened 026, the preflight
baseline, and the writer receipt. The inspected rig record and recorded
material-node/light/world/color settings match their independent source.
There are no unexpected visible meshes or missing expected visible meshes.

All 15 new meshes use non-null, non-evaluated persistent database materials;
no face refers to a missing material. Every base vertex has exactly one
unit Head weight. These are rest-frame and recorded-state checks, not
animated-rig acceptance or exhaustive serialization of all Blender settings.

## Geometry and contact

All new base and evaluated meshes have finite coordinates, no edges shorter
than 1e-10 m, and no faces smaller than 1e-14 square meters. Every evaluated
mesh has two incident faces per edge, consistent winding and positive
signed volume. The hood's 382 and fringe's 680 base boundary edges close
under Solidify as expected; open base cloth is not itself a failure.

All 628 source core polygons wholly below world Z = 0.089 m are preserved
exactly after reopening: coordinate hash, vertex order/winding and resolved
material names match. Maximum corresponding world-coordinate displacement
is zero. Coordinate hash:
`7d5fbbdacad948065129cee279cbb3b1aea51b7a023446960c56a7e371aa158b`.

Strict root samples include the original evaluated boundary vertices and
quarter samples along their boundary edges. All 639 sample locations and
nearest-head gaps are exactly unchanged from the 026 preflight:

| Root | Samples | Minimum gap, mm | Maximum gap, mm |
| --- | ---: | ---: | ---: |
| Left collar | 121 | 0.374947 | 0.490014 |
| Right collar | 121 | 0.374945 | 0.489154 |
| Tie | 397 | 0.297082 | 0.810112 |

No collar/tie triangulated-edge crossings into the new core were found.
Of the new hair-cloth/collar pairs, only the right cheek-lock/right-collar
pair had an AABB separation lower bound within 2 mm. Both edge directions
had no intersections; the minimum vertex-to-triangle distance was
2.136873 mm. Other pairs were separated by an AABB lower bound above
2 mm. These are bounded contact checks, not global collision guarantees.

## Hood layers: distinguish outer surface from cloth thickness

The evaluated hood has 18,228 vertices. Its first 9,114 vertices exactly
match a disposable evaluation without Solidify, identifying the outer
surface independently of index assumptions. Their nearest-normal signed
core distances range from zero to 0.550003 mm. There are 101 samples within
0.01 mm of the core and no outer samples with negative distance beyond
0.01 mm. No tested upper-front outer sample was behind the core.

The inward Solidify thickness is 0.7 mm, offset -1. All 9,114 identified
inner-layer samples lie inside the core by 0.042135 to 0.701589 mm. This
is cloth-thickness overlap and is reported separately from outer-surface
clipping; the assembly must not be described as intersection-free.

## Fixed face-plane comparison

The exact 026 9 by 7 core-only frontal ray grid hit 63 of 63 samples in 027c.

| Witness | 026 | 027c |
| --- | ---: | ---: |
| Y residual RMS, mm | 2.122327 | 0.945571 |
| Maximum absolute Y residual, mm | 5.147697 | 4.134376 |
| Orthogonal residual RMS, mm | 2.005124 | 0.885176 |
| Maximum absolute orthogonal residual, mm | 4.863421 | 3.870307 |
| Triangle-normal RMS angle to fitted plane, degrees | 15.968907 | 10.185668 |
| Maximum pairwise triangle-normal angle, degrees | 51.973939 | 50.964433 |

These are same-region geometric comparison witnesses, not reference normals,
differential curvature, new acceptance thresholds or a visual-quality pass.

[Durable JSON bundle](head_027c_machine_evidence.md) contains all object checks,
contact witnesses, layer extremes and fixed-grid samples. No visual,
reference-likeness, full-stage or final acceptance is claimed.
