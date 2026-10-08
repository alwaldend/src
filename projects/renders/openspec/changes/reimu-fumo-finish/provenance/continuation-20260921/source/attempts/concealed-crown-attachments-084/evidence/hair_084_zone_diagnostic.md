# Hair 084 contact-zone localization

The observed Z134-mm gate failure is a false conservative whole-triangle classification. Every flagged intersection segment remains above the frozen boundary. This conclusion does not change the boundary or grant visual acceptance.

One pinned Blender 5.2.1 LTS (`9e2066aef7ef`) process replayed the exact 083-repair/084-bound construction and exited 0. It stopped before `check_root_contacts`, piles, exception/report writes and model saves. All input hashes remained unchanged; no candidate was saved and no parameters were varied.

## Actual segments, not triangle extents

| Same-side layer pair | Total BVH pairs | Pairs whose target triangle reaches below 134 mm | Resolved minimum segment Z | Unresolved flagged pairs |
|---|---:|---:|---:|---:|
| Left overlap 0 / overlap 1 | 483 | 6 | 134.749490 mm | 0 |
| Right overlap 0 / overlap 1 | 452 | 6 | 134.755465 mm | 0 |

All twelve flagged pairs are resolved noncoplanar intersections. The conservative minimum target-vertex Z is 133.685529 mm on both sides, but that low portion of the triangle is not part of the intersection segment. The flagged segments therefore remain at least 0.749490 mm above the 134-mm boundary.

Lowest-segment witnesses, world millimetres:

- Left target triangle 3010 / other triangle 3728: `(-60.076711, 3.915556, 134.749490)`; complete segment Z range 134.749490–134.755866 mm.
- Right target triangle 3010 / other triangle 3728: `(58.424877, 3.915556, 134.755465)`; segment Z range 134.755465–134.769549 mm.

All target/head intersection pairs have target-triangle vertices entirely at or above Z134 mm. Cross-side target pairs are empty. Consequently there are no additional target/head or cross-side below-zone candidates requiring segment resolution in this audit.

## Below-zone depth sampling

Four implicated overlap-0 triangles per side extend below the boundary. Their portions below Z133.9999 mm were clipped for sampling, with triangle centroids and denominator-eight barycentric grids: 230 samples per side. All 460 samples are outside the opposing overlap-1 shell by unanimous three-skew-ray parity. No indeterminate or inside samples occurred; maximum sampled inside depth is zero and there is no penetration witness to report. The implicated overlap-1 triangles themselves do not extend below the zone.

The finite samples supplement the resolved segment locations; they are not a continuous maximum-depth proof or a global audit of every free-hair vertex against every receiver. Contacts wholly above the boundary were counted, not independently localized. Piles, front locks, root burial depth and visual appearance were not reviewed.

## Evidence hashes

- Source 076: `c7aeaf157f7d451d658c302c6a9300125ab145a7551b5fa8713288052c747050`.
- Replayed 084 builder: `06a21d80abaa7c291917bb7de11feae1c95f2aac6d9d053450f2c1b935a2e97d`.
- Diagnostic code: `66d39c76ca8166d06eceb15081a4eb506b4d666c3c9df5097a4cdca45b316164`.
- Diagnostic JSON: `04ec912d498e2d992e42133ff83ca9d81869e2bc2fb0d016bc0d555f46e3d672`.

The JSON binds the inherited source and intersection-helper chain, retains every flagged pair and its actual segment-Z bounds, and records sample counts and limitations. No unresolved flag was silently classified as safe.
