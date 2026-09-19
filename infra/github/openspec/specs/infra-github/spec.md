# Infrastructure GitHub

## Purpose

Describe Terraform ownership of the catalog-defined GitHub organization,
repositories, member access, default-branch restrictions, and Pages settings.
These requirements describe source guarantees; they do not assert that a plan
has been applied or that published sites are healthy.

Sources: [component documentation](../../../README.md),
[Terraform workflow](../../../tf/README.md),
[shared catalog](../../../../repos/README.md),
[packaged inputs](../../../tf/BUILD.bazel),
[repository resources](../../../tf/repositories.tf),
[organization access](../../../tf/organization.tf),
[branch rules](../../../tf/branch_rules.tf), and
[provider configuration](../../../tf/provider.tf).

## Requirements

### Requirement: Consume the shared organization and repository catalog

The package SHALL consume organization identity, administrator and developer
lists, repository identity, default branches, repository settings, and Pages
configuration through `infra/repos/tf`. It SHALL NOT generate a second inventory
from the build-project registry. Its provider instance SHALL require exactly
one configured GitHub organization.

#### Scenario: Consume a catalog repository

- **WHEN** the package renders its repository resources
- **THEN** every catalog GitHub repository is managed from its catalog record

### Requirement: Publish project landings from the apex site

Dedicated per-project landing repositories SHALL NOT be managed. The main site
publishes every project landing at `/projects/<name>/`, so the catalog carries
no `landing_project` key and no per-project Pages repository, custom domain,
`github-pages` environment, or landing default branch is declared.

#### Scenario: Publish a project landing

- **WHEN** a project's landing page is published
- **THEN** it is served by the apex site at `/projects/<name>/` from the apex
  Pages repository
- **AND** no per-project repository, Pages configuration, or custom domain is
  managed

### Requirement: Adopt existing resources without recreating them

The package SHALL provide import blocks for existing organization settings,
repositories, default branches, memberships, developer collaborators, and
the two existing `src` rulesets. Repositories with observed catalog IDs SHALL
retain those IDs. Managed resources SHALL reject Terraform destruction, and
any deletion or replacement SHALL require explicit user approval before apply.

The organization settings resource SHALL own only catalog-defined base
repository access. Its mandatory billing argument SHALL be an import-only
placeholder whose imported value is ignored; billing, profile, and other
organization defaults SHALL remain outside this module's ownership. Existing
Dependabot alert settings SHALL be preserved.

#### Scenario: An existing repository cannot be imported

- **WHEN** import cannot resolve an observed catalog repository
- **THEN** the operation fails rather than creating a replacement repository

### Requirement: Restrict developers to contribution branches

Catalog administrators SHALL be organization owners. Catalog developers SHALL
receive explicit write access on every GitHub catalog repository through
non-authoritative collaborator resources. Repository default branches SHALL
have active rulesets restricting creation, updates, deletion, and force
pushes, with bypass restricted to organization administrators. The update
restriction SHALL NOT permit fork syncing as a developer bypass. Ruleset
application SHALL follow the managed default-branch selection. The two
pre-existing `src` rulesets SHALL be imported unchanged, retaining their
rules and bypass actors alongside the new default-branch ruleset.

#### Scenario: A developer contributes a change

- **WHEN** a catalog developer writes to a non-protected contribution branch
- **THEN** repository write access permits the contribution and opening a pull
  request
- **AND** the developer cannot push or merge into the default branch

#### Scenario: Publish the main site

- **WHEN** the site publisher writes content to the apex repository's `pages`
  branch
- **THEN** repository write access permits publication
- **AND** its `master` default branch remains protected from developer pushes
  and merges

### Requirement: Keep repository provisioning separate from site and DNS publication

This package SHALL own repository resources and Pages settings. The documented
workflow SHALL use `//projects/alwaldend.com:deploy` for site content and
`//infra/dns` for DNS records.

#### Scenario: Review a site rollout

- **WHEN** an operator prepares an authorized site rollout
- **THEN** repository and Pages changes are reviewed in the Terraform workflow
- **AND** built content publication and DNS changes use their separately owned
  workflows

### Requirement: Provide an explicit certificate reprovision procedure

GitHub issues the custom-domain certificate as part of its Pages build and
offers no direct reissue action, so the package SHALL document a supported,
separately reviewed way to omit the apex site's Pages block and then restore
it, forcing a new certificate. The procedure SHALL use two filtered applies
and SHALL include an HTTPS verification step.

#### Scenario: Reprovision the apex site certificate

- **WHEN** an authorized apply removes the apex site's Pages block
- **THEN** the reviewed plan removes only that configuration
- **AND** the documented procedure restores it with a later apply
- **AND** the operator verifies the certificate served for the custom domain
