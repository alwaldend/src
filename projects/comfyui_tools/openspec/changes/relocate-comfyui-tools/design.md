## Context

The operator requested moving ComfyUI nodes and workflows into a reusable
project. tools/comfyui continues to own repository executables and dependencies.
infra/image_gen continues to own rental configuration and model profiles.

## Decisions

Use nodes/<name> for each first-party node and workflows for workflow documents.
Each node owns its implementation, test, browser script and README. Workflows
are grouped by model under workflows/flux, workflows/qwen and workflows/chroma,
with their native filenames and JSON contents retained. Native
installation remains under ComfyUI/custom_nodes/<name>.
The relocation preserves the node registry name and workflow JSON bytes.
The separate add-prompt-sequence-editor change owns the requested UI extension.
Terraform reads the relocated files through declared Bazel data dependencies;
its relative paths resolve within the same staged repository tree.

Use one asset_consumers package group for the project aggregate and rental
Terraform package. Register documentation and landing content through the
existing shared site; no separate deployment or release target is introduced.

## Validation

Run the existing node behavior test and offline Terraform mocked-provider
checks, build the site, inspect representative HTML and check generated
inventories are stable. Verify tools/comfyui and its external references remain
unchanged from the previous candidate. No rental access or apply is authorized.

The subsequent normalize-model-workflows change owns the requested filename
normalization and added Juggernaut assets. The relocation's byte-preservation
guarantee applies to the existing Flux, Qwen and Chroma graph contents.
