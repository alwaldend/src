# Retained094 actual assembly audit

Observed 2026-09-06 02:23 UTC. One repository-pinned Blender 5.2.1 LTS
(`9e2066aef7ef`) background process, four fixed threads, exit 0. Only
`hair_094_candidate.blend` was opened. No object/model edits, frame changes,
rendering, candidate creation, or saves occurred.

## Retained versus rejected geometry

The retained094 sleeves are **exactly the 067 loft plus the 068 cuff return**.
Both actual base meshes have 3,600 vertices and 3,520 source-grid quads.
Indexed quad connectivity matches, and the maximum coordinate residual is
**0 m** for every row 0–40 against 067 and every row 41–44 against 068.
069 changes four hair sheets, not sleeves or arms.

The first 80 vertices therefore are the verified retained cloth-center-surface
root, centered nominally at **(±24, −3, 70) mm**. They are not an evaluated
Solidify wall ring. The 35/38 mm attachment constructions belonged to rejected
trials; the 38 mm root was never retained087/094. Earlier feasibility reasoning
that called the 38 mm root “current” must not be applied to this baseline.

| Source row / indices, inclusive | Verified meaning |
|---|---|
| 0 / 0–79 | Retained root, 80 points exported per sleeve |
| 40 / 3200–3279 | 067 outgoing cuff lip, `t=1` |
| 41 / 3280–3359 | 068 return, 45° |
| 42 / 3360–3439 | 068 return, 90° |
| 43 / 3440–3519 | 068 return, 135° |
| 44 / 3520–3599 | 068 return, 180°, terminal base-mesh boundary |

The four return rows use 068's measured row39→40 tangent and 0.9 mm turn
radius. All five cuff rings' actual world arrays are in the JSON. These index
labels intentionally do not claim a mapping into evaluated Solidify vertices.

Measured bilateral root bounds: |X| **19.423–28.577 mm**, Y **−13.000–7.000 mm**,
Z **63.439–76.561 mm**. Outgoing cuff row40 center is approximately
**(±59.231, −2.391, 45.347) mm**, with |X| 47.662–73.248 mm,
Y −20.702–16.000 mm, Z 28.148–61.734 mm.

Each evaluated sleeve has 7,200 vertices / 14,400 triangles, zero nonmanifold
edges and positive signed volume. The existing centered Solidify thickness is
0.8 mm. These observations do not assert fresh sleeve self/collision clearance.

## Actual arms and receiver identity

Both arm meshes have 1,106 vertices / 2,208 triangles, no modifiers, centers
**(±39.500, −6.000, 56.000) mm**, and actual principal radii approximately
**9.000, 10.000, 34.518 mm**. Their distal principal-axis endpoints are
**(±68.978, −13.526, 39.693) mm**; proximal endpoints are
**(±10.022, 1.526, 72.307) mm**.

Actual sampled arm bounds are approximately |X| **9.611–69.389 mm**,
Y **−17.562–5.562 mm**, Z **37.471–74.529 mm**. The JSON separately labels
continuous ellipsoid extrema derived from each actual matrix, after checking
the underlying sampled unit sphere (maximum radial residual 1.386e−6) and
axis orthogonality. These analytic extrema are not additional measured mesh
vertices.

All five evaluated coordinate/triangle digests **and object fingerprints**
match the hash-bound prior084 and source087 evidence exactly: both arms, body,
red skirt, ivory hem. The current084/087 files were also SHA-verified, without
opening them in Blender. Retained094's protected087 fingerprints all match.

Inherited receiver hygiene remains explicit: arm/body nonadjacent self-pairs
0; red skirt **272**; ivory hem **1,688**. These are prior measurements tied to
the exactly unchanged receiver geometry, not newly rerun collision tests or
waived defects. The cloth parity guards remain failed.

## Export and unchanged guards

JSON contains **22 meshes**, **168,383 evaluated world vertices** and
**336,550 triangles**, per-triangle material indices, material slot names,
material/node metadata, world/local matrices, and modifier properties.
Included: both arms/sleeves/dashes, body, red skirt, ivory hem, both collars,
current head, `Hair094 continuous fitted textile`, four released side cloths,
and the five visible bow meshes: buried center, both gathered loops, both soft
tails. Full bow-name inventory records excluded hidden/nonmesh matches.

Scene inventory stayed **413 objects / 282 meshes / 103 curves / 48 materials**.
All object fingerprints, object metadata, material metadata, inventory,
frame1/subframe0, camera and scene state were identical before/after. All bound
input hashes stayed unchanged. JSON `guards` contains only `true` values.
Modifier RNA scalar/vector/ID-pointer values are exported; nested non-ID
pointers and RNA collections are not recursively expanded. This bounded audit
does not claim full Blender serialization coverage or a new acceptance gate.

## Artifact identity

All files are alongside this report under the assigned scratch directory.

| Artifact | SHA-256 |
|---|---|
| Retained `hair_094_candidate.blend`, before and after | `af5a61921ae309a69cfe98b0da092d206b2d0ff29c9123ac2a4886de5a0f9add` |
| `sleeve_095_actual_assembly.py` | `1f803b7ae3bdb48b8d6edbb07b91499dee8861e22120547a1160706e40b48ea4` |
| `sleeve_095_actual_assembly.json` (48,553,803 bytes) | `4a23e0161fa80d992aedce2598ac7943b887bf83c2d5f6b1f3f10cd74783c3bf` |
| `sleeve_095_actual_assembly.log` | `7ce9766b7ec74d354ffff0d300117bb25a97156d80cc15837a93fad6cf9be396` |

Full input bindings and the source assembly fingerprint are embedded in JSON.
Script syntax was checked before the sole native call. No second process,
station search, candidate, geometry tuning, model/goal/Git write, or reference
redesign was performed.
