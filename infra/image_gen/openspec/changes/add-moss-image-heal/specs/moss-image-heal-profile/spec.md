## Purpose

Install the MOSS-Image-heal model and its workflows through the image-generation profile registry.

## ADDED Requirements

### Requirement: Selectable MOSS-Image-heal installation

The image-generation playbook SHALL accept `moss-image-heal` as a model profile.

#### Scenario: Install the Heal profile

- **WHEN** an operator selects `moss-image-heal`
- **THEN** Ansible installs the pinned Heal checkpoint, MOSS text encoder, small FLUX.2 VAE, Prompt Sequence node and all five Heal workflow and App files.
