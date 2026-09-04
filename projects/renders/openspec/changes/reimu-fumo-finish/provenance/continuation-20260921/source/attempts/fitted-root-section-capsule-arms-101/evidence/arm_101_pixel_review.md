# Arm 101 independent bounded pixel review

Reviewer ID: `reimu-fumo-hair-094-pixel-review-agent-v1`.
Role: implementation-blind pixel reviewer, `/root/hair_094_pixel_review`.
Observation time: `2026-09-06 04:55:45 UTC`.
Record status: immutable review of the exact image/receipt bytes listed below.
Candidate identified by receipt: `40ec32f8d2832f4a42e36c4ee4f27bdf213d1940ba0114a4fd2efcc476b9408b`.

**Bounded arms verdict: RETAIN.** ARM101 combines front containment with a substantial rounded hand high inside the cuff in side and both three-quarter views. It improves the reference relationship over retained BOW098 and avoids ARM099's excessive concealment. The complete view set does not reveal a material new arm-root pinch, open gap, or whole-image regression that outweighs this local gain. This is intermediate arm retention only; the whole asset remains unaccepted.

## Independence and limits

I inspected all five fixed views and the uncropped ARM101 presentation, then paired all six with retained `bow_098_all_review`. I reopened canonical front and canonical-turn frames 10, 12, 16, and 21. I retain my earlier pixel observations of ARM099 and the failed ARM100 shape-direction diagnostic. No builder, model, geometry arrays, preflight, implementation messages, or other-agent review was read. No Blender operation was performed.

After reaching the pixel verdict, I read only image-digest fields, the candidate identifier, and the unchanged flag from the render receipt. All six independently calculated image hashes match the receipt. The receipt reports `candidate_unchanged: true`; the saved model itself was not inspected or hashed by this reviewer.

The images cannot certify hidden attachment, contact, or mesh validity. Absence of a visible crossing is not evidence of valid geometry. Occluded root regions remain uncertain. The turn references are oblique, and I make no exact overlay or landmark-tolerance claim. No absolute asset score or acceptance claim is issued.

## Evidence across the complete view set

| View | Bounded arm finding |
| --- | --- |
| Front | Hands remain within the sleeve outline. BOW098's low pale protrusions disappear, closer to the canonical front. The familiar sleeve span, face, and overall framing remain visually consistent. |
| Rear | Low hand protrusions also disappear. Small pale areas at the upper shoulder roots are slightly more exposed than in 098, but they remain localized beside the torso/head junction. I do not see a conspicuous detached root, open separation, or sharp exterior pinch. |
| Side | A rounded hand occupies a meaningful part of the upper cuff opening. It reads as a volume rather than ARM099's thin top-rim crescent, preserving the useful direction seen in ARM100. This agrees better with the high contained hand in turn10/12. |
| Three-quarter | The upper hand bulb is visible beneath the cuff's top edge while staying inside the overall sleeve outline. The sleeve no longer reads as vacant to the degree it did in 099. |
| Mirrored three-quarter | The same balance holds on the opposite side: a clear rounded upper-cuff hand with cloth visible below, without the old low protrusion. No material asymmetrical placement regression is apparent. |
| Presentation | The higher hand is partly occluded but still visible inside the near cuff. Front containment and the cleaner sleeve silhouette survive the presentation angle; the head, bow, dress, feet, and image readability remain visually consistent with 098. |

## Disconfirming checks

The closest upper-rim contact in side and the slightly more visible rear root areas are the strongest reasons to scrutinize this change. The hand still passes behind the top cuff edge, so these images cannot show the entire root or prove that it is correctly attached. However, unlike 099, enough rounded volume remains visible in side and both quarters to establish the intended contained-hand read. I do not see a new severe shape pinch or a visibly floating upper bulb.

The small irregular red area inside the lower sleeve remains visible, especially in side/mirror. It was already present in retained 098, so I do not classify it as newly introduced clipping. The lower sleeve remains open, but the restored upper hand now provides meaningful volume rather than leaving the largely vacant opening that drove the 099 reset verdict.

Reference turn10/12 supports a substantial hand within the upper opening; canonical front supports containing it behind the front sleeve silhouette. Rear controls 16/21 do not justify inventing an exact shoulder measurement from their occluded views. ARM101's gain is that the two clearly observable requirements now coexist across the complete rendered set. The minor root-exposure concern does not provide sufficient visual evidence to reset this bounded arm change.

## Exact image and receipt binding

Paths are relative to `out/reimu_fumo_finish/desktop_astra/`. SHA-256 values in this table were independently calculated for this review.

| File | SHA-256 |
| --- | --- |
| `arm_101_all_review/front.png` | `060de15cbdca4d5ba100c541f06c2a9a46f86e4e17e4250198207d2c624a858d` |
| `arm_101_all_review/rear.png` | `198ad2ec8073bc429531ae7d878f037ef2c747e219c83dba496fb1844ee02826` |
| `arm_101_all_review/side.png` | `917cf76c96d2e4c89114364c72f0d84a84c724c6819cbc6e92fb99bd6381191c` |
| `arm_101_all_review/three_quarter.png` | `5d0c1e2d4c1aca8a348abbc01698c25bb3fc6e1da6b0c7edf3193d1338a0acca` |
| `arm_101_all_review/three_quarter_mirror.png` | `c93683504de6cda5a2119f7399618b8012ba0765fc1d2f88d916e602533ede8b` |
| `arm_101_all_review/presentation.png` | `0bc915f803734fc587b04a75cf89310886d9daf4399d808cf63614c81391b199` |
| `arm_101_all_review/receipt.json` | `eaf7def05a74cc42947f79f1e6931624bf2f66b6cb547b609a85fd50c5739801` |
| `bow_098_all_review/receipt.json` | `afed4b96d051bcf50dd00c269500b943ce3d690f21fe5febc8f7d25dc346f592` |

The six retained comparison image hashes are recorded in this reviewer's `bow_098_pixel_review.md`; all six comparison images were reopened here. Canonical front and turn-extract 10/12/16/21 hashes are recorded in `hair_094_pixel_review.md`. The earlier visual observations are bound to images in `arm_099_pixel_review.md` and `arm_100_pixel_direction.md`. Those prior hash records are reused, not represented as freshly recalculated in this review.
