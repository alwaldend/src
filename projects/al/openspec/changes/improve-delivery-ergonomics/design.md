## Context

See proposal.md. Validation and continuation already share receipt-bound guards. Final publication currently discards its report when review inspection fails.

## Goals / Non-Goals

Provide one supported validation/publication invocation and recoverable partial results. Keep review completion and remote consistency requirements unchanged.

## Decisions

Proceed with `validate --publish`, composing the existing validation and continuation operations after validation releases its receipt lock. Continuation reacquires the lock and rechecks all evidence. This avoids custom agent wrappers without weakening environment identity. Excluding additional environment variables would conceal potentially meaningful input changes; retaining wrappers would preserve recurring overhead.

Construct a pending publication report after push and PR binding, before final verification. Return it with the original error on final verification failure. Preserve the existing attempted state so continuation remains read-only. Treating review failure as complete success would weaken the final guarantee.

## Risks / Trade-offs

- Verification still returns a nonzero exit on incomplete review reads; the structured result makes the completed mutation and safe next action discoverable.
- External changes between validation and continuation remain possible; existing receipt and candidate checks still reject them.

## Failure cases before implementation

- Failing, truncated, or dirty validation must never publish.
- `--no-pull-request` without `--publish` must fail before execution.
- Passing combined execution must validate once and publish the recorded head.
- Review epoch failure after push must retain PR URL and head, report incomplete verification, and permit read-only recovery without another push.
- Repeated explicit publication must remain refused.
