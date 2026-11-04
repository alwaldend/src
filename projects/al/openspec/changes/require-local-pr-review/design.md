## Context

See [proposal.md](proposal.md) for motivation. The canonical delivery skill owns
operator procedure; AL receipts guard publication and remote mutations, but do
not certify human or model review quality.

## Decisions

- Extend the existing prepublication correctness step with explicit local
  `review-agent` invocation and complete candidate scope. A separate CLI review
  ledger would duplicate caller-owned readiness decisions and is unnecessary.
- Repeat local review when the reviewed diff changes, including publish-time
  rebases. An unchanged candidate can reuse its review evidence during recovery.
- Keep remote review handling in the existing guarded workflow and its
  capability-limited fallback. Tests and remote comments cannot replace local
  review.
- Add scenarios to the existing offline evaluation suite before changing the
  procedure. Regenerate discovery through its owning target.

## Risks / Trade-offs

- Instructions cannot prove review execution or quality. Offline validation
  checks configuration only; inspect this task's complete diff locally and
  report actual publication and remote review state separately.
- Reviewing inert records at the same depth as behavioral changes adds noise.
  Include every changed file while keeping scrutiny proportional to its content.
