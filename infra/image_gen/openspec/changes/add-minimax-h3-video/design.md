## Context

The 48 GB Blackwell rental profile uses a ComfyUI image with native support for
H3 Base text-to-video and first/last-frame image-to-video.

## Goals / Non-Goals

Make H3 Base available through the existing Ansible model selection. The
requested rollout selects H3 on the existing 48 GB hardware profile and
uses the owning Terraform and Ansible targets for rental and installation.

## Decisions

- Use Comfy-Org's pruned INT8 FL2VA denoiser and NVFP4 text encoder to reduce
  offloading pressure. Use the standard audio and video VAEs from the same
  repository revision.
- Register the prepared graphs from
  [ComfyUI tools](../../../../../projects/comfyui_tools/openspec/changes/add-minimax-h3-video/design.md)
  for installation with the model.

## Risks / Trade-offs

Peak memory depends on resolution, duration, and host RAM for offloading.
Document a short first render in the operator guide.

## Deployment on 2026-11-21

The reviewed Terraform apply created Vast.ai instance `61823418` with zero
changes or destructions. Vast.ai reported one RTX PRO 5000 Blackwell with
48,935 MiB of VRAM, a 100 GB disk, and a rate of USD 0.9151 per hour.
The reviewed refresh-only apply recorded the connection outputs. Direct and
proxy SSH presented the same Ed25519 host key, which was pinned for the
installation. The authenticated ComfyUI endpoint reported version 0.44.0 and
the expected GPU.

Ansible installed all four weights and both saved workflows and restarted
ComfyUI. The play recap reported 19 successful tasks, five changed tasks,
and zero failures or unreachable hosts. The running service listed the four
weights under their expected loader choices. Both saved workflow files served
by ComfyUI matched the packaged files by SHA-256.

The authenticated CLI ran a shortened copy of the saved text-to-video
workflow. All 14 nodes executed and produced a 640×352, 39-frame MP4 at
24 fps with stereo AAC audio. Both streams decoded on the rental, and the
downloaded file matched the decoded source by SHA-256. First/last-frame
verification covered workflow installation and model availability.

## Invocation-time profile selection

Terraform accepts `hardware_profile` from the existing hardware registry.
Ansible accepts the `imagegen_model_profile` extra variable and uses the
selected profile for both asset lookup and workflow paths. Each stage uses
the default in `config.json` when its override is omitted.

### Deployment verification on 2026-11-23

The reviewed apply with `TF_VAR_hardware_profile=48gb` added one rental,
Vast.ai instance `62131385`, with no changes or destructions. The provider
reported an RTX PRO 5000 Blackwell with 48,935 MiB of VRAM and a rate of
USD 0.8263 per hour. A read-only plan rejected an unregistered profile name.

The playbook with `imagegen_model_profile=minimax-h3-video` installed all
four weights and both workflows with zero failures. The running service
listed the installed weight files and served workflows that matched the
packaged files by SHA-256. ComfyUI accepted a first/last-frame job and
began GPU execution.
