## Purpose

Bootstrap a dedicated test team and Terraform-generated project credential
without manual management of team tokens or any compute rental charges.

## ADDED Requirements

### Requirement: Declarative credential bootstrap

The root SHALL create a test_ team without credit transfer, generate a
project API key using the Vast.ai provider, and store it directly in Vault
for the project's dedicated AppRole. The personal key SHALL remain in its
original context and SHALL NOT be delivered to the project.

#### Scenario: Bootstrap and consume

- **WHEN** the bootstrap root is applied with the personal token
- **THEN** the project injector can obtain its generated team-scoped token
  from its own secret namespace, with no compute resources rented.
