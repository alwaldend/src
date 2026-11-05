## 1. Failure evidence

- [x] 1.1 Write lagged-head, partial-push, advanced/dirty local-state, and refusal cases before implementation; capture baseline failure and repeatable JSON artifacts.

## 2. Guarded recovery

- [x] 2.1 Preserve observed update expectations and partial-push reports; verify the existing and new provider/publication cases.
- [x] 2.2 Add receipt-bound PR-only recovery with repository, historical candidate, identity, metadata, and remote-head guards; verify local state preservation and idempotence.
- [x] 2.3 Update delivery reference and canonical skill with exact recovery steps; validate CLI help, structured output, and skill packaging.

## 3. Independent delivery

- [x] 3.1 Run owner tests, formatting, strict OpenSpec, root quality, and semantic lint; review and prepare the independent trunk-based repair candidate. Delivery receipts own final checks and publication state.

All 31 selected tests and affected semantic lint passed before preparation.
The isolated JSON artifact records unchanged local work, receipt, and remote
ref, with one metadata mutation across recovery and its idempotent retry.
