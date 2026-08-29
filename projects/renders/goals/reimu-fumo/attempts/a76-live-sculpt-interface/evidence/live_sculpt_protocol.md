# A76 live `VIEW_3D` sculpt-interface protocol

## Verdict before execution

Do not run another `bpy.ops.sculpt.brush_stroke(stroke=...)` probe.  A live
foreground Blender process and a valid `VIEW_3D` override do not make that
call an interactive gesture; it is the same synthetic operator path already
falsified by A68 and A71.

The smallest discriminating test is one real pointer-driven Grab drag on a
deep copy of the exact A71 uniform-remesh coupon.  A temporary pass-through
modal event sentinel must observe the `LEFTMOUSE` press, native mouse motion,
and release.  This distinguishes an actual viewport interaction from a Python
stroke list.  Numeric support and matching before/after viewport pixels then
establish whether that interaction is artistically useful.

The Blender MCP capabilities exposed on 2026-08-31 can prepare and inspect the
live scene, execute `bpy`, focus a `VIEW_3D`, capture area/window screenshots,
render to a path, and read bundled API/manual documentation.  They expose no
mouse press/move/release or general desktop-input tool.  The connected session
was Blender 5.1.1, foreground, clean, in the Sculpting workspace, with a
3150-by-2049 `VIEW_3D`; it had the A72 output open in Object mode with no
active object.  Therefore the present MCP surface alone cannot execute this
discriminating probe.  `execute_blender_code` may prepare and measure it, but
must not impersonate the gesture.  Proceed only when a human or a purpose-
built input bridge can deliver genuine pointer events to that Blender window.

This is a tooling/interface gate, not a visual-fidelity cycle and not approval
of sculpting as the model strategy.

## Exact control and isolation

Use the exact A71 raw P0 artifact instead of the protected whole-plush rung.
That changes only the input mechanism relative to the failed probe and avoids
opening any protected asset in the mutation host.

- Control file:
  `out/reimu_fumo_attempt_071_macro_sculpt/p0_raw/`
  `reimu_fumo_a71_p0_raw.blend`
- Required control SHA-256:
  `0004a3f0bc4987a250f7028b8697a7f740c70866ad8b60f7181e2b4eafa96400`
- Control object: `A71_MacroClay`
- Required control geometry: 3,863 vertices, 3,861 polygons, approximately
  `130.9 x 101.3 x 129.1 mm`
- Prior ordinary-Grab result: 11/3,863 vertices changed (`0.285%`), about a
  10 mm support patch, maximum displacement 1.80 mm.
- Prior Elastic-Grab result: only 40 vertices moved more than 0.10 mm and only
  six moved more than 0.50 mm, despite epsilon changes on all vertices.

Before launching Blender, make a byte copy at:

`out/reimu_fumo_attempt_076_live_sculpt/live/session_working.blend`

Record hashes for both paths and require that the source hash remains exact.
The live Blender may open only `session_working.blend`; it must never open the
control path for saving.  In that working file:

1. Deep-copy `A71_MacroClay` to `A76_LiveGrab_Coupon` with both
   `object.copy()` and `object.data.copy()`.
2. Link the copy to a task-owned `A76_LiveProbe` collection.
3. Hide and disable selection on `A71_MacroClay`.  Keep its coordinate digest
   as an in-file control.
4. Require `A76_LiveGrab_Coupon.data.users == 1`, an identity world matrix,
   applied scale, no modifiers, no shape keys, and the required vertex and
   polygon counts.
5. Fingerprint every non-target object's transforms, mesh coordinates, and
   visibility before the gesture.  No non-target change is permitted.

The protected inputs remain read-only throughout:

- whole-subject rung SHA-256:
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- tracked reusable asset SHA-256:
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`

Rehash both after the live phase.  A mismatch is an immediate hard failure.

## Four-phase settled-state sequence

Use separate MCP calls for every phase.  Never retain an area, region, or
`RegionView3D` reference across a workspace or file transition.

### 1. Open and identify

Open the absolute working-copy path with UI loading disabled.  In a later
read-only call, require all of the following:

- Blender 5.1.1, `bpy.app.background is False`;
- `bpy.data.filepath` is the absolute A76 working path;
- `bpy.data.is_dirty is False` before preparation;
- the Sculpting workspace exists;
- exactly one foreground window and a usable `VIEW_3D` `WINDOW` region exist;
- the control and target objects satisfy the isolation checks above.

If any identity differs, stop before entering Sculpt mode.

### 2. Prepare the exact viewport and brush

In one explicit `temp_override`, set:

- workspace: `Sculpting`;
- active and sole selected object: `A76_LiveGrab_Coupon`;
- mode: `SCULPT`;
- view: right orthographic via `view3d.view_axis(type='RIGHT')` followed by
  `view3d.view_selected(use_all_regions=False)`;
- `RegionView3D.view_distance = 0.28`;
- solid shading, overlays enabled, no clipping plane;
- X, Y, and Z sculpt symmetry all disabled;
- Essentials brush asset: `Grab`;
- locked size: `SCENE`;
- unprojected radius: `0.050 m`;
- strength: `0.40`;
- normal direction, no automasking, no dyntopo, no remesh, no texture, and no
  pressure-dependent radius or strength.

Redraw twice, return from the MCP call, and let the UI settle.  In a later
call reconstruct the override and require:

- area type `VIEW_3D`, region type `WINDOW`, orthographic right view;
- the target is active, selected, visible, single-user, and in Sculpt mode;
- `bpy.ops.sculpt.brush_stroke.poll()` is true;
- the active brush name, radius, strength, symmetry, and tool settings match
  the values above.

`poll() == True` proves only readiness.  It is not an effect gate.

### 3. Arm a real-event sentinel and capture the baseline

Register a temporary modal operator only in live memory.  It must return
`PASS_THROUGH` so Blender's normal sculpt keymap still receives the events.
While armed it records, with monotonic timestamps:

- event `type` and `value`;
- window and region mouse coordinates;
- tablet pressure when available;
- the target window, screen, area, and region identities.

It records only the interval from the first in-region `LEFTMOUSE/PRESS` through
the matching `LEFTMOUSE/RELEASE`.  A direct
`bpy.ops.sculpt.brush_stroke(stroke=...)` call does not create this raw event
sequence and therefore cannot pass the provenance gate.

Before arming, store the target's local and world coordinates and coordinate
digest.  Select a surface point using the same deterministic A71 rule:

```text
max(x - 5 * abs(y - 0.018) - 5 * abs(z - 0.171))
```

Let that world point be `P0`.  Let `P1 = P0 + (0, 0.005, 0)` meters.  Project
both points into the settled region using
`location_3d_to_region_2d`.  Require both to lie at least one 50 mm brush
radius inside the viewport.  Persist the region-local and window-local pixel
coordinates in `gesture_plan.json`; an input bridge using top-left screen
coordinates must explicitly convert Blender's bottom-left Y coordinate and
record the window origin and display scale.

Capture before evidence twice:

1. call the focused MCP `get_screenshot_of_area_as_image(VIEW_3D)` so the
   screenshot is attached to the session log;
2. in the exact area override, run `bpy.ops.screen.screenshot_area` to save
   `evidence/before_view3d.png`.

Also save a clean immutable pre-stroke snapshot at
`snapshots/p0_before.blend`, hash it, and do not render or inspect that file
from another Blender process while the live save is active.

### 4. Deliver exactly one genuine viewport drag

No Python or MCP code may call `sculpt.brush_stroke` in this phase.  The input
provider must target the verified Blender window and deliver:

1. pointer move to the projected `P0`;
2. one left-button press;
3. at least six monotonic motion events along the straight segment to the
   projected `P1`, over 0.20--0.60 seconds;
4. one left-button release at `P1`.

Use an ordinary mouse-pressure value of one.  Do not add a second stroke,
hold Shift, change radius/strength, rotate the view, or repair the result.

The modal sentinel must observe exactly one in-region press, at least six
motion events while held, and exactly one release.  Start and end pixels must
be within four pixels of the planned coordinates; path length must be within
10% of the plan, and at least 90% of motion projected onto the planned screen
direction must be positive.  Otherwise classify the input provenance as
invalid and stop, even if vertices changed.

## Quantitative effect gate

Measure after the release in a later MCP call, after a dependency-graph
update.  Use the stored pre-stroke world coordinates.  Report every threshold,
not just a single `pass` field.

Define:

- `material`: displacement greater than `0.00010 m` (0.10 mm);
- `strong`: displacement greater than `0.00050 m` (0.50 mm);
- `d = normalize(P1 - P0)`, the expected positive world `+Y` direction;
- support extents from the pre-stroke world coordinates of `material` vertices;
- adjacency from the unchanged mesh edge graph.

All hard gates must pass:

1. At least 8% of all vertices are `material` (at least 310 of 3,863).
2. At least 1% are `strong` (at least 39 of 3,863).
3. At least 90% of `material` vertices belong to one connected component.
4. The material support spans at least 30 mm in both world Y and world Z.
   This rejects A71's approximately 8-by-10 mm local patch even if the count
   is inflated by epsilon noise.
5. At least 80% of `material` vertices have `delta dot d > 0`, and the mean
   signed displacement along `d` is positive.
6. Mean absolute displacement orthogonal to `d` is no more than 35% of the
   mean positive displacement along `d`.
7. Maximum displacement is greater than 0.75 mm and strictly below 6.0 mm.
8. Vertices whose pre-stroke right-view distance from `P0` exceeds 60 mm have
   no displacement above 0.10 mm.  This rejects an accidental global edit.
9. Vertex count, polygon count, edge connectivity, and finite coordinates are
   unchanged; no zero-area or inverted face is introduced.
10. The hidden source object's coordinate digest and every non-target
    fingerprint remain exact.

The old A71 ordinary result fails gates 1 and 4.  The Elastic result fails
gates 1 and 2 when epsilon changes are excluded.  A reported `FINISHED`
operator result is deliberately absent from the gate.

## Pixel and snapshot evidence

Without changing the view, shading, selection, or tool, capture:

- a second focused MCP `VIEW_3D` screenshot for the session log;
- `evidence/after_view3d.png` through `screen.screenshot_area`;
- `evidence/event_receipt.json` from the modal sentinel;
- `evidence/displacement_metrics.json` with counts, quantiles, support bounds,
  connected-component sizes, direction ratios, hashes, and every gate;
- `evidence/window_state_before.json` and
  `evidence/window_state_after.json`;
- a same-size absolute pixel difference image and its changed-pixel fraction.

The before and after screenshots must show the same right orthographic crop.
The visible deformation must occupy the same broad region reported by the
material-support mask.  A blank, shifted, differently shaded, or cropped pair
fails.  Pixel change alone cannot pass without raw input provenance and
geometry support; geometry support cannot pass without matching pixels.

After evidence capture, save exactly one immutable live result:

`snapshots/p0_after_live.blend`

Save it whether the effect gate passes or fails, label the verdict in a scene
custom property, hash it, and perform no further mutation of that snapshot.
If and only if all interface gates pass, clean-reopen this exact snapshot with
the repository-pinned Blender 5.2.1 in background mode, with automatic scripts
disabled.  Verify object identities, geometry metrics, missing resources, and
the target/source digests.  Render only the frozen snapshot.  The live 5.1.1
session is never the renderer or deliverable authority.

## Stop conditions and authorization boundary

Stop before the gesture when:

- no genuine pointer-event provider is available;
- the live session, file, source hash, object, version, foreground state,
  viewport, or brush differs from this protocol;
- the candidate is dirty before preparation or is outside the A76 `out/`
  subtree;
- the target mesh is shared with any source object;
- the event sentinel cannot coexist with and pass events to normal viewport
  sculpting.

Stop after the single gesture when:

- the raw press/move/release receipt is absent or invalid;
- any quantitative, topology, isolation, or pixel gate fails;
- before/after evidence disagrees;
- the save or pinned clean-reopen boundary fails;
- any protected input hash changes.

Do not retry with more strength, another brush, more strokes, remeshing, or a
different coupon inside A76.  A failure closes the interface strategy.  A pass
authorizes only one later bounded artistic correction on a fresh disposable
whole-plush candidate; it does not approve that correction, the sculpt, or the
reusable asset.
