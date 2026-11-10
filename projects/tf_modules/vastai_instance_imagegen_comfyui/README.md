---
title: Vast.ai ComfyUI instance module
description: Loopback ComfyUI with native model provisioning and authenticated HTTPS
---

This module composes [the base rental](../vastai_instance_imagegen/README.md)
with the native Vast ComfyUI service at `127.0.0.1:18188`. Supply a reviewed
`vastai/comfy` image through `image`. It forwards the base's access, credential,
offer, storage and connection interfaces, retaining SSH-only access by default.

ComfyUI uses `--disable-all-custom-nodes --offline`; Manager and third-party
custom nodes are not enabled, partner API nodes are disabled, and the frontend
receives ComfyUI's restrictive offline Content-Security-Policy. No wildcard CORS
flag is supplied. The bundled API wrapper remains on loopback `127.0.0.1:18288` for scripted
jobs. It is included in the local portal inventory but never routed by the
public proxy. Portal, Jupyter and TensorBoard also remain local; tunnels, sync
and worker services remain disabled. These controls limit exposure; they are not an outbound firewall.

Model kinds map below `/workspace/ComfyUI`:

| Kind           | Directory                 |
| -------------- | ------------------------- |
| `checkpoint`   | `models/diffusion_models` |
| `lora`         | `models/loras`            |
| `vae`          | `models/vae`              |
| `embedding`    | `models/embeddings`       |
| `text_encoder` | `models/text_encoders`    |

Here `checkpoint` means the standalone diffusion model, loaded with ComfyUI's
Load Diffusion Model node. This profile does not assume an all-in-one checkpoint.
Select its encoder and VAE separately in a native workflow. Downloads are served
by Vast's provisioner inside the rented instance, with no model cache retained
after destruction.

```hcl
module "imagegen" {
  source = "../../../projects/tf_modules/vastai_instance_imagegen_comfyui"
  image = "vastai/comfy:v0.44.0-cuda-12.9-py312"
}
```

The user approved this Vast image specifically. Its public source and bundled
services still require maintenance review when upgrading; no registry checker
runs during planning. [Vast's source](https://github.com/vast-ai/base-image/tree/main/derivatives/pytorch/derivatives/comfyui)
owns the boot contract. [ComfyUI's native arguments](https://github.com/Comfy-Org/ComfyUI/blob/v0.44.0/comfy/cli_args.py)
own the offline and custom-node behavior. Offline Terraform acceptance covers
both access modes, native model paths and proxy/forwarding ports; actual image
boot and GPU generation require separately authorized live acceptance.
