# DNS Terraform

This root packages the owner's canonical [DNS declarations](../dnsconfig.json)
through the shared DNS module and keeps state in the owner AppRole's Vault HTTP
backend. The root contains DNS resources only. `dns_enabled` defaults to `true`
after verified adoption. Keep it enabled to retain existing records; disabling
it would propose deletion.

The [DNS migration procedure](../../dns/README.md) owns credential provisioning,
reconciliation, and recovery. The [adoption record](../openspec/changes/archive/2026-09-13-adopt-dns-records/design.md)
contains import and verification evidence. Use the `//infra/nas/tf:tf` wrappers
for operational commands and `//infra/nas/tf:tf_tests.fmt_test` for offline
source validation.

Operational Terraform calls require the DNS AppRole policy grants and credential
fields described by the migration procedure. Provider configuration and
authentication require real credentials and live access.
