## Why

Run Codex CLI with local inference and local tools on a dedicated GPU VM,
without sending prompts, workspace content, or tool calls to hosted services.
Prepare declarative infrastructure now; deployment is explicitly excluded.

## What Changes

- Package Ansible deployment for a separately provisioned GPU host.
- Deliver VM provisioning in the independent Terraform PR.
- Configure loopback-only llama.cpp and a coding account shared with T3.
- Disable hosted Codex tools and pin its local inference provider.
- Document artifact inputs, prerequisites, offline checks, and future acceptance.

## Capabilities

### New Capabilities

- `local-inference-host`: Encrypted GPU infrastructure and local Codex runtime.

### Modified Capabilities

None.

## Impact

Owns the host deployment under `infra/local_inference/ansible`, using
Vault-backed SSH. Terraform provisioning is delivered in a separate PR. Codex and llama.cpp are the requested runtimes;
model files are operator-supplied, checksum-verified artifacts. No cloud
resources, credentials, or current-host configuration are changed.
