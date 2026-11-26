## Context

The standalone DNS package already uses a dedicated backend path and the
existing `src_infra_flux` AppRole in the shared Vault root.

## Decisions

- Keep the DNS package and authentication target in place. Use the MLflow
  pattern for a local AppRole root, administrative authentication, and a
  dedicated state backend. Preserve the role, policies, and operator
  membership so existing DNS and shared Vault callers can keep using them.
- Resolve the transferred roles' entities and groups by name in shared Vault.
  Keep shared grants there, including the Forgejo and Harbor grants.
  Preserve the existing addresses for the cluster authentication resources
  so the Transit key and OIDC clients can be adopted without replacement.
  Guard state relinquishment with non-destructive removal blocks. Move the
  indexed DNS policy to a temporary address in the old root so the removal
  selects only that policy.
- Remove the Flux host from the shared security inventory with its source
  package. Keep Harbor's own operator commands and use them in its README.

## Risks / Trade-offs

- A controller still reconciling from the deleted paths would lose its
  configuration source. Coordinate any live decommissioning separately.
- If DNS state still needs to be transferred, the shared DNS procedure uses
  the prior revision for the historical source root.

## State transfer

On 2026-11-26, the Flux root adopted 53 bindings from shared Vault state.
Terraform imported 47 resources. Vault provider v5.8.0 lacked import support
for the two OIDC assignments, clients, and providers, so the reviewed apply
adopted those six objects through their same-name update endpoints. Readback
confirmed that the settings and both clients' credentials were unchanged.

The scoped shared Vault apply relinquished all 53 bindings with zero resource
additions, changes, or deletions. The destination retained all 53 bindings.
An unrelated update to the TrueNAS user group was excluded from the apply.
