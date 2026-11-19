## Context

The shared catalog at `infra/repos/config.json` owns GitHub repository
settings, including the three merge-method booleans. Its defaults enable merge
commits, squash merges, and rebase merges. `orgs/alwaldend/repos/src.json`
overrides `allow_merge_commit` and `allow_squash_merge` to false, leaving only
`allow_rebase_merge` enabled from the defaults.

That combination is the opposite of the intended policy. Squash and rebase
merges produce a single-parent commit and rewrite or reattribute the feature
commit, while a merge commit preserves the reviewed commit as a parent.

The deployed state is consistent with the checked-in source: `alwaldend/src`
currently reports `allow_merge_commit: true`, `allow_squash_merge: false`, and
`allow_rebase_merge: false`. Source and deployed state disagree with the
intended policy, and the remaining catalog repositories still permit squash and
rebase merges.

## Goals

- Express one merge-commit-only policy for GitHub repositories in the source
  of truth, rather than per-repository exceptions.
- Make the policy verifiable offline, without contacting GitHub.
- Keep the reviewed feature commit as a parent of the default-branch commit,
  so the delivered history identifies the exact reviewed candidate.

## Non-Goals

- Changing GitLab or Forgejo merge behavior. Their consumers do not manage
  merge methods, and this catalog does not own those settings.
- Changing branch protection, rulesets, or required signatures.
- Enforcing the merge method at the branch level. GitHub repository settings
  govern the methods offered by the merge button; this change does not add a
  ruleset rule that does not exist in the GitHub API.

## Decisions

### One shared default rather than per-repository overrides

The policy belongs in `config.json` defaults because it expresses an
organization-wide intent. Per-repository negation would recreate the drift
that produced this problem: 29 records each required to restate the same
policy, with any omission silently re-enabling a disallowed method.

### Remove the `src` overrides

The `src` record's `allow_merge_commit: false` and `allow_squash_merge: false`
become redundant once defaults are merge-only, and `allow_merge_commit: false`
directly contradicts the new policy. Removing both keys lets `src` inherit the
shared default. The record's remaining settings are unrelated and stay.

`delete_branch_on_merge: true` remains on `src` as an existing intentional
difference.

### Validate in the module, not only in tests

A catalog precondition and a plan-test assertion both check the policy. The
precondition fails plan generation for any consumer, so a contradictory
repository record cannot reach a GitHub apply. The plan test pins the projected
output so the offline catalog test catches a regression without credentials.

## Risks

- **Live drift beyond `src`.** Applying the reviewed plan updates merge
  settings on the 28 catalog GitHub repositories. Some are forks and landing
  repositories whose settings were adopted from their observed state; the plan
  will show those as in-place changes. The apply is a live operation requiring
  separate authorization, and its plan must be reviewed before execution.
- **Merge-commit requirement and protected branches.** Contributors without
  default-branch write access open pull requests; enabling merge commits and
  disabling the alternatives does not grant additional write access. Existing
  rulesets continue to restrict direct pushes and merges.
- **Reverting squash history.** Squash merges already accepted before this
  change remain in history unchanged. This change is not retroactive and does
  not rewrite existing commits.
