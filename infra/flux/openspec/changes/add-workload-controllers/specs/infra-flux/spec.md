## ADDED Requirements

### Requirement: Terraform ingress and certificate controllers

Flux SHALL install Traefik and cert-manager as Terraform-managed Helm releases
using pinned charts and images mirrored in Forgejo.

#### Scenario: Apply the workload root

- **WHEN** an operator applies the workload Terraform root
- **THEN** Terraform installs Traefik with a cluster-local Service and cert-manager with its CRDs
