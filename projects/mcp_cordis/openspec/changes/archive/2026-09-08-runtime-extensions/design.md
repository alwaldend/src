# Accepted runtime-extension design history

## Context and status

This is an import of completed work, not a new implementation proposal.
The [legacy status](provenance/source/README.md#status),
[final attempt](provenance/source/attempts/011.md), and
[final review evidence](provenance/source/evidence.md#attempt-11-final-review-evidence)
record completion, publication on PR 32, and review reconciliation. The goal's
last stated next action was to merge PR 32 when desired; this archive does not
assert that it was merged. The final published commit identifier and mutable
delivery receipt are not present as a complete final receipt in the maintained
records. Machine-readable execution state is unavailable.

Historical tests and publication claims remain attributed to those records.
The archive date is the migration date, not the original implementation date.
The delta in `specs/project-mcp-cordis/spec.md` was imported directly as
accepted history and was not applied to the current baseline. Consult the
[current project README](../../../../README.md) for later
changes, including the current task/run scratch layout and cached launcher.

## Goals and boundaries

The [final acceptance criteria](provenance/source/acceptance.md) require a
client-neutral standalone MCP server, real Cordis lifecycle ownership, a
complete package lifecycle over one connection, two storage scopes, restart
and promotion, three tested starter packages, stable invocation, deterministic
validation, supervised Linux execution, and worktree-local registration.

Trusted runtime modules may use Node built-ins. Process supervision is a
reliability boundary, not a security sandbox. Deliberate new process sessions
are outside the original-process-group cleanup contract. Importing these
requirements grants no new infrastructure or remote-publication authority.

## Accepted decisions

### Cordis owns the runtime lifecycle

The accepted design uses official Loader, Include, HMR, and Timer packages,
standard `cordis.yaml` entries, and ordinary ESM modules. Explicit scope and
name distinguish reusable from disposable packages. Git owns reusable source
history. Fixed MCP list and invoke gateways support clients that cache their
initial tool catalog. The project is an ordinary root-workspace Bazel package
with pinned JavaScript dependencies, not a Codex plugin bundle.

Sources: [Attempt 7](provenance/source/attempts/007.md), refined by
[Attempt 9](provenance/source/attempts/009.md) and
[Attempt 10](provenance/source/attempts/010.md).

### Persistence and activation have separate outcomes

Source tools validate syntax and atomically replace bytes. An update to an
already running entry returns persisted source with activation pending; it
does not promise a synchronous activation transaction or broad on-disk
rollback. Public Include refresh owns entry-list changes. Normal module
imports, top-level await, and asynchronous `apply()` remain supported under
trusted Cordis semantics.

Attempt 9 removed injected markers, acknowledgement waiters, Loader cache
inspection, and wrapper-owned source rollback. Independent review then proved
that unmodified HMR 1.0.16 could lose a write during an in-flight reload.
Attempt 10 applies a reproducibly pinned pnpm patch at the owning HMR scheduler
to serialize reloads, snapshot observed changes, and drain later writes.
Explicit release gates cover overlapping module evaluation and asynchronous
activation. Refreshing a disabled entry also obtains its latest persisted
module before enabling it.

Sources: [Attempt 9 verdict](provenance/source/attempts/009.md#current-verdict),
[Attempt 10 evidence](provenance/source/attempts/010.md#current-evidence), and
[HMR failure](provenance/source/failure_ledger.md#cordis-hmr-loses-a-write-during-an-in-flight-reload).

### Execution is supervised and output loss is explicit

`ctx.exec()` rejects overflow and invalid UTF-8 by default. Explicit partial
output opt-in returns a bounded valid prefix with truncation and output-limit
signals. Linux process-group cleanup precedes settlement after results,
timeouts, and disposal; unsupported platforms fail closed. Invocation timeout
bounds the gateway response rather than cancelling an admitted JavaScript
handler, which retains its Fiber lease. Invocation-owned subprocesses are
cancelled and joined before its timeout response settles.

The final hosted review bounded permanent process-inspection failures with
`EXEC_CLEANUP` and replaced a process-start timing assumption in coverage.
Package stdout is isolated from the stdio JSON-RPC stream, and shutdown
disposes the root Fiber even when its numeric identifier is zero.

Sources: [final acceptance](provenance/source/acceptance.md),
[Attempt 6](provenance/source/attempts/006.md), and
[Attempt 11](provenance/source/attempts/011.md).

### Starter packages preserve bounded-read and path semantics

`repo_context`, `git_worktree`, and `network_probe` implement recurring generic
workflows without copying session transcripts or protected data. Review
corrections made limits distinguish actual omission, preserve UTF-8 byte
offsets, retain partial-startup errors, and report retained read endpoints.
Repository reads, directory inspection, search subprocesses, and Git commands
use verified file or directory handles throughout use. The final fallback
supports bounded fixed-string search on explicit valid UTF-8 files; directory
search, regex search, and malformed UTF-8 fail closed without ripgrep.

Sources: [preflight evidence](provenance/source/evidence.md#preflight-evidence),
[Attempt 8](provenance/source/attempts/008.md), and
[final review evidence](provenance/source/evidence.md#attempt-11-final-review-evidence).

### Delivery follows the exact accepted candidate

The historical user authorized importing only PR 24's `projects/agents`
changes and later a guarded delivery-adapter consolidation path. The adapter
retained ownership, linearity, exact-head, projection, signature, and remote
lease checks. Rebase preserved the incoming `projects/goal` and skill
discovery layouts as they existed then. Review findings invalidated prior
correctness verdicts and were followed by source scrutiny and affected tests.
These are historical decisions, not instructions to restore the goal system
being replaced by this migration.

Sources: [requirement changes](provenance/source/requirements.md#requirement-changes),
[Attempt 10](provenance/source/attempts/010.md), and
[Attempt 11 reconciliation](provenance/source/attempts/011.md#reconciliation).

## Superseded approaches

| Attempts | Historical approach                                                                                  | Accepted replacement                                                                  |
| -------- | ---------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| 1–6      | Worker generations, per-package manifests, immutable hash-named source versions, and active pointers | Official Cordis entries, ordinary ESM files, and Git history in Attempt 7             |
| 2–3      | Elapsed time inferred an in-flight lifecycle call                                                    | Explicit started/release handshakes                                                   |
| 4–5      | All overflow returned partial success                                                                | Explicit partial-output opt-in preserving rejection by default in Attempt 6           |
| 6        | Worker-side process spawning and one-shot teardown                                                   | Supervisor-owned admission and original-process-group cleanup                         |
| 7–8      | Synchronous activation restrictions, correlated HMR tokens, and wrapper rollback                     | Normal Cordis semantics and atomic persistence in Attempt 9                           |
| 9        | Unmodified HMR for overlapping writes                                                                | Dependency-owned serialized reload scheduling in Attempt 10                           |
| 2–8      | JavaScript fallback attempted broad ripgrep parity                                                   | Restricted explicit-file fixed search and fail-closed unsupported cases by Attempt 11 |

The full [attempt index](provenance/source/attempts/README.md),
[failure ledger](provenance/source/failure_ledger.md),
[evidence manifest](provenance/source/evidence.md), and
[artifact log](provenance/source/artifacts.md) preserve rejected candidates,
intermediate verdicts, exact identifiers, and historical validation receipts.
Earlier claims of rollback or regex parity do not override the final contract.

## Validation and remaining uncertainty

Attempt 11 records exact-candidate tests and builds, a real stdio launcher
probe concurrent with Bazel, publication receipt verification, and repeated
focused regressions after hosted corrections. Its final evidence includes
the complete seven-test and sixteen-target MCP packets after the last UTF-8
correction. The legacy status reports no open implementation criteria.

No historical command was rerun merely to migrate these records. Runtime
execution state, final mutable receipts, and current PR state are unavailable
to this historical import. Future changes should start a new OpenSpec change
against the independently maintained current baseline.
