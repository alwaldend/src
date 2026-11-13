## Purpose

Provide reproducible public MOSS-Image downloads and matching examples for the
operator's ComfyUI rental without changing its access or spending controls.

## ADDED Requirements

### Requirement: MOSS bundle and examples

The configuration SHALL provide a selected `moss-image` profile with immutable
publisher URLs for both SFT and Base models, a compatible local text encoder,
VAE, edit adapter and color/depth adapter, and six matching example workflows.

#### Scenario: Provision the requested profile

- **WHEN** the operator deploys the selected profile
- **THEN** the native provisioner downloads its models to their matching directories
  and stages all six examples under the ComfyUI user's workflows

### Requirement: Sequence example catalog

The MOSS profile SHALL also list six sequence examples for authenticated API
staging, preserving the original six first-boot workflows so adding examples
does not require replacing an existing rental.

#### Scenario: Extend the running instance's examples

- **WHEN** the operator stages the sequence catalog on the existing instance
- **THEN** ComfyUI lists all twelve examples and the original instance is retained

### Requirement: Retain 48 GB rental controls

The selected hardware SHALL retain the existing `48gb` requirements, authenticated
HTTPS, Vault-managed credentials and USD 1/hour ceiling.

#### Scenario: Deploy the MOSS rental

- **WHEN** the reviewed Terraform apply creates the rental
- **THEN** it requests the existing named 48 GB GPU profile and preserves its
  storage, price and security constraints
