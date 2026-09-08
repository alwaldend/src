# Reimu Fumo finish process

This contract replaces the ambiguous attempt semantics used by PR #24 and the
later local continuation.

## Roles

- The coordinator is the only writer of canonical goal state.
- One explicitly named local Blender MCP session or sponsor-approved native
  Blender artist is the sole writer of each model candidate.
- Measurement and technical QA may run in parallel on frozen copies.
- Two implementation-blind reviewers independently judge final-stage pixels.

## State vocabulary

- `technique valid` means an operation behaved as intended.
- `module retained` means a bounded component improved without a regression.
- `stage pass` means one exact candidate passed every gate for that stage.
- `goal accepted` means one exact candidate passed every current required
  criterion and exact-byte delivery was verified.

Only `stage pass` advances modeling tiers. Only `goal accepted` may populate
the goal's accepted attempt and result digest. An `accept` review decision for
a technique or module is not asset acceptance.

## Workspace handoff and user feedback

`out/reimu_fumo_finish/CURRENT.md` is the local artifact-discovery convenience
projection when the originating worktree is available. It is ignored, may be
absent in a fresh worktree, and never overrides the canonical goal or frozen
receipts. The coordinator updates it after a decision-bearing execution,
attempt, baseline, candidate, render, blocker, or acceptance change. It must
record:

- its observation and durable source identity;
- goal resource version, outcome, execution, and accepted attempt or `none`;
- protected baseline path and SHA-256;
- latest saved Reimu candidate path, SHA-256, and `accepted`, `rejected`, or
  `baseline` label;
- latest reviewable Reimu render path with the same label;
- terminal capability result, blocker, resume condition, and next authorized
  action; and
- the durable goal and pull-request links.

Never point a “latest Reimu” field at an unsaved viewport state or a technique
coupon. A direct request for status, a render, or a model path is answered in
the next visible message from the last frozen state before unrelated work
continues.

Local artifact lookup in the originating worktree requires only this asset
README and the workspace index. Normal task orientation still follows the
repository instructions. Before mutation, read the canonical goal, current
criteria, active attempt when one exists, and the skills for the current
phase. Do not load historical attempts or later-phase procedures unless the
current decision needs them.

Routine attempt starts, setup repairs, and ordinary continuation use local
goal-tool checkpoints with the required frozen plans and evidence. Verify the
checkpoint succeeded and its artifacts are accessible before continuing. A
pause or same-worktree handoff does not require a commit or push. Preserve
local records and update `CURRENT.md`; distinguish a local checkpoint from a
verified remote copy.

Use repository delivery at meaningful review-ready milestones, final
implementation handoff, an explicitly requested remote backup, or before a
handoff to another checkout that cannot access the required local artifacts.
For that remote-dependent handoff, commit and push the necessary records and
permitted artifacts, verify remote reachability, and establish that the
recipient can access the exact required bytes before it resumes. Publication
does not transfer writer authority or replace live binding checks.

## Modeling gates

Work advances in this order:

1. authoring capability;
2. whole-character clay silhouette and proportions;
3. sewn construction, attachment, contact, and overlap;
4. secondary compression, seams, gathers, and controlled asymmetry;
5. final materials and facial graphics;
6. reusable structure, rigging, and animation;
7. exact-byte repository delivery.

The authoring-capability gate requires a bounded visible edit made through the
writer named by the active plan, a save to a new candidate, a clean reopen and
render in pinned Blender, and unchanged protected-input hashes. A failure
blocks modeling until a new capability or an explicit sponsor-approved writer
supersedes the plan.

## Preflight and attempt boundaries

A genuinely different authoring route may run disposable setup checks while
execution remains blocked, but they never open Reimu bytes or exercise the
decision-bearing organic operation. A sponsor-approved artist satisfies the
writer-authority branch directly. Provisioning a materially different route
satisfies only the external-availability condition: it permits the coordinator
to reopen execution for one bounded capability investigation per named route
and hypothesis, not to accept that capability.

Before exercising the organic capability question, create a new active portable
goal plan or use the already-active plan whose strategy matches the route and
hypothesis. If execution is `blocked`, use a separate plan-only, no-attempt
checkpoint to return it to `active`; do not combine that transition with
attempt start. If both the intended plan and active execution are already
current, emit no checkpoint. Then bind the plan and start a durable
investigation attempt. A pass closes the investigation with an `accept` review
decision, accepts the settled capability plan, keeps goal outcome `open` and
execution `active`, and may then authorize modeling through that route. A
failure closes the attempt with its negative evidence, rejects the settled
route plan with the failure reason, and returns execution to `blocked` only
when no other authorized route is immediately actionable because a named
authority, input, or external state is unavailable; otherwise execution
remains `active` for the ready route. Causally justified setup repairs remain
checkpoints in the same attempt while the route, hypothesis, and gate remain
unchanged.

Every checkpoint that changes execution to `blocked` must record the named
unavailable authority, input, or external state and a concrete, testable resume
condition in its durable result or evidence and in the ignored `CURRENT.md`
projection. When those values are known before a newly created terminal
investigation starts, also publish them through its structured `blocker` and
`resumeCondition` fields. Those attempt fields are immutable: when a blocker is
learned while closing an existing capability or modeling attempt, preserve it
in the close result and evidence instead of pretending to update the frozen
attempt specification.

After a route satisfies the resume condition, run checkpoint zero before a
model attempt. Use disposable bytes and check only contracts the planned work
will consume:

1. bind stable writer and verifier identities and the writer session and
   workspace; record the Blender version and build, process, launch flags,
   add-on or listener, and authoring tool;
2. open or reset the disposable input in one writer call without editing it;
3. in a later call, reacquire and check the window, area, region, active
   object, selection, mode, and operator context needed for that operation;
4. discover version-dependent enums and nullable scene data only when the
   writer path consumes them, and activate and read back only required brushes
   or operators;
5. for geometry operations, keep the fixture below 0.5 m, hash coordinates
   after Blender stores them using packed float32 values, and prove one
   localized nonzero edit plus the planned undo path; and
6. when saved-byte interchange is part of the route, save to a fresh path,
   clean-open it with the existing repository-pinned Blender, compare exact
   file and geometry digests, and perform the smallest required smoke render.

Before checkpoint zero, record a session-local preflight plan with the exact
checks, a numeric causal-repair budget, and a terminal stop condition. A repair
consumes that budget only after a checkpoint names the failed check, causal
change, and predicted result. If the budget is exhausted or a check is
irreparable, the terminal failure becomes a decision about the route. Prepare
uniquely named Markdown plan, result, and evidence inputs that preserve the
decision-relevant failure receipt and digest. First create its new active
portable goal plan, or use its already-active portable goal plan
whose strategy matches the frozen preflight plan. Use a plan-only, no-attempt
checkpoint for plan creation and/or a required
`blocked`-to-`active` transition. If both the intended plan and active execution
are already current, emit no checkpoint. Then bind that plan ID and start and
close one bounded `investigation` attempt with a `reset` review decision,
reject the plan with the terminal reason, and keep outcome `open`. Set execution
to `blocked` only when no other authorized route is immediately actionable
because a named authority, input, or external state is unavailable; otherwise
keep execution `active` for the ready route. Leave acceptance pointers unset
and do not open Reimu bytes. Verify the closed local investigation and update
`CURRENT.md` with the retired route and its evidence. Publish according to the
milestone and handoff rules above, and report whether the retirement record is
local or verified remote. A publication failure does not reopen the settled
route decision. Do not rerun checkpoint zero unless a new route or hypothesis
is authorized.

Write a passing receipt last under
`out/reimu_fumo_finish/preflight/<session-id>/receipt.json`. Before attempt
start, create a new active portable goal plan or use its already-active portable
goal plan whose strategy matches the frozen attempt plan. Use a plan-only,
no-attempt checkpoint for plan creation and/or a required
`blocked`-to-`active` transition. If both the intended plan and active execution
are already current, emit no checkpoint. Do not combine a required execution
transition with attempt start because attempt publication checks the
pre-mutation execution state. Then bind the receipt to work in this order:

1. Prepare a durable evidence copy containing its decision-relevant bindings,
   checks, results, source writer workspace, and digest. The ignored receipt is
   a local diagnostic; its path and hash alone are not durable evidence.
2. Freeze an attempt plan that references the receipt's exact path and SHA-256
   and the durable evidence copy, and bind its portable goal plan ID. Goal
   execution must already be `active`.
3. Start the attempt with a goal-tool checkpoint that records the frozen plan
   and evidence copy in the attempt. Verify the checkpoint succeeded and its
   exact evidence is accessible before the first decision-bearing operation.
4. Before handing off an open attempt that used any preauthorized setup repair,
   append every contemporaneous repair checkpoint to tracked attempt evidence
   with the goal tool. Verify the local record before ordinary continuation or
   a handoff with access to the same artifacts. For a remote-dependent handoff,
   also complete the publication and accessibility checks above; the recipient
   must wait if the required evidence is unavailable.
5. Treat checkpoints and publication as evidence preservation, not session
   authority. Before a recipient operation, compare the bound writer and
   verifier identities, writer session and workspace, process, version, build,
   launch flags, add-on, listener, authoring tool, and procedure against the
   live environment.

Any mismatch or change in a stable binding invalidates the receipt, including
across a handoff before an attempt opens. A publication failure alone does not
invalidate an unchanged local receipt or require repeating preflight. Preserve
the local checkpoint and resolve publication before a remote-dependent
handoff; a recipient without the required artifacts must wait.

When no attempt is open, discard an invalid receipt without an attempt
mutation, rerun checkpoint zero in the recipient writer workspace, then follow
the passing or terminal flow above. If a receipt bound to an open attempt
becomes invalid for any reason, immediately suspend
decision-bearing work. Do not rewrite its frozen plan or resume that attempt.
Close it through a goal-tool checkpoint with a `reset` review decision, keep
goal outcome `open`, leave acceptance pointers unset, and record
`setup-invalidated`, the invalidator, and the relevant failure evidence in setup
evidence. Keep execution `active` when replacement checkpoint zero is
immediately actionable. Set it to `blocked` only when a named authority, input,
or external state makes replacement preflight unavailable, and record that
dependency as the resume condition. If execution stays active, rerun checkpoint
zero immediately without opening Reimu bytes. If it becomes blocked, wait for
the resume condition, use a separate no-attempt goal checkpoint to return
execution to `active`, and then rerun checkpoint zero. After the replacement
receipt passes, perform the common plan and execution transition and all five
checkpoint and binding steps above before work resumes.

Recheck dynamic window, selection, mode, and operator context immediately
before each operation instead of treating them as durable session properties.
Event simulation, live rendering, World setup, keymap checks, and undo tests are
required only when the planned writer path consumes them.

Before an attempt opens, a repairable preflight or setup fault that touches no
decision subject receives no attempt ID, criterion verdict, or acceptance
status. Repair it causally within the frozen preflight budget; the eventual
attempt plan binds the passing receipt. Budget exhaustion or an irreparable
fault decides the route and uses the immediately closed investigation attempt
defined above. After an attempt opens, bound-receipt invalidation closes that
existing attempt as the preceding paragraph requires, even if no
decision-bearing operation occurred. Before opening the attempt, bind an exact
subject, decision question, hypothesis, writer, acceptance gate, and stop
condition. The attempt-start checkpoint described above opens the durable
attempt; decision-bearing work begins only after that local checkpoint
succeeds and its exact evidence is accessible.

An attempt plan may preauthorize a bounded number or class of causal setup
repairs without predicting an unknown failure. Immediately before using that
budget, the writer appends a checkpoint naming the observed failure, exact
causal change, and predicted observable result. The original plan plus that
contemporaneous checkpoint authorizes the repair; the attempt freezes the
checkpoint into its evidence at close. Do not reconstruct this authorization
after the repair or rewrite the frozen plan.

Once an attempt has started, a setup repair stays inside that attempt when
those bindings remain unchanged. Allocate a new attempt only when the
hypothesis, subject, writer, acceptance gate, or strategy materially changes.
Bound-receipt invalidation is the forced-lifecycle exception: it closes the
existing attempt and requires a replacement attempt as specified above, even
when every decision binding is otherwise unchanged.

Preregistered dose blocks and other repeated measurements of one unchanged
question remain checkpoints in one attempt. If a repair changes the strategy,
close the old attempt honestly before starting another.

Rigging, UV work, and material polish do not begin before the whole-character
clay candidate passes. A subsystem checkpoint may be retained, but it does not
advance the whole-character stage while another macro component fails.

## One fidelity cycle

Each cycle starts from a frozen candidate digest and targets one named dominant
failure. Before editing, record the representation, the expected multi-view
effect, the stop condition, and the likely regression view. Then:

1. have the named sole writer edit one coherent subsystem;
2. save a new candidate path and hash it;
3. clean-reopen it with pinned Blender;
4. render front and the highest-risk view into a fresh directory;
5. measure the aligned comparison against the exact landmarks and reference
   bytes bound by `review_contract.json`;
6. render the full fixed-view set only if the fast pair passes; and
7. record `retain` (structured close decision `accept`), `refine`, or `reset`
   with the exact candidate digest.

Relative improvement is regression evidence only. It cannot pass an absolute
stage or goal criterion.

The model writer's self-review may reset a candidate but may not retain one.
Every `retain` decision requires at least one implementation-blind reviewer of
the frozen pixels. A whole-character stage pass and final acceptance require
two independently identified reviewers whose stable IDs and roles are recorded
in the candidate packet as required by `review_contract.json`. If independent
review overturns a provisional call, the structured result follows the
independent review and records the disconfirmation.

For each frozen checkpoint, use only roles whose result can change the next
decision: one coordinator, the sole writer, at most one quantitative verifier
for a decision fact, and at most one adversarial reviewer when the decision is
material. Spawn visual reviewers only after frozen pixels exist; retain the two
independent reviewers required for a whole-character stage pass or final
acceptance. Do not recompute a settled fact more than once unless the second
computation uses genuinely different evidence or method. Stop nested fan-out
when its outputs would enter the same review queue or restate an existing
receipt.

## Stop and reopen rules

Two reviewed failures retire a representation family only when causally
distinct repairs both support the same representation-level limitation.
Repeating the same visible defect is not enough while its cause remains
unresolved. A later cycle may reopen the family only through a new named plan
that records a new capability, new evidence, or a falsifiable causal repair and
supersedes the old plan. A failed gate cannot receive an unplanned correction
in the same cycle.

A settled result forbids rerunning the same operative setup: the same target,
relevant inputs, causal path, and decision question. Changing labels, command
spelling, flag order, tools, or parameters does not make a new setup unless the
change could repair, bypass, or discriminate the named failure mechanism or
hypothesis.

After a failure mechanism is identified or narrowed to a falsifiable causal
hypothesis, a changed input or execution path is not the same setup when it
could repair, bypass, or discriminate that mechanism if the hypothesis is
true. Its plan must reference the prior failure, name the causal change, and
state the observable difference it predicts. Retain only the evidence needed
to establish the prior outcome and the hypothesis; preserving every failed
candidate byte is unnecessary. If the proposed change cannot explain or test
why the result should differ, stop and treat it as the same setup.

For example, resending equivalent simulated number keys under a new attempt ID
after Blender ignored them is the same setup. Supplying the measured vector
through Blender's native operator value after diagnosing that input-parser
failure is a permitted causal repair, even when it uses the same source model,
selection, and desired displacement.

Every third closed cycle, or immediately after a stalled result, the
coordinator reviews the complete trend and changes the process only when the
change shortens feedback or improves acceptance reliability.

When every authorized representation or authoring route for the dominant
failure has met its stop condition, the coordinator sets goal execution to
`blocked` with one concrete blocker and resume condition. Renaming the same
geometry family or adding harness work is not a new capability.

## Milestone packet

Stage-pass and final-acceptance packets become durable and must be immutable
and self-contained:

- the exact candidate `.blend`;
- front, side, rear, both three-quarter views, and one presentation render;
- controlling reference copies or repository-relative reference identities;
- aligned comparisons or masks used for measurements;
- a manifest binding candidate, Blender, scene, cameras, settings, and every
  output digest;
- the scorecard and concise verdict; and
- for final acceptance, two independent reviewer matrices.

The render producer writes into a new directory and publishes readiness last.
Reused or partially populated output directories are invalid.

A compact terminal-failure packet may also be retained when its pixels explain
a representation-family reset. It must say that it is rejected, whether the
candidate bytes are published, and whether the renders are independently
reproducible. Failure evidence cannot satisfy or inherit an acceptance
criterion.

## Goal closure

The repository goal tool is the sole goal-state writer. A result narrative
cannot override structured verdicts. Final closure requires every current
criterion to pass against one candidate digest, followed by verification that
the committed Git LFS object is byte-identical to that reviewed candidate. The
public controlling references remain tracked at the user's explicit direction
under the fan-work notice in the asset README. Their per-image provenance stays
informational and conservative: an acquisition locator or content hash does
not claim ownership or relicensing authority, and an `unverified` license field
is neither a model failure nor an acceptance pass.
