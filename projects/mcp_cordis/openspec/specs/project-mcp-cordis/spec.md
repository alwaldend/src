# MCP Cordis Specification

## Purpose

Describe the standalone stdio MCP server that exposes workspace-local Cordis
JavaScript packages through a stable gateway and bounded repository tools.

Sources: [project README](../../../README.md),
[runtime targets](../../../BUILD.bazel),
[pinned dependency integration](../../../include.MODULE.bazel),
and [runtime integration tests](../../../test/runtime_test.mjs).
Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`,
observed 2026-09-08. Runtime availability requires a session observation;
maintained work now follows the migrated OpenSpec change.

## Requirements

### Requirement: Resolve packages within the active workspace

The server SHALL use an explicit workspace root or
`BUILD_WORKSPACE_DIRECTORY`, load reusable packages from that workspace's
`cordis.yaml`, and address every package by scope and name. Disposable packages
SHALL use task/run-namespaced scratch so they cannot silently shadow reusable
packages with the same name.

#### Scenario: Reusable and scratch packages share a name

- **WHEN** both scopes contain a package with the same name
- **THEN** invocation selects the package using its explicit scope and name.

### Requirement: Support package changes through a stable MCP gateway

The fixed `cordis_*` interface SHALL discover and invoke live package handlers
without requiring the MCP client to refresh its initial tool list.
`cordis_define` SHALL reject invalid module syntax before changing the file and
use native Cordis activation and HMR for valid source.

#### Scenario: Define and invoke a new package during a session

- **WHEN** a valid package is persisted and activated through `cordis_define`
- **THEN** `cordis_list_tools` exposes its live handlers and `cordis_invoke` can
  call them through the existing MCP connection.

### Requirement: Bound subprocess output and join process cleanup

`ctx.exec()` SHALL apply its combined stdout/stderr byte budget and reject
overflow or invalid UTF-8 unless the package explicitly permits truncated
output. On Linux, execution results SHALL settle after the child and live
members of its original process group stop. Unsupported platforms SHALL fail
closed; trusted code that creates a new session lies outside this boundary.

#### Scenario: A command exceeds the default output budget

- **WHEN** a package command produces output above `maxBytes` without opting
  into truncation
- **THEN** execution fails with `EXEC_OUTPUT_LIMIT` and follows the supervised
  process-cleanup path before settling.

### Requirement: Distinguish invocation timeout from handler cancellation

An invocation timeout SHALL bound gateway waiting without claiming cancellation
of an admitted JavaScript handler. The handler SHALL retain its Fiber lease
until completion; package shutdown and replacement SHALL wait for draining.
Commands launched by a timed-out invocation SHALL be cancelled and reaped before
the timeout response settles.

#### Scenario: A JavaScript handler remains active after timeout

- **WHEN** gateway waiting expires before the handler completes
- **THEN** the gateway reports the timeout while the handler's Fiber remains
  leased and can delay package reload or shutdown.
