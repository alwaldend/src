# 027 execution recovery

The initial console result was unavailable after context compaction. Read-only
recovery found no candidate, receipt, renderer packet, or live writer. Bazel's
command log confirmed the build/run launch, not the model result.

A separate no-save diagnostic captured the causal failure in
head_027_diagnose.log: fringe point X=-58.100685mm, Z=116.634071mm lies
outside the head's projected receiver; nearest fallback was4.722305mm,
exceeding the2.5mm contact-query guard.

Changed only fringe overhang transfer: if either projection misses, use the
nearest original 3D surface point and transfer its panel displacement while
preserving the hair's XZ and depth residual. This is explicitly not contact.
Actual facial/cheek receiver guards remain unchanged. Predicted difference:
the fringe overhang no longer requests a nonexistent projected receiver.
The next execution may save only after all original preservation guards pass.

The changed execution passed fringe and cheek construction, then failed on
rear hair's same unsupported overhang assumption: X=-55.002160mm,
Z=115.469113mm, lateral3.948882mm (head_027_execution.log). Extended the
displacement-transfer rule to rear hair only, using the rear panel map.
Facial and cheek contact guards remain strict. Prediction: rear overhangs
can retain their residual without falsely claiming receiver contact.

Next execution crashed in native Blender/TBB with no Python backtrace before
saving (head_027_execution_rear_fix.log and tmp/hem_026_candidate.crash.txt).
Narrowing hypothesis: dependency-graph evaluation while creating/copying
objects may hit a threaded native failure. Add flushed phase markers and run
the no-save diagnostic with one thread. This distinguishes construction
phase and a parallel-evaluation dependency without overwriting any candidate.

Single-thread run also crashed after attachments but before visibility-update
completion, ruling out thread count alone. Next no-save diagnostic evaluates
after each attachment, logging its name. This discriminates one invalid
attachment/mesh from the batched visibility transition; no shape changes.

All15 attachment evaluations passed separately. The crash is specifically
in hiding the originals, not evaluating the new geometry. Next diagnostic
evaluates after each original visibility change to identify the dependency.

The per-original diagnostic reached the old right-eye applique and crashed
when disabling it. Keep source dependency evaluation enabled (hide_viewport
false), but hide originals via hide_set and hide_render. This is a targeted
scene-level workaround for the observed dependency-disable crash, not an
animation/integrity pass. Prediction: same candidate meshes evaluate and
originals remain visually absent without disabling their graph nodes.

That workaround saved027 (8c38a46141eae10d44a00d61de1d745e134717a5ce2fe42a9ad3433b487f9ae4)
and rendered all5 views, but copied hair/eye materials rendered white. It is
not retained. Revised causal hypothesis: _snapshot retained evaluated mesh
material IDs and _make linked these transient IDs into persistent new meshes.
This may explain both invalid shading and the disable-time native crash.

Prepared separate027b scripts, leaving saved027 and its exact sources intact.
Snapshot material slots now resolve by name through bpy.data.materials while
the temporary evaluated mesh is alive. The writer asserts every new object's
material is non-null, not evaluated, and a member of the persistent database.
Restored ordinary original-object disabling to test whether correcting the
material ownership removes the crash. No geometric or material-node changes.
Prediction: persistent material binding, original shader appearance, and no
disable-time crash. Independent material-ID audit is running concurrently.
