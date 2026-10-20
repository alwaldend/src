## Why

Callers need to execute Anki unary operations inside the Go process while retaining generated gRPC clients, native errors and explicit cleanup.

## What Changes

- Share Rust dispatch between the existing gRPC server and an owned command C ABI.
- Add a native Client in the opt-in ankigrpc/native subpackage, exposing the existing grpc.ClientConnInterface.
- Select fresh native state with inprocess:// labels or inprocess+cgo://. Labels do not imply sharing.
- Provide the minimal address-only constructor and application hook required to select this feature without linking native code into generic libraries.
- Build on the Client lifecycle naming, build switches and executable composition merged through PR #186. Preserve pure-library imports and optional native linking. WASM remains deferred.
- Configure private factory registrations and explicit defaults without wrapper functions. Keep the minimal lifecycle contract, expose asynchronous cleanup separately and validate network endpoints once during construction.

## Capabilities

### New Capabilities

- `in-process-backend`: Direct native execution, lifetime, errors and optional address selection.

### Modified Capabilities

None.

## Impact

Native Rust/C ABI, optional Go adapter and feature E2E in projects/anki_grpc. Application integration has its own owner-local change. This PR targets master after the merge of #186 and introduces no external dependencies.
