## Deployment boundary

Use the owning `//projects/sri/tf` wrappers and the existing canonical DNS
declaration. Import the exact existing provider ID, then require a plan with
zero record additions, changes, replacements, or deletions before applying.
Compare the public zone before and after adoption and verify public resolution.

The shared [cutover procedure](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns migration ordering, writer coordination, credential handling, and recovery.
Keep record ownership enabled after import; disabling it would propose deletion.

## Evidence

At source revision `246c4e111b4c69f442a0530f1342fb886a7074c6`, the scoped
Vault plan added nine resources under `module.src_projects_sri` with no existing
resource changes or deletions. Applying that saved plan succeeded, and the
operator's SecretID capability changed from denied to read/update.

At `2026-09-13T01:17:49Z`, source revision
`9ec1bf2d11863518d7a02c57693ca9c016403f65` imported Cloudflare record
`b0617da7b67b2c43bcecdab35939eb5b` at
`module.dns.cloudflare_dns_record.records["sri_landing/CNAME/global"]`.
The reviewed adoption plan and apply changed no managed resources. All 80
public zone records remained byte-for-byte equivalent as JSON objects keyed
by provider ID; no record was added, removed, or changed. Public CNAME
resolution remained `alwaldend.github.io`, with TTL 600 and proxying disabled.

Raw plans and provider inventories stay in private ignored task scratch; they
are not publication artifacts. The root Vault and RouterOS recovery probes
also succeeded with direct IP connections and verified hostname TLS identities.
