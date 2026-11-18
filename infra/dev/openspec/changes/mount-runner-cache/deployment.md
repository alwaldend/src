# Authorized rollout · 2026-11-18

The user authorized NAS/cache deployment, Vault state recovery and moving
AppRoles to independent states, while explicitly prohibiting a Vault restart.
Vault was neither restarted nor reloaded. All state pushes used the owning
targets without force. The [Vault migration receipt](../../../../vault/openspec/changes/isolate-approle-state/validation.md)
records the conserved 800 resources and four no-change resource plans.

## Deployed stages

- Vault: the three access additions and three updates completed. State persistence
  was recovered by splitting NAS, development and AL into independent states.
- Shared NAS: 54 accounts were created and the existing src_infra_nas account
  (ID 71) was imported. It now belongs to truenas_admins with FULL_ADMIN.
  All 55 accounts have coordinated 64-character passwords and Vault records.
  Groups, administrator privilege, encrypted parent/ACL and required SMB
  encryption were deployed. The adopted pool, seven datasets and Nextcloud
  were unchanged. The completed import block was removed; a full plan reports
  no changes. The provider uses its injected short-lived token directly.
- Development NAS: eight resources created the namespace, cache, private
  ACLs and shares. The cache inherits encryption, uses copies=1 and a 1-TiB
  quota, and restricts SMB clients to 192.168.1.213. Its full plan has no changes.
- XCP-ng: one bootstrap resource created Runner cache and one resource-set
  update delegated it to development. The scoped apply used
  -target=xenorchestra_resource_set.approles["src_infra_dev"]. A broader apply
  was canceled because it included unrelated OIDC replacement and AppRole sets.
  The final full plan contains only that OIDC replacement and the two unrelated
  AppRole sets: three additions, zero updates, one deletion/replacement. No cache
  drift remains. Those changes remain outside this rollout.
- VM: the existing VM and 500-GiB boot disk were preserved. The 800-GiB cache
  disk attached at position 2. The pinned provider plans new computed disk fields
  as null; a temporary built-in projection also encountered its inconsistent
  final plan. It was removed. A second complete apply published the observed
  /dev/xvdc device and capacity to Vault. The full VM plan now has no changes.
- Ansible: the runner-only selection completed with 48 successful tasks,
  seven changed, nine skipped, zero failed and zero unreachable. It created
  ext4 without forced overwrite, mounted by UUID and required the mount before
  starting the worker. Dev1 was outside the host limit. Repeating the same
  selection returned 47 successful tasks, zero changed, nine skipped, zero
  failed and zero unreachable.

## Runtime acceptance

The guest reports /dev/xvdc, 858993459200 bytes, ext4 at /var/cache/bazel with
nosuid,nodev,noexec. Its disk_cache directory belongs to forgejo-runner with
mode 0700 and is writable by that account. The runner service is active and
RequiresMountsFor includes the cache. Its private bazelrc selects this action
cache and the 700G garbage-collection limit; output bases remain local.

A disposable runner-account probe built the same action with Bazel 9.3.1 in two
distinct output bases. The second build reported one disk cache hit. Both
artifacts had SHA256 2d176c508b13679ad753f15a515fbc465a9140a25406c37bf5603493d6a477bc.
The fixture and output bases were removed. The repeatable script is retained
in ignored task scratch at out/runner-cache/verify_runner_cache.py.

The encrypted SMB repository is operational against mandatory NAS encryption.
Direct negotiation inspection through the owning XCP-ng SSH target was unavailable
because this controller lacks sshpass; no client package was installed for that
inspection. No negotiated cipher or protocol-version claim is made here.

## Diagram and execution notes

CI renders as light/dark SVG and a dark WebP from docs/ci.mmd through the same
site preset as the development-host diagram. Visual inspection verifies the
cache inside the VM, runner-to-disk-to-NAS and runner-to-Forgejo arrows, a
VM-to-VPN/proxy edge and one External deps node. Rendered arrow markers were
checked after the pinned Mermaid renderer discarded a reverse-arrow spelling;
the final source uses a normal directed edge and an invisible layout constraint.

An automatic approval rejection for the initial NAS-storage confirmation was
resolved by canceling it and reviewing a fresh, complete interactive plan before
confirmation. An output-relay signal was initially rejected for uncertain process
identity, then allowed after elevated process inspection established its identity.
No raw plan or state content was retained as documentation. Protected state
recovery copies were removed after verified persistence. These results establish
runner cache readiness, not a new complete repository CI run.
