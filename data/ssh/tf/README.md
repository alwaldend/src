---
title: SSH CA snapshots
description: Refresh the data-owned SSH CA public keys from Vault
---

This independent Terraform root reads `ssh/clients/config/ca` and
`ssh/servers/config/ca` through the existing Vault provider, then writes only
their `public_key` fields to `ca_clients.pub` and `ca_servers.pub`. It uses the
existing Local provider; no additional provider or copied Vault KV secret is
needed. The parent BUILD targets consume checked-in files, not this root.

The repository AL wrapper supplies the Vault connection, trusted CA, and
operator authentication. The operator needs read access to both paths and
`auth/token/lookup-self`. The provider uses that existing token without creating
a child token. No infrastructure targets or Terraform state are read. The local
backend path below keeps state under ignored `out/` and shares it between the
plan and apply launchers. Provider locks stay in Bazel runfiles. This root
manages only local snapshots and needs no shared state backend.

From the repository root, review a refresh before applying it:

```sh
mkdir -p out/ssh-ca
export AL_TF_BACKEND_CONFIG_path="path=$PWD/out/ssh-ca/terraform.tfstate"
bazel_agent bazel run //data/ssh/tf:tf.plan -- \
  -var="output_directory=$PWD/data/ssh"
bazel_agent bazel run //data/ssh/tf:tf.apply -- \
  -var="output_directory=$PWD/data/ssh"
```

These are explicit operator commands that contact Vault. The apply command
asks for confirmation. Review and commit the resulting public-key changes.
The absolute output directory avoids writing only into Bazel runfiles. The
existing issuer resources remain owned by infrastructure Terraform, which no
longer writes these snapshots.

Offline checks use `:tf_tests.fmt_test` and `:validate_test`; neither refreshes
the files or contacts Vault.
