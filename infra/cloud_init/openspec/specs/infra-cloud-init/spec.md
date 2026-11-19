# Infrastructure cloud-init Specification

## Purpose

Specify the shared virtual-machine bootstrap inputs owned by `infra/cloud_init`
and their PVE, Yandex Cloud, and Xen consumers. The baseline is checked-in
source at revision `550d7e79b1f5fdbc2b6017b75178471d6914082f`, observed on
2026-09-08. It describes source and generated inputs, without asserting the
configuration of any deployed VM.

## Requirements

### Requirement: Shared bootstrap foundation

The project SHALL expose `assets/cloud_init.yaml` as the canonical bootstrap
foundation consumed directly by PVE and Yandex Cloud. This foundation SHALL
own the Ansible sudo user, authorized SSH keys, local TLS root CA, and common
bootstrap policy. The project SHALL also expose the minimal PVE configuration
with QEMU guest tools as `//infra/cloud_init:cloud_init_min`.

Sources: [ownership and consumers](../../../README.md),
[bootstrap foundation](../../../assets/cloud_init.yaml),
[minimal configuration](../../../assets/cloud_init_min.yaml),
and [targets](../../../BUILD.bazel).

#### Scenario: Package a bootstrap input

- **WHEN** an infrastructure consumer requests `//infra/cloud_init` or
  `//infra/cloud_init:cloud_init_min`
- **THEN** Bazel supplies the corresponding canonical YAML file.

### Requirement: Xen bootstrap derives identity semantically

The `xen_linux` target SHALL derive `xen_linux.json` from the shared foundation.
It SHALL select exactly one user named `ansible` and exactly one SSH CA key
whose options identify the `ansible` certificate-authority principal. Missing
or duplicate matching users or keys SHALL fail rendering. The derivation SHALL
use Fedora's `users,wheel` groups and Xen guest tools, retain the other shared
base fields, and omit the explicit package-reboot and public metadata-key
flags documented as PVE-specific behavior.

Sources: [derivation contract](../../../README.md),
[Go template](../../../assets/xen_linux.json.gotmpl), and
[template target](../../../BUILD.bazel).

#### Scenario: Reorder bootstrap users or keys

- **WHEN** the foundation's users or authorized keys are reordered while the
  unique Ansible user and matching CA remain present
- **THEN** the generated Xen configuration selects that same user and CA.

#### Scenario: Ansible bootstrap identity is ambiguous

- **WHEN** the foundation contains duplicate Ansible users or duplicate matching CA keys
- **THEN** the template fails instead of choosing an arbitrary identity.

### Requirement: Preserve consumer packaging and configuration boundaries

PVE Ansible packaging SHALL retain the `files/cloud_init.yaml` and
`files/cloud_init_min.yaml` aliases for the shared inputs, preserving existing
snippet references. Generated Xen bootstrap SHALL leave consumer-specific
hostname and static network composition to its Terraform consumer. Subsequent
persistent VM configuration SHALL remain owned by Ansible. These inputs SHALL
remain repository-internal runtime inputs.

Sources: [consumer boundaries](../../../README.md) and
[PVE packaging](../../../../pve/ansible/BUILD.bazel).

#### Scenario: Package shared snippets for PVE

- **WHEN** the PVE Ansible source package is assembled
- **THEN** it includes the shared YAML files under the established `files/`
  names used by the PVE snippet deployment.
