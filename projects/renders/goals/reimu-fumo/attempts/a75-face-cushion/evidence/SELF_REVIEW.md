# A75 self-review — one-object face module rejected

## Verdict

Reset A75. The one-object boundary is mechanically valid only for the exposed
beige aperture, but its quick salvage probe is worse than rung 003 and a fresh
purpose-built retopology cannot address the dominant frozen hair-shell defect.
No protected or tracked Blender asset changed.

## Exact tested subject

- Candidate blend:
  `sha256:8e10d4cf11b78b79b3d2f013f606a6086bcbdf19bf9d5fe6994da39ced072be5`
- Protected parent:
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Salvage source:
  `sha256:2f0ee05ee60b6b66b3015d524a93c26424e303e4ead48e0a843f9b18408f3751`
- Render manifest:
  `sha256:ffcf79621c68c9c36d089f61daf53d95ffe52772fecd23ef4ff1904355bf086d`
- A/B sheet:
  `sha256:3cbcecf7255fac1b25808b86b1c1062b41650b995117bfa585ba8f83c8b68b47`

Pinned Blender 5.2.1 built and reopened the candidate. Exactly
`Head_Cushion_Manual_Target` changed source visibility. Every other source
object was fingerprint-identical. The candidate added one single-user cushion
mesh and one single-user beige material.

## Reference and interface findings

- The canonical face exposure is an observed `.6033 x .6033 Wh` pixel field;
  hair owns its top and side boundaries. The directly observed lower contour
  is a broad shallow cheek/chin arc, not a full cushion outline.
- The turn sequence shows beige exposure collapsing to a small front crescent;
  it does not expose a complete beige profile or rear.
- The frozen cap is only about `0.54-1.02 mm` outside the old receiver and the
  eyes are only `0.076/0.080 mm` clear. A one-object replacement can change the
  visible face field, but cannot honestly change the macro head, profile, rear,
  or hair seating.
- Side and rear views are therefore leakage/regression tests only.

## C0 result

The reused A66 cushion is numerically close to the frontal exposure width but
visually wrong. It creates a small spherical face ball, a dark cheek-to-chin
cavity ring, bulbous three-quarter exposure, and detached beige islands in
component masks. Implementation-blind A/B review prefers rung 003 in front
and both three-quarter views and scores the affected region identity `4/10`,
construction `3/10`, and contact `3/10`.

The failure is the visible front field itself, not a hidden depth course. Do
not spend the allowed correction by shifting it forward or widening it.

## Process and next decision

A75 was deliberately small, produced a decisive A/B packet quickly, and
preserved the whole model as context. That decomposition worked. Its outcome
also shows the limit of module size: a technically isolated part is not useful
when the unchanged hair shell remains the dominant failure and tightly fixes
the interface.

A purpose-built 226-vertex all-quad carrier is documented for future use, but
building it now would optimize a secondary region without changing the
whole-plush read. The automatic color visual-hull alternative also fails
preflight because occlusion fragments and contaminates component masks.

The next attempt changes the interaction method rather than another analytic
shape family. Test a broad native sculpt stroke through the live Blender MCP
in a real `VIEW_3D` context on a disposable copy. A71 tested background
synthetic strokes and moved too few vertices; the live-context probe can
directly confirm or falsify whether the interface, rather than sculpting
itself, was the blocker.

