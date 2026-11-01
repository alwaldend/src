## Why

The standard T3 Code release provides a standalone Linux distribution. Its
existing shared role only supports npm, whose current launcher requires a
newer Node version than the local inference VM's base distribution supplies.

## What Changes

- Add an opt-in checksum-verified standalone installation mode.
- Preserve the existing npm mode and ordinary service account behavior.
- Run the preinstalled standalone executable directly and restart on changes.
- Protect authentication environment writes from Ansible logs.

## Capabilities

### New Capabilities

- `t3code-standalone`: Install the official T3 Code distribution without npm.

### Modified Capabilities

None.

## Impact

The T3 Code role and its local-inference caller. No host deployment is run.
