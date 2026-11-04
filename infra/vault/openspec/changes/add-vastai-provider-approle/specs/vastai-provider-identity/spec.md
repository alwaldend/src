## Purpose

Isolate Vast.ai provider validation credentials and state while allowing its
bootstrap infrastructure to deliver a Terraform-generated project token.

## ADDED Requirements

### Requirement: Isolated project identity

The project SHALL have its own AppRole and secret subtree, with no general
infrastructure or cloud permissions. Bootstrap SHALL publish project credentials
using this role's existing namespace policy, with no additional cross-role grant.

#### Scenario: Project token delivery

- **WHEN** the bootstrap IaC stores the generated token
- **THEN** the project can read it using its own role, without reading the
  personal bootstrap credential or unrelated namespaces.
