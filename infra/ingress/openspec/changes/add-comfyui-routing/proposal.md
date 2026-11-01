## Why

ComfyUI on Host Bot needs the same authenticated ingress path as T3 Code.

## What Changes

- Add ComfyUI to the existing Traefik service map using its dedicated site-local backend.
- Preserve the existing mTLS, HTTPS, and backend SNI behavior.

## Capabilities

### New Capabilities

- `comfyui-ingress`: Forward ComfyUI requests to Host Bot over HTTPS.

### Modified Capabilities

None.

## Impact

Ingress Ansible service declarations. Host-local routing and DNS are owned by
`users/simeonwarren/host_bot/openspec/changes/add-comfyui-routing`.
No deployment or dependency changes.
