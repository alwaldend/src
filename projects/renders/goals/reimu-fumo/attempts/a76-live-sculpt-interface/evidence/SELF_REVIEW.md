# A76 self-review — live sculpt novelty claim rejected

## Verdict

Reset before any sculpt stroke. The connected Blender MCP was inspected and a
single-user disposable probe was prepared, but the available mutation surface
does not provide real pointer press/move/release events. Its only sculpt path
is the same Python-authored `bpy.ops.sculpt.brush_stroke` method already
falsified by A68 and A71.

## Safe setup evidence

- Live Blender: `5.1.1`, with one real `VIEW_3D` area.
- Before-probe file:
  `out/reimu_fumo_attempt_076_live_sculpt/probe/live_probe_before.blend`
- Before-probe SHA-256:
  `d1484be1c9acfa0fc92e0c67abd83e7425eb0a2bd0d0e404e75cec403381ac65`
- Probe: `A76_Sculpt_Interface_Probe`, 11,184 vertices and 11,182 polygons.
- Source/probe coordinate digest at duplication:
  `54b08a9aa81b4a79d4da9de49e68cd406dbf915041af5c897950dfe0b2d193c9`
- Protected rung remained:
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.

The probe lives in its own scene and mesh datablock. No stroke, coordinate
mutation, render, promotion, or tracked-asset change occurred.

## Disconfirming evidence

A68 already used pinned Blender 5.2.1, a foreground `VIEW_3D`, the interactive
MCP bridge, separated settled calls, and the native Essentials Grab brush.
Useful support remained 35 vertices and larger strokes created dimples. A71
repeated a foreground live-context test on fresh uniform topology: ordinary
Grab moved 11 of 3,863 vertices and Elastic Grab moved only 40 vertices more
than 0.10 mm.

The current MCP exposes code execution and screenshots, but no genuine GUI
pointer-drag capability. Repeating `brush_stroke` would change transport state,
not the authoring method. The intended A76 discriminator is therefore
unavailable and the approach/module is wrong before mutation.

## Process decision

The early-return rule worked: capability inspection plus historical review
prevented another fake novelty cycle. The next method is sparse direct
BMesh/Edit-mode proportional macro editing, where selected vertices and
displacements are explicit and predictable. It must remain a small macro-form
probe in full-plush context and must not repeat dense formula fields,
legacy-stack deformation, another receiver loft, or scripted panels.

