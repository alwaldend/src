# Reimu Fumo attempt 14

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 14 — reference-derived asymmetric cloth panel sack

**Candidate:** planned
`out/reimu_fumo_attempt_014_cloth_panel_sack/checkpoint_setup.blend`, SHA-256
pending, untouched frame-90 cloth-settled panel stage, review packet
`attempt-014-panel-sack-reference-and-fixed-views`.

**Failure targeted:** Attempt 13's valid topology still encoded every rim ring
as one uniformly scaled copy. That construction forced a constant bezel,
parallel sidewall, repeated front/rear cap, and molded-foam regularity even
after subdivision and fine voxel remesh.

**Hypothesis:** A welded three-panel pattern whose front, rear, crown, temple,
lower cheek, underside, and depth are independently derived from the physical
and turn references can express sewn construction before identity details.
Low-pressure Blender cloth settling can distribute curvature through that
already-near-target cage without turning it into a free balloon. One later
bounded native Grab family can then add stuffing compression and break the
remaining manufactured symmetry before the full plush-medium gate.

**Plan written before implementation:**

1. Start factory-empty under a byte-new temporary directory. Import no Attempt
   13 object, coordinate, mesh, modifier, material, image, world, collection,
   scene, or cached simulation. Explicitly purge and assert the permitted
   datablock set. Reuse only audited Blender context, evidence, and checkpoint
   functions whose own bytes are hash-locked.
2. Build one closed welded manifold from a distinct front panel, rear panel,
   and nonuniform gusset: 48 perimeter stations, eight concentric rings plus a
   center on each panel, and five intermediate gusset rows. Front/rear seam
   loops and gusset faces share vertices. Use no loose sewing edge, open
   pressure boundary, bevel, separate rim strip, homothetic full-outline belt,
   or prior candidate coordinate.
3. Derive the mirrored front trace in `Wh` units from these controlling points:
   top `(0,+0.50)`, shoulder `(±0.30,+0.47)`, upper side
   `(±0.44,+0.32)`, widest cheek `(±0.46,+0.02)`, lower cheek
   `(±0.44,-0.27)`, lower corner `(±0.32,-0.46)`, and bottom
   `(0,-0.50)`. Sample a smooth closed spline, retain `0.97` of its outline at
   the front and rear seam loops, and permit at most `0.008 Wh` declared
   asymmetry. The front panel must own nearly the full silhouette; an inset
   face plate or constant bezel is forbidden.
4. Vary seam depth with normalized height `u=z/0.50`:
   `d_front=0.18+0.14(1-u²)-0.012u` and
   `d_rear=0.18+0.14(1-u²)+0.012u`. Place the front seam at `-d_front` and the
   rear at `+d_rear`. Make the front and rear panel centers reach only
   approximately `y=-0.34` and `+0.34`, using a smooth radial-square falloff.
   Front fullness must peak lower and rear fullness higher. Mid-height depth is
   approximately `0.64 Wh`; crown and underside depth fall to approximately
   `0.36–0.38 Wh` instead of forming parallel walls.
5. Interpolate five gusset rows between the shared seam loops. Roll their XZ
   coordinates outward by `sin(pi*t)` with independent maximum amplitudes:
   `0.010 Wh` at the lateral sides, `0.018 Wh` at the crown, and `0.012 Wh`
   downward at the underside. Require a continuous side taper and distinct
   crown/underside arcs. A constant-Y strip, flat roof, straight underside, or
   single perimeter radius is forbidden.
6. Verify the initial cage has one connected component, no boundary or
   non-manifold edge, no zero-area face, coherent orientation, positive signed
   volume, no self-intersection, and approximately 1,000 vertices. Preserve its
   exact coordinates and metrics only as temporary proof; it is not candidate
   geometry.
7. Add exactly one Cloth modifier to the already-near-target closed cage with
   `quality=12`, `time_scale=0.65`, `mass=0.20`, angular bending, stiffness
   `35/25/20/0.18` for tension/compression/shear/bending, damping
   `10/10/10/2`, `air_damping=3`, and zero gravity influence. Enable pressure
   with `uniform_pressure_force=1.5`, `use_pressure_volume=True`,
   `target_volume=1.035 × initial_closed_volume`, `pressure_factor=8`, zero
   fluid density, no internal springs, and no sewing springs. Enable
   Cloth collision quality four, enable self-collision with minimum distance
   `0.006 Wh` and friction five, and change no parameter after the simulation
   begins.
8. Evaluate frames 1 through 90 interactively. Frame 90 is the only candidate.
   Reject the simulation if RMS motion from frame 80 to 90 exceeds
   `0.0015 Wh`, volume changes more than `0.5%`, bounding-box dimensions change
   more than `0.003 Wh`, a self-intersection appears, or normalization would
   need more than `3%` scale. Recenter and uniformly normalize width only; use
   no axis-specific rescue scale.
9. Copy the evaluated frame-90 mesh into one candidate object, remove Cloth and
   every cache/modifier, and keep the pre-simulation cage hidden in a temporary
   proof collection only. Do not voxel-remesh or subdivide before the seed
   pixel gate. Mark provenance honestly as generated panel construction plus
   native cloth settling, not direct sculpt.
10. Render fixed 640-pixel neutral, silhouette, and readable grazing front,
    side, rear, canonical three-quarter, and mirrored-three-quarter views.
    Build a labeled comparison packet containing the exact controlling physical
    front, physical oblique side, turn three-quarter, turn side, and turn rear
    reference frames beside the candidate. Bind every review to the exact
    script, audited base, probe, checkpoint, seed metrics, reference, render,
    and contact-sheet hashes. Reject any stale byte and every non-finite or
    out-of-range score.
11. Seed numeric bands are width `0.90–0.95 Wh`, height `0.98–1.04 Wh`, depth
    `0.64–0.74 Wh`, and front controlling-landmark error at most `0.05 Wh`.
    Hard-reject a cap/rim shoulder, constant side band, front/rear parallel run
    over `0.12 Wh`, side chord over `0.12 Wh`, top/bottom roof over `0.10 Wh`,
    seam trench deeper than `0.008 Wh`, radial pinch, self-intersection, box,
    mattress, egg, balloon, helmet, or molded-foam read.
12. A context-light reviewer must inspect the controlling references and exact
    candidate pixels while blind to implementation. Require reference
    macroform at least `6/10`, continuous convex volume at least `7/10`, panel
    transition at least `6/10`, soft stuffing at least `6/10`, and non-box/
    non-egg read at least `7/10`, with no major defect. Failure ends Attempt 14
    before cleanup.
13. Only after seed PASS, copy the exact approved frame-90 mesh and apply one
    topology-normalization operation: one Catmull-Clark level if brush
    neighborhoods are even, otherwise one `0.018 Wh` voxel remesh. Predeclare
    and execute exactly six standard Grab strokes with no Smooth. In front view
    with X symmetry: shoulder `(0.34,+0.37)→(0.365,+0.365)`, radius `0.22 P`,
    strength `0.70`; lower cheek `(0.405,-0.18)→(0.425,-0.19)`, radius
    `0.24 P`, strength `0.70`. In side view without symmetry: front middle
    `(-0.335,-0.08)→(-0.355,-0.10)`, radius `0.25 P`, strength `0.70`; rear
    middle `(+0.33,+0.10)→(+0.35,+0.12)`, radius `0.24 P`, strength `0.70`;
    lower front `(-0.22,-0.43)→(-0.20,-0.45)`, radius `0.18 P`, strength
    `0.65`; and lower rear `(+0.21,-0.42)→(+0.19,-0.445)`, radius `0.18 P`,
    strength `0.65`. Stop and rerender; no conditional repetition is allowed.
14. The family gate requires every canonical view at least `7/10`, soft stuffed
    volume at least `7/10`, panel plausibility at least `6/10`, non-molded/
    non-primitive read at least `7/10`, numeric bands intact, and no hard-fail
    defect. A PASS accepts only a macro sculpt and later retopology target, not
    final plush, material, face, hair, bow, animation, or reusable delivery.

**Structural difference from prior failures:** A sphere would repeat the
documented egg/balloon failures and establish no sewn construction. This cage's
front and rear seams retain `97%` of the silhouette instead of making an inset
cap; depth varies independently with height and between front/rear; crown,
lateral, and underside gusset rolls have different amplitudes; cloth only
settles an already-shaped closed pattern; and a fixed convergence veto prevents
simulation tuning from silently creating multiple candidates.

**Pre-implementation API correction:** Blender 5.1 exposes
`ClothCollisionSettings.collision_quality` but no separate
`self_collision_quality` property. Step 7 therefore means
`collision_quality=4` with `use_self_collision=True`; the self-distance and
self-friction values remain exact. This correction was verified read-only from
the running Blender RNA before candidate geometry existed.

**Attempt 14 work actually performed:** The dormant driver passed syntax and
architecture checks, purged the candidate scene, created one coherent closed
1,010-vertex, 2,064-edge, 1,056-face panel cage with signed volume
`0.476913 Wh³`, and saved the fixed pre-simulation checkpoint. It added the
exact recorded Cloth settings and evaluated frames 1 through 90 interactively.
The numeric gate stopped execution during uniform normalization. No candidate
evidence packet, seed state, seed review, subdivision, remesh, or Sculpt Mode
stroke was created.

**Attempt 14 raw evidence:** The exact driver, simulation-setup checkpoint,
terminal numeric-failure checkpoint, and simulation probe hashes are
`8bee397b83d2f2661e57c04e3be32465421893423a0d7abc129e8aecbd25fa06`,
`fffab0271668ed7b6a78ddf0d0ed8bcf95d7d091164de9eefcfcc548c0022c75`,
`defdf74884f017a9684645485854e673d5254c1019fff8151cfc884949d71004`,
and
`02e8604aa71f990103acf9e3da8e58af320ac9ec417929eaf10d7200ccdd35ca`.
Frame 80 measured `1.076 × 0.998 × 0.996 Wh`; frame 90 measured
`1.082 × 1.053 × 0.992 Wh`. RMS vertex motion was `0.069870 Wh`, relative
volume change was only `0.000026`, and the maximum bounding-box change was
`0.054189 Wh` in depth. Width normalization would have been `14.987%`, five
times the frozen upper limit. The terminal checkpoint preserves the live Cloth
cage, copied frame-90 mesh, and hidden original pattern without changing them.

**Attempt 14 diagnostic evidence:** After terminal rejection, the untouched
pre-simulation proof pattern was rendered solely to evaluate the failed
approach. Its `0.914 × 0.680 × 1.000 Wh` fixed-view contact sheet hash is
`e31831239ae73515d6d29582acf2d8a7ef888c7802b7b45061935d6f4655b779`.
It visibly retains concentric radial facets, a broad rounded-box front, flat
top/bottom side roofs, and a hexagonal side profile. It is not an alternate
candidate and cannot enter a later attempt.

**Attempt 14 criterion result and decision:** Initial closure, connectivity,
orientation, positive volume, declared vertex tier, native Blender execution,
and volume stability pass. Convergence, translation stability, scale,
bounding-box stability, pixel eligibility, and presentation evidence fail.
Reject before render and native sculpt. The fixed plan forbids pressure,
stiffness, damping, anchor, frame, or scale tuning inside this attempt.

### Progress and approach audit after Attempt 14

- **Improved:** the simulation interlock caught an unstable cloth solve before
  a distorted mesh could be normalized into apparent compliance. The initial
  mesh had explicit reference-owned sections and a valid closed topology.
- **Regressed or unchanged:** the frame-90 solve was less controlled than every
  static seed: it drifted, expanded roughly `18%` from initial width, remained
  in large-amplitude motion, and became almost as deep as it was wide. The
  unsimulated pattern also exposed another faceted rounded box.
- **Absolute result:** no reviewable seed exists and no visual criterion passes.
  Attempt 14 is a technical and artistic failure, not progress toward approval.
- **Evidence against continuing:** pressure settling without spatial anchors
  permits global drift and broad inflation; the pattern's radial rings also
  remain visually encoded when settling is insufficient. Adding anchors or
  tuning pressure would define a new simulation family, while the starting
  pattern itself already fails the absolute pixels-only diagnostic.
- **Workflow defects found by static audit:** the dormant driver's successful-
  state idempotence would not reload the bound scene; numeric gates omitted
  finite-value assertions; expected reference hashes were not hard-coded;
  one zero-user pattern mesh would remain; montage columns were not fully
  aligned; convergence sampled only endpoints; and the simulation probe was
  not included in the intended review binding. None caused the observed
  failure because execution stopped earlier, but no later driver may repeat
  them.
- **Highest-leverage unresolved problem:** stop asking an untouched generated
  solid or simulation to look like final plush. Use a deliberately disposable,
  under-scale, topology-neutral tool blank and make the first judgeable object
  only through a bounded family of verified native Sculpt Mode strokes against
  the controlling references.
- **Approach decision:** discard pressure simulation, radial panel rings, and
  untouched-seed plush gates. Attempt 15 will start byte-new, prove only
  technical validity of the tool blank, then execute one predeclared direct-
  sculpt macro family before any likeness claim. It must hard-code reference
  hashes, reject every non-finite metric, load exact checkpoints on idempotent
  return, remove every orphan, and align candidate/reference evidence by view.
