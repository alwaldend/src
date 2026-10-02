# Chin 072 direction review

Observed 2026-09-05 19:21:47 UTC. Reviewed exact `chin_071_fast_review/{front,side,three_quarter}.png` and `chin_071_close_review/chin_close.png`, against previously inspected physical front/side, canonical front, and turntable frames 06/08/10/12/26/28. No model, code, or scene inspection. This is a proposal critique, not an implementation-blind acceptance review: the coordinator supplied the intended change and dimensions. Camera calibration and hidden chin-center depth remain unavailable.

**Recommendation: test mostly forward redistribution of the rounded lower-face support; omit or substantially reduce the 2.5 mm downward component initially.**

071 disconfirms the earlier need for a substantial downward extension. Its front mouth is around y320 and central cream underside around y342, against frontal head width about 210 px: approximately 0.10–0.11 Wh. This matches the roughly 0.10–0.12 Wh span in the frontal controls within visual measurement uncertainty. The close view now shows a convex cream cushion bottom. The former short visible span was largely occlusion. This supersedes the downward-extension recommendation in `chin_071_reference_diagnosis.md`.

At the supplied Wh of 116.195 mm, 2.5 mm downward is 0.0215 Wh, adding about four to five pixels in the frozen front render. That could move a currently plausible lower-face height toward the long end of the reference range and crowd the white collar. The reference supplies little evidence that this lengthening is necessary.

5 mm forward is 0.043 Wh. A smooth, broad redistribution below the mouth may make the cheek-to-underside transition less abruptly tucked in profile, but the full amount is a hypothesis, not a measured correction. Turntable frames 08/10 visibly roll backward under the mouth; hair occludes the center chin in the true side views. Do not flatten away all rollback, create a jutting lip/shelf, or claim a measured reference chin depth. The thin dark/red line above the collar is visible, but these pixels do not establish a physical gap.

Acceptance checks at the unchanged cameras:

- Front: mouth-to-underside remains approximately 0.10–0.12 Wh; preserve a broad, gently convex bottom, cheek width, and facial landmarks. Reject a longer or pointed chin.
- Side: smooth transition from mouth region through lower cheek into underside, with no local bulge ahead of the mouth plane, shelf, or sudden curvature break. Respect the reference's remaining natural rollback.
- Three-quarter and close: support should read as one stuffed cushion, with continuous cream surface and clean seating above the collar. Reject collar clipping, a new bright band/double contour, exposed core, or fibers bridging the contact.
- Judge the modified lower-face shape against references before consulting deformation values. A subtler change is acceptable only if the recessed impression actually improves; unchanged pixels do not validate the hypothesis.

The coordinator retains the final choice and gate. No claim of full-asset acceptance is made here.
