## Context

See [proposal.md](proposal.md) for the requested ownership change.

## Goals / Non-Goals

Implement and verify checked-in definitions, then execute the separately
authorized provisioning and state adoption while preserving live printer
objects. Printer device configuration and firmware remain outside this change.

## Decisions

**Verdict: proceed.** Create an accounts Terraform root separate from DNS; use the MikroTik-owned routeros_admin secret for account provisioning. Use the user-approved Random provider at 3.9.1 with its publisher checksum. Its sensitive result feeds the RouterOS user and Vault KV resource. Grant printers read on this exact secret. RouterOS permissions are coarse: read/write are router-wide; include api and rest-api together for RouterOS REST compatibility, but omit policy, sensitive, shell and UI login privileges. A full-access group is unnecessary. The existing DNS root stays DNS-only. A generic Vault random endpoint was rejected because the pinned provider exposes response fields without sensitive markings.

Use the shared infra.mikrotik helper to inject MIKROTIK_HOST, MIKROTIK_USER
and MIKROTIK_PASSWORD directly. The default provider needs no Terraform
credential variables; DNS aliases keep their separate explicit inputs.

## Risks / Trade-offs

Moving files does not move state. Follow the printer project's migration guide
before running operational roots. Offline checks cannot establish provider
permissions, deployed adoption or Wi-Fi connectivity.
