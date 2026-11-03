## Why

VM provisioning needs independent review and delivery from the Ansible host
configuration in PR #233. DNS and Vault have already been delivered separately.

## What Changes

- Extract the existing tf_setup stage into a master-based Terraform-only PR.
- Give that stage its own AL configuration for Yandex credentials and backend.
- Preserve resource addresses, inputs, provider pin and existing state path.
- Remove provisioning source and credential calls from the host PR.

## Capabilities

### New Capabilities

- `local-inference-provisioning`: Independent declarative GPU provisioning.

## Impact

The initial phase extracts and publishes the source without live deployment.
The user subsequently authorized applying this root, and explicitly allowed
one initial VM boot followed immediately by a stop because the pinned provider
cannot create a stopped VM. Host deployment and changes to public DNS remain
outside this operation. The separate folder/account prerequisite requires its
own exact scope approval; image and administrator-network inputs are required.
The operator approved the nine-resource folder/account prerequisite, selected
the existing Fedora image, and requested public SSH. Provisioning therefore
uses shared Fedora cloud-init and accepts an explicitly supplied IPv4 /0.
No new external dependency or compatibility command is introduced.
