## Why

Provide a shared MLflow tracking and artifact service for projects executed
on existing Vast.ai rentals.

## What Changes

- Provision a standard infrastructure VM and authenticated HTTPS tracking endpoint.
- Deploy the pinned MLflow CLI with the repository's Vast.ai backend and run the
  job coordinator on the same VM.

## Capabilities

### New Capabilities

- `mlflow-deployment` covers VM provisioning, service deployment and access.

## Impact

The change uses the existing MLflow AppRole and DNS records. It adds Vault OIDC for
browser users, workload bearer authentication, and an ingress route.
