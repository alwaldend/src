## 1. Initial owning-client implementation (superseded by section 3)

- [x] 1.1 Write real-worker failure cases before implementation: unprepared access, successful native access, failed collection startup, Stop-before-Start, repeated Stop, worker shutdown and final Done evidence.
- [x] 1.2 Implement the initial owning-client composition with independent cleanup budgets. Section 3 replaces that composition with one factory-created client.

## 2. Integration and verification

- [x] 2.1 Update the linked Anki as Code ColOpen consumers to retain Client and verify cancellation, publication, source exclusion and asynchronous reservation lifetime.
- [x] 2.2 Verify owning and borrowed lifetimes with reusable-client, lifecycle and native/WASM MCP E2E; inspect repeatable native artifacts and document ownership. Exact-candidate semantic lint, repository quality and strict OpenSpec checks are delivery gates recorded in receipts.

## 3. Unified client lifecycle

- [x] 3.1 Write real-worker cases for opening and closing multiple collections on one client, retained-connection invalidation, raw access between collections, failed opening, repeated Stop and final native completion before replacing the implementation.
- [x] 3.2 Replace collection Client implementations with private collection state on the factory-created client. Preserve archive, native-call uncertainty and cleanup guarantees; keep publication settings in collection options.
- [x] 3.3 Migrate callers and E2E fixtures to the unified client, remove the application reservation-client wrapper, and verify all transport and native/WASM workflows. Exact-candidate delivery checks remain in receipts.

- [x] 3.4 Verify with a real active native call that terminal Stop remains bounded when another CloseCollection is already waiting; let cleanup joiners honor their own contexts without repeating the close attempt.

- [x] 3.5 Simplify redundant private collection-state lifecycle plumbing while preserving readiness, one close attempt, bounded concurrent cleanup, restore retention and final native completion. Verify the existing real-worker lifecycle and publication cases.
