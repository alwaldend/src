## ADDED Requirements

### Requirement: Airflow identity and DNS prerequisites

The Airflow infrastructure SHALL declare a component-owned AppRole and DNS root.
Global DNS SHALL target shared ingress for the public name and DC1 DNS SHALL
resolve the same name directly to the guest.

#### Scenario: An operator prepares the Airflow deployment

- **WHEN** the AppRole and DNS roots and the shared Vault memberships are reconciled
- **THEN** the Airflow component has its own Vault identity and DNS access
- **AND** the browser OIDC client is assigned to the global administrators
- **AND** the public and DC1 views resolve the names to their respective targets
