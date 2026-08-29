# A70 C0 v4 implementation-blind visual review

## Review context

I reviewed only the five immutable candidate renders and the controlling
appearance references: the canonical 25 cm front, canonical 180-degree turn
frames (both profiles, both three-quarter neighborhoods, and rear), physical
front, and physical side. I did not inspect the mesh, object names, topology,
measurements, or author verdict before making the visual decision.

## Verdict

**Reject C0 and reset the receiver/cap representation. Do not continue to
C1.** The candidate is recognizable as a rough Reimu model because of its
colors and bow, but not as the same physical Reimu Fumo variant. The new head
still reads as a rigid, deep, low-poly box/helmet rather than a shallow stuffed
cushion under thin layered hair. It misses the bounded internal gate by a wide
margin and triggers multiple early vetoes, so another parameter correction on
this representation is not justified.

## Absolute image review

- Unlabeled same-subject and same-variant recognition: **no**. The palette and
  motifs suggest Reimu, but the proportions, face framing, and construction do
  not identify the supplied Fumo variant.
- Intended-medium read: **rigid stylized primitives / helmet**, not a sewn
  plush neutral sculpt.
- Major visible failure present: **yes**.
- Absolute decision: **reject and rebuild the subsystem**.

| Category | Score | Pixel evidence |
| --- | ---: | --- |
| Overall reference likeness | 2/10 | Reimu colors are present, but the head, face, bow, and body silhouette are far from the controlling plush. |
| Macro silhouette and proportions | 2/10 | The head is tall, rectangular, and excessively deep; the canonical head is broad, rounded, and shallow with a compact seated body. |
| Constructed-plush logic | 1/10 | Hard planar breaks, a filled rear envelope, and disconnected panel-like slabs replace stuffed cushion and thin fabric-layer behavior. |
| Identity-defining features | 2/10 | The face opening is oversized and angular, eyes are tiny and misplaced, fringe framing is wrong, and the bow is too long, flat, and rigid. |
| Contact, attachment, and occlusion | 3/10 | Bow/hair roots and rear layers meet through abrupt steps and accidental-looking overlaps; side layers do not seat as soft fabric. |
| Intended-medium read | 1/10 | No broad compression, soft edge roll, sewn thickness, or stuffing response is visible. |
| Presentation readability | 7/10 | All five views are clear and neutrally lit enough to expose the failures. |

## Five largest visible discrepancies

1. **Head construction and depth:** both three-quarter views expose a deep
   cuboid shell, while the side shows a long flat crown and near-vertical rear
   wall. The references show a shallow stuffed receiver whose crown and
   underside taper before separately draping rear hair.
2. **Face/hair framing:** the beige opening rises as a broad, angular cutout
   into the crown and wraps around the side like a mask. The canonical front
   has a compact rounded face framed by a continuous soft brown cap and a
   controlled three-part fringe.
3. **Identity scale:** the eyes are far too small and low relative to the face;
   the head is too tall/narrow; and the body is too conical and exposed. These
   changes prevent same-variant recognition even before material detail.
4. **Rear layer construction:** the rear is a filled helmet/dome with large
   hard-edged petal shapes. The turn and physical-side references show thin,
   overlapping, independently draping hair panels with readable edge thickness
   and asymmetry.
5. **Bow and attachments:** the bow is extremely wide, straight, and boardlike,
   with abrupt roots and little folded-fabric compression. It dominates the
   silhouette differently from the compact, floppy reference bow and makes the
   blocky head read worse in profile.

## Early-veto assessment

- **S1/S2 side veto fired:** long flat crown, near-vertical rear wall, filled
  rear envelope, and no visibly separate free rear-leaf ownership.
- **Q1/Q2 three-quarter veto fired:** abrupt cap-to-lock steps, mask-like face
  opening, block silhouette, and inconsistent layered-hair ordering.
- **Front identity/aperture veto fired:** the beige domain and fringe clefts
  visibly drift far beyond the allowed tolerance.
- **Rigid helmet/mask veto fired:** this is the dominant read in every angled
  view.
- **New-defect veto fired:** the replacement creates obvious planar breaks and
  a boxier temple/rear transition rather than a soft crown-to-temple flow.

No C1 macro correction should be attempted on this geometry. The next test
needs a genuinely different coarse representation: first prove one shallow,
rounded, low-frequency stuffed receiver with a compact canonical face aperture,
then add a separate thin crown/temple hair skin and free rear leaves. Keep the
whole subject visible, but do not preserve frozen accessory contacts by filling
the rear volume; if the existing bow and lock roots cannot meet that receiver
without a bridge, widen the module boundary instead of recreating a helmet.
