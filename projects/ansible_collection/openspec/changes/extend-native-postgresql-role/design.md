## Context

See [proposal.md](proposal.md) for motivation and the linked Nexus consumer.
At `35fda4e1eca61c3dcc7cb1e0be1a929185ac3086`, the
[role tasks](../../../roles/postgresql/tasks/main.yml) only install packages
and start the service using `postgresql`, while
[defaults](../../../roles/postgresql/defaults/main.yml) declare
`postgresql_service`. No cluster initialization, authentication, or database
object management is present. The collection's existing specification covers
packaging rather than this role's behavior.

## Goals / Non-Goals

**Goals:** Provide one reusable role contract with explicit opt-in management,
safe reruns, and compatibility evidence independent of Nexus deployment.

**Non-Goals:** Replication, automatic major-version upgrades, database backup
automation, VM/disk provisioning, Vault writes, and Nexus-specific defaults.
This change remains a proposal until implementation is explicitly requested.

## Decisions

### Keep optional management in the existing role

Extend `roles/postgresql` instead of adding a Nexus-specific database
installer. Keep new cluster, listener/authentication, and database-object
management disabled or empty by default. Consumers opt in explicitly; Nexus
provides its inputs through its own deployment. This avoids silently changing
existing hosts when their collection dependency updates.

Use `postgresql_service` as the documented service selection. Accept the
legacy `postgresql` override during migration, with validation for conflicting
explicit selections and a deprecation message that contains no credentials.
Capture precedence in fixtures before changing the task. Simply dropping the
old variable could switch an existing consumer to the default service.

### Guard initialization before invoking distribution tools

Use the supported distribution's initialization mechanism only after checking
the configured data directory, required mount identity, and cluster version.
Empty validated storage permits first initialization; an existing compatible
cluster is reused. Ambiguous data or incompatible versions fail without repair
or reinitialization. The consumer owns storage provisioning and supplies the
expected path and mount prerequisites. Service startup must respect them.

### Separate host authentication from application ownership

When explicitly enabled, render listener configuration and ordered
authentication rules with bounded reload/restart handlers. Validate required
inputs before replacing files. Use locally privileged administration for
database setup, then verify access through the application's declared
authentication path; a successful administrative connection alone is not
acceptance. Keep secret-bearing tasks and diffs suppressed.

Use idempotent PostgreSQL modules from the pinned collection dependencies for
declared roles, database ownership, and extensions. Verify module availability
and the Python driver against the selected host version during implementation;
add any missing dependency through the collection's lock workflow. Repeated
raw SQL shell commands would make quoting and no-change detection harder to
validate. Do not prune undeclared objects or grant application superuser
access to install an extension.

### Keep acceptance with the collection owner

The [capability delta](specs/ansible-postgresql/spec.md) owns the reusable
contract. Collection checks cover service selection, rendered authentication,
storage failure paths, ownership, and unchanged reruns. Disposable-host
acceptance covers actual initialization and application authentication only
when separately authorized; offline fixtures do not establish those results.
The Nexus change consumes this contract and owns its service integration.

## Risks / Trade-offs

- Distribution initialization and service layouts vary: support and test the
  selected distribution/version explicitly; reject unsupported combinations.
- Authentication changes can remove administrative access: preserve the
  declared local administrative path and validate application success and
  rejection paths in a disposable environment before production use.
- Legacy service overrides can be ambiguous: inventory known consumers and
  test precedence and conflicting inputs before publishing the implementation.
- Data-version mismatch is not recoverable by changing binaries: stop and
  require a separately planned migration with verified backups.

## Migration Plan

First implement and validate the collection-owned tasks after authorization.
Document the service-variable compatibility path and optional inputs. Verify
an initialized consumer with only existing inputs remains unchanged. Then
integrate the Nexus consumer through its linked change and validate the
combined configuration. No live playbook execution is authorized here.

Rollback of source changes does not roll back a database schema or cluster
version. Preserve existing data; any restore or major-version transition
requires its own explicit procedure and deployment authorization.

## Open Questions

- Exact PostgreSQL package version, initialization command, and Python/module
  pins for the supported Nexus guest are selected by the linked deployment's
  compatibility work before implementation finalizes those inputs.
