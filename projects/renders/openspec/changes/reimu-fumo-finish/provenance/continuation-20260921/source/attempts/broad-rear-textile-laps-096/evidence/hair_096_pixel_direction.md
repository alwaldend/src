# Hair 096 image-only direction review

Reviewer ID: `reimu-fumo-hair-094-pixel-review-agent-v1`.
Role: independent pixel reviewer, `/root/hair_094_pixel_review`.
Observation time: `2026-09-06 02:53:43 UTC`.

Scope: the four 096 diagnostic images (rear, side, front, and mirrored three-quarter), paired with corresponding retained 094 images and canonical-turn extracts 10, 12, 16, and 21. The coordinator identifies 096 as a failed-geometry diagnostic; this review does not inspect or independently establish that geometry status. **This is direction evidence only: no retention, acceptance, stage advancement, or numeric pass decision is issued.**

I read no new code, plans, preflight, model state, or implementation explanation. I inspected 096 pixels before reopening the paired 094/reference images. This reviewer retains the visual context of the previous 094 review. The image-only scope cannot determine surface validity or independently calibrate landmark tolerances. The turn references are oblique; findings concern visible form and layer order.

## Direction finding

The broad rear boundaries improve panel-scale readability relative to the nearly uninterrupted rear mass in 094. They do **not** yet establish the long overlapping cloth construction in reference frames 16/21. The candidate still looks like several fitted plates closing over a round cap. The side and mirror expose a clear silhouette regression: a hanging rear tip visible in 094 has been shortened or absorbed into the rounded back contour. Front expression and tied-lock silhouette appear preserved at 512 px.

## Ranked visible changes

| Rank | Change | Direction and evidence |
| --- | --- | --- |
| 1 | Rear tip silhouette | **Regression.** In 096 side, the back closes into a smooth rounded underside around neck level; 094 has a distinct tapered rear lock below that contour. In the mirrored three-quarter, 094's narrow hanging tip behind the near cheek lock is likewise no longer distinct in 096. Turn12/16/21 show free tapered cloth ends forming the lower and rear silhouette. |
| 2 | Broad rear boundaries | **Gain, limited to readability.** Long curved divisions now traverse the large rear mass, making broad panel-sized regions visible. Their scale and descending direction are closer to the reference's large hair pieces than 094's mostly blank rear surface. |
| 3 | Rear overlap and cloth read | **Major gap remains.** The new regions follow the same bulging cap envelope and meet near a shared rounded bottom. In turn16/21, one broad piece visibly lies over another, has a sustained exposed edge, and ends independently over the underlying layer. Candidate dark boundary lines alone do not show that depth order or free panel behavior. |
| 4 | Side attachment transition | **Mixed, not a demonstrated construction gain.** More rear boundary information is visible, but the near-side edge reads as a long close-fitting crease/ridge along a curved shell. The line high behind the crown also looks slightly irregular. Neither view yet shows a broad cloth panel whose root, exposed lapping edge, and free tip form one clear object. |
| 5 | Front, expression, and cheek opening | **No obvious regression found in the available view.** Central fringe, eye expression, two tied cheek locks, and the 094 side cheek opening remain visually consistent. This is a bounded visual observation, not pixel equality or a technical contact certificate. |

## Remaining construction gaps

The references read as long lapping cloth: broad roots near the crown, sustained descending faces, exposed side edges, and staggered tapered ends that help define the hair silhouette. Their layers have different visible depths and directions. In 096 the broad regions are legible, but they remain fitted to a rounded volume; their bottom edges largely disappear into or terminate on the cap contour. This reads closer to a segmented helmet or plates than sewn fabric hanging over a stuffed head.

The tip issue is more significant than the added boundary detail. In reference frame21, the broad foreground piece extends to a clear oblique tip while lower pieces produce separate descending points. Frame16 similarly shows multiple overlapping ends above the collar. The diagnostic's rear view remains dominated by a broad smooth bottom arc; its side and mirror remove an existing free tip instead of adding the needed long panel silhouette. More dark seams on that envelope would not supply this missing evidence.

Disconfirming evidence against dismissing the whole direction: the rear boundaries are a real visual gain, and the recognizable front has not visibly suffered. Disconfirming evidence against calling the direction successful: side and mirror reveal less hanging cloth silhouette, and the rear still cannot demonstrate sustained overlap or independent ends. The next visual experiment should preserve the broader panel readability while making the long exposed edges and staggered tapered tips visible in rear, side, and mirror together. This states the required visible result without inferring or prescribing the unseen implementation.

## Image binding

Paths are relative to `out/reimu_fumo_finish/desktop_astra/`. The four 096 hashes were independently calculated for this review. The unchanged reference and retained-094 hashes are recorded in `hair_094_pixel_review.md`; the corresponding image bytes were reopened here.

| Diagnostic image | SHA-256 |
| --- | --- |
| `hair_096_failed_review/rear.png` | `cd6e6956b7e4857d61507d60d0c8195f9103c656f36810a7d5d3620d8532b777` |
| `hair_096_failed_review/side.png` | `2467fe19410f165caffcd0d2158bc06ce6af7b0e558cac0e8204ae821ea3da3f` |
| `hair_096_failed_review/front.png` | `69850f0b5915a50b6e962753f62157883de1ff9efdce520bc358905b11b5e283` |
| `hair_096_failed_review/three_quarter_mirror.png` | `44b346d53fbf36f12ea8ffed2d575b20c0d57993ed8e3c85e3715291e9b774bb` |
