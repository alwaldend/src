## ADDED Requirements

### Requirement: Imagegen access

The implementation SHALL provide the following behavior: Retain ssh_only by default; exposed publishes only 7860/tcp through a dedicated TLS and Basic-auth Caddy configuration. Forge stays on loopback. Reject unsafe template options and weak credentials; disable outbound tunnels and synchronization while retaining loopback-only portal, Jupyter, and TensorBoard services.

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

The personal consumer SHALL select the newest published Vast Forge Neo/CUDA 12.9 image tag and SHALL inherit no template by default. Optional reviewed templates SHALL
contain only environment assignments. The module SHALL preserve loopback-only
management tools and expose only the authenticated Forge HTTPS proxy when selected.

#### Scenario: Default personal rental

- **WHEN** the personal root provisions without an optional template
- **THEN** it SHALL use the newest selected image tag without stock template ports or provisioning scripts
