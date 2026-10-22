## Why

The Go browser app needs a reusable JavaScript command API for Anki. Reuse the
existing Rust dispatcher and generated Go clients while removing the browser
pipe, HTTP/2 framing and JSPI readiness machinery. Working storage remains
Afero and the application continues to own IndexedDB snapshots.

## What Changes

- Share the native direct Go connection's protobuf, error and lifecycle handling
  with the browser adapter.
- Build a callable Rust WASI reactor around the existing command dispatcher and
  shared native backend ABI. Wrap owned buffers in a small JavaScript API.
- Replace browser stream hosting with command calls and explicit backend cleanup.
- Keep the browser WASI filesystem/timezone host and injected collection storage.
- Verify real commands, buffer ownership, failures and native/browser collection
  compatibility without requiring JSPI.

## Capabilities

### New Capabilities

- `browser-wasm-backend`: direct browser command access to Rust Anki with shared storage.

### Modified Capabilities

None.

## Impact

Implementation belongs to `pkg/ankigrpc/browser`, shared backend dispatch/ABI
and a shared internal direct-client adapter. Native public APIs and the standalone
WASI CLI remain supported. No new external dependency or shared memory is needed.
Application workflows and IndexedDB persistence belong to `projects/anki_as_code`.
