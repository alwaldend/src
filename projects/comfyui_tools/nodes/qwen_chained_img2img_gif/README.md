---
title: Qwen chained image to image GIF
---

This node refreshes Qwen's positive and negative image-edit conditioning from
each preceding decoded image, and VAE-encodes those same pixels for the next
frame. The shared chained_img2img_gif node owns sampling, decoding, progress
and GIF publication. Both first-party node directories must be enabled; the
native registry resolves the shared node when the workflow is inspected or run,
so custom-node import order does not matter.

Only Qwen evaluation rentals install this model-specific node. Other models
retain the shared renderer without carrying Qwen code in their startup payload.
This uses ComfyUI's installed native Qwen encoder and adds no external package.
The shared real-GIF fixture exercises its reference and latent feedback.
