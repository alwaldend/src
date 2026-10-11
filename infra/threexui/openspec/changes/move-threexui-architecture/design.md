## Context

The architecture was a page in infra/arch's multi-page Drawio source. The shared
site Mermaid renderer already provides hermetic SVG and WebP output.

## Goals / Non-Goals

Keep one architecture source, show datacenter placement, and explain node service
relationships. This change does not configure or deploy hosts.

## Decisions

Use infra/threexui/architecture.mmd as the canonical source. Show njalla1 and
work1 as simple cells within njal.la and yc1 within Yandex Cloud. A separate
node-detail subgraph applies to every data-plane node: Traefik routes to 3x-ui
and Xray, 3x-ui manages Xray, and Xray connects to Mullvad using WireGuard.
Keep control-plane ACME on Vault and node ACME on Let's Encrypt.

Publish light/dark site SVGs and a 2x WebP from maintained Bazel rules. Remove
the old Drawio page and regenerate the other pages through their existing owner.

## Risks / Trade-offs

The overview is checked-in desired architecture, not a live-health claim.
The work1 implementation is delivered independently in the node PR.
Removing one Drawio page changes the source hash and derived IDs on the retained
pages; inspect their generated diff and run the owner freshness checks.
