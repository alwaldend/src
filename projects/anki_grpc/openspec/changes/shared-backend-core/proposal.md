## Why

Extract the shared backend core from PR #191 so native and WASI foundations
can be reviewed and merged independently before the browser follow-ups.
Native startup keeps initialization and failure cleanup in the caller while
retaining explicit ownership under concurrent shutdown.

## What Changes

- Extract direct unary calls, error decoding and lifecycle into internal/directclient.
- Keep native Start synchronous, checking context before and after initialization;
  concurrent Stop and Done join startup and cleanup.
- Move the Rust command C ABI to internal/backend and reuse it from native hosting.
- Move timezone calculations to internal/timezone, preserving native WASI behavior.
- Retain native/WASI documentation and the shared core's narrow consumer grants.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `in-process-backend`: Startup waits for initialization and failure cleanup;
  concurrent Stop retains its bounded wait and completion observation.

## Impact

Only projects/anki_grpc shared implementation, native adaptation and WASI
timezone delegation change. Existing generated-client, worker and application
consumers validate compatibility. No dependency, deployment or browser package
is introduced. Browser follow-ups depend on this independently based change.
