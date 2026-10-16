## 1. Shared packages

- [x] 1.1 Relocate the unchanged codec and update existing callers and Bazel labels.
- [x] 1.2 Write native E2E failure cases before implementation for read-only previews, detached inputs/results, moves, invalid ownership and provisional identities; observe the missing API fail.
- [x] 1.3 Add Store.Preview and Resources protobuf inputs, privatize binding/move helpers, and preserve validated persisted writes; verify the native preview cases.
- [x] 1.4 Document the store boundary and preview contract, with no exported raw-state mutation helper; inspect public APIs and caller source.

## 2. Verification

- [x] 2.1 Verify state/config/app/CLI E2E and inspect emitted preview/state/native database artifacts.
- [x] 2.2 Validate OpenSpec and select exact-candidate quality, semantic lint and documentation/workflow checks; receipts under out/store-boundary own final check/publication results.

## Resources protobuf preview input

- [x] Define native E2E cases before implementation for whole Resources replacement, nil retention, detached snapshots and invalid/nil entries; observe the old API fail.
- [x] Remove Bindings, accept Resources in Store.Preview, update native callers and keep state transformations/writes inside Store.
- [x] Update API docs, verify native workflows/representative artifacts/OpenSpec, and select exact-candidate delivery checks; receipts under out/resources-proto own final results.
