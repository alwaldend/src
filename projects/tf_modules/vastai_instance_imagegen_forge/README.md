---
title: Vast.ai Forge instance module
description: Native Forge settings on the shared image-generation rental
---

This module composes [the base rental](../vastai_instance_imagegen/README.md)
with the native Vast Forge service at `127.0.0.1:17860`. Supply a reviewed Forge
image through `image`. It forwards the base's access, credential, offer, storage,
and connection interfaces, retaining SSH-only access by default.

Model kinds map below `/workspace/stable-diffusion-webui-forge`:

| Kind           | Directory                 |
| -------------- | ------------------------- |
| `checkpoint`   | `models/Stable-diffusion` |
| `lora`         | `models/Lora`             |
| `vae`          | `models/VAE`              |
| `embedding`    | `embeddings`              |
| `text_encoder` | `models/text_encoder`     |

Optional `model_selection` selects a declared checkpoint, native `preset`,
encoder/VAE `modules`, integer `steps` and `cfg_scale`. Native JSON config sets
both global and preset-specific keys through `--ui-settings-file`. Omitting it
preserves the image defaults. These creation-time settings replace a rental.

Consumers previously using the base's Forge-specific interface must change
source to this wrapper. A live instance gains a nested `module.imagegen` state
address and requires an explicitly reviewed migration or replacement; changing
source is not authorization to perform either. The personal consumer's previous
rental was destroyed before this refactor.

```hcl
module "imagegen" {
  source = "../../../projects/tf_modules/vastai_instance_imagegen_forge"
  image = var.forge_image_reference
}
```
