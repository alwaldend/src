# RouterOS accounts

This Terraform root provisions the `src_infra_printers` RouterOS group and user
on router1. Its group permits `api`, `read`, `rest-api`, and `write`. These are
router-wide configuration permissions: RouterOS does not scope them to the
printer firewall rule. The group excludes account-policy administration,
sensitive-value access, and shell or UI login permissions. RouterOS REST
authentication on the deployed router requires both `api` and `rest-api`.

The pinned Random provider generates a 32-character alphanumeric password once
per resource lifecycle. Ordinary plans preserve it. Deliberate replacement of
`random_password.printers` rotates the router password and the Vault KV v2
credential together. There are no credential outputs. Terraform state contains
the password and remains in the protected Vault backend.

The credential is stored at
`secrets/alwaldend.com/vault1/approles/src_infra_mikrotik/printers`, with fields
`mikrotik_username` and `mikrotik_password`. The MikroTik AppRole owns this
prefix; the printers AppRole receives read access to this exact secret through
the [Vault definition](../../vault/tf/approles/src_infra_printers/main.tf).

## Provisioning prerequisites

The `tf=accounts` wrapper authenticates as `src_infra_mikrotik`, uses its
`tf_backend/accounts` backend, and injects its token for the Vault provider.
RouterOS provisioning reads the operator-supplied administrator credential at
`secrets/alwaldend.com/vault1/approles/src_infra_mikrotik/routeros_admin`, with
fields `mikrotik_username` and `mikrotik_password`. Its existing owner-prefix
policy grants access; no cross-owner credential read is needed. The account
must permit user and group management (`policy` as well as read/write and REST
API access). Printers consumes only its separate generated credential.

Operational calls require explicit authority for their exact scope. With that
authority, review the owning `//infra/mikrotik/accounts:tf.plan` and interactive
`tf.apply` targets. These four resources are provisioned; an ordinary reconciliation now
reports no changes. The [deployment receipt](../../printers/openspec/changes/move-printer1-infrastructure/deployment.md)
records their creation and the completed state adoption.

The shared injection helper supplies `MIKROTIK_HOST`, `MIKROTIK_USER`, and
`MIKROTIK_PASSWORD` directly to the default RouterOS provider. No Terraform
credential variables are needed.

## Offline validation

```sh
bazel test //infra/mikrotik/accounts:tf_tests.fmt_test \
  //infra/mikrotik/accounts:tf_tests.validate_test
```

These checks package the pinned providers and omit only `backend.tf` so schema
validation needs no Vault authentication or router access. They do not create
an account or test its live permissions.
