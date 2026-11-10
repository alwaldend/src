## Context

Juggernaut XL v9 is an all-in-one SDXL checkpoint with bundled encoders and VAE.
The approved RES4LYF revision includes an SDXL style example and ReSDPatcher;
the archived Chroma graph uses ReChromaPatcher and cannot be reused unchanged.

## Decisions

Use <profile>-<purpose>.json everywhere: file, profile manifest and saved workflow.
Use full text-to-image and image-to-image purposes; keep single-reference,
character-style, prompt-sequence, style-reference and style-prompt-sequence
where those accurately describe the graph. The folders retain model ownership.
Preserve existing Flux, Qwen and Chroma JSON bytes during renaming.

Native Juggernaut graphs use CheckpointLoaderSimple, CLIPTextEncode, KSampler,
VAE decode and SaveImage. Image-to-image scales and center-crops the reference
to 1024 square and encodes it as the starting latent. Text-to-image uses fresh
noise. Prompt sequences use the existing first-party node and shared text.

For style-reference graphs, encode the reference with the checkpoint VAE and
route it through ClownGuide_Style_Beta to ClownsharKSampler_Beta. Patch the SDXL
model with ReSDPatcher in float32. Start from an empty SDXL latent rather than
the reference's composition. Expose guide weight, use a neutral illustrated
prompt, and keep a native negative prompt. The sequence variant reuses the same
guide for every panel. This is style guidance, not guaranteed style isolation or
character identity preservation.

Decision review: proceed with the previously approved pack and exact dependency
lock. Plain img2img carries composition and does not meet the reference-style
request. A new adapter and vision-model dependency are unnecessary for this
source option. Native-only graphs remain available if the pack fails on the
next rental. Do not claim inference or quality acceptance from offline checks.

## Validation

Before acceptance, check malformed/nonreciprocal links, mismatched slot types,
wrong checkpoint/VAE paths, a missing model patcher or guide edge, disconnected
sequence outputs, missing widget metadata, unpinned dependencies, loss of old
profile bytes and startup-payload constraints. Existing browser and Terraform
CLI fixtures cover editor duplication and provisioner behavior. Retain a graph
and dependency verification receipt under ignored task scratch. No rental or
private content is used by these checks.
