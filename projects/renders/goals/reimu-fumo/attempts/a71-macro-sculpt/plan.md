# A71 — coherent macro head/hair sculpt

## Immutable bindings

- Goal: `reimu-fumo`
- Goal resource version / checkpoint CAS token: `7`
- Goal generation: `1`
- Lifecycle generation: `1`
- Criteria revision: `1`
- Exact parent context:
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Protected reusable asset:
  `sha256:489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`

## Decision review

The strongest objection is that sculpt-first work already produced a smooth,
fat, reference-poor head. A domain expert would also reject a freeform sculpt
without measurable panels or reusable topology. Those failures came from
sculpting a cap-only copy on top of the old boxy receiver and judging brush
motion rather than a new multi-view macro silhouette. A70 establishes the
opposite failure: premature panel construction creates masks and cards before
a coherent rest surface exists.

Verdict: use sculpting only to establish one new, coherent clay macro surface
from scratch. It is a disposable form-finding stage, not deliverable topology.
Freeze all identity context, render every controlling view after each bounded
checkpoint, and reject before retopology if the silhouette is wrong. Only a
passing macro surface may become the rest surface from which constructed
fabric panels are derived.

## Bounded plan

1. Append exact rung 003 into an empty disposable live Blender session. Hide
   the old receiver/cap but preserve them as rollback context; freeze all other
   objects and cameras.
2. Create one low-resolution, watertight head/hair macro volume from scratch,
   positioned against the existing face, fringe, bow, locks, and body. Do not
   add seams, thickness, fibers, materials, or rigging.
3. Use native Blender sculpt operators for broad silhouette only: Grab/Elastic
   Deform for width and depth distribution, Clay Strips/Inflate for stuffing,
   and Smooth sparingly. Record every stroke or deterministic macro operation.
4. S0 blockout: front, both profiles, rear, and both three-quarter views. Veto
   a sphere, egg, cube, mask, long flat crown, vertical rear wall, exposed
   support, or loss of the canonical face aperture.
5. S1 allows one bounded macro correction after pixel review. S2 is allowed
   only if S1 materially improves the controlling silhouettes without a new
   major defect. Stop after S2 and change strategy if still below the gate.
6. If and only if clay macro silhouette passes, preserve it as a geometry rung
   and begin a later attempt that derives receiver, crown/temple skin, and rear
   leaves from that shared surface.

## Macro acceptance gate

- Whole-subject front, rear, both profiles, and both three-quarter views are
  visibly closer to all controlling references than rung 003.
- Broad rounded head width, compact face aperture, mid-height depth maximum,
  tapered crown/underside, and independent rear reserve are readable without
  material or seam cues.
- No old helmet surface is visible or used to fill accessory contacts.
- Front identity and frozen context do not drift.
- An implementation-blind review scores macro silhouette at least 6/10,
  reference likeness at least 5/10, reports no new major defect, and agrees the
  form is viable for panel derivation. This is not final model approval.

## Execution controls

- One coordinator owns the live `.blend`; independent agents may design the
  workflow and review immutable packets only.
- Blender 5.1.1 Flatpak is a disposable live sculpt host. Every checkpoint is
  saved under `out/`, reopened and rendered in pinned Blender 5.2.1, and never
  overwrites the parent or tracked asset.
- Render the first coherent macro volume before adding any secondary part.
- Return early if the volume/interface/brush support is wrong. Record the
  artifact, exact hashes, visual verdict, and whole-process audit.
