# Attempt 017 Blender-writer result

Status: failed at the first shape gate and stopped before root or tip work.

## Frozen fixture

- Live authoring host: Blender 5.1.1, foreground, existing loopback MCP.
- The setup did not set a render engine and did not access or create a World.
- Connected fixture: 9,506 vertices, 9,504 polygons, zero loose edges.
- Named regions: plane, root, tip, mid-panel, receiver witness, three
  operation-specific controls, and width witnesses.
- Corrected coordinate digest:
  `41ee23670f67335ac070d95bd782436f53405034f1e24efdebdad709f7d47df2`.
- Embedded baseline-coordinate maximum delta from live mesh: `0.0`.
- Frozen baseline:
  `out/reimu_fumo_finish/attempt_017_localized_sculpt_coupon/localized_sculpt_fixture_baseline_frozen.blend`.
- Frozen baseline SHA-256:
  `6f49bd4e0a8af6b45870d9d4224a520c1398e52e6e9c42f8fb5bee7b8c17118e`.

The earlier file `localized_sculpt_fixture_baseline.blend` remains as setup
evidence at SHA-256
`313f392bf295109584eae28fdeb2e4a39273c63c31fa5f050304277f9bc82f55`.
Its stored digest described the pre-import Python doubles, not Blender's
float32 mesh coordinates. Geometry did not change; the corrected baseline was
saved to a new path after deriving the digest from the actual mesh.

## Brush and settings receipt

All four assets came from
`essentials_brushes-mesh_sculpt.blend`, were activated through Blender's
native asset operator, and read back with scene-locked size:

| Asset | Blender 5.1.1 type |
| --- | --- |
| `Grab` | `GRAB` |
| `Smooth` | `SMOOTH` |
| `Flatten/Contrast` | `PLANE` |
| `Clay Strips` | `CLAY_STRIPS` |

- Unified scene size: `0.5799999833106995`.
- High strength: `0.6499999761581421`.
- Low strength: `0.3499999940395355`.
- Predicted response: high-strength primary strokes should exceed
  low-strength smoothing displacement.
- Blender 5.1.1 uses `sculpt_brush_type` and `unprojected_size`; setting pixel
  size before scene size prevents the pixel setter from overwriting the
  scene-space value. Settled activations were not repeated.

## First native shape operation

- Fixed view: front orthographic, distance `3.9000000953674316`, rotation
  quaternion `(0.7071067690849304, 0.7071067690849304, 0, 0)`.
- Isolation: 1,548 plane vertices unmasked and 7,958 vertices masked.
- Brush: `Flatten/Contrast`, type `PLANE`, strength `0.65`, scene size `0.58`.
- Events: 33 delivered with no error: three hover moves, three presses,
  24 drag moves, and three releases.
- Each stroke contained eight moves after its press; actual inter-event
  spacing ranged from `0.8655732629995327` to
  `1.1363116559950868` seconds.
- Changed plane-target vertices: 958.
- Changed vertices outside the target: 0.
- Maximum named control displacement: `0.0`.
- Maximum displacement: `0.03233039258016114`.

## Failed plane gate

- Plane variance: `0.014215173048885853` to
  `0.01352948526060966`, a `4.8236330709314545%` reduction.
- Pillow-height range: `0.45792272686958313` to
  `0.4257808327674866`, a `7.019065055325502%` reduction.
- Required reduction: at least `35%` in either accepted metric.

The event transport, brush activation, and isolation worked, but the planned
Flatten operation lacked the required magnitude. This is a partial sculpt
result rather than a transport or configuration failure. The stop rule barred
another stroke and barred root, tip, and three-quarter operations.

## Exact partial artifact and native undo

- Partial failure snapshot:
  `out/reimu_fumo_finish/attempt_017_localized_sculpt_coupon/localized_sculpt_coupon_017_partial_plane_fail.blend`.
- Snapshot SHA-256:
  `2428de0a0b65e572de9437a8d3ef35f1ee21c18bd9dbf27ff01de1816418c0bd`.
- The snapshot embeds the exact 33-event log, brush receipts, settings
  receipt, and plane metrics as object custom properties.
- After saving it, three native simulated Ctrl-Z press/release pairs were
  delivered without error.
- Live coordinate digest after undo:
  `41ee23670f67335ac070d95bd782436f53405034f1e24efdebdad709f7d47df2`.
- Baseline-digest match: true.
- Changed vertices after undo: 0; maximum delta: `0.0`.

The live scene now contains the exact baseline geometry in memory while the
saved partial artifact retains the failed geometry. No Reimu blend was opened,
saved, or modified during this attempt. Pinned reopen and rendering were left
to the coordinator as required.

