# Frozen leg031 technical record — visually rejected

The minimal technical checks reproduce the writer's geometry and contact
measurements. They do not reverse the root and blind-review rejection:
cream sections remain unreadable in the five review views. Geometric reveal
bounds and small hem distances are not visible-exposure certification.
Root retains030b; no031 repair was performed.

Candidate031 SHA256:
`ece85247dc07e9ac59388c20321b992c8638d4f3294ac4d4ef6436e975489b71`.
Independent source030b SHA256:
`d69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6`.
Both clean-reopened in pinned Blender 5.2.1 LTS, build `9e2066aef7ef`,
background with automatic file scripts disabled. Both hashes remain intact;
no model save, render or canonical change was made.

## Preservation and mesh checks

- All 84 visible geometry controls match independently reopened030b and the
  writer receipt. All 214 non-target object records also match independently.
  Recorded rig/rest pose and appearance are unchanged. No missing or
  unexpected visible geometry was found.
- Both retired source roots retain exact intrinsic and evaluated geometry;
  they are hidden from render and the view layer, not globally disabled.
- Each replacement preserves exactly 737 source vertices in the distal
  `q <= 0.42` band, with zero world displacement. Its 680 wholly protected
  faces retain exact coordinates, indices and winding.
- Both roots have 1,962 vertices, 3,976 edges, 2,016 faces and Euler 2.
  Every edge has two incident faces, with consistent winding, no nonfinite
  vertices, no edges below 1e-10 m and no faces below 1e-14 m².
  Evaluated signed volumes are positive: left +4.676603927e-6 m³,
  right +4.676604442e-6 m³.
- Material regions match source; `Dress warm white cloth.002` is a valid
  persistent database ID. Every vertex has exactly unit `Leg_L` or `Leg_R`
  weight, respectively, attached to `ReimuFumoRig`.

## Requested contact confirmation

| Measurement | Left | Right |
| --- | ---: | ---: |
| Source root/pod BVH polygon-pair overlaps | 148 | 150 |
| Candidate root/pod BVH polygon-pair overlaps | 154 | 154 |
| Root/Hem026 polygon-pair overlaps | 0 | 0 |
| Root/Skirt022 polygon-pair overlaps | 0 | 0 |
| Bidirectional sampled root/hem distance | 0.494444 mm | 0.333350 mm |

All candidate contact counts and nearest samples match the writer receipt
exactly. The nearest witnesses are hem vertices 16051 and16215. These are
polygon-pair intersection and vertex-to-surface measurements, not penetration
depth, exact minimum-distance bounds or seam-visibility guarantees.
Black pods and Hem026 are exact controls, so their existing geometric
relationship carries without another pod/hem scan.

The bounds extend inward by 8.986213/8.782196 mm and above each pod by
5.872937 mm, yet this does not establish visible cream area. No additional
occlusion, global collision or animation scan was added after visual rejection.

Machine evidence: [complete compact report](leg_031_machine_evidence.md), SHA256
`f4125d9ce38654d7498a08634577cf6adfb498bf57b562ca503ad45ace2d4695`.
It contains the complete compact report, protected-coordinate hashes and
contact measurements. No image-based retention is inferred from it.
