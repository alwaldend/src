## Why

Make the MLflow service deployment reusable through the Ansible collection.

## What Changes

- Add the `alwaldend.main.mlflow` role for MLflow and its Vast.ai coordinator.
- Accept a caller-supplied runtime archive and configuration.

## Capabilities

### New Capabilities

- `ansible-mlflow` installs and manages the MLflow service and coordinator.

## Impact

The role adds source to the Ansible collection. Infrastructure callers supply
the archive, credentials, and reverse proxy configuration.
