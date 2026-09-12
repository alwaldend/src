# DNS Terraform

This root packages the owner's canonical [DNS declaration](../dnsconfig.json)
through the shared DNS module and keeps state in the owner AppRole's Vault HTTP
backend. The root contains DNS resources only. `dns_enabled` defaults to `true`
after verified adoption. Keep it enabled to retain existing records; disabling
it would propose deletion.

Use the normal `//infra/mikrotik/tf:tf.<operation>` wrappers for this DNS-only
root. The [DNS migration procedure](../../dns/README.md) owns credential
provisioning, reconciliation, and recovery. The
[adoption record](../openspec/changes/archive/2026-09-13-adopt-dns-records/design.md)
contains the import and inventory-preservation evidence. The package exposes
`//infra/mikrotik/tf:tf_tests.fmt_test` for offline formatting checks.

Operational Terraform calls require the DNS AppRole policy grants and credential
fields described by the migration procedure. Provider configuration requires
real credentials and live access independently of the activation setting.
