## 1. Browser application

- [x] 1.1 Implement the Go browser entry point with shared Afero storage, direct App configuration methods and binary file operations; verify with the browser cross-build and E2E round trip.
- [x] 1.2 Implement the JavaScript module and worker with readiness, ordered requests, errors and disposal; verify startup failure, recoverable errors, isolated instances and disposal in Chromium.
- [x] 1.3 Package the browser artifact with its matching Go runtime and a local serving target; verify the E2E uses the packaged files and the archive contains stable names.

## 2. Verification and usage

- [x] 2.1 Add a repeatable browser E2E suite covering the design's acceptance cases and emitting configuration files and a structured report; run it against the built distribution.
- [x] 2.2 Document the module API, build/serve commands and storage lifetime; verify the documented usage against the browser suite.
- [x] 2.3 Run native filesystem workflows, rebuild the WASI CLI, validate OpenSpec and run repository formatting; record the verification results. Exact-candidate delivery gates and publication are tracked by the delivery receipt.

## 3. Collection integration and persistent snapshots

The following collection and persistence extension implements the user's approved
Afero working-filesystem and IndexedDB snapshot design. The completed tasks above
record the initial configuration-only distribution.

- [x] 3.1 Extend the browser E2E with collection/media save-reload, failed saves and conflicting-writer cases before implementation.
- [x] 3.2 Integrate managed pipe/WASI hosting, asynchronous ordered App collection operations and lazy server loading in the reusable browser SDK.
- [x] 3.3 Add transactional IndexedDB save/restore for complete package/config workspaces with observable failures and conflicting-writer protection.
- [x] 3.4 Verify the actual browser bundle and native Anki round trip, preserve config-only behavior and update documentation. Final delivery gates remain owned by the delivery receipt.

## 4. Review fixes and simplification

- [x] 4.1 Extend browser E2E for rejected directory renames preserving files across save/restore, and explicit revision conflicts before implementation.
- [x] 4.2 Limit workspace rename to regular files and simplify persistence to caller-supplied expected revisions without hidden worker revision state.
- [x] 4.3 Verify browser/native workflows and update the API documentation and delivery evidence.

## 5. Browser architecture diagram

- [x] 5.1 Document the current page, worker, Go/Rust modules, pipe, shared filesystem and IndexedDB connections in one Mermaid source; publish themed SVG siblings and a dark WebP through the existing rendering rules, inspect the image and link it from the README.
- [x] 5.2 Simplify the diagram to the UI, Go app, Rust server, Afero working files and IndexedDB saves while retaining the worker boundary and key connection labels; inspect the dark WebP.

## 6. Reject invalid workspace writes

- [x] 6.1 Extend the browser E2E before implementation to reject writes to directories or beneath regular files and preserve existing files through save/restore.
- [x] 6.2 Validate workspace write targets before creating parents or replacing bytes; document the regular-file contract.
- [x] 6.3 Run the packaged browser workflow and inspect its preservation artifacts before completing delivery.

## 7. Reject conflicting snapshot paths

- [x] 7.1 Extend the browser E2E before implementation with file/ancestor collisions in both entry orders and verify rejected restores preserve current files and saved bytes.
- [x] 7.2 Validate snapshot file and directory types before adopting the replacement workspace, then verify the packaged workflow and record its artifacts.

## 8. Direct browser backend integration

- [x] 8.1 Update browser acceptance cases before production changes to require real collection operations without JSPI and recoverable server startup failures.
- [x] 8.2 Package and lazily load the callable Rust backend through the existing app/worker API, retaining Afero and IndexedDB behavior.
- [x] 8.3 Verify SDK and rendered-site workflows with native package inspection and update the diagram and documentation. Exact-candidate gates and delivery remain owned by the delivery receipt.

## 9. Rebase onto application-owned collection scopes

- [x] 9.1 Adapt the native E2E fixture and browser factory to the current client interfaces before resolving production compatibility.
- [x] 9.2 Preserve the upstream collection scope owner and move the existing shared-filesystem directory operations into it, retaining native and browser storage behavior.
- [x] 9.3 Review the combined diff and verify browser, native collection and lifecycle workflows. Exact-candidate validation and publication remain owned by the delivery receipt.
