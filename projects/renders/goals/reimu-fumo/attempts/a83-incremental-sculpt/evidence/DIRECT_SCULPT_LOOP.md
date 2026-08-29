# C19 direct multi-view sculpt loop

Status: workflow design only. No Blender process or `.blend` mutation was
performed while writing it.

## Reconciled verdict

Follow `BASELINE_VERDICT.md`: build a new coarse modular neutral sculpt from
an empty task-owned scene. Keep exact C1b only as a hidden comparator. Do not
copy, append, retopologize, or sculpt any C1b visible-form mesh.

The **first sculpt owner is the shallow beige head cushion**, named
`C19_HeadCushion`. It precedes hair and body because it establishes head width,
front/rear planes, depth, crown support, face field, collar interface, and the
receiving surface for hair and bow. Starting with hair would tune another
shell around an unproven receiver; starting with the body would leave the
primary identity scale unresolved.

This replaces the earlier idea of sculpting the C18 crown owner. C18 remains
useful evidence that owner separation matters, but its dome and abrupt side
boundary are not a viable starting form.

## Immutable inputs

- C1b comparator:
  `out/reimu_fumo_attempt_083_incremental_sculpt/live_author/`
  `a83_C1b_coupled_cap_receiver_narrow.blend`, SHA-256
  `d2357588b42b18285f31fcf780f2be5e76111a002a25b9ac25cd569be6cbf8d1`.
- Canonical front:
  `projects/renders/blender/fumo/reimu_fumo/references/`
  `canonical_front_25cm.png`, SHA-256
  `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`.
- Canonical depth and layer order:
  `canonical_turn_180.gif`, SHA-256
  `0d774eaa7f75828e388df4fb886cda7c563ce3bcd4ccb38d9885997a0846af30`.
- Construction checks only: `physical_front.png`, `physical_side.png`, and
  `sofa.gif` from the same reference directory.

Copy C1b byte-for-byte to `C19_process_reset/comparator/C1b_exact.blend`, hash
both files, and never save either. Extract or reuse canonical-turn front,
profile, rear, and three-quarter frames under `C19_process_reset/references/`.

## Empty-scene setup

Create `C19_process_reset/live/C19_head_working.blend` from an empty scene.
Set metres and real scale from the canonical 25 cm complete-plush height. Use
the tracked front image to measure the complete visible height and `Wh = 368
+/- 4 px`; record the resulting `Wh` in metres instead of hard-coding a
remembered estimate.

Create fixed orthographic review cameras facing `-Y` with `Z` up:

- `C19_front`, `C19_side`, `C19_rear`;
- `C19_three_quarter` and `C19_three_quarter_mirror`.

Register camera backgrounds, not renderable geometry:

- front controls width, vertical landmarks, and symmetry;
- the corresponding canonical-turn frames control depth, rear silhouette,
  and three-quarter volume;
- supporting photos may veto implausible cushion construction but may not
  override the canonical silhouette.

Align by head width, head center, crown, eye-line proxy, and later ground
plane. Keep reference opacity near `0.35`; use a `45%` silhouette overlay for
the decision packet. Do not judge an overlay when registration residual is
above `.015 Wh`. Change cameras rather than rotating the candidate.

Put all reference backgrounds, C1b render cards, cameras, and lights in
`C19_REVIEW_ONLY`. They must remain outside the reusable asset collection.

## Build only the first owner

1. Add one cube and name it `C19_HeadCushion`.
2. Scale it to the measured head width, a shallow front/rear depth within the
   canonical-turn envelope, and the lock-excluded head-height witness near
   `.986 +/- .03 Wh`.
3. Apply scale. Bevel only enough to remove literal cube corners, then apply
   that setup modifier once. Add a Multires modifier with just enough levels
   for broad sculpting; do not voxel-remesh, use Dyntopo, or add surface noise.
4. Assign one neutral beige clay material. Add no hair, bow, eyes, embroidery,
   seams, fibers, body, or rig in this cycle.
5. Sculpt at the lowest useful Multires level with broad Grab, Elastic
   Deform, Flatten, and light Smooth. Preserve a readable front plane, rear
   plane, crown/underside transition, side gusset volume, and slight stuffing
   asymmetry. Reject sphere, egg, box, helmet, or globally inflated-balloon
   reads.

The one-cycle hypothesis is: **a single directly sculpted, shallow cushion can
match the canonical front field and turntable depth while retaining flattened
fabric planes in every view**.

Do not add a second owner to rescue the first packet. Scripts may create
cameras, register references, save copies, hash, and render. They must not
interpolate or generate the artistic surface.

## Reversible micro-cycle

1. Reload the last kept snapshot; confirm active object, Multires level,
   camera registration, and source hashes.
2. Name one visible discrepancy, one controlling view, and one regression
   view. Make at most three broad strokes and one light smoothing stroke.
3. Inspect front, both profiles, rear, and both three-quarters in the live
   viewport. Undo immediately on a categorical sphere, box, twist, or depth
   failure.
4. Use **Save Copy** to publish an immutable
   `snapshots/C19_head_sNN.blend`. Hash it. Never render the live save target.
5. Clean-reopen that frozen snapshot in repository-pinned Blender 5.2.1 with
   automatic file scripts disabled. Render only the controlling camera and
   its highest-risk view at `512 x 512`.
6. Produce aligned reference/candidate overlays and an unlabeled C1b/candidate
   A/B. Review the candidate against references before revealing method or
   C1b identity.
7. Keep only a visible reference-directed improvement with no protected-view
   regression. Otherwise reopen the preceding snapshot; do not accumulate a
   failed delta.
8. A kept controlling pair earns one six-direction neutral-clay packet. Hair
   work begins only after this head-cushion packet passes.

Whole-file snapshots are the rollback boundary. Do not use topology-changing
tools between snapshots, and do not rely on Blender's numbered backups.

## Blender 5.1 live host versus pinned 5.2.1

Use one authoring host only.

- If the user-selected Flatpak Blender 5.1 MCP host is available, create the
  empty C19 scene there and save only 5.1-authored snapshots. C1b is a hidden
  image/render comparator, not a file to open or append. Never open and
  overwrite a newer 5.2.1 file in 5.1.
- If MCP is unavailable, launch one foreground Blender. Prefer the
  repository-pinned 5.2.1 host; the already-installed Flatpak is acceptable
  only under the explicit live-host exception. Do not leave both running.
- Repository-pinned Blender 5.2.1 always clean-reopens and renders the frozen
  snapshot of record. It may open an older 5.1 snapshot. A 5.1 authoring host
  must not reopen a file after 5.2.1 has saved it.

## Keep/undo gate for `C19_HeadCushion`

Keep only when all are true:

- critical front landmarks are within `3% Wh` and major silhouette/depth
  extrema within `5% Wh` after valid camera alignment;
- front and rear remain visibly flattened, with a shallow gusseted cushion
  rather than an uninterrupted dome;
- both profiles and three-quarters agree with the canonical-turn depth and do
  not reveal a sphere, egg, box, helmet, twist, or pinched edge;
- the lower head supplies a plausible future collar seat and the crown a
  plausible future hair/bow receiver without modeling either attachment; and
- an implementation-blind review recognizes a stuffed cushion construction,
  not plastic or a smooth primitive.

Undo on any categorical failure. If the first two bounded sculpt cycles fail
for different views, stop and change the head-cushion base proportions or
coarse topology before more strokes. Do not use hair to conceal a rejected
receiver.

## Feedback latency

Target `60 s`, hard maximum `90 s`, from immutable snapshot publication to an
inspected pinned controlling-pair sheet:

- Save Copy and hash: `<= 10 s`;
- pinned clean reopen and two 512 px renders: `<= 35 s`;
- overlay/contact sheet: `<= 15 s`;
- pixel verdict: complete by the `90 s` wall.

C17's pinned pair took `19.3 s`; C18's comparable pair took `25.67 s`.
Batch both views in one warm pinned process. If the pair misses `90 s`, repair
the snapshot/render harness before another stroke. Limit each stroke batch to
two minutes, so normal edit-to-decision latency stays below four minutes.

## First decisive artifact

Publish one immutable head-only packet containing neutral-clay and owner-ID
front, both profiles, rear, and both three-quarters; aligned silhouettes; one
uncropped normal-scale sheet; candidate/tool/hash manifest; and an absolute
review. It passes only the head-cushion owner. The full C19-0 macro blockout
still requires separate crown/rear hair, foreground leaf, bow, seated body,
skirt, sleeves, and feet owners before any material or rig work.
