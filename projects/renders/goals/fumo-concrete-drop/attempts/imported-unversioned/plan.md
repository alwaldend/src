# Fumo concrete-drop scaffold goal

## Goal

Create a reusable Blender scene scaffold for a future approved Fumo asset
being dropped onto a concrete floor.  The current stage uses only a clearly
labeled neutral rigid placeholder to prove scale, framing, timing, collision,
lighting, and the replacement interface.

## Status

`in progress` — Attempt 2 passes the scripted mechanics but fails composition:
frames `1` and `12` crop most of the falling proxy. This goal does not authorize
changes to the reusable Reimu model or claim plush deformation, final
animation, or final integration.

## Current state

- Current candidate: Attempt 2 under
  `out/fumo_concrete_drop_scaffold/attempt_02/`, SHA-256
  `a9488c220c5076a3202e61c9897cf3710f24b1abe74fb9edfc4750bfaebfdc26`.
- Parent checkpoint: none; the reusable Reimu Fumo remains unapproved.
- Stage: neutral placeholder scaffold.
- Last accepted checkpoint: none.
- Attempt 2 preserves the measured contact at frame `22`, maximum sampled
  penetration below `.000687 m`, zero late-motion span, metric scale, interface
  contract, and protected source hashes. Its warning and settled proxy are
  readable, but the initial fall is not completely framed.
- Failing or unverified criteria: complete start-to-floor framing and exact
  clean-reopen audit.
- Exact next action: widen or tilt the camera without changing mechanics, then
  render the same samples and audit the exact saved candidate.

## Current plan

1. **Completed:** build Attempt 2 with only the three documented visual
   corrections.
2. **Completed:** render the same sampled frames and compare them with Attempt
   1.
3. **Rejected:** implementation-blind review found the initial proxy cropped.
4. Widen the next camera, repeat the packet, and preserve exact evidence hashes
   without promoting anything
   into a tracked Blender asset.

## Records

- [Requirements and acceptance](requirements.md)
- [Current attempt](current_attempt.md)
- [Failure ledger](failures.md)
- [Evidence manifest](evidence.md)
- [Artifact log](artifacts.md)
- [Attempt 1 history](attempts/attempt_01.md)
- [Attempt 2 history](attempts/attempt_02.md)
