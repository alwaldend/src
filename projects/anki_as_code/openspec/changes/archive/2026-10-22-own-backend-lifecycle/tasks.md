## 1. Scope implementation

- [x] 1.1 Replace manager-based lifecycle fixtures before implementation with native process observations covering success, callback error, cancellation, panic, startup failure and joined crash errors; emit a native collection artifact.
- [x] 1.2 Construct both backends directly and retain one private bounded scope helper; verify the native lifecycle and source/publication E2E.
- [x] 1.3 Update current collection/API documentation and inspect the full diff for stale manager references.

## 2. Verification and publication

- [x] 2.1 Validate the exact candidate with native store/application/backend E2E, repository quality, affected semantic lint and CLI/documentation builds; inspect emitted artifacts.
- [x] 2.2 Review and publish the independent master-based PR with both owning OpenSpec records and verified delivery receipts.
