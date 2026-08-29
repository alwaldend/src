# Reimu Fumo attempt 12

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 12 — under-depth rolled-rim pillow puck

**Candidate:** planned
`out/reimu_fumo_attempt_012_head_pillow_puck/checkpoint_family_00.blend`,
SHA-256 pending, neutral bare-cushion stage, review packet
`attempt-012-pillow-puck-fixed-views`.

**Failure targeted:** Two sculpt families could move local vertices but could
not remove the sharp cube's six planar masses. Starting at final depth made
every edit subtract from a mattress and caused front, rear, and side tucks to
meet as a vertical three-quarter cleft.

**Hypothesis:** A shallow eight-sided pillow puck with two broad quiet caps and
one continuous rolled perimeter begins with the correct cushion construction
without becoming an ellipsoid. Voxel remeshing removes its cylinder and bevel
topology. Two low-strength side pulls can then build deliberately missing
front/rear depth outward at different heights, and three mirrored front pulls
can create the broad cheek/chin silhouette without subtractive corner
trenches.

**Plan written before implementation:**

1. Start factory-empty in a new temporary directory and copy no rejected
   object, mesh, modifier, coordinate, or checkpoint. Preserve the exact
   protected source hashes.
2. Create one native eight-sided cylinder along `Y`, rotated `22.5°` so its
   front outline has horizontal and vertical sides. Apply dimensions about
   `0.91 Wh` wide, `1.02 Wh` high, and only `0.52 Wh` deep.
3. Bevel front/rear perimeter rings broadly by `0.16–0.18 Wh` with four or
   five segments; bevel depth-running edges by `0.06–0.08 Wh` with three
   segments. Apply the construction, voxel-remesh once at about `0.025 Wh`,
   and shade smooth. Use no later subdivision, global Smooth, mesh filter,
   Dyntopo, or remesh.
4. Verify one connected closed mesh, no modifier-driven silhouette, and no
   surviving polygonal facet in neutral, silhouette, or grazing seed renders.
   Setup bounds are vetoes only, never likeness evidence.
5. In fixed side orthographic view with X symmetry disabled, use one broad
   ordinary Grab near front lower-middle from approximately
   `(y,z)=(-0.26,-0.06)` toward `(-0.345,-0.08)`, radius `0.29–0.32 P`,
   strength `0.78–0.85`. Pull rear middle slightly above it from about
   `(+0.26,+0.04)` toward `(+0.325,+0.01)`, radius `0.27–0.30 P`, at the same
   strength tier. Render side, front, and both three-quarter diagnostics.
6. Only if depth remains one continuous volume, enable X symmetry in fixed
   front view and execute at most three standard Grab strokes: upper shoulder
   `(0.39,+0.30)→(0.35,+0.34)`, lower corner
   `(0.39,-0.31)→(0.34,-0.37)`, and lower cheek
   `(0.37,-0.08)→(0.42,-0.10)`. Use radius tiers `0.18–0.22 P` and no
   repetition.
7. Do not conditionally smooth the first candidate. A long ridge or cleft is a
   seed/family rejection, not a cleanup task. Stop after the five Grab strokes
   and render fixed neutral, silhouette, and readable grazing front, side,
   rear, three-quarter, and mirrored-three-quarter views.
8. Hard-reject any full-height crease, central notch, waist, two-lobe split,
   dimple deeper than `0.01 Wh`, side roof or wall straight beyond
   `0.15–0.18 Wh`, parallel front/rear profile beyond `0.20 Wh`, hard face
   plate, visible octagonal facet, rounded-box, mattress, egg, balloon,
   helmet, molded-foam read, broken three-quarter highlight, or missing broad
   lower cheek and shallow chin arc.
9. Require a context-light implementation-blind review of the exact final
   pixels. Every fixed view and each of silhouette, sewn-cushion construction,
   plush-medium read, and presentation must reach at least `7/10`, with no
   major defect, before this clay can become a retopology target. Otherwise
   reject it and revise the seed or sculpt tool, not its detail.

**Pre-implementation calibration correction:** The strategy audit initially
suggested `0.22–0.30` strength for the two depth pulls. Exact Attempt 11 Grab
records show maximum center displacement was approximately stroke travel times
brush strength: `0.034–0.043 Wh` for about `0.085–0.108 Wh` travel at `0.40`
strength. The lower tier could add only about `0.05 Wh` total depth to the
`0.52 Wh` seed and would fail the frozen `0.62 Wh` minimum before visual
review. Step 5 therefore supersedes that suggestion with `0.78–0.85`, which
predicts about `0.12 Wh` total added front/rear depth while retaining broad
falloff, one shot per side, and the `0.76 Wh` upper veto.

**Attempt 12 work actually performed:** A hash-locked wrapper reused the exact
Attempt 11 interactive viewport, native-stroke, render, checkpoint, and review
machinery. It created one native eight-sided cylinder along `Y`, applied broad
five-segment cap-ring bevels and narrow three-segment depth-edge bevels,
removed the bevel attributes and modifiers, voxel-remeshed once at
`0.025 Wh`, and rendered the untouched seed. The seed was
`0.910 × 0.520 × 1.020 Wh`, one closed 5,648-vertex mesh with zero boundary or
non-manifold edges. The driver stopped before either sculpt family.

**Attempt 12 raw evidence:** The exact seed checkpoint and contact sheet hashes
are
`90d93bd9221b37914a9f47bc7da2d6cfe6e41abd2f7f6b9cbeb2d32a2ce0c8d6`
and
`9200edaf6243240a9d8a009702eba0a0415d2f125fb6ef87dc2fa7b6bb16d034`.
The context-light seed reviewer scored continuous volume `4/10`, quiet caps
`3/10`, rolled perimeter `1/10`, and facet-free `1/10`. It saw a regular
octagon in front/rear, a flattened parallel-wall capsule from the side, a
separate cap and chamfer belt, planar wedge facets, a lower-center V ridge, and
straight top/bottom plateaus.

**Attempt 12 criterion result and decision:** Factory-empty provenance,
under-depth bounds, closure, and modifier removal pass. The facet-free,
continuous-volume, quiet-cap, rolled-perimeter, sewn-cushion, and plush-medium
seed gates fail. Reject before native strokes. No coordinate, mesh, modifier,
or datablock may enter Attempt 13.

### Progress and approach audit after Attempt 12

- **Improved:** the seed-only interlock prevented five sculpt strokes from
  being spent on a visibly invalid base. The narrower proof reduced feedback
  time and isolated the exact representation defect before candidate work.
- **Regressed or unchanged:** the visual seed was worse than the smooth-cube
  block in facet freedom. Its regular octagonal silhouette and chamfer belt
  made the constructed-medium read explicit at only `1/10`.
- **Absolute result:** this is not a viable clay base and provides no visual
  progress toward the user-visible goal.
- **Evidence against continuing:** beveling eight longitudinal sectors creates
  discrete planar fields by construction. Voxel remesh changes topology but
  preserves those occupied planes, so Grab would only deform a faceted drum.
- **Highest-leverage unresolved problem:** produce two broad, subtly domed caps
  and one tangent-continuous rim with no resolvable chord, shoulder, flat
  gusset, radial egg curvature, or surface split before any native stroke.
- **Approach decision:** discard the eight-sector cylinder and bevel system.
  Use a dense shallow-domed cushion drum whose disposable seed geometry is
  analytically continuous, then destroy its topology before visual review and
  native sculpting.
