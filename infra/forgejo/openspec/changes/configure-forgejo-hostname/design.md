## Context

The DNS split and public GitHub redirect already exist. Local Traefik currently
proxies both names directly, and Forgejo advertises the Git alias. Vault already
allows the canonical callback and certificates for both names.

## Decisions

Use a dedicated local Git router and permanent redirect middleware. Keep the
canonical and inventory routers attached to Forgejo. Reuse the existing Git
DNS record and Vault certificate resolver; neither VM provisioning nor DNS
needs a change. Keep certificate SANs for the Git alias after redirecting it.

Use the canonical URL directly in shared AL authentication and the runner
scaffold. Keep the existing Flux SSH URL: its local DNS and port remain valid,
and the HTTP redirect does not affect SSH.

## Failure cases before implementation

Check root, repository, escaped path and query redirects; canonical and inventory
requests must reach the upstream, and unknown or suffix-matching hosts must
return 404. Verify the public GitHub destination and local/public separation
against rendered templates. Inspect rendered HTTP, SSH, root URL, TLS names
and Vault callback derivation. Use only placeholder credentials and isolated
loopback services; no ACME enrollment or live backend calls.

## Migration limits

This task changes desired state only. A later authorized rollout must apply
the Vault and Ansible owners and verify TLS, OIDC login and client URLs.
It must preserve Forgejo state through the guarded upgrade role. Redirects
do not migrate saved remotes or authorize any live operation.
