## Why

The library's manager API duplicates application scope composition. Its sole production consumer already chooses a fresh stdin/stdout backend for each collection.

## What Changes

- **BREAKING** Remove Manager, ManagerSTD, ManagerCollection and their constructors.
- Remove the shared manager helper and tests for those removed APIs.
- Retain BackendSTD, BackendCollection and explicit Start/Connection/Stop ownership.
- Move scope orchestration and its failure coverage to the application store.

## Capabilities

### New Capabilities

- `explicit-backends`: Reusable explicit backend lifecycle without an application scope-manager layer.

### Modified Capabilities

None.

## Impact

The paired anki_as_code change migrates the only production caller. Existing direct backend and collection E2E stay in this project. Native APIs, server protocol and dependencies remain unchanged.
