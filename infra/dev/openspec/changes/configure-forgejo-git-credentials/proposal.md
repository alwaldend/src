## Why

The dev1 Git configuration has a GitHub credential helper. Forgejo Git
operations over HTTPS need a helper that can reuse the existing native fj login.

## What Changes

- Install a Git credential helper and scope it to the Forgejo HTTPS origin.
- Use the managed user's existing fj login for alwaldend/src.

## Capabilities

### New Capabilities

- `forgejo-git-credentials` configures authentication for Git over HTTPS.

### Modified Capabilities

None.

## Impact

The dev1 Ansible playbook packages the helper and references it in Git config.
