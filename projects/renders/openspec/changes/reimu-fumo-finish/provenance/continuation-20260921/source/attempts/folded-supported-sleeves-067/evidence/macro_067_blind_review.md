# Candidate 067 independent pixel review

Observed 2026-09-05 18:27 UTC. Stage: macro plush construction, before final textile detail. Reviewer: `/root/sleeve_067_blind_review`. The coordinator remains the sole acceptance owner.

**Whole asset: REJECT. Bounded sleeve module: REJECT the fast gate.** The character is immediately recognizable as Reimu Fumo, but the sleeve opening still has a conspicuously regular tube-like construction in side and three-quarter views. This is a shape issue at the current stage. Missing fibers or stitching are not reasons for rejection.

## Evidence and limits

Reviewed only the three supplied candidate images and four supplied reference images. No builder, geometry, previous candidate, prior verdict, or intended fix was inspected. Administrative inspection established the feature branch and linked worktree; a context listing exposed filenames, but no implementation or previous-review contents were read.

- `canonical_front_25cm.png` controls front identity, graphic arrangement, and broad component proportions.
- `physical_front.png` controls actual front fabric construction. Its more upright bow and hem markings are variant evidence, not requirements to combine with the canonical image.
- `physical_side.png` controls actual cloth thickness, cuff/hand relationship, and foot covering depth.
- `side_062_frame_12.png` controls the open cuff's oblique profile and how the sleeve hangs beside the body. It is visibly a photographic reference.
- Candidate evidence is `macro_067_fast_review/{front,side,three_quarter}.png`.

Camera calibration, object measurements, a rear candidate, and a separate presentation image were unavailable. These are direct pixel observations, not calibrated landmark measurements or proof of a numerical tolerance. The available views cannot establish full retention or hidden contact correctness.

## Scores

Scores apply to the current macro construction stage. An 8 is the minimum absolute gate threshold.

| Category | Whole asset /10 | Sleeve module /10 |
| --- | ---: | ---: |
| Reference likeness | 7.3 | 7.3 |
| Silhouette and proportions | 7.7 | 7.6 |
| Manufactured construction | 6.8 | 6.6 |
| Identity features | 8.1 | 8.0 |
| Contact, attachment, and occlusion | 7.4 | 7.2 |
| Intended plush medium | 6.9 | 6.7 |
| Presentation readability | 8.2 | 8.0 |

Recognition: yes for the character; exact physical-variant fidelity is less certain. Medium: recognizable stylized soft toy overall, with several parts still reading as smooth manufactured shells. Major visible construction failure: yes, the sleeve/cuff form. No numerical average overrides it.

## Five largest visible differences

1. **Sleeve opening is too regular and stiff.** In the candidate side view, the opening has a nearly horizontal upper lip, nearly vertical sides, and a broad rounded rectangular lower edge. Its uniform bright rim and smooth inner wall resemble a short molded tube. The two side photos show an oblique, uneven oval opening whose upper and lower cloth edges pull in different directions. The supplied turntable frame particularly shows a long tilted opening with a locally collapsed edge. This discrepancy remains visible at three-quarter view.
2. **Bow panels lack the photographed depth and root gathering.** The side render shows very thin upper bow blades and a narrow trailing ribbon profile. The physical photos show cloth turning over itself with broader folds and a gathered, less uniformly thin root. The candidate's front bow span and red/white identity are convincing; the failure is the side construction, not a demand for the physical-front photo's different upright pose.
3. **Hair has overly continuous inflated volume.** The side and three-quarter crown is a smooth rounded shell ending in thick, rounded lower locks. The physical side references have flatter panel areas, more compressed boundaries, and thinner hanging ends. This is visible before fiber detail is considered.
4. **Hem layering looks untidy in depth.** The side and three-quarter candidate show a secondary white ruffle mass extending below/behind the main skirt edge as a separate irregular flap. The photos also have ruffles, but their dominant edge remains a coherent band following the skirt. I cannot diagnose intersections from these pixels; the visible layering itself is unclear.
5. **Black foot covering is too shallow in side profile.** The candidate side shows a long white pod ending in a very shallow black face. Both supplied side photos show black material wrapping visibly farther around the end. The front black circles preserve recognition, but they conceal this depth difference.

## Bounded sleeve finding

The front view succeeds at the broad white flared sleeve placement, red broken trim, and separation from the red bodice. The sleeve meets the shoulder region without an obvious floating gap in the supplied views. These successes do not resolve the opening geometry.

The hand sits as a small smooth pale lobe high inside the candidate opening. The high-resolution physical-side photo shows a more prominent rounded hand close to the opening, with the sleeve cloth folding around it. Pose and viewing direction introduce uncertainty here, so I rank this below the cuff profile. A small angular red patch is visible behind the candidate hand; it makes the interior depth order difficult to read, but these pixels do not establish clipping.

Recommended next bounded hypothesis: represent the cuff as an asymmetric fabric opening around the arm, with a tilted major axis, uneven rim tension, and broad flattened cloth faces. Preserve the front flare and trim placement, and validate the hand/cuff relationship from both available side reference directions. This recommendation describes the required visible construction without assuming the current implementation.

No material or microdetail work can substitute for that macro correction. A future fast pass would only authorize a full regression review; it would not establish rear retention or final approval.

## Input SHA256

| Input | SHA256 |
| --- | --- |
| Candidate front | `b02fc2dc7c37ec65bbbdb200506d373c96ded02ad251bf9fd85d13aa16eb412a` |
| Candidate side | `edad764eb07305bb4198164674ec40597ba31dec417698fa5f38fdb2b887453a` |
| Candidate three-quarter | `fb121ec08b77294f3cfe800d889819ffc1255220637e20165228b66b8fc9b67b` |
| Canonical front 25 cm | `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c` |
| Physical front | `f8c7d0f9911dbff1ef7f5d75601f9b10825015aecb367381971c076a5a3e7b51` |
| Physical side | `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d` |
| Side frame 12 | `4ebed9233acb351f2a8093496eb3e544771f1ee01e503cec84a1ce86f22b7365` |
