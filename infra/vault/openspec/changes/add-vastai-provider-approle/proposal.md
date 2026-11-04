## Why

The Vast.ai provider project needs its own authentication and state namespace.
The bootstrap infrastructure must be able to deliver its generated team token
to that namespace without sharing the bootstrap token with the project.

## What Changes

- Declare src_projects_vastai_terraform_provider with operator bootstrap access.
- Reuse the role's standard secret namespace policy for generated token storage;
  the caller selects this role directly, with no cross-role policy grant.

## Capabilities

### New Capabilities

- `vastai-provider-identity`: Dedicated provider project identity and token handoff.

## Impact

Uses the existing AppRole module and KV mount. It adds no general infra group,
cloud policy, host configuration, or committed secret values.
