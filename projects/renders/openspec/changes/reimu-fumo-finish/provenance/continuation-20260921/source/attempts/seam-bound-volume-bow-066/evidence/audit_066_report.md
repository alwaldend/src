# Read-only 064-to-066 bow regression audit

No covered non-target regression was found. The four replacement red panels
are closed with positive volume, and the zigzag centerlines meet the requested
front offset. Crown embedding and sampled ruffle seam clearances remain
explicit findings. This report grants no visual, module, or whole-asset
acceptance.

Candidate `macro_066_candidate.blend` SHA256:
`39cf9b6384b375f78b15ab65ac3a266494a20f2683a4722e68c3ddb42746b98e`.
Protected comparison `macro_064_candidate.blend` SHA256:
`502a52cbd825460f20f49d15d9b6443edfc68a8a9883c6d139aead7741db0211`.
Pinned Blender 5.2.1 LTS, build `9e2066aef7ef`; background, factory startup,
automatic scripts disabled, four threads. One general audit and two seating
audits exited 0. No model save or render was performed. Observations are
recorded in the JSON at 2026-09-05 18:08:03, 18:08:25 and 18:10:45 UTC.
Input hashes remained exact during the audits and in fresh post-exit checks.

## Frozen controls and replacement panels

| Check | Bounded result |
| --- | --- |
| Inventory | 405 scene objects in both files; no additions/removals after Macro064/Macro066 normalization. |
| Non-target surfaces | All 31 protected visible mesh/curve controls match, including head/graphics, hair/locks/ties, arms, sleeves and lower body. No unexpected changed object fields. |
| Arm comparison | Indexed vertices and canonical edge/face topology match; mesh storage ordering is normalized while winding, material/smooth and seam/sharp flags remain checked. |
| Transforms and visibility | All object transforms and recorded visibility fields match. |
| Materials | All 48 material property/node graph fingerprints match; none added or removed. |
| Target positions | All 18 intended objects retain corresponding world X/Z exactly. Each closed panel's two vertex blocks correspond to the original single surface. |
| Finite/state checks | All base mesh coordinates finite; no live CLOTH/COLLISION modifiers. Both scenes remain at frame 1 with no linked libraries. |
| Modifier scope | Exactly the four red-panel SOLIDIFY modifiers were removed; no other modifier changes detected. |

All four red panels have zero boundary, nonmanifold, or wire edges and no
remaining modifiers. Independently measured signed world volumes are:

| Panel | Vertices | Faces | Volume (cm³) |
| --- | ---: | ---: | ---: |
| Left upper wing | 3,234 | 3,232 | 6.014851 |
| Right upper wing | 3,234 | 3,232 | 6.008492 |
| Left tail | 6,034 | 11,609 | 1.937869 |
| Right tail | 6,034 | 11,609 | 1.937869 |

## Contacts

Contact tests use evaluated panel surfaces, including 064's old SOLIDIFY
result, against the actual closed head/crown cushion and the buried bow knot.
All four 066 panels have zero detected nonadjacent self-overlap triangle
pairs; no cross-panel overlap pairs were found among the six red-panel pairs.
Both tails have zero detected head crossings or unanimously inside vertices.
Their root samples remain approximately 6.90–7.81 mm from the head; this
does not assess attachment to other objects.

Both wings remain embedded in the crown and intersect the buried knot.
For diagnostic separation only, the wing root band is |X| <= 19.2 mm
(u <= 0.20 using u=(|X|-3 mm)/81 mm). It is not an established physical seam.

| Wing/crown metric, both sides combined | 064 | 066 |
| --- | ---: | ---: |
| Inside vertices in root band | 843 | 872 |
| Inside vertices beyond root band | 251 | 248 |
| Maximum root-band depth (mm) | 9.541767 | 9.522811 |
| Maximum depth beyond band (mm) | 6.700152 | 7.590459 |
| Head-crossing triangle pairs | 576 | 507 |
| Knot-overlap triangle pairs | 1,077 | 858 |

The 066 depth outside this narrow root band reaches its maximum near
X=19.875 mm, Z=175.506 mm. The measurements retain crown-root embedding and
its changed local depths; they are not a collision-free claim or a visual
severity judgment.

## Graphics, trim and ruffle seating

All 210 zigzag controls hit their respective front panels at Y offsets
0.249999575–0.250000507 mm. All 1,610 dense centerline samples hit, spanning
0.237572–0.263747 mm without entering a panel. The specified centerline
seating is supported. With 0.55 mm bevel radius, 3,927 of 9,760 evaluated
graphic-shell samples lie inside the panels, at a maximum nearest depth of
0.403424 mm; none lies beyond the rear face. Partial embedding is consistent
with radius exceeding centerline offset and is not automatically accidental
clipping.

Upper trim has detected panel intersections at 118 left/117 right sampled
seam thickness sections. Ten sections per side remain outside the panel's XZ
footprint throughout nine samples without a section intersection; maximum
section minimum sampled distance is 0.045279 mm. The focused seam check
classifies nearest distances <=1 micrometer as surface contact, correcting
the broad Y-interval check's treatment of sidewall contact.

For each tail ruffle's 384 seam thickness sections, 226 intersect the red
panel and 39 are inside at all nine sampled thickness stations. At 118
sections, all nine stations lie behind the rear panel and the thickness
segment has no panel intersection. Their minimum sampled gap ranges from
0.003360 to 0.266236 mm, mean approximately 0.123093 mm. These are local seam
clearance/overlap findings; they do not prove whole-ruffle detachment or
continuous attachment between samples. Full seating methods, witnesses and
limitations are in `audit_066_graphic_report.md`.

## Coverage and evidence

No pixels, likeness, animation, rig, export, or whole-scene collision checks
were performed. General mesh fingerprints omit arbitrary attributes/UVs;
material checks omit arbitrary custom ID properties, external image bytes
and exhaustive datablock serialization. Ray/BVH contacts are bounded vertex
and triangle measurements; tangencies, coplanar contacts and unsampled
between-vertex depths limit interpretation. Same-panel triangles sharing a
vertex are excluded from self-overlap counts. Seating checks are discrete;
free trim overhang is not treated as detachment. Hair fingerprints were
checked, while settled 064 hair self-contact measurements were not repeated.

Raw evidence is preserved unchanged:

- `audit_066_regression.json`, SHA256
  `dd212fa34a0f0b362532e7263747ec6a1f1e271b0dc03daac74820d1ebfb07a2`.
- `audit_066_graphic_seating.json`, SHA256
  `b6b5c3009a95816e96e4cb9df21179498af7ccdc3dcaa6951583ab9543e9c631`.
- `audit_066_graphic_seams.json`, SHA256
  `4508a8a7bfe385def368d832fa63d96a9e6e818352376dc1b0d0bd11b3786e18`.
- Matching `audit_066_*.py` and `.log` files record the read-only operations.

The requested audit is complete. No further refinement or Blender run is
pending from this audit worker or its seating worker.
