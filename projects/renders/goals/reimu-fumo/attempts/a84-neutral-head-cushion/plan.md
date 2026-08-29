# A84 plan — coarse modular neutral sculpt, head-cushion owner

## Decision and hypothesis

Build a new coarse modular neutral sculpt from an empty task-owned scene.
Keep exact C1b only as a hidden comparator; do not copy, append, retopologize,
or sculpt its visible-form meshes.

This attempt is bounded to the first owner, `C19_HeadCushion`. The hypothesis
is that one directly sculpted, shallow cushion can match the canonical front
field and turntable depth while preserving readable front and rear planes,
side-gusset volume, and soft stuffing asymmetry in every required view.

The owner has causal reach because it directly controls head width, front and
rear depth, crown and underside transitions, face field, collar seat, and the
future receiver for hair and bow. Hair cannot rescue a rejected receiver.

## Immutable inputs

- C1b comparator:
  `out/reimu_fumo_attempt_083_incremental_sculpt/live_author/`
  `a83_C1b_coupled_cap_receiver_narrow.blend`, SHA-256
  `d2357588b42b18285f31fcf780f2be5e76111a002a25b9ac25cd569be6cbf8d1`.
- Canonical front:
  `projects/renders/blender/fumo/reimu_fumo/references/`
  `canonical_front_25cm.png`, SHA-256
  `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`.
- Canonical depth and layer order: `canonical_turn_180.gif`, SHA-256
  `0d774eaa7f75828e388df4fb886cda7c563ce3bcd4ccb38d9885997a0846af30`.
- `physical_front.png`, `physical_side.png`, and `sofa.gif` may veto
  implausible cushion construction but may not override the canonical
  silhouette.
- A83 close evidence: `BASELINE_VERDICT.md` and `DIRECT_SCULPT_LOOP.md`.

## Scene and interface contract

Create an empty task-owned Blender scene at the measured 25 cm complete-plush
scale. Put reference backgrounds, comparator render cards, fixed cameras, and
lights in `C19_REVIEW_ONLY`, outside the reusable asset collection. Register
fixed front, rear, both profiles, and both front three-quarter cameras by head
width, head center, crown, eye-line proxy, and future ground plane. Reject an
overlay when registration residual exceeds `.015 Wh`; correct the camera,
not the candidate.

Create only one transformable artistic owner, `C19_HeadCushion`, from a
low-resolution cube with applied scale, minimal applied bevel, and a Multires
modifier. Assign neutral beige clay. Do not add hair, bow, face embroidery,
body, seams, fibers, surface noise, rig, voxel remesh, Dyntopo, or generated
artistic geometry in this attempt.

Preserve a broad front plane, a readable rear plane, shallow depth, side
gusset volume, a crown receiver, an underside/collar seat, and slight stuffing
asymmetry. The object must remain separate and transformable for later rigging.

## Reversible sculpt loop

1. Preflight source hashes, 25 cm scale, camera registration, active object,
   and Multires level.
2. Name one visible discrepancy, one controlling view, and one regression
   view. Make at most three broad sculpt strokes and one light smoothing
   stroke using direct Grab, Elastic Deform, Flatten, and Smooth tools.
3. Inspect front, rear, both profiles, and both three-quarters. Undo
   immediately on a categorical sphere, egg, box, helmet, twist, pinched edge,
   or depth failure.
4. Save Copy to one immutable `snapshots/C19_head_sNN.blend`, hash it, and
   never render the mutable live save target.
5. Clean-reopen the snapshot in repository-pinned Blender 5.2.1 with file
   scripts disabled. Batch-render the controlling camera and highest-risk
   view at `512 x 512` in one warm process.
6. Produce registered reference overlays and an unlabeled C1b/candidate A/B.
   Review the pixels against references before revealing implementation.
7. Keep only a visible reference-directed improvement with no protected-view
   regression. Otherwise restore the preceding snapshot. Do not accumulate a
   failed delta.
8. A kept controlling pair earns one six-direction neutral-clay and owner-ID
   packet from the same unchanged geometry state.

Target at most 90 seconds from immutable snapshot publication to an inspected
controlling pair and at most two minutes per stroke batch. If the render pair
misses 90 seconds, repair the existing snapshot/render harness before another
stroke. Do not render the same snapshot in both live and pinned hosts.

## Review packet and stop gate

The first decisive artifact contains:

- neutral-clay and owner-ID front, rear, both profiles, and both
  three-quarters from one immutable snapshot;
- registered silhouette overlays aligned by head width, crown, eye line, and
  head center;
- one uncropped normal-scale contact sheet;
- source, candidate, tool-version, and render hashes; and
- an implementation-blind absolute review.

Keep the head owner only when all critical landmarks are within `3% Wh`, all
major silhouette and depth extrema are within `5% Wh`, registration is valid,
front and rear remain visibly flattened, and every profile and three-quarter
reads as a stuffed gusseted cushion rather than a smooth primitive. It must
also provide plausible future collar and crown attachment surfaces.

Close early as reset on any categorical sphere, egg, box, helmet, twist,
pinched-edge, or plastic read. If two bounded sculpt cycles fail for different
views, change the base proportions or coarse topology in a new attempt rather
than adding more strokes. Do not begin hair, bow, or body work until this
owner's full packet passes.

## Criteria and fixed regressions

- Criterion 001 is affected only for the head-cushion proportion and
  silhouette subset; the complete criterion remains unverified.
- Criterion 002 is affected by the registered head landmarks and depth bands;
  the complete criterion remains unverified until the full model exists.
- Criterion 003 is affected by the cushion-construction read; the complete
  criterion remains unverified.
- Criterion 005 is affected by owner naming, separation, scale, transforms,
  and review/reusable collection boundaries; complete reusable structure
  remains unverified.
- Criteria 004 and 006--008 remain unverified in this bounded attempt.

Fixed regressions are exact source hashes, valid camera registration,
25 cm complete-plush scale, one unchanged candidate across all packet views,
no C1b visible-form import, a separate transformable head owner, and no
hair/body/detail used to conceal a failed receiver.

## Parallelism and next module

The canonical Blender state has one artistic owner and no stable mergeable
geometry interface, so authoring remains with one coordinator. Reference-frame
preparation, manifest/hash verification, and implementation-blind image review
may run independently against immutable inputs or a frozen snapshot. They may
not mutate canonical goal state or the live blend.

If and only if this head owner passes, freeze it and open a new bounded attempt
for the compact crown/rear-hair scaffold. The seated torso/pelvis and garment
module remains later work in the same whole-plush scale and camera contract.
