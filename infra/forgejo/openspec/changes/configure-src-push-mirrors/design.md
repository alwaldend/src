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

Adopt the existing matching HTTPS entry without an API write. Fail closed
on conflicting mirrors or incomplete discovery. Before activating a new
native mirror, require each destination branch and tag to have the same
object ID in Forgejo. That protects destination-only refs and nonmatching
histories from Forgejo's native force-push behavior. Forgejo schedules
normal synchronization.

The Terraform bootstrap is creation-only. It does not observe live drift on
refresh or delete an entry when the Terraform resource is removed. The
GitHub branches were reconciled before HTTPS activation.
