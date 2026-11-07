# Printer DNS

This DNS-only root manages the canonical [printer declaration](../dnsconfig.json)
through the shared DNS module. Printer1 retains `printer1.dc1.alwaldend.com`,
its private IPv4 address, and both dc1 and global DNS views. Public resolution
does not provide a route to the printer. Records remain unproxied.

The `tf=dns` wrapper selects `src_infra_printers`, the shared DNS provider
credentials, and the owner's `tf_backend/dns` Vault HTTP backend. It does not
use the printers firewall account. DNS management defaults to enabled to
preserve the existing records adopted through the completed
[state transfer](../docs/state-migration.md). Both views retain their original
provider identities.

Run offline formatting, pinned-provider schema validation and declaration checks:

```sh
bazel test //infra/printers/dns:tf_tests.fmt_test \
  //infra/printers/dns:tf_tests.validate_test //infra/dns:config_test
```

Operational targets use `//infra/printers/dns:tf.<operation>` and require explicit
live authority. Expected reconciliation after adoption is no DNS changes.
