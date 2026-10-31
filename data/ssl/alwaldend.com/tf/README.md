---
title: PKI CA snapshots
description: Refresh the data-owned PKI CA certificates from Vault
---

This independent Terraform root reads `pki/ica_clients/cert/ca` and
`pki/ica_servers/cert/ca` through the existing Vault provider, then writes only
their `certificate` fields to `pki_ca_clients.crt` and `pki_ca_servers.crt`.
It uses the existing Local provider; no additional provider or copied Vault KV
secret is needed. The parent BUILD targets consume checked-in files, not this
root. The manually maintained root and intermediate CA1 certificates are not
managed here.

The repository AL wrapper supplies the Vault connection, trusted CA, and
operator authentication. The provider uses that existing token without creating
a child token; it needs `auth/token/lookup-self` access. No infrastructure
targets or Terraform state are read. The local
backend path below keeps state under ignored `out/` and shares it between the
plan and apply launchers. Provider locks stay in Bazel runfiles. This root
manages only local snapshots and needs no shared state backend.

From the repository root, review a refresh before applying it:

```sh
mkdir -p out/pki-ca
export AL_TF_BACKEND_CONFIG_path="path=$PWD/out/pki-ca/terraform.tfstate"
al bazel run //data/ssl/alwaldend.com/tf:tf.plan -- \
  -var="output_directory=$PWD/data/ssl/alwaldend.com"
al bazel run //data/ssl/alwaldend.com/tf:tf.apply -- \
  -var="output_directory=$PWD/data/ssl/alwaldend.com"
```

These are explicit operator commands that contact Vault. The apply command
asks for confirmation. Review and commit the resulting certificate changes.
The absolute output directory avoids writing only into Bazel runfiles. The
existing issuer resources remain owned by infrastructure Terraform, which no
longer writes these snapshots.

Offline checks use `:tf_tests.fmt_test` and `:validate_test`; neither refreshes
the files or contacts Vault.
