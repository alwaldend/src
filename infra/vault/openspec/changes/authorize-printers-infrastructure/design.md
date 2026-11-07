## Context

See [proposal.md](proposal.md) for the requested ownership change.

## Goals / Non-Goals

Implement and verify checked-in definitions, then execute the separately
authorized provisioning and state adoption while preserving live printer
objects. Printer device configuration and firmware remain outside this change.

## Decisions

**Verdict: proceed.** Reuse the existing vault_approle and dns_access modules and operator membership. Register printers with the existing scoped AppRole identity group. MikroTik uses the operator-supplied routeros_admin credential within its own prefix; remove the former ingress-secret policy and its AppRole attachment. Printer state stays within its own AppRole prefix. Related source movement is maintained in infra/printers/openspec/changes/move-printer1-infrastructure.

## Risks / Trade-offs

Moving files does not move state. Follow the printer project's migration guide
before running operational roots. Offline checks cannot establish provider
permissions, deployed adoption or Wi-Fi connectivity.
