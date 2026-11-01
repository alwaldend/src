## Why

Run Codex CLI with local inference and local tools on a dedicated GPU VM,
without sending prompts, workspace content, or tool calls to hosted services.
Prepare declarative infrastructure now; deployment is explicitly excluded.

## What Changes

- Add a standard Terraform provisioning stage and packaged Ansible deployment.
- Declare one Broadwell V100 VM with 8 vCPUs and 96 GiB system RAM.
- Encrypt its boot disk with a component-owned, rotating KMS key.
- Configure loopback-only llama.cpp and a coding account shared with T3.
- Disable hosted Codex tools and pin its local inference provider.
- Document artifact inputs, prerequisites, offline checks, and future acceptance.

## Capabilities

### New Capabilities

- `local-inference-host`: Encrypted GPU infrastructure and local Codex runtime.

### Modified Capabilities

None.

## Impact

Owns `infra/local_inference`. Uses the repository's Yandex provider and
Vault-backed AL workflow. Codex and llama.cpp are the requested runtimes;
model files are operator-supplied, checksum-verified artifacts. No cloud
resources, credentials, or current-host configuration are changed.
