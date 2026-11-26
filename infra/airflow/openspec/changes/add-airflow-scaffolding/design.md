## Context

The Airflow guest and shared ingress consume a component-owned DNS declaration.
The AppRole root follows the MLflow bootstrap model.

## Goals / Non-Goals

- Make the component's identity and DNS available independently of the workload.

## Decisions

Use separate Terraform roots and state for the Vault bootstrap and DNS. Register
the browser client against the Airflow hostname and assign it to the global
administrators group.

## Risks / Trade-offs

The names can be published before deployment. The guest and ingress must be
ready before the published routes serve requests.
