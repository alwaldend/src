## 1. CLI implementation

- [x] 1.1 Write native CLI failure/success scenarios before production changes, covering missing/conflicting sources, schema errors, installed execution, JSON/TOML, both network transports and persistence flags; emit repeatable native/JSON artifacts.
- [x] 1.2 Assemble all service/method commands from compiled Go descriptors with one executor and verify the native CLI scenarios, including tags and inherited backend routes.
- [x] 1.3 Supply explicit backend configuration from the embedded worker, remove discovery and application-owned publication, and verify existing application/store E2E against the revised persistence contract.

## 2. Delivery

- [x] 2.1 Update current docs, validate all three OpenSpec changes and inspect native CLI artifacts.
- [x] 2.2 Select exact-candidate package tests/builds, repository quality, affected semantic lint and guarded trunk-based delivery through repo-delivery; receipts own the results and publication state.

## 3. Explicit write-back policy

- [x] 3.1 Explicitly enable backend WriteBack in application scopes, preserve the CLI boolean pair and validate app/CLI persistence.
- [x] 3.2 Replace the custom two-executable bundle rule with standard Go embedding, verify standalone installed execution, temporary executable cleanup, and remote/help operation without a writable temporary directory.

## 4. Collection scope options

- [x] 4.1 Replace WithCollectionWriteBack and positional collection paths with one WithCollection method accepting CollectionOptions{Path, WriteBack}; migrate all callers and documentation, then verify existing native persistence and application workflows.
- [x] 4.2 Link all compiled schemas through one catalog dependency; verify complete native CLI discovery.
- [x] 4.3 Remove WithConnection from collection.Store and let App.APICall own remote client startup and cleanup; verify TCP/Unix remote process and collection lifetime.
- [x] 4.4 Consolidate backend executable configuration validation into one exactly-one check and verify native store and CLI workflows before resolving review feedback.

## 5. CLI-owned registration and prerequisite split

- [x] 5.1 Generate the schema registration file directly in internal/rpc through a CLI-owned Go tool and an automatic Bazel action; remove the backend schema library, checked-in imports and manual generation step, then verify complete native CLI discovery.
- [x] 5.2 After the separate backend/bindings PR merges into master, rebase and finish separating the CLI PR, validate the exact candidate and publish the CLI-only diff against master.

## 6. Final integration with current master

- [x] 6.1 Rebase onto linked-server and shared caller-owned collection APIs; preserve both accepted capabilities and remove obsolete embedding/runtime inputs.
- [x] 6.2 Define stalled-input timeout and interrupt cases before fixing request-input cancellation, retaining repeatable diagnostics and an unopened collection.
- [x] 6.3 Integrate and exercise offline publication independently of request outcomes, borrowed lifetime, network modes, complete descriptors, current docs and installed execution. Delivery receipts record final exact-candidate validation and publication.

- [x] 6.4 Preserve linked server lifetime after repeated Unix client disconnects by adapting SIGPIPE in the executable; verify existing network E2E.

- [x] 6.5 Rebase onto the merged architecture overview/dependency split, retain API relationships in the detailed diagram and document built-in external backend selection. Task delivery receipts own final render, validation and review results.

- [x] 6.6 Move the schema registration generator from cmd to project-local tools, retain its declared Bazel action and verify generated imports and native CLI discovery.
