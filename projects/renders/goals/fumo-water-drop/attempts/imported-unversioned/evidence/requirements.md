# Fumo water-drop requirements and constraints

[Back to goal](README.md)

## Supplied requirements

- Create a durable water-drop animation scaffold subordinate to the main
  Reimu Fumo goal.
- Keep goal records under `projects/renders/goals/fumo_water_drop/` and all
  task-owned generated work under `out/fumo_water_drop_scaffold/`.
- Do not modify a tracked Blender file.
- Use a clearly labeled neutral placeholder and freeze the future Fumo
  collection/rig interface.
- Define falsifiable tank, world, camera, lighting, drop-timing,
  liquid-domain, and inexpensive proxy gates.
- Use Blender MCP for an inspectable scaffold when available. If it is not
  available, stop at a decisive source-grounded preflight artifact.
- Do not claim a final simulation or approved-plush integration.

## Inferred operating requirements

- The approved reusable Fumo remains `0.25 m` tall and uses meters, `Z` up,
  and viewer-facing `-Y`, matching the parent goal's scale contract.
- The shot scene drives the fall through a scene-owned control; it does not
  write shot motion into the reusable asset or its source `.blend`.
- The fall to first water contact is ballistic at `9.81 m/s^2` during the
  cheap proxy. Post-contact motion remains intentionally unsolved until a
  liquid/collision test exists.
- A future liquid bake must use a simplified deforming collider, not the
  render mesh, and must first survive a low-resolution cache test.

## Scope boundary

Attempt 00 may freeze interfaces and create neutral diagnostic evidence. It
does not authorize final fluid settings, foam/spray/mesh polish, materials,
sound, final character performance, or changes to the standalone Fumo.
