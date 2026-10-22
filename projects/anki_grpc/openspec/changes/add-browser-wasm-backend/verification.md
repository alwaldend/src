## Rebase verification

On 2026-10-24, the rebased backend passed native adapter and browser command
E2Es, application collection scopes, native application workflows and both
browser application suites. The private native server adapts to the same shared
direct client; the boundary fixture uses StartableClient. Generic collection
ownership remains removed from this library, with shared-filesystem preparation
owned by the application scope.

Read-only backend review found no actionable defects in the resolved code.
Evidence and exact-candidate delivery receipts are under ignored
`out/browser-wasm-app/rebase-review-20261024/`.

## Behavioral evidence

On 2026-10-23, the direct JavaScript command implementation passed:

- `//projects/anki_grpc/test/browser_backend:browser_backend_test`: the actual
  Rust reactor and Go generated clients in pinned Node 24.13.0, without JSPI.
- `//projects/anki_grpc/pkg/ankiclient:ankiclient_test`: native direct-client
  ownership, concurrent progress, cancellation, startup and shutdown behavior.
- `//projects/anki_as_code/test/browser:browser_test`: fifteen packaged SDK
  workflow groups, including real collection operations with JSPI disabled.
- `//projects/anki_as_code/test/browser:site_test`: fourteen rendered-site
  workflow groups, including collection editing, package download and reload.

The boundary report records eleven Go workflow cases and four host checks:
Unicode and empty payloads, native error details and unknown routes, supported
call options, buffer ownership, failed and cancelled startup, late-handle
cleanup, retained clients, repeatable shutdown and a real guest memory trap.
Sixteen 1,032,192-byte requests after warmup retain Rust memory at 16,580,608
bytes. Earlier returned bytes remain unchanged after later calls and cleanup.
The trapped instance never enters another Wasm export.

The compiled module exports `_initialize` and command/buffer functions,
has no `_start`, and imports only WASI and timezone services. Disassembly
confirms one-time standard reactor initialization and no per-command
constructor/destructor wrappers. Rust 1.98 omits the reactor CRT for cdylib;
the build selects `crt1-reactor.o` from the existing pinned Rust toolchain.

Application E2Es verify the nested Go/JavaScript/Rust/filesystem callback path
with real collections and binary media in shared Afero storage. Native Anki
reopens browser-produced packages and checks note/card identity, scheduling,
review history, the edited Unicode field and unchanged media. The SDK exercises
collections with JSPI disabled and configuration editing after failed backend
startup. IndexedDB snapshots remain owned by the application change.

## Evidence and limits

The boundary artifact is
`bazel-testlogs/projects/anki_grpc/test/browser_backend/browser_backend_test/test.outputs/browser-backend.json`.
Application and site reports, packages, native inspections and screenshots live
in their respective browser test output directories.

Task-local evidence is under ignored `out/browser-wasm-app/direct-command/`:
`integration-test.log` records the passing native adapter and initial browser
initialization failures; `reactor-integration-test.log` records all three
passing browser suites after the CRT fix. `reactor-artifact.json` and
`reactor-constructor-disassembly.txt` record the final library boundary.
Exact-candidate delivery receipts own final quality, lint, build, test and
publication results.

The browser suite runs Chromium, not a cross-browser matrix. The boundary
fixture uses Node WASI; the app tests establish the actual browser filesystem
host. Synchronous Rust execution cannot be preempted within the worker.
Persistence failure coverage injects transaction aborts rather than forcing
actual browser quota exhaustion. No relative performance claim is made.

## Session ergonomics

Independent agents owned Go adaptation, Rust exports, test fixtures and review;
Bazel invocations were serialized. Failure cases preceded production edits.
The read-only review found no additional actionable defect.

`wasi-reactor-initialization` was caught by the real-module startup assertion;
linking the toolchain's standard CRT corrected the lifetime model.
`browser-site-stale-crate-cache` recurred as a broken absolute Cargo.toml
symlink; one targeted repository refetch restored it without changing
dependencies or shared configuration. The initial guessed native test label was
corrected by reading its owning BUILD file. Bounded evidence and resume details
are retained in the task-local `session.md`. No new dependency or live
deployment was introduced.
