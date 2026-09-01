# A94 plan — genuine pointer-driven sculpt-control coupon

## Bound state and immutable inputs

- Goal: `reimu-fumo`, generation 1, lifecycle generation 5.
- Expected goal resource version at start: 63; criteria revision 1.
- Exact A71 source:
  `out/reimu_fumo_attempt_071_macro_sculpt/p0_raw/`
  `reimu_fumo_a71_p0_raw.blend`.
- Verified source SHA-256:
  `0004a3f0bc4987a250f7028b8697a7f740c70866ad8b60f7181e2b4eafa96400`.
- Source object: `A71_MacroClay`, 3,863 vertices.
- Frozen tracked asset SHA-256:
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.
- Frozen rung003 SHA-256:
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.

The governing input protocol is A76's tracked
`evidence/live_sculpt_protocol.md`. A94 changes the one capability A76 lacked:
real X11 pointer events delivered through XTest in a private Xvfb.

## Hypothesis

One real pointer-driven native Grab gesture can move a broad, connected,
low-frequency region of uniform macro-clay in its intended direction without
the tiny support or dimples produced by Python-authored stroke lists. Three
isolated non-collinear passes are required before this interface may author a
whole-plush sculpt.

## Execution

1. Build and resolve repository-pinned Blender 5.2.1. Launch it foreground in
   a task-private Xvfb with all configuration, cache, and temporary paths under
   `out/reimu_fumo_attempt_094_real_pointer_sculpt_coupon/author/`.
2. Make three byte-isolated working copies from the exact A71 source. Deep-copy
   `A71_MacroClay` to a single-user target and fingerprint source, target,
   non-targets, cameras, transforms, topology, and coordinates.
3. In separate settled Blender interactions, activate Sculpt mode and the
   native Essentials Grab brush in a verified `VIEW_3D`: scene-locked 50 mm
   radius, strength 0.40, all symmetry, dyntopo, automasking, textures, and
   pressure-dependent settings off.
4. Register a temporary pass-through modal sentinel that records actual
   `LEFTMOUSE` press, motion, and release events. Blender code may prepare the
   view, project exact endpoints, fingerprint geometry, and measure afterward;
   it must never call `bpy.ops.sculpt.brush_stroke`.
5. A task-local XTest injector delivers exactly one 0.20--0.60 second drag with
   at least six monotonic motion events to the verified region. Test isolated
   intended world directions `+Y`, `-X`, then `-Z`; each branch starts from
   identical source coordinates. Stop after the first failed branch.
6. Freeze pre/post snapshots, raw event receipt, window/area geometry, exact
   displacement metrics, topology/camera hashes, fixed front/right/three-
   quarter renders, pixel differences, manifests, and one comparison sheet.

## Hard gate for every branch

- Exactly one in-region press and release, at least six monotonic held motions,
  endpoints within four pixels, and path length within 10% of the plan.
- At least 8% of vertices move more than 0.10 mm and at least 1% move more than
  0.50 mm; at least 90% of material support is one connected component.
- Material support spans at least 60 mm in both tangent axes.
- At least 80% of material deltas project positively onto the requested world
  direction; mean orthogonal displacement is at most 35% of mean intended
  displacement; maximum displacement is greater than 0.75 mm and below 6 mm.
- Nothing beyond 60 mm from the target moves more than 0.10 mm.
- Vertex/edge/polygon counts, connectivity, finiteness, source/non-target
  digests, object transforms, and camera matrices remain exact.
- A fixed-camera controlling view changes at least 1% of subject pixels by
  more than 1% RGB and shows one coherent broad outline movement without a
  dimple, hard curvature ring, pinch, or global shift.

Across all three passing branches, achieved mean directions must have pairwise
absolute dot products at most 0.25 and dot their requested direction by at
least 0.8.

## Stop and authorization boundary

Invalid event provenance, local/zero/dimpled effect, changed protected state,
or any numerical/pixel failure closes the modality immediately. Do not tune a
second stroke, radius, strength, topology, or camera inside A94. A full pass
authorizes exactly one fresh low-resolution whole-plush manual sculpt; it does
not satisfy a final goal criterion or authorize tracked-asset promotion.

## Workstreams

- One author owns the task-private Blender/Xvfb/injector session and artifacts.
- The coordinator independently binds hashes and assembles the evidence.
- Exactly one fresh reviewer judges frozen before/after pixels without method
  context before the coordinator adjudicates the coupon.
