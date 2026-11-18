# Validation

Shared NAS and development storage roots pass pinned-provider checks independently.
The shared root owns accounts, encrypted parent and appliance-wide settings;
the development project owns its namespace, private ACLs, shares and runner cache.
Backend paths, authentication selection and packaged source boundaries were
verified. Neither root reads the other's Terraform state.

The [authorized rollout](../../../../dev/openspec/changes/mount-runner-cache/deployment.md)
created 54 accounts, imported the existing NAS administrator and published all
55 account records. The NAS provider uses the short-lived injected Vault token
without requesting child-token creation rights. Completed import source was
removed. Both NAS roots' final full plans report no changes.

The cache share is src_infra_dev_runner_cache, matching the XCP-ng UNC path.
TrueNAS forbids slashes in SMB names. Namespace ownership, XCP-ng traversal and
cache MODIFY grants, client host restriction, inherited encryption, copies=1 and
1-TiB quota were applied. No adopted data resource was replaced or destroyed.
The cache is operational through XCP-ng and the runner's private guest mount.
