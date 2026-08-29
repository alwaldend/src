---
name: repo-secrets
description: Safely add, route, rotate, or review secrets and authentication in this repository using Vault, AppRole, the al.lua DSL, injectors, and Bazel wrappers. Use whenever work involves credentials, tokens, private keys, sensitive Terraform values, Ansible secrets, Vault policy, or secret-bearing configuration.
---

# Handle Secrets

## Follow the repository trust boundary

- Store secret values in Vault, never in Git, Terraform source, Ansible vars,
  command history, test fixtures, logs, screenshots, or PR text.
- Commit only secret **references**: Vault mount/path, environment variable name,
  plugin label, policy, role, or non-sensitive metadata.
- Treat `infra/` as private infrastructure material as required by
  `infra/README.md`; do not publish or reuse it as build input.
- Treat an unexpected credential already in the tree as compromised. Do not
  repeat it in output. Stop exposing it, preserve evidence minimally, and tell
  the user it needs revocation/rotation.

## Reuse the injection flow

1. Find the nearest component `al.lua` and its parent `al_config` dependency.
2. Reuse `lib.vault_auth` for AppRole authentication.
3. Use `lib.plugin_call` with the existing component plugin (`injector`,
   `tf_backend`, provider login, or another established plugin).
4. Select calls through labels such as `tf=main`, `tf=setup`, or `ansible=1`.
5. Ensure the corresponding plugin binary is present in the Bazel target's
   `data` and its label is present in `run_args`.
6. Inject the minimum values and grant the minimum Vault policy needed.

Follow sibling `al.lua` files exactly for the DSL shape. Keep AppRole names and
paths aligned with the component label; do not silently share a broader role.

## Prevent disclosure while working

- Do not print environment variables, rendered templates, Terraform state,
  Vault responses, token files, or generated credentials.
- Use placeholders in examples. Prefer stdin or `@file` input over a literal
  value on a command line, and keep temporary secret files outside the repo with
  restrictive permissions.
- Do not run Vault writes, token generation, unseal, apply, or deployment merely
  to test a code change.
- Review diffs with `git diff --check` and a targeted search for newly introduced
  private keys, literal passwords/tokens, and suspicious high-entropy values.
  Do not use a scan that echoes suspected values into the final response.

## Validate structure

Build or test the narrow Bazel targets that package the changed `al.lua`, Vault
policy, Terraform, or Ansible configuration. Explain when live authentication is
required and leave the live operation to an authorized operator unless the user
explicitly requests it and the environment is clearly non-production.
