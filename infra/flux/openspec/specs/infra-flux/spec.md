# Flux infrastructure specification

## Purpose

Describe the Flux management cluster's VM definition, K3s configuration,
and Git reconciliation entry points owned by `infra/flux`. The `tf_setup`,
`ansible`, and `cl` packages implement stages of this owner. This is a
checked-in source baseline; no cluster health or successful deployment was
observed for this specification.

Baseline revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`.
Observed: 2026-09-08. Sources: [owner README](../../../README.md),
[BUILD](../../../BUILD.bazel), and the implementation links below.

## Requirements

### Requirement: Dedicated K3s management host

Terraform setup SHALL declare a Proxmox VM in the `src_infra_flux` pool with
two cores, 4096 MiB memory, and separate 20 GiB boot, K3s data, and K3s
storage disks. The deployment playbook SHALL apply the shared host and K3s
roles, and the K3s configuration SHALL enable secrets encryption and disable
the bundled Traefik component.

Sources: [VM definition](../../../tf_setup/vms.tf),
[deployment playbook](../../../ansible/playbook_deploy.yaml), and
[K3s configuration](../../../ansible/files/k3s.yaml).

#### Scenario: Inspect the management cluster deployment definition

- **WHEN** a maintainer evaluates the setup and deployment source
- **THEN** VM provisioning, host configuration, and K3s configuration are
  separate stages with the recorded compute and disk allocation
- **AND** the K3s configuration enables secrets encryption at rest

### Requirement: Git reconciliation for Flux and Harbor

The cluster's root Kustomization SHALL include the Flux and Harbor
reconciliation trees. Flux's reconciliation resources SHALL read the
`flux-git-src` GitRepository and use SOPS decryption through
`flux-sops-secret`, with pruning enabled and forced replacement disabled.
The Harbor tree SHALL include its remote-cluster configuration and Harbor's
Flux manifests.

Sources: [root Kustomization](../../../cl/kustomization.yaml),
[Flux reconciliation](../../../cl/src-infra-flux-cl/manifests.yaml),
and [Harbor reconciliation tree](../../../cl/src-infra-harbor-cl/kustomization.yaml).

#### Scenario: Resolve the reconciliation source tree

- **WHEN** the checked-in root Kustomization is traversed
- **THEN** both `src-infra-flux-cl` and `src-infra-harbor-cl` are included
- **AND** the Flux reconciliation resources retain their Git source,
  decryption reference, and explicit prune and force settings

### Requirement: Packaged operator and controller configuration

The cluster package SHALL provide repository-configured Flux, Flux
operator, kubectl, Helm, certificate-manager, and SOPS command wrappers.
Its Flux instance SHALL declare source, Kustomize, Helm, notification,
image-reflector, image-automation, and source-watcher controllers, with
network policy enabled.

Sources: [cluster BUILD](../../../cl/BUILD.bazel) and
[Flux instance](../../../cl/flux-system/flux-instance.yaml).

#### Scenario: Select a repository operator entry point

- **WHEN** an operator selects a declared cluster command target
- **THEN** the target supplies the packaged tool and the Flux owner's `al`
  configuration
- **AND** the controller source records the enabled components without
  asserting that they are currently available in a live cluster
