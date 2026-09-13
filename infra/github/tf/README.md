---
title: GitHub Terraform
description: GitHub organization, repository, access, and Pages management
---

This package consumes the [shared repository catalog](../../repos/README.md)
through its provider-free Terraform module. The catalog supplies the GitHub
owner, repositories, default branches, repository settings, organization
administrators, and developers. This root supports one GitHub organization per
provider instance and rejects multiple configured owners.

Existing repositories are imported, including forks and published landing
repositories whose source project has been retired. The `landing_project`
catalog key preserves the existing Terraform addresses. Published repository names and Pages
domains come directly from the catalog; they are not regenerated from the
current build-project registry.

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
settings for newly adopted non-landing repositories are similarly preserved;
landing repositories retain their existing enabled-alert policy.

Catalog administrators are organization owners. Developers receive explicit
write access on every catalog repository, while organization base permissions
remain at the catalog setting. The grants are non-authoritative, preserving
other collaborators. Each repository has an active default-branch ruleset
restricting creation, updates, deletion, and force pushes, with an organization
administrator bypass. Developers can push other branches and open pull
requests, but cannot directly push or merge into the default branch. Forks
also disable the upstream fetch-and-merge exception. Both pre-existing `src`
rulesets retain their separate release-branch protections and bypass actors;
the default-branch ruleset also applies alongside them. Repository rulesets are
available for these public repositories on
[GitHub Free](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/about-rulesets).

Landing repositories use the catalog default branch independently of their
Pages publication branch. Terraform creates a missing default branch from the
existing Pages branch, then changes the default in place. It never renames or
deletes the Pages branch. Default-branch rules protect `master`, while the bot
retains write access to `pages` for the existing deployment workflow.

## Publish a landing site

Each landing repository receives its built Hugo site from
`//projects:deploy_landings` in the root workspace. Terraform manages the
catalog-defined custom domain and Pages source, plus the existing
`github-pages` environment with custom branch policies.

For a new site, first plan and apply with `bootstrap_projects` containing only
the new projects whose `pages` branches do not yet exist. For example, pass
`-var='bootstrap_projects=["new_project"]'` to `//infra/github/tf:tf.plan` and
the reviewed `tf.apply`. This creates the repository without enabling Pages or
creating its separate default branch before the Pages branch has been
published. Do not include an existing live site: doing so
would remove its Pages configuration.

Publish the built site with `//projects:deploy_landings`, then plan and apply
with the default empty `bootstrap_projects` to enable Pages. Keep all applies
limited to reviewed, authorized repository and Pages changes. DNS is managed
by each site's owning Terraform root; [DNS documentation](../../dns/README.md)
describes the workflow. Validate the public custom domain after rollout.

## Reprovision a landing-site certificate

GitHub issues the custom-domain certificate as part of its Pages build and
offers no direct reissue action. To force a new certificate for one live site,
plan and apply with only that project in `certificate_reprovision_projects`:

```sh
bazel_agent bazel run //infra/github/tf:tf.plan -- \
  -var='certificate_reprovision_projects=["affected_project"]'
```

Inspect that the plan removes only that project's Pages block, then apply the
reviewed changes. Restore the site with a second plan and apply that leaves both
`bootstrap_projects` and `certificate_reprovision_projects` empty. GitHub
schedules a new certificate once the custom domain returns.

Republish the site with `//projects:deploy_landings` only if the first apply
removed its `pages` branch content; normally the branch is unchanged. Verify the
pushed revision, Pages build state, DNS, and the certificate served for the
custom domain before closing the recovery.

Keep this variable empty outside an active recovery. Never name an existing
served site in `bootstrap_projects` to force a reissue.
