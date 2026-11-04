## Context

The existing GitHub AppRole declaration already provides the needed pattern.
See [proposal](proposal.md) for the requested identity.

## Goals / Non-Goals

Reuse the shared module's identity, namespace, and bootstrap policy. No new
policy implementation or live Vault operation is required.

## Decisions

Copy the GitHub declaration shape with `src_infra_vastai` as its name and
`simeonwarren` as operator. Disable Yandex policy attachment and omit the general
infrastructure group because neither capability is needed for this foundation.

## Risks / Trade-offs

Checked-in declarations do not establish deployed access. Deployment remains
a separately authorized apply through the existing Vault Terraform root.
