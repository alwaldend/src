# Kustomization Specification

## Purpose

Provide in-progress Kubernetes resource definitions for Flux, Traefik, and
cert-manager. This baseline records declared resources at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`, observed on 2026-09-08. It does not
assert successful rendering, reconciliation, or a deployed cluster state.

Sources: [project README](../../../README.md),
[Flux source](../../../flux-repo/flux-git-src.yaml),
[Traefik release](../../../traefik/flux/release.yaml),
[Traefik remote resources](../../../traefik/remote/kustomization.yaml),
and [cert-manager release](../../../cert-manager/flux/release.yaml).

## Requirements

### Requirement: Scoped Flux repository source

The Flux GitRepository declaration SHALL select the repository's `master`
branch, refer to the `flux-git-src` credential secret, and include only `infra`
and `projects/kustomization` from the repository through its ignore rules.

#### Scenario: Inspect repository source scope

- **WHEN** the checked-in Flux GitRepository is inspected
- **THEN** its branch, secret reference, and inclusion paths SHALL match the
  declared source scope.

### Requirement: Traefik Gateway API configuration

The Traefik definitions SHALL declare an OCI chart source and HelmRelease that
enable the Kubernetes Gateway provider, disable the Kubernetes Ingress provider,
and configure HTTP redirection to HTTPS with a named TLS certificate secret.

#### Scenario: Inspect gateway configuration

- **WHEN** the Traefik HelmRelease values are inspected
- **THEN** they SHALL enable HTTP and HTTPS gateway listeners and reference
  `traefik-gateway-websecure-tls` for HTTPS termination.

### Requirement: cert-manager reconciliation configuration

The cert-manager definitions SHALL declare an OCI chart source and HelmRelease
with Gateway API support, CRD installation, Helm tests, and drift detection
enabled.

#### Scenario: Inspect certificate controller configuration

- **WHEN** the cert-manager HelmRelease is inspected
- **THEN** its values SHALL enable Gateway API support and CRDs, and its release
  settings SHALL enable tests and drift detection.
