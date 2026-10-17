## Why

Collection lifecycle and state archive integration can be reviewed independently of reconciliation. Extract the existing #128 implementation onto master so state commands can use the same native database/package workflow without bringing export or reconciliation logic with them.

## What Changes

- Compose the existing Anki process and collection backends per operation, with bounded cleanup and native package publication for changed archives.
- Route state commands through App.WithCollection and construct namespaced stores directly from its prepared connection, preserving direct database writes and unchanged archive bytes.
- Rename collection input flags to --collection and keep the state namespace in application options.
- Extract lifecycle/state E2E coverage and current documentation; collection export, plan, apply and build remain placeholders in this PR.

## Capabilities

### New Capabilities

- collection-workflow: Per-operation native lifecycle and state commands over offline databases or collection packages.

### Modified Capabilities

None.

## Impact

Only projects/anki_as_code changes. Reuse the merged backend, state, config and codec APIs with no new dependency. This independent master-based PR is extracted from #128; the original retains the overlap until this prerequisite merges.
