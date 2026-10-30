## Why

X Go SDK is requested and has no existing repository implementation.

## What Changes

- Fetch the authenticated X OpenAPI schema into source and generate a Go client.

## Capabilities

### New Capabilities

- `x-sdk`: Fetch the authenticated X OpenAPI schema into source and generate a Go client.

### Modified Capabilities

None.

## Impact

internal/x_sdk; shared Bazel declarations and documentation integration. The SDK uses tools/openapi_generator and reuses the article uploader credential reference and signer without adding secret values.
