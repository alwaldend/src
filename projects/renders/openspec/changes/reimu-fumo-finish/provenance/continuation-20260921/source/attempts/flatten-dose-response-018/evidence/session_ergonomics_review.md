# Reimu continuation agent-ergonomics review

Observed on 2026-09-04 at branch head
`e47f0102b19bbc45da7cb804a832d0f5b40450b1`, with uncommitted goal
resource version 57 and `flatten-dose-response-018` still active. The stable
task is `reimu-fumo-finish`. The exact delivered subject remains the protected
A157 scaffold plus recovered references and evidence; there is no accepted
Reimu asset or retained improved model.

## Outcome

The continuation materially improved truthfulness, provenance, native Blender
control, and resume safety. It did not yet improve the product. The largest
ergonomic failure is that internal experiment throughput made the user ask for
an immediate status, render, and file link instead of keeping those pointers
prominent throughout the work. The current process can now distinguish a
causal retry from a renamed replay, but its goal-record and multi-agent
granularity are still too fine.

## Compact scorecard

- Resume integrity, 4/5: exact branch, hashes, goal version, active attempt,
  writer, and next action survived compaction and interruption.
- Provenance and verification, 4/5: frozen inputs, candidate hashes, pinned
  clean-open/render, control isolation, and exact undo are recorded.
- Causal retry discipline, 4/5: later retries name a mechanism and prediction;
  attempt 018 changes only cumulative dose.
- Tool routing, 3/5: Cordis, the goal CLI, Blender MCP, and pinned Blender were
  correctly separated, but version/setup mismatches caused avoidable failures.
- Discovery and context cost, 2/5: a minimum resume set plus the 11 selected
  skills totals 2,146 lines and 128,147 bytes.
- Multi-agent efficiency, 2/5: useful role separation coexisted with duplicate
  analysis and six interrupted branches.
- Artifact discoverability, 2/5: exact artifacts exist, but 382 scratch files
  have no single current-result pointer.
- User-feedback latency, 2/5: the render and file link were answered in the
  next visible messages, but only after the user interrupted to request them.

## Measured path

- The minimum current-state set (`AGENTS.md`, `PROCESS.md`, goal README, and
  active plan) is 770 lines and 50,116 bytes. The 11 selected skill entry
  files add 1,376 lines and 78,031 bytes.
- The successor goal now contains 18 attempt directories: 17 closed and one
  open. Its 81 files occupy 174,962 bytes, including 1,883 Markdown lines.
- Since the committed recovery checkpoint, the worktree added ten attempt
  IDs, 009 through 018. Nine are closed and 018 is active.
- Task scratch contains 382 files in 23 top-level directories and occupies
  235,026,682 bytes, about 224 MiB. Across goal and scratch there are 34
  `.blend` files, 121 PNGs, 38 JSON files, and 132 Markdown files.
- `out/reimu_fumo_finish/agents` alone occupies 84,072,622 bytes. At the
  review snapshot the collaboration tree exposed 20 non-coordinator agents:
  12 completed, six interrupted, and two still working, including this
  reviewer.
- The prior closeout review recorded seven bounded attempts and about 142 MiB
  of scratch. The new work therefore improved capability evidence while
  adding ten goal attempts and roughly 82 MiB of scratch.

## Ranked findings

### 1. `REIMU-FEEDBACK-002` — internal progress hid the usable answer

Tier: `live`.

The user had to interrupt an active capability experiment to ask for current
progress, a render, and then a file link. The coordinator answered both direct
artifact requests in the next visible messages, using the clearly labelled
rejected attempt-014 render. That response was timely once requested, but the
same usable artifact and the fact that attempt 015 produced no saved candidate
were already known and should have remained prominent without prompting.

Owner change: the coordinator should answer an artifact/status request from
the last frozen state immediately, then continue background work. Add a
session convention that `CURRENT.md` is updated whenever the active candidate
or attempt changes, with status, exact hash, render link, model link, and an
explicit `accepted`, `rejected`, or `baseline` label.

Measure: one bounded read should yield a clickable current render and model
path; a direct user artifact request should be answered in the next visible
message, without waiting for an unrelated experiment.

### 2. `REIMU-ARTIFACT-002` — exact evidence is difficult to find

Tier: `live`.

The scratch tree is internally well hashed but externally shaped like an
archive: 382 files, 23 top-level directories, 34 blends, and 121 images. The
goal README identifies the active attempt but not the best currently viewable
Reimu render or file. Hashes provide integrity after discovery; they do not
provide discovery.

Owner change: generate one task-local `CURRENT.md` and `artifact_index.json`
from goal state. They should distinguish protected baseline, latest saved
candidate, latest rejected visual, active work, and accepted artifact. Never
point `latest` at an unsaved state or a technique coupon as if it were Reimu.

Measure: a cold-start agent should locate the correct model and render in one
read and no directory search. Retire the index only when the goal tool exposes
equivalent fields in its bounded README projection.

### 3. `REIMU-ATTEMPT-002` — setup events are too often promoted to attempts

Tier: `fixture-tested`.

The goal record is honest, but ten new attempt IDs for this continuation is a
high coordination tax. Attempt 016 became a complete goal attempt despite
stopping before mesh creation on an irrelevant Blender enum and absent World.
Attempts 017 and 018 are a valid hypothesis progression, but their evidence
is spread across separate plans, results, receipts, and reviews.

Owner change: reserve a goal attempt for a changed hypothesis, target, writer,
or acceptance gate. Put preregistered dose blocks and repaired setup events in
append-only checkpoints under the same attempt. A setup fault that produces no
subject or decision evidence should not consume an attempt ID unless it changes
the strategy.

Measure: one capability question should normally require one attempt record;
its setup checkpoints should be discoverable from that record. Revisit this
rule if checkpoint histories become harder to audit than separate attempts.

### 4. `REIMU-PREFLIGHT-001` — Blender compatibility was tested late

Tier: `live`.

The Flatpak Blender 5.1.1 authoring host and pinned Blender 5.2.1 verifier are
a sound boundary, but several failures were avoidable:

- opening a file and editing it in one MCP request invalidated the operator
  context;
- stale selection required an explicit BMesh deselection repair;
- Blender 5.1 rejected a 5.2-only render enum;
- an empty scene had no World datablock;
- a baseline digest initially described Python doubles rather than Blender's
  stored float32 coordinates; and
- one combined primitive transform created a roughly two-metre grid.

These corrections were causally legitimate and were not equivalent retries.
They were still setup debt rather than evidence about Reimu fidelity.

Owner change: add a small `repo-blender` compatibility preflight that proves
separate open/edit calls, empty-scene invariants, brush asset names, stored
float32 digesting, save/reopen, and a disposable render on both authoring and
pinned versions. Candidate plans should consume its receipt instead of
rediscovering these contracts.

Measure: no candidate or capability attempt should close before touching its
subject because of a preflight-covered condition.

### 5. `REIMU-FANOUT-001` — analytical delegation duplicated settled evidence

Tier: `live`.

Sole-writer isolation, implementation-blind visual review, reference
measurement, and independent pinned verification were valuable. In contrast,
attempt 017's same numerical result was separately summarized, recomputed,
and causally reviewed by nested agents after the durable writer and pinned
receipts already contained it. Six agent branches were interrupted, and the
agent scratch directory is the largest single scratch subtree at about 80 MiB.

Owner change: for each checkpoint use at most one writer, one verifier, and
one adversarial reviewer. Spawn visual reviewers only after frozen pixels
exist. Give the verifier the quantitative recomputation and let the adversary
consume that receipt rather than spawning another recomputation.

Measure: no decision fact should be independently recomputed more than once
unless the two computations use genuinely different evidence or methods.

### 6. `REIMU-CONTEXT-002` — the resume route is correct but overweight

Tier: `configured`.

The generated goal README and detailed compaction handoff preserved exact
state well. A new agent can resume without reading PR #24. However, the
minimum state plus selected procedure files was already 128 KiB before scene
inventories, attempt evidence, agent results, or image review. The goal README
leads with the entire plan history rather than the active decision.

Owner change: put a generated `Resume now` block first: protected source hash,
latest retained candidate, active attempt and owner, blocker, next causal
action, and artifact-index path. Load phase-specific skills only when that
phase begins; retain their conclusions in receipts rather than reloading
unrelated procedures after compaction.

Measure: a fresh agent should identify authority, current bytes, current gate,
and next action in at most five bounded reads and under 32 KiB.

## Necessary causal retries versus avoidable churn

Necessary and decision-relevant:

- attempt 012 to 013 replaced ignored modal key parsing with a native explicit
  transform value and proved the predicted displacement;
- attempt 015's model failure correctly changed the question from another
  geometry variant to whether richer localized sculpt controls exist;
- attempt 016 to 017 removed two named, irrelevant fixture dependencies;
- attempt 017 produced a real, isolated 4.8236% variance reduction, so attempt
  018's fixed-variable cumulative-dose test is a falsifiable final retry; and
- split open/edit calls, explicit deselection, and stored-float digesting were
  causal repairs to distinct observed mechanisms.

Avoidable or too expensive for the information gained:

- opening a durable attempt before the empty-scene fixture reached its first
  mesh;
- including renderer and World setup in a live sculpt-only fixture;
- discovering known Blender-version and storage-precision differences inside
  candidate work;
- three bow cycles that retained no module before returning to the dominant
  head defect; and
- multiple nested agents restating the same attempt-017 metrics and decision.

## What should happen next

1. Freeze or stop attempt 018 exactly at its preregistered terminal condition.
2. Publish a single current-artifact index and immediately surface the
   protected A157 file plus a clearly labelled current or rejected render.
3. If 018 fails, mark execution blocked and stop autonomous modeling. If it
   passes, allow exactly one richer-sculpt whole-head attempt; a repeated
   helmet, card, pillow, or floating-root failure then blocks the route.
4. Commit and push the rule clarification, process contract, goal projection,
   and completed attempt records together, while leaving failed binary scratch
   ignored.
5. Route the larger proposals above through normal review; this report does
   not itself authorize edits to goal schemas, shared skills, or runtime
   configuration.

The revised no-loop wording is a genuine improvement: equivalence is now
defined by decision question, target, relevant inputs, and causal execution
path, and a retry must predict an observable difference. That clarification
should reduce literal misreadings without forbidding repairs such as attempt
013 or the bounded dose-response test.
