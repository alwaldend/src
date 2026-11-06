## Why

Consumers need a bounded template lookup with validated identity and an explicit
template selector when renting an offer.

## What Changes

- Read one template by exact hash through the pinned generated search operation.
- Report whether its Docker options contain only literal environment assignments.
- Pass an optional template hash through the existing generated create operation.

## Capabilities

### New Capabilities

- `template-selection`: Validated template reads and instance creation selection.

### Modified Capabilities

None.

## Impact

Shared client and documentation, verified through provider Terraform CLI fixtures.
No live API operations or new dependencies.
