## 1. Previous browser pipe boundary

- [x] 1.1 Write the real Go/Rust pipe E2E for the recorded failure cases before replacement production code; emit a structured artifact using pinned tools.
- [x] 1.2 Add the Go pipe owner with existing `net.Pipe` and `ClientGRPC`; expose asynchronous JavaScript I/O without duplicating client lifecycle or native dispatch.
- [x] 1.3 Implement bounded `anki_io` imports and JSPI readiness suspension for the existing Rust WASI command.
- [x] 1.4 Verify actual generated RPCs, Unicode/empty/large transfers, native errors, cancellation and joined shutdown; preserve native implementation unchanged.

## 2. Integration and delivery

- [x] 2.1 Document lifecycle ownership, JSPI and host requirements, and deferred browser storage/application support; validate OpenSpec and formatting.
- [x] 2.2 Inspect the real-module report and select affected browser tests, build and lint gates for publication through the existing pull request. Delivery receipts own final gate and publication results.

## 3. Shared browser storage

The following storage extension is authorized by the user's continuation of the
Afero working-filesystem and IndexedDB snapshot design. The completed transport
tasks above record the earlier boundary; this section adds collection hosting.

- [x] 3.1 Record failure cases and extend real-module E2E for shared filesystem collection/package/media operations before implementation.
- [x] 3.2 Inject Afero into the existing collection client with native-compatible defaults and pass App storage through collection scopes.
- [x] 3.3 Implement WASI and timezone hosting over the shared Afero workspace and managed browser execution using the existing pipe.
- [x] 3.4 Verify package write-back, descriptor/lifecycle cleanup and native compatibility; update documentation and delivery evidence. Final delivery gates remain owned by the delivery receipt.

## 4. Direct JavaScript command boundary

The user approved replacing the pipe with a JavaScript API over shared Rust
dispatch. This section supersedes the earlier transport implementation.

- [x] 4.1 Define the command and lifetime contract and write real-module failure cases before replacing production transport code.
- [x] 4.2 Share the existing direct Go unary adapter between native and browser clients, preserving generated clients, error details, call options and cleanup ownership.
- [x] 4.3 Build callable Rust Wasm exports over the existing dispatcher, expose them through a JavaScript command API, and remove browser pipe/HTTP2/JSPI hosting.
- [x] 4.4 Verify direct commands, buffer ownership, startup failures, cancellation, traps and repeated cleanup; preserve shared filesystem collection/media round trips and native compatibility.
- [x] 4.5 Update current documentation and acceptance evidence and select validation through the existing PR. Exact-candidate gates and publication remain owned by the delivery receipt.

## 5. Rebase onto current transport interfaces

- [x] 5.1 Adapt the boundary fixture to StartableClient and retain upstream's private native implementation through the shared direct adapter.
- [x] 5.2 Preserve removal of generic collection ownership, review the resolved backend and validate native/browser compatibility. Exact-candidate validation and publication remain owned by the delivery receipt.
