## Why

Export and reconciliation duplicate native client construction and enumeration. Extract one native-resource boundary so both operations can share reads without sharing configuration or reconciliation policy.

## What Changes

- Add internal/resource.Store with generated resource clients on a borrowed prepared connection.
- Provide a fresh Snapshot containing upstream deck, note-type, note and card messages plus existing preset names.
- Keep state bindings, config projection, planning, write ordering and native lifecycle outside resource access.
- Publish the library independently against master; consumer integration remains in #128 and #160.

## Capabilities

### New Capabilities

- native-resource-access: Reusable generated resource clients and complete native snapshots without lifecycle or policy ownership.

### Modified Capabilities

None.

## Impact

Add a project-internal Go package and native E2E using existing Anki schemas and backend dependencies. No dependency additions, SQL, custom CRUD protocol or storage changes.
