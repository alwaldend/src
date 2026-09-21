# secure-runner-ci Specification

## Purpose

Define the secure runner's narrowly bound, short-lived Vault CI identity and
explicit branch trust limits independently of repository build/test checks.

## Requirements

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
