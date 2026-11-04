## Why

The user requested a shared internal home for the Vast.ai API client.

## What Changes

- Move schema generation and bounded API adaptation into internal/vastai_sdk.
- Keep Terraform-specific schemas, state, and fixtures with the provider.
- Preserve the pinned upstream schema and generator without source duplication.

## Capabilities

### New Capabilities

- `shared-vastai-client`: Shared generated operations and validated API transport.

## Impact

Updates first-party imports, Bazel visibility, documentation, and owner discovery.
No independent product release, landing page, or external dependency is added.
