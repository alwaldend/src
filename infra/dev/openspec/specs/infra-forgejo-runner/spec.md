# Secure runner infrastructure specification

## Purpose

Describe the host-only secure1 baseline under infra/dev. The checked-in source
does not install or register an Actions worker and does not establish live state.

## Requirements

### Requirement: No legacy Proxmox allocation

The owner SHALL NOT declare runner1/VM 1200, its old DNS entry or its inventory
entry. Secure1 SHALL retain its host-only configuration. Its DNS SHALL belong
to the common dns root without a Proxmox provider or login.

#### Scenario: Inspect runner allocation

- **WHEN** the current source is inspected
- **THEN** no legacy Proxmox allocation or separate tf_setup root remains
- **AND** secure1 DNS and host configuration remain available

### Requirement: Host-only deployment baseline

The [secure1 play](../../../ansible/playbook_secure1.yaml) SHALL apply only the
shared host role to the secure1 inventory group. It SHALL NOT require runner
registration credentials, an Actions instance URL or a runner disk setting.

#### Scenario: Determine the playbook effect

- **WHEN** the current playbook is inspected
- **THEN** only the shared host role is enabled
- **AND** runner installation and registration remain outside this baseline
