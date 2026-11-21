## Context

The 48 GB Blackwell rental profile uses a ComfyUI image with native support for
H3 Base text-to-video and first/last-frame image-to-video.

## Goals / Non-Goals

Make H3 Base available through the existing Ansible model selection. Rental
lifecycle and GPU acceptance are separate operator actions.

## Decisions

- Use Comfy-Org's pruned INT8 FL2VA denoiser and NVFP4 text encoder to reduce
  offloading pressure. Use the standard audio and video VAEs from the same
  repository revision.
- Register the prepared graphs from
  [ComfyUI tools](../../../../../projects/comfyui_tools/openspec/changes/add-minimax-h3-video/design.md)
  for installation with the model.

## Risks / Trade-offs

Peak memory depends on resolution, duration, and host RAM for offloading.
Document a short first render and the model's territorial license restrictions
in the operator guide.
