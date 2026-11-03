## ADDED Requirements

### Requirement: Independent VM provisioning

The tf_setup stage SHALL own the GPU VM, network, SSH ingress, encrypted
persistent disk, rotating KMS key and Yandex child DNS zone independently of
host configuration. It SHALL use the existing pinned Yandex provider and
component AppRole/backend. It MUST NOT require Ansible or Cloudflare credentials.

#### Scenario: Review the provisioning source

- **WHEN** an operator inspects the independent Terraform PR
- **THEN** it declares the selected GPU, vCPUs and host RAM, encrypted
  storage and explicit administrator CIDRs without deploying any resources

### Requirement: Preserve provisioning identity and inputs

Extraction MUST preserve resource addresses, the tf_setup state backend, and
required immutable image ID and administrator network inputs. It MUST retain
storage destruction guards outside an explicitly authorized replacement
attempt, restore them after that attempt, and MUST NOT add migration or
compatibility roots.

#### Scenario: Validate the extracted root

- **WHEN** the extracted source is validated with its backend disabled
- **THEN** its configuration is valid without contacting live infrastructure
  or changing existing resource ownership

### Requirement: Honor explicit deployment inputs and power state

The authorized deployment SHALL use the existing immutable Fedora 43 image
and explicitly permit public SSH when supplied as 0.0.0.0/0. Provisioning
SHALL NOT imply that the NVIDIA driver or host applications are configured.
After the authorized initial boot, the instance MUST be stopped immediately
and its stopped status verified through the control plane.

#### Scenario: Provision the requested stopped Fedora VM

- **WHEN** the approved Terraform plan is applied with the operator's Fedora
  image and public SSH input
- **THEN** the persistent resources remain Terraform-owned, the new instance
  is stopped after creation, and a fresh plan reconciles the observed state

### Requirement: Use the selected A100 configuration in zone a

The current deployment SHALL use ru-central1-a and one 80 GB A100 on
platform gpu-standard-v3 with 28 vCPUs and 119 GiB RAM. It SHALL preserve
the encrypted 200 GiB disk, existing authentication and explicit SSH inputs.
The boot disk SHALL use legacy PCI topology V1, as tested in the final
compatibility attempt. Supported zone inputs SHALL remain restricted to
the selected platform's zones a/b.

#### Scenario: Retain the last attempted larger GPU configuration

- **WHEN** the operator requests larger VRAM alternatives and increases
  zone-a CPU quota for the A100 trial
- **THEN** the source declares the A100 shape without claiming that quota
  headroom or a successful plan establishes allocation capacity

### Requirement: Operator-selected preemptible capacity

The current deployment SHALL set preemptible to true for the selected
A100 VM. The persistent encrypted disk SHALL remain independent
of VM interruption. The initial-boot-and-immediate-stop requirement SHALL
continue to apply.

#### Scenario: Try the approved preemptible alternative

- **WHEN** the operator selects preemptible capacity after regular allocation
  failures and receiving the interruption and 24-hour lifetime explanation
- **THEN** the plan creates only the preemptible VM, preserving its existing
  disk and network, and any created VM is immediately stopped and verified
