## Context

Forgejo owns both views of git.alwaldend.com. The global CNAME currently points
to ingress.alwaldend.com; dc1 has an A record for 192.168.10.40. The public
download name resolves to the deployed Yandex host.

## Goals / Non-Goals

Route public web requests to GitHub while preserving local Forgejo. Do not
move DNS ownership, change local clients, or alter the public apex.

## Decisions

Proceed with an in-place CNAME target update to download.alwaldend.com. Reusing
that service alias avoids copying the Yandex IP. Keep the record in this owner
and use its DNS-only Terraform targets; a full setup apply is unnecessary.

The download owner uses a Yandex-only Traefik Host router, redirectRegex
middleware, and noop internal service. Preserve request paths and queries for browser requests. Go discovery
requests with go-get=1 instead use the download owner's loopback Nginx
metadata route for git.alwaldend.com/alwaldend/src and its descendants.
Let’s Encrypt handles the new public hostname; the local Vault issuer is not
asked to issue a git certificate. The existing Nginx backend serves only the fixed discovery metadata.

## Risks / Trade-offs

DNS propagation and initial HTTP-01 issuance have a bootstrap interval. Verify
public DNS, certificate trust, and root/path/query redirects after applying.
The local A record remains unchanged, so local requests still reach Forgejo.
Public web/API/SSH clients using the old endpoint no longer reach that host;
the requested public hostname now identifies the redirect server.

## Failure cases before implementation

Reject plans with changes beyond the global git CNAME. Render both inventory
views and verify no git router locally, no accidental apex enablement, exact
GitHub authority, path/query preservation, unknown-host rejection, and working
download routing. Exercise the actual pinned Traefik in an isolated fixture
before deployment, then check public TLS and redirects after cutover.
