## Why

GitHub, GitLab, and Forgejo need one catalog of organizations, repositories,
and named access roles. Existing repositories and state must survive adoption.

## What Changes

- Add the shared catalog under `infra/repos` with separate organization and
  repository JSON files and consume it in all three stages.
- Adopt the existing GitHub organization and every repository through imports.
- Adopt the GitLab group, import GitHub repositories once, and fork F-Droid's
  metadata repository into the group.
- Give the configured administrator administrative access and the configured
  developer repository/branch/PR access without default-branch write access.
- Use `master` as the landing repositories' protected default branch while
  preserving `pages` as their Pages source and deployment branch.
- Preserve existing resources, Terraform addresses, and service identities.
- Retain first-party names across forges; name organization-owned external
  forks and mirrors from their original upstream's reversed hostname and
  full owner/repository path.
- Record ongoing synchronization as deferred work; the user will configure it
  separately. Native GitLab pull mirrors are outside this implementation.
- Reconcile the catalog with the user's subsequent
  [rule-landing retirement](../../../../../src/openspec/changes/archive/2026-09-13-retire-bazel-rule-landings/proposal.md)
  before GitLab adoption. The original 39-repository GitHub adoption is
  verified; the remaining cohort is 28 GitHub repositories and those 28
  GitLab imports plus the F-Droid fork. [Design evidence](design.md) records
  completed adoption and links the verified retirement outcomes.

## Capabilities

### New Capabilities

- `repository-catalog`: Shared repository inventory, access roles, and consumer
  adoption boundaries.

### Modified Capabilities

None in this new workspace. Consumer documentation and existing Forgejo
requirements will be reconciled with the catalog's ownership of named roles.

## Impact

`infra/repos`, `infra/github/tf`, `infra/gitlab/tf`, `infra/forgejo/tf`, and their
documentation and OpenSpec integration. Live adoption is authorized for these
organizations and repositories. The subsequent user request explicitly
authorizes the selected eleven landing retirements, recorded in the linked
change; unrelated deletion or replacement remains outside that authorization.
Credentials remain in Vault.
