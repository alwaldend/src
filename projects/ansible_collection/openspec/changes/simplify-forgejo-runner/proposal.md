## Why

The runner role exposes independently configurable internal paths and parses
registration failures into categories. These interfaces have no active callers
in this checkout and obscure the install, register and service lifecycle.

## What Changes

- Replace registration error classification with a fixed, credential-safe error.
- **BREAKING**: derive binary, configuration, home, work and registration paths
  from `forgejo_runner_dir`; use the packaged binary and standard systemd unit path.
- **BREAKING**: consolidate advanced runner settings in `forgejo_runner_config`;
  retain labels, capacity, service account and package-management inputs.
- Remove Podman from the default host-executor prerequisites.
- Preserve one-time registration, bounded execution, stdin token input, ownership
  and change-triggered service restarts.
- Document the authentication boundary and exercise registration and rendering
  locally without infrastructure or real credentials.

## Capabilities

### New Capabilities

- `forgejo-runner`: Minimal installation interface and credential-safe registration lifecycle.

### Modified Capabilities

None.

## Impact

Changes affect `projects/ansible_collection` and the narrow `tools/BUILD.bazel`
test-consumer allowlist entry needed by its new acceptance package. No new dependencies,
deployment, workflow activation or live identity changes. `infra/dev` currently
does not enable the runner role. Its old runner-specific identity and registration
lookup were already retired; shared Vault-backed deployment authentication and
the service-wide Forgejo OIDC login helper remain outside this role.
