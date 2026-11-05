## Why

A successful Git push can leave GitHub's PR object reporting the previous head.
AL substitutes the pushed head into its expected PR metadata, misclassifies the
old observed head as a human edit, and then refuses receipt-bound recovery.

## What Changes

- Preserve the observed pre-push PR metadata and allow only its exact transition
  to the already-pushed candidate during a guarded update.
- Add `delivery recover-pull-request --receipt-file` to synchronize the PR for
  an exactly receipt-bound, already-pushed commit, even after local work advances.
- Retain partial-push identity on update failure and route the delivery skill to
  explicit recovery rather than repeated pushes or waiting indefinitely.

## Capabilities

### New Capabilities

- `delivery-publication-recovery`: Recover an exact interrupted PR projection
  without changing Git refs or overwriting independently edited metadata.

### Modified Capabilities

None.

## Impact

AL delivery implementation, CLI, isolated Git/provider fixtures, delivery
reference, and canonical delivery skill. No new dependencies, host configuration,
merge, deployment, or renderer source changes.
