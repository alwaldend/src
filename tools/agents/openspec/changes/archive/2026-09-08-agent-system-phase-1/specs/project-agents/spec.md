# Historical project-agents delta

These requirements describe the archived change at its recorded acceptance point. They were imported without applying them to the current `project-agents` baseline. Current source and baseline requirements own present behavior.

## ADDED Requirements

### Requirement: Shared operation contracts

Shared v1alpha1 contracts MUST round-trip deterministically and reject malformed identity, unknown effectful operations, widened authority, and incompatible information flows.

#### Scenario: An invalid effectful declaration is submitted

- **WHEN** a declaration contains an unknown operation or wider authority than its owner grants
- **THEN** contract validation rejects the declaration before it can represent an authorized action

### Requirement: Registered universe reporting

The completeness report MUST derive operation and lifecycle classifications from the closed set of owner-local registrations and identify every missing or unclassified entry.

#### Scenario: A registration is absent or unclassified

- **WHEN** the registered-universe checker evaluates the owner declarations
- **THEN** the report names the missing classification and does not silently omit the entry

### Requirement: Independent information policy

The policy model MUST keep public source, secrets, and personal information independent from target visibility, build consumers, and publication boundaries.

#### Scenario: A repository-internal target is reviewed

- **WHEN** a target is restricted to internal build consumers
- **THEN** that restriction alone does not classify ordinary source as secret or authorize disclosure of credentials

### Requirement: Isolated scratch and publication

Task and run scratch MUST remain namespaced, and expected-revision publication MUST reject a stale writer.

#### Scenario: Two tasks share a workspace

- **WHEN** two tasks execute concurrently and one later submits a stale publication revision
- **THEN** their scratch does not collide and the stale write fails
