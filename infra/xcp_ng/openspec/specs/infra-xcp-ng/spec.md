# XCP-ng infrastructure Specification

## Purpose

Specify Xen Orchestra resource delegation, scoped authentication, and
certificate deployment owned by `infra/xcp_ng`. The baseline is checked-in
source at revision `550d7e79b1f5fdbc2b6017b75178471d6914082f`, observed on
2026-09-08. Existing operational observations in project documentation are
historical context; this baseline does not verify current deployment, edition,
identity synchronization, or service health.

## Requirements

### Requirement: AppRole resource sets with named inventory

Terraform SHALL declare one Xen Orchestra resource set per entity in Vault's
`approles` group. Each set SHALL contain its explicitly assigned template,
storage repository, and network resolved by names within the selected pool;
unassigned sets SHALL contain no inventory objects. Each set SHALL have a
positive integer CPU quota, defaulting to 32. These resource sets SHALL
represent delegation groups rather than new physical host pools.

Sources: [resource sets](../../../tf/pools.tf),
[named inventory](../../../tf/inventory.tf), and
[project contract](../../../README.md).

#### Scenario: An AppRole has no assigned inventory

- **WHEN** an AppRole exists in Vault's group but has no `resource_set_inventory` entry
- **THEN** its resource set contains no assigned inventory objects and retains
  the configured CPU quota.

### Requirement: Immutable OIDC subjects and explicit existing-VM ownership

Resource-set membership SHALL select the exact synchronized OIDC user by Vault
issuer and immutable AppRole entity UUID. AppRoles without synchronized users
SHALL be reported in `approles_pending_oidc_login` and have no resource-set
subjects. Existing-VM owner ACLs SHALL require exactly one VM matching the
hostname and owning AppRole tag within its named pool, and SHALL grant `admin`
to the owner's synchronized user. Unknown owners, ambiguous VMs, or missing
owner synchronization SHALL fail planning. Native Terraform ACL resources
SHALL separately manage the synchronized administrators group's pool grants.

Sources: [identity contract](../../../tf/README.md),
[owner ACLs and preconditions](../../../tf/authorization.tf),
[resource-set subjects](../../../tf/pools.tf), and
[administrator ACLs](../../../tf/admin_acl.tf).

#### Scenario: A new AppRole has not logged in to XO

- **WHEN** identity discovery cannot find its synchronized OIDC subject
- **THEN** its resource set has no subjects and its name appears in the pending
  login output; an existing-VM ownership grant for it fails its precondition.

#### Scenario: An existing VM name does not identify a unique owned VM

- **WHEN** hostname and AppRole tag matching in the named pool returns zero or multiple VMs
- **THEN** Terraform rejects the existing-VM ownership configuration.

### Requirement: Invocation-scoped tenant authentication

The XO login plugin SHALL authenticate using the calling component's Vault
identity and export `XOA_TOKEN`, `XOA_URL`, and `XOA_INSECURE=false` to the
invoked process. XO HTTPS verification SHALL be mandatory for this plugin.
Session cookies and tokens SHALL remain in process memory, and normal plugin
shutdown SHALL revoke only the XO token issued by that invocation. Tenant
authentication SHALL remain separate from the infrastructure administrator
token. Forced termination or a lost callback can leave a session until XO's
configured expiry and SHALL not be described as guaranteed revocation.

Sources: [XO login contract](../../../cmd/xo_login/README.md) and
[authentication ownership](../../../README.md).

#### Scenario: A component finishes its authenticated invocation

- **WHEN** normal plugin shutdown follows a successful XO login
- **THEN** the plugin revokes its issued session using `token.deleteOwn`
  without revoking unrelated sessions.

### Requirement: Declarative certificate issuance and renewal

The XO appliance playbook SHALL use Certbot with Vault external account
binding, retain the inspected listener certificate paths, install
service-specific Vault CA trust, and configure twice-daily renewal checks.
The host playbook SHALL use its packaged compatibility helper and the
checksum-verified Lego release to install certificates through XAPI's
supported certificate-install command. Temporary EAB registration material
SHALL be removed after use. Appliance renewal may interrupt XO management
access; host renewal SHALL not stop XAPI or running VMs.

Sources: [certificate deployment contract](../../../ansible/README.md)
and [packaged playbooks and runtime inputs](../../../ansible/BUILD.bazel).

#### Scenario: Renew an XCP-ng host certificate

- **WHEN** the renewal helper obtains a certificate different from the installed certificate
- **THEN** it uses `xe host-server-certificate-install`, with installation
  retry tracked independently of issuance, while preserving running VMs.
