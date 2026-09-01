# A98 plan — minimal no-sentinel real-pointer coupon

## Binding

- Goal: `projects/renders/goals/reimu-fumo`
- Goal generation: 1
- Lifecycle generation: 5
- Criteria revision: 1
- Expected resource version at start: 71
- Work type: investigation
- Exact input:
  `out/reimu_fumo_attempt_094_real_pointer_sculpt_coupon/author/plus_y/`
  `snapshots/plus_y_before.blend`
- Required input SHA-256:
  `02dd81b24a23a135462044c8b15a7498f743442f71d4de05ae21dae8ba9a1331`
- Target object: `A94_Y_LiveGrab_Coupon`
- Original A71 source SHA-256:
  `0004a3f0bc4987a250f7028b8697a7f740c70866ad8b60f7181e2b4eafa96400`

A98 is a diagnostic authoring-control coupon, not a likeness candidate. Work
only beneath `out/reimu_fumo_attempt_098_minimal_pointer_coupon/`. The exact
input, rung003, and tracked reusable Blend are read-only.

## Objective

Observe exactly one genuine native Blender Grab drag with valid live pixels,
post-release geometry, a saved Blend, and fixed renders. Remove both A94
failure points: the competing modal sentinel and the internal black screenshot.
Do not build another provenance system.

## Pre-mutation mechanism gate

The harness implementation must be reviewed before launch. It must:

- use repository-pinned Blender 5.2.1 foreground on a private Xvfb;
- open only a byte-identical task-local copy of the exact input;
- register no modal operator and contain no call to
  `bpy.ops.sculpt.brush_stroke`;
- recreate and verify the right-orthographic view, normalized target, native
  Grab brush, 50 mm scene radius, 0.40 strength, no symmetry/Dyntopo/remesh;
- externally capture the mapped X11 window before and after input;
- use the existing XTest delivery mechanism for one recomputed 5 mm world
  `+Y` path, eight held motions, 320 ms;
- use only a plain timer watching an injector-complete marker to measure,
  save, and render after release; and
- enforce a one-drag latch and bounded timeout.

Reject the mechanism before input if it adds a sentinel, XRecord, alternate
axis, retry/tuning loop, synthetic stroke list, or another evidence subsystem.

## Execution and evidence

1. Verify and copy the exact snapshot; freeze source/non-target coordinates,
   topology, matrices, and tool state.
2. Start private Xvfb and pinned foreground Blender.
3. Wait for an explicit readiness receipt and mapped window.
4. Capture the actual X11 window externally. The coordinator must open and
   inspect it; black, blank, wrong-window, or wrong-view pixels stop before
   input.
5. Deliver exactly one authorized XTest drag.
6. Let the non-modal timer observe the completion marker, settle, write post
   metrics, save one immutable post Blend, render fixed right and three-quarter
   images, and signal completion.
7. Capture and inspect the external post window. Stop all task-owned processes.
8. Reopen the post Blend in pinned background Blender and verify the complete
   state.

Launch to inspected before/after pair is capped at 15 minutes. There is no
harness repair or second input after the first launch crosses the input gate.

## Pass gate

- External live before/after captures are valid, non-black, same-window, and
  consistently framed.
- At least 8% of target vertices move more than 0.10 mm; at least 1% move more
  than 0.50 mm; at least 90% of material vertices form one component.
- Displacement is predominantly world `+Y`; maximum is above 0.75 mm and below
  6 mm; vertices outside the 60 mm support remain below 0.10 mm.
- At least 1% of fixed-right subject pixels change as a broad connected
  silhouette arc, not a dimple, ring, pinch, tear, or global affine shift.
- Fixed three-quarter pixels corroborate the same deformation.
- Topology, target matrix, source coordinates, and every non-target fingerprint
  remain exact.

## Stop and consequence

Stop before input on invalid baseline pixels or mechanism noncompliance. After
the sole drag, stop on missing post evidence, zero/local/wrong-direction effect,
failed pixels, topology/source drift, or timeout. Do not repair or retry A98.

A pass authorizes a separate model-bearing incremental native-sculpt attempt in
exact rung003 context; it does not verify any model criterion. A failure marks
current autonomous authoring infeasible. Resumption then requires a stable live
pointer/sculpt capability or human DCC authoring, not another generator.

