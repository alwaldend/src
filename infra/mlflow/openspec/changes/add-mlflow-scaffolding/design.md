## Context

The MLflow deployment will use a standalone Vault AppRole and separate
Terraform state for DNS. This change establishes those contracts first.

## Decisions

- Use the existing standalone AppRole module and shared Vault groups.
- Own public and DC1 DNS declarations in the MLflow project. The public name
  targets shared ingress and the DC1 records identify the intended guest.
