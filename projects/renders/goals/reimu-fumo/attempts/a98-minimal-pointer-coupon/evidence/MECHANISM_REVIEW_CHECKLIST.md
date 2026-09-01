# A98 pre-input mechanism review checklist

The coordinator completes the static section after the harness is frozen and
before starting Xvfb/Blender. The runtime section is completed after the live
baseline is inspected and before input. Crossing the input boundary with any
unchecked item is a reset.

- [x] Exact snapshot is copied, not opened for overwrite, and its source hash
      is `02dd81b24a23a135462044c8b15a7498f743442f71d4de05ae21dae8ba9a1331`.
- [x] Target object is exactly `A94_Y_LiveGrab_Coupon`.
- [x] Harness contains no `bpy.ops.sculpt.brush_stroke` call.
- [x] Harness registers no modal Blender operator or event sentinel.
- [x] No XRecord, alternate axis, retry loop, parameter sweep, or automatic
      repair path exists.
- [x] Pinned Blender 5.2.1 is the only Blender executable.
- [x] Private Xvfb is configured for a task-owned display and process receipt.
- [x] Grab is explicitly configured and attested at 50 mm world radius and
      strength 0.40; symmetry, Dyntopo, remesh, texture,
      and pressure variation are off.
- [x] XTest plan contains one press, eight held motions over 320 ms, one
      release, 5 mm world `+Y`, and an irreversible one-shot latch.
- [x] Plain timer watches only the injector-complete marker and produces post
      coordinates, immutable Blend, fixed right/3Q renders, and `DONE`.
- [x] External capture targets the mapped Blender X11 window; it does not call
      Blender's screenshot operator.
- [x] Readiness, pre-input pixel inspection, post timeout, and process-stop
      paths are bounded and task-local.
- [x] Source/rung003/tracked hashes were rechecked before execution.

## Runtime pre-input gate

- [ ] `READY.json` attests exact Blender, input, target, topology, source,
      brush, and one projected `+Y` gesture plan.
- [ ] The settled view proves world `+Y` projects visibly and world `+X` is
      the depth axis.
- [ ] The external mapped-window capture is nonblack and visibly shows the
      complete coupon in the expected Sculpting viewport.
- [ ] The mapped X11 window ID and dimensions agree between the ready plan and
      external capture.
- [ ] No input latch, injector receipt, or injector-complete marker exists.
