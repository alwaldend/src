## Why

The existing PostgreSQL role installs packages and starts a service but does
not initialize a cluster or manage application database access. The planned
[Nexus deployment](../../../../../infra/nexus/openspec/changes/add-nexus-deployment/proposal.md)
needs those reusable guarantees recorded and validated by the collection owner.

## What Changes

- Correct service-variable handling and document compatibility with existing
  inventories.
- Add explicit, idempotent cluster initialization and configurable service
  settings with guards against reinitializing or replacing existing data.
- Add opt-in authentication, application roles, owned databases, and extension
  management using supplied credentials without recording secrets in source.
- Preserve existing consumers through conservative defaults, migration
  documentation, and acceptance scenarios for already initialized hosts.

This handoff contains planning artifacts only. It does not change the role,
install packages, run playbooks, or authorize a deployment.

## Capabilities

### New Capabilities

- `ansible-postgresql`: The reusable native PostgreSQL role's initialization,
  authentication, database ownership, extension, and compatibility contracts.
  This newly specifies behavior of an existing role.

### Modified Capabilities

None. The existing `project-ansible-collection` capability owns packaging and
locked dependency assembly; those requirements remain unchanged.

## Impact

Implementation belongs in `projects/ansible_collection/roles/postgresql`,
its documentation and checks, and the collection dependency lock if additional
modules are required. Nexus supplies its deployment-specific database,
extension, storage, listener, and credential inputs through its own owner.
The collection's registered OpenSpec validation includes this companion
change; its implementation acceptance remains separate from Nexus acceptance.
