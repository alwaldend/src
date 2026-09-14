## Why

The existing runner scaffold does not install a worker and targets Proxmox.
Deploy a usable Forgejo Actions worker on XCP-ng with Vault CI authentication.

## What Changes

- Replace the inactive Proxmox scaffold with an extensible XCP-ng runner map,
  initially containing `secure` with 8 vCPUs, 16 GiB RAM and a 500 GiB disk.
- Configure the VM through the shared host, dev_vm and Forgejo runner roles.
- Reuse the component AppRole, establish its XCP-ng access and manage DNS
  through the existing owner state.
- Configure CI authentication with Vault and verify a simple Forgejo job.
- Enforce protected-branch runner execution in IaC where Forgejo supports it;
  verify the actual enforcement boundary rather than relying on job conditions.

## Capabilities

### New Capabilities

- `secure-runner-ci`: Vault authentication and protected-branch CI execution.

### Modified Capabilities

- `infra-forgejo-runner`: Active XCP-ng provisioning and dev_vm deployment.

## Impact

Changes affect this owner, the shared runner role, Vault and XCP-ng access
definitions, and Forgejo CI configuration. The user authorized provisioning,
Vault setup, deployment through the packaged Ansible target, and a live CI
test. Unrelated infrastructure changes and destruction are outside that scope.
