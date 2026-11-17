## Why

Secure1 needs a persistent action cache outside its boot disk. The user selected a NAS-backed virtual disk through XCP-ng.

## What Changes

- Attach a dedicated cache disk while preserving the existing VM and boot disk.
- Keep namespace/cache configuration under infra/nas/approles/src_infra_dev with independent Terraform state and targets; the VM root only consumes the provisioned repository.
- Publish Terraform disk metadata to Vault and load it through the existing Ansible injector automatically; retain guest-only checks, persistent mounting and the worker mount requirement.
- Configure the action cache and bounded GC in the managed runner bazelrc.

## Capabilities

### New Capabilities

- `runner-cache-mount`: Persistent virtual-disk action cache for secure1.

### Modified Capabilities

None.

## Impact

Changes infra/dev VM provisioning and secure1 Ansible configuration. The NAS and XCP-ng own their storage declarations. No guest network-filesystem dependency is added. Deployment remains separate.
