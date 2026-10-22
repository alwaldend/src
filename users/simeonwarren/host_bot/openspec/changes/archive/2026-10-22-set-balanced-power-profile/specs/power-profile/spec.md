## Purpose

Keep the host's selected power profile reproducible through its owning
Ansible configuration while retaining an explicit desired profile.

## ADDED Requirements

### Requirement: Managed host power profile

The host's Ansible configuration SHALL declare `balanced` as its desired
Tuned profile. Deployment SHALL enable and start the existing Tuned service
and persist the declared profile with manual selection. Deployment SHALL
avoid reselecting an already-correct profile in manual mode and SHALL
report unexpected profile-query failures.

#### Scenario: Selected profile differs

- **WHEN** the active profile differs from the declared profile
- **THEN** deployment selects and persists the declared profile

#### Scenario: Automatic selection is enabled

- **WHEN** the active profile matches but selection is automatic
- **THEN** deployment persists the declared profile with manual selection

#### Scenario: No profile is active

- **WHEN** Tuned reports that no profile is active
- **THEN** deployment selects the declared profile instead of failing on that expected query result

#### Scenario: Desired state already matches

- **WHEN** the active profile matches and selection is manual
- **THEN** deployment does not repeat the profile-selection command

#### Scenario: Query fails unexpectedly

- **WHEN** querying the active profile or selection mode fails unexpectedly
- **THEN** deployment reports the failure instead of concealing it
