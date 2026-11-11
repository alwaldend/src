---
title: ComfyUI tools
description: Reusable ComfyUI custom nodes and image-generation workflows
---

This project owns reusable ComfyUI nodes and workflow JSON. Each node has its
own directory under nodes/, with its implementation, test and
README. Workflow definitions are grouped by model under workflows/. Deployment,
model downloads, credentials and rental state belong to
[image-generation infrastructure](../../infra/image_gen/README.md). The
[repository CLI and local launchers](../../tools/comfyui/README.md) remain in
`tools/comfyui`.

## Workflows

The workflow files, including the [Flux panel sequence](workflows/flux/flux-prompt-sequence.json), contain single-reference, character/style and
prompt-sequence graphs for FLUX.2 klein Base and Qwen-Image-Edit, plus
Juggernaut XL v9 text-to-image, img2img and prompt-sequence graphs and archived
Chroma Flash examples. All files use `<profile>-<purpose>.json`; the saved
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

## Workflow source and generation

The maintained sources are [workflows.jsonnet](workflows/workflows.jsonnet),
model-specific `workflows.libsonnet` files and shared graph/node libraries in
`workflows/internal/`. The sequence transformation is shared across model families.
Connections are declared once; generation derives port metadata and reciprocal
link lists. The public rental configuration supplies model filenames from the first matching
download of each required kind. Additional downloads remain available without
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
sequential panels named `comic_001`, `comic_002`, and so on. Sampling batch size
remains one; native list mapping reuses the shared conditioning. Panels do not inherit
previous-panel composition or story state. The node requires no external package.
The `JSON` separator remains available for lists saved by the former custom editor.

Terraform packages the node source and selected workflow files through their
Bazel filegroups and installs them with Vast's native provisioner. This project
does not own a CLI, server image or infrastructure state.

## Offline validation

```sh
bazel test //projects/comfyui_tools/nodes/prompt_sequence:prompt_sequence_test
bazel test //infra/image_gen/tf_setup:tf_tests.test_test
```

These checks exercise node behavior and Terraform's packaged workflow inputs.
The node uses the standard ComfyUI editor without a custom browser extension.
They do not access
a rental or establish compatibility with its frontend version or GPU image quality.
