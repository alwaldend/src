# Forgejo infrastructure specification

## Purpose

Describe the Forgejo service owned by `infra/forgejo`, including Xen
Orchestra provisioning, Ansible service configuration, and Terraform account
and repository management. The `tf_setup`, `ansible`, and `tf` packages are
implementation stages of this owner. This baseline records source guarantees
and prerequisites; it does not verify live service health or restore earlier
Forgejo data.

Baseline revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`.
Observed: 2026-09-08. Sources: [owner README](../../../README.md),
[BUILD](../../../BUILD.bazel), and the implementation links below.

## Requirements

### Requirement: Scoped Xen Orchestra provisioning and disk checks

Setup SHALL provision Forgejo through Xen Orchestra in the
`src_infra_dc1_forgejo1` resource set using the owner's Vault AppRole and
packaged XO OIDC login flow. The VM SHALL declare 20 GiB boot, 40 GiB
Forgejo, and 5 GiB Traefik disks. Ansible SHALL verify the expected sizes
of `xvdb` and `xvdc` before applying the service roles that use them.

Sources: [setup contract](../../../tf_setup/README.md),
[setup BUILD](../../../tf_setup/BUILD.bazel),
[VM definition](../../../tf_setup/vms.tf), and
[deployment playbook](../../../ansible/playbook_deploy.yaml).

#### Scenario: Detect unexpected Xen disk attachments

- **WHEN** Ansible gathers disks that do not match the expected Forgejo and
  Traefik device names and sizes
- **THEN** its pre-task assertion fails before the service roles run

### Requirement: Canonical service routing with legacy hostname support

Forgejo SHALL configure `https://git.alwaldend.com/` as its canonical root
URL, bind its HTTP service to loopback port 3000, and enable the built-in
SSH service on port 3005. Traefik SHALL route the canonical hostname, the
inventory hostname, and `forgejo.alwaldend.com` to that HTTP service with
the Vault certificate resolver.

Sources: [service variables](../../../ansible/group_vars/all.yaml),
[Forgejo configuration](../../../ansible/files/forgejo.ini), and
[Traefik routing](../../../ansible/files/traefik_dynamic.toml).

#### Scenario: Inspect support for existing Forgejo clients

- **WHEN** the Traefik routing template is rendered with the owner's variables
- **THEN** the canonical and legacy hostnames resolve to the same configured
  loopback service
- **AND** generated Forgejo URLs use the canonical Git hostname

### Requirement: Guarded Vault OIDC bootstrap

The deployment playbook SHALL create a missing `vault` OIDC authentication
source after configuring the Forgejo service. It SHALL preserve one existing
active OAuth2 source and reject multiple matches, an inactive source, or a
source of another type. Creation SHALL obtain client credentials from the
injected environment and suppress secret-bearing task output.

Source: [deployment playbook](../../../ansible/playbook_deploy.yaml).

#### Scenario: Preserve an existing active authentication source

- **WHEN** Forgejo reports exactly one active OAuth2 source named `vault`
- **THEN** the bootstrap accepts it and skips source creation

#### Scenario: Reject an ambiguous authentication source

- **WHEN** more than one `vault` source is discovered, or its type or active
  state is incompatible
- **THEN** the bootstrap assertion fails rather than selecting or replacing
  an arbitrary source

### Requirement: Vault identities and service grants with catalog named roles

Service Terraform SHALL require a verified positive integer Vault OAuth
source ID, discover login users from the owning Vault group through at most
two nested group levels, and use entity UUIDs as external login names.
Managed accounts SHALL be protected from deletion. The shared repository
catalog SHALL own named organization administrator and developer assignments
and organization-owned repository identities. Catalog members and Vault
access-group members SHALL belong to the discovered login population.
Existing Vault service-administrator, package-writer, and automation-writer
grants SHALL be retained; the `src` automation writer group SHALL resolve to
exactly one user. A catalog developer SHALL NOT also receive administrator
access through the retained Vault groups.

Sources: [service Terraform contract](../../../tf/README.md),
[users](../../../tf/users.tf),
[access validation](../../../tf/access.tf), and
[repository access](../../../tf/alwaldend_repos.tf).
The [shared catalog contract](../../../../repos/README.md) owns repository
naming and named-role assignments; its
[adoption change](../../../../repos/openspec/changes/archive/2026-09-13-adopt-shared-repository-catalog/design.md)
records the source migration and pending verification.

#### Scenario: Reject unsupported group membership

- **WHEN** a login group extends beyond two nested levels, an access-group
  member is outside the login population, or the automation writer count is
  not one
- **THEN** the corresponding Terraform condition rejects the configuration
- **AND** Terraform does not silently omit the unsupported membership

#### Scenario: Resolve a catalog member through Vault

- **WHEN** a catalog administrator or developer is assigned organization access
- **THEN** that assignment uses the discovered account's existing Vault entity
  UUID and verified OAuth source mapping
- **AND** the assignment does not create a password-based replacement account

#### Scenario: Retain service access while restricting a named developer

- **WHEN** shared named roles are adopted
- **THEN** existing Vault service administration, package-writing, and
  automation-writer grants remain configured
- **AND** the catalog developer receives feature-branch and pull-request
  access without a default-branch push or merge bypass
- **AND** an overlapping administrator grant for that developer fails
  validation

### Requirement: Organization-owned CI action pull mirrors

Forgejo SHALL support catalog-declared organization-owned pull mirrors sourced
from each record's original upstream URL. Terraform SHALL own creation and
mirror synchronization settings. The checkout action mirror SHALL have its own
workflows disabled, and repository CI SHALL consume it at an immutable commit.
Mirror refreshes SHALL NOT implicitly update the workflow's action pin.

#### Scenario: Create the checkout mirror

- **WHEN** the authorized mirror deployment is applied
- **THEN** `alwaldend/com_github_actions_checkout` is a public pull mirror of `https://github.com/actions/checkout`
- **AND** the existing workflow action commit is available from the mirror
- **AND** the deployment does not recreate existing repositories or create copies on unselected forges

#### Scenario: Consume a mirrored action

- **WHEN** repository CI checks out its source revision
- **THEN** its checkout implementation comes from the Forgejo mirror at the declared immutable action commit
- **AND** the source revision remains the workflow event SHA
