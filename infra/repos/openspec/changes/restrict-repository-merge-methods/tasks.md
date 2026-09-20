## 1. Catalog source

- [x] 1.1 Set the shared GitHub defaults in `config.json` to permit merge
      commits only; verify `allow_squash_merge` and `allow_rebase_merge` are
      false while `allow_merge_commit` is true.
- [x] 1.2 Remove the merge-method overrides from the `alwaldend/src` record so
      it inherits the shared policy; verify its unrelated settings, including
      `delete_branch_on_merge`, are unchanged.

## 2. Validation and documentation

- [x] 2.1 Add a catalog precondition rejecting any GitHub repository that
      permits squash or rebase merges; verify it names the repository and method.
- [x] 2.2 Extend the catalog plan test to assert the merge-commit-only policy
      for every projected GitHub repository; verify the test fails when a
      record enables a disallowed method.
- [x] 2.3 Document the merge-method policy in the catalog and GitHub consumer
      READMEs; verify the text names the owning source of truth.

## 3. Offline validation

- [x] 3.1 Run the configured formatter and the affected package tests; verify
      the formatted source and projected output against the tested candidate.

## 4. Delivery

- [ ] 4.1 Commit, push, and offer the pull request; verify the publication
      receipt identifies the tested candidate and resolve applicable review
      findings.

## 5. Authorized live apply

- [x] 5.1 After the user's explicit authorization, review the complete GitHub
      plan and apply the merge-settings changes; verify each repository reports
      merge commits permitted with squash and rebase disabled. The reviewed plan
      proposed 0 to add, 28 to change, and 0 to destroy; the apply reported the
      same totals. All 28 repositories were then read back from the GitHub API
      as `allow_merge_commit=true`, `allow_squash_merge=false`, and
      `allow_rebase_merge=false`, and a follow-up plan reported no changes.
      The plan also corrected pre-existing drift on `alwaldend/src`, whose live
      `allow_update_branch` was true while the checked-in catalog already
      declared false; that change is unrelated to merge methods and was
      included in the reviewed plan.
