## Why

The local inference deployment embeds reusable llama.cpp host setup. Extract
that behavior into its own role and let Bazel own the existing runtime pins.

## What Changes

- Move account, runtime and CUDA installation, verified model staging, service,
  restart, and readiness into the llama role.
- Accept caller-packaged archives rather than downloading in the role.
- Keep image provisioning and Codex integration in the deployment.

## Capabilities

No behavioral delta; this extracts existing host configuration.

## Impact

Ansible collection and infra/local_inference. No deployment or new dependency.
