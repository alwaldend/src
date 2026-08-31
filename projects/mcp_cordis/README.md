---
title: MCP Cordis
description: Workspace-local runtime packages behind a stable MCP server
---

`mcp_cordis` is a standalone stdio MCP server that mounts runtime JavaScript
packages through [Cordis](https://github.com/deepseek-ai/deepseek-harness).
It is intentionally an MCP server, not a Codex plugin bundle.

Reusable definitions are ordinary ESM files in
`projects/mcp_cordis/plugins`, listed by `projects/mcp_cordis/cordis.yaml`.
Disposable definitions use the same layout under `out/mcp_cordis`. Every
package is addressed by both scope and name, so a scratch package never
silently shadows a reusable package.

## Run

The repository's `.codex/config.toml` registers `mcp_cordis` as a
project-scoped stdio server. Codex loads that file for a trusted workspace and
finds the active Git worktree before starting the server. Separate clones and
worktrees therefore use their own source, `projects/mcp_cordis` packages, and
`out/mcp_cordis` scratch packages. Trusting the repository's root checkout
also covers its linked worktrees; a glob trust entry is neither needed nor
supported. A new Codex session is needed after the MCP registration itself is
first added; package changes after that do not require another session.

The registration calls `cmd/mcp_cordis/launch.sh`, which asks Bazel to
generate a launch script under `out/mcp_cordis` and then executes that script.
Bazel exits before the stdio server starts, so the long-lived MCP does not hold
the worktree's Bazel output-base lock. Builds and tests can run normally while
Codex remains connected.

To build and run the server directly from the repository root:

```sh
bazel_agent run //projects/mcp_cordis:mcp_cordis -- \
  --workspace-root "$PWD"
```

The workspace root is mandatory unless `BUILD_WORKSPACE_DIRECTORY` is
present. The checked-in launcher resolves the current Git worktree explicitly
and supplies that path to the server.

The fixed `cordis_*` tools define, start, inspect, invoke, update, stop,
remove, and promote packages without reconnecting the MCP client. Package
handlers are called through `cordis_invoke`; this remains reliable even when
an MCP client caches its initial tool list.

## Cordis configuration and plugin format

`cordis.yaml` uses the standard Cordis Include entry-list format:

```yaml
- id: hello
  name: ./plugins/hello.mjs
```

The referenced file is a normal ESM Cordis plugin:

```js
const plugin = {
  description: "Provide a greeting.",
  apply(ctx) {
    ctx.tool(
      {
        name: "hello_world",
        description: "Return a greeting.",
        inputSchema: {
          type: "object",
          properties: { name: { type: "string" } },
          additionalProperties: false,
        },
      },
      ({ name = "world" }) => ({ greeting: `Hello, ${name}!` }),
    );
  },
};

plugin.apply.description = plugin.description;

export default plugin;
```

Cordis normalizes an object plugin to its `apply` callback. Attaching the
optional package description to that callback exposes it through
`cordis_list` and `cordis_inspect`; tool descriptions remain part of each
`ctx.tool()` definition.

The server mounts the official Cordis Loader, Include, and HMR services.
`cordis_define` syntax-checks and atomically persists the ordinary module.
Creating or enabling an entry refreshes any cached module through Cordis HMR,
then uses the public Include refresh API and waits for activation. Updating an
already-running entry returns `activation: "pending"`; poll `cordis_invoke` or
`cordis_list_tools` until the new behavior is visible.
The reproducibly pinned HMR package carries a focused pnpm patch that
serializes module reloads and drains source changes arriving during an
in-flight reload, so the latest persisted source is not lost.

Syntax errors are rejected before the file changes. Evaluation and `apply()`
failures follow native Cordis HMR behavior; the wrapper does not add a second
activation transaction around them. It also does not inject source markers,
inspect Loader caches, correlate watcher events, or maintain its own source
rollback/version store. Reusable history is normal Git history. Manual edits
to watched plugin files are also picked up by Cordis HMR.

Runtime modules use normal Cordis semantics, including static imports,
top-level `await`, and asynchronous `apply(ctx, config)`. Package code is
trusted: a never-settling module evaluation or activation can therefore stall
Cordis lifecycle work. The stdio launcher reserves its protocol stream and
redirects package stdout to stderr, keeping accidental `console.log()` calls
off the JSON-RPC wire.

The package context exposes `ctx.workspaceRoot`, `ctx.resolveWorkspace()`,
`ctx.readText()`, and structured `ctx.exec()` in addition to `ctx.tool()`.
`ctx.exec()` returns `code`, `signal`, `stdout`, `stderr`, `truncated`, and
`outputLimitExceeded`; `maxBytes` is a combined stdout/stderr budget. By
default, exceeding that budget or producing invalid UTF-8 rejects with
`EXEC_OUTPUT_LIMIT` or `EXEC_INVALID_UTF8`. Packages that explicitly set
`allowTruncatedOutput: true` instead receive the valid retained prefix with
`truncated` set; `outputLimitExceeded` distinguishes the byte cap from UTF-8
loss. A Fiber-owned supervisor admits each launch atomically, and results
settle only after the direct child and every live member of its original Linux
process group have stopped. Limits, timeouts, and plugin disposal use the same
cleanup path. A process that deliberately creates a new session escapes that
group and is outside this trusted-package contract. `ctx.exec()` therefore
fails closed with `EXEC_UNSUPPORTED_PLATFORM` away from Linux. Package code
also has normal Node built-ins; this host is a reliability boundary, not a
security sandbox.

`cordis_invoke.timeout_ms` bounds how long the gateway waits for a result; it
does not cancel an already admitted JavaScript handler. The handler keeps its
Fiber lease until it finishes, so stop, remove, and shutdown wait for it.
Cordis HMR waits for a retired Fiber to finish draining before it activates
and publishes the replacement, so a live invocation can delay a reload. Any
`ctx.exec()` launched by a timed-out invocation is cancelled and its process
group is confirmed stopped before the timeout response settles.

## Included packages

- `repo_context`: bounded repository reads and searches.
- `git_worktree`: read-only branch, status, log, and comparison snapshots.
- `network_probe`: DNS, TCP/TLS, and HTTP diagnostics.

These were selected from aggregate recurring task categories in recent local
sessions. No transcript content, credentials, or private outputs are included.

## Durable goal

Future runtime-extension work should resume from the maintained
[runtime extensions goal](goals/runtime_extensions/), which records acceptance
criteria, decisions, failed attempts, and supporting evidence.

## Dependency notices

The runtime directly uses `@deepseek-ai/cordis`, its official Loader, Include,
HMR, and Timer plugins, and the Model Context Protocol TypeScript SDK. Their
license texts are retained in the resolved package artifacts by the pinned
pnpm/Bazel dependency graph.
