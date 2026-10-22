## Why

Deferred; implementation has not started. Review finding 4 identifies a
potential indefinitely pending read: `bindPath` opens a path before inspecting
the opened file type, so a FIFO without a writer can block. A gateway timeout
does not finish that handler or release its runtime lease.

## What Changes

- Reject unsupported special files without waiting for a peer to open them.
- Preserve descriptor-based containment, symlink, and file-type checks through
  the complete open/read/close sequence.
- Exercise FIFO rejection through the MCP gateway and verify subsequent
  ordinary reads, stop, and reload finish within explicit test deadlines.
- Accept only with a repeatable E2E result and descriptor cleanup evidence.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-mcp-cordis`: bound special-file rejection in repository reads.

## Impact

Planning only is authorized. Future work belongs to `repo_context` and its
runtime integration tests. The scope does not promise general JavaScript
cancellation or change trusted-package lifecycle semantics. No dependency or
live runtime change is authorized now.

Baseline: `e0d24cb3a`, inspected 2026-10-22. Evidence:
[descriptor binding](../../../plugins/repo_context.mjs) and
[gateway timeout and Fiber lifetime](../../../README.md#cordis-configuration-and-plugin-format).
This is a static finding; no hanging FIFO reproduction was run during planning.
