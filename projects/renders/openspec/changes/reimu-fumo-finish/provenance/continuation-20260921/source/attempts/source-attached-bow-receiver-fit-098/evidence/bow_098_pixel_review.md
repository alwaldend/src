# Bow 098 independent bounded pixel review

Reviewer: `reimu-fumo-hair-094-pixel-review-agent-v1`, image-only reviewer `/root/hair_094_pixel_review`.
Observed: `2026-09-06 03:50:20 UTC`.
Candidate recorded by render receipt: `2bee72c54ba571662e36fead903c9516ecb0f13b04d9138fbe897e0c1ffb6193`; receipt reports `candidate_unchanged: true`.

**Bounded bow verdict: RETAIN.** The wings now gather into the central knot more legibly, especially in rear and presentation. The sharper root pucker is a minor residual regression; it is not materially worse enough to reject this local connection gain. This retains an intermediate bow improvement only. The whole asset remains unaccepted, and this review does not authorize any stage advancement.

## Scope and independence

I inspected all five fixed views and the uncropped presentation from `bow_098_all_review`, then paired every view with `hair_094_all_review`. I also reopened the canonical front and physical front/side references. I did not inspect model state, source, implementation, preflight, or another agent's review. After reaching the visual verdict, I read only the candidate hash and unchanged flag from the render receipt to bind this report. Candidate-file hashing itself was not repeated by this image reviewer. The six image hashes below were independently calculated.

This reviewer retains the visual context of earlier 094 and failed-diagnostic reviews, including the earlier pucker concern. No rejected hair diagnostic is used as the comparison baseline here. The judgment is limited to bow connection, knot size, gathers, and newly visible whole-image regressions. No absolute asset score is assigned.

## Bow findings

| Aspect | Finding from the paired pixels |
| --- | --- |
| Connection | **Useful gain.** In 094 rear/presentation, the central knot reads more like a small smooth cap above the head while the broad wings meet behind it. In 098, the folds converge more clearly into the knot and the wings read as parts of one tied bow. The references support visible fabric gathering at this center. |
| Knot size | **Fuller, still proportionate for this bounded step.** The knot is more prominent in presentation and oblique views, but it remains small relative to the wings and does not become the dominant feature. The fixed front preserves the familiar bow/crown relationship. It could be softer and less rounded in a later refinement; I do not see a major size regression. |
| Gather quality | **Net improvement with a local cost.** The converging folds better explain the tied center. The closest wing root in side/three-quarter has a sharper concentrated pucker than 094. The physical references show pronounced gathering too, but with broader, softer transitions. This candidate is locally tighter and more angular. |
| Overall bow outline | **Preserved at the available resolution.** Wing span, outer tips, white edge, trailing-ribbon positions, and their relationship to the head remain visually consistent across the paired views. The added root detail is localized. |

## Disconfirming check: is the pucker a reason to reset?

I looked specifically at side, both three-quarter views, and the larger presentation. Side shows a small raised root bump; the oblique views expose the tighter fold against the knot. This is the strongest evidence against retention. It does not grow into a large spike, break the bow outline, create an apparent detached part, or read as a new conspicuous crossing in these images. In presentation, the root folds still read as compressed gathered cloth. The clearer tied connection is visible at ordinary viewing size, while the sharper pucker remains a local refinement issue.

**Answer: no, the pucker is not materially worse enough to outweigh the bow gain.** A later bow refinement should spread that concentrated pinch into a softer, broader fold while preserving the visible knot-to-wing connection. This is a visual recommendation, not an inference about the unseen implementation.

## Whole-image regression check

All six paired views were reviewed. I found no material new whole-image regression: the front expression and silhouette, visible hair outline, clothing, feet, framing, and overall readability remain consistent with the retained baseline. This is a bounded visual observation, not pixel equality, a surface-validity certificate, or approval of existing asset features. No new major bow floating, clipping, or accidental tangency is established from these pixels.

## Image binding

Paths are relative to `out/reimu_fumo_finish/desktop_astra/`.

| Image | SHA-256 |
| --- | --- |
| `bow_098_all_review/front.png` | `06deac13b6a400265cdd94bb5d09364fdfd478843d8c161752bcd0915d4cadcb` |
| `bow_098_all_review/rear.png` | `ee6ae827bd1113c4e6cbc02bcec5e8f3b0aa9bb105553428f1c6feeff98be05e` |
| `bow_098_all_review/side.png` | `61fd8c49fbddb97d74ff6464e504638686fd23936e76a70486a82c0af1196379` |
| `bow_098_all_review/three_quarter.png` | `031a36bbb8f34b84a3037b89896fa72918c4ff6e5a35f747328c3da6ba6cf4e6` |
| `bow_098_all_review/three_quarter_mirror.png` | `8d0b23421790584a9234443bf3e006d17c050eaba9b4deb22d495783cab48a3c` |
| `bow_098_all_review/presentation.png` | `eda0cc7cbe399e6b96bd6b0ce4bb8cd88199277f9b39f3a53aa20daa2840bb2e` |

The six retained baseline images are `hair_094_all_review/{front,rear,side,three_quarter,three_quarter_mirror,presentation}.png`. Their recorded hashes and those of `references/canonical_front_25cm.png`, `references/physical_front.png`, and `references/physical_side.png` are in this reviewer's `hair_094_pixel_review.md`. All named baseline/reference images were reopened for this comparison; no overlay-tolerance claim is made.
