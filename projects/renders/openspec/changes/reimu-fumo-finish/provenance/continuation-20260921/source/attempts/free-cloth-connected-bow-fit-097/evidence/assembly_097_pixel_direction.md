# Assembly 097 independent visual direction review

Reviewer: `reimu-fumo-hair-094-pixel-review-agent-v1`, image-only reviewer `/root/hair_094_pixel_review`.
Observed: `2026-09-06 03:36:40 UTC`.

**Verdict: RESET the current fitted-cap hair direction.** The rear boundaries and bow gathering show local gains, but 097 has not materially established broad sewn hair with long, independently hanging ends across the available views. This is a bounded direction judgment, not a recommendation to discard the whole model. It does not accept or retain this failed-build diagnostic as a candidate.

I inspected only rendered/reference pixels for the visual judgment, without implementation, preflight, model, or other-agent review evidence. I retain the visual context of my earlier 094/096 reviews. The coordinator disclosed the failed-build status and stale pile layer. Fine stray contours, fuzz, and surface pile are excluded; the verdict concerns the primary colored forms. No numeric scores are assigned. I used `hair_094_all_review`, the existing retained comparison set; the requested shorthand `hair_094_review` does not exist.

| Aspect | Independent visual finding |
| --- | --- |
| Broad rear hair | **Limited gain over 096; clear boundary gain over 094.** 097 has more legible large descending panel shapes, with outer lower corners extending beyond the central rounded mass. This is closer to panel-scale hair than 094's mostly uninterrupted back. It still reads as broad plates fitted around the same domed cap. |
| Free ends and overlap | **Main gap remains; insufficient reason to retain this direction.** Rear corners are longer than 096, but side still closes into a rounded back/underside. The mirror does not establish a broad hanging panel with a sustained exposed edge and independent tip. I do not count its fine hanging contour because the pile is stale. Reference frames 16/21 show broad lapping faces, visible depth order, and several staggered cloth ends forming the lower silhouette. |
| Relative to retained 094 | **The lost hanging silhouette is not convincingly restored.** 094 shows a definite filled tapered rear lock in side and mirror. 097 improves broad rear divisions but does not replace that silhouette with the wider hanging cloth ends required by the references. More visible boundary lines are not an equivalent gain. |
| Bow connection | **Modest gain.** In rear, the wings gather more visibly into the central knot and the gathering meets the hair crown more clearly. That gives a more legible connection than the previous smooth meeting under a small rounded knot. Front remains similar; the side still does not establish the full attached cloth construction by itself. |
| Bow folds | **Local regression/risk.** The added root gathering is tighter and more angular. Side shows a sharper puckered bump at the root, and mirror shows a more abrupt fold near the knot. The references support gathering, but theirs spreads into softer, broader folds. I cannot infer clipping from this alone. |
| Front silhouette and expression | **No obvious new regression at 512 px.** Crown envelope, central fringe, face, tied cheek locks, and the overall bow span remain visually consistent with retained 094 and diagnostic 096. This is not a pixel-equality or contact-validity claim. |

The strongest evidence in favor of continuing is the rear view: broad panels and their gathering now read more clearly. The strongest evidence against it is the side/mirror pair: the same rounded, close-fitting cap still controls the primary silhouette. Canonical-turn frame21's foreground hair piece has a long free edge and an oblique end lying over lower pieces; frame16 also shows several distinct lower ends. 097 mainly changes divisions and lower corners on a rounded mass, so the key manufactured-cloth mismatch survives.

The reset target should be a visibly different hair envelope: broad descending cloth faces whose exposed overlap continues into staggered free ends, legible together in rear, side, and mirror. Keep the front likeness and use the clearer bow-root gathering as useful visual evidence, while softening the new concentrated pucker. This describes the visual result needed; it does not infer or prescribe the unseen construction method. No claim is made that this four-view failed-build diagnostic passes any acceptance gate.

## Evidence binding

The four 097 diagnostic files were independently SHA-256 hashed. Paths are relative to `out/reimu_fumo_finish/desktop_astra/`.

| File | SHA-256 |
| --- | --- |
| `assembly_097_failed_review/front.png` | `6a81c023539df4202ea205804d6ff9e25ea9b1e71ae4f0f1a752ed89ad419853` |
| `assembly_097_failed_review/rear.png` | `176edc17b2fbca16b9bdab2678f6d48613eef002ce003c11ab298712d15876e9` |
| `assembly_097_failed_review/side.png` | `84a0dead6496ee9cf24259494af31d7a024d4f423e65e3ac22a91151a3256275` |
| `assembly_097_failed_review/three_quarter_mirror.png` | `1764254933fc1ef7586f401cab2808d37b3acadcd0acc67838130f939fd886f1` |

Reopened comparisons: the corresponding four views from `hair_094_all_review` and `hair_096_failed_review`, `side_062_frame_16.png`, `side_062_frame_21.png`, and `projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png`. Their hashes are recorded in this reviewer's earlier `hair_094_pixel_review.md` and `hair_096_pixel_direction.md`. The reference turn views are oblique; no exact overlay or landmark-tolerance claim is made.
