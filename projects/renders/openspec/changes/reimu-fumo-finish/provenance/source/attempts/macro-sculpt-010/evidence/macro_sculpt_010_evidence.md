# Macro sculpt capability receipt

## Configuration

- Authoring host: installed Blender 5.1.1 Flatpak, forced to the existing
  XWayland display for a mapped `5120x2880` window.
- Bridge: existing Blender MCP add-on, loopback port 9876, started by a
  deferred timer after the GUI mapped.
- Blender safety gate: `--enable-event-simulate`.
- Brush asset:
  `brushes/essentials_brushes-mesh_sculpt.blend/Brush/Grab` from the
  `ESSENTIALS` library. Runtime checks reported brush `Grab`, type `GRAB`, and
  workspace tool `builtin.brush`.
- Effective radius owner:
  `tool_settings.sculpt.unified_paint_settings`, with unified size enabled,
  lock mode `SCENE`, and unprojected diameter `2.0`.
- Input: one projected right-silhouette target followed by nine native events
  over timer ticks 0.12 seconds apart: initial move, button press, six outward
  moves, and button release.

## Geometry result

The target was projected from source vertex 1952 at world coordinate
`(1.9566307, 0.3921607, 0.1334401)` to window coordinate `(3181, 1418)`.
The event sequence ended at `(3541, 1418)`.

- Vertices changed: 161 of 2,562 (6.284%).
- Maximum displacement: `0.6282851620288168`.
- Mean displacement among changed vertices: `0.18258674961513435`.
- X silhouette bounds: source `[-1.9999999, 2.0]`; candidate
  `[-1.9999999, 2.5007470]`.
- Native simulated `Ctrl-Z` after saving restored changed vertex count `0` and
  maximum delta `0.0` in the live session.

## Exact artifacts

| Artifact | SHA-256 |
| --- | --- |
| `out/reimu_fumo_finish/live_blender/outputs/sculpt_coupon_source.blend` | `23af3ada0a3145aa20bab03b2686144a978c42203049ab2fba14f9a6955cdb99` |
| `out/reimu_fumo_finish/live_blender/outputs/sculpt_coupon_event_grab_silhouette.blend` | `45ab2fae5bbb2112e3432391f199f0cc4fc494a853eacf50e104e42c1daaf14d` |
| `out/reimu_fumo_finish/live_blender/outputs/sculpt_coupon_source_pinned.png` | `4f0e72cba28fb0834536df09b0756eb94ff7d132cbed1a3d2996c5e09e570578` |
| `out/reimu_fumo_finish/live_blender/outputs/sculpt_coupon_grab_silhouette_pinned.png` | `e67c2e968a3389f43db39d4825e943edc94b3a3c32c9f42a3112e4ac2ecc21a2` |

Repository-pinned Blender 5.2.1 clean-opened both files. Candidate and source
geometry digests were respectively
`b7c1e11f8348b05fde6fa132b9dadb953e5cfb5644af1c84712b7b04b6a2f855`
and
`d297854f139a00951ffc4d96d4808f8e951a82fee09d4773886acaffbe4536`.
Both retained 2,562 vertices and 5,120 polygons. The pinned renders visibly
show the localized right-side pull and have normalized RMSE `0.0863399`.

The source blend hash remains identical to the prior attempt's recorded hash.
No Reimu asset was opened or modified.

## Diagnosis trail

Each correction changed an evidenced mechanism:

1. XWayland replaced the unmapped zero-size native-Wayland window.
2. Deferred bridge startup replaced a premature pre-map start.
3. Explicit port 9876 replaced the saved profile's port 9985.
4. `--enable-event-simulate` satisfied Blender's explicit safety gate.
5. The unified paint settings replaced ineffective brush-datablock radius
   edits.
6. A projected silhouette target replaced an internally displaced center
   patch that could not prove visible outline control.

