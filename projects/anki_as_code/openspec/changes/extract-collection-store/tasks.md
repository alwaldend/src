## 1. Native behavior cases

- [x] 1.1 First define native E2E cases for reused connections/services, namespace isolation, config export, aliases, retained clients after cleanup, cancellation and failed publication; emit a native collection artifact.

## 2. Collection ownership

- [x] 2.1 Extract Store and the injected Collection into internal/collection with documented package, declarations and fields; preserve native backend ownership and successful scope publication.
- [x] 2.2 Expose prepared resource/state/export services and native explicit export. Delegate the existing master app workflow to Store while retaining placeholder commands and CLI compatibility.
- [x] 2.3 Simplify callback injection per user feedback: remove Use, Ref, path selectors and ID lookup. Keep the prerequisite independent of unmerged CLI and reconciliation implementations.

## 3. Verification and delivery

- [x] 3.1 Document collection.Store and its boundaries, validate OpenSpec, and verify native/API behavior and representative emitted artifacts.
- [x] 3.2 Select exact-candidate quality, semantic lint, build and E2E gates for the new master-based PR. Delivery receipts own publication and review evidence.

## 4. Retain exporters within the collection scope

- [x] 4.1 Extend native service reuse coverage before implementation: repeated namespace lookups must return the same exporter, empty/default keys must share it, and different namespaces must remain separate. Reuse that exporter after native state updates and preserve the round-trip artifacts. Both database and package cases reproduced exporter reconstruction before the cache was added.
- [x] 4.2 Cache exporters by normalized namespace alongside the existing state stores, document the contract, and select exact-candidate validation gates. Delivery receipts own publication evidence.

## 5. Publish only the source archive

- [x] 5.1 Adjust native E2E before implementation to explicitly export caller-selected destinations through Anki's API, retaining source-byte and native state round-trip checks. Keep successful source publication and failed/cancelled scope cleanup coverage.
- [x] 5.2 Remove SetOutput, output reservations and destination selection from collection.Store. Preserve automatic publication to the store's own source archive; callers export other destinations directly as their final operation and return false to skip source publication.
- [x] 5.3 Update lifecycle documentation and OpenSpec contracts and select aggregate validation gates. Delivery receipts own exact-head verification and publication.

## 6. Application workflow boundary

- [x] 6.1 Define direct native App state-workflow E2E cases before implementation, preserving database/package publication, namespace isolation, validation failures and repeatable artifacts; adapt lifecycle coverage to the collection API.
- [x] 6.2 Remove App.WithCollection and application flag options; expose typed ColStateDump, ColStateImport, ColStateUpdate and ColStateReplace workflows. Keep CLI flags, serialization and output at the command boundary.

## 7. Use backend managers

- [x] 7.1 Add native coverage before implementation for an injected transport manager, shared native storage, error/cancellation/panic cleanup and reservation release.
- [x] 7.2 Delegate collection scopes to backend.ManagerCollection over an optional injected backend.Manager. Remove manual backend Start/Stop from Collection; retain zero-value Store defaults, reservations, service reuse and archive publication.
- [x] 7.3 Update documentation and specifications, inspect native output artifacts and select exact-candidate publication gates; delivery receipts own execution and publication results.

## 8. Direct prepared connections

- [x] 8.1 Adapt native workflow cases before production changes to reuse locally constructed services on the injected connection, including namespace isolation, expired connections, source aliases, publication and manager cleanup.
- [x] 8.2 Remove Collection and its service caches. Keep a private source-reservation map in Store, delegate native lifetime to ManagerCollection and pass grpc.ClientConnInterface to callbacks. Adapt application callers and carry the same boundary to #164 and #128.
- [x] 8.3 Update current documentation and specifications and select exact-candidate native workflow, artifact, quality, lint and build gates. Delivery receipts own check execution and publication evidence.
