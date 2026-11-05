# Arm 099 independent bounded pixel review

Reviewer: `reimu-fumo-hair-094-pixel-review-agent-v1`, image-only reviewer `/root/hair_094_pixel_review`.
Observed: `2026-09-06 04:12:53 UTC`.
Candidate recorded by render receipt: `925350ae56377cea29f22fbf55045f26e6df1b5040ca38e8efd4d4526ed1290a`; receipt reports `candidate_unchanged: true`.

**Bounded arm verdict: RESET this placement.** Front containment improves, but the hand becomes too hidden in side and oblique views, leaving a more conspicuously vacant sleeve opening. The controlling references support a contained hand that remains substantially visible high inside the opening. ARM099 does not yet preserve that read or clearly show the shoulder connection. The whole asset remains unaccepted; this report does not authorize stage advancement.

## Scope and independence

I reviewed all five fixed images and the uncropped presentation from `arm_099_all_review`, pairing every image with retained `bow_098_all_review`. I reopened canonical front, physical side, and canonical-turn extracts 10/12. No model, source, construction implementation, preflight, or other-agent review was inspected. After forming the visual verdict, I read only the render receipt's candidate hash and unchanged flag. The six image hashes below were independently calculated; I did not hash or inspect the saved model.

This is a visual arm/cloth placement judgment, not a geometry-validity judgment. Occlusion can hide crossings and attachment defects. Failure to see a crossing does not establish valid geometry, and inability to see the shoulder join does not prove physical detachment. The reference views are oblique; no calibrated overlay or landmark-tolerance claim is made.

## Paired findings

| Aspect | Evidence and finding |
| --- | --- |
| High contained hand | **Front gain, but excessive concealment elsewhere.** The low pale protrusions below the sleeves in 098 disappear in 099 front/rear, closer to canonical front. In exact side, however, the visible hand is reduced to a small pale crescent at the top rim. Both three-quarter views and the presentation largely lose a readable hand. Physical side and turn10/12 show a more substantial rounded contained hand near the upper opening. |
| Shoulder attachment | **Not visually established.** The sleeve remains seated beside the torso, but the upper arm/hand is too occluded to trace a convincing visible connection from shoulder into the sleeve. I see no definite new detached object, but cannot award an attachment improvement from concealment alone. |
| Arm and cloth appearance | **Material local regression.** Relative to 098, the opening reads more as an empty smooth bell/tube, especially in the two three-quarter views. The larger presentation preserves the cloth outline but removes the readable stuffed hand. This amplifies the cavity rather than showing a hand seated within soft fabric. |
| Pinch risk | **Unresolved image-level concern.** In side, the small visible pale arc is pressed close to the upper rim and can read as a clipped-off cap. Pixels alone do not distinguish deliberate occlusion from pinching or crossing. No definite new cloth crease or exterior pinch is established. |
| Gaping and red inner opening | **Empty-sleeve read is worse; the red opening is not new.** The small irregular red area inside the sleeve is already visible in 098. I do not label it newly introduced clipping. With less hand visible, the open interior and red patch are more exposed and visually distracting. |
| Whole-image effects | **Front simplification is useful; oblique arm readability is worse.** Face, bow, outer sleeve shape, dress, feet, framing, and general silhouette otherwise remain visually consistent at the available resolution. No material unrelated regression was found. |

## Disconfirming evidence and verdict rationale

The strongest argument for retention is the canonical front: the reference does not show the conspicuous low hands protruding from behind the sleeves that 098 has. ARM099 removes them cleanly from that view. I therefore support the goal of a higher, contained placement.

The strongest argument against retaining this exact result is the side/oblique set. The photographed hand occupies a meaningful part of the upper sleeve interior; it does not disappear into the top edge while the rest of the sleeve becomes a mostly vacant bowl. ARM099 moves too far toward that latter appearance. In presentation, the arms no longer read clearly enough as stuffed limbs associated with the sleeves. That local loss is significant enough to reset this placement despite the front gain.

The next visible result should combine the improved front containment with a substantial rounded hand inside the upper opening and a readable shoulder-to-arm relationship in side and both obliques. This states a visual target without inferring the unseen implementation. It also does not claim that 098 is an absolute arm-quality pass.

## Image binding

Paths are relative to `out/reimu_fumo_finish/desktop_astra/`.

| Image | SHA-256 |
| --- | --- |
| `arm_099_all_review/front.png` | `8e93a5214056bf1dfcd5a31c64cc6dd69a53f70c950d9dbf8fe3f3304132c594` |
| `arm_099_all_review/rear.png` | `817e9b63762907c48873788553eaf7771a8f3b80c6b1fcdcfb6f20afadc6c032` |
| `arm_099_all_review/side.png` | `db69d9f6f240bb2a5f242a4a2902013682cc2b8f30dca4b55252fd341aacd69b` |
| `arm_099_all_review/three_quarter.png` | `ee8428d0fec427dde2dd9c7f035958ee3264d440c06b7e3b24cc15da79383021` |
| `arm_099_all_review/three_quarter_mirror.png` | `f558848dee78969602a7fab2a9cf8755e9cc4ba31cac90967a85af66688dd936` |
| `arm_099_all_review/presentation.png` | `5729a3ab057b8ee66267d796d365fb842f26f7fadf0d42f8115c43ae4f40d14c` |

The corresponding six `bow_098_all_review` baseline image hashes are recorded in this reviewer's `bow_098_pixel_review.md`. Reference paths are `projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png`, `projects/renders/assets/reimu_fumo/references/physical_side.png`, and `side_062_frame_10.png` / `side_062_frame_12.png` in this scratch directory. Their recorded hashes are in `hair_094_pixel_review.md`. All named baseline/reference images were reopened for this comparison.
