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

The workflow files, including the [Flux panel sequence](workflows/flux/flux-klein-prompt-sequence.json), contain single-reference, character/style and
prompt-sequence graphs for FLUX.2 klein Base and Qwen-Image-Edit, plus the
existing Chroma Flash graphs. Load a matching JSON in ComfyUI or pass it to the
authenticated infrastructure CLI. Each graph requires its matching models and
nodes; the rental's [model profile](../../infra/image_gen/tf_setup/README.md#download-models)
selects which files to install. Keeping a workflow here does not install its
dependencies or establish runtime acceptance.

## Prompt sequences

The [Prompt Sequence node](nodes/prompt_sequence/README.md) shows one text box per panel
and an empty box at the end. Filling that box adds another; clear a box to omit
it from generation. Blank lines within a box stay part of that panel. Shared context describes image 1
as the character and image 2 as style. Run a `prompt-sequence` workflow once for
sequential panels named `comic_001`, `comic_002`, and so on. Sampling batch size
remains one; native list mapping reuses the references. Panels do not inherit
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
