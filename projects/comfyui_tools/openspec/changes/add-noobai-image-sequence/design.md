## Context

Existing SDXL sequences concatenate reference embeddings but generate independent
panels. The operator explicitly selects previous-image carry-forward, scene breaks,
and native ComfyUI Apps. Existing weights and IP-Adapter Plus remain the conditioning
owner; native ComfyUI owns encoding, sampling, cancellation and PNG serialization.

## Goals / Non-Goals

Provide one action editor, persistent references, scene resets, partial continuation
and native App controls. No lettering, page layout, video, new dependency, live
deployment or identity-quality guarantee is part of this implementation.

## Decisions

Proceed with a dedicated first-party sequence output node and native graph. An
independently mapped prompt list cannot carry the previous result forward. Reusing
the GIF writer would impose unrelated animation semantics and delay still outputs.
Keep the existing Prompt Sequence parser as the text-splitting owner where practical.
Use native sampler/encoder/save calls; do not implement diffusion or PNG serialization.

A `---` line resets the scene. Each action describes the complete desired image;
NoobAI is not an instruction-edit model. Scene starts use denoise one. Within scenes,
encode the previous decoded image and use editable change strength (initially 0.45).
Keep all original reference conditioning on the model throughout. Partial runs use
absolute action indices for seeds and require a continuation image when necessary.

Generate paired native graph and `.app.json` projections from one definition, with
the same selected controls and explicit graph/App mode. Verify their serialization
against frontend v1.57.6, bundled with the configured ComfyUI v0.44.0.

## Risks / Trade-offs

- Previous-image drift and restricted motion → retain original anchors, expose
  change strength and explicit fresh-scene boundaries; document visual limits.
- Partial completion on cancellation → save each completed PNG immediately, with
  scene/action/seed metadata and native progress messages.
- Startup payload growth → stage the new node only for workflows that use it;
  stage the sequence App as NoobAI's sole startup workflow and retain the complete
  importable catalog. Exercise configured-bundle acceptance against the real manifest.
- Offline evidence cannot establish artistic quality → distinguish interface,
  packaging and sampler-contract checks from separately authorized GPU acceptance.

## Failure cases before implementation

Reject empty actions, empty scenes, invalid ranges, missing continuation, unsafe
prefixes, wrong image batches/dimensions and invalid strengths before sampling.
Exercise fresh-scene resets, actual previous-result feedback, unchanged references,
stable range seeds, repeat runs, individual saved PNG metadata and interruption after
a completed output. Preserve existing workflow behavior and the 16 KiB startup limit.

## Packaging decision and failure evidence

`NOOBAI-STARTUP-SIZE`: the configured-profile fixture rejected the complete ten-file
NoobAI startup after adding the new module. Proceed with the existing optional
`startup_workflows` selection, saving the new App on boot and retaining every
generated workflow for import. This avoids another loader/dependency or weakening
the provider limit. Other profiles and the selected MOSS deployment remain unchanged.

`NOOBAI-APP-BASENAME`: both the ComfyUI wrapper and base runtime-file validator
rejected `.app.json`. Their owning validators now permit only the native optional
suffix, retaining directory and traversal restrictions.

## GIF output extension

Retain NoobAI's description-driven image-to-image behavior. Assemble GIFs from
the exact PNGs already saved by the native writer, preserving generated order and
canvas. Pillow owns GIF serialization, palette conversion and duplicate-frame
coalescing; playback duration, rather than a forced frame count for identical
images, is the contract. This avoids duplicating the existing animation node's
custom low-level encoder, whose one-prompt/one-frame contract differs.

Append an optional `gif_frame_duration_ms` widget (default 500, 20–10000 in multiples
of 10) so older saved nodes retain their sampling widget positions. Save each
selected scene when its last selected action completes, then save the combined
sequence. All outputs share the run identifier. Return GIF descriptors alongside
PNGs with native animated preview enabled. Partial selection covers only generated
actions; cancellation preserves completed PNGs and scene GIFs.

Before implementation, exercise two-scene grouping and decoded frame order,
duration/loop/dimensions, single-frame partial runs, repeated-frame hold duration,
invalid timing before sampling and preservation of completed scenes on interruption.
Recheck the narrow startup headroom after adding serialization and the App control.

The operator subsequently authorizes deploying the existing 96 GB NoobAI profile
and showing the default two-scene result. Keep the renderer and description-based
approach. The example has three small pose changes in a snowy park and three book
poses in a library, repeating camera, location and lighting tags inside each scene.
The infrastructure owner records deployment and retention. The operator explicitly
requires leaving the new rental running and providing its link after completion.

`NOOBAI-GIF-STARTUP-SIZE`: the GIF-enabled module plus saved App exceeded the
unchanged provider limit. Use the existing API-staging workflow (already used by
MOSS/HiDream) for NoobAI's generated documents: `startup_workflows` is empty, and
the deployment saves the generated native App through the owning authenticated
CLI once ComfyUI is ready. Model/node pins and declarative native provisioning
remain unchanged. No source is installed imperatively after boot.

## Optional Canny guidance

The operator approved `Eugeoter/noob-sdxl-controlnet-canny` FP16 and requests
installation without replacing the existing rental. Use native ComfyUI Canny
and ControlNet conditioning on the resized preceding image within a scene.
Start every scene without the previous scene's guide. A zero strength disables
guidance; missing optional ControlNet retains compatibility with existing graphs.
Expose strength and end percentage, initially 0.35 and 0.8, retaining all earlier
widget positions. Guidance constrains edges and can resist large pose changes;
it does not guarantee temporal stability or supply instruction editing.

The infrastructure owner supplies an idempotent Ansible update for the pinned
checkpoint, first-party node and generated App. Restart only ComfyUI after node
changes, preserving the rental and output files. Future bootstrap uses the same
model manifest. Do not apply a replacement plan to the current rental.

## Optional GIF frame interpolation

Reuse the repository-pinned Video2X 6.4.0 AppImage and its bundled RIFE 4.25 model;
no new interpolation dependency is introduced. The runtime requires Vulkan.
Interpolation is applied separately to each scene's saved keyframes, at 1x
(disabled), 2x, 4x (App default) or 8x. Combine the resulting scene GIFs with a
hard cut. Preserve the original PNGs and their metadata. Distribute GIF timing in
10-millisecond units and retain the final keyframe's hold, preserving scene
duration. Repeated GIF frames may coalesce. Fail before sampling when the engine
is absent or timing cannot represent the requested multiplier. Runtime failures
retain all completed PNGs and earlier scene GIFs. Temporary engine extraction and
media remain in a run-owned output subdirectory and are cleaned after encoding.

`NOOBAI-RIFE-CANCELLATION`: automated review correctly identified that the
blocking RIFE wait could keep the GPU/queue occupied after native cancellation.
Check interruption before launch, throughout bounded subprocess waits and after
completion. Own a new process session; on cancellation, timeout or failure send
termination to that process group, escalate to kill after a short grace period
and reap the parent. Preserve existing PNGs and completed scene GIFs.

Failure cases cover cancellation before launch and while the child is running,
an uncooperative descendant with closed output pipes, parent exit with inherited output pipes, timeout,
nonzero diagnostics and subprocess launch failure. A real child/descendant
fixture exercises cancellation, saved-output preservation and temporary cleanup;
the existing interpolation fixture retains successful decoded playback coverage.
The infrastructure updater permits skipping workflow copies for this node-only
fix so later edits on the live App remain untouched.

## Reference and composer variants

The operator supplied Tache's Scenes Composer (`scene-composer.json`, 68 root
nodes and 27 subgraphs) and requested two versions. Its generation path shares
resolved character, clothing, render and environment text, loops over independent
prompts with a common sampler seed, and refines each sampled latent with a second
KSampler before decoding. It contains no image loader, IP-Adapter, ControlNet or
previous-frame feedback. Its optional LoRA stack is unselected. The supplied
serialization contains quarantined subgraph widget bindings, so derive clean
native graphs rather than republishing that file or its preset library.

Proceed with two generated graph/App pairs: `noobai-reference-sequence` retains
the existing reference/carry-forward contract, while `noobai-scene-composer`
uses explicit shared character, clothing and style fields, one environment per
scene and one desired-image action per line. A single environment may apply to
every scene. Reject any other environment-count mismatch before sampling.
Reuse the existing output node with appended optional generation/seed/refinement
controls; old serialized graphs keep their behavior and input order. Composer
images start from fresh noise with a fixed base seed and same-prompt/seed latent
refinement. Independent partial runs require no preceding image. Reference runs
retain absolute incrementing seeds and require continuation inside a scene.

Use a small native LoRA selector around ComfyUI's existing loader, defaulting to
none and accepting only installed filenames with finite strengths. Download no
model or extension. Keep the supplied workflow's external prompt/utility packs
out of the runtime. The two-pass path checks native interruption between passes
and saves only completed images. Both variants retain per-scene and combined
GIFs, partial ranges and optional RIFE; their new presets disable smoothing for
initial inspection. The reference preset also disables previous-pose Canny by
default. No setting establishes exact identity fidelity or successful pose edits.

Before implementation, exercise failures for malformed scene settings, invalid
modes/refinement settings, unknown LoRAs, interrupted refinement, accidental
feedback in independent mode and seed changes across partial runs. Produce real
PNG/GIF/JSON evidence through the existing native-boundary fixture, then inspect
the generated Apps' exact widget bindings and reciprocal links. Preserve the
live App edits and all existing jobs; source publication does not authorize a
new rental, model installation or live configuration replacement.

The operator subsequently removed per-scene GIF outputs. Both variants now save
PNGs and exactly one combined GIF. Interpolate scene groups separately in memory
then concatenate them for publication, retaining hard cuts without publishing
intermediate scene GIFs. Cancellation preserves PNGs and publishes no partial GIF.
