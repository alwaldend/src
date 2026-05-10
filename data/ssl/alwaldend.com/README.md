---
title: Alwaldend.com
description: Alwaldend.com certificates
tags:
  - pki
---

## Links

- Generation: https://arminreiter.com/2022/01/create-your-own-certificate-authority-ca-using-openssl/
- Vault: https://developer.hashicorp.com/vault/tutorials/pki/pki-engine-external-ca

## Root CA

Update: manually

{{< readfile file="root_ca.crt" code="true" >}}

## Intermediate CA1

Update: manually

{{< readfile file="ica1.crt" code="true" >}}

## Client certificate CA

Update:
```sh
bazel run //infra/dc1/vault/tf:tf.update_ica_clients
```

{{< readfile file="ica_clients.crt" code="true" >}}
