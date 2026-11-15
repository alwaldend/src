---
title: Image sequences
---

`Image Sequence` generates ordered still images in one ComfyUI job. Each scene
starts from fresh noise. Further actions edit the previous decoded image, while
the graph's original IP-Adapter references remain active. It saves each completed
PNG immediately, before starting another action. Cancellation leaves completed
files in ComfyUI's output directory.

Use one complete desired-image prompt per line, with `---` alone on a line to
start a new scene. Blank lines provide spacing and do not create images. Scene
breaks reset composition, not references or shared appearance text. NoobAI is a
tag-oriented generator: describe the desired result rather than writing only an
edit command such as “move the hand”.

`Sequence References` provides one required appearance reference, two optional
views and an optional continuation image. `(none)` omits an optional input.
Reference views are fitted inside 512-square canvases with edge padding, preserving
their full contents. All references influence the whole image; this is not
per-character region assignment. Continuation retains its original pixels until
the sampler resizes it to the requested canvas.

For a partial run, keep the complete action list, set `start_image` to its absolute
action number and `image_count` to the wanted number (`0` means all remaining).
Starting inside a scene requires its accepted preceding image in `continuation`.
Starting at a scene boundary ignores continuation and creates a fresh composition.
The seed is the base seed plus the absolute action index, modulo 2^64. Keep the
seed control fixed to repeat a run; extending the list does not change earlier
action seeds. Inserting or deleting earlier actions changes later indices.

Change strength starts at `0.45`. Lower values retain more composition; higher
values allow larger changes. Scene starts always use `1.0`. References reduce
appearance drift but do not guarantee identity, anatomy or action fidelity.

Each PNG includes the native workflow/prompt metadata and an `image_sequence`
record of its actual scene/action, seed, effective prompts and sampler settings.
Filenames share a fresh run identifier and contain ordered scene/action numbers.
The native App gallery receives all completed results when the job finishes;
files from interrupted runs remain on disk even when the job has no final gallery.

After the run finishes, one combined GIF plays all generated images in action
order. The App exposes `gif_frame_duration_ms`, initially 500 milliseconds per
image; values from 20 to 10000 in multiples of 10 are accepted. PNGs and GIFs share
the run identifier. GIFs use the already-saved PNGs, keep their dimensions and
add no intermediate motion when interpolation is disabled. Identical adjacent images can share an encoded frame
with a combined hold time. Partial runs contain only their newly generated images.
Per-scene GIFs are not published. Interruption preserves completed PNGs and
does not publish an incomplete GIF. The native gallery includes the stills and
the combined animation.

Optional `control_net` receives the approved NoobAI Canny model. Within a scene,
native Canny extracts edges from the resized preceding image and native
ControlNet conditions both positive and negative prompts. `control_strength`
defaults to 0.35; zero disables it. `control_end` defaults to 0.8. Scene starts
skip this guide. Lower the strength for larger pose changes; edge guidance can
resist them and does not guarantee identity or temporal consistency.

`interpolation_multiplier` accepts 1 (disabled), 2, 4 and 8. The App defaults to
4x and requires the configured Video2X 6.4.0/RIFE 4.25 installation and Vulkan
device. RIFE smooths each scene separately; the combined GIF retains hard scene
cuts. Scene playback duration is preserved with GIF's 10-millisecond precision.
The original PNGs remain untouched. Interpolation can distort inconsistent poses.
Native cancellation also stops interpolation and its child processes, removes
run-owned temporary media and retains completed PNGs.

The separate `noobai-reference-sequence` preset retains this reference/continuation
behavior, with Canny and smoothing disabled initially for inspection. The
`noobai-scene-composer` preset uses the same sampler/output node with
`generation_mode=independent`, `seed_mode=fixed`, 25 first-pass steps and 20
refinement steps at denoise 0.4. Every first pass starts from fresh noise; its
second pass consumes only that image's sampled latent, with the same seed and
conditioning. Zero refinement steps disables the second pass. The legacy node
defaults remain previous-image generation, incrementing seeds and no refinement.

`Scene Prompts` combines fixed character, clothing and style descriptions with
one setting per scene and one action per line. A single setting applies across
all scenes; any other count mismatch fails before sampling. `---` still groups
scenes for settings and hard playback cuts. An independent partial run needs no
continuation image. The composer has no image-reference adapter or edge control.

`Optional LoRA` defaults to `(none)` and passes the model/CLIP pair unchanged.
Selecting an already-installed filename delegates to native `LoraLoader` with
separate model and text-encoder strengths. No model is downloaded. Shared tags,
fixed seeds and refinement do not guarantee character identity or better quality.

The [NoobAI workflow instructions](../../README.md#noobai-image-sequence-app)
describe importing the generated App. The runtime source is staged as one native
custom-node module. Interpolation additionally uses the repository-pinned
Video2X installation managed by the infrastructure owner.

`ImagegenKeyframeGIF` assembles an approved IMAGE batch without loading a model.
Use 2 to 128 equally sized poses, a frame duration divisible by 10 ms and an
interpolation multiplier of 1, 2, 4 or 8. Reversible playback traverses the poses
and returns to the first, allowing RIFE to smooth both directions. The node
retains each source as a PNG and publishes one looping GIF using the same writer
as image sequences. Interpolation requires the deployment's existing Video2X
engine. Validate the poses first; neither reversal nor RIFE fixes incorrect edits.
