---
title: Runner credential bootstrap
description: Creation-only Terraform bootstrap with credentials cached in Vault
---

The [service Terraform](../../tf/README.md#runner-credentials) packages this Go
tool and invokes `runner_registration ensure` once for each new bootstrap
resource. It registers a repository-scoped runner through Forgejo's HTTP API and
stores the returned UUID/token directly in Vault KV v2. It does not install a
runner, configure a host, rotate credentials, or delete registrations.

## Execution and cache

`terraform_data.runners` has stable keys and no replacement triggers. Normal
plans and applies do not run the command again. Editing the executable or
provisioner does not invalidate the registration. Explicit replacement invokes
the command again, which reuses a complete matching Vault record without
contacting Forgejo. Removing the resource retains the runner and Vault record.

The tool receives the invocation's `FORGEJO_HOST`, `FORGEJO_API_TOKEN`,
`VAULT_ADDR`, `VAULT_TOKEN`, and optional Vault client configuration through
the existing AL flow. Non-secret environment inputs select the target:

| Variable                     | Meaning                          |
| ---------------------------- | -------------------------------- |
| `FORGEJO_RUNNER_REPOSITORY`  | Forgejo `owner/repository` scope |
| `FORGEJO_RUNNER_NAME`        | Stable runner name               |
| `FORGEJO_RUNNER_VAULT_MOUNT` | KV v2 mount                      |
| `FORGEJO_RUNNER_VAULT_PATH`  | Credential path below that mount |

Credentials never appear in arguments, stdout, Terraform resource attributes,
or temporary files. The cached record contains `instance`, `repository`, `name`,
`id`, `uuid`, and `token`. Runner configuration consumers use `instance`, `uuid`
and `token`; obtain them through their own scoped Vault injection. The tool
requires validated HTTPS, refuses redirects, limits API response sizes, bounds
the entire operation to one minute, and suppresses upstream error bodies.

## Failure and reconciliation

Use this command only through the owning Terraform stage, whose backend lock
serializes registration. Concurrent standalone invocations are unsupported.
A cached record with a different target or missing fields fails validation;
permission and transport failures never count as an absent cache.

When the cache is absent, the tool checks every runner-list page for a
same-name runner before creation. It stores credentials with Vault CAS=0 so it
cannot overwrite an existing secret. If creation succeeds but its response or
Vault write fails, Terraform marks the provisioner failed. A subsequent attempt
detects the existing runner and stops rather than creating another one.
The token cannot be recovered through a runner GET. Preserve the remote runner
and Vault history until an operator explicitly authorizes credential restoration,
adoption, or cleanup. This tool does none of those operations automatically.

Changing a declaration's repository or name does not reconfigure an existing
runner. Cached credentials also remain cached if somebody deletes or changes
the remote runner. Reconciliation and rotation are explicit operations.

## Acceptance

```sh
bazel test //infra/forgejo/cmd/runner_registration:acceptance_test
```

The fixture executes the actual binary against private TLS Forgejo/Vault
services with public dummy credentials. It also executes the production
Terraform provisioner declaration with built-in resources and a private local
state directory: plan has no side effects, initial apply creates once, a changed
tool does not rerun, explicit replacement reuses Vault, and removal retains the
registration. It checks credential exclusion from output and fixture state.
Sanitized JSON receipts are written to Bazel's undeclared test outputs. No live
backend, credential, registration or host is used.
