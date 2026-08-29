# A70 C0 self-review

## Verdict

Reject the constructed receiver/cap coupon at C0. Do not continue to C1,
softness, materials, attachment transfer, or detail.

## Visible evidence

- V1 exposes a nearly square beige receiver because the crown/rear yokes do
  not own the frontal hair silhouette.
- V2 adds planar front panels, but they become a literal rectangular mask in
  front, side, and both three-quarter views.
- V3 rounds the outer front boundary and improves the front outline, but the
  aperture creates two tall beige horn/ear shapes and the side remains a beige
  egg with only a thin brown arc.
- V4 adds broad side coverage, but the side panels become large rigid cards.
  They recreate the filled helmet envelope and produce abrupt panel steps in
  both three-quarter views.
- The rear remains visually dominated by the old locks and bow, while the new
  assembly supplies no convincing soft seam transition into them.

The early vetoes for a block/mask/card read, exposed support outside the face
opening, abrupt cap-to-lock steps, and whole-subject regression all fire.

## Technical evidence

- The coarse closed receiver is 482 vertices / 480 polygons and measures
  approximately 121.1 x 97.0 x 115.0 mm.
- Each closed yoke panel is 234 vertices / 160 polygons and approximately
  3.8 mm thick by construction.
- All four immutable live snapshots reopen and render in pinned Blender 5.2.1.
- The source parent remains byte-exact at
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.
- The rejected V4 is
  `sha256:616276f7dbf0530c8c35a0f80f8a14d3672b9e2d8ec8914277a7f5d11a3d09aa`.

## Strategy finding

The receiver and hair cannot be authored as several independently extruded
silhouette polygons and expected to become a soft plush through bevels. The
panels lack a shared curved rest surface, so every attempt to add coverage
creates a card or bridge. This is a representation failure, not a parameter
failure.

The next strategy should establish one coherent sculpted macro volume against
the fixed multi-view silhouettes, then derive constructed panel boundaries and
thickness from that approved surface. It must not reuse the old helmet shell,
and it must not detail or retopologize until clay-only front, side, rear, and
both three-quarter views pass.

## Process audit

The live Flatpak loop reduced each geometry edit to a few seconds and the
pinned five-view render to roughly 33 seconds, which is a substantial feedback
improvement. The waste was conceptual: three panel refinements were attempted
after the first render showed that independent planar patches lacked a shared
3D rest surface. Future C0 work should render a single silhouette-owning macro
surface before adding any secondary panel. Independent read-only contracts and
reviews remain useful; multiple concurrent `.blend` writers would not.
