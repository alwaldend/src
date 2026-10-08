# Chin 087 independent pixel review

Observed 2026-09-06 00:19 UTC. Reviewer: `/root/chin_087_pixel_review`.

Verdict: **retain 087 as a bounded intermediate improvement; refine remains necessary. Absolute chin acceptance: reject.** This paired review supports no immediate reset from this pair alone. It does not establish that the user's recessed-chin concern is resolved.

## Independence and limits

Read the Blender reference-fidelity skill, scorecard, plush-construction guide, and absolute visual-quality gate. Viewed controlling references first, then 087 full renders, and recorded an absolute initial impression before opening 086. Viewed 087 close/isolated renders before their 086 matches. The task named the chin concern and candidate chronology, so this is implementation-blind but not fully hypothesis-blind. No scene, source script, diagnostic, topology, receipt content, prior review, or modeling measurements were inspected. No model, goal, or Git writes were made.

The coordinator identifies the close cameras as exact matches. This review confirms obvious paired framing agreement from the pixels but does not independently certify camera parameters. No calibrated reference overlay or measured geometry claim is made. All scores apply only to the visible chin, adjacent cushion construction, and local contacts. Hair quality, the whole head, rig, complete reference likeness, rear regression, presentation approval, and overall goal completion are outside this verdict.

## Reference authority

- `projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png` controls the frontal lower-face outline: a broad soft chin, with a gently curved center and cheeks that gather into it.
- `side_062_frame_10.png` and `side_062_frame_12.png` constrain head depth and the visible face-to-body relationship. Hair, side accessories, and viewing angle partly hide the real chin. Their hidden underside contour cannot be recovered from these pixels and was not inferred.
- The material and perspective differences between physical references and rendered candidates limit exact visual landmark comparisons. Candidate-to-candidate observations are stronger than claims about absolute unseen chin depth.

## Paired findings

1. In `isolated_three_quarter`, 087 carries a broader continuous roll from the lower cheek into the underside. The 086 lower face has a straighter descending retreat into a comparatively flat base. The added fullness in 087 is broad enough to read as part of one cushion.
2. In assembled `chin_close`, 087 modestly fills the space immediately above the collar. The chin's lower center sits visibly lower than 086 while the cheek corners remain similar. The difference is small, and does not substantially alter the front-view impression.
3. The fast side view exposes a rounder pale patch behind the front hair lock. The isolated view connects this exposed region to the lower-face surface; it is not sufficient evidence of an independent pouch. Its exact likeness to the real underside remains uncertain because that reference region is hidden.
4. Front and three-quarter full views show modest change. The frontal lower edge remains broad and nearly straight through its center compared with the gently rounded reference chin. The lower face still turns away quite shallowly, so some recessed appearance remains.
5. No new detached bulge, projecting ledge, abrupt shelf, accidental head silhouette distortion, obvious clipping, or floating chin is visible in these paired pixels. The dark separation above the collar persists but is smaller centrally; the visible red collar/neck region supplies continuous depth order. Hidden contact and surface intersections are not certified by this image-only review.

## Scoped scores

These are visual ratings, not measured tolerances. Scores below 8 prevent absolute acceptance under the visual-quality gate even when the pair improves.

| Category | 086 | 087 | Evidence and decision |
| --- | ---: | ---: | --- |
| Chin likeness | 7.0 | 7.5 | Fuller and less abruptly retreated underneath, but the central frontal contour remains too straight and the user's concern is only partly answered. |
| Cushion construction at chin/underside | 7.0 | 7.5 | The isolated lower face becomes a more continuous broad roll. It still needs a clearer soft transition through the front lower quarter into the underside to match the constructed cushion read. |
| Visible local contact and occlusion | 8.0 | 8.0 | No new visible contact failure; the assembled collar gap modestly improves. Occluded surfaces remain unavailable. |

The useful next target is the broad lower-face transition, considered together with the underside: preserve the continuous cushion roll, make the front chin's center and cheek transitions read softly rounded, and keep the collar contact clean. These images do not justify simply adding an isolated protruding chin or claiming an exact hidden reference contour. If this same category has already failed two reviewed cycles, the coordinator should apply the skill's subsystem-reset rule using that history; cycle history was not provided to this reviewer.

## Image provenance

Paths below are relative to `out/reimu_fumo_finish/desktop_astra/`, except the canonical reference, which is repository-relative. All supplied images were inspected at original tool detail. No images were generated, edited, cropped, or truncated for this review.

| Image | SHA-256 |
| --- | --- |
| `projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png` | `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c` |
| `side_062_frame_10.png` | `37c1e2866fdbe97ce79a0b5bdddf216a7f151f8e61d2a94c7286c22eaf41fd07` |
| `side_062_frame_12.png` | `4ebed9233acb351f2a8093496eb3e544771f1ee01e503cec84a1ce86f22b7365` |
| `chin_087_fast_review/front.png` | `d45b20fcaffcf736da031999ee2eb08c00dd105d5a8991fe0b14520eeb506746` |
| `chin_087_fast_review/side.png` | `67fdb670070cca3748a8b415cd59754501d7b09037ce5ca22e0ee101e16d03e8` |
| `chin_087_fast_review/three_quarter.png` | `fde60cf57c3e58d349b4d83f2bb6b7d7bea0a8bdc9970dd359cc84d9be3710d8` |
| `hair_086_all_review/front.png` | `68821187bbe6dd7e44c7925a849ba3223e53c2a80116df082174aa480e97de96` |
| `hair_086_all_review/side.png` | `c6f545c7ff4bb4fa519c5522b328aee45a0fdee74feaffdf20f0f27d4e168705` |
| `hair_086_all_review/three_quarter.png` | `ade42bed6e0dc3e083cb2d5e2fdc855fbc65411a5c0698542b511ec766c0a812` |
| `chin_087_close_review/chin_close.png` | `6debd8bbbca560652ca27091fd3223ce9df31ccf10b5a51c28e845e4787e72b6` |
| `chin_087_close_review/isolated_three_quarter.png` | `acfedc8cbe29597124e56d20b70b180008ea89273390398261011c8b71e6224c` |
| `hair_086_close_review/chin_close.png` | `6c6da84f13bc4a24783286afcba627d5c5bb7361c2f4fa638bcf3575e1386801` |
| `hair_086_close_review/isolated_three_quarter.png` | `17f463c67bef2da058dfe4bc71585f5e605f2d195b490c730f4f3282447765c4` |

The report was written only under ignored task scratch after confirming feature branch `t3code/continue-fumo-desktop-use` in linked worktree `t3code-a13ca48d`.
