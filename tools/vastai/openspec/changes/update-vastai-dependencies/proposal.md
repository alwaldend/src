## Why

The November 4 Dependabot screenshot reports 63 open alerts, with the displayed
alerts affecting `tools/vastai/uv.lock`. Upstream Vast.ai pins vulnerable Pillow
and cryptography releases, and the local manifest unnecessarily pins Vast.ai.

## What Changes

- Remove the exact Vast.ai requirement and version-specific README text.
- Override vulnerable upstream dependency pins with secure releases and regenerate
  the component-local native uv lock through its owning Bazel target.
- Raise the environment's Python minimum to 3.11 for current Pillow support.
- Verify resolved dependencies and representative offline CLI behavior.

## Capabilities

No specification requirements change. The CLI remains reproducible through
`uv.lock`; this is dependency maintenance with `skip_specs: true`.

## Impact

Only the Vast.ai manifest, lock, documentation, and component-local change record
are in scope. No resource provisioning, host configuration, or alert dismissal.
Acceptance requires secure resolved dependencies, working offline CLI commands,
and repository quality and affected-target lint checks.
