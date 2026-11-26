## Why

Give the Airflow deployment its own Vault identity and DNS zone contributions.

## What Changes

- Establish the component's Bazel and OpenSpec workspace.
- Declare split DNS and a component-owned Terraform state for DNS.
- Declare the AppRole, SSH and ACME roles, browser OIDC client, and shared Vault memberships.

## Capabilities

### New Capabilities

- `airflow-infrastructure`: Provision the Airflow identity and DNS prerequisites.

### Modified Capabilities

None.

## Impact

Airflow's AppRole and DNS roots and the shared Vault membership configuration.
