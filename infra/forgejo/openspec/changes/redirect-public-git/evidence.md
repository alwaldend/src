# Public Git redirect acceptance

On 2026-10-26 the owning `tf_setup:dns.plan` produced exactly one in-place
Cloudflare CNAME update: git.alwaldend.com from ingress.alwaldend.com to
download.alwaldend.com, TTL 600, DNS-only. The saved plan was reviewed and
applied through `tf_setup:dns.apply`: zero additions, one change, zero destroys.
The local A record and other Forgejo records were unchanged.

After the linked Yandex Traefik deployment, public resolver 1.1.1.1 returned
51.250.85.209; the controller's local resolver returned 192.168.10.40.
Address-pinned HTTPS probes retained the correct SNI and used the controller's
default public CA trust. Git requests returned 301 to github.com for `/`, a
repository URL with a query and encoded query value, and an encoded-space
path. HTTP first returned 301 to the corresponding HTTPS git URL.

The issued certificate covers git.alwaldend.com, is issued by Let's Encrypt
YR1, and expires on 2027-01-24 at 14:57:53 UTC. Four Git checks and ten existing
download checks passed. Local/public rendering and routing also passed twelve
isolated checks using the pinned Traefik binary.

Private plan digest, sanitized change summary, deployment logs, and HTTP
receipts are under `out/gr/`; the isolated fixture is under
`out/git-redirect/`. The download owner records proxy deployment and the
repository owner records the shared Traefik version upgrade. This verifies
issuance and current service, not future certificate renewal. External resolver
caches can retain the former record for its TTL.

## Go discovery correction (2026-10-30)

The original browser probes did not establish Go module discovery. The
public redirect also matched go-get=1 and GitHub returned a different import
prefix. The download owner's [discovery evidence](../../../../download/openspec/changes/add-static-hosting/evidence.md#go-discovery-before-browser-redirection-2026-10-30)
records the reproduced 301, 42 real Traefik/Nginx fixture checks, and a real
nested-module lookup with the pinned Go client after the source fix. Public
DNS and local Forgejo remain unchanged. This correction has not been deployed.
