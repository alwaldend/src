# Rig037 pixel-witness diagnosis

Retessellation is rejected for this reconstruction. All 85 objects retain
identical ordered loop triangles: 360,276 triangles per source/probe set,
zero changed polygon diagonals. All 225 paired rays at the 25 recorded worst
pixels hit the same first object, polygon and triangle, at exactly the same
double-precision depth. No geometry/triangulation correction is justified by
this evidence.

## Bound execution

One unsaved reconstruction used pinned Blender 5.2.1 LTS build
`9e2066aef7ef`, immutable retained033 and unchanged helper
`d9d6287af3bf7800341da688190a4b5b033b68d788f90103d42f42c8821ea439`.
The existing fixed-camera contract and all ten source/probe PNG hashes were
verified. No render, save, shape/material tuning, action or pose change was
performed. The source blend and existing images were rehashed unchanged.
The process exited successfully; no second reconstruction was run.

Rays use pixel centers plus eight quarter-pixel offsets. These are exact
orthographic CPU triangle intersections for the 85 mapped character surfaces,
not EEVEE raster/depth-buffer decisions or its TAA sample sequence. Floor,
lighting and shadow rays were outside this bounded inspection.

## Close competing surfaces

All depths and competing materials below match between source and probe.
The full JSON records triangle indices, vertex/loop indices, barycentrics,
normals, material names and up to six ordered hits at each center witness.

| View / pixel | First surface and material | First depth (m) | Competing surface and material | Depth behind first |
| --- | --- | ---: | --- | ---: |
| Front (313,270) | Right eye applique, polygon 38, vertices 38/103/102; `Eye dark embroidery.002` | 0.750625232291 | Eye back, polygon 1006; same dark material | 14.9745 micrometres |
| Front (313,270) | Same eye hit | 0.750625232291 | Crown/back hood, polygon 5989; `Hair brown clay` | 33.0193 micrometres |
| Side (242,328) | Crown/back hood, polygon 252, vertices 857/3233/6761; `Hair brown clay` | 0.767118404425 | Sewn cushion, polygon 5220; `Face fabric clay` | 1.86181 micrometres |
| Mirror 3/4 (231,278) | Left cheek lock, polygon 1220, vertices 648/666/665; `Hair brown clay.002` | 0.677311232172 | Sewn cushion, polygon 3572; `Face fabric clay` | 32.8132 micrometres |

Across the five side center witnesses, hood-to-face gaps are
1.8618–4.8316 micrometres; across the five rear centers they are
1.0782–7.4400 micrometres. Both material pairs are consistently
`Hair brown clay` in front of `Face fabric clay`. In all views combined,
117/225 sampled rays have a second hit within 10 micrometres, and 6/225
within 1 micrometre. These are near-coincident layers, not measured exact
coplanarity or a demonstrated CPU hit-order reversal.

The largest observed saved colors are consistent with different visible
layers: front source brown versus probe dark eye; side/rear source dark hair
versus probe warmer face-like color. This is a rendering-sensitivity
hypothesis, not an object-ID proof of which EEVEE fragment won.

## Coordinates and normals

Maximum first-hit interpolated normal-vector difference is `9.7144e-8`.
Independent source Generated reconstruction is available on 108/225 sampled
first hits using source-rest coordinates for topology-preserving, only-deform
mesh stacks. Their maximum source/probe field difference is `2.6493e-8`.
The same probe vertex field interpolated over both surfaces agrees exactly
at every paired first hit; no triangulation-induced interpolation drift was
found.

The crown/back hood includes Solidify. Its true source ORCO was not exposed
by this non-rendering inspection, so no circular comparison with the helper's
constructive companion is presented as independent field evidence. Actual
Generated/implicit-shader equivalence on that hood remains unresolved. CPU
field agreement elsewhere does not establish GPU derivative or shading
identity.

## Handoff

Do not claim unchanged appearance or approve a save from these numeric
results. Close-layer raster/depth sensitivity remains plausible but unproven.
Root owns the next isolated side/rear crown/hood Generated-emission coupon
to discriminate the remaining Solidify ORCO uncertainty. This worker stops
here; no further reconstruction or render is authorized or performed.

Artifacts:

- `rig_037_witness_diagnostic.py`, SHA256
  `feb0be1f4d420acbf42d888da270b4a8f62c383d631994e262b60667d1414f58`.
- `rig_037_witness_diagnostic.json`, SHA256
  `545e6f8bc34402edb968d1807755827264a2a9b734c0c350082b7c8d0d4580c6`.
- `rig_037_witness_diagnostic.log`: successful single pinned execution.
- Preserved donor SHA256:
  `98e92ee9a73ff49be32695dc06518ff885e5d91016278d16fb5a8771fd8fed48`.
