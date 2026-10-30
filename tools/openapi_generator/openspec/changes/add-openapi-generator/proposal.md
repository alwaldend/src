## Why

Pinned OpenAPI generation is requested and has no existing repository implementation.

## What Changes

- Expose the requested OpenAPI Generator CLI and an offline generation rule.

## Capabilities

### New Capabilities

- `openapi-generation`: Expose the requested OpenAPI Generator CLI and an offline generation rule.

### Modified Capabilities

None.

## Impact

tools/openapi_generator; shared Bazel declarations and documentation integration. The SDK uses tools/openapi_generator and reuses the article uploader credential reference and signer without adding secret values.
