## 1. Specify native failure cases first

- [x] 1.1 Add native E2E before implementation for prepared-backend borrowing, direct persistence, connection gating, repeated lifecycle calls, temporary/persistent package restores, consuming export, failed/cancelled startup and retained backend availability; produce a reusable native archive and database.

## 2. Implement the collection backend

- [x] 2.1 Add the options object and Backend implementation using upstream services and run the native E2E, including cleanup and consuming-call behavior.
- [x] 2.2 Document ownership, restore paths and failure cleanup; verify public APIs and all declaration/field comments, strict OpenSpec validation and representative native artifacts.

## 3. Deliver and integrate

- [x] 3.1 Select exact-candidate quality, semantic lint, backend/transport E2E and build gates for a separate master-based PR; delivery receipts own validation, publication and feedback status.
- [x] 3.2 Update anki_as_code's existing PR work to consume the library and verify its collection/state/planning E2E without adding application logic to the library PR. Keep publication gated by the library prerequisite; delivery receipts own publication status.

## 4. Serialize shutdown with startup

- [x] 4.1 Reproduce shutdown winning before startup's lifecycle token using the real Anki backend, then reject stopped owners immediately after acquiring that token; verify no restore is allocated and the borrowed backend opens another collection. Baseline and corrected E2E logs are under out/anki-collection-backend; delivery receipts own updated exact-head gates and review-thread resolution.

## 5. Simplify lifecycle ownership

- [x] 5.1 Remove implicit startup shutdown and its independent timeout; callers explicitly stop the collection after startup failure. Keep startup/shutdown serialization and verify the race reproduction still passes.
- [x] 5.2 Limit the connection adapter to unary lifecycle tracking, delegate the unchanged transport interface, and reproduce then fix cleanup after a completed native export failure. Explain why native consuming operations require tracking.

## 6. Verify native cancellation synchronization

- [x] 6.1 Investigate Anki's dispatch and close mutex and distinguish running work from unconfirmed queued calls. Add a blocked-media E2E fixture with safe cleanup after backend reaping; extend its call accounting in task 7.

## 7. Account for all native calls

- [x] 7.1 Extend native cancellation E2E to require retention after an unconfirmed ordinary call and bounded shutdown while an ordinary call is active; reproduce the missing uncertainty tracking before implementation.
- [x] 7.2 Track active unary calls and sticky uncertainty without serializing ordinary RPCs, drain active callers before cleanup, and correct the narrower native-lock explanation in docs and review feedback.

## 8. Keep collection details behind the backend interface

- [x] 8.1 Remove RestoreDirectory() from the collection backend. Verify default restore ownership through native media paths, keep an explicit destination in options for caller-owned cleanup, and update anki_as_code to own that destination without another backend method.
