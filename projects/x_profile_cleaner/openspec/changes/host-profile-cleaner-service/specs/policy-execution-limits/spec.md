## Purpose

Keep hostile or accidentally expensive policies and archives from exhausting the
shared web service while preserving failure without a partial executable plan.

## ADDED Requirements

### Requirement: Bounded untrusted inputs and work

The service SHALL enforce configurable limits before and during parsing and
evaluation for policy bytes, request bodies, archive bytes and entries, records,
rule registrations, explicit identifiers, callbacks, graph visits, and output.
All bounds SHALL apply to the owner as well. Syntax validation SHALL NOT be treated
as proof of bounded execution. Limit violations SHALL fail the operation without
publishing a partial snapshot or executable plan or dispatching a mutation.

#### Scenario: Small source produces excessive allocation

- **WHEN** a syntactically valid policy allocates or traverses beyond its allowance
- **THEN** the operation fails within its resource limits and publishes no executable partial plan

#### Scenario: Compressed input exceeds its decoded allowance

- **WHEN** an archive is within its upload bound but exceeds entry or decoded-byte limits
- **THEN** parsing stops before exhausting shared storage or backend memory

### Requirement: Worker failure isolation

Untrusted parsing and evaluation SHALL execute under separate hard memory, CPU,
task, and elapsed-time controls. Workers SHALL receive no X credentials or general
network or host-file access. A worker crash or OOM SHALL NOT terminate the web
backend or grant its output authority without validation. Cancellation SHALL end
the worker and release its slot after bounded cleanup.

#### Scenario: Worker exceeds its memory limit

- **WHEN** a policy worker is killed for memory exhaustion
- **THEN** login, status, and unrelated bounded jobs remain available
- **AND** the failed operation reports resource exhaustion without mutating X

### Requirement: Bounded fair job admission

The service SHALL enforce one active account operation, bounded account/global
waiting queues, bounded global worker count, and request-rate admission independent
of monetary budgets. Owner exemption SHALL NOT bypass these controls. Excess work
SHALL receive an explicit capacity response without unbounded goroutine or queue growth.

#### Scenario: Repeated cached evaluations attempt to monopolize the VM

- **WHEN** one account submits more work than its active and queued allowance
- **THEN** excess submissions are rejected while other accounts retain admission opportunities
