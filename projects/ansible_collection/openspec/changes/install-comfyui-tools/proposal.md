## Why

ComfyUI application installation is split between Terraform startup commands
and a NoobAI-specific update playbook. A reusable role should own the selected
workflows, model files (including LoRAs), node packs and locked dependencies.

## What Changes

- Add the caller-configured `comfyui_tools` role to the shared collection.
- Install pinned model/node inputs and packaged workflows at configurable paths.
- Preserve aggregate-idle checks, workflow-copy opt-out and notified restarts.

## Capabilities

### New Capabilities

- `comfyui-assets`: Declarative installation of selected ComfyUI assets.

## Impact

Collection role and local acceptance fixtures; image-generation deployment
consumes it. No dependency additions or live operations.
