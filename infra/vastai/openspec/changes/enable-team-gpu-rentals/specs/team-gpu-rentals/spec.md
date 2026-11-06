## Purpose

Rent GPU compute in the funded alwaldend team through declarative infrastructure,
and verify the lifecycle with a minimal-duration paid rental.

## ADDED Requirements

### Requirement: Independent team authentication

Rental commands SHALL use a full-access alwaldend team credential retrieved from
Vault through the default provider's environment variable. They MUST NOT switch
credential context or modify personal credentials or team ownership.

#### Scenario: Configure rental access

- **WHEN** the flat rental root is configured with the dedicated key
- **THEN** its provider uses that team's credential without generating keys,
  switching context, creating teams, or transferring credit

### Requirement: Cheapest on-demand rental

The rental root SHALL select the lowest total hourly price among currently
rentable, unoccupied on-demand GPU offers, using the configured disk allocation.
An empty marketplace SHALL fail before creating a rental. Market refreshes
MUST NOT replace an existing rental merely because a cheaper offer appears.

#### Scenario: Select and retain a rental

- **WHEN** a new rental is applied
- **THEN** the cheapest matching offer is selected, subject to the stated price ceiling
- **AND** later offer changes do not replace the rented contract

### Requirement: Separate account and team ownership

The account root SHALL own team identities and account settings without creating
project API keys. A separate alwaldend root SHALL use an operator-provided team
credential through its provider, create full-access keys in that team, and store
them in Vault. The seed credential MUST remain operator-owned. Each root SHALL
use its own backend; rentals MUST remain in their independent root and state.

#### Scenario: Provision rental access

- **WHEN** the alwaldend root is applied with its verified team provider key
- **THEN** Terraform creates the rental key and Vault entry without creating a team or transferring credit
- **AND** the generated key selects alwaldend before compute is rented

### Requirement: Shared environment defaults

Account and team configuration SHALL read one environment-variable section.
ENABLE_HTTPS, WEB_ENABLE_HTTPS and WEB_ENABLE_AUTH SHALL be true in both scopes.
Settings MUST NOT depend on inheritance between accounts and teams.

#### Scenario: Configure both scopes

- **WHEN** the account and team roots are applied
- **THEN** both receive the same declared environment map from one configuration owner

### Requirement: Retire other teams explicitly

The desired account configuration SHALL retain only alwaldend. Retirement SHALL
clean up test resources through their owning state before revoking credentials.
A failed deletion MUST remain tracked and MUST NOT be reported as completed or
silently detached. Temporary migration blocks SHALL be removed after state
reconciliation is verified.

#### Scenario: A retired team has no management credential

- **WHEN** the provider cannot destroy an imported identity
- **THEN** the team remains tracked and the remaining deletion is reported explicitly
- **AND** the funded alwaldend team is preserved

### Requirement: Explicit rent-and-release test

The flat rentals_test root SHALL provide separate rent and release commands.
The operator SHALL rent the cheapest matching offer and start releasing its
contract immediately after creation, without waiting for boot or SSH readiness.
Creation and release SHALL use reviewed interactive Terraform applies without
saved plans or automatic approval. Cleanup SHALL also be attempted after an
apply failure. State listing and CLI inventory SHALL confirm no test-owned rental
remains; unrelated contracts MUST be preserved.

#### Scenario: Successful lifecycle test

- **WHEN** the operator runs and approves rent followed immediately by release
- **THEN** one contract is rented and immediately destroyed
- **AND** state listing and CLI inventory confirm cleanup

#### Scenario: Creation fails

- **WHEN** the creation apply fails after possibly changing state
- **THEN** cleanup is attempted and the test reports failure rather than retrying rental creation
