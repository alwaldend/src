# Failed096 bow contacts against retained094

**The new contacts occupy an upper lateral patch zone, not the Z186 top
closure and not the whole long root-to-collar boundary. The retained094 cover
and core already intersect the same bow loops in this zone.** This is evidence
of an existing bow/head assembly conflict, not only a new lap-edge defect.

Pure-data diagnosis only: no Blender/native process, model/goal writes,
parameters, geometry variants, or candidate. Inputs are the frozen096 plan,
exact failed repair arrays/preflight, and the actual retained094 assembly JSON.

## Actual096 contact loci

All values below are world millimeters; bounds are resolved edge/plane hits,
not broad contacting-triangle boxes. The scalar contact counts exactly match
native preflight: 165 left, 152 right. Of these, 162/148 are resolved
noncoplanar crossings; 3/4 sub-0.1 µm contacts remain explicitly unresolved.

| Actual segment measurement | Left lap / left bow | Right lap / right bow |
|---|---:|---:|
| X range | −30.717–−25.299 | 28.472–36.079 |
| Y range | 20.174–27.159 | 19.175–26.493 |
| Z range | 169.467–180.153 | 166.605–179.196 |
| Contact-zone Z span | 10.686 | 12.592 |
| Distance below chosen Z186 cap, by Z | 5.847–16.533 | 6.804–19.395 |
| Outer / inner / rim pairs | 84 / 74 / 7 | 71 / 74 / 7 |
| Z186 cap-rim pairs | **0** | **0** |
| Maximum endpoint XZ distance inward from boundary | 8.876 | 10.095 |

Exact array ordering verifies the outer, inner and rim partition. These are
Y-offset construction surfaces, not visible/hidden classifications. Each
lap's contacts form one endpoint-connected curve group at the explicit 200 nm
joining tolerance. The whole failed lap height is 95.352/99.000 mm: contacts
are localized to about 11–13% of that height near its upper side.

The intersected rim edges lie on the upper lateral taper, not the horizontal
cap. Left edge indices **254, 255, 256** intersect at Z179.669–180.153;
right indices **222, 223, 224, 236** intersect at Z178.643–179.196. Actual rim
source endpoints are recorded in JSON. Most crossings continue through outer
and inner **face interiors**, reaching approximately 9–10 mm from the nearest
projected boundary. Therefore adjusting only the seven crossing rim faces
would not address the measured failure.

This rules out “only the top closing seam touches.” It does **not** establish
whether this upper lateral zone is bow-occluded in the controlling views.
That is a pixel/occlusion question for the exact failed rear/side/front renders.
No claim of a long visible-boundary intersection is justified by these ranges.

## Retained094 already intersects the bow

| Existing contact, measured from095 snapshot | Left | Right |
|---|---:|---:|
| Cover094 / own bow pairs | 518 resolved | 490 resolved + 2 unresolved |
| Core094 / own bow pairs | 250 resolved | 238 resolved |
| Cover contact X range, mm | −29.894–−3.000 | 3.000–34.259 |
| Cover contact Y range, mm | 2.263–26.038 | 1.529–24.626 |
| Cover contact Z range, mm | 170.053–194.185 | 167.559–194.005 |
| Core contact Z range, mm | 170.721–193.211 | 168.206–193.099 |
| Cover pairs whose segment XZ boxes overlap096 contact zone | 138 | 126, including 2 unresolved |
| Bow triangles shared by096 and cover094 contact sets | 30 | 15 |
| Core pairs whose segment XZ boxes overlap096 contact zone | 62 | 58 |

Thus the baseline problem is in the **same upper lateral bow receiver area**,
not merely elsewhere on the crown. “Same zone” does not mean identical
intersection coordinates: the new laps are displaced in Y. Shared actual bow
triangle IDs and overlapping XZ loci supply the comparison evidence.
The baseline cover has two contact-curve groups per side and the core one;
pair counts and segment-length sums are tessellation-dependent, not measures
of physical penetration severity. Existing defects do not exempt new096
contacts from its frozen zero-contact gate.

## Consequence for the next decision

A receiver-shaped **lap-only trim** could target this bounded upper lateral
footprint, but cannot fix the already intersecting094 underlay/core. It must
also reshape or exclude face interiors, create a new closed rim and recheck
all resulting contacts; the present rim intersections alone are insufficient
trim instructions. Whether such an upper cut preserves the intended long
overlap requires the actual rendered views.

**Bow relocation** would address a receiver involved with all three surfaces,
but these data do not prove a safe displacement or preserved bow likeness.
It changes a protected096 object and is a new design/scope decision, not an
API/data repair. This diagnosis selects no parameters and grants no waiver:
it separates the new-lap failure from the inherited assembly defect so the
coordinator can choose the next bounded plan explicitly.

## Evidence and limitations

All segment endpoints, triangle IDs,096 source-face ancestry, contacted rim
edges, projected boundary distances and classifications are in JSON. Native
165/152 counts are reproduced independently. Existing cover/core counts are
pure scalar observations, not fresh native BVH receipts. Endpoint joining and
XZ distance are finite geometric diagnostics, not sewing, hiddenness,
containment or guaranteed three-dimensional clearance.

| Artifact | SHA-256 |
|---|---|
| `hair_096_bow_diagnosis.py` | `3f2ddc9569122eee3ade7b198b12dcf1d0b69859850d48c12b33882585ee9b90` |
| `hair_096_bow_diagnosis.json` | `7f8b7e99cad5eadcbf3d91f3fd706962030726769681e842cf12ebfa2c4e13a8` |
| Exact failed `hair_096_repair_arrays.json` | `bee174d66dcf5f765dcdfd99b1351439c387676e822b9d4ded67eb04df2bed1b` |
| `hair_096_repair_preflight.json` | `835208b3a352b40c42a930d8edd09f872fe17d6ff3f5ab1d15f723a989841cf1` |
| Retained `sleeve_095_actual_assembly.json` | `4a23e0161fa80d992aedce2598ac7943b887bf83c2d5f6b1f3f10cd74783c3bf` |

Fully read scalar functions are AST-extracted under SHA guards; no old native
builder is executed. All inputs stayed byte-identical. Standard-library Python
computation exited 0. The previous arm diagnosis was not rerun.
