---
title: GitLab Terraform
description: GitLab groups, memberships, repository imports, and forks
---

This module manages the GitLab groups, memberships, repository imports, and
forks declared in the shared [repository catalog](../../repos/README.md).
Its [catalog module](../../repos/tf/README.md) supplies names, source URLs,
and repository settings. Provider versions and checksums are recorded in the shared
[provider declarations](../../../third_party/terraform/include.MODULE.bazel).

Catalog administrators receive the Owner role; catalog users receive the
Developer role. Group defaults permit Maintainers and Owners to push and
merge into default branches, with force pushes and developer initial pushes
disabled. Developers can read repositories, create feature branches, and
open merge requests. Existing higher direct or inherited grants must be
reviewed separately because GitLab retains a user's highest access level.

GitHub copies use one-time Git imports, preserving their source default
branches. They have no ongoing synchronization. GitLab forks retain their
upstream fork relationship and target upstream when opening merge requests.
The catalog owns destination names and upstream identities.

[AL configuration](../al.lua) authenticates with the
[`src_infra_gitlab` AppRole](../../vault/tf/approles/src_infra_gitlab/main.tf),
injects the `gitlab_token` field as `GITLAB_TOKEN`, and supplies the existing
[Vault HTTP state backend](../../../tools/agents/skills/repo-infra/references/flow.md).
The configuration owns the Vault paths and plugin labels.

The provider uses GitLab.com by default. For another instance, set
`GITLAB_BASE_URL` to its API endpoint with a trailing slash, such as
`https://gitlab.example.com/api/v4/`; see the
[provider configuration](https://registry.terraform.io/providers/gitlabhq/gitlab/19.2.1/docs).

The [BUILD file](BUILD.bazel) exposes the standard `tf.*` commands and
`tf_tests.fmt_test`; a plan uses live Vault and GitLab access.

## Import and apply

The import blocks adopt existing catalog groups, projects, and their direct memberships.
Membership discovery reads the catalog's existing group IDs independently of
the managed group resource, allowing a newly added member to be created.
GitLab.com top-level groups must exist before Terraform can manage them.
Group, membership, and project resources have `prevent_destroy` enabled.

For an existing destination project, record its observed numeric ID as
`gitlab.id` in its catalog file. The project import block adopts that identity,
and a postcondition rejects an unexpected ID. Do not create a replacement for
an existing project or fork. Review the full plan for deletions and replacements
before every apply; either requires explicit authorization for that scope.

Adopt default-branch protections in two phases:

1. A new project starts without `gitlab.id`. Apply the reviewed full plan for
   groups, memberships, and projects. Wait
   for each repository import or fork to finish, then verify its default
   branch and effective protection against the catalog and group defaults.
2. Record the observed project ID in the catalog after verifying its protection.
   The existing protection declarations select only projects with recorded IDs
   and import `<project-id>:<default-branch>`. This keeps new project creation
   independent of import IDs that are unknown until creation finishes.
   Review and apply the full plan that imports or updates those protections.
   Existing protections must never enter the provider's Create operation:
   it can unprotect and recreate an existing default-branch rule.

Keep the original protection in place if it cannot be imported or its access
levels need a replacement. The pinned provider supports in-place
`allowed_to_push` and `allowed_to_merge` updates on GitLab.com; its scalar
`push_access_level` and `merge_access_level` fields can require replacement.
See the [provider's protection behavior](https://registry.terraform.io/providers/gitlabhq/gitlab/19.2.1/docs/resources/branch_protection).

Confirm that no overlapping project or group rule grants developers access
to a default branch. A restrictive rule does not override a more permissive
matching rule; see [GitLab protected branches](https://docs.gitlab.com/user/project/repository/branches/protected/).
