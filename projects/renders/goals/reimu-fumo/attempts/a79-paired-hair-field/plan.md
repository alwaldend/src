# A79 complete paired hair-field module

## Objective

Replace the complete brown crown/rear silhouette owner with constructed,
layered plush hair panels that remove the helmet read while preserving the
existing face identity and all non-hair context. Produce a candidate that is
visibly preferred to rung 003 in front, both three-quarters, the worse profile,
and rear before any material-detail pass.

## Controlling references

- `canonical_front_25cm.png` controls exact-variant front width, crown/fringe
  transition, face exposure, lock widths, and 25 cm scale.
- All frames of `canonical_turn_180.gif` are authoritative for continuous
  layer order; frames 6--12, 17--20, and 25--26 are the primary front-3Q,
  profile, rear, and opposite-profile brackets.
- `physical_front.png` and `physical_side.png` control the thin padded hair
  panels, edge roll, overlap, and non-helmet construction.
- `clean_front.png`, `turn.gif`, and `sofa.gif` support pile softness,
  compression, and contact only; they cannot override the canonical pair.
- Rung 003 is the immutable fixed-camera baseline. A74 and A77 control the
  canopy/card/bald-rear failure categories.

## Frozen boundary

- Parent:
  `out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/reimu_fumo_working_rung_003.blend`
  at SHA-256
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.
- Tracked reusable asset remains SHA-256
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.
- Hide exactly:
  - `A44 continuous hair cap with smooth opening`
  - `A42 Left asymmetric rear lock`
  - `A42 Off-center main rear lock`
  - `A42 Short right rear lock`
  - `A42 Main lock left seated seam`
  - `A42 Main lock right seated seam`
- Preserve the receiver, seven face witnesses, five front fringe/temple
  panels, two cheek locks, crown seam, bow, body, garments, feet, materials,
  lights, cameras, animation, transforms, and every other source object.
- Work only in disjoint disposable copies under
  `out/reimu_fumo_attempt_079_paired_hair_field/`.

## Representation contract

Create a small complete assembly with these visual roles:

1. a front/crown padded field that seats behind the retained fringe and turns
   toward the side seams without becoming a full concentric shell;
2. an independent compact rear base field that closes every rear bald region
   and narrows under the bow before widening and tucking at the nape; and
3. an independent broad asymmetric rear leaf that owns the long profile depth
   and diagonal rear overlap.

Each visible fabric part uses explicitly authored paired outer and inner
surfaces with bridged boundaries and controlled edge roll. The hidden inner
surface may register to the receiver; the visible outer surface must be
authored from reference silhouette, camber, tension, and layer-order evidence.
Uniform Solidify, near-concentric receiver copies, sphere/cube receivers,
single-surface cards, generated per-height lofts, and materials-as-shape are
forbidden.

## Work graph

1. Freeze a source/landmark/interface contract from every controlling reference
   and the exact receiver/fringe/root relationships.
2. Preflight two independent representation variants in disjoint copies. Each
   must publish its topology and root-band model before any Blender candidate.
   Select at most one path by the mechanical and representation gates; do not
   merge two weak candidates.
3. Build the selected neutral-clay P0 with fresh objects/datablocks and full
   protected-source fingerprints. Publish builder source early.
4. Validate finite geometry, winding, manifoldness, duplicate/degenerate/self-
   intersecting faces, paired-skin thickness, root-band signed clearance and
   coverage, baseline-relative crossings, rear occlusion, modifier state, and
   clean reopen in pinned Blender 5.2.1.
5. Render a fast whole-context packet: front, both three-quarters, worse
   profile, and rear. If it survives, batch both profiles plus rear and an
   isolated neutral-clay construction view, then produce aligned
   reference/baseline/candidate boards and semantic masks.
6. Permit at most one correction to one named panel/control group when every
   other gate passes. Otherwise reset and change the representation.

## Mechanical gates before rendering

- Exactly the six declared source objects change visibility; all other source
  structure, transforms, visibility, materials, modifiers, custom data, and
  animation remain equivalent to the parent.
- Every new object is single-user, identity-transformed, consistently oriented,
  finite, and free of duplicate, degenerate, non-manifold, wire, or unintended
  self-intersecting elements.
- Evaluated paired-skin thickness remains inside a declared `2.0--5.0 mm`
  component band over at least 95% of sampled visible area; no Subdivision step
  may expand it into A77's canopy.
- At least 80% of every declared root band has positive signed receiver or
  retained-panel clearance inside its component-specific range; no minimum-
  vertex shortcut. Baseline-relative triangle crossings do not increase.
- Fixed-camera semantic masks show complete brown rear coverage, no beige leak,
  no disconnected island wider than one pixel, and no root gap wider than one
  pixel at review resolution.
- Protected hashes remain exact through clean reopen with no missing library,
  modifier, topology, or visibility drift.

## Visual gates

- Immediate self-veto: helmet, egg, mattress, canopy, curtain, card, blade,
  bald rear, beige leak, floating root, hard shell, or uniform plastic cap.
- Front crown widths, face exposure, asymmetric fringe relationship, profile
  compact-field/leaf separation, rear maximum-width height, unequal lower
  lobes, and long diagonal overlap stay inside the frozen reference bands.
- The candidate must be unambiguously preferred to rung 003 in every affected
  view. Macro silhouette scores at least `6/10`, constructed plush at least
  `5/10`, and contact/integration at least `6/10` in implementation-blind
  review, with no categorical defect.
- Final goal acceptance remains two independent reviews with every applicable
  category at least `8/10`; A79 passing is a working-rung gate, not final
  approval.

## Process checkpoint

The primary agent is the sole goal-record and selected-candidate writer.
Parallel workers use independent copies for reference measurement,
representation variants, validation, and implementation-blind review. Stop a
path immediately when its module or interface is wrong; do not spend render
time or add materials to rescue invalid geometry. Optimize for visible model
progress, not artifact volume.
