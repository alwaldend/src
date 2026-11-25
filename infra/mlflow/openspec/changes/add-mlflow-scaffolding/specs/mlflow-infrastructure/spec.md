## ADDED Requirements

### Requirement: MLflow infrastructure identity and DNS

The repository SHALL declare a standalone MLflow Vault AppRole and a
component-owned DNS Terraform root.

#### Scenario: An operator bootstraps the component

- **WHEN** the AppRole and shared Vault memberships are applied
- **THEN** the component identity can manage its DNS records through its own
  Vault-backed Terraform state
