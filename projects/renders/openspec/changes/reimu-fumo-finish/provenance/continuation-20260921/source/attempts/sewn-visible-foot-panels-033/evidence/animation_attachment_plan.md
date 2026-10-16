# Draft attachment, action, and reusable-collection plan

This is a proposed later-stage plan, not an implementation or stage pass.
The coordinator owns the choice and all saved-model changes. Keep the
existing `ReimuFumoRig` and its controls. Preserve the evaluated neutral world
shape before adding motion or packaging the asset.

## Exact evidence

Both inspections opened only
`out/reimu_fumo_finish/desktop_astra/head_032_candidate.blend`, SHA-256
`6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`.
The digest remained unchanged after the processes exited.

- Prior inventory: [animation_readiness_032.json](animation_readiness_032.json),
  SHA-256
  `a7a3a77dcdcd1718e5f8801428084ebb3cc7fb8851657196b746933e461c9bab`.
- Dependency evidence:
  [animation_attachment_dependencies_032.json](animation_attachment_dependencies_032.json),
  SHA-256
  `07f278ff740db773f753f8330da24e5c52cc788be3add9029665f754c31ef7e7`.
- Follow-up script:
  [trace_animation_attachments_032.py](trace_animation_attachments_032.py),
  SHA-256
  `72dc3ec32035cb2db17593fa554f0997ba1b1ac64418e07f00fd06bd26f6ae7b`.
- Observation: `2026-09-05T01:15:04.431837+00:00`.
- Pinned background Blender `5.2.1 LTS`, build `9e2066aef7ef`;
  Bazel invocation `986e1be5-a93c-4c7c-8e95-a9c0d1eac824`, exit 0.
- Scope: 44 decoration curves and six reachable lattice objects. No
  truncation, pose changes, frame changes, renders, or saves. Frame stayed 1.

The original JSON omitted non-armature modifiers. The follow-up therefore
asked a different, bounded question: whether those modifiers lead to a rig
dependency. It did not repeat the whole-scene inventory.

## What drives the 44 decorations now

| Decoration set | Count | Ordered active lattice modifiers |
| --- | ---: | --- |
| Left sleeve seams, cuff edges, stitches, pleats | 19 | `022 body proportion cage`, then `023 sleeve root L` |
| Right sleeve seams, cuff edges, stitches, pleats | 19 | `022 body proportion cage`, then `023 sleeve root R` |
| Left bow root folds and zigzag | 3 | `A154 Left loop macro cage`, then `022 bow proportion cage` |
| Right bow root folds and zigzag | 3 | `A155 Right loop macro cage`, then `022 bow proportion cage` |

All these lattice modifiers have strength 1 and are enabled for viewport and
render. Each of the six lattices has no parent, modifier, object constraint,
object/data action, or object/data driver. The decoration objects themselves
also have no parent, object constraint, object/data action, or object/data
driver. They are indirectly shaped by lattices, but there is no rig path in
the inspected graph. Moving an arm or bow bone has no attachment mechanism
through these object/modifier/data routes.

One narrow dependency channel remains uninspected: drivers on a separate
shape-key datablock. Do not call the entire dependency closure proven until
that channel is checked on the eventual frozen implementation input. The
single additional metadata inspection permitted for this task is complete.

## Minimal proposed implementation

1. **Freeze the eventual input and neutral output.** The plan is based on
   head032, not candidates 033 or 034. The sole writer first binds the actual
   retained input digest. Record the current armature/bone rest matrices and
   evaluated world-space positions and materials of affected meshes and
   curve tessellations. Use the existing neutral pose and review cameras.
   Check the unresolved shape-key driver channel as part of this preflight.

2. **Make the existing Bow bone a child of Head without moving rest geometry.**
   Keep Bow's armature-space rest matrix, length, and world rest placement;
   use an unconnected parent relationship. Keep Head, body, and arm controls
   otherwise unchanged. Do not add a second head/bow controller or compensate
   the bow with redundant action keys. Recheck the neutral evaluated bow and
   its root contact immediately after the parent edit.

3. **Attach the 44 curves after their existing rest-shape modifiers.** The
   smallest candidate is one fully compensated Hook modifier per curve,
   targeting the existing rig bone: the 19 left-sleeve curves to `Arm_L`,
   the 19 right-sleeve curves to `Arm_R`, and all six bow curves to `Bow`.
   Include every relevant curve point and Bezier handle, use strength 1 and
   no falloff, and bind the inverse at the preserved neutral transform.
   Append the deformation after the two existing lattice modifiers so the
   lattices continue to define the neutral shape. This is a proposal to
   test, not an assertion that the binding or modifier ordering already
   matches the adjacent mesh motion.

   Do not simply bone-parent the already deformed curves: changing object
   space under fixed lattice objects can change their shape. Do not parent
   shared rest-shape lattices to arm or bow bones: they can affect multiple
   components and create double deformation. Prove a zero neutral delta
   before exercising one affected arm/bow movement. If the Hook route does
   not preserve the receiver relationship, stop at that causal failure;
   evaluate a rest-baked, bone-bound decoration copy as a new bounded choice,
   keeping the editable source immutable.

4. **Store the four required actions on the existing rig.** Reuse the empty
   neutral action as `ReimuFumo_SeatedNeutral`. Add `ReimuFumo_HeadYaw`,
   `ReimuFumo_ArmWave`, and `ReimuFumo_Combined`. Preserve unused actions with
   explicit retention and valid action-slot binding; NLA tracks are not
   necessary for this validation workflow. Explicitly key neutral pose
   channels at clip boundaries so switching actions cannot inherit stale
   values. Keep armature object transforms fixed and finish each motion at
   neutral.

   A practical initial test envelope, not a new acceptance criterion, is
   frames 1 through 49 at 24 fps with approximately 20 degrees of head yaw
   in each direction and a modest outward wave. In this rig's recorded
   neutral basis, Head local Y maps to world vertical; local Z does not.
   Establish the arm's outward world-space swing from its actual pivot and
   basis before choosing its keyed Euler component. The combined action
   must overlap head yaw and an arm wave in time. Bow inherits Head motion;
   its local neutral offset remains unchanged.

5. **Package one named reusable collection without duplicating the rig.**
   Create `ReimuFumo` containing the intentional current character objects,
   `ReimuFumoRig`, and the complete set of still-required deformation inputs.
   The six lattices are live shape dependencies; they are not disposable
   review helpers merely because they are named cages. Include every other
   required dependency found in the final collection audit, or bake it on
   the export copy with verified zero evaluated-shape change. Keep historical
   alternatives, review floor, reference objects, cameras, and lights outside
   this collection. Retain material/action datablocks through their proper
   users. Avoid global transform application, origin changes, or deletion of
   unused controls unless a specific acceptance failure requires them.

6. **Validate the exact final candidate.** After the appropriate visual gate,
   exercise the stored actions and inspect fixed-view renders of neutral,
   yaw extrema, wave extrema, and combined motion. Check sleeve decoration
   alignment, shoulder contact, head/bow root seating, hair/collar overlap,
   and hem/feet contacts. Append only `ReimuFumo` into a blank file, save a
   new test file, clean-reopen it with pinned Blender, and confirm that every
   required object, action, material, and dependency survives and produces
   the same evaluated neutral and animated outputs. These are the existing
   criterion-005/006 obligations, not evidence already supplied here.

## Unresolved risks

- Shape-key datablock driver links were not captured in the bounded follow-up.
- Current deformation modifier ordering and compensation on the neighboring
  sleeve/bow meshes may differ from a final Hook's motion. A zero neutral
  delta alone cannot prove that their animated paths agree.
- Bow's current rest pivot may differ from its visible seated root. Preserving
  its rest matrix avoids a neutral jump but does not prove head-turn contact.
- Arm pivots and shoulder compression have not been exercised. No wave
  amplitude is certified by this report.
- Other required dependencies and resources across the full reusable
  collection have not been closed or tested by a blank-file append.
- Candidates 033/034 may change the final input. This plan must be rebound to
  the retained bytes before implementation; it grants neither candidate a
  visual or animation pass.

The asset still fails stored-action coverage and reusable packaging. The
traced curves have only observed rest-shape dependencies; their eventual
motion and all tearing/clipping claims remain unverified until rendered.
