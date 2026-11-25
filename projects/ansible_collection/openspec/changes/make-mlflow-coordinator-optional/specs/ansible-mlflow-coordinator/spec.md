## ADDED Requirements

### Requirement: Make the MLflow coordinator optional

The MLflow role SHALL support installing tracking with the job coordinator disabled while preserving the existing coordinator-enabled default for role callers.

#### Scenario: Tracking without a coordinator

- **WHEN** A caller disables the coordinator
- **THEN** The role installs tracking without evaluating worker credentials or starting the coordinator
