## Context

See proposal.md. The source contains 898 frames at 30000/1001 fps and 1280x720.
The host H.264 decoder is a stub; pinned Blender successfully reads the clip.
The task has no exposed GPU device.

## Goals / Non-Goals

**Goals:** Editable geometry, character motion and shot cuts, self-contained
scene, original soundtrack, practical background rendering.

**Non-Goals:** Recovering the original production meshes or exact facial rigs.

## Decisions

Use procedural geometry and keyframes in pinned Blender 5.2.1, with no external
asset dependency. Build faces, helmets, uniform construction and base props
from editable meshes. Bake motion as ordinary object animation,
so reopening does not depend on Python handlers. Use CPU Cycles with eight samples, denoising, and persistent render data.
The measured Eevee preview cost exceeded the Cycles preview cost on this host. Preserve the complete source audio
by stream copying into the final H.264 MP4. Keep source media outside the asset;
the generator checks its SHA-256 identity.

## Risks / Trade-offs

- Original character likeness requires visual comparison; procedural topology
  and facial timing may remain approximate. Record deviations explicitly.
- CPU graphics may be slow. Measure preview renders before rendering the full
  sequence and retain resumable frame output.
- H.264 host stubs require Blender's bundled encoder and decoder for verification.

## Reference-fidelity revision

The procedural draft's head, clothing, and poses failed independent review.
Retuning its isolated primitives is unlikely to correct the underlying
construction. Build a continuous head and face with a separately modeled,
deformable mouth and eye region, using the original footage to set landmarks,
expressions, camera framing, and timing. Rebuild the garment and limb joins
as continuous forms. Compare fixed reference views and reject visible defects
while developing the animation. Preserve the published draft and source
attachment while candidates are evaluated under ignored task output. The user
authorized delivery without the fidelity skill's numeric acceptance gate;
the review findings and remaining differences must remain visible.

Keeping the existing draft would retain known major likeness defects. A flat
video projection would reproduce pixels but would not provide an editable 3D
remake. A reconstructed and animated model is the selected approach. Hidden
rear anatomy remains inferred from the available footage.

The revised Blender file is the maintained source for the revised animation.
The original procedural scene command remains available as the initial draft
generator; it does not reproduce the later manual scene revisions. The revised
file contains its geometry, materials, keyframes, and packed audio and can be
rendered with the repository-pinned Blender.

The current visual revision rebuilds the shirts as continuous surfaces driven
by native armatures, and gives each leg a knee and independently oriented foot.
Camera assignments and framing are checked against each source cut, including
the officer's speaking close-ups and the diagonal group views. Foreground
blocking uses the reference foot positions; the distant set is composed for
each shot. A camera-aligned sky gradient preserves the source's illustrated
background treatment across focal lengths. All deformation and shot changes
are baked into ordinary Blender data. Still-frame comparisons precede short
motion previews and a complete replacement render.

The leg bake calculates each shin against its explicitly computed thigh pose,
using the native segment lengths. This avoids dependency-graph evaluation order
changing a joint's position. Full-timeline and nonsequential-frame checks measure
knee and ankle continuity before the rendered walking sequence is encoded.
The opening uses planted world-space footsteps and continuous leg reach. Shin
scale inheritance is disabled, and Blender computes each bone basis from the
explicit desired parent pose, avoiding shear from a stretched thigh. The
opening camera push, body turn, and hand pose are native keyframes.

The wide dialogue shot uses shot-specific head turns, an officer stance with
hands behind his back, and connected leg reach for the main recruit. The main
uniform has shorter jacket and sleeve hems, fuller trousers, and a curved
collar. Ordinary and salute hands follow the revised cuffs. Both the opening
and wide shot have separate short motion previews; the full render remains
stopped while character likeness and performance are revised.

The mouth uses the modeled cavity's interior surfaces, with native shape keys
closing its aperture. This removes the blue interior wall that appeared as a
raised lip. The main recruit's jacket has a rebuilt continuous surface with
separate chest and arm weights and an animated collar allowance during the
salute. The officer's crown, planted stance, and jacket details were revised;
surface projection seats buttons and ribbons against the jacket. A separate
character preview checks dialogue and the salute without restarting the full
render. Remaining likeness and performance differences are recorded in
verification.md.

## Project ownership

The campaign project owns the army spot assets, commands, and change record.
Its landing page and project inventory entry identify the separate owner;
other reusable render assets remain in the renders project.

The project groups each scene under `assets/<scene>/` and keeps individual
scene generators separate. The frame encoder accepts explicit timing and size
so it can serve additional scenes without changing the army spot.
