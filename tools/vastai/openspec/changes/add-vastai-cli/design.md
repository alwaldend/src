## Context

Existing tools use py_console_script_binary and tools/py/uv.lock. Upstream v1.10.0 publishes a Python wheel; its installer also resolves Python packages.

## Goals / Non-Goals

**Goals:** Run the unmodified pinned CLI through Bazel with repository-internal visibility.

**Non-Goals:** Host installation, credentials provisioning, resource rental, and SDK consumers.

## Decisions

Use a component-local Python manifest and native lock instead of executing upstream's installer. Shared resolution failed because upstream pins cryptography==49.0.0 while the shared manifest requires >=50.0.2. Pin 1.10.0, the current stable release. Use the existing console-script rule with no first-party Python wrapper. Validate help output without credentials or API operations.

## Risks / Trade-offs

The upstream package introduces transitive dependencies; review the generated lock for unrelated changes and compatibility. Credentials and charged operations remain explicit CLI user actions.

## Execution evidence

- `bazel run //tools/vastai -- --help` passed; representative output is under ignored `out/add-vastai/help.txt`.
- The shared lock experiment failed with an exact cryptography version conflict; no shared manifest or lock change remains.
- `al delivery provider` confirmed the configured GitHub repository and SSH transport. `al delivery inspect` failed with HTTP 401 during GitHub repository preflight. Device login subsequently succeeded; API and SSH authentication were verified, and delivery inspection against master passed.
- Session review: inspecting the release manifest earlier would have avoided the initial standalone-binary assumption. Pip registration remains with tools/py so the runtime version retains one owner.

- Validation passed: CLI `--version` reports 1.10.0; 27 repository quality tests, affected-target semantic lint, and strict OpenSpec validation passed.
