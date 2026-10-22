## Why

GitHub currently permits squash and rebase merges for nearly every catalog
repository. Those methods do not preserve the reviewed feature commit as a
parent of the default-branch commit, so the default-branch history does not
record the exact candidate that was reviewed and delivered. The catalog
already declares merge-method settings, but its shared defaults enable all
three methods, and the `alwaldend/src` overrides disable merge commits while
inheriting rebase merges, so the intended policy is neither expressed in the
source of truth nor consistent with the deployed repository.

## What Changes

- Set the shared GitHub catalog defaults to permit merge commits only:
  `allow_merge_commit` remains true while `allow_squash_merge` and
  `allow_rebase_merge` become false.
- Remove the `alwaldend/src` merge-method overrides, which currently disable
  merge commits and inherit rebase merges, so every repository inherits one
  shared policy instead of contradicting it.
- Add a catalog precondition and plan-test coverage that reject any GitHub
  repository permitting squash or rebase merges.
- Document the merge-method policy where the catalog and its GitHub consumer
  describe repository settings.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `repository-catalog`: The catalog gains a merge-commit-only repository
  policy for GitHub repositories, enforced across shared defaults and
  per-repository records, and the existing adoption requirement is modified
  to exempt these catalog-owned merge settings from the settings-preservation
  guarantee that otherwise applies to adopted repositories.

## Impact

`infra/repos` (catalog defaults, the `src` repository record, the Terraform
module and its plan test, and its README) and `infra/github/tf` documentation.
The GitHub, GitLab, and Forgejo consumers keep their existing structure and
state addresses. GitLab and Forgejo merge methods are not managed by this
catalog and remain outside this change. Applying the reviewed plan updates
repository settings on the 28 catalog GitHub repositories; that live apply
requires explicit separate authorization. Credentials remain in Vault.
