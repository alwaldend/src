## Context

RouterOS 1.99.1 supports MIKROTIK_HOST, MIKROTIK_USER and MIKROTIK_PASSWORD.
The shared infra.mikrotik helper already injects those variables from Vault.
The new account and printer firewall roots each use a single RouterOS identity;
the existing ingress root demonstrates an empty provider block.

## Decisions

Place the general preference in repo-infra's Terraform section. Preserve
explicit provider arguments for different aliased identities, which cannot all
use the same process-wide credential variables. Verify support at the pinned
provider version and preserve the existing Vault and backend boundaries.

The offline Promptfoo target validates case loading and packaging only; it does
not prove agent behavior or authenticate to a provider. No live evaluation or
infrastructure operation is required.

## Source verification

Seven focused tests passed, covering both RouterOS roots, skill discovery,
Promptfoo configuration loading and the tools/agents OpenSpec workspace.
The skill validation aspect and affected semantic lint passed, as did all
27 repository quality tests. Final candidate checks and publication are
recorded under ignored out/printer1-migration. No live evaluation or
infrastructure operation was performed.
