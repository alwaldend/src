## Context

Terraform consumers need sensitive creation options and observed TCP port mappings.

## Goals / Non-Goals

Add replacement-only sensitive docker_options with an empty compatible default and computed mapped_ports. Preserve imported and legacy instances. Exercise actual Terraform lifecycle fixtures.

## Decisions

Reuse the existing provider, pinned Terraform providers, and AL/Vault configuration flow.
Secrets remain sensitive in Terraform and reside in the Vault-backed backend; they are never public outputs.
The caller owns template trust and deployment authority. The 32 GB requirement is interpreted as 32 GiB VRAM.

## Risks / Trade-offs

Offline fixtures verify request and lifecycle behavior, not actual image startup or marketplace availability.
HTTPS uses the Vast-provided instance certificate; clients must verify its trust before submitting credentials.

## Migration Plan

Default SSH behavior is retained; changed startup settings intentionally replace a rental.
Deploy identity prerequisites separately before using the new rental root.

## Open Questions

None for implementation. A live rental requires an operator-selected reviewed template.
