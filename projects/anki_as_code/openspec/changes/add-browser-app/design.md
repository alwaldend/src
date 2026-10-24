## Context

See proposal.md for motivation. The merged App constructor accepts `Options.FS`; all configuration I/O already uses that storage. The existing WASI CLI has a command lifecycle, whereas a browser caller needs a long-lived object. The repository already pins Go, Afero, rules_pkg, Node, Puppeteer and Chromium.

## Goals / Non-Goals

**Goals:** Reuse App configuration and collection methods, codecs and lifecycle;
retain one shared Afero workspace per instance; make browser startup and teardown
observable; explicitly save and restore completed work across reloads.

**Non-Goals:** A second command schema, shared memory, live SQLite persistence,
multi-user editing, automatic background saves, new dependencies or an end-user
editor UI. The existing site change owns the visitor-facing controls.

## Decisions

- Add `cmd/anki_as_code_browser`, built with the Go `js_wasm` toolchain. Keep native and WASI entry points separate. Package the SDK's own `wasm_exec.js`; do not copy or patch its implementation.
- Put the reusable JavaScript module and worker in `pkg/browser`. `createAnkiApp({wasmURL?, workerURL?})` returns `{app, workspace, dispose}`. App offers `configExport({output, configs})` and `configDump({text, format = "json"})`. Workspace offers `writeFile(path, Uint8Array)`, `readFile(path)`, `listFiles(path)` and `remove(path)`; listing returns sorted absolute file paths recursively, writing creates parents, and removal recursively removes the named path.
- Each worker constructs one `afero.NewMemMapFs()` and one App supplied with that
  storage and a managed browser client factory. Workspace, App, collection
  manager and Rust WASI host share its virtual paths. Boundary paths must be
  absolute, nonempty virtual paths; no host/home expansion or extra subtree
  confinement is promised. Transfer uses byte copies.
- The Go entry point exposes a private asynchronous dispatch callback. It
  returns immediately and runs App work in a goroutine so asynchronous module
  loading and backend creation can progress through JavaScript events. Rust
  commands then run synchronously inside the worker. The worker awaits each
  complete operation before starting the next, including workspace access,
  saves and restores.
  Existing App methods, protobuf codecs and collection lifecycle own business
  behavior; the browser boundary adds no CLI parser or reconciliation logic.
  A separate ready message signals initialization; `Go.run` completion means
  termination.
- Load and compile the callable Rust WASI reactor only when a collection
  operation first needs it. The `anki_grpc` browser client invokes its JavaScript
  command API over the existing Anki dispatcher, using the shared direct-client
  adapter for protobuf calls, errors and lifetime. All workflows run without
  JSPI; configuration-only startup, editing and JSON operations do not fetch
  the backend. Backend creation remains asynchronous. Collection startup errors
  are observable and leave configuration work available. Preserve collection
  closure and joined shutdown before the next workspace operation.
- Cancellation can stop waiting for asynchronous initialization and joins any
  late-created backend during cleanup. It cannot preempt a synchronous Rust
  command executing in the worker. Ordinary backend errors remain recoverable;
  trapped instances are discarded after cleanup. Explicit application disposal
  terminates the worker and rejects pending operations.
- Expose the existing App's collection import/export, planning and application
  behavior through the promise-based module boundary. Pass binary packages
  through workspace byte operations; Anki owns package decoding, SQLite and media
  semantics. Preserve native application error and write-back behavior.
- Persist explicit snapshots through IndexedDB. Copy the completed `.colpkg`
  and editable configuration files from a quiescent workspace after collection
  work and package write-back finish. Configuration-only workspaces are also
  saveable and carry no collection package. Commit the optional package, config bytes and
  snapshot revision in one transaction; do not persist temporary collection
  databases or backend staging. Reserve `/tmp` for runtime temporary files and
  exclude it from saved/restored workspace data. Preserve exact editable bytes
  and path identity.
- A save succeeds only after transaction completion. Quota errors, aborts and
  serialization failures leave the previous saved snapshot intact and retain
  current in-memory work. Check the expected saved revision inside the same
  read-write transaction, or provide equivalent exclusive writer coordination,
  so a stale tab cannot silently overwrite newer saved work. The caller supplies
  the expected revision from its last save or restore; zero creates a new key.
  The worker keeps no duplicate revision state, allowing an application to replace
  its in-memory workspace while retaining its last observed persistence revision.
- Restore a complete saved snapshot into a fresh workspace and replace the
  active workspace only after restoration succeeds. Restore neither merges
  stale files nor destroys unsaved work on failure. Track the restored revision
  in the caller for later conflict checks. Saving and restoring are explicit caller actions;
  IndexedDB availability does not affect ordinary configuration downloads.
  Use the shared filesystem resolver to reject regular-file ancestors and
  inspect each destination's type before creation. Explicit directory entries
  may follow their children, but a file cannot replace an implicit directory.
- Preserve existing export behavior: schema errors are detected before writes, export may partially write on filesystem error, and semantic validation belongs to the subsequent load. Use protobuf JSON conventions, including string-encoded int64 fields.
- Publish stable bundle names with existing package rules and an existing static server target. Use `test/browser` for the documented, repository-pinned Puppeteer E2E workflow; no new dependency is needed.

## Risks / Trade-offs

- Worker/runtime startup or termination can strand promises → handle worker errors, message errors, runtime exit and disposal explicitly; bound startup with a timeout.
- Long synchronous Rust commands delay worker events and cancellation until they
  return; keep the backend in the dedicated worker so the main page stays responsive.
- JavaScript type mistakes can panic Go bindings → validate argument types and byte arrays before conversions, and report ordinary request errors through the worker protocol.
- Unsaved memory storage is lost on disposal/reload → distinguish active edits,
  completed local saves and portable downloads. Browser storage can be cleared
  or evicted; explicit downloads remain available.
- Interrupted saves or stale writers could corrupt a saved workspace → commit
  all snapshot components atomically, report transaction completion and reject
  conflicts while preserving both the prior snapshot and current edits.
- An active SQLite database is not a safe snapshot boundary → wait for the
  existing collection lifecycle and package publication to finish before save.
- Worker copies incur transfer cost → use byte arrays, avoid base64, and test arbitrary binary bytes. Preserve caller ownership of input buffers.
- Browser platform code can break native wildcard builds → mark platform-only Go targets appropriately while exposing a dedicated cross-build target.
- Afero's memory backend permits byte writes to directory handles and can convert
  a regular file when creating a directory over it. Inspect a workspace write's
  resolved destination before any parent creation or byte replacement. Existing
  destinations must be regular files; the shared filesystem resolver rejects
  regular-file ancestors. Do not duplicate its path traversal in the browser.

## Acceptance Cases Before Implementation

Exercise the actual static bundle in Chromium: Unicode/multiline export and both dump formats; exported TOML retrieval; imported bytes and fresh reads after edits; ordered operations; independent workers; binary bytes and unchanged input ownership; malformed schema, invalid paths, missing files and recovery; startup fetch/compile errors; pending/future rejection on disposal. Save the TOML tree, JSON dump and structured browser report as test artifacts. Run existing native memory-filesystem workflows and rebuild the existing WASI entry point for compatibility.

Before production storage integration, extend this E2E to import a real
collection with binary media, export configs, plan and apply an edit, export a
completed package, save locally, reload, restore and verify the edit and media.
Reopen the exported package with native Anki to check content, scheduling and
review history. Exercise a failed/aborted save, quota failure, stale-writer
conflict, failed restore, recoverable backend startup and operation errors, and
resource cleanup. Verify config-only usage never fetches the guest. Emit the
completed package, restored config tree and a structured persistence/lifecycle
report as repeatable artifacts.
Run collection workflows without JSPI. Also save and restore a configuration-only
workspace without fetching the Rust guest; verify runtime `/tmp` files never
enter its snapshot. Verify asynchronous backend creation and its late cleanup
after cancellation, observable trapped-command failure, and recoverable native
errors through the direct adapter's owning E2E.
Before changing workspace writes, add E2E cases for empty and populated directory
targets, the root directory, and immediate or deeper regular-file ancestors.
Verify failures preserve the entire file tree and bytes before and after a fresh
worker restores the saved snapshot; verify ordinary create/replace still works.
Before changing snapshot restoration, seed conflicting file/ancestor records in
both entry orders, including explicit child directories and root-as-file. Verify
rejection preserves unsaved work and stored records, then save and restore the
retained workspace in a fresh worker. Compatible child-before-parent directory
entries must still restore successfully.

## Migration Plan

This is an additive build and module API. Consumers opt in by serving the bundle and importing its module. Removing those imports rolls back usage; native consumers require no migration.

## Review simplification decision

Proceed with explicit `session.save(key, expectedRevision = 0)` and regular-file
`workspace.rename`. The page already owns its last observed revision; passing it
to the existing transactional comparison removes the worker's second owner and
the import-time restore. Keeping implicit revisions would require extra state
transfer or restoring unrelated saved files into every replacement worker.
Only Apply consumes rename, to replace a completed package. Reject directories
instead of duplicating the backend's directory-rename validation. Browser E2E
must preserve stale-writer protection and prove rejected renames retain files
through persistence. This refines the new, unmerged SDK; its declarations,
documentation and callers change together.
