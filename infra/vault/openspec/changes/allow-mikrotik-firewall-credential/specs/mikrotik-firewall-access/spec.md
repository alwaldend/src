## Purpose

Allow the MikroTik firewall workflow to reuse a named existing RouterOS credential
while retaining its own Vault identity and limiting the new secret-access grant.

## ADDED Requirements

### Requirement: Exact-secret firewall credential access

The MikroTik AppRole SHALL receive read-only access to
`secrets/data/alwaldend.com/vault1/approles/src_infra_ingress/mikrotik`.
The grant SHALL not include other ingress secrets, wildcard paths, or write
capabilities. Existing DNS, state, and bootstrap access SHALL be retained.
The existing DNS policy SHALL additionally grant read on the named Cloudflare
DNS secret for the global view, retaining the dc1 grant.

#### Scenario: Obtain the existing RouterOS account

- **WHEN** the authorized policy has been deployed and the MikroTik workflow authenticates
- **THEN** it can read the named RouterOS secret using its own AppRole
- **AND** no credential copy or new RouterOS account is required

#### Scenario: Preserve unrelated secret boundaries

- **WHEN** the additional policy is inspected
- **THEN** it grants only read on the exact named secret
- **AND** it grants no general ingress identity, policy management, or secret writes

#### Scenario: Publish the printer globally

- **WHEN** the owner selects both dc1 and global DNS views
- **THEN** its existing DNS policy grants named provider-secret reads for both views
- **AND** its existing state and bootstrap access remain intact
