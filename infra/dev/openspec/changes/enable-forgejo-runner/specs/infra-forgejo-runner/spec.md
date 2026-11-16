## MODIFIED Requirements

### Requirement: No legacy Proxmox allocation

The owner SHALL NOT declare runner1/VM 1200, its old DNS entry or its inventory
entry. Secure1 SHALL retain its common DNS root and XO allocation without a
Proxmox provider or login.

#### Scenario: Inspect runner allocation

- **WHEN** the current source is inspected
- **THEN** no legacy Proxmox allocation or separate DNS root remains
- **AND** secure1 DNS and host configuration remain available

### Requirement: Host-only deployment baseline

The secure1 deployment SHALL preserve the shared host baseline and enable one
repository-scoped Forgejo host worker using the existing cached identity.
The runner account SHALL have repository build prerequisites and SHALL NOT
receive deployment credentials or sudo access.

#### Scenario: Determine the playbook effect

- **WHEN** the authorized secure1 runner deployment completes
- **THEN** the service is started and enabled with one host-executor slot
- **AND** API-issued credentials are installed privately without new registration
- **AND** Node.js, Bazel, AL and Android prerequisites are available to that account

#### Scenario: Inspect CI activation

- **WHEN** only the runner deployment is applied
- **THEN** repository workflow activation remains independently controlled
