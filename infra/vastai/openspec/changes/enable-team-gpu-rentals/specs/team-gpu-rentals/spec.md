## Purpose

Rent GPU compute in the funded alwaldend team through declarative infrastructure,
and verify the lifecycle with a minimal-duration paid rental.

## ADDED Requirements

### Requirement: Independent team authentication

Rental commands SHALL use a full-access alwaldend team credential retrieved from
Vault with an aliased provider passed into a child module. They MUST NOT switch
credential context or modify personal credentials or team ownership.

#### Scenario: Configure rental access

- **WHEN** the rental module is configured with the dedicated key
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

### Requirement: Declarative full-access team credential

The existing team catalog SHALL create a full-access key through its team
module for the owned alwaldend team and store it directly in Vault. It MUST preserve the
personal issuer's context and verify the resulting team binding before rental.
The root SHALL reuse the shared team module and configure issuer authentication
through its provider, with Vault storage remaining in the caller. Each secret
SHALL use its namespace owner's existing Vault identity. The team and its keys
SHALL share the existing catalog state; no separate credential root is required.

#### Scenario: Provision rental access

- **WHEN** the reconciled team catalog is applied for the owned alwaldend team
- **THEN** Terraform creates its key and Vault entry without creating a team or transferring credit
- **AND** the generated key is verified to select alwaldend before compute is rented

### Requirement: Explicit rent-and-release test

The live test SHALL rent the cheapest matching offer and start releasing its
contract immediately after creation, without waiting for boot or SSH readiness.
Creation and release SHALL use reviewed interactive Terraform applies without
saved plans or automatic approval. Cleanup SHALL also be attempted after an
apply failure. The test SHALL produce a sanitized receipt and confirm no
test-owned rental remains; unrelated contracts MUST be preserved.

#### Scenario: Successful lifecycle test

- **WHEN** the operator runs and approves the explicit live test
- **THEN** one contract is rented and immediately destroyed
- **AND** a receipt records the result and cleanup verification

#### Scenario: Creation fails

- **WHEN** the creation apply fails after possibly changing state
- **THEN** cleanup is attempted and the test reports failure rather than retrying rental creation
