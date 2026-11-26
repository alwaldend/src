## Why

The image-generation catalog needs a selectable MOSS-Image-heal installation.

## What Changes

- Register the public INT8 ConvRot checkpoint and reuse the MOSS encoder and VAE.
- Package the Heal workflows and Apps through the existing Ansible installation.

## Capabilities

### New Capabilities

- `moss-image-heal-profile`: Select and install the MOSS-Image-heal bundle.

## Impact

The model registry selects the assets. ComfyUI tools owns the generated graphs.
