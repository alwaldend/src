## Why

Binary locks can silently discard execution constraints because obsolete and unknown fields are accepted.

## What Changes

- Reject unknown fields in owned lock structures, including obsolete archive rules metadata.
- Preserve execution constraints under the supported toolchain field.
- Migrate repository consumers without changing binary versions or archive hashes.

## Capabilities

### New Capabilities

- `binary-lock-schema`: Strict lock structure and execution-platform constraints.

### Modified Capabilities

None.

## Impact

Binary toolchain lock parsing and repository-owned locks. Coordinated with infra/src/openspec/changes/repair-repository-organization.
