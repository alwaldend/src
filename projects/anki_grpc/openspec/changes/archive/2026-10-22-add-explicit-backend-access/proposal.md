## Why

CLI callers need a uniform backend interface for a bundled stdio worker and an existing TCP/Unix server. Package persistence belongs to the collection owner rather than each application callback.

## What Changes

- Add a gRPC connection backend with explicit endpoint configuration and client-only shutdown.
- Add opt-in `WriteBack` to collection constructor options; its zero value leaves source archives unchanged. When enabled, shutdown writes packages back with media regardless of earlier request errors.
- Preserve recoverable restores on publication failure and report cleanup errors alongside operation errors.
- Remove the server runtime dependency from the reusable backend library; application bundles supply exact executable paths.

## Capabilities

### New Capabilities

- `backend-access`: Explicit network connections and writable collection lifetime.

### Modified Capabilities

None.

## Impact

Go backend library, native E2E and documentation. The server protocol and upstream native library stay unchanged. The paired anki_as_code change consumes these APIs and owns bundling. No new dependencies.
