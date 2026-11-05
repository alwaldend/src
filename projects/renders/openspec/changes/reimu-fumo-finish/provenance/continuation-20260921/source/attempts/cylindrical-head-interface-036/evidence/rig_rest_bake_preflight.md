# Retained033 rig rest-bake preflight

All 41 live character meshes have uniform unit weight to one existing bone.
There are no genuine weight blends to preserve. Eighteen meshes evaluate
static lattices after their armature modifier, so the earlier proposal to
add a Hook after the curves' lattices does not give decorations the same
operation order as their neighboring meshes.

A rest-baked export copy is a credible smaller coordination model: preserve
85 separate visible objects, freeze each current evaluated rest shape, and
give each one final rigid deformation from its existing bone. It would
replace the inherited motion behavior, not preserve that behavior by proof.
The coordinator must make the material decision before implementation. No
Hook, bake, reparenting, action, export, or collection change occurred here.

## Exact evidence and scope

- Input: `out/reimu_fumo_finish/desktop_astra/foot_033_candidate.blend`.
- Input SHA-256:
  `98e92ee9a73ff49be32695dc06518ff885e5d91016278d16fb5a8771fd8fed48`.
  It matches `foot_033_writer_receipt.json` and remained unchanged before,
  during, and after the background inspection.
- Observation: `2026-09-05T02:03:42.792401+00:00`.
- Verifier: pinned background Blender `5.2.1 LTS`, build `9e2066aef7ef`;
  Bazel invocation `5f63337e-8225-4700-b974-1b6914c90e98`, exit 0.
- Script: [inspect_rig_rest_bake_033.py](inspect_rig_rest_bake_033.py), SHA-256
  `4813d4c37c7e97fcaa67694adc1f739fd4e0f783498a8c74190ea619ef3b1e83`.
- JSON: [rig_rest_bake_preflight.json](rig_rest_bake_preflight.json), SHA-256
  `d2909e9bee2fc6d48cf8df3ddc846d0bca8340ec7154ce5178948218ae2f489e`.
- Scope: render-enabled character MESH objects, excluding `Review floor`;
  the previously identified 44 decoration curves; and their 18 used
  materials' texture-coordinate metadata. No arrays were truncated.
- One pinned metadata inspection; saved frame 1 unchanged. No pose sweep,
  geometry edits, renders, saves, or delegation. The settled 44-curve lattice
  dependency graph was reused, not retraced.

## Weight purity and bone mapping

Every live mesh has one enabled armature modifier targeting `ReimuFumoRig`,
vertex-group deformation enabled, and bone envelopes disabled. Every base
vertex has one positive deform-bone weight of 1.0. No vertex is unweighted,
blended, or assigned a nonunit weight. No inspected mesh or decoration curve
has shape keys.

| Existing bone | Live meshes | Decoration curves | Proposed separate output objects |
| --- | ---: | ---: | ---: |
| `Head` | 15 | 0 | 15 |
| `Bow` | 9 | 6 | 15 |
| `Body` | 7 | 0 | 7 |
| `Arm_L` | 3 | 19 | 22 |
| `Arm_R` | 3 | 19 | 22 |
| `Leg_L` | 2 | 0 | 2 |
| `Leg_R` | 2 | 0 | 2 |
| Total | 41 | 44 | 85 |

The JSON maps every mesh name to `uniform_bone` and every decoration name to
`intended_bone`. The full decoration map is: all 19 `Sleeve44P L ...` curves
to `Arm_L`; all 19 `Sleeve44P R ...` curves to `Arm_R`; and each left/right
`A42 ... root fold 1`, `root fold 2`, and `white zigzag applique` to `Bow`.
No new rig or control is needed to express these assignments.

## Exact modifier ordering findings

The JSON's `meshes[].modifiers` records each complete ordered stack, exact
modifier names, target objects, enabled flags, and relevant resolution and
deformation settings. `post_armature_modifiers` isolates the ordering risk.

| Objects | Count | Suffix after Armature |
| --- | ---: | --- |
| All bow meshes | 9 | `022 bow proportion cage` lattice |
| Four sleeve panels | 4 | `022 body proportion cage`, then `023 sleeve root L/R` lattices |
| Bodice, seat pad, two cream leg roots, skirt | 5 | `022 body proportion cage`, then `023 narrow waist field` lattices |
| Hem026 strip | 1 | `Folded cotton thickness` Solidify |
| Head/hair/face, hands, collar/cravat, Foot033 pods | 22 | None |

Some bow loops additionally have the A154/A155 macro lattice before
Armature. Bow loop/tail ruffles have Solidify before Armature; the knot has
Bevel before Armature. The four sleeve panels are ordered Subsurf, Solidify,
Bevel, Armature, body lattice, side-specific root lattice. The skirt is
Subsurf, Solidify, Armature, body lattice, waist lattice. These are not
equivalent to a final curve Hook applied after both rest-shape lattices.

In a simplified notation, the present mesh path can be `L(R(x))`, whereas a
final Hook path is `R(L(x))`. Static nonuniform deformation and rotation do
not generally commute. This establishes a mechanism for disagreement; no
specific posed gap or clipping was measured in this task.

All inspected mesh modifier viewport/render enable flags agree, and every
Subsurf viewport level equals its render level. The six bow curves have
resolution 1, and the 38 sleeve curves resolution 2; all have render
resolution 0. These metadata checks reduce one mismatch risk but do not
certify a baked render surface or curve tessellation.

## Minimal proposed rest-bake/export-copy route

1. Preserve retained033 and its editable donors as immutable, hash-bound
   source files. Build a separately named candidate from the actual retained
   input. Keep the original rig as the sole rig in that candidate; do not
   append another rig. Do not keep duplicate editable donor meshes inside
   the reusable export collection. A donor manifest and the frozen source
   file provide recovery.

2. Capture each object's complete evaluated neutral surface separately,
   including current static lattice effects, subdivision, bevel, thickness,
   curve bevel/tessellation, and existing neutral attachment compensation.
   Keep one resulting object per original: 41 meshes plus 44 curve-derived
   meshes. Do not join or weld components. Preserve the world matrix and
   transform the captured geometry into that object's corresponding local
   space; avoid global transform application or origin changes. Record an
   original-object to output-object mapping.

3. Replace the export copy's baked stacks with one final armature modifier
   using the same bone assignments above. Give every output vertex a single
   unit weight to its assigned bone. Remove inherited duplicate deformation
   paths only on those export objects. There are no actual blends that need
   interpolation or transfer. Do not assume that copying the old modifier
   alone preserves the captured neutral output: bind against the chosen
   neutral rig/rest matrices and verify compensated evaluated world points.

4. Make the existing unconnected `Bow` bone a child of `Head`, retaining its
   armature-space rest matrix, length, and neutral world placement. Keep
   other controls unchanged unless a concrete pose failure demands a bounded
   change. Bow and its six decorations then share the same final rigid
   transform and inherit actual head yaw. Do not animate two parallel bow
   follow mechanisms.

5. Place the 85 intentional objects and the existing rig in one clearly
   named reusable collection, such as `ReimuFumo`. Preserve intentional
   material and action datablocks and all remaining dependencies. Once
   their shape contribution is captured and equality checked, the static
   lattices need not be in this baked collection; they remain recoverable
   in the editable donor file. Exclude historical alternatives, reference
   objects, the floor, cameras, lights, and other review content.

6. Only after neutral shape and shading checks succeed, reuse the empty
   neutral action and add the required head-yaw, arm-wave, and combined
   actions on this rig. Keep action switching deterministic with explicit
   neutral boundaries. Use Head's actual local-Y/world-Z yaw axis. No NLA
   layer, extra control bone, or duplicate rig is required for these four
   stored tests. Their safe amplitudes and resulting contacts remain to be
   established from posed renders.

This route keeps the editable manufacturing representation in its frozen
source while making the reusable copy's motion depend on one final rigid
transform per object. It sacrifices direct editing of the baked static
modifiers in the export copy and can increase geometry size. It does not
provide soft deformation across component boundaries; that was also absent
from these uniform single-bone weights.

## Shading and material preservation are separate gates

The coordinate risk is present in this exact candidate. Fourteen used
materials explicitly send Texture Coordinate `Generated` output to Noise
and/or Wave textures. Four sleeve materials have unlinked Noise Vector
inputs and need their implicit-coordinate behavior preserved and checked.
Only six of the 41 source meshes have a `UVMap`; the other 35 have no UV
layer. Existing UVs alone do not cover the procedural material behavior.

Baking changes the base mesh and can change its Generated coordinates,
texture-space bounds, interpolation, and normals. Keeping material names,
node trees, or numeric shader settings unchanged does not prove identical
color, weave, bump, or shading. Copying the old texture-space location/size
alone cannot generally invert a nonuniform lattice deformation.

Before baking, establish a coordinate-preservation method on the smallest
representative object with a Generated-driven material and a nonlinear
rest-shape stack. A candidate method is to carry the original rest-coordinate
field through the same surface evaluation into an explicit named vector
attribute on the baked output, then use that attribute in export-only
material copies wherever the old shader sampled Generated. This needs an
actual supported correspondence and numeric/render validation; it is not
implemented or guaranteed by this plan. Curve conversion, new Solidify rim
faces, and subdivision interpolation need explicit coverage. If that field
cannot be reproduced, report and review the shading deviation rather than
claiming exact preservation or silently changing the fidelity requirement.

Preserve slot order and evaluated face material assignment for each object,
not just the set of material names. In particular, each Foot033 pod has
slots `[Feet black velour.002, Dress warm white cloth.002]` and currently
2,894 black faces plus 1,314 white faces. Cheek locks, eyes, the head cushion,
and sleeves also have multiple slots. Source sleeve base faces use slot 0,
while modifier-generated surfaces can require the other slot; copying only
the base polygon material indices is insufficient. Preserve evaluated
normals/smoothing and needed UV/attribute layers as well.

Keep the original materials intact in the donor. If export-only coordinate
changes require material copies, retain one consistent copy per original
material where possible, preserve slot references, and verify that the new
named attributes survive append/reopen. Shared unchanged material datablocks
may remain shared. No texture image node was observed in the inspected used
material trees, but that is not a complete resource or export audit.

## Evidence needed before either route can pass

- A fresh candidate binding if the retained model changes beyond 033.
- Zero unintended evaluated neutral world-shape change, including curves,
  thickness, silhouette, and bow-root contact; preserved object separation.
- Correct material slots, evaluated material indices, normals, and
  coordinate fields. Same-camera rendered comparison is required for shading;
  unchanged nodes are only metadata evidence.
- Actual stored head-yaw, arm-wave, and combined-action render evidence,
  checking rigid-boundary contacts, shoulder gaps, bow seating, and feet/hem.
- One-collection append into a blank file followed by pinned clean reopen,
  with no missing rig, action, material, attribute, or geometry dependency.

No animation, visual, or reusable-structure criterion is passed by this
metadata preflight. Root retains the final choice between a validated
rest-bake candidate and any alternative that resolves the measured ordering
conflict without expanding controls unnecessarily.
