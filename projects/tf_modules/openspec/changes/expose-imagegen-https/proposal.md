## Why

Image generation needs SSH-only and password-authenticated HTTPS access with a 32 GiB VRAM floor.

## What Changes

Retain ssh_only by default; exposed publishes only 7860/tcp through a dedicated TLS and Basic-auth Caddy configuration. The selected runtime stays on loopback. Reject unsafe template options and weak credentials; disable outbound tunnels and synchronization while retaining loopback-only portal, Jupyter, and TensorBoard services.

## Capabilities

### New Capabilities

- `imagegen-access`: Image generation needs SSH-only and password-authenticated HTTPS access with a 32 GiB VRAM floor.

## Impact

Owner-local source, packaging, documentation, and offline validation. Reuses the existing pinned Terraform providers; the specifically approved Vast ComfyUI container is the only new runtime dependency, with no live deployment.

The requested follow-up separates the base rental from Forge and ComfyUI
profiles. Keep runtime-specific model paths and settings in their wrappers while
sharing exposure and offer constraints. A new ComfyUI container requires specific
approval; this source refactor does not authorize another rental.
