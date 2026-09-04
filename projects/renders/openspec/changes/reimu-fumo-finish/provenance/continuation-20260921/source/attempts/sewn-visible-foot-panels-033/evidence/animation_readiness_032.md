# Head032 animation and reuse readiness

The retained candidate has a basic rig but is not animation-ready or packaged
as the required reusable asset. The sole stored action is empty, and none of
the named collections includes the rig. This report grants no stage pass.

## Evidence binding and scope

- Candidate: `out/reimu_fumo_finish/desktop_astra/head_032_candidate.blend`.
- Candidate SHA-256:
  `6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`.
- The same digest was verified before loading, after inspection, and after
  Blender exited. No save, frame change, pose sweep, render, or model mutation
  occurred. Only this immutable candidate was opened.
- Observation: `2026-09-05T01:09:39.523907+00:00`.
- Verifier: pinned Blender `5.2.1 LTS`, build `9e2066aef7ef`, background mode,
  launched by `bazel_agent bazel run //tools/blender:blender` with
  `--factory-startup --disable-autoexec --python-exit-code 2`.
- Bazel invocation: `787e187e-bb45-4e26-9ad4-17eb402de631`; process exit 0.
- Source workspace: linked feature worktree branch
  `t3code/continue-fumo-desktop-use`, observed commit
  `1e6704af1277420c21d34c4b48f0f86a7a5239b8`.
- Inspection script: [inspect_animation_readiness_032.py](inspect_animation_readiness_032.py),
  SHA-256
  `27585239b8467d3dac7ed6933d3f2a9f2afdae45ac1134eafd137747029c5401`.
- Machine evidence: [animation_readiness_032.json](animation_readiness_032.json),
  SHA-256
  `a7a3a77dcdcd1718e5f8801428084ebb3cc7fb8851657196b746933e461c9bab`.
  No bone, action-channel, or attachment arrays were truncated.
- Authority: `projects/renders/goals/reimu-fumo-finish/criteria.yaml`, criteria
  revision 4, especially criterion-005 r2 and criterion-006 r2; owning asset
  README and PROCESS distinguish module retention from stage acceptance.

This was one metadata inspection. It did not evaluate animated geometry,
append into a blank file, inspect every modifier/dependency type, or audit
resources, manifoldness, contact, or final pixels. The coordinator retains
all model and canonical-goal write authority; no work was delegated further.

## Rig and action inventory

`ReimuFumoRig` is the only armature. It is linked directly to `Scene
Collection`, has identity object transforms, and is in `POSE` mode. All nine
bones are deforming bones using XYZ Euler rotation, with zero current local
translation/rotation, unit scale, and no pose constraints. There are no rig
drivers.

| Bone | Parent | Observed role in render-enabled mesh groups |
| --- | --- | --- |
| `root` | none | Parent of every other bone |
| `Body` | `root` | Seven torso, collar, cravat, skirt, and hem meshes |
| `Head` | `root` | Fifteen current head, hair, and facial meshes |
| `Hair` | `root` | No render-enabled mesh group in this inspection |
| `Arm_L` | `root` | Two sleeve panels and one hand insert |
| `Arm_R` | `root` | Two sleeve panels and one hand insert |
| `Leg_L` | `root` | Foot pod and cream root |
| `Leg_R` | `root` | Foot pod and cream root |
| `Bow` | `root` | Nine bow, tail, knot, and ruffle meshes |

The only action is `ReimuFumoNeutral.002`, assigned to the rig with one user
and no fake user. It has zero F-curves, zero action slots, and frame range
`0..0`. The rig has no NLA tracks or strips. No inspected geometry object has
its own active action.

The scene is `Attempt41_Manual_Head_Maquette`, saved at frame 1, with playback
range `1..250` at 24 fps. That range is a scene setting, not evidence of an
animation. There is no stored head-yaw, arm-wave, or combined validation
action, and the named neutral action contains no keyed neutral pose.

## Attachment findings and blockers

The inventory covers 191 mesh/curve objects, including 105 marked
`hide_render`. Of 95 objects with armature modifiers, none has a base-mesh
vertex without a positive weight. Forty-one are render-enabled meshes bound
to the groups listed above. This is useful attachment inventory, not proof
of correct pivots, attachment compensation, or deformation.

1. **Required action coverage is absent.** Criterion-006 cannot be exercised
   from this file's stored actions. A clean open and zero rest transforms
   cannot substitute for its seated-neutral, actual head-yaw, arm-wave, and
   combined pose-render evidence.
2. **Head and bow have independent controls.** `Head`, `Hair`, and `Bow`
   are siblings under `root`; there is no inherited `Head` to `Bow`
   relationship or pose constraint. The current visible head/hair meshes
   follow `Head`, while the bow meshes follow `Bow`. A head-only turn does
   not provide a hierarchical bow-follow mechanism. Future stored actions
   need a deliberate head/bow relationship and contact validation.
3. **Actual yaw uses the bone's local Y axis at this neutral transform.**
   `Head` has local X = world +X, local Y = world +Z, and local Z = world -Y.
   With the identity armature transform, a local-Z Euler rotation would not
   be a turn around the world's vertical axis. This is a metadata inference;
   no rotation was applied.
4. **Forty-four render-enabled decoration curves have no direct rig route
   in the inspected parent/armature/constraint channels.** These comprise
   six bow root-fold/zigzag curves and 38 sleeve cuff, seam, stitch, and pleat
   curves. They have no parent, armature modifier, or object constraint.
   Representative names are `A42 Left white zigzag applique` and
   `Sleeve44P L front red running stitch 1`. Other modifier dependency routes
   were not inventoried, so stationary behavior or detachment is a risk to
   verify, not a proven pose failure. The unbound review-floor mesh is
   intentional review content, not character rigging.
5. **Reusable collection packaging is absent.** The 13 named datablock
   collections are separate historical/current modules. None contains an
   armature; the rig is only in the scene master collection. The file retains
   hidden historical geometry, a deformer collection, and a review floor.
   There is no observed named reusable collection that contains the complete
   current character and rig while excluding review content. The required
   append-only blank-file test was not run.

## Criterion disposition and next evidence

| Criterion | Disposition for head032 |
| --- | --- |
| 005 r2: reusable structure | Fails observed collection packaging; append and dependency behavior unverified |
| 006 r2: animation readiness | Fails stored-action coverage; tearing/clipping and pose quality unverified |
| 007 r2: technical integrity | Unverified by this bounded inspection |

After the model reaches the applicable visual stage, the next relevant work
is to establish one intentional reusable collection, resolve decoration and
head/bow attachment routes, and store the four required validation actions.
The same final bytes then need action-driven fixed-view pose renders and a
blank-file append/reopen test. This inspection does not authorize or perform
that later work and does not assess candidates 033 or 034.
