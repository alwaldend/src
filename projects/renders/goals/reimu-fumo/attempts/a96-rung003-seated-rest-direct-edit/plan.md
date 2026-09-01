# A96 plan — rung003 seated-rest direct-edit coupon

## Binding

- Goal: `projects/renders/goals/reimu-fumo`
- Goal generation: 1
- Lifecycle generation: 5
- Criteria revision: 1
- Expected resource version at start: 67
- Work type: change
- Parent Blend: `out/reimu_fumo_working_ladder/`
  `rung_003_eyes_locks_sleeves/reimu_fumo_working_rung_003.blend`
- Parent SHA-256:
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Canonical front reference:
  `projects/renders/blender/fumo/reimu_fumo/references/`
  `canonical_front_25cm.png`
- Boundary audit:
  `out/reimu_fumo_attempt_053_rung3_file_audit/FILE_AUDIT.md`

The tracked canonical Blend is read-only throughout this attempt. All output
is isolated beneath `out/reimu_fumo_attempt_096_rung003_seated_rest/`.

## Objective

Improve the proven rung003 model's weakest coupled region without rebuilding
the character: reshape its existing lower assembly into a compact, pooled,
believably seated plush rest. The skirt must stop reading as a cone or tent;
two stuffed feet must sit partly beneath the hem with broad contact; the
recognizable upper model must remain pixel-identical apart from occlusion.

## Frozen and editable ownership

Freeze by digest, object identity, topology, vertex count, modifier order, and
material assignment:

- head, hair, face, eyes, locks, and bow;
- bodice, collars, and cravat;
- retained `Sleeve44P` padded panels, roots, pleats, and stitches;
- cameras, lights, world, materials, and render settings.

Edit only the existing seat pad, skirt panels, hem/ruffles, feet, and leg
roots identified by the boundary audit. Use sparse coordinate changes stored
as reversible per-object shape-key deltas. Pin upper attachment rows and carry
owned stitch/ruffle curves with their surfaces. No generator, global affine,
new visible owner, remesh, join, topology change, material work, or rigging.

## Parallel workstreams

1. The author works on one isolated parent copy and produces the first coherent
   lower-body candidate plus fixed front, exact-side, and worst-three-quarter
   renders.
2. A verification worker independently fingerprints frozen/editable ownership,
   renders the exact baseline views, and supplies a reusable comparison and
   topology/frozen-state check. It does not modify the author Blend.
3. The coordinator validates provenance and actually inspects all controlling
   pixels. Only after a candidate passes the category gate will one fresh
   independent reviewer receive a blinded comparison.

## Ordered execution

1. Verify the parent and canonical tracked Blend hashes.
2. Inventory the lower-stack objects and freeze all protected fingerprints.
3. Render and inspect exact-parent front, side, and three-quarter controls.
4. Copy the exact parent into the author workspace.
5. Add reversible shape keys and make one coupled sparse edit:
   broaden the seated side footprint toward `1.15–1.20 Wh`, pool the skirt,
   reduce lower-body height, tuck and soften both feet, and establish broad
   contact while preserving upper attachment rows.
6. Render front, exact side, and worst three-quarter immediately.
7. Verify non-black/non-blank pixels, source/candidate hashes, frozen-object
   fingerprints, topology, modifier order, materials, object identities, and
   separation of feet/panels/leg roots.
8. Inspect absolute candidate pixels before comparison. Reject immediately on
   a cone/tent/ramp/cape/rail silhouette, missing or floating feet, clipping,
   gaps, or upper-character regression.
9. If the absolute gate passes, make a compact baseline/candidate board and
   route exactly one fresh implementation-blind reviewer. Permit at most one
   correction only when the category already passes and one localized contact
   defect remains.

## Acceptance gate

Keep the candidate only if every controlling view:

- clearly reads as a compact pooled seated skirt;
- shows two partly occluded stuffed foot pods with broad believable contact;
- preserves Reimu recognition and all frozen interfaces;
- has no sleeve/collar gap, clipping, floating element, or floor rail;
- reaches at least internal `6/10` for lower silhouette, construction, and
  contact; and
- is clearly preferred to exact rung003 in blinded comparison.

Numeric measurements can veto regressions but cannot qualify weak pixels.
Animation readiness requires stable object identities and topology, separate
panels/feet/leg roots, pinned attachment loops, unchanged modifier/material
state, and reversible shape-key edits. Rigging is deferred until visual
survival.

## Stop and strategy conditions

- Stop early on invalid/black baseline pixels, wrong source hash, protected
  mutation, absent broad coordinate effect, or any categorical silhouette
  recurrence.
- A tie or invisible change is a reset, not progress.
- If broad direct control works but the complete pixels still fail, close A96
  and make an explicit feasibility decision before any further rebuild.
- If the route fails only because sparse direct control cannot produce the
  required displacement, the simplified pointer-sculpt discriminator becomes
  the only remaining authoring-control test.

