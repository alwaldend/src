## 1. Configuration

- [x] 1.1 Add tool and tool-method messages to AL configuration and verify Lua
      parsing and merging with unit tests.
- [x] 1.2 Add root `al.lua` declarations for the current cached tools and verify
      the merged configuration exposes unique names.

## 2. Tool command and cache

- [x] 2.1 Implement `al tool` dispatch, default root configuration discovery,
      XDG cache selection, source-key hashing, and executable validation; verify
      with Go unit and command tests.
- [x] 2.2 Implement locked atomic Bazel build, install, and direct child
      execution; verify miss, hit, invalid output, concurrent miss, and child
      exit-status behavior.

## 3. Runner migration

- [x] 3.1 Remove `bazel_agent tool`, its registry, cache, and tests; verify the
      subcommand fails closed and owning Bazel checks pass.
- [x] 3.2 Migrate the MCP launcher and tool documentation, remove obsolete
      compatibility probes, and verify relevant package checks.

## 4. Validation

- [x] 4.1 Validate the OpenSpec change and run the owning AL and Bazel-agent
      package tests.
- [x] 4.2 Exercise an actual `al tool` hit and miss in the repository workspace
      and record evidence in the change.
