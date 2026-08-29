# A69 result — coupled deformation falsified

## Outcome

Reset. A69 proved that a coupled low-frequency head, cap, and rear-lock
deformation can move the correct rear silhouette while preserving the front,
but the best 4 mm checkpoint fails the promotion attachment gate and leaves
the dominant monolithic helmet construction intact. No candidate is promoted,
and no C2 amplitude increase is allowed.

## Work performed

- Audited the exact rung-003 head/cap/lock interface and established that the
  smallest silhouette-owning module is the head cushion, cap, three rear-lock
  meshes, and two seam curves.
- Rejected an invalid initial lattice design before use: Curve objects do not
  support vertex groups, and the proposed lattice bounds were wrong.
- Built an exact C0 copy with five analytic mesh shape keys and a correctly
  bounded lattice for the two seam curves. Pinned Blender 5.2.1 clean reopen
  found all seven copies identical to their source and all 177 frozen objects
  unchanged.
- Used the explicitly requested Blender 5.1.1 Flatpak as a disposable live
  MCP host. The verified 5.2.1 C0 was appended into an empty scene rather than
  opened or overwritten across versions.
- Tested root pinning, cap bridges, shared radial pins, bow-distance
  attenuation, and amplitude sweeps. The retained hybrid applied the intended
  4 mm rear correction without a visible bald patch or front regression.
- Saved an immutable candidate, reopened and measured it with pinned Blender
  5.2.1, rendered it with the repository packet tool, and obtained an
  implementation-independent diagnostic review.

## Evidence verdict

- The side and both three-quarter views move modestly in the intended
  direction; the front remains stable.
- The main rear-lock attachment band retains only 72.44% of baseline contact;
  the other two retain 82.58% and 78.12%. The required 75% promotion gate
  therefore fails.
- The reviewer scored overall likeness 4/10, macro silhouette 3/10,
  construction 3/10, contact 3/10, and plush-medium read 3/10.
- The head still reads as one deep, hard helmet with a long crown and vertical
  rear wall, not a shallow stuffed receiver plus thin layered hair.
- Criteria 001–004 fail for this candidate attempt. Criteria 005–008 remain
  unverified and unchanged.

## Exact identities

- Protected rung 003:
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Verified isolated C0:
  `sha256:fd48da9a0a6b2f919477da0809a6a6fee987513d2a33841a084d8bad2b528b6d`
- Rejected C1 v4:
  `sha256:12f2a09d376d754c838fc2dddb9c528e59709e60119882bbf5563572903bcfe9`
- Protected tracked reusable asset:
  `sha256:489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`

## Dominant failure and next decision

The failure is representational, not an amplitude problem. Repeated sculpt and
deformation attempts preserve the current dense, boxy receiver and can only
trade a small silhouette change against dents or detached locks. The next
attempt must restart from rung 003 and construct a shallow, soft head receiver
and separate layered hair cap as a small isolated module, preserving the face
aperture and surrounding parts as frozen context. Exact side and both
three-quarter views remain early veto gates.

## Whole-process audit

The useful loop is now faster: one live Flatpak mutation owner supports
seconds-scale parameter probes, while immutable snapshots are reopened,
measured, and rendered with pinned Blender 5.2.1. Independent review begins
only after a saved checkpoint exists. The wasted work came from trying several
contact heuristics after the current shell representation had already shown
that it could not satisfy both depth and attachment. The next attempt will
build and veto a coarse receiver/cap coupon before transferring locks, seams,
materials, or details.
