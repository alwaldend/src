## 1. Extraction

- [x] 1.1 Adapt the existing lifecycle/state E2E cases before production changes and inspect their native fixture helpers to exclude reconciliation/export behavior.
- [x] 1.2 Extract the collection workflow and state adapters; verify native lifecycle/state, CLI and configuration E2E targets pass while feature commands remain placeholders.

## 2. Documentation and delivery

- [x] 2.1 Update current lifecycle/state documentation and its diagram, validate OpenSpec strictly, and inspect rendered output.
- [x] 2.2 Inspect the independent master-based diff and select exact-candidate quality, affected-target lint, E2E and integrated site checks. Delivery receipts own final validation and publication results.

- [x] 2.3 Put lifecycle ownership in WithCollection, construct state stores directly inside its callback, and synchronize the native E2E, callers and documentation. Delivery receipts own the updated candidate checks.
