## Why

Make the MLflow coordinator optional so browser tracking can start ahead of Vast.ai enrollment.

## What Changes

Separate tracking authentication from worker credentials.

## Capabilities

### New Capabilities

- `ansible-mlflow-coordinator` defines independent tracking startup.

## Impact

Changes MLflow Ansible deployment entry points and role configuration.
