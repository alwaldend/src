# TrueNAS dataset

Creates a protected child dataset and its private account ownership from
`parent` and `config`. The module owns lz4 compression, atime OFF, owner-only NFSv4 ACLs and inherited encryption. Owner defaults to root; copies
follow the parent and quota is unset unless overridden. Set `traverse = true` for a shared traversal parent.
The caller supplies the provider and deployment authorization.

Optional `smb` creates a non-browsable, authenticated DEFAULT_SHARE after
filesystem ACLs are applied. A share ACL allows only its owner and MODIFY
delegates. Its name defaults to the dataset child name; `hosts`
restricts access to those clients and denies everyone else when nonempty.
Global SMB transport encryption remains the NAS root's responsibility.

See the [runner configuration](../../../infra/nas/approles/src_infra_dev/runner_cache.tf)
for a complete consumer. That AppRole project has its own state and targets.

No separate owner lookup, permission resource or share resource is needed.
Passing a parent module's name establishes ordering. This module does not create
missing ancestors or accept nested paths as child names.

`access` grants named accounts MODIFY (inherited) or TRAVERSE (this path only).
Traversal-only accounts cannot connect to the dataset's SMB share. Named owner
access is inherited so files created by delegates remain accessible to the owner.
MODIFY ACL changes apply recursively within this dataset, without crossing child
dataset mounts. The parent namespace owner must separately allow traversal to
reach a shared child; the AppRole storage project owns that grant for the runner cache.
