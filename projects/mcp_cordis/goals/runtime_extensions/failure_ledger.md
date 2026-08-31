# Failure ledger

[Back to durable goal](./)

## Attempt 8 hosted review finds starter fallback and encoding gaps

- Candidate: published standard-Cordis commit `deaa2dcd`.
- Result: hosted review found three valid defects after the earlier local
  whole-diff review had accepted the runtime architecture.
- Causes: the no-ripgrep fallback recursively traversed ignored and hidden
  files, JavaScript string indexes were exposed as byte offsets, and a bounded
  textual HTTP body decoded an incomplete trailing UTF-8 sequence.
- Strategy delta: fail closed for directory fallback, retain explicit-file
  fallback with true UTF-8 offsets, and decode only a complete UTF-8 prefix.
- Regression guard: explicit hidden-file versus directory cases, fixed and
  regex non-ASCII offsets, and a one-byte multibyte HTTP preview must pass with
  the complete MCP test/build packet.
- Follow-up causes: lowercasing a line could change UTF-16 length before offset
  projection, fallback glob filtering contradicted ripgrep's explicit-path
  precedence, and unfiltered listing did not preserve a healthy scope after
  partial startup.
- Follow-up guard: length-changing `İx`, explicit file plus excluding glob, and
  malformed-project/healthy-scratch listing cases pass with the complete MCP
  packet.
- Final-review causes: config parsing alone did not identify a scope whose
  Include never mounted, and bounded reads reused the requested endpoint after
  clipping content.
- Fresh-scrutiny cause: idempotent initialization returned an empty error list
  after partial startup because only the first call retained local failures.
- Additional fresh-scrutiny cause: UTF-8 suffix trimming did not distinguish a
  byte-cap split from a naturally completed malformed response body.
- Final guard: missing-module and malformed-config scopes are both skipped by
  unfiltered listing, clipped reads report retained endpoints, and repeated
  initialization preserves structured scope errors. Textual bodies trim an
  incomplete suffix only when locally truncated. `repo-delivery` now makes this
  correctness scrutiny an explicit post-edit gate.

## Exact aggregate runtime evidence was timing-sensitive

- Candidate: post-review Attempt 8 aggregate validation.
- Result: the complete packet failed only when system load delayed a
  wall-clock `<400 ms` assertion and allowed an expected 300 ms rejection to
  occur before its assertion was attached.
- Cause: the test inferred synchronous admission from elapsed time and created
  a temporary unhandled-rejection window while waiting for a PID file.
- Strategy delta: assert the wrapped synchronous-admission error semantics and
  attach the expected rejection before awaiting the fixture handshake.
- Regression guard: the complete runtime target passes three consecutive
  concurrent runs without loosening production deadlines.

## Failed-import rollback waited for an event Cordis does not emit

- Candidate: loaded exact aggregate validation after the evidence correction.
- Result: a failed source update occasionally reported
  `reload_rollback_failed` even though Cordis had already restored the prior
  module cache.
- Cause: failed HMR import restores its cache and returns without emitting
  `hmr/reload`; the host always required a second reload event after restoring
  prior bytes.
- Strategy delta: after restoring bytes, accept the exact prior cache marker
  immediately; retain correlated HMR waiting when the failed candidate marker
  actually reached the cache.
- Regression guard: three concurrent runtime runs prove failed-update rollback,
  source restoration, and subsequent healthy invocation.

## Attempt 1, check 1: npm repository analysis

- Command: `bazel_agent query //projects/mcp_cordis:all`
- Result: failed before target analysis.
- Cause: rules_js requires pnpm v10 workspaces to declare
  `onlyBuiltDependencies`, including when lifecycle actions are disabled.
- Evidence: repository fetch failed in
  `verify_lifecycle_hooks_specified` with that exact requirement.
- Strategy delta: declare an empty lifecycle allowlist in the project-owned
  workspace and regenerate its exact lock before rerunning the same query.
- Regression guard: successful focused query and build of the translated npm
  repository.

## Attempt 1, check 2: package query after npm correction

- Command: `bazel_agent query //projects/mcp_cordis:all`
- Result: npm translation succeeded; target loading stopped because the
  deliberately non-empty starter-package glob had not yet been populated.
- Cause: implementation was queried while the parallel starter-package draft
  was still outstanding.
- Strategy delta: retain the non-empty invariant and add the three accepted
  packages before rerunning the same query.
- Regression guard: the `starter_packages` target must contain real files and
  the focused query must succeed.

## Attempt 1, check 3: runtime test process did not terminate

- Command: `bazel_agent test //projects/mcp_cordis:runtime_test`
- Result: both test bodies completed in under 400 ms, one failed, but a worker
  retained by the failing test kept the process alive until the run was
  interrupted at 144 seconds.
- Cause: the test registered cleanup only along its success path, so the first
  assertion failure leaked its runtime and obscured the underlying defect.
- Strategy delta: register unconditional `node:test` teardown before the first
  assertion, then rerun to expose the real behavioral failure promptly.
- Regression guard: the target must terminate normally on both passing and
  failing assertions.

## Attempt 1, check 4: activation error assertion mismatch

- Command: focused runtime test after unconditional teardown.
- Result: target terminated in 0.8 seconds; the MCP gateway test passed and
  the lifecycle test stopped at its syntax-rollback assertion.
- Cause: the runtime correctly exposed `activation_failed` as the error's
  machine-readable `code`, while the test searched only its human message.
- Strategy delta: assert the stable code and separately check the underlying
  syntax diagnostic.
- Regression guard: failed activation retains the working v2 generation and
  reports both the stable wrapper code and candidate cause.

## Starter search depends on undeclared ripgrep (occurrence 1)

- Command: `bazel_agent test
//projects/mcp_cordis:starter_packages_test`.
- Result: failed in 0.4 seconds with `spawn rg ENOENT`.
- Cause: `repo_context_search` invoked the preferred ripgrep engine without a
  fallback, while the Bazel test PATH intentionally does not expose the host
  installation.
- Attempted strategy: none before this measurement.
- Strategy delta: Attempt 2 adds a bounded in-process fallback while retaining
  ripgrep when available.
- Regression guard: the unchanged hermetic starter test must exercise a
  successful search and all remaining package handlers.
- Latest result: resolved in Attempt 2. The unchanged hermetic test passes with
  the new bounded JavaScript fallback and executes all eight tools.

## Lifecycle drain test uses an elapsed-time race (occurrence 1)

- Command: `bazel_agent test //projects/mcp_cordis:all` on rebased commit
  `e3e74cb1e573867825347292bf17220a5b9a4a0c`.
- Result: three test targets pass; `runtime_test` reports actual drain count
  `0`, expected `1`.
- Cause: the test delays the old request for 150 ms but must start and validate
  a replacement worker before retirement. Nothing guarantees the swap occurs
  before the fixed request delay expires.
- Attempted strategy: a 20 ms sleep before starting replacement; this proves
  that the old request started, not that it remains active at the later swap.
- Strategy delta: Attempt 3 uses an explicit started marker and release latch.
- Regression guard: run the focused lifecycle test and complete suite with no
  wall-clock assumption controlling the drain-count assertion.

## BUILD data labels not in Buildifier order (occurrence 1)

- Command: `bazel_agent test //:buildifier_test` during Attempt 3.
- Result: failed with an exact three-line ordering diff in
  `projects/mcp_cordis/BUILD.bazel`.
- Cause: the external-style `:node_modules/...` label followed the two shorter
  local labels in `runtime_test.data`.
- Strategy delta: apply Buildifier's exact lexical ordering and rerun the same
  repository formatter target.
- Regression guard: `//:buildifier_test` must pass on the frozen candidate.
- Latest result: resolved. The forced Buildifier test passes on commit
  `7cfef071`.

## Remote publication lacks explicit authorization (occurrence 1)

- Intended command: repository delivery `publish` using the exact preparation
  receipt and validated head `7cfef071`.
- Result: the approval gate rejected execution before any push or pull request
  mutation.
- Cause: the latest user instruction explicitly requested a rebase and did not
  authorize the separate consequential remote publication operation.
- Safe attempts: preparation, exact-tree validation, and local commit are
  complete; no workaround or indirect mutation was attempted.
- Exact unblocker: explicit user authorization to push this branch and create
  or update its pull request.
- Latest result: resolved. The user explicitly requested push and continued
  goal execution; PR 32 was published at `7cfef071` and remains the authorized
  delivery vehicle for the corrected candidate.

## Published candidate has three valid review defects (occurrence 1)

- Candidate: PR 32 commit `7cfef0719075ad372c3bb257ad216b35770356b2`.
- Result: automated review found three independently reproducible defects.
- Causes: output overflow rejects instead of returning bounded data with a
  truncation marker; branch records bypass `max_changes`; shutdown snapshots
  active workers before an admitted activation finishes.
- Strategy delta: Attempt 4 changes each controlling mechanism and adds a
  black-box regression for each, rather than suppressing or merely replying to
  the review.
- Regression guard: republish only when all focused tests, the full MCP Cordis
  package, imported skill validation, and Buildifier pass on one exact commit.
- Latest result: resolved in the Attempt 4 working tree. Four focused tests,
  all ten integrated tests, the complete package build, three skill-validation
  aspects, and root Buildifier pass. Exact-commit rerun remains required.

## Attempt 4 independent review rejects completeness claims (occurrence 1)

- Candidate: uncommitted Attempt 4 working tree after all recorded tests
  passed.
- Result: four directly related semantic and evidence gaps survived.
- Causes: invalid UTF-8 loss did not set `truncated`; several starter fields
  ignored host truncation; Git status inferred truncation from result length
  instead of actual omission and used newline-unsafe path regexes; imported
  skill eval cases did not exercise their new contracts.
- Strategy delta: Attempt 5 makes loss explicit, propagates truncation by
  field, parses one-record lookahead semantics, adds newline/exact-limit tests,
  and expands offline eval cases before another integrated run.
- Regression guard: independent review must find no correctness issue before
  delivery preparation; a green test suite alone is insufficient.

## Unavailable worker teardown escapes shutdown tracking (occurrence 1)

- Evidence: independent shutdown review of Attempt 4.
- Cause: `#handleUnavailable()` removes the activation from `#active`, while
  its worker termination is asynchronous and was not added to `#retirements`.
- Strategy delta: track the activation's idempotent `dispose()` promise as a
  retirement at the same moment it is removed.
- Regression guard: a deterministic unavailable-then-shutdown scenario must
  prove shutdown does not finish before the teardown promise.

## Main references an absent decision-review skill (occurrence 1)

- Evidence: current base `7ad2704c` mentions `decision-review` in `AGENTS.md`,
  but the referenced package is absent from that tree.
- Cause: the skill exists in open PR 24, not in the rebased `master` commit.
- Strategy delta: import only PR 24's four `projects/agents` changes, merged
  against current guidance, and supply validation assets required by current
  repository policy.
- Regression guard: build and validate the imported skill through its Bazel
  `skill_library` and offline Promptfoo target.
- Latest result: resolved in the working tree. `decision-review` matches PR
  24's instruction blob and passes quick validation, offline Promptfoo loading,
  and the repository skill-validation aspect.

## Attempt 5 independent review rejects lifecycle and compatibility

- Candidate: green uncommitted Attempt 5 working tree.
- Result: focused tests passed 8/8, integrated tests passed 12/12, the affected
  build and skill aspects passed, and Buildifier passed; independent review
  nevertheless found release-blocking defects.
- Causes: immediate `Worker.terminate()` can bypass detached-child cleanup;
  already-admitted handlers can spawn after one-shot disposal cleanup;
  promotion is outside the shutdown lock snapshot; an unavailable candidate's
  one-shot notification can be ignored during persistence; and always-success
  output truncation silently changes retained package-version semantics.
- Additional gaps: command signals/failures, exact-bound search and omitted
  records, Git optional writes/history framing, LF hash portability, weak
  integration assertions, and contradictory or domain-specific imported-skill
  rules.
- Strategy delta: Attempt 6 versions partial-output behavior explicitly,
  closes process and runtime admission, creates sole final package candidates,
  and binds every completeness claim to a direct regression.
- Regression guard: no delivery preparation until a fresh review accepts the
  new candidate after all invalidated gates pass.

## Attempt 6 first final review rejects worker-side process ownership

- Candidate: green Attempt 6 working tree before parent-owned supervision.
- Result: focused tests, the 14/14 integrated packet, build, skill validation,
  and Buildifier passed; independent reviews still rejected release.
- Causes: outer timeouts settled before cleanup; signaling a process group was
  called reaping; the worker's native-spawn/PID-publication window could not be
  both bounded and orphan-safe; direct exit could leave inherited pipes open;
  shutdown awaited but did not dispose retired generations; ripgrep byte fields
  could be returned as empty complete text; and durable records were stale.
- Strategy delta: the parent activation now owns process spawning and group
  cleanup, runtime tracks retired activations rather than only their promises,
  byte fields set explicit truncation, and the durable goal records every
  verdict and invalidated gate.
- Regression guard: immediate non-liveness assertions, inherited-pipe cleanup,
  fatal cleanup ordering, retired-generation shutdown, byte-field cases, and a
  fresh adversarial review must pass before final integrated validation.

## Custom manifests and hash versions are the wrong extension model

- Candidate: published PR 32 checkpoint `7cfef071` and its Attempt 6
  descendants.
- Result: the user rejected the per-package `manifest.json`, immutable
  `versions/`, and active-pointer design as nonstandard and temporary-looking.
- Cause: the runtime had grown a second package manager instead of using
  Cordis Loader entries, Include-backed configuration, HMR, and Git history.
- Strategy delta: Attempt 7 deletes the custom store and workers, pins the
  official services, and uses `cordis.yaml` plus ordinary ESM modules.
- Regression guard: no package manifest, hash-named source snapshot, custom
  storage layer, or version-file `.gitattributes` rule may remain.

## In-process standard Cordis needs reliability admission guards

- Candidate: independently reviewed Attempt 7 working tree.
- Result: review found that accidental stdout writes could corrupt stdio,
  Promise/async-iterator activation or top-level await could wedge lifecycle
  mutation, filename-only HMR events could acknowledge the wrong write, and
  response timeout could leave invocation-owned children running.
- Strategy delta: reserve a private protocol stream, use Node's ESM parser to
  reject top-level await, require synchronous object activation, correlate
  managed source with a named-export token, and use an invocation-scoped
  supervisor that cancels and joins `ctx.exec()`.
- Regression guard: real stdio logging, hung activation, source round-trip,
  invalid/expired invocation, and descendant non-liveness tests all pass.

## Official HMR requires Node internal-module access

- Command: first Attempt 7 `runtime_test` initialization.
- Result: HMR rejected startup with
  `--expose-internals is required for HMR service`.
- Cause: the optional native fallback peer was not reliably visible through
  Bazel's strict pnpm layout.
- Strategy delta: every runtime-bearing Bazel launcher passes the official
  package's supported `--expose-internals` Node flag. Automatic peer install is
  disabled because all required peers are pinned explicitly and the optional
  native fallback is unnecessary.
- Regression guard: the complete project build and all runtime tests must boot
  through the Bazel launchers.

## Root Fiber uid zero skips watcher disposal

- Command: runtime and starter tests without Node's force-exit option.
- Result: both test bodies passed but timed out with live `FSEventWrap`
  resources after `runtime.shutdown()`.
- Cause: shutdown guarded `root.fiber.dispose()` with `root.fiber.uid`;
  Cordis assigns the root Fiber uid `0`, so the truthiness check skipped every
  root-owned cleanup effect, including HMR watchers.
- Strategy delta: dispose whenever the root Fiber exists, independent of its
  numeric uid. Remove force-exit workarounds from tests and keep stdio shutdown
  graceful.
- Regression guard: runtime, starter, and real stdio tests must exit normally
  without `--test-force-exit` or `process.exit()`.

## Review completion is not thread reconciliation

- Candidate: published commit `cfab0fb5` after a completed hosted review.
- Result: the review summary completed, but GraphQL thread inspection exposed
  one unresolved older regex-parity finding and one new annotation finding.
- Cause: regex fallback omitted JavaScript Unicode mode and matched surrogate
  halves; source-overwriting MCP tools declared `destructiveHint: false`.
- Strategy delta: enable Unicode scalar matching, mark both potentially
  overwriting tools destructive, and treat the review-thread graph—not the
  summary state—as the authoritative review ledger.
- Regression guard: an astral `.` fallback match must be one UTF-8 span, and
  real stdio tool discovery must expose both destructive annotations.

## Fallback event and empty-range boundaries diverge

- Candidate: published review correction `bc4e5ae9`.
- Result: the next hosted review found duplicate, out-of-order fallback context
  around adjacent matches and a `null` endpoint for an existing empty line.
- Cause: context filtering knew only the current matching line, while bounded
  reads inferred range existence from non-empty selected text.
- Strategy delta: precompute match classifications, emit the union of match and
  context lines once in source order, and track whether the requested range
  exists separately from its content.
- Regression guard: adjacent matches followed by context must emit
  match/match/context in line order; a selected empty line reports its line
  number while a start past EOF reports `null`.

## Delivery adapter refuses the multi-commit feature range

- Candidate: local Attempt 9 changes above published head `bc4e5ae9`, with
  fetched base `d29f9d47`.
- Result: read-only delivery inspection reports nine feature commits and
  refuses preparation; version 1 will not infer consolidation ownership.
- Cause: only the first commit carries the adapter's ownership disclaimer,
  while the adapter supports preparation of at most one feature commit.
- Rejected workaround: direct rebase, reset, cherry-pick, or a replacement
  branch would bypass the GitHub adapter's explicit safety refusal.
- Required strategy delta: obtain scope to add a guarded exact-head,
  merge-base-aware consolidation path to `repo_delivery`, or have that support
  land separately before resuming the rebase.
- Resolution: the user explicitly authorized a guarded adapter extension.
  `repo_delivery prepare --consolidate <exact-inspected-head>` now verifies a
  merge-free linear range, identical author and committer identities, the
  oldest commit's ownership marker, unchanged pull-request projection,
  signature preservation, and every unrelated refusal before replacing the
  range. Its integration test also proves the prior remote tip remains the
  receipt-bound publication lease.

## Cordis HMR loses a write during an in-flight reload

- Candidate: Attempt 9 with unmodified
  `@deepseek-ai/cordis-plugin-hmr` 1.0.16.
- Result: a deterministic slow top-level-await replacement followed by a
  second write left the latest source on disk but the first replacement live.
- Cause: debounced `partialReload()` work was not serialized; each successful
  run reset one shared stash even when a later change arrived during import.
- Rejected wrapper workaround: public HMR emits neither import-failure nor
  settled-activation events, so a wrapper gate would either reopen unsafely on
  a timeout or wedge permanently after failure.
- Strategy delta: Attempt 10 uses standard pnpm patching to serialize the
  owning HMR task, snapshot each change set, and drain changes observed while
  it runs.
- Regression guard: overlapping writes during both slow module evaluation and
  slow asynchronous activation must converge to the latest persisted source.
