# A72 — deterministic low-frequency receiver rest surface

## Immutable bindings

- Goal: `reimu-fumo`
- Goal resource version / checkpoint CAS token: `9`
- Goal generation: `1`
- Lifecycle generation: `1`
- Criteria revision: `1`
- Exact parent context:
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Protected reusable asset:
  `sha256:489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`

## Decision review

The strongest objection is that another procedural surface may repeat A70's
rigid, low-poly result. That would happen if silhouette polygons or separate
panels are extruded before a shared curved surface exists. A72 instead owns
only the hidden stuffed receiver: one smooth low-frequency watertight surface
with explicit front plane, width, crown taper, mid-height rear maximum, and
underside turn-in. It does not claim hair, panels, seams, or final topology.

Verdict: use deterministic control because A71 proved the synthetic native
brush interface cannot move a broad region reliably. Explicit formulas make
support and landmarks auditable. Reject at the first clay packet if the form
reads as an egg, cube, mattress, mask, or human head.

## Bounded plan

1. Build one closed receiver from analytic low-frequency cross-sections in an
   isolated scene; no old head/cap geometry is reused.
2. Keep `Wh = 132 mm`, a broad near-planar face, rounded-square frontal
   envelope, mid-height depth maximum, tapered crown, and lower-rear undercut.
3. Use only 17 height rings by 32 section samples plus poles. Preserve the
   explicit parameters and vertex indices; subdivision is review-only.
4. P0: render front, both profiles, rear, and both three-quarter views as
   neutral clay. Compare against every controlling reference and the exact
   numeric bands.
5. Allow at most one parameter correction if P0 is directionally viable.
   Stop if a second packet repeats the same macro category failure.
6. A passing receiver becomes a shared rest-surface input for a later thin
   hair-skin/rear-leaf attempt. It is not promoted as a complete head or model.

## Acceptance and vetoes

- Width, height, supporting aperture, base depth, crown apex, and undercut are
  inside the frozen A70/A71 bands.
- Front reads as a broad stuffed Fumo cushion, not a sphere or human skull.
- Both profiles show a shallow face plane, mid-height depth maximum, and crown
  plus underside turn-in without a long plateau or vertical wall.
- Rear and three-quarter views have continuous curvature with no box corner,
  pole pinch, slab, shield, or asymmetry trick.
- An independent reviewer scores the receiver macro silhouette at least 6/10
  and finds it viable as a panel-rest surface. This bounded gate does not claim
  hair, identity, plush construction, or final approval.

## Execution controls

- One coordinator owns the live `.blend`; independent agents may inspect
  immutable packets only.
- Flatpak Blender 5.1.1 is the disposable live host; pinned Blender 5.2.1
  reopens and renders every checkpoint.
- Render the receiver alone first, then optionally with frozen whole-subject
  context. Do not add hair or accessory bridges to conceal a bad receiver.
- Record exact parameters, hashes, pixel verdict, and process audit.
