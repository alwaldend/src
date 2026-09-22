# Blender asset safety and verification

## Establish identities before mutation

Resolve and record:

- workspace root and task-owned `out/<task>/` root;
- source and candidate absolute paths;
- source hash when preservation matters;
- Blender toolchain identity;
- scene and linked-library inventory;
- missing external resources; and
- the exact scene, objects, cameras, and outputs in scope.

Use a new candidate path. Do not rely on Blender's numbered backup as a
rollback mechanism, and direct backup/temp/cache outputs under the task root.
After opening, verify that Blender loaded the intended file rather than a
startup scene or stale session.

## Separate editing from evidence rendering

Never render a file concurrently with a foreground save. Publish a complete
checkpoint copy, hash it, and use the pinned background Blender to reopen that
immutable snapshot. Give every render packet a fresh output directory and a
manifest containing the candidate, toolchain, scene, cameras, settings, and
output hashes.

Keep one regression series on one recorded Blender and graphics environment.
Backend changes may still yield useful qualitative images, but byte hashes
and small pixel differences are not comparable as one deterministic series.

Render or viewport success is not visual acceptance. Inspect the final image,
including framing, blank or missing subjects, clipping, lighting failures, and
the views most likely to expose regressions. Apply a specialist visual skill
when the request has artistic or reference-fidelity acceptance criteria.

## Verify the deliverable

For a changed `.blend`:

1. confirm the save command completed and the expected path exists;
2. hash the candidate;
3. reopen it in pinned background Blender with automatic file scripts
   disabled;
4. inspect required scenes, objects, links, modifiers, dependencies, and task
   postconditions;
5. render or export a minimal discriminating artifact when structure alone
   cannot establish the result; and
6. rehash protected inputs and confirm they remain unchanged.

For an export, verify the exact exported file and material request-dependent
properties, not merely the presence of an export operator result. For a
render, verify dimensions, nonempty meaningful pixels, expected camera, and
manifest hashes. Keep diagnostics separate from the user-facing presentation
artifact.

Promotion of a candidate into a tracked or protected path requires the user's
existing mutation authority and the owning project workflow. This safety
procedure grants no permission to replace canonical assets.
