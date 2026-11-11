---
title: Chained image to image GIF
---

This output node uses native ComfyUI sampling and VAE decoding to create one
frame for each Prompt Sequence line. Each generated latent becomes the input
to the next frame. The initial image is a starting condition, not an extra frame.
Negative conditioning and sampling settings are shared; the seed increments
once per frame. Start with denoise 0.45, 40 steps and 6 fps for full Chroma.

The complete chain is one dependent queue job. Ordinary panel sequences also
stay in one job to keep their media grouped. GIF frame prompts describe successive changes to
the same scene; this is discrete img2img chaining, without motion interpolation.
Lower denoise preserves more of the preceding frame; higher denoise changes it
more. GIF uses a 256-color palette per frame and frame intervals rounded to
10 milliseconds. All frames are retained, including identical adjacent frames.
The GIF is assembled inside the ComfyUI output node and uses its native animated
preview. No external assembly step is required.

The node uses Torch, NumPy and Pillow already installed by ComfyUI. It does not
download software. A single initial image is required. Sampling errors and
interruptions do not publish a partial GIF. Files use native output naming,
stay within the output directory and never overwrite an existing GIF.
