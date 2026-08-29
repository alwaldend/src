# Reimu Fumo acceptance and verification

[Back to goal](README.md)

## Acceptance criteria

1. **Reference likeness:** fixed front, side, rear, and three-quarter renders
   preserve the characteristic broad head, compact seated body, layered brown
   hair, asymmetric red-and-white bow, half-lidded red eyes, detached sleeves,
   gathered red skirt, yellow necktie, and small black plush feet.
2. **Measured silhouette:** every controlled front landmark error is at most
   `0.05 Wh`, and every depth landmark lies inside its frozen uncertainty band,
   where `Wh` is the calibrated head width. The landmark set, normalization,
   alignment, uncertainty, and exclusions are frozen in `LANDMARKS.md` before
   the first candidate; there are no post-hoc exceptions.
3. **Plush construction:** visible forms read as stuffed fabric panels with
   plausible thickness, compression, seams, gathers, ruffles, and contact;
   there is no helmet, cone, mask, robot, hard-plastic, or disconnected-
   primitive read in any required view.
4. **Presentation quality:** two implementation-blind reviews independently
   identify the same Reimu Fumo variant without a label, give every fidelity
   category at least 8/10, and report no major visible defect.
5. **Reusable structure:** the `.blend` has clearly named collections, meshes,
   materials, armature, actions, and controls; transforms are intentional; no
   required resource is missing; reusable content is separated from review
   cameras, lights, and reference images.
6. **Animation readiness:** an armature controls the head, body, arms, legs,
   hair, and bow; weights or explicit parent controls deform without visible
   tearing or clipping in a seated pose, a head turn, an arm wave, and a short
   validation animation.
7. **Technical integrity:** meshes have no unintended non-manifold boundaries,
   degenerate geometry, duplicate object names, broken modifiers, or dependency
   cycles; textures are packed or deliberately procedural.
8. **Repository delivery:** the reusable asset is stored at
   `projects/renders/blender/fumo/reimu_fumo/reimu_fumo.blend`, tracked by Git
   LFS, and exposed as `//projects/renders:reimu_fumo`.
## Evidence plan

- Maintain the frozen measured contract in `LANDMARKS.md`; keep generated
  silhouette overlays under `out/`.
- Freeze the landmark list, normalization method, camera calibration, alignment
  calculation, uncertainty, and exclusions before the first model candidate.
  Do not remove a controlled landmark after seeing a candidate.
- Render the same orthographic front, side, rear, and three-quarter cameras for
  every candidate, plus a perspective presentation render.
- Give each implementation-blind reviewer only the frozen raw references,
  candidate renders, and linked scorecard. Do not disclose history, object
  names, hypotheses, previous scores, or implementation details. Both reviewers
  must pass every category independently.
- Inspect the final Blender data model, object bounds, modifiers, materials,
  packed resources, armature, actions, and mesh integrity through Blender MCP.
- Require named actions `Fumo_Seated` frames 1–24, `Fumo_HeadTurn` frames 1–48,
  `Fumo_Wave` frames 1–64, and `Fumo_Validation` frames 1–120. Evaluate every
  frame of every action for tearing, skirt/foot clipping, hair/head collision,
  bow deformation, and unintended intersections, then render fixed-view
  contact sheets and a short animation.
- Prove reuse by appending only the reusable collection into a new blank
  Blender file, playing every action, and rendering it without the source
  scene, review collection, missing paths, or name collisions.
- Validate the Bazel target, Git LFS tracking, protected-scene checksum, and
  repository diff before delivery.
- Freeze the final candidate, rerun visual, deformation, integrity,
  clean-import, packed-resource, Bazel, LFS, and protected-scene checks, then
  verify the committed LFS object hash exactly matches the tested candidate
  hash.
## Fixed regression set

Every attempt must check its targeted criteria plus these recurring defects:

- helmet-like, fat, tall, or excessively deep head;
- overly human body or limb proportions;
- flat, disconnected, slab-like, cone-like, or hard-plastic construction;
- incorrect hair silhouette, layering, or side/rear read;
- bald rear head or a partially dressed front-hair/bald-rear review proxy;
- rigid, wing-like, planar, or block-mounted bow;
- tall torso or dress instead of a compact seated plush body;
- foot/skirt, body/terrain, or accessory clipping;
- missing feet, human anatomical feet, or tangent ball feet;
- weak side or rear silhouette;
- loss of the half-lidded Reimu expression; and
- bland lighting or materials that conceal construction and contact.
## Phase gates

1. **Dossier and pattern gate:** freeze named outlines, controlling sources,
   camera contracts, exclusions, and uncertainty before geometry.
2. **Black silhouette gate:** review only the head cushion, compact torso,
   skirt envelope, sleeve envelopes, feet, and bow pockets. Reject any egg,
   cone, tube, wing, tall-body, or human-proportion read before adding identity
   detail.
3. **Neutral sewn-construction gate:** require seams and panel transitions to
   affect tension and silhouette; require every component root to be physically
   seated; run the four-view blind absolute-quality gate.
4. **Identity gate:** add the continuous fringe, cheek locks, eyes, bow trim,
   ruffle, and garment graphics without using color or fuzz to rescue form.
5. **Rig and contact gate:** at neutral clay, test head yaw and squash, body
   lean, sleeve swing and cuff twist, foot movement, bow curl, and skirt
   spread.
   Reject seam opening, root floating, clipping, volume collapse, inversion,
   applique drift, or identity loss.
6. **Material gate:** add final fabric response and surface construction, then
   rerun the visual, deformation, and contact regressions.
7. **Packaging gate:** reopen the saved candidate, append only its asset
   collection into a blank file, play its actions, render it, and run the final
   integrated regression on the exact bytes intended for delivery.
