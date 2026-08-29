# Attempts 60–64 — direct complete-model and upper-bow modules

[Back to attempt index](README.md) | [Back to goal](../README.md)

Every attempt started from an ignored byte-verified copy of working-ladder
rung 3, SHA-256
`c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.
Neither tracked Blender deliverable changed.

## Attempt 60 — complete vertical allocation

The lower body was compressed to `80%` in world Z while the head/hair/bow
group was expanded to `110%` about the fixed `0.25 m` top. The larger head
share is directionally closer in front, but it exposes detached/crossing sleeve
outlines, buries the feet, and turns the skirt into a thinner rigid cone with
long floor tangencies. Both blind reviews reject it. The complete cross-junction
affine interface is retired.

- [Five-view parent/candidate matrix](../../../../../out/reimu_fumo_attempt_060_direct_module/parent_candidate_matrix.png)
- [First blind review](../../../../../out/reimu_fumo_attempt_060_direct_module/BLIND_REVIEW.md)
- [Second blind review](../../../../../out/reimu_fumo_attempt_060_direct_module/SECOND_REVIEW.md)

## Attempt 61 — head-only vertical allocation

Only the 28 head/hair/face/bow objects were expanded to `110%` in world Z.
One reviewer sees a modest relative allocation win; the stricter calibrated
review measures side head/torso overlap growing from about `6 px` to `42 px`,
with the collar entering the lower face and rear hair swallowing the upper
back. Both reject it absolutely. A relative-review disagreement cannot promote
a candidate with a measured contact veto. One-sided affine allocation is
retired.

- [Five-view parent/candidate matrix](../../../../../out/reimu_fumo_attempt_061_head_allocation/parent_candidate_matrix.png)
- [Calibrated blind review](../../../../../out/reimu_fumo_attempt_061_head_allocation/BLIND_REVIEW.md)
- [Independent relative review](../../../../../out/reimu_fumo_attempt_061_head_allocation/SECOND_REVIEW.md)

## Attempt 62 — bow affine

The upper pockets were narrowed and made taller while the tails and hard tie
were reduced separately. The upper macro span improves, but the bow becomes an
upright cat-ear/card assembly with a cuboid point-seated center. More
importantly, the builder transformed nine meshes from a fifteen-object module:
two zigzag and four root-fold curves stayed at their original world positions,
directly causing detached graphics. The entire affine family is rejected, not
retuned.

- [Five-view parent/candidate matrix](../../../../../out/reimu_fumo_attempt_062_bow_silhouette/parent_candidate_matrix.png)
- [First blind review](../../../../../out/reimu_fumo_attempt_062_bow_silhouette/BLIND_REVIEW.md)
- [Second blind review](../../../../../out/reimu_fumo_attempt_062_bow_silhouette/SECOND_REVIEW.md)

## Attempt 63 — continuous pinched upper pocket

A complete eleven-object upper/root replacement tested one connected closed
hourglass pocket. Its ruffles and appliques were derived from the same surface,
so the A62 detachment bug was removed. The representation itself failed: front
and rear read as one rectangular bow-tie pillow; exact side collapses into a
thick upright block; the central pinch is not a gathered crown saddle. Two
blind reviews agree that one bounded correction cannot pass. The continuous
hourglass representation is retired after one packet.

- [Full five-view sheet](../../../../../out/reimu_fumo_attempt_063_constructed_upper_bow/five_view_sheet.png)
- [First blind review](../../../../../out/reimu_fumo_attempt_063_constructed_upper_bow/BLIND_REVIEW.md)
- [Second blind review](../../../../../out/reimu_fumo_attempt_063_constructed_upper_bow/SECOND_REVIEW.md)
- [Complete parent inventory](../../../../../out/reimu_fumo_attempt_063_constructed_upper_bow/PARENT_BOW_INVENTORY.md)

## Attempt 64 — separate sewn fan lobes

Two independent `29 × 17` front/back fan patterns used corrugated throats,
low-frequency relief, and white turned strips sharing the red boundary
vertices. A material-pixel audit corrected the upper-pair authority before the
first review from a halo-inclusive `1.55 Wh` to `1.497 ± .014 Wh`; the complete
bow remains `1.989 ± .019 Wh`. The candidate reaches that front envelope and
its trim is mechanically integrated, but exact side shows two upright leaves
opening into a V/cup. Obliques remain paddles, the center is a pale gap, folds
do not affect silhouette, and trim reads as piping. It fails the written early
gate, so no saddle, applique, or material polish is added.

- [Full five-view sheet](../../../../../out/reimu_fumo_attempt_064_sewn_fan_lobes/five_view_sheet.png)
- [Implementation-blind review](../../../../../out/reimu_fumo_attempt_064_sewn_fan_lobes/BLIND_REVIEW.md)
- [Second blind review](../../../../../out/reimu_fumo_attempt_064_sewn_fan_lobes/SECOND_REVIEW.md)
- [Exact candidate manifest](../../../../../out/reimu_fumo_attempt_064_sewn_fan_lobes/candidate_v2_manifest.json)

## Result and next strategy

Attempts 60–64 shorten the feedback loop and remove several false interfaces,
but none advances the working parent. The reusable gains are the complete
fifteen-object bow inventory, surface-bound graphics requirement, corrected
material-only span, fixed five-view batch packet, and proof that whole-subject
pixels must veto local numeric success.

Attempt 65 ran a deliberately disjoint portfolio on separate parent copies:
low-resolution sculptable cages, an explicit front/turn/return-leaf
construction, and dense native surface relief. Its adversarial strategy review
finds that bow work no longer dominates the complete-model error budget and
restores the rejected-crown dependency. Its terminal record is
[Attempt 65](attempt_65.md); no branch may select a least-bad result.
