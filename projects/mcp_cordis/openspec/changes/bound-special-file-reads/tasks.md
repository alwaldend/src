Deferred; implementation has not started. These unchecked tasks are future
work; this planning request authorizes no runtime or production-code change.

## 1. Reproduce and define bounded rejection

- [ ] 1.1 Inventory shared path-binding callers and failure modes, including
      FIFO without writer, target replacement, directory use, and cleanup;
      record the expected file-type and containment outcomes.
- [ ] 1.2 Add an isolated MCP E2E exercise before implementation with an outer
      process deadline; demonstrate the FIFO gap without blocking the agent's
      live runtime and emit a repeatable reproduction result.

## 2. Implement and verify cleanup

- [ ] 2.1 Implement nonblocking open and descriptor validation or an equally
      bounded approach using existing dependencies; verify ordinary files and
      valid directory operations retain their existing behavior.
- [ ] 2.2 Run repeated special-file rejection, containment/replacement, normal
      read, and stop/reload scenarios against the exact candidate; preserve a
      result artifact proving prompt errors and completed resource cleanup.
