## Why

AL's private executable cache can reuse stale tools and bypasses Bazel's source
tracking. Explicit installation gives host commands a clear update boundary.

## What Changes

- Remove AL's cached-tool command, Lua declarations and configuration schema.
- Add a standalone binary installer and explicit install targets for repo_delivery and bazel_agent.
- Keep ordinary local Bazel builds and direct installed command execution.
- **BREAKING**: remove `al tool`, `tool`, and `tool_method` configuration.

## Capabilities

### New Capabilities

### Modified Capabilities

- `tool-cache`: remove private caching in favor of explicit host installation.

## Impact

AL command/configuration, the file installer, repo_delivery targets, and agent
invocation guidance. No download target, flag, upload or host deployment.
