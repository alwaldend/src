# A78 purpose-built flattened head-rest module

## Objective

Create and judge one reusable, shallow stuffed head rest that supports the
existing Reimu face graphics and hair interfaces without a human-like cranium,
spherical muzzle, mattress profile, square card, dark cavity, or detached
beige islands. Work only in disposable copies of rung 003.

## Controlling evidence

- `canonical_front_25cm.png` controls the exact variant, exposed face field,
  frontal proportions, and 25 cm scale.
- Every frame of `canonical_turn_180.gif`, especially frames 6--12, 17--20,
  and 25--26, controls oblique/profile/rear depth and layer order.
- `physical_front.png`, `physical_side.png`, `clean_front.png`, `turn.gif`,
  and `sofa.gif` support sewn-panel, compression, and plush-medium judgments
  but cannot override the canonical pair.
- Rung 003 is the immutable whole-plush baseline. A75 controls the explicit
  failure categories; A77 controls the interface and process failures.

## Hypothesis

A shallow sewn pillow made from an explicitly authored front panel, rear
panel, and perimeter gusset can satisfy the frozen face witnesses and existing
hair aperture while reading as a stuffed fabric receiver. A denser local front
control field can resolve the A77 anchor conflict without inflating the whole
head or creating A75's spherical muzzle.

## Frozen boundary

- Source:
  `out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/reimu_fumo_working_rung_003.blend`
  at SHA-256
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.
- Tracked reusable asset remains SHA-256
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.
- Hide only `Head_Cushion_Manual_Target`.
- Freeze all hair, bow, garments, feet, facial witnesses, lighting, cameras,
  materials, transforms, and animation data.
- Create fresh mesh datablocks and save only under
  `out/reimu_fumo_attempt_078_head_rest/`.

## Work graph

1. Measure the visible face field, frozen witness surfaces, hair aperture, and
   existing receiver depth in all controlling views. Reconcile A75/A77 gates.
2. Preflight an explicit panel topology: front grid with dedicated eye, mouth,
   cheek-roll, and boundary controls; independent rear panel; connected gusset;
   consistent winding; no Solidify or Geometry Nodes.
3. Build one deterministic neutral-clay P0. Check finite geometry, components,
   manifoldness, duplicate/degenerate faces, boundary loops, winding, modifier
   stack, object transforms, datablock ownership, and protected hashes.
4. Measure signed witness clearances and coverage, evaluated front planarity,
   cheek roll, head depth, seam seating, and baseline-relative crossings. Clean
   reopen in pinned Blender 5.2.1.
5. If mechanical gates pass, render the fixed front, both three-quarters, both
   profiles, and rear in whole-plush context, plus an isolated neutral-clay
   panel/gusset diagnostic and aligned reference/baseline/candidate board.
6. Apply at most one correction to one named front-panel control group only if
   every other gate passes. Otherwise reset and record the representation
   failure.

## Mechanical gates before rendering

- One connected, consistently oriented, closed/manifold receiver with explicit
  front, rear, and gusset regions; zero duplicate or degenerate faces.
- No sphere primitive, shrinkwrap-generated shell, receiver loft, Solidify,
  Geometry Nodes, or post-hoc whole-form smoothing.
- Every frozen eye and mouth witness has positive signed surface clearance in
  its declared band with no crossing; at least 80% of each attachment/root
  band lies inside its band rather than passing on one nearest vertex.
- No new baseline-relative triangle crossing, exposed gap, or beige island at
  retained hair/fringe interfaces.
- Evaluated support depth is bounded and visibly shallower than rung 003 while
  retaining a broad stuffed cheek roll; front field is not planar cardboard.
- Protected hashes exact before and after save; clean reopen reports no missing
  libraries or topology/interface drift.

## Visual gates

- Absolute self-veto: no muzzle, sphere, egg, mattress, square card, cavity
  ring, detached island, anatomical skull, or taut plastic surface.
- The affected exposed face must be clearly preferred to rung 003 in front and
  both three-quarters, with no regression in either profile or rear.
- Macro/identity, constructed plush, and contact/integration must each score at
  least `6/10` in an implementation-blind review to become a working rung.
- Final-goal approval remains every frozen category at least `8/10`; passing
  this module is necessary evidence, not final acceptance.

## Process checkpoint

Publish the builder early. Run pure topology/interface probes before Blender,
batch all authorized renders in one pinned process, and stop immediately when
the representation or interface is disproved. Keep the primary agent as the
only goal-record writer; parallel workers use disjoint scratch paths.
