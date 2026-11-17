## Purpose

Provide protected filesystem datasets below caller-selected parents while
keeping their configuration and authentication with the infrastructure owner.

## ADDED Requirements

### Requirement: Explicit parent and child configuration

The module SHALL accept parent and config inputs, create one direct filesystem
child, inherit parent encryption, prohibit destruction and return its name and
mountpoint. The module SHALL supply compression/atime/ACL defaults, account ownership and
private NFSv4 ACLs. Named MODIFY access SHALL inherit; named TRAVERSE access
SHALL apply only to that path and grant no namespace share connection. SMB share
ACLs SHALL allow only the owner and MODIFY delegates. Optional SMB config SHALL create an authenticated
DEFAULT_SHARE after permissions, applying an explicit host allow list when given.
Callers SHALL not repeat these resource declarations or defaults.

#### Scenario: Compose namespaces

- **WHEN** a caller supplies one module's dataset name as another's parent
- **THEN** Terraform creates the parent before its configured child

#### Scenario: Reject unsafe names

- **WHEN** the parent is invalid or the child name is empty, a dot component or contains a slash
- **THEN** validation fails before creating a dataset
