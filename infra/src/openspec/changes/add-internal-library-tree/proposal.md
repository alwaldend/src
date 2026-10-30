## Why

The user requested a repository-only home for shared runtime libraries. The existing products and developer-tools trees do not express that boundary.

## What Changes

- Define internal/ for shared first-party runtime libraries with repository-internal visibility.
- Move the new X SDK integration into internal/x_sdk and retain generated client source as Bazel outputs.
- Integrate reference documentation, transitive visibility mapping and validation.

## Capabilities

### New Capabilities

- `internal-libraries`: Shared repository-only runtime library ownership.

### Modified Capabilities

None.

## Impact

Root documentation and build integration, internal/, tools generation contract, and the new X SDK integration. No infrastructure is deployed. SDK behavior is maintained in internal/x_sdk/openspec/changes/add-x-sdk.
