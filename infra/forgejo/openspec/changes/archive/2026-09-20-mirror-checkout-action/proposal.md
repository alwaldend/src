## Why

The user requested an organization-owned `com_github_actions_checkout` mirror and use of that mirror in repository CI, preserving the shared Bazel CI entry point.

## What Changes

- Declare the public Forgejo pull mirror in the shared repository catalog.
- Extend Forgejo's catalog projection and resource wiring for upstream pull mirrors.
- Create the mirror through the owning Terraform workflow and verify the existing action SHA is available.
- Switch the checkout action URL while preserving its SHA and the `ci` job's local action.

## Capabilities

### New Capabilities

### Modified Capabilities

- `infra-forgejo`: declarative organization-owned pull mirrors for pinned CI actions.

## Impact

Repository catalog and tests, Forgejo service Terraform, workflow URL and documentation. Authorized live scope is creation of this mirror and its normal settings. No unrelated infrastructure mutation or policy deployment is included.
