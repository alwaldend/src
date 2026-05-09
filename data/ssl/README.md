---
title: Ssl
description: SSL certificates
---

## Links

- Generation: https://arminreiter.com/2022/01/create-your-own-certificate-authority-ca-using-openssl/
- Vault: https://developer.hashicorp.com/vault/tutorials/pki/pki-engine-external-ca

## Root CA (crt)

{{< readfile file="alwaldend.com_root_ca.crt" code="true" >}}

## Root CA (pem)

{{< readfile file="alwaldend.com_root_ca.pem" code="true" >}}

## Intermediate CA1

{{< readfile file="alwaldend.com_ica1.crt" code="true" >}}
