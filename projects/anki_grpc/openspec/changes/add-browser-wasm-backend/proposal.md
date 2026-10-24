## Why

Browser Go clients need reusable Anki commands backed by Rust and the caller's
Node-compatible JavaScript filesystem, independently of any application SDK.

## What Changes

- Add the browser Go client, callable Rust reactor, JavaScript command bridge,
  and production WASI filesystem/timezone host.
- Reuse PR1's shared direct-client adapter, backend ABI and timezone helpers.
- Join factory completion in Start, checking context before and after startup;
  retain command cancellation, concurrent Stop deadlines and Done cleanup.
- Retain the Node command-boundary E2E and add independent Chromium coverage
  through production `wasi.js`, generated Go clients and shared provider storage.
- Supply a provider-independent Go filesystem initializer; Rust filesystem calls
  stay in JavaScript, and Go uses NewOsFs/BasePathFs over the same provider.
- Supply the approved Emscripten classic MEMFS provider with a pinned official
  SDK and hermetic build rule, then exercise it through an explicit Chromium
  fixture; ZenFS and the npm `memfs` package are not selected.

## Capabilities

### New Capabilities

- `browser-wasm-backend`: direct browser commands and shared working storage.

### Modified Capabilities

None.

## Impact

PR2 targets master with the shared core merged through PR195. It owns `pkg/ankigrpc/browser`,
`test/browser_backend`, the browser-test schema visibility grant and the existing
Go SDK host-label export. It also preserves the original standalone WASI guest
visibility grant to the existing Anki client consumers. PR1 owns shared native
adaptation, `internal/directclient`,
`internal/backend/ffi.rs`, both shared BUILD files (including browser grants),
and timezone extraction. Application collection scopes,
SDK packaging, IndexedDB and site code belong to subsequent PRs. The current
provider-independent host retains its Node-compatible contract. The user-approved
extension adds Emscripten and its required official SDK components, a generated
standalone MEMFS module and the separately specified Node-compatible adapter.
It adds no deployment or IDBFS; Rust continues to use the existing WASI target.
Chromium and composed SDK acceptance precede publication.
