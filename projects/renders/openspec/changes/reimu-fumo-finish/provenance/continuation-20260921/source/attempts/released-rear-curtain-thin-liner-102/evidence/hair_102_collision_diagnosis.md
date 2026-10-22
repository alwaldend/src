# Hair 102 actual-contact diagnosis

## Verdict

The **473 resolved core/body crossings are exact inherited source contacts**, not new damage from the 102 core fit. Every complete involved core triangle is unchanged, both source/current canonical pair sets and actual segments are identical, and no moved core vertex participates. The region is the retained lower head/body junction, Z **79.448–84.601 mm**. Calling it a concealed attachment is consistent with its retained role, but concealment in the final image remains a render judgment.

The cover failures are different: **721 of 726 cover/body pairs involve changed triangles**, and both overlap-0 receivers acquire new lower crossings. All those non-inherited contacts are on **mixed front/rear transition faces**, not the pure rear graph. Recognizing the inherited core attachment would not make 102 pass.

No Blender/native run, candidate, geometry change, parameter trial, model save, or goal/Git edit was performed. Only this diagnostic's PY/JSON/MD were written.

## Exact source binding and identity rule

102's source is retained `arm_101_candidate.blend`, SHA-256 `40ec32f8d2832f4a42e36c4ee4f27bdf213d1940ba0114a4fd2efcc476b9408b`. The 101 protected fingerprints match the 095 actual-094 export exactly for HEAD, COVER, BODY and both overlap-0 meshes. This is an individual-object binding; intervening arm/bow changes do not justify treating the entire assembly as unchanged.

The 094 construction packet agrees with the 095 actual head/cover coordinates after the exact native float32 conversion used by the builder. Its unconverted doubles differ by at most 7.44 nm. Triangle topology matches exactly. The 102 target topology and vertex indexing also match the actual source.

Canonical pair keys consist of both object names and each triangle's sorted vertex IDs. **Exact inherited** additionally requires equality of every actual world-coordinate component of both complete triangles, not rounded proximity or just reused IDs. Exact inherited segment endpoints are independently asserted equal. The JSON retains source/current records, canonical keys, indices, endpoint coordinates, ancestry, moved vertex IDs and triangle bounds.

## Pair accounting

Table counts are resolved noncoplanar crossings. Additional source numerical endpoint contacts are disclosed below rather than waived.

| Target / receiver | Source | Actual 102/native | Exact inherited | Reused pair IDs, changed geometry | Added pair IDs |
| --- | ---: | ---: | ---: | ---: | ---: |
| Core / body | 473 | 473 | 473 | 0 | 0 |
| Core / left overlap 0 | 181 | 181 | 181 | 0 | 0 |
| Core / right overlap 0 | 182 | 182 | 182 | 0 | 0 |
| Cover / body | 816 | 726 | 5 | 5 | 716 |
| Cover / left overlap 0 | 333 | 371 | 159 | 163 | 49 |
| Cover / right overlap 0 | 344 | 488 | 168 | 159 | 161 |

The four overlap-0 comparisons reproduce the native pair **indices**, not only counts. Cover/body reproduces 726 native pairs as 726 scalar resolved crossings. Core/body yields 473 scalar resolved crossings, matching native, plus one separately retained unresolved endpoint contact: core triangle `13319`, body triangle `60769`, span `1.0666408198256236e-14 m`, near `(0, 11.109771, 84.569363) mm`. That extra contact is also exactly inherited. Native body-pair indices were not exported, so this reconciles the count without claiming native index identity for that one endpoint case.

The source cover/body scan also retains four unresolved contacts in addition to its 816 resolved crossings; none survives in 102. Thus the JSON's raw source totals are 474 core/body and 820 cover/body. These finite-arithmetic distinctions are not acceptance waivers.

## Core/body: unchanged attachment, not fit regression

The resolved core/body segments occupy X **-13.078–14.708 mm**, Y **-13.500–11.110 mm**, Z **79.448–84.601 mm**. All 64 participating core triangles and their complete coordinates are exact source geometry. Source and 102 have no added or removed pair keys, and maximum involved vertex displacement is zero. Core/overlap-0 contacts likewise remain exact: left Z **141.206–148.000 mm**, right Z **141.160–148.000 mm**.

The 102 receiver loop does not actually classify core/body contacts as new: it rejects every such contact because its only inherited-contact exception is for overlap-0 receivers. Consequently the reported core/body failure lacks baseline classification. It is evidence of an inherited junction, not evidence that 102 introduced 473 new crossings. Any acceptance-rule clarification remains the coordinator's decision; this report does not add an exemption.

## Cover/body: new lower transition crossings

The source's resolved cover/body junction lay at Z **80.074–84.170 mm**. The 102 crossings lie lower and wider:

| Actual 102 cover/body measure | Value |
| --- | --- |
| Segment X range | -19.378 to 20.100 mm |
| Segment Y range | -9.364 to 13.188 mm |
| Segment Z range | 74.830 to 81.471 mm |
| Distinct cover / body triangles | 51 / 552 |
| Outer / inner / cut-rim pairs | 351 / 350 / 25 |
| Mixed front/rear transition pairs | 726 of 726 |
| Maximum displacement of an involved triangle vertex | 43.911 mm |

Only five cut-rim pairs remain exact inherited, all on three unchanged cover triangles near X -10.913 to -10.543 mm, Z 80.262–80.338 mm. Five other pair IDs persist but their cover triangles moved. The other 716 pair IDs are new; together these give **721 non-exact-inherited pairs on 48 changed cover triangles**.

The lower total count is therefore not an improvement claim: most old contacts disappeared while different lower transition surfaces entered the fixed body. Ancestry identifies these as shared faces connecting front-supported nodes to released rear-supported nodes. The rear PL liner/core ordering checks do not constrain those mixed faces against the body; the final full-surface gate exposed that separate construction failure.

## Cover/overlap-0: new lower side-transition crossings

All source overlap-0 contact segments were above Z 135.56 mm. The new below-134 mm groups are spatially separate:

| New lower contact group | Left overlap 0 | Right overlap 0 |
| --- | --- | --- |
| Added pairs below 134 mm | 38 | 144 |
| X range | -54.676 to -53.088 mm | 47.759 to 53.135 mm |
| Y range | 0.737 to 5.301 mm | 2.231 to 4.513 mm |
| Z range | 114.145 to 116.951 mm | 98.612 to 109.111 mm |
| Cover triangles / receiver triangles | 3 / 34 | 20 / 106 |
| Outer / inner / rim pairs | 38 / 0 / 0 | 81 / 63 / 0 |
| Maximum involved vertex displacement | 9.613 mm | 14.999 mm |

Every pair in these lower groups is an added pair key on a changed mixed-transition triangle. None is an inherited upper attachment, pure-rear-face contact, or cover-rim contact. All resolve noncoplanarly.

Upper contact totals must not be mistaken for exact preservation: 163 left and 159 right retained pair IDs now involve changed mixed-transition triangles; an additional 11 left / 17 right pair IDs appear above 134 mm while the same numbers disappear there. The exact inherited subsets are only 159 left and 168 right pairs, all on unchanged pure-front-ancestry faces. The native exception's object-level source-overlap test does not establish this stricter pair identity.

## Scope and receipts

These measurements establish inherited versus new geometry and actual segment location. They do not determine whether a crossing is visible through the finished cloth/pile, approve a redesigned transition, or rerun any lap/full-asset gate. Root's exact failed-array renders supply the visual evidence. No radius, displacement, topology or receiver-fit alternative was tested.

| Artifact | SHA-256 |
| --- | --- |
| `hair_102_collision_diagnosis.py` | `82957049a3746f9a987c2bf88593e81c7164b742e25c98b06dcad33beb879b1f` |
| `hair_102_collision_diagnosis.json` | `f1cd902fdb672054469938c8ba968b3077a213e801cea6890bd04e087dc04218` |
| `hair_102_arrays.json` | `5e3b39f6dc4776ad13d38d4102d689e1cbcbff9a5a2aa6d854ab9c30f67aa5fe` |
| `hair_102_preflight.json` | `962e93a7f0aea65906acaca1d36b94247cb6ff6f529b837a480cb9f315f9cebe` |
| `sleeve_095_actual_assembly.json` | `4a23e0161fa80d992aedce2598ac7943b887bf83c2d5f6b1f3f10cd74783c3bf` |

All additional input/helper hashes and the observation time are in the diagnostic JSON. It extracts only the bound pure scalar definitions; no old builder is executed.
