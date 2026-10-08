# Hair 094 independent pixel review

Reviewer ID: `reimu-fumo-hair-094-pixel-review-agent-v1`  
Reviewer role: implementation-blind visual reviewer, `/root/hair_094_pixel_review`  
Observed: `2026-09-06 02:10:01 UTC`  
Stage: intermediate construction/sculpt; final fibers and material finishing are excluded from scoring.  
Candidate: `af5a61921ae309a69cfe98b0da092d206b2d0ff29c9123ac2a4886de5a0f9add`  
Comparison candidate 087: `c171fea554811e825e2cd9d7f1d5a069c6b305e85888c6058d2390e1a20738e1`

Bounded retention verdict: **retain 094 as a small intermediate improvement over 087**. Absolute sculpt/likeness verdict: **reject**. The retained improvement does not resolve the major hair-construction failures or authorize stage advancement.

## Independence and evidence limits

I inspected the canonical front, canonical-turn extracts 10/12, physical construction references, and 094 fast pixels first. I then inspected all five complete fixed views, canonical-turn rear extracts 16/21, and the complete uncropped 094 presentation. I fixed my absolute scores before opening any 087 image. I subsequently compared all six 087 images. The task disclosed the candidate name, hash, bounded head/hair/chin scope, and comparison identity, but no intended implementation. I read the fidelity skill, its three review references, LANDMARKS.md, review_contract.json, the reference dossier, and render receipts. I did not inspect model objects, builder code, mesh diagnostics, preflight reports, plans, goals, or another review.

The render receipt binds the six completed views to the candidate hash and reports the saved candidate unchanged. I independently hashed the images below. Candidate-file hashing itself was not repeated by this image-only reviewer. Camera metadata agrees with the declared fixed-view setup, but I did not produce a calibrated reference overlay or independently certify landmark tolerances. Reference side/rear views are oblique and support construction and depth-order findings rather than exact orthographic pixel distances.

## Absolute review, fixed before baseline comparison

Unlabeled recognition: yes for Reimu Fumo; the exact photographed variant is not convincingly reproduced in construction. The eyes and central fringe preserve the recognizable sleepy expression. At this stage the head/hair assembly still reads partly as a smooth molded shell rather than assembled stuffed fabric. This judgment concerns the visible shape and layer order, not missing fibers.

| Category | Score / 10 | Visible evidence |
| --- | ---: | --- |
| Overall reference likeness | 6.0 | Recognizable character, but side/rear construction weakens exact-variant likeness. |
| Head silhouette and proportions | 6.5 | Front envelope is plausible; lower face is broad and flat, and rear hair forms an overly continuous oval. |
| Fabric construction | 4.5 | Crown, side hair, and rear remain smooth contiguous volumes without the referenced padded-panel overlap. |
| Hair layering and attachment | 4.0 | Rear lacks the broad overlapping cloth locks; side divisions resemble shallow grooves and tabs attached to a shell. |
| Face and expression | 7.0 | Heavy lids, small mouth, and fringe preserve identity; the surrounding face/cheek construction prevents an exact likeness. |
| Chin and neck contact | 5.5 | Contact is visually present, but the broad underside and narrow horizontal collar/red-neck strip do not show convincing stuffed compression. |
| Intended plush-medium read | 5.0 | Major head/hair forms remain molded in appearance even without requiring final material detail. |
| Presentation readability | 8.5 | Complete subject is readable under clear light; the presentation does not hide the hair-root and lower-face issues. |

Five largest visible head/hair/chin discrepancies, ordered by impact:

1. **Major: rear hair construction.** In `hair_094_all_review/rear.png`, most of the back is one round continuous cap, ending in a broad curved bottom with narrow side pieces. Canonical-turn frames 16/21 show large, separately overlapping, tapered fabric panels descending across the back and forming the lower silhouette. The missing overlap occupies a large part of the identity-defining hair mass.
2. **Major: hair roots and panel construction.** Both three-quarter views show broad smooth crown-to-temple transitions and thin groove-like divisions. In the references, the side locks have readable cloth edges, overlap, and a root-to-tip panel relationship. Simply seeing more boundary lines does not establish that construction.
3. **Major: lower-face cushion form.** Side and presentation views expose a broad pale side cheek and a fairly flat lower face/underside. The photographs show a softer compressed cushion with the cheek, hair opening, and underside turning into one another. The candidate's pale side region has improved visibility, but its shape is still too block-like.
4. **Moderate: tied cheek-lock outline.** The two front locks have smooth, strongly pointed leaf shapes with very regular curves. The canonical front has fuller, softer fabric-cut ends and less uniform thickness through the lock. This is a shape issue, not a request for surface texture.
5. **Moderate: chin-to-collar construction.** The lower face sits over a narrow, quite straight dark/red strip and clean collar edges. It lacks the reference's convincing seated compression and soft transitions at this contact. I do not infer floating or a mesh intersection from this image alone.

Major visible failure present: **yes**. No average score or relative gain overrides the absolute rejection. Within the bounded head/hair/chin inspection I did not find a newly visible clipping/floating defect that can be established from pixels; the major failures above are reference construction failures.

## Paired 087 comparison and disconfirming evidence

The useful gain is visible in the side cheek opening. In 087, dark hair almost encloses the pale side cheek, leaving a small oval-looking patch behind the bound front lock. In 094, the opening is wider and the pale cushion reads more continuously beneath the shortened hair edge. The physical side reference and canonical-turn frame 10 both show a meaningful cheek opening beneath the cloth hair, so this is supported by reference evidence rather than by an assumption about the intended edit. The front fringe/face boundary also appears slightly softer.

I specifically looked for evidence against retention. The newly exposed pale area makes the flat underside easier to see and can read as a larger white slab in exact side view. It therefore does not solve cushion construction and should not be treated as an absolute contact pass. Against that risk, the reference does support exposing this area, and 087's small isolated patch is less convincing. I found no new critical head/hair/chin regression in the front, either three-quarter view, rear, or presentation. The expression, front lock silhouette, crown envelope, and collar relationship are visually preserved at the available resolution. This is a bounded visual finding, not a claim of pixel equality or technical surface validity.

The rear view is particularly useful disconfirming evidence against a broader success claim: its large smooth hair mass is essentially unchanged in visible construction. A better cheek window does not materially repair rear layering. The old and new presentations also remain very similar, so the magnitude of progress is small.

Retain 094 only for this local improvement. The next substantive hair work should establish the broad overlapping rear cloth panels and their roots as visible forms, using the canonical rear controls, before further surface detail. Preserve the improved cheek opening while restoring the softer compressed cushion and chin/collar relationship. This is a visual recommendation, not a description of the unseen implementation.

## Image and contract hashes

Paths below are relative to `out/reimu_fumo_finish/desktop_astra/` unless stated otherwise. All SHA-256 values were read from the file bytes for this review.

| Image | SHA-256 |
| --- | --- |
| `hair_094_all_review/front.png` | `de0e0d25e0079a3e5243c380ef214c8a90083bf413e5138b6186cd1f549bc336` |
| `hair_094_all_review/side.png` | `8c6dae2cf78e1cce98b0385b11010c032a60ff762a00484d3b50a7a27d402535` |
| `hair_094_all_review/rear.png` | `88d95ba8ad6001e4e714c488933cd4bf06daa62d1e2d1df81de5fc7b7d54d9bd` |
| `hair_094_all_review/three_quarter.png` | `12fab2150d6e9c2b79ace4a07bae2da0d2ed2510c137c214957b27e56bf60b1b` |
| `hair_094_all_review/three_quarter_mirror.png` | `9a5c790506fc8ff41c96539239002f01315f4e7707732a567d7e8526c0154d66` |
| `hair_094_all_review/presentation.png` | `9cd4ca43ba44be18cc38161e4ddccf647e0d8a608f6245008b190f6757e13ff0` |
| `hair_094_fast_review/front.png` | `ffd3134f30ff7b86edf24bcb12e9a1d1ffa383e683e86a3d179522eda2e12fc2` |
| `hair_094_fast_review/side.png` | `b1a6bdb99507277e2958a97228cdc9069a69a58d41648af33f51b50bce3dd002` |
| `hair_094_fast_review/three_quarter.png` | `5d1f82aed69ff1373e7c9bebb0d342c304287ffd76617de91b020e02df88b848` |
| `chin_087_all_review/front.png` | `18610642926ca2be36a12ff0b843ae7b7d69fd2bb0a437df93634374babeeca9` |
| `chin_087_all_review/side.png` | `1e1e6d50be55d2d3528106bc39650932305d95e4bf31c321be71776552aff987` |
| `chin_087_all_review/rear.png` | `e041f94be6630fa2c1954e0d1231285f6b666fe2503bd3e0b777d12069f357d5` |
| `chin_087_all_review/three_quarter.png` | `dca74e014f7abf74918ec6d3329b0f4d479e4c1633bc3c39670a0dcac9dda197` |
| `chin_087_all_review/three_quarter_mirror.png` | `fc8deee7ea53c0e070f056afff2fe995e979fac4f71604220ca130000ae1a3d9` |
| `chin_087_all_review/presentation.png` | `72bd99787d9eae58c1b9c8e258d52b0f7072da9efb8ad1c356b6c32f021df016` |
| `side_062_frame_10.png` | `37c1e2866fdbe97ce79a0b5bdddf216a7f151f8e61d2a94c7286c22eaf41fd07` |
| `side_062_frame_12.png` | `4ebed9233acb351f2a8093496eb3e544771f1ee01e503cec84a1ce86f22b7365` |
| `side_062_frame_16.png` | `209fb44b435093dc1ce7e7339cd05503091d6d7815c471cda1d4e5b4843fc432` |
| `side_062_frame_21.png` | `685de63c429493eda0047253f875e1b20bd1ae7847d9610ed57ebf5ef4a1a333` |

Reference and contract paths below are relative to `projects/renders/assets/reimu_fumo/`.

| Source | SHA-256 |
| --- | --- |
| `references/canonical_front_25cm.png` | `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c` |
| `references/canonical_turn_180.gif` | `0d774eaa7f75828e388df4fb886cda7c563ce3bcd4ccb38d9885997a0846af30` |
| `references/physical_front.png` | `f8c7d0f9911dbff1ef7f5d75601f9b10825015aecb367381971c076a5a3e7b51` |
| `references/physical_side.png` | `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d` |
| `LANDMARKS.md` | `133d741d9252f34cee4e27990b0ea8d0710dc7408f43eb2a38b47d3c52a98140` |
| `review_contract.json` | `4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4` |
