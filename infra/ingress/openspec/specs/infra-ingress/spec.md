# Infrastructure ingress

## Purpose

Describe Yandex Cloud ingress provisioning, the RouterOS WireGuard connection,
and the Ansible-managed Traefik configuration. The baseline covers checked-in
desired state, not a live routing or availability observation.

Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`.
Observation date: 2026-09-08. Sources are linked in full; no excerpts are used.

Sources: [component contract](../../../README.md),
[VM resources](../../../tf/vms.tf),
[WireGuard resources](../../../tf/wireguard.tf),
[deployment playbook](../../../ansible/playbook_deploy.yaml),
[Traefik entry points](../../../ansible/files/traefik.toml), and
[Traefik routes](../../../ansible/files/traefik_dynamic.toml).

## Requirements

### Requirement: Provision ingress hosts from the declared host map

Terraform SHALL create addresses, encrypted boot disks, and instances for the
active `local.vpc` entries, using the shared cloud-init source. The baseline
active map SHALL contain `host1` in `ru-central1-d`; the commented `host2` entry
SHALL NOT be treated as an enabled resource declaration.

#### Scenario: Evaluate the baseline host configuration

- **WHEN** the checked-in ingress Terraform configuration is evaluated
- **THEN** instance resources are declared for `host1`
- **AND** its boot disk references the declared KMS symmetric key
- **AND** instance metadata includes the shared cloud-init content

### Requirement: Connect declared hosts through WireGuard

The RouterOS configuration SHALL declare the `ingress-vpc` interface and a peer
for each active ingress host. Peer endpoints SHALL use the corresponding cloud
address, and allowed addresses SHALL come from the Ansible inventory's `wg_ip`.

#### Scenario: Configure an ingress peer

- **WHEN** an ingress host is present in the active Terraform host map
- **THEN** its WireGuard peer uses the host's declared public and preshared key
  variables and inventory address
- **AND** the interface belongs to the configured LAN, forwarding, and ICMP lists

### Requirement: Route to backend addresses that do not return to ingress

Each `traefik_services` backend target MUST resolve from the ingress hosts
without resolving back to ingress. The service owner SHALL declare the dedicated
site-local address. Generated HTTPS backend transports SHALL use the service
hostname as TLS server name.

#### Scenario: Add a service behind ingress

- **WHEN** a service is added to `traefik_services`
- **THEN** its host route forwards to `https://` followed by its backend target
- **AND** that target is a dedicated site-local address defined by the service
  owner
- **AND** backend TLS uses the service's configured hostname

### Requirement: Deploy Traefik with HTTPS and client-certificate authentication

The ingress playbook SHALL apply the host, WireGuard, and Traefik roles. The
Traefik templates SHALL redirect the web entry point to HTTPS and configure
default TLS client authentication as `RequireAndVerifyClientCert` using the
configured client CA file.

#### Scenario: Render the ingress TLS configuration

- **WHEN** the Ansible role renders the checked-in Traefik templates
- **THEN** port 80 redirects to the secure entry point
- **AND** the default TLS options require a certificate verified against the
  configured client CA
