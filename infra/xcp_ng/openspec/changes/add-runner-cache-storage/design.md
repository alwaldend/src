## Context

The installed XCP-ng 8.3 SMB driver uses SMB3.0; cifs-utils is already installed.
The pinned XO provider can attach disks but cannot create storage repositories.

## Decisions

Proceed with TrueNAS mandatory SMB encryption, a dedicated Vault-backed account,
and an XCP-ng SMB repository. This meets the user's selected virtual-disk path
without exposing plaintext AUTH_SYS NFS or adding guest mount clients.

Use creation-only `terraform_data` and the existing XO helper for the provider
gap. Identify the host unambiguously; accept an existing repository only with
the expected ownership marker, SMB type and pool. Never delete storage or rotate
its credentials automatically. Reject conflicts; retry discovery after a lost
creation response before any operator-authorized rerun. Terraform's backend lock
serializes this workflow; concurrent standalone bootstrap is unsupported.

Terraform owns disk attachment identity and publishes metadata to Vault for Ansible injection.
Remove the custom cache-device command, its tests and Ansible packaging/authentication.
The XO helper remains only for the provider gap in SR creation.

## Risks / Trade-offs

Creation-only bootstrap cannot continuously reconcile PBD configuration or remove
the SR. Retain it when Terraform forgets the bootstrap resource. Changes to its
target require explicit reconciliation; credentials are not replacement triggers.
Mandatory NAS encryption rejects clients unable to encrypt. Live mount, protocol
and cache verification remain deployment acceptance checks.

Read the hypervisor's own published NAS credential. The independent development
NAS project owns the namespace traversal and cache grants through the dataset
module. Apply shared NAS accounts/settings and then that project before XCP-ng.
