---
title: Command and plugin lifecycle
description: Startup, shutdown, and secret cleanup guarantees
---

`al run` starts selected plugins, passes their environment to a command, and
keeps plugin resources available until that command exits. On cancellation it
asks the command to terminate and waits up to ten seconds before killing it.
Plugin shutdown then drains requests, releases resources in reverse registration
order, and waits for plugin exit. Startup failure cancels sibling starts and
rolls back partially initialized resources. Cleanup errors fail the invocation.

The injector creates secret files with owner-only permissions and removes them
after consumers stop. Vault stores revoke AppRole tokens they issue; failed
AppRole login attempts destroy their unused single-use SecretID. The role's
bootstrap policy must permit `secret-id-accessor/destroy`, as declared by
`projects/tf_modules/vault_approle`. The user's existing token-helper credential
is never erased or revoked. Secret inputs and remote response bodies are omitted
from diagnostic messages. OIDC authorization uses the configured Vault HTTPS
origin and TLS transport, rejecting alternate origins and redirects.
Explicit config dumps still contain the requested configuration; file outputs
use owner-only permissions and truncate previous contents.

`no_auth` prevents AL's Vault client from loading credentials. For environment
injection it clears inherited `VAULT_TOKEN`; it does not sandbox the command or
prevent a command from independently reading the user's token helper. Plugins
and invoked commands remain trusted programs with the user's filesystem access.

Cleanup is observable best effort, not secure erasure. Memory copies are not
zeroized. Forced termination, host failure, unresponsive cleanup code, or an
unavailable service can prevent deletion or revocation. Subprocess signaling
targets direct children, not arbitrary descendants. Credential expiry remains
a fallback; service-specific logout requirements are documented with each plugin.
