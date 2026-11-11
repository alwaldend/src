---
title: Chained image to image GIF
---

This output node uses native ComfyUI sampling and VAE decoding to create one
frame for each Prompt Sequence line. Each decoded frame is VAE-encoded as the input
to the next frame. The initial image is a starting condition, not an extra frame.
Negative conditioning and sampling settings are shared. `frame_source` selects
the previous decoded image or the original reference for each frame;
`frame_seed_mode` selects incrementing seeds or fixed noise. Saved nodes without
these controls retain previous-frame feedback and incrementing seeds. Original
reference plus fixed noise reduces accumulated drift and makes identical
first/last prompts suitable for a controlled loop boundary. Neither control
supplies a motion trajectory. Start with denoise 0.45, 40 steps and 6 fps for full Chroma.

The complete chain is one dependent queue job. Ordinary panel sequences also
stay in one job to keep their media grouped. GIF frame prompts describe successive changes to
the same scene; this is discrete img2img chaining, without motion interpolation.
Lower denoise preserves more of the preceding frame; higher denoise changes it
more. GIF uses a 256-color palette per frame and frame intervals rounded to
10 milliseconds. All frames are retained, including identical adjacent frames.
The GIF is assembled inside the ComfyUI output node and uses its native animated
preview. No external assembly step is required.

The `image-to-image-gif-masked` workflows add native ComfyUI masks so only an
editable rectangle receives denoising. Adjust its size, position and feathering
around the intended moving subject. Eval cases own their rectangle in output
pixels. Qwen's GIF workflows condition every frame on the selected image
through its native image-edit encoder; they start at denoise 1, CFG 4 and 40
Euler/simple steps.

The node uses Torch, NumPy and Pillow already installed by ComfyUI. It does not
download software. A single initial image is required. Sampling errors and
interruptions do not publish a partial GIF. Files use native output naming,
stay within the output directory and never overwrite an existing GIF.

Frame decoding uses ComfyUI's native VAEDecode normalization, including Qwen's single-frame temporal VAE output. CPU frames are converted to float32 for
NumPy/Pillow encoding. This compatibility path is covered by the real-GIF
fixture; successful GPU generation remains separate acceptance.

The separately packaged [Qwen GIF node](../qwen_chained_img2img_gif/README.md) consumes raw frame instructions, a Qwen
text encoder and the initial image. In previous-frame mode it feeds each
decoded result through the native Qwen image-edit encoder for the next
positive/negative conditioning, and through VAE encoding for that frame's
sampler input. This remains genuine image feedback at denoise 1. The
original-reference control remains available for independent-edit comparisons.
The existing generic node retains its saved-input interface.
