## Why

ComfyUI needs the same ingress and host-local HTTPS access pattern as T3 Code.

## What Changes

- Declare ComfyUI DNS names for local access, the ingress backend, and public access.
- Add a local Traefik route to ComfyUI on loopback port 8188.
- Document the routing contract without installing or deploying ComfyUI.
- Add a ComfyUI architecture diagram beside the T3 Code diagram, using the current native site style for both default renders and showing local CPU generation.

## Capabilities

### New Capabilities

- `comfyui-routing`: Host Bot DNS and HTTPS forwarding to the local ComfyUI server.

### Modified Capabilities

None.

## Impact

Host Bot DNS declarations, Ansible variables and Traefik configuration. The
coordinated ingress change is `infra/ingress/openspec/changes/add-comfyui-routing`.
No new dependencies or live operations.
