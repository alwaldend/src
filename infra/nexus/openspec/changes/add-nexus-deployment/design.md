## Context

See [proposal.md](proposal.md) for motivation and scope. Source discovery at
`a988d4c3561ca6d5e5c0f079f387a3ae8fd76160` found no `infra/nexus` owner or Nexus
provider pin. This change currently contains planning artifacts only.

Relevant existing owners and evidence:

- [Forgejo](../../../../forgejo/README.md) demonstrates native deployment on
  XCP-ng and the setup/Ansible/service split.
- [XO authorization](../../../../xcp_ng/tf/README.md) binds a resource set to
  the exact synchronized Vault AppRole subject, not a broad group.
- [PostgreSQL tasks](../../../../../projects/ansible_collection/roles/postgresql/tasks/main.yml)
  currently only install packages and start a service. They use `postgresql`
  as the service variable while defaults declare `postgresql_service`; they
  do not initialize a cluster or configure a Nexus database.
- [Provider declarations](../../../../../third_party/terraform/README.md)
  own reproducible provider archives and checksums.

The user confirmed PostgreSQL on the Nexus VM and limited the repository set
to PyPI, npm, and Docker Hub proxies. Git proxying was explicitly dropped.

## Goals / Non-Goals

**Goals:** Give each provisioning and configuration action one owner, make
fresh deployment and reruns predictable, and make the three caches usable by
their native clients. Preserve data and credentials through routine reruns.

**Non-Goals:** This planning handoff does not implement the deployment. The
eventual source change does not authorize live operations. It introduces no
new reverse-proxy implementation, database cluster, or automatic backup service.
Upgrade and consistent backup/restore procedures remain necessary documentation.

## Decisions

### 1. Retain the repository's stage and ownership boundaries

Use `infra/nexus/tf_setup`, `infra/nexus/ansible`, and `infra/nexus/tf`.
The reusable service role belongs in
`projects/ansible_collection/roles/nexus`. Extend the existing PostgreSQL role
for reusable database behavior rather than creating a second installer.

The companion
[extend-native-postgresql-role](../../../../../projects/ansible_collection/openspec/changes/extend-native-postgresql-role/proposal.md)
change owns the shared role's initialization, configuration, and compatibility
contract and acceptance tasks. Implement and validate that contract before
integrating the role with Nexus. Both proposals remain planning-only changes;
the companion defines reusable guarantees while this change owns their
deployment-specific consumption.

Declare `src_infra_nexus` in `infra/vault/tf`, including group membership,
scoped DNS access, and SSH signing policy. Add the named pool/template/storage/
network inventory to `infra/xcp_ng/tf`. Do not add an existing-VM ACL lookup
for a VM that does not exist yet; fresh provisioning and later adoption have
different prerequisites. Follow the existing XO first-login synchronization
procedure before provisioning with the new identity.

`al.lua` uses `tf=setup`, `ansible=1`, and `tf=main` selectors, with matching
packaged plugin dependencies and invocation labels. Use separate Vault HTTP
state entries under the component's AppRole path. Avoid a combined Terraform
root: the Nexus provider cannot configure an API before that API exists.

### 2. Use native Nexus Community Edition and PostgreSQL

Plan for Community Edition, which needs no Pro license for these proxies.
Pin an official Linux distribution and checksum together with the Nexus
Terraform provider. Use the distribution's bundled compatible Java runtime
and the repository's supported Fedora cloud-template pattern; check the exact
OS/runtime/provider combination before implementation fixes the pins.

Run the Nexus service under a dedicated user, with explicit working/data
directories, resource limits, bounded readiness retries, and restart handlers.
Use supported Nexus REST APIs rather than enabling arbitrary Groovy scripts.

Configure PostgreSQL on loopback before the first Nexus start, including an
owned database, authentication rules, and `pg_trgm`. Render connection settings
through one supported configuration mechanism into a restricted file. Reuse
Ansible PostgreSQL modules if available in the pinned collection; any missing
dependency follows the collection's lock workflow.

H2 would reduce setup work, but the user selected PostgreSQL and Sonatype
recommends it. An external database would break the agreed single-VM scope.
See [database installation](https://help.sonatype.com/en/install-nexus-repository-with-a-postgresql-database.html)
and [service installation](https://help.sonatype.com/en/run-as-a-service.html).

### 3. Make data and ingress dependencies explicit

Keep OS/application binaries separate from Nexus and PostgreSQL data. Use
dedicated guest data disks/mounts with identity and filesystem checks, safe
initialization, and service mount dependencies. Protect the stateful VM from
accidental replacement. Additional disks inside a VM resource do not by
themselves guarantee persistence when that VM is destroyed.

Assume same-VM Traefik using its existing role. Nexus binds to loopback at
the main HTTP backend; PostgreSQL remains loopback-only. Component inputs
describe the HTTPS hostname and route, while the Nexus role changes no proxy
configuration. The service API route must work before service Terraform runs.

For Docker, use one dedicated HTTPS registry hostname routed by Traefik to
a private Nexus Docker HTTP connector. Service Terraform owns the connector
and Docker Bearer Token Realm setting. Verify the listener's actual bind and
firewall behavior and keep it inaccessible directly from client networks.
This avoids depending on newer path-routing fields before provider support
is verified; one connector is sufficient for the single requested registry.
See [Docker routing](https://help.sonatype.com/en/docker-registry.html).

### 4. Break the credential-bootstrap dependency before Terraform

Store the intended administrator and database credentials in Vault before an
authorized deployment. Ansible first checks whether the intended Nexus
credential authenticates. If it does, bootstrap is already complete. If it
does not, use a valid initial password file for the supported password-change
API, then verify the intended credential. Fail closed when neither works.

Use `no_log` for secret-bearing tasks and leave no secret-bearing task output.
Never substitute a default password or reset credentials by modifying the
database. Ansible owns only this bootstrap transition; routine credential
rotation requires an explicit procedure and does not occur on every rerun.
The bootstrap administrator is not also a Terraform-managed user resource.
AL injects its credential for initial service configuration without committing
it or declaring a Terraform resource for it. See
[initial setup](https://help.sonatype.com/en/installation-methods.html).

### 5. Let Terraform own the three caches and access configuration

Use the community `datadrivers/nexus` provider, verified against the selected
Nexus release. Pin it through the shared provider catalog. Configure the URL
and CA trust explicitly; TLS verification must not depend on provider defaults.
Maintain separate filesystem blob stores for the three proxy caches so their
usage and future retention settings can be controlled independently.

| Name              | Format       | Remote storage                  |
| ----------------- | ------------ | ------------------------------- |
| `pypi-proxy`      | PyPI proxy   | `https://pypi.org/`             |
| `npm-proxy`       | npm proxy    | `https://registry.npmjs.org/`   |
| `dockerhub-proxy` | Docker proxy | `https://registry-1.docker.io/` |

Select Docker Hub index behavior for the Docker repository. Public upstreams
can be used without embedding credentials; allow Vault-injected Docker Hub
credentials when needed for authenticated upstream access. Do not promise
that caching removes upstream rate limits for uncached content.

Disable anonymous access as the initial security default and declare a
read-only role/account for consumers. Keep its credential in Vault; prefer
provider write-only fields when the pinned version supports them. Otherwise
document that sensitive user/upstream attributes may reside in the protected
Terraform state, not merely be hidden in CLI output. Credential rotation
versions must be explicit where required by the provider.

Inventory the selected release's built-in resources. Adopt singleton security
settings through explicit imports where necessary. Keep unrelated pre-existing
repositories outside management and report them; never add a delete-all step.
Do not add hosted/group repositories merely to expose each lone proxy.

See [provider source](https://github.com/datadrivers/terraform-provider-nexus),
[PyPI proxy setup](https://help.sonatype.com/en/create-a-pypi-repository.html),
[npm proxy setup](https://help.sonatype.com/en/npm-registry.html), and
[Docker proxy setup](https://help.sonatype.com/en/proxy-repository-for-docker.html).
These are design references observed on 2026-09-20, not compatibility-test
results for a pinned deployment candidate.

### 6. Validate contracts at the appropriate boundary

This planning workspace is registered with repository OpenSpec validation and
has a component README and documentation/source packaging. During
implementation, add runtime packaging through the existing Bazel rules.
The current BUILD targets do not deploy services or introduce provider pins.

Offline checks cover Terraform structure/provider schema, YAML and template
rendering, systemd configuration, plugin selection, package contents, DNS
ownership, and bootstrap branches using non-secret fixtures. Include the
fresh-host PostgreSQL initialization path and the corrected service variable.
Do not run a normal playbook target as a syntax check or assume check mode is
safe for every task.

When separately authorized, use a disposable environment for installation,
reboot, unchanged rerun, authentication failures, artifact downloads, and cache
reuse. Use fresh client caches and inspect upstream artifact-body requests;
metadata freshness checks and authentication traffic are not cache misses.

## Risks / Trade-offs

- **Single failure domain:** Nexus and PostgreSQL share CPU, RAM, storage,
  and VM availability. Size both services together and document a consistent
  database/blob/configuration backup and restore procedure.
- **Data-loss boundary:** VM replacement or an incompatible database upgrade
  can lose data. Require reviewed replacement, verified backups, and a restore
  path; selecting an older binary is insufficient rollback evidence.
- **Provider drift:** Upstream provider documentation has differing tested
  Nexus versions. Validate the chosen release pair and required resources
  before committing pins; do not infer support from the provider name.
- **Cache growth and rate limits:** Docker layers can exhaust disk rapidly.
  Expose quotas and cache ages, document capacity monitoring, and leave
  destructive retention disabled until criteria are selected.
- **Bootstrap ordering:** XO identity synchronization and HTTPS API ingress
  precede their consumers. Document these prerequisites explicitly and report
  missing access rather than retrying with a broader identity.

## Migration Plan

There is no existing Nexus state to migrate in this checkout. First implement
and validate the source under a later apply request. A separately authorized
deployment follows this sequence:

1. Provision the Vault AppRole, scoped grants, and referenced secret entries.
2. Synchronize its XO identity and apply its resource-set inventory/grants.
3. Review and apply `tf_setup` for the VM and owned DNS records.
4. Run Ansible host/database/service setup and the separate ingress role;
   verify persistent mounts, readiness, bootstrap, and HTTPS API reachability.
5. Review and apply service Terraform, including reviewed adoption where needed.
6. Verify authorized and unauthorized clients, repeated artifact retrieval,
   reboot persistence, and an unchanged convergence run.

Stop at a failed prerequisite. Retain data during recovery and correct the
owning stage; destruction is not a routine rollback. Later upgrades require
compatible backups of PostgreSQL, blob data, configuration, and encryption
material before allowing Nexus to migrate its schema.

## Open Questions

These are deployment values or compatibility selections within the agreed
architecture, not unresolved product scope:

- Exact Nexus/provider/PostgreSQL versions and verified distribution digest.
- Unique VM/service/registry hostnames, address, DNS views, and selected XO
  template/network/storage names. Do not guess an unused address.
- CPU/RAM and OS/Nexus/PostgreSQL disk capacities appropriate to expected use.
- Approved consumer account identifier, cache age/quota values, future retention
  criteria, and whether Docker Hub upstream authentication is needed.

Next action after proposal review: begin task 1 only when the user explicitly
requests implementation. No live state has been inspected or changed here.
