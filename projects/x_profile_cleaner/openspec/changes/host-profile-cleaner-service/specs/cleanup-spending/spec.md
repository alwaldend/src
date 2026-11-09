## Purpose

Bound shared X API expenditure with durable account and cleanup allowances while
letting users choose stricter budgets and protecting the owner's separate access.

## ADDED Requirements

### Requirement: Shared service and personal budgets

The default service limits SHALL be $5 USD per UTC calendar month shared by all
non-owner accounts, $0.25 per non-owner account per month, and $0.10 per non-owner
cleanup. Users SHALL be able to choose lower cleanup and personal monthly budgets.
The configured stable owner X ID SHALL be exempt from service monetary caps but
SHALL remain subject to explicitly chosen personal budgets and resource limits.
The effective allowance SHALL enforce every applicable remaining budget together.
Before identity is known, initial login lookup costs SHALL consume the public pool
under bounded anonymous admission rather than an unverifiable account/run allowance.
After identity verification, account and cleanup ceilings SHALL apply normally.

#### Scenario: Pool exhausted and owner continues

- **WHEN** non-owner usage exhausts the shared pool
- **THEN** further paid non-owner requests are denied while the verified owner remains eligible
- **AND** public requests cannot claim owner status through a supplied ID or handle

#### Scenario: User selects a smaller budget

- **WHEN** the user's chosen allowance is smaller than service allowances
- **THEN** requests are bounded by that allowance across reads, apply, and resume
- **AND** changing a budget never erases committed spend or reservations

#### Scenario: Login before account identity is established

- **WHEN** an initial OAuth flow needs a paid lookup to discover the X identity
- **THEN** its cost is admitted against the remaining public pool and anonymous rate controls
- **AND** owner exemption requires verified identity or protected operator recovery admission

### Requirement: Durable request reservations

Every potentially billable X request SHALL reserve a conservative upper-bound cost
against all applicable budgets atomically before dispatch. This includes identity,
pagination, expanded resources, mutations, and explicit reattempts. Unknown charges
SHALL remain reserved; requests proven undispatched can release reservations. Budget
or storage failure SHALL prevent dispatch. Request-window pacing SHALL remain a
separate limit and SHALL NOT imply available monetary allowance.

#### Scenario: Two requests compete for the last allowance

- **WHEN** concurrent workers request costs whose sum exceeds the remaining pool
- **THEN** reservation admits only requests whose combined exposure fits all applicable limits

#### Scenario: Lost response and process crash

- **WHEN** a request might have reached X and no trustworthy settlement is available
- **THEN** its reservation survives restart without an automatic refund

### Requirement: Budgets survive lifecycle changes

A cleanup's accumulated spending SHALL include discovery, evaluation-related paid
work if any, apply, cancellation, and resume. Cached evaluation SHALL perform no
paid work. New cleanups, login, reauthorization, account recreation, and restart
SHALL NOT reset account or shared allowance. Rollover SHALL preserve unresolved
exposure and re-admit undispatched work against the dispatch period.

#### Scenario: Reconnect and recreate an account

- **WHEN** a user disconnects or recreates their service account with the same X ID
- **THEN** their applicable usage and outstanding reservations remain enforced

#### Scenario: Queued request crosses month end

- **WHEN** work waits through a UTC month boundary before dispatch
- **THEN** admission uses the new period and retained unresolved exposure without double spending

### Requirement: Explain costs before spending

The interface SHALL show selected budgets, snapshot acquisition cost, estimated
selected-action cost, spent and reserved amounts, and remaining allowance. Refresh
SHALL require explicit spending intent and SHALL reserve for maximum page resources.
Apply SHALL confirm an affordable exact selection before mutation. The service SHALL
NOT silently fetch missing data, assume deduplication, or treat unknown tariffs as free.
Confirmation SHALL reserve the selected batch so concurrent work cannot consume
its allowance; per-request reservations SHALL allocate from that hold without
double counting.

#### Scenario: Reads would leave insufficient deletion headroom

- **WHEN** a requested refresh cannot fit the user's total budget and chosen deletion headroom
- **THEN** the service refuses additional spending and offers a narrower window or archive import

#### Scenario: Final selection exceeds the remaining allowance

- **WHEN** selected deletions cost more than the remaining effective allowance
- **THEN** apply is rejected before mutation and the user can reduce selection or explicitly revise a permitted budget
