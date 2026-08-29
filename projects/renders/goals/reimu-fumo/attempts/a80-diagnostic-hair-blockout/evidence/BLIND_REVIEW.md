# A80 implementation-blind hair review

## Scope and evidence

This review used rendered pixels and the frozen A80 visual controls only. I did
not inspect either `.blend`, builder code, manifests, topology, object names, or
self-reviews before reaching the verdict.

The controlling evidence was the exact 25 cm front, the complete 30-frame
canonical turn sheet (with full-size checks at front, both 3Qs, both profiles,
rear brackets, and both rear arcs), the clean/physical front and physical-side
stills, all four older turn views, the selected sofa construction views, and
the five fixed renders of rung 003.

The constructed-cap board has no direct profile or opposite-side 3Q. The
sculpted-cap board has no 3Q or opposite profile. Those omissions prevent a
five-view pass, but the shown pixels already contain categorical vetoes.

## Absolute verdict

**Reject both A80 boards. Neither is better than rung 003.** Both preserve
enough unchanged Reimu context for subject recognition, but neither depicts
the reference hair assembly. Each has a major identity/coverage failure, so no
average score or localized front improvement can authorize promotion.

Hair-scope scores, from the pixels alone:

| Board | Likeness | Silhouette | Construction | Identity | Contact | Medium |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Constructed cap | 2/10 | 2/10 | 1/10 | 2/10 | 1/10 | 2/10 |
| Sculpted cap | 2/10 | 2/10 | 1/10 | 2/10 | 1/10 | 2/10 |

## View-by-view diagnosis

### Front

- **Constructed cap:** a clearly visible beige crescent remains between the
  raised brown crown band and the retained fringe. The band creates a second
  hairline and reads as a padded headband/hood around a round opening. Its
  narrow, steep arch does not follow the reference's gradual crown widening.
  This regresses rung 003's complete brown coverage.
- **Sculpted cap:** brown coverage is complete in this one view, but the crown
  remains a nearly symmetric constant-radius dome. It lacks the reference's
  offset, progressively widening crown and restrained planar stuffing. The
  deepest retained bang remains near the center rather than the canonical
  roughly `.101 Wh` viewer-right offset. This is at best a small front-only
  shape improvement over rung 003, not a solved front.

### Three-quarter

- **Constructed cap:** the upper side exposes a large beige patch, on the order
  of a third of the visible side-crown region rather than a one-pixel seam.
  The foreground band, beige support, round rear mass, and small lower pieces
  read as separate shells. There is no continuous crown-to-rear field and no
  broad bowed return leaf. The rear contribution is a ball, not layered felt.
- **Sculpted cap:** no 3Q was supplied. The profile's straight termination and
  bald rear make a passing 3Q impossible, but its exact 3Q silhouette is not
  measurable from this board.

### Profile

- **Constructed cap:** no direct profile was supplied. Its shown 3Q already
  exposes a bald side and spherical rear, so it cannot earn the profile gate.
- **Sculpted cap:** the brown front field ends at a ruler-straight vertical
  meridian. Roughly the rear half of the projected head is exposed beige
  support. The reference instead shows two brown contributions: a compact
  crown/rear field and an independent thin bowed leaf, with total depth about
  `1.14 Wh` on side A and `1.19--1.23 Wh` on side B. This board shows neither
  contribution and supplies no evidence of the required side asymmetry.

### Rear

- **Constructed cap:** one near-circular brown ball dominates the center. Its
  maximum width is near mid-height, it ends in a smooth horizontal arc, and it
  has no four/five unequal nape lobes or moving diagonal overlap. A separate
  capsule/banana sits on viewer-right while beige leaks around the remaining
  perimeter. The target is broad lower down, narrows under the bow, and reads
  as overlapping thin padded panels. Rung 003's rear is too symmetric, but it
  at least has complete brown coverage and multiple lower lobes; this is worse.
- **Sculpted cap:** most of the rear head is beige. A small top patch, side
  slivers, and a bottom teardrop do not form a connected hair mass. There is no
  rear width, lobe relief, diagonal T-junction, or free leaf to compare with
  frames 17--24. This is the strongest categorical failure of either board.

The physical side, older turn, and sofa views also veto both constructions:
the reference edges are thin, softly rolled, compressed at buried roots, and
visibly draped. A80's round solids and hard meridian/rim read as molded or
disconnected parts; material and fuzz would not repair that form.

## Representation versus parameter failure

- **Constructed cap: representation failure.** Moving or resizing the arch,
  ball, and capsule cannot create the missing continuous coverage, irregular
  nape boundary, diagonal overlap, and thin return leaf. Those are missing
  silhouette events and layer relationships, not scalar errors.
- **Sculpted cap: local parameter failure in front, representation failure as
  a complete assembly.** Its front dome could be flattened and made more
  asymmetric, but the visible surface simply terminates before the rear.
  Extending the same dome around the support would recreate rung 003's helmet;
  it would not produce the reference's draped rear panels and independent
  leaf.
- **Rung 003 remains only the comparison baseline.** It has full coverage, but
  its profile is an egg/helmet and its rear is a centered, near-symmetric petal
  mass. Neither A80 board replaces those failures with the required reference
  construction.

## Smallest next representation

**Decision: revise. Do not tune either A80 cap and do not add another sphere or
rear card.** The smallest defensible next module is one low-density,
**asymmetric crown-to-nape drape coupon** represented as a paired open-panel
assembly:

1. one continuous thin padded base panel owns the crown, both side returns,
   and the rear down to an explicitly traced four/five-lobe nape edge; and
2. one broad tapered return leaf shares a bow-hidden root, bows outward, and
   crosses the rear diagonally before returning toward the base.

This is one diagnostic module, not two unrelated add-ons. Shape it from the
front, both profile, and rear silhouette traces; keep the visible edge in the
`.015--.03 Wh` band. It can remain coarse and untextured. Render front plus the
weaker profile at 320--512 px immediately, then frame-18 rear only if that pair
has zero beige exposure and no second rim.

The strongest objection is that an open panel can become another curtain or
card. The first render should therefore reject it unless non-coplanar crown and
rear sections create soft volume, the nape changes the silhouette, and the
leaf produces a real moving T-junction. This coupon is still preferable to a
third cap: it directly represents every missing silhouette event and can fail
within one fast render cycle without promotion machinery.
