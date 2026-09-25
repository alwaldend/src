## Why

The download component needs its own infrastructure identity for state,
provisioning, DNS, SSH, and certificate automation. Existing components'
credentials should not become the new service's deployment identity.

## What Changes

- Declare `src_infra_download` in a standalone Terraform root with its own
  backend, using the existing AppRole composition and name-based data lookups.
- Grant its own state access, required provider-secret references, and the
  SSH and certificate permissions selected for the download component.
- Add the existing group membership and outputs needed by Yandex folder
  provisioning and XCP-ng identity/resource-set assignment.
- Keep credential values in Vault and preserve unrelated identities.

## Capabilities

### New Capabilities

- `download-identity`: Component-scoped authentication and infrastructure
  permissions for the download deployment.

### Modified Capabilities

None. Existing DNS-owner isolation requirements continue to apply.

## Impact

The declaration and policies belong to
`infra/vault/approles/src_infra_download`, consuming reusable modules in
`projects/tf_modules`. The core `infra/vault/tf` stage looks up the created
identity/group by name and retains ownership of shared memberships. The [hosting plan](../../../../../infra/download/openspec/changes/add-static-hosting/proposal.md)
owns AL injection and the consuming infrastructure. Yandex Cloud and XCP-ng
retain ownership of their provider-side assignments.

Source checks must establish policy composition and isolation. Authenticating,
creating roles, issuing credentials, or applying assignments is a later live
operation and is not authorized by this planning delivery.
