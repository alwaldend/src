# Public apex restoration

Observed on 2026-10-25 at 21:21 UTC (2026-10-26 locally). The user explicitly
requested restoring the public root site to GitHub Pages and deferring public
download deployment. The rollback used committed candidate
`dbc0756a6ea4bb6d9950081faa970499ba677cdb` through the owning
`//infra/dns/tf:tf.plan`, `tf.show`, and `tf.apply` targets.

## Applied changes

1. A targeted saved plan deleted only the public apex CNAME. The complete
   Cloudflare inventory confirmed its removal before the next plan.
2. A fresh full-root plan added exactly four A and four AAAA records, matching
   the pre-cutover GitHub Pages values with TTL 300 and proxying disabled.
3. A subsequent full-root plan returned exit code zero with no changes.

The final Cloudflare inventory contains 59 records. Comparing it with the
pre-cutover inventory, accounting for the separately authorized download aliases
and new IDs for the restored apex records, confirms every unrelated record's
identity and attributes are preserved. The 84-record RouterOS inventory is
unchanged, including the local apex and download aliases.

HTTPS to `alwaldend.com` returned HTTP 200 with certificate verification enabled
and an address override to `185.199.108.153`, confirming GitHub Pages serves
the hostname independently of local split DNS or resolver caches. Cached former
CNAME answers may remain until their 300-second TTL expires.

## Evidence and limits

Sanitized results and command logs are under ignored `out/dltf/`, notably
`rollback-restore-summary.json`, `rollback-verification.json`,
`rollback-public-https.status`, and `rollback-final-plan.log`. Raw plans and
provider state projections are private temporary artifacts and are removed
after summarization. The delivery receipt binds the final published source.

The local download VM remains running. The Yandex prerequisite folder and
identity were created, and six download resources were provisioned, but its
content disk and VM were not created because the cloud HDD quota was insufficient.
The three download DNS additions remain. No Ansible deployment or quota change
was performed; this change restores only public apex routing.
