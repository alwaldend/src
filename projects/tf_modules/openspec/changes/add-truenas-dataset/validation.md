# Validation

The shared NAS root validates the shared encrypted parent. The independent NAS
AppRole root validates namespace and cache calls, using its own state and targets.
Both package the reusable dataset module directly, without referencing each
other's Terraform files. The VM root retains only Vault/XO providers.

The reusable implementation is unchanged. Its provider-free synthetic render in
out/runner-cache/rendered_acl.json still covers owner inheritance, non-inherited
namespace traversal, inherited cache MODIFY, owner-only namespace shares and
owner/delegate cache shares. Namespace and cache ownership now belongs solely
to the AppRole project's state; shared appliance addresses are unchanged.
No live storage, ACL or state operation ran.
