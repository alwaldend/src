# A99 plan — terminal direct-open pointer coupon

## Binding

- Goal: `projects/renders/goals/reimu-fumo`
- Goal generation: 1
- Lifecycle generation: 5
- Criteria revision: 1
- Expected resource version at start: 73
- Work type: investigation
- Exact input:
  `out/reimu_fumo_attempt_094_real_pointer_sculpt_coupon/author/plus_y/`
  `snapshots/plus_y_before.blend`
- Required SHA-256:
  `02dd81b24a23a135462044c8b15a7498f743442f71d4de05ae21dae8ba9a1331`
- Target: `A94_Y_LiveGrab_Coupon`

A99 is the one terminal correction authorized after A98 failed before `READY`
because deferred `open_mainfile(load_ui=False)` invalidated
`bpy.context.window`. A98 remains immutable and closed.

## Only permitted implementation delta

Copy the frozen A98 harness into the A99 output root, then:

1. pass the exact task-local working Blend positionally to Blender before
   `--python`; and
2. remove the deferred `bpy.ops.wm.open_mainfile` block, retaining only exact
   loaded-filepath and byte-hash assertions.

No other executable behavior may change. Keep pinned Blender 5.2.1, private
Xvfb, invalid Wayland endpoint, Grab `SCENE / 0.050 m / 0.40`, external
same-window X11 capture, plain timers, exact fingerprints, one-shot latches,
one `+Y` 5 mm XTest drag with eight motions over 320 ms, fixed renders, metrics,
timeouts, and stop behavior identical to A98. Add no instrumentation,
recorder, alternate axis, repair path, or model generator.

## Gates

1. Independently inspect the exact A98-to-A99 source diff and source hashes.
2. Launch exactly once. A launch/readiness failure closes A99 without patch or
   relaunch.
3. Require `READY.json`, a valid projected plan, and an externally captured,
   visually inspected complete nonblack Sculpting viewport before input.
4. If and only if the pre-input gate passes, deliver the sole unchanged drag.
5. Require post Blend, geometry metrics, fixed right/three-quarter renders,
   same-window external after capture, independent clean-reopen verification,
   and the exact A98 numeric/visual gates.
6. Stop all task-owned processes after success or failure.

Any launch, readiness, capture, input, post-evidence, timeout, numerical, or
visual failure ends A99. There is no A100 harness. Failure blocks the main goal
pending a stable external live-pointer/sculpt interface or human DCC authoring.
A pass is diagnostic only and leads immediately to a model-bearing incremental
native sculpt on exact rung003; criteria 001–008 remain unverified.

