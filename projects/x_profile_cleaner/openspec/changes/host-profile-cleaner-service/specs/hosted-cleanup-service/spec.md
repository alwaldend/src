## Purpose

Provide an account-isolated web cleanup workflow with durable review and recovery,
while retaining the existing planner and exact-action execution safety contract.

## ADDED Requirements

### Requirement: Hosted account boundary

The service SHALL permit registration through verified X login and SHALL isolate
policies, snapshots, plans, jobs, selections, progress, and logs by authenticated
account. User-supplied identifiers SHALL NOT establish ownership. The hosted mode
SHALL NOT substitute the operator's personal X credential for a user's permission.

#### Scenario: Another account's operation identifier

- **WHEN** a logged-in user requests another account's plan, job, or log stream
- **THEN** the service denies access without disclosing its content or dispatching X requests

### Requirement: Server-rendered review and actions

The browser SHALL receive server-rendered pages and bounded interactive fragments
for policy editing, snapshot selection/import/refresh, budget selection, evaluation,
review, apply, cancellation, and resume. Post and policy text SHALL be escaped data.
State-changing requests SHALL require authenticated ownership and CSRF protection.

#### Scenario: Malicious post text and forged action

- **WHEN** a post contains script markup and a foreign origin submits an apply request
- **THEN** the post renders as text and the forged request is rejected before mutation

### Requirement: Immutable confirmed execution

Evaluation SHALL preserve the baseline cleanup selection and preservation semantics
and SHALL never mutate X. Each evaluation SHALL create a new server-owned plan
bound to its account, snapshot, policy, and evaluation time. Apply SHALL atomically
consume a valid current plan and accept only unique action indexes from it. It
SHALL execute exactly the confirmed subset, without reevaluation, recollection,
fabricated targets, deletion cascades, or implicit expansion.

#### Scenario: Duplicate apply and modified editor

- **WHEN** two apply requests target the same plan or the editor changes after confirmation
- **THEN** only one execution is admitted and its actions remain the original confirmed subset

### Requirement: Durable explicit recovery

The service SHALL persist execution selection, dispatch intent, confirmed outcomes,
and uncertainty. Restart SHALL interrupt jobs without automatically dispatching
mutations. Resume SHALL require explicit user action, matching X authorization,
valid budgets, and the original selection. Confirmed completions SHALL NOT replay;
unknown outcomes SHALL be shown before explicit reattempt. Cancelling SHALL stop
new dispatch without claiming already accepted X operations were reversed.

#### Scenario: Crash after dispatch and before receipt

- **WHEN** the service restarts after dispatching an action without recording its response
- **THEN** it retains that action as unknown, preserves its spending exposure, and requires explicit resume
- **AND** a later confirmed absent-target response can complete the action without adding other actions

### Requirement: Bounded retention and account deletion

The service SHALL publish and enforce retention and aggregate storage limits for
account data, uploads, plans, execution records, and logs. Explicit account deletion
SHALL cancel work and remove retained user content. It SHALL NOT reset current or
unresolved spending exposure; minimal ledger and backup retention SHALL be disclosed.

#### Scenario: Storage full of retained execution data

- **WHEN** storage cannot admit another snapshot without discarding protected execution progress
- **THEN** new work is rejected and existing progress and budget accounting remain intact

### Requirement: Recoverable single-VM operation

The service SHALL operate on one VM with HTTPS ingress, supervised backend and
isolated workers, durable local state, bounded logs, health reporting, and consistent
backups. Deployment SHALL be declarative and SHALL require its separately authorized
environment. Backup restoration SHALL NOT silently restore spendable public quota.

#### Scenario: Restore an older database

- **WHEN** an operator restores state predating dispatched paid requests
- **THEN** public spending remains closed until exposure is reconciled or conservatively accounted
- **AND** no X authorization or automatic mutation is restored
