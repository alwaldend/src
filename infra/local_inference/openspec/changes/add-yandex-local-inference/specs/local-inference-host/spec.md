## Purpose

Provide a declaratively managed GPU host for Codex CLI with local inference
and tools, encrypted persistent storage, and explicit deployment boundaries.

## ADDED Requirements

### Requirement: Requested GPU shape and encrypted storage

The provisioning definition SHALL select Intel Broadwell V100 with one GPU,
eight vCPUs, and 96 GiB RAM. All declared persistent disks MUST use KMS
encryption. SSH ingress MUST require explicitly supplied administrator CIDRs.

#### Scenario: Inspect the provisioning configuration

- **WHEN** an operator reviews the source before provisioning
- **THEN** the requested shape, KMS key, encrypted disk, and restricted SSH
  rules are declared without needing a live apply

### Requirement: Local inference and local tools

Codex SHALL use only the host's loopback inference endpoint. Hosted search,
apps, plugins, remote MCP tools, and image generation MUST be disabled.
Local command execution, editing, and the approved loopback T3 MCP MUST
remain available. Auto review
SHALL use the selected local provider. Full access SHALL remain selectable.
Runtime accounts SHALL retain normal networking; Codex permissions govern
commands rather than an additional account-level egress firewall.

#### Scenario: Select approval behavior

- **WHEN** an operator starts a T3 session
- **THEN** Auto uses local inference by default and Full access is available

#### Scenario: Override project configuration

- **WHEN** a workspace attempts to enable a remote integration or provider
- **THEN** managed requirements keep the provider local and remote tools off

### Requirement: Reproducible runtime installation

Deployment SHALL verify checksums of runtime archives and model files before
using them. The inference service MUST require a working V100 CUDA backend
and MUST NOT silently substitute CPU-only inference.

#### Scenario: An artifact is missing or invalid

- **WHEN** a required model, checksum, or compatible driver is unavailable
- **THEN** configuration fails before starting the inference service

### Requirement: Preparation does not deploy

Source validation MUST NOT provision cloud resources, write Vault data, or
connect to inventory hosts. Documentation SHALL distinguish offline checks
from future live acceptance and deployment prerequisites.

#### Scenario: Validate this change

- **WHEN** an agent validates the prepared source
- **THEN** only offline checks run and no deployment is reported
