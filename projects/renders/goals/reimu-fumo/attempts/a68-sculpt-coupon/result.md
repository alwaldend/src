# A68 result — cap-only native sculpt falsified

## Outcome

Reset. A68 established a genuine pinned Blender 5.2.1 foreground sculpt
workflow and preserved all protected model context, but two native Grab
checkpoints failed to improve `HD-01-rear-cranium-balloon`. Neither candidate
is promoted. The exact rung-003 parent and tracked reusable asset remain
unchanged.

## Work performed

- Identified visible `A44 continuous hair cap with smooth opening` as the
  cap-only coupon and deep-copied it to `A68_BackCap_Sculpt` in an isolated C0
  candidate. The original cap stayed hidden and recoverable.
- Froze and fingerprinted 177 surrounding objects, including face, fringe,
  bow, locks, rear curtain, garment, cameras, and lights.
- Rendered a pinned five-view C0 baseline. It confirmed full rear hair coverage
  but an oversized helmet-like cap and near-vertical rear wall.
- C1 applied a verified native Grab stroke to 35 vertices, with a maximum
  displacement of 1.194 mm. It moved the patch rearward and produced no
  material pixel improvement.
- C2 restarted from C0 and applied matched native Grab strokes from both side
  views. Each moved 35 vertices by at most 8.356 mm. It created paired dimples
  without changing the broad silhouette.
- Clean background reopen checks passed for both checkpoints with no changed,
  added, or removed frozen object fingerprints.

## Evidence verdict

- C1 side-view normalized RMSE from C0 was `0.0001235`, with zero side pixels
  differing by more than one percent. It was rejected.
- C2 changed more local pixels, but more-than-one-percent RGB deltas remained
  only `0.1157%` in side, `0.0706%` in three-quarter, and `0.1238%` in mirrored
  three-quarter. The visible result was two dents, not a shallower stuffed
  cushion.
- The implementation-blind C2 review scored likeness `5/10`, silhouette
  `4/10`, construction `3/10`, identity `6/10`, contact `3/10`, medium `3/10`,
  and presentation `7/10`. It rejected C2 and found no new bald patch.
- Criteria 001–004 therefore fail for this attempt. Criteria 005–008 remain
  unverified and unchanged; this bounded coupon did not claim them.

## Exact identities

- Protected rung 003:
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- C0 isolated candidate:
  `sha256:26c8613fe3eb17a1ddfcf7c8b596ed2aa264162b86d2b1e81acf7033d1fa75ba`
- Rejected C1:
  `sha256:67c9dbf7787749038ca168215647991f6b4df422f081c35a1b410852b9931557`
- Rejected C2:
  `sha256:eac856bbdf2de7942bc2e41bb4ad92a58cf192963a7d141031fe254bf5f895e5`
- Protected tracked reusable asset:
  `sha256:489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`

## Dominant failure and next decision

The failure is structural, not a need for more stroke tuning. The current cap
is a dense shell over a boxy underlying head cushion; a large Grab radius still
affected only a small side-surface patch. Stronger strokes produce dents while
the coupled head/cap silhouette remains fixed. Per the two-checkpoint stop
condition, no C3 is allowed on this representation.

The next attempt must reset from rung 003 and test the smallest coupled
rear-cranium interface: head cushion plus cap, with the front face opening,
fringe, locks, rear curtain, bow, and body frozen. It should begin with
reference-measured macro depth and a clean low-frequency deformation or
representation gate, not detail or another local Grab.

## Process audit

The reliable parts were isolation, foreground native-stroke proof, clean
reopen verification, frozen-context hashing, pinned batch rendering, and
parallel blind review. The waste came from treating brush radius as proof of
broad geometric support and from not measuring displacement direction before
C1. C1 then supplied both measurements, allowing C2 to reverse cursor
direction and bilateralize the test; the remaining failure falsified the
cap-only interface promptly. Future sculpt attempts should run a cheap
single-stroke support/direction probe before a full checkpoint and must verify
the coupled silhouette-owning objects before choosing a module.
