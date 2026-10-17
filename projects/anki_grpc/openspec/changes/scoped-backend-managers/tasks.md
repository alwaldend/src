## 1. Native behavioral coverage

- [x] 1.1 Before implementation, add native manager E2E cases for reaping, independent scopes, copied options, failure/cancellation/panic cleanup, combined errors, delegated collection cleanup order and persistence; retain repeatable JSON/database/archive artifacts and observe the missing-API failure.

## 2. Scoped resource managers

- [x] 2.1 Add the Manager contract and ManagerSTD using the existing backend, with bounded independent cleanup; verify native transport scope and error cases.
- [x] 2.2 Add ManagerCollection over an injected Manager using the existing borrowed collection backend; verify ordered cleanup, database/package round trips, source preservation and temporary restore removal.

## 3. Documentation and delivery

- [x] 3.1 Document family-first naming, scoped borrowing, cancellation and layer boundaries; validate strict OpenSpec and inspect emitted native artifacts.
- [x] 3.2 Select exact-candidate quality, semantic lint, native/backend/consumer and documentation gates for a separate master-based PR; verify isolated branch scope and record the delivery entry point. Delivery receipts own publication and review evidence.

## 4. Preserve custom cancellation causes

- [x] 4.1 Extend the native cancellation cases before the fix to require both context cancellation and the custom WithCancelCause identity; reproduce the loss and preserve both through manager cancellation checks.
- [x] 4.2 Document the cancellation guarantee, inspect native artifacts and select fresh exact-candidate delivery gates; delivery receipts own publication and the reasoned review-thread resolution.

## 5. Shared backend scope

- [x] 5.1 Select the existing native manager E2E cases for invalid inputs, startup failure, cancellation, callback errors, panic unwinding, cleanup ordering and process reaping; keep the shared helper private rather than extending the public API.
- [x] 5.2 Add private withBackend with an explicit cleanup timeout and delegate both managers' lifecycle sequences to it; preserve their existing cleanup ordering and timeout policies.
- [x] 5.3 Document the shared scope contract, inspect native output artifacts and select fresh exact-candidate delivery gates; delivery receipts own validation and publication results.
