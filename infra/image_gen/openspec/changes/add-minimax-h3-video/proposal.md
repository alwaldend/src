## Why

The image-generation rental needs a selectable MiniMax H3 Video model for
animating keyframes on the 48 GB GPU profile.

## What Changes

- Add public, revision-pinned H3 Base weights to the model registry.
- Stage prepared text-to-video and first/last-frame workflows.
- Expose invocation-time profile selectors for provisioning and installation.

## Capabilities

### New Capabilities

- `minimax-h3-video-profile` provides selectable H3 Base video inference assets.

### Modified Capabilities

None.

## Impact

The existing image-generation Ansible playbook installs the selected weights
and the generated workflow assets owned by ComfyUI tools.
