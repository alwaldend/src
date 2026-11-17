# TrueNAS

This root declares the existing `default` RAIDZ2 pool, seven datasets, and
Nextcloud's catalog version and running state. Edit the settings directly in
[`storage.tf`](storage.tf) and [`apps.tf`](apps.tf). There is no separate
configuration document or variable interface.

The small inline dataset table and one-element pool/app indexes preserve the
addresses already imported into Terraform. Pools, datasets, and the app have
deletion protection. Application configuration and encryption secrets remain
on the appliance.

The official provider is pinned to 1.9.3. The observed appliance version is
25.10.6. Credentials come from the NAS owner's Vault `truenas` entry, whose
`endpoint` and `api_key` fields populate the provider environment. TLS
verification is enabled. State uses the owner's `tf_backend/truenas` reference;
the existing DNS root and state are separate.

Review the live state with:

```sh
bazel run //infra/nas/truenas:tf.plan
```

Apply is a separately authorized operation through `:tf.apply`. The existing
resources were imported without a NAS configuration apply; the
[adoption receipt](../openspec/changes/add-truenas-management/validation.md)
records verification.

Offline checks require neither Vault nor the NAS:

```sh
bazel test //infra/nas/truenas:tf_tests.fmt_test //infra/nas/truenas:tf_tests.validate_test
```
