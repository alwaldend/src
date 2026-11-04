## Purpose

Define the Vast.ai infrastructure project's authentication boundary and
component-local Terraform state configuration in checked-in source.

## ADDED Requirements

### Requirement: Dedicated component authentication

The project SHALL expose repository-internal AL and Vault targets using
the `src_infra_vastai` AppRole and inherited infrastructure configuration.

#### Scenario: Select a component Vault command

- **WHEN** an operator selects `//infra/vastai:vault`
- **THEN** its AL configuration selects the `src_infra_vastai` identity
- **AND** source validation requires no live Vault authentication

### Requirement: Component-local state configuration

The AL configuration SHALL register the Vault HTTP backend for `tf=main`
under `secrets/alwaldend.com/vault1/approles/src_infra_vastai/tf_backend/tf`.

#### Scenario: Inspect future service configuration wiring

- **WHEN** the component AL configuration is loaded
- **THEN** the backend registration is scoped to `tf=main` and its own namespace
- **AND** the project does not provision Vast.ai resources merely by building it

### Requirement: API-key injection for the component CLI

The component CLI SHALL run the pinned upstream CLI with `VAST_API_KEY`
injected from field `api_key` at
`secrets/alwaldend.com/vault1/approles/src_infra_vastai/vastai`.
The key SHALL be passed through the command environment rather than arguments
or a persistent key file. Standard component Vault commands SHALL NOT select
this API-key injection.

#### Scenario: Select the authenticated CLI

- **WHEN** an operator runs `//infra/vastai:vastai`
- **THEN** its wrapper selects `vastai=1` and packages the CLI and injector
- **AND** the injector reads the named field using the component AppRole

#### Scenario: Set up or inspect the component secret

- **WHEN** an operator selects a component Vault command
- **THEN** it does not fetch or require the Vast.ai API key

#### Scenario: Build without credentials

- **WHEN** the component CLI launcher is built
- **THEN** no live Vault or Vast.ai request is made
