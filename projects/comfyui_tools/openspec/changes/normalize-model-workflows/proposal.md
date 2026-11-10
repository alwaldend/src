## Why

The operator requested Juggernaut as the next default, a reference-style workflow,
and consistent workflow names. Reusable workflow graphs belong to this project;
infrastructure selects the profile and installs its declared dependencies.

## What Changes

- Name every workflow <profile>-<purpose>, preserving the existing model folders.
- Add Juggernaut XL v9 text-to-image, image-to-image and prompt-sequence graphs.
- Add single-prompt and prompt-sequence style-reference graphs using the previously
  approved RES4LYF revision and its SDXL model patcher.
- Update configuration references and documentation without accessing a rental.

## Capabilities

### New Capabilities

- `model-workflows`: consistent native graph names and Juggernaut generation paths.

## Impact

Existing Flux, Qwen and Chroma graph contents remain unchanged. No new external
software is approved here. Infrastructure owns restoration of the already approved
checkpoint, node pack and locked requirements. GPU execution and style quality
remain unverified until a separately authorized rental test.
