---
title: Vault
description: Setup for vault.dc1.alwaldend.com
---

## Deployment

```sh
bazel run //infra/dc1/vault:deploy
```

## Certificates

Vault certificates (`tls_cert_file`, `tls_key_file`) should be updated by hand
