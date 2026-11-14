---
name: comfyui-experiments
description: >-
  Develop ComfyUI experiment directions through evaluation, refinement,
  forking and explicit rejection on an existing endpoint. Preserve original
  outputs, lineage and job ownership. Deployment and model installation
  belong to their infrastructure owner.
---

# Develop ComfyUI experiments

Read `projects/comfyui_tools/docs/experiments.md` for the generic command and
producer contract. Keep subjects, prompts, models and experiment directions in
caller-owned inputs rather than runner source.

## Establish the contract and directions

Read the request, existing direction record, run journal and endpoint
capabilities. Preserve authority for endpoint use, spending, batch size and
rental retention. Experimentation does not authorize provisioning, replacing
or destroying an instance.

Maintain a task-local record of independent directions. Give each a stable ID,
a descriptive name beginning with its primary theme or hypothesis, acceptance
criteria and a revision. Record parent direction and source job IDs for forks.
The direction record owns hypotheses, evaluations and decisions; the runner
journal owns accepted jobs, exact workflows and native prompt IDs. Link those
records instead of maintaining competing copies of execution state.

Define what each proposed job tests, what changes and what result would support
or undermine the idea. Queue one job for each ready direction. Preserve the
authorized batch size. A direction may have only one submitted revision awaiting
evaluation; seed changes alone do not create an independent direction.

## Evaluate before submitting successors

Inspect each finished job's original outputs against its direction's criteria
and the latest user feedback. Execution success is not acceptance. Record the
observed strengths, defects, decision and reason before preparing a successor.
Choose explicitly:

- **Scrap:** retire the direction. Preserve its results and reason, but never
  regenerate it or cycle back to it automatically. Reopen it only when the user
  explicitly requests that.
- **Refine:** keep the direction ID, increment its revision and link the exact
  source job. Modify prompts, settings or conditioning to address a specific
  observed defect or test a stated hypothesis, then queue one revised job.
- **Fork:** create several children with separate IDs, names, hypotheses and
  changes, each linked to the evaluated parent and source job. Queue one job
  per child. Evaluate and develop each independently; do not silently keep
  resubmitting the parent or couple future sibling changes.

If the direction already meets its goal, retain the accepted result and stop
that branch. If the evidence is inconclusive, record that uncertainty rather
than inventing a positive verdict. Independent directions need not finish
simultaneously: review a completed direction while others execute.

Do not endlessly queue experiments, invent quota-filling successors or use a
cyclic producer that revisits earlier ideas. The loop is direction, job,
evaluation, decision, then a justified successor. Enlarging a batch or queue
must not bypass evaluation. A utilization target does not authorize automatic
follow-ups to unevaluated results; let the queue drain when no reviewed next
step is ready and report that state.

## Execute a finite set of ready jobs

Use finite `max_experiments` for the ready directions and a producer that
returns `{"done":true}` when that set is exhausted. A new reviewed set uses a
new finite run linked by direction and source IDs. Resume the same run only
to recover its existing work; do not evade ambiguous ownership with a new run.
The runner's continuous mode is not the direction-development procedure and
must not be inferred from an old request to keep a rental busy.

Build the runner and CLI launchers first. Use compiled executables or
`bazel run --script_path=...` launchers; do not put nested Bazel invocations
inside a long-lived controller. Configure lookahead only within the finite
ready set. One executing job plus a pending buffer is useful when independent
ready directions exist, not an obligation to manufacture more work.

Supply candidates through the external producer contract. Include direction
ID, revision, lineage, purpose and changed controls in opaque notes, alongside
requested prompt/settings provenance. Declare JSON-pointer identity exclusions
for seeds, output prefixes and irrelevant bookkeeping, keeping meaningful
controls in the identity. Canonical identity detects declared duplicates; the
agent must judge whether a variation answers a useful question.

Inspect actual endpoint state to distinguish executing from pending work.
Local `queued` means an owned job awaiting a terminal observation. Never claim
permanent supervision because a local controller was started; retain its ID,
native job IDs and the time and scope of queue observations.

## Preserve results and recover within ownership

Keep original media, hashes, exact workflow snapshots, prompts, seeds, settings
and the evaluation history. Make every result accessible in the requested
queue or gallery. Honor the user's display preference; do not duplicate images
in chat when they inspect the queue. Collapse verbose metadata when requested.
Inspect private inputs before publishing; credentials and raw diagnostics stay
private.

A valid successful native envelope followed by wrapper cleanup failure remains
acknowledged; record the warning and do not resubmit. A missing submission
acknowledgment stops automatic submission: inspect its private receipt and
endpoint ownership rather than blindly retrying. Observer transport errors and
timeouts retain ownership until a native outcome is known. Native terminal
errors may be carried in unsuccessful envelopes and must not be treated as
mere transport failures. Best-effort cancellation is not confirmed termination;
reconcile the exact native ID and retain completed outputs.

When the user stops or pauses experimentation, stop new submissions immediately
and retain that instruction across restarts. Drain already submitted work unless
cancellation is requested. `stop --cancel-owned` targets recorded IDs through
the native guard; an idle run consumes the request on resume. Interruption
preserves remote jobs, and the advisory lock prevents concurrent controllers.
Never clear a global queue, interrupt foreign work, restart a service or change
rental lifecycle to manage experiments. Report stopped or uncertain states
accurately and wait for the user's instruction before resuming a paused stream.
