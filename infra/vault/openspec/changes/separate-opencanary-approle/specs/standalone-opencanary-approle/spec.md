## ADDED Requirements

### Requirement: Independent OpenCanary AppRole state

The OpenCanary AppRole and its SSH signing policy SHALL have a dedicated Terraform root and backend.

#### Scenario: Resolve shared membership

- **WHEN** the standalone root has provisioned the AppRole
- **THEN** the shared Vault root discovers the entity and group by name for Xen Orchestra and Ansible membership
