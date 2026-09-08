# Historical completed work

Checked items record completion supported by the maintained legacy records;
they do not claim that migration repeated the implementation or its tests.
The [legacy status](provenance/source/README.md#status) says Complete and
reports no failing or unverified criteria. Machine-readable execution state
is unavailable. This completed history was imported into the archive without
applying its delta to the current baseline.

## 1. Final acceptance criteria

- [x] 1.1 Deliver a documented Bazel-built standalone stdio MCP server,
      independent of a Codex plugin manifest (acceptance 1; Attempts 7 and 11).
- [x] 1.2 Use official Cordis Loader, Include, HMR, and Timer lifecycle
      services with ordinary modules (acceptance 2; Attempts 7, 9, and 10).
- [x] 1.3 Exercise define, start, inspect, invoke, update, stop, and remove
      over one MCP connection (acceptance 3; real stdio evidence in Attempt 7,
      followed by final complete packet in Attempt 11).
- [x] 1.4 Syntax-check and atomically persist source, use native eventual
      HMR activation, and serialize overlapping reloads in the pinned dependency
      (acceptance 4; Attempts 9 and 10, carried into Attempt 11).
- [x] 1.5 Separate reusable project packages from workspace-local disposable
      output with explicit scope identities (acceptance 5; Attempts 7 and 10).
- [x] 1.6 Recover reusable packages after restart and support explicit
      disposable-to-project promotion (acceptance 6; runtime and stdio evidence).
- [x] 1.7 Deliver and execute the three non-sensitive starter packages
      `repo_context`, `git_worktree`, and `network_probe` (acceptance 7; preflight,
      starter integration, and final review evidence).
- [x] 1.8 Keep fixed MCP discovery and invocation gateways usable by clients
      that cache dynamic schemas (acceptance 8; stdio lifecycle evidence).
- [x] 1.9 Pass focused tests, full affected builds, formatting, and real
      stdio coverage for the delivered candidate, rerunning invalidated gates
      after review corrections (acceptance 9; Attempt 11 final review evidence).
- [x] 1.10 Supervise Linux original-process-group cleanup before settlement,
      retain handler Fiber leases across response timeouts, and fail closed on
      unsupported execution platforms (acceptance 10; Attempts 6, 10, and 11).
- [x] 1.11 Register the MCP in trusted project configuration, resolve the
      active clone or linked worktree, and release the Bazel output-base lock
      before serving stdio (acceptance 11; launcher probes in Attempts 10 and 11).

Sources: [final acceptance](provenance/source/acceptance.md),
[evidence manifest](provenance/source/evidence.md),
[Attempt 7](provenance/source/attempts/007.md),
[Attempt 10](provenance/source/attempts/010.md), and
[Attempt 11](provenance/source/attempts/011.md).

## 2. Accepted review and delivery work

- [x] 2.1 Replace the rejected worker/version-store model with standard
      Cordis configuration, then remove wrapper-side HMR acknowledgement and
      rollback machinery (legacy plan steps 12–13 and 19; Attempts 7 and 9).
- [x] 2.2 Preserve strict default output behavior and explicit partial-output
      opt-in, then correct exact limits, UTF-8 handling, and process cleanup
      (legacy plan steps 7–11 and 16–18; Attempts 4–6, 8, and 11).
- [x] 2.3 Import only the authorized PR 24 agent-tree changes, preserve
      advancing base policy, and reconcile the then-current goal and discovery
      layouts (legacy plan steps 6, 9–11, and 21; Attempts 4–6 and 11).
- [x] 2.4 Add a pinned HMR scheduler patch and deterministic overlapping
      evaluation/activation regressions (legacy plan step 20; Attempt 10).
- [x] 2.5 Correct verified-handle reads and subprocess paths, bound permanent
      process-inspection failures, and fail closed for unsupported fallback
      search semantics (final hosted review in Attempt 11).
- [x] 2.6 Consolidate the authorized owned range through the guarded delivery
      adapter, rebase, validate the candidate, publish PR 32, reconcile review
      threads, and verify the receipt-bound remote state (legacy plan steps
      14–15 and 21; Attempt 11 and legacy completed status).

Sources: [completed plan](provenance/source/README.md#current-plan),
[failure ledger](provenance/source/failure_ledger.md), and
[final review evidence](provenance/source/evidence.md#attempt-11-final-review-evidence).

The recorded optional next action, merging PR 32 when desired, is not a
completed implementation task and is not asserted here. No pending task is
invented from missing execution metadata or a mutable receipt.
