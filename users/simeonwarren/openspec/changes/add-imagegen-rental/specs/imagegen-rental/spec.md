## ADDED Requirements

### Requirement: Imagegen rental

The implementation SHALL provide the following behavior: Add imagegen/tf_setup with random_password, Vault KV storage, and the reusable exposed-mode module. Reuse the existing user_simeonwarren AppRole, regular user secret namespace, and independent Vault-backed state. Use the newest published Vast Forge Neo/CUDA 12.9 image tag without a template by default; provide optional declarative public model downloads; allow an optional reviewed environment-only template; retain the USD 1/hour cap. Do not execute a live rental as implementation validation.

#### Scenario: Authorized configuration

- **WHEN** an operator supplies valid reviewed configuration
- **THEN** the declared behavior SHALL be available through repository-owned targets

#### Scenario: Offline validation

- **WHEN** validation runs without live credentials
- **THEN** it SHALL NOT rent capacity or write live secrets

### Requirement: Rental-owned default checkpoint

The rental root SHALL default to the approved Juggernaut XL v9 standalone safetensors checkpoint from a fixed publisher revision. The shared template SHALL remain model-configurable.

#### Scenario: Default rental configuration

- **WHEN** model_downloads is omitted
- **THEN** the native provisioner SHALL receive the public Juggernaut XL v9 checkpoint download before starting Forge

#### Scenario: Explicit model selection

- **WHEN** the caller supplies model_downloads
- **THEN** that list SHALL replace the rental default, including an empty list to disable configured downloads
