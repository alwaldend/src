## ADDED Requirements

### Requirement: Imagegen rental

The implementation SHALL provide the following behavior: Add imagegen/tf_setup with random_password, Vault KV storage, and the reusable exposed-mode module. Reuse the existing user_simeonwarren AppRole, regular user secret namespace, and independent Vault-backed state. Declare a direct published Vast ComfyUI/CUDA 12.9 image tag without a template by default; provide optional declarative public model downloads; allow an optional reviewed environment-only template; retain the USD 1/hour cap. Offline validation SHALL NOT rent capacity; the user separately authorized one live rental after security review.

#### Scenario: Authorized configuration

- **WHEN** an operator supplies valid reviewed configuration
- **THEN** the declared behavior SHALL be available through repository-owned targets

#### Scenario: Offline validation

- **WHEN** validation runs without live credentials
- **THEN** it SHALL NOT rent capacity or write live secrets

### Requirement: Rental-owned default checkpoint

The rental root SHALL call the ComfyUI module and default to the full Chroma1-HD-Flash safetensors checkpoint and its required T5 encoder and VAE from fixed publisher revisions. The shared modules SHALL remain model-configurable.

#### Scenario: Default rental configuration

- **WHEN** model_downloads is omitted
- **THEN** the native provisioner SHALL receive the three public model downloads in native ComfyUI directories before ComfyUI serves generation requests

#### Scenario: Explicit model selection

- **WHEN** the caller supplies model_downloads
- **THEN** that list SHALL replace the rental default, including an empty list to disable configured downloads
