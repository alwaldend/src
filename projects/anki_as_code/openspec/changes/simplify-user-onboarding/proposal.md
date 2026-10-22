## Why

Deferred; implementation has not started. Review finding 6 identifies a long
product README whose optional engine reference precedes the ordinary user
workflow. Readers should reach a supported first collection operation before
needing to understand the backend matrix.

## What Changes

- Put a concise default-engine quickstart near the start of the README.
- Move detailed engine and advanced reference material into owner-local
  documentation, preserving links and existing behavioral caveats.
- Link the generated CLI reference and shared repository setup as the
  authoritative sources for their respective subjects.
- Accept through a read-through and execution of the documented workflow on
  a disposable collection, plus a rendered documentation/link check.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. This is documentation reorganization with `skip_specs: true`; command
semantics and supported engines remain unchanged.

## Impact

Planning only is authorized. Future edits belong to the product README and
documentation package. Repository prerequisites and truststore portability
remain in the linked
[contributor setup change](../../../../../infra/src/openspec/changes/improve-contributor-onboarding/proposal.md).
No CLI implementation or dependency change is proposed.

Baseline: `e0d24cb3a`, inspected 2026-10-22. Evidence:
[README and optional-engine ordering](../../../README.md#optional-server-engines).
The 1,115-line review observation measures the candidate README, not reader
performance; acceptance should demonstrate usability instead of a line quota.
