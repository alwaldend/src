## Decisions

Keep parent/config as the caller interface. The module owns lz4 compression,
atime OFF, NFSv4 ACLs, inherited encryption, protected lifecycle,
owner lookup and private permissions. Callers provide the child name, account
owner and genuine property overrides. Optional smb config creates one private
DEFAULT_SHARE export; name defaults to the child name and hosts defaults to
unrestricted authenticated access. A nonempty hosts list denies all others.
The shared parent enables public traversal. Named access maps support MODIFY
(inherited by children) and TRAVERSE (this path only). Owner access is inherited
by name as well as owner@ so delegate-created files stay accessible to the owner.
SMB share ACLs allow the owner and MODIFY delegates only. TRAVERSE never grants
connection to the namespace share. One filesystem ACL owns each path; no chmod
resource competes with it.

## Failure cases

Before implementing the refactor, retain rejection of empty/dot/nested names,
check that disabling SMB creates no export, host restrictions survive the
module move, owner resolution follows account creation, and a child share
cannot precede permissions. Keep the runner's cache name, quota, copies and
SMB path unchanged. Preserve adopted pool/dataset/app addresses. Validate both
real provider consumers; no live storage operation is needed.

## Delegation failure cases

Before implementation, verify the hypervisor reads only its own credential;
parent traversal grants neither listing nor share connection; MODIFY is limited
to the cache and inherited by its files; ordinary group membership confers no
access; owner access survives delegate-created files; and delegate account
creation precedes UID lookup. Reject unsupported access names. Apply recursive
cache ACL changes only inside that dataset, never across child dataset mounts.
The adopted pool, seven datasets and app must remain unchanged.
