## Purpose

Manage Vast.ai account and team resource state and verify it without live
credentials, compute rental, or credit transfer.

## ADDED Requirements

### Requirement: Stable launch-mode state

Refresh SHALL normalize supported compound API launch modes before saving state.
Equivalent API modes SHALL retain a compatible configured value; imports SHALL
select a deterministic supported mode. Unsupported or absent modes SHALL report
a read error and preserve existing resource ownership.

#### Scenario: Equivalent launch representation

- **WHEN** the API returns compound access modes equivalent to the configured launch mode
- **THEN** refreshing and planning unchanged configuration does not replace the instance.

#### Scenario: Imported launch mode

- **WHEN** an instance with compound launch modes is imported
- **THEN** state uses a supported configuration value while preserving genuine launch-family drift.

### Requirement: Non-compute lifecycle

The provider SHALL create, refresh, update, import, and delete SSH public keys,
environment variables, and team roles. Confirmed missing objects SHALL be
removed from state; malformed responses and API failures SHALL preserve it.

#### Scenario: Multiple resource types

- **WHEN** Terraform applies, updates, imports, and destroys fixture resources
- **THEN** all three resource types follow their API lifecycle and produce repeatable evidence without live infrastructure.

#### Scenario: Encoded role identity

- **WHEN** a role name contains a forward slash
- **THEN** refresh, import, and deletion address that name as one escaped API path segment.

### Requirement: Team and API-key resources

The provider SHALL create teams with zero credit transfer, generate an
independent management key without changing the caller's credential context,
and expose generated project keys only through sensitive attributes.

#### Scenario: Generated credentials

- **WHEN** the fixture creates a team and API key
- **THEN** generated values are retained in sensitive state and hidden from normal CLI output.

### Requirement: Restricted test project

The test project SHALL use its own AppRole, generated-token injection, isolated
state, and clearly identified test_ objects (TEST_ environment names) in its
team. It SHALL contain no compute rental, volume allocation, credit transfer,
or membership invitation.

#### Scenario: Team validation

- **WHEN** an explicitly authorized live test-root operation runs after bootstrap
- **THEN** only its SSH public key, non-secret marker variable, and read-only role are managed, without renting compute.

### Requirement: Compound team ownership

The provider SHALL track confirmed team creation independently of management-key
ownership. Cleanup after failed creation SHALL revoke only owned credentials
and SHALL NOT adopt or destroy pre-existing teams. Refresh SHALL preserve key
cleanup ownership after external team deletion and plan explicit replacement.

#### Scenario: Duplicate team name

- **WHEN** creation fails because the requested team name already exists
- **THEN** refresh and destroy leave that team intact and revoke only the newly generated management key.

#### Scenario: External team deletion

- **WHEN** a created team disappears while its management key remains valid
- **THEN** refresh performs no mutation, replacement or destroy revokes the original key, and recreation does not abandon a live credential.
