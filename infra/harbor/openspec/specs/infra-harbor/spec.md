# Harbor infrastructure specification

## Purpose

Describe the Harbor registry infrastructure owned by `infra/harbor`, spanning
Proxmox setup, K3s host configuration, Flux-managed chart declarations, and
Harbor service Terraform. The `tf_setup`, `ansible`, `cl`, and `tf` packages
are implementation stages of this owner. This is a checked-in source
baseline, without an observation of live registry availability or deployment
success.

Baseline revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`.
Observed: 2026-09-08. Sources: [owner README](../../../README.md),
[BUILD](../../../BUILD.bazel), and the implementation links below.

## Requirements

### Requirement: Dedicated K3s registry host

Setup SHALL declare a Proxmox VM in the `src_infra_harbor` pool with two
cores, 4096 MiB memory, and 20 GiB boot, 20 GiB K3s data, and 40 GiB K3s
storage disks. Ansible SHALL apply the shared host and K3s roles with K3s
secrets encryption enabled. Its firewall variables SHALL scope API port
6443 access to the declared Flux host address.

Sources: [VM definition](../../../tf_setup/vms.tf),
[deployment playbook](../../../ansible/playbook_deploy.yaml),
[K3s configuration](../../../ansible/files/k3s.yaml), and
[host variables](../../../ansible/group_vars/all.yaml).

#### Scenario: Inspect management API access

- **WHEN** the Harbor host configuration is evaluated
- **THEN** the API-port firewall entry identifies the Flux source address
- **AND** the K3s configuration enables secrets encryption and disables
  bundled Traefik

### Requirement: Remote Flux chart reconciliation

The Harbor Flux manifests SHALL declare a remote-cluster Kustomization and
HelmRelease using `src-infra-harbor-cl-kubeconfig`. The HelmRelease SHALL
pin its chart version, enable Helm tests and drift detection, set
`https://harbor.alwaldend.com` as the external URL, and use the `local-path`
storage class for its declared persistent components.

Source: [Harbor Flux manifests](../../../cl/flux/harbor.yaml).

#### Scenario: Inspect the Harbor release declaration

- **WHEN** a maintainer evaluates the Harbor HelmRelease
- **THEN** the source selects a specific chart version and the remote Harbor
  cluster configuration
- **AND** tests, drift detection, external URL, and persistent storage choices
  are explicit

### Requirement: Vault-backed registry identity configuration

Harbor service Terraform SHALL configure Vault OIDC authentication using
client credentials from the Vault provider, the `groups` claim, and the
`src_infra_harbor_admins` administrator group. It SHALL enable automatic
onboarding and use the `username` claim for users.

Source: [authentication configuration](../../../tf/config.tf).
The baseline sets `oidc_verify_cert = false` with a CA-issues comment; this
specification does not claim that OIDC server certificate verification is
enabled.

#### Scenario: Resolve OIDC identity settings

- **WHEN** the service Terraform authentication resource is inspected
- **THEN** it references Vault-provided client credentials and the configured
  user and group claims
- **AND** its certificate-verification setting remains visible as a baseline
  limitation

### Requirement: Registry projects and guest access

Service Terraform SHALL declare an `alwaldend` project with vulnerability
scanning and automatic SBOM generation, plus a `dockerhub` proxy-cache
project linked to the Docker Hub registry. Both projects SHALL grant the
configured `src_infra_harbor_users` OIDC group the guest role.

Sources: [first-party project](../../../tf/alwaldend.tf) and
[Docker Hub proxy](../../../tf/dockerhub.tf).

#### Scenario: Inspect registry consumer permissions

- **WHEN** the managed registry projects are evaluated
- **THEN** `alwaldend` enables scanning and SBOM generation, while `dockerhub`
  references the Docker Hub proxy registry
- **AND** both group membership resources specify guest access
