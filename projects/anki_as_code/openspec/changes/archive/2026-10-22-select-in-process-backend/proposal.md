## Why

Callers need to select native in-process execution through existing address options while preserving collection preparation, archive publication and source reservations.

## What Changes

- Accept inprocess:// and inprocess+cgo:// labels for managed collections and native API calls.
- Allow any selected client alongside a managed API path, including TCP/Unix servers with shared paths. Keep borrowed connections exclusive and ownership internal to client implementations.
- Delegate address interpretation to an address-only library constructor, with native constructors registered by the executable's optional package import.
- Preserve configured subprocess defaults, external TCP/Unix behavior and the Client API and optional build composition merged through PR #186.
- Use configured factories for private worker registrations and defaults, use the shared asynchronous completion interface for source reservations, and remove discarded scope callback change results.

## Capabilities

### New Capabilities

- `in-process-selection`: Application and CLI selection of direct native execution.

### Modified Capabilities

None.

## Impact

Application scopes, API source validation, CLI help and native E2E. The native implementation is owned by projects/anki_grpc/openspec/changes/add-in-process-backend. This candidate targets master after the merge of PR #186.
