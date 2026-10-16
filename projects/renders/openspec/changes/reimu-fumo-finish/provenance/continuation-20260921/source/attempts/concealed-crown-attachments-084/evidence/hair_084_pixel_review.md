# Hair 084 independent pixel review

Observed 2026-09-05 23:35 UTC. Review scope: hair/head pixels only. References and candidate were viewed before the baseline. No scene, builder, topology, implementation receipt, prior audit, or modeling method was inspected. The review task disclosed the requested hair targets, so isolation was not fully blind to subsystem intent.

**Verdict: retain the bounded lower-layer improvement as an intermediate, then refine the hair module. Absolute hair approval: reject.** This is not whole-asset, rig, or export acceptance.

Compared with `face_076_fast_review`, the side and three-quarter lower hair is less blunt and bulbous, with a clearer tapered overlap. Front hair/face framing and crown silhouette have no important new regression visible in the matched views. Root emergence is not convincingly improved: small angular starts and narrow vertical separations remain visible. These pixels support a modest lower-end improvement, not a completed hair construction.

The additional rear and mirrored three-quarter images were reviewed. The rear cap has four short horizontal slit/crease-like marks around x200–305, y247–263 in the 512 px render. Their provenance is unresolved: no immediate `face_076` rear/mirror images were available in the bounded PNG search. The older `collar_050` rear has a materially different panel pattern and cannot establish whether 084 introduced these marks. Therefore **no-regression certification is limited to the matched front, side, and three-quarter views**, not the full rear.

| Hair-only category | Score / 10 | Evidence |
| --- | ---: | --- |
| Reference likeness | 6 | Recognizable framing and lock placement; rear profile remains too compact. |
| Construction | 5 | Ends separate more clearly, but cap, roots, and rear slit marks do not establish convincing fabric panels. |
| Intended plush medium | 5 | Smooth cap and sharp/narrow terminal pieces still read more like a molded shell than soft sewn hair. |

Top three remaining defects:

1. The broad continuous rear cap and short horizontal slits lack readable sewn-panel flow. The cap still dominates as a smooth helmet-like mass. This is a major absolute construction failure; its attribution to this cycle is unknown.
2. The side/back silhouette contracts inward near the neck. Frame 12 shows a broader, flatter rear taper extending outward from behind the bow. Bow occlusion and unmatched turntable angle limit exact dimensional claims, but the direction of the silhouette discrepancy is visible.
3. Root starts and ends remain inconsistent: angular upper starts/vertical lines in three-quarter view, a relatively blunt inner lower lobe in side view, and a narrow sharp outer rear tip. The reference has broader overlapping fabric-shaped tapers.

Reference control: canonical front for hair framing and front lock shapes; frame 12 for the exposed rear profile; frame 10 for overlapping side construction. No directly matched rear reference or calibrated overlay was supplied, so no numeric tolerance pass is claimed. The later mirrored view exposes no additional major hair/head failure. Further whole-asset approval evidence is outside this review.

Image SHA-256 (paths relative to repository root; `S` = `out/reimu_fumo_finish/desktop_astra`):

```text
864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c  projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png
37c1e2866fdbe97ce79a0b5bdddf216a7f151f8e61d2a94c7286c22eaf41fd07  S/side_062_frame_10.png
4ebed9233acb351f2a8093496eb3e544771f1ee01e503cec84a1ce86f22b7365  S/side_062_frame_12.png
03604bf1c43c7d0d3a0cdbd94b2a164bc4328637928516b44e144a519e79ab29  S/hair_084_fast_review/front.png
35b1a5428659157f5948350fe2ba8eeb685cd6ecc67f6fd1c56a241e8657f10f  S/hair_084_fast_review/side.png
a823dd9f6c3befa497703b7bafcdb81e497646eeb75707eaeb009dec5ef7a7d4  S/hair_084_fast_review/three_quarter.png
7e5c319814d501cb84bd97b935bac780a8727766fa05b4ea455693a996717713  S/hair_084_all_review/rear.png
9ba0a889990859e497ec35fd97ee758a9c5ba5b98ce46ebb8e1c71bf95403ecf  S/hair_084_all_review/three_quarter_mirror.png
edcdb6422426415cb84185fa6368405f4e598d44200a7e893ea9ec6386ff9810  S/face_076_fast_review/front.png
38fd498b087061527561cc0db39b84a4fc5316f25f9ced2fe64a3b88d1f3e734  S/face_076_fast_review/side.png
d7dd80519fc4bc01246b7d48eb8bf3696803b2699d56747437720dd86d2dcfcb  S/face_076_fast_review/three_quarter.png
813c26344bb08e1807ef1e030480d5584a99f8097ca65556d76facebace295a7  S/collar_050_all_review/rear.png
9d28083ca2e1bfbb403640601da99e29c3c9997a9e9084d580fef63546b686fc  S/collar_050_all_review/three_quarter_mirror.png
```

## Matched rear addendum

Observed 2026-09-05 23:37 UTC. The newly supplied `face_076_baseline084_review/rear.png` was compared directly with `hair_084_all_review/rear.png`; the coordinator identifies these as the same frozen rear camera and lighting. Its independently checked SHA-256 is `416f2ff4f58c60df807693366ccf039b8ac2f6a2ca9e9a1f6fd6843d996a26b6`.

The short horizontal marks are already visible at the same locations and with the same apparent shapes in baseline 076. They are **visually inherited defects, not a new 084 rear regression**. The broad smooth central cap is also inherited. The visible 084 changes concentrate on the lower lateral layer contours; these do not create an important new rear hair/head defect.

This resolves the earlier provenance limit and supersedes the restricted no-regression verdict: **retain 084's bounded lower-layer improvement; no important new hair/head regression is visible across the matched front, side, three-quarter, and rear images.** Continue refining the hair module. The three defect priorities, scores (6/5/5), and absolute rejection remain unchanged because inheritance does not make the cap and horizontal marks acceptable. No broader asset or rig acceptance is implied.
