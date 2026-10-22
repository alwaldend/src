## Why

Connection construction belongs in the backend library, but importing its generic API must not pull in optional native engines. Explicit platform subpackages and a Client lifecycle owner make that boundary clear.

## What Changes

- Rename lifecycle owners to Client and update callers without retaining compatibility aliases for the existing Backend API.
- Construct clients from a connection string without a separate public constructor options object; continue exposing grpc.ClientConnInterface for RPC access.
- Add a pure Go connection-string API and one optional cgo integration subpackage.
- Keep core and application imports independent of cgo even when the build enables native hosting.
- Keep build-selected imports and native engine selection in the Go command, preserve exact configured subprocess arguments and external transports, and independently disable cgo in the executable.
- Exclude the direct native command C ABI and in-process execution, which remain in PR #184. Do not implement WASM.

## Capabilities

### New Capabilities

- `client-connection-platforms`: Explicit optional platform imports and library-owned connection construction.

### Modified Capabilities

None.

## Impact

Go client API, optional engine composition and build settings in projects/anki_grpc. Application integration has its own change. No external dependencies or native command implementation are added. This PR targets master independently of PR #184.
