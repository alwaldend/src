## Why

Keep application command composition independent of the linked server adapter before adding another implementation.

## What Changes

- Bind the native adapter in server_native.go and retain shared application composition.
- Keep one public executable target, backed by an explicitly named native target.
- Describe signal ownership and forward server arguments unchanged.

## Capabilities

No specification changes: this extracts native composition and build preparation while preserving runtime behavior. `skip_specs: true`.

## Impact

Application command composition and existing CLI workflows; collection behavior is unchanged.
