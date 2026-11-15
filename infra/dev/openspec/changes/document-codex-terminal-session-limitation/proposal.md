## Why

Repeated Codex terminal-continuation failures can strand running commands and
Terraform locks. Operators need to distinguish the upstream argument-parser
defect from a failure of the command being controlled.

## What Changes

- Document the upstream issue, checked source versions, and failing numeric form
  in the development host README.
- Explain how to pass an integer session ID and inspect the original command
  before retrying an operation or recovering a lock.

## Capabilities

### New Capabilities

None. This documentation-only change declares `skip_specs: true`.

### Modified Capabilities

None.

## Impact

Only development-host documentation and its change record are updated. The
user selected documentation after investigation identified an upstream Codex
defect; executable, provider, authentication, and deployment settings remain
outside this change.
