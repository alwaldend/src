## Context

See proposal.md for motivation. AL is an independent nested module, and the root
repository already pins `github.com/modelcontextprotocol/go-sdk` at v1.8.0.
The existing Bazel package owns enclosing-workspace resolution.

## Goals / Non-Goals

**Goals:** Provide a stdio MCP adapter with explicit schemas, process cancellation,
bounded output, and a fixed workspace and executable selected at server startup.

**Non-Goals:** Remote hosting, host configuration, arbitrary command execution,
`bazel run`, dependency upgrades, and replacing affected-test selection.

## Decisions

- Add the adapter to `internal/bazelcmd`, exposed as `al bazel mcp serve`. Use the
  official SDK already present at its existing pin instead of introducing a
  community Bazel wrapper or implementing JSON-RPC.
- Register six tools: `bazel_query`, `bazel_cquery`, `bazel_aquery`, `bazel_build`,
  `bazel_test`, and `bazel_info`. Use typed inputs for expressions, targets,
  output formats, and option arrays. Startup options are operator-controlled;
  tool options cannot select another command. Targets must use absolute labels.
- Use the normal Bazel executable on PATH, preserving Bazelisk pins and rc files.
  Select one workspace with the existing workspace resolver. Permit an explicit
  executable path for integration fixtures and installations.
- Give every subprocess a server-selected deadline, bounded stdout and stderr
  buffers, and inherited request cancellation. Return nonzero exits as MCP tool
  errors while preserving structured diagnostics and partial output indicators.
  Reject invalid inputs before launching Bazel. Never connect Bazel to MCP stdin
  or stdout directly. Serialize calls to avoid unnecessary Bazel lock queues.
- Use SDK-managed stdio framing, discovery, validation, and cancellation. Borrow
  Cobra's command streams; write diagnostics only to stderr.
- Write an executable E2E fixture before implementation. Exercise MCP discovery,
  successful operations, invalid input, failed processes, bounded output, timeout,
  request cancellation, and EOF cleanup. Validate representative queries and
  builds against a real isolated Bazel workspace outside the sandboxed test.

## Risks / Trade-offs

- Bazel rc files, build rules, tests, and options can execute code or contact
  configured services. This interface requires a trusted local workspace and
  client; MCP annotations are advisory and are not an authorization boundary.
- Bounded output can contain an incomplete query result. Report truncation
  explicitly so consumers cannot mistake it for complete graph coverage.
- Query tools can update Bazel caches. Read-only annotations describe source
  intent, not filesystem immutability or remote cache effects.
- Cancellation stops the Bazel client; Bazel owns its server and action cleanup.
  No rollback or persistent-server termination is promised.

## Migration Plan

The command is opt-in and requires no host change. Clients start `al bazel mcp serve`
with an explicit workspace. Removing that client entry disables the integration.
