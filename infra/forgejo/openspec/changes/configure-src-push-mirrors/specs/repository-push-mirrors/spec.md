## Purpose

Maintain GitHub and GitLab copies of the primary Forgejo src repository
through native HTTPS push mirrors.

## ADDED Requirements

### Requirement: Terraform bootstrap of native HTTPS push mirrors

Terraform SHALL bootstrap the configured push mirrors through the Forgejo
API. HTTPS entries SHALL transfer Git LFS with dedicated Vault-backed
credentials, synchronize on commits, and use an eight-hour interval.

#### Scenario: Adopt the active GitLab mirror

- **WHEN** the expected HTTPS mirror is already configured
- **THEN** Terraform bootstrap verifies the entry without mutating the API
- **AND** the runtime keeps its existing mirror identity

#### Scenario: Reject unsafe native activation

- **WHEN** any destination branch or tag is missing from Forgejo or differs from it
- **THEN** bootstrap rejects creation before sending a mutating API request

#### Scenario: Keep mirror credentials out of state

- **WHEN** the bootstrap needs a destination password
- **THEN** it reads the dedicated Vault field inside the packaged process
- **AND** the password does not appear in Terraform state or command arguments
