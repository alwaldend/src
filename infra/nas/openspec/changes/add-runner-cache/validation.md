# Validation

The shared NAS and standalone development storage roots pass pinned-provider
formatting/validation independently. The shared root owns accounts, the encrypted
parent and appliance-wide settings; it has no AppRole namespace loop or cache
module call. The AppRole root owns its namespace, ACLs, shares and cache without
input variables or remote-state reads. Cache creation depends on namespace output.

The cache share is src_infra_dev_runner_cache, matching the XCP-ng UNC path.
The [TrueNAS 25.10 API](https://api.truenas.com/v25.10/api_methods_sharing.smb.create.html)
forbids slashes in share names, so the AppRole separator uses an underscore.
Affected NAS/XCP-ng provider checks and the existing storage integration fixture
pass. No live share rename or retained SR mount reconciliation was run.

Generated AL configuration selects exactly one backend per root: the unchanged
core path or src_infra_nas/tf_backend/approles/src_infra_dev. Both select the NAS
operator and shared connection injector; DNS injection is excluded. Both plan
and apply launchers were regenerated without executing Terraform. Their manifests
and materialized runfiles match source and exclude the other root's .tf files.
Sanitized backend evidence is in out/runner-cache/independent-backends.json.

The first core validation caught a stale SMB service dependency on the removed
namespace loop. It now depends on mandatory transport encryption; the corrected
root validates. Account publication no longer depends on namespace provisioning.
Vault membership, 64-character password rotation, private ACL defaults and the
adopted pool, seven dataset and app declarations remain unchanged.

Apply Vault groups/grants, shared NAS settings/accounts, the selected AppRole
storage project, XCP-ng, VM provisioning and Ansible in order. Namespace creation
belongs to its project, not group membership. No live plan, apply or state move
ran. If an earlier shared-state revision was deployed elsewhere, reconcile its
namespace/cache ownership into AppRole states before applying these definitions;
see the transition recorded in design.md. Runtime access and SMB negotiation
remain deployment acceptance checks.
