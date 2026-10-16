# User-requested save — 2026-09-06

Attempt independent-original-macro-115 remains unfinished. No full model or
stage acceptance; all eight final criteria remain unverified. User requested
saving, so no additional modeling or rendering was started at this checkpoint.

## Saved files and observed completion

All paths below are relative to out/reimu_fumo_finish/desktop_astra/.

- Hair: hair_parallel_115_a.blend, SHA256
  46020796b7bb913e2acd1dd895e70eb1a1ff8efad1d0293d2d764801b05b335d.
  Clean-open five-view packet hair_parallel_115_a_fast/ is complete; review log
  contains HAIR115_REOPEN_RENDERED and Blender quit. Root has NOT yet inspected
  these pixels; no visual verdict. Builder and review scripts are preserved.
- Body: body_parallel_114_candidate_b.blend, SHA256
  c0c41fd23e4e5fc862c756d19844046329ba6704e618f11dd421da741e7bf048.
  Eight-view packet body_parallel_114_review_b/ is complete; review log contains
  BODY114_REVIEW_OK and Blender quit. Root's prior front/side/three-quarter
  inspection found visible torso/skirt intersections, a tent-like skirt, a
  broad flat support, and mechanical hem. Not retained as an acceptable body.
  A/B exhausted the worker's two coherent body candidates; do not sweep more
  parameters. A remains rejected, unbound diagnostic evidence only.
- Face: face_115_layout.py, SHA256
  be9c9ef6151cc1fc730ca6238da39741ab3e524d81c9bcb3e114a1ca1c5dd99e,
  and face_115_layout_plan.md saved. Script NOT executed. It creates hidden,
  review-only face-placement footprints, not final facial graphics or materials.

These hashes were rechecked during this save. git diff --check passed before
the checkpoint write. All seven root workers and the body review child were
observed completed; none is actively modifying task files. Native render logs
show both Blender processes exited. No new native job was launched.

## Resume exactly here

1. Inspect hair_parallel_115_a_fast/close_front.png, close_side.png, and
   close_three_quarter.png against the controlling references. Send the completed
   packet to fresh_hair_model for bounded diagnosis; its last response predates
   root's native run. Hair budget: at most two coherent candidates A/B plus one
   causal fit repair. Do not reuse rejected 110–114 generators.
2. Have seated_body_draft record B's actual pixel failures and stop that sweep;
   its last response likewise predates root execution. Resolve the construction
   problem before another body strategy, not by adding final materials.
3. Root may execute face_115_layout.py with -- build then -- render using the
   existing pinned native route, and inspect the guide placement only.
4. Preserve original head109 and fallback arm101. Integrate only visually viable
   pieces, with explicit collar/arm/bow fitting and whole-character review.

Root is sole native/model/goal/Git writer under root-native-batch-115; workers
only draft isolated scripts and review returned pixels. Passing preflight:
../preflight/root_native_115/receipt.json, SHA256
5be64219e985051e0601a78be66bd3c908d436ad651cad4fbfcb10af9d146257.
Reuse while binding remains unchanged. Blender 5.2.1 LTS, build 9e2066aef7ef,
background/factory-startup/disable-autoexec, four threads, python-exit-code 2,
through bazel_agent bazel run //tools/blender:blender. Task config/cache/TMPDIR
remain in desktop_astra. No GUI dependency.

Local checkpoint only: branch t3code/continue-fumo-desktop-use, HEAD
c7601f0fc80e0a94b585910459639ffccbdbdbd4. No commit, push, or remote backup.
Use canonical receipt for resourceVersion. Existing legacy plan cardinality
limit is settled: keep the open ordinary attempt; do not add portable plans.
