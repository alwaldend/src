# Public website activation — 2026-10-26

The user explicitly authorized local publication, verification, authoritative
sync to Yandex, public verification, and root DNS cutover, including deletion
of the eight Pages apex records. The download owner's
[rollout record](../../../../download/openspec/changes/add-static-hosting/evidence.md#versioned-site-publication-and-mirror-2026-10-26)
owns the release, mirror, retained archives, and release-ref publication limit.

Before DNS changes, Nginx was enabled through `ansible.yandex --tags nginx`.
Five pages and four assets fetched over HTTP on Yandex matched the exact
2026.44.0 archive. Public archive downloads already passed trusted direct
HTTPS and byte-range validation. The Ansible mirror's repeat preview reported
no public projects/sites changes; its three changed controller tasks only
allocated, populated, and removed the private capture.

The owning `//infra/dns/tf:tf.plan`, `tf.show`, and `tf.apply` wrappers used
the existing Vault backend and provider identities. The first saved plan
deleted exactly four Pages A records and four AAAA records for alwaldend.com.
Apply reported 0 added, 0 changed, 8 destroyed. The second saved plan created
only the unproxied 300-second apex CNAME to download.alwaldend.com; apply
reported 1 added, 0 changed, 0 destroyed. Mail, local apex, www, staging, and
all unrelated records were unchanged. A final plan reported no changes.

Public DNS-over-HTTPS returned apex A 51.250.85.209 with TTL 300, no apex
AAAA records, both existing Proton Mail MX records, and www resolving to the
same public host. The controller's dc1 DNS still resolves the apex to
192.168.10.66.
An independent TCP query to terry.ns.cloudflare.com returned an authoritative
answer (`aa`) containing 51.250.85.209 with TTL 300.

After public DNS reached Yandex, `ansible.yandex --tags traefik` completed
with ok=18, changed=2, failed=0. Direct, proxy-free requests to 51.250.85.209
passed normal certificate validation for apex, www, download, and git:

- Five apex pages and four CSS/JavaScript assets matched the release archive
  byte-for-byte, including tags and project/tool/infrastructure documentation.
- www redirected permanently to the apex, preserving path and query.
- HTTP redirected permanently to HTTPS, preserving path and query.
- git redirected permanently to github.com, preserving path and query.
- The release archive returned HTTP 206 with the correct range, total length,
  and bytes. Both hosts retain the old and new archive checksums and only the
  selected extracted site.

Private saved plans, provider logs, source snapshots, and machine-readable
HTTP/DNS receipts remain under ignored `out/site-rollout`. Source declarations
and the public acceptance summary are the maintained records. Certificate
renewal and future OS update transactions remain outside this observation.
