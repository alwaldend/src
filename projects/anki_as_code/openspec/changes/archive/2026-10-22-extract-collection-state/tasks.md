## State boundary

- [x] Keep typed State.Resources maps and the borrowed-connection Dump, Import, Update and Replace APIs; preserve namespace isolation and batch validation.
- [x] Keep the Go store independent of lifecycle; use the existing managed backend for application wiring.
- [x] Write focused native E2E before implementation; preserve resource and state behavior.
- [x] Keep collection/archive implementation and reconciliation outside #156; compose the existing backend library in the application.
- [x] Replace direct CLI dialing with App.WithState and managed BackendSTD/BackendCollection connection ownership.
- [x] Replace raw codec format arguments with typed Format and named TOML/JSON values; update fixed-format and CLI callers.
- [x] Verify state/config/CLI E2E, inspect emitted TOML/JSON/native database and document the reduced boundary.
- [x] Review the narrowed diff and select exact-candidate quality, lint, E2E, CLI/site and OpenSpec checks. Delivery receipts own validation/publication and review-thread resolution evidence.

## State recovery and command errors

- [x] Add native/API/CLI E2E cases before implementation for replacing undecodable or invalid saved state, preserving replacement validation and no-op behavior, and rejecting missing/unknown state subcommands. Both failures reproduced before production changes; logs are under out/anki-state-review-fixes/.
- [x] Make replacement independent of the discarded state's validity while propagating native read errors, and give the state command explicit argument/subcommand validation.
- [x] Verify native and config/CLI E2E artifacts, review the fixes, and select exact-candidate quality, lint and build checks for #156. Delivery receipts under out/anki-state-review-fixes own gate results and publication status. Keep the separate diagram previews unpublished.

## Managed backend wiring

- [x] Observe the updated E2E reject the old --target-only CLI before production changes; both native and CLI fixtures fail on unknown --input, recorded in out/anki-state-managed-backend/red.log.
- [x] Accept existing database inputs, delegate lifecycle to the managed backend and preserve errors during bounded cleanup.
- [x] Remove transport dependencies from the command package and document managed connection ownership.
- [x] Verify native/config/CLI E2E artifacts, review changed wiring and select exact-candidate quality, lint and build gates. Delivery receipts under out/anki-state-managed-backend own gate and publication evidence. Keep the separate diagram previews unpublished.

## Application state options

- [x] Keep the state namespace in documented app.StateOptions and collection input in app.Options; bind CLI flags to application-owned types.
- [x] Document option ownership, review the refactor and select existing native/CLI/config E2E plus exact-candidate quality, lint and build gates. Delivery receipts under out/anki-state-options own validation and publication evidence.

## Application lifecycle separation

- [x] Write native E2E before implementation for shared state access, unavailable/stopped applications, repeated startup, cancellation, operation errors and persistence after explicit shutdown. The old application API lacks Options and Start/Stop; the expected compile failure is recorded in out/anki-state-lifecycle/red.log.
- [x] Move input to app.Options and give App explicit Start/Stop lifecycle. Make WithState borrow the prepared backend and keep StateOptions limited to Key.
- [x] Initialize once per CLI invocation after argument/input decoding and stop the application at the command boundary, preserving operation and cleanup errors.
- [x] Document lifecycle ownership, inspect E2E artifacts, review changed behavior and select exact-candidate quality, lint and build gates. Delivery receipts under out/anki-state-lifecycle own validation/publication evidence.
