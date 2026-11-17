## 1. AppRole storage

- [x] 1.1 Discover AppRoles from their canonical Vault groups and declare NAS accounts; use protected encrypted namespaces and private ownership in AppRole storage projects.
- [x] 1.2 Publish credentials through the existing per-AppRole shared-read paths; verify source policy boundaries and administrator/user group boundaries and coordinated 64-character password rotation.
- [x] 1.3 Give infra/nas/approles/src_infra_dev independent state and plan/apply targets with no Terraform inputs; retain NAS operator authentication, verify backend selection and document bootstrap order.
- [x] 1.4 Remove the shared namespace loop and runner_storage call; put the namespace, traversal/share ACLs and runner_cache solely in the AppRole project, preserving adopted appliance resource addresses.
- [x] 1.5 Namespace the cache SMB share with src_infra_dev, using underscores because TrueNAS forbids slashes; keep the XCP-ng UNC path aligned and validate without deployment.

## 2. Verification

- [x] 2.1 Run formatting, provider, packaging, lint and repository quality checks without deployment; record live prerequisites.
