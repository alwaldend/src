# Workspace-local Cordis runtime extensions

## Why

The original task required reusable and disposable runtime packages behind a
standalone MCP connection, with Cordis owning package lifecycle. Recurring
repository, Git, and network work justified a small initial package catalog.
The accepted outcome is recorded in the
[legacy goal](provenance/source/README.md) and its
[requirements](provenance/source/requirements.md).

This change imports that completed work as accepted history on 2026-09-08.
The legacy goal explicitly says **Complete**, reports no failing or unverified
criteria, and identifies PR 32 as delivered and review-reconciled. A separate
machine-readable execution record is unavailable in the maintained source;
its absence does not imply that the work is paused or incomplete. The import
does not claim to have repeated historical tests or checked live PR state.

## What Changes

- Provide a standalone Bazel-built stdio MCP server with fixed package
  lifecycle and invocation tools.
- Mount official Cordis Loader, Include, HMR, and Timer services using normal
  ESM modules and `cordis.yaml`, with project and disposable workspace scopes.
- Syntax-check and atomically persist package source, then let Cordis own
  activation; serialize overlapping reloads in the pinned HMR dependency.
- Preserve reusable packages across restart and support explicit promotion
  from disposable storage.
- Include tested repository-context, Git-worktree, and network-probe packages
  with bounded results and explicit completeness signals.
- Supervise Linux subprocess groups and launch the MCP from the active
  worktree without holding Bazel's output-base lock.

## Capabilities

### New Capabilities

- `project-mcp-cordis`: historically added workspace-local runtime package
  management, fixed MCP gateways, package persistence, supervised execution,
  starter packages, and project-scoped startup.

### Modified Capabilities

None are applied by this import. The archived delta expresses the final
historical acceptance contract, including corrections that superseded earlier
attempts. It was not applied to the current baseline. The
[current project spec](../../../specs/project-mcp-cordis/spec.md) is
independently sourced from the present repository.

## Impact

The original work affected `projects/mcp_cordis`, project-scoped MCP
registration, and the narrowly authorized agent-policy and delivery-adapter
changes documented in the [attempt history](provenance/source/attempts/README.md).
This migration preserves that history without reauthorizing old operations.
The [project README](../../../../README.md) owns current
runtime behavior; [design.md](design.md) records accepted and superseded
decisions, and [tasks.md](tasks.md) maps historical completion to its evidence.
