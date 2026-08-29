# A74 P0 visible-envelope geometry contract

## Result

`geometry_contract.json` is a source-only contract for the next visible brown
head/hair-envelope candidate. It preserves the canonical front asymmetry and
does not reuse A73's rejected hidden receiver cage.

The contract contains:

- all 14 exact canonical-front scanline rows from A73, including their
  separate ownership and uncertainty;
- exact pixel points for the independent left and right upper contours, the
  seven-point canonical hairline, and both non-mirrored cheek locks from
  `head_hair_curves.json`;
- the canonical-front and canonical-turn source hashes;
- the supported `.77-.85 Wh` compact visible field depth, phase-specific
  complete-profile targets, `.36-.38 Wh` observed rear-leaf extension, and
  `.02-.04 Wh` hair-panel core-thickness band;
- the six-row A73 same-crown counterexample, explicitly marked as a rejection
  of that registration rather than a source for hidden geometry.

## Validation

- The five extracted visible curves are exact point-for-point copies of the
  corresponding observations in `head_hair_curves.json`.
- Every scanline value is copied at full CSV precision from A73's
  `front_outer_scanlines.csv`; the first six continuous-field rows also match
  A73's written measurement table.
- The `Wh` scale, reference hashes, aggregate outer-profile band, rear-leaf
  target, and construction thickness match `head_hair_curves.json`.
- The compact-field and phase-specific profile bands match A73's written
  measurement audit. The profile frame sides are deliberately not averaged.
- A73's profile coordinates and its coordinate hash appear only in the
  exclusion/validation record; no naked-receiver coordinates, hidden seams,
  crown offset, or receiver-only depth split were invented.

## Builder boundary

This contract is sufficient to build and compare a visible-envelope coupon.
It is not permission to infer the occluded cushion or its seams. Rows below
roughly `v=.60 Wh` belong to overlapping front, cheek-lock, and rear/nape hair
layers and must remain separate in any later construction.
