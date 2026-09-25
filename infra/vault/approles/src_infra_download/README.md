---
title: Download AppRole
description: Independent Vault state for the download identity
---

This standalone Terraform root owns the download AppRole, DNS read policy,
SSH signing role, and internal ACME/EAB role. It has its own HTTP backend at
`alwaldend.com/vault1/approles/src_infra_dc1_vault/tf_backend/approles/src_infra_download`.
It authenticates through the existing Vault administrative AppRole, so creating
the download identity does not require that identity to exist first. Its
`tf=approle_download` label selects only this backend and Vault authentication.

The root looks up the AppRole auth mount and operator entity by path/name.
The shared mounts remain owned by `infra/vault/tf`; this root references the
`secrets`, `ssh/servers`, and `pki/ica_servers` paths. No parent-module context
or remote state is consumed. The central stage looks up the created download
entity and group by name for its own infrastructure and Ansible memberships;
those membership lists keep one Terraform owner.

`src_infra_download` receives its own KV subtree and Yandex folder integration.
SSH host signing covers `download.alwaldend.com` and its descendants. The
internal ACME role permits only `alwaldend.com`, `download.alwaldend.com`, and
`www.alwaldend.com`, with no descendant names and a seven-day maximum TTL.
The AppRole group can request EAB credentials for this role.

## Bootstrap order

Under separately authorized live scope:

1. Establish the existing shared Vault mounts and administrative/operator
   identities through the core Vault deployment.
2. Review/apply `//infra/vault/approles/src_infra_download:tf.plan`/`tf.apply`.
3. Review/apply `//infra/vault/tf:tf.plan`/`tf.apply` to resolve the new identity
   and add its shared group memberships.
4. Apply the XCP-ng/Yandex provider assignments, then the download deployment.

Run the named targets through `bazel_agent bazel run`; AL injects credentials
and backend configuration. The new identity is undeployed, so no state move or
import is required. A core-stage plan using the new name lookups requires
step 2 first; do not use the core stage to bootstrap this identity.

The test-only `:offline` target supports backend-disabled initialization,
validation, and the mocked standalone plan in task-owned scratch without AL
authentication. `:tf_tests.fmt_test` checks formatting. Offline checks do not
establish deployed permissions or successful live lookups.
