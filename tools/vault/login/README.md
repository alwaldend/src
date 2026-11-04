---
title: Vault certificate login
description: Go command for Vault login with a PKCS#11 hardware certificate
tags:
  - vault
  - login
---

The Go command authenticates once using the selected hardware certificate, then
atomically replaces `~/.vault-token` with mode `0600`. Authentication or response
validation failures preserve the previous token. Tokens and server response
bodies are never printed. Curl's normal TLS diagnostics and provider PIN prompts
remain visible on stderr; verbose HTTP logging is not enabled.

## Usage

```sh
bazel run //tools/vault/login -- <role-name>
```

The existing AL wrapper injects `VAULT_ADDR` and `VAULT_CACERT` without requiring
an existing Vault token. The address must use HTTPS. The role name is the name
configured in Vault's certificate authentication method.

Automatic discovery ignores trust modules and examines public certificates on
hardware tokens. It excludes CA certificates, expired or not-yet-valid
certificates, and certificates unsuitable for TLS client authentication. It
prefers certificates whose common name or DNS identity matches the role name
(exactly or as the first DNS label), choosing the newest valid matching
certificate. A sole usable hardware certificate can also be selected when the
role has a different name. Selection does not depend on token enumeration order.
The private key is selected by the certificate's object ID; a certificate's
label is not reused as a private-key label.
To override automatic selection, provide the certificate's PKCS#11 URI:

```sh
bazel run //tools/vault/login -- <role-name> \
  --certificate-uri 'pkcs11:token=example;id=%01;type=cert'
```

Provider mode (`--certificate-type PROV`) is the default and supports OpenSSL 4.
Use `--certificate-type ENG` only for an older engine-based installation. Both
certificate and private-key loading use the selected mode. PIN handling belongs
to the native PKCS#11 provider and its controlling terminal; do not put a PIN in
the certificate URI. No global OpenSSL configuration is changed.

The Go binary invokes the existing host `p11tool` for automatic discovery and
`curl` for PKCS#11 TLS signing. It does not invoke a shell, jq, grep, awk, or tail.
Curl configuration files are disabled for this request; certificate settings,
HTTPS, JSON encoding, bounded response size, and the request timeout are explicit.
Login does not follow redirects or retry. Raw trailing curl arguments have been
replaced by typed options. `--timeout` defaults to one minute; `--token-file`
selects a different destination. Existing symbolic-link destinations are rejected.

## Offline validation

```sh
bazel test //tools/vault/login:login_test
bazel build //tools/vault/login
```

The command test uses native-tool fixtures written in Go and synthetic tokens.
It verifies provider/engine selection, hardware and role selection, certificate
validity and usage filtering, response rejection, timeout,
secret-free diagnostics, token preservation, and private atomic replacement.
Its undeclared test outputs contain `login-results.json`. These checks do not
prove a live hardware login; that requires the real certificate and Vault role.
