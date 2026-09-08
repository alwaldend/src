# Acceptance criteria

[Back to durable goal](./)

1. `projects/mcp_cordis` is a documented, Bazel-built standalone MCP stdio
   server and is not coupled to a Codex plugin manifest.
2. The server demonstrably uses Cordis lifecycle primitives rather than a
   separately invented plugin framework.
3. One connected MCP client can define, start, inspect, invoke, update, stop,
   and remove a runtime package without restarting the server.
4. MCP source writes are syntax-checked and atomically persisted. Existing
   running entries are handed to Cordis HMR without a private
   acknowledgement protocol. The MCP does not promise synchronous activation,
   broad activation rollback, or restoration of prior on-disk bytes. Native
   Cordis failure behavior remains intact. The pinned HMR dependency serializes
   reload work and drains later observed writes so an in-flight replacement
   cannot discard the latest persisted source. Git owns reusable source
   history; the runtime does not invent a second version store.
5. Reusable package source is stored beneath
   `projects/mcp_cordis/plugins`; disposable source and configuration resolve
   beneath the current workspace's `out/mcp_cordis`.
6. Reusable packages reload after an MCP server restart. Disposable packages
   have an explicit promotion path into the reusable library.
7. At least three useful, non-sensitive starter packages are justified by
   recurring prior-session tasks and have executable tests.
8. Stable gateway tools allow immediate invocation even if a particular MCP
   client does not refresh dynamically registered schemas.
9. Focused tests, builds, repository formatting checks, and an end-to-end MCP
   transcript pass for the exact delivered candidate.
10. On Linux, `ctx.exec()` results, execution timeouts, and plugin disposal
    settle only after the direct child and live members of its original process
    group stop. A trusted package that deliberately creates a new session is
    explicitly outside that guarantee; unsupported platforms fail closed.
    `cordis_invoke.timeout_ms` is a response deadline, not handler cancellation;
    the admitted handler retains its Fiber lease until completion.
11. Trusted Codex sessions discover the MCP from the checked-in
    `.codex/config.toml` in the active clone or linked worktree. Server startup
    releases Bazel's output-base lock before serving stdio, and source plus
    disposable state remain scoped to that worktree.

## Evidence plan

- Unit tests cover validation, standard Cordis entries, lifecycle disposal,
  persistence roots, eventual reload, and promotion.
- An MCP integration test drives initialize/list/define/run/invoke/update/stop
  over stdio without restarting the process.
- A restart test proves reusable package recovery.
- Package tests execute every starter package through the same runtime path.
- Process regressions cover normal results, output limits, valid UTF-8 prefix
  retention, timeout, disposal, and rejection of launches after disposal.
- HMR regressions write a second generation during slow top-level evaluation
  and slow asynchronous activation, then require convergence to the latest
  persisted generation.
- Bazel query, test, build, Buildifier, and `git diff --check` cover repository
  integration.
- A launcher probe starts the generated-script MCP and runs a second Bazel
  command concurrently in the same linked worktree.

## Fixed regression set

- `git diff --check`
- Focused `//projects/mcp_cordis:all` Bazel tests and build
- `//:buildifier_test` after every BUILD or Bzlmod change
- End-to-end stdio lifecycle and restart tests
