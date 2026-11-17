# Development NAS storage

This standalone Terraform root owns `default/encrypted/approles/src_infra_dev`,
its private namespace share and the runner cache dataset/share. The namespace
permits XCP-ng to traverse to the cache without connecting to the namespace
share. The cache grants XCP-ng MODIFY access and limits SMB clients to its host.
The reusable dataset module owns encryption inheritance and private ACL defaults.
The cache share is named `src_infra_dev_runner_cache`: SMB names are NAS-wide,
and [TrueNAS forbids slashes in them](https://api.truenas.com/v25.10/api_methods_sharing.smb.create.html).

The root has its own HTTP state backend at
`alwaldend.com/vault1/approles/src_infra_nas/tf_backend/approles/src_infra_dev`.
It authenticates as the existing NAS operator and reuses the shared NAS API
connection. Its `tf=approle_src_infra_dev` selection chooses this backend;
`truenas` selects the connection injector. The development AppRole does not gain
NAS administrator credentials. There are no Terraform input variables or
remote-state reads, and the shared NAS root does not load this project.

## Operation order

Under separately authorized live scope:

1. Apply Vault NAS groups/grants, then `//infra/nas/truenas:tf.apply` for accounts,
   the shared encrypted parent and mandatory SMB encryption. Unlock the parent.
2. Review/apply `//infra/nas/approles/src_infra_dev:tf.plan`/`tf.apply` for this
   namespace and cache.
3. Apply `//infra/xcp_ng/tf:tf.apply` for the storage repository, then the runner
   VM root and Ansible for disk attachment and mounting.

Run the named targets through `bazel run`. Each apply covers its complete root;
no target selection is needed. Other AppRole states are outside this project.
Formatting and offline provider validation use `:tf_tests.fmt_test` and
`:tf_tests.validate_test`.

The account/cache definitions have not been deployed by this task. If an earlier
shared-state revision was applied elsewhere, reconcile its namespace and cache
resources into this state before applying either root. Do not recreate existing
data or treat their removal from the shared source as permission to delete them.
No live plan, apply or state transfer was run for this change.
