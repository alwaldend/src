---
title: Vault resource provider
description: Invocation-scoped Vault resources for AL's built-in injector
---

`@rules_al//plugins/vault:config` registers the provider executable and its
runfiles. Compose it with AL configuration when resources need Vault. The
repository's shared `//tools/al:config` already includes this registration.
Builtin-only calls do not discover or launch an otherwise unselected provider.

Startup initializes handlers and advertises `kv`, `op`, `oidc`, `vault_env` and `vault_ssh` without reading
secrets in `PluginStartResponse.resource_types`; its optional result is absent. `ResourceFetch` returns `ResourceResult` directly from opaque settings and declared dependencies;
Vault-specific schema and dispatch remain in this package. The provider uses
generic resource `data` as the write payload for `op` resources; `args` holds
the operation method, path and Vault connection/authentication settings.
The provider uses
AL's configured Vault connection/authentication, including per-resource
`vault_conn` and `vault_auth` settings. Authentication is lazy and issued tokens
remain invocation-owned.

Provider shutdown drains RPCs, revokes owned tokens and removes temporary SSH
material. OIDC requests honor cancellation, verify the configured TLS transport
and reject redirects. Template, authentication and endpoint errors omit secret
values and response contents. `no_auth` explicitly exports an empty `VAULT_TOKEN`
so an inherited token is overridden; it does not remove token-helper files.

The AL-owned local Vault and generic provider fixtures verify selected fetches,
resource templates, dependency order, failures, cancellation and cleanup without
contacting live infrastructure. The parent repository separately tests the pinned
Vault dev-server executable at `//tools/al:vault_dev_server_test`.
