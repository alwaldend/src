## Purpose

Maintain GitHub and GitLab copies of the primary Forgejo src repository
through native HTTPS push mirrors.

## ADDED Requirements

### Requirement: Terraform bootstrap of native HTTPS push mirrors

Terraform SHALL bootstrap the configured push mirrors through the Forgejo
API. HTTPS entries SHALL transfer Git LFS with dedicated Vault-backed
credentials, synchronize on commits, and use an eight-hour interval.

#### Scenario: Adopt a matching mirror

- **WHEN** the expected HTTPS mirror is already configured
- **THEN** Terraform bootstrap verifies the entry without mutating the API
- **AND** the runtime keeps its existing mirror identity

#### Scenario: Reject unsafe native activation

- **WHEN** the configured trunk is missing or any destination branch or tag is missing from Forgejo or differs from it
- **THEN** bootstrap rejects creation or replacement before sending a mutating API request

#### Scenario: Keep mirror credentials out of state

- **WHEN** the bootstrap needs a destination password
- **THEN** it reads the dedicated Vault field inside the packaged process
- **AND** the password does not appear in Terraform state or command arguments

### Requirement: Restrict mirrored refs

Configured push mirrors SHALL synchronize the catalog's default branch,
branches under `releases/`, and all tags.

#### Scenario: Reconcile an unfiltered mirror

- **WHEN** an existing HTTPS mirror has the expected destination and schedule but no branch filter
- **THEN** Terraform replaces the mirror with the configured filter after the ref preflight passes
- **AND** the next Forgejo synchronization prunes previously mirrored development branches from the destination
