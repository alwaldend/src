---
name: blender-reference-fidelity
description: >-
  Build or revise Blender models against supplied image references using
  measured landmarks, calibrated multi-view comparisons, implementation-blind
  image review, regression gates, and modeling-strategy resets. Use when
  likeness or reference fidelity is the acceptance criterion, especially after
  repeated visual iterations have stalled; do not use for freeform designs
  without controlling references.
---

# Blender Reference Fidelity

Treat fidelity to the supplied references as the outcome. A technically clean,
plausible, or attractive model does not pass when its proportions,
construction, or silhouette differ visibly from the references.

## Establish evidence before editing

1. Preserve the user's scope, known-good files, and approval stages.
2. Build a reference dossier with the best available front, side, rear, and
   three-quarter frames. State which reference controls each dimension; do not
   average incompatible poses or camera perspectives.
   When sources conflict, prefer a high-resolution physical photo for
   construction and volume, the closest turntable view for otherwise hidden
   silhouettes, and a clean product image only for graphic landmarks. Name one
   primary source per target and record variant differences rather than
   combining incompatible extremes.
3. Record normalized landmarks before changing geometry. Use a stable unit such
   as head width and include overall bounds, component bounds, centers, angles,
   overlaps, and visible gaps. Record uncertainty for perspective or occlusion.
4. Create fixed review cameras that correspond to the canonical references.
   Record projection type, lens or orthographic scale, object rotation, crop,
   and alignment landmarks. Do not use an overlay as evidence until the camera
   match is within the target's tolerance. Keep framing, projection,
   resolution, and lighting unchanged between candidates.
5. Read [references/scorecard.md](references/scorecard.md) and keep a scorecard
   in the task's temporary-output directory. For plush or constructed-fabric
   subjects, also read
   [references/plush-construction.md](references/plush-construction.md).
6. Read [references/visual-quality-gate.md](references/visual-quality-gate.md)
   before judging or presenting any candidate. Its absolute image-quality gate
   overrides relative improvement, measurements, and technical cleanliness.

## Evaluate the representation before each cycle

Write down the following before editing:

- the largest reference mismatch, supported by a landmark or overlay;
- the current modeling representation for that form;
- whether that representation can produce the referenced construction;
- one bounded change and its expected effect in every review view;
- an objective acceptance condition and the likely regression risk.

Replace the representation when it is structurally wrong. Do not keep resizing
an ellipsoid when the reference shows a gusseted cushion, smooth a cone when it
shows sewn panels, or add surface detail to hide a silhouette error.

## Run a complete fidelity cycle

A cycle is complete only after all of these steps:

1. Save the last accepted checkpoint and its fixed-camera renders.
2. Diagnose one dominant failure from the references and scorecard.
3. Apply one coherent hypothesis at one detail tier.
4. Make fast renders at 512--640 px from the controlling view and the view most
   likely to expose a regression.
5. Align the candidate to the reference using the recorded landmarks. Produce
   a side-by-side comparison and a 30--50% silhouette overlay or edge
   difference. Do not judge an unaligned contact sheet by impression alone.
6. Render the full front, side, rear, and three-quarter regression set and one
   uncropped presentation render when the fast comparison passes.
7. Run the implementation-blind review from the visual-quality gate. Review
   against the references before revealing the prior candidate, modeling
   method, measurements, topology, or intended fix.
8. Update every affected scorecard row and explicitly accept or reject the
   candidate. Revert rejected geometry rather than building on it.

An edit batch without comparison and an accept/reject decision is not a cycle.

## Gate changes from large to small

Advance in this order:

1. macro silhouette and component proportions;
2. manufactured construction, planar profile, attachment, and overlap;
3. secondary stuffing compression, folds, seams, and controlled asymmetry;
4. materials, fibers, stitches, and microdetail.

Do not advance while an earlier tier fails. Unless the references justify a
wider tolerance, use these initial limits:

- critical canonical-view landmarks: at most 3% of head width from target;
- silhouette extrema and major gaps: at most 5% of head width from target;
- uncertain perspective views: at most the recorded uncertainty, never an
  unbounded visual judgment;
- no new critical regression in any other view.

Technical checks such as manifold geometry, subdivision, parenting, or polygon
count are required hygiene, not evidence of likeness.

Do not pass a candidate because it is better than the previous one. A major
visible mismatch in an identity-defining form is an automatic failure even
when averages, overlays, or technical checks pass.

## Reset a stalled approach

- If the same high-level complaint survives two reviewed cycles, stop tuning
  that subsystem and rebuild it from the reference construction or restore the
  last accepted checkpoint.
- If three or more identity-defining components fail together, return to a
  macro blockout instead of patching them independently.
- If a front-view fix detaches, intersects, or distorts a part in another view,
  replace the placement method with surface-relative construction.
- Do not declare success from an internal critique, agent vote, or cleaner
  topology when the overlays still fail.
- Do not use the user as the first visual-quality reviewer. Inspect the final
  pixels and reject visible failures before presenting or exporting them.

## Keep review and deliverables separate

- Store extracted frames, measurements, overlays, candidate renders, and
  working files under the task's temporary-output directory unless the user
  requests otherwise.
- Reference planes may exist in a hidden review-only collection in the working
  file. Strip them, review cameras, lights, and temporary data from a reusable
  asset export.
- Preserve explicit approval gates. Sculpt approval does not authorize
  materials; material approval does not authorize posing or integration.
- Keep diagnostic renders separate from the full presentation render. Never
  substitute viewport geometry, a contact sheet, or topology output for the
  image the user is being asked to approve.
- Report the controlling references, measured changes, remaining deviations,
  and the evidence used to pass a gate.
