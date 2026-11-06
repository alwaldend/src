## Purpose

Generate an independent credential for an existing owned team while preserving
the shared issuer and the team's resources, credit, and ownership.

## ADDED Requirements

### Requirement: Bind only the generated credential

An API-key resource configured with an existing team ID SHALL resolve its owner
context, generate its independent credential, and verify that credential selects
the requested team. It MUST NOT change the issuer's context or create a team.
This owner-binding mode MUST require explicit default full access and reject
scoped permission input before making API requests.

An unresolved computed team ID MAY defer that check until team identity is
read. Once resolved, scoped owner binding MUST be rejected before key creation.

#### Scenario: Bind an owned team

- **WHEN** the issuer owns the requested team and activation succeeds
- **THEN** only the generated credential selects that team
- **AND** normal destruction revokes only that generated credential

#### Scenario: Scoped input cannot widen access

- **WHEN** an existing-team owner binding is combined with scoped permissions
- **THEN** configuration is rejected before any API request

### Requirement: Preserve failed activation ownership

The provider SHALL retain a generated key's identity and sensitive value in
state before activation. Failed activation MUST fail apply and remain
recoverable through an explicit resource destruction.

#### Scenario: Activation fails

- **WHEN** the API rejects activation after returning a generated key
- **THEN** apply fails and retains that key in Terraform state
- **AND** destruction revokes the key without changing the existing team

### Requirement: Provider-owned issuer authentication

API-key resources SHALL use their configured provider for lifecycle operations
when no bootstrap override is supplied. Removing a redundant override SHALL
verify visibility through the provider and preserve the generated key's identity
and value. A generated same-apply management key MAY override the issuer. A new
explicit bootstrap issuer SHALL retain replacement behavior.

#### Scenario: Remove a redundant issuer input

- **WHEN** the resource credential is removed and its provider can access the key
- **THEN** Terraform preserves the key without creation or revocation
- **AND** subsequent reads and destruction use the provider credential
