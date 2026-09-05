---
name: goal
description: >-
  Persist resumable research or implementation work with local checkpoints
  and evidence. Use when work benefits from state that survives interruption;
  skip simple tasks and one-response questions.
---

# Goal

Keep the objective, acceptance check, current candidate, evidence, and next
action inspectable. The CLI owns record identity, revisions, digests, locking,
and interrupted-write recovery. Use it instead of editing canonical YAML or
generated README projections.

## Ordinary tasks

Use an ignored workspace goal under `out/<task>/goals/<goal-id>/`. A project
goal under `<owner-root>/goals/<goal-id>/` is appropriate only when the user
requests, or the project establishes, maintained history. Length alone does
not justify committed records. A question grants no new mutation authority.

Select an existing goal by stable ID and path before initializing another.
When identity is uncertain, use `list --goals-root <root>`. Initialize once:

```sh
bazel_agent bazel run //projects/goal/cmd/goal -- init \
  --goals-root "out/<task>/goals" \
  --title "Requested outcome" --criterion "Observable acceptance check"
```

For routine progress, supply a short summary, the exact candidate identity,
and the next action. Mark missing checks and inferred criteria honestly:

```sh
bazel_agent bazel run //projects/goal/cmd/goal -- checkpoint \
  --goal-dir "<record>" --expected-resource-version "<observed-version>" \
  --subject "<exact candidate>" \
  --summary "<work performed, evidence, and checks still needed>" \
  --next-action "<next concrete step>"
bazel_agent bazel run //projects/goal/cmd/goal -- show --goal-dir "<record>"
```

`--summary` accepts up to 8192 bytes of inline Markdown and creates ordinary
plan/result artifacts internally. Optional `--evidence <file.md>` imports
immutable evidence. Later checkpoints replace the open attempt's result and
update its candidate and next action. The initial plan and imported evidence
remain intact. Routine progress needs no separate plan ID, payload file,
process audit, or new attempt.

`show` returns objective, criteria, current resource version, and active
attempt progress with a bounded result preview, source paths, and digests.
Inspect the referenced result when truncated. Candidate and prose fields are
caller declarations, not live observations or proof of acceptance. Revalidate
only mutable inputs required by the recorded next action after interruption.

Carry the literal observed `--expected-resource-version` into each write.
If stale, reread and reconcile the intervening changes; never blindly replace
it with the latest version. Keep one coordinator for canonical goal writes.
A changed candidate needs a new summary or result; old evidence applies only
to the candidate it tested. Neither a checkpoint nor a successful command
establishes acceptance.

## Finish or continue

Keep outcome (`open`, `achieved`, `abandoned`, `superseded`) separate from
execution (`active`, `paused`, `waiting`, `blocked`). Questions and scheduling
interruptions do not silently pause a goal or expand its authority. Honor
explicit pause or cancellation and resume authorized work after incidental
interruptions. Difficulty alone is not a blocker.

At acceptance, verify the exact result against every required criterion and
applicable regression. Close the attempt with `checkpoint --close-attempt
--review-file <review.yaml>` and the final result/evidence; use `--outcome
achieved` only when all required criteria pass. The review records criterion
IDs, revisions, verdicts, and evidence references. Read the
[close review format](references/record-format-v1alpha1.md) when preparing it.
Closed attempts are immutable; corrections use a new attempt.

Checkpoint locally when work needs a reliable continuation point. A
conversational yield does not require committing, pushing, or running the
full delivery workflow. Use `repo-delivery` at a meaningful review-ready
milestone or final implementation handoff, and for an explicitly requested
remote backup. Deliver all legitimate source changes before final handoff;
keep scratch and workspace records ignored. Report a blocked or incomplete
publication honestly and distinguish a local checkpoint from a verified
remote copy.

## Detailed and conditional workflows

For maintained project goals, repeated failures, or complex coordination,
read [lifecycle and evidence](references/lifecycle-and-evidence.md) for bounded
plans, attempt histories, review, and strategy changes. After the same defect
survives two attempts, change the approach and preserve its failure evidence.
Process reviews support that work; they are not routine checkpoint fields.

- For checkpoint payloads, criteria changes, or validation errors, read
  [record format](references/record-format-v1alpha1.md).
- For registered goals, catalog-backed `resume`, session switching, or worker
  coordination, read
  [sessions and concurrency](references/sessions-and-concurrency.md).
- For related goals and dependency-aware dispatch, read
  [graph organization](references/graph-organization.md).
- For promotion, imports, or public evidence retention, read
  [promotion and migration](references/promotion-and-migration.md).
