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
NoobAI's publisher restricts commercial use, including generated products;
the [model configuration](../../infra/image_gen/tf_setup/README.md#additional-model-options)
records the licenses and component choices.

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

The `moss-image` family supplies six examples, adapted from the
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

These are editable examples, not guarantees of exact edits, layout or reference
fidelity. Actual GPU execution and its limitations belong to the
[deployment record](../../infra/image_gen/openspec/changes/add-moss-image-examples/design.md).

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

Terraform packages the node source and selected workflow files through their
Bazel filegroups and installs them with Vast's native provisioner. This project
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
