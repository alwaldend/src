# Capability receipt

## Diagnosed launch failures

- Repository-pinned Blender 5.2.1 accepted the MCP bridge but exposed one
  `0x0` window. Window screenshot capture therefore failed before input.
- Calling `bpy.ops.wm.window_new_main()` from that unmapped context crashed
  Blender. The ignored crash receipt has SHA-256
  `af29820ec175aa24361648eb5abc3840881e13b271c69ed17806c41db896ca28`.
- Flatpak Blender 5.1.1 on native Wayland also exposed `0x0`; forcing only
  this task-owned process to XWayland produced a `5120x2880` window with a
  `4204x2620`-class mapped `VIEW_3D` region.
- The first mapped input request was rejected explicitly because Blender was
  not started with `--enable-event-simulate`. Relaunching with that safety
  flag resolved the rejection.

## Native event result

Blender's active Sculpt-mode Draw brush received one
`MOUSEMOVE`, `LEFTMOUSE PRESS`, four further `MOUSEMOVE` events, and
`LEFTMOUSE RELEASE` through `Window.event_simulate`. The stroke changed 3 of
2,562 vertices. Maximum displacement was `0.028769294878297935`; mean
displacement among changed vertices was `0.022315126317176134`.

The changed file was saved before recovery. A native simulated `Ctrl-Z` then
restored all 2,562 live coordinates to the recorded baseline: changed vertex
count `0`, maximum delta `0.0`.

## Exact artifacts

| Artifact | SHA-256 |
| --- | --- |
| `out/reimu_fumo_finish/live_blender/outputs/sculpt_coupon_source.blend` | `23af3ada0a3145aa20bab03b2686144a978c42203049ab2fba14f9a6955cdb99` |
| `out/reimu_fumo_finish/live_blender/outputs/sculpt_coupon_event_draw.blend` | `a75599a35ba272477943373c131f29add0883907188620f57c8fcdc577e7469a` |
| `out/reimu_fumo_finish/live_blender/outputs/sculpt_coupon_pinned.png` | `5ad5d935eebf3721123344e251428eaecd96dd86dab89a6f2864b24f3e8f202e` |

Repository-pinned Blender 5.2.1 clean-opened the changed file, rendered the
512-by-512 PNG, and then clean-opened the source. Its coordinate digests were
respectively
`7c5df5f40cb7b3c221b7d6e6579b3c0034341c70f170c6151f13bd5cab293b2e`
and
`d297854f139a00951ffc4d96d4808f8e951a82fee09d4773886acaffbe4536`.
Both files retained 2,562 vertices and 5,120 polygons. The render contained
130,661 bytes.

## Interpretation

This proves mapped, event-driven sculpt delivery, distinct saved bytes,
native undo, and pinned clean-open rendering. It does not yet prove useful
organic authoring: the three-vertex deformation is not visibly reviewable in
the fixed render. No Reimu asset was opened or modified.

