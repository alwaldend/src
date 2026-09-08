# Evidence manifest

[Back to durable goal](./)

## Attempt 8 review-correction evidence

- Recursive JavaScript fallback fails closed when ripgrep is unavailable;
  explicitly selected files remain supported.
- Fallback fixed and regular-expression submatches use UTF-8 byte offsets.
- Bounded textual HTTP previews omit an incomplete UTF-8 suffix.
- Focused tests pass 2/2 in invocation
  `b73ce3b4-bfb6-4081-b84f-29c3f763b3a4`; the complete MCP package passes 7/7
  tests in `93845f08-cd10-4dfb-88c5-034497dea58a` and builds all 16 targets in
  `d7716a48-9509-44fe-9e37-b5c44904fffd`.
- Buildifier passes in invocation `1fb561cc-d336-49db-813f-de26b7fedbe4`.
- Follow-up review corrections preserve original-line indexes for Unicode
  case-insensitive matching, make explicit fallback files override globs, and
  return healthy scopes plus structured errors from unfiltered listing after
  partial startup.
- Follow-up focused tests pass 2/2 in
  `c07704be-cdf8-4045-a0c0-4626ebc0d1e7`; complete MCP tests pass 7/7 in
  `3fbaa4e1-8ea1-4b3f-a533-51e1ad01f87b`, all 16 targets build in
  `f87e45b6-765e-49c8-a618-cd2eb1efa1ad`, and Buildifier passes in
  `8542838d-92ae-4999-b1e8-fac631629f6f`.
- Final review corrections skip every unavailable Include scope and report only
  retained line endpoints. Fresh correctness scrutiny also preserves startup
  errors across repeated initialization and limits UTF-8 suffix removal to
  locally truncated HTTP bodies.
- Focused tests pass 2/2 in `99866999-c31d-457f-b64b-c7a15073e7a2`.
  Combined MCP and `repo-delivery` tests pass 8/8 in
  `e86a5d1c-cadf-4ff6-9647-3e1050e461e6`; all 19 affected targets build in
  `083fbf55-851e-4db4-b8e5-e5871c874faa`. Skill quick validation and
  Buildifier pass.
- Focused HTTP evidence passes in
  `e80439df-bddd-4a76-b9e1-be6c7f1ed649`.
- Runtime evidence uses semantic admission assertions and eagerly attaches the
  expected timeout rejection; three consecutive runs pass in
  `e7d69598-b843-4da4-833b-e024d406b8ca`.
- HMR rollback accepts an already-restored prior cache marker after failed
  import; otherwise it still waits for an exact correlated reload. Three
  runtime runs pass in `5f9b49e0-d3d3-4f9e-84ec-602dfbe38c77`.
- Regex fallback uses Unicode scalar mode, so `.` reports one four-byte match
  for `😀`, and the real MCP tool catalog marks `cordis_define` and
  `cordis_promote` as potentially destructive. Both focused regressions pass
  in `132e2f33-837c-40b2-83b5-4b06ceadfd0f`.
- The complete affected packet passes 8/8 tests in
  `b19d303e-0aca-4097-a191-e81015dc2982`, builds all 19 targets in
  `1c8aa274-dd72-4e6f-9491-de5e283b2c5c`, and passes Buildifier in
  `34f6f5c7-7f44-478b-b939-63ac03c3bbb1`.
- Fallback context is emitted once in line order and never reclassifies a
  matching line as context. Bounded reads distinguish a retained empty line
  from EOF. Focused invocation `2528b6ee-7d82-48c9-a607-298a4cef0b9b`
  passes both regressions.

## Attempt 7 working-tree evidence

- The custom `manifest.json`, hash-named `versions/`, storage layer, activation
  worker, and package worker have been removed.
- Official `@deepseek-ai/cordis-plugin-loader`, `-include`, `-hmr`, and
  `-timer` packages are pinned. Bazel launches Node with the HMR package's
  documented `--expose-internals` requirement.
- Project entries use `projects/mcp_cordis/cordis.yaml` and ordinary modules
  under `plugins/`; scratch entries use the same layout under
  `out/mcp_cordis`.
- Runtime integration proves create, live HMR update, failed activation with
  on-disk/live rollback, stop, run, promotion, removal, and restart recovery.
- Direct starter-module tests and the complete starter runtime test pass. A
  real stdio client proves the complete lifecycle, failed-update rollback,
  package-log isolation, and restart recovery without replacing the MCP
  process during an update.
- Runtime regressions prove exact source-limit round-trip, invalid-timeout
  side-effect exclusion, never-settling activation rejection, handler-lease
  draining, and direct/descendant process-group non-liveness at settlement.
- Affected test and Buildifier invocation
  `39711671-375a-413f-8a72-e6f9ff892bd3` passed 11/11 on the rebased
  implementation after the final import-boundary corrections. Affected build
  invocation `153467fd-e62d-4eab-b260-6754d17fe8e2` passed all 26 targets.
- These full receipts bind implementation commit `0a93e487`. Later amendments
  are confined to the durable goal record and require proportional diff and
  formatting validation before publication.
- Fresh independent review accepted durable-record commit `c05bd45a` with no
  actionable findings.
- PR 32 was republished at the verified rebased head. Its description now
  records the standard Cordis architecture, and all three obsolete review
  threads are resolved.

## Preflight evidence

- Official OpenAI documentation confirms local Codex clients can connect
  directly to stdio MCP servers and read server instructions.
- Official DeepSeek documentation states that its dynamic Cordis definitions
  are process-local and memory-only, establishing the need for the requested
  persistence layer.
- Cordis `4.0.1` exposes the required `Context`, `plugin`, `Fiber.await`,
  `Fiber.dispose`, and effect-scoped cleanup primitives without depending on
  DeepSeek Harness's agent/session/browser packages.
- MCP SDK v2 provides a stable stdio server and fixed tool registration. Codex
  does not reliably refresh dynamically added tool schemas, so the accepted
  design keeps a fixed list/invoke gateway.
- Repository review selected an ordinary root-workspace Bazel package with a
  project-owned pnpm lock and Bzlmod dependency fragment.
- Safe aggregate analysis of 40 top-level recent sessions selected
  `repo_context`, `git_worktree`, and `network_probe` as the initial reusable
  packages. No transcript or secret-bearing content will be copied.

## Attempt 1 verdict

- Architecture: one worker and Cordis root/fiber per active package generation.
- Storage: immutable content-addressed source and atomic manifests in explicit
  `project` or `scratch` scopes.
- Update rule: start and validate a candidate, atomically swap the active
  generation, then drain and dispose the previous generation.
- MCP rule: stdout is protocol-only; package output is redirected to stderr and
  `out/mcp_cordis/logs`.
- Candidate hash:
  `c9300c9887104777c8915e3d4f390196604e9bd18497bbec319415d1a4ad057f`.
- Focused query, lifecycle/in-memory MCP, and subprocess stdio tests pass.
- Starter execution test fails at `repo_context_search` with
  `spawn rg ENOENT`; candidate rejected and Attempt 2 opened.

## Attempt 2 verdict

- Candidate commit:
  `e3e74cb1e573867825347292bf17220a5b9a4a0c`.
- Candidate tree: `079a0c27b86527c6950cc75b0c8b9dbf572d3e4b`.
- Fetched base and direct parent:
  `7ad2704cd27757355ab36ec8eb1bb27ef9e1d91d`.
- The delivery preparation receipt confirms an exact, conflict-free rebase,
  task-only path scope, and absence of a remote feature ref or pull request.
- Focused post-rebase query passes.
- Integrated post-rebase result: three of four tests pass. Starter packages,
  stdio, and build coverage pass; lifecycle evidence fails at an elapsed-time
  drain precondition with actual `0`, expected `1`.
- Verdict: refine the test evidence in Attempt 3; do not change the runtime
  architecture based on this measurement.

## Attempt 3 verdict

- Final local commit:
  `7cfef0719075ad372c3bb257ad216b35770356b2`.
- Final local tree: `34153eca0f582af5c641f81bf8c7209b0045ab9a`.
- Direct parent and fetched base:
  `7ad2704cd27757355ab36ec8eb1bb27ef9e1d91d`.
- Forced project tests: four of four pass on the exact commit.
- Complete project build: nine of nine targets pass.
- Forced root Buildifier and exact commit diff check: pass.
- Review verdict: accept as the final local candidate. Publication was stopped
  before execution because a rebase-only request does not authorize remote
  push/PR mutation.

## Attempt 4 preflight

- PR 32 was published at exact commit `7cfef071`; remote branch and PR tree
  matched the local receipt.
- Current remote `master` is `7ad2704c`, the direct parent of the published
  task commit; repository inspection reports `needs_rebase: false`.
- PR 24 head is `da2085f1`, based on `ada3ed90`. Its `projects/agents` diff is
  exactly four files: one `bazel-agent` update, one `goal` update, and two new
  `decision-review` files.
- Three PR 32 review threads were independently diagnosed as valid. Their
  controlling fixes are bounded-success execution, pre-record change limits,
  and shutdown admission closure followed by lock draining.
- Verdict: reject `7cfef071` as final; proceed with the scoped three-way import
  and three behavior-changing corrections in Attempt 4.

## Attempt 4 working-tree evidence

- PR 24 import: `bazel-agent` batching guidance, result-first `goal` guidance,
  and the exact `decision-review` instruction blob are present. Newer
  throwaway-record and bounded-delegation guidance remains intact.
- Execution overflow now retains a combined arrival-order byte prefix, returns
  only valid UTF-8, stops live process-group members, and reports `truncated`
  rather than rejecting. Timeout and spawn errors retain rejection semantics.
- `git_worktree` keeps `70d8...` as immutable history and activates new exact
  content hash `de978...`; its record limit is checked before every porcelain
  record matcher.
- Shutdown closes admission, joins the exact admitted package-lock snapshot,
  then disposes the final active set and awaits retirements through one
  memoized promise. Initialization rechecks closure after awaited storage
  boundaries.
- Four focused Bazel tests pass, including all new regression targets and
  `decision-review` offline validation.
- Integrated Bazel test invocation passes ten of ten tests: the whole MCP
  package plus all three imported/updated skill configurations.
- Full affected build passes and validates `bazel-agent`, `goal`, and
  `decision-review`; root Buildifier passes.
- Verdict: behaviorally acceptable as a working-tree candidate. Independent
  diff review and exact commit-bound validation are still required.

## Attempt 4 independent-review verdict

- Max-change review proved the bound placement but rejected the completeness
  flag when exactly `maximum` records exist, the missing second-record test,
  and newline-unsafe pathname capture.
- Execution review accepted byte bounding, process-group cleanup, and runfiles,
  but
  rejected silent invalid-UTF-8 loss and incomplete propagation of host
  truncation through reusable package result fields.
- PR 24 provenance review accepted the three-way import and packaging, but
  found no eval cases for compatible Bazel batching, immutable candidate
  promotion, or durable-versus-throwaway push behavior.
- Verdict: refine; passing Attempt 4 checks do not qualify it for commit.

## Attempt 5 working-tree evidence

- Execution tests cover ASCII overflow, multi-byte boundary overflow,
  malformed UTF-8 both below and above the byte cap, combined stdout/stderr
  budgeting, timeout rejection, normal nonzero exit, and process-group
  non-liveness at settlement.
- `git_worktree` active hash is
  `8853aa20665778aeec43e03f2fe975445002d56ed33ea1cf38bef3946381f60d`;
  its manifest retains only that version and unchanged historical hash
  `70d8f28dc947d19410b8e79bad90cb416303107243d72d062dd70e30f97a2c3b`.
- `repo_context` active hash is
  `abd0db3e26ec970dcf5cc3ec21b9f2b2c452f9302edc0e10c5231a140b92fbc0`;
  all three manifest version hashes match their exact source bytes.
- Focused MCP regressions and three skill eval configurations pass 8/8.
- Full affected test packet passes 12/12; full affected build and three
  `rules_skill` validation aspects pass; root Buildifier passes 1/1.
- Exact commit and post-rebase evidence remain pending, so this is a green
  working-tree candidate rather than a delivered checkpoint.

## Attempt 6 working-tree evidence

- Current immutable starter hashes are
  `04b06a7d6277c4a6e8513d970f549ad980a780b68755f28d7b402fe8be26c279`
  for `git_worktree` and
  `94131e058f82328f091613dc68d2717484378066a9c64940d99522c14b48b4d7`
  for `repo_context`; line wrapping does not introduce whitespace into the
  first identifier.
- Historical version bytes retain their exact filename hashes. The checkout
  enforces LF for every hash-addressed JavaScript source.
- The parent activation now owns every `ctx.exec()` child handle. Focused tests
  prove immediate live-process absence after inner timeout, outer timeout,
  startup timeout, output overflow, normal completion with a background child,
  and a background child inheriting output pipes.
- Runtime regressions prove shutdown admits complete promotion, rolls back a
  candidate that fails during active-version persistence, and actively disposes
  a retired generation whose admitted handler remains gated.
- Bazel invocation `9b1b2d34-8490-474a-b12c-e2052bf2d90b` passes the current
  process, unavailable-shutdown, and runtime-admission targets 3/3.
- Bazel invocation `7c680be8-e2a6-4b5b-9b45-4a0390f39a5a` passes the current
  process, ripgrep-byte-field, and runtime-admission targets 3/3.
- Fresh whole-diff and adversarial supervisor reviews are running. Full
  integrated, build, skill-aspect, Buildifier, exact-commit, rebase, and remote
  verification remain unverified after the latest changes.

## Attempt 10 focused HMR evidence

- Subject: dirty Attempt 10 tree using
  `@deepseek-ai/cordis-plugin-hmr` 1.0.16 with pnpm patch hash
  `ec800d86298faacc86c7717ffa1dce7c28116ab1393b8abc198be6ac02c38489`.
- Primary dependency evidence: the npm registry lists 1.0.16 as the latest
  release, and current upstream `vendor/hmr/src/index.ts` retains the same
  untracked debounced `partialReload()` plus shared-stash reset behavior.
- Independent reproduction: a slow top-level-await generation followed by a
  second write reached `{live: "slow", diskLatest: true}` with the unpatched
  package.
- Patch applicability: `git apply --check` accepts
  `patches/hmr@1.0.16.patch` against the exact resolved 1.0.16 package files.
- Lock generation: Bazel-managed pnpm invocation
  `a1c50484-20c6-4935-8832-92029d0de3c6` completed successfully with no
  unrelated dependency resolution changes.
- Causal runtime regressions: Bazel invocation
  `ac79a8e4-f5d9-4dd7-a821-29bce3d8ece6` passed explicit release-gated
  top-level evaluation and asynchronous activation overlaps, failed apply
  rollback and recovery, manual editing, and deterministic disabled activation.
- Repetition: Bazel invocation `65c1932a-896b-470d-9462-086dd93beaff`
  passed ten runs each of `runtime_test` and `starter_packages_test`.
- Complete project packet: invocations
  `6490dd57-49e0-4558-b280-f5625db07208`,
  `218e5797-b2d7-44bc-aded-f5df8139ca1c`, and
  `09d4ec2c-aac1-4509-8579-5ef8c5eebe39` passed all project tests, all project
  builds, and root Buildifier respectively.
- Project-layout adaptation: Bazel invocation
  `ea2c6c6e-b2e1-425d-bedb-2ac9679de6c5` passed `runtime_test` after moving
  the command, internal implementation, launcher, and test suite to their
  role-based directories.
- Independent final HMR review accepted patch identity, stashed-change
  draining, complete Fiber cleanup join, declarations, causal overlap tests,
  rollback and recovery, and disabled-entry activation. Its only finding was
  a README sentence describing the superseded publication order; that contract
  text is corrected in the current tree.
- The real `cmd/mcp_cordis/launch.sh` completed MCP initialization and listed
  all ten tools. While that stdio server remained live, Bazel query invocation
  `7739c14e-1b22-40c7-94af-b81143e84d4a` completed successfully, proving the
  launcher releases the workspace Bazel lock before serving; SIGINT shut the
  server down cleanly.
- Preliminary current-tree Bazel invocations
  `95d37231-63f9-44f6-9eee-3f0fe8fb4107`,
  `d70f61a2-57e5-4ce9-b8fd-27b06bd02a6d`, and
  `1e186c82-3bcd-4048-a534-fe747e1cf79c` pass the complete affected test
  packet (10/10), affected build packet, and root Buildifier check.
- Guarded delivery inspection bound same-repository PR 32 to local and remote
  feature OID `bc4e5ae97ef9ea968c01b1b2a55403ae032a6a8d`, fetched base OID
  `d29f9d471ea467e8dfc75db4eedeedbbae43dc2d`, and SSH transport. All nine
  linear feature commits have the task-bot author and committer identity; the
  only refusal is version 1's unconditional multi-commit consolidation guard,
  which is unchanged on the fetched base.
- The user explicitly authorized the narrow adapter extension. Its exact-head
  consolidation path verifies linearity, identity, oldest ownership marker,
  pull-request metadata matching the requested aggregate projection,
  signature requirements, and every other inspection refusal before creating
  one aggregate commit. Integration
  coverage proves that the pre-consolidation remote head remains the
  receipt-bound publication lease.
- Bazel invocations `574dadd5-51b1-443e-b8c2-50ca5d257eb2` and
  `030a8bb4-2a22-4ece-b632-b3c75572bcee` pass the complete adapter suite once,
  then its Go, skill-validation, and root Buildifier targets three times.
- Independent adapter review rejected the first implementation because clean
  ranges had no staged delta. The correction permits an unchanged index only
  behind validated exact consolidation evidence, while parent-to-tree scope
  validation continues to reject an empty aggregate. A clean `--path`
  regression proves tree preservation and one final commit; Bazel invocation
  `0d715375-6a8d-4269-ad3e-f8a002888808` passes, and independent re-review
  accepts the corrected adapter with no remaining findings.
- After adding upstream-compatible root-consumer visibility to the
  branch-owned `decision-review` skill, Bazel invocations
  `6099d90d-0ade-43b2-b50b-8f7050c26c32` and
  `7f4b37a4-b970-46ea-bc76-b4a60aeeab59` pass its focused validation and the
  root Buildifier check; `git diff --check` also passes.
- Guarded consolidation and rebase produced one candidate commit on fetched
  base `63e7b9f0be1e054373415914ff3d2ea2282aa3da`. The adapter now stages
  deletions and symlink-to-directory changes and permits a rebased path to
  vanish only when the old candidate and new base have the exact same Git tree
  entry.
- Exact code-candidate Bazel invocations
  `13588b19-1f8a-43a7-8006-f9d5d4652670` and
  `5863ac24-3bb2-47be-99fb-50141e933018` passed all 12 affected tests and all
  29 affected builds, including Buildifier and discovery-link validation.
  Live launcher initialization listed all ten gateway tools while concurrent
  query invocation `4637f187-d05e-4905-8153-18fca8644ea1` passed.
- Verdict: the focused race, recovery, layout, repeated-run, project-wide,
  consolidation, and exact code-candidate evidence pass. The final
  adapter-and-record rewrite, remote publication, hosted review
  reconciliation, and final
  receipt verification remain open.

## Attempt 11 final review evidence

- Guarded publication and receipt verification established a clean,
  single-commit feature branch on base
  `63e7b9f0be1e054373415914ff3d2ea2282aa3da`, with the local and remote tree
  identical and PR 32 synchronized.
- The exact-head hosted review completed against the published aggregate and
  reported three actionable findings: lexical symlink reopening, unbounded
  `/proc` inspection retries, and a 100 ms process-start assumption.
- Repository reads now open the canonical target with `O_NOFOLLOW`, validate
  the opened descriptor through `/proc/self/fd`, inspect and read through that
  same handle, and retain the existing byte bound. A regression proves an
  internal symlink read never reopens the lexical alias through `ctx.readText`.
- Process-group verification now skips inaccessible per-PID entries and turns
  three consecutive inspection failures into `EXEC_CLEANUP`; a direct
  regression proves the retry bound. The timeout cleanup fixture now allows a
  five-second startup window before exercising forced group cleanup.
- Bazel invocation `d32cccf7-b5b9-427c-93a1-6b612e32a0aa` passed all seven MCP
  tests. Invocation `9d11ce11-1f79-4077-b61a-efa95a5fc3de` built all 16 MCP
  targets. Invocation `aadcc7f6-e2b3-4480-8b76-0607b2017415` passed three runs
  each of the process-supervisor and repository-context regression targets.
- Final publication, exact-head review completion, thread reconciliation, and
  receipt verification are performed by the guarded delivery workflow; the
  ignored receipt is the authoritative mutable delivery record.
- Follow-up exact-head review identified lexical reopening in `git_worktree`.
  Repository selection and discovered-root use now retain verified directory
  handles for the complete tool call, and all Git `-C` arguments address those
  handles through `/proc/<pid>/fd`. The focused command-contract test and real
  Cordis starter-package integration pass together in Bazel invocation
  `4c3eba8d-32c9-47e9-a257-f8af92ef0c19`; invocation
  `6c638120-90cb-446c-b76e-6d39df3640e5` repeats both targets three times.
  Invocations `ef83ec38-84ed-466e-abbd-ffd6cef892c1` and
  `16e0c322-adc7-4763-8ff4-11d5f38a4172` pass the complete seven-test and
  sixteen-target MCP packets.
- The next exact-head pass identified the same lexical reopening in
  `repo_context`'s ripgrep and Git metadata branches. Selected search paths and
  repository directories now remain open while subprocesses address them
  through `/proc/<pid>/fd`; reported ripgrep paths are mapped back to stable
  workspace-relative names. Directory kind and entry inspection also use the
  selected handle. Bazel invocation
  `a3f4ed56-a44b-46d6-8717-c38ecc7f05eb` passes the focused context contract
  and real ripgrep/Git starter integration together. Invocation
  `e7af6e8f-823e-40ef-a833-dcb0a704f46f` repeats both three times, while
  invocations `9fd4d465-957a-424f-9bb8-f6e10eb93009` and
  `4a464ddd-f646-4b9a-aaab-a411bf930ee0` pass the complete seven-test and
  sixteen-target MCP packets.
- The terminal exact-head review found that the JavaScript regex fallback
  could monopolize the MCP event loop through backtracking and could not match
  ripgrep's Unicode regex semantics. Regex search now requires ripgrep when
  the executable is unavailable, while the bounded fixed-string fallback is
  preserved. Bazel invocation `b98534c0-1667-462b-81d5-ee393fca343b`
  passes the focused context and real starter integration targets;
  invocations `4689e7cb-9600-4b35-ad7a-97e02cd6730c` and
  `557f322c-4fad-4358-bf5f-98ad9b8972b0` pass the complete seven-test and
  sixteen-target MCP packets.
- The next exact-head pass found that invalid UTF-8 replacement decoding
  changed fixed-search byte offsets when ripgrep was unavailable. The fallback
  now verifies that decoded text round-trips to the original bytes and fails
  closed otherwise. Bazel invocation
  `449ea2c4-4d52-44d2-b7d1-b1e9745359bd` passes the focused context and real
  starter integration targets; invocations
  `77bae252-af20-4242-89d8-7d6bc242b64f` and
  `3a0bd136-a327-416a-8128-aaaec171bd9c` pass the complete seven-test and
  sixteen-target MCP packets.
