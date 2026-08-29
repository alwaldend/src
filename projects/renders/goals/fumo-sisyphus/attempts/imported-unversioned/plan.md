# Fumo Sisyphus scene scaffold goal

## Goal

Build a reference-faithful, reusable Blender scaffold for the original
Sisyphus composition: believable slope and boulder, pale sepia atmosphere,
strong readable shadows, fixed cameras, and the original image available at
the side for comparison. Leave the character slot as an explicitly neutral
placeholder until the standalone Reimu Fumo is approved.

## Status

`in progress` — two Attempt 01 variants were rendered and rejected. V01 keeps
rock relief but loses the pale sky; V02 adds a faceted right wall and also
misses the source hierarchy.

## Current state

- Protected tracked source:
  `blender/fumo/fumo_sisyphus/fumo_sisyphus.blend`, SHA-256
  `c5bd58ed9b29a6d67c398136eaec7ed34e227934c464662dfcb61f61f8e6f591`.
- The prior render reverses the source hierarchy: a near-black foreground
  wedge dominates while the boulder is too small and the slope contour is too
  smooth and graphic.
- Character fidelity is outside this subordinate goal. The candidate may
  contain only `PLACEHOLDER_FUTURE_APPROVED_FUMO`, which is not a Fumo model,
  pose, or integration claim.
- Rejected candidates and review frames are preserved under
  `out/fumo_sisyphus_scene_scaffold/`; neither is authorized for promotion.

## Current plan

1. **Completed:** start from a candidate copy under
   `out/fumo_sisyphus_scene_scaffold/`; never save over either tracked blend.
2. **Completed twice, rejected:** reset only slope, boulder, background, fixed camera, lighting, and the
   viewport reference-at-side.
3. **Completed twice, rejected:** render a quote-free review frame, compare it beside the controlling image,
   and reject if the macro landmarks or light/shadow hierarchy miss the gates.
4. Preserve both failures. The next attempt must solve pale negative sky and
   the diagonal slope/right-face hierarchy before adding detail; promotion is
   not authorized.

## Acceptance criteria for this scaffold

- Boulder silhouette occupies `0.40–0.55` of frame width, sits in the upper
  left/middle, and reads as a rounded weathered mass rather than a faceted or
  melted blob.
- Primary slope rises from lower-left to upper-right and occupies roughly the
  right half of the frame without becoming an opaque black wedge.
- Pale amber/cream sky remains visible on the left and supplies separation;
  directional light produces a readable boulder cast shadow and rock relief.
- Fixed review camera is quote-free. The original image is packed in a
  render-disabled `REFERENCE_ONLY` collection and displayed at the side in
  the viewport.
- The neutral placeholder is visibly labeled by its object name and remains
  replaceable; no Reimu asset is linked, copied, edited, or evaluated.
- Both protected tracked blend hashes are unchanged after the attempt.

## Records

- [References and landmark contract](references.md)
- [Current attempt and evaluation](current_attempt.md)
- [Artifact log](artifacts.md)
- [Evidence manifest](evidence.md)
- [Failure ledger](failures.md)
