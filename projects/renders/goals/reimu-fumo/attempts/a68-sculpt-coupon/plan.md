# A68 — native-sculpt rear/crown cap coupon

## Immutable bindings

- Goal: `reimu-fumo`
- Goal resource version / checkpoint CAS token: `1`
- Goal generation: `1`
- Lifecycle generation: `1`
- Criteria revision: `1`
- Criteria digest:
  `sha256:c5522700389e76975e7978515c586433ca2058a6d5012ef45fbbadcb78a5740c`
- Goal-state digest:
  `sha256:69d29116e14d282349e6ad5c073453839be4872922c886d9c74a7b5161fbae35`
- Exact parent rung:
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Protected reusable asset:
  `sha256:489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`

## Target and hypothesis

Primary stable defect: `HD-01-rear-cranium-balloon`. Fixed regressions are
`HF-01-crown-face-aperture`, `BW-01-bow-identity`,
`LK-01-detached-wire-locks`, `RC-01-collapsed-rear-curtain`, and
`PM-01-hard-shell-medium`.

Test whether genuinely brush-sculpting only the rear/crown cap can replace the
parent rung's balloon-like depth with the photographed shallow, flatter-sided
stuffed cushion while preserving its stronger Fumo identity frame. Native
`bpy.ops.sculpt.brush_stroke` calls must move the tested cap geometry; direct
scripted vertex formulas do not satisfy the hypothesis.

The original fresh whole-head proposal is rejected before implementation.
Pixels-first review showed that the previous whole-head replacement regressed
every head/hair-sensitive category to `2-3/10`, chiefly by reopening the bald
crown, omitting the bow, detaching locks, and collapsing the rear curtain. A
joined head/hair surface would also conflict with the references' separately
readable cushion, cap, fringe, locks, and rear panels. The smallest credible
test is therefore the cap profile alone.

## Bounded result

- One isolated candidate copied from the exact parent rung under this `out/`
  directory; never overwrite the parent rung or protected reusable asset.
- Preserve the parent rung's bow, front fringe and face opening, cheek locks,
  rear curtain, body, camera, and lighting as immutable comparison context.
- Sculpt only one replacement or duplicate rear/crown cap module. Do not
  rebuild the face, fringe, locks, bow, body, or garment.
- No production topology, fibers, new seams, or material-detail work.
- One early voxel remesh is allowed; later geometry must be shaped by native
  sculpt strokes in a live `VIEW_3D` context.
- Render the controlling and regression views after every one to three broad
  strokes or one coherent stroke family.

## Parallel workstreams

1. `workstream_reference/`: all-reference normalized head/hair contract.
2. `workstream_foreground/`: Blender 5.2.1 foreground/MCP and real-stroke
   recovery.
3. `workstream_blind_review/`: pixels-first baseline failure review.
4. `workstream_sculpt_design/`: independent stroke-family and checkpoint
   design.

Each worker writes only its disjoint directory and cannot publish goal state.
The coordinator integrates the four reports before touching coupon geometry.

## Controlling evidence

- `projects/renders/blender/fumo/reimu_fumo/references/canonical_front_25cm.png`
  controls front identity, silhouette, landmarks, and 25 cm complete scale.
- `canonical_turn_180.gif` controls depth, rear/side silhouette, and layer
  order.
- `physical_front.png` and `physical_side.png` control plausible stuffed-fabric
  construction.
- Other tracked references provide qualitative continuity checks only, as
  documented in the reference packet README.

## Gates and stop conditions

1. Stop before candidate geometry if a live `VIEW_3D` context and a verified
   coordinate-changing native brush stroke cannot be established promptly.
2. Reject before a full packet if any preserved identity component moves,
   disappears, clips, or scores below the exact parent-rung baseline.
3. Require visible `HD-01` improvement in side and both three-quarter views:
   shallower cap depth, flatter stuffed side profile, and no new floating seam,
   hard ridge, or loss of root contact.
4. If `HD-01` survives two distinct broad stroke-family checkpoints, close the
   workflow as falsified rather than rebuilding adjacent hair.
5. A baseline-beating internal coupon requires at least `6/10` macro
   silhouette and `5/10` construction with no category below the parent rung.
   It is not an approval candidate; later approval still requires every
   applicable category at `8/10` and no major failure.

## Intended evidence

- All four workstream reports and any annotated comparisons.
- Native-stroke records with brush identifiers, views, coordinate hashes, and
  moved-vertex counts.
- Checkpoint `.blend` copies and fixed front/side/rear/three-quarter renders,
  including direct parent-versus-candidate comparisons.
- One implementation-blind absolute review and an explicit accept, refine, or
  reset decision against criteria 001–004, with criteria 005–008 unchanged and
  unclaimed.
