# Infrastructure GitHub

## Purpose

Describe Terraform ownership of GitHub repositories and Pages settings for
project landing sites. This baseline records desired configuration; it does not
assert that repositories or published sites were inspected live.

Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`.
Observation date: 2026-09-08. Sources are linked in full; no excerpts are used.

Sources: [component documentation](../../../README.md),
[Terraform workflow](../../../tf/README.md),
[registry-derived inputs](../../../tf/BUILD.bazel),
[repository resources](../../../tf/alwaldend_pages_repos.tf), and
[provider configuration](../../../tf/provider.tf).

## Requirements

### Requirement: Derive landing repositories from the project registry

The package SHALL generate its Terraform project-to-repository mapping from
`PROJECTS`. Each entry SHALL name a public repository by replacing underscores
in the project name with hyphens and appending `-landing`, under the configured
`alwaldend` GitHub owner.

#### Scenario: Generate the mapping for a project with underscores

- **WHEN** a registered project is named `example_project`
- **THEN** its generated repository name is `example-project-landing`
- **AND** its configured homepage is `https://example-project.alwaldend.com/`

### Requirement: Enable Pages only after its source branch exists

Terraform SHALL configure Pages to use the `pages` branch at `/` with the
project's hyphenated `alwaldend.com` subdomain, except for entries explicitly in
`bootstrap_projects`. That set SHALL default to empty.

#### Scenario: Bootstrap a new landing repository

- **WHEN** an authorized apply includes a new project in `bootstrap_projects`
- **THEN** the repository resource omits its Pages block
- **AND** the documented workflow publishes the site before a later apply removes
  that project from the bootstrap set and enables Pages

### Requirement: Keep repository provisioning separate from site and DNS publication

This package SHALL own repository resources, Pages settings, and the
`github-pages` environment with custom branch policies. The documented workflow
SHALL use `//projects:deploy_landings` for site content and `//infra/dns` for DNS
records.

#### Scenario: Review a landing-site rollout

- **WHEN** an operator prepares an authorized project landing rollout
- **THEN** repository and Pages changes are reviewed in the Terraform workflow
- **AND** built content publication and DNS changes use their separately owned
  workflows
