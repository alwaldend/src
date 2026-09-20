## Purpose

Define the initial Terraform-managed caches for public PyPI and npm packages
and Docker Hub images, including access and client verification behavior.

## ADDED Requirements

### Requirement: Declare three format-specific proxy repositories

Service Terraform SHALL manage `pypi-proxy`, `npm-proxy`, and
`dockerhub-proxy` using the native PyPI, npm, and Docker proxy formats.
Their respective upstreams SHALL be `https://pypi.org/`,
`https://registry.npmjs.org/`, and `https://registry-1.docker.io/`, with Docker
Hub index behavior selected for the Docker proxy. This initial managed set
SHALL NOT include Git proxies, hosted repositories, or repository groups.

#### Scenario: Initial configuration is applied

- **WHEN** Nexus is ready, bootstrap is complete, and the service root is applied
- **THEN** all three named proxies exist with the intended format, upstream,
  and persistent filesystem blob-store association

### Requirement: Keep API configuration under one declarative owner

Service Terraform SHALL own proxy settings, blob-store definitions, client
access settings, and any declared cleanup policy. Ansible SHALL own service
installation and bootstrap without independently reconciling these resources.
Built-in or pre-existing resources SHALL have an explicit adoption or exclusion
decision before management; unrecognized data SHALL NOT be silently deleted.

#### Scenario: Existing configuration is reconciled

- **WHEN** a managed proxy setting drifts and service Terraform is applied
- **THEN** its declared value is restored and the next plan reports no drift
  for the managed configuration

#### Scenario: An unexpected repository already exists

- **WHEN** implementation or deployment encounters a pre-existing repository
  outside the declared managed set
- **THEN** it records or reports that resource for explicit handling instead
  of deleting it as a side effect of the initial configuration

### Requirement: Authenticate clients and protect upstream credentials

Anonymous repository access SHALL be disabled by default. The deployment
SHALL document authenticated read access for pip, npm, and Docker clients
through HTTPS. Docker clients SHALL use Nexus's Docker authentication flow.
Optional upstream credentials SHALL come from Vault and SHALL NOT appear in
tracked configuration, client examples, or ordinary logs.

#### Scenario: A client lacks repository permission

- **WHEN** an unauthenticated or unauthorized client requests a protected artifact
- **THEN** Nexus refuses access while an authorized read-only client can fetch
  the same artifact without gaining configuration or upload privileges

### Requirement: Verify actual artifact retrieval and cache reuse

Acceptance SHALL exercise a pinned Python package download, a pinned npm
package download, and a Docker Hub image pull by digest through the documented
Nexus endpoints. Cache reuse SHALL be verified with a new client cache and
upstream request evidence for immutable artifact bytes. Metadata revalidation
and registry token requests SHALL be distinguished from artifact downloads.

#### Scenario: Client downloads an uncached artifact

- **WHEN** an authorized client requests an available artifact absent from Nexus
- **THEN** Nexus retrieves it from the configured upstream, stores it in the
  persistent blob store, and returns the expected content or digest

#### Scenario: Another client requests the cached artifact

- **WHEN** a client with an empty local cache requests the same immutable
  artifact while the Nexus cache entry remains valid
- **THEN** Nexus serves its cached artifact bytes without downloading those
  bytes again from upstream

### Requirement: Make retention and storage behavior explicit

Proxy cache ages, storage quotas, and cleanup behavior SHALL be explicit
deployment settings. Destructive cleanup SHALL remain disabled until its
retention criteria are selected and documented. Service configuration SHALL
preserve cached data across ordinary Ansible reruns and service restarts.

#### Scenario: No retention policy has been selected

- **WHEN** the initial proxy configuration is applied
- **THEN** it does not schedule an implicit destructive cleanup policy and
  documents how operators observe and manage available storage
