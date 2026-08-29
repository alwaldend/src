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

## Choose a pixel-visible delta before each cycle

Start from the last accepted model when one exists. Otherwise explicitly
freeze the best viable current model as the comparison baseline and record
that it is not accepted. Preserve that baseline and prefer a small, reversible
edit over regenerating a component or rebuilding the model. Before editing,
record:

- one dominant pixel-visible defect, supported by a landmark or overlay;
- the smallest geometry or material owner that can change that defect;
- a lightweight causal-reach check: estimate the measured gap, the largest
  safe visible effect this owner can produce, and whether the edit can create
  the required kind of cue such as a silhouette, free edge, overlap, or
  material response;
- the controlling view and the view most likely to expose a regression;
- landmarks already within tolerance, which the edit must not disturb;
- the expected visible result and a binary keep-or-undo condition.

Edit only that owner. When the visible boundary depends on touching or nested
surfaces, move the surfaces as a coupled set when needed to preserve coverage,
contact, and depth order. For example, moving an opening without its underlying
surface must not expose a gap or cause a crossing. Do not regenerate a whole
subsystem to solve a local boundary error.

Cancel or redirect the cycle before authoring when the bounded edit cannot
produce a meaningful fraction of the measured correction, or when it changes
the wrong visual category—for example, a shallow dent cannot create a missing
independent fabric leaf. Keep this preflight proportional: use existing
measurements, a quick projection, or a disposable probe, not a new analysis
project. If new evidence invalidates causal reach while work is in progress,
stop at that evidence boundary and preserve only useful diagnostics.

## Run a complete fidelity cycle

A cycle is complete only after all of these steps:

1. Save the exact accepted checkpoint, or the explicitly frozen unaccepted
   baseline when no accepted model exists, and its fixed-camera renders.
2. Diagnose one dominant failure from the references and scorecard.
3. Verify that the owner and safe effect size can reach the named pixels, then
   apply the smallest reversible edit to its owner, including paired
   deformation of dependent surfaces when required.
4. Make fast renders at 512--640 px only from the controlling view and the view
   most likely to expose a regression.
5. Align the candidate to the reference using the recorded landmarks. Produce
   a side-by-side comparison and a 30--50% silhouette overlay or edge
   difference. Do not judge an unaligned contact sheet by impression alone.
6. Run an implementation-blind absolute review of the candidate against the
   references before revealing the modeling method, measurements, topology,
   intended fix, or previous model. Then make a fixed-camera baseline-versus-
   candidate A/B comparison for the keep-or-undo decision. Keep the edit only
   when its named defect visibly improves and protected landmarks and the
   regression-risk view do not regress. Numeric or topological change without
   a visible improvement is churn: undo it.
7. Render the full front, side, rear, and three-quarter regression set and one
   uncropped presentation render only when promoting a fast-pass candidate to
   a checkpoint or approval stage.
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

- Do not replace a representation after one failed edit. Rebuild or introduce
  a new representation only after multiple reviewed, isolated edits show that
  the current owner cannot correct the target pixels without violating a
  protected landmark, contact, or silhouette in another view. Record that
  evidence and start the replacement as a bounded, reversible comparison
  against the frozen baseline.
- If repeated edits fail for different reasons, improve the diagnosis or split
  the owner more narrowly before concluding that topology is the limit.
- If three or more identity-defining components fail together, first determine
  whether the frozen baseline is still viable. Return to a macro
  blockout only when direct evidence makes that baseline categorically
  nonviable or reviewed local edits establish that the failures share a
  structural owner.
- If a front-view fix detaches, intersects, or distorts a part in another view,
  undo it. Use a surface-relative local edit next when that directly addresses
  the cross-view failure; replace the placement representation only after
  repeated controlled evidence establishes that the current one is the limit.
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
