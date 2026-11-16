## Why

NoobAI's existing prompt sequences generate independent images and scatter routine
controls across a graph. The operator wants one streamlined sequence of related
images, with previous-image carry-forward as the default and persistent references.

## What Changes

- Add `noobai-image-sequence` and its native `.app.json` projection, with shared subject/setting text, an action list,
  persistent IP-Adapter references and ordered image-to-image sampling.
- Keep reference selection, change strength, dimensions, seeds and partial-run
  controls directly accessible; save each completed image immediately.
- Supply a generated native workflow, operating instructions and repeatable offline
  execution evidence. Retain the existing profiles and workflows.
- Add one looping GIF for the selected sequence,
  with adjustable image duration and preserved individual PNGs.
- Provide separate reference-sequence and scene-composer Apps. The first retains
  persistent uploaded references and carry-forward edits. The second adapts the
  supplied composer's shared character/clothing/style fields, per-scene settings,
  optional locally installed LoRA, fixed sampling seed and independent two-pass
  generation. Both retain scene grouping and playback outputs.
- Provide separate source-and-sketch edit graphs and Apps, including a painted-mask
  variant, using the already approved NoobAI Canny checkpoint and native nodes.

## Capabilities

### New Capabilities

- `image-sequences`: Reference-anchored, sequential still-image generation and
  deterministic continuation through a native ComfyUI workflow.

### Modified Capabilities

None.

## Impact

Reusable nodes and workflows belong to `projects/comfyui_tools`. The infrastructure
owner will separately record NoobAI packaging changes. No new external dependency,
model, deployment, paid rental or live inference is authorized by this change.
