# Development DNS

This root manages every record in the canonical [declarations](dnsconfig.json):
secure1, dev1 and dev1's services. It uses the shared `src_infra_dev` AppRole,
DNS-only provider credentials and one Vault HTTP backend:
`alwaldend.com/vault1/approles/src_infra_dev/tf_backend/dns`.

Validate offline with `//infra/dev/dns:tf_tests.fmt_test`. Operational entry
points are `:tf.plan` and interactive `:tf.apply`.

The authorized deployment adopted seven RouterOS and five Cloudflare records
from the previous owners, preserving their provider identities, then renamed
them in place. Temporary import inputs and backend access targets have been
removed. Both previous states are empty and their access identities have been retired.
Secure1 now uses 192.168.10.102; its VM is owned by the separate tf_setup root.

The [rollout record](../openspec/changes/consolidate-development-hosts/deployment.md)
contains the sanitized handoff evidence. Rollback requires freezing writers
and reconciling state ownership through the [shared procedure](../../dns/README.md#migration-and-recovery);
reverting source alone does not restore the previous ownership.
