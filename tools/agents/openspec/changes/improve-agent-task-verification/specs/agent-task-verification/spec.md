## ADDED Requirements

### Requirement: Verify the reported symptom

Agent procedures SHALL check the specific property under revision through
the component that owns it.

#### Scenario: Select the owner of a service setting

- **WHEN** a comparable service already runs on another host
- **THEN** the agent traces the sibling playbook into its roles and checks which role consumes the setting

#### Scenario: Correct the initial presentation of a job log

- **WHEN** a user reports expanded diagnostics in a Forgejo job
- **THEN** the agent checks the candidate's rendered job or records that only the emitted markers were verified

#### Scenario: Correct a movement at a loop boundary

- **WHEN** a user reports unexpected movement at a loop boundary
- **THEN** the agent inspects the motion across that transition at output speed
