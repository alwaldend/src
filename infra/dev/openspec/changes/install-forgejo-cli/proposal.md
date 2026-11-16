## Why

The user requested a host installation of the repository-pinned fj CLI and
launch of its browser login. Host provisioning must reproduce the installation.

## What Changes

- Package the existing pinned fj binary and install it in the managed dev1
  user's `.local/bin` through a separately selectable Ansible task.
- Provide a local-user installation entry point using the same packaged playbook
  and task, without unrelated service changes or Vault authentication.
- Launch native fj browser login only as the user requested; no credentials or
  OAuth state are stored in source or Ansible variables.

## Capabilities

### New Capabilities

- `forgejo-cli-installation`: Declarative, idempotent installation of the pinned
  Forgejo CLI for the development host's managed account.

### Modified Capabilities

None.

## Impact

`infra/dev` Ansible packaging, playbook and documentation. Reuses an existing
dependency pin, exposing its standalone archive in the dependency's existing
`include.MODULE.bazel`. User authorization covers installation on the current dev1 host
and starting its native browser login; it does not authorize unrelated deployment.
