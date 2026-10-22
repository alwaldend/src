## 1. Backend implementation

- [x] 1.1 Specify and write native behavior coverage before production changes for default/disabled write-back, prior RPC failure, consuming calls, publication failure retention and remote connection ownership.
- [x] 1.2 Implement collection-owned write-back and explicit gRPC connections, then verify backend and CLI E2E preserve native data and errors.
- [x] 1.3 Remove implicit server runtime data when the paired CLI embedding change lands; retain the existing input in the independent backend prerequisite so current master application consumers keep working.

## 2. Validation

- [x] 2.1 Validate the owning OpenSpec change and select affected native tests, build and lint for exact-candidate delivery; receipts own passing results.

## 3. Explicit write-back policy

- [x] 3.1 Replace the negative constructor option with opt-in WriteBack and verify omitted, false and true policies through native E2E.
- [x] 3.2 Build optimized workers with final executable symbol stripping and verify the native backend and embedded CLI in that configuration.

## 4. Independent backend delivery

- [x] 4.1 Extract backend and complete third-party bindings into a master-based candidate without CLI descriptor registration or its generator; verify native behavior and select required independent delivery gates. Delivery receipts own publication and exact-candidate results.

- [x] 4.2 Document the retained transitive worker runtime input in the current backend README and select exact-candidate checks before resolving the extracted PR review.

## Archive reconciliation

Task 1.3 remained unchecked after `expose-linked-server` completed its task 1.5.
Source inspection on 2026-10-22 confirms that
`pkg/ankiclient/BUILD.bazel` declares worker runtime data only for its test,
not its reusable library, and the current README documents caller-supplied
executables. The checkbox now records that completed integration; no product
code changed and no runtime test was repeated for this task-state correction.
