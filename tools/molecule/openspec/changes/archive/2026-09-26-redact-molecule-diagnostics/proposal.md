## Why

PR #112 review identified rendered Ansible task names in retained diagnostics.
Templated names can contain sensitive values even when the task uses `no_log`.

## What Changes

- Retain generated task identifiers instead of rendered task names in events.
- Keep only strictly matched recap counters from private Ansible logs.
- Exercise templated play, task, and handler headings with a public sentinel
  in the real VM fixture and reject that sentinel in retained artifacts.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. This corrects implementation of the existing credential-isolation and
sanitized-evidence requirements; it does not change those guarantees.

## Impact

The Molecule callback, diagnostic collector, smoke fixture, acceptance assertions,
and README change. No host configuration, dependency, or VM lifecycle changes.
