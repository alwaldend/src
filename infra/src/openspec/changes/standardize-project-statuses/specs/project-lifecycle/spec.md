## Purpose

Provide consistent, reviewable lifecycle information for every first-party
project without making release, support, or deployment guarantees.

## ADDED Requirements

### Requirement: Canonical project lifecycle

Every direct project README SHALL declare exactly one lifecycle status from
the shared catalog. The catalog SHALL define active, maintenance, completed,
and archived. Active SHALL mean current substantive feature development.
Historical development, unfinished plans, or routine upkeep alone SHALL NOT
establish active status. Maintenance SHALL cover retained projects without
current feature work, including projects with unfinished scope. Lifecycle SHALL
remain independent of release readiness.

#### Scenario: Invalid declaration

- **WHEN** a project omits a status, declares multiple statuses, or uses an
  unknown value
- **THEN** inventory generation and freshness checks fail with project context

### Requirement: Derived public status

The project index and landing status SHALL derive from README declarations.
The inventory freshness check SHALL reject stale landing status. Generation
SHALL preserve unrelated landing metadata and content. Previous taxonomy URLs
SHALL redirect to their replacement terms.

#### Scenario: Status changes

- **WHEN** a README status changes and inventory update runs
- **THEN** the project index and existing landing status reflect the new value
- **AND** a second generation produces identical bytes

#### Scenario: Stale projection

- **WHEN** landing status differs from its README
- **THEN** freshness validation fails without changing source files

### Requirement: Evidence-based migration

The migration SHALL review every direct project using repository descriptions
and implementation history, record its disposition, and preserve known limits.

#### Scenario: Historical feature work without current development

- **WHEN** historical feature additions exist but the owner confirms no current
  feature development
- **THEN** the review records maintenance rather than inferring active status
  from the history or unfinished plans
