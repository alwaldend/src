## 1. Sequence execution

- [x] 1.1 Write failure-first behavioral coverage for scene resets, feedback, range seeds, invalid input and interruption, producing repeatable PNG/JSON evidence.
- [x] 1.2 Implement reference preparation and a scene-aware native-sampler output node; verify the execution contract and saved metadata with the offline fixture.

## 2. Workflow and App

- [x] 2.1 Generate `noobai-image-sequence` with persistent IP-Adapter references and a native ComfyUI App; verify generated links, widget bindings and App schema against upstream.
- [x] 2.2 Package the workflow and node through the owning infrastructure change; verify every configured startup remains within its payload limit.
- [x] 2.3 Document scene syntax, continuation, controls and import instructions; inspect the generated workflow and preserve the existing selected profile.

## 3. Delivery

- [x] 3.1 Run required formatting, semantic lint, affected behavior/packaging checks and local review against the exact candidate; record evidence and limits.
- [x] 3.2 Prepare the task-owned candidate for guarded publication; delivery receipts own the final commit, push, pull request and remote-feedback state.

## 4. GIF playback outputs

- [x] 4.1 Define failure-first GIF grouping, decoded playback, partial-run, timing and interruption cases; retain actual GIFs and a repeatable receipt.
- [x] 4.2 Add per-scene and combined GIF serialization with an App timing control; verify real media and preserve existing PNG and sampling behavior.
- [x] 4.3 Prepare the updated App candidate and its validation scope for guarded publication through the existing PR; delivery receipts own the final checks, review and publication state.

## 5. Default demonstration

- [x] 5.1 Include a default six-image, two-scene example that repeats camera/location tags within each scene and demonstrates small pose changes; inspect its generated App values.
- [x] 5.2 Execute the generated default example on the explicitly authorized NoobAI 96 GB deployment; inspect all PNGs and scene/combined GIFs, preserve artifacts and provide the retained endpoint.

## 6. Optional scene structure guidance

- [x] 6.1 Exercise failure-first continuation, scene-reset, disabled-control and invalid-setting cases for previous-image Canny conditioning.
- [x] 6.2 Add the approved NoobAI Canny checkpoint and native guidance, exposing strength and end percentage in the App.
- [x] 6.3 Deploy through checked-in in-place configuration on the same rental, verify runtime compatibility and inspect the default result.

## 7. Scene interpolation

- [x] 7.1 Verify failure-first decoded GIF timing and scene separation across an interpolation boundary.
- [x] 7.2 Reuse the repository-pinned Video2X/RIFE engine for optional intermediate GIF frames, preserving keyframe PNGs and scene cuts.
- [x] 7.3 Verify actual RIFE output on the retained rental and show the resulting scene and combined GIFs.

## 8. Review feedback: interpolation cancellation

- [x] 8.1 Exercise a real interpolation process and descendant that ignores termination; require native cancellation to stop both, retain completed PNGs and clean run-owned temporary media.
- [x] 8.2 Poll native interruption during interpolation, enforce the existing deadline and terminate the owned process group on cancellation, timeout or failure.
- [x] 8.3 Validate and install only the corrected node while preserving current workflow edits; prepare guarded delivery of the review fix. Delivery receipts own publication and review-thread resolution.

## 9. Reference and composer Apps

- [x] 9.1 Capture the supplied composer's actual sampling/prompt contract and write failure-first end-to-end cases for both variants, optional LoRA, partial ranges and interrupted refinement.
- [x] 9.2 Add explicit prompt composition, optional native LoRA loading and backward-compatible independent/two-pass sampling controls.
- [x] 9.3 Generate separate reference and composer graph/App pairs, with a two-scene example for each; package them through their infrastructure owner without new dependencies.
- [x] 9.4 Verify actual PNG/GIF/metadata artifacts, generated App bindings and affected behavior/packaging; prepare both variants for guarded publication. Delivery receipts own publication and review state; preserve live workflow edits and rental retention.

- [x] 9.5 Remove per-scene GIF publication from both variants; verify one combined GIF, preserved scene cuts and PNG-only recovery after interruption.

## Installation ownership

- [x] Document the shared Ansible role as the owner of selected workflow and node installation; retain existing generated graphs and native behavior.

## 10. Separate sketch-guided edits

- [x] 10.1 Record source/sketch/mask roles, native-node design and disconfirming validation cases before implementation.
- [x] 10.2 Generate whole-image and painted-mask sketch-edit graphs and matching native Apps using the approved Canny model.
- [x] 10.3 Select the documents in the NoobAI profile and document alignment, controls, mask behavior and verification limits.
- [x] 10.4 Inspect exact generated graphs, run offline installation and preserve repeatable receipts; delivery receipts own final publication state.
