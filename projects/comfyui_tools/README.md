---
title: ComfyUI tools
description: Reusable ComfyUI custom nodes and image-generation workflows
---

This project owns reusable ComfyUI nodes and workflow JSON. Each node has its
own directory under nodes/, with its implementation, test, browser script and
README. Workflow documents are grouped by model under workflows/. Deployment,
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

## Prompt sequences

The [Prompt Sequence node](nodes/prompt_sequence/README.md) shows one text box per panel
and an empty box at the end. Filling that box adds another; clear a box to omit
it from generation. Blank lines within a box stay part of that panel. For Flux/Qwen, shared context describes image 1
as the character and image 2 as style. Juggernaut's native sequence uses shared textual character and style descriptions;
its style-sequence variant also reuses one reference image for style. Run a `prompt-sequence` workflow once for
sequential panels named `comic_001`, `comic_002`, and so on. Sampling batch size
remains one; native list mapping reuses the shared conditioning. Panels do not inherit
previous-panel composition or story state. The node requires no external package.
Existing line and blank-line lists remain supported. The widget stores edited
lists as JSON; without its browser extension, the standard text box and explicit
JSON separator remain available.

Terraform packages the node source and selected workflow files through their
Bazel filegroups and installs them, including the node's browser script, with Vast's native provisioner. This project
does not own a CLI, server image or infrastructure state.

## Offline validation

```sh
bazel test //projects/comfyui_tools/nodes/prompt_sequence:prompt_sequence_test
bazel test //projects/comfyui_tools/test/browser:prompt_sequence_test
bazel test //infra/image_gen/tf_setup:tf_tests.test_test
```

These checks exercise node behavior, a browser fixture for ComfyUI's DOM-widget
contract and Terraform's packaged workflow inputs. The browser check preserves
a screenshot and serialization receipt in its test outputs. They do not access
a rental or establish compatibility with its frontend version or GPU image quality.
