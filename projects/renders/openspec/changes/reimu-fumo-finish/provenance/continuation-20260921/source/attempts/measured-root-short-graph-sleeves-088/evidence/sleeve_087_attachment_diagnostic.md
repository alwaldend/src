# Sleeve087 attachment feasibility diagnostic

Saved084 SHA256: `cb6bff1fa0e654f1720236a5a90dab56346f41e36f999dfc08fe746ea26294cb`; unchanged: True.
Pinned Blender 5.2.1 LTS; status: measurement_complete.

160 equal-arclength samples per actual arm section, moved radially outward 1.2mm. Counts below are point-inside / midpoint-inside / crossing-edge; ambiguous is total point, midpoint, and edge classifications.

| Side | Station center XYZ (mm) | Ring bounds min → max (mm) | Arm counts | Body counts | Ambiguous | Clear |
|---|---|---|---|---|---|---|
| left | -33.0000, -4.3404, 59.5957 | (-38.568, -14.066, 49.957) → (-27.432, 5.385, 69.234) | 0 / 0 / 0 | 11 / 10 / 2 | 0 | False |
| left | -34.0000, -4.5957, 59.0426 | (-39.602, -14.381, 49.345) → (-28.398, 5.189, 68.740) | 0 / 0 / 0 | 0 / 0 / 0 | 0 | True |
| left | -35.0000, -4.8511, 58.4894 | (-40.636, -14.695, 48.732) → (-29.364, 4.993, 68.247) | 0 / 0 / 0 | 0 / 0 / 0 | 0 | True |
| left | -36.0000, -5.1064, 57.9362 | (-41.662, -15.002, 48.140) → (-30.338, 4.790, 67.732) | 0 / 0 / 0 | 0 / 0 / 0 | 1 | False |
| right | 33.0000, -4.3404, 59.5957 | (27.432, -14.066, 49.957) → (38.568, 5.385, 69.234) | 0 / 0 / 0 | 11 / 10 / 2 | 0 | False |
| right | 34.0000, -4.5957, 59.0426 | (28.398, -14.381, 49.345) → (39.602, 5.189, 68.740) | 0 / 0 / 0 | 0 / 0 / 0 | 1 | False |
| right | 35.0000, -4.8511, 58.4894 | (29.364, -14.695, 48.732) → (40.636, 4.993, 68.247) | 0 / 0 / 0 | 0 / 0 / 0 | 0 | True |
| right | 36.0000, -5.1064, 57.9362 | (30.338, -15.002, 48.140) → (41.662, 4.790, 67.732) | 0 / 0 / 0 | 0 / 0 / 0 | 0 | True |

First bilaterally clear predefined station: |X| = 35mm.

No Blender objects constructed; no model save, simulation, or render. Mesh guards and all sampled coordinates are in the JSON.

Limits: Finite geometry diagnostic for four fixed stations; no construction choice or visual acceptance. The 1.2mm displacement is radial in the section plane, not a guaranteed normal clearance. Measurements cover 160 sampled points, all straight polygon edges, and all edge midpoints. Connecting equal-arclength samples uses chords and can omit original contour corners. Three agreeing rays corroborate parity; disagreement, near-surface/tangent/edge hits fail clear status. Mesh self-overlap guard excludes triangles sharing source vertices; adjacent foldovers are not separately proved absent. No guarantee about geometry between stations, full garment construction, shell thickness, or attachment visibility.

Post-run method review: ray origins advance 20nm and surface events within 100nm are merged. Those tolerances can skip or merge near-coincident distinct events without marking ambiguity; three agreeing rays are bounded numerical evidence, not an exact containment proof. No such distinct near-coincident events were established in this run. The reported minimum point/midpoint distances at 35mm are approximately 1.1815mm to the arm and 1.4058mm to the body; these are sampled distances, not a continuous clearance lower bound.
