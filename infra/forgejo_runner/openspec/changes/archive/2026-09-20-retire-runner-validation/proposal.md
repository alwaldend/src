## Why

The user requested full removal of the completed runner acceptance harness from the PR. Normal CI will build and test the repository through `//tools/ci`.

## What Changes

- Remove the smoke package, CI publication helper/playbook, and validation-branch configuration.
- Narrow the declared Vault role to the protected default branch and retain the existing issuer, workflow, event, repository, audience, TTL and policy restrictions.
- Update current runner documentation without rewriting historical acceptance records.

## Capabilities

### New Capabilities

### Modified Capabilities

- `secure-runner-ci`: retain signed CI identity without a permanent runner smoke job or validation ref.

## Impact

Runner packaging/docs and Vault desired policy. Related repository change: `infra/src/openspec/changes/repository-ci-entrypoint`. No live policy deployment or branch deletion is included.
