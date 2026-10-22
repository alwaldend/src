## Context

BackendSTD already owns process lifetime and BackendCollection owns restore/open/close. The latter currently discards temporary packages. Applications perform explicit native exports. Native package export consumes the open collection, including native export failures.

## Goals / Non-Goals

Keep one archive publication owner and preserve explicit resource boundaries. Reuse existing grpc/protobuf dependencies and native APIs. No server protocol or upstream source changes.

## Decisions

- Add BackendGRPC implementing Start, Connection and Stop with a configured endpoint. It owns only the client connection, uses the same native error decoding as stdio, and never retries application RPCs or closes remote collections.
- Add WriteBack bool to BackendCollectionOptions; its zero value disables archive write-back. Application owners explicitly enable publication. Resolve and retain the source path during Start. Stop uses native ExportCollectionPackage with media instead of CloseCollection for an opened writable package.
- Do not use RPC success/failure or change classification to decide persistence. Existing uncertainty tracking still governs whether native work can safely complete; failures are reported rather than deleting active or recoverable files.
- Retain restores when publication fails. If an explicit client export already consumed package storage, reopen its restored database for source publication; shutdown must not close consumed storage twice.
- Factor native error decoding into a shared helper and remove runtime server data from the reusable library. Tests and application bundle declare the executable explicitly.

- Optimized worker builds strip debugging information and static symbols from the final executable; development builds retain them. Validate the embedded CLI and native backend against the optimized worker as well as the normal delivery configuration.

## Risks / Trade-offs

- A failed export can close storage: retain restore files and report their path.
- Read commands can rewrite archives: this is intentional writable-source behavior.
- Remote endpoint ownership: readiness is bounded by startup context; Stop closes clients only.

## Delivery

Extract this backend and the complete third-party Go bindings into an independent PR targeting master. CLI descriptor registration and its build-time generator belong to anki_as_code. After this prerequisite merges, rebase and finish separating the CLI PR; both PRs target master.

The independent prerequisite retains the backend library's existing worker runtime input for compatibility with current master consumers. Removing that input accompanies the CLI embedding change after the prerequisite merges. Native application E2E reproduces the missing-worker failure when that removal is extracted too early.
