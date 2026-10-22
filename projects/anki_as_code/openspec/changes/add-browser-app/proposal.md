## Why

The application supports injected Afero storage, so browser callers can reuse
its configuration and collection workflows. A browser-specific Go build and
reusable JavaScript API retain an editable workspace between calls. A callable
Rust WASI reactor exposes Anki's existing dispatcher through a JavaScript command
API for offline collection operations. Explicit local snapshots preserve
completed packages and configurations across reloads.

## What Changes

- Add a pure Go `js/wasm` entry point and a static distribution containing the matching Go runtime.
- Expose promise-based configuration and collection operations and binary
  workspace file access through a dedicated worker, calling the existing App.
- Load the callable Rust backend lazily, invoke its JavaScript command API and
  share one Afero filesystem with its WASI host, without requiring JSPI.
- Add atomic IndexedDB snapshots of completed collection packages and editable
  configurations, with explicit save/restore and conflicting-writer protection.
- Document browser loading, storage lifetime and supported operations; verify
  browser collection save/reload and native Anki compatibility.

## Capabilities

### New Capabilities

- `browser-app`: Browser loading, direct App methods, shared Afero workspace,
  worker lifecycle and explicit local snapshots.

### Modified Capabilities

None.

## Impact

Behavior changes are owned by `projects/anki_as_code`. One import in
`tools/go/include.MODULE.bazel` exposes the existing SDK label repository so the
bundle can consume its matching runtime. This adds no SDK or dependency version.
The implementation reuses the existing App, Afero dependency, Go SDK, package
rules, browser test tools and `anki_grpc` direct browser client/WASI host. IndexedDB uses
browser APIs without a new external dependency. Native and WASI CLI entry points
retain their current contracts; site presentation has its own change record.
