---
title: GitHub Terraform
description: GitHub organization, repository, access, and Pages management
---

This package consumes the [shared repository catalog](../../repos/README.md)
through its provider-free Terraform module. The catalog supplies the GitHub
owner, repositories, default branches, repository settings, organization
administrators, and developers. This root supports one GitHub organization per
provider instance and rejects multiple configured owners.

Existing repositories are imported, including forks. The apex site repository
keeps its Pages configuration and custom domain. Published repository names and
Pages domains come directly from the catalog; they are not regenerated from the
current build-project registry.

Dedicated per-project landing repositories were retired: the main site now
publishes every project landing at `/projects/<name>/`, so the catalog carries
no `landing_project` key and this root manages no per-project Pages repository,
custom domain, `github-pages` environment, or landing default branch.

## Adoption and access

The checked-in import blocks adopt the organization, existing repositories,
default branches, memberships, developer collaborators, and both existing
`src` rulesets. A repository with an observed catalog `github.id` must retain
that numeric identity. An import failure must be resolved before applying;
never remove an import merely to create an existing object again. Repositories,
memberships, access grants, environments, and rulesets have deletion guards.
Review the complete plan and obtain explicit authorization for any deletion
or replacement. Previously authorized retirement may proceed within its exact
reviewed scope.

Organization settings own only catalog-defined repository base permissions.
Billing, profile fields, and other organization defaults remain controlled by
GitHub and are explicitly ignored. The required empty `billing_email` argument
is an import-only placeholder: the existing billing address is retained from
state and must never be committed or exposed in plan output. Do not create
this organization-settings resource without its import block. Dependabot alert
settings for newly adopted repositories are similarly preserved.

Catalog administrators are organization owners. Developers receive explicit
write access on every catalog repository, while organization base permissions
remain at the catalog setting. The grants are non-authoritative, preserving
other collaborators. Each repository has an active default-branch ruleset
restricting creation, updates, deletion, and force pushes, with an organization
administrator bypass. Developers can push other branches and open pull
requests, but cannot directly push or merge into the default branch. Forks
also disable the upstream fetch-and-merge exception. Catalog repositories
permit merge commits only; squash and rebase merges are disabled, so an
accepted pull request records a merge commit whose parents include the
reviewed feature commit. Both pre-existing `src` rulesets retain their
separate release-branch protections and bypass actors; the default-branch
ruleset also applies alongside them. Repository rulesets are available for
these public repositories on
[GitHub Free](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/about-rulesets).

## Publish the main site

The apex site repository receives its built Hugo site from
`//projects/alwaldend.com:deploy` in the root workspace. Terraform manages the
catalog-defined custom domain and Pages source. Keep all applies limited to
reviewed, authorized repository and Pages changes. DNS is managed by the
owning Terraform root; [DNS documentation](../../dns/README.md) describes the
workflow. Validate the public custom domain after rollout.

GitHub issues the custom-domain certificate as part of its Pages build. If a
certificate must be reissued, remove the `pages` block from the apex catalog
record in a reviewed plan, apply it, then restore it in a second plan; GitHub
schedules a new certificate once the custom domain returns. Verify the pushed
revision, Pages build state, DNS, and the served certificate before closing the
recovery.
