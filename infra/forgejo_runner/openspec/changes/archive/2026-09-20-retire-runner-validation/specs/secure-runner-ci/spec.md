## REMOVED Requirements

### Requirement: Signed CI identity

**Reason**: The completed smoke workflow and its validation branch are removed by request.
**Migration**: The default-branch CI identity requirement retains production authentication constraints without an acceptance-test job or ref.

## ADDED Requirements

### Requirement: Default-branch CI identity

Vault SHALL verify the Forgejo Actions issuer, audience, repository, branch,
workflow and event claims. Authentication SHALL accept only the IaC-protected
default branch. CI tokens SHALL expire within five minutes and permit only
self-lookup and self-revocation. Normal repository CI SHALL build and test the
repository without running runner infrastructure acceptance probes.

#### Scenario: Authorized CI identity

- **WHEN** a job on the accepted default branch requests the configured Vault identity
- **THEN** Vault permits only the restricted short-lived token
- **AND** normal build/test orchestration does not request this token merely to test runner infrastructure

#### Scenario: Unauthorized identity

- **WHEN** a job presents a different repository, ref, workflow, audience or unsupported event
- **THEN** Vault denies authentication
