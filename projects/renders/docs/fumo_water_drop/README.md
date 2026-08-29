# Fumo water-drop scaffold goal

## Goal

Create a reusable Blender shot scaffold for dropping the approved `0.25 m`
Fumo into water. The scaffold must freeze the asset/rig handoff, tank and
liquid scale, camera, lighting, and shot timing without treating a placeholder
as the character or an unbaked proxy as a finished simulation.

This is a subordinate successor to the
[standalone Reimu Fumo goal](../reimu_fumo/README.md). It may prove the scene
interface in parallel, but it may not approve, replace, or visually normalize
an unfinished Fumo.

## Status

`in progress` — Attempt 00's validator reports `15/15`, but absolute review
rejects the packet: its contact drawing contradicts the bottom-center root
contract and it omits protected-blend hashes. No Blender scene, fluid bake,
final animation, or plush integration has been accepted.

## Current state

- The scene owns the tank, liquid domain, cameras, lights, world, shot timing,
  and the scene-level drop controller.
- The reusable asset remains immutable behind the exact `FUMO` /
  `Fumo_Root` / `Fumo_Rig` interface in
  [current_attempt.md](current_attempt.md).
- A deliberately non-likeness placeholder may occupy the interface only for
  scale, framing, path, and clearance checks. It must be visibly labeled
  `NEUTRAL PLACEHOLDER — NOT FUMO`.
- Blender MCP was queried on 2026-08-30 but no connected Blender instance was
  available. The resulting source-grounded, machine-checkable preflight packet
  is recorded in [evidence.md](evidence.md) as a rejected preflight; it is not
  a `.blend` artifact or evidence of simulated water.

## Exact next action

Correct the contact-root diagram and add before/after protected-blend hashes,
then regenerate and re-review the preflight. Blender Stage B remains blocked
until that packet passes.

## Records

- [Requirements and constraints](requirements.md)
- [Acceptance and evidence plan](acceptance.md)
- [Sources and scale basis](references.md)
- [Current attempt](current_attempt.md)
- [Evidence](evidence.md)
- [Artifact log](artifacts.md)
- [Failures and lessons](failures.md)
