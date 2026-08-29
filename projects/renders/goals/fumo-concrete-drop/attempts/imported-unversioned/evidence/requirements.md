# Fumo concrete-drop scaffold requirements

[Back to goal](README.md)

## Acceptance criteria

1. One saved candidate opens cleanly and contains a `FUMO` collection plus a
   named world-transform replacement interface for the future reusable asset.
2. The visible proxy is unmistakably labeled as neutral and placeholder-only;
   no render or metadata claims Reimu likeness or plush behavior.
3. Scene units are metric. The collision envelope is `0.25 m` tall, its visible
   inset proxy is `0.245 m` tall, and the floor contains the sampled motion.
4. The scene records `24 fps`, frames `1..72`, a held pre-roll, rigid-body
   release, floor collision, and a post-impact settling interval.
5. A fixed camera and neutral lighting keep the proxy, floor, scale witness,
   and impact region readable in every sampled frame.
6. The floor is a passive collider and the proxy uses a separate active
   collision object.  Sampled transforms show descent, contact without floor
   tunneling, and materially reduced late motion.
7. The protected tracked Reimu and Sisyphus blends remain byte-identical.
8. An inspectable render artifact and machine-readable audit are bound to the
   exact candidate hash.

## Constraints

- Generated scripts, blends, renders, reports, and caches belong only under
  `out/fumo_concrete_drop_scaffold/`.
- Tracked changes are limited to this goal directory.
- Do not edit, append into, or save over a tracked Blender file.
- The active reusable-Reimu goal remains authoritative for character form and
  rig approval.  This goal may consume only its future stable interface.
- No cloth, soft-body, shape-key, deformation, facial, material-fidelity, or
  final integration claim is in scope.
- Use one candidate writer and one falsifiable Attempt 1 hypothesis.

## Evidence plan

- Reopen the saved blend in background Blender and report collections,
  objects, custom properties, units, render settings, rigid bodies, animation
  range, and sampled world transforms.
- Render frames `1, 12, 20, 28, 40, 56, 72` from one fixed camera and combine
  them into a labeled contact sheet.
- Require the proxy bottom to remain no more than `0.005 m` below the floor
  top in sampled post-impact frames.
- Require a descent of at least `0.40 m` from held pose to minimum sampled
  center height and late-frame center movement below `0.02 m`.
- Hash the builder, candidate, report, and visual artifact; compare protected
  blend hashes before and after.

## Fixed regression set

- `git diff --check` for the goal records.
- Candidate clean reopen.
- Exact `FUMO`/interface names and placeholder-only metadata.
- Metric scale, camera, frame range, floor/active rigid-body roles.
- Protected Reimu and Sisyphus hashes.
- Pixels visibly retain the placeholder warning and do not imply a finished
  character.

## Requirement changes

- `2026-08-30` — user authorized this scaffold as a parallel subordinate
  module, with a neutral placeholder only and explicit exclusion of plush
  deformation and final integration.
- `2026-08-30` — inferred first-gate thresholds were added for sampled floor
  penetration and late motion; they judge scaffold mechanics only.
