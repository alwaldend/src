# Failed arm095: actual collision locations

**The new crossings are in the lifted transition at the upper/rear sleeve
throat-to-mid region, not at the cuff or distal pole.** A small second cluster
occurs at the lower/front edge of the nearly frozen transition. No geometry
parameters were changed or proposed in this diagnosis.

This is pure-data computation from the exact failed arrays, native preflight,
and retained094 assembly export. No Blender/native process, model/goal write,
candidate, or parameter sweep was used.

## Counts: source versus failed arrays

| Own-sleeve triangle pairs | Left | Right |
|---|---:|---:|
| Source retained094 | 310 | 296 |
| Failed095 | 338 | 324 |
| Added | 63 | 63 |
| Removed | 35 | 35 |
| Same pair indices retained | 275 | 261 |
| Failed pairs touching changed arm triangles | 94 | 96 |
| Failed pairs on wholly frozen arm triangles | 244 | 228 |

All source and failed pairs resolve as actual noncoplanar crossing segments.
The independent scalar spatial-grid search reproduces the native counts and
the **exact recorded added and changed-triangle pair sets**. New segments are
at least 4.264 µm long, above the helper's 0.1 µm unresolved-span threshold.
Pair counts are tessellation-dependent, not independent penetration events.

## New-contact zones

Construction quadrants mean upper/lower from the sleeve cross-section and
rear/front toward +Y/−Y. They do not mean every point has positive/negative
absolute Y; the sleeve is centered near Y−3 mm. Both new clusters are on
evaluated half **B**; this report does not assume B means inner or outer wall.

| New cluster, per arm | Upper/rear | Lower/front |
|---|---:|---:|
| Pair count | 56 | 7 |
| Actual sleeve row interval | 11.630–14.366 | 12.986–15.000 |
| Loft parameter `row/40` | 0.291–0.359 | 0.325–0.375 |
| Actual segment abs(X), mm | 40.682–43.398 | 31.200–31.619 |
| Actual segment Y, mm | −0.760–5.127 | −12.829–−10.477 |
| Actual segment Z, mm | 69.180–72.460 | 51.911–56.306 |
| Interpolated applied blend weight | 0.757–0.954 | 0.000416–0.005807 |

The dominant upper/rear cluster comes from source-arm surface points at
|X| **49.771–54.845 mm**, Y **−8.033–−2.068 mm**, Z **56.630–60.936 mm**.
Under the existing failed displacement field, that transition surface moves
upward/rearward and inward to the measured upper/rear sleeve crossing zone.
This directly implicates the **transition deformation**, not the fully
translated distal endpoint. The seven lower/front additions arise in triangles
with very small but nonzero transition motion; their existence cannot be
dismissed because nearby vertices were frozen.

All failed crossings, including inherited ones, end by row **15.036**.
The actual root boundary remains row0 and the cuff lip is row40; new contacts
never reach either boundary. “Throat-to-mid” is only a descriptive label for
the explicitly reported 29–38% loft coordinate, not a newly invented seam.

## Distal target and cuff

The source had **28 cuff-return/rim crossing pairs per side**, on rows42–44.
All 28 disappear; the other seven removed pairs are around rows13/15. Failed095
has **zero cuff-region crossing pairs** and **zero crossing pairs on arm
triangles whose vertices are all in the fully translated source |X|≥58 mm
region**. Every new crossing's interpolated weight remains below 0.954.

Actual distal pole vertex0 moves from **(±68.978, −13.526, 39.693) mm** to
**(±56.978, −5.026, 53.993) mm**. Its unsigned point-to-entire-sleeve distance
is **11.054 mm** bilaterally. The nearest point is on cuff return row44.
Thus the distal target is not itself on the sleeve wall; the observed failure
is proximal to it. This is not a global minimum-gap or containment certificate
for the whole distal arm, nor an acceptance claim about the new silhouette.

## Mapping proof and reproducibility

The 7,200 evaluated sleeve vertices pair as indices `i` and `i+3600`.
Midpoints match the retained067 source formula for all rows0–40 within
**8.674 nm**, and all exported actual root/cuff ring anchors within **1.863 nm**.
All **7,040 half-sheet quads** and **160 end-rim quads** have exactly one valid
two-triangle triangulation, with complete accounting of all 14,400 triangles.
This proves indexed construction-row correspondence; it does not prove a
self-intersection-free shell. Seam columns79→0 are unwrapped before
barycentric localization. Rows41–44 describe the 068 return, not extrapolated
067 loft `t` values.

JSON records every source/failed/added/removed/changed segment endpoint, both
triangle indices, chart coordinates, source-arm corresponding points, blend
weights, and classification. No unresolved native pair was silently dropped.
Input SHA guards passed before and after computation.

| Artifact | SHA-256 |
|---|---|
| `arm_095_collision_diagnosis.py` | `ab4800f0c80d6529bbcbe927105055b91d33a3afe797f150c46170ed23e0e0db` |
| `arm_095_collision_diagnosis.json` | `269759de8e6a8ba1961d1748afb23ff14dc5fcc94c0cd3979ad3e5ad5a1e0aa3` |
| Failed `arm_095_arrays.json` | `f15d5b4300f8ae9f4603bcc8fced173893cb456d442ae57fe5e203628fd1fbff` |
| Native `arm_095_preflight.json` | `321f50c35836f539f98fb3b50491bdfe60e8ef060789b76154f9d576c0e248d8` |
| Retained assembly JSON | `4a23e0161fa80d992aedce2598ac7943b887bf83c2d5f6b1f3f10cd74783c3bf` |

The scalar helpers are individually AST-extracted from fully read, hash-bound
`face_076_hair_diagnostic.py` and `hair_084_zone_diagnostic.py`; their native
builders are not executed. Standard-library Python computation exited 0.
