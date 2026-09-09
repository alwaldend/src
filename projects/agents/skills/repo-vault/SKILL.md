---
name: repo-vault
description: Operate Vault through this repository's Bazel-managed vault targets, including token refresh, KV access, login, and status. Use for vault token or authentication problems and Vault CLI operations; do not use for Terraform, Ansible, or secret-value routing handled by repo-secrets.
---

# Work with Vault

## Use repository targets, not a host vault binary

Vault access on this host is defined by `//tools/vault/cmd/workspace:vault`,
which wraps the `al` runner and a Vault injector plugin. The injector resolves
`VAULT_ADDR`, `VAULT_CACERT`, and `VAULT_TOKEN` from the repository's Vault
configuration; the host token file supplies the default token. Prefer the
generated aliases from the root `BUILD.bazel`:

```sh
bazel_agent bazel run //:vault.status
bazel_agent bazel run //:vault.token_lookup
bazel_agent bazel run //:vault -- <vault-cli-arguments>
```

The `//:vault` passthrough runs the Vault CLI with arbitrary arguments, such as
KV, policy, or auth subcommands, without defining new Bazel targets for one-off
operations.

## Refresh the host token

The host token is issued by certificate authentication. Prefer re-issuing it
through the certificate login target:

```sh
bazel_agent bazel run //tools/vault/login -- <username>
```

This target runs with the `default_no_auth` Vault environment and writes a
fresh token to the host token file; it never needs or prints the current
token. It requires the client certificate to be reachable through PKCS#11, so
a YubiKey or equivalent token must be present. Do not confuse it with
`//:vault.login`, which uses userpass authentication.

When the certificate is unavailable but `//:vault.token_lookup` shows a valid,
renewable token, extend it with:

```sh
bazel_agent bazel run //:vault -- token renew
```

The normal `//:vault` wrapper supplies the current token-helper bearer token;
token renewal does not require the issuing client certificate.

## Handle authentication failures

When a Vault operation fails with an authentication or client-certificate
error:

1. Confirm the server is reachable and unsealed with `//:vault.status`.
2. Inspect token validity and TTL with `//:vault.token_lookup`; this does not
   reveal the token itself, but summarizes policies, renewability, and expiry.
3. Determine the issuing authentication method from the `path` field.
4. Extend a valid renewable token with `//:vault -- token renew`; do not
   re-authenticate as a troubleshooting step unless the user explicitly
   authorizes that exact login operation and scope.

Do not write the token value to new files, echo it, or pass it as a literal
command-line argument. The injector already supplies it through the
environment.

## Safety

- Never run Vault writes, token generation, unseal, or deployment merely to
  test a code change.
- Certificate login overwrites the host token file; require explicit user
  authorization for that operation before running it.
- Treat Vault output as potentially sensitive: summarize rather than paste
  responses, and omit token material entirely.
- For secret values and routing decisions, defer to the `repo-secrets` skill.
