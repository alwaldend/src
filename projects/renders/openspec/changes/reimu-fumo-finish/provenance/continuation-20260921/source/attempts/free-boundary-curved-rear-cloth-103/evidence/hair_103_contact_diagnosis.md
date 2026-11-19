# Hair 103: new front-boundary rim contacts

## Finding

**All 22 supplied new pairs are new rim contacts. None involves a deformed pure-rear sheet face.** They occur on two adjacent front-boundary edges per side, near Y = -10 mm and Z = 137.3–141.8 mm. Every implicated rim vertex equals its original paired source coordinate exactly: measured displacement **0 m**.

The topology change is causal. Each original edge joined one retained pure-front face to one mixed front/rear face. Removing the mixed face created a free boundary; closing that boundary outer-to-inner introduced new triangles through the retained overlap-0 side cloth. This is the upper lateral front-cut closure, not the released rear sheet, the Z=160 mm crown cut, or the lower rear hem.

The report examines only the 22 pairs exported by the coordinator. No broadphase reconstruction, Blender/native execution, model mutation, shape trial, or additional artifact was performed.

## Measured actual intersection locus

All 22 provided triangle pairs independently resolve to noncoplanar line segments using the existing hash-bound scalar edge/plane helpers. No near-coplanar or sub-0.1-micrometre case occurs in this set. Dimensions below are millimetres.

| Measure | Left overlap 0 | Right overlap 0 |
| --- | ---: | ---: |
| Provided new pairs | 10 | 12 |
| Distinct new cover rim triangles | 3 | 4 |
| Segment X bounds | -58.254570 to -56.937306 | 55.532918 to 56.689820 |
| Segment Y bounds | -9.992784 to -9.868186 | -9.988931 to -9.868400 |
| Segment Z bounds | 137.334934 to 141.755255 | 137.335665 to 141.667822 |
| Receiver outer / inner sheet pairs | 6 / 4 | 8 / 4 |
| Sum of triangle-pair segment lengths | 5.104719 | 4.901545 |

Segment-length sums describe these supplied pair segments, not a physical seam length or an independent curve-connectivity proof. The shortest segment is still resolved: right cover `21101` / receiver `1559`, length 0.007217 mm.

## Exact rim and original-edge mapping

The finalized 103 array has 5,364 paired outer nodes and 10,434 outer triangles. Its shell labels and canonical boundary indexing verify every exported rim. Adjacent outer-face ancestry is `front`, zone `low`, for all four implicated boundary edges.

| Side | New cover triangles | New boundary edge | Original outer edge | Retained source front face | Removed source mixed face |
| --- | --- | --- | --- | ---: | ---: |
| Left | 21078, 21079 | 319–318 | 5–4 | 25 | 3567 |
| Left | 21081 | 320–319 | 6–5 | 60 | 3573 |
| Right | 21100, 21101 | 371–370 | 57–56 | 3100 | 3773 |
| Right | 21102, 21103 | 372–371 | 58–57 | 3086 | 3759 |

Original-node correspondence has one source weight exactly equal to 1 for each edge endpoint; these are not interpolated root/hem cut nodes. The complete implicated rim triangles compare exactly with the correspondence's original paired outer/inner coordinates. Their *vertices* are inherited, but their closing *triangles* are new, so they are not inherited whole-triangle contact identities.

The small export has null `source_kind`, `source_zone`, and `source_cover_face` for rims because a new closure triangle has no donor surface face. The classification above comes from the hash-bound 103 boundary correspondence and the original 094 edge incidence, not by treating null as rear ancestry.

For exact pair lookup, the supplied receiver triangle indices are:

| Cover triangle | Receiver triangles |
| --- | --- |
| Left 21078 | 1556, 1557, 1558, 4996, 4997 |
| Left 21079 | 4997, 4998 |
| Left 21081 | 1558, 1559, 1560 |
| Right 21100 | 1558, 1559 |
| Right 21101 | 1559, 1560 |
| Right 21102 | 1556, 4996, 4997 |
| Right 21103 | 1556, 1557, 1558, 4997, 4998 |

## Receiver rows and construction role

The actual receiver is the **083 retessellated side lobe**, not an untouched 066 grid. Its 3,612 vertices form 1,806 outer/inner pairs; dividing current indices by 37 globally would be unsupported after retessellation.

For only the implicated receiver points, their unchanged upper Y/Z coordinates uniquely match original rows **4–8**, columns **15–16** of the 48-by-36-cell source pattern. The tested upper-grid formula is

`t=i/48; u=.56*j/36; Y=-.022+.049*u; Z=(.145+.003*sin(pi*u))*(1-t)+(.084+.009*(2*j/36-1)^2)*t+.008*t^3`.

Maximum Y/Z correspondence residual is **13.37 nm**, consistent with native float arithmetic; the nearest competing upper-grid point is at least **0.765 mm** away. This identifies the source parameter neighborhood without claiming unchanged global triangulation. Barycentric localization of the actual segments gives:

- Left source row **4.080919–7.460194**, column **15.752908–15.916375**.
- Right source row **4.147738–7.459906**, column **15.757962–15.916094**.

Thus the receiver triangles occupy source row cells **4–7**, column cell **15**, in the upper attachment region, well before the hanging hem. The pairs hit both actual receiver sheets, not merely its rim.

The export reports left source/current counts 333/169 and right 344/180. With the coordinator's exact-key classification, the current sets consist of **159 inherited + 10 new rims** left and **168 inherited + 12 new rims** right. Those global sets were not reconstructed here. The new rims occupy the old upper attachment height band, but being above 134 mm does not make a newly created closing triangle an exact inherited contact. No contact exception or acceptance change is proposed.

Visibility remains for the coordinator's exact failed-array views. The measured construction role is a newly sealed front-side cut boundary crossing existing upper side cloth; this evidence does not support changing the rear depth field or treating these as the 102 lower mixed-face failures.

## Bindings

| Input | SHA-256 |
| --- | --- |
| `hair_103_failed_review/contacts.json` | `c2393e4a7513baa8df7484884dcfb4e6894ce5c6df1ea1f29290a09d96b926c7` |
| `hair_103_arrays.json` | `0bbfad0aafc68aa1e3a745294c0f03685cad2a76b7c92ac45b72f9a0caf3bbae` |
| Retained source 101 | `40ec32f8d2832f4a42e36c4ee4f27bdf213d1940ba0114a4fd2efcc476b9408b` |
| Original 094 packet | `492c89b457110d01273a1d8de0616dc266d18c5c441aa634c871395c54946b7f` |
| `hair_103_free_boundary.py` | `db656380842ede1d19f2581fc0a45a6c732ff7d415165499a8a9b0cf0a2b1c5d` |

Only pure tuple-vector/intersection definitions were extracted from `face_076_hair_diagnostic.py` (`680c070f…`) and `hair_084_zone_diagnostic.py` (`66d39c76…`), with full digest guards. No old builder or clipping helper was executed. Original receiver-grid roles were read from 066 construction and 083's retessellation/rear-lobe update; current pair geometry came exclusively from the finalized export.
