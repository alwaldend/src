## Why

The injector mixes generic graph evaluation with hardcoded Vault resource types.
The user requests the previously deferred built-in engine and plugins that
advertise the resources they provide.

## What Changes

- Move generic injection, templates, files, environment exports and process ownership into AL.
- Add provider capability advertisements and generic resource-fetch RPCs.
- Keep Vault-backed operations in a standalone provider plugin.
- Preserve the public res_<type> authoring API, explicit graph semantics and cleanup.
- BREAKING: use generic type/args resource definitions and replace injector plugin registration with Vault provider registration.

## Capabilities

### New Capabilities

- builtin-resource-injection: core graph evaluation and provider protocol.

### Modified Capabilities

None.

## Impact

AL core/plugin APIs, the Vault plugin, Bazel registrations and existing process
fixtures. No new dependencies or live infrastructure operations.
