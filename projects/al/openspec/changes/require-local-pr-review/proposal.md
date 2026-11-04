## Why

Delivery already requires general correctness scrutiny and guarded remote review
handling, but does not explicitly require a local code review for each PR
publication. Make both review stages mandatory for initial and updated candidates.

## What Changes

- Require the `review-agent` workflow on the complete merge-base-to-candidate
  diff before each PR creation or update.
- Recheck a changed candidate after local or remote review fixes, and retain
  guarded remote feedback resolution after every publication.
- Extend delivery evaluation scenarios to cover both review stages.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `review-delivery`: require local candidate review and remote feedback handling
  throughout the delivery cycle.

## Impact

Changes AL's delivery skill, its generated discovery projection, and evaluation
cases. No CLI, dependency, host configuration, or deployment changes.
