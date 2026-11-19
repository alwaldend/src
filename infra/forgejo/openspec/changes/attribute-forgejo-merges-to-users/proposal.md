## Why

Forgejo records the instance as the committer of signed pull request merges.
The Git committer field should identify the account that performs the merge.

## What Changes

- Use the collaborator trust model for repositories that inherit the instance default.
- Retain SSH signing with the Forgejo instance key.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `infra-forgejo` adds a requirement for signed merge attribution.

## Impact

The Ansible configuration changes the default signature verification policy
for repositories that inherit it. Repository-specific settings take precedence.
