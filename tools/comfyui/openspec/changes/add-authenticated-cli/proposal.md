## Why

Remote ComfyUI workflow execution currently needs task-specific HTTP scripts.
The user requested the official CLI and a rental target with automatic authentication.

## What Changes

- Package the approved Comfy-Org comfy-cli 1.25.0 with an independent Python lock.
- Provide a short-lived HTTPS authentication bridge for its HTTP-only local target.
- Keep Vault and rental endpoint discovery in the infrastructure owner.

## Capabilities

### New Capabilities

- `authenticated-cli`: Packaged CLI execution against an authenticated HTTPS ComfyUI.

## Impact

Changes tools/comfyui packaging and adds a first-party command using existing Go dependencies.
No host packages or persistent services are installed. The approved CLI introduces its locked transitive packages.
