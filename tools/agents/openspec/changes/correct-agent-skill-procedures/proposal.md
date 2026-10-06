## Why

The agent-instruction audit found six procedural defects: destructive Android
troubleshooting, host edits outside the declarative deployment flow, invalid
Bazel runner examples, redundant OpenSpec approval, an incorrect blog output
path, and inaccurate Hugo draft guidance.

## What Changes

- Diagnose Android installation failures before proposing data removal.
- Route persistent Codex configuration changes through owning Ansible state.
- Correct the runner syntax in Bazel and Gazelle skills.
- Preserve implementation authority across required OpenSpec planning.
- Align blog output paths and Hugo draft guidance with the build declaration.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. These documentation corrections restore existing policy and executable
contracts; `skip_specs: true` avoids inventing new requirements.

## Impact

Canonical skills under tools/agents, tools/android, tools/openspec, and
projects/alwaldend.com. No upstream projections, dependencies, live host
configuration, application data, or website deployment change.
