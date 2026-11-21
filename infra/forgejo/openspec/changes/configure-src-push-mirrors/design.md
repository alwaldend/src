## Context

Forgejo owns ongoing native push synchronization for `alwaldend/src`.
HTTPS transfers Git LFS objects. The pinned Forgejo Terraform provider
supports pull mirroring, and Forgejo's API supports push-mirror discovery
and creation.

## Design

Use the existing Terraform bootstrap pattern with a packaged Go command.
The catalog supplies the source repository and destinations. The command
uses the existing Forgejo controller credential and reads destination
credentials from Vault only when it needs to create a mirror. Terraform
state carries no secret inputs.

Set the native branch filter from the catalog's default branch and the
`releases/*` namespace. Forgejo adds its tag refspec to the filtered mirror.
An exact match is adopted without an API write. Fail closed on conflicting
mirrors or incomplete discovery.

Proceed with native filter reconciliation. Forgejo 17.0.1 exposes only
create and delete operations for push mirrors through its API. Replacing the
entry keeps native scheduling and LFS transfer. A separate synchronization
job would also need to own those functions. An unfiltered mirror continues
publishing development refs.

Before creating or replacing a mirror, require the source trunk and verify
that each destination branch and tag has the same object ID in Forgejo.
Forgejo's filtered mirror push prunes excluded destination branches, so
this check requires their tips to remain recoverable from their source refs.
Read the destination credential before removing the old mirror. A retry can
create a missing mirror or adopt a completed replacement. The Terraform
backend lock serializes reconciliation.

The Terraform provisioner runs on creation. The desired filter triggers
replacement of the Terraform resource to reconcile previously applied mirrors.
Refresh does not detect live drift. Removing the Terraform resource alone
leaves the live mirror in place.
