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
[landing repositories](../../../tf/alwaldend_pages_repos.tf),
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

Landing resources SHALL retain their existing addresses, using the catalog's
`landing_project` field. Retired projects SHALL retain their published
repositories while their catalog entries remain. Explicitly retired landing
repositories SHALL be removed from the catalog after the authorized retirement
workflow has established their exact deletion scope.

#### Scenario: Preserve a published identity after project retirement

- **WHEN** a source project leaves the build registry but remains in the
  repository catalog
- **THEN** its landing repository and environment remain configured at their
  existing resource addresses

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
non-landing Dependabot alert settings SHALL be preserved.

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

#### Scenario: Publish to a landing repository's Pages branch

- **WHEN** a catalog developer publishes site content to the landing
  repository's `pages` branch
- **THEN** repository write access permits publication
- **AND** its separate `master` default branch remains protected from developer
  pushes and merges

### Requirement: Separate landing defaults from Pages publication

Landing repositories SHALL use the catalog default branch independently of
their Pages source branch. A missing default branch SHALL be created from the
existing Pages branch before Terraform selects it as the default. Existing
default branches SHALL be imported. The migration SHALL NOT rename or delete
the Pages branch or change its content.

#### Scenario: A landing repository only has its Pages branch

- **WHEN** Terraform changes its default branch from `pages` to `master`
- **THEN** it creates `master` from the current Pages branch tip
- **AND** it selects `master` as the default after creation
- **AND** Pages continues to publish from the original `pages` branch

### Requirement: Enable Pages only after its source branch exists

Terraform SHALL configure landing Pages from the catalog, except for entries
explicitly in `bootstrap_projects`. That set SHALL default to empty. New
bootstrapping repositories SHALL omit default-branch creation and selection
until their Pages branch has been published.

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

### Requirement: Provide an explicit certificate reprovision procedure

GitHub issues the custom-domain certificate as part of its Pages build, so the
package SHALL provide a supported, separately named way to omit a live site's
Pages block and then restore it, forcing a new certificate. This intent SHALL
NOT be expressed through `bootstrap_projects`, whose documented meaning remains
projects whose `pages` branch does not yet exist. When a project is named in
either set, its Pages block SHALL be omitted; an empty set for both SHALL enable
Pages for every project. The documented procedure SHALL include the two filtered
applies and an HTTPS verification step.

#### Scenario: Reprovision a live site's certificate

- **WHEN** an authorized apply names one live site in the certificate
  reprovision set
- **THEN** the plan removes only that project's Pages block
- **AND** the documented procedure restores it with a later apply that leaves
  both sets empty
- **AND** the operator verifies the certificate served for the custom domain

#### Scenario: Bootstrap a new site without affecting served sites

- **WHEN** a new project is named in `bootstrap_projects`
- **THEN** no project that is absent from both sets loses its Pages block
- **AND** an already-served site is never named in `bootstrap_projects` to force
  recovery
