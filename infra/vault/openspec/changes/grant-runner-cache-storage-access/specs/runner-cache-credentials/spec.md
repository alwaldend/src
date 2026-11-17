## Purpose

Publish NAS accounts through existing AppRole credential boundaries without
sharing administrative NAS credentials with service accounts.

## ADDED Requirements

### Requirement: Scoped NAS credential publication

NAS provisioning SHALL be allowed to discover canonical AppRole identities and
publish NAS accounts only below shared/read/approle/+/nas. Existing shared-read
policies SHALL give each AppRole its matching record. XCP-ng SHALL read its own published NAS account through the existing self-scoped
policy, with no added development-account credential grant.

#### Scenario: Inspect account grants

- **WHEN** the source policies are reviewed
- **THEN** publication is restricted to NAS credential records
- **AND** the NAS administrator API credential remains readable only through the NAS provisioning identity

### Requirement: Vault owns NAS access membership

Vault SHALL declare truenas_admins and truenas_user groups. The NAS AppRole SHALL
belong to the administrator group; ordinary AppRoles SHALL belong to the user
group. NAS provisioning SHALL create accounts from those groups' entity members
and mirror group membership rather than deciding access from account names.

#### Scenario: Reconcile membership

- **WHEN** an identity belongs to either NAS access group
- **THEN** NAS provisioning declares its account and matching NAS group membership
- **AND** only administrator membership grants FULL_ADMIN
