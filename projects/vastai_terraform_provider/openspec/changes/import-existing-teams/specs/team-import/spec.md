## Purpose

Adopt existing owned Vast.ai teams into Terraform without retrieving credentials
or changing live team or account configuration.

## ADDED Requirements

### Requirement: Read-only team adoption

The provider SHALL import an owned team by canonical positive decimal account
ID using read-only API calls. It SHALL retain the observed name and account ID
without generating credentials or changing any credential's context.
Unrecoverable management credential fields SHALL remain null.

#### Scenario: Import an existing owned team

- **WHEN** an owned team is imported and configured with its observed name
- **THEN** a subsequent plan reports no changes and the API receives only reads

#### Scenario: Reject invalid or unavailable identity

- **WHEN** the ID is invalid, not owned, absent, or cannot be verified
- **THEN** import fails without creating state or mutating the account

### Requirement: Preserve imported teams without management credentials

The provider SHALL refuse destruction of an imported team that still exists.
Refresh failures SHALL retain state and SHALL NOT revoke or generate keys.
Terraform-created teams SHALL retain their existing lifecycle behavior.

#### Scenario: Refuse destruction of an imported team

- **WHEN** Terraform attempts to destroy an imported team that still exists
- **THEN** destruction fails with an actionable diagnostic and retains its state
- **AND** no live team or credential is changed

#### Scenario: Observe external deletion

- **WHEN** the imported team's account metadata confirms it is absent
- **THEN** refresh removes the imported identity from state without mutation
