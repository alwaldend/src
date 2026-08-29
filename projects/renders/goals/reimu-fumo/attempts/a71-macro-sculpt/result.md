# A71 result — native macro sculpt rejected at P0

## Outcome

Reset. A71 created the planned fresh, single-remesh macro volume, but both
ordinary and Elastic Grab failed the direction/support gate. No artistic
sculpt, panel construction, material work, or parent integration was attempted.

## Evidence verdict

- Ordinary Grab materially affects only a tiny local patch.
- Elastic Grab reports global floating-point changes but materially affects
  only 40 of 3,863 vertices above 0.1 mm.
- This repeats A68's support failure on new uniform topology, falsifying the
  native synthetic-stroke method rather than the old cap alone.
- The raw Icosphere render is only a method baseline and fails reference
  likeness, macro silhouette, and construction by design. Criteria 001–004
  fail for this candidate attempt; criteria 005–008 remain unverified.

## Exact identities

- Protected rung 003:
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Raw P0:
  `sha256:0004a3f0bc4987a250f7028b8697a7f740c70866ad8b60f7181e2b4eafa96400`
- Rejected ordinary Grab probe:
  `sha256:3fd572256ecc319d5e8b1255f8045fb56b3574879d5bbf965ae54adacb7368b8`
- Rejected Elastic Grab probe:
  `sha256:049ba3fa37a4beb5fef5729a52040a46a8c307f1600cd453d6b74fd3c8f08991`
- Protected tracked reusable asset:
  `sha256:489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`

## Next strategy

Use one deterministic low-frequency macro rest surface with explicit support
and multi-view landmark parameters. Render it before panels. Reserve native
sculpting for local organic breakup after a macro silhouette passes.
