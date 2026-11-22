## Why

MLflow needs a way to execute projects on existing Vast.ai rentals. Job
execution and result tracking must survive the submitting MLflow process and
temporary loss of SSH connectivity.

## What Changes

- Define a Python MLflow adapter backed by a Go job server.
- Define server, worker service, and worker client modes of one Go binary.
- Define gRPC interfaces, SQLite persistence, and job reconciliation across
  the SSH boundary.
- Require all server and worker business logic to run through Stateless state
  machines.

## Capabilities

### New Capabilities

- `vastai-project-execution` for submitting and managing MLflow jobs on
  Terraform-managed Vast.ai rentals.

## Impact

The implementation will affect the project Python package, Go binary, and
protocol definitions. Terraform remains the operator's rental control plane.
