# Historical project-agents delta

These requirements describe the archived change at its recorded acceptance point. They were imported without applying them to the current `project-agents` baseline. Current source and baseline requirements own present behavior.

## ADDED Requirements

### Requirement: Discoverable continuation

The maintained work discovery surface MUST expose each open objective and its current attempt, stable defect, next action, and resume condition without requiring prior record-path knowledge.

#### Scenario: A fresh agent discovers open work

- **WHEN** the agent requests a bounded continuation view from the repository root
- **THEN** the view identifies the work and its explicit next action with source references

### Requirement: Candidate and review joins

Delivery MUST verify exact candidate validation before publication, and review evidence MUST retain traversable work, delivery, defect, fix, and regression references.

#### Scenario: A caller supplies an unvalidated head

- **WHEN** publication receives a head without matching candidate-bound validation
- **THEN** the operation refuses and does not treat the caller assertion as evidence

### Requirement: Immutable release identity

The version/channel handoff, bundle head, immutable artifact names, remote refs, and release manifest MUST agree; a published immutable release tag MUST never move.

#### Scenario: A release tag already exists

- **WHEN** a reviewed release plan targets an existing immutable tag at different bytes
- **THEN** the guarded publisher refuses to rewrite the tag
