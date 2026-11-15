## Purpose

Allow the development deployment to create its secure host using its own scoped virtualization identity.

## ADDED Requirements

### Requirement: Development virtualization allocation

The src_infra_dev resource set SHALL contain the existing Fedora 44 cloud
template, local storage and wired pool network. Its subject SHALL be the exact
OIDC user for the development Vault entity, following the existing owner model.

#### Scenario: Create secure1

- **WHEN** the development OIDC identity is bootstrapped and its resource set is applied
- **THEN** infra/dev can create secure1 using its own identity and invocation-scoped XO token
- **AND** the identity is not a global XO administrator

### Requirement: Verified XO certificate prerequisite

The certificate deployment SHALL verify an existing managed certificate against
the Vault CA and every configured hostname. Missing or invalid certificates
SHALL trigger issuance and management-service reload without stopping guest VMs.

#### Scenario: Managed file contains a self-signed placeholder

- **WHEN** certificate reconciliation finds a hostname-less self-signed managed certificate
- **THEN** it obtains a Vault certificate and verifies HTTPS after restarting XO
