## Why

The rental lifecycle command can cancel its packaged launcher while leaving the
actual Terraform child running and holding its output pipe. That prevents cleanup.

## What Changes

- Propagate SIGINT and SIGTERM through a launcher-owned cancellation context.
- Interrupt and wait for Terraform before returning to the caller.
- Verify interruption across the actual nested launcher with a process fixture.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `terraform-execution`: Graceful cancellation reaches the Terraform child.

## Impact

The standalone launcher's Go command and its infrastructure consumer. No new
dependencies, provider changes, or live infrastructure operations.
