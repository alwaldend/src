---
title: XO login
description: Invocation-scoped Xen Orchestra authentication using Vault OIDC
---

This AL plugin logs in to Xen Orchestra using the calling component's Vault
identity. Configure `xoa_url` as an HTTPS origin and `discovery_url` as the
Vault OIDC provider discovery URL. Optional `vault_conn` and `vault_auth`
select the component's existing AL Vault connection and authentication.

The plugin exports `XOA_TOKEN`, `XOA_URL`, and `XOA_INSECURE=false` to the
invoked process. XO HTTPS verification is mandatory. Vault authorization uses
the configured Vault HTTP transport; credentials are never sent to the XO
authorization callback. Session cookies and tokens remain in process memory.

Plugin shutdown revokes only its issued token through `token.deleteOwn`.
Forced termination or loss of the callback response can leave a session until
XO's configured session expiry. Resource-set membership and any existing VM
ACLs must be provisioned separately by the infrastructure administrator.

Bootstrap an approved AppRole's XO identity and synchronized groups before
binding its resource set:

```sh
XO_BOOTSTRAP_APPROLE=src_infra_dc1_forgejo1 bazel_agent bazel run //infra/xcp_ng/cmd/xo_login:bootstrap
```

The target authenticates through existing Vault AppRole policies, reports only
callback completion, then revokes its XO session. Repeat for each approved
AppRole. It does not load the administrator XO token or Terraform backend.
`bootstrap.lua` is loaded at runtime through AL's `--config` flag; the selected
role is not baked into a Bazel configuration artifact. This operation creates
external identities/groups on first login and requires deployment authority.

After bindings exist, add `XO_EXPECT_RESOURCE_SET` and optionally
`XO_EXPECT_VM_ID` to the bootstrap environment. The report then verifies the
authenticated user is not a global administrator, sees exactly the named
resource set, and has no explicit object permissions outside the specified VM.
When a VM is specified, it must have an administration permission and belong
to that resource set. Without a VM, no object permissions are allowed. API
errors fail verification; they are never treated as evidence of denied access.
Only verification booleans and counts are printed.

For an AppRole whose SecretIDs can only be issued by another existing AppRole,
set `XO_BOOTSTRAP_ISSUER_APPROLE` explicitly, for example `src_infra_flux` when
bootstrapping `src_infra_flux_git`. The plugin authenticates the issuer using
the normal AL policy, creates a one-use target SecretID, then logs into XO as
the target identity. It revokes both invocation-owned Vault tokens on shutdown
and cleans unused SecretIDs on login failure. No issuance policy is broadened.
