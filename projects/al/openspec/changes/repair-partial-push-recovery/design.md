## Context

The original preparation receipt already binds the candidate, Git endpoints,
remote lease, and exact prior PR projection. The new recovery operation uses
that owner rather than accepting a new loosely captured expectation.

## Goals / Non-Goals

Recover an already-pushed candidate's PR while preserving later local work.
Do not authorize another push, a rebase, a replacement PR, human-edit overwrite,
or final validation of the currently checked-out source.

## Decisions

Preserve the real pre-push expectation; the GitHub update guard accepts only
its exact old-to-pushed head transition while retaining every other field check.
Return the observed push report before any subsequent fallible PR operation.

Use explicit `recover-pull-request --receipt-file` because ordinary publish
must still require a clean exact candidate. A metadata-only operation can
validate the historical commit and receipt bindings without requiring local
HEAD or working files to match that historical candidate. Require the current
branch and configured base to match the receipt refs and freshly fetch the exact
remote head before metadata mutation. Verify the original PR identity and
allow only the prior or desired exact projection. Recheck the remote immediately
before mutation; preserve the receipt unchanged. Do not reinterpret generic
inspection refusals as broad permission to publish.

Alternatives rejected: indefinite waiting, another Git push, deleting or
recreating the receipt, and blanket removal of head/metadata guards. An automatic
historical recovery inside ordinary publish would obscure which local candidate
is being validated; the explicit command makes its restricted scope reviewable.

## Risks / Trade-offs

- GitHub has no metadata compare-and-swap: retain immediate guarded reads and
  post-mutation verification; concurrent changes can still cause pending status.
- A final provider read can fail after mutation: retain known identities and
  document inspection before a retry.
- The extra command is justified by recovery after local work advances; existing
  publish and verify candidate gates remain unchanged.

## Failure cases written before implementation

Exercise lagged and already-current PR heads; genuine title/body/identity/draft/
head changes; an interrupted update after a real Git push; dirty and advanced
local state plus an advanced base; remote-head changes; historical receipt tree
and repository mismatches; idempotent retry; and unknown post-mutation outcomes.
Retain repeatable JSON evidence from the isolated real-Git fixture.
