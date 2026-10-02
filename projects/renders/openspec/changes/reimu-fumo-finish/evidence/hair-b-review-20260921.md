# Hair B independent visual review

Reviewer `/root/hair_b_blind_review`, 2026-09-21, inspected only the controlling
reference images and candidate renders. No implementation, metrics, history,
prior images, or intended correction were provided. Scope: neutral head/hair
macro cloth study; missing face, body, bows, ties, colors and texture were not
penalized.

Candidate SHA-256: `64dff253091b32548caefe17d11e133dacd0f462d4a7d59bd71e1d409104f236`.

| Category                     | Score / 10 |
| ---------------------------- | ---------: |
| Overall head/hair silhouette |          8 |
| Head and cheek volume        |          7 |
| Fringe shape and placement   |          8 |
| Side-lock shape              |          7 |
| Rear hair form               |          6 |
| Soft sewn construction       |          6 |
| Attachment and contact       |          7 |

Verdict: **reject for retention**. Five largest discrepancies:

1. Rear hem kicks outward like a rigid shell: close side x365–450/y380–445
   and close three-quarter x390–511/y350–440, in 512 px images.
2. Front-cap boundary forms an excessive ridge: close three-quarter
   x310–400/y75–285, producing a helmet impression.
3. Side-lock roots have abrupt planar transitions and hard contact shadows:
   close front x52–112/y275–324 and x398–456/y275–324.
4. Exposed lower head is too deep and shelf-like in profile: close side
   x150–370/y352–450.
5. Rear panel divisions resemble rigid overlapping tongues: rear
   x195–216/y172–296 and x285–314/y173–296. Exact rear likeness has lower
   confidence because the supplied physical reference angles are incomplete.

The front silhouette and asymmetric fringe are recognizable and close. Root
separately confirms disappearance of the broad pale intersection arcs, but
agrees the construction failures prevent retention or a modeling-stage pass.
The fit technique is supported within sampled coverage; the mesh is rejected.

## Frozen local evidence

Working bytes remain under ignored `out/reimu_fumo_finish/desktop_astra/`.
They are not a published reusable asset or guaranteed available in a fresh
checkout. This review binds to the local saved-byte packet below.

- Render receipt SHA-256: `48ba91c44dbec5e6c903ea137290e575db2f9758e8519bcfac6af195ce62ed8f`.
- `front`: `e04272ad31b551f061d057f35981cd60b2f0828b825b8699536ddb8f7da30947`.
- `side`: `617371a76b2abac029196a7ecbfebc4ca194a03da860392da5621070d6a90fd1`.
- `rear`: `c9cd4f7a3fbffc3103a4017f070e28fbf493abbcd54c731317181764ed0b639a`.
- `three_quarter`: `65f4766a7b8e97a11ef622a5d1c2859d422f62c8eab4c671b22789f7ea2452d8`.
- `close_front`: `bbd1922510a59d48c2fa59d077aa65894c9d3c31af0a4942118fcb030d87e060`.
- `close_side`: `82eeda258fefd3c717a45cc2369ffadd309a249f38a735b925f21d87fdfa09fd`.
- `close_three_quarter`: `4c5a1ae73d4aeb8845f018cd95c77f31816809008eba0470f9181e3d63925de7`.
