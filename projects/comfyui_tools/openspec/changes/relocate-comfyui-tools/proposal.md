## Why

Custom nodes and workflow JSON are reusable ComfyUI assets, currently owned by
one rental component. Move them into projects/comfyui_tools so infrastructure
consumes their source without owning it.

## What Changes

- Move nodes into projects/comfyui_tools/nodes/<name> and workflow JSON into
  projects/comfyui_tools/workflows. Each node directory owns its source, test,
  browser script and README.
- Update Terraform file paths, Bazel data inputs and documentation ownership.
- Register the assets project and its landing page.
- Keep the CLI and local launchers in tools/comfyui.

## Capabilities

No behavioral delta; this is a source ownership and packaging change.

## Impact

No dependencies, model settings, state addresses or deployed systems change.
Offline checks must verify the relocated node and packaged Terraform inputs.
