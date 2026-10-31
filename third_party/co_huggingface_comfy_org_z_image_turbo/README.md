---
title: Z-Image-Turbo model inputs
description: Integrity-checked upstream weights for a local ComfyUI CPU trial
---

[Z-Image](https://github.com/Tongyi-MAI/Z-Image) is owned by Tongyi-MAI.
These [ComfyUI-packaged weights](https://huggingface.co/Comfy-Org/Z-Image-Turbo)
retain their upstream Apache-2.0 license and ownership. The diffusion model,
Qwen3-4B text encoder, and VAE total about 20.7 GB. Immutable URLs and publisher
SHA-256 checksums are declared in `include.MODULE.bazel`; no runtime fetches.

Targets are manual to exclude the downloads from wildcard builds. The local
consumer is `//tools/comfyui:z_image_turbo_benchmark`.
