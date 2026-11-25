## Why

Establish the MLflow infrastructure project so its identity and DNS can be
reviewed independently of the service deployment.

## What Changes

- Declare the component AppRole, shared Vault memberships, and DNS records.
- Add the project documentation and OpenSpec workspace.

## Capabilities

### New Capabilities

- `mlflow-infrastructure` defines the component's infrastructure identity and DNS.

## Impact

Vault and DNS gain an MLflow component. The server and its dependencies are
introduced by the service deployment change.
