## Context

See [proposal.md](proposal.md) for the requested outcome. This change follows
the [shared repository-catalog adoption](../../../../../repos/openspec/changes/archive/2026-09-13-adopt-shared-repository-catalog/design.md).
The user subsequently explicitly requested deletion of the rule-project
landing repositories and DNS records, removal of their landing configuration,
and relocation of every Bazel rule project from `projects/` to `tools/`.
That newer authorization applies to this cohort; unrelated resources and
AppRoles remain outside the retirement scope.

The source inventory on 2026-09-13 contains the following modules and live
landing identities. Each source module moves from `projects/<module>/` to
`tools/<module>/`. The existing `rules_skill` publication spelling is retained
in the retirement inventory for the `rules_skills` module.

| Module                    | GitHub landing repository         | DNS hostname                            |
| ------------------------- | --------------------------------- | --------------------------------------- |
| `rules_binary_toolchain`  | `rules-binary-toolchain-landing`  | `rules-binary-toolchain.alwaldend.com`  |
| `rules_dnscontrol`        | `rules-dnscontrol-landing`        | `rules-dnscontrol.alwaldend.com`        |
| `rules_docs`              | `rules-docs-landing`              | `rules-docs.alwaldend.com`              |
| `rules_docs_gazelle`      | `rules-docs-gazelle-landing`      | `rules-docs-gazelle.alwaldend.com`      |
| `rules_hugo`              | `rules-hugo-landing`              | `rules-hugo.alwaldend.com`              |
| `rules_iso`               | `rules-iso-landing`               | `rules-iso.alwaldend.com`               |
| `rules_openspec`          | None                              | None                                    |
| `rules_promptfoo`         | `rules-promptfoo-landing`         | `rules-promptfoo.alwaldend.com`         |
| `rules_promptfoo_gazelle` | `rules-promptfoo-gazelle-landing` | `rules-promptfoo-gazelle.alwaldend.com` |
| `rules_skill_gazelle`     | `rules-skill-gazelle-landing`     | `rules-skill-gazelle.alwaldend.com`     |
| `rules_skills`            | `rules-skill-landing`             | `rules-skill.alwaldend.com`             |
| `rules_template`          | `rules-template-landing`          | `rules-template.alwaldend.com`          |

GitHub adoption completed before this retirement request; the linked
[adoption evidence](../../../../../repos/openspec/changes/archive/2026-09-13-adopt-shared-repository-catalog/design.md#github-adoption-evidence)
owns that baseline. Retirement and surviving-repository verification are
complete. The reduced catalog contains 28 GitHub repositories and 29 GitLab
projects, including F-Droid. The linked
[GitLab evidence](../../../../../repos/openspec/changes/archive/2026-09-13-adopt-shared-repository-catalog/design.md#gitlab-adoption-evidence)
records both completed adoption phases.

## Goals / Non-Goals

**Goals:** Preserve the relocated modules' public build interfaces and
standalone resolution; remove only the selected landing infrastructure;
keep the shared catalog and downstream imports consistent with retirement.

**Non-Goals:** Change provider or credential ownership, remove service
identities or AppRoles, retire unrelated sites or repositories, or add
ongoing forge synchronization.

## Decisions

### Preserve module identity while changing source ownership

Move complete module owners, including their specifications, and update
root overrides, nested relative paths, ignore boundaries, generated catalogs,
documentation, and skill discovery through their owning mechanisms. Module
names and supported public `.bzl` interfaces remain stable. The `tools/`
boundary documentation must account for these reusable rule modules without
changing unrelated tool publication boundaries. Remove dedicated landing
configuration while retaining generated rule documentation in the main site.

### Use two GitHub phases to retire managed default branches

The coordinator's decision is **proceed** with two-phase retirement after
verification of each complete plan. Provider 6.10.2's
[branch deletion](https://github.com/integrations/terraform-provider-github/blob/v6.10.2/github/resource_github_branch.go#L159-L175)
calls GitHub's reference deletion endpoint directly. GitHub rejects deletion
of a default branch with HTTP 422, and the provider's
[default-branch deletion](https://github.com/integrations/terraform-provider-github/blob/v6.10.2/github/resource_github_branch_default.go#L119-L130)
does not choose another branch. Deleting the managed `master` branches while
they remain defaults would therefore fail before repository deletion.

The first phase changes only the eleven selected `github_branch_default`
instances to `pages`. Keep the catalog and managed branch resources on
`master`; changing their branch field would invoke branch renaming rather
than prepare deletion. The second phase removes the selected repository
entries and destroys their managed objects after the defaults have changed.

A direct `removed` block for each selected branch was rejected because pinned
Terraform 1.14.8
[disallows indexed removal addresses](https://github.com/hashicorp/terraform/blob/v1.14.8/internal/addrs/remove_target.go#L13-L17).
No imperative reference deletion or state editing is needed for this plan.

### Restrict destruction by the complete plan address set

The selected eleven repositories each have six managed object types:
repository, Pages environment, created `master` branch, default-branch
selection, default-branch ruleset, and developer collaborator grant.
The expected second GitHub phase is 66 selected object deletions with no
unrelated resource action. Temporarily relax literal `prevent_destroy`
settings only on these six generic resource blocks, then restore them after
retirement. Organization, membership, nonlanding repository, and original
`src` ruleset guards remain enabled throughout.

This temporarily relaxes guards on surviving instances sharing those blocks.
The coordinator therefore verifies the complete plan's address set against
the explicit cohort before applying and restores guards before delivery.
The new user request authorizes the listed retirement; it does not authorize
incidental deletions discovered in a plan.

## Risks / Trade-offs

- Default-branch reference deletion fails while `master` remains the default.
  Verify the first phase changed exactly the selected defaults before planning
  destruction.
- Moving rule owners can break relative module paths, generated references,
  or specification discovery. Validate root consumers and each moved module
  with the pinned tooling and inspect representative documentation output.
- A deleted repository could leave a dangling custom domain. Retire the
  selected DNS records before deleting their GitHub repositories.
- Generic deletion guards cannot distinguish `for_each` keys. Admit only the
  exact selected destruction set and restore guards after the apply.
- Repository deletion is destructive. Failure recovery retains surviving
  resources and retries only corrected declarative plans; it must not recreate
  repositories or claim that their original remote IDs can be restored.

## Migration Plan

1. Preserve the exact cohort and current baseline, prepare source relocation,
   and validate the intended owner and public-interface mapping.
2. Prepare and verify the first GitHub plan: eleven selected default switches
   to `pages`, no additions, deletions, or replacements. Apply and verify.
3. Remove dedicated landing configuration and the selected catalog records.
   Regenerate owning catalogs and references. Prepare the DNS plan for exactly
   eleven selected record deletions, apply it, and verify the DNS outcome.
4. Prepare the second GitHub plan for the selected 66 managed deletions.
   Verify the complete address set, apply, and verify all eleven repository
   IDs are absent while surviving repository identities and settings remain.
5. Restore deletion guards and remove temporary retirement controls. Verify
   a clean GitHub follow-up plan and remaining catalog projections. Retain
   the DNS empty-state and record-ID absence receipts before removing its
   owning roots. Discard the old GitLab plan and review the reduced cohort
   before resuming the linked adoption.
6. Complete root and moved-module checks, inspect representative documentation
   and discovery outputs, record exact candidate and deployment receipts, and
   deliver through the repository workflow.

## Current Evidence and Next Action

Planning began on feature branch `t3code/add-infra-gitlab-terraform-project`
at source revision `aee0863bccddc2a46479cef7f0495aeb4279cb47` in its linked
worktree. The already-built pinned OpenSpec CLI reports version 1.11.0.
The GitHub baseline was observed at 2026-09-13T19:44:42.867198Z through 263
successful read-only API requests. The coordinator reports a complete
GitHub post-adoption Terraform plan with no changes. These facts establish
the pre-retirement baseline, not completion of this change.

The pinned CLI scaffold and artifact instructions were used for this change.
Initial strict change validation passed for this retirement and the linked
catalog adoption, and the artifacts' relative links resolved. Updated
artifacts still require validation; structural validation does not establish
implementation or deployment acceptance.

Completed retirement observations on 2026-09-13:

- The first GitHub apply changed exactly eleven selected defaults, with no
  additions or deletions. Receipt:
  `out/repos/private/github_retire_prepare_apply.log`.
- The second GitHub apply deleted exactly the 66 selected managed objects,
  with no additions or updates. Deletion guards are restored and temporary
  retirement controls are removed. The complete follow-up plan reports no
  changes. Receipts: `out/repos/private/github_retire_delete_apply_v3.log`
  and `out/repos/private/github_retirement_post.log`.
- GitHub API verification observed at 2026-09-13T20:09:35.019543Z completed
  208 requests with the expected results: 186 HTTP 200 and 22 HTTP 404.
  All eleven retired repository names and immutable IDs returned HTTP 404.
  The 28 surviving repository IDs, commit SHAs, and access settings match
  the preserved baseline. All twenty surviving landing defaults remain
  `master` with unchanged `pages` commits; all 22 surviving Pages settings
  and both original `src` rulesets remain preserved. Receipt:
  `out/repos/private/github_post_apply_summary.json`.
- Eleven reviewed DNS plans each deleted exactly one selected CNAME and
  contained no other resource action. All eleven owning states are now empty,
  and GET requests for the original Cloudflare record IDs returned HTTP 404.
  Receipts: `out/repos/private/rules_dns/retirement_inventory.json` and
  `out/repos/private/rules_dns/retirement_postconditions.json`. Initial states
  were preserved alongside post-apply states during verification.
- All twelve module owners now reside under `tools/`; their old project
  roots, dedicated landing directories, DNS Terraform roots, and DNS
  configuration files are absent. Owning DNS roots were removed after their
  live postconditions were verified, so DNS completion uses those receipts
  rather than a follow-up plan against removed configuration.
- AppRoles and their credential ownership remain preserved. GitLab adoption
  used the reduced catalog; its outcomes belong to the linked adoption
  evidence above.
- All 24 build/test commands across the twelve moved modules passed, covering
  46 passing tests, 23 cached. The sixteen focused package tests and 54 root
  consumer tests also passed. Receipts:
  `out/repos/private/nested_checks/receipt.json`,
  `out/repos/relocation_package_tests.log`, and
  `out/repos/relocation_consumer_tests.log`.
- Rendered documentation was inspected for all twelve `tools/` owners.
  Standard `.bazeliskrc` links now select the shared configuration; the
  `rules_dnscontrol` and `rules_openspec` nested locks were regenerated with
  the repository's Bazel 8.7 pin before their successful checks.

The implementation candidate below passed mandatory quality and lint checks
and was published before this archive. Private rollout artifacts were used
only during execution; the sanitized outcomes above preserve their scope and
postconditions.

## Validated Implementation Candidate

Implementation candidate `69ae7d11b3ae43fbb5621fd0b521ea8199d2102a`
(tree `259877706d04dfca1a21c4efe86fe712cd0bc5c4`) passed all 27
selected delivery checks, including root repository quality, Buildifier,
root consumer semantic lint, representative builds, and explicit build/test
checks in all twelve standalone modules. The standalone wrapper has no lint
profile; its supported builds and tests supplement the root checks.

The candidate was pushed and its [pull request](https://github.com/alwaldend/src/pull/87)
was verified by the owning delivery workflow. Infrastructure work is complete.
The pinned OpenSpec CLI archived this change and merged its two requirements
into the repository specification.
Private rollout artifacts were used for the observations above and are
removed before handoff; these sanitized outcomes and catalog identities are
the durable evidence. Final publication and review identity belong to the
repository delivery receipt and Git history.
