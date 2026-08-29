# A90 plan — full-history progress and recovery review

## Binding and scope

This decision work unit begins after the requested rebase, against Reimu Fumo
goal generation 1, lifecycle generation 3, criteria revision 1, and Goal
resource version 55. It reviews the entire recorded effort through closed A89.
It does not edit, promote, texture, rig, or integrate any Blender candidate.

## Decision to test

The working proposal is that recent stagnation is not primarily a lack of
attempt count or compute. It is a failure to maintain one acceptance-visible
whole-plush baseline and to put sufficiently early absolute visual vetoes on
the authoring critical path. The review must try to disprove this and compare
credible alternatives before choosing a recovery strategy.

## Outputs

1. `FUMO_GOAL_PROGRESS_REVIEW.md`: an evidence-backed review of measurable
   progress, regressions, repeated defects, process latency, record quality,
   and the current honest acceptance state.
2. `FUMO_RECOVERY_PLAN.md`: a bounded sequence of decision-bearing work units
   with exact baselines, owners, views, time targets, keep/undo gates, and
   conditions for resetting the strategy.
3. Independent supporting audits covering attempt history, pixels/reference
   fidelity, end-to-end process, and recovery alternatives.

## Workstreams

- The coordinator owns the canonical record, conflict integration, final
  decision, and synthesis.
- A history reviewer quantifies A68–A89 and legacy attempt outcomes.
- A visual reviewer judges the recent exact artifacts against all controlling
  references without author rationale.
- A process reviewer identifies latency and workflow failures.
- An adversarial reviewer compares recovery strategies and seeks evidence
  against the coordinator's initial diagnosis.

All workers write only to `out/fumo_goal_review/`. Their reports are evidence,
not authority, and the coordinator must inspect and reconcile them.

## Acceptance and stop conditions

The review passes only if it:

- covers the complete goal history rather than extrapolating from A89;
- separates infrastructure, technical, module, and whole-result progress;
- states which exact candidate is the viable baseline and whether it is
  accepted;
- identifies the dominant systemic cause of low visible progress with direct
  evidence, not generic advice;
- compares at least three credible recovery strategies and issues an explicit
  `proceed`, `revise`, `ask`, or `refuse` verdict;
- produces a small, falsifiable next plan that prevents another long sequence
  of locally cleaner but globally poor fragments; and
- leaves every full-goal criterion honest (`fail` or `unverified` unless exact
  evidence proves otherwise).

Close early if the record lacks enough candidate-bound visual evidence to
support a strategy choice; in that case the plan must begin with a bounded
evidence-recovery step rather than inventing confidence.

Target: first synthesized decision-bearing draft within 45 minutes of worker
start. No modeling or archival expansion is allowed before that draft.
