## Why

Recent sessions repeatedly reran passing validation because separate commands inherited different environments (ERG-001). Publication also lost its useful result when concurrent review updates interrupted final verification (ERG-002).

## What Changes

- Add explicit publication to the validation command, retaining all existing guards.
- Return the known candidate and pull request when final verification fails, with an explicit pending status and verification-only recovery.
- Document the combined invocation and distinguish publication evidence from complete review verification.
- Route the delivery skill to the clarified root standing authorization and validate a case that avoids redundant destination approval.

## Capabilities

### New Capabilities

- `delivery-continuation`: Stable validation/publication execution and recoverable publication reporting.

### Modified Capabilities

None.

## Impact

The delivery command, integration fixtures, and owning README. No dependencies or live host changes.
