---
name: goal
description: >-
  Persist and pursue a concrete multi-step goal through written plans,
  evidence-backed attempts, honest acceptance checks, and strategy-changing
  retries. Use when success requires iterative work rather than a single edit;
  do not use for a simple one-step task or a read-only answer.
---

# Goal

Keep the objective, decisions, evidence, and failures in the repository so a
long-running task does not lose direction between iterations.

## Leave an artifact every turn

Here, a turn is one user-request/assistant-response cycle, not each internal
tool call, attempt, stage, or commentary update. Every turn that works on an
active goal must produce at least one new or materially updated, inspectable
artifact before yielding. Depending on the stage, this may be a candidate
model or export, render or comparison sheet, measurement or audit report, test
result, deterministic script, goal-record entry, or other concrete evidence of
work. A progress message by itself is not an artifact.

Surface that artifact in the user-visible session log during the same turn.
Always put a clickable local-file or external artifact link in the assistant's
own commentary or final message. An embedded visual preview is useful in
addition, but do not rely on a tool call or tool-result image alone because the
interface may collapse it. Resolve the link before yielding: verify that a
local target exists and use an absolute target in the session message when a
relative path would be ambiguous. Include a one-line label and honest verdict
so the user can inspect progress without opening the goal records or searching
the output directory. Do this for rejected artifacts too: they are useful
evidence when their visible failure changes the next plan.

Also record each material output and each reviewed candidate in a human-
oriented goal artifact log before the turn ends. Exclude caches, transient
implementation intermediates, the goal records themselves, and the artifact
log from logging themselves.
Put the inspectable path first and make it a Markdown link when the artifact is
on the filesystem. Give it a short label, say what it shows, and record its
review verdict. Keep the newest turn easy to find. A path written only inside
prose, a content hash, or a claim that an artifact exists is not an artifact
log entry the user can inspect.

Hashes are optional integrity metadata, not the primary progress report. Put
them after the artifact link or in the detailed evidence manifest; never use a
hash instead of a link to the render, comparison sheet, model, report, or test
output. Do not fill the live README with long hash lists. Put disposable
artifacts under the repository's designated temporary-output directory and
keep required deliverables in their owning project directory. Do not create
meaningless files or cosmetic churn merely to satisfy this rule: the artifact
must advance the goal, test an assumption, preserve a decision, or expose a
failure that changes the next plan.

## Start the goal record

Before implementation, identify the narrowest subproject root that owns the
intended result. Create a goal directory at
`<subproject-root>/goals/<goal-name>/`, following repository naming rules for
`<goal-name>`. The required entrypoint is `README.md`; do not place a monolithic
`GOAL.md` beside the deliverable.

Keep `README.md` short enough to serve as the live control plane. It must hold
the objective, status, current state, current plan, and a map of every
supporting record. Split durable detail into purpose-specific Markdown files
inside the goal directory, such as requirements, references and criteria,
failure ledger, artifact log, evidence manifest, current-attempt audit, and
numbered attempt history. Use an `attempts/` subdirectory when the history is
substantial. Keep the artifact log newest-first and optimized for opening the
actual outputs; keep dense hashes and machine-oriented measurements in the
evidence manifest.

Link every supporting record from `README.md`, and link back to the entrypoint
from each supporting file. Do not duplicate mutable current state across files:
`README.md` is authoritative. When resuming, read the entire `README.md` plus
the linked acceptance, current-attempt, failure, and evidence records before
choosing the next action. Re-split a record when its length makes current state
or evidence difficult to find.

Write or link these sections before beginning work:

- **Goal:** the user-visible outcome, not an implementation activity.
- **Acceptance criteria:** observable conditions that together prove success.
- **Constraints:** scope, permissions, protected artifacts, required tools, and
  explicit exclusions.
- **Evidence plan:** how each criterion will be tested or reviewed.
- **Fixed regression set:** checks that run after every attempt even when they
  are outside the targeted criterion.
- **Current plan:** ordered actions, with one action marked in progress.
- **Current state:** current candidate identifier and hash, parent checkpoint,
  stage, last accepted checkpoint, failing or unverified criteria, dominant
  failure, and exact next action.
- **Failure ledger:** stable defect names, occurrence counts, attempted
  strategies, and the latest measured result.
- **Requirement changes:** dated changes that distinguish user instructions
  from inferred acceptance gates.
- **Evidence manifest:** durable hashes, measurements, commands, and verdicts
  for temporary or external evidence.
- **Artifact log:** newest-first, directly linked, human-viewable outputs with
  a one-line description and honest verdict for each active-goal turn.
- **Attempt history:** initially empty; never erase failed attempts.
- **Status:** `in progress`, `achieved`, or `blocked`, with a short reason.

Do not weaken the user's goal to make it easier to pass. When the goal is
subjective, translate it into concrete comparisons and review gates while
preserving the user's stated quality bar. Record any inferred criterion as an
assumption, not as a user requirement.

## Execute one attempt

One attempt is one falsifiable hypothesis that produces one reviewed candidate.
Do not hide many variants inside one attempt. Before work, record the
candidate path or identifier, intended stage, parent checkpoint, every input
and parameter that can define the candidate, and the planned review packet.
After creation, record its SHA-256 or equivalent immutable content identifier.

Before freezing or committing an attempt plan, perform a read-only preflight
of assumptions that can be checked without producing the candidate. Keep this
proportional to the work, but resolve deterministic failures early:

- validate geometry, projection, coordinate, path, and numeric-gate math;
- inspect controlling references and any proposed target or guide pixels;
- check required APIs, tools, permissions, and runtime availability; and
- verify that the planned artifact can physically satisfy its own acceptance
  gates without copying a known rejected representation.

For subjective work, review a proposed guide independently before treating it
as a frozen target. For fragile automation, statically audit or compile dormant
code only after the plan's inputs and geometry are feasible. Do not build a
large driver to discover a contradiction that a short calculation or target
image would have exposed.

A draft rejected by this preflight may be corrected before the attempt is
frozen, because no candidate or reviewed attempt artifact exists yet. Record
material preflight evidence in the goal's evidence manifest, but do not inflate
the attempt ledger with discarded prose drafts. Once implementation begins, a
candidate artifact is produced, or the plan is declared frozen, preserve the
plan; a material correction then requires an honest new attempt.

For every attempt, append a dated or numbered entry containing:

1. the failure or uncertainty being targeted;
2. the hypothesis and why it should address that failure;
3. the exact plan written before implementation;
4. the work actually performed;
5. raw verification evidence, including commands, renders, measurements,
   reviewer findings, or other artifacts;
6. the result for every affected acceptance criterion; and
7. the decision: accept, refine, or reset the strategy.

A pass belongs only to the recorded candidate bytes. Any later change that
can affect a criterion invalidates that pass until it is rerun. Never combine
visual
evidence from one candidate with technical evidence from another.

After every attempt, update the current state, current plan, failure ledger,
and evidence manifest. An attempt is incomplete while any of those summaries
is stale. Temporary artifacts belong in the repository's designated
  temporary-output directory, not beside the goal record, unless the repository
explicitly says otherwise.

## Evaluate honestly

After an attempt, check the produced artifact itself. Successful execution of
commands proves only that the commands ran; it does not prove the goal was
achieved.

- Evaluate every acceptance criterion separately as `pass`, `fail`, or
  `unverified`, and link or name the evidence.
- Treat `unverified` as not achieved.
- Check the affected criteria and a fixed regression set after every attempt.
  Recheck all criteria at every milestone and for the final integrated
  candidate.
- Inspect likely failure views and edge cases, not only the strongest output.
- For subjective or visual work, use consistent comparisons and an independent
  review when available. Do not make the user the first quality-control pass.
- Prefer absolute quality against the target over relative improvement from the
  previous attempt.
- State visible defects plainly. Do not promote a candidate because substantial
  effort was spent on it.

Set status to `achieved` and write **Final result** only when every required
criterion passes and no required work remains. Summarize the delivered artifact
and the evidence that proves each criterion.

Before finalizing, freeze one candidate, rerun the entire evidence plan against
that exact content, perform any required commit or export, and verify that the
delivered content hash is identical to the fully tested hash. If delivery
changes the content, the final regression must run again.

## Audit progress and the current approach

After every attempt, and before choosing the next plan, add a **Progress,
approach, and process audit** to the goal record. Evaluate the whole delivery
loop, not only the candidate or hypothesis. It must answer:

- Which acceptance criteria measurably improved, regressed, or did not move?
- Did the artifact improve in absolute terms, or only relative to a poor prior
  attempt?
- Which evidence supports continuing the current approach?
- Which repeated defects suggest that the representation, tools, assumptions,
  task ordering, or validation method is wrong?
- What is now the highest-leverage unresolved problem?
- Should the next attempt continue, revise, or discard the current approach?
- Where did the iteration's wall time and avoidable rework go: setup, tool or
  process startup, implementation, rendering or generation, tests, review,
  coordination, evidence handling, or recovery from a late rejection?
- What is the current critical-path bottleneck, and which concrete process
  change could shorten the next feedback cycle without weakening any
  acceptance criterion or regression check?

Use measurements, comparisons, test results, or reviewer findings rather than
effort spent or intuition about future potential. For a long-running goal,
also re-read the full goal record every three attempts and whenever a defect
survives twice, so recent local improvements do not hide a stalled or
regressing overall result. Verify the current candidate and checkpoint hashes
when resuming after a pause or context change.

Continuously optimize the loop as part of pursuing the goal. Prefer the
cheapest artifact that can reliably falsify the current hypothesis: for
example, a calculation or guide before implementation, a low-resolution
failure view before a full render packet, a focused package test before a
repository-wide test, or a reusable persistent tool session before repeated
process startup. Cache stable inputs, make repeated review steps deterministic,
reuse validated fixtures, and parallelize independent work when shared mutable
state and review order remain safe. Instrument elapsed time or record a useful
estimate when the bottleneck is uncertain; do not optimize a step merely
because it is easy to automate.

Use progressive fidelity asymmetrically: a fast checkpoint may reject a bad
candidate, but it may not establish final acceptance unless it exercises the
full recorded evidence plan. Preserve the controlling references, fixed
regression set, final-quality settings, and implementation-blind review while
shortening earlier cycles. Treat false confidence, missed defects, stale
artifacts, mismatched candidate bytes, and review overload as quality losses,
not speed improvements. Record the process change and check in the following
cycle whether it actually reduced feedback time or rework. If the same
critical-path bottleneck survives two cycles, revise the workflow or tooling
rather than merely noting it again.

## Retry with learning

If any criterion fails or remains unverified:

1. keep status `in progress`;
2. evaluate the entire attempt history, identifying repeated failure patterns,
   invalid assumptions, and work that did not change the judged result;
3. update the goal record with that diagnosis;
4. replace the current plan with the next attempt; and
5. start again from the highest-leverage unresolved criterion.

A retry must change the hypothesis, representation, workflow, or validation
method when the same failure survives two attempts. Do not accumulate detail on
a structurally wrong foundation. Preserve working parts only when evidence says
they are not contributing to the failure.

Use the failure ledger's canonical defect names when counting recurrence. Do
not evade a strategy reset by renaming the same visible or functional failure.

Continue through the attempt/evaluation/retry loop while the goal remains in
progress and useful in-scope work is possible. Do not stop merely because an
iteration completed, the task is long, or the current result is better than the
previous one.

If progress requires authority, input, or an external state change that is not
available, record the exact blocker and safe attempts made. Do not label
ordinary difficulty, poor quality, or an incomplete iteration as blocked.
