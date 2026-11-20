## Context

VictoriaMetrics and ingress share a private routed path through WireGuard.

## Goals / Non-Goals

Collect ingress metrics without adding a separate exporter or a public route.

## Decisions

Use Traefik's built-in Prometheus handler on a separate entry point. UFW
permits the collector address from the VictoriaMetrics DNS declaration.
VictoriaMetrics derives targets from the ingress inventory's WireGuard
addresses. The collector uses the port defined in the shared Traefik
role defaults.

## Risks / Trade-offs

The offline checks can verify rendered restrictions and scrape behavior.
Live reachability through WireGuard needs verification after deployment.
