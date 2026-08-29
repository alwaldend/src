# A89 S05 implementation-blind review

## Review conditions

- Stage reviewed: neutral-sculpt lower-rear hair module on the frozen head,
  front hair, and crown state.
- Candidate evidence: the fixed `profile_left`, `profile_right`, and `rear`
  beauty renders for S00 through S05. The S05 beauty renders were judged
  absolutely before comparing the earlier states.
- Controlling reference: `canonical_turn_180.gif`, especially its true side,
  rear-three-quarter, and rear frames.
- Supporting construction evidence: `physical_side.png`, `turn.gif`, and
  `sofa.gif`. The latter is useful only for general soft-fabric behavior and
  cannot control rear dimensions because its rear hair is not visible.
- Context isolation: I was not given the modeling code, topology, object
  names, author's diagnosis, or intended S05 correction. This review is based
  on rendered pixels and the stated module boundary only.
- Limitation: the review views are not camera-calibrated overlays. Numerical
  landmark tolerances therefore remain unverified; the categorical failures
  below are large enough to be visible without an overlay.

## Absolute S05 review

S05 does **not** read as the referenced constructed plush hair. In the rear
beauty render it reads first as one broad, rigid shield or cape laid over the
back of the head. In both profiles, the lateral pieces read as nearly planar,
vertical fins. The canonical turn instead shows three separately stuffed
fabric leaves: a broad center leaf and two overlapping side leaves. Their free
edges remain readable through depth separation and soft shadow, their surfaces
bow outward, and their roots gather beneath the bow/crown interface.

The reference's rear mass is not a continuous shell. Its center and side
leaves have different depths, bottom lengths, and curvatures. The lower edge
is softly scalloped and mildly asymmetric. S05 replaces this with a nearly
continuous wall terminating in three hard, similar triangular points. Closing
the previous gaps has removed the construction cue rather than resolving it.

### Recognition and medium

- Unlabeled same-subject/variant recognition from this module: **no**. It is
  recognizable only as generic stylized brown/salmon hair on a head block,
  not specifically as the canonical Reimu Fumo hair construction.
- Intended-medium read: **hard molded plates / shield**, not lightly stuffed
  sewn fabric panels.
- Major visible failure present: **yes**.

### Five largest discrepancies

1. **The three-leaf construction collapses into one shield.** The canonical
   rear view preserves three overlapping free edges and separate volumes;
   S05's beauty render visually fuses them into one uninterrupted surface.
2. **Profile depth is too flat and vertical.** The side leaves appear as long
   fins with straight axial flow and acute tips. The references show softly
   bowed leaves whose thickness, curl, and changing separation are visible in
   profile and rear three-quarter views.
3. **Roots do not look gathered or seated.** The candidate's upper transitions
   are hard, broad corners. The physical references show the leaves entering
   a compressed crown/bow-root region, with overlap and cloth tension rather
   than an orthogonal hinge.
4. **The lower contour is too geometric.** Three similarly sharp points form
   a regular sawtooth edge. The canonical center leaf is broader and rounder,
   while the side leaves differ in length and angle and end with softened,
   fabric-like tips.
5. **Side-leaf placement and overlap are wrong.** In S05 the side pieces hug
   the outer wall and become edge strips in profile. In the turntable they
   overlap the center mass at different depths and remain substantial,
   independently readable leaves in rear three-quarter views.

### Scores

These scores apply to the A89 hair module at the neutral-sculpt stage, not to
the unfinished whole character.

| Category | Score | Pixel evidence |
| --- | ---: | --- |
| Overall reference likeness | 3/10 | Generic hair wall, not the canonical layered Fumo rear hair |
| Macro silhouette and proportions | 4/10 | Broad coverage is plausible, but side and lower silhouettes are too planar and regular |
| Manufactured construction | 2/10 | Shield/fin read; gathered roots and overlapping padded leaves are absent |
| Identity-defining hair features | 3/10 | The characteristic broad center leaf plus distinct flanking leaves is not readable |
| Contact, attachment, and occlusion | 3/10 | Coverage exists, but roots and panel-on-panel seating look hinged or fused |
| Intended-medium read | 2/10 | Hard plates rather than constructed stuffed fabric |
| Presentation readability | 7/10 | Neutral lighting and fixed views clearly expose the failures |

Absolute decision: **reject**. S05 is not safe as a provisional checkpoint.

## S00-S05 progression

| State | Review |
| --- | --- |
| S00 | Reject. The central rear piece reads as a detached rectangular apron; the lateral pieces are thin, unrelated strips, including an implausible hooked lower curl. |
| S01 | Reject. Coverage improves, but the center remains a large flat hanging card and the side pieces read as isolated capsules/fins. |
| S02 | Reject. The profile is cleaner, but the rear remains dominated by a single broad plate with weak overlap and hard vertical side flow. |
| S03 | Reject. It produces the most continuous rear coverage, but continuity is the wrong cue: the result is effectively one smooth shield with almost no three-panel construction. |
| S04 | Reject, but it is the strongest diagnostic state. Three panel volumes become visibly distinct, yet large open slits expose the head, the side panels flare as broad shields, and the tips remain overly hard. |
| S05 | Reject. It closes S04's open slits, but does so by merging the panel masses. It removes the only strong construction cue and retains the fin-like profiles. |

No state from S00-S05 clears the absolute plush-construction gate. Relative
improvement is not sufficient to preserve S05.

## Reset verdict and smallest viable next test

**Reset S05 as the active candidate.** Do not reset the entire three-panel
representation: the canonical reference itself supports that representation.
Use S04 only as a diagnostic starting point, not as a survivor, because it at
least exposes the three independently readable leaves.

The next bounded test should change the three panels as a coupled overlap set:

1. preserve a broad center leaf but round and slightly asymmetrize its bottom;
2. reduce and bow the side leaves, move their inner edges behind the center
   leaf, and overlap enough to remove the white slits without making the free
   edges disappear;
3. stagger all three leaves in depth so a soft shadow line survives in rear and
   rear-three-quarter beauty renders;
4. curve and narrow the roots into the existing crown surface, avoiding hard
   top corners or hinge-like transitions; and
5. soften the tips and give each leaf a different length/angle matching the
   canonical rear and side silhouettes.

Binary keep condition: keep the next state only if all three free edges remain
readable in the rear beauty render, neither profile shows a vertical fin, no
head-colored slit is exposed, and the roots look compressed into the crown
rather than attached by hard corners. Otherwise undo it and reassess the
panel-root/depth representation before another contour edit.
