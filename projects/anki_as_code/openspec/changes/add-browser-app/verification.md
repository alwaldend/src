## Rebase verification

On 2026-10-24, compatibility with the current application-owned collection scope
passed the packaged SDK and rendered-site E2Es (fifteen and fourteen workflow
groups), the native application suite, collection scope tests, native adapter
tests and browser command-boundary E2E. The browser reports contain no runtime
errors. Native package inspection still verifies the edited field, identities,
scheduling, review history and binary media.

The rebase preserves the new scope owner and private native implementation.
Temporary restore, media and cleanup operations use its injected filesystem;
the browser factory and native fixtures use the current client interfaces.
The UI and SDK files are unchanged by conflict resolution.

Compatibility logs, the three-way resolution record and review evidence are
under ignored `out/browser-wasm-app/rebase-review-20261024/`.
Exact-candidate delivery receipts there own final gates and publication.

## Behavioral evidence

On 2026-10-23,
`//projects/anki_as_code/test/browser:browser_test` passed fifteen workflow
groups against the packaged Go app, callable Rust WASI backend and browser host
in pinned Chromium 150.0.7871.24.

- Configuration export/dump, Unicode and multiline text, fresh reads, binary
  byte ownership, ordered requests and independent workers.
- Rejected writes to directories or beneath regular files preserve the complete
  tree through save/restore; ordinary create and replacement still work.
- Regular-file publication and rejected renames preserve bytes. Conflicting
  snapshot paths reject in either entry order without replacing unsaved work
  or modifying saved records; compatible child-before-parent entries restore.
- Native collection import, configuration export/edit, plan/apply, completed
  package write-back, ordered save and restoration into a fresh worker.
- IndexedDB revision conflicts and an injected transaction abort preserve the
  previous complete snapshot.
- Real collection export succeeds with JSPI disabled. Invalid server WebAssembly
  produces an observable startup error while configuration editing stays usable.
- Input errors, startup fetch/compile failures and idempotent disposal preserve
  the documented worker and promise behavior.

Native Anki independently reopened `browser-output.colpkg`. Its baseline and
verification JSON confirm retained note/card identity, scheduling and review
history, the applied `browser edited 你好` field, and unchanged binary media
SHA-256 `7e3691790cd64b19d4edb1a80e988214515abeb53aa0f34ffbfe4b4bf405d120`.
The report records saved revision 1, no browser errors and no console messages.

The rendered-site suite also passed fourteen workflow groups with the same
backend, including editor and collection saves across page reload, stale-save
protection, escaped-path rejection, operation without IndexedDB, native package
inspection, and light/dark/mobile presentation. Site evidence remains owned by
`integrate-browser-app-site/verification.md`.

The direct boundary E2E and native adapter tests pass; their evidence is owned
by `projects/anki_grpc/openspec/changes/add-browser-wasm-backend/verification.md`.

## Artifacts and limits

The SDK suite writes `browser-report.json`, `browser-output.colpkg`,
native baseline/verification JSON, configuration dumps and restored TOML under
`bazel-testlogs/projects/anki_as_code/test/browser/browser_test/test.outputs/`.
The report records the exact archive SHA-256. Passing browser run evidence is
`out/browser-wasm-app/direct-command/reactor-integration-test.log`.

The browser matrix currently covers Chromium. The suite injects IndexedDB
transaction aborts rather than actual quota exhaustion. Runtime exit handling
exists, but a post-readiness Go runtime crash is not injected. Synchronous Rust
commands block their worker until completion; cancellation cannot preempt them.

## Diagram and delivery

`docs/browser-app.mmd` remains the sole diagram source. The light/dark SVG
siblings and dark WebP build successfully through `:browser_app` and
`:browser_app_preview`. The inspected WebP retains the shared palette, readable
labels, contained titles and unobstructed connections. It shows the UI, Go app,
Rust backend, Afero working files and IndexedDB saves, with the JavaScript
command API replacing the pipe connection.

Exact-candidate receipts under ignored `out/browser-wasm-app/direct-command/`
own final build, test, semantic lint, repository-quality and publication
results. Earlier pipe implementation receipts are historical. No new dependency,
host configuration change or live deployment was needed.
