## ADDED Requirements

### Requirement: Imagegen access

The implementation SHALL provide the following behavior: Retain ssh_only by default; exposed publishes only 7860/tcp through a dedicated TLS and Basic-auth Caddy configuration. The selected runtime stays on loopback. Reject unsafe template options and weak credentials; disable outbound tunnels and synchronization while retaining loopback-only portal, Jupyter, and TensorBoard services.

#### Scenario: Authorized configuration

- **WHEN** an operator supplies valid reviewed configuration
- **THEN** the declared behavior SHALL be available through repository-owned targets

#### Scenario: Offline validation

- **WHEN** validation runs without live credentials
- **THEN** it SHALL NOT rent capacity or write live secrets

### Requirement: GPU and hosting constraints

The module SHALL require at least 32 GiB VRAM per GPU, request verified datacenter
capacity, and reject offers without returned hosting_type=1 and verification=verified.
It SHALL preserve the configured hourly price ceiling and SHALL NOT fall back to
community hosts or smaller GPUs.

#### Scenario: Missing certified tier metadata

- **WHEN** an offer has insufficient GPU memory or lacks certified tier metadata
- **THEN** Terraform SHALL fail before creating a rental

### Requirement: Approved native runtime

Each runtime consumer SHALL declare a direct published Vast image tag and SHALL inherit no template by default. Optional reviewed templates SHALL
contain only environment assignments. The module SHALL preserve loopback-only
management tools and expose only the authenticated image-generation HTTPS proxy when selected.

#### Scenario: Default personal rental

- **WHEN** the personal root provisions without an optional template
- **THEN** it SHALL use the explicitly configured image tag without stock template ports or provisioning scripts

### Requirement: External model components

The module SHALL support public safetensors text-encoder files and an optional
selection of a declared checkpoint and encoder/VAE modules in native Forge settings.
It SHALL preserve the existing HTTPS authentication and loopback bindings.

#### Scenario: Chroma components

- **WHEN** the caller selects a checkpoint, encoder, and VAE declared in model_downloads
- **THEN** the provisioner SHALL stage them in their Forge directories and startup SHALL select them globally and for the chosen UI preset

### Requirement: Runtime module boundaries

The base module SHALL own offer selection, rental lifecycle, native downloads and
access controls without Forge-specific paths or settings. Separate Forge and
ComfyUI modules SHALL configure their native model directories and loopback
listeners while preserving base constraints and connection outputs. ComfyUI SHALL
retain its bundled API wrapper on loopback and disable unapproved custom nodes and partner APIs.

#### Scenario: ComfyUI access

- **WHEN** an exposed ComfyUI rental is configured
- **THEN** only the authenticated HTTPS proxy and account SSH SHALL be published, with ComfyUI on loopback and model files in native ComfyUI directories

### Requirement: Approved ComfyUI node provisioning

The ComfyUI wrapper SHALL pass explicitly approved custom-node repositories at
immutable revisions and version- and hash-locked Python requirements to Vast's
native provisioner. Only approved node directories SHALL be allowlisted;
public access controls SHALL remain unchanged. The wrapper SHALL compact named native workflow documents and stage them after
workspace initialization. The base SHALL reject rendered startup exceeding
16,384 bytes before creating a rental.

#### Scenario: Approved style-reference node pack

- **WHEN** a caller supplies a reviewed node repository and requirements
- **THEN** Vast's native provisioner SHALL fetch the pinned sources and install requirements with mandatory hashes and no dependency resolution
