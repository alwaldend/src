## Purpose

Define the provisioning, authentication, and lifecycle guarantees for a
single XCP-ng VM running native Nexus Repository 3 and PostgreSQL.

## ADDED Requirements

### Requirement: Separate infrastructure and application configuration

The deployment SHALL provide `infra/nexus/tf_setup` for VM provisioning and
DNS, `infra/nexus/ansible` for host services, and `infra/nexus/tf` for Nexus API
configuration. The two Terraform roots SHALL use distinct state backends.

#### Scenario: Application configuration changes

- **WHEN** only a Nexus repository or access setting changes
- **THEN** the service Terraform root manages that change without provisioning
  a VM, changing DNS, or installing host software

### Requirement: Provision one identifiable VM and owned DNS records

Setup Terraform SHALL declare exactly one Nexus VM with explicit template,
network, storage, sizing, and address inputs. DNS declarations SHALL have one
owner in `infra/nexus/dnsconfig.json` and feed the shared DNS module. Guest
hostname and address configuration SHALL agree with those declarations.

#### Scenario: Provisioning inputs are complete

- **WHEN** valid deployment inputs and the Nexus XO resource set are available
- **THEN** setup configuration describes one VM and its owned DNS records
  without copying address values into independently maintained inventories

#### Scenario: Required inventory is absent or ambiguous

- **WHEN** the selected template, network, storage, or identity cannot be
  resolved unambiguously
- **THEN** provisioning fails without selecting an arbitrary resource or
  falling back to shared administrator credentials

### Requirement: Run Nexus and PostgreSQL as native services

The Ansible deployment SHALL install Nexus 3 and PostgreSQL natively on the
same VM. Nexus SHALL run as a dedicated unprivileged service account under
systemd with a verified, explicitly pinned distribution and compatible Java
runtime. Installation SHALL require no container engine. PostgreSQL SHALL
be ready with its Nexus-owned database and required extensions before Nexus
starts against it.

#### Scenario: Fresh installation and reboot

- **WHEN** a fresh supported guest is configured and subsequently rebooted
- **THEN** PostgreSQL and Nexus start as native services and Nexus passes an
  application readiness check against its configured database

#### Scenario: Artifact verification fails

- **WHEN** the Nexus distribution does not match its expected checksum
- **THEN** deployment fails before installing or activating that distribution

### Requirement: Preserve persistent data during routine convergence

Application binaries SHALL be separate from persistent Nexus and PostgreSQL
data. Routine reruns SHALL preserve database contents, cached artifacts, and
completed initialization. Missing or incorrectly attached required storage
SHALL stop deployment before service startup or destructive formatting.
VM replacement and upgrades SHALL have documented data-loss and recovery
boundaries; a binary downgrade SHALL NOT be claimed as a database rollback.

#### Scenario: Unchanged deployment is rerun

- **WHEN** the role runs again with unchanged inputs against initialized data
- **THEN** it preserves the data and credentials and avoids unnecessary
  service restarts and database initialization

#### Scenario: Required data storage is missing

- **WHEN** a configured data disk or mount is missing or fails validation
- **THEN** deployment stops before Nexus or PostgreSQL writes to a fallback
  directory on the operating-system disk

### Requirement: Keep ingress outside the Nexus role

The Nexus role SHALL NOT install or configure a load balancer. The component
SHALL expose its backend address and port for the separate Traefik role.
For same-VM Traefik, Nexus and PostgreSQL SHALL listen on loopback and client
access SHALL pass through the declared HTTPS ingress.

#### Scenario: The Nexus role is applied independently

- **WHEN** only the Nexus role is selected
- **THEN** it configures the Nexus backend without changing Traefik routers,
  certificates, or another load balancer's configuration

### Requirement: Use a dedicated deployment identity and scoped secrets

The component SHALL authenticate through its own Vault AppRole and obtain XO,
SSH, DNS, database, and Nexus credentials through the repository's AL flow
only where each stage needs them. AppRole and resource-set grants SHALL remain
owned by the existing Vault and XCP-ng projects. Tracked configuration and
ordinary logs SHALL contain no secret values.

#### Scenario: A stage is invoked

- **WHEN** an operator invokes a packaged setup, Ansible, or service target
- **THEN** its plugin labels and packaged dependencies select the required
  identity, state backend, and credentials for that stage

### Requirement: Bootstrap API credentials without repeated resets

The deployment SHALL complete Nexus's initial administrator password change
using a Vault-supplied credential before service Terraform needs API access.
Normal reruns SHALL verify completed bootstrap without resetting credentials.
Unexpected authentication failures SHALL fail with a redacted diagnostic
instead of bypassing authentication or editing database password records.

#### Scenario: Nexus is starting for the first time

- **WHEN** Nexus is ready and its initial password file is available
- **THEN** bootstrap changes the administrator password through the supported
  API and verifies the intended credential without logging either password

#### Scenario: Bootstrap already completed

- **WHEN** the initial password file is absent and the intended credential works
- **THEN** bootstrap succeeds without changing the password

#### Scenario: Neither supported credential works

- **WHEN** the intended credential fails and no valid initial credential exists
- **THEN** deployment reports an authentication failure and performs no reset

### Requirement: Distinguish offline validation from deployment evidence

The component SHALL provide offline packaging, formatting, configuration, and
fixture checks plus documented client acceptance and recovery procedures.
Offline checks SHALL NOT invoke a live apply, inventory playbook, or Vault
write. Any deployment or live acceptance run SHALL require separate authority
for its named environment.

#### Scenario: Source validation runs without infrastructure access

- **WHEN** repository checks validate the deployment definitions
- **THEN** they use declared artifacts and non-secret fixture inputs, report
  unsupported checks honestly, and do not claim that Nexus is deployed
