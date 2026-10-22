## Behavioral evidence

On 2026-10-23,
`//projects/anki_as_code/test/browser:site_test` passed all fourteen workflow groups
in pinned Chromium 150.0.7871.24. The retained run is
`out/browser-wasm-app/direct-command/reactor-integration-test.log`. This run
uses the callable Rust backend through its JavaScript command API and the
shared Afero host.

- Landing navigation, the automatically opened example, exact edit retention
  across files, JSON download/reopening and malformed-input recovery retain the
  existing configuration workflow and shared site presentation.
- Save locally includes current editor text and restores it on page reload.
  A stale tab cannot replace a newer saved revision; its edits remain available,
  including after reopening JSON, and Restore saved recovers the committed snapshot.
- With IndexedDB unavailable, startup recovers to the example, JSON and native
  collection imports succeed, Plan converges and both downloads remain usable.
  Save reports the storage failure. The completed package is retained as
  `unavailable-storage-collection.colpkg`.
- Opening a native collection exports editable configurations with installed
  bindings. The unchanged import converges; editing a Unicode note produces a
  plan, Apply updates the collection and the next plan converges again.
- A failed Apply retains the prior package byte-for-byte and the invalid editor
  text. Invalid collection imports and JSON replacement preserve the selected
  package. Saving unapplied edits preserves those config bytes and the last
  completed package without implicitly applying; both restore after reload.
- JSON paths targeting `/collection.colpkg`, descendants of that package path,
  reserved `/tmp` storage, or another file outside `/workspace/` are rejected before adoption. The page
  retains exact editor text and collection bytes through local save and reload.
  A configuration-only session also rejects a JSON path that would create a
  collection or hide notes in temporary storage. `site-report.json` records the
  four escaped-path errors.
- The downloaded `site-browser-output.colpkg` is independently reopened with
  native Anki. Native baseline/verification JSON confirms retained identities,
  scheduling, review history and binary media with the applied
  `browser edited 你好` field. The unchanged media SHA-256 is
  `7e3691790cd64b19d4edb1a80e988214515abeb53aa0f34ffbfe4b4bf405d120`.
- Light/dark theme switching and preference persistence pass. At a 390-pixel
  viewport the page width is 390 pixels. Desktop light/dark and mobile captures
  retain the existing site shell and added collection/persistence controls.

`site-report.json` records no uncaught errors, failed requests or undeclared
external requests and includes worker lifecycle and request diagnostics. The
local test server enforces the network allowlist with Content Security Policy;
passive request auditing verifies the observed requests. The only permitted
external scripts are the site's existing SRI-pinned scripts.

The suite writes `site-report.json`, `site-desktop-light.png`,
`site-desktop-dark.png`, `site-mobile-light.png`, configuration downloads, the
completed collection package and native verification artifacts under
`bazel-testlogs/projects/anki_as_code/test/browser/site_test/test.outputs/`.
An earlier live T3 preview of the pipe candidate completed Save locally and
page reload with restored edits. The current command API is verified by the
automated workflow above. No live deployment occurred.

Exact-candidate build, test, lint, repository-quality and publication results
are owned by the delivery receipts under ignored
`out/browser-wasm-app/direct-command/`; individual behavioral results above do not
replace those gates.

## Diagnosed build issues

- `browser-site-fetch-pause-on-reload`: the previous Puppeteer Fetch request
  interception left worker requests paused during the longer workflow.
  A controlled comparison of the full workflow sequence isolated interception
  as the cause, while live T3 reload also succeeded. Replacing interception with
  a test-server CSP allowlist and passive request auditing made the full twelve
  groups pass. Worker/request/error diagnostics remain in the test. No production
  workaround was added.
- `browser-site-grpc-docs-target`: the aggregate documentation graph required a missing `docs` target for the existing Go gRPC README. Added that mapping without modifying the dependency.
- `browser-site-stale-crate-cache`: the cached Rust `percent-encoding-iri` repository contained broken absolute symlinks into a deleted output base. A targeted Bazel forced refetch with repository-contents caching disabled restored readable inputs. Normal builds then succeeded; no shared cache deletion or dependency/configuration change was needed.
- `browser-site-resource-mapping`: the site's explicit webfont static mount suppresses the implicit static mount. The first rendered-site run found the missing controller. Packaging the app as a Hugo leaf bundle with `preserve_paths = True` publishes the runtime through the existing content mount. The complete site workflow passed after this correction.
- `browser-site-roundtrip-source-comparison`: the first simplification E2E compared regenerated TOML byte-for-byte. Reopening JSON correctly makes inherited deck references explicit. The corrected assertion compares the complete configuration after reopening; switching files still requires exact editor-text preservation. The initial failure report is retained in task scratch.

## Session ergonomics

The collection/storage extension reused this change, the SDK and the existing
rendered-site E2E. Parent-owned Bazel execution remained serialized while agents
owned disjoint source and test files. The reload failure was investigated with
retained request diagnostics and a controlled full-sequence comparison rather
than changing working production startup. No new dependency, shared
configuration or live deployment was needed.
