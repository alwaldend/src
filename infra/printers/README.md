---
title: Printers
description: Printer DNS and network access in dc1
---

This project owns printer DNS declarations and network access. The
[DNS root](dns/README.md) publishes the canonical [declaration](dnsconfig.json).
The [service configuration](tf/README.md) permits Wi-Fi access to printer1.
The [MikroTik accounts root](../mikrotik/accounts/README.md) owns its RouterOS
account and Vault credential. Printer firmware and print queues are unmanaged.

The existing printer records and firewall rule have been adopted into this
project's state. The [migration receipt](docs/state-migration.md) records their
preserved identities and the owning validation commands.

## Network

The [diagram source](docs/network.mmd) describes the checked-in topology.
The printer connects to the router over the wired LAN; its DNS name resolves
directly to the printer. Wi-Fi access uses the printer-only firewall exception.

[![Printer network](docs/assets/network.svg)](docs/assets/network.svg)
