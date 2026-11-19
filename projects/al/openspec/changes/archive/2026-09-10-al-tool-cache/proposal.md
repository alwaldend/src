## Why

Tool invocation is currently bound to `bazel_agent tool`, whose tool catalog is
hard-coded in the Bazel runner. AL owns command configuration in `al.lua`, but
cannot yet expose repository tools as cacheable commands.

## What Changes

- Add `al tool` as a cached command runner.
- Load tool declarations from AL configuration, including the root `al.lua` by
  default.
- **BREAKING**: remove `bazel_agent tool` and its source-hash cache contract.
- Install a cache miss with a configured Bazel label and move the selected
  executable into a user-local XDG cache.

## Capabilities

### New Capabilities

- `tool-cache`: Load Bazel-backed tools from AL configuration and run them from
  a private source-keyed executable cache.

### Modified Capabilities

## Impact

- `projects/al` command dispatch and configuration schema.
- Root and component `al.lua` files that declare tools.
- `projects/bazel_agent`, `projects/mcp_cordis/cmd/mcp_cordis/launch.sh`,
  `tools/repo_delivery/README.md`, and `projects/bazel_agent/README.md`.
