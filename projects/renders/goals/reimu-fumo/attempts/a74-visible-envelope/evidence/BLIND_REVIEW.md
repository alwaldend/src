# A74 P0 implementation-blind visual review

## Review isolation

- Reviewer saw only the immutable A74 P0 beauty renders, contact sheet, the
  complete supplied reference packet, and the reference-source role map before
  reaching this verdict.
- The builder, topology, measurements, strategy notes, intended change, and
  prior-candidate verdicts were not inspected before scoring.
- Stage judged: visible head/hair envelope P0 in full-body context.
- Controlling sources: `canonical_front_25cm.png` for frontal identity and
  proportions; `canonical_turn_180.gif` for profile, rear, depth, and layer
  order. Physical front/side and sofa references supplied construction and
  soft-contact cross-checks.

## Absolute verdict from rendered pixels

- Unlabeled same exact subject and variant: **no**. The intact bow, eyes, and
  outfit suggest a Reimu-inspired doll, but the dominant head/face/hair mass is
  not recognizably the referenced Fumo construction.
- Intended-medium read: **no**. The new peach head covering reads as a smooth,
  low-poly hard shell or canopy around a flat graphic face, not as layered,
  softly stuffed brown fabric panels.
- Major visible failure present: **yes**.
- Decision under the stated internal `>= 6` gates: **RESET**.
- Absolute visual-quality-gate decision: **REJECT**. It also falls well below
  the approval threshold of 8/10 per applicable category.

## Scores

| Gate | Score | Result | Pixel evidence |
|---|---:|---|---|
| Macro identity | 2/10 | FAIL | The canonical compact brown hair-and-face mass is replaced by a tall peach dome and a large white square face. The profile is a projecting canopy rather than the reference's deep, layered plush head. The rear silhouette is almost absent. |
| Plush construction | 1/10 | FAIL | The crown is a uniformly smooth shell with polygonal facets; side locks are thin pointed blades. There are no readable filled fabric panels, sewn roots, broad compression, or layered rear locks. |
| Contact / integration | 1/10 | FAIL | Crown, square face, cheek blades, and bow do not form one supported plush assembly. The side view contains large empty separations and implausibly thin parts; the rear view exposes the front face and eyes instead of occluding them behind the rear hair mass. |

## Decisive discrepancies, highest impact first

1. **The face/head silhouette is fundamentally wrong.** The reference has a
   rounded beige plush face framed by dark-brown hair. A large, perfectly
   rectangular white face plane instead dominates the candidate, with abrupt
   square corners and no stuffed cheek/chin volume.
2. **The hair mass is the wrong color, topology, and proportion.** The
   reference shows a compact, dark-brown, panel-built cap with a low blunt
   center fringe and side locks. The candidate shows a peach/tan hemispherical
   canopy extending too high and too far backward, with no credible brown
   frontal envelope.
3. **Rear construction is missing.** Canonical rear views are almost fully
   occupied by overlapping dark-brown lobes descending toward the torso. The
   candidate rear view exposes the face, eyes, and white facial field and has
   no identity-defining rear hair silhouette.
4. **The profile cannot be a stuffed Fumo head.** The face reads as an
   extremely thin vertical card, while the crown floats as a separate arch.
   The unsupported air gaps, sharp lower edge, and abrupt depth break are
   incompatible with the soft, continuous front/side/rear cushion seen in the
   turntable and physical-side references.
5. **Cheek locks and roots are disconnected blade forms.** They taper to long
   hard points, sit outside the face without a believable seam/root, and do not
   reproduce the reference's broad, softly filled side-hair panels or their
   overlap with the lower face and collar.
6. **Surface character reinforces the wrong medium.** Broad unbroken glossy
   gradients and visible planar faceting make the new envelope read as molded
   plastic. The references show fleece/felt pile, panel thickness, soft
   irregularity, seam tension, and asymmetric stuffing.

## Required representation change

Do not tune this shell by adjusting vertices, smoothing, or adding texture.
Reset the P0 representation. The next candidate must establish, before rear
detail or microtexture:

1. a real shallow stuffed face/head volume rather than a visible square plane;
2. a dark-brown front hair field constructed from overlapping, lightly filled
   fabric panels with a broad blunt center fringe;
3. continuous profile depth and physically seated roots at crown and temples;
4. an occluding rear brown mass so no front facial graphics are visible from
   the rear; and
5. broad, soft cheek locks with finite panel thickness rather than pointed
   blades.

Render front, one profile, one three-quarter, and rear immediately after the
macro blockout. Reject again if any view reads as a helmet, card, bald rear,
floating canopy, or disconnected primitives. Only after those four views all
clear 6/10 should the approach advance to seam, material, or rear-lock detail.
