# Face 076 bounded contact audit

One read-only pinned Blender process exited 0. No model was saved or modified on disk. This is quantitative evidence, not visual acceptance.

- Source 071 SHA256: `129180df3be51e140d53f3b0bc4119327fb8993adcd9b079fad9807d396d0abc`.
- Candidate 076 SHA256: `c7aeaf157f7d451d658c302c6a9300125ab145a7551b5fa8713288052c747050`.
- Both hashes verified before and after execution. Blender 5.2.1 LTS, build `9e2066aef7ef`.
- Raw JSON SHA256: `c625f69817850fd481f1c80285875a2fa2066f5242d1f9b14f0d3b9f8fcabcc0`.

## Side panels: new intrusion, not just changed pair counts

All four evaluated side-panel world-coordinate arrays and triangle arrays are exactly identical between 071 and 076. Each panel was tested at all 3,626 evaluated vertices plus 1,812 index-strided triangle centroids, or 5,438 finite surface samples per panel per file. The head was finite and closed in both files. Inside classification requires unanimous parity from three skew rays; points within 1 micron are near-surface and disagreements remain indeterminate.

The two `overlap 1` panels retain exactly the same classifications and penetration depths. The two `overlap 0` panels do not:

| Diagnostic Z band | Left inside samples 071 → 076 | Right inside samples 071 → 076 | 076 max sampled depth left / right |
|---|---:|---:|---:|
| Below 110 mm | 0 → 0 | 0 → 0 | None found |
| 110–118 mm | 0 → 62 of 714 | 0 → 72 of 718 | 1.051 / 1.184 mm |
| 118–123 mm | 0 → 44 of 451 | 0 → 45 of 446 | 0.969 / 1.027 mm |
| At least 123 mm | 557 → 729 | 560 → 710 | 2.081 / 1.931 mm |

New lower intrusion witnesses, world millimetres:

- Left vertex 2886: `(-47.829, -21.870, 113.684)`; nearest head `(-48.660, -22.112, 113.087)`; depth 1.051 mm.
- Right vertex 2812: `(51.068, -21.869, 115.833)`; nearest head `(52.170, -22.127, 115.486)`; depth 1.184 mm.

The deepest overall contacts remain near the upper attachment, Z144.93 mm. Their maxima increased from 1.135 to 2.081 mm left and 1.131 to 1.931 mm right. The new Z110–123 mm samples lie below every previously inside sample, so they should not be silently classified as unchanged upper-root embedding. Visibility, physical support-region boundaries, and continuous penetration maxima were not established. Lower tips below Z110 mm show no sampled regression.

## Front-lock root seating

Root measurements use midsurface vertices, recovering candidate midsurfaces as paired front/back shell midpoints. Positive front-Y gap means the head is behind the lock at the same X/Z; it is not normal distance.

| Root band | Samples per lock | 071 front-Y gap | 076 max front-Y gap left / right |
|---|---:|---:|---:|
| Z ≥ 120.85 mm | 10 | 0.700 mm throughout | 0.700 / 0.700 mm |
| Top 0.5 mm, Z ≥ 120.3894 mm | 212 | 0.700 mm throughout | 1.115 / 1.072 mm |
| Top 1 mm, Z ≥ 119.8894 mm | 324 | 0.700 mm throughout | 3.129 / 2.969 mm |

The candidate constant-Y shell extends 0.425 mm toward the head from its midsurface, so the tip samples imply about 0.275 mm back-shell Y clearance. The tips have not newly detached by this measure. Just below them, however, the gap increases. In the top 1 mm band, maximum nearest-head distances are 1.179 mm left and 1.299 mm right, versus 0.520 / 0.553 mm in 071. No inside root midsurface samples were found; three left samples remain indeterminate.

Worst directional-gap witnesses, world millimetres:

- Left vertex 1204: `(-49.901, -34.315, 119.916)`; head front `(-49.901, -31.186, 119.916)`; 3.129 mm Y gap.
- Right vertex 1204: `(49.901, -36.348, 119.916)`; head front `(49.901, -33.379, 119.916)`; 2.969 mm Y gap.

This widening below the seated tip is consistent with a geometric wedge, but does not identify the dark image patch, establish daylight visibility, or distinguish tie occlusion from shading.

## Representative pile-root sample

Candidate samples: 254 of 18,000 strands per front lock, and 256 of 130,000 head strands. All sampled coordinates are finite and all sampled tip displacements point outward relative to the nearest receiver normal. Root nearest distances are at most 10.007 microns on the locks and 10.040 microns on the head. One or more sampled head roots have a negative local signed normal offset, minimum -5.355 microns; this small local embedding is retained as a raw finding, not waived as a full strand-contact pass. This checks rooting and tip direction only, not complete strand geometry or collisions.

## Coverage limits

No continuous-surface penetration bound, exact contact area, whole-scene/material audit, self-contact rerun, visibility/occlusion test, or visual acceptance. Triangle-pair counts are recorded only in the JSON and are not the severity criterion. Diagnostic coordinate bands are not independently identified sewing seams. Root samples do not cover the entire hanging lock.
