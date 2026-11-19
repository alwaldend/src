# Historical runtime-extension acceptance delta

This archived delta records the final accepted historical contract from
[acceptance.md](../../provenance/source/acceptance.md) and
[Attempt 11](../../provenance/source/attempts/011.md). It was imported as
completed history on 2026-09-08 and was not applied to the current baseline.
Earlier attempt requirements superseded by the final contract are retained
only in provenance.

## ADDED Requirements

### Requirement: Standalone Cordis MCP server

The project SHALL provide a documented, Bazel-built standalone stdio MCP
server using official Cordis Loader, Include, HMR, and Timer services, without
requiring a Codex plugin manifest.

#### Scenario: Start the server through its repository target

- **WHEN** the documented MCP target starts with an explicit workspace root
- **THEN** it serves MCP over stdio and mounts normal Cordis entries and
  modules through the official lifecycle services
- **AND** package logging does not corrupt the JSON-RPC protocol stream.

### Requirement: Stable package lifecycle gateway

The server SHALL expose fixed tools that can define, start, inspect, invoke,
update, stop, and remove packages without reconnecting the MCP client.
Invocation MUST remain available when the client caches its initial schemas.

#### Scenario: Change a package through one client connection

- **WHEN** one connected client defines and starts a package, invokes it,
  updates its source, and eventually observes the new behavior
- **THEN** the same client can inspect, stop, and remove that package
- **AND** the server process and MCP connection remain in use throughout.

### Requirement: Atomic persistence with Cordis-owned activation

MCP source writes SHALL reject invalid syntax before changing stored bytes
and SHALL atomically persist valid source. Updates to running entries SHALL
request Cordis HMR with activation reported as pending. The wrapper MUST NOT
promise synchronous activation or broad restoration of prior on-disk bytes.
Git SHALL own reusable source history.

#### Scenario: Persist an update to a running entry

- **WHEN** a client supplies syntactically valid replacement module source
- **THEN** source-tool success means the bytes were atomically persisted and
  official reload was requested
- **AND** the client observes eventual live behavior through the gateway
- **AND** evaluation or activation failures retain native Cordis semantics.

#### Scenario: Reject invalid syntax

- **WHEN** replacement source fails syntax validation
- **THEN** the source tool rejects the update before replacing the stored
  module bytes.

### Requirement: Reload scheduling preserves later writes

The pinned HMR dependency SHALL serialize reload work and drain source changes
observed during an in-flight reload so that later persisted source is not
discarded by an earlier reload.

#### Scenario: A second update overlaps slow module evaluation

- **WHEN** a replacement module is held at an explicit evaluation gate and a
  second source version is persisted before that gate is released
- **THEN** HMR eventually activates the latest persisted version.

#### Scenario: A second update overlaps asynchronous activation

- **WHEN** a replacement is held inside asynchronous `apply()` and another
  source version is persisted before activation resumes
- **THEN** HMR drains the later observed write and eventually exposes its
  behavior after the earlier lifecycle work settles.

### Requirement: Scoped persistence and promotion

The server SHALL distinguish reusable project packages from disposable
workspace-output packages by explicit scope and name. Reusable source SHALL
live beneath `projects/mcp_cordis/plugins`, and disposable source SHALL live
beneath the active workspace's ignored output. Reusable packages SHALL reload
after server restart, and disposable packages SHALL have an explicit promotion
path into the project library.

#### Scenario: Restart and promote scoped packages

- **WHEN** a disposable package is promoted and the server restarts against
  the same workspace
- **THEN** the promoted reusable package is recovered from project files
- **AND** a disposable package with the same name does not silently shadow it.

### Requirement: Useful bounded starter packages

The project SHALL include executable tests for at least three non-sensitive
starter packages: repository context, read-only Git worktree operations, and
network diagnostics. Results MUST expose truncation or unavailable fields
rather than silently claiming completeness after data loss.

#### Scenario: Execute the starter catalog

- **WHEN** the starter integration suite loads the reusable catalog through
  the runtime and invokes its handlers
- **THEN** all three packages execute their documented operations
- **AND** bounded output distinguishes complete results from actual omission
  or unsupported operations.

### Requirement: Verified repository handles and restricted search fallback

Repository reads, directory inspection, search subprocesses, and Git commands
SHALL use verified file or directory handles for the duration of access.
Without ripgrep, fallback search SHALL fail closed for directory traversal,
regular expressions, and files whose text does not round-trip as valid UTF-8.

#### Scenario: Use a verified repository selection

- **WHEN** a handler verifies a selected repository file or directory
- **THEN** its later reads and subprocess arguments address the retained
  verified handle rather than reopening a replaceable lexical alias.

#### Scenario: Search an explicit file without ripgrep

- **WHEN** ripgrep is unavailable and a caller requests bounded fixed-string
  search of an explicit valid UTF-8 file
- **THEN** the fallback returns ordered events with raw UTF-8 byte offsets
- **AND** unsupported directory, regex, or malformed-text requests fail
  explicitly instead of weakening the search contract.

### Requirement: Supervised Linux execution

On Linux, `ctx.exec()` SHALL settle normal results, timeouts, output-limit
outcomes, and plugin-disposal outcomes only after its direct child and live
members of the original process group have stopped, or report bounded cleanup
failure when inspection cannot establish cleanup. Deliberately created new
sessions are outside this trusted-package contract. Unsupported execution
platforms MUST fail closed.

#### Scenario: Timeout an execution with an original-group descendant

- **WHEN** an execution deadline expires while a descendant remains in the
  original process group
- **THEN** the supervisor stops and verifies that group before settling the
  timeout response.

#### Scenario: Process inspection fails permanently

- **WHEN** process-group inspection fails repeatedly
- **THEN** cleanup reports bounded `EXEC_CLEANUP` failure instead of retrying
  indefinitely.

### Requirement: Explicit partial output and handler lease semantics

`ctx.exec()` SHALL reject output overflow or invalid UTF-8 by default. Explicit
partial-output opt-in SHALL return only a bounded valid prefix with truncation
and output-limit signals. A gateway response deadline SHALL retain an admitted
handler's Fiber lease until it finishes, while invocation-owned subprocess
cleanup MUST precede timeout-response settlement.

#### Scenario: Opt into bounded output

- **WHEN** a package explicitly enables partial output and a process exceeds
  the combined stdout/stderr budget
- **THEN** the result contains a bounded valid prefix and reports truncation
  plus the output-limit reason after supervised cleanup.

#### Scenario: A response deadline expires before a handler completes

- **WHEN** an admitted JavaScript handler exceeds the gateway response
  deadline
- **THEN** the response deadline does not itself cancel that JavaScript
  handler or release its Fiber lease
- **AND** any invocation-owned execution is cancelled and joined before the
  timeout response settles.

### Requirement: Worktree-local MCP startup

Trusted project-scoped Codex configuration SHALL resolve the active clone or
linked worktree for MCP source and disposable state. The launcher SHALL
release Bazel's output-base lock before the long-lived stdio server serves.

#### Scenario: Build while the registered MCP is connected

- **WHEN** the worktree's registered MCP launcher initializes and lists the
  gateway tools
- **THEN** a second Bazel command in that same worktree completes while the
  MCP remains connected
- **AND** the MCP uses the explicitly selected worktree's package state.

### Requirement: Candidate-bound validation

The delivered candidate SHALL pass focused tests, affected package builds,
repository formatting, and an end-to-end stdio lifecycle and restart
transcript. Behavior-changing review corrections MUST invalidate earlier
correctness verdicts and receive affected validation plus source scrutiny.

#### Scenario: Reconcile a correctness finding after publication

- **WHEN** review identifies a valid runtime or starter-package defect
- **THEN** the corrected candidate receives focused regression coverage,
  affected complete checks, and renewed correctness scrutiny
- **AND** delivery verification binds the published candidate and review
  reconciliation to the same accepted work.
