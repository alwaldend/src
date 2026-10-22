## Context

Deferred; implementation has not started. Authority is planning only. See the
[proposal](proposal.md) for source evidence and the static finding's limits.
The current open uses `O_RDONLY | O_NOFOLLOW` before descriptor inspection.
Gateway timeouts do not cancel an admitted JavaScript handler, so an open that
never settles can retain its Fiber lease.

## Goals / Non-Goals

**Goals:** Reject a FIFO without a writer promptly, preserve path confinement,
and finish handler cleanup so stop and reload can drain normally.

**Non-Goals:** Add general JavaScript cancellation, sandbox trusted plugins,
or change lifecycle semantics for unrelated long-running handlers.

## Decisions

- Evaluate adding nonblocking open while retaining descriptor validation and
  regular-file checks as the smallest candidate. A pathname type check alone
  cannot protect against a replacement between check and open; a gateway
  timeout alone leaves the pending open and lease behind.
- Keep opened-descriptor containment authoritative and close every acquired
  descriptor on rejection. Audit every caller of the shared path-binding
  helper because some legitimate operations target directories.
- Reproduce through a disposable MCP server process with an outer deadline
  before changing the implementation. Exercise the real handler and lifecycle
  paths without placing a never-settling request in the agent's live runtime.
- Reuse Node built-ins and existing runtime test infrastructure; no additional
  dependency is proposed. Confirm platform behavior on the currently supported
  Linux runtime before choosing final open flags.

## Risks / Trade-offs

- Nonblocking flags differ across platforms and file types -> keep supported
  platform behavior explicit and test regular files and directories too.
- An error path leaks a descriptor -> verify repeated rejection and cleanup,
  followed by a successful ordinary read and lifecycle transition.
- Timing-only assertions are unstable -> set an outer process deadline and
  record semantic rejection and successful stop/reload as the primary evidence.

## Migration Plan

Next action after implementation is requested: reproduce the FIFO case in the
isolated E2E harness and enumerate path replacement and cleanup failures.
Implement the bounded open/validation change, rerun relevant existing
containment coverage, and save a repeatable rejection/lifecycle result.
No package reload or code mutation is performed by this deferred plan.
