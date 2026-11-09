## Why

Image generation needs SSH-only and password-authenticated HTTPS access with a 32 GiB VRAM floor.

## What Changes

Retain ssh_only by default; exposed publishes only 7860/tcp through a dedicated TLS and Basic-auth Caddy configuration. Forge stays on loopback. Reject unsafe template options and weak credentials; disable outbound tunnels and synchronization while retaining loopback-only portal, Jupyter, and TensorBoard services.

## Capabilities

### New Capabilities

- `imagegen-access`: Image generation needs SSH-only and password-authenticated HTTPS access with a 32 GiB VRAM floor.

## Impact

Owner-local source, packaging, documentation, and offline validation. Reuses the existing external-data and other pinned Terraform providers; no new external library or live deployment.
