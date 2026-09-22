# Vault operations in this repository

## Choose the authentication context

The root `//:vault` and `//:vault.*` aliases use
`//tools/vault/cmd/workspace` with `//tools/al:config`. The wrapper injects
`VAULT_ADDR`, `VAULT_CACERT`, and `VAULT_TOKEN`; the default authentication
loads the host token helper. `//:vault -- <arguments>` passes arguments to the
pinned Vault CLI without requiring a new target for a one-off operation.

Component targets such as `//infra/<component>:vault.kv_put` use the package's
`al_config`. Inspect its authentication before selecting the target. The
command map in `tools/vault/defs.bzl` adds `-mount secrets` to the KV aliases;
pass the logical path beneath that mount, without the KV v2 `data/` segment.
Use the full API path, including `<mount>/data/`, for capability checks.

## Separate AppRole bootstrap from KV permission

The AppRole flow in `projects/al/pkg/al/vault.go` has two credentials: the host
token creates a single-use SecretID, then AL exchanges it for an AppRole token
and injects that token into the command. The lifecycle and cleanup guarantees
are owned by `projects/al/docs/README.md`.

`projects/tf_modules/vault_approle/main.tf` owns the policy split: the bootstrap
identity can request a role's SecretID, while the role token receives the
policy for its own KV subtree. A host token's denial on that subtree does not
establish the component token's permissions. Permission to create tokens also
does not establish permission to write a particular KV path.

A failure saying `could not create secret id for the approle` occurs before
the requested Vault CLI command runs. Inspect the declared role and membership
under `infra/vault/tf` and the bootstrap token's capability on
`auth/approle/role/<role>/secret-id`. A 403 alone does not prove the role is
missing: undeployed configuration and insufficient bootstrap policy are
separate possibilities. Even a component `vault.status` target can fail here;
use `//:vault.status` to isolate server reachability from component login.

When implementing a new component, register its role and operator bootstrap
membership in the [owning Vault stage](../../../../../infra/vault/tf/README.md)
using the [shared AppRole module](../../../../../projects/tf_modules/vault_approle/main.tf).
Its own KV subtree covers component credentials and Terraform state. Disable
the optional Yandex Cloud policy when unused, and do not add general
infrastructure group membership merely to enable a credential write.

A checked-in role is not evidence that it has been deployed. If the role or
its policy has not been deployed, prepare and validate the owning
`infra/vault/tf` change. An authorized KV write does not also authorize
that Terraform apply; retain the exact operation and scope boundary from
`AGENTS.md`. Once bootstrap works, inspect the selected role token's capability
on the exact destination through the component passthrough, for example:

```sh
bazel_agent bazel run //infra/<component>:vault -- token capabilities secrets/data/<logical-path>
```

These checks report capabilities without reading a secret. Use the target
whose injected identity owns the destination; do not copy or broaden a policy
merely to bypass a failing authentication context. Keep source validation
moving while a live prerequisite is unavailable.

## Set up a component credential

Derive the logical secret path and field name from the component's `al.lua`
and select its `vault.kv_put` target from the owning `BUILD.bazel`. Once the
AppRole and bootstrap access above are deployed, this target writes using
the component's identity. Its KV alias already selects the `secrets` mount.
For a new single-field credential, replace the placeholders below with those
source-owned names and run this Bash example:

```bash
(
  set +x
  read -rsp 'API token: ' component_token || exit 1
  printf '\n' >&2
  [ -n "$component_token" ] || exit 1
  printf '%s' "$component_token" |
    bazel_agent bazel run '//infra/<component>:vault.kv_put' -- \
      '<logical-path>' '<field>=-'
)
```

The hidden prompt and Bash builtins keep the value out of shell history and
process arguments; the subshell releases the variable when it exits.
[`<field>=-` reads the value from stdin](https://developer.hashicorp.com/vault/docs/commands/kv/put).
`kv put` replaces the secret's fields, so preserve any existing fields when
updating a shared record; use the owning `vault.kv_patch` target when only one
field should change. Confirm the command's result without reading the token
back.

A request for the command authorizes explanation, not execution. Execute a
write or prerequisite apply only within the exact authorization already
granted; do not ask again when it covers that operation and scope. Preparing
the role or explaining this example does not authorize either live operation.

## Inspect token metadata without exposing credentials

The `token_lookup` alias is an unfiltered `token lookup`, whose response can
include the token in `data.id`. JSON output alone does not redact it. For a
single TTL observation, the passthrough can select the field before output:

```sh
bazel_agent bazel run //:vault -- read -field=ttl auth/token/lookup-self
```

When several fields are needed, capture `token_lookup -- -format=json` in
memory and emit only an explicit allowlist from `.data`, such as `ttl`,
`renewable`, `expire_time`, and `policies`, before returning any tool output.
Do not print or retain the original response. Inspect the issuing `path` only
when needed, since it can include a username. For the pinned CLI, selecting
`read -field=renewable` returns the response envelope's flag; use
`.data.renewable` from the filtered lookup for token renewability. See the
[lookup API](https://developer.hashicorp.com/vault/api-docs/auth/token#lookup-a-token-self)
and the pinned CLI's
[field selection](https://github.com/hashicorp/vault/blob/v1.21.4/command/util.go).

## Refresh the host token

When a refresh is authorized, distinguish an expired token from a valid
renewable token. A valid renewable token can be renewed through the root
passthrough without replacing the host token file. This command requests the
same self-renewal endpoint as `token renew` and prints only the resulting TTL:

```sh
bazel_agent bazel run //:vault -- write -force -field=token_duration auth/token/renew-self
```

Renewal does not fix missing roles or insufficient policies. An expired token
needs a fresh login; use the certificate login target when that is the intended
authentication method:

```sh
bazel_agent bazel run //tools/vault/login -- <username>
```

That target uses `default_no_auth`, obtains the client certificate through
PKCS#11, and writes a fresh token to the host token file without printing it.
The certificate must be available, for example through a YubiKey. The generated
`//:vault.login` alias instead uses userpass authentication. Certificate login
overwrites shared host authentication state: preserve any authorization
already given for that refresh, and do not expand an unrelated troubleshooting
request into login or token creation. Follow `repo-secrets` for secret handling
and `AGENTS.md` for live-operation authority.
