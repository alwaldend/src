---
title: Forgejo signing-key bootstrap acceptance
description: Verify the shared Vault tool through Forgejo's actual Terraform provisioner
---

`bazel test //infra/forgejo/test/signing_key:acceptance_test` runs the shared
[Vault SSH key tool](../../../../tools/vault/cmd/ssh_key/README.md) and the
actual `terraform_data.commit_signing` declaration against an isolated TLS Vault
fixture. It uses public dummy credentials and private disposable local state.

The fixture checks initial creation, cache reuse, access/history/redirect/CAS
failure, corrupt-record rejection, no effects during planning or normal repeat
apply, and cache reuse on explicit replacement. Private values must remain out
of command output and Terraform state. A sanitized `signing-bootstrap.json`
receipt is emitted into Bazel's undeclared test outputs. No live backend, key,
Vault write, or service is used.
