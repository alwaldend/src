## ADDED Requirements

### Requirement: Shared Grafana deployment memberships

The shared Vault root SHALL resolve the Grafana AppRole by name and register
it with the existing AppRole and Ansible groups.

#### Scenario: Reconcile the shared Vault groups

- **WHEN** the component-local Grafana AppRole exists
- **THEN** the shared Vault root can assign its entity and group to the deployment memberships
