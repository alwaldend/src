# Upstream impact review for the Fumo rebase

## Observation boundary

- Compared old base
  `21c299a7193e1a4fe4467d8e0e7d713650958399` with fetched base
  `21edc3f227bfa7355848f184adf65129b53947f0`.
- Observed the rebased feature worktree at
  `054028a1ac7325cefc60f4d896d5ce1077dda69e` on
  `t3code/replace_man_with_fumo` at `2026-09-01T11:25:01+03:00`.
- Sources were bounded Git comparisons and bounded reads through the live
  project Cordis handlers. No tracked file, Git ref, or remote state was
  mutated by this review.
- This is an integration-impact review, not a build, goal-store validation, or
  remote-delivery receipt. The observed worktree was clean, but its remote
  tracking relationship was still `ahead 18, behind 8` and therefore was not
  proof of a completed remote checkpoint.

## Executive finding

The fetched base introduces several workflow contracts that materially affect
the Fumo goal. They should be preserved during conflict resolution rather than
replaced by older local wording:

1. long-running goals must prioritize the requested result, survive ordinary
   interruptions, use bounded modules, review the whole process on a stalled
   trend or every three attempts, and remotely checkpoint durable progress;
2. independent workstreams should be actively delegated, while canonical goal
   state and monolithic Blender assets retain one writer;
3. repository Blender work now has a canonical execution skill that requires
   the pinned toolchain, isolated candidates, immutable render snapshots, and
   clean-reopen verification; and
4. repository inspection should prefer the available bounded Cordis handlers,
   while every modifying task remains in a dedicated feature worktree.

The observed rebased tree retains these upstream contracts. Its Fumo-specific
additions are additive rather than replacements: the current diffs add
delta-first iteration, causal-reach checks, staged feedback, direct authoring
guidance, and artifact-link requirements without deleting the fetched-base
text. No conflict markers were found in the reviewed policy, goal, Blender,
reference-fidelity, or Fumo-goal paths.

## Upstream changes and required integration

| Area | Upstream change | Required Fumo integration |
|---|---|---|
| `AGENTS.md` | Adds proactive bounded delegation, one coordinator, mandatory dedicated feature worktrees, bounded tooling-bug policy, and preference for live Cordis handlers. | Keep all clauses. Enumerate ready workstreams on resume and strategy changes. Delegate only decision-changing variants, review, measurement, or artifact work. Keep the default checkout untouched and use Cordis for bounded repository reads, searches, and Git inspection when available. |
| Goal skill | Makes the requested result the priority; interruptions do not terminate a goal; requires small high-leverage modules, early falsification, remote durable checkpoints, a full process review at least every three attempts or immediately on stall, and explicit parallelism checks. | The requested review must cover the history through A89 and name that cutoff. The next attempt must target one acceptance-visible defect, retain whole-head context, have an early-stop artifact, and avoid turning record maintenance into the critical path. Durable records and reusable sources must be committed and pushed through repository delivery before yielding. |
| Goal concurrency reference | Requires isolated copies of one frozen baseline for a monolithic artifact rather than splitting the canonical file. Workers cannot publish canonical goal state. | Continue using hash-bound `.blend` copies under `out/`, one writer per copy. The coordinator alone checkpoints the project goal and promotes a proven candidate. Do not split `reimu_fumo.blend` merely to create parallel edit lanes. |
| Goal lifecycle/evidence reference | Adds module ownership, minimum whole-result context, early termination on interface failure, and a mandatory end-to-end process audit cadence. | A rear-hair probe must show that its bounded operation can materially alter the failed rear/profile evidence before authoring. Review rear, both profiles, and enough front/three-quarter context to expose fins, capes, or lost crown handoff. Stop at the first categorical recurrence. |
| Goal CLI/store | Repeated `--criterion` flags now preserve literal values instead of comma-splitting. Imported plan, result, review, criteria, and evidence paths must resolve inside the workspace. A missing `attempts/` directory is now valid until the first attempt and is created safely. | Existing Fumo records appear structurally compatible: their imported evidence lives inside the workspace and `attempts/` already exists. Future checkpoint inputs must stay under this worktree, normally `out/<task>/`. Do not rely on comma-separated criteria flags. Run the goal validator after the rebase before writing the next canonical attempt. |
| Operation contracts | Adds owner-local declarations for goal checkpoint/promote/migrate and repository-delivery prepare/publish/review mutations. Phase 1 remains report-only; declarations do not create authority. | Treat `goal checkpoint` as source/task-state mutation guarded by exact resource version. Treat delivery prepare as history/source mutation and publish/review as remote mutation. Do not infer permission from the declaration itself. |
| New `repo-blender` skill | Establishes the pinned Bazel Blender as the batch/deliverable toolchain, background execution as default, explicit task scratch, protected-source hashing, candidate copies, immutable checkpoint rendering, final-pixel inspection, and clean-reopen verification. | Use it together with `blender-reference-fidelity`. Never render the file while the foreground editor may be saving it. Freeze and hash a snapshot, then render it with pinned background Blender. The Flatpak may remain the explicitly selected live MCP authoring host, but it is not the renderer or verifier of record. |
| Blender execution reference | Requires explicit GUI/MCP capability checks, one foreground host, `--factory-startup`, `--disable-autoexec`, `--python-exit-code`, batched related operations, and task-contained configuration and temporary paths. | Do not treat a missing listener as a terminal blocker. Launch one authorized host only when live UI is required; otherwise use background Blender. Batch the fixed review packet in one load where practical and record why any safety flag must be omitted. |
| Blender toolchain manifest | Keeps Blender `5.2.1` but adds the complete extracted distribution as runtime files. | This is not a version change. Use the current Bazel wrapper so the full runtime tree is present; do not call the extracted ELF directly or substitute PATH/Flatpak Blender for reproducible output. |
| `bazel_agent` | Stops injecting a shared workspace `out/tmp` through `TMPDIR`, `TMP`, and `TEMP`; adds a bounded `doctor` command. Bazel actions/tests use Bazel-owned temporary storage. | Fumo host-side scripts and foreground Blender must receive explicit task/run scratch paths. Do not assume `bazel_agent` supplies `TMPDIR`. A task-local doctor report may verify scratch classification when diagnosing the environment, but it is not required for every render. |
| Cordis runtime | Moves scratch to `out/<task>/mcp_cordis/runs/<run>/` with portable task/run/worker identities and a bounded manifest. | Any new disposable Cordis helper for the Fumo process must use that namespaced layout. Existing live handlers should be reused before creating another helper. |
| Agent-system Phase 1 | Adds shared effect, authority, information, retention, and operation vocabularies plus a report-only registered-universe check. | This is descriptive architecture, not a new acceptance gate for the plush. Do not derail the Fumo goal to adopt all Phase 1 machinery. Use the owner-local goal and operation contracts already relevant to the workflow. |

## Delivery and rebase findings

There is no direct source change between the two bases in either
`git-rebase-remote/SKILL.md` or `repo-delivery/SKILL.md`. The established
procedures therefore remain authoritative:

- establish exact base and feature refs and ownership before rewriting;
- preserve unrelated or human-owned work;
- revalidate changes invalidated by conflict resolution or rewrite;
- use the supported repository-delivery adapter for GitHub operations rather
  than reproducing them with raw Git or `gh`;
- publish only an exact validated candidate with its receipt and exact lease;
  and
- fetch and verify the final base and feature identities after publication.

The new `tools/repo_delivery/agent_operations.json` describes prepare,
publish, and review effects but does not relax or replace the delivery skill.
The goal skill's new remote-checkpoint rule makes successful delivery part of
an advancing goal turn; if transport, conflicts, or a changed lease prevent
it, the goal record and user report must state that exact blocker instead of
claiming remote preservation.

## Current Fumo-record compatibility snapshot

- The project goal remains open but explicitly paused at resource version
  `54`, lifecycle generation `2`, with no active attempt.
- A89 is closed as `refine`; no A89 lower-rear-hair candidate was promoted.
- A87 front hair and A88 crown are provisional modules only. No full-goal
  criterion passes.
- The existing goal directory already follows the versioned v1alpha1 layout,
  has an `attempts/` directory, uses durable retention, and records exact
  attempt/evidence digests.
- Nothing in the upstream store changes requires a record migration based on
  structural inspection. This must still be confirmed by the current goal
  validator before the next checkpoint.
- Canonical goal files and closed attempt directories must not be edited by
  hand. Import the requested process review and later plans through the goal
  tool using the latest compare-and-swap resource version.

## Conflict-resolution acceptance checklist

The integrated candidate should not be accepted until all of the following are
true:

- `AGENTS.md` retains the fetched-base parallel-work, dedicated-worktree,
  tooling-bug, and Cordis-handler clauses.
- The goal skill retains result priority, interruption continuity, bounded
  modules, remote checkpoints, process-review cadence, proactive delegation,
  and isolated-copy guidance.
- Local goal-skill additions remain additive: evidenced deltas, causal reach,
  staged feedback, time-to-artifact control, and standing adversarial review
  must not contradict the upstream rules or make process work outrank results.
- The full upstream `repo-blender` safety/toolchain contract remains present;
  local direct-authoring and failed-listener guidance may extend it but not
  turn Flatpak/MCP into the renderer of record.
- The `5.2.1` Blender runtime-files addition and the removal of ambient Bazel
  temporary-variable injection remain intact.
- The existing Reimu references, tracked blend, LFS declarations, goal
  resource versions, and closed-attempt digests remain unchanged by conflict
  resolution unless an exact authorized tool operation intentionally changes
  them.
- No conflict marker or generated-file hand edit remains.
- Goal validation, task-scoped repository checks, and the exact-head delivery
  gate are rerun after the final resolved candidate is prepared.

## Implication for the next Fumo plan

The upstream changes support a narrower and faster next iteration, but they do
not themselves solve the artistic failure. The next plan should:

1. preserve A88 S07 as the immutable baseline and treat A89 S00-S06 only as
   negative evidence;
2. name the stable failure as lower-rear hair collapsing into a rigid cape or
   side fin while exposing the rear receiver;
3. establish causal reach for one direct, locally sculpted lock or small lock
   cluster before building a complete rear arrangement;
4. render a cheap frozen rear/profile coupon with component IDs immediately,
   with an implementation-blind reviewer already waiting;
5. reject on the first cape, fin, exposed-band, or crown-overlap recurrence;
6. expand only a passing coupon into staggered independent locks, then rerun
   front and three-quarter regressions against the exact same snapshot;
7. keep render/review infrastructure invariant and change only the geometry
   under test; and
8. record the whole-goal/process review through A89, then remotely checkpoint
   the durable review and plan through the exact repository-delivery workflow.

This plan incorporates the fetched-base contracts while keeping the requested
plush result, rather than process expansion, on the critical path.
