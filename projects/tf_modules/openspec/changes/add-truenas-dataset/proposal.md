## Why

Dataset lifecycle and inheritance should be reusable while each consuming
infrastructure owner declares its own dataset beneath an explicit parent.

## What Changes

- Add a TrueNAS filesystem module with parent and inline config inputs.
- Own shared filesystem defaults, private account ownership and optional authenticated DEFAULT_SHARE exports.
- Inherit encryption, protect deletion and return dataset identity/mountpoint.

## Capabilities

### New Capabilities

- `truenas-dataset`: Protected filesystem datasets under an explicit parent.

### Modified Capabilities

None.

## Impact

Uses the existing pinned TrueNAS provider. NAS and runner roots select their own
provider credentials and consume the module without JSON configuration.
