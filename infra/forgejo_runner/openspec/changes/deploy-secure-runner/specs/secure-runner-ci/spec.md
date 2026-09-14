## Purpose

Verify the secure runner's resources and authenticate CI through a narrowly
bound, short-lived Vault identity with explicit branch trust limits.

## ADDED Requirements

### Requirement: Signed CI identity

Vault SHALL verify the Forgejo Actions issuer, audience, repository, branch,
workflow and event claims. Authentication SHALL accept only the IaC-protected
default branch and the exact declared protected validation branch. CI tokens
SHALL expire within five minutes and permit only self-lookup and self-revocation.

#### Scenario: Authorized smoke job

- **WHEN** the smoke workflow runs on an accepted protected branch
- **THEN** it obtains and inspects a short-lived Vault token and revokes it
- **AND** verifies the requested VM resources without printing credentials

#### Scenario: Unauthorized identity

- **WHEN** a job presents a different repository, ref, workflow, audience or
  unsupported event
- **THEN** Vault denies authentication

### Requirement: Honest protected-branch boundary

IaC and documentation SHALL distinguish workflow branch selection from
server-enforced runner admission. The deployment SHALL scope registration to
one repository and document the pinned Forgejo version's missing protected-only
runner setting and unusable `ref_protected` claim.

#### Scenario: Inspect enforcement

- **WHEN** a maintainer reviews the security boundary
- **THEN** the supported guarantee is signed-claim restriction at Vault
- **AND** workflow filtering is not represented as preventing another branch
  from targeting the runner
