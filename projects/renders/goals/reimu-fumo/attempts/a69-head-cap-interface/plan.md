# A69 — coupled rear-cranium interface

## Immutable bindings

- Goal: `reimu-fumo`
- Goal resource version / next checkpoint CAS token: `3`
- Goal generation: `1`
- Lifecycle generation: `1`
- Criteria revision: `1`
- Criteria digest:
  `sha256:c5522700389e76975e7978515c586433ca2058a6d5012ef45fbbadcb78a5740c`
- Goal-state digest:
  `sha256:69d29116e14d282349e6ad5c073453839be4872922c886d9c74a7b5161fbae35`
- Exact parent rung:
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Protected reusable asset:
  `sha256:489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`

## Decision review

The cap-only hypothesis is falsified by A68: the visible cap is an almost exact
0.780 mm offset shell over the skin cushion, and three rear locks extend beyond
both. A domain expert would reject further cap-only sculpting because the
underlying cushion and frozen locks continue to own the side/rear silhouette;
stronger cap strokes can only dent or intersect them. A whole-head rebuild is
also rejected because it needlessly risks the approved face opening, fringe,
bow, eyes, and garment.

Verdict: test the smallest coupled silhouette-owning module, seven objects:
`Head_Cushion_Manual_Target`, `A44 continuous hair cap with smooth opening`,
three A42 rear-lock meshes, and their two seated-seam curves. Use a shared
low-frequency deformation for the head/cap/locks and rigidly carry the seam
curves. Freeze everything else.

## Reference contract

- Use all tracked relevant references, led by canonical front 25 cm,
  canonical 180-degree turn, physical front, and physical side.
- Visible crown/rear base depth target: `0.77–0.85 Wh`.
- Crown apex: `0.33–0.37 Wh` behind the face plane.
- Preserve a separate thin rear layer supplying `0.36–0.38 Wh` additional
  overhang. Do not fill the full `1.14–1.23 Wh` side envelope as solid cap.
- Controlling turntable neighborhoods: frames 10–12 and 24–26 for exact side,
  6–8 for three-quarter, and 18–20 for rear.
- Front identity, face aperture, eyes, cheek locks, fringe, bow attachment,
  garment, cameras, and lights are fixed regression context.

## Bounded implementation

1. Create an isolated deep-copy candidate from exact rung 003; never overwrite
   the parent or tracked asset.
2. Inventory and fingerprint all context. The seven allowed objects must have
   unique data before deformation.
3. Build one low-frequency cage or equivalent non-destructive macro control
   around only the coupled rear-cranium module. Pin the face-facing half-space
   (`Y <= 0`) and the bow attachment island (`Z >= 199 mm`).
4. C1: apply at most 4 mm broad world-negative-Y rear-depth correction. Carry
   the lock seams rigidly with their locks. Render exact side and both
   three-quarter views before a full packet.
5. C2 is allowed only if C1 moves the broad silhouette in the correct direction
   without dents, intersections, aperture drift, or bow-root separation. Its
   cumulative correction is capped at 6–8 mm.
6. Clean-reopen every saved candidate, recheck frozen fingerprints, contacts,
   and exact parent hashes, then run a pinned five-view packet and blind review.

## Gates and stop conditions

- Reject before rendering if the deformation touches a frozen object, changes
  the front face opening, breaks head/cap offset ordering, or detaches a rear
  lock/seam.
- Reject C1 if the broad rear-wall silhouette does not visibly move, even if
  operators or geometry checks pass.
- Reject immediately for local dimples, a bald crown, skin breakthrough,
  bow-root drift, detached locks, or a filled solid rear envelope.
- Stop after C2 if reference depth is still missed; change representation
  rather than tuning the same cage.
- A bounded success must improve side and both three-quarter silhouette and
  plush construction without lowering any fixed parent-rung category. It is
  not final model approval.

## Intended evidence and process controls

- Preserve `workstream_reference/REPORT.md` and
  `workstream_interface/REPORT.md` as preflight evidence.
- Record exact allowed/frozen object fingerprints, cage parameters, displacement
  support, contact gaps, and candidate hashes.
- Render cheap controlling views before full packets; review pixels, not only
  mesh statistics.
- Parallelize read-only rendering, numeric verification, and blind review only
  after each immutable checkpoint exists. The coordinator remains the sole
  `.blend` and goal-state writer.
- Return early when the interface or deformation support is wrong. The goal is
  the model result; process work is retained only when it accelerates or makes
  that result more reliable.
