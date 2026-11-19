## Why

The repository has no Nexus infrastructure owner for caching Python packages,
npm packages, and Docker Hub images. A reproducible native Nexus 3 deployment
will provide those three caches on one XCP-ng VM with PostgreSQL, using the
repository's existing provisioning, authentication, and configuration stages.

## What Changes

- Add `infra/nexus` with component documentation, AL configuration, Ansible
  deployment inputs, canonical DNS declarations, and separate Terraform roots.
- Add a reusable Nexus Ansible role for a pinned native distribution, its
  bundled Java runtime, systemd service, persistent data, readiness checks,
  and idempotent initial administrator bootstrap.
- Extend the existing PostgreSQL role as needed to initialize and configure
  a native database on the same VM through the linked collection-owned change.
- Define `tf_setup` for one XCP-ng VM, its disks and network, and DNS through
  the shared DNS module. Wire a dedicated Vault AppRole and XO resource set
  through their existing infrastructure owners.
- Define `tf` for Nexus API configuration, authenticated client access,
  filesystem blob storage, and exactly three initial proxy definitions:
  `pypi-proxy`, `npm-proxy`, and `dockerhub-proxy`.
- Keep TLS termination and ingress in the separate Traefik role. The Nexus
  role does not install or configure a load balancer or container engine.
- Add pinned Nexus distribution/provider inputs, offline validation, and
  documented deployment, upgrade, and recovery procedures during implementation.

Git/GitHub proxying, hosted or group repositories, application SSO, custom
Nexus plugins, multi-node availability, and live deployment are outside this
change. The present handoff contains planning artifacts only.

## Capabilities

### New Capabilities

- `infra-nexus`: Reproducible single-VM provisioning, native Nexus/PostgreSQL
  deployment, Vault-backed stage wiring, bootstrap, and lifecycle boundaries.
- `nexus-proxy-repositories`: Terraform-managed PyPI, npm, and Docker Hub
  proxy configuration and authenticated client acceptance criteria.

### Modified Capabilities

None in the Nexus workspace. Existing Vault, XO, DNS, and provider-pinning
contracts remain with their current owners. The companion
[extend-native-postgresql-role](../../../../../projects/ansible_collection/openspec/changes/extend-native-postgresql-role/proposal.md)
change introduces the `ansible-postgresql` capability under the collection
owner, with its own requirements and acceptance tasks for the existing role.
The Nexus specs describe the deployment's consumption of that contract.

## Impact

Implementation will touch `infra/nexus`, the reusable Ansible collection,
`infra/vault/tf`, `infra/xcp_ng/tf`, and the owning external-dependency pins.
The planning workspace is registered with repository OpenSpec validation.
Secrets stay in Vault and reach each stage through AL injection. Setup and
service Terraform use separate Vault HTTP backend state paths.

Acceptance requires structural validation, rendered service/configuration
fixtures, idempotent bootstrap coverage, Terraform validation against pinned
providers, and successful client/cache scenarios in an explicitly authorized
disposable environment. Offline checks and the proposal itself do not prove
live deployment or authorize any apply, playbook execution, or secret write.
