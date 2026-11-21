---
statuses:
  - maintenance
title: ComfyUI tools
description: Reusable ComfyUI custom nodes and image-generation workflows
---

This project owns reusable ComfyUI nodes, workflows and evaluation tooling. Each node has its
own directory under nodes/, with its implementation, test and
README. Workflow definitions are grouped by model under workflows/. Evaluation code and reusable cases live here; run configuration and deployment,
model downloads, credentials and rental state belong to
[image-generation infrastructure](../../infra/image_gen/README.md). The
[repository CLI and local launchers](../../tools/comfyui/README.md) remain in
`tools/comfyui`.

## Workflows

The workflow files, including the [Flux panel sequence](workflows/flux/flux-prompt-sequence.json), contain single-reference, character/style and
prompt-sequence graphs for FLUX.2 klein Base and Qwen-Image-Edit, plus
Juggernaut XL v9, NoobAI XL 1.1, Illustrious XL v2.0-STABLE, HiDream-I1 Full,
full Chroma, MOSS-Image and archived Chroma Flash examples. All files use `<profile>-<purpose>.json`; the saved
workflow name uses that same basename. Load a matching JSON in ComfyUI or pass it to the
authenticated infrastructure CLI. Each graph requires its matching models and
nodes; the rental's [model profile](../../infra/image_gen/tf_setup/README.md#download-models)
selects which files to install. Keeping a workflow here does not install its
dependencies or establish runtime acceptance.

The [Juggernaut txt2img](workflows/juggernaut/juggernaut-text-to-image.json),
[img2img](workflows/juggernaut/juggernaut-image-to-image.json) and
[prompt sequence](workflows/juggernaut/juggernaut-prompt-sequence.json) use
the native SDXL checkpoint loader with the encoder and VAE bundled in the
checkpoint. Start with 1024-square output, 30 DPM++ 2M SDE/Karras steps and CFG
6.5. Img2img starts at denoise 0.65 and preserves some reference composition;
it does not extract only its style. Use [style reference](workflows/juggernaut/juggernaut-style-reference.json) for
a new scene guided by a reference image, or
[style prompt sequence](workflows/juggernaut/juggernaut-style-prompt-sequence.json)
for several panels sharing that style. Both use the previously approved
RES4LYF pack, its SDXL patcher and a VAE-encoded style guide while sampling
from fresh noise. The guide's weight starts at 1; zero disables it for a
fixed-seed comparison. Describe the broad medium and new scene in the prompt.
These graphs start at 30 RES 2M/beta57 steps and CFG 6.5. They guide style;
they do not guarantee exact style isolation or character preservation.
Offline graph checks do not establish GPU execution or style quality.

For several references that contain both characters and style, use
`juggernaut-multi-reference`, or `juggernaut-multi-reference-prompt-sequence`
for a sequence. Select three reference files; reuse a file if you have fewer.
Each reference is resized and center-cropped to 512 square, then all three feed
the approved SDXL Plus IP-Adapter with concatenated reference embeddings. The
standard SDXL sampler starts from fresh noise. The adapter's `weight` starts at
0.8; raise it for more reference influence or lower it for more prompt freedom.
For more scene adherence, try weight 0.5 and `end_at` 0.8; this retained the
neutral references' geometric style while adding the requested reading scene
in the live comparison. Keep `linear` for general subject/style conditioning. Describe the wanted scene,
characters and medium in text as well. These workflows guide appearance and style;
they do not guarantee exact identities or isolate every character in a crowded
reference. The sequence uses the same references for each panel.

The adapter uses the approved
[ComfyUI_IPAdapter_plus revision](https://github.com/cubiq/ComfyUI_IPAdapter_plus/tree/a0f451a5113cf9becb0847b92884cb10cbdec0ef)
and publisher-pinned SDXL Plus and CLIP-ViT-H weights. Upstream is in
maintenance-only mode; future runtime upgrades need compatibility validation.
The older RES4LYF style graphs retain the upstream `float64` style-patcher
default to avoid mixing float32 reference features with its float64 SDXL branch.

`noobai` and `illustrious` provide the same seven SDXL workflow purposes
as Juggernaut. NoobAI also provides
`noobai-image-to-image-prompt-sequence`: enter one prompt per nonempty line in
Prompt Sequence; every image starts from the same original reference. It saves
separate PNGs with shared negative text, a panel counter and unique grouped
filenames. Adjust output dimensions in ImageScale and denoise in KSampler
(default 0.65). They reuse the shared SDXL builder, style patcher and approved
IP-Adapter/CLIP vision weights. Native graphs start at 30 Euler ancestral/normal
steps and CFG 5.5; NoobAI uses neutral tag-style starting prompts. These are
starting presets, with no new GPU acceptance or reference-quality claim.

`hidream-i1-full` provides text-to-image, img2img, style-reference, all three panel
sequences and ordinary/masked chained img2img GIFs. The img2img sequence reuses
the same original reference for each prompt, with editable output dimensions
and denoise. Native graphs use 50 UniPC/simple steps,
CFG 5 and shift 3, following the official full-model example. They load four
local text encoders: CLIP-L, CLIP-G, T5 and Llama 3.1 8B. Style graphs reuse
RES4LYF with its dedicated HiDream patcher and style sampler. This is style
guidance; it is not a trained character-reference adapter. Its GIF workflow
completed the approved plush case on a verified datacenter CMP 170HX 40 GiB
on 2026-11-11. The output drifts in appearance and does not establish convincing
petting motion. Other HiDream workflows still have offline packaging checks
only; their GPU execution and output quality remain unverified.

`uso` provides native subject-only, style-only and combined subject/two-style
reference workflows, each with a prompt-sequence variant. `umo-uno` provides
single and three-identity reference workflows and their sequences, without
USO's style-projector branch. They share the native FLUX.1 graph builder,
load matching model adapters and preserve the publisher's reference positioning.
All sequences reuse the same original reference inputs per prompt and keep
progress plus uniquely grouped filenames. Both families expose output size
and native sampling controls. At CFG 1 their sequence negative output is
unused. See the [model options](../../infra/image_gen/tf_setup/README.md#additional-model-options)
for workflow names, reference roles and starting settings. Source and offline
packaging are available; GPU compatibility and reference quality are untested.

All seven Juggernaut workflows completed on the replacement rental on
2026-11-10, producing 13 neutral 1024-square images at 30 steps. Both style paths
completed without the reported dtype error. This establishes execution for the
selected runtime; reference fidelity and scene adherence still depend on inputs
and guidance settings. The original WCT paths produced softer images in this
check than the adapter path. Deployment evidence belongs to the
[owning infrastructure record](../../infra/image_gen/openspec/changes/add-imagegen-rental/design.md#replacement-gpu-acceptance).

For FLUX.2 dev 32B, use `flux2-dev-text-to-image` without a reference,
`flux2-dev-single-reference` with one image, or `flux2-dev-multi-reference`
with three images. Both reference paths have a matching `-prompt-sequence`
variant. Reference images supply character and style through native VAE latents. Start at
1024 square, 28 Euler steps and FluxGuidance 4. This model uses distilled
positive guidance and has no negative branch; the sequence node's negative
output is unconnected in this family. Its Mistral text encoder runs locally.

The `chroma-hd` profile provides text-to-image, image-to-image,
style-reference, prompt-sequence and style-prompt-sequence workflows with the
full Chroma1-HD checkpoint, FP16 T5 and original VAE, plus a chained img2img GIF
workflow with one prompt per frame. Native sampling starts at
40 Euler/beta steps and CFG 3; style sampling uses the already approved
RES4LYF Chroma patcher and 40 RES 2M/beta57 steps with CFG 3. Its style guide
starts at weight 0.5 and ends at step 15, which preserved more detail than
weight 1 through all steps in the controlled neutral comparison. The two sequences
use the existing multiline node and shared negative prompt. Img2img preserves
composition; the style path guides style from fresh noise and does not promise
character identity. The original three `chroma-*` workflows retain their Flash
defaults as archived examples.

The [full Chroma GPU acceptance](../../infra/image_gen/openspec/changes/add-imagegen-rental/design.md#full-chroma-deployment-on-rtx-pro-5000)
verifies execution and ordered outputs, with significant visual limits: img2img
can preserve the source too strongly, the style guide can blur detail, and comic
wording can elicit extra panels and lettering. Start with native text-to-image;
reference and sequence paths require prompt/strength tuning for useful results.

## MOSS-Image examples

The `moss-image` family supplies six examples and six prompt-sequence variants, adapted from the
[official templates](https://github.com/Comfy-Org/workflow_templates/tree/36ecd0663961c5e1e83f7c0ea521658b33b36bbc/templates)
as compact native graphs without the optional LLM prompt enhancer. Their model
files come from the infrastructure profile; generation uses local Qwen3.5 text
encoding and the small FLUX.2 VAE.

- `moss-image-text-to-image` and `moss-image-base-text-to-image` use SFT and Base,
  respectively, with 1024-square output, 25 ER-SDE/simple steps, CFG 2.5 and
  layer-22 skip guidance through the first half of sampling.
- `moss-image-layout-guided` and `moss-image-base-layout-guided` add an editable
  native bounding-box canvas. Edit the box description and overall caption;
  canvas size drives the prompt coordinates and output latent dimensions.
- `moss-image-image-edit` requires `reference.png` selected in Load Image. Edit
  its instruction; the example uses the SFT edit LoRA, `Moss-Image-Edit` reference
  conditioning on both branches, 30 Euler/normal steps and CFG 4. The reference
  is scaled to about one megapixel and determines the output dimensions.
- `moss-image-color-depth-control` requires prepared `color-map.png` and
  `depth-map.png` selected in their Load Image nodes. Start with both maps at
  1024 square. Region colors are identifiers; depth is white near the viewer.
  Disconnect an optional map input to use only the other map. Change the latent
  dimensions when using another size. The example uses the control LoRA,
  layer-24 skip guidance, shift 5, temporal score rescaling, 25 ER-SDE/simple
  steps, CFG 3 and control ending at 60% of sampling. It does not estimate depth.

The sequence variants use the existing multiline Prompt Sequence node:

- `moss-image-prompt-sequence` and `moss-image-base-prompt-sequence`
- `moss-image-layout-prompt-sequence` and `moss-image-base-layout-prompt-sequence`
- `moss-image-image-edit-prompt-sequence`
- `moss-image-color-depth-prompt-sequence`

Enter one panel prompt per nonempty line, shared context and shared negative
text in Prompt Sequence. One Run saves ordered PNGs with a fresh grouped filename
prefix and displays panel progress. Layout variants reuse the same box captions
and positions; edit variants reuse the original reference for every instruction;
control variants reuse the selected maps. Panels are independent generations,
with no joint attention or guaranteed character continuity. Sampling defaults
and reference requirements match the corresponding single-image examples.

The image-generation Ansible stage installs all twelve selected examples.
You can also save a generated workflow through the authenticated API, for example:

```sh
bazel build //projects/comfyui_tools/workflows:workflows
bazel run //infra/image_gen:comfy -- workflow save \
  "$PWD/bazel-bin/projects/comfyui_tools/workflows/moss-image/moss-image-prompt-sequence.json" \
  --name moss-image-prompt-sequence --where local
```

Repeat the save command for the other five sequence variants. This stages user
workflows without replacing the running rental or changing its startup command.

These are editable examples, not guarantees of exact edits, layout or reference
fidelity. Actual GPU execution and its limitations belong to the
[deployment record](../../infra/image_gen/openspec/changes/add-moss-image-examples/design.md).

## MiniMax H3 Video

- [Text to video](workflows/minimax-h3-video/minimax-h3-video-text-to-video.json)
  starts from an editable example prompt.
- [First and last frame](workflows/minimax-h3-video/minimax-h3-video-first-last-frame.json)
  takes two uploaded images. Select one in each labeled image loader and edit
  the prompt to describe the transition. Match the images' aspect ratio to the
  configured canvas.

Both workflows start at 864×480 with 107 frames at 24 fps and save an MP4
with stereo audio. The deployment profile installs the four preselected model
files and both workflows. The generated graphs use native ComfyUI nodes.

## Workflow source and generation

The maintained sources are [workflows.jsonnet](workflows/workflows.jsonnet),
model-specific `workflows.libsonnet` files and shared graph/node libraries in
`workflows/internal/`. The sequence transformation is shared across model families.
Connections are declared once; generation derives port metadata and reciprocal
link lists. SDXL model families share one graph builder. The public rental
configuration owns an asset registry keyed by filename; profiles reference
those assets and approved node packs by name. Workflows select the first
download of each required kind, with explicit component roles for HiDream's
four text encoders. Additional downloads remain available without
changing those selections. Recorded workflow defaults remain available when
downloads are disabled.

```sh
bazel build //projects/comfyui_tools/workflows:workflows
bazel test //projects/comfyui_tools/workflows:format_test
```

Native JSON is generated beneath
`bazel-bin/projects/comfyui_tools/workflows/<profile>/` at the documented names.
Edit Libsonnet source rather than generated JSON. Terraform and documentation
consume those declared outputs directly. The pinned Jsonnet compiler and
formatter are build tools and are never installed on the rental. To format an
edited source, pass its absolute path to
`bazel run //tools/jsonnet:format -- --in-place /absolute/path/to/source.libsonnet`.

## NoobAI sketch-guided edits

Import `noobai/noobai-sketch-edit.app.json` for a whole-image edit, or
`noobai/noobai-sketch-edit-masked.app.json` to restrict the edit with a painted
mask. Matching `.json` files open in Graph Mode. Both reuse NoobAI XL 1.1 and
the profile's existing Canny ControlNet; no additional weights or node pack is
needed. These are still-image editors, separate from the sequence Apps.

Choose the original picture in **Source image** and an aligned outline drawing
in **Separate sketch**. The source supplies starting pixels; only the sketch
supplies structural edges. Draw the desired final outline, removing obsolete
lines. Arrows, crossed-out parts and written annotations are not edit commands.
Keep both inputs on the same framing and aspect ratio. The output canvas starts
at 1024 square; change its width and height to multiples of eight matching the
source aspect ratio. Source and sketch are resized to that canvas without cropping.

Describe the complete finished image in NoobAI tags, including appearance,
setting and the desired pose. **Denoise** starts at 0.65: raise it for larger
changes or lower it to retain more source composition. **Strength** starts at
0.7 and controls sketch influence; **end_percent** starts at 0.8. Set strength
to zero with the same seed for an unguided comparison. The two Canny thresholds
start at 0.1/0.3; lower them if faint sketch lines are missing. Clear outlines
are the intended input; rough scribbles are not a dedicated scribble-control path.

For the masked variant, open the source loader in **Graph Mode**, use its native
**Mask Editor**, paint the allowed edit area and save the edited input before
running. White permits edits; black preserves the source. Cover both the old
shape and the space the new shape will occupy. An empty mask returns the resized
source unchanged. The same resized mask controls sampling and final compositing,
preserving source pixels wherever it is zero. The unmasked App edits the whole
canvas. Each run saves one PNG with native generation metadata.

Offline graph and deployment checks do not establish sketch fidelity, anatomical
correctness or precise pose edits. GPU inference for these new graphs has not
been performed. Profile selection stages the files through the normal Ansible
deployment; importing an App does not install its model weights.

## NoobAI image sequence App

Import the generated `noobai/noobai-image-sequence.app.json` from the workflow build
into ComfyUI. It opens in native App mode with a reference uploader, shared
appearance/style text, an action editor and a results gallery. Use **Graph Mode**
to inspect its editable graph. The App targets frontend 1.57.6, bundled with the
configured ComfyUI v0.44.0. Its `.app.json` filename also identifies it in ComfyUI's
Apps list. NoobAI workflows are saved through the authenticated workflow API after
boot, keeping its startup within the provider limit. The complete catalog remains
generated and importable; model and node installation still uses normal provisioning.
The matching `noobai-image-sequence.json` opens the same graph in Graph
Mode. It requires the NoobAI checkpoint, approved SDXL Plus
adapter/CLIP vision weights and the first-party `image_sequence` module; an imported
JSON alone does not install those dependencies. Build with:

```sh
bazel build //projects/comfyui_tools/workflows:workflows
```

Select the main appearance reference and optionally two additional views; keep
unused inputs at `(none)`. These references remain active for the whole sequence.
Use the shared prompt for stable appearance and style. Enter complete desired-image
tags, one image per nonempty line:

```text
standing, snowy park, beside a bench
sitting, snowy park, on a bench
---
standing, library, beside a bookshelf
reading, library, holding an open book
```

The first image in each scene starts from fresh noise. The following images edit
the preceding decoded result at **change strength**, initially `0.6`. Raise it
for larger changes, lower it to retain more composition, and insert `---` for a
new location or camera setup. No page layout, speech bubbles or comic style is
added automatically. Describe the wanted medium in the shared prompt.

**Reference weight** starts at `0.05`; stronger generic SDXL conditioning can
distort NoobAI instead of preserving identity.
Choose dimensions once for the whole sequence. Keep the seed fixed while tuning.
Set **image count** to `1` for a single preview, or `0` for all remaining actions.
For continuation, keep the complete action list, select the accepted preceding
image in **continuation**, and set **start image** to the next action number.
Starting inside a scene requires that image. Starting at a scene boundary creates
a fresh composition. Completed PNGs are saved immediately with grouped ordered
filenames and effective prompt/seed metadata; the App gallery populates on completion.

The default action list demonstrates two scenes with three images each: a robot
raises its hand in a snowy park, then opens a book in a library. Camera, location
and lighting tags repeat within each scene; `---` changes the scene. Keep shared
appearance tags in the shared prompt and describe the complete desired image on
each action line.

The selected run produces one combined GIF and no per-scene GIF files.
**GIF frame duration** defaults to 500 milliseconds per
keyframe. **Interpolation multiplier** defaults to 4x using the repository-pinned
RIFE engine; 1 disables interpolation. Each scene is smoothed separately and the
combined GIF retains hard scene cuts. Original PNGs remain available after
interruption; no partial GIF is published. Partial runs include only new images.

**Control strength** defaults to 0.2 in the App and guides same-scene edits from the preceding
image's Canny edges. Zero disables it. **Control end** defaults to 0.75. Lower the
strength for larger pose changes; edges can resist those changes. New scenes
start without the old scene's guide.

The App starts IP-Adapter reference influence at 0.05 for NoobAI. The generic
SDXL adapter can overwhelm this checkpoint at stronger weights. Shared appearance
tags and previous-image continuation provide additional continuity; exact identity
across fresh scenes and precise pose edits are not guaranteed.

After an authorized NoobAI deployment, use the infrastructure owner's
[in-place update](../../infra/image_gen/ansible/README.md) to install the node,
approved Canny checkpoint, App and RIFE engine. NoobAI defers this node's source
and workflow documents to that stage to keep startup below its payload limit.
The App can also be refreshed through the authenticated workflow API:

```sh
bazel run //infra/image_gen:comfy -- workflow save \
  "$PWD/bazel-bin/projects/comfyui_tools/workflows/noobai/noobai-image-sequence.app.json" \
  --name noobai-image-sequence.app --where local
```

The [node contract](nodes/image_sequence/README.md) describes partial runs,
interruption and drift limits. Offline checks verify transitions and saved media.
The [live acceptance](openspec/changes/add-noobai-image-sequence/acceptance.md)
records the completed two-scene GPU run and its remaining pose/identity limits.

### Reference and scene-composer versions

Two separate native Apps make the different generation strategies explicit:

| App                                  | Reused information                                                                 | Generation                                                               |
| ------------------------------------ | ---------------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| `noobai-reference-sequence.app.json` | Uploaded references, shared appearance text and preceding pixels within a scene    | Fresh scene starts, then previous-image edits; incrementing action seeds |
| `noobai-scene-composer.app.json`     | Fixed character/clothing/style text, per-scene setting and optional installed LoRA | Independent images with a fixed seed and a second latent refinement pass |

Both have matching editable `.json` graphs and six-image, two-scene examples.
Both save PNGs and one combined GIF. Their new presets start with interpolation
disabled; the reference preset also starts Canny at zero to avoid locking the
previous pose. Existing `noobai-image-sequence` files remain import-compatible.

The composer adapts the shared-prompt organization and two-pass sampler in the
operator-supplied Tache Scenes Composer, without its external prompt/utility
node packs or preset library. Enter a fixed character description, clothing and
render style. Enter one setting per scene on separate lines, or one setting for
all scenes. Enter one desired-image action per line, using `---` between scenes.
Settings are expanded into the effective prompts before sampling. The default
uses an adult woman in a red coat in a snowy park, then a library.

The optional LoRA selector lists only locally installed weights and defaults to
`(none)`. It installs nothing. The second pass uses the same model, prompt and
seed to resample each first-pass latent at lower denoise; it does not transfer
information between images. Set refinement steps to zero to disable it. Shared
tags and a fixed seed are not an identity lock, and refinement can worsen an image.
This version does not use uploaded reference images. Inspect individual PNGs
before enabling interpolation; smoothing cannot correct missing actions.

## Prompt sequences

The [Prompt Sequence node](nodes/prompt_sequence/README.md) uses ComfyUI's native
multiline text box. Select `line` and enter one prompt per line; blank and
whitespace-only lines are ignored, so they can separate prompts visually.
Select `blank line` only when each panel description spans several lines.
The sequence node also contains a shared `negative_prompt` text box. In the
Juggernaut, Flux Base and Qwen sequence workflows, its value
applies to every panel and is not split into lines or combined with positive
context. The generated sequence graphs connect it to the model's native negative
encoder, replacing the separate editable negative-prompt box. Juggernaut retains
its existing negative default; Flux and Qwen start empty. The encoder remains
in the graph to provide conditioning to the sampler.

Earlier saved nodes retain their first four input positions and
their prompt/filename output slots; an omitted negative prompt defaults to empty.
For Flux/Qwen, shared context describes image 1
as the character and image 2 as style. Juggernaut's native sequence uses shared textual character and style descriptions;
its style-sequence variant also reuses one reference image for style. Run a `prompt-sequence` workflow once for
sequential panels named `comic_<run UUID>_001`, `comic_<run UUID>_002`, and so on.
Each run has a unique prefix, including repeated runs with identical text. The
sequence remains one native job so its media outputs remain grouped. The native
sampler displays `Panel 2 of 3` alongside its step progress; current workflows
connect the sequence node's optional MODEL path to enable that counter. Sampling batch size
remains one; native list mapping reuses the shared conditioning. Panels do not inherit
previous-panel composition or story state. The node requires no external package.
The `JSON` separator remains available for lists saved by the former custom editor.

Use `chroma-hd-image-to-image-gif` to animate from an initial image, with one
prompt per generated frame in the same multiline box. Each decoded image is VAE-encoded for the
next img2img step. The result is one looping GIF and a native frame batch; start
at denoise 0.45 and 6 fps. This makes discrete chained images without motion
interpolation. The [GIF node](nodes/chained_img2img_gif/README.md) documents
timing and output behavior.

For controlled frame edits, choose `frame_source=original` and
`frame_seed_mode=fixed`; each pose then starts from the same reference and
noise. Chroma, HiDream and Qwen each have an `image-to-image-gif-masked` variant
with an adjustable soft edit rectangle. Qwen's GIF paths re-encode the preceding decoded image in the native
instruction-based edit encoder when `frame_source=previous`; original mode
keeps the initial reference for independent edits. Identical first/last pose
prompts provide a controlled loop boundary with original-reference/fixed-noise
sampling. These workflows still require visual motion evaluation.

For Qwen animation, start with `qwen-keyframe-edit` or
`qwen-keyframe-sequence`. Both edit a crop independently from the shared original,
then feather the edited crop onto the original scene. The crop's X/Y controls
also position the composite; include the old and desired hand positions inside
the crop. The separate movement rectangle controls which part of that crop is
blended back; include both old and intended hand positions inside it, away from
its feathered boundary. Qwen receives only the crop as its image reference.
The crop is sampled without a latent noise mask. The raw crop and final PNG
are both saved for inspection. Original scene pixels outside the movement area
are retained before image encoding, but the workflow cannot guarantee a correct
pose or alignment inside it. Reject edits that cross or misalign with the mask.

The initial crop/full-scene-reference trial produced redraw and seams. A later
single-reference localized press with a narrower composite was steadier;
a precise hand lift remained unreliable. These are inspectable still-edit tools,
not a validated natural-petting animation preset.

For rough-pose repair, use `qwen-pose-guide-edit`: image 1 is a roughly composed
target frame and image 2 is the original appearance reference. Supply equally
sized, aligned images with only one hand in the guide. The prompt assigns pose
to image 1 and appearance to image 2. The guide is scaled to about one megapixel
with its aspect preserved and dimensions rounded to eight-pixel steps. Sampling
starts at denoise 0.7 to retain visible guide geometry. The guide determines the
sampling canvas;
the output is restored to the original dimensions, and a feathered region is
composited onto the original. Shared X/Y and width/height controls keep the crop,
composite and mask matched; the region must include old and target
poses plus the exposed background. Inspect both the full repair and composite.
This supplies visible geometry but does not enforce the pose or temporal
consistency. Evaluate a single repaired pose before producing a sequence.

The four-edit rough-pose trial preserved a lifted hand and improved its anatomy
at denoise 0.7, but retained rough background placeholders. Full-denoise repair
returned to the original contact pose; reversing the reference order produced a
waving gesture. A separate background-only masked repair removed most gray fill
while preserving all pixels outside its mask, but left a cutout artifact and
changed the plush top. No finished frame or GIF passed visual acceptance. Avoid
large flat placeholders in guides; repairing anatomy and rebuilding vacated
background may require separate, independently inspected edits.

A later [photographic-guide trial](openspec/changes/add-parallel-evaluations/design.md#photographic-pose-guides-and-visible-motion)
uses image-editor-prepared background/hand layers and FFmpeg to construct target
poses, then refines every selected pose with Qwen and interpolates with RIFE.
Increasing the guide's hand travel from 14 to 154 pixels makes the movement
visible. Five inspected Qwen poses produce a 64-frame loop, but the forearm still
deforms unnaturally and textures vary. Lowering denoise from 0.7 to 0.45 removes
an extra digit in one pose but does not repair another pose's artificial skin
texture. These are case-specific findings, not new defaults or a validated
natural-animation preset. Guide preparation is a separate image-editing step;
the complete process is not solely Qwen and RIFE.

For structural references, first use `qwen-sketch-reconstruction` with a coherent
target sketch as image 1. It produces a raw full-frame reconstruction from that
single reference. Then use `qwen-structure-guide-edit`: image 1 remains the sketch
and image 2 adds the original appearance photograph. Both graphs share the
prompt, blank negative text, seed 42, 40 steps, CFG 4 and full denoise. Both
retain normal native semantic and VAE reference conditioning on both branches.
The sketch determines the sampling canvas. Compare raw outputs to distinguish
sketch comprehension from competition with the appearance reference.

The structural-reference workflow also restores the original dimensions and
saves an optional feathered composite onto the original. Inspect its raw output
first. Earlier exports used the opposite reference order; import the corrected
graph and assign the inputs using their titles. Native conditioning has no
separate ControlNet strength and does not guarantee pose or identity preservation.
Use the [prepared sketch and comparison procedure](evals/plush_petting/README.md).

The earlier structural-reference trial followed the moved-hand outline closely, but
raised the plush into the palm instead of preserving a gap. Full-scene edges
introduced a dark smear; geometry-only conditioning lost the pose, and a white
sketch produced a doubled hand. None of the four outputs passed finished-frame
acceptance. This remains an experimental guide workflow; no usable petting GIF
was obtained. Subsequent input review found conflicting head contours in the
scene sketch. In the [corrected-guide GPU trial](openspec/changes/add-parallel-evaluations/design.md#retained-gpu-trial),
the sketch-only default reconstructs a clear hand gap and a coherent scene.
Adding the original photo restores its appearance but introduces unwanted crown
folds. These results demonstrate useful guidance, not a finished petting animation.

Inspect each pose before loading selected, equally sized PNGs into
`qwen-keyframe-gif`. This model-free workflow combines three approved poses,
optionally reverses them back to the first pose, and uses the existing RIFE
engine at 1x (disabled), 2x, 4x or 8x. Reversal closes the loop; interpolation
smooths accepted poses and cannot repair malformed hands. The assembly node
also accepts any native IMAGE batch of 2 to 128 keyframes. Each input is retained
as a PNG. The older chained Qwen GIF workflows remain available for comparison.

The image-generation Ansible deployment packages the selected node sources and
workflow files through Bazel and installs them with the shared `comfyui_tools` role. This project
does not own the upstream ComfyUI CLI, server image or infrastructure state.

## Offline validation

```sh
bazel test //projects/comfyui_tools/nodes/prompt_sequence:prompt_sequence_test
bazel test //infra/image_gen/tf_setup:tf_tests.test_test
```

These checks exercise node behavior and Terraform's packaged workflow inputs.
The node uses the standard ComfyUI editor without a custom browser extension.
They do not access
a rental or establish compatibility with its frontend version or GPU image quality.

## Continuous experiments

The [experiment runner](docs/experiments.md) maintains a configurable owned job
buffer on an existing endpoint. External producers supply arbitrary native
workflows and meaningful variations. It preserves original media, workflow
snapshots, hashes and a resumable private journal, with no instance lifecycle
operations. The [ComfyUI experiments skill](skills/comfyui-experiments/SKILL.md)
covers queue observation, feedback and result review.

## Evaluations

The evaluator runs independent, bounded rental workloads using the existing
Vast SDK and ComfyUI CLI. Reusable [cases](evals/README.md) include the supplied
thorn-landscape artwork prompt/reference and the plush-petting GIF reference.
[Infrastructure composition](../../infra/image_gen/evals/README.md) owns the
run matrix, credentials, limits and final/manual cleanup targets. No personal
Terraform state is read or changed. All executable code stays in this project.

Matrix groups expand hardware/model/workflow/case coordinates. All compatible
workflows, cases and seeds sharing a hardware/model profile reuse one rental.
Each rental starts as soon as its own endpoint is ready, keeps native jobs queued
ahead of retrieval and releases capacity as soon as all its configured outputs
are downloaded. Local media validation follows
deletion. Inputs, native workflows, fixed seeds, model definitions, timings,
output hashes and cleanup status are retained in private run artifacts.
Artistic quality and reference adherence remain separate, unscored review
criteria; execution success is not a quality verdict.

For repeated iterations, retain the first run with `--cleanup=false`, then pass
`--reuse-run /absolute/path/to/that/run --cleanup=false`. Reuse requires matching
declared installations, recorded provisioning fingerprints and exact live
ownership; it never allocates a fallback. Fingerprints cover installed node
source and requirements contents, runtime settings and the startup template.
Runs without fingerprints cannot be reused.
Each iteration records new cases, workflows and sampling in a separate result
directory while preserving the original instance and prior results. Retention
continues billing outside controller deadlines; release requires explicit cleanup.

Loopback E2E checks exercise matrix pooling, offer reuse, namespace exclusion
across output directories, duplicate-seed rejection, parallel completion,
failures, interruption, lost
creation acknowledgments, cleanup recovery and the real authenticated CLI's
upload/conversion/native-queue/download path, plus descendant cancellation. They never rent capacity or prove
GPU model quality. Runtime composition currently supports the maintained
CLIPTextEncode-based SDXL and HiDream graphs; unsupported workflow bindings
fail before allocation. The authenticated CLI bridge requires Linux procfs.
