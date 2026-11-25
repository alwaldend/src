## Why

Deploy tracking before worker enrollment so browser tracking can start ahead of Vast.ai enrollment.

## What Changes

Separate tracking authentication from worker credentials.

## Capabilities

### New Capabilities

- `mlflow-tracking-deployment` defines independent tracking startup.

## Impact

Changes MLflow Ansible deployment entry points and role configuration.
