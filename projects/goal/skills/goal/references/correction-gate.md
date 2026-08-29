# Correction gate

Use this contract only when the goal skill opens or resumes a correction
episode. It turns independent review into a live control on the next decision,
not a retrospective report or a new source of authority.

## Bind one review episode

Bind the episode to the goal generation, lifecycle generation, planned or
active attempt ID, named stable defect or assumption, frozen immediate
comparison, best relevant historical comparison, and first decision-bearing
artifact. Coalesce overlapping triggers that can affect the same decision.
Record whether the slot is `open`, `complete`, or `degraded`, plus the accepted
reviewer result identity. While it is open, also record the worker identity,
packet digest, and requested and effective model, reasoning effort, and fork
mode from the routing receipt. On resume, reconcile that receipt with the
runtime before spawning anything; do not reopen a completed or active slot
merely because the session changed.

Exactly one reviewer may be active and exactly one result may be accepted for
an episode. A crashed, timed-out, or malformed worker may be replaced once in
the same slot. One focused follow-up to that same reviewer and result may
clarify missing evidence; it must not spawn another worker. Neither case
creates another review episode or a delegation cascade. This process review
does not replace specialist or independent reviews required by the acceptance
criteria.

## Route independent judgment

When the coordinator runs below the strongest reasoning tier available for
the task, use exactly one fresh reviewer at that stronger tier if the runtime
supports a per-agent override. Keep authoring and ordinary workers at the
coordinator's tier. Use a no-history or minimally bounded fork so the override
and implementation blindness survive; provide an explicit packet rather than
the author's conversation.

When the coordinator already runs at the strongest available tier, use one
fresh independent reviewer at that same tier with the same context-isolation
and authority boundaries. Strongest-tier coordination does not bypass a gate
that the goal skill has opened.

If the stronger route or a reviewer slot is unavailable, use an independent
reviewer at the inherited tier when possible and mark the episode `degraded`.
If independent review is required for safety or acceptance, wait or block
honestly. Otherwise continue only with a stronger deterministic discriminator
and record the lost independence. Never claim a stronger-tier review that did
not occur.

## Supply a bounded packet

Give the reviewer immutable, task-authorized inputs:

- objective, acceptance criteria, and authority and safety boundaries;
- exact goal, attempt, baseline, candidate, and artifact identities;
- controlling requirements or references and the stable-defect history;
- factual outcomes of the relevant prior attempts;
- the frozen immediate comparison, the best acceptance-visible historical
  comparison under the same criterion revision, calibrated view, and evidence
  method for each affected criterion, and any prior module that passed a role
  gate relevant to the current decision;
- the proposed hypothesis or operation without its desired verdict;
- the first-artifact target and the produced artifact when available; and
- the fixed comparison or regression evidence needed for the decision.

Do not include the author's rationale, preferred verdict, mutable canonical
access, unrelated history, or instructions to improve or publish the result.
Tell the reviewer not to mutate, publish, interrupt workers directly, expand
scope, or delegate.

If an older artifact was judged under a different criterion revision, view,
or evidence method, re-evaluate it under the current calibrated method before
calling it the historical high-water mark. Otherwise mark that comparison
unavailable and do not infer a rank from incomparable evidence.

## Require a compact result

The reviewer returns only:

1. the bound episode and artifact identity;
2. one verdict: `CONTINUE`, `STOP`, `RESET`, `SMALLER_PROBE`, or `UNVERIFIED`;
3. the decisive acceptance-relevant evidence, including deltas against both
   the immediate comparison and the historical high-water mark;
4. material limitations; and
5. the next discriminator or resume condition.

An improvement over a recent regression is not progress when it remains below
a viable historical high-water mark for the same criterion. Treat it as
recovery unless current acceptance evidence shows net advancement across the
active required criteria, including any explicit bounded trade-off. Finding or
reusing an older passing module may justify the next bounded integration probe,
but it is not result progress until the integrated candidate is measured.

`CONTINUE` permits only the already planned bounded work. `SMALLER_PROBE`
replaces expansion with the named cheaper discriminator. `STOP` and `RESET`
suspend discretionary expansion at the next safe checkpoint. `UNVERIFIED`
does not count as approval and requires the named evidence recovery or an
honestly degraded fallback.

## Adjudicate before continuing

Resolve the result before the earliest of:

- expanding beyond the first decision-bearing artifact;
- taking a costly, destructive, or irreversible step;
- closing or promoting the current attempt; or
- starting the next attempt.

The coordinator records `follow`, `override`, or `degraded-fallback` with the
reviewer result. It may override `STOP`, `RESET`, or `SMALLER_PROBE` only with
cited contrary evidence that addresses the decisive finding; schedule pressure,
sunk cost, author confidence, or a successful command is not contrary
evidence. The coordinator remains responsible for the verdict, criterion
state, canonical writes, and delivery.

Record only the trigger, slot and routing receipt, packet digest, compact
result, and disposition in the existing attempt plan, result, or evidence. Do
not create a review-only attempt, dashboard, general infrastructure, or long
report, and do not count review activity as result progress.
