# A83 next pixel-visible delta

## Evidence reviewed

- Current best: `C1b_packet/front.png`, `C1b_packet/three_quarter.png`,
  and `C3_baseline/rear.png`.
- Exact-variant authorities: `canonical_front_25cm.png` and all 30 frames
  of `canonical_turn_180.gif` (especially three-quarter frames 06--10 and
  rear frames 16--22).
- Construction cross-checks: `physical_front.png`, `physical_side.png`, all
  four frames of `turn.gif`, and the README-selected frames of `sofa.gif`.
- Graphic-only cross-check: `clean_front.png`.

## One defect

The camera-facing side/back hair still reads as one inflated spherical helmet:
the broad uninterrupted highlight rolls continuously from crown to rear, and
the visible lower lock reads as a radial petal cut from that sphere. In the
canonical turn and physical side, the same area is dominated by a broad,
thinly stuffed **hanging fabric panel** with a flatter middle plane, a readable
edge/overlap, and a slightly asymmetric lower point. This construction error
is the largest remaining hair mismatch visible in both three-quarter and rear
views; the front face and fringe do not need to move.

## Smallest connected owner and edit

Edit only the existing camera-facing rear lock (the outer lock visible beside
the cap in three-quarter and as one of the broad rear layers). Reprofile or
replace that single lock as a thin front/back fabric panel:

- keep its root seated under the bow/crown;
- preserve the current outer head bound and lowest hair bound;
- make the central span visibly flatter than the cap;
- give it finite soft edge thickness and one clean overlap/occlusion boundary
  against the cap;
- use the canonical turn's broad, mildly asymmetric taper rather than a
  centered rounded teardrop.

Do not scale or reshape the whole cap, head receiver, bow, fringe, face,
garment, or body. In particular, protect the accepted front crown/opening,
central fringe tip, eye positions, side-lock endpoints, overall height, and
head width.

## Views and binary decision

- **Controlling view:** fixed three-quarter render, compared with canonical
  turn frames 06--10 and `physical_side.png`.
- **Regression-risk view:** fixed rear render, compared with canonical turn
  frames 16--22. Render front only after a fast pass to confirm no visible
  face/fringe or overall-bound change.

**KEEP** only if the three-quarter pixels unambiguously show a separate,
slightly stuffed planar lock across a substantial side/back area, the rear
pixels show a clean layered overlap, and neither view gains floating edges,
intersection, double silhouette, or card-like flatness. The external hair
bounds must remain effectively unchanged, and the front face/fringe must be
pixel-stable.

**UNDO** if the result is merely a drawn seam or thin highlight on the same
round dome, if the new panel is not obvious at 512 px, or if it changes the
outer proportions/accepted frontal landmarks.

This is not another broad deformation: it changes one bounded owner and must
create a new visible depth discontinuity plus a planar highlight over a large
three-quarter region. If those two pixel cues are absent, the edit fails even
when topology or measurements changed.
