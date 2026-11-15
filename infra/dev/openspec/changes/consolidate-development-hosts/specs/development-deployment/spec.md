## Purpose

Define unified development host ownership and deployment while preserving state and machine-specific configuration.

## ADDED Requirements

### Requirement: Unified development deployment

The repository SHALL own both hosts under infra/dev and expose one packaged Ansible inventory and entry point with host-specific groups and variables. Existing runner activation, physical storage, account and package-backend behavior SHALL remain explicit and unchanged.

#### Scenario: Select either host

- **WHEN** the packaged deployment inventory is evaluated
- **THEN** secure1.dev.alwaldend.com and dev1.dev.alwaldend.com are present
- **AND** limiting to either host selects its appropriate configuration without applying the other host's storage or service tasks

#### Scenario: Deploy both host groups

- **WHEN** the shared Bazel deployment target runs
- **THEN** one standard Ansible invocation runs both host-specific plays under src_infra_dev
- **AND** that identity has no legacy service-secret or Terraform-state access

### Requirement: Consistent host and service identity

Secure1 SHALL resolve to 192.168.10.102 and dev1 SHALL resolve to 192.168.1.210. Dev1 application routes, DNS, ingress targets, and certificate scopes SHALL consistently use the new dev1.dev.alwaldend.com suffix while retaining ingress mTLS and loopback application listeners.

#### Scenario: Trace an application request

- **WHEN** a maintainer evaluates T3 Code or ComfyUI routing
- **THEN** the public name routes through ingress to its new dc1 target and matching backend TLS hostname

### Requirement: Unified DNS ownership and preserved work state

The consolidation SHALL expose one dns root consuming the entire canonical DNS document through src_infra_dev and its owner-named backend. It SHALL preserve existing OpenSpec history and active work. Legacy state and access SHALL remain available until a separately authorized handoff reconciles old resource identities and adopts existing records without destructive recreation. Source checks SHALL be offline and SHALL NOT establish live deployment success.

#### Scenario: Inspect migration boundaries

- **WHEN** a maintainer reviews the consolidated roots
- **THEN** both hosts and dev1 services use one DNS root, identity and backend
- **AND** the old backends are identified by their reconciled state ownership
- **AND** no live state-changing operation is required to validate source

### Requirement: Current deployment names and inputs

Active inventory groups and playbooks SHALL identify secure1 and dev1. The host-only secure1 play SHALL NOT inject or require an unused Actions registration token. Current operating-system account names, reusable role APIs and persistent application profiles SHALL remain unchanged.

#### Scenario: Package the host-only runner baseline

- **WHEN** the shared deployment configuration is built
- **THEN** it includes no OpenHands service secret reference
- **AND** secure1 retains its shared host role without an Actions worker or registration-token lookup

### Requirement: No legacy Proxmox runner declaration

The owner SHALL NOT declare the old runner1 VM, its DNS record or the host1.forgejo-runner inventory entry. Secure1 and dev1 SHALL remain configured. This cleanup SHALL NOT mutate live infrastructure or Terraform state.

#### Scenario: Inspect runner desired state

- **WHEN** the current owner configuration is inspected
- **THEN** no Proxmox VM or Proxmox authentication/provider is declared
- **AND** only secure1 and dev1 remain in the deployment inventory

### Requirement: Verified SSH identity through host renames

The rollout SHALL retain strict SSH host verification. Dev1's managed known_hosts SHALL trust secure1 through the canonical server CA and declare dev1's public host key, verified from the local host filesystem and SSH handshake. Temporary old-name connection aliases SHALL NOT become persistent deployment defaults.

#### Scenario: Repeat deployment after cutover

- **WHEN** the normal shared deployment connects through the new DNS names
- **THEN** both host identities pass verification using the declared controller trust

### Requirement: OpenHands removal

Dev1 SHALL NOT install or start the OpenHands server or read its secret. The
completed one-time service, executable and dependency removal SHALL NOT remain
in the normal deployment playbook.

#### Scenario: Deploy after completed removal

- **WHEN** dev1 is deployed after OpenHands retirement
- **THEN** no OpenHands installation or removal task runs
- **AND** the retired service stays disabled and its executables remain absent

### Requirement: Recreated secure1 VM

The owner SHALL provision secure1 using its own XO identity and backend with
8 CPUs, 16 GiB RAM and a 500 GiB boot disk. It SHALL use the canonical DNS
address and hostname and preserve the host-only Ansible role. The authorized
replacement SHALL discard the old VM without migrating its VM state.

#### Scenario: Replace the legacy allocation

- **WHEN** the authorized recreation completes
- **THEN** secure1 runs as a new VM owned by infra/dev and the old VM and deployment identity are absent
- **AND** dev1 has not been rebooted and its application ingress remains available
