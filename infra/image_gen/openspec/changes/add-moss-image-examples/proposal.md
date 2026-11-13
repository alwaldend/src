## Why

The operator requested MOSS-Image examples and a personal 48 GB ComfyUI rental.
The model catalog currently has no MOSS profile.

## What Changes

- Add and select a pinned public MOSS-Image bundle with SFT/Base INT8 models,
  Qwen3.5 INT8 text encoding, the small FLUX.2 VAE and edit/control adapters.
- Stage the six native example workflows owned by ComfyUI tools.
- Deploy through the existing reviewed Terraform rental workflow on `48gb`.

## Capabilities

### New Capabilities

- `moss-image-profile`: Reproducible MOSS model provisioning and example staging.

### Modified Capabilities

None.

## Impact

Changes `config.json` and component documentation. Reuses the existing runtime,
Vault login, native downloader, 100 GB disk and USD 1/hour ceiling. Workflow
source and its corresponding change belong to `projects/comfyui_tools`.
