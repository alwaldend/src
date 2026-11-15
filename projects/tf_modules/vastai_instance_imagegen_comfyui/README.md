---
title: Vast.ai ComfyUI instance module
description: Provision a loopback ComfyUI runtime with authenticated HTTPS
---

This module composes [the base rental](../vastai_instance_imagegen/README.md)
with the native Vast ComfyUI service at `127.0.0.1:18188`. Supply a reviewed
`vastai/comfy` image. It forwards access, credential, offer, storage and connection
interfaces, retaining SSH-only access by default. `workspace_directory` defaults
to `/workspace` and is exported before native workspace initialization. For
example, `/srv/imagegen` places the bundled application at `/srv/imagegen/ComfyUI`.
Paths must be absolute and contain no traversal, whitespace or shell syntax.
Changing this boot setting is a provisioning change, not a migration of existing
models, uploads or outputs.

Terraform owns rental constraints, the initial restrictive runtime and access
bootstrap. It does not install model files, LoRAs, node packs, Python dependencies,
first-party sources or workflows. Use the shared
[comfyui_tools Ansible role](../../ansible_collection/roles/comfyui_tools/README.md)
for those inputs and subsequent updates without changing paid-instance startup.
The generic base retains its native provisioner interface for the separate Forge
consumer; this wrapper supplies an empty application manifest.

ComfyUI initially uses `--disable-all-custom-nodes --offline`; Manager and
partner API nodes are disabled, no wildcard CORS flag is supplied, and the
frontend receives ComfyUI's restrictive offline Content-Security-Policy.
Ansible appends the installed node directories to its whitelist and manages the
`COMFYUI_ARGS` line in `/etc/environment`, marked `# imagegen-ansible`. Boot
preserves that marked line; otherwise it supplies the restrictive defaults.
These controls limit exposure and are not an outbound firewall.

The bundled API wrapper stays on loopback `127.0.0.1:18288`, included in the
local portal inventory but never routed by the public proxy. Portal, Jupyter
and TensorBoard remain local; tunnels, synchronization and workers remain disabled.
Only the dedicated TLS and Basic-auth proxy is published in exposed mode.

```hcl
module "imagegen" {
  source = "../../../projects/tf_modules/vastai_instance_imagegen_comfyui"
  image = "vastai/comfy:v0.44.0-cuda-12.9-py312"
  workspace_directory = "/workspace"
}
```

`runtime.json` owns the baseline service, arguments, model-directory map and
application name. The deployment passes that map to Ansible, selecting
`models/checkpoints` for SDXL all-in-one checkpoints and the default
`models/diffusion_models` for standalone denoisers. LoRAs use `models/loras`,
ControlNet uses `models/controlnet`, and other kinds retain their declared paths.
The separate ephemeral evaluator also consumes this declaration; its evaluation
lifecycle continues to own its native provisioner recipe.

[Vast's source](https://github.com/vast-ai/base-image/tree/main/derivatives/pytorch/derivatives/comfyui)
owns the image boot contract. [ComfyUI's native arguments](https://github.com/Comfy-Org/ComfyUI/blob/v0.44.0/comfy/cli_args.py)
own offline and node behavior. Offline Terraform fixtures cover both access
modes, alternate workspace validation, empty application provisioning and
proxy/forwarding ports. Actual image boot and GPU generation require separately
authorized live acceptance.
